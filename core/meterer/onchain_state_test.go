package meterer_test

import (
	"context"
	"testing"

	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOnchainPaymentState : TO BE REPLACED WITH ACTUAL IMPLEMENTATION
type MockOnchainPaymentState struct {
	mock.Mock
}

func (m *MockOnchainPaymentState) GetCurrentBlockNumber() (uint, error) {
	args := m.Called()
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockOnchainPaymentState) GetActiveReservations(ctx context.Context, blockNumber uint) (*meterer.ActiveReservations, error) {
	args := m.Called(ctx, blockNumber)
	return args.Get(0).(*meterer.ActiveReservations), args.Error(1)
}

func (m *MockOnchainPaymentState) GetActiveReservationByAccount(ctx context.Context, blockNumber uint, accountID string) (*meterer.ActiveReservation, error) {
	args := m.Called(ctx, blockNumber, accountID)
	return args.Get(0).(*meterer.ActiveReservation), args.Error(1)
}

func (m *MockOnchainPaymentState) GetOnDemandPayments(ctx context.Context, blockNumber uint) (*meterer.OnDemandPayments, error) {
	args := m.Called(ctx, blockNumber)
	return args.Get(0).(*meterer.OnDemandPayments), args.Error(1)
}

func (m *MockOnchainPaymentState) GetOnDemandPaymentByAccount(ctx context.Context, blockNumber uint, accountID string) (*meterer.OnDemandPayment, error) {
	args := m.Called(ctx, blockNumber, accountID)
	return args.Get(0).(*meterer.OnDemandPayment), args.Error(1)
}

func (m *MockOnchainPaymentState) GetIndexedActiveReservations(ctx context.Context, blockNumber uint) (*meterer.ActiveReservations, error) {
	args := m.Called(ctx, blockNumber)
	return args.Get(0).(*meterer.ActiveReservations), args.Error(1)
}

func (m *MockOnchainPaymentState) GetIndexedActiveReservationByAccount(ctx context.Context, blockNumber uint, accountID string) (*meterer.ActiveReservation, error) {
	args := m.Called(ctx, blockNumber, accountID)
	return args.Get(0).(*meterer.ActiveReservation), args.Error(1)
}

func (m *MockOnchainPaymentState) GetIndexedOnDemandPayments(ctx context.Context, blockNumber uint) (*meterer.OnDemandPayments, error) {
	args := m.Called(ctx, blockNumber)
	return args.Get(0).(*meterer.OnDemandPayments), args.Error(1)
}

func (m *MockOnchainPaymentState) GetIndexedOnDemandPaymentByAccount(ctx context.Context, blockNumber uint, accountID string) (*meterer.OnDemandPayment, error) {
	args := m.Called(ctx, blockNumber, accountID)
	return args.Get(0).(*meterer.OnDemandPayment), args.Error(1)
}

func (m *MockOnchainPaymentState) Start(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestGetCurrentBlockNumber(t *testing.T) {
	mockState := new(MockOnchainPaymentState)
	mockState.On("GetCurrentBlockNumber").Return(uint(1000), nil)

	blockNumber, err := mockState.GetCurrentBlockNumber()
	assert.NoError(t, err)
	assert.Equal(t, uint(1000), blockNumber)
}

func TestGetActiveReservations(t *testing.T) {
	mockState := new(MockOnchainPaymentState)
	ctx := context.Background()
	expectedReservations := &meterer.ActiveReservations{
		Reservations: map[string]*meterer.ActiveReservation{
			"account1": {
				DataRate:       100,
				StartTimestamp: 1000,
				EndTimestamp:   2000,
				QuorumSplit:    []byte{50, 50},
			},
		},
	}
	mockState.On("GetActiveReservations", ctx, uint(1000)).Return(expectedReservations, nil)

	reservations, err := mockState.GetActiveReservations(ctx, 1000)
	assert.NoError(t, err)
	assert.Equal(t, expectedReservations, reservations)
}

func TestGetOnDemandPaymentByAccount(t *testing.T) {
	mockState := new(MockOnchainPaymentState)
	ctx := context.Background()
	accountID := "account1"
	expectedPayment := &meterer.OnDemandPayment{
		CumulativePayment: meterer.TokenAmount(1000000),
	}
	mockState.On("GetOnDemandPaymentByAccount", ctx, uint(1000), accountID).Return(expectedPayment, nil)

	payment, err := mockState.GetOnDemandPaymentByAccount(ctx, 1000, accountID)
	assert.NoError(t, err)
	assert.Equal(t, expectedPayment, payment)
}

func TestStart(t *testing.T) {
	mockState := new(MockOnchainPaymentState)
	ctx := context.Background()
	mockState.On("Start", ctx).Return(nil)

	err := mockState.Start(ctx)
	assert.NoError(t, err)
}
