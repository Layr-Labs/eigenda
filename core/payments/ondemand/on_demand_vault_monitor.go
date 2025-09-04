package ondemand

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/Layr-Labs/eigenda/core/payments"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// Checks for updates to the PaymentVault contract, and updates ledgers with the new state
type OnDemandVaultMonitor struct {
	logger logging.Logger
	// fetches data from the PaymentVault
	paymentVault payments.PaymentVault
	// how frequently to fetch state from the PaymentVault to check for updates
	updateInterval time.Duration
	// function to get accounts that need to be updated
	getAccountsToUpdate func() []gethcommon.Address
	// function to update the total deposit for an account
	updateTotalDeposit func(accountID gethcommon.Address, newTotalDeposit *big.Int) error
}

// Creates a new OnDemandVaultMonitor and starts a routine to periodically check for updates
func NewOnDemandVaultMonitor(
	ctx context.Context,
	logger logging.Logger,
	paymentVault payments.PaymentVault,
	updateInterval time.Duration,
	getAccountsToUpdate func() []gethcommon.Address,
	updateTotalDeposit func(accountID gethcommon.Address, newTotalDeposit *big.Int) error,
) (*OnDemandVaultMonitor, error) {
	if updateInterval <= 0 {
		return nil, errors.New("updateInterval must be > 0")
	}

	monitor := &OnDemandVaultMonitor{
		logger:              logger,
		paymentVault:        paymentVault,
		updateInterval:      updateInterval,
		getAccountsToUpdate: getAccountsToUpdate,
		updateTotalDeposit:  updateTotalDeposit,
	}

	go monitor.runUpdateLoop(ctx)
	return monitor, nil
}

// Fetches the latest state from the PaymentVault, and updates the ledgers with it
func (vm *OnDemandVaultMonitor) refreshTotalDeposits(ctx context.Context) error {
	accountIDs := vm.getAccountsToUpdate()
	if len(accountIDs) == 0 {
		return nil
	}

	// Add timeout to prevent hanging if the RPC node is unresponsive.
	// This timeout is higher than it needs to be, but at least if we are unable to access
	// the eth node, then we will time out before the next refresh try.
	ctxWithTimeout, cancel := context.WithTimeout(ctx, vm.updateInterval)
	defer cancel()

	newDeposits, err := vm.paymentVault.GetTotalDeposits(ctxWithTimeout, accountIDs)
	if err != nil {
		return fmt.Errorf("get total deposits: %w", err)
	}

	if len(newDeposits) != len(accountIDs) {
		// this shouldn't be possible
		return fmt.Errorf("deposit count mismatch: got %d deposits for %d accounts", len(newDeposits), len(accountIDs))
	}

	// This loop could theoretically be parallelized, but none of the current use cases (either a cache, or an
	// individual account) require it. Therefore, a loop is used for simplicity.
	for i, newDeposit := range newDeposits {
		accountID := accountIDs[i]
		err := vm.updateTotalDeposit(accountID, newDeposit)
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

	vm.logger.Debugf("Starting OnDemandPaymentVault background update thread with updateInterval %v", vm.updateInterval)

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
