package ondemand

import (
	"context"
	"fmt"
	"sync"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/payments/ondemand"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// OnDemandPaymentVaultState manages on-demand payment state from the payment vault
type onDemandPaymentVaultState struct {
	logger           logging.Logger
	ethReader        *eth.Reader
	onDemandPayments map[gethcommon.Address]*core.OnDemandPayment
	lock             sync.Mutex
}

// NewOnDemandPaymentVaultState creates a new OnDemandPaymentVaultState
func NewOnDemandPaymentVaultState(logger logging.Logger, ethReader *eth.Reader) ondemand.OnDemandPaymentVaultState {
	return &onDemandPaymentVaultState{
		logger:           logger,
		ethReader:        ethReader,
		onDemandPayments: make(map[gethcommon.Address]*core.OnDemandPayment),
	}
}

// RefreshOnDemandPayments updates the cached on-demand payments and returns a list of detected changes
func (vs *onDemandPaymentVaultState) RefreshOnDemandPayments(
	ctx context.Context,
) ([]ondemand.TotalDepositUpdate, error) {
	vs.lock.Lock()
	defer vs.lock.Unlock()

	if len(vs.onDemandPayments) == 0 {
		vs.logger.Info("No on-demand payments to refresh")
		return nil, nil
	}

	accountIDs := make([]gethcommon.Address, 0, len(vs.onDemandPayments))
	for accountID := range vs.onDemandPayments {
		accountIDs = append(accountIDs, accountID)
	}

	newOnDemandPayments, err := vs.ethReader.GetOnDemandPayments(ctx, accountIDs)
	if err != nil {
		return nil, fmt.Errorf("get on-demand payments: %w", err)
	}

	updates := vs.buildTotalDepositUpdates(newOnDemandPayments)

	vs.onDemandPayments = newOnDemandPayments
	return updates, nil
}

// buildTotalDepositUpdates compares old and new payments to detect changes
// and creates TotalDepositUpdate objects for accounts with payment changes
func (vs *onDemandPaymentVaultState) buildTotalDepositUpdates(
	newPayments map[gethcommon.Address]*core.OnDemandPayment,
) []ondemand.TotalDepositUpdate {
	var updates []ondemand.TotalDepositUpdate
	for accountID, newPayment := range newPayments {
		oldPayment, exists := vs.onDemandPayments[accountID]
		if !exists {
			vs.logger.Errorf("received an unrequested value for account %v. This shouldn't be possible",
				accountID.Hex())
			continue
		}

		if oldPayment.CumulativePayment.Cmp(newPayment.CumulativePayment) == 0 {
			// no update necessary, since new value is the same as the old
			continue
		}

		update, err := ondemand.NewTotalDepositUpdate(accountID, newPayment.CumulativePayment)
		if err != nil {
			vs.logger.Error("new total deposit update", "error", err, "accountID", accountID.Hex())
			continue
		}
		updates = append(updates, *update)
	}

	return updates
}

// GetOnDemandPaymentByAccount retrieves on-demand payment info for a specific account
func (vs *onDemandPaymentVaultState) GetOnDemandPaymentByAccount(
	ctx context.Context,
	accountID gethcommon.Address,
) (*core.OnDemandPayment, error) {
	onDemandPayment, err := vs.ethReader.GetOnDemandPaymentByAccount(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("get on-demand payment by account: %w", err)
	}

	vs.lock.Lock()
	defer vs.lock.Unlock()

	vs.onDemandPayments[accountID] = onDemandPayment

	return onDemandPayment, nil
}
