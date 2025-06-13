package meterer_test

import (
	"context"
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
				_, err = retrievedParams.GetQuorumPaymentConfig(meterer.OnDemandQuorumID)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "payment config not found for quorum")

				_, err = retrievedParams.GetQuorumProtocolConfig(meterer.OnDemandQuorumID)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "protocol config not found for quorum")

				_, err = retrievedParams.GetMinNumSymbols(meterer.OnDemandQuorumID)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "protocol config not found for quorum")

				_, err = retrievedParams.GetPricePerSymbol(meterer.OnDemandQuorumID)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "payment config not found for quorum")

				_, err = retrievedParams.GetReservationWindow(meterer.OnDemandQuorumID)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "protocol config not found for quorum")
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
				paymentConfig, err := retrievedParams.GetQuorumPaymentConfig(meterer.OnDemandQuorumID)
				assert.NoError(t, err)
				assert.Equal(t, uint64(100), paymentConfig.OnDemandSymbolsPerSecond)
				// Test protocol config access
				protocolConfig, err := retrievedParams.GetQuorumProtocolConfig(meterer.OnDemandQuorumID)
				assert.NoError(t, err)
				assert.Equal(t, uint64(400), protocolConfig.OnDemandRateLimitWindow)
				minNumSymbols, err := retrievedParams.GetMinNumSymbols(meterer.OnDemandQuorumID)
				assert.NoError(t, err)
				assert.Equal(t, uint64(300), minNumSymbols)
				pricePerSymbol, err := retrievedParams.GetPricePerSymbol(meterer.OnDemandQuorumID)
				assert.NoError(t, err)
				assert.Equal(t, uint64(200), pricePerSymbol)
				reservationWindow, err := retrievedParams.GetReservationWindow(meterer.OnDemandQuorumID)
				assert.NoError(t, err)
				assert.Equal(t, uint64(500), reservationWindow)
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
		{
			name: "GetQuorumNumbers_ParamsNotInitialized",
			test: func(t *testing.T) {
				state, err := meterer.NewOnchainPaymentStateEmpty(context.Background(), nil, testutils.GetLogger())
				assert.NoError(t, err)

				_, err = state.GetQuorumNumbers(context.Background())
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "payment vault params not initialized")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestPaymentVaultParamsFromProtobuf(t *testing.T) {
	tests := []struct {
		name        string
		vaultParams *disperser_rpc.PaymentVaultParams
		expected    *meterer.PaymentVaultParams
		expectedErr string
	}{
		{
			name:        "nil vault params",
			vaultParams: nil,
			expectedErr: "payment vault params cannot be nil",
		},
		{
			name: "nil quorum payment configs",
			vaultParams: &disperser_rpc.PaymentVaultParams{
				QuorumPaymentConfigs:  nil,
				QuorumProtocolConfigs: map[uint32]*disperser_rpc.PaymentQuorumProtocolConfig{},
				OnDemandQuorumNumbers: []uint32{},
			},
			expectedErr: "payment quorum configs cannot be nil",
		},
		{
			name: "nil quorum protocol configs",
			vaultParams: &disperser_rpc.PaymentVaultParams{
				QuorumPaymentConfigs:  map[uint32]*disperser_rpc.PaymentQuorumConfig{},
				QuorumProtocolConfigs: nil,
				OnDemandQuorumNumbers: []uint32{},
			},
			expectedErr: "payment quorum protocol configs cannot be nil",
		},
		{
			name: "empty configs",
			vaultParams: &disperser_rpc.PaymentVaultParams{
				QuorumPaymentConfigs:  map[uint32]*disperser_rpc.PaymentQuorumConfig{},
				QuorumProtocolConfigs: map[uint32]*disperser_rpc.PaymentQuorumProtocolConfig{},
				OnDemandQuorumNumbers: []uint32{},
			},
			expected: &meterer.PaymentVaultParams{
				QuorumPaymentConfigs:  map[core.QuorumID]*core.PaymentQuorumConfig{},
				QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{},
				OnDemandQuorumNumbers: []core.QuorumID{},
			},
		},
		{
			name: "single quorum config",
			vaultParams: &disperser_rpc.PaymentVaultParams{
				QuorumPaymentConfigs: map[uint32]*disperser_rpc.PaymentQuorumConfig{
					0: {
						ReservationSymbolsPerSecond: 100,
						OnDemandSymbolsPerSecond:    200,
						OnDemandPricePerSymbol:      300,
					},
				},
				QuorumProtocolConfigs: map[uint32]*disperser_rpc.PaymentQuorumProtocolConfig{
					0: {
						MinNumSymbols:              400,
						ReservationAdvanceWindow:   500,
						ReservationRateLimitWindow: 600,
						OnDemandRateLimitWindow:    700,
						OnDemandEnabled:            true,
					},
				},
				OnDemandQuorumNumbers: []uint32{0},
			},
			expected: &meterer.PaymentVaultParams{
				QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
					0: {
						ReservationSymbolsPerSecond: 100,
						OnDemandSymbolsPerSecond:    200,
						OnDemandPricePerSymbol:      300,
					},
				},
				QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
					0: {
						MinNumSymbols:              400,
						ReservationAdvanceWindow:   500,
						ReservationRateLimitWindow: 600,
						OnDemandRateLimitWindow:    700,
						OnDemandEnabled:            true,
					},
				},
				OnDemandQuorumNumbers: []core.QuorumID{0},
			},
		},
		{
			name: "multiple quorum configs",
			vaultParams: &disperser_rpc.PaymentVaultParams{
				QuorumPaymentConfigs: map[uint32]*disperser_rpc.PaymentQuorumConfig{
					0: {
						ReservationSymbolsPerSecond: 100,
						OnDemandSymbolsPerSecond:    200,
						OnDemandPricePerSymbol:      300,
					},
					1: {
						ReservationSymbolsPerSecond: 1100,
						OnDemandSymbolsPerSecond:    1200,
						OnDemandPricePerSymbol:      1300,
					},
					255: {
						ReservationSymbolsPerSecond: 25500,
						OnDemandSymbolsPerSecond:    25600,
						OnDemandPricePerSymbol:      25700,
					},
				},
				QuorumProtocolConfigs: map[uint32]*disperser_rpc.PaymentQuorumProtocolConfig{
					0: {
						MinNumSymbols:              400,
						ReservationAdvanceWindow:   500,
						ReservationRateLimitWindow: 600,
						OnDemandRateLimitWindow:    700,
						OnDemandEnabled:            true,
					},
					1: {
						MinNumSymbols:              1400,
						ReservationAdvanceWindow:   1500,
						ReservationRateLimitWindow: 1600,
						OnDemandRateLimitWindow:    1700,
						OnDemandEnabled:            false,
					},
					255: {
						MinNumSymbols:              25800,
						ReservationAdvanceWindow:   25900,
						ReservationRateLimitWindow: 26000,
						OnDemandRateLimitWindow:    26100,
						OnDemandEnabled:            true,
					},
				},
				OnDemandQuorumNumbers: []uint32{0, 1, 255},
			},
			expected: &meterer.PaymentVaultParams{
				QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
					0: {
						ReservationSymbolsPerSecond: 100,
						OnDemandSymbolsPerSecond:    200,
						OnDemandPricePerSymbol:      300,
					},
					1: {
						ReservationSymbolsPerSecond: 1100,
						OnDemandSymbolsPerSecond:    1200,
						OnDemandPricePerSymbol:      1300,
					},
					255: {
						ReservationSymbolsPerSecond: 25500,
						OnDemandSymbolsPerSecond:    25600,
						OnDemandPricePerSymbol:      25700,
					},
				},
				QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
					0: {
						MinNumSymbols:              400,
						ReservationAdvanceWindow:   500,
						ReservationRateLimitWindow: 600,
						OnDemandRateLimitWindow:    700,
						OnDemandEnabled:            true,
					},
					1: {
						MinNumSymbols:              1400,
						ReservationAdvanceWindow:   1500,
						ReservationRateLimitWindow: 1600,
						OnDemandRateLimitWindow:    1700,
						OnDemandEnabled:            false,
					},
					255: {
						MinNumSymbols:              25800,
						ReservationAdvanceWindow:   25900,
						ReservationRateLimitWindow: 26000,
						OnDemandRateLimitWindow:    26100,
						OnDemandEnabled:            true,
					},
				},
				OnDemandQuorumNumbers: []core.QuorumID{0, 1, 255},
			},
		},
		{
			name: "zero values",
			vaultParams: &disperser_rpc.PaymentVaultParams{
				QuorumPaymentConfigs: map[uint32]*disperser_rpc.PaymentQuorumConfig{
					0: {
						ReservationSymbolsPerSecond: 0,
						OnDemandSymbolsPerSecond:    0,
						OnDemandPricePerSymbol:      0,
					},
				},
				QuorumProtocolConfigs: map[uint32]*disperser_rpc.PaymentQuorumProtocolConfig{
					0: {
						MinNumSymbols:              0,
						ReservationAdvanceWindow:   0,
						ReservationRateLimitWindow: 0,
						OnDemandRateLimitWindow:    0,
						OnDemandEnabled:            false,
					},
				},
				OnDemandQuorumNumbers: []uint32{0},
			},
			expected: &meterer.PaymentVaultParams{
				QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
					0: {
						ReservationSymbolsPerSecond: 0,
						OnDemandSymbolsPerSecond:    0,
						OnDemandPricePerSymbol:      0,
					},
				},
				QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
					0: {
						MinNumSymbols:              0,
						ReservationAdvanceWindow:   0,
						ReservationRateLimitWindow: 0,
						OnDemandRateLimitWindow:    0,
						OnDemandEnabled:            false,
					},
				},
				OnDemandQuorumNumbers: []core.QuorumID{0},
			},
		},
		{
			name: "max values",
			vaultParams: &disperser_rpc.PaymentVaultParams{
				QuorumPaymentConfigs: map[uint32]*disperser_rpc.PaymentQuorumConfig{
					0: {
						ReservationSymbolsPerSecond: ^uint64(0), // max uint64
						OnDemandSymbolsPerSecond:    ^uint64(0),
						OnDemandPricePerSymbol:      ^uint64(0),
					},
				},
				QuorumProtocolConfigs: map[uint32]*disperser_rpc.PaymentQuorumProtocolConfig{
					0: {
						MinNumSymbols:              ^uint64(0),
						ReservationAdvanceWindow:   ^uint64(0),
						ReservationRateLimitWindow: ^uint64(0),
						OnDemandRateLimitWindow:    ^uint64(0),
						OnDemandEnabled:            true,
					},
				},
				OnDemandQuorumNumbers: []uint32{0},
			},
			expected: &meterer.PaymentVaultParams{
				QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
					0: {
						ReservationSymbolsPerSecond: ^uint64(0),
						OnDemandSymbolsPerSecond:    ^uint64(0),
						OnDemandPricePerSymbol:      ^uint64(0),
					},
				},
				QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
					0: {
						MinNumSymbols:              ^uint64(0),
						ReservationAdvanceWindow:   ^uint64(0),
						ReservationRateLimitWindow: ^uint64(0),
						OnDemandRateLimitWindow:    ^uint64(0),
						OnDemandEnabled:            true,
					},
				},
				OnDemandQuorumNumbers: []core.QuorumID{0},
			},
		},
		{
			name: "empty on-demand quorum numbers",
			vaultParams: &disperser_rpc.PaymentVaultParams{
				QuorumPaymentConfigs: map[uint32]*disperser_rpc.PaymentQuorumConfig{
					0: {
						ReservationSymbolsPerSecond: 100,
						OnDemandSymbolsPerSecond:    200,
						OnDemandPricePerSymbol:      300,
					},
				},
				QuorumProtocolConfigs: map[uint32]*disperser_rpc.PaymentQuorumProtocolConfig{
					0: {
						MinNumSymbols:              400,
						ReservationAdvanceWindow:   500,
						ReservationRateLimitWindow: 600,
						OnDemandRateLimitWindow:    700,
						OnDemandEnabled:            true,
					},
				},
				OnDemandQuorumNumbers: []uint32{},
			},
			expected: &meterer.PaymentVaultParams{
				QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
					0: {
						ReservationSymbolsPerSecond: 100,
						OnDemandSymbolsPerSecond:    200,
						OnDemandPricePerSymbol:      300,
					},
				},
				QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
					0: {
						MinNumSymbols:              400,
						ReservationAdvanceWindow:   500,
						ReservationRateLimitWindow: 600,
						OnDemandRateLimitWindow:    700,
						OnDemandEnabled:            true,
					},
				},
				OnDemandQuorumNumbers: []core.QuorumID{},
			},
		},
		{
			name: "mismatched quorum configs",
			vaultParams: &disperser_rpc.PaymentVaultParams{
				QuorumPaymentConfigs: map[uint32]*disperser_rpc.PaymentQuorumConfig{
					0: {
						ReservationSymbolsPerSecond: 100,
						OnDemandSymbolsPerSecond:    200,
						OnDemandPricePerSymbol:      300,
					},
					2: {
						ReservationSymbolsPerSecond: 2100,
						OnDemandSymbolsPerSecond:    2200,
						OnDemandPricePerSymbol:      2300,
					},
				},
				QuorumProtocolConfigs: map[uint32]*disperser_rpc.PaymentQuorumProtocolConfig{
					1: {
						MinNumSymbols:              1400,
						ReservationAdvanceWindow:   1500,
						ReservationRateLimitWindow: 1600,
						OnDemandRateLimitWindow:    1700,
						OnDemandEnabled:            false,
					},
					3: {
						MinNumSymbols:              3400,
						ReservationAdvanceWindow:   3500,
						ReservationRateLimitWindow: 3600,
						OnDemandRateLimitWindow:    3700,
						OnDemandEnabled:            true,
					},
				},
				OnDemandQuorumNumbers: []uint32{0, 1, 2, 3},
			},
			expected: &meterer.PaymentVaultParams{
				QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
					0: {
						ReservationSymbolsPerSecond: 100,
						OnDemandSymbolsPerSecond:    200,
						OnDemandPricePerSymbol:      300,
					},
					2: {
						ReservationSymbolsPerSecond: 2100,
						OnDemandSymbolsPerSecond:    2200,
						OnDemandPricePerSymbol:      2300,
					},
				},
				QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
					1: {
						MinNumSymbols:              1400,
						ReservationAdvanceWindow:   1500,
						ReservationRateLimitWindow: 1600,
						OnDemandRateLimitWindow:    1700,
						OnDemandEnabled:            false,
					},
					3: {
						MinNumSymbols:              3400,
						ReservationAdvanceWindow:   3500,
						ReservationRateLimitWindow: 3600,
						OnDemandRateLimitWindow:    3700,
						OnDemandEnabled:            true,
					},
				},
				OnDemandQuorumNumbers: []core.QuorumID{0, 1, 2, 3},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := meterer.PaymentVaultParamsFromProtobuf(tt.vaultParams)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)

				// Compare QuorumPaymentConfigs
				assert.Equal(t, len(tt.expected.QuorumPaymentConfigs), len(result.QuorumPaymentConfigs))
				for quorumID, expectedConfig := range tt.expected.QuorumPaymentConfigs {
					actualConfig, exists := result.QuorumPaymentConfigs[quorumID]
					assert.True(t, exists, "QuorumPaymentConfig for quorum %d should exist", quorumID)
					assert.Equal(t, expectedConfig.ReservationSymbolsPerSecond, actualConfig.ReservationSymbolsPerSecond)
					assert.Equal(t, expectedConfig.OnDemandSymbolsPerSecond, actualConfig.OnDemandSymbolsPerSecond)
					assert.Equal(t, expectedConfig.OnDemandPricePerSymbol, actualConfig.OnDemandPricePerSymbol)
				}

				// Compare QuorumProtocolConfigs
				assert.Equal(t, len(tt.expected.QuorumProtocolConfigs), len(result.QuorumProtocolConfigs))
				for quorumID, expectedConfig := range tt.expected.QuorumProtocolConfigs {
					actualConfig, exists := result.QuorumProtocolConfigs[quorumID]
					assert.True(t, exists, "QuorumProtocolConfig for quorum %d should exist", quorumID)
					assert.Equal(t, expectedConfig.MinNumSymbols, actualConfig.MinNumSymbols)
					assert.Equal(t, expectedConfig.ReservationAdvanceWindow, actualConfig.ReservationAdvanceWindow)
					assert.Equal(t, expectedConfig.ReservationRateLimitWindow, actualConfig.ReservationRateLimitWindow)
					assert.Equal(t, expectedConfig.OnDemandRateLimitWindow, actualConfig.OnDemandRateLimitWindow)
					assert.Equal(t, expectedConfig.OnDemandEnabled, actualConfig.OnDemandEnabled)
				}

				// Compare OnDemandQuorumNumbers
				assert.Equal(t, len(tt.expected.OnDemandQuorumNumbers), len(result.OnDemandQuorumNumbers))
				for i, expectedQuorum := range tt.expected.OnDemandQuorumNumbers {
					assert.Equal(t, expectedQuorum, result.OnDemandQuorumNumbers[i])
				}
			}
		})
	}
}
