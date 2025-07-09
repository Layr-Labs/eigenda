package meterer

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	disperser_v2 "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer/payment_logic"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// LocalAccountLedger implements AccountLedger for client-side dispersal usage tracking.
//
// Architecture and Threading:
//   - Uses sync.RWMutex: multiple concurrent readers, exclusive writers
//   - All state mutations are atomic via deep-copy-and-commit pattern
//   - Read operations (GetAccountStateProtobuf) use read locks for efficiency
//   - Write operations (Debit, RevertDebit) acquire exclusive write locks
//
// Reservation Rate Limiting:
//   - Tracks symbol usage in time-based period bins
//   - Implements overflow bin logic: when current period exceeds limit for the first time, uses overflow period
//   - Overflow period calculated as: current_period + 2 * reservation_window
//   - MinNumBins bins per quorum with modulo-based circular indexing
//   - Rate limits enforce: symbols_per_second * window_size as maximum per period
//
// Multi-Quorum Transaction Semantics:
//   - All-or-nothing: either ALL quorums succeed with reservations or ALL fall back to on-demand
//   - Atomic rollback: partial failures revert all state changes via deep copy restoration
//   - Validation hierarchy: config errors prevent fallback, usage errors allow fallback
//
// Payment State Synchronization:
//   - onDemand.CumulativePayment: on-chain balance from PaymentVault contract
//   - cumulativePayment: local consumed amount, must not exceed on-chain balance
//   - Protobuf serialization matches disperser GetPaymentStateForAllQuorums RPC format
//   - Compatible with Accountant behavior for seamless client-server state transitions
type LocalAccountLedger struct {
	// reservations stores the per-quorum reserved payment settings from on-chain state.
	// Protected by mutex for concurrent access.
	reservations map[core.QuorumID]*core.ReservedPayment

	// onDemand stores the on-demand payment settings from on-chain state.
	// Protected by mutex for concurrent access.
	onDemand *core.OnDemandPayment

	// periodRecords tracks the local usage history across time periods for each quorum.
	// Protected by mutex for concurrent access.
	periodRecords QuorumPeriodRecords

	// cumulativePayment tracks the local cumulative payment amount consumed.
	// Protected by mutex for concurrent access.
	cumulativePayment *big.Int

	// mutex protects concurrent access to all ledger state.
	// Uses RWMutex to allow concurrent reads while ensuring exclusive writes.
	mutex sync.RWMutex
}

// NewLocalAccountLedger creates a new LocalAccountLedger with zero-initialized state.
func NewLocalAccountLedger() *LocalAccountLedger {
	return &LocalAccountLedger{
		reservations:      make(map[core.QuorumID]*core.ReservedPayment),
		onDemand:          &core.OnDemandPayment{CumulativePayment: big.NewInt(0)},
		periodRecords:     make(QuorumPeriodRecords),
		cumulativePayment: big.NewInt(0),
	}
}

