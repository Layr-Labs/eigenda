package mock

import (
	"context"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/core/payment"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
)

type MockOnchainPaymentState struct {
	mock.Mock
}

var _ meterer.OnchainPayment = (*MockOnchainPaymentState)(nil)

func (m *MockOnchainPaymentState) RefreshOnchainPaymentState(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockOnchainPaymentState) GetReservedPaymentByAccountAndQuorums(ctx context.Context, accountID gethcommon.Address, quorumNumbers []core.QuorumID) (map[core.QuorumID]*payment.ReservedPayment, error) {
	args := m.Called(ctx, accountID, quorumNumbers)
	var value map[core.QuorumID]*payment.ReservedPayment
	if fn, ok := args.Get(0).(func(context.Context, gethcommon.Address, []core.QuorumID) map[core.QuorumID]*payment.ReservedPayment); ok {
		value = fn(ctx, accountID, quorumNumbers)
	} else if args.Get(0) != nil {
		value = args.Get(0).(map[core.QuorumID]*payment.ReservedPayment)
	}
	var err error
	if len(args) > 1 {
		err = args.Error(1)
	}
	return value, err
}

func (m *MockOnchainPaymentState) GetOnDemandPaymentByAccount(ctx context.Context, accountID gethcommon.Address) (*payment.OnDemandPayment, error) {
	args := m.Called(ctx, accountID)
	var value *payment.OnDemandPayment
	if args.Get(0) != nil {
		value = args.Get(0).(*payment.OnDemandPayment)
	}
	return value, args.Error(1)
}

func (m *MockOnchainPaymentState) GetQuorumNumbers(ctx context.Context) ([]core.QuorumID, error) {
	args := m.Called(ctx)
	var value []core.QuorumID
	if args.Get(0) != nil {
		value = args.Get(0).([]core.QuorumID)
	}
	return value, args.Error(1)
}

func (m *MockOnchainPaymentState) GetPaymentGlobalParams() (*payment.PaymentVaultParams, error) {
	args := m.Called()
	var value *payment.PaymentVaultParams
	if args.Get(0) != nil {
		value = args.Get(0).(*payment.PaymentVaultParams)
	}
	return value, args.Error(1)
}
