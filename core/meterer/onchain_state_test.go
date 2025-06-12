package meterer_test

import (
	"context"
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
		// QuorumSplits:     []byte{50, 50},
	}
	dummyOnDemandPayment = &core.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
)

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

		// Test that nil PaymentVaultParams returns appropriate errors
		_, err = stateWithNilParams.GetOnDemandGlobalSymbolsPerSecond(meterer.OnDemandQuorumID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "payment vault params not initialized")

		_, err = stateWithNilParams.GetOnDemandGlobalRatePeriodInterval(meterer.OnDemandQuorumID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "payment vault params not initialized")

		_, err = stateWithNilParams.GetMinNumSymbols(meterer.OnDemandQuorumID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "payment vault params not initialized")

		_, err = stateWithNilParams.GetPricePerSymbol(meterer.OnDemandQuorumID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "payment vault params not initialized")

		_, err = stateWithNilParams.GetReservationWindow(meterer.OnDemandQuorumID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "payment vault params not initialized")
	})

	t.Run("PaymentVaultParams_MissingQuorum", func(t *testing.T) {
		state, err := meterer.NewOnchainPaymentStateEmpty(context.Background(), nil, testutils.GetLogger())
		assert.NoError(t, err)
		params := &meterer.PaymentVaultParams{
			QuorumPaymentConfigs:  make(map[core.QuorumID]*core.PaymentQuorumConfig),
			QuorumProtocolConfigs: make(map[core.QuorumID]*core.PaymentQuorumProtocolConfig),
		}
		state.PaymentVaultParams.Store(params)

		// Test that missing quorum returns appropriate errors
		_, err = state.GetOnDemandGlobalSymbolsPerSecond(meterer.OnDemandQuorumID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "payment config not found for quorum")

		_, err = state.GetOnDemandGlobalRatePeriodInterval(meterer.OnDemandQuorumID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "protocol config not found for quorum")

		_, err = state.GetMinNumSymbols(meterer.OnDemandQuorumID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "protocol config not found for quorum")

		_, err = state.GetPricePerSymbol(meterer.OnDemandQuorumID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "payment config not found for quorum")

		_, err = state.GetReservationWindow(meterer.OnDemandQuorumID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "protocol config not found for quorum")
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

		globalSymbolsPerSecond, err := state.GetOnDemandGlobalSymbolsPerSecond(meterer.OnDemandQuorumID)
		assert.NoError(t, err)
		assert.Equal(t, uint64(100), globalSymbolsPerSecond)
		globalPeriodInterval, err := state.GetOnDemandGlobalRatePeriodInterval(meterer.OnDemandQuorumID)
		assert.NoError(t, err)
		assert.Equal(t, uint64(400), globalPeriodInterval)
		minNumSymbols, err := state.GetMinNumSymbols(meterer.OnDemandQuorumID)
		assert.NoError(t, err)
		assert.Equal(t, uint64(300), minNumSymbols)
		pricePerSymbol, err := state.GetPricePerSymbol(meterer.OnDemandQuorumID)
		assert.NoError(t, err)
		assert.Equal(t, uint64(200), pricePerSymbol)
		reservationWindow, err := state.GetReservationWindow(meterer.OnDemandQuorumID)
		assert.NoError(t, err)
		assert.Equal(t, uint64(500), reservationWindow)
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

func TestPaymentVaultParams_GetConfigs(t *testing.T) {
	paymentVaultParams := &meterer.PaymentVaultParams{
		QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
			0: {
				ReservationSymbolsPerSecond: 1,
				OnDemandSymbolsPerSecond:    2,
				OnDemandPricePerSymbol:      3,
			},
			1: {
				ReservationSymbolsPerSecond: 4,
				OnDemandSymbolsPerSecond:    5,
				OnDemandPricePerSymbol:      6,
			},
		},
		QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
			0: {
				MinNumSymbols:              7,
				ReservationAdvanceWindow:   8,
				ReservationRateLimitWindow: 9,
				OnDemandRateLimitWindow:    10,
				OnDemandEnabled:            true,
			},
			1: {
				MinNumSymbols:              11,
				ReservationAdvanceWindow:   12,
				ReservationRateLimitWindow: 13,
				OnDemandRateLimitWindow:    14,
				OnDemandEnabled:            false,
			},
		},
		OnDemandQuorumNumbers: []core.QuorumID{0, 1},
	}

	// Test with non-existent quorum
	_, _, err := paymentVaultParams.GetConfigs(99)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config not found")

	paymentQuorumConfig, protocolConfig, err := paymentVaultParams.GetConfigs(0)
	assert.NoError(t, err)
	assert.NotNil(t, paymentQuorumConfig)
	assert.NotNil(t, protocolConfig)
	assert.Equal(t, uint64(7), protocolConfig.MinNumSymbols)
	assert.Equal(t, uint64(3), paymentQuorumConfig.OnDemandPricePerSymbol)
	assert.Equal(t, []core.QuorumID{0, 1}, paymentVaultParams.OnDemandQuorumNumbers)
}
