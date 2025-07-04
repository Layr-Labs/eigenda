package meterer

import (
	"context"
	"fmt"
	"math/big"
	"time"

	disperser_v2 "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer/payment_logic"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// ServerAccountLedger handles validation for a single specific account
// It replicates the meterer logic but for account-specific validation only
type ServerAccountLedger struct {
	// Account identity
	accountID gethcommon.Address

	// Account-specific state from on-chain
	reservations    map[core.QuorumID]*core.ReservedPayment
	onDemandPayment *core.OnDemandPayment

	// Account-specific state from metering store
	periodRecords     QuorumPeriodRecords
	cumulativePayment *big.Int
	
	// Payment rollback state - stores the previous payment amount for potential rollback
	lastOldPayment *big.Int

	// Dependencies
	meteringStore MeteringStore
	config        Config
	logger        logging.Logger
}

// NewServerAccountLedger creates a new account-specific ledger
func NewServerAccountLedger(
	ctx context.Context,
	accountID gethcommon.Address,
	chainPaymentState OnchainPayment,
	meteringStore MeteringStore,
	config Config,
	logger logging.Logger,
) (*ServerAccountLedger, error) {
	// Fetch fresh on-chain state for this account
	reservations, err := chainPaymentState.GetReservedPaymentByAccountAndQuorums(ctx, accountID, []core.QuorumID{})
	if err != nil {
		return nil, fmt.Errorf("failed to get reservations: %w", err)
	}

	onDemandPayment, err := chainPaymentState.GetOnDemandPaymentByAccount(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get on-demand payment: %w", err)
	}

	// Get cumulative payment from storage
	cumulativePayment, err := meteringStore.GetLargestCumulativePayment(ctx, accountID)
	if err != nil {
		logger.Warn("Failed to get cumulative payment from store, using zero", "accountID", accountID.Hex(), "error", err)
		cumulativePayment = big.NewInt(0)
	}

	sal := &ServerAccountLedger{
		accountID:         accountID,
		reservations:      reservations,
		onDemandPayment:   onDemandPayment,
		periodRecords:     make(QuorumPeriodRecords),
		cumulativePayment: cumulativePayment,
		meteringStore:     meteringStore,
		config:            config,
		logger:            logger,
	}

	// Load period records for quorums with reservations
	// Note: We'll load them lazily when needed since we need payment params

	return sal, nil
}

// loadPeriodRecords loads recent period records from the metering store
func (sal *ServerAccountLedger) loadPeriodRecords(ctx context.Context, params *PaymentVaultParams) error {
	if len(sal.reservations) == 0 {
		return nil // No reservations, no period records needed
	}

	// Get all quorum IDs that have reservations and calculate their current periods
	quorumIDs := make([]core.QuorumID, 0, len(sal.reservations))
	periods := make([]uint64, 0, len(sal.reservations))
	currentTime := time.Now().UnixNano()

	for quorumID := range sal.reservations {
		quorumIDs = append(quorumIDs, quorumID)
		
		// Get the proper reservation window for this quorum
		_, protocolConfig, err := params.GetQuorumConfigs(quorumID)
		if err != nil {
			sal.logger.Warn("Failed to get quorum config for period records", "quorumID", quorumID, "error", err)
			// Use a default if we can't get the config
			periods = append(periods, payment_logic.GetReservationPeriodByNanosecond(currentTime, 3600))
			continue
		}
		
		// Calculate current period using the proper reservation window
		currentPeriod := payment_logic.GetReservationPeriodByNanosecond(currentTime, protocolConfig.ReservationRateLimitWindow)
		periods = append(periods, currentPeriod)
	}

	// Fetch period records from storage
	periodRecordsProto, err := sal.meteringStore.GetPeriodRecords(ctx, sal.accountID, quorumIDs, periods, MinNumBins)
	if err != nil {
		return fmt.Errorf("failed to get period records: %w", err)
	}

	// Convert protobuf format to local format
	for quorumID, protoRecords := range periodRecordsProto {
		if protoRecords != nil && len(protoRecords.Records) > 0 {
			localRecords := make([]*PeriodRecord, len(protoRecords.Records))
			for i, record := range protoRecords.Records {
				localRecords[i] = &PeriodRecord{
					Index: record.Index,
					Usage: record.Usage,
				}
			}
			sal.periodRecords[quorumID] = localRecords
		}
	}

	return nil
}

// Debit validates and records usage for this specific account
// This replicates the meterer logic but account-specific
func (sal *ServerAccountLedger) Debit(
	ctx context.Context,
	header core.PaymentMetadata,
	numSymbols uint64,
	quorumNumbers []uint8,
	params *PaymentVaultParams,
	receivedAt time.Time,
) (*big.Int, error) {
	// Use the same payment method decision logic as meterer
	if !payment_logic.IsOnDemandPayment(&header) {
		// Reservation path - validate against account's reservations
		return sal.processReservationRequest(ctx, header, numSymbols, quorumNumbers, params, receivedAt)
	} else {
		// On-demand path - validate payment and record usage
		return sal.processOnDemandRequest(ctx, header, numSymbols, quorumNumbers, params, receivedAt)
	}
}

// processReservationRequest handles reservation-based requests
// This replicates meterer.serveReservationRequest logic
func (sal *ServerAccountLedger) processReservationRequest(
	ctx context.Context,
	header core.PaymentMetadata,
	numSymbols uint64,
	quorumNumbers []uint8,
	params *PaymentVaultParams,
	receivedAt time.Time,
) (*big.Int, error) {
	sal.logger.Debug("Recording and validating reservation usage", "header", header, "reservations", sal.reservations)

	// Load period records if needed (lazy loading with proper params)
	if len(sal.periodRecords) == 0 && len(sal.reservations) > 0 {
		if err := sal.loadPeriodRecords(ctx, params); err != nil {
			sal.logger.Warn("Failed to load period records", "accountID", sal.accountID.Hex(), "error", err)
		}
	}

	// Convert uint8 quorum numbers to QuorumID for consistency with existing reservations map
	quorumIDs := make([]core.QuorumID, len(quorumNumbers))
	for i, q := range quorumNumbers {
		quorumIDs[i] = core.QuorumID(q)
	}

	// Validate reservations exactly like meterer does
	if err := payment_logic.ValidateReservations(sal.reservations, params.QuorumProtocolConfigs, quorumNumbers, header.Timestamp, receivedAt.UnixNano()); err != nil {
		return nil, fmt.Errorf("invalid reservation: %w", err)
	}

	// Use meterer's incrementBinUsage logic (adapted for single account)
	if err := sal.incrementBinUsage(ctx, header, sal.reservations, params, numSymbols); err != nil {
		return nil, fmt.Errorf("failed to increment bin usages: %w", err)
	}

	// Reservation successful - no payment charged
	return nil, nil
}

// incrementBinUsage replicates meterer's incrementBinUsage logic for single account
func (sal *ServerAccountLedger) incrementBinUsage(
	ctx context.Context,
	header core.PaymentMetadata,
	reservations map[core.QuorumID]*core.ReservedPayment,
	globalParams *PaymentVaultParams,
	numSymbols uint64,
) error {
	charges := make(map[core.QuorumID]uint64)
	quorumNumbers := make([]core.QuorumID, 0, len(reservations))
	reservationWindows := make(map[core.QuorumID]uint64, len(reservations))
	requestReservationPeriods := make(map[core.QuorumID]uint64, len(reservations))

	for quorumID := range reservations {
		_, protocolConfig, err := globalParams.GetQuorumConfigs(quorumID)
		if err != nil {
			return fmt.Errorf("failed to get quorum config for quorum %d: %w", quorumID, err)
		}
		charges[quorumID] = payment_logic.SymbolsCharged(numSymbols, protocolConfig.MinNumSymbols)
		quorumNumbers = append(quorumNumbers, quorumID)
		reservationWindows[quorumID] = protocolConfig.ReservationRateLimitWindow
		requestReservationPeriods[quorumID] = payment_logic.GetReservationPeriodByNanosecond(header.Timestamp, protocolConfig.ReservationRateLimitWindow)
	}

	// Batch increment all quorums for the current quorums' reservation period
	updatedUsages := make(map[core.QuorumID]uint64)
	usage, err := sal.meteringStore.IncrementBinUsages(ctx, header.AccountID, quorumNumbers, requestReservationPeriods, charges)
	if err != nil {
		return err
	}
	for _, quorumID := range quorumNumbers {
		updatedUsages[quorumID] = usage[quorumID]
	}

	overflowAmounts := make(map[core.QuorumID]uint64)
	overflowPeriods := make(map[core.QuorumID]uint64)

	for quorumID, reservation := range reservations {
		reservationWindow := reservationWindows[quorumID]
		requestReservationPeriod := requestReservationPeriods[quorumID]
		usageLimit := payment_logic.GetBinLimit(reservation.SymbolsPerSecond, reservationWindow)
		newUsage, ok := updatedUsages[quorumID]
		if !ok {
			return fmt.Errorf("failed to get updated usage for quorum %d", quorumID)
		}
		prevUsage := newUsage - charges[quorumID]
		if newUsage <= usageLimit {
			continue
		} else if prevUsage >= usageLimit {
			// Bin was already filled before this increment
			return fmt.Errorf("bin has already been filled for quorum %d", quorumID)
		}
		overflowPeriod := payment_logic.GetOverflowPeriod(requestReservationPeriod, reservationWindow)
		if charges[quorumID] <= usageLimit && overflowPeriod <= payment_logic.GetReservationPeriod(int64(reservation.EndTimestamp), reservationWindow) {
			// Needs to go to overflow bin
			overflowAmount := newUsage - usageLimit
			overflowAmounts[quorumID] = overflowAmount
			overflowPeriods[quorumID] = overflowPeriod
		} else {
			return fmt.Errorf("overflow usage exceeds bin limit for quorum %d", quorumID)
		}
	}
	if len(overflowAmounts) != len(overflowPeriods) {
		return fmt.Errorf("overflow amount and period mismatch")
	}

	// Batch increment overflow bins for all overflown reservation candidates
	if len(overflowAmounts) > 0 {
		sal.logger.Debug("Utilizing reservation overflow period", "overflowAmounts", overflowAmounts, "overflowPeriods", overflowPeriods)
		overflowQuorums := make([]core.QuorumID, 0, len(overflowAmounts))
		for quorumID := range overflowAmounts {
			overflowQuorums = append(overflowQuorums, quorumID)
		}
		_, err := sal.meteringStore.IncrementBinUsages(ctx, header.AccountID, overflowQuorums, overflowPeriods, overflowAmounts)
		if err != nil {
			// Rollback the increments for the current periods
			rollbackErr := sal.meteringStore.DecrementBinUsages(ctx, header.AccountID, quorumNumbers, requestReservationPeriods, charges)
			if rollbackErr != nil {
				return fmt.Errorf("failed to increment overflow bins: %w; rollback also failed: %v", err, rollbackErr)
			}
			return fmt.Errorf("failed to increment overflow bins: %w; successfully rolled back increments", err)
		}
	}

	return nil
}

// processOnDemandRequest handles on-demand requests
// This replicates meterer.serveOnDemandRequest logic
func (sal *ServerAccountLedger) processOnDemandRequest(
	ctx context.Context,
	header core.PaymentMetadata,
	numSymbols uint64,
	quorumNumbers []uint8,
	params *PaymentVaultParams,
	receivedAt time.Time,
) (*big.Int, error) {
	sal.logger.Debug("Recording and validating on-demand usage", "header", header, "onDemandPayment", sal.onDemandPayment)

	if err := payment_logic.ValidateQuorum(quorumNumbers, params.OnDemandQuorumNumbers); err != nil {
		return nil, fmt.Errorf("invalid quorum for On-Demand Request: %w", err)
	}

	// Verify that the claimed cumulative payment doesn't exceed the on-chain deposit
	if header.CumulativePayment.Cmp(sal.onDemandPayment.CumulativePayment) > 0 {
		return nil, fmt.Errorf("request claims a cumulative payment greater than the on-chain deposit")
	}

	paymentConfig, protocolConfig, err := params.GetQuorumConfigs(OnDemandQuorumID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment config for on-demand quorum: %w", err)
	}

	symbolsCharged := payment_logic.SymbolsCharged(numSymbols, protocolConfig.MinNumSymbols)
	paymentCharged := payment_logic.PaymentCharged(symbolsCharged, paymentConfig.OnDemandPricePerSymbol)
	oldPayment, err := sal.meteringStore.AddOnDemandPayment(ctx, header, paymentCharged)
	if err != nil {
		return nil, fmt.Errorf("failed to update cumulative payment: %w", err)
	}

	// Store the old payment for potential rollback by RevertDebit
	sal.lastOldPayment = oldPayment

	// Update bin usage atomically and check against bin capacity
	// Note: This is where we would call global bin usage validation
	// but that should be handled at the ServerLedger level

	// Update local cached state - only update our internal cumulative payment tracking
	// NOTE: sal.onDemandPayment.CumulativePayment represents the on-chain deposit and should not be modified
	sal.cumulativePayment = header.CumulativePayment

	sal.logger.Debug("Processed on-demand payment",
		"accountID", sal.accountID.Hex(),
		"symbolsCharged", symbolsCharged,
		"paymentCharged", paymentCharged.String(),
		"oldPayment", oldPayment.String(),
		"newCumulative", header.CumulativePayment.String())

	return paymentCharged, nil
}

// RevertDebit rollback a previous debit operation
func (sal *ServerAccountLedger) RevertDebit(
	ctx context.Context,
	header core.PaymentMetadata,
	numSymbols uint64,
	quorumNumbers []uint8,
	params *PaymentVaultParams,
	receivedAt time.Time,
	payment *big.Int,
) error {
	// Convert uint8 quorum numbers to QuorumID for consistency
	quorumIDs := make([]core.QuorumID, len(quorumNumbers))
	quorumPeriods := make(map[core.QuorumID]uint64)
	quorumSizes := make(map[core.QuorumID]uint64)

	for i, q := range quorumNumbers {
		quorumID := core.QuorumID(q)
		quorumIDs[i] = quorumID

		// Get the proper reservation window for this quorum
		_, protocolConfig, err := params.GetQuorumConfigs(quorumID)
		if err != nil {
			sal.logger.Warn("Failed to get quorum config for revert", "quorumID", quorumID, "error", err)
			// Use a default if we can't get the config
			quorumPeriods[quorumID] = payment_logic.GetReservationPeriodByNanosecond(header.Timestamp, 3600)
			quorumSizes[quorumID] = numSymbols // Fallback to raw symbols
		} else {
			// Calculate period using the proper reservation window
			quorumPeriods[quorumID] = payment_logic.GetReservationPeriodByNanosecond(header.Timestamp, protocolConfig.ReservationRateLimitWindow)
			// Calculate charged symbols exactly like meterer does
			quorumSizes[quorumID] = payment_logic.SymbolsCharged(numSymbols, protocolConfig.MinNumSymbols)
		}
	}

	// Rollback usage tracking
	err := sal.meteringStore.DecrementBinUsages(ctx, sal.accountID, quorumIDs, quorumPeriods, quorumSizes)
	if err != nil {
		return fmt.Errorf("failed to rollback bin usages: %w", err)
	}

	// If this was a payment, rollback the payment
	if payment != nil && payment.Cmp(big.NewInt(0)) > 0 && sal.lastOldPayment != nil {
		newPayment := header.CumulativePayment  // The payment that was written to database during Debit
		oldPayment := sal.lastOldPayment

		err = sal.meteringStore.RollbackOnDemandPayment(ctx, sal.accountID, newPayment, oldPayment)
		if err != nil {
			return fmt.Errorf("failed to rollback on-demand payment: %w", err)
		}

		// Update local cached state - only update our internal cumulative payment tracking
		// NOTE: sal.onDemandPayment.CumulativePayment represents the on-chain deposit and should not be modified
		sal.cumulativePayment = oldPayment
		
		// Clear the stored old payment since we've used it
		sal.lastOldPayment = nil
	}

	return nil
}

// GetAccountStateProtobuf returns the account state in protobuf format
func (sal *ServerAccountLedger) GetAccountStateProtobuf() (
	reservations map[uint32]*disperser_v2.QuorumReservation,
	periodRecords map[uint32]*disperser_v2.PeriodRecords,
	onchainCumulativePayment []byte,
	cumulativePayment []byte,
) {
	// Convert reservations to protobuf format
	reservationsProto := make(map[uint32]*disperser_v2.QuorumReservation)
	for quorumID, reservation := range sal.reservations {
		if reservation != nil {
			reservationsProto[uint32(quorumID)] = &disperser_v2.QuorumReservation{
				SymbolsPerSecond: reservation.SymbolsPerSecond,
				StartTimestamp:   uint32(reservation.StartTimestamp),
				EndTimestamp:     uint32(reservation.EndTimestamp),
			}
		}
	}

	// Convert period records to protobuf format
	periodRecordsProto := make(map[uint32]*disperser_v2.PeriodRecords)
	for quorumID, records := range sal.periodRecords {
		if len(records) > 0 {
			protoRecords := make([]*disperser_v2.PeriodRecord, len(records))
			for i, record := range records {
				protoRecords[i] = &disperser_v2.PeriodRecord{
					Index: record.Index,
					Usage: record.Usage,
				}
			}
			periodRecordsProto[uint32(quorumID)] = &disperser_v2.PeriodRecords{
				Records: protoRecords,
			}
		}
	}

	// Convert payments to bytes
	var onchainPaymentBytes []byte
	var cumulativePaymentBytes []byte

	if sal.onDemandPayment != nil && sal.onDemandPayment.CumulativePayment != nil {
		onchainPaymentBytes = sal.onDemandPayment.CumulativePayment.Bytes()
	}

	if sal.cumulativePayment != nil {
		cumulativePaymentBytes = sal.cumulativePayment.Bytes()
	}

	return reservationsProto, periodRecordsProto, onchainPaymentBytes, cumulativePaymentBytes
}

// RefreshFromChain refreshes the account's on-chain state
func (sal *ServerAccountLedger) RefreshFromChain(ctx context.Context, chainPaymentState OnchainPayment) error {
	// Fetch fresh on-chain state
	reservations, err := chainPaymentState.GetReservedPaymentByAccountAndQuorums(ctx, sal.accountID, []core.QuorumID{})
	if err != nil {
		return fmt.Errorf("failed to refresh reservations: %w", err)
	}

	onDemandPayment, err := chainPaymentState.GetOnDemandPaymentByAccount(ctx, sal.accountID)
	if err != nil {
		return fmt.Errorf("failed to refresh on-demand payment: %w", err)
	}

	// Update cached state
	sal.reservations = reservations
	sal.onDemandPayment = onDemandPayment

	sal.logger.Debug("Refreshed account state from chain", "accountID", sal.accountID.Hex())
	return nil
}
