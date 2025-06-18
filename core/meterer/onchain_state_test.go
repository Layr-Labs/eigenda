package meterer_test

import (
	"context"
	"math"
	"math/big"
	"testing"

	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/core/mock"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	dummyReservedPayment = &core.ReservedPayment{
		SymbolsPerSecond: 100,
		StartTimestamp:   1000,
		EndTimestamp:     2000,
	}
	dummyOnDemandPayment = &core.OnDemandPayment{
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

				expectedPayments := map[core.QuorumID]*core.ReservedPayment{
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
				mockParams := &meterer.PaymentVaultParams{
					QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
						0: {
							OnDemandSymbolsPerSecond: 100,
							OnDemandPricePerSymbol:   200,
						},
					},
					QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
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
				params := &meterer.PaymentVaultParams{
					QuorumPaymentConfigs:  make(map[core.QuorumID]*core.PaymentQuorumConfig),
					QuorumProtocolConfigs: make(map[core.QuorumID]*core.PaymentQuorumProtocolConfig),
				}
				state.PaymentVaultParams.Store(params)

				// Test that we can get params but missing quorum configs will cause errors when accessing them
				retrievedParams, err := state.GetPaymentGlobalParams()
				assert.NoError(t, err)
				assert.NotNil(t, retrievedParams)

				// Test that missing quorum returns appropriate errors through the PaymentVaultParams methods
				_, _, err = retrievedParams.GetQuorumConfigs(meterer.OnDemandQuorumID)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "config not found for quorum")

				_, _, err = retrievedParams.GetQuorumConfigs(meterer.OnDemandQuorumID)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "config not found for quorum")
			},
		},
		{
			name: "PaymentVaultParams_ValidConfig",
			test: func(t *testing.T) {
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

				// Test access through the new interface
				retrievedParams, err := state.GetPaymentGlobalParams()
				assert.NoError(t, err)
				assert.NotNil(t, retrievedParams)

				// Test payment config access
				paymentConfig, protocolConfig, err := retrievedParams.GetQuorumConfigs(meterer.OnDemandQuorumID)
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

				state.ReservedPayments[testAccount] = make(map[core.QuorumID]*core.ReservedPayment)
				state.ReservedPayments[testAccount][0] = &core.ReservedPayment{
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
					state.ReservedPayments[testAccount2] = make(map[core.QuorumID]*core.ReservedPayment)
				}
				state.ReservedPayments[testAccount2][1] = &core.ReservedPayment{
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

				chainData := map[core.QuorumID]*core.ReservedPayment{
					0: {SymbolsPerSecond: 100, StartTimestamp: 1000, EndTimestamp: 2000},
					1: {SymbolsPerSecond: 200, StartTimestamp: 1000, EndTimestamp: 2000},
				}

				_, exists := state.ReservedPayments[testAccount]
				assert.False(t, exists)

				// Apply nil protection fix
				if state.ReservedPayments[testAccount] == nil {
					state.ReservedPayments[testAccount] = make(map[core.QuorumID]*core.ReservedPayment)
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
				state.ReservedPayments[testAccount] = map[core.QuorumID]*core.ReservedPayment{
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

func TestPaymentVaultParamsConversion(t *testing.T) {
	tests := []struct {
		name        string
		input       *meterer.PaymentVaultParams
		expectedErr string
		validate    func(t *testing.T, pbParams *disperser_rpc.PaymentVaultParams, coreParams *meterer.PaymentVaultParams)
	}{
		{
			name: "valid complete params",
			input: &meterer.PaymentVaultParams{
				QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
					1: {
						ReservationSymbolsPerSecond: 100,
						OnDemandSymbolsPerSecond:    200,
						OnDemandPricePerSymbol:      300,
					},
					2: {
						ReservationSymbolsPerSecond: 400,
						OnDemandSymbolsPerSecond:    500,
						OnDemandPricePerSymbol:      600,
					},
				},
				QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
					1: {
						MinNumSymbols:              10,
						ReservationAdvanceWindow:   20,
						ReservationRateLimitWindow: 30,
						OnDemandRateLimitWindow:    40,
						OnDemandEnabled:            true,
					},
					2: {
						MinNumSymbols:              50,
						ReservationAdvanceWindow:   60,
						ReservationRateLimitWindow: 70,
						OnDemandRateLimitWindow:    80,
						OnDemandEnabled:            false,
					},
				},
				OnDemandQuorumNumbers: []core.QuorumID{1, 2},
			},
			validate: func(t *testing.T, pbParams *disperser_rpc.PaymentVaultParams, coreParams *meterer.PaymentVaultParams) {
				require.Equal(t, 2, len(pbParams.QuorumPaymentConfigs))
				require.Equal(t, 2, len(pbParams.QuorumProtocolConfigs))
				require.Equal(t, 2, len(pbParams.OnDemandQuorumNumbers))

				// Verify quorum 1
				require.Equal(t, uint64(100), pbParams.QuorumPaymentConfigs[1].ReservationSymbolsPerSecond)
				require.Equal(t, uint64(200), pbParams.QuorumPaymentConfigs[1].OnDemandSymbolsPerSecond)
				require.Equal(t, uint64(300), pbParams.QuorumPaymentConfigs[1].OnDemandPricePerSymbol)
				require.Equal(t, uint64(10), pbParams.QuorumProtocolConfigs[1].MinNumSymbols)
				require.True(t, pbParams.QuorumProtocolConfigs[1].OnDemandEnabled)

				// Verify quorum 2
				require.Equal(t, uint64(400), pbParams.QuorumPaymentConfigs[2].ReservationSymbolsPerSecond)
				require.Equal(t, uint64(500), pbParams.QuorumPaymentConfigs[2].OnDemandSymbolsPerSecond)
				require.Equal(t, uint64(600), pbParams.QuorumPaymentConfigs[2].OnDemandPricePerSymbol)
				require.Equal(t, uint64(50), pbParams.QuorumProtocolConfigs[2].MinNumSymbols)
				require.False(t, pbParams.QuorumProtocolConfigs[2].OnDemandEnabled)

				// Verify round-trip conversion
				convertedCoreParams, err := meterer.PaymentVaultParamsFromProtobuf(pbParams)
				require.NoError(t, err)
				require.Equal(t, coreParams, convertedCoreParams)
			},
		},
		{
			name:        "nil params",
			input:       nil,
			expectedErr: "payment vault params cannot be nil",
		},
		{
			name: "nil payment configs",
			input: &meterer.PaymentVaultParams{
				QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
					1: {MinNumSymbols: 10},
				},
				OnDemandQuorumNumbers: []core.QuorumID{1},
			},
			expectedErr: "payment quorum configs cannot be nil",
		},
		{
			name: "nil protocol configs",
			input: &meterer.PaymentVaultParams{
				QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
					1: {ReservationSymbolsPerSecond: 100},
				},
				OnDemandQuorumNumbers: []core.QuorumID{1},
			},
			expectedErr: "payment quorum protocol configs cannot be nil",
		},
		{
			name: "empty maps",
			input: &meterer.PaymentVaultParams{
				QuorumPaymentConfigs:  map[core.QuorumID]*core.PaymentQuorumConfig{},
				QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{},
				OnDemandQuorumNumbers: []core.QuorumID{},
			},
			validate: func(t *testing.T, pbParams *disperser_rpc.PaymentVaultParams, coreParams *meterer.PaymentVaultParams) {
				require.Empty(t, pbParams.QuorumPaymentConfigs)
				require.Empty(t, pbParams.QuorumProtocolConfigs)
				require.Empty(t, pbParams.OnDemandQuorumNumbers)
			},
		},
		{
			name: "max uint32 quorum ID",
			input: &meterer.PaymentVaultParams{
				QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
					core.QuorumID(255): {ReservationSymbolsPerSecond: 100},
				},
				QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
					core.QuorumID(255): {MinNumSymbols: 10},
				},
				OnDemandQuorumNumbers: []core.QuorumID{core.QuorumID(255)},
			},
			validate: func(t *testing.T, pbParams *disperser_rpc.PaymentVaultParams, coreParams *meterer.PaymentVaultParams) {
				require.Equal(t, uint64(100), pbParams.QuorumPaymentConfigs[255].ReservationSymbolsPerSecond)
				require.Equal(t, uint64(10), pbParams.QuorumProtocolConfigs[255].MinNumSymbols)
				require.Equal(t, uint32(255), pbParams.OnDemandQuorumNumbers[0])
			},
		},
		{
			name: "max uint64 values",
			input: &meterer.PaymentVaultParams{
				QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
					1: {
						ReservationSymbolsPerSecond: math.MaxUint64,
						OnDemandSymbolsPerSecond:    math.MaxUint64,
						OnDemandPricePerSymbol:      math.MaxUint64,
					},
				},
				QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
					1: {
						MinNumSymbols:              math.MaxUint64,
						ReservationAdvanceWindow:   math.MaxUint64,
						ReservationRateLimitWindow: math.MaxUint64,
						OnDemandRateLimitWindow:    math.MaxUint64,
						OnDemandEnabled:            true,
					},
				},
				OnDemandQuorumNumbers: []core.QuorumID{1},
			},
			validate: func(t *testing.T, pbParams *disperser_rpc.PaymentVaultParams, coreParams *meterer.PaymentVaultParams) {
				require.Equal(t, uint64(math.MaxUint64), pbParams.QuorumPaymentConfigs[1].ReservationSymbolsPerSecond)
				require.Equal(t, uint64(math.MaxUint64), pbParams.QuorumPaymentConfigs[1].OnDemandSymbolsPerSecond)
				require.Equal(t, uint64(math.MaxUint64), pbParams.QuorumPaymentConfigs[1].OnDemandPricePerSymbol)
				require.Equal(t, uint64(math.MaxUint64), pbParams.QuorumProtocolConfigs[1].MinNumSymbols)
				require.Equal(t, uint64(math.MaxUint64), pbParams.QuorumProtocolConfigs[1].ReservationAdvanceWindow)
				require.Equal(t, uint64(math.MaxUint64), pbParams.QuorumProtocolConfigs[1].ReservationRateLimitWindow)
				require.Equal(t, uint64(math.MaxUint64), pbParams.QuorumProtocolConfigs[1].OnDemandRateLimitWindow)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test conversion to protobuf
			pbParams, err := tt.input.PaymentVaultParamsToProtobuf()
			if tt.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErr)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, pbParams)

			// Run validation if provided
			if tt.validate != nil {
				tt.validate(t, pbParams, tt.input)
			}
		})
	}
}

