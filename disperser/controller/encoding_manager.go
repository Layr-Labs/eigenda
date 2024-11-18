package controller

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
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
}

// EncodingManager is responsible for pulling queued blobs from the blob
// metadata store periodically and encoding them. It receives the encoder responses
// and creates BlobCertificates.
type EncodingManager struct {
	EncodingManagerConfig

	// components
	blobMetadataStore *blobstore.BlobMetadataStore
	pool              common.WorkerPool
	encodingClient    disperser.EncoderClientV2
	chainReader       core.Reader
	logger            logging.Logger

	// state
	lastUpdatedAt uint64
}

func NewEncodingManager(
	config EncodingManagerConfig,
	blobMetadataStore *blobstore.BlobMetadataStore,
	pool common.WorkerPool,
	encodingClient disperser.EncoderClientV2,
	chainReader core.Reader,
	logger logging.Logger,
) (*EncodingManager, error) {
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

		lastUpdatedAt: 0,
	}, nil
}

func (e *EncodingManager) Start(ctx context.Context) error {
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
	blobMetadatas, err := e.blobMetadataStore.GetBlobMetadataByStatus(ctx, v2.Queued, e.lastUpdatedAt)
	if err != nil {
		return err
	}

	if len(blobMetadatas) == 0 {
		return errNoBlobsToEncode
	}

	for _, blob := range blobMetadatas {
		blob := blob
		blobKey, err := blob.BlobHeader.BlobKey()
		if err != nil {
			e.logger.Error("failed to get blob key", "err", err, "requestedAt", blob.RequestedAt, "paymentMetadata", blob.BlobHeader.PaymentMetadata)
			continue
		}
		e.lastUpdatedAt = blob.UpdatedAt

		// Encode the blobs
		e.pool.Submit(func() {
			for i := 0; i < e.NumEncodingRetries+1; i++ {
				encodingCtx, cancel := context.WithTimeout(ctx, e.EncodingRequestTimeout)
				fragmentInfo, err := e.encodeBlob(encodingCtx, blobKey, blob)
				cancel()
				if err != nil {
					e.logger.Error("failed to encode blob", "blobKey", blobKey.Hex(), "err", err)
					continue
				}
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
					e.logger.Error("failed to put blob certificate", "err", err)
					continue
				}

				storeCtx, cancel = context.WithTimeout(ctx, e.StoreTimeout)
				err = e.blobMetadataStore.UpdateBlobStatus(storeCtx, blobKey, v2.Encoded)
				cancel()
				if err == nil || errors.Is(err, dispcommon.ErrAlreadyExists) {
					// Successfully updated the status to Encoded
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

	return nil
}

func (e *EncodingManager) encodeBlob(ctx context.Context, blobKey corev2.BlobKey, blob *v2.BlobMetadata) (*encoding.FragmentInfo, error) {
	encodingParams, err := blob.BlobHeader.GetEncodingParams()
	if err != nil {
		return nil, fmt.Errorf("failed to get encoding params: %w", err)
	}
	return e.encodingClient.EncodeBlob(ctx, blobKey, encodingParams)
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
