package meterer

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer/payment_logic"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

type UsageRequest struct {
	AccountID gethcommon.Address
	QuorumID  core.QuorumID
	Period    uint64
	Usage     uint64
}

type AccountUsage struct {
	PeriodRecords map[core.QuorumID][]*pb.PeriodRecord
	Lock          sync.RWMutex
}

type BatchLedger struct {
	state    OnchainPayment
	accounts sync.Map
	logger   logging.Logger
}

func NewBatchLedger(state OnchainPayment, _ Config, logger logging.Logger) *BatchLedger {
	return &BatchLedger{
		state:    state,
		accounts: sync.Map{},
		logger:   logger.With("component", "BatchLedger"),
	}
}

func (bl *BatchLedger) MeterBatch(ctx context.Context, batch *corev2.Batch, receivedAt time.Time) error {
	if batch == nil || len(batch.BlobCertificates) == 0 {
		return fmt.Errorf("batch is nil or empty")
	}

	params, err := bl.state.GetPaymentGlobalParams()
	if err != nil {
		return fmt.Errorf("failed to get payment global params: %w", err)
	}

	requests, err := bl.parseRequests(params, batch)
	if err != nil {
		return fmt.Errorf("failed to parse batch: %w", err)
	}

	return bl.processBatchRequests(ctx, params, requests, receivedAt)
}

func (bl *BatchLedger) parseRequests(params *PaymentVaultParams, batch *corev2.Batch) ([]*UsageRequest, error) {
	var requests []*UsageRequest
	requestMap := make(map[string]*UsageRequest)

	for _, cert := range batch.BlobCertificates {
		if cert.BlobHeader == nil {
			return nil, fmt.Errorf("blob certificate has nil header")
		}

		accountID := cert.BlobHeader.PaymentMetadata.AccountID
		numSymbols := uint64(cert.BlobHeader.BlobCommitments.Length)

		for _, quorumID := range cert.BlobHeader.QuorumNumbers {
			quorumConfig, exists := params.QuorumProtocolConfigs[quorumID]
			if !exists {
				return nil, fmt.Errorf("no protocol config found for quorum %d", quorumID)
			}
			if quorumConfig.ReservationRateLimitWindow == 0 {
				return nil, fmt.Errorf("invalid zero ReservationRateLimitWindow for quorum %d", quorumID)
			}

			period := payment_logic.GetReservationPeriodByNanosecond(cert.BlobHeader.PaymentMetadata.Timestamp, quorumConfig.ReservationRateLimitWindow)
			usage := payment_logic.SymbolsCharged(numSymbols, quorumConfig.MinNumSymbols)

			key := fmt.Sprintf("%s_%d_%d", accountID.Hex(), quorumID, period)
			if existing, exists := requestMap[key]; exists {
				existing.Usage += usage
			} else {
				request := &UsageRequest{
					AccountID: accountID,
					QuorumID:  quorumID,
					Period:    period,
					Usage:     usage,
				}
				requestMap[key] = request
				requests = append(requests, request)
			}
		}
	}

	return requests, nil
}

func (bl *BatchLedger) processBatchRequests(ctx context.Context, params *PaymentVaultParams, requests []*UsageRequest, receivedAt time.Time) error {
	reservations, err := bl.fetchAccountReservations(ctx, requests)
	if err != nil {
		return fmt.Errorf("failed to fetch reservations: %w", err)
	}

	// Validate all requests first without modifying state
	for _, req := range requests {
		if err := bl.validateUsageRequest(params, req, reservations[req.AccountID], receivedAt); err != nil {
			return fmt.Errorf("validation failed for account %s quorum %d: %w", req.AccountID.Hex(), req.QuorumID, err)
		}
	}

	// Apply all changes after validation passes
	for _, req := range requests {
		if err := bl.commitUsageRequest(params, req, reservations[req.AccountID]); err != nil {
			return fmt.Errorf("commit failed for account %s quorum %d: %w", req.AccountID.Hex(), req.QuorumID, err)
		}
	}

	return nil
}

func (bl *BatchLedger) fetchAccountReservations(ctx context.Context, requests []*UsageRequest) (map[gethcommon.Address]map[core.QuorumID]*core.ReservedPayment, error) {
	accountQuorums := make(map[gethcommon.Address]map[core.QuorumID]struct{})
	for _, req := range requests {
		if accountQuorums[req.AccountID] == nil {
			accountQuorums[req.AccountID] = make(map[core.QuorumID]struct{})
		}
		accountQuorums[req.AccountID][req.QuorumID] = struct{}{}
	}

	reservations := make(map[gethcommon.Address]map[core.QuorumID]*core.ReservedPayment)
	for accountID, quorumSet := range accountQuorums {
		quorumIDs := make([]core.QuorumID, 0, len(quorumSet))
		for quorumID := range quorumSet {
			quorumIDs = append(quorumIDs, quorumID)
		}
		slices.Sort(quorumIDs)

		accountReservations, err := bl.state.GetReservedPaymentByAccountAndQuorums(ctx, accountID, quorumIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to get reservations for account %s: %w", accountID.Hex(), err)
		}
		reservations[accountID] = accountReservations
	}

	return reservations, nil
}

