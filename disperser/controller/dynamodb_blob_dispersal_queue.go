package controller

import (
	"context"
	"fmt"
	"math"
	"time"

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
) (BlobDispersalQueue, error) {

	if dynamoClient == nil {
		return nil, fmt.Errorf("dynamoClient cannot be nil")
	}
	if requestBatchSize == 0 {
		return nil, fmt.Errorf("requestBatchSize must be greater than 0")
	}
	if requestBatchSize > math.MaxInt32 {
		// This is annoying, but I'd rather not mess with the types of pre-existing interfaces right now.
		return nil, fmt.Errorf("requestBatchSize cannot be greater than %d", math.MaxInt32)
	}

	bdq := &dynamodbBlobDispersalQueue{
		ctx:                  ctx,
		logger:               logger,
		dynamoClient:         dynamoClient,
		queue:                make(chan *v2.BlobMetadata, queueSize),
		requestBatchSize:     requestBatchSize,
		requestBackoffPeriod: requestBackoffPeriod,
	}

	go bdq.run()

	return bdq, nil
}

func (bdq *dynamodbBlobDispersalQueue) GetNextBlobForDispersal(ctx context.Context) <-chan *v2.BlobMetadata {
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
		bdq.queue <- blobMetadata
	}

	return len(blobMetadatas) > 0, nil
}
