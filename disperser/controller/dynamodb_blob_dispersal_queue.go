package controller

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/Layr-Labs/eigenda/common/replay"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

var _ BlobDispersalQueue = (*dynamodbBlobDispersalQueue)(nil)

// An implementation of BlobDispersalQueue that uses DynamoDB as the backend communication mechanism between the
// encoder and the controller.
type dynamodbBlobDispersalQueue struct {
	ctx    context.Context
	logger logging.Logger

	// used to interact with the DynamoDB table storing blob metadata
	dynamoClient blobstore.MetadataStore

	// cursor for iterating through blobs ready for dispersal
	cursor *blobstore.StatusIndexCursor

	// channel for delivering blobs ready for dispersal
	queue chan *v2.BlobMetadata

	// When requesting blobs from DynamoDB, the number of blobs to request in each batch.
	requestBatchSize uint32

	// If we query dynamo and it has no blobs ready for dispersal, wait this long before trying again.
	requestBackoffPeriod time.Duration

	// Prevents the same blob from being returned multiple times, regardless of backend dynamo shenanigans.
	replayGuardian replay.ReplayGuardian

	// Encapsulated metrics for the controller.
	metrics *ControllerMetrics
}

// NewDynamodbBlobDispersalQueue creates a new instance of DynamodbBlobDispersalQueue.
func NewDynamodbBlobDispersalQueue(
	ctx context.Context,
	logger logging.Logger,
	dynamoClient blobstore.MetadataStore,
	// The maximum number of blobs to keep in the queue at any time.
	queueSize uint32,
	// When requesting blobs from DynamoDB, the number of blobs to request in each batch.
	requestBatchSize uint32,
	// How long to wait before re-querying DynamoDB if no blobs are found.
	requestBackoffPeriod time.Duration,
	// For each blob, compare the blob's timestamp to the current time. If it's this far in the future, ignore it.
	maxFutureAge time.Duration,
	// For each blob, compare the blob's timestamp to the current time. If it's older than this, ignore it.
	maxPastAge time.Duration,
	// Encapsulated metrics for the controller. No-op if nil.
	metrics *ControllerMetrics,
) (BlobDispersalQueue, error) {

	if dynamoClient == nil {
		return nil, fmt.Errorf("dynamoClient cannot be nil")
	}
	if requestBatchSize == 0 {
		return nil, fmt.Errorf("requestBatchSize must be greater than 0")
	}
	if requestBatchSize > math.MaxInt32 {
		// This is annoying, but I'd rather not mess with the types of pre-existing interfaces right now.
		return nil, fmt.Errorf("requestBatchSize cannot be greater than %d, got %d", math.MaxInt32, requestBatchSize)
	}
	if requestBackoffPeriod < 0 {
		return nil, fmt.Errorf("requestBackoffPeriod must not be negative, got %v", requestBackoffPeriod)
	}
	if maxFutureAge < 0 {
		return nil, fmt.Errorf("maxFutureAge must not be negative, got %v", maxFutureAge)
	}
	if maxPastAge < 0 {
		return nil, fmt.Errorf("maxPastAge must not be negative, got %v", maxPastAge)
	}

	replayGuardian := replay.NewReplayGuardian(time.Now, maxPastAge, maxFutureAge)

	bdq := &dynamodbBlobDispersalQueue{
		ctx:                  ctx,
		logger:               logger,
		dynamoClient:         dynamoClient,
		queue:                make(chan *v2.BlobMetadata, queueSize),
		requestBatchSize:     requestBatchSize,
		requestBackoffPeriod: requestBackoffPeriod,
		replayGuardian:       replayGuardian,
		metrics:              metrics,
	}

	go bdq.run()

	return bdq, nil
}

func (bdq *dynamodbBlobDispersalQueue) GetBlobChannel() <-chan *v2.BlobMetadata {
	return bdq.queue
}

// A function that runs in the background to fetch blobs ready for dispersal and push them onto the queue.
func (bdq *dynamodbBlobDispersalQueue) run() {
	for {
		select {
		case <-bdq.ctx.Done():
			close(bdq.queue)
			return
		default:
			foundData, err := bdq.fetchBlobs()
			if err != nil {
				bdq.logger.Errorf("Error fetching blobs for dispersal: %v", err)
			}

			if !foundData {
				// No data found, back off for a bit
				select {
				case <-time.After(bdq.requestBackoffPeriod):
				case <-bdq.ctx.Done():
					// cleanup will happen in the outer select
				}
			}
		}
	}
}

// Fetch a batch of blobs ready for dispersal from DynamoDB and push them onto the queue. Returns true
// if at least one blob was fetched, false otherwise.
func (bdq *dynamodbBlobDispersalQueue) fetchBlobs() (bool, error) {
	blobMetadatas, cursor, err := bdq.dynamoClient.GetBlobMetadataByStatusPaginated(
		bdq.ctx,
		v2.Encoded,
		bdq.cursor,
		int32(bdq.requestBatchSize),
	)

	if err != nil {
		return false, fmt.Errorf("failed to fetch blobs from DynamoDB: %w", err)
	}

	bdq.cursor = cursor

	for _, blobMetadata := range blobMetadatas {
		if blobMetadata == nil {
			bdq.logger.Errorf("Fetched nil blob metadata, skipping.")
			continue
		}
		if blobMetadata.BlobHeader == nil {
			bdq.logger.Errorf("Fetched blob metadata with nil BlobHeader, skipping.")
			continue
		}

		hash, err := blobMetadata.BlobHeader.BlobKey()
		if err != nil {
			bdq.logger.Errorf("Failed to compute blob header hash, skipping: %v", err)
			continue
		}
		timestamp := time.Unix(0, blobMetadata.BlobHeader.PaymentMetadata.Timestamp)

		status := bdq.replayGuardian.DetailedVerifyRequest(hash[:], timestamp)
		switch status {
		case replay.StatusValid:
			bdq.queue <- blobMetadata
		case replay.StatusTooOld:
			bdq.metrics.reportStaleDispersal()
			bdq.markBlobAsFailed(hash)
		case replay.StatusTooFarInFuture:
			bdq.metrics.reportTimeTravelerDispersal()
			bdq.markBlobAsFailed(hash)
		case replay.StatusDuplicate:
			// Ignore duplicates
		default:
			bdq.logger.Errorf("Unknown replay guardian status %d for blob %s, skipping.", status, hash.Hex())
		}
	}

	return len(blobMetadatas) > 0, nil
}

func (bdq *dynamodbBlobDispersalQueue) markBlobAsFailed(blobKey corev2.BlobKey) {
	err := bdq.dynamoClient.UpdateBlobStatus(
		bdq.ctx,
		blobKey,
		v2.Failed,
	)
	if err != nil {
		bdq.logger.Errorf("Failed to mark blob %s as failed: %v", blobKey.Hex(), err)
	}
}
