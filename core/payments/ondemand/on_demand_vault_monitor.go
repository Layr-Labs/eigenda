package ondemand

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/core/payments/vault"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

type OnDemandVaultMonitor struct {
	logger         logging.Logger
	paymentVault   *vault.PaymentVault
	totalDeposits  map[gethcommon.Address]*big.Int
	lock           sync.Mutex
	updateCallback onDemandPaymentUpdateCallback

	// Background update configuration
	updateInterval time.Duration
	cancelFunc     context.CancelFunc
}

// onDemandPaymentUpdateCallback is called when payment vault state detects changes in on-demand payments
// accountID is the account that had a payment update
// newTotalDeposit is the new total deposit amount for that account
type onDemandPaymentUpdateCallback func(accountID gethcommon.Address, newTotalDeposit *big.Int) error

// NewOnDemandVaultMonitor creates a new OnDemandVaultMonitor and immediately starts background monitoring
func NewOnDemandVaultMonitor(ctx context.Context, logger logging.Logger, paymentVault *vault.PaymentVault, updateCallback onDemandPaymentUpdateCallback, updateInterval time.Duration) *OnDemandVaultMonitor {
	monitor := &OnDemandVaultMonitor{
		logger:         logger,
		paymentVault:   paymentVault,
		totalDeposits:  make(map[gethcommon.Address]*big.Int),
		updateCallback: updateCallback,
		updateInterval: updateInterval,
	}
	
	// Auto-start the background monitoring thread
	monitor.Start(ctx)
	
	return monitor
}

func (vs *OnDemandVaultMonitor) RefreshTotalDeposits(ctx context.Context) error {
	vs.lock.Lock()
	defer vs.lock.Unlock()

	if len(vs.totalDeposits) == 0 {
		// no accounts are being tracked, so there isn't anything to refresh
		return nil
	}

	accountIDs := make([]gethcommon.Address, 0, len(vs.totalDeposits))
	for accountID := range vs.totalDeposits {
		accountIDs = append(accountIDs, accountID)
	}

	newOnDemandPayments, err := vs.paymentVault.GetTotalDeposits(ctx, accountIDs)
	if err != nil {
		return fmt.Errorf("get on-demand payments: %w", err)
	}

	err = vs.processUpdates(newOnDemandPayments)
	if err != nil {
		return fmt.Errorf("process payment updates: %w", err)
	}

	vs.totalDeposits = newOnDemandPayments
	return nil
}

// processUpdates compares old and new payments to detect changes
// and invokes the callback for accounts with payment changes
func (vs *OnDemandVaultMonitor) processUpdates(
	updates map[gethcommon.Address]*big.Int,
) error {
	for accountID, newValue := range updates {
		oldValue, exists := vs.totalDeposits[accountID]
		if !exists {
			vs.logger.Errorf("received an unrequested value for account %v. This shouldn't be possible",
				accountID.Hex())
			continue
		}

		if oldValue.Cmp(newValue) == 0 {
			// no update necessary, since new value is the same as the old
			continue
		}

		err := vs.updateCallback(accountID, newValue)
		if err != nil {
			return fmt.Errorf("update callback for account %v: %w", accountID.Hex(), err)
		}
	}

	return nil
}

func (vs *OnDemandVaultMonitor) GetTotalDeposit(ctx context.Context, accountID gethcommon.Address) (*big.Int, error) {
	onDemandPayment, err := vs.paymentVault.GetTotalDeposit(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("get total deposit: %w", err)
	}

	vs.lock.Lock()
	defer vs.lock.Unlock()

	vs.totalDeposits[accountID] = onDemandPayment

	return onDemandPayment, nil
}

// Start starts the background update thread
func (vs *OnDemandVaultMonitor) Start(ctx context.Context) {
	ctxWithCancel, cancel := context.WithCancel(ctx)
	vs.cancelFunc = cancel

	go vs.runUpdateLoop(ctxWithCancel)
}

// Stop stops the background update thread
func (vs *OnDemandVaultMonitor) Stop() {
	if vs.cancelFunc != nil {
		vs.cancelFunc()
	}
}

// runUpdateLoop runs the background update loop to periodically consume updates made to the PaymentVault
func (vs *OnDemandVaultMonitor) runUpdateLoop(ctx context.Context) {
	ticker := time.NewTicker(vs.updateInterval)
	defer ticker.Stop()

	vs.logger.Info("Starting OnDemandPaymentVault background update thread", "updateInterval", vs.updateInterval)

	for {
		select {
		case <-ticker.C:
			if err := vs.RefreshTotalDeposits(ctx); err != nil {
				vs.logger.Error("perform on-demand payment updates", "error", err)
			}
		case <-ctx.Done():
			vs.logger.Info("OnDemandPaymentVault background update thread stopped")
			return
		}
	}
}
