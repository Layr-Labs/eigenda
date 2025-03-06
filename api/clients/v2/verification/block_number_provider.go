package verification

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// BlockNumberProvider is a utility for interacting with the ethereum block number
type BlockNumberProvider struct {
	logger    logging.Logger
	ethClient common.EthClient
	// duration of interval when periodically polling the block number
	pollIntervalDuration time.Duration

	// storage shared between goroutines, containing the most recent block number observed by calling ethClient.BlockNumber()
	latestBlockNumber atomic.Uint64
	// atomic bool, so that only a single goroutine is polling the internal client with BlockNumber() calls at any given time
	pollingActive atomic.Bool
}

// NewBlockNumberProvider creates a new block number provider
func NewBlockNumberProvider(
	logger logging.Logger,
	ethClient common.EthClient,
	pollIntervalDuration time.Duration,
) *BlockNumberProvider {
	if pollIntervalDuration <= time.Duration(0) {
		logger.Warn(
			`Poll interval duration is <= 0. Therefore, any method calls made with this object that 
					rely on the internal client having reached a certain block number will fail if
					the internal client is too far behind.`,
			"pollIntervalDuration", pollIntervalDuration)
	}

	return &BlockNumberProvider{
		logger:               logger,
		ethClient:            ethClient,
		pollIntervalDuration: pollIntervalDuration,
	}
}

// MaybeWaitForBlockNumber waits until the internal eth client has advanced to a certain targetBlockNumber, unless
// configured pollInterval is <= 0, in which case this method will NOT wait for the internal client to advance.
//
// This method will check the current block number of the internal client every pollInterval duration.
// It will return nil if the internal client advances to (or past) the targetBlockNumber. It will return an error
// if the input context times out, or if any error occurs when checking the block number of the internal client.
//
// This method is synchronized in a way that, if called by multiple goroutines, only a single goroutine will actually
// poll the internal eth client for the most recent block number. The goroutine responsible for polling at a given time
// updates an atomic integer, so that all goroutines may check the most recent block without duplicating work.
func (bnp *BlockNumberProvider) MaybeWaitForBlockNumber(ctx context.Context, targetBlockNumber uint64) error {
	if bnp.pollIntervalDuration <= 0 {
		// don't wait for the internal client to advance
		return nil
	}

	if bnp.latestBlockNumber.Load() >= targetBlockNumber {
		// immediately return if the local client isn't behind the target block number
		return nil
	}

	ticker := time.NewTicker(bnp.pollIntervalDuration)
	defer ticker.Stop()

	polling := false
	if bnp.pollingActive.CompareAndSwap(false, true) {
		// no other goroutine is currently polling, so assume responsibility
		polling = true
		defer bnp.pollingActive.Store(false)
	}

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf(
				"timed out waiting for block number %d (latest block number observed was %d): %w",
				targetBlockNumber, bnp.latestBlockNumber.Load(), ctx.Err())
		case <-ticker.C:
			if bnp.latestBlockNumber.Load() >= targetBlockNumber {
				return nil
			}

			if bnp.pollingActive.CompareAndSwap(false, true) {
				// no other goroutine is currently polling, so assume responsibility
				polling = true
				defer bnp.pollingActive.Store(false)
			}

			if polling {
				fetchedBlockNumber, err := bnp.FetchLatestBlockNumber(ctx)
				if err != nil {
					bnp.logger.Debug(
						"ethClient.BlockNumber returned an error",
						"targetBlockNumber", targetBlockNumber,
						"latestBlockNumber", bnp.latestBlockNumber.Load(),
						"error", err)

					// tolerate some failures here. if failure continues for too long, it will be caught by the timeout
					continue
				}

				if fetchedBlockNumber >= targetBlockNumber {
					return nil
				}
			}

			bnp.logger.Debug(
				"local client is behind the reference block number",
				"targetBlockNumber", targetBlockNumber,
				"actualBlockNumber", bnp.latestBlockNumber.Load())
		}
	}
}

// FetchLatestBlockNumber fetches the latest block number from the eth client, and returns it.
//
// This method atomically stores the latest block number for internal use.
//
// Calling this method doesn't have an impact on the cadence of the standard block number polling that occurs
// in MaybeWaitForBlockNumber.
func (bnp *BlockNumberProvider) FetchLatestBlockNumber(ctx context.Context) (uint64, error) {
	blockNumber, err := bnp.ethClient.BlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("get block number from eth client: %w", err)
	}

	bnp.latestBlockNumber.Store(blockNumber)

	return blockNumber, nil
}
