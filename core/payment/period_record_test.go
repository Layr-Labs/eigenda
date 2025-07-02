package payment_test

import (
	"testing"

	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/payment"
	"github.com/stretchr/testify/assert"
)

func TestQuorumPeriodRecords_GetRelativePeriodRecord(t *testing.T) {
	tests := []struct {
		name               string
		initialRecords     payment.QuorumPeriodRecords
		index              uint64
		quorumNumber       core.QuorumID
		expectedIndex      uint32
		expectedUsage      uint64
		shouldCreateQuorum bool
		shouldCreateRecord bool
	}{
		{
			name:               "new quorum and record",
			initialRecords:     make(payment.QuorumPeriodRecords),
			index:              5,
			quorumNumber:       core.QuorumID(1),
			expectedIndex:      5,
			expectedUsage:      0,
			shouldCreateQuorum: true,
			shouldCreateRecord: true,
		},
		{
			name: "existing quorum, new record",
			initialRecords: payment.QuorumPeriodRecords{
				core.QuorumID(1): make([]*payment.PeriodRecord, 3),
			},
			index:              7,
			quorumNumber:       core.QuorumID(1),
			expectedIndex:      7,
			expectedUsage:      0,
			shouldCreateQuorum: false,
			shouldCreateRecord: true,
		},
		{
			name: "existing quorum and record",
			initialRecords: payment.QuorumPeriodRecords{
				core.QuorumID(1): []*payment.PeriodRecord{
					nil,
					{Index: 4, Usage: 100},
					nil,
				},
			},
			index:              4,
			quorumNumber:       core.QuorumID(1),
			expectedIndex:      4,
			expectedUsage:      100,
			shouldCreateQuorum: false,
			shouldCreateRecord: false,
		},
		{
			name:               "index wraps around (modulo operation)",
			initialRecords:     make(payment.QuorumPeriodRecords),
			index:              10, // 10 % 3 = 1
			quorumNumber:       core.QuorumID(2),
			expectedIndex:      10,
			expectedUsage:      0,
			shouldCreateQuorum: true,
			shouldCreateRecord: true,
		},
		{
			name:               "zero index",
			initialRecords:     make(payment.QuorumPeriodRecords),
			index:              0,
			quorumNumber:       core.QuorumID(0),
			expectedIndex:      0,
			expectedUsage:      0,
			shouldCreateQuorum: true,
			shouldCreateRecord: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record := tt.initialRecords.GetRelativePeriodRecord(tt.index, tt.quorumNumber)

			assert.NotNil(t, record)
			assert.Equal(t, tt.expectedIndex, record.Index)
			assert.Equal(t, tt.expectedUsage, record.Usage)

			// Verify quorum exists after call
			_, quorumExists := tt.initialRecords[tt.quorumNumber]
			assert.True(t, quorumExists)

			// Verify record exists in expected position
			relativeIndex := uint32(tt.index % 3) // MinNumBins = 3
			assert.NotNil(t, tt.initialRecords[tt.quorumNumber][relativeIndex])
		})
	}
}

