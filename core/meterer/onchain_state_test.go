package meterer_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/mock"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

var (
	dummyActiveReservation = core.ActiveReservation{
		SymbolsPerSecond: 100,
		StartTimestamp:   1000,
		EndTimestamp:     2000,
		QuorumSplits:     []byte{50, 50},
	}
	dummyOnDemandPayment = core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
)

func TestRefreshOnchainPaymentState(t *testing.T) {
	mockState := &mock.MockOnchainPaymentState{}
	ctx := context.Background()
	mockState.On("RefreshOnchainPaymentState", testifymock.Anything, testifymock.Anything).Return(nil)

	err := mockState.RefreshOnchainPaymentState(ctx, &eth.Transactor{})
	assert.NoError(t, err)
}

func TestGetCurrentBlockNumber(t *testing.T) {
	mockState := &mock.MockOnchainPaymentState{}
	mockState.On("GetCurrentBlockNumber").Return(uint32(1000), nil)
	ctx := context.Background()
	blockNumber, err := mockState.GetCurrentBlockNumber(ctx)
	assert.NoError(t, err)
	assert.Equal(t, uint32(1000), blockNumber)
}

func TestGetActiveReservations(t *testing.T) {
	mockState := &mock.MockOnchainPaymentState{}
	ctx := context.Background()
	expectedReservations := map[string]core.ActiveReservation{
		"account1": dummyActiveReservation,
	}
	mockState.On("GetActiveReservations", testifymock.Anything, testifymock.Anything).Return(expectedReservations, nil)

	reservations, err := mockState.GetActiveReservations(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedReservations, reservations)
}

func TestGetActiveReservationByAccount(t *testing.T) {
	mockState := &mock.MockOnchainPaymentState{}
	ctx := context.Background()
	mockState.On("GetActiveReservationByAccount", testifymock.Anything, testifymock.Anything).Return(dummyActiveReservation, nil)

	reservation, err := mockState.GetActiveReservationByAccount(ctx, "account1")
	assert.NoError(t, err)
	assert.Equal(t, dummyActiveReservation, reservation)
}

func TestGetOnDemandPayments(t *testing.T) {
	mockState := &mock.MockOnchainPaymentState{}
	ctx := context.Background()
	expectedPayments := map[string]core.OnDemandPayment{
		"account1": dummyOnDemandPayment,
	}
	mockState.On("GetOnDemandPayments", testifymock.Anything, testifymock.Anything).Return(expectedPayments, nil)

	payments, err := mockState.GetOnDemandPayments(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedPayments, payments)
}

func TestGetOnDemandPaymentByAccount(t *testing.T) {
	mockState := &mock.MockOnchainPaymentState{}
	ctx := context.Background()
	accountID := "account1"
	mockState.On("GetOnDemandPaymentByAccount", testifymock.Anything, testifymock.Anything, testifymock.Anything).Return(dummyOnDemandPayment, nil)

	payment, err := mockState.GetOnDemandPaymentByAccount(ctx, accountID)
	assert.NoError(t, err)
	assert.Equal(t, dummyOnDemandPayment, payment)
}

func TestGetOnDemandQuorumNumbers(t *testing.T) {
	mockState := &mock.MockOnchainPaymentState{}
	ctx := context.Background()
	mockState.On("GetOnDemandQuorumNumbers", testifymock.Anything, testifymock.Anything).Return([]uint8{0, 1}, nil)

	quorumNumbers, err := mockState.GetOnDemandQuorumNumbers(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []uint8{0, 1}, quorumNumbers)
}
