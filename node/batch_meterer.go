package node

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	meterer "github.com/Layr-Labs/eigenda/core/meterer"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// AccountUsage tracks usage for a specific account across quorums
type AccountUsage struct {
	// Circular buffer of usage records, with minimum length MinNumBins
	PeriodRecords map[core.QuorumID][]*pb.PeriodRecord
	// Lock to protect concurrent access to usage records
	Lock sync.RWMutex
}

// BatchMeterer handles metering for batches of requests that may contain multiple accounts
type BatchMeterer struct {
	// Configuration for the batch meterer
	Config meterer.Config

	// ChainPaymentState reads on-chain payment state
	ChainPaymentState meterer.OnchainPayment

	// AccountUsages tracks in-memory usage for accounts and quorums
	AccountUsages sync.Map // map[gethcommon.Address]*AccountUsage

	// Number of bins to track for each account's reservation usage
	NumBins uint32

	logger logging.Logger
}

// NewBatchMeterer creates a new batch meterer
func NewBatchMeterer(
	config meterer.Config,
	paymentChainState meterer.OnchainPayment,
	numBins uint32,
	logger logging.Logger,
) *BatchMeterer {
	return &BatchMeterer{
		Config:            config,
		ChainPaymentState: paymentChainState,
		NumBins:           max(numBins, meterer.MinNumBins),
		logger:            logger.With("component", "BatchMeterer"),
	}
}

// Start starts to periodically refresh the on-chain state
func (b *BatchMeterer) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(b.Config.UpdateInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := b.ChainPaymentState.RefreshOnchainPaymentState(ctx); err != nil {
					b.logger.Error("Failed to refresh on-chain state", "error", err)
				}
				b.logger.Debug("Refreshed on-chain state")
			case <-ctx.Done():
				return
			}
		}
	}()
}

// MeterBatch meters a batch of blobs directly from a corev2.Batch
// Returns an error if any part of the batch is invalid
func (b *BatchMeterer) MeterBatch(
	ctx context.Context,
	batch *corev2.Batch,
	batchReceivedAt time.Time,
) error {
	b.logger.Info("MeterBatch", "batch", batch, "batchReceivedAt", batchReceivedAt)
	// Convert the batch to a map of account IDs to quorum IDs to symbols usage, aggregated over the batch's blobs
	usageMap, err := b.BatchRequestUsage(batch)
	if err != nil {
		b.logger.Error("BatchRequestUsage", "error", err)
		return err
	}

	// validate the usages against the accounts' reservations
	return b.BatchMeterRequest(ctx, usageMap, batchReceivedAt)
}

// BatchRequestUsage aggregates the usage for each account and quorum in the batch
// Returns a map of account IDs to quorum IDs to period index to symbols usage
// TODO(hopeyen): use batch blob header payment metadata to get reservation period to charge against
func (b *BatchMeterer) BatchRequestUsage(batch *corev2.Batch) (map[gethcommon.Address]map[core.QuorumID]map[uint64]uint64, error) {
	b.logger.Info("BatchRequestUsage", "batch", batch)
	if batch == nil || len(batch.BlobCertificates) == 0 {
		return nil, fmt.Errorf("batch is nil or empty")
	}

	// Preallocate map with capacity for all unique accounts in the batch
	usageMap := make(map[gethcommon.Address]map[core.QuorumID]map[uint64]uint64, len(batch.BlobCertificates))

	for _, cert := range batch.BlobCertificates {
		if cert.BlobHeader == nil {
			return nil, fmt.Errorf("blob certificate has nil header")
		}

		accountID := cert.BlobHeader.PaymentMetadata.AccountID
		numSymbols := uint64(cert.BlobHeader.BlobCommitments.Length)
		symbolsCharged := b.symbolsCharged(numSymbols)

		if _, exists := usageMap[accountID]; !exists {
			usageMap[accountID] = make(map[core.QuorumID]map[uint64]uint64, len(cert.BlobHeader.QuorumNumbers))
		}

		// Add usage for each quorum in the blob
		for _, quorumID := range cert.BlobHeader.QuorumNumbers {
			if _, exists := usageMap[accountID][quorumID]; !exists {
				usageMap[accountID][quorumID] = make(map[uint64]uint64)
			}
			// Use current period as the index
			currentPeriod := meterer.GetReservationPeriodByNanosecond(cert.BlobHeader.PaymentMetadata.Timestamp, b.ChainPaymentState.GetReservationWindow())
			usageMap[accountID][quorumID][currentPeriod] += symbolsCharged
		}
	}

	b.logger.Info("BatchRequestUsage", "usageMap", usageMap)
	return usageMap, nil
}

