package reservation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/core/payments"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// Checks for updates to the PaymentVault contract, and updates ledgers with the new state
type ReservationVaultMonitor struct {
	logger logging.Logger
	// fetches data from the PaymentVault
	paymentVault payments.PaymentVault
	// how frequently to fetch state from the PaymentVault to check for updates
	updateInterval time.Duration
	// function to get accounts that need to be updated
	getAccountsToUpdate func() []gethcommon.Address
	// function to update the reservation for an account
	updateReservation func(accountID gethcommon.Address, newReservation *Reservation) error
}

// Creates a new ReservationVaultMonitor and starts a routine to periodically check for updates
func NewReservationVaultMonitor(
	ctx context.Context,
	logger logging.Logger,
	paymentVault payments.PaymentVault,
	updateInterval time.Duration,
	getAccountsToUpdate func() []gethcommon.Address,
	updateReservation func(accountID gethcommon.Address, newReservation *Reservation) error,
) (*ReservationVaultMonitor, error) {
	if updateInterval <= 0 {
		return nil, errors.New("updateInterval must be > 0")
	}

	monitor := &ReservationVaultMonitor{
		logger:              logger,
		paymentVault:        paymentVault,
		updateInterval:      updateInterval,
		getAccountsToUpdate: getAccountsToUpdate,
		updateReservation:   updateReservation,
	}

	go monitor.runUpdateLoop(ctx)
	return monitor, nil
}

// Fetches the latest state from the PaymentVault, and updates the ledgers with it
func (vm *ReservationVaultMonitor) refreshReservations(ctx context.Context) error {
	accountIDs := vm.getAccountsToUpdate()
	if len(accountIDs) == 0 {
		return nil
	}

	// Add timeout to prevent hanging if the RPC node is unresponsive.
	// This timeout is higher than it needs to be, but at least if we are unable to access
	// the eth node, then we will time out before the next refresh try.
	ctxWithTimeout, cancel := context.WithTimeout(ctx, vm.updateInterval)
	defer cancel()

	newReservations, err := vm.paymentVault.GetReservations(ctxWithTimeout, accountIDs)
	if err != nil {
		return fmt.Errorf("get reservations: %w", err)
	}

	if len(newReservations) != len(accountIDs) {
		// this shouldn't be possible
		return fmt.Errorf(
			"reservation count mismatch: got %d reservations for %d accounts", len(newReservations), len(accountIDs))
	}

	for i, newReservationData := range newReservations {
		accountID := accountIDs[i]
		if newReservationData == nil {
			err := vm.updateReservation(accountID, nil)
			if err != nil {
				vm.logger.Errorf("update nil reservation for account %v failed: %v", accountID.Hex(), err)
			}
			continue
		}

		newReservation, err := NewReservationFromBindings(newReservationData)
		if err != nil {
			vm.logger.Errorf("convert reservation for account %v failed: %v", accountID.Hex(), err)
			continue
		}

		err = vm.updateReservation(accountID, newReservation)
		if err != nil {
			vm.logger.Errorf("update reservation for account %v failed: %v", accountID.Hex(), err)
		}
	}

	return nil
}

// Runs the background update loop to periodically consume updates made to the PaymentVault
func (vm *ReservationVaultMonitor) runUpdateLoop(ctx context.Context) {
	ticker := time.NewTicker(vm.updateInterval)
	defer ticker.Stop()

	vm.logger.Debugf(
		"Starting ReservationVaultMonitor background update thread with updateInterval %v", vm.updateInterval)

	for {
		select {
		case <-ticker.C:
			if err := vm.refreshReservations(ctx); err != nil {
				vm.logger.Errorf("refresh reservations: %v", err)
			}
		case <-ctx.Done():
			vm.logger.Debug("ReservationVaultMonitor background update thread stopped")
			return
		}
	}
}
