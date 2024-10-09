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

func (m *MockOnchainPaymentState) CurrentOnchainPaymentState(ctx context.Context, tx *eth.Transactor) (meterer.OnchainPaymentState, error) {
	args := m.Called()
	var value meterer.OnchainPaymentState
	if args.Get(0) != nil {
		value = args.Get(0).(meterer.OnchainPaymentState)
	}
	return value, args.Error(1)
}

func (m *MockOnchainPaymentState) GetActiveReservations(ctx context.Context, blockNumber uint32) (map[string]core.ActiveReservation, error) {
	args := m.Called()
	var value map[string]core.ActiveReservation
	if args.Get(0) != nil {
		value = args.Get(0).(map[string]core.ActiveReservation)
	}
	return value, args.Error(1)
}

func (m *MockOnchainPaymentState) GetActiveReservationsByAccount(ctx context.Context, blockNumber uint32, accountID string) (core.ActiveReservation, error) {
	args := m.Called()
	var value core.ActiveReservation
	if args.Get(0) != nil {
		value = args.Get(0).(core.ActiveReservation)
	}
	return value, args.Error(1)
}

func (m *MockOnchainPaymentState) GetOnDemandPayments(ctx context.Context, blockNumber uint32) (map[string]core.OnDemandPayment, error) {
	args := m.Called()
	var value map[string]core.OnDemandPayment
	if args.Get(0) != nil {
		value = args.Get(0).(map[string]core.OnDemandPayment)
	}
	return value, args.Error(1)
}

func (m *MockOnchainPaymentState) GetOnDemandPaymentByAccount(ctx context.Context, blockNumber uint32, accountID string) (core.OnDemandPayment, error) {
	args := m.Called()
	var value core.OnDemandPayment
	if args.Get(0) != nil {
		value = args.Get(0).(core.OnDemandPayment)
	}
	return value, args.Error(1)
}