func TestQuorumPeriodRecords_UpdateUsage(t *testing.T) {
	tests := []struct {
		name                  string
		initialRecords        payment.QuorumPeriodRecords
		quorumNumber          core.QuorumID
		timestamp             int64
		numSymbols            uint64
		reservation           *payment.ReservedPayment
		protocolConfig        *payment.PaymentQuorumProtocolConfig
		expectedError         string
		expectedCurrentUsage  uint64
		expectedOverflowUsage uint64
		setupCurrentRecord    bool
		setupOverflowRecord   bool
		currentRecordUsage    uint64
		overflowRecordUsage   uint64
	}{
		{
			name:           "symbol usage exceeds bin limit",
			initialRecords: make(payment.QuorumPeriodRecords),
			quorumNumber:   core.QuorumID(1),
			timestamp:      1000000000000, // 1 second in nanoseconds
			numSymbols:     550,
			reservation: &payment.ReservedPayment{
				SymbolsPerSecond: 50, // This will create bin limit of 50 * 10 = 500
				StartTimestamp:   0,
				EndTimestamp:     2000,
			},
			protocolConfig: &payment.PaymentQuorumProtocolConfig{
				MinNumSymbols:              100, // min symbols is 100, so 550 -> 600
				ReservationRateLimitWindow: 10,
			},
			expectedError: "symbol usage exceeds bin limit",
		},
		{
			name:           "usage within bin limit",
			initialRecords: make(payment.QuorumPeriodRecords),
			quorumNumber:   core.QuorumID(1),
			timestamp:      1000000000000,
			numSymbols:     50,
			reservation: &payment.ReservedPayment{
				SymbolsPerSecond: 100,
				StartTimestamp:   0,
				EndTimestamp:     2000,
			},
			protocolConfig: &payment.PaymentQuorumProtocolConfig{
				MinNumSymbols:              10,
				ReservationRateLimitWindow: 10,
			},
			expectedCurrentUsage: 50,
		},
		{
			name:           "usage with minimum symbols applied",
			initialRecords: make(payment.QuorumPeriodRecords),
			quorumNumber:   core.QuorumID(1),
			timestamp:      1000000000000,
			numSymbols:     5, // Below min symbols
			reservation: &payment.ReservedPayment{
				SymbolsPerSecond: 100,
				StartTimestamp:   0,
				EndTimestamp:     2000,
			},
			protocolConfig: &payment.PaymentQuorumProtocolConfig{
				MinNumSymbols:              20, // Min symbols enforced
				ReservationRateLimitWindow: 10,
			},
			expectedCurrentUsage: 20, // Should use min symbols
		},
		{
			name:               "usage exceeds limit but overflow available",
			initialRecords:     make(payment.QuorumPeriodRecords),
			quorumNumber:       core.QuorumID(1),
			timestamp:          1000000000000,
			numSymbols:         80,
			setupCurrentRecord: true,
			currentRecordUsage: 30,
			reservation: &payment.ReservedPayment{
				SymbolsPerSecond: 10, // bin limit = 10 * 10 = 100
				StartTimestamp:   0,
				EndTimestamp:     2000,
			},
			protocolConfig: &payment.PaymentQuorumProtocolConfig{
				MinNumSymbols:              1,
				ReservationRateLimitWindow: 10,
			},
			expectedCurrentUsage:  100, // capped at bin limit
			expectedOverflowUsage: 10,  // 30 + 80 - 100 = 10
		},
		{
			name:               "current usage already at limit",
			initialRecords:     make(payment.QuorumPeriodRecords),
			quorumNumber:       core.QuorumID(1),
			timestamp:          1000000000000,
			numSymbols:         10,
			setupCurrentRecord: true,
			currentRecordUsage: 100,
			reservation: &payment.ReservedPayment{
				SymbolsPerSecond: 10, // bin limit = 100
				StartTimestamp:   0,
				EndTimestamp:     2000,
			},
			protocolConfig: &payment.PaymentQuorumProtocolConfig{
				MinNumSymbols:              1,
				ReservationRateLimitWindow: 10,
			},
			expectedError: "reservation limit exceeded for quorum 1",
		},
		{
			name:               "current usage exceeds limit",
			initialRecords:     make(payment.QuorumPeriodRecords),
			quorumNumber:       core.QuorumID(1),
			timestamp:          1000000000000,
			numSymbols:         10,
			setupCurrentRecord: true,
			currentRecordUsage: 150,
			reservation: &payment.ReservedPayment{
				SymbolsPerSecond: 10, // bin limit = 100
				StartTimestamp:   0,
				EndTimestamp:     2000,
			},
			protocolConfig: &payment.PaymentQuorumProtocolConfig{
				MinNumSymbols:              1,
				ReservationRateLimitWindow: 10,
			},
			expectedError: "reservation limit exceeded for quorum 1",
		},
		{
			name:                "overflow bin already in use",
			initialRecords:      make(payment.QuorumPeriodRecords),
			quorumNumber:        core.QuorumID(1),
			timestamp:           1000000000000,
			numSymbols:          80,
			setupCurrentRecord:  true,
			setupOverflowRecord: true,
			currentRecordUsage:  30,
			overflowRecordUsage: 50,
			reservation: &payment.ReservedPayment{
				SymbolsPerSecond: 10, // bin limit = 100
				StartTimestamp:   0,
				EndTimestamp:     2000,
			},
			protocolConfig: &payment.PaymentQuorumProtocolConfig{
				MinNumSymbols:              1,
				ReservationRateLimitWindow: 10,
			},
			expectedError: "reservation limit exceeded for quorum 1",
		},
		{
			name:           "exactly at bin limit",
			initialRecords: make(payment.QuorumPeriodRecords),
			quorumNumber:   core.QuorumID(1),
			timestamp:      1000000000000,
			numSymbols:     100,
			reservation: &payment.ReservedPayment{
				SymbolsPerSecond: 10, // bin limit = 100
				StartTimestamp:   0,
				EndTimestamp:     2000,
			},
			protocolConfig: &payment.PaymentQuorumProtocolConfig{
				MinNumSymbols:              1,
				ReservationRateLimitWindow: 10,
			},
			expectedCurrentUsage: 100,
		},
		{
			name:               "zero usage (enforces min symbols)",
			initialRecords:     make(payment.QuorumPeriodRecords),
			quorumNumber:       core.QuorumID(1),
			timestamp:          1000000000000,
			numSymbols:         0,
			setupCurrentRecord: true,
			currentRecordUsage: 50,
			reservation: &payment.ReservedPayment{
				SymbolsPerSecond: 100,
				StartTimestamp:   0,
				EndTimestamp:     2000,
			},
			protocolConfig: &payment.PaymentQuorumProtocolConfig{
				MinNumSymbols:              5, // Min symbols enforced even for 0 input
				ReservationRateLimitWindow: 10,
			},
			expectedCurrentUsage: 55, // 50 + 5 (min symbols)
		},
		{
			name:           "negative timestamp",
			initialRecords: make(payment.QuorumPeriodRecords),
			quorumNumber:   core.QuorumID(1),
			timestamp:      -1000000000000,
			numSymbols:     50,
			reservation: &payment.ReservedPayment{
				SymbolsPerSecond: 100,
				StartTimestamp:   0,
				EndTimestamp:     2000,
			},
			protocolConfig: &payment.PaymentQuorumProtocolConfig{
				MinNumSymbols:              1,
				ReservationRateLimitWindow: 10,
			},
			expectedCurrentUsage: 50, // Should handle negative timestamp gracefully
		},
		{
			name:           "large reservation window",
			initialRecords: make(payment.QuorumPeriodRecords),
			quorumNumber:   core.QuorumID(1),
			timestamp:      1000000000000,
			numSymbols:     1000,
			reservation: &payment.ReservedPayment{
				SymbolsPerSecond: 1,
				StartTimestamp:   0,
				EndTimestamp:     2000,
			},
			protocolConfig: &payment.PaymentQuorumProtocolConfig{
				MinNumSymbols:              1,
				ReservationRateLimitWindow: 10000, // Large window = large bin limit
			},
			expectedCurrentUsage: 1000,
		},
		{
			name:           "zero min symbols",
			initialRecords: make(payment.QuorumPeriodRecords),
			quorumNumber:   core.QuorumID(1),
			timestamp:      1000000000000,
			numSymbols:     50,
			reservation: &payment.ReservedPayment{
				SymbolsPerSecond: 100,
				StartTimestamp:   0,
				EndTimestamp:     2000,
			},
			protocolConfig: &payment.PaymentQuorumProtocolConfig{
				MinNumSymbols:              0, // Zero min symbols
				ReservationRateLimitWindow: 10,
			},
			expectedCurrentUsage: 50, // Should use actual symbols since min is 0
		},
		{
			name:           "zero symbols per second (zero bin limit)",
			initialRecords: make(payment.QuorumPeriodRecords),
			quorumNumber:   core.QuorumID(1),
			timestamp:      1000000000000,
			numSymbols:     1,
			reservation: &payment.ReservedPayment{
				SymbolsPerSecond: 0, // Zero symbols per second
				StartTimestamp:   0,
				EndTimestamp:     2000,
			},
			protocolConfig: &payment.PaymentQuorumProtocolConfig{
				MinNumSymbols:              1,
				ReservationRateLimitWindow: 10,
			},
			expectedError: "symbol usage exceeds bin limit", // bin limit would be 0
		},
		{
			name:           "different quorum numbers",
			initialRecords: make(payment.QuorumPeriodRecords),
			quorumNumber:   core.QuorumID(255), // Max quorum ID
			timestamp:      1000000000000,
			numSymbols:     50,
			reservation: &payment.ReservedPayment{
				SymbolsPerSecond: 100,
				StartTimestamp:   0,
				EndTimestamp:     2000,
			},
			protocolConfig: &payment.PaymentQuorumProtocolConfig{
				MinNumSymbols:              1,
				ReservationRateLimitWindow: 10,
			},
			expectedCurrentUsage: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Calculate expected periods for setup
			currentPeriod := payment.GetReservationPeriodByNanosecond(tt.timestamp, tt.protocolConfig.ReservationRateLimitWindow)
			overflowPeriod := payment.GetOverflowPeriod(currentPeriod, tt.protocolConfig.ReservationRateLimitWindow)

			// Setup initial records if needed
			if tt.setupCurrentRecord {
				currentRecord := tt.initialRecords.GetRelativePeriodRecord(currentPeriod, tt.quorumNumber)
				currentRecord.Usage = tt.currentRecordUsage
			}
			if tt.setupOverflowRecord {
				overflowRecord := tt.initialRecords.GetRelativePeriodRecord(overflowPeriod, tt.quorumNumber)
				overflowRecord.Usage = tt.overflowRecordUsage
			}

			err := tt.initialRecords.UpdateUsage(
				tt.quorumNumber,
				tt.timestamp,
				tt.numSymbols,
				tt.reservation,
				tt.protocolConfig,
			)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)

				// Check current record usage
				currentRecord := tt.initialRecords.GetRelativePeriodRecord(currentPeriod, tt.quorumNumber)
				assert.Equal(t, tt.expectedCurrentUsage, currentRecord.Usage)

				// Check overflow record usage if expected
				if tt.expectedOverflowUsage > 0 {
					overflowRecord := tt.initialRecords.GetRelativePeriodRecord(overflowPeriod, tt.quorumNumber)
					assert.Equal(t, tt.expectedOverflowUsage, overflowRecord.Usage)
				}
			}
		})
	}
}