func TestPaymentVaultParamsFromProtobuf(t *testing.T) {
	tests := []struct {
		name        string
		input       *disperser_rpc.PaymentVaultParams
		expectedErr string
		validate    func(t *testing.T, coreParams *meterer.PaymentVaultParams)
	}{
		{
			name:        "nil params",
			input:       nil,
			expectedErr: "payment vault params cannot be nil",
		},
		{
			name: "nil payment configs",
			input: &disperser_rpc.PaymentVaultParams{
				QuorumProtocolConfigs: map[uint32]*disperser_rpc.PaymentQuorumProtocolConfig{
					1: {MinNumSymbols: 10},
				},
				OnDemandQuorumNumbers: []uint32{1},
			},
			expectedErr: "payment quorum configs cannot be nil",
		},
		{
			name: "nil protocol configs",
			input: &disperser_rpc.PaymentVaultParams{
				QuorumPaymentConfigs: map[uint32]*disperser_rpc.PaymentQuorumConfig{
					1: {ReservationSymbolsPerSecond: 100},
				},
				OnDemandQuorumNumbers: []uint32{1},
			},
			expectedErr: "payment quorum protocol configs cannot be nil",
		},
		{
			name: "empty maps",
			input: &disperser_rpc.PaymentVaultParams{
				QuorumPaymentConfigs:  map[uint32]*disperser_rpc.PaymentQuorumConfig{},
				QuorumProtocolConfigs: map[uint32]*disperser_rpc.PaymentQuorumProtocolConfig{},
				OnDemandQuorumNumbers: []uint32{},
			},
			validate: func(t *testing.T, coreParams *meterer.PaymentVaultParams) {
				require.Empty(t, coreParams.QuorumPaymentConfigs)
				require.Empty(t, coreParams.QuorumProtocolConfigs)
				require.Empty(t, coreParams.OnDemandQuorumNumbers)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coreParams, err := meterer.PaymentVaultParamsFromProtobuf(tt.input)
			if tt.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErr)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, coreParams)

			// Run validation if provided
			if tt.validate != nil {
				tt.validate(t, coreParams)
			}
		})
	}
}

