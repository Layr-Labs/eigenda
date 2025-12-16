package controller

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/config"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/common/replay"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser"
	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/metadata"
)

var errNoBlobsToEncode = errors.New("no blobs to encode")

// EncodingManagerConfig contains configuration parameters for the EncodingManager.
// The EncodingManager is responsible for pulling queued blobs from the blob metadata store,
// sending them to the encoder service for encoding, and creating blob certificates.
type EncodingManagerConfig struct {
	// PullInterval is how frequently the EncodingManager polls for new blobs to encode.
	// Must be positive.
	PullInterval time.Duration

	// EncodingRequestTimeout is the maximum time to wait for a single encoding request to complete.
	// Must be positive.
	EncodingRequestTimeout time.Duration
	// StoreTimeout is the maximum time to wait for blob metadata store operations.
	// Must be positive.
	StoreTimeout time.Duration
	// NumEncodingRetries is the number of times to retry encoding a blob after the initial attempt fails.
	// A value of 0 means no retries (only the initial attempt).
	// Must be non-negative.
	NumEncodingRetries int
	// NumRelayAssignment is the number of relays to assign to each blob.
	// Must be at least 1 and cannot exceed the length of AvailableRelays.
	NumRelayAssignment uint16
	// AvailableRelays is the list of relay keys that can be assigned to blobs.
	// Must not be empty.
	AvailableRelays []corev2.RelayKey
	// EncoderAddress is the network address of the encoder service (e.g., "localhost:50051").
	// Must not be empty.
	EncoderAddress string
	// MaxNumBlobsPerIteration is the maximum number of blobs to pull and encode in each iteration.
	// Must be at least 1.
	MaxNumBlobsPerIteration int32
	// OnchainStateRefreshInterval is how frequently the manager refreshes blob version parameters from the chain.
	// Must be positive.
	OnchainStateRefreshInterval time.Duration
	// NumConcurrentRequests is the size of the worker pool for processing encoding requests concurrently.
	// Must be at least 1.
	NumConcurrentRequests int
	// If true, accounts that DON'T have a human-friendly name remapping will be reported as their full account ID
	// in metrics.
	//
	// If false, accounts that DON'T have a human-friendly name remapping will be reported as "0x0" in metrics.
	//
	// NOTE: No matter the value of this field, accounts that DO have a human-friendly name remapping will be reported
	// as their remapped name in metrics. If you must reduce metric cardinality by reporting ALL accounts as "0x0",
	// you shouldn't define any human-friendly name remappings.
	EnablePerAccountBlobStatusMetrics bool
}

var _ config.VerifiableConfig = &EncodingManagerConfig{}

func DefaultEncodingManagerConfig() *EncodingManagerConfig {
	return &EncodingManagerConfig{
		PullInterval:                2 * time.Second,
		EncodingRequestTimeout:      5 * time.Minute,
		StoreTimeout:                15 * time.Second,
		NumEncodingRetries:          3,
		MaxNumBlobsPerIteration:     128,
		OnchainStateRefreshInterval: 1 * time.Hour,
		NumConcurrentRequests:       250,
		NumRelayAssignment:          1,
	}
}

func (c *EncodingManagerConfig) Verify() error {
	if c.PullInterval <= 0 {
		return fmt.Errorf("PullInterval must be positive, got %v", c.PullInterval)
	}
	if c.EncodingRequestTimeout <= 0 {
		return fmt.Errorf("EncodingRequestTimeout must be positive, got %v", c.EncodingRequestTimeout)
	}
	if c.StoreTimeout <= 0 {
		return fmt.Errorf("StoreTimeout must be positive, got %v", c.StoreTimeout)
	}
	if c.NumEncodingRetries < 0 {
		return fmt.Errorf("NumEncodingRetries must be non-negative, got %d", c.NumEncodingRetries)
	}
	if c.NumRelayAssignment < 1 {
		return fmt.Errorf("NumRelayAssignment must be at least 1, got %d", c.NumRelayAssignment)
	}
	if len(c.AvailableRelays) == 0 {
		return fmt.Errorf("AvailableRelays cannot be empty")
	}
	if int(c.NumRelayAssignment) > len(c.AvailableRelays) {
		return fmt.Errorf(
			"NumRelayAssignment (%d) cannot be greater than the number of available relays (%d)",
			c.NumRelayAssignment, len(c.AvailableRelays))
	}
	if c.MaxNumBlobsPerIteration < 1 {
		return fmt.Errorf("MaxNumBlobsPerIteration must be at least 1, got %d", c.MaxNumBlobsPerIteration)
	}
	if c.OnchainStateRefreshInterval <= 0 {
		return fmt.Errorf("OnchainStateRefreshInterval must be positive, got %v", c.OnchainStateRefreshInterval)
	}
	if c.NumConcurrentRequests < 1 {
		return fmt.Errorf("NumConcurrentRequests must be at least 1, got %d", c.NumConcurrentRequests)
	}
	if c.EncoderAddress == "" {
		return fmt.Errorf("EncoderAddress cannot be empty")
	}
	return nil
}

