package mock

import (
	"context"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
)

type MockOnchainPaymentState struct {
	mock.Mock
}

var _ meterer.OnchainPayment = (*MockOnchainPaymentState)(nil)

func (m *MockOnchainPaymentState) GetCurrentBlockNumber(ctx context.Context) (uint32, error) {
	args := m.Called()
	var value uint32
	if args.Get(0) != nil {
		value = args.Get(0).(uint32)
	}
	return value, args.Error(1)
}

func (m *MockOnchainPaymentState) RefreshOnchainPaymentState(ctx context.Context) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockOnchainPaymentState) GetReservedPaymentByAccount(ctx context.Context, accountID gethcommon.Address) (*core.ReservedPayment, error) {
	args := m.Called(ctx, accountID)
	var value *core.ReservedPayment
	if args.Get(0) != nil {
		value = args.Get(0).(*core.ReservedPayment)
	}
	return value, args.Error(1)
}

func (m *MockOnchainPaymentState) GetOnDemandPaymentByAccount(ctx context.Context, accountID gethcommon.Address) (*core.OnDemandPayment, error) {
	args := m.Called(ctx, accountID)
	var value *core.OnDemandPayment
	if args.Get(0) != nil {
		value = args.Get(0).(*core.OnDemandPayment)
	}
	return value, args.Error(1)
}

func (m *MockOnchainPaymentState) GetOnDemandQuorumNumbers(ctx context.Context) ([]uint8, error) {
	args := m.Called()
	var value []uint8
	if args.Get(0) != nil {
		value = args.Get(0).([]uint8)
	}
	return value, args.Error(1)
}

func (m *MockOnchainPaymentState) GetGlobalSymbolsPerSecond() uint64 {
	args := m.Called()
	return args.Get(0).(uint64)
}

func (m *MockOnchainPaymentState) GetGlobalRatePeriodInterval() uint32 {
	args := m.Called()
	return args.Get(0).(uint32)
}

func (m *MockOnchainPaymentState) GetMinNumSymbols() uint32 {
	args := m.Called()
	return args.Get(0).(uint32)
}

func (m *MockOnchainPaymentState) GetPricePerSymbol() uint32 {
	args := m.Called()
	return args.Get(0).(uint32)
}

func (m *MockOnchainPaymentState) GetReservationWindow() uint32 {
	args := m.Called()
	return args.Get(0).(uint32)
}