func TestQuorumPeriodRecords_DeepCopy(t *testing.T) {
	tests := []struct {
		name            string
		originalRecords payment.QuorumPeriodRecords
	}{
		{
			name:            "empty records",
			originalRecords: make(payment.QuorumPeriodRecords),
		},
		{
			name: "single quorum with records",
			originalRecords: payment.QuorumPeriodRecords{
				core.QuorumID(1): []*payment.PeriodRecord{
					{Index: 0, Usage: 100},
					{Index: 1, Usage: 200},
					nil,
				},
			},
		},
		{
			name: "multiple quorums with mixed records",
			originalRecords: payment.QuorumPeriodRecords{
				core.QuorumID(1): []*payment.PeriodRecord{
					{Index: 0, Usage: 100},
					nil,
					{Index: 2, Usage: 300},
				},
				core.QuorumID(2): []*payment.PeriodRecord{
					nil,
					{Index: 4, Usage: 400},
					{Index: 5, Usage: 500},
				},
			},
		},
		{
			name: "quorum with all nil records",
			originalRecords: payment.QuorumPeriodRecords{
				core.QuorumID(3): []*payment.PeriodRecord{
					nil,
					nil,
					nil,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			copied := tt.originalRecords.DeepCopy()

			// Verify structure is copied
			assert.Equal(t, len(tt.originalRecords), len(copied))

			for quorumID, originalSlice := range tt.originalRecords {
				copiedSlice, exists := copied[quorumID]
				assert.True(t, exists)
				assert.Equal(t, len(originalSlice), len(copiedSlice))

				for i, originalRecord := range originalSlice {
					if originalRecord == nil {
						assert.Nil(t, copiedSlice[i])
					} else {
						assert.NotNil(t, copiedSlice[i])
						assert.Equal(t, originalRecord.Index, copiedSlice[i].Index)
						assert.Equal(t, originalRecord.Usage, copiedSlice[i].Usage)

						// Verify it's a deep copy (different memory addresses)
						assert.NotSame(t, originalRecord, copiedSlice[i])
					}
				}
			}

			// Verify modifying copy doesn't affect original
			if len(copied) > 0 {
				for quorumID, copiedSlice := range copied {
					for i, record := range copiedSlice {
						if record != nil {
							// Modify the copy
							record.Usage = 9999
							record.Index = 8888

							// Verify original is unchanged
							originalRecord := tt.originalRecords[quorumID][i]
							if originalRecord != nil {
								assert.NotEqual(t, 9999, originalRecord.Usage)
								assert.NotEqual(t, 8888, originalRecord.Index)
							}
							break
						}
					}
					break
				}
			}
		})
	}
}

func TestQuorumPeriodRecords_FromProtoRecords(t *testing.T) {
	tests := []struct {
		name         string
		protoRecords map[uint32]*disperser_rpc.PeriodRecords
		expected     payment.QuorumPeriodRecords
	}{
		{
			name:         "empty proto records",
			protoRecords: make(map[uint32]*disperser_rpc.PeriodRecords),
			expected:     make(payment.QuorumPeriodRecords),
		},
		{
			name: "single quorum with records",
			protoRecords: map[uint32]*disperser_rpc.PeriodRecords{
				1: {
					Records: []*disperser_rpc.PeriodRecord{
						{Index: 0, Usage: 100},
						{Index: 1, Usage: 200},
					},
				},
			},
			expected: payment.QuorumPeriodRecords{
				core.QuorumID(1): []*payment.PeriodRecord{
					{Index: 0, Usage: 100},
					{Index: 1, Usage: 200},
					{Index: 2, Usage: 0}, // Default initialized
				},
			},
		},
		{
			name: "multiple quorums",
			protoRecords: map[uint32]*disperser_rpc.PeriodRecords{
				1: {
					Records: []*disperser_rpc.PeriodRecord{
						{Index: 5, Usage: 500}, // 5 % 3 = 2
					},
				},
				2: {
					Records: []*disperser_rpc.PeriodRecord{
						{Index: 3, Usage: 300}, // 3 % 3 = 0
						{Index: 7, Usage: 700}, // 7 % 3 = 1
					},
				},
			},
			expected: payment.QuorumPeriodRecords{
				core.QuorumID(1): []*payment.PeriodRecord{
					{Index: 0, Usage: 0},   // Default
					{Index: 1, Usage: 0},   // Default
					{Index: 5, Usage: 500}, // Overwritten at index 2
				},
				core.QuorumID(2): []*payment.PeriodRecord{
					{Index: 3, Usage: 300}, // Overwritten at index 0
					{Index: 7, Usage: 700}, // Overwritten at index 1
					{Index: 2, Usage: 0},   // Default
				},
			},
		},
		{
			name: "index wrapping with modulo",
			protoRecords: map[uint32]*disperser_rpc.PeriodRecords{
				0: {
					Records: []*disperser_rpc.PeriodRecord{
						{Index: 10, Usage: 1000}, // 10 % 3 = 1
						{Index: 11, Usage: 1100}, // 11 % 3 = 2
						{Index: 12, Usage: 1200}, // 12 % 3 = 0
					},
				},
			},
			expected: payment.QuorumPeriodRecords{
				core.QuorumID(0): []*payment.PeriodRecord{
					{Index: 12, Usage: 1200}, // index 0
					{Index: 10, Usage: 1000}, // index 1
					{Index: 11, Usage: 1100}, // index 2
				},
			},
		},
		{
			name: "empty records for quorum",
			protoRecords: map[uint32]*disperser_rpc.PeriodRecords{
				1: {
					Records: []*disperser_rpc.PeriodRecord{},
				},
			},
			expected: payment.QuorumPeriodRecords{
				core.QuorumID(1): []*payment.PeriodRecord{
					{Index: 0, Usage: 0},
					{Index: 1, Usage: 0},
					{Index: 2, Usage: 0},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := payment.FromProtoRecords(tt.protoRecords)

			assert.Equal(t, len(tt.expected), len(result))

			for quorumID, expectedSlice := range tt.expected {
				resultSlice, exists := result[quorumID]
				assert.True(t, exists)
				assert.Equal(t, len(expectedSlice), len(resultSlice))

				for i, expectedRecord := range expectedSlice {
					assert.NotNil(t, resultSlice[i])
					assert.Equal(t, expectedRecord.Index, resultSlice[i].Index)
					assert.Equal(t, expectedRecord.Usage, resultSlice[i].Usage)
				}
			}
		})
	}
}
