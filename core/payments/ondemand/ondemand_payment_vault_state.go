package ondemand

import (
	"context"

	"github.com/Layr-Labs/eigenda/core"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// OnDemandPaymentVaultStateInterface defines the interface for on-demand payment vault state operations
type OnDemandPaymentVaultState interface {
	RefreshOnDemandPayments(ctx context.Context) ([]TotalDepositUpdate, error)
	GetOnDemandPaymentByAccount(ctx context.Context, accountID gethcommon.Address) (*core.OnDemandPayment, error)
}