// BatchMeterRequest validates and tracks usage for a batch of requests
// For each account and quorum:
// 1. Validates against reservation limits
// 2. Checks reservation period validity
// 3. Tracks usage in period records
func (b *BatchMeterer) BatchMeterRequest(
	ctx context.Context,
	usageMap map[gethcommon.Address]map[core.QuorumID]map[uint64]uint64,
	batchReceivedAt time.Time,
) error {
	b.logger.Info("BatchMeterRequest", "usageMap", usageMap, "batchReceivedAt", batchReceivedAt)
	// Process each account's usage
	for accountID, quorumUsages := range usageMap {
		// Get active reservations for this account's quorums
		reservations, err := b.ChainPaymentState.GetReservedPaymentByAccountAndQuorums(
			ctx, accountID, mapQuorumIDsToSortedSlice(quorumUsages),
		)
		if err != nil {
			return fmt.Errorf("failed to get reservations for account %s: %w", accountID.Hex(), err)
		}

		// Validate each quorum's usage against its reservation
		for quorumID, newUsage := range quorumUsages {
			reservation, ok := reservations[quorumID]
			if !ok {
				return fmt.Errorf("account %s has no reservation for quorum %d", accountID.Hex(), quorumID)
			}

			// Check reservation validity
			if !reservation.IsActive(uint64(batchReceivedAt.Unix())) {
				return fmt.Errorf("account %s has inactive reservation for quorum %d", accountID.Hex(), quorumID)
			}

			for periodIndex, usage := range newUsage {
				// TODO(hopeyen): current period must become quorum-specific
				if !b.validateReservationPeriod(reservation, periodIndex, b.ChainPaymentState.GetReservationWindow(), batchReceivedAt) {
					return fmt.Errorf("account %s has invalid reservation period for quorum %d", accountID.Hex(), quorumID)
				}

				// Track and validate usage
				if err := b.incrementAndValidateUsage(accountID, quorumID, usage, reservation, periodIndex); err != nil {
					return fmt.Errorf("account %s failed usage validation for quorum %d: %w", accountID.Hex(), quorumID, err)
				}
			}
		}
	}

	//TODO(hopeyen): rollback to init state if any update was not successful

	return nil
}

// incrementAndValidateUsage validates and increments usage for an account and quorum
// Uses a write-first approach for atomicity:
// 1. Writes changes first
// 2. Validates the new state
// 3. Rolls back if validation fails
func (b *BatchMeterer) incrementAndValidateUsage(
	accountID gethcommon.Address,
	quorumID core.QuorumID,
	usage uint64,
	reservation *core.ReservedPayment,
	currentPeriod uint64,
) error {
	accountUsage := b.getOrCreateAccountUsage(accountID)
	accountUsage.Lock.Lock()
	defer accountUsage.Lock.Unlock()

	periodRecord := b.getOrCreatePeriodRecord(accountUsage, quorumID, currentPeriod)
	prevUsage := periodRecord.Usage
	binLimit := reservation.SymbolsPerSecond * b.ChainPaymentState.GetReservationWindow()

	b.logger.Info("incrementAndValidateUsage",
		"accountID", accountID.Hex(),
		"quorumID", quorumID,
		"prevUsage", prevUsage,
		"newUsage", usage,
		"binLimit", binLimit,
		"currentPeriod", currentPeriod,
	)

	// Write-first: increment usage
	periodRecord.Usage += usage
	newUsage := periodRecord.Usage

	b.logger.Info("after increment",
		"newTotalUsage", newUsage,
		"isWithinBinLimit", newUsage <= binLimit,
	)

	// If new usage is within bin limit, done
	if newUsage <= binLimit {
		return nil
	}

	// If bin was already full before this increment, revert and error
	if prevUsage >= binLimit {
		b.logger.Info("bin already full, reverting",
			"prevUsage", prevUsage,
			"binLimit", binLimit,
		)
		periodRecord.Usage = prevUsage
		return fmt.Errorf("bin has already been filled for quorum %d", quorumID)
	}

	// If new usage is more than 2x bin limit, revert and error
	if newUsage > 2*binLimit {
		b.logger.Info("usage exceeds 2x bin limit, reverting",
			"newUsage", newUsage,
			"2xBinLimit", 2*binLimit,
		)
		periodRecord.Usage = prevUsage
		return fmt.Errorf("overflow usage exceeds bin limit for quorum %d", quorumID)
	}

	// Check if overflow period is within reservation window
	nextPeriod := currentPeriod + b.ChainPaymentState.GetReservationWindow()
	endPeriod := meterer.GetReservationPeriod(int64(reservation.EndTimestamp), b.ChainPaymentState.GetReservationWindow())

	b.logger.Info("checking overflow period",
		"nextPeriod", nextPeriod,
		"endPeriod", endPeriod,
		"isWithinWindow", nextPeriod < endPeriod,
	)

	if nextPeriod >= endPeriod {
		b.logger.Info("overflow period outside window, reverting")
		periodRecord.Usage = prevUsage
		return fmt.Errorf("overflow period exceeds reservation window for quorum %d", quorumID)
	}

	// Move overflow to next period
	overflow := newUsage - binLimit
	periodRecord.Usage = binLimit
	overflowRecord := b.getOrCreatePeriodRecord(accountUsage, quorumID, nextPeriod)
	overflowPrev := overflowRecord.Usage
	overflowRecord.Usage += overflow

	b.logger.Info("handling overflow",
		"overflow", overflow,
		"overflowPrev", overflowPrev,
		"overflowNew", overflowRecord.Usage,
		"overflowBinLimit", binLimit,
	)

	// If overflow in next period exceeds bin limit, revert both and error
	if overflowRecord.Usage > binLimit {
		b.logger.Info("overflow in next period exceeds limit, reverting both")
		// revert both writes
		periodRecord.Usage = prevUsage
		overflowRecord.Usage = overflowPrev
		return fmt.Errorf("overflow usage exceeds bin limit for quorum %d in next period", quorumID)
	}

	return nil
}

