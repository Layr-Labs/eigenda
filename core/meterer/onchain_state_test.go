package meterer_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/core/mock"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

var (
	dummyActiveReservation = core.ActiveReservation{
		DataRate:       100,
		StartTimestamp: 1000,
		EndTimestamp:   2000,
		QuorumSplit:    []byte{50, 50},
	}
	dummyOnDemandPayment = core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
)

func TestGetCurrentOnchainPaymentState(t *testing.T) {
	mockState := &mock.MockOnchainPaymentState{}
	ctx := context.Background()
	mockState.On("CurrentOnchainPaymentState", testifymock.Anything, testifymock.Anything).Return(meterer.OnchainPaymentState{
		ActiveReservations: map[string]core.ActiveReservation{
			"account1": dummyActiveReservation,
		},
		OnDemandPayments: map[string]core.OnDemandPayment{
			"account1": dummyOnDemandPayment,
		},
	}, nil)

	state, err := mockState.CurrentOnchainPaymentState(ctx, &eth.Transactor{})
	assert.NoError(t, err)
	assert.Equal(t, meterer.OnchainPaymentState{
		ActiveReservations: map[string]core.ActiveReservation{
			"account1": dummyActiveReservation,
		},
		OnDemandPayments: map[string]core.OnDemandPayment{
			"account1": dummyOnDemandPayment,
		},
	}, state)
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

	reservations, err := mockState.GetActiveReservations(ctx, 1000)
	assert.NoError(t, err)
	assert.Equal(t, expectedReservations, reservations)
}

func TestGetActiveReservationByAccount(t *testing.T) {
	mockState := &mock.MockOnchainPaymentState{}
	ctx := context.Background()
	mockState.On("GetActiveReservationsByAccount", testifymock.Anything, testifymock.Anything, testifymock.Anything).Return(dummyActiveReservation, nil)

	reservation, err := mockState.GetActiveReservationsByAccount(ctx, 1000, "account1")
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

	payments, err := mockState.GetOnDemandPayments(ctx, 1000)
	assert.NoError(t, err)
	assert.Equal(t, expectedPayments, payments)
}

func TestGetOnDemandPaymentByAccount(t *testing.T) {
	mockState := &mock.MockOnchainPaymentState{}
	ctx := context.Background()
	accountID := "account1"
	mockState.On("GetOnDemandPaymentByAccount", testifymock.Anything, testifymock.Anything, testifymock.Anything).Return(dummyOnDemandPayment, nil)

	payment, err := mockState.GetOnDemandPaymentByAccount(ctx, 1000, accountID)
	assert.NoError(t, err)
	assert.Equal(t, dummyOnDemandPayment, payment)
}
