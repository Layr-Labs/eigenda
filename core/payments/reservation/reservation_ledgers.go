package reservation

import (
	"fmt"
	"sync"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// ReservationLedgers manages and validates reservation payments for multiple accounts
type ReservationLedgers struct {
	// map from account ID to ReservationLedger
	ledgers map[gethcommon.Address]*ReservationLedger

	// lock protects concurrent access to the ledgers map
	lock sync.Mutex

	// timeSource is a function that returns the current time
	// This allows for easier testing and consistent time handling
	timeSource func() time.Time
}

// NewReservationPaymentValidator creates a new ReservationPaymentValidator
func NewReservationPaymentValidator(timeSource func() time.Time) *ReservationLedgers {
	return &ReservationLedgers{
		ledgers:    make(map[gethcommon.Address]*ReservationLedger),
		timeSource: timeSource,
	}
}

// Debit validates a reservation payment for a blob dispersal
// The caller is responsible for verifying the signature before calling this method
func (rl *ReservationLedgers) Debit(
	accountID gethcommon.Address,
	symbolCount uint32,
	quorumNumbers []uint8,
	dispersalTime time.Time,
) error {
	ledger, err := rl.getOrCreateLedger(accountID)
	if err != nil {
		return fmt.Errorf("get or create reservation ledger: %w", err)
	}

	// For reservation payments, first call CheckInvariants, then Debit
	// First check invariants
	err = ledger.CheckInvariants(quorumNumbers, dispersalTime)
	if err != nil {
		return fmt.Errorf("reservation payment failed invariant checks: %w", err)
	}

	// For Debit, use current time (not the dispersal time)
	now := rl.timeSource()
	success, err := ledger.Debit(now, symbolCount)
	if err != nil {
		return fmt.Errorf("debit reservation payment: %w", err)
	}

	if !success {
		return fmt.Errorf("reservation debit failed: insufficient capacity")
	}

	// TODO: Consider in what cases we should remove the ledger from the map
	// Possible cases:
	// - Reservation has expired
	// - Account has been inactive for a certain period
	// - Explicit cleanup request

	return nil
}

// getOrCreateLedger gets an existing reservation ledger or creates a new one if it doesn't exist
func (rl *ReservationLedgers) getOrCreateLedger(accountID gethcommon.Address) (*ReservationLedger, error) {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	if ledger, exists := rl.ledgers[accountID]; exists {
		return ledger, nil
	}

	// TODO: These are placeholder values - need to get actual reservation from chain or config
	// Creating placeholder reservation config
	// NOTE: ReservationLedger requires a ReservationLedgerConfig and a time
	// When implementing, use rl.timeSource() for the time parameter
	// This is a placeholder implementation
	return nil, fmt.Errorf("reservation ledger for account %s not found - placeholder creation not implemented", accountID.Hex())
}
