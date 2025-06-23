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

// BatchMeterError represents a standardized error from the batch meterer
type BatchMeterError struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	AccountID string `json:"account_id,omitempty"`
	QuorumID  uint8  `json:"quorum_id,omitempty"`
}

// Error implements the error interface
func (e *BatchMeterError) Error() string {
	if e.AccountID != "" && e.QuorumID != 0 {
		return fmt.Sprintf("[%s] %s (account: %s, quorum: %d)", e.Code, e.Message, e.AccountID, e.QuorumID)
	}
	if e.AccountID != "" {
		return fmt.Sprintf("[%s] %s (account: %s)", e.Code, e.Message, e.AccountID)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Batch meter error codes
const (
	ErrCodeBatchEmpty          = "BATCH_EMPTY"
	ErrCodeBlobHeaderNil       = "BLOB_HEADER_NIL"
	ErrCodePaymentParamsFailed = "PAYMENT_PARAMS_FAILED"
	ErrCodeReservationNotFound = "RESERVATION_NOT_FOUND"
	ErrCodeReservationLookup   = "RESERVATION_LOOKUP_FAILED"
	ErrCodeReservationInactive = "RESERVATION_INACTIVE"
	ErrCodeReservationPeriod   = "RESERVATION_PERIOD_INVALID"
	ErrCodeBinAlreadyFull      = "BIN_ALREADY_FULL"
	ErrCodeUsageExceedsLimit   = "USAGE_EXCEEDS_LIMIT"
	ErrCodeOverflowPeriodLimit = "OVERFLOW_PERIOD_LIMIT"
	ErrCodeOverflowWindowLimit = "OVERFLOW_WINDOW_LIMIT"
)

// Error constructors for standardized error creation
func newBatchMeterError(code, message string) *BatchMeterError {
	return &BatchMeterError{Code: code, Message: message}
}

func newAccountError(code, message string, accountID gethcommon.Address) *BatchMeterError {
	return &BatchMeterError{Code: code, Message: message, AccountID: accountID.Hex()}
}

func newAccountQuorumError(code, message string, accountID gethcommon.Address, quorumID core.QuorumID) *BatchMeterError {
	return &BatchMeterError{Code: code, Message: message, AccountID: accountID.Hex(), QuorumID: uint8(quorumID)}
}

// IsBatchMeterError checks if an error is a BatchMeterError and returns it
func IsBatchMeterError(err error) (*BatchMeterError, bool) {
	if bmErr, ok := err.(*BatchMeterError); ok {
		return bmErr, true
	}
	return nil, false
}

// AccountUsage tracks usage for a specific account across quorums
type AccountUsage struct {
	// Circular buffer of usage records, with minimum length MinNumBins
	PeriodRecords map[core.QuorumID][]*pb.PeriodRecord
	// Each account has a lock to protect concurrent access to period records
	Lock sync.RWMutex
}

// UpdateRecord tracks a successful usage update for potential rollback
type UpdateRecord struct {
	accountID gethcommon.Address
	quorumID  core.QuorumID
	// period is the reservation rate limit period
	period uint64
	usage  uint64
}

// newUpdateRecord creates a new UpdateRecord with validation
func newUpdateRecord(accountID gethcommon.Address, quorumID core.QuorumID, period, usage uint64) *UpdateRecord {
	return &UpdateRecord{
		accountID: accountID,
		quorumID:  quorumID,
		period:    period,
		usage:     usage,
	}
}

// BatchMeterer handles metering for batches of requests that may contain multiple accounts, quorums, and periods
type BatchMeterer struct {
	// Configuration for the batch meterer
	Config meterer.Config

	ChainPaymentState meterer.OnchainPayment

	// AccountUsages tracks in-memory usage for a map of accounts to AccountUsage
	AccountUsages sync.Map

	logger logging.Logger
}

// NewBatchMeterer creates a new batch meterer
func NewBatchMeterer(
	config meterer.Config,
	paymentChainState meterer.OnchainPayment,
	logger logging.Logger,
) *BatchMeterer {
	return &BatchMeterer{
		Config:            config,
		ChainPaymentState: paymentChainState,
		AccountUsages:     sync.Map{},
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
	params, err := b.ChainPaymentState.GetPaymentGlobalParams()
	if err != nil {
		return newBatchMeterError(ErrCodePaymentParamsFailed, fmt.Sprintf("failed to get payment global params: %v", err))
	}
	// convert batch into usage records tracking account, quorum, reservation period, and usage
	usages, err := b.batchRequestUsage(params, batch)
	if err != nil {
		b.logger.Error("batchRequestUsage", "error", err)
		return err
	}

	// validate the usages against the accounts' reservations
	return b.processBatch(ctx, params, usages, batchReceivedAt)
}

// batchRequestUsage aggregates the usage for each account and quorum in the batch
// Returns an array of UpdateRecord representing the usage updates
func (b *BatchMeterer) batchRequestUsage(params *meterer.PaymentVaultParams, batch *corev2.Batch) ([]*UpdateRecord, error) {
	b.logger.Info("batchRequestUsage", "batch", batch)
	if batch == nil || len(batch.BlobCertificates) == 0 {
		return nil, newBatchMeterError(ErrCodeBatchEmpty, "batch is nil or empty")
	}

	// Use a map to track existing updates by account, quorum, and period
	updatesMap := make(map[string]*UpdateRecord)

	for _, cert := range batch.BlobCertificates {
		if cert.BlobHeader == nil {
			return nil, newBatchMeterError(ErrCodeBlobHeaderNil, "blob certificate has nil header")
		}

		accountID := cert.BlobHeader.PaymentMetadata.AccountID
		numSymbols := uint64(cert.BlobHeader.BlobCommitments.Length)

		// Add usage for each quorum in the blob
		for _, quorumID := range cert.BlobHeader.QuorumNumbers {
			// Use current period as the index
			currentPeriod := meterer.GetReservationPeriodByNanosecond(cert.BlobHeader.PaymentMetadata.Timestamp, params.QuorumProtocolConfigs[quorumID].ReservationRateLimitWindow)
			symbolsCharged := meterer.SymbolsCharged(numSymbols, params.QuorumProtocolConfigs[quorumID].MinNumSymbols)

			// Create a unique key for this account/quorum/period combination
			key := fmt.Sprintf("%s_%d_%d", accountID.Hex(), quorumID, currentPeriod)

			if existingRecord, exists := updatesMap[key]; exists {
				// If record exists, add to its usage
				existingRecord.usage += symbolsCharged
			} else {
				// Create new record if none exists
				record := newUpdateRecord(accountID, quorumID, currentPeriod, symbolsCharged)
				updatesMap[key] = record
			}
		}
	}

	// Convert map to slice
	updates := make([]*UpdateRecord, 0, len(updatesMap))
	for _, record := range updatesMap {
		updates = append(updates, record)
	}

	return updates, nil
}

// processBatch validates and tracks usage for a batch of requests
// For each account and quorum:
// 1. Validates against reservation limits
// 2. Checks reservation period validity
// 3. Tracks usage in period records
func (b *BatchMeterer) processBatch(
	ctx context.Context,
	params *meterer.PaymentVaultParams,
	updates []*UpdateRecord,
	batchReceivedAt time.Time,
) error {
	b.logger.Info("processBatch", "numUpdates", len(updates), "batchReceivedAt", batchReceivedAt)

	successfulUpdates := make([]*UpdateRecord, 0)

	// Group updates by account for efficient reservation lookup
	accountUpdates := make(map[gethcommon.Address][]*UpdateRecord)
	for _, update := range updates {
		accountUpdates[update.accountID] = append(accountUpdates[update.accountID], update)
	}

	for accountID, accountUpdates := range accountUpdates {
		b.logger.Info("processing account updates", "accountID", accountID.Hex(), "numUpdates", len(accountUpdates))

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

		// Get active reservations for this account's quorums
		reservations, err := b.ChainPaymentState.GetReservedPaymentByAccountAndQuorums(ctx, accountID, quorumIDs)
		if err != nil {
			b.rollbackUpdates(params, successfulUpdates)
			return newAccountError(ErrCodeReservationLookup, fmt.Sprintf("failed to get reservations: %v", err), accountID)
		}

		// Process each update for this account
		for _, update := range accountUpdates {
			reservation, ok := reservations[update.quorumID]
			if !ok {
				b.rollbackUpdates(params, successfulUpdates)
				return newAccountQuorumError(ErrCodeReservationNotFound, "no reservation", accountID, update.quorumID)
			}

			// TODO: Validating reservations can be refactored after the other refactoring PR
			// Check if reservation is active first
			if !reservation.IsActive(uint64(batchReceivedAt.Unix())) {
				b.rollbackUpdates(params, successfulUpdates)
				return newAccountQuorumError(ErrCodeReservationInactive, "reservation is inactive", accountID, update.quorumID)
			}
			if !meterer.ValidateReservationPeriod(reservation, update.period, params.QuorumProtocolConfigs[update.quorumID].ReservationRateLimitWindow, batchReceivedAt) {
				b.rollbackUpdates(params, successfulUpdates)
				return newAccountQuorumError(ErrCodeReservationPeriod, "reservation period is invalid", accountID, update.quorumID)
			}

			b.logger.Info("incrementing and validating usage", "accountID", accountID.Hex(), "quorumID", update.quorumID, "usage", update.usage, "period", update.period)
			prevUsage, err := b.incrementAndValidateUsage(params, accountID, update.quorumID, update.usage, reservation, update.period)
			if err != nil {
				b.logger.Error("Failed to increment and validate usage", "error", err)
				b.rollbackUpdates(params, successfulUpdates)
				return err
			}
			// Record successful update for potential rollback
			successfulUpdates = append(successfulUpdates, newUpdateRecord(accountID, update.quorumID, update.period, prevUsage))
		}
	}

	return nil
}

// rollbackUpdates reverts all successful updates in reverse order
func (b *BatchMeterer) rollbackUpdates(params *meterer.PaymentVaultParams, updates []*UpdateRecord) {
	b.logger.Info("Rolling back updates", "numUpdates", len(updates))
	if len(updates) == 0 {
		return
	}

	// Process updates in reverse order to maintain consistency
	for i := len(updates) - 1; i >= 0; i-- {
		update := updates[i]
		accountUsage := b.getOrCreateAccountUsage(update.accountID)

		// Lock only for the duration of this update
		accountUsage.Lock.Lock()
		periodRecord := b.getOrRefreshPeriodRecord(params, accountUsage, update.quorumID, update.period)
		periodRecord.Usage = update.usage
		accountUsage.Lock.Unlock()
	}
}

// incrementAndValidateUsage validates and increments usage for an account and quorum
// Uses a write-first approach for atomicity:
// 1. Writes changes first
// 2. Validates the new state
// 3. Rolls back if validation fails
func (b *BatchMeterer) incrementAndValidateUsage(
	params *meterer.PaymentVaultParams,
	accountID gethcommon.Address,
	quorumID core.QuorumID,
	usage uint64,
	reservation *core.ReservedPayment,
	currentPeriod uint64,
) (uint64, error) {
	accountUsage := b.getOrCreateAccountUsage(accountID)
	accountUsage.Lock.Lock()
	defer accountUsage.Lock.Unlock()

	periodRecord := b.getOrRefreshPeriodRecord(params, accountUsage, quorumID, currentPeriod)
	binLimit := reservation.SymbolsPerSecond * params.QuorumProtocolConfigs[quorumID].ReservationRateLimitWindow
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
		return 0, newAccountQuorumError(ErrCodeBinAlreadyFull, "rate limit bin is already full", accountID, quorumID)
	}

	// If new usage is more than 2x bin limit, revert and error
	if newUsage > 2*binLimit {
		b.logger.Debug("usage exceeds 2x bin limit, reverting", "accountID", accountID.Hex(), "quorumID", quorumID, "period", currentPeriod, "prevUsage", prevUsage, "newUsage", newUsage, "2xBinLimit", 2*binLimit)
		periodRecord.Usage = prevUsage
		return 0, newAccountQuorumError(ErrCodeUsageExceedsLimit, "usage exceeds bin limit", accountID, quorumID)
	}

	// Check if overflow period is within reservation window
	nextPeriod := currentPeriod + params.QuorumProtocolConfigs[quorumID].ReservationRateLimitWindow
	endPeriod := meterer.GetReservationPeriod(int64(reservation.EndTimestamp), params.QuorumProtocolConfigs[quorumID].ReservationRateLimitWindow)

	if nextPeriod >= endPeriod {
		b.logger.Debug("overflow period outside window, reverting", "accountID", accountID.Hex(), "quorumID", quorumID, "period", currentPeriod, "prevUsage", prevUsage, "nextPeriod", nextPeriod, "endPeriod", endPeriod)
		periodRecord.Usage = prevUsage
		return 0, newAccountQuorumError(ErrCodeOverflowWindowLimit, "overflow period exceeds reservation window", accountID, quorumID)
	}

	// Move overflow to next period
	overflow := newUsage - binLimit
	periodRecord.Usage = binLimit
	overflowRecord := b.getOrRefreshPeriodRecord(params, accountUsage, quorumID, nextPeriod)
	overflowPrev := overflowRecord.Usage
	overflowRecord.Usage += overflow
	// If overflow in next period exceeds bin limit, revert both and error
	if overflowRecord.Usage > binLimit {
		b.logger.Debug("overflow in next period exceeds limit, reverting both", "accountID", accountID.Hex(), "quorumID", quorumID, "period", currentPeriod, "prevUsage", prevUsage, "overflow", overflow, "overflowPrev", overflowPrev, "overflowNew", overflowRecord.Usage, "overflowBinLimit", binLimit)
		periodRecord.Usage = prevUsage
		overflowRecord.Usage = overflowPrev
		return 0, newAccountQuorumError(ErrCodeOverflowPeriodLimit, "overflow usage exceeds next period limit", accountID, quorumID)
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
		return actual.(*AccountUsage)
	}

	return accountUsage
}

// getOrRefreshPeriodRecord gets or creates a period record for an account and quorum
func (b *BatchMeterer) getOrRefreshPeriodRecord(
	params *meterer.PaymentVaultParams,
	accountUsage *AccountUsage,
	quorumID core.QuorumID,
	period uint64,
) *pb.PeriodRecord {
	if _, exists := accountUsage.PeriodRecords[quorumID]; !exists {
		accountUsage.PeriodRecords[quorumID] = make([]*pb.PeriodRecord, meterer.MinNumBins)
	}

	// Calculate relative index in circular buffer
	relativeIndex := uint32((period / params.QuorumProtocolConfigs[quorumID].ReservationRateLimitWindow) % uint64(meterer.MinNumBins))
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
