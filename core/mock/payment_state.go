package mock

import (
	"context"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/meterer"
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

func (m *MockOnchainPaymentState) RefreshOnchainPaymentState(ctx context.Context, tx *eth.Reader) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockOnchainPaymentState) GetActiveReservationByAccount(ctx context.Context, accountID string) (core.ActiveReservation, error) {
	args := m.Called(ctx, accountID)
	var value core.ActiveReservation
	if args.Get(0) != nil {
		value = args.Get(0).(core.ActiveReservation)
	}
	return value, args.Error(1)
}

func (m *MockOnchainPaymentState) GetOnDemandPaymentByAccount(ctx context.Context, accountID string) (core.OnDemandPayment, error) {
	args := m.Called(ctx, accountID)
	var value core.OnDemandPayment
	if args.Get(0) != nil {
		value = args.Get(0).(core.OnDemandPayment)
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
