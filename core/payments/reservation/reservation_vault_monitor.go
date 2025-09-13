package reservation

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	bindings "github.com/Layr-Labs/eigenda/contracts/bindings/v2/PaymentVault"
	"github.com/Layr-Labs/eigenda/core/payments"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"golang.org/x/sync/errgroup"
)

// Checks for updates to the PaymentVault contract, and updates ledgers with the new state
type ReservationVaultMonitor struct {
	logger logging.Logger
	// fetches data from the PaymentVault
	paymentVault payments.PaymentVault
	// how frequently to fetch state from the PaymentVault to check for updates
	updateInterval time.Duration
	// maximum number of accounts to fetch in a single RPC call (0 = unlimited batch size)
	rpcBatchSize uint32
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
	rpcBatchSize uint32,
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
		rpcBatchSize:        rpcBatchSize,
		getAccountsToUpdate: getAccountsToUpdate,
		updateReservation:   updateReservation,
	}

	go monitor.runUpdateLoop(ctx)
	return monitor, nil
}

// Refreshes reservation ledgers with the latest state from the PaymentVault
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

	reservationsMap, err := vm.fetchReservations(ctxWithTimeout, accountIDs)
	if err != nil {
		return fmt.Errorf("fetch reservations: %w", err)
	}

	for accountID, newReservationData := range reservationsMap {
		if newReservationData == nil {
			err := vm.updateReservation(accountID, nil)
			if err != nil {
				vm.logger.Errorf("update nil reservation for account %v failed: %v", accountID.Hex(), err)
			}
			continue
		}

		newReservation, err := FromContractStruct(newReservationData)
		if err != nil {
			vm.logger.Errorf("reservation from contract struct for account %v failed: %v", accountID.Hex(), err)
			continue
		}

		err = vm.updateReservation(accountID, newReservation)
		if err != nil {
			vm.logger.Errorf("update reservation for account %v failed: %v", accountID.Hex(), err)
		}
	}

	return nil
}

// Fetches reservations from the PaymentVault. If number of accountIDs exceeds configured rpcBatchSize, multiple RPC
// calls will be made in parallel to fetch all reservation data. If rpcBatchSize is configured to be 0, all data
// will be fetched in a single call, no matter how many accounts are passed in.
func (vm *ReservationVaultMonitor) fetchReservations(
	ctx context.Context,
	accountIDs []gethcommon.Address,
) (map[gethcommon.Address]*bindings.IPaymentVaultReservation, error) {
	// Split accounts into accountBatches to avoid RPC size limits
	var accountBatches [][]gethcommon.Address

	// Special case: 0 means unlimited batch size, i.e. all accounts are included in a single batch
	if vm.rpcBatchSize == 0 {
		accountBatches = [][]gethcommon.Address{accountIDs}
	} else {
		// Create batches of the specified size
		for i := 0; i < len(accountIDs); i += int(vm.rpcBatchSize) {
			end := min(i+int(vm.rpcBatchSize), len(accountIDs))
			accountBatches = append(accountBatches, accountIDs[i:end])
		}
	}

	results := make(map[gethcommon.Address]*bindings.IPaymentVaultReservation, len(accountIDs))
	var resultsMutex sync.Mutex

	errorGroup, groupCtx := errgroup.WithContext(ctx)

	for batchIndex, batchAccounts := range accountBatches {
		errorGroup.Go(func() error {
			newReservations, err := vm.paymentVault.GetReservations(groupCtx, batchAccounts)
			if err != nil {
				return fmt.Errorf("get reservations for batch %d: %w", batchIndex, err)
			}

			if len(newReservations) != len(batchAccounts) {
				// this shouldn't be possible
				return fmt.Errorf(
					"reservation count mismatch in batch %d: got %d reservations for %d accounts",
					batchIndex, len(newReservations), len(batchAccounts))
			}

			resultsMutex.Lock()
			defer resultsMutex.Unlock()
			// Store results in the map
			for i, accountID := range batchAccounts {
				results[accountID] = newReservations[i]
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
