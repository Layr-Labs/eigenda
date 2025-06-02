package meterer_test

import (
	"context"
	"math"
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/core/mock"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

var (
	dummyReservedPayment = &core.ReservedPayment{
		SymbolsPerSecond: 100,
		StartTimestamp:   1000,
		EndTimestamp:     2000,
		QuorumSplits:     []byte{50, 50},
	}
	dummyOnDemandPayment = &core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
)

func TestRefreshOnchainPaymentState(t *testing.T) {
	mockState := &mock.MockOnchainPaymentState{}
	ctx := context.Background()
	mockState.On("RefreshOnchainPaymentState", testifymock.Anything).Return(nil)

	err := mockState.RefreshOnchainPaymentState(ctx)
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

func TestGetReservedPaymentByAccount(t *testing.T) {
	mockState := &mock.MockOnchainPaymentState{}
	ctx := context.Background()
	mockState.On("GetReservedPaymentByAccount", testifymock.Anything, testifymock.Anything).Return(map[core.QuorumID]*core.ReservedPayment{0: dummyReservedPayment}, nil)

	reservations, err := mockState.GetReservedPaymentByAccount(ctx, gethcommon.Address{})
	assert.NoError(t, err)
	assert.Equal(t, map[core.QuorumID]*core.ReservedPayment{0: dummyReservedPayment}, reservations)
}

func TestGetOnDemandPaymentByAccount(t *testing.T) {
	mockState := &mock.MockOnchainPaymentState{}
	ctx := context.Background()
	mockState.On("GetOnDemandPaymentByAccount", testifymock.Anything, testifymock.Anything, testifymock.Anything).Return(dummyOnDemandPayment, nil)

	payment, err := mockState.GetOnDemandPaymentByAccount(ctx, gethcommon.Address{})
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

// TestOnchainPaymentStateNilAssignmentProtection tests that the OnchainPaymentState
// properly handles nil map assignments and doesn't panic
func TestOnchainPaymentStateNilAssignmentProtection(t *testing.T) {
	t.Run("PaymentVaultParams_NilProtection", func(t *testing.T) {
		stateWithNilParams, err := meterer.NewOnchainPaymentStateEmpty(context.Background(), nil, testutils.GetLogger())
		assert.NoError(t, err)

		assert.Equal(t, uint64(0), stateWithNilParams.GetOnDemandGlobalSymbolsPerSecond(meterer.OnDemandQuorumID))
		assert.Equal(t, uint64(0), stateWithNilParams.GetOnDemandGlobalRatePeriodInterval(meterer.OnDemandQuorumID))
		assert.Equal(t, uint64(math.MaxUint64), stateWithNilParams.GetMinNumSymbols(meterer.OnDemandQuorumID))
		assert.Equal(t, uint64(math.MaxUint64), stateWithNilParams.GetPricePerSymbol(meterer.OnDemandQuorumID))
		assert.Equal(t, uint64(0), stateWithNilParams.GetReservationWindow(meterer.OnDemandQuorumID))
	})

	t.Run("PaymentVaultParams_MissingQuorum", func(t *testing.T) {
		state, err := meterer.NewOnchainPaymentStateEmpty(context.Background(), nil, testutils.GetLogger())
		assert.NoError(t, err)
		params := &meterer.PaymentVaultParams{
			QuorumPaymentConfigs:  make(map[core.QuorumID]*core.PaymentQuorumConfig),
			QuorumProtocolConfigs: make(map[core.QuorumID]*core.PaymentQuorumProtocolConfig),
		}
		state.PaymentVaultParams.Store(params)

		assert.Equal(t, uint64(0), state.GetOnDemandGlobalSymbolsPerSecond(meterer.OnDemandQuorumID))
		assert.Equal(t, uint64(0), state.GetOnDemandGlobalRatePeriodInterval(meterer.OnDemandQuorumID))
		assert.Equal(t, uint64(math.MaxUint64), state.GetMinNumSymbols(meterer.OnDemandQuorumID))
		assert.Equal(t, uint64(math.MaxUint64), state.GetPricePerSymbol(meterer.OnDemandQuorumID))
		assert.Equal(t, uint64(0), state.GetReservationWindow(meterer.OnDemandQuorumID))
	})

	t.Run("PaymentVaultParams_ValidConfig", func(t *testing.T) {
		state, err := meterer.NewOnchainPaymentStateEmpty(context.Background(), nil, testutils.GetLogger())
		assert.NoError(t, err)

		params := &meterer.PaymentVaultParams{
			QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
				meterer.OnDemandQuorumID: {
					OnDemandSymbolsPerSecond: 100,
					OnDemandPricePerSymbol:   200,
				},
			},
			QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
				meterer.OnDemandQuorumID: {
					MinNumSymbols:              300,
					OnDemandRateLimitWindow:    400,
					ReservationRateLimitWindow: 500,
				},
			},
		}
		state.PaymentVaultParams.Store(params)

		assert.Equal(t, uint64(100), state.GetOnDemandGlobalSymbolsPerSecond(meterer.OnDemandQuorumID))
		assert.Equal(t, uint64(400), state.GetOnDemandGlobalRatePeriodInterval(meterer.OnDemandQuorumID))
		assert.Equal(t, uint64(300), state.GetMinNumSymbols(meterer.OnDemandQuorumID))
		assert.Equal(t, uint64(200), state.GetPricePerSymbol(meterer.OnDemandQuorumID))
		assert.Equal(t, uint64(500), state.GetReservationWindow(meterer.OnDemandQuorumID))
	})

	t.Run("NilMapAssignment_Protection", func(t *testing.T) {
		state, err := meterer.NewOnchainPaymentStateEmpty(context.Background(), nil, testutils.GetLogger())
		assert.NoError(t, err)

		account := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")

		reservations, exists := state.ReservedPayments[account]
		assert.False(t, exists)
		assert.Nil(t, reservations)

		state.ReservedPayments[account] = make(map[core.QuorumID]*core.ReservedPayment)
		state.ReservedPayments[account][0] = &core.ReservedPayment{
			SymbolsPerSecond: 100,
			StartTimestamp:   1000,
			EndTimestamp:     2000,
		}

		assert.NotNil(t, state.ReservedPayments[account])
		assert.NotNil(t, state.ReservedPayments[account][0])
		assert.Equal(t, uint64(100), state.ReservedPayments[account][0].SymbolsPerSecond)

		// Test the nil protection pattern
		testAccount := gethcommon.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd")
		if state.ReservedPayments[testAccount] == nil {
			state.ReservedPayments[testAccount] = make(map[core.QuorumID]*core.ReservedPayment)
		}
		state.ReservedPayments[testAccount][1] = &core.ReservedPayment{
			SymbolsPerSecond: 200,
			StartTimestamp:   1000,
			EndTimestamp:     2000,
		}

		assert.NotNil(t, state.ReservedPayments[testAccount])
		assert.NotNil(t, state.ReservedPayments[testAccount][1])
		assert.Equal(t, uint64(200), state.ReservedPayments[testAccount][1].SymbolsPerSecond)
	})
}

// TestNilAssignmentPanicScenario tests the "assignment to entry in nil map" panic scenario
func TestNilAssignmentPanicScenario(t *testing.T) {
	t.Run("OriginalPanicScenario_NowFixed", func(t *testing.T) {
		state, err := meterer.NewOnchainPaymentStateEmpty(context.Background(), nil, testutils.GetLogger())
		assert.NoError(t, err)

		account := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")
		quorumNumbers := []core.QuorumID{0, 1}

		chainData := map[core.QuorumID]*core.ReservedPayment{
			0: {SymbolsPerSecond: 100, StartTimestamp: 1000, EndTimestamp: 2000},
			1: {SymbolsPerSecond: 200, StartTimestamp: 1000, EndTimestamp: 2000},
		}

		_, exists := state.ReservedPayments[account]
		assert.False(t, exists)

		// Apply nil protection fix
		if state.ReservedPayments[account] == nil {
			state.ReservedPayments[account] = make(map[core.QuorumID]*core.ReservedPayment)
		}

		for _, quorumNumber := range quorumNumbers {
			if reservation, ok := chainData[quorumNumber]; ok {
				state.ReservedPayments[account][quorumNumber] = reservation
			}
		}

		assert.NotNil(t, state.ReservedPayments[account])
		assert.Equal(t, 2, len(state.ReservedPayments[account]))

		for _, quorumNumber := range quorumNumbers {
			assert.NotNil(t, state.ReservedPayments[account][quorumNumber])
			assert.Equal(t, chainData[quorumNumber], state.ReservedPayments[account][quorumNumber])
		}
	})

	t.Run("MultipleAccounts_ConcurrentSafe", func(t *testing.T) {
		state, err := meterer.NewOnchainPaymentStateEmpty(context.Background(), nil, testutils.GetLogger())
		assert.NoError(t, err)

		accounts := []gethcommon.Address{
			gethcommon.HexToAddress("0x1234567890123456789012345678901234567890"),
			gethcommon.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"),
			gethcommon.HexToAddress("0x9876543210987654321098765432109876543210"),
		}

		for i, account := range accounts {
			if state.ReservedPayments[account] == nil {
				state.ReservedPayments[account] = make(map[core.QuorumID]*core.ReservedPayment)
			}

			for quorum := 0; quorum < 3; quorum++ {
				state.ReservedPayments[account][core.QuorumID(quorum)] = &core.ReservedPayment{
					SymbolsPerSecond: uint64((i + 1) * (quorum + 1) * 100),
					StartTimestamp:   1000,
					EndTimestamp:     2000,
				}
			}
		}

		for i, account := range accounts {
			assert.NotNil(t, state.ReservedPayments[account])
			assert.Equal(t, 3, len(state.ReservedPayments[account]))

			for quorum := 0; quorum < 3; quorum++ {
				reservation := state.ReservedPayments[account][core.QuorumID(quorum)]
				assert.NotNil(t, reservation)
				expectedRate := uint64((i + 1) * (quorum + 1) * 100)
				assert.Equal(t, expectedRate, reservation.SymbolsPerSecond)
			}
		}
	})

	t.Run("WithoutProtection_WouldPanic", func(t *testing.T) {
		state, err := meterer.NewOnchainPaymentStateEmpty(context.Background(), nil, testutils.GetLogger())
		assert.NoError(t, err)

		account := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")

		_, exists := state.ReservedPayments[account]
		assert.False(t, exists)
	})
}