// NewLocalAccountLedgerFromProtobuf creates a LocalAccountLedger from protobuf state components.
// Deserializes account state received from disperser GetPaymentStateForAllQuorums RPC responses.
// Converts wire-format protobuf data directly to in-memory structures without intermediate objects.
// Returns error if big.Int byte arrays are malformed or if required protobuf fields are invalid.
func NewLocalAccountLedgerFromProtobuf(
	reservations map[uint32]*disperser_v2.QuorumReservation,
	periodRecords map[uint32]*disperser_v2.PeriodRecords,
	onchainCumulativePayment []byte,
	cumulativePayment []byte,
) (*LocalAccountLedger, error) {
	// Convert protobuf reservations to Go structures
	goReservations := make(map[core.QuorumID]*core.ReservedPayment)
	for quorumID, protoReservation := range reservations {
		if protoReservation != nil {
			goReservations[core.QuorumID(quorumID)] = &core.ReservedPayment{
				SymbolsPerSecond: protoReservation.SymbolsPerSecond,
				StartTimestamp:   uint64(protoReservation.StartTimestamp),
				EndTimestamp:     uint64(protoReservation.EndTimestamp),
			}
		}
	}

	// Convert on-chain cumulative payment
	var onchainPayment *big.Int
	if len(onchainCumulativePayment) > 0 {
		onchainPayment = new(big.Int).SetBytes(onchainCumulativePayment)
	} else {
		onchainPayment = big.NewInt(0)
	}

	// Convert cumulative payment
	var localPayment *big.Int
	if len(cumulativePayment) > 0 {
		localPayment = new(big.Int).SetBytes(cumulativePayment)
	} else {
		localPayment = big.NewInt(0)
	}

	// Convert period records
	goPeriodRecords := make(QuorumPeriodRecords)
	for quorumID, protoPeriodRecords := range periodRecords {
		if protoPeriodRecords != nil && len(protoPeriodRecords.Records) > 0 {
			records := make([]*PeriodRecord, len(protoPeriodRecords.Records))
			for i, protoRecord := range protoPeriodRecords.Records {
				records[i] = &PeriodRecord{
					Index: protoRecord.Index,
					Usage: protoRecord.Usage,
				}
			}
			goPeriodRecords[core.QuorumID(quorumID)] = records
		}
	}

	return &LocalAccountLedger{
		reservations:      goReservations,
		onDemand:          &core.OnDemandPayment{CumulativePayment: onchainPayment},
		periodRecords:     goPeriodRecords,
		cumulativePayment: localPayment,
	}, nil
}

// Debit implements AccountLedger.Debit with atomic multi-quorum transaction semantics.
// Acquires exclusive write lock, validates all quorums, applies deep-copy-commit pattern.
// Uses all-or-nothing logic: either ALL quorums use reservations or ALL use on-demand.
func (lal *LocalAccountLedger) Debit(
	ctx context.Context,
	accountID gethcommon.Address,
	timestampNs int64,
	numSymbols uint64,
	quorumNumbers []core.QuorumID,
	params *PaymentVaultParams,
) (*big.Int, error) {
	lal.mutex.Lock()
	defer lal.mutex.Unlock()

	// Store reservation error for potential fallback error message (matches Accountant behavior)
	var reservationError error

	// First, try reservation usage for ALL quorums (matches Accountant behavior)
	// Validate all reservations exist and are valid for this timestamp
	reservationValidationErr := payment_logic.ValidateReservations(lal.reservations, params.QuorumProtocolConfigs, quorumNumbers, timestampNs, time.Now().UnixNano())
	if reservationValidationErr == nil {
		// Create a deep copy for atomic updates
		periodRecordsCopy := lal.periodRecords.DeepCopy()

		var reservationUsageErr error
		for _, quorumNumber := range quorumNumbers {
			reservation := lal.reservations[quorumNumber]
			_, protocolConfig, err := params.GetQuorumConfigs(quorumNumber)
			if err != nil {
				reservationUsageErr = err
				break
			}
			// UpdateUsage includes overflow logic - this is where overflow bins are used!
			if err := periodRecordsCopy.UpdateUsage(quorumNumber, timestampNs, numSymbols, reservation, protocolConfig); err != nil {
				reservationUsageErr = err
				break
			}
		}

		if reservationUsageErr == nil {
			// All quorums used reservations successfully, commit and return nil (no payment)
			lal.periodRecords = periodRecordsCopy
			return nil, nil
		}

		// Reservations exist and are valid, but usage failed (limit exceeded, etc.)
		// Store the reservation error and allow fallback to on-demand (matches Accountant behavior)
		reservationError = reservationUsageErr
	} else {
		// Reservations are unavailable/invalid for ALL quorums
		// Check if this is a config error (should not try on-demand fallback)
		if isConfigError(reservationValidationErr) {
			return nil, fmt.Errorf("no payment method available: reservation failed and %v", reservationValidationErr)
		}
		// Store the reservation validation error for potential combined error message
		reservationError = reservationValidationErr
	}

	// Try on-demand for ALL quorums (either reservation validation failed or usage failed)
	onDemandErr := payment_logic.ValidateQuorum(quorumNumbers, params.OnDemandQuorumNumbers)
	if onDemandErr != nil {
		// Both reservation and on-demand failed, return combined error (matches Accountant)
		return nil, fmt.Errorf("cannot create payment information for reservation or on-demand. Account: %s, Reservation Error: %w, On-demand Error: %w", accountID.Hex(), reservationError, onDemandErr)
	}

	paymentQuorumConfig, protocolConfig, err := params.GetQuorumConfigs(OnDemandQuorumID)
	if err != nil {
		// On-demand config error, return combined error
		return nil, fmt.Errorf("cannot create payment information for reservation or on-demand. Account: %s, Reservation Error: %w, On-demand Error: %w", accountID.Hex(), reservationError, err)
	}

	symbolsCharged := payment_logic.SymbolsCharged(numSymbols, protocolConfig.MinNumSymbols)
	paymentCharged := payment_logic.PaymentCharged(symbolsCharged, paymentQuorumConfig.OnDemandPricePerSymbol)

	// Calculate the increment required to add to the cumulative payment
	resultingPayment := new(big.Int).Add(lal.cumulativePayment, paymentCharged)
	if resultingPayment.Cmp(lal.onDemand.CumulativePayment) <= 0 {
		lal.cumulativePayment.Add(lal.cumulativePayment, paymentCharged)
		return new(big.Int).Set(lal.cumulativePayment), nil
	}

	// On-demand payment insufficient, return combined error (matches Accountant)
	return nil, fmt.Errorf("cannot create payment information for reservation or on-demand. Account: %s, Reservation Error: %w, On-demand Error: %w", accountID.Hex(), reservationError, fmt.Errorf("insufficient ondemand payment"))
}

