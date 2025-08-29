package reservation

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// ReservationPaymentVaultState manages reservation payment state from the payment vault
type ReservationPaymentVaultState struct {
	tx               *eth.Reader
	logger           logging.Logger
	reservedPayments map[gethcommon.Address]*core.ReservedPayment
	lock             sync.RWMutex
}

// NewReservationPaymentVaultState creates a new ReservationPaymentVaultState
func NewReservationPaymentVaultState(tx *eth.Reader, logger logging.Logger) *ReservationPaymentVaultState {
	return &ReservationPaymentVaultState{
		tx:               tx,
		logger:           logger.With("component", "ReservationPaymentVaultState"),
		reservedPayments: make(map[gethcommon.Address]*core.ReservedPayment),
	}
}

// RefreshReservedPayments updates the cached reserved payments and returns a list of detected changes
// TODO: Replace periodic polling with event subscription from PaymentVault contract
// Should subscribe to ReservationCreated/ReservationUpdated events instead of polling GetReservedPayments
func (vs *ReservationPaymentVaultState) RefreshReservedPayments(ctx context.Context) ([]ReservationUpdate, error) {
	vs.lock.Lock()
	defer vs.lock.Unlock()

	if len(vs.reservedPayments) == 0 {
		vs.logger.Info("No reserved payments to refresh")
		return nil, nil
	}

	accountIDs := make([]gethcommon.Address, 0, len(vs.reservedPayments))
	for accountID := range vs.reservedPayments {
		accountIDs = append(accountIDs, accountID)
	}

	newReservedPayments, err := vs.tx.GetReservedPayments(ctx, accountIDs)
	if err != nil {
		return nil, fmt.Errorf("get reserved payments: %w", err)
	}

	// Detect changes and create updates
	var updates []ReservationUpdate
	for accountID, newPayment := range newReservedPayments {
		oldPayment, exists := vs.reservedPayments[accountID]
		if !exists || !reservedPaymentsEqual(oldPayment, newPayment) {
			// Create a Reservation object
			reservationObj, err := NewReservation(
				newPayment.SymbolsPerSecond,
				time.Unix(int64(newPayment.StartTimestamp), 0),
				time.Unix(int64(newPayment.EndTimestamp), 0),
				newPayment.QuorumNumbers,
			)
			if err != nil {
				vs.logger.Error("Failed to create reservation object", "error", err, "accountID", accountID.Hex())
				continue
			}

			update, err := NewReservationUpdate(accountID, reservationObj)
			if err != nil {
				vs.logger.Error("Failed to create reservation update", "error", err, "accountID", accountID.Hex())
				continue
			}
			updates = append(updates, *update)
		}
	}

	vs.reservedPayments = newReservedPayments
	return updates, nil
}

// GetReservedPaymentByAccount retrieves reservation payment info for a specific account
func (vs *ReservationPaymentVaultState) GetReservedPaymentByAccount(
	ctx context.Context,
	accountID gethcommon.Address,
) (*core.ReservedPayment, error) {
	vs.lock.RLock()
	if reservation, ok := vs.reservedPayments[accountID]; ok {
		vs.lock.RUnlock()
		return reservation, nil
	}
	vs.lock.RUnlock()

	// pulls the chain state
	res, err := vs.tx.GetReservedPaymentByAccount(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("get reserved payment by account: %w", err)
	}

	vs.lock.Lock()
	vs.reservedPayments[accountID] = res
	vs.lock.Unlock()

	return res, nil
}

// reservedPaymentsEqual compares two ReservedPayment objects for equality
func reservedPaymentsEqual(old, new *core.ReservedPayment) bool {
	if old.SymbolsPerSecond != new.SymbolsPerSecond ||
		old.StartTimestamp != new.StartTimestamp ||
		old.EndTimestamp != new.EndTimestamp {
		return false
	}

	// Compare quorum arrays
	if len(old.QuorumNumbers) != len(new.QuorumNumbers) {
		return false
	}
	for i, q := range old.QuorumNumbers {
		if q != new.QuorumNumbers[i] {
			return false
		}
	}
	return true
}