func (bl *BatchLedger) validateUsageRequest(params *PaymentVaultParams, req *UsageRequest, reservations map[core.QuorumID]*core.ReservedPayment, receivedAt time.Time) error {
	reservation, ok := reservations[req.QuorumID]
	if !ok {
		return fmt.Errorf("no reservation for quorum %d", req.QuorumID)
	}
	if !reservation.IsActive(uint64(receivedAt.Unix())) {
		return fmt.Errorf("inactive reservation for quorum %d", req.QuorumID)
	}

	quorumConfig := params.QuorumProtocolConfigs[req.QuorumID]
	if !payment_logic.ValidateReservationPeriod(reservation, req.Period, quorumConfig.ReservationRateLimitWindow, receivedAt.UnixNano()) {
		return fmt.Errorf("invalid reservation period")
	}

	accountUsage := bl.getAccount(req.AccountID)
	accountUsage.Lock.RLock()
	defer accountUsage.Lock.RUnlock()

	return bl.checkUsageLimits(params, accountUsage, req.QuorumID, req.Usage, reservation, req.Period)
}

func (bl *BatchLedger) commitUsageRequest(params *PaymentVaultParams, req *UsageRequest, reservations map[core.QuorumID]*core.ReservedPayment) error {
	reservation := reservations[req.QuorumID]
	accountUsage := bl.getAccount(req.AccountID)
	accountUsage.Lock.Lock()
	defer accountUsage.Lock.Unlock()

	_, err := bl.applyUsage(params, accountUsage, req.QuorumID, req.Usage, reservation, req.Period)
	return err
}

func (bl *BatchLedger) checkUsageLimits(params *PaymentVaultParams, accountUsage *AccountUsage, quorumID core.QuorumID, usage uint64, reservation *core.ReservedPayment, period uint64) error {
	quorumConfig := params.QuorumProtocolConfigs[quorumID]
	binLimit := reservation.SymbolsPerSecond * quorumConfig.ReservationRateLimitWindow

	if usage > binLimit*2 {
		return fmt.Errorf("overflow usage exceeds bin limit for quorum %d", quorumID)
	}

	periodRecord := bl.getPeriodRecord(params, accountUsage, quorumID, period)
	originalUsage := periodRecord.Usage
	newUsage := originalUsage + usage

	if newUsage <= binLimit {
		return nil
	}

	if originalUsage >= binLimit {
		return fmt.Errorf("bin already full for quorum %d", quorumID)
	}

	nextPeriod := period + quorumConfig.ReservationRateLimitWindow
	endPeriod := payment_logic.GetReservationPeriod(int64(reservation.EndTimestamp), quorumConfig.ReservationRateLimitWindow)
	if nextPeriod >= endPeriod {
		return fmt.Errorf("overflow would exceed reservation window")
	}

	overflow := newUsage - binLimit
	overflowRecord := bl.getPeriodRecord(params, accountUsage, quorumID, nextPeriod)
	if overflowRecord.Usage+overflow > binLimit {
		return fmt.Errorf("overflow would exceed next period limit")
	}

	return nil
}

func (bl *BatchLedger) applyUsage(params *PaymentVaultParams, accountUsage *AccountUsage, quorumID core.QuorumID, usage uint64, reservation *core.ReservedPayment, period uint64) (uint64, error) {
	quorumConfig := params.QuorumProtocolConfigs[quorumID]
	binLimit := reservation.SymbolsPerSecond * quorumConfig.ReservationRateLimitWindow

	if usage > binLimit*2 {
		return 0, fmt.Errorf("overflow usage exceeds bin limit for quorum %d", quorumID)
	}

	periodRecord := bl.getPeriodRecord(params, accountUsage, quorumID, period)
	originalUsage := periodRecord.Usage
	newUsage := originalUsage + usage

	if newUsage <= binLimit {
		periodRecord.Usage = newUsage
		return originalUsage, nil
	}

	if originalUsage >= binLimit {
		return 0, fmt.Errorf("bin already full for quorum %d", quorumID)
	}

	nextPeriod := period + quorumConfig.ReservationRateLimitWindow
	endPeriod := payment_logic.GetReservationPeriod(int64(reservation.EndTimestamp), quorumConfig.ReservationRateLimitWindow)
	if nextPeriod >= endPeriod {
		return 0, fmt.Errorf("overflow would exceed reservation window")
	}

	overflow := newUsage - binLimit
	overflowRecord := bl.getPeriodRecord(params, accountUsage, quorumID, nextPeriod)
	if overflowRecord.Usage+overflow > binLimit {
		return 0, fmt.Errorf("overflow would exceed next period limit")
	}

	periodRecord.Usage = binLimit
	overflowRecord.Usage += overflow
	return originalUsage, nil
}

func (bl *BatchLedger) getAccount(accountID gethcommon.Address) *AccountUsage {
	value, ok := bl.accounts.Load(accountID)
	if ok {
		return value.(*AccountUsage)
	}

	accountUsage := &AccountUsage{
		PeriodRecords: make(map[core.QuorumID][]*pb.PeriodRecord),
	}

	actual, loaded := bl.accounts.LoadOrStore(accountID, accountUsage)
	if loaded {
		return actual.(*AccountUsage)
	}

	return accountUsage
}

func (bl *BatchLedger) getPeriodRecord(params *PaymentVaultParams, accountUsage *AccountUsage, quorumID core.QuorumID, period uint64) *pb.PeriodRecord {
	if accountUsage.PeriodRecords[quorumID] == nil {
		accountUsage.PeriodRecords[quorumID] = make([]*pb.PeriodRecord, MinNumBins)
	}

	quorumConfig := params.QuorumProtocolConfigs[quorumID]
	relativeIndex := uint32((period / quorumConfig.ReservationRateLimitWindow) % uint64(MinNumBins))

	if accountUsage.PeriodRecords[quorumID][relativeIndex] == nil {
		accountUsage.PeriodRecords[quorumID][relativeIndex] = &pb.PeriodRecord{
			Index: uint32(period),
			Usage: 0,
		}
	}

	record := accountUsage.PeriodRecords[quorumID][relativeIndex]
	if record.Index < uint32(period) {
		record.Index = uint32(period)
		record.Usage = 0
	}

	return record
}