// RevertDebit implements AccountLedger.RevertDebit for undoing failed or cancelled operations.
// Distinguishes between reservation reverts (subtracts from period usage) and on-demand reverts
// (subtracts from cumulative payment). Validates revert amounts to prevent negative balances.
func (lal *LocalAccountLedger) RevertDebit(
	ctx context.Context,
	accountID gethcommon.Address,
	timestampNs int64,
	numSymbols uint64,
	quorumNumbers []core.QuorumID,
	params *PaymentVaultParams,
	payment *big.Int,
) error {
	lal.mutex.Lock()
	defer lal.mutex.Unlock()

	// If payment is nil, this was a reservation usage - revert period records
	if payment == nil {
		// Validate reservations exist for reversal
		if err := payment_logic.ValidateReservations(lal.reservations, params.QuorumProtocolConfigs, quorumNumbers, timestampNs, time.Now().UnixNano()); err != nil {
			return fmt.Errorf("cannot revert reservation usage: %v", err)
		}

		// Create a deep copy for atomic updates
		periodRecordsCopy := lal.periodRecords.DeepCopy()

		for _, quorumNumber := range quorumNumbers {
			_, exists := lal.reservations[quorumNumber]
			if !exists {
				return fmt.Errorf("cannot revert: reservation not found for quorum %d", quorumNumber)
			}
			_, protocolConfig, err := params.GetQuorumConfigs(quorumNumber)
			if err != nil {
				return fmt.Errorf("cannot revert: %v", err)
			}

			// Calculate the period for this timestamp
			reservationPeriod := payment_logic.GetReservationPeriodByNanosecond(timestampNs, protocolConfig.ReservationRateLimitWindow)

			// Get period records for this quorum
			if periodRecordsCopy[quorumNumber] == nil {
				return fmt.Errorf("cannot revert: no period records found for quorum %d", quorumNumber)
			}

			records := periodRecordsCopy[quorumNumber]

			// Find the record with matching period index
			found := false
			for _, record := range records {
				if record.Index == uint32(reservationPeriod) {
					if record.Usage < numSymbols {
						return fmt.Errorf("cannot revert: insufficient usage to subtract (%d < %d) for quorum %d", record.Usage, numSymbols, quorumNumber)
					}
					record.Usage -= numSymbols
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf("cannot revert: no usage record found for period %d, quorum %d", reservationPeriod, quorumNumber)
			}
		}

		// All reversals succeeded, commit the changes
		lal.periodRecords = periodRecordsCopy
		return nil
	}

	// Payment is not nil, this was on-demand usage - revert cumulative payment
	if payment.Cmp(big.NewInt(0)) <= 0 {
		return fmt.Errorf("cannot revert: invalid payment amount %v", payment)
	}

	// Validate the payment can be reverted
	if lal.cumulativePayment.Cmp(payment) < 0 {
		return fmt.Errorf("cannot revert: insufficient cumulative payment (%v < %v)", lal.cumulativePayment, payment)
	}

	// Subtract the payment from cumulative payment
	lal.cumulativePayment.Sub(lal.cumulativePayment, payment)
	return nil
}

// GetAccountStateProtobuf implements AccountLedger.GetAccountStateProtobuf for state extraction.
// Uses read lock for concurrent access. Converts internal types to wire-format protobuf structures.
// Filters out nil period records and converts big.Int values to byte arrays for transmission.
func (lal *LocalAccountLedger) GetAccountStateProtobuf() (
	reservations map[uint32]*disperser_v2.QuorumReservation,
	periodRecords map[uint32]*disperser_v2.PeriodRecords,
	onchainCumulativePayment []byte,
	cumulativePayment []byte,
) {
	lal.mutex.RLock()
	defer lal.mutex.RUnlock()

	// Convert reservations to protobuf format
	protoReservations := make(map[uint32]*disperser_v2.QuorumReservation)
	for quorumID, reservation := range lal.reservations {
		protoReservations[uint32(quorumID)] = &disperser_v2.QuorumReservation{
			SymbolsPerSecond: reservation.SymbolsPerSecond,
			StartTimestamp:   uint32(reservation.StartTimestamp),
			EndTimestamp:     uint32(reservation.EndTimestamp),
		}
	}

	// Convert period records to protobuf format
	protoPeriodRecords := make(map[uint32]*disperser_v2.PeriodRecords)
	for quorumID, records := range lal.periodRecords {
		if len(records) > 0 {
			protoRecords := make([]*disperser_v2.PeriodRecord, 0, len(records))
			for _, record := range records {
				if record != nil {
					protoRecords = append(protoRecords, &disperser_v2.PeriodRecord{
						Index: record.Index,
						Usage: record.Usage,
					})
				}
			}
			protoPeriodRecords[uint32(quorumID)] = &disperser_v2.PeriodRecords{
				Records: protoRecords,
			}
		}
	}

	// Convert big.Int values to byte arrays
	var onchainCumulativePaymentBytes []byte
	if lal.onDemand != nil && lal.onDemand.CumulativePayment != nil {
		onchainCumulativePaymentBytes = lal.onDemand.CumulativePayment.Bytes()
	}

	var cumulativePaymentBytes []byte
	if lal.cumulativePayment != nil {
		cumulativePaymentBytes = lal.cumulativePayment.Bytes()
	}

	return protoReservations, protoPeriodRecords, onchainCumulativePaymentBytes, cumulativePaymentBytes
}

// isConfigError determines if an error is related to configuration issues
// (which should not allow on-demand fallback) vs timing/availability issues (which should allow fallback)
func isConfigError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := err.Error()
	configErrorIndicators := []string{
		"config not found",
		"no quorum numbers provided",
		"payment config not found",
		"protocol config not found",
	}

	for _, indicator := range configErrorIndicators {
		if strings.Contains(errMsg, indicator) {
			return true
		}
	}

	return false
}
