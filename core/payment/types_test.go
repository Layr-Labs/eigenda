package payment_test

import (
	"math"
	"math/big"
	"testing"
	"time"

	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/payment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConversionPaymentVaultParams(t *testing.T) {
	tests := []struct {
		name        string
		input       *payment.PaymentVaultParams
		expectedErr string
		validate    func(t *testing.T, pbParams *disperser_rpc.PaymentVaultParams, coreParams *payment.PaymentVaultParams)
	}{
		{
			name: "valid complete params",
			input: &payment.PaymentVaultParams{
				QuorumPaymentConfigs: map[core.QuorumID]*payment.PaymentQuorumConfig{
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
				QuorumProtocolConfigs: map[core.QuorumID]*payment.PaymentQuorumProtocolConfig{
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
			validate: func(t *testing.T, pbParams *disperser_rpc.PaymentVaultParams, coreParams *payment.PaymentVaultParams) {
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
				convertedCoreParams, err := payment.PaymentVaultParamsFromProtobuf(pbParams)
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
			input: &payment.PaymentVaultParams{
				QuorumProtocolConfigs: map[core.QuorumID]*payment.PaymentQuorumProtocolConfig{
					1: {MinNumSymbols: 10},
				},
				OnDemandQuorumNumbers: []core.QuorumID{1},
			},
			expectedErr: "payment quorum configs cannot be nil",
		},
		{
			name: "nil protocol configs",
			input: &payment.PaymentVaultParams{
				QuorumPaymentConfigs: map[core.QuorumID]*payment.PaymentQuorumConfig{
					1: {ReservationSymbolsPerSecond: 100},
				},
				OnDemandQuorumNumbers: []core.QuorumID{1},
			},
			expectedErr: "payment quorum protocol configs cannot be nil",
		},
		{
			name: "empty maps",
			input: &payment.PaymentVaultParams{
				QuorumPaymentConfigs:  map[core.QuorumID]*payment.PaymentQuorumConfig{},
				QuorumProtocolConfigs: map[core.QuorumID]*payment.PaymentQuorumProtocolConfig{},
				OnDemandQuorumNumbers: []core.QuorumID{},
			},
			validate: func(t *testing.T, pbParams *disperser_rpc.PaymentVaultParams, coreParams *payment.PaymentVaultParams) {
				require.Empty(t, pbParams.QuorumPaymentConfigs)
				require.Empty(t, pbParams.QuorumProtocolConfigs)
				require.Empty(t, pbParams.OnDemandQuorumNumbers)
			},
		},
		{
			name: "max uint32 quorum ID",
			input: &payment.PaymentVaultParams{
				QuorumPaymentConfigs: map[core.QuorumID]*payment.PaymentQuorumConfig{
					core.QuorumID(255): {ReservationSymbolsPerSecond: 100},
				},
				QuorumProtocolConfigs: map[core.QuorumID]*payment.PaymentQuorumProtocolConfig{
					core.QuorumID(255): {MinNumSymbols: 10},
				},
				OnDemandQuorumNumbers: []core.QuorumID{core.QuorumID(255)},
			},
			validate: func(t *testing.T, pbParams *disperser_rpc.PaymentVaultParams, coreParams *payment.PaymentVaultParams) {
				require.Equal(t, uint64(100), pbParams.QuorumPaymentConfigs[255].ReservationSymbolsPerSecond)
				require.Equal(t, uint64(10), pbParams.QuorumProtocolConfigs[255].MinNumSymbols)
				require.Equal(t, uint32(255), pbParams.OnDemandQuorumNumbers[0])
			},
		},
		{
			name: "max uint64 values",
			input: &payment.PaymentVaultParams{
				QuorumPaymentConfigs: map[core.QuorumID]*payment.PaymentQuorumConfig{
					1: {
						ReservationSymbolsPerSecond: math.MaxUint64,
						OnDemandSymbolsPerSecond:    math.MaxUint64,
						OnDemandPricePerSymbol:      math.MaxUint64,
					},
				},
				QuorumProtocolConfigs: map[core.QuorumID]*payment.PaymentQuorumProtocolConfig{
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
			validate: func(t *testing.T, pbParams *disperser_rpc.PaymentVaultParams, coreParams *payment.PaymentVaultParams) {
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

func TestConversionPaymentVaultParamsFromProtobuf(t *testing.T) {
	tests := []struct {
		name        string
		input       *disperser_rpc.PaymentVaultParams
		expectedErr string
		validate    func(t *testing.T, coreParams *payment.PaymentVaultParams)
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
			validate: func(t *testing.T, coreParams *payment.PaymentVaultParams) {
				require.Empty(t, coreParams.QuorumPaymentConfigs)
				require.Empty(t, coreParams.QuorumProtocolConfigs)
				require.Empty(t, coreParams.OnDemandQuorumNumbers)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coreParams, err := payment.PaymentVaultParamsFromProtobuf(tt.input)
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

func TestConversionPaymentVaultParamsToProtobuf(t *testing.T) {
	tests := []struct {
		name        string
		input       *payment.PaymentVaultParams
		expectedErr string
		validate    func(t *testing.T, pbParams *disperser_rpc.PaymentVaultParams, coreParams *payment.PaymentVaultParams)
	}{
		{
			name: "valid complete params",
			input: &payment.PaymentVaultParams{
				QuorumPaymentConfigs: map[core.QuorumID]*payment.PaymentQuorumConfig{
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
				QuorumProtocolConfigs: map[core.QuorumID]*payment.PaymentQuorumProtocolConfig{
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
			validate: func(t *testing.T, pbParams *disperser_rpc.PaymentVaultParams, coreParams *payment.PaymentVaultParams) {
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
				convertedCoreParams, err := payment.PaymentVaultParamsFromProtobuf(pbParams)
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
			input: &payment.PaymentVaultParams{
				QuorumProtocolConfigs: map[core.QuorumID]*payment.PaymentQuorumProtocolConfig{
					1: {MinNumSymbols: 10},
				},
				OnDemandQuorumNumbers: []core.QuorumID{1},
			},
			expectedErr: "payment quorum configs cannot be nil",
		},
		{
			name: "nil protocol configs",
			input: &payment.PaymentVaultParams{
				QuorumPaymentConfigs: map[core.QuorumID]*payment.PaymentQuorumConfig{
					1: {ReservationSymbolsPerSecond: 100},
				},
				OnDemandQuorumNumbers: []core.QuorumID{1},
			},
			expectedErr: "payment quorum protocol configs cannot be nil",
		},
		{
			name: "empty maps",
			input: &payment.PaymentVaultParams{
				QuorumPaymentConfigs:  map[core.QuorumID]*payment.PaymentQuorumConfig{},
				QuorumProtocolConfigs: map[core.QuorumID]*payment.PaymentQuorumProtocolConfig{},
				OnDemandQuorumNumbers: []core.QuorumID{},
			},
			validate: func(t *testing.T, pbParams *disperser_rpc.PaymentVaultParams, coreParams *payment.PaymentVaultParams) {
				assert.Empty(t, pbParams.QuorumPaymentConfigs)
				assert.Empty(t, pbParams.QuorumProtocolConfigs)
				assert.Empty(t, pbParams.OnDemandQuorumNumbers)
			},
		},
		{
			name: "max uint32 quorum ID",
			input: &payment.PaymentVaultParams{
				QuorumPaymentConfigs: map[core.QuorumID]*payment.PaymentQuorumConfig{
					core.QuorumID(255): {ReservationSymbolsPerSecond: 100},
				},
				QuorumProtocolConfigs: map[core.QuorumID]*payment.PaymentQuorumProtocolConfig{
					core.QuorumID(255): {MinNumSymbols: 10},
				},
				OnDemandQuorumNumbers: []core.QuorumID{core.QuorumID(255)},
			},
			validate: func(t *testing.T, pbParams *disperser_rpc.PaymentVaultParams, coreParams *payment.PaymentVaultParams) {
				assert.Equal(t, uint64(100), pbParams.QuorumPaymentConfigs[255].ReservationSymbolsPerSecond)
				assert.Equal(t, uint64(10), pbParams.QuorumProtocolConfigs[255].MinNumSymbols)
				assert.Equal(t, uint32(255), pbParams.OnDemandQuorumNumbers[0])
			},
		},
		{
			name: "max uint64 values",
			input: &payment.PaymentVaultParams{
				QuorumPaymentConfigs: map[core.QuorumID]*payment.PaymentQuorumConfig{
					1: {
						ReservationSymbolsPerSecond: math.MaxUint64,
						OnDemandSymbolsPerSecond:    math.MaxUint64,
						OnDemandPricePerSymbol:      math.MaxUint64,
					},
				},
				QuorumProtocolConfigs: map[core.QuorumID]*payment.PaymentQuorumProtocolConfig{
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
			validate: func(t *testing.T, pbParams *disperser_rpc.PaymentVaultParams, coreParams *payment.PaymentVaultParams) {
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

func TestConversionReservationsFromProtobuf(t *testing.T) {
	tests := []struct {
		name     string
		input    map[uint32]*disperser_rpc.QuorumReservation
		expected map[core.QuorumID]*payment.ReservedPayment
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty map",
			input:    map[uint32]*disperser_rpc.QuorumReservation{},
			expected: map[core.QuorumID]*payment.ReservedPayment{},
		},
		{
			name: "single reservation",
			input: map[uint32]*disperser_rpc.QuorumReservation{
				0: {
					SymbolsPerSecond: 100,
					StartTimestamp:   1000,
					EndTimestamp:     2000,
				},
			},
			expected: map[core.QuorumID]*payment.ReservedPayment{
				0: {
					SymbolsPerSecond: 100,
					StartTimestamp:   1000,
					EndTimestamp:     2000,
				},
			},
		},
		{
			name: "multiple reservations with max values",
			input: map[uint32]*disperser_rpc.QuorumReservation{
				0: {
					SymbolsPerSecond: 100,
					StartTimestamp:   1000,
					EndTimestamp:     2000,
				},
				255: {
					SymbolsPerSecond: math.MaxUint64,
					StartTimestamp:   math.MaxUint32,
					EndTimestamp:     math.MaxUint32,
				},
			},
			expected: map[core.QuorumID]*payment.ReservedPayment{
				0: {
					SymbolsPerSecond: 100,
					StartTimestamp:   1000,
					EndTimestamp:     2000,
				},
				255: {
					SymbolsPerSecond: math.MaxUint64,
					StartTimestamp:   math.MaxUint32,
					EndTimestamp:     math.MaxUint32,
				},
			},
		},
		{
			name: "reservation with nil entry",
			input: map[uint32]*disperser_rpc.QuorumReservation{
				0: {
					SymbolsPerSecond: 100,
					StartTimestamp:   1000,
					EndTimestamp:     2000,
				},
				1: nil, // This should be skipped
			},
			expected: map[core.QuorumID]*payment.ReservedPayment{
				0: {
					SymbolsPerSecond: 100,
					StartTimestamp:   1000,
					EndTimestamp:     2000,
				},
			},
		},
		{
			name: "zero values",
			input: map[uint32]*disperser_rpc.QuorumReservation{
				0: {
					SymbolsPerSecond: 0,
					StartTimestamp:   0,
					EndTimestamp:     0,
				},
			},
			expected: map[core.QuorumID]*payment.ReservedPayment{
				0: {
					SymbolsPerSecond: 0,
					StartTimestamp:   0,
					EndTimestamp:     0,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := payment.ReservationsFromProtobuf(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConversionCumulativePaymentFromProtobuf(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected *big.Int
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty bytes",
			input:    []byte{},
			expected: big.NewInt(0),
		},
		{
			name:     "small positive value",
			input:    big.NewInt(123).Bytes(),
			expected: big.NewInt(123),
		},
		{
			name:     "large positive value",
			input:    new(big.Int).SetUint64(math.MaxUint64).Bytes(),
			expected: new(big.Int).SetUint64(math.MaxUint64),
		},
		{
			name: "very large value",
			input: func() []byte {
				val, _ := new(big.Int).SetString("123456789012345678901234567890", 10)
				return val.Bytes()
			}(),
			expected: func() *big.Int {
				val, _ := new(big.Int).SetString("123456789012345678901234567890", 10)
				return val
			}(),
		},
		{
			name:     "multi-byte value",
			input:    []byte{0x01, 0x00, 0x00},
			expected: big.NewInt(65536), // 2^16
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := payment.CumulativePaymentFromProtobuf(tt.input)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				require.NotNil(t, result)
				assert.Equal(t, 0, tt.expected.Cmp(result), "expected %s, got %s", tt.expected.String(), result.String())
			}
		})
	}
}

func TestConversionPaymentStateFromProtobuf(t *testing.T) {
	tests := []struct {
		name        string
		input       *disperser_rpc.GetPaymentStateForAllQuorumsReply
		expectedErr string
		validate    func(t *testing.T, vaultParams *payment.PaymentVaultParams, reservations map[core.QuorumID]*payment.ReservedPayment, cumulative, onchainCumulative *big.Int, periodRecords payment.QuorumPeriodRecords)
	}{
		{
			name:        "nil input",
			input:       nil,
			expectedErr: "payment state cannot be nil",
		},
		{
			name: "invalid payment vault params",
			input: &disperser_rpc.GetPaymentStateForAllQuorumsReply{
				PaymentVaultParams: nil,
			},
			expectedErr: "error converting payment vault params",
		},
		{
			name: "complete valid input",
			input: &disperser_rpc.GetPaymentStateForAllQuorumsReply{
				PaymentVaultParams: &disperser_rpc.PaymentVaultParams{
					QuorumPaymentConfigs: map[uint32]*disperser_rpc.PaymentQuorumConfig{
						0: {
							ReservationSymbolsPerSecond: 100,
							OnDemandSymbolsPerSecond:    200,
							OnDemandPricePerSymbol:      300,
						},
					},
					QuorumProtocolConfigs: map[uint32]*disperser_rpc.PaymentQuorumProtocolConfig{
						0: {
							MinNumSymbols:              10,
							ReservationAdvanceWindow:   20,
							ReservationRateLimitWindow: 30,
							OnDemandRateLimitWindow:    40,
							OnDemandEnabled:            true,
						},
					},
					OnDemandQuorumNumbers: []uint32{0},
				},
				Reservations: map[uint32]*disperser_rpc.QuorumReservation{
					0: {
						SymbolsPerSecond: 500,
						StartTimestamp:   1000,
						EndTimestamp:     2000,
					},
				},
				CumulativePayment:        big.NewInt(12345).Bytes(),
				OnchainCumulativePayment: big.NewInt(67890).Bytes(),
			},
			validate: func(t *testing.T, vaultParams *payment.PaymentVaultParams, reservations map[core.QuorumID]*payment.ReservedPayment, cumulative, onchainCumulative *big.Int, periodRecords payment.QuorumPeriodRecords) {
				require.NotNil(t, vaultParams)
				assert.Len(t, vaultParams.QuorumPaymentConfigs, 1)
				assert.Equal(t, uint64(100), vaultParams.QuorumPaymentConfigs[0].ReservationSymbolsPerSecond)

				require.NotNil(t, reservations)
				assert.Len(t, reservations, 1)
				assert.Equal(t, uint64(500), reservations[0].SymbolsPerSecond)

				require.NotNil(t, cumulative)
				assert.Equal(t, 0, big.NewInt(12345).Cmp(cumulative))
				require.NotNil(t, onchainCumulative)
				assert.Equal(t, 0, big.NewInt(67890).Cmp(onchainCumulative))
			},
		},
		{
			name: "minimal valid input",
			input: &disperser_rpc.GetPaymentStateForAllQuorumsReply{
				PaymentVaultParams: &disperser_rpc.PaymentVaultParams{
					QuorumPaymentConfigs:  map[uint32]*disperser_rpc.PaymentQuorumConfig{},
					QuorumProtocolConfigs: map[uint32]*disperser_rpc.PaymentQuorumProtocolConfig{},
					OnDemandQuorumNumbers: []uint32{},
				},
			},
			validate: func(t *testing.T, vaultParams *payment.PaymentVaultParams, reservations map[core.QuorumID]*payment.ReservedPayment, cumulative, onchainCumulative *big.Int, periodRecords payment.QuorumPeriodRecords) {
				require.NotNil(t, vaultParams)
				assert.Empty(t, vaultParams.QuorumPaymentConfigs)
				assert.Nil(t, reservations)
				assert.Nil(t, cumulative)
				assert.Nil(t, onchainCumulative)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vaultParams, reservations, cumulative, onchainCumulative, periodRecords, err := payment.ConvertPaymentStateFromProtobuf(tt.input)

			if tt.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErr)
				return
			}

			require.NoError(t, err)
			if tt.validate != nil {
				tt.validate(t, vaultParams, reservations, cumulative, onchainCumulative, periodRecords)
			}
		})
	}
}

func TestWithinTime(t *testing.T) {
	tests := []struct {
		name        string
		start       time.Time
		end         time.Time
		checkTime   time.Time
		wantInRange bool
	}{
		{
			name:        "in range - current time in middle of range",
			start:       time.Unix(100, 0),
			end:         time.Unix(200, 0),
			checkTime:   time.Unix(150, 0),
			wantInRange: true,
		},
		{
			name:        "in range - current time at start",
			start:       time.Unix(100, 0),
			end:         time.Unix(200, 0),
			checkTime:   time.Unix(100, 0),
			wantInRange: true,
		},
		{
			name:        "in range - current time at end",
			start:       time.Unix(100, 0),
			end:         time.Unix(200, 0),
			checkTime:   time.Unix(200, 0),
			wantInRange: true,
		},
		{
			name:        "out of range - current time before start",
			start:       time.Unix(100, 0),
			end:         time.Unix(200, 0),
			checkTime:   time.Unix(99, 0),
			wantInRange: false,
		},
		{
			name:        "out of range - current time after end",
			start:       time.Unix(100, 0),
			end:         time.Unix(200, 0),
			checkTime:   time.Unix(201, 0),
			wantInRange: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isInRange := payment.WithinTime(tt.checkTime, tt.start, tt.end)
			assert.Equal(t, tt.wantInRange, isInRange)
		})
	}
}
