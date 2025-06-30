package meterer

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// AccountLedger defines the standard interface for tracking account payment state,
// including both on-chain settings and local usage tracking. This interface abstracts
// the differences between client-side in-memory tracking and server-side persistent storage.
//
// Implementation Patterns:
//   - LocalAccountLedger: In-memory tracking for clients
//   - Future: DatabaseAccountLedger for persistent server-side tracking
//   - Future: DistributedAccountLedger for consensus-based tracking
type AccountLedger interface {
	// RecordReservationUsage attempts to record symbol usage against the ledger's stored
	// reservations for the specified quorums. It atomically checks usage against bin limits,
	// handles overflows, and updates the underlying storage.
	//
	// The method performs the following operations atomically:
	//   1. Validates that reservations exist and are currently active for all quorums
	//   2. Calculates the time period bin for the given timestamp
	//   3. Checks if the requested symbols fit within rate limits for each quorum
	//   4. Handles overflow to future periods if current period is full
	//   5. Updates all period records or rolls back entirely on any failure
	//
	// Parameters:
	//   ctx: Context for cancellation and timeouts (currently unused in LocalAccountLedger)
	//   accountID: The account identifier (used for logging and future database implementations)
	//   timestampNs: The timestamp in nanoseconds when the usage occurred
	//   numSymbols: The number of symbols to be charged against reservations
	//   quorumNumbers: The list of quorums that will be charged for this usage
	//   params: Payment vault parameters containing quorum configurations
	//
	// Returns:
	//   nil if the usage was successfully recorded
	//   error if the usage cannot be accommodated:
	//     - "reservation not found for quorum X": No reservation exists for quorum
	//     - "reservation limit exceeded": Usage would exceed rate limits
	//     - "reservation expired": Current time is outside reservation window
	//     - "quorum config not found for X": Missing protocol configuration
	RecordReservationUsage(
		ctx context.Context,
		accountID gethcommon.Address,
		timestampNs int64,
		numSymbols uint64,
		quorumNumbers []core.QuorumID,
		params *PaymentVaultParams,
	) error

	// RecordOnDemandUsage attempts to record on-demand payment usage against the ledger's
	// stored on-demand payment state and returns the new cumulative payment if successful.
	//
	// The method performs the following operations:
	//   1. Validates that the requested quorums are enabled for on-demand usage
	//   2. Calculates the payment amount based on symbols and pricing configuration
	//   3. Checks if the payment would exceed the available on-demand balance
	//   4. Updates the cumulative payment if sufficient balance exists
	//
	// Payment Calculation:
	//   symbolsCharged = max(numSymbols, minNumSymbols)
	//   paymentCharged = symbolsCharged * onDemandPricePerSymbol
	//   newCumulative = currentCumulative + paymentCharged
	//
	// Parameters:
	//   ctx: Context for cancellation and timeouts (currently unused in LocalAccountLedger)
	//   accountID: The account identifier (used for logging and future database implementations)
	//   numSymbols: The number of symbols to be charged
	//   quorumNumbers: The list of quorums for this usage (must be on-demand enabled)
	//   params: Payment vault parameters containing pricing and quorum configurations
	//
	// Returns:
	//   (*big.Int, nil): The new cumulative payment amount if successful
	//   (nil, error): If the usage cannot be accommodated:
	//     - "quorum X not enabled for on-demand": Requested quorum doesn't support on-demand
	//     - "insufficient ondemand payment": Would exceed available balance
	//     - "quorum config not found for X": Missing payment configuration
	RecordOnDemandUsage(
		ctx context.Context,
		accountID gethcommon.Address,
		numSymbols uint64,
		quorumNumbers []core.QuorumID,
		params *PaymentVaultParams,
	) (*big.Int, error)

	// SetAccountState updates the ledger with complete account state including
	// on-chain settings (reservations, on-demand) and local tracking data.
	//
	// This method completely replaces the current ledger state with the provided state.
	// It performs deep copying of all data to ensure isolation between the provided
	// state and the internal ledger state.
	//
	// Use Cases:
	//   - Initial setup: Loading account state from on-chain data
	//   - Synchronization: Updating local state with fresh on-chain data
	//   - Testing: Setting up known state for test scenarios
	//
	// Parameters:
	//   state: The complete account state to set. Nil fields will be initialized
	//          with appropriate zero values (empty maps, zero big.Int, etc.)
	//
	// Thread Safety:
	// This method is safe for concurrent use. It acquires an exclusive write lock
	// for the entire operation to ensure atomic state replacement.
	SetAccountState(state AccountState)

	// GetAccountState returns the complete current account state including on-chain
	// settings and local tracking data for inspection and synchronization.
	//
	// The returned state is a deep copy of the internal ledger state, ensuring that
	// modifications to the returned data will not affect the ledger's internal state.
	// This prevents accidental mutations while allowing safe inspection and analysis.
	//
	// Use Cases:
	//   - State inspection: Examining current payment usage and limits
	//   - Synchronization: Extracting state for storage or transmission
	//   - Testing: Verifying state after operations
	//   - Debugging: Analyzing account state for troubleshooting
	//
	// Returns:
	//   AccountState: A deep copy of the current account state including:
	//     - All reservation settings and their current usage
	//     - On-demand payment balance and cumulative usage
	//     - Complete period records history
	GetAccountState() AccountState
}

