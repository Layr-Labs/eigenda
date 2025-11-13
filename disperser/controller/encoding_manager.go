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

	// MaxDispersalAge is the maximum age a dispersal request can be before it is discarded.
	// Dispersals older than this duration are marked as Failed and not processed.
	//
	// Age is determined by the BlobHeader.PaymentMetadata.Timestamp field, which is set by the
	// client at dispersal request creation time (in nanoseconds since Unix epoch).
	MaxDispersalAge time.Duration
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
		MaxDispersalAge:             45 * time.Second,
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
	if c.MaxDispersalAge <= 0 {
		return fmt.Errorf("MaxDispersalAge must be positive, got %v", c.MaxDispersalAge)
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
	// blobSet keeps track of blobs that are currently being encoded
	// This is used to deduplicate blobs to prevent the same blob from being encoded multiple times
	// blobSet is shared with Dispatcher which removes blobs from this queue as they are packaged for dispersal
	blobSet BlobSet

	metrics                *encodingManagerMetrics
	controllerLivenessChan chan<- healthcheck.HeartbeatMessage
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
	blobSet BlobSet,
	controllerLivenessChan chan<- healthcheck.HeartbeatMessage,
) (*EncodingManager, error) {
	if err := config.Verify(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &EncodingManager{
		EncodingManagerConfig:  config,
		getNow:                 getNow,
		blobMetadataStore:      blobMetadataStore,
		pool:                   pool,
		encodingClient:         encodingClient,
		chainReader:            chainReader,
		logger:                 logger.With("component", "EncodingManager"),
		cursor:                 nil,
		metrics:                newEncodingManagerMetrics(registry),
		blobSet:                blobSet,
		controllerLivenessChan: controllerLivenessChan,
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
	now := e.getNow()

	for _, metadata := range inputMetadatas {
		blobKey, err := metadata.BlobHeader.BlobKey()
		if err != nil {
			e.logger.Errorf("compute blob key: %w", err)
			// we must discard if we cannot compute key, since it's used for deduplication
			continue
		}

		if e.checkAndHandleStaleBlob(ctx, blobKey, now, metadata.BlobHeader.PaymentMetadata.Timestamp) {
			// discard stale blob
			continue
		}

		if e.blobSet.Contains(blobKey) {
			// discard duplicate blob
			continue
		}

		outputMetadatas = append(outputMetadatas, metadata)
	}

	return outputMetadatas
}

// Checks if a blob is older than MaxDispersalAge and handles it accordingly.
// If the blob is stale, it increments metrics, logs a warning, and updates the database status to Failed.
// Returns true if the blob is stale, otherwise false.
func (e *EncodingManager) checkAndHandleStaleBlob(
	ctx context.Context,
	blobKey corev2.BlobKey,
	now time.Time,
	dispersalTimestamp int64,
) bool {
	dispersalTime := time.Unix(0, dispersalTimestamp)
	dispersalAge := now.Sub(dispersalTime)

	if dispersalAge <= e.MaxDispersalAge {
		return false
	}

	e.metrics.reportStaleDispersal()

	e.logger.Warnf(
		"discarding stale dispersal: blobKey=%s dispersalAge=%s maxAge=%s dispersalTime=%s",
		blobKey.Hex(),
		dispersalAge.String(),
		e.MaxDispersalAge.String(),
		dispersalTime.Format(time.RFC3339),
	)

	storeCtx, cancel := context.WithTimeout(ctx, e.StoreTimeout)
	err := e.blobMetadataStore.UpdateBlobStatus(storeCtx, blobKey, v2.Failed)
	cancel()
	if err != nil {
		e.logger.Errorf("update stale blob status to Failed: blobKey=%s err=%w", blobKey.Hex(), err)
	} else {
		// we need to remove the blobKey from the blobSet once the BlobStatus is set to FAILED
		// the Dispatcher removes the blobKey from the blobSet when batching, but blobs that are set to FAILED
		// never are batched, and therefore must be removed manually
		e.blobSet.RemoveBlob(blobKey)
	}

	return true
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
	blobMetadatas, cursor, err := e.blobMetadataStore.GetBlobMetadataByStatusPaginated(ctx, v2.Queued, e.cursor, e.MaxNumBlobsPerIteration)
	if err != nil {
		return err
	}

	blobMetadatas = e.filterStaleAndDedupBlobs(ctx, blobMetadatas)
	e.metrics.reportBlobSetSize(e.blobSet.Size())
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
			e.logger.Error("failed to get blob key", "err", err, "requestedAt", blob.RequestedAt, "paymentMetadata", blob.BlobHeader.PaymentMetadata)
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
				e.metrics.reportCompletedBlob(int(blob.BlobSize), v2.Encoded)
			} else {
				e.metrics.reportFailedSubmission()
				storeCtx, cancel := context.WithTimeout(ctx, e.StoreTimeout)
				err = e.blobMetadataStore.UpdateBlobStatus(storeCtx, blobKey, v2.Failed)
				cancel()
				if err != nil {
					e.logger.Error("failed to update blob status to Failed", "blobKey", blobKey.Hex(), "err", err)
					return
				}
				// we need to remove the blobKey from the blobSet once the BlobStatus is set to FAILED
				// the Dispatcher removes the blobKey from the blobSet when batching, but blobs that are set to FAILED
				// never are batched, and therefore must be removed manually
				e.blobSet.RemoveBlob(blobKey)
				e.metrics.reportCompletedBlob(int(blob.BlobSize), v2.Failed)
			}
		})
	}

	e.metrics.reportBatchSubmissionLatency(time.Since(submissionStart))

	e.cursor = cursor

	for _, blob := range blobMetadatas {
		key, err := blob.BlobHeader.BlobKey()
		if err != nil {
			e.logger.Error("failed to get blob key", "err", err, "requestedAt", blob.RequestedAt)
			continue
		}
		e.blobSet.AddBlob(key)
	}

	e.logger.Debug("successfully submitted encoding requests", "numBlobs", len(blobMetadatas))
	return nil
}

func (e *EncodingManager) encodeBlob(ctx context.Context, blobKey corev2.BlobKey, blob *v2.BlobMetadata, blobParams *core.BlobVersionParameters) (*encoding.FragmentInfo, error) {
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
		return nil, fmt.Errorf("numAssignment (%d) cannot be greater than numRelays (%d)", numAssignment, len(availableRelays))
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