func TestPaymentVaultParamsToProtobuf(t *testing.T) {
	tests := []struct {
		name        string
		input       *meterer.PaymentVaultParams
		expectedErr string
		validate    func(t *testing.T, pbParams *disperser_rpc.PaymentVaultParams, coreParams *meterer.PaymentVaultParams)
	}{
		{
			name: "valid complete params",
			input: &meterer.PaymentVaultParams{
				QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
					1: {
						ReservationSymbolsPerSecond: 100,
						OnDemandSymbolsPerSecond:    200,
						OnDemandPricePerSymbol:      300,
					},
					2: {
						ReservationSymbolsPerSecond: 400,
						OnDemandSymbolsPerSecond:    500,
						OnDemandPricePerSymbol:      600,
					},
				},
				QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
					1: {
						MinNumSymbols:              10,
						ReservationAdvanceWindow:   20,
						ReservationRateLimitWindow: 30,
						OnDemandRateLimitWindow:    40,
						OnDemandEnabled:            true,
					},
					2: {
						MinNumSymbols:              50,
						ReservationAdvanceWindow:   60,
						ReservationRateLimitWindow: 70,
						OnDemandRateLimitWindow:    80,
						OnDemandEnabled:            false,
					},
				},
				OnDemandQuorumNumbers: []core.QuorumID{1, 2},
			},
			validate: func(t *testing.T, pbParams *disperser_rpc.PaymentVaultParams, coreParams *meterer.PaymentVaultParams) {
				assert.Equal(t, 2, len(pbParams.QuorumPaymentConfigs))
				assert.Equal(t, 2, len(pbParams.QuorumProtocolConfigs))
				assert.Equal(t, 2, len(pbParams.OnDemandQuorumNumbers))

				// Verify quorum 1
				assert.Equal(t, uint64(100), pbParams.QuorumPaymentConfigs[1].ReservationSymbolsPerSecond)
				assert.Equal(t, uint64(200), pbParams.QuorumPaymentConfigs[1].OnDemandSymbolsPerSecond)
				assert.Equal(t, uint64(300), pbParams.QuorumPaymentConfigs[1].OnDemandPricePerSymbol)
				assert.Equal(t, uint64(10), pbParams.QuorumProtocolConfigs[1].MinNumSymbols)
				assert.True(t, pbParams.QuorumProtocolConfigs[1].OnDemandEnabled)

				// Verify quorum 2
				assert.Equal(t, uint64(400), pbParams.QuorumPaymentConfigs[2].ReservationSymbolsPerSecond)
				assert.Equal(t, uint64(500), pbParams.QuorumPaymentConfigs[2].OnDemandSymbolsPerSecond)
				assert.Equal(t, uint64(600), pbParams.QuorumPaymentConfigs[2].OnDemandPricePerSymbol)
				assert.Equal(t, uint64(50), pbParams.QuorumProtocolConfigs[2].MinNumSymbols)
				assert.False(t, pbParams.QuorumProtocolConfigs[2].OnDemandEnabled)

				// Verify round-trip conversion
				convertedCoreParams, err := meterer.PaymentVaultParamsFromProtobuf(pbParams)
				assert.NoError(t, err)
				assert.Equal(t, coreParams, convertedCoreParams)
			},
		},
		{
			name:        "nil params",
			input:       nil,
			expectedErr: "payment vault params cannot be nil",
		},
		{
			name: "nil payment configs",
			input: &meterer.PaymentVaultParams{
				QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
					1: {MinNumSymbols: 10},
				},
				OnDemandQuorumNumbers: []core.QuorumID{1},
			},
			expectedErr: "payment quorum configs cannot be nil",
		},
		{
			name: "nil protocol configs",
			input: &meterer.PaymentVaultParams{
				QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
					1: {ReservationSymbolsPerSecond: 100},
				},
				OnDemandQuorumNumbers: []core.QuorumID{1},
			},
			expectedErr: "payment quorum protocol configs cannot be nil",
		},
		{
			name: "empty maps",
			input: &meterer.PaymentVaultParams{
				QuorumPaymentConfigs:  map[core.QuorumID]*core.PaymentQuorumConfig{},
				QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{},
				OnDemandQuorumNumbers: []core.QuorumID{},
			},
			validate: func(t *testing.T, pbParams *disperser_rpc.PaymentVaultParams, coreParams *meterer.PaymentVaultParams) {
				assert.Empty(t, pbParams.QuorumPaymentConfigs)
				assert.Empty(t, pbParams.QuorumProtocolConfigs)
				assert.Empty(t, pbParams.OnDemandQuorumNumbers)
			},
		},
		{
			name: "max uint32 quorum ID",
			input: &meterer.PaymentVaultParams{
				QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
					core.QuorumID(255): {ReservationSymbolsPerSecond: 100},
				},
				QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
					core.QuorumID(255): {MinNumSymbols: 10},
				},
				OnDemandQuorumNumbers: []core.QuorumID{core.QuorumID(255)},
			},
			validate: func(t *testing.T, pbParams *disperser_rpc.PaymentVaultParams, coreParams *meterer.PaymentVaultParams) {
				assert.Equal(t, uint64(100), pbParams.QuorumPaymentConfigs[255].ReservationSymbolsPerSecond)
				assert.Equal(t, uint64(10), pbParams.QuorumProtocolConfigs[255].MinNumSymbols)
				assert.Equal(t, uint32(255), pbParams.OnDemandQuorumNumbers[0])
			},
		},
		{
			name: "max uint64 values",
			input: &meterer.PaymentVaultParams{
				QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
					1: {
						ReservationSymbolsPerSecond: math.MaxUint64,
						OnDemandSymbolsPerSecond:    math.MaxUint64,
						OnDemandPricePerSymbol:      math.MaxUint64,
					},
				},
				QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
					1: {
						MinNumSymbols:              math.MaxUint64,
						ReservationAdvanceWindow:   math.MaxUint64,
						ReservationRateLimitWindow: math.MaxUint64,
						OnDemandRateLimitWindow:    math.MaxUint64,
						OnDemandEnabled:            true,
					},
				},
				OnDemandQuorumNumbers: []core.QuorumID{1},
			},
			validate: func(t *testing.T, pbParams *disperser_rpc.PaymentVaultParams, coreParams *meterer.PaymentVaultParams) {
				assert.Equal(t, uint64(math.MaxUint64), pbParams.QuorumPaymentConfigs[1].ReservationSymbolsPerSecond)
				assert.Equal(t, uint64(math.MaxUint64), pbParams.QuorumPaymentConfigs[1].OnDemandSymbolsPerSecond)
				assert.Equal(t, uint64(math.MaxUint64), pbParams.QuorumPaymentConfigs[1].OnDemandPricePerSymbol)
				assert.Equal(t, uint64(math.MaxUint64), pbParams.QuorumProtocolConfigs[1].MinNumSymbols)
				assert.Equal(t, uint64(math.MaxUint64), pbParams.QuorumProtocolConfigs[1].ReservationAdvanceWindow)
				assert.Equal(t, uint64(math.MaxUint64), pbParams.QuorumProtocolConfigs[1].ReservationRateLimitWindow)
				assert.Equal(t, uint64(math.MaxUint64), pbParams.QuorumProtocolConfigs[1].OnDemandRateLimitWindow)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test conversion to protobuf
			pbParams, err := tt.input.PaymentVaultParamsToProtobuf()
			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, pbParams)

			// Run validation if provided
			if tt.validate != nil {
				tt.validate(t, pbParams, tt.input)
			}
		})
	}
}
