package verification

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// BlockNumberMonitor is a utility for waiting for a certain ethereum block number
//
// TODO: this utility is not currently in use, but DO NOT delete it. It will be necessary for the upcoming
//  CertVerifierRouter effort
type BlockNumberMonitor struct {
	logger    logging.Logger
	ethClient common.EthClient
	// duration of interval when periodically polling the block number
	pollIntervalDuration time.Duration

	// storage shared between goroutines, containing the most recent block number observed by calling ethClient.BlockNumber()
	latestBlockNumber atomic.Uint64
	// atomic bool, so that only a single goroutine is polling the internal client with BlockNumber() calls at any given time
	pollingActive atomic.Bool
}

// NewBlockNumberMonitor creates a new block number monitor
func NewBlockNumberMonitor(
	logger logging.Logger,
	ethClient common.EthClient,
	pollIntervalDuration time.Duration,
) *BlockNumberMonitor {
	if pollIntervalDuration <= time.Duration(0) {
		logger.Warn(
			`Poll interval duration is <= 0. Therefore, any method calls made with this object that 
					rely on the internal client having reached a certain block number will fail if
					the internal client is too far behind.`,
			"pollIntervalDuration", pollIntervalDuration)
	}

	return &BlockNumberMonitor{
		logger:               logger,
		ethClient:            ethClient,
		pollIntervalDuration: pollIntervalDuration,
	}
}

// WaitForBlockNumber waits until the internal eth client has advanced to a certain targetBlockNumber.
//
// This method will check the current block number of the internal client every pollInterval duration.
// It will return nil if the internal client advances to (or past) the targetBlockNumber. It will return an error
// if the input context times out, or if any error occurs when checking the block number of the internal client.
//
// If configured pollInterval is <= 0, this method will NOT wait for the internal client to advance, and instead will
// simply return nil.
//
// This method is synchronized in a way that, if called by multiple goroutines, only a single goroutine will actually
// poll the internal eth client for the most recent block number. The goroutine responsible for polling at a given time
// updates an atomic integer, so that all goroutines may check the most recent block without duplicating work.
func (bnm *BlockNumberMonitor) WaitForBlockNumber(ctx context.Context, targetBlockNumber uint64) error {
	if bnm.pollIntervalDuration <= 0 {
		// don't wait for the internal client to advance. duration <= 0 would cause the ticker to panic
		return nil
	}

	if bnm.latestBlockNumber.Load() >= targetBlockNumber {
		// immediately return if the local client isn't behind the target block number
		return nil
	}

	ticker := time.NewTicker(bnm.pollIntervalDuration)
	defer ticker.Stop()

	polling := false
	if bnm.pollingActive.CompareAndSwap(false, true) {
		// no other goroutine is currently polling, so assume responsibility
		polling = true
		defer bnm.pollingActive.Store(false)
	}

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf(
				"timed out waiting for block number %d (latest block number observed was %d): %w",
				targetBlockNumber, bnm.latestBlockNumber.Load(), ctx.Err())
		case <-ticker.C:
			if bnm.latestBlockNumber.Load() >= targetBlockNumber {
				return nil
			}

			if bnm.pollingActive.CompareAndSwap(false, true) {
				// no other goroutine is currently polling, so assume responsibility
				polling = true
				defer bnm.pollingActive.Store(false)
			}

			if polling {
				blockNumber, err := bnm.ethClient.BlockNumber(ctx)
				if err != nil {
					return fmt.Errorf("get block number from eth client: %w", err)
				}

				bnm.latestBlockNumber.Store(blockNumber)

				if err != nil {
					bnm.logger.Debug(
						"ethClient.BlockNumber returned an error",
						"targetBlockNumber", targetBlockNumber,
						"latestBlockNumber", bnm.latestBlockNumber.Load(),
						"error", err)

					// tolerate some failures here. if failure continues for too long, it will be caught by the timeout
					continue
				}

				if blockNumber >= targetBlockNumber {
					return nil
				}
			}

			bnm.logger.Debug(
				"local client is behind the reference block number",
				"targetBlockNumber", targetBlockNumber,
				"actualBlockNumber", bnm.latestBlockNumber.Load())
		}
	}
}
