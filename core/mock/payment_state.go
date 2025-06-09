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

func (m *MockOnchainPaymentState) RefreshOnchainPaymentState(ctx context.Context) (*meterer.PaymentVaultParams, error) {
	args := m.Called(ctx)
	var value *meterer.PaymentVaultParams
	if args.Get(0) != nil {
		value = args.Get(0).(*meterer.PaymentVaultParams)
	}
	return value, args.Error(1)
}

func (m *MockOnchainPaymentState) GetReservedPaymentByAccount(ctx context.Context, accountID gethcommon.Address) (map[core.QuorumID]*core.ReservedPayment, error) {
	args := m.Called(ctx, accountID)
	var value map[core.QuorumID]*core.ReservedPayment
	if args.Get(0) != nil {
		value = args.Get(0).(map[core.QuorumID]*core.ReservedPayment)
	}
	return value, args.Error(1)
}

func (m *MockOnchainPaymentState) GetReservedPaymentByAccountAndQuorum(ctx context.Context, accountID gethcommon.Address, quorumId uint8) (*core.ReservedPayment, error) {
	args := m.Called(ctx, accountID, quorumId)
	var value *core.ReservedPayment
	if args.Get(0) != nil {
		value = args.Get(0).(*core.ReservedPayment)
	}
	return value, args.Error(1)
}

func (m *MockOnchainPaymentState) GetReservedPaymentByAccountAndQuorums(ctx context.Context, accountID gethcommon.Address, quorumNumbers []core.QuorumID) (map[core.QuorumID]*core.ReservedPayment, error) {
	args := m.Called(ctx, accountID, quorumNumbers)
	var value map[core.QuorumID]*core.ReservedPayment
	if fn, ok := args.Get(0).(func(context.Context, gethcommon.Address, []core.QuorumID) map[core.QuorumID]*core.ReservedPayment); ok {
		value = fn(ctx, accountID, quorumNumbers)
	} else if args.Get(0) != nil {
		value = args.Get(0).(map[core.QuorumID]*core.ReservedPayment)
	}
	var err error
	if len(args) > 1 {
		err = args.Error(1)
	}
	return value, err
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

func (m *MockOnchainPaymentState) GetOnDemandGlobalSymbolsPerSecond(quorumID core.QuorumID) uint64 {
	args := m.Called(quorumID)
	return args.Get(0).(uint64)
}

func (m *MockOnchainPaymentState) GetOnDemandGlobalRatePeriodInterval(quorumID core.QuorumID) uint64 {
	args := m.Called(quorumID)
	return args.Get(0).(uint64)
}

func (m *MockOnchainPaymentState) GetMinNumSymbols(quorumID core.QuorumID) uint64 {
	args := m.Called(quorumID)
	return args.Get(0).(uint64)
}

func (m *MockOnchainPaymentState) GetPricePerSymbol(quorumID core.QuorumID) uint64 {
	args := m.Called(quorumID)
	return args.Get(0).(uint64)
}

func (m *MockOnchainPaymentState) GetReservationWindow(quorumID core.QuorumID) uint64 {
	args := m.Called(quorumID)
	return args.Get(0).(uint64)
}

func (m *MockOnchainPaymentState) GetQuorumNumbers(ctx context.Context) ([]uint8, error) {
	args := m.Called()
	var value []uint8
	if args.Get(0) != nil {
		value = args.Get(0).([]uint8)
	}
	return value, args.Error(1)
}

func (m *MockOnchainPaymentState) GetQuorumPaymentConfig(quorumID core.QuorumID) (*core.PaymentQuorumConfig, error) {
	args := m.Called(quorumID)
	return args.Get(0).(*core.PaymentQuorumConfig), args.Error(1)
}

func (m *MockOnchainPaymentState) GetQuorumProtocolConfig(quorumID core.QuorumID) (*core.PaymentQuorumProtocolConfig, error) {
	args := m.Called(quorumID)
	return args.Get(0).(*core.PaymentQuorumProtocolConfig), args.Error(1)
}
