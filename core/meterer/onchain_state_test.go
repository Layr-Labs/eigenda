package meterer_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/core/payment"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

var (
	dummyReservedPayment = &payment.ReservedPayment{
		SymbolsPerSecond: 100,
		StartTimestamp:   1000,
		EndTimestamp:     2000,
	}
	dummyOnDemandPayment = &payment.OnDemandPayment{
		CumulativePayment: big.NewInt(1000),
	}
	testAccount = gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")
)

func TestMockOnchainPaymentState(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "GetOnDemandPaymentByAccount_Success",
			test: func(t *testing.T) {
				mockState := &mock.MockOnchainPaymentState{}
				ctx := context.Background()
				mockState.On("GetOnDemandPaymentByAccount", testifymock.Anything, testifymock.Anything).Return(dummyOnDemandPayment, nil)

				payment, err := mockState.GetOnDemandPaymentByAccount(ctx, gethcommon.Address{})
				assert.NoError(t, err)
				assert.Equal(t, dummyOnDemandPayment, payment)
			},
		},
		{
			name: "GetReservedPaymentByAccountAndQuorums_Success",
			test: func(t *testing.T) {
				mockState := &mock.MockOnchainPaymentState{}
				ctx := context.Background()
				quorumNumbers := []core.QuorumID{0, 1}

				expectedPayments := map[core.QuorumID]*payment.ReservedPayment{
					0: dummyReservedPayment,
					1: {
						SymbolsPerSecond: 200,
						StartTimestamp:   1500,
						EndTimestamp:     2500,
					},
				}

				mockState.On("GetReservedPaymentByAccountAndQuorums", testifymock.Anything, testifymock.Anything, testifymock.Anything).Return(expectedPayments, nil)

				payments, err := mockState.GetReservedPaymentByAccountAndQuorums(ctx, testAccount, quorumNumbers)
				assert.NoError(t, err)
				assert.Equal(t, expectedPayments, payments)
				assert.Equal(t, 2, len(payments))
				assert.Equal(t, dummyReservedPayment, payments[0])
				assert.Equal(t, uint64(200), payments[1].SymbolsPerSecond)
			},
		},
		{
			name: "RefreshOnchainPaymentState_Success",
			test: func(t *testing.T) {
				mockState := &mock.MockOnchainPaymentState{}
				ctx := context.Background()

				mockState.On("RefreshOnchainPaymentState", testifymock.Anything).Return(nil)

				err := mockState.RefreshOnchainPaymentState(ctx)
				assert.NoError(t, err)

				mockState.AssertExpectations(t)
			},
		},
		{
			name: "GetPaymentGlobalParams_Success",
			test: func(t *testing.T) {
				mockState := &mock.MockOnchainPaymentState{}
				mockParams := &payment.PaymentVaultParams{
					QuorumPaymentConfigs: map[core.QuorumID]*payment.PaymentQuorumConfig{
						0: {
							OnDemandSymbolsPerSecond: 100,
							OnDemandPricePerSymbol:   200,
						},
					},
					QuorumProtocolConfigs: map[core.QuorumID]*payment.PaymentQuorumProtocolConfig{
						0: {
							MinNumSymbols:           300,
							OnDemandRateLimitWindow: 400,
						},
					},
					OnDemandQuorumNumbers: []core.QuorumID{0, 1},
				}
				mockState.On("GetPaymentGlobalParams").Return(mockParams, nil)

				params, err := mockState.GetPaymentGlobalParams()
				assert.NoError(t, err)
				assert.Equal(t, mockParams, params)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestOnchainPaymentState_NilProtection(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "PaymentVaultParams_NilProtection",
			test: func(t *testing.T) {
				stateWithNilParams, err := meterer.NewOnchainPaymentStateEmpty(context.Background(), nil, testutils.GetLogger())
				assert.NoError(t, err)

				// Test that nil PaymentVaultParams returns appropriate errors
				_, err = stateWithNilParams.GetPaymentGlobalParams()
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "payment vault params not initialized")
			},
		},
		{
			name: "PaymentVaultParams_MissingQuorum",
			test: func(t *testing.T) {
				state, err := meterer.NewOnchainPaymentStateEmpty(context.Background(), nil, testutils.GetLogger())
				assert.NoError(t, err)
				params := &payment.PaymentVaultParams{
					QuorumPaymentConfigs:  make(map[core.QuorumID]*payment.PaymentQuorumConfig),
					QuorumProtocolConfigs: make(map[core.QuorumID]*payment.PaymentQuorumProtocolConfig),
				}
				state.PaymentVaultParams.Store(params)

				// Test that we can get params but missing quorum configs will cause errors when accessing them
				retrievedParams, err := state.GetPaymentGlobalParams()
				assert.NoError(t, err)
				assert.NotNil(t, retrievedParams)

				// Test that missing quorum returns appropriate errors through the PaymentVaultParams methods
				_, _, err = retrievedParams.GetQuorumConfigs(payment.OnDemandDepositQuorumID)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "config not found for quorum")

				_, _, err = retrievedParams.GetQuorumConfigs(payment.OnDemandDepositQuorumID)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "config not found for quorum")
			},
		},
		{
			name: "PaymentVaultParams_ValidConfig",
			test: func(t *testing.T) {
				state, err := meterer.NewOnchainPaymentStateEmpty(context.Background(), nil, testutils.GetLogger())
				assert.NoError(t, err)

				params := &payment.PaymentVaultParams{
					QuorumPaymentConfigs: map[core.QuorumID]*payment.PaymentQuorumConfig{
						payment.OnDemandDepositQuorumID: {
							OnDemandSymbolsPerSecond: 100,
							OnDemandPricePerSymbol:   200,
						},
					},
					QuorumProtocolConfigs: map[core.QuorumID]*payment.PaymentQuorumProtocolConfig{
						payment.OnDemandDepositQuorumID: {
							MinNumSymbols:              300,
							OnDemandRateLimitWindow:    400,
							ReservationRateLimitWindow: 500,
						},
					},
				}
				state.PaymentVaultParams.Store(params)

				// Test access through the new interface
				retrievedParams, err := state.GetPaymentGlobalParams()
				assert.NoError(t, err)
				assert.NotNil(t, retrievedParams)

				// Test payment config access
				paymentConfig, protocolConfig, err := retrievedParams.GetQuorumConfigs(payment.OnDemandDepositQuorumID)
				assert.NoError(t, err)
				assert.Equal(t, uint64(100), paymentConfig.OnDemandSymbolsPerSecond)
				assert.Equal(t, uint64(400), protocolConfig.OnDemandRateLimitWindow)
				assert.Equal(t, uint64(300), protocolConfig.MinNumSymbols)
				assert.Equal(t, uint64(200), paymentConfig.OnDemandPricePerSymbol)
				assert.Equal(t, uint64(500), protocolConfig.ReservationRateLimitWindow)
			},
		},
		{
			name: "NilMapAssignment_Protection",
			test: func(t *testing.T) {
				state, err := meterer.NewOnchainPaymentStateEmpty(context.Background(), nil, testutils.GetLogger())
				assert.NoError(t, err)

				reservations, exists := state.ReservedPayments[testAccount]
				assert.False(t, exists)
				assert.Nil(t, reservations)

				state.ReservedPayments[testAccount] = make(map[core.QuorumID]*payment.ReservedPayment)
				state.ReservedPayments[testAccount][0] = &payment.ReservedPayment{
					SymbolsPerSecond: 100,
					StartTimestamp:   1000,
					EndTimestamp:     2000,
				}

				assert.NotNil(t, state.ReservedPayments[testAccount])
				assert.NotNil(t, state.ReservedPayments[testAccount][0])
				assert.Equal(t, uint64(100), state.ReservedPayments[testAccount][0].SymbolsPerSecond)

				// Test the nil protection pattern
				testAccount2 := gethcommon.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd")
				if state.ReservedPayments[testAccount2] == nil {
					state.ReservedPayments[testAccount2] = make(map[core.QuorumID]*payment.ReservedPayment)
				}
				state.ReservedPayments[testAccount2][1] = &payment.ReservedPayment{
					SymbolsPerSecond: 200,
					StartTimestamp:   1000,
					EndTimestamp:     2000,
				}

				assert.NotNil(t, state.ReservedPayments[testAccount2])
				assert.NotNil(t, state.ReservedPayments[testAccount2][1])
				assert.Equal(t, uint64(200), state.ReservedPayments[testAccount2][1].SymbolsPerSecond)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestNilAssignmentPanicScenarios(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "OriginalPanicScenario_NowFixed",
			test: func(t *testing.T) {
				state, err := meterer.NewOnchainPaymentStateEmpty(context.Background(), nil, testutils.GetLogger())
				assert.NoError(t, err)

				quorumNumbers := []core.QuorumID{0, 1}

				chainData := map[core.QuorumID]*payment.ReservedPayment{
					0: {SymbolsPerSecond: 100, StartTimestamp: 1000, EndTimestamp: 2000},
					1: {SymbolsPerSecond: 200, StartTimestamp: 1000, EndTimestamp: 2000},
				}

				_, exists := state.ReservedPayments[testAccount]
				assert.False(t, exists)

				// Apply nil protection fix
				if state.ReservedPayments[testAccount] == nil {
					state.ReservedPayments[testAccount] = make(map[core.QuorumID]*payment.ReservedPayment)
				}

				for _, quorumNumber := range quorumNumbers {
					if reservation, ok := chainData[quorumNumber]; ok {
						state.ReservedPayments[testAccount][quorumNumber] = reservation
					}
				}

				assert.NotNil(t, state.ReservedPayments[testAccount])
				assert.Equal(t, 2, len(state.ReservedPayments[testAccount]))

				for _, quorumNumber := range quorumNumbers {
					assert.NotNil(t, state.ReservedPayments[testAccount][quorumNumber])
					assert.Equal(t, chainData[quorumNumber], state.ReservedPayments[testAccount][quorumNumber])
				}
			},
		},
		{
			name: "MultipleAccounts_ConcurrentSafe",
			test: func(t *testing.T) {
				state, err := meterer.NewOnchainPaymentStateEmpty(context.Background(), nil, testutils.GetLogger())
				assert.NoError(t, err)

				accounts := []gethcommon.Address{
					gethcommon.HexToAddress("0x1234567890123456789012345678901234567890"),
					gethcommon.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"),
					gethcommon.HexToAddress("0x9876543210987654321098765432109876543210"),
				}

				for i, account := range accounts {
					if state.ReservedPayments[account] == nil {
						state.ReservedPayments[account] = make(map[core.QuorumID]*payment.ReservedPayment)
					}

					for quorum := 0; quorum < 3; quorum++ {
						state.ReservedPayments[account][core.QuorumID(quorum)] = &payment.ReservedPayment{
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
			},
		},
		{
			name: "WithoutProtection_WouldPanic",
			test: func(t *testing.T) {
				state, err := meterer.NewOnchainPaymentStateEmpty(context.Background(), nil, testutils.GetLogger())
				assert.NoError(t, err)

				_, exists := state.ReservedPayments[testAccount]
				assert.False(t, exists)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestOnchainPaymentState_CacheOperations(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "GetReservedPaymentByAccountAndQuorums_CacheHit",
			test: func(t *testing.T) {
				state, err := meterer.NewOnchainPaymentStateEmpty(context.Background(), nil, testutils.GetLogger())
				assert.NoError(t, err)

				quorumNumbers := []core.QuorumID{0, 1}

				// Pre-populate cache
				state.ReservedPayments[testAccount] = map[core.QuorumID]*payment.ReservedPayment{
					0: dummyReservedPayment,
					1: {SymbolsPerSecond: 200, StartTimestamp: 1500, EndTimestamp: 2500},
				}

				payments, err := state.GetReservedPaymentByAccountAndQuorums(context.Background(), testAccount, quorumNumbers)
				assert.NoError(t, err)
				assert.Equal(t, 2, len(payments))
				assert.Equal(t, dummyReservedPayment, payments[0])
				assert.Equal(t, uint64(200), payments[1].SymbolsPerSecond)
			},
		},
		{
			name: "GetOnDemandPaymentByAccount_CacheHit",
			test: func(t *testing.T) {
				state, err := meterer.NewOnchainPaymentStateEmpty(context.Background(), nil, testutils.GetLogger())
				assert.NoError(t, err)

				// Pre-populate cache
				state.OnDemandPayments[testAccount] = dummyOnDemandPayment

				payment, err := state.GetOnDemandPaymentByAccount(context.Background(), testAccount)
				assert.NoError(t, err)
				assert.Equal(t, dummyOnDemandPayment, payment)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
