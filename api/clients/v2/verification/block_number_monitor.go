package verification

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var blockNumberMonitorTracer = otel.Tracer(
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification/block_number_monitor")

// BlockNumberMonitor is a utility for waiting for a certain ethereum block number
//
// This utility is used by the CertVerifierAddressProvider implementations to ensure that the client
// has reached a sufficient block height before making queries about block-specific state
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
) (*BlockNumberMonitor, error) {
	if pollIntervalDuration <= time.Duration(0) {
		return nil, fmt.Errorf("input pollIntervalDuration (%v) must be greater than zero", pollIntervalDuration)
	}

	return &BlockNumberMonitor{
		logger:               logger,
		ethClient:            ethClient,
		pollIntervalDuration: pollIntervalDuration,
	}, nil
}

// WaitForBlockNumber waits until the internal eth client has advanced to a certain targetBlockNumber.
//
// This method will check the current block number of the internal client every pollInterval duration.
// It will return nil if the internal client advances to (or past) the targetBlockNumber. It will return an error
// if the input context times out, or if any error occurs when checking the block number of the internal client.
//
// This method is synchronized in a way that, if called by multiple goroutines, only a single goroutine will actually
// poll the internal eth client for the most recent block number. The goroutine responsible for polling at a given time
// updates an atomic integer, so that all goroutines may check the most recent block without duplicating work.
func (bnm *BlockNumberMonitor) WaitForBlockNumber(ctx context.Context, targetBlockNumber uint64) error {
	ctx, span := blockNumberMonitorTracer.Start(ctx, "BlockNumberMonitor.WaitForBlockNumber",
		trace.WithAttributes(
			attribute.Int64("target_block_number", int64(targetBlockNumber))))
	defer span.End()

	if bnm.pollIntervalDuration <= 0 {
		err := fmt.Errorf(
			"pollIntervalDuration is <= 0: you ought to be using the provided constructor, which checks this")
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid poll interval")
		return err
	}

	if bnm.latestBlockNumber.Load() >= targetBlockNumber {
		// immediately return if the local client isn't behind the target block number
		span.SetStatus(codes.Ok, "block number already reached")
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
			err := fmt.Errorf(
				"timed out waiting for block number %d (latest block number observed was %d): %w",
				targetBlockNumber, bnm.latestBlockNumber.Load(), ctx.Err())
			span.RecordError(err)
			span.SetStatus(codes.Error, "timeout waiting for block number")
			return err
		case <-ticker.C:
			if bnm.latestBlockNumber.Load() >= targetBlockNumber {
				span.SetStatus(codes.Ok, "block number reached")
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
					bnm.logger.Debug(
						"ethClient.BlockNumber returned an error",
						"targetBlockNumber", targetBlockNumber,
						"latestBlockNumber", bnm.latestBlockNumber.Load(),
						"error", err)

					// tolerate some failures here. if failure continues for too long, it will be caught by the timeout
					continue
				}

				bnm.latestBlockNumber.Store(blockNumber)
				span.SetAttributes(attribute.Int64("actual_block_number", int64(blockNumber)))

				if blockNumber >= targetBlockNumber {
					span.SetStatus(codes.Ok, "block number reached")
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
