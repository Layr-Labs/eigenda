package ondemand

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/core/payments"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

type OnDemandVaultMonitor struct {
	logger         logging.Logger
	paymentVault   payments.PaymentVault
	totalDeposits  map[gethcommon.Address]*big.Int
	lock           sync.Mutex
	updateCallback onDemandPaymentUpdateCallback

	updateInterval time.Duration
	cancelFunc     context.CancelFunc
}

// onDemandPaymentUpdateCallback is called when payment vault state detects changes in on-demand payments
// accountID is the account that had a payment update
// newTotalDeposit is the new total deposit amount for that account
type onDemandPaymentUpdateCallback func(accountID gethcommon.Address, newTotalDeposit *big.Int) error

// NewOnDemandVaultMonitor creates a new OnDemandVaultMonitor and immediately starts background monitoring
func NewOnDemandVaultMonitor(
	ctx context.Context,
	logger logging.Logger,
	paymentVault payments.PaymentVault,
	updateCallback onDemandPaymentUpdateCallback,
	updateInterval time.Duration,
) *OnDemandVaultMonitor {
	ctxWithCancel, cancel := context.WithCancel(ctx)

	monitor := &OnDemandVaultMonitor{
		logger:         logger,
		paymentVault:   paymentVault,
		totalDeposits:  make(map[gethcommon.Address]*big.Int),
		updateCallback: updateCallback,
		updateInterval: updateInterval,
		cancelFunc:     cancel,
	}

	go monitor.runUpdateLoop(ctxWithCancel)
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

	newDeposits, err := vs.paymentVault.GetTotalDeposits(ctx, accountIDs)
	if err != nil {
		return fmt.Errorf("get total deposits: %w", err)
	}

	if len(newDeposits) != len(accountIDs) {
		return fmt.Errorf("deposit count mismatch: got %d deposits for %d accounts", len(newDeposits), len(accountIDs))
	}

	for i, newDeposit := range newDeposits {
		accountID := accountIDs[i]
		oldDeposit := vs.totalDeposits[accountID]

		if oldDeposit.Cmp(newDeposit) != 0 {
			err := vs.updateCallback(accountID, newDeposit)
			if err != nil {
				vs.logger.Error("update callback failed", "error", err, "accountID", accountID.Hex())
				continue
			}
		}

		vs.totalDeposits[accountID] = newDeposit
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

// Stop stops the background update thread
func (vs *OnDemandVaultMonitor) Stop() {
	vs.cancelFunc()
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
