package paymentvault

import (
	"context"
	"fmt"
	"sync"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// OnDemandPaymentVaultState manages on-demand payment state from the payment vault
type OnDemandPaymentVaultState struct {
	logger           logging.Logger
	ethReader        *eth.Reader
	onDemandPayments map[gethcommon.Address]*core.OnDemandPayment
	lock             sync.RWMutex
}

// NewOnDemandPaymentVaultState creates a new OnDemandPaymentVaultState
func NewOnDemandPaymentVaultState(logger logging.Logger, ethReader *eth.Reader) *OnDemandPaymentVaultState {
	return &OnDemandPaymentVaultState{
		logger:           logger,
		ethReader:        ethReader,
		onDemandPayments: make(map[gethcommon.Address]*core.OnDemandPayment),
	}
}

// RefreshOnDemandPayments updates the cached on-demand payments and returns a list of detected changes
// TODO: Replace periodic polling with event subscription from PaymentVault contract
// Should subscribe to OnDemandDepositReceived events instead of polling GetOnDemandPayments
func (vs *OnDemandPaymentVaultState) RefreshOnDemandPayments(ctx context.Context) ([]TotalDepositUpdate, error) {
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

	// Detect changes and create updates
	var updates []TotalDepositUpdate
	for accountID, newPayment := range newOnDemandPayments {
		oldPayment, exists := vs.onDemandPayments[accountID]
		if !exists || oldPayment.CumulativePayment.Cmp(newPayment.CumulativePayment) != 0 {
			update, err := NewTotalDepositUpdate(accountID, newPayment.CumulativePayment)
			if err != nil {
				vs.logger.Error("Failed to create total deposit update", "error", err, "accountID", accountID.Hex())
				continue
			}
			updates = append(updates, *update)
		}
	}

	vs.onDemandPayments = newOnDemandPayments
	return updates, nil
}

// GetOnDemandPaymentByAccount retrieves on-demand payment info for a specific account
func (vs *OnDemandPaymentVaultState) GetOnDemandPaymentByAccount(
	ctx context.Context,
	accountID gethcommon.Address,
) (*core.OnDemandPayment, error) {
	vs.lock.RLock()
	if payment, ok := vs.onDemandPayments[accountID]; ok {
		vs.lock.RUnlock()
		return payment, nil
	}
	vs.lock.RUnlock()

	// pulls the chain state
	res, err := vs.ethReader.GetOnDemandPaymentByAccount(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("get on-demand payment by account: %w", err)
	}

	vs.lock.Lock()
	vs.onDemandPayments[accountID] = res
	vs.lock.Unlock()

	return res, nil
}