// EncodingManager is responsible for pulling queued blobs from the blob
// metadata store periodically and encoding them. It receives the encoder responses
// and creates BlobCertificates.
type EncodingManager struct {
	*EncodingManagerConfig

	// components
	blobMetadataStore blobstore.MetadataStore
	pool              common.WorkerPool
	encodingClient    disperser.EncoderClientV2
	chainReader       core.Reader
	logger            logging.Logger
	getNow            func() time.Time

	// state
	cursor                *blobstore.StatusIndexCursor
	blobVersionParameters atomic.Pointer[corev2.BlobVersionParameterMap]

	metrics                *encodingManagerMetrics
	controllerMetrics      *ControllerMetrics
	controllerLivenessChan chan<- healthcheck.HeartbeatMessage

	// Prevents the same blob from being processed multiple times, regardless of dynamo shenanigans.
	replayGuardian replay.ReplayGuardian
}

func NewEncodingManager(
	config *EncodingManagerConfig,
	getNow func() time.Time,
	blobMetadataStore blobstore.MetadataStore,
	pool common.WorkerPool,
	encodingClient disperser.EncoderClientV2,
	chainReader core.Reader,
	logger logging.Logger,
	registry *prometheus.Registry,
	controllerLivenessChan chan<- healthcheck.HeartbeatMessage,
	userAccountRemapping map[string]string,
	// For each blob, compare the blob's timestamp to the current time. If it's this far in the future, ignore it.
	// This is used by a replay guardian to prevent double-processing of blobs.
	maxFutureAge time.Duration,
	// For each blob, compare the blob's timestamp to the current time. If it's older than this, ignore it.
	// This is used by a replay guardian to prevent double-processing of blobs.
	maxPastAge time.Duration,
	controllerMetrics *ControllerMetrics,
) (*EncodingManager, error) {

	if err := config.Verify(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	replayGuardian, err := replay.NewReplayGuardian(getNow, maxPastAge, maxFutureAge)
	if err != nil {
		return nil, fmt.Errorf("failed to create replay guardian: %w", err)
	}

	return &EncodingManager{
		EncodingManagerConfig: config,
		getNow:                getNow,
		blobMetadataStore:     blobMetadataStore,
		pool:                  pool,
		encodingClient:        encodingClient,
		chainReader:           chainReader,
		logger:                logger.With("component", "EncodingManager"),
		cursor:                nil,
		metrics: newEncodingManagerMetrics(
			registry, config.EnablePerAccountBlobStatusMetrics, userAccountRemapping),
		controllerLivenessChan: controllerLivenessChan,
		replayGuardian:         replayGuardian,
		controllerMetrics:      controllerMetrics,
	}, nil
}

func (e *EncodingManager) Start(ctx context.Context) error {
	// Refresh blob version parameters
	err := e.refreshBlobVersionParams(ctx)
	if err != nil {
		return fmt.Errorf("failed to refresh blob version parameters: %w", err)
	}

	go func() {
		ticker := time.NewTicker(e.OnchainStateRefreshInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				e.logger.Info("refreshing blob version params")
				if err := e.refreshBlobVersionParams(ctx); err != nil {
					e.logger.Error("failed to refresh blob version params", "err", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// Start the encoding loop
	go func() {
		ticker := time.NewTicker(e.PullInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := e.HandleBatch(ctx)
				if err != nil {
					if errors.Is(err, errNoBlobsToEncode) {
						e.logger.Debug("no blobs to encode")
					} else {
						e.logger.Error("failed to process a batch", "err", err)
					}
				}
			}
		}
	}()

	return nil
}

// Iterates over the input metadata slice, and returns a new slice with stale and duplicate metadatas filtered out
func (e *EncodingManager) filterStaleAndDedupBlobs(
	ctx context.Context,
	inputMetadatas []*v2.BlobMetadata,
) []*v2.BlobMetadata {
	outputMetadatas := make([]*v2.BlobMetadata, 0, len(inputMetadatas))

	for _, metadata := range inputMetadatas {
		blobKey, err := metadata.BlobHeader.BlobKey()
		if err != nil {
			e.logger.Errorf("compute blob key: %w", err)
			// we must discard if we cannot compute key, since it's used for deduplication
			continue
		}

		timestamp := time.Unix(0, metadata.BlobHeader.PaymentMetadata.Timestamp)

		status := e.replayGuardian.DetailedVerifyRequest(blobKey[:], timestamp)
		switch status {
		case replay.StatusValid:
			outputMetadatas = append(outputMetadatas, metadata)
		case replay.StatusTooOld:
			e.controllerMetrics.reportDiscardedBlob("encodingManager", true)
			e.markBlobAsFailed(ctx, blobKey)
		case replay.StatusTooFarInFuture:
			e.controllerMetrics.reportDiscardedBlob("encodingManager", false)
			e.markBlobAsFailed(ctx, blobKey)
		case replay.StatusDuplicate:
			// Ignore duplicates
		default:
			e.logger.Errorf("Unknown replay guardian status %d for blob %s, skipping.", status, blobKey.Hex())
		}
	}

	return outputMetadatas
}

func (e *EncodingManager) markBlobAsFailed(ctx context.Context, blobKey corev2.BlobKey) {
	err := e.blobMetadataStore.UpdateBlobStatus(
		ctx,
		blobKey,
		v2.Failed,
	)
	if err != nil {
		e.logger.Errorf("Failed to mark blob %s as failed: %v", blobKey.Hex(), err)
	}
}

// HandleBatch handles a batch of blobs to encode
// It retrieves a batch of blobs from the blob metadata store, encodes them, and updates their status
// It also creates BlobCertificates and stores them in the blob metadata store
//
// WARNING: This method is not thread-safe. It should only be called from a single goroutine.
func (e *EncodingManager) HandleBatch(ctx context.Context) error {
	// Signal Liveness to indicate no stall
	healthcheck.SignalHeartbeat(e.logger, "encodingManager", e.controllerLivenessChan)

	// Get a batch of blobs to encode
	blobMetadatas, cursor, err := e.blobMetadataStore.GetBlobMetadataByStatusPaginated(
		ctx, v2.Queued, e.cursor, e.MaxNumBlobsPerIteration)
	if err != nil {
		return err
	}

	blobMetadatas = e.filterStaleAndDedupBlobs(ctx, blobMetadatas)
	if len(blobMetadatas) == 0 {
		return errNoBlobsToEncode
	}

	blobVersionParams := e.blobVersionParameters.Load()
	if blobVersionParams == nil {
		return fmt.Errorf("blob version parameters is nil")
	}

	e.metrics.reportBatchSize(len(blobMetadatas))
	batchSizeBytes := uint64(0)
	for _, blob := range blobMetadatas {
		batchSizeBytes += blob.BlobSize
	}
	e.metrics.reportBatchDataSize(batchSizeBytes)

	submissionStart := time.Now()

	e.logger.Debug("request encoding", "numBlobs", len(blobMetadatas))
	for _, blob := range blobMetadatas {
		blob := blob
		blobKey, err := blob.BlobHeader.BlobKey()
		if err != nil {
			e.logger.Error("failed to get blob key",
				"err", err,
				"requestedAt", blob.RequestedAt,
				"paymentMetadata", blob.BlobHeader.PaymentMetadata)
			continue
		}

		blobParams, ok := blobVersionParams.Get(blob.BlobHeader.BlobVersion)
		if !ok {
			e.logger.Error("failed to get blob version parameters", "version", blob.BlobHeader.BlobVersion)
			continue
		}

		// Encode the blobs
		e.pool.Submit(func() {
			start := time.Now()

			var i int
			var finishedEncodingTime time.Time
			var finishedPutBlobCertificateTime time.Time
			var finishedUpdateBlobStatusTime time.Time
			var success bool

			for i = 0; i < e.NumEncodingRetries+1; i++ {
				encodingCtx, cancel := context.WithTimeout(ctx, e.EncodingRequestTimeout)
				fragmentInfo, err := e.encodeBlob(encodingCtx, blobKey, blob, blobParams)
				cancel()
				if err != nil {
					e.logger.Error("failed to encode blob", "blobKey", blobKey.Hex(), "err", err)
					continue
				}

				finishedEncodingTime = time.Now()

				relayKeys, err := GetRelayKeys(e.NumRelayAssignment, e.AvailableRelays)
				if err != nil {
					e.logger.Error("failed to get relay keys", "err", err)
					// Stop retrying
					break
				}
				cert := &corev2.BlobCertificate{
					BlobHeader: blob.BlobHeader,
					Signature:  blob.Signature,
					RelayKeys:  relayKeys,
				}

				storeCtx, cancel := context.WithTimeout(ctx, e.StoreTimeout)
				err = e.blobMetadataStore.PutBlobCertificate(storeCtx, cert, fragmentInfo)
				cancel()
				if err != nil && !errors.Is(err, blobstore.ErrAlreadyExists) {
					e.logger.Error("failed to put blob certificate", "err", err)
					continue
				}

				finishedPutBlobCertificateTime = time.Now()

				storeCtx, cancel = context.WithTimeout(ctx, e.StoreTimeout)
				err = e.blobMetadataStore.UpdateBlobStatus(storeCtx, blobKey, v2.Encoded)
				finishedUpdateBlobStatusTime = time.Now()
				cancel()
				if err == nil || errors.Is(err, blobstore.ErrAlreadyExists) {
					// Successfully updated the status to Encoded
					success = true
					break
				}

				e.logger.Error("failed to update blob status to Encoded", "blobKey", blobKey.Hex(), "err", err)
				sleepTime := time.Duration(math.Pow(2, float64(i))) * time.Second
				time.Sleep(sleepTime) // Wait before retrying
			}

			e.metrics.reportBatchRetryCount(i)

			if success {
				e.metrics.reportEncodingLatency(finishedEncodingTime.Sub(start))
				e.metrics.reportPutBlobCertLatency(finishedPutBlobCertificateTime.Sub(finishedEncodingTime))
				e.metrics.reportUpdateBlobStatusLatency(
					finishedUpdateBlobStatusTime.Sub(finishedPutBlobCertificateTime))
				e.metrics.reportBlobHandleLatency(time.Since(start))

				requestedAt := time.Unix(0, int64(blob.RequestedAt))
				e.metrics.reportE2EEncodingLatency(time.Since(requestedAt))
				e.metrics.reportCompletedBlob(
					int(blob.BlobSize), v2.Encoded, blob.BlobHeader.PaymentMetadata.AccountID.Hex())
			} else {
				e.metrics.reportFailedSubmission()
				storeCtx, cancel := context.WithTimeout(ctx, e.StoreTimeout)
				err = e.blobMetadataStore.UpdateBlobStatus(storeCtx, blobKey, v2.Failed)
				cancel()
				if err != nil {
					e.logger.Error("failed to update blob status to Failed", "blobKey", blobKey.Hex(), "err", err)
					return
				}
				e.metrics.reportCompletedBlob(
					int(blob.BlobSize), v2.Failed, blob.BlobHeader.PaymentMetadata.AccountID.Hex())
			}
		})
	}

	e.metrics.reportBatchSubmissionLatency(time.Since(submissionStart))

	e.cursor = cursor

	e.logger.Debug("successfully submitted encoding requests", "numBlobs", len(blobMetadatas))
	return nil
}

func (e *EncodingManager) encodeBlob(
	ctx context.Context,
	blobKey corev2.BlobKey,
	blob *v2.BlobMetadata,
	blobParams *core.BlobVersionParameters,
) (*encoding.FragmentInfo, error) {
	// Add headers for routing
	md := metadata.New(map[string]string{
		"content-type": "application/grpc",
		"x-blob-size":  fmt.Sprintf("%d", blob.BlobSize),
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	encodingParams, err := corev2.GetEncodingParams(blob.BlobHeader.BlobCommitments.Length, blobParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get encoding params: %w", err)
	}
	return e.encodingClient.EncodeBlob(ctx, blobKey, encodingParams, blob.BlobSize)
}

func (e *EncodingManager) refreshBlobVersionParams(ctx context.Context) error {
	e.logger.Debug("Refreshing blob version params")
	blobParams, err := e.chainReader.GetAllVersionedBlobParams(ctx)
	if err != nil {
		return fmt.Errorf("failed to get blob version parameters: %w", err)
	}

	e.blobVersionParameters.Store(corev2.NewBlobVersionParameterMap(blobParams))
	return nil
}

func GetRelayKeys(numAssignment uint16, availableRelays []corev2.RelayKey) ([]corev2.RelayKey, error) {
	if int(numAssignment) > len(availableRelays) {
		return nil, fmt.Errorf(
			"numAssignment (%d) cannot be greater than numRelays (%d)", numAssignment, len(availableRelays))
	}
	relayKeys := make([]corev2.RelayKey, len(availableRelays))
	copy(relayKeys, availableRelays)
	// shuffle relay keys
	for i := len(relayKeys) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		relayKeys[i], relayKeys[j] = relayKeys[j], relayKeys[i]
	}

	return relayKeys[:numAssignment], nil
}