// AccountState represents the complete payment state for an account,
// including both on-chain settings and local usage tracking.
//
// This struct serves as the canonical representation of all payment-related
// data for a single account, enabling atomic state transitions and consistent
// snapshots for synchronization between different system components.
//
// State Categories:
//  1. On-chain Settings: Immutable data from blockchain (reservations, on-demand)
//  2. Local Tracking: Mutable usage counters maintained locally (period records, cumulative payment)
type AccountState struct {
	// Reservations contains the per-quorum reserved payment settings from on-chain state.
	// Each reservation specifies the symbols per second rate limit and the time window
	// during which the reservation is valid. Nil entries indicate no reservation for that quorum.
	Reservations map[core.QuorumID]*core.ReservedPayment

	// OnDemand contains the on-demand payment settings from on-chain state.
	// This includes the cumulative payment balance available for on-demand usage.
	// Must not be nil - use zero balance if no on-demand payment is available.
	OnDemand *core.OnDemandPayment

	// PeriodRecords tracks the local usage history across time periods for each quorum.
	// This enables rate limiting by maintaining a sliding window of usage data.
	// The records are updated locally as symbols are consumed against reservations.
	PeriodRecords QuorumPeriodRecords

	// CumulativePayment tracks the local cumulative payment amount that has been
	// consumed from the on-demand balance. This value should never exceed the
	// OnDemand.CumulativePayment amount.
	CumulativePayment *big.Int
}

// LocalAccountLedger implements AccountLedger for client-side in-memory payment tracking.
// It manages complete account state including on-chain settings and local usage tracking.
//
// This implementation is designed for client applications that need to track their own
// payment usage locally without requiring persistent storage. It provides:
//   - Thread-safe concurrent access using RWMutex
//   - Atomic multi-quorum operations with rollback on failure
//   - Deep copying to prevent accidental state mutations
//   - Efficient read operations with read locks
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

// NewLocalAccountLedger creates a new LocalAccountLedger with empty state.
//
// The returned ledger is initialized with:
//   - Empty reservations map
//   - Zero on-demand payment balance
//   - Empty period records
//   - Zero cumulative payment
//
// Returns:
//
//	A new LocalAccountLedger instance with empty state
func NewLocalAccountLedger() *LocalAccountLedger {
	return &LocalAccountLedger{
		reservations:      make(map[core.QuorumID]*core.ReservedPayment),
		onDemand:          &core.OnDemandPayment{CumulativePayment: big.NewInt(0)},
		periodRecords:     make(QuorumPeriodRecords),
		cumulativePayment: big.NewInt(0),
	}
}

// RecordReservationUsage implements the AccountLedger interface for local in-memory tracking
func (lal *LocalAccountLedger) RecordReservationUsage(
	ctx context.Context,
	accountID gethcommon.Address,
	timestampNs int64,
	numSymbols uint64,
	quorumNumbers []core.QuorumID,
	params *PaymentVaultParams,
) error {
	lal.mutex.Lock()
	defer lal.mutex.Unlock()

	// Validate reservations using the stored state
	if err := ValidateReservations(lal.reservations, params.QuorumProtocolConfigs, quorumNumbers, timestampNs, time.Now().UnixNano()); err != nil {
		return err
	}

	// Create a deep copy for atomic updates - if any quorum fails, rollback entire operation
	periodRecordsCopy := lal.periodRecords.DeepCopy()

	for _, quorumNumber := range quorumNumbers {
		reservation, exists := lal.reservations[quorumNumber]
		if !exists {
			return fmt.Errorf("reservation not found for quorum %d", quorumNumber)
		}
		_, protocolConfig, err := params.GetQuorumConfigs(quorumNumber)
		if err != nil {
			return err
		}
		if err := periodRecordsCopy.UpdateUsage(quorumNumber, timestampNs, numSymbols, reservation, protocolConfig); err != nil {
			return err
		}
	}

	// All updates succeeded, commit the changes
	lal.periodRecords = periodRecordsCopy
	return nil
}

