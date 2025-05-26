package node

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// Minimum number of bins to track for each account's reservation usage
const MinNumBins = 3

// UsageRecord represents the usage for a specific reservation period
type UsageRecord struct {
	Index uint64 // The reservation period index
	Usage uint64 // The usage amount for this period
}

// AccountUsage tracks usage for a specific account and quorum
type AccountUsage struct {
	// Circular buffer of usage records, with minimum length MinNumBins
	UsageRecords map[core.QuorumID][]UsageRecord
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
		NumBins:           max(numBins, MinNumBins),
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
	// Convert the batch to a slice of BatchRequestInfo
	requests, err := b.BatchToRequestInfos(batch)
	if err != nil {
		return err
	}

	// Process the batch requests
	return b.BatchMeterRequest(ctx, requests, batchTimestamp)
}

// BatchToRequestInfos converts a corev2.Batch to a slice of BatchRequestInfo
func (b *BatchMeterer) BatchToRequestInfos(batch *corev2.Batch) ([]BatchRequestInfo, error) {
	if batch == nil {
		return nil, fmt.Errorf("batch is nil")
	}

	if len(batch.BlobCertificates) == 0 {
		return nil, fmt.Errorf("batch has no blob certificates")
	}

	requests := make([]BatchRequestInfo, 0, len(batch.BlobCertificates))

	for _, cert := range batch.BlobCertificates {
		if cert.BlobHeader == nil {
			return nil, fmt.Errorf("blob certificate has nil header")
		}

		// Extract the account ID from the payment metadata
		accountID := cert.BlobHeader.PaymentMetadata.AccountID

		// Extract the quorum numbers from the blob header
		quorumIDs := make([]core.QuorumID, len(cert.BlobHeader.QuorumNumbers))
		copy(quorumIDs, cert.BlobHeader.QuorumNumbers)

		// Get the number of symbols from the blob commitment
		numSymbols := uint64(cert.BlobHeader.BlobCommitments.Length)

		// Add the request to the slice
		requests = append(requests, BatchRequestInfo{
			AccountID:  accountID,
			QuorumIDs:  quorumIDs,
			NumSymbols: numSymbols,
		})
	}

	return requests, nil
}

// BatchMeterRequest tracks and validates a batch of requests
// Each request includes an account ID, quorum numbers, and symbols used
// Returns an error if any account or quorum in the batch is invalid
func (b *BatchMeterer) BatchMeterRequest(
	ctx context.Context,
	requests []BatchRequestInfo,
	batchTimestamp time.Time,
) error {
	// Aggregate usage by account and quorum
	aggregatedUsage := b.AggregateRequests(requests)

	// Validate all accounts and quorums against their reservations
	for accountID, quorumUsages := range aggregatedUsage {
		// Get reservations for this account
		reservations, err := b.ChainPaymentState.GetReservedPaymentByAccountAndQuorums(
			ctx, accountID, mapQuorumIDsToSortedSlice(quorumUsages),
		)
		if err != nil {
			b.logger.Error("Failed to get reservations", "account", accountID, "error", err)
			return fmt.Errorf("failed to get reservations for account %s: %w", accountID.Hex(), err)
		}

		// Get current reservation period
		currentPeriod := meterer.GetReservationPeriod(
			batchTimestamp.Unix(),
			b.ChainPaymentState.GetReservationWindow(),
		)

		// Validate each quorum's usage
		for quorumID, usage := range quorumUsages {
			reservation, ok := reservations[quorumID]
			if !ok {
				return fmt.Errorf("account %s has no reservation for quorum %d", accountID.Hex(), quorumID)
			}

			// Check if reservation is active
			if !reservation.IsActive(uint64(batchTimestamp.Unix())) {
				return fmt.Errorf("account %s has inactive reservation for quorum %d", accountID.Hex(), quorumID)
			}

			// Validate reservation period (TODO: quorum specific)
			if !b.validateReservationPeriod(reservation, currentPeriod, b.ChainPaymentState.GetReservationWindow(), batchTimestamp) {
				return fmt.Errorf("account %s has invalid reservation period for quorum %d", accountID.Hex(), quorumID)
			}

			// Increment and validate usage
			err := b.incrementAndValidateUsage(
				accountID,
				quorumID,
				usage,
				reservation,
				currentPeriod,
			)
			if err != nil {
				return fmt.Errorf("account %s failed usage validation for quorum %d: %w", accountID.Hex(), quorumID, err)
			}
		}
	}

	return nil
}

