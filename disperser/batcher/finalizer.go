package batcher

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/ethereum/go-ethereum/core/types"

	gcommon "github.com/ethereum/go-ethereum/common"
)

const maxRetries = 3
const baseDelay = 1 * time.Second

// Finalizer runs periodically to finalize blobs that have been confirmed
type Finalizer interface {
	Start(ctx context.Context)
	FinalizeBlobs(ctx context.Context) error
}

type finalizer struct {
	timeout      time.Duration
	loopInterval time.Duration
	blobStore    disperser.BlobStore
	ethClient    common.EthClient
	rpcClient    common.RPCEthClient
	logger       common.Logger
}

func NewFinalizer(timeout time.Duration, loopInterval time.Duration, blobStore disperser.BlobStore, ethClient common.EthClient, rpcClient common.RPCEthClient, logger common.Logger) Finalizer {
	return &finalizer{
		timeout:      timeout,
		loopInterval: loopInterval,
		blobStore:    blobStore,
		ethClient:    ethClient,
		rpcClient:    rpcClient,
		logger:       logger,
	}
}

func (f *finalizer) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(f.loopInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := f.FinalizeBlobs(ctx); err != nil {
					f.logger.Error("failed to finalize blobs", "err", err)
				}
			}
		}
	}()
}

// FinalizeBlobs checks the latest finalized block and marks blobs in `confirmed` state as `finalized` if their confirmation
// block number is less than or equal to the latest finalized block number.
// If it failes to process some blobs, it will log the error, skip the failed blobs, and will not return an error. The function should be invoked again to retry.
func (f *finalizer) FinalizeBlobs(ctx context.Context) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, f.timeout)
	defer cancel()
	finalizedHeader, err := f.getLatestFinalizedBlock(ctxWithTimeout)
	if err != nil {
		return fmt.Errorf("FinalizeBlobs: error getting latest finalized block: %w", err)
	}

	metadatas, err := f.blobStore.GetBlobMetadataByStatus(ctx, disperser.Confirmed)
	if err != nil {
		return fmt.Errorf("FinalizeBlobs: error getting blob headers: %w", err)
	}

	f.logger.Info("FinalizeBlobs: finalizing blobs", "numBlobs", len(metadatas), "finalizedBlockNumber", finalizedHeader.Number)

	for _, m := range metadatas {
		blobKey := m.GetBlobKey()
		confirmationMetadata, err := f.blobStore.GetBlobMetadata(ctx, blobKey)
		if err != nil {
			f.logger.Error("FinalizeBlobs: error getting confirmed metadata", "blobKey", blobKey.String(), "err", err)
			continue
		}

		// Leave as confirmed if the confirmation block is after the latest finalized block (not yet finalized)
		if uint64(confirmationMetadata.ConfirmationInfo.ConfirmationBlockNumber) > finalizedHeader.Number.Uint64() {
			continue
		}

		// confirmation block number may have changed due to reorg
		confirmationBlockNumber, err := f.getTransactionBlockNumber(ctxWithTimeout, confirmationMetadata.ConfirmationInfo.ConfirmationTxnHash)
		if err != nil {
			f.logger.Error("FinalizeBlobs: error getting transaction block number", "err", err)
			continue
		}

		// Leave as confirmed if the reorged confirmation block is after the latest finalized block (not yet finalized)
		if uint64(confirmationBlockNumber) > finalizedHeader.Number.Uint64() {
			continue
		}

		confirmationMetadata.ConfirmationInfo.ConfirmationBlockNumber = uint32(confirmationBlockNumber)
		err = f.blobStore.MarkBlobFinalized(ctx, blobKey)
		if err != nil {
			f.logger.Error("FinalizeBlobs: error marking blob as finalized", "blobKey", blobKey.String(), "err", err)
			continue
		}
	}
	f.logger.Info("FinalizeBlobs: successfully processed all finalized blobs")
	return nil
}

func (f *finalizer) getTransactionBlockNumber(ctx context.Context, hash gcommon.Hash) (uint64, error) {
	var txReceipt *types.Receipt
	var err error
	for i := 0; i < maxRetries; i++ {
		txReceipt, err = f.ethClient.TransactionReceipt(ctx, hash)
		if err == nil {
			break
		}

		retrySec := math.Pow(2, float64(i))
		f.logger.Error("Finalizer: error getting transaction", "err", err, "retrySec", retrySec, "hash", hash.Hex())
		time.Sleep(time.Duration(retrySec) * baseDelay)
	}

	if err != nil {
		return 0, fmt.Errorf("Finalizer: error getting transaction receipt after retries: %w", err)
	}

	return txReceipt.BlockNumber.Uint64(), nil
}

func (f *finalizer) getLatestFinalizedBlock(ctx context.Context) (*types.Header, error) {
	var header = types.Header{}
	var err error
	for i := 0; i < maxRetries; i++ {
		err := f.rpcClient.CallContext(ctx, &header, "eth_getBlockByNumber", "finalized", false)
		if err == nil {
			break
		}

		retrySec := math.Pow(2, float64(i))
		f.logger.Error("Finalizer: error getting latest finalized block", "err", err, "retrySec", retrySec)
		time.Sleep(time.Duration(retrySec) * baseDelay)
	}

	if err != nil {
		return nil, fmt.Errorf("Finalizer: error getting latest finalized block after retries: %w", err)
	}

	return &header, nil
}