// RecordOnDemandUsage implements the AccountLedger interface for local on-demand payment tracking
func (lal *LocalAccountLedger) RecordOnDemandUsage(
	ctx context.Context,
	accountID gethcommon.Address,
	numSymbols uint64,
	quorumNumbers []core.QuorumID,
	params *PaymentVaultParams,
) (*big.Int, error) {
	// Validate quorums for on-demand usage
	if err := ValidateQuorum(quorumNumbers, params.OnDemandQuorumNumbers); err != nil {
		return nil, err
	}

	paymentQuorumConfig, protocolConfig, err := params.GetQuorumConfigs(OnDemandQuorumID)
	if err != nil {
		return nil, err
	}

	symbolsCharged := SymbolsCharged(numSymbols, protocolConfig.MinNumSymbols)
	paymentCharged := PaymentCharged(symbolsCharged, paymentQuorumConfig.OnDemandPricePerSymbol)

	lal.mutex.Lock()
	defer lal.mutex.Unlock()

	// Calculate the increment required to add to the cumulative payment
	resultingPayment := new(big.Int).Add(lal.cumulativePayment, paymentCharged)
	if resultingPayment.Cmp(lal.onDemand.CumulativePayment) <= 0 {
		lal.cumulativePayment.Add(lal.cumulativePayment, paymentCharged)
		return new(big.Int).Set(lal.cumulativePayment), nil
	}

	return nil, fmt.Errorf("insufficient ondemand payment")
}

// GetAccountState returns the complete current account state including on-chain settings
func (lal *LocalAccountLedger) GetAccountState() AccountState {
	lal.mutex.RLock()
	defer lal.mutex.RUnlock()

	// Deep copy reservations map
	reservationsCopy := make(map[core.QuorumID]*core.ReservedPayment)
	for quorumID, reservation := range lal.reservations {
		reservationsCopy[quorumID] = &core.ReservedPayment{
			SymbolsPerSecond: reservation.SymbolsPerSecond,
			StartTimestamp:   reservation.StartTimestamp,
			EndTimestamp:     reservation.EndTimestamp,
		}
	}

	// Deep copy on-demand payment
	var onDemandCopy *core.OnDemandPayment
	if lal.onDemand != nil {
		onDemandCopy = &core.OnDemandPayment{
			CumulativePayment: new(big.Int).Set(lal.onDemand.CumulativePayment),
		}
	}

	return AccountState{
		Reservations:      reservationsCopy,
		OnDemand:          onDemandCopy,
		PeriodRecords:     lal.periodRecords.DeepCopy(),
		CumulativePayment: new(big.Int).Set(lal.cumulativePayment),
	}
}

// SetAccountState updates the ledger with complete account state
func (lal *LocalAccountLedger) SetAccountState(state AccountState) {
	lal.mutex.Lock()
	defer lal.mutex.Unlock()

	// Set on-chain reservations
	if state.Reservations != nil {
		lal.reservations = make(map[core.QuorumID]*core.ReservedPayment)
		for quorumID, reservation := range state.Reservations {
			lal.reservations[quorumID] = &core.ReservedPayment{
				SymbolsPerSecond: reservation.SymbolsPerSecond,
				StartTimestamp:   reservation.StartTimestamp,
				EndTimestamp:     reservation.EndTimestamp,
			}
		}
	} else {
		lal.reservations = make(map[core.QuorumID]*core.ReservedPayment)
	}

	// Set on-chain on-demand payment state
	if state.OnDemand != nil {
		lal.onDemand = &core.OnDemandPayment{
			CumulativePayment: new(big.Int).Set(state.OnDemand.CumulativePayment),
		}
	} else {
		lal.onDemand = &core.OnDemandPayment{
			CumulativePayment: big.NewInt(0),
		}
	}

	// Set local tracking state
	if state.PeriodRecords != nil {
		lal.periodRecords = state.PeriodRecords
	} else {
		lal.periodRecords = make(QuorumPeriodRecords)
	}

	if state.CumulativePayment != nil {
		lal.cumulativePayment = new(big.Int).Set(state.CumulativePayment)
	} else {
		lal.cumulativePayment = big.NewInt(0)
	}
}
