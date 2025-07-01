package payment_logic_test

import (
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer/payment_logic"
	"github.com/stretchr/testify/assert"
)

func TestGetBinLimit(t *testing.T) {
	tests := []struct {
		name             string
		symbolsPerSecond uint64
		binInterval      uint64
		expected         uint64
	}{
		{
			name:             "normal case",
			symbolsPerSecond: 100,
			binInterval:      60,
			expected:         6000,
		},
		{
			name:             "zero symbols per second",
			symbolsPerSecond: 0,
			binInterval:      60,
			expected:         0,
		},
		{
			name:             "zero bin interval",
			symbolsPerSecond: 100,
			binInterval:      0,
			expected:         0,
		},
		{
			name:             "both zero",
			symbolsPerSecond: 0,
			binInterval:      0,
			expected:         0,
		},
		{
			name:             "large values without overflow",
			symbolsPerSecond: math.MaxUint32,
			binInterval:      100,
			expected:         math.MaxUint32 * 100,
		},
		{
			name:             "overflow case - should handle gracefully",
			symbolsPerSecond: math.MaxUint64,
			binInterval:      2,
			expected:         math.MaxUint64, // Overflow protection returns max value
		},
		{
			name:             "edge case - near overflow boundary",
			symbolsPerSecond: math.MaxUint64 / 2,
			binInterval:      2,
			expected:         math.MaxUint64 - 1, // Just under max value
		},
		{
			name:             "overflow case - large factors",
			symbolsPerSecond: math.MaxUint32 + 1,
			binInterval:      math.MaxUint32 + 1,
			expected:         math.MaxUint64, // Would overflow
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := payment_logic.GetBinLimit(tt.symbolsPerSecond, tt.binInterval)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetReservationPeriod(t *testing.T) {
	tests := []struct {
		name        string
		timestamp   int64
		binInterval uint64
		expected    uint64
	}{
		{
			name:        "normal case - start of period",
			timestamp:   1800, // exactly divisible by 60
			binInterval: 60,
			expected:    1800,
		},
		{
			name:        "normal case - middle of period",
			timestamp:   1830, // 30 seconds into period
			binInterval: 60,
			expected:    1800, // should round down
		},
		{
			name:        "zero timestamp",
			timestamp:   0,
			binInterval: 60,
			expected:    0,
		},
		{
			name:        "zero bin interval",
			timestamp:   1000,
			binInterval: 0,
			expected:    0,
		},
		{
			name:        "small timestamp",
			timestamp:   30,
			binInterval: 60,
			expected:    0, // rounds down to 0
		},
		{
			name:        "large timestamp",
			timestamp:   math.MaxInt32,
			binInterval: 1000,
			expected:    uint64(math.MaxInt32/1000) * 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := payment_logic.GetReservationPeriod(tt.timestamp, tt.binInterval)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetReservationPeriodByNanosecond(t *testing.T) {
	tests := []struct {
		name                string
		nanosecondTimestamp int64
		binInterval         uint64
		expected            uint64
	}{
		{
			name:                "negative timestamp",
			nanosecondTimestamp: -1000,
			binInterval:         60,
			expected:            0,
		},
		{
			name:                "zero timestamp",
			nanosecondTimestamp: 0,
			binInterval:         60,
			expected:            0,
		},
		{
			name:                "normal timestamp - 1 second",
			nanosecondTimestamp: 1 * time.Second.Nanoseconds(),
			binInterval:         60,
			expected:            0, // 1 second rounds down to 0 in 60-second bins
		},
		{
			name:                "normal timestamp - 2 minutes",
			nanosecondTimestamp: 120 * time.Second.Nanoseconds(),
			binInterval:         60,
			expected:            120,
		},
		{
			name:                "timestamp between periods",
			nanosecondTimestamp: 90 * time.Second.Nanoseconds(), // 1.5 minutes
			binInterval:         60,
			expected:            60, // rounds down to 60
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := payment_logic.GetReservationPeriodByNanosecond(tt.nanosecondTimestamp, tt.binInterval)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetOverflowPeriod(t *testing.T) {
	tests := []struct {
		name              string
		reservationPeriod uint64
		reservationWindow uint64
		expected          uint64
	}{
		{
			name:              "normal case",
			reservationPeriod: 1800,
			reservationWindow: 60,
			expected:          1920, // 1800 + 2*60
		},
		{
			name:              "zero reservation period",
			reservationPeriod: 0,
			reservationWindow: 60,
			expected:          120, // 0 + 2*60
		},
		{
			name:              "zero reservation window",
			reservationPeriod: 1800,
			reservationWindow: 0,
			expected:          1800, // 1800 + 2*0
		},
		{
			name:              "large values",
			reservationPeriod: math.MaxUint32,
			reservationWindow: 1000,
			expected:          math.MaxUint32 + 2000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := payment_logic.GetOverflowPeriod(tt.reservationPeriod, tt.reservationWindow)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSymbolsCharged(t *testing.T) {
	tests := []struct {
		name       string
		numSymbols uint64
		minSymbols uint64
		expected   uint64
	}{
		{
			name:       "symbols less than minimum",
			numSymbols: 5,
			minSymbols: 10,
			expected:   10,
		},
		{
			name:       "symbols equal to minimum",
			numSymbols: 10,
			minSymbols: 10,
			expected:   10,
		},
		{
			name:       "symbols greater than minimum - exact multiple",
			numSymbols: 20,
			minSymbols: 10,
			expected:   20,
		},
		{
			name:       "symbols greater than minimum - round up",
			numSymbols: 15,
			minSymbols: 10,
			expected:   20, // rounds up to next multiple of 10
		},
		{
			name:       "zero minimum symbols",
			numSymbols: 15,
			minSymbols: 0,
			expected:   15, // returns numSymbols when minSymbols is 0
		},
		{
			name:       "zero symbols with minimum",
			numSymbols: 0,
			minSymbols: 10,
			expected:   10,
		},
		{
			name:       "large numbers requiring rounding",
			numSymbols: 1025,
			minSymbols: 1024,
			expected:   2048, // rounds up to next multiple of 1024
		},
		{
			name:       "overflow protection case",
			numSymbols: math.MaxUint64 - 100,
			minSymbols: 1000,
			expected:   math.MaxUint64, // should return MaxUint64 on overflow
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := payment_logic.SymbolsCharged(tt.numSymbols, tt.minSymbols)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPaymentCharged(t *testing.T) {
	tests := []struct {
		name           string
		numSymbols     uint64
		pricePerSymbol uint64
		expected       *big.Int
	}{
		{
			name:           "normal case",
			numSymbols:     100,
			pricePerSymbol: 5,
			expected:       big.NewInt(500),
		},
		{
			name:           "zero symbols",
			numSymbols:     0,
			pricePerSymbol: 5,
			expected:       big.NewInt(0),
		},
		{
			name:           "zero price",
			numSymbols:     100,
			pricePerSymbol: 0,
			expected:       big.NewInt(0),
		},
		{
			name:           "both zero",
			numSymbols:     0,
			pricePerSymbol: 0,
			expected:       big.NewInt(0),
		},
		{
			name:           "large numbers",
			numSymbols:     1000000,
			pricePerSymbol: 999999,
			expected:       new(big.Int).Mul(big.NewInt(1000000), big.NewInt(999999)),
		},
		{
			name:           "max uint64 values",
			numSymbols:     math.MaxUint64,
			pricePerSymbol: 1,
			expected:       new(big.Int).SetUint64(math.MaxUint64),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := payment_logic.PaymentCharged(tt.numSymbols, tt.pricePerSymbol)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateQuorum(t *testing.T) {
	tests := []struct {
		name           string
		headerQuorums  []uint8
		allowedQuorums []uint8
		expectError    bool
		errorContains  string
	}{
		{
			name:           "valid quorums - all allowed",
			headerQuorums:  []uint8{0, 1},
			allowedQuorums: []uint8{0, 1, 2},
			expectError:    false,
		},
		{
			name:           "valid quorums - subset",
			headerQuorums:  []uint8{1},
			allowedQuorums: []uint8{0, 1, 2},
			expectError:    false,
		},
		{
			name:           "empty header quorums",
			headerQuorums:  []uint8{},
			allowedQuorums: []uint8{0, 1},
			expectError:    true,
			errorContains:  "no quorum numbers provided",
		},
		{
			name:           "invalid quorum - not in allowed list",
			headerQuorums:  []uint8{0, 3},
			allowedQuorums: []uint8{0, 1, 2},
			expectError:    true,
			errorContains:  "quorum number mismatch: 3",
		},
		{
			name:           "multiple invalid quorums",
			headerQuorums:  []uint8{0, 3, 4},
			allowedQuorums: []uint8{0, 1, 2},
			expectError:    true,
			errorContains:  "quorum number mismatch: 3", // should fail on first invalid
		},
		{
			name:           "empty allowed quorums",
			headerQuorums:  []uint8{0},
			allowedQuorums: []uint8{},
			expectError:    true,
			errorContains:  "quorum number mismatch: 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := payment_logic.ValidateQuorum(tt.headerQuorums, tt.allowedQuorums)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateReservationPeriod(t *testing.T) {
	now := time.Now()
	nowNano := now.UnixNano()
	reservationWindow := uint64(60) // 60 seconds

	tests := []struct {
		name                string
		reservation         *core.ReservedPayment
		requestPeriod       uint64
		reservationWindow   uint64
		receivedTimestampNs int64
		expected            bool
	}{
		{
			name: "valid - current period within reservation",
			reservation: &core.ReservedPayment{
				StartTimestamp: uint64(now.Add(-time.Hour).Unix()),
				EndTimestamp:   uint64(now.Add(time.Hour).Unix()),
			},
			requestPeriod:       payment_logic.GetReservationPeriodByNanosecond(nowNano, reservationWindow),
			reservationWindow:   reservationWindow,
			receivedTimestampNs: nowNano,
			expected:            true,
		},
		{
			name: "valid - previous period within reservation",
			reservation: &core.ReservedPayment{
				StartTimestamp: uint64(now.Add(-time.Hour).Unix()),
				EndTimestamp:   uint64(now.Add(time.Hour).Unix()),
			},
			requestPeriod:       payment_logic.GetReservationPeriodByNanosecond(nowNano, reservationWindow) - reservationWindow,
			reservationWindow:   reservationWindow,
			receivedTimestampNs: nowNano,
			expected:            true,
		},
		{
			name: "invalid - future period",
			reservation: &core.ReservedPayment{
				StartTimestamp: uint64(now.Add(-time.Hour).Unix()),
				EndTimestamp:   uint64(now.Add(time.Hour).Unix()),
			},
			requestPeriod:       payment_logic.GetReservationPeriodByNanosecond(nowNano, reservationWindow) + reservationWindow,
			reservationWindow:   reservationWindow,
			receivedTimestampNs: nowNano,
			expected:            false,
		},
		{
			name: "invalid - too old period",
			reservation: &core.ReservedPayment{
				StartTimestamp: uint64(now.Add(-time.Hour).Unix()),
				EndTimestamp:   uint64(now.Add(time.Hour).Unix()),
			},
			requestPeriod:       payment_logic.GetReservationPeriodByNanosecond(nowNano, reservationWindow) - 2*reservationWindow,
			reservationWindow:   reservationWindow,
			receivedTimestampNs: nowNano,
			expected:            false,
		},
		{
			name: "invalid - before reservation start",
			reservation: &core.ReservedPayment{
				StartTimestamp: uint64(now.Add(time.Hour).Unix()), // starts in future
				EndTimestamp:   uint64(now.Add(2 * time.Hour).Unix()),
			},
			requestPeriod:       payment_logic.GetReservationPeriodByNanosecond(nowNano, reservationWindow),
			reservationWindow:   reservationWindow,
			receivedTimestampNs: nowNano,
			expected:            false,
		},
		{
			name: "invalid - after reservation end",
			reservation: &core.ReservedPayment{
				StartTimestamp: uint64(now.Add(-2 * time.Hour).Unix()),
				EndTimestamp:   uint64(now.Add(-time.Hour).Unix()), // ended in past
			},
			requestPeriod:       payment_logic.GetReservationPeriodByNanosecond(nowNano, reservationWindow),
			reservationWindow:   reservationWindow,
			receivedTimestampNs: nowNano,
			expected:            false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := payment_logic.ValidateReservationPeriod(tt.reservation, tt.requestPeriod, tt.reservationWindow, tt.receivedTimestampNs)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsOnDemandPayment(t *testing.T) {
	tests := []struct {
		name            string
		paymentMetadata *core.PaymentMetadata
		expected        bool
	}{
		{
			name: "on-demand payment - positive amount",
			paymentMetadata: &core.PaymentMetadata{
				CumulativePayment: big.NewInt(100),
			},
			expected: true,
		},
		{
			name: "not on-demand - zero payment",
			paymentMetadata: &core.PaymentMetadata{
				CumulativePayment: big.NewInt(0),
			},
			expected: false,
		},
		{
			name: "not on-demand - negative payment",
			paymentMetadata: &core.PaymentMetadata{
				CumulativePayment: big.NewInt(-1),
			},
			expected: false,
		},
		{
			name: "on-demand payment - large amount",
			paymentMetadata: &core.PaymentMetadata{
				CumulativePayment: new(big.Int).Mul(big.NewInt(1000000), big.NewInt(1000000)),
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := payment_logic.IsOnDemandPayment(tt.paymentMetadata)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateReservations(t *testing.T) {
	now := time.Now()
	nowNano := now.UnixNano()

	// Helper to create valid reservation
	validReservation := &core.ReservedPayment{
		SymbolsPerSecond: 100,
		StartTimestamp:   uint64(now.Add(-time.Hour).Unix()),
		EndTimestamp:     uint64(now.Add(time.Hour).Unix()),
	}

	// Helper to create expired reservation
	expiredReservation := &core.ReservedPayment{
		SymbolsPerSecond: 100,
		StartTimestamp:   uint64(now.Add(-2 * time.Hour).Unix()),
		EndTimestamp:     uint64(now.Add(-time.Hour).Unix()),
	}

	validQuorumConfigs := map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
		0: {
			ReservationRateLimitWindow: 60,
			MinNumSymbols:              1,
		},
		1: {
			ReservationRateLimitWindow: 60,
			MinNumSymbols:              1,
		},
	}

	tests := []struct {
		name                     string
		reservations             map[core.QuorumID]*core.ReservedPayment
		quorumConfigs            map[core.QuorumID]*core.PaymentQuorumProtocolConfig
		quorumNumbers            []uint8
		paymentHeaderTimestampNs int64
		receivedTimestampNs      int64
		expectError              bool
		errorContains            string
	}{
		{
			name: "valid reservations",
			reservations: map[core.QuorumID]*core.ReservedPayment{
				0: validReservation,
				1: validReservation,
			},
			quorumConfigs:            validQuorumConfigs,
			quorumNumbers:            []uint8{0, 1},
			paymentHeaderTimestampNs: nowNano,
			receivedTimestampNs:      nowNano,
			expectError:              false,
		},
		{
			name: "missing quorum config",
			reservations: map[core.QuorumID]*core.ReservedPayment{
				0: validReservation,
				2: validReservation, // quorum 2 not in configs
			},
			quorumConfigs:            validQuorumConfigs,
			quorumNumbers:            []uint8{0, 2},
			paymentHeaderTimestampNs: nowNano,
			receivedTimestampNs:      nowNano,
			expectError:              true,
			errorContains:            "quorum config not found for quorum 2",
		},
		{
			name: "invalid quorum in header",
			reservations: map[core.QuorumID]*core.ReservedPayment{
				0: validReservation,
				1: validReservation,
			},
			quorumConfigs:            validQuorumConfigs,
			quorumNumbers:            []uint8{0, 3}, // quorum 3 not in reservations
			paymentHeaderTimestampNs: nowNano,
			receivedTimestampNs:      nowNano,
			expectError:              true,
			errorContains:            "quorum number mismatch: 3",
		},
		{
			name: "inactive reservation",
			reservations: map[core.QuorumID]*core.ReservedPayment{
				0: expiredReservation,
			},
			quorumConfigs:            validQuorumConfigs,
			quorumNumbers:            []uint8{0},
			paymentHeaderTimestampNs: nowNano,
			receivedTimestampNs:      nowNano,
			expectError:              true,
			errorContains:            "reservation not active",
		},
		{
			name: "invalid reservation period",
			reservations: map[core.QuorumID]*core.ReservedPayment{
				0: {
					SymbolsPerSecond: 100,
					StartTimestamp:   uint64(now.Add(-2 * time.Hour).Unix()),
					EndTimestamp:     uint64(now.Add(2 * time.Hour).Unix()), // active but wrong period
				},
			},
			quorumConfigs:            validQuorumConfigs,
			quorumNumbers:            []uint8{0},
			paymentHeaderTimestampNs: now.Add(-10 * time.Minute).UnixNano(), // too far in past for valid period
			receivedTimestampNs:      nowNano,
			expectError:              true,
			errorContains:            "invalid reservation period",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := payment_logic.ValidateReservations(
				tt.reservations,
				tt.quorumConfigs,
				tt.quorumNumbers,
				tt.paymentHeaderTimestampNs,
				tt.receivedTimestampNs,
			)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
