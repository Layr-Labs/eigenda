package ondemand

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/core/payments"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// Checks for updates to the PaymentVault contract, and updates ledgers with the new state
type OnDemandVaultMonitor struct {
	logger logging.Logger
	// fetches data from the PaymentVault
	paymentVault payments.PaymentVault
	// the ledgers that need to be updated
	ledgers UpdatableOnDemandLedgers
	// how frequently to fetch state from the PaymentVault to check for updates
	updateInterval time.Duration
	// cancels the periodic update routine
	cancelFunc context.CancelFunc
}

// Creates a new OnDemandVaultMonitor and starts a routine to periodically check for updates
func NewOnDemandVaultMonitor(
	ctx context.Context,
	logger logging.Logger,
	paymentVault payments.PaymentVault,
	ledgers UpdatableOnDemandLedgers,
	updateInterval time.Duration,
) *OnDemandVaultMonitor {
	ctxWithCancel, cancel := context.WithCancel(ctx)

	monitor := &OnDemandVaultMonitor{
		logger:         logger,
		paymentVault:   paymentVault,
		ledgers:        ledgers,
		updateInterval: updateInterval,
		cancelFunc:     cancel,
	}

	go monitor.runUpdateLoop(ctxWithCancel)
	return monitor
}

func (vm *OnDemandVaultMonitor) Stop() {
	vm.cancelFunc()
}

// Fetches the latest state from the PaymentVault, and updates the ledgers with it
func (vm *OnDemandVaultMonitor) refreshTotalDeposits(ctx context.Context) error {
	accountIDs := vm.ledgers.GetAccountsToUpdate()
	if len(accountIDs) == 0 {
		return nil
	}

	newDeposits, err := vm.paymentVault.GetTotalDeposits(ctx, accountIDs)
	if err != nil {
		return fmt.Errorf("get total deposits: %w", err)
	}

	if len(newDeposits) != len(accountIDs) {
		// this shouldn't be possible
		return fmt.Errorf("deposit count mismatch: got %d deposits for %d accounts", len(newDeposits), len(accountIDs))
	}

	for i, newDeposit := range newDeposits {
		accountID := accountIDs[i]
		err := vm.ledgers.UpdateTotalDeposit(accountID, newDeposit)
		if err != nil {
			vm.logger.Errorf("update total deposit for account %v failed: %v", accountID.Hex(), err)
		}
	}

	return nil
}

// Runs the background update loop to periodically consume updates made to the PaymentVault
func (vm *OnDemandVaultMonitor) runUpdateLoop(ctx context.Context) {
	ticker := time.NewTicker(vm.updateInterval)
	defer ticker.Stop()

	vm.logger.Infof("Starting OnDemandPaymentVault background update thread with updateInterval %d", vm.updateInterval)

	for {
		select {
		case <-ticker.C:
			if err := vm.refreshTotalDeposits(ctx); err != nil {
				vm.logger.Errorf("refresh total deposits: %v", err)
			}
		case <-ctx.Done():
			vm.logger.Info("OnDemandPaymentVault background update thread stopped")
			return
		}
	}
}
