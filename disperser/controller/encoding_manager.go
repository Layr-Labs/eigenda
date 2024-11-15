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
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser"
	dispcommon "github.com/Layr-Labs/eigenda/disperser/common"
	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

var errNoBlobsToEncode = errors.New("no blobs to encode")

type EncodingManagerConfig struct {
	PullInterval time.Duration

	EncodingRequestTimeout time.Duration
	StoreTimeout           time.Duration
	// NumEncodingRetries defines how many times the encoding will be retried
	NumEncodingRetries int
	// NumRelayAssignment defines how many relays will be assigned to a blob
	NumRelayAssignment uint16
	// AvailableRelays is a list of available relays
	AvailableRelays []corev2.RelayKey
	// EncoderAddress is the address of the encoder
	EncoderAddress string
	// MaxNumBlobsPerIteration is the maximum number of blobs to encode per iteration
	MaxNumBlobsPerIteration int32
	// OnchainStateRefreshInterval is the interval at which the onchain state is refreshed
	OnchainStateRefreshInterval time.Duration
}

// EncodingManager is responsible for pulling queued blobs from the blob
// metadata store periodically and encoding them. It receives the encoder responses
// and creates BlobCertificates.
type EncodingManager struct {
	*EncodingManagerConfig

	// components
	blobMetadataStore *blobstore.BlobMetadataStore
	pool              common.WorkerPool
	encodingClient    disperser.EncoderClientV2
	chainReader       core.Reader
	logger            logging.Logger

	// state
	cursor                *blobstore.StatusIndexCursor
	blobVersionParameters atomic.Pointer[corev2.BlobVersionParameterMap]
}

func NewEncodingManager(
	config *EncodingManagerConfig,
	blobMetadataStore *blobstore.BlobMetadataStore,
	pool common.WorkerPool,
	encodingClient disperser.EncoderClientV2,
	chainReader core.Reader,
	logger logging.Logger,
) (*EncodingManager, error) {
	if config.NumRelayAssignment < 1 ||
		len(config.AvailableRelays) == 0 ||
		config.MaxNumBlobsPerIteration < 1 {
		return nil, fmt.Errorf("invalid encoding manager config")
	}
	if int(config.NumRelayAssignment) > len(config.AvailableRelays) {
		return nil, fmt.Errorf("NumRelayAssignment (%d) cannot be greater than NumRelays (%d)", config.NumRelayAssignment, len(config.AvailableRelays))
	}
	return &EncodingManager{
		EncodingManagerConfig: config,
		blobMetadataStore:     blobMetadataStore,
		pool:                  pool,
		encodingClient:        encodingClient,
		chainReader:           chainReader,
		logger:                logger.With("component", "EncodingManager"),

		cursor: nil,
	}, nil
}

func (e *EncodingManager) Start(ctx context.Context) error {
	// Refresh blob version parameters
	err := e.refreshBlobVersionParams(ctx)
	if err != nil {
		return fmt.Errorf("failed to refresh blob version parameters: %w", err)
	}

	go func() {
		ticker := time.NewTicker(e.EncodingManagerConfig.OnchainStateRefreshInterval)
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
						e.logger.Warn("no blobs to encode")
					} else {
						e.logger.Error("failed to process a batch", "err", err)
					}
				}
			}
		}
	}()

	return nil
}

func (e *EncodingManager) HandleBatch(ctx context.Context) error {
	// Get a batch of blobs to encode
	blobMetadatas, cursor, err := e.blobMetadataStore.GetBlobMetadataByStatusPaginated(ctx, v2.Queued, e.cursor, e.MaxNumBlobsPerIteration)
	if err != nil {
		return err
	}

	if len(blobMetadatas) == 0 {
		return errNoBlobsToEncode
	}

	blobVersionParams := e.blobVersionParameters.Load()
	if blobVersionParams == nil {
		return fmt.Errorf("blob version parameters is nil")
	}

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
			for i := 0; i < e.NumEncodingRetries+1; i++ {
				e.logger.Debug("encoding blob", "blobKey", blobKey.Hex())
				encodingCtx, cancel := context.WithTimeout(ctx, e.EncodingRequestTimeout)
				fragmentInfo, err := e.encodeBlob(encodingCtx, blobKey, blob, blobParams)
				cancel()
				if err != nil {
					e.logger.Error("failed to encode blob", "blobKey", blobKey.Hex(), "err", err)
					continue
				}
				e.logger.Debug("successfully encoded blob", "blobKey", blobKey.Hex(), "fragmentInfo", fragmentInfo)
				relayKeys, err := GetRelayKeys(e.NumRelayAssignment, e.AvailableRelays)
				if err != nil {
					e.logger.Error("failed to get relay keys", "err", err)
					// Stop retrying
					break
				}
				cert := &corev2.BlobCertificate{
					BlobHeader: blob.BlobHeader,
					RelayKeys:  relayKeys,
				}

				storeCtx, cancel := context.WithTimeout(ctx, e.StoreTimeout)
				err = e.blobMetadataStore.PutBlobCertificate(storeCtx, cert, fragmentInfo)
				cancel()
				if err != nil && !errors.Is(err, dispcommon.ErrAlreadyExists) {
					e.logger.Error("failed to put blob certificate", "blobKey", blobKey.Hex(), "err", err)
					continue
				} else if err != nil {
					e.logger.Warn("blob certificate already exists", "blobKey", blobKey.Hex(), "err", err)
				}

				storeCtx, cancel = context.WithTimeout(ctx, e.StoreTimeout)
				err = e.blobMetadataStore.UpdateBlobStatus(storeCtx, blobKey, v2.Encoded)
				cancel()
				if err == nil || errors.Is(err, dispcommon.ErrAlreadyExists) {
					// Successfully updated the status to Encoded
					e.logger.Debug("successfully updated blob to encoded status", "blobKey", blobKey.Hex())
					return
				}

				e.logger.Error("failed to update blob status to Encoded", "blobKey", blobKey.Hex(), "err", err)
				time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second) // Wait before retrying
			}

			storeCtx, cancel := context.WithTimeout(ctx, e.StoreTimeout)
			err = e.blobMetadataStore.UpdateBlobStatus(storeCtx, blobKey, v2.Failed)
			cancel()
			if err != nil {
				e.logger.Error("failed to update blob status to Failed", "blobKey", blobKey.Hex(), "err", err)
				return
			}
		})
	}

	if cursor != nil {
		e.cursor = cursor
	}
	return nil
}

func (e *EncodingManager) encodeBlob(ctx context.Context, blobKey corev2.BlobKey, blob *v2.BlobMetadata, blobParams *core.BlobVersionParameters) (*encoding.FragmentInfo, error) {
	encodingParams, err := blob.BlobHeader.GetEncodingParams(blobParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get encoding params: %w", err)
	}
	return e.encodingClient.EncodeBlob(ctx, blobKey, encodingParams)
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
	relayKeys := availableRelays
	// shuffle relay keys
	for i := len(relayKeys) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		relayKeys[i], relayKeys[j] = relayKeys[j], relayKeys[i]
	}

	return relayKeys[:numAssignment], nil
}
