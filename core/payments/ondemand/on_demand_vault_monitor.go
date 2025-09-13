package ondemand

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/core/payments"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"golang.org/x/sync/errgroup"
)

// Checks for updates to the PaymentVault contract, and updates ledgers with the new state
type OnDemandVaultMonitor struct {
	logger logging.Logger
	// fetches data from the PaymentVault
	paymentVault payments.PaymentVault
	// how frequently to fetch state from the PaymentVault to check for updates
	updateInterval time.Duration
	// maximum number of accounts to fetch in a single RPC call (0 = no batching)
	rpcBatchSize uint32
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
	rpcBatchSize uint32,
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
		rpcBatchSize:        rpcBatchSize,
		getAccountsToUpdate: getAccountsToUpdate,
		updateTotalDeposit:  updateTotalDeposit,
	}

	go monitor.runUpdateLoop(ctx)
	return monitor, nil
}

// Refreshes total deposits with the latest state from the PaymentVault
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

	depositsMap, err := vm.fetchTotalDeposits(ctxWithTimeout, accountIDs)
	if err != nil {
		return fmt.Errorf("fetch total deposits: %w", err)
	}

	for accountID, newDeposit := range depositsMap {
		err := vm.updateTotalDeposit(accountID, newDeposit)
		if err != nil {
			vm.logger.Errorf("update total deposit for account %v failed: %v", accountID.Hex(), err)
		}
	}

	return nil
}

// Fetches total deposits from the PaymentVault. If number of accountIDs exceeds configured rpcBatchSize, multiple RPC
// calls will be made in parallel to fetch all deposit data. If rpcBatchSize is configured to be 0, all data
// will be fetched in a single call, no matter how many accounts are passed in.
func (vm *OnDemandVaultMonitor) fetchTotalDeposits(
	ctx context.Context,
	accountIDs []gethcommon.Address,
) (map[gethcommon.Address]*big.Int, error) {
	// Split accounts into accountBatches to avoid RPC size limits
	var accountBatches [][]gethcommon.Address

	// Special case: 0 means no batching
	if vm.rpcBatchSize == 0 {
		accountBatches = [][]gethcommon.Address{accountIDs}
	} else {
		// Create batches of the specified size
		for i := 0; i < len(accountIDs); i += int(vm.rpcBatchSize) {
			end := min(i+int(vm.rpcBatchSize), len(accountIDs))
			accountBatches = append(accountBatches, accountIDs[i:end])
		}
	}

	results := make(map[gethcommon.Address]*big.Int, len(accountIDs))
	var resultsMutex sync.Mutex

	errorGroup, groupCtx := errgroup.WithContext(ctx)

	for index, batch := range accountBatches {
		// Capture loop variables for goroutine
		batchIndex := index
		batchAccounts := batch

		errorGroup.Go(func() error {
			newDeposits, err := vm.paymentVault.GetTotalDeposits(groupCtx, batchAccounts)
			if err != nil {
				return fmt.Errorf("get total deposits for batch %d: %w", batchIndex, err)
			}

			if len(newDeposits) != len(batchAccounts) {
				// this shouldn't be possible
				return fmt.Errorf(
					"deposit count mismatch in batch %d: got %d deposits for %d accounts",
					batchIndex, len(newDeposits), len(batchAccounts))
			}

			resultsMutex.Lock()
			defer resultsMutex.Unlock()
			// Store results in the map
			for i, accountID := range batchAccounts {
				results[accountID] = newDeposits[i]
			}

			return nil
		})
	}

	if err := errorGroup.Wait(); err != nil {
		return nil, fmt.Errorf("error group wait: %w", err)
	}

	return results, nil
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