// AggregateRequests takes a slice of BatchRequestInfo and returns a map of account IDs to quorum usage
// This is useful for pre-computing the aggregated usage before metering
func (b *BatchMeterer) AggregateRequests(requests []BatchRequestInfo) map[gethcommon.Address]map[core.QuorumID]uint64 {
	aggregatedUsage := make(map[gethcommon.Address]map[core.QuorumID]uint64)

	for _, req := range requests {
		if _, exists := aggregatedUsage[req.AccountID]; !exists {
			aggregatedUsage[req.AccountID] = make(map[core.QuorumID]uint64)
		}

		symbolsCharged := b.symbolsCharged(req.NumSymbols)
		for _, quorumID := range req.QuorumIDs {
			aggregatedUsage[req.AccountID][quorumID] += symbolsCharged
		}
	}

	return aggregatedUsage
}

// incrementAndValidateUsage increments the usage for an account and quorum and validates against limits
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

	usageRecord := b.getOrCreateUsageRecord(accountUsage, quorumID, currentPeriod)
	prevUsage := usageRecord.Usage
	newUsage := prevUsage + usage

	// Calculate the bin limit based on the reservation's symbols per second
	binLimit := reservation.SymbolsPerSecond * b.ChainPaymentState.GetReservationWindow()

	if newUsage <= binLimit {
		// Usage is within the bin limit, update the record
		usageRecord.Usage = newUsage
		return nil
	}

	// Usage exceeds the bin limit, check if overflow is possible
	if prevUsage >= binLimit {
		// Bin was already filled before this increment
		return fmt.Errorf("bin already filled for period %d", currentPeriod)
	}

	// Check if we can overflow to a future bin
	overflowPeriod := currentPeriod + 2*b.ChainPaymentState.GetReservationWindow()
	if overflowPeriod >= reservation.EndTimestamp {
		// Overflow period would be after reservation end, can't overflow
		return fmt.Errorf("bin usage exceeds limit and cannot overflow")
	}

	// Calculate overflow amount
	overflowAmount := newUsage - binLimit

	// Check if overflow amount exceeds bin limit
	if overflowAmount > binLimit {
		// Overflow amount exceeds a single bin
		return fmt.Errorf("overflow amount %d exceeds bin limit %d", overflowAmount, binLimit)
	}

	// Update the current bin to be at capacity
	usageRecord.Usage = binLimit

	// Increment the overflow bin
	overflowRecord := b.getOrCreateUsageRecord(accountUsage, quorumID, overflowPeriod)
	if overflowRecord.Usage+overflowAmount > binLimit {
		// Overflow bin would exceed capacity
		return fmt.Errorf("overflow bin would exceed capacity")
	}
	overflowRecord.Usage += overflowAmount

	return nil
}

// getOrCreateAccountUsage gets or creates a usage record for an account
func (b *BatchMeterer) getOrCreateAccountUsage(accountID gethcommon.Address) *AccountUsage {
	value, ok := b.AccountUsages.Load(accountID)
	if ok {
		return value.(*AccountUsage)
	}

	accountUsage := &AccountUsage{
		UsageRecords: make(map[core.QuorumID][]UsageRecord),
	}
	b.AccountUsages.Store(accountID, accountUsage)
	return accountUsage
}

// getOrCreateUsageRecord gets or creates a usage record for a specific account, quorum, and period
func (b *BatchMeterer) getOrCreateUsageRecord(
	accountUsage *AccountUsage,
	quorumID core.QuorumID,
	period uint64,
) *UsageRecord {
	if _, ok := accountUsage.UsageRecords[quorumID]; !ok {
		accountUsage.UsageRecords[quorumID] = make([]UsageRecord, b.NumBins)
		for i := range accountUsage.UsageRecords[quorumID] {
			accountUsage.UsageRecords[quorumID][i] = UsageRecord{
				Index: uint64(i),
				Usage: 0,
			}
		}
	}

	// Find the existing record for this period or the least recently used slot
	records := accountUsage.UsageRecords[quorumID]
	oldestIndex := uint32(0)
	oldestTime := uint64(0xFFFFFFFFFFFFFFFF)

	for i, record := range records {
		if record.Index == period {
			// Found the record for this period
			return &records[i]
		}
		if record.Index < oldestTime {
			oldestTime = record.Index
			oldestIndex = uint32(i)
		}
	}

	// Reuse the oldest slot for the new period
	records[oldestIndex] = UsageRecord{
		Index: period,
		Usage: 0,
	}
	return &records[oldestIndex]
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

// BatchRequestInfo represents a single request within a batch
type BatchRequestInfo struct {
	AccountID  gethcommon.Address
	QuorumIDs  []core.QuorumID
	NumSymbols uint64
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