// getOrCreateAccountUsage gets or creates a usage record for an account
func (b *BatchMeterer) getOrCreateAccountUsage(accountID gethcommon.Address) *AccountUsage {
	value, ok := b.AccountUsages.Load(accountID)
	if ok {
		return value.(*AccountUsage)
	}

	accountUsage := &AccountUsage{
		PeriodRecords: make(map[core.QuorumID][]*pb.PeriodRecord),
	}
	b.AccountUsages.Store(accountID, accountUsage)
	return accountUsage
}

// getOrCreatePeriodRecord gets or creates a period record for an account and quorum
func (b *BatchMeterer) getOrCreatePeriodRecord(
	accountUsage *AccountUsage,
	quorumID core.QuorumID,
	period uint64,
) *pb.PeriodRecord {
	// Initialize quorum records if needed
	if _, exists := accountUsage.PeriodRecords[quorumID]; !exists {
		accountUsage.PeriodRecords[quorumID] = make([]*pb.PeriodRecord, b.NumBins)
	}

	// Calculate relative index in circular buffer
	relativeIndex := uint32((period / b.ChainPaymentState.GetReservationWindow()) % uint64(b.NumBins))
	b.logger.Info("getOrCreatePeriodRecord", "relativeIndex", relativeIndex, "period", period, "reservationWindow", b.ChainPaymentState.GetReservationWindow(), "numBins", b.NumBins)

	// Initialize record if needed
	if accountUsage.PeriodRecords[quorumID][relativeIndex] == nil {
		accountUsage.PeriodRecords[quorumID][relativeIndex] = &pb.PeriodRecord{
			Index: uint32(period),
			Usage: 0,
		}
	}

	// Reset usage if this is a new period
	if accountUsage.PeriodRecords[quorumID][relativeIndex].Index < uint32(period) {
		accountUsage.PeriodRecords[quorumID][relativeIndex].Index = uint32(period)
		accountUsage.PeriodRecords[quorumID][relativeIndex].Usage = 0
	}

	return accountUsage.PeriodRecords[quorumID][relativeIndex]
}

// symbolsCharged returns the number of symbols charged for a given data length
func (b *BatchMeterer) symbolsCharged(numSymbols uint64) uint64 {
	minSymbols := b.ChainPaymentState.GetMinNumSymbols()
	if numSymbols <= minSymbols {
		return minSymbols
	}
	// Round up to the nearest multiple of MinNumSymbols
	return core.RoundUpDivide(numSymbols, minSymbols) * minSymbols
}

// validateReservationPeriod checks if the provided reservation period is valid
// Based on meterer.ValidateReservationPeriod
// requestReservationPeriod is the period the request indicates the period it is using
// receivedAt is the timestamp the request was received at and will be used to validate requestReservationPeriod
func (b *BatchMeterer) validateReservationPeriod(reservation *core.ReservedPayment, requestReservationPeriod uint64, reservationWindow uint64, receivedAt time.Time) bool {
	currentReservationPeriod := meterer.GetReservationPeriod(receivedAt.Unix(), reservationWindow)

	// Valid reservation periods are either the current bin or the previous bin
	isCurrentOrPreviousPeriod := requestReservationPeriod == currentReservationPeriod ||
		requestReservationPeriod == (currentReservationPeriod-reservationWindow)

	startPeriod := meterer.GetReservationPeriod(int64(reservation.StartTimestamp), reservationWindow)
	endPeriod := meterer.GetReservationPeriod(int64(reservation.EndTimestamp), reservationWindow)

	isWithinReservationWindow := startPeriod <= requestReservationPeriod && requestReservationPeriod < endPeriod

	return isCurrentOrPreviousPeriod && isWithinReservationWindow
}

// mapQuorumIDsToSortedSlice converts map keys of type QuorumID to a sorted slice
// This ensures consistent ordering for deterministic test behavior
func mapQuorumIDsToSortedSlice(m map[core.QuorumID]map[uint64]uint64) []core.QuorumID {
	keys := make([]core.QuorumID, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	// Sort QuorumIDs in ascending order using the standard library
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	return keys
}
