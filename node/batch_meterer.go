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
	batchTimestamp time.Time,
) error {
	// Convert the batch to a map of account IDs to quorum IDs to symbols usage, aggregated over the batch's blobs
	usageMap, err := b.BatchRequestUsage(batch)
	if err != nil {
		return err
	}

	// validate the usages against the accounts' reservations
	return b.BatchMeterRequest(ctx, usageMap, batchTimestamp)
}

// BatchRequestUsage aggregates the usage for each account and quorum in the batch
// Returns a map of account IDs to quorum IDs to symbols usage
func (b *BatchMeterer) BatchRequestUsage(batch *corev2.Batch) (map[gethcommon.Address]map[core.QuorumID]uint64, error) {
	if batch == nil || len(batch.BlobCertificates) == 0 {
		return nil, fmt.Errorf("batch is nil or empty")
	}

	// Preallocate map with capacity for all unique accounts in the batch
	usageMap := make(map[gethcommon.Address]map[core.QuorumID]uint64, len(batch.BlobCertificates))

	for _, cert := range batch.BlobCertificates {
		if cert.BlobHeader == nil {
			return nil, fmt.Errorf("blob certificate has nil header")
		}

		accountID := cert.BlobHeader.PaymentMetadata.AccountID
		numSymbols := uint64(cert.BlobHeader.BlobCommitments.Length)
		symbolsCharged := b.symbolsCharged(numSymbols)

		if _, exists := usageMap[accountID]; !exists {
			usageMap[accountID] = make(map[core.QuorumID]uint64, len(cert.BlobHeader.QuorumNumbers))
		}

		// Add usage for each quorum in the blob
		for _, quorumID := range cert.BlobHeader.QuorumNumbers {
			usageMap[accountID][quorumID] += symbolsCharged
		}
	}

	return usageMap, nil
}

// BatchMeterRequest validates and tracks usage for a batch of requests
// For each account and quorum:
// 1. Validates against reservation limits
// 2. Checks reservation period validity
// 3. Tracks usage in period records
func (b *BatchMeterer) BatchMeterRequest(
	ctx context.Context,
	usageMap map[gethcommon.Address]map[core.QuorumID]uint64,
	batchTimestamp time.Time,
) error {
	currentPeriod := meterer.GetReservationPeriod(
		batchTimestamp.Unix(),
		b.ChainPaymentState.GetReservationWindow(),
	)

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
			if !reservation.IsActive(uint64(batchTimestamp.Unix())) {
				return fmt.Errorf("account %s has inactive reservation for quorum %d", accountID.Hex(), quorumID)
			}

			if !b.validateReservationPeriod(reservation, currentPeriod, b.ChainPaymentState.GetReservationWindow(), batchTimestamp) {
				return fmt.Errorf("account %s has invalid reservation period for quorum %d", accountID.Hex(), quorumID)
			}

			// Track and validate usage
			if err := b.incrementAndValidateUsage(accountID, quorumID, newUsage, reservation, currentPeriod); err != nil {
				return fmt.Errorf("account %s failed usage validation for quorum %d: %w", accountID.Hex(), quorumID, err)
			}
		}
	}

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

	if prevUsage >= binLimit {
		return fmt.Errorf("bin has already been filled for quorum %d", quorumID)
	}

	// Calculate how much can be added to the current bin
	toCurrent := usage
	if prevUsage+usage > binLimit {
		toCurrent = binLimit - prevUsage
	}
	periodRecord.Usage = prevUsage + toCurrent
	if periodRecord.Usage > binLimit {
		periodRecord.Usage = binLimit
	}

	// If there is overflow, handle it
	overflow := usage - toCurrent
	if overflow > 0 {
		// Only allow overflow if within 2x bin limit and within reservation window
		if usage+prevUsage > 2*binLimit || currentPeriod+b.ChainPaymentState.GetReservationWindow() > meterer.GetReservationPeriod(int64(reservation.EndTimestamp), b.ChainPaymentState.GetReservationWindow()) {
			// Rollback current period write
			periodRecord.Usage = prevUsage
			return fmt.Errorf("overflow usage exceeds bin limit for quorum %d", quorumID)
		}
		overflowPeriod := currentPeriod + b.ChainPaymentState.GetReservationWindow()
		overflowRecord := b.getOrCreatePeriodRecord(accountUsage, quorumID, overflowPeriod)
		overflowRecord.Usage += overflow
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

// getOrCreatePeriodRecord gets or creates a usage record for a specific account, quorum, and period
func (b *BatchMeterer) getOrCreatePeriodRecord(
	accountUsage *AccountUsage,
	quorumID core.QuorumID,
	period uint64,
) *pb.PeriodRecord {
	if _, ok := accountUsage.PeriodRecords[quorumID]; !ok {
		accountUsage.PeriodRecords[quorumID] = make([]*pb.PeriodRecord, b.NumBins)
		for i := range accountUsage.PeriodRecords[quorumID] {
			accountUsage.PeriodRecords[quorumID][i] = &pb.PeriodRecord{
				Index: 0,
				Usage: 0,
			}
		}
	}

	// Calculate the relative index in the circular buffer
	relativeIndex := uint32(period % uint64(b.NumBins))
	records := accountUsage.PeriodRecords[quorumID]

	// If the record at this index is for a different period, reuse it
	if records[relativeIndex].Index != uint32(period) {
		records[relativeIndex] = &pb.PeriodRecord{
			Index: uint32(period),
			Usage: 0,
		}
	}

	return records[relativeIndex]
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
func mapQuorumIDsToSortedSlice(m map[core.QuorumID]uint64) []core.QuorumID {
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
