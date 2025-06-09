package node

import (
	"context"
	"fmt"
	"slices"
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

// UpdateRecord tracks a successful usage update for potential rollback
type UpdateRecord struct {
	accountID gethcommon.Address
	quorumID  core.QuorumID
	period    uint64
	usage     uint64 // The usage value to restore during rollback
}

// newUpdateRecord creates a new UpdateRecord with validation
func newUpdateRecord(accountID gethcommon.Address, quorumID core.QuorumID, period, usage uint64) (UpdateRecord, error) {
	record := UpdateRecord{
		accountID: accountID,
		quorumID:  quorumID,
		period:    period,
		usage:     usage,
	}
	return record, nil
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
// Returns an array of UpdateRecord representing the usage updates
func (b *BatchMeterer) BatchRequestUsage(batch *corev2.Batch) ([]UpdateRecord, error) {
	b.logger.Info("BatchRequestUsage", "batch", batch)
	if batch == nil || len(batch.BlobCertificates) == 0 {
		return nil, fmt.Errorf("batch is nil or empty")
	}

	// Use a map to track existing updates by account, quorum, and period
	updatesMap := make(map[string]*UpdateRecord)

	for _, cert := range batch.BlobCertificates {
		if cert.BlobHeader == nil {
			return nil, fmt.Errorf("blob certificate has nil header")
		}

		accountID := cert.BlobHeader.PaymentMetadata.AccountID
		numSymbols := uint64(cert.BlobHeader.BlobCommitments.Length)
		symbolsCharged := b.symbolsCharged(numSymbols)

		// Add usage for each quorum in the blob
		for _, quorumID := range cert.BlobHeader.QuorumNumbers {
			// Use current period as the index
			currentPeriod := meterer.GetReservationPeriodByNanosecond(cert.BlobHeader.PaymentMetadata.Timestamp, b.ChainPaymentState.GetReservationWindow())

			// Create a unique key for this account/quorum/period combination
			key := fmt.Sprintf("%s_%d_%d", accountID.Hex(), quorumID, currentPeriod)

			if existingRecord, exists := updatesMap[key]; exists {
				// If record exists, add to its usage
				existingRecord.usage += symbolsCharged
			} else {
				// Create new record if none exists
				record, err := newUpdateRecord(accountID, quorumID, currentPeriod, symbolsCharged)
				if err != nil {
					return nil, fmt.Errorf("failed to create update record: %w", err)
				}
				updatesMap[key] = &record
			}
		}
	}

	// Convert map to slice
	updates := make([]UpdateRecord, 0, len(updatesMap))
	for _, record := range updatesMap {
		updates = append(updates, *record)
	}

	b.logger.Info("BatchRequestUsage", "numUpdates", len(updates))
	return updates, nil
}

// BatchMeterRequest validates and tracks usage for a batch of requests
// For each account and quorum:
// 1. Validates against reservation limits
// 2. Checks reservation period validity
// 3. Tracks usage in period records
func (b *BatchMeterer) BatchMeterRequest(
	ctx context.Context,
	updates []UpdateRecord,
	batchReceivedAt time.Time,
) error {
	b.logger.Info("BatchMeterRequest", "numUpdates", len(updates), "batchReceivedAt", batchReceivedAt)

	// Track successful updates for potential rollback
	successfulUpdates := make([]UpdateRecord, 0)

	// Group updates by account for efficient reservation lookup
	accountUpdates := make(map[gethcommon.Address][]UpdateRecord)
	for _, update := range updates {
		accountUpdates[update.accountID] = append(accountUpdates[update.accountID], update)
	}

	// Process each account's updates
	for accountID, accountUpdates := range accountUpdates {
		b.logger.Info("Processing account updates", "accountID", accountID.Hex(), "numUpdates", len(accountUpdates))

		// Get unique quorum IDs for this account
		quorumSet := make(map[core.QuorumID]struct{})
		for _, update := range accountUpdates {
			quorumSet[update.quorumID] = struct{}{}
		}
		quorumIDs := make([]core.QuorumID, 0, len(quorumSet))
		for quorumID := range quorumSet {
			quorumIDs = append(quorumIDs, quorumID)
		}
		slices.Sort(quorumIDs)
		b.logger.Info("Sorted quorum IDs", "quorumIDs", quorumIDs)

		// Get active reservations for this account's quorums
		reservations, err := b.ChainPaymentState.GetReservedPaymentByAccountAndQuorums(ctx, accountID, quorumIDs)
		if err != nil {
			b.rollbackUpdates(successfulUpdates)
			return fmt.Errorf("failed to get reservations for account %s: %w", accountID.Hex(), err)
		}
		b.logger.Info("Got reservations", "reservations", reservations)

		// Process each update for this account
		for _, update := range accountUpdates {
			b.logger.Debug("Processing update", "update", update)
			reservation, ok := reservations[update.quorumID]
			if !ok {
				b.rollbackUpdates(successfulUpdates)
				return fmt.Errorf("account %s has no reservation for quorum %d", accountID.Hex(), update.quorumID)
			}

			// Check if reservation is active first
			if !reservation.IsActive(uint64(batchReceivedAt.Unix())) {
				b.rollbackUpdates(successfulUpdates)
				return fmt.Errorf("account %s has inactive reservation for quorum %d", accountID.Hex(), update.quorumID)
			}

			// Then validate reservation period
			if !b.validateReservationPeriod(reservation, update.period, b.ChainPaymentState.GetReservationWindow(), batchReceivedAt) {
				b.rollbackUpdates(successfulUpdates)
				return fmt.Errorf("account %s has invalid reservation period for quorum %d", accountID.Hex(), update.quorumID)
			}

			// Track and validate usage
			b.logger.Info("Incrementing and validating usage", "accountID", accountID.Hex(), "quorumID", update.quorumID, "usage", update.usage, "period", update.period)
			prevUsage, err := b.incrementAndValidateUsage(accountID, update.quorumID, update.usage, reservation, update.period)
			if err != nil {
				b.logger.Error("Failed to increment and validate usage", "error", err)
				b.rollbackUpdates(successfulUpdates)
				return fmt.Errorf("account %s failed usage validation for quorum %d: %w", accountID.Hex(), update.quorumID, err)
			}
			b.logger.Info("Incremented and validated usage", "accountID", accountID.Hex(), "quorumID", update.quorumID, "usage", update.usage, "period", update.period, "prevUsage", prevUsage)
			// Record successful update for potential rollback
			record, err := newUpdateRecord(accountID, update.quorumID, update.period, prevUsage)
			if err != nil {
				b.rollbackUpdates(successfulUpdates)
				return fmt.Errorf("failed to create update record: %w", err)
			}
			b.logger.Info("Created update record", "accountID", accountID.Hex(), "quorumID", update.quorumID, "period", update.period, "prevUsage", prevUsage)
			successfulUpdates = append(successfulUpdates, record)
		}
	}

	return nil
}

// rollbackUpdates reverts all successful updates in reverse order
func (b *BatchMeterer) rollbackUpdates(updates []UpdateRecord) {
	b.logger.Info("Rolling back updates", "numUpdates", len(updates))
	if len(updates) == 0 {
		return
	}

	b.logger.Info("Starting rollback", "numUpdates", len(updates))

	// Process updates in reverse order to maintain consistency
	for i := len(updates) - 1; i >= 0; i-- {
		update := updates[i]
		accountUsage := b.getOrCreateAccountUsage(update.accountID)

		// Lock only for the duration of this update
		accountUsage.Lock.Lock()
		periodRecord := b.getOrCreatePeriodRecord(accountUsage, update.quorumID, update.period)

		// Log the state before rollback
		b.logger.Info("Rolling back update",
			"accountID", update.accountID.Hex(),
			"quorumID", update.quorumID,
			"period", update.period,
			"currentUsage", periodRecord.Usage,
			"restoredUsage", update.usage,
		)

		periodRecord.Usage = update.usage
		accountUsage.Lock.Unlock()
	}

	b.logger.Info("Completed rollback", "numUpdates", len(updates))
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
) (uint64, error) {
	accountUsage := b.getOrCreateAccountUsage(accountID)
	accountUsage.Lock.Lock()
	defer accountUsage.Lock.Unlock()

	periodRecord := b.getOrCreatePeriodRecord(accountUsage, quorumID, currentPeriod)
	binLimit := reservation.SymbolsPerSecond * b.ChainPaymentState.GetReservationWindow()
	prevUsage := periodRecord.Usage

	b.logger.Debug("incrementAndValidateUsage",
		"accountID", accountID.Hex(),
		"quorumID", quorumID,
		"prevUsage", prevUsage,
		"newUsage", usage,
		"binLimit", binLimit,
		"currentPeriod", currentPeriod,
	)

	periodRecord.Usage += usage
	newUsage := periodRecord.Usage

	// If new usage is within bin limit, done
	if newUsage <= binLimit {
		return prevUsage, nil
	}

	// If bin was already full before this increment, revert and error
	if prevUsage >= binLimit {
		b.logger.Debug("bin already full, reverting", "accountID", accountID.Hex(), "quorumID", quorumID, "period", currentPeriod, "prevUsage", prevUsage, "binLimit", binLimit)
		periodRecord.Usage = prevUsage
		return 0, fmt.Errorf("bin has already been filled for quorum %d", quorumID)
	}

	// If new usage is more than 2x bin limit, revert and error
	if newUsage > 2*binLimit {
		b.logger.Debug("usage exceeds 2x bin limit, reverting", "accountID", accountID.Hex(), "quorumID", quorumID, "period", currentPeriod, "prevUsage", prevUsage, "newUsage", newUsage, "2xBinLimit", 2*binLimit)
		periodRecord.Usage = prevUsage
		return 0, fmt.Errorf("overflow usage exceeds bin limit for quorum %d", quorumID)
	}

	// Check if overflow period is within reservation window
	nextPeriod := currentPeriod + b.ChainPaymentState.GetReservationWindow()
	endPeriod := meterer.GetReservationPeriod(int64(reservation.EndTimestamp), b.ChainPaymentState.GetReservationWindow())

	if nextPeriod >= endPeriod {
		b.logger.Debug("overflow period outside window, reverting", "accountID", accountID.Hex(), "quorumID", quorumID, "period", currentPeriod, "prevUsage", prevUsage, "nextPeriod", nextPeriod, "endPeriod", endPeriod)
		periodRecord.Usage = prevUsage
		return 0, fmt.Errorf("overflow period exceeds reservation window for quorum %d", quorumID)
	}

	// Move overflow to next period
	overflow := newUsage - binLimit
	periodRecord.Usage = binLimit
	overflowRecord := b.getOrCreatePeriodRecord(accountUsage, quorumID, nextPeriod)
	overflowPrev := overflowRecord.Usage
	overflowRecord.Usage += overflow
	// If overflow in next period exceeds bin limit, revert both and error
	if overflowRecord.Usage > binLimit {
		b.logger.Debug("overflow in next period exceeds limit, reverting both", "accountID", accountID.Hex(), "quorumID", quorumID, "period", currentPeriod, "prevUsage", prevUsage, "overflow", overflow, "overflowPrev", overflowPrev, "overflowNew", overflowRecord.Usage, "overflowBinLimit", binLimit)
		periodRecord.Usage = prevUsage
		overflowRecord.Usage = overflowPrev
		return 0, fmt.Errorf("overflow usage exceeds bin limit for quorum %d in next period", quorumID)
	}

	return prevUsage, nil
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

	// Use LoadOrStore to ensure thread-safe creation
	actual, loaded := b.AccountUsages.LoadOrStore(accountID, accountUsage)
	if loaded {
		// Another goroutine created the account usage first
		return actual.(*AccountUsage)
	}

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

	b.logger.Info("validateReservationPeriod debug",
		"requestReservationPeriod", requestReservationPeriod,
		"currentReservationPeriod", currentReservationPeriod,
		"startPeriod", startPeriod,
		"endPeriod", endPeriod,
		"isCurrentOrPreviousPeriod", isCurrentOrPreviousPeriod,
		"isWithinReservationWindow", isWithinReservationWindow,
	)

	return isCurrentOrPreviousPeriod && isWithinReservationWindow
}
