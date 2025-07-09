package meterer_test

import (
	"context"
	"math/big"
	"sync"
	"testing"
	"time"

	disperser_v2 "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/core/meterer/payment_logic"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLocalAccountLedger_Construction tests object creation and protobuf serialization
func TestLocalAccountLedger_Construction(t *testing.T) {
	t.Run("empty constructor", func(t *testing.T) {
		ledger := meterer.NewLocalAccountLedger()
		assert.NotNil(t, ledger)

		// Verify empty state
		reservations, periodRecords, onchainPayment, cumulativePayment := ledger.GetAccountStateProtobuf()
		assert.Empty(t, reservations)
		assert.Empty(t, periodRecords)
		assert.Equal(t, big.NewInt(0).Bytes(), onchainPayment)
		assert.Equal(t, big.NewInt(0).Bytes(), cumulativePayment)
	})

	t.Run("protobuf deserialization", func(t *testing.T) {
		// Test with real data including edge cases
		reservations := map[uint32]*disperser_v2.QuorumReservation{
			1: {
				SymbolsPerSecond: 100,
				StartTimestamp:   1000,
				EndTimestamp:     2000,
			},
			2: nil, // Test nil handling
		}

		periodRecords := map[uint32]*disperser_v2.PeriodRecords{
			1: {
				Records: []*disperser_v2.PeriodRecord{
					{Index: 0, Usage: 50},
					{Index: 1, Usage: 75},
				},
			},
			3: nil, // Test nil handling
		}

		onchainPayment := big.NewInt(1000)
		cumulativePayment := big.NewInt(500)

		ledger, err := meterer.NewLocalAccountLedgerFromProtobuf(
			reservations,
			periodRecords,
			onchainPayment.Bytes(),
			cumulativePayment.Bytes(),
		)
		require.NoError(t, err)

		// Verify round-trip serialization
		protoReservations, protoPeriodRecords, protoOnchainPayment, protoCumulativePayment := ledger.GetAccountStateProtobuf()

		// Only non-nil entries should be preserved
		assert.Len(t, protoReservations, 1)
		assert.Equal(t, reservations[1].SymbolsPerSecond, protoReservations[1].SymbolsPerSecond)

		assert.Len(t, protoPeriodRecords, 1)
		assert.Len(t, protoPeriodRecords[1].Records, 2)

		assert.Equal(t, onchainPayment.Bytes(), protoOnchainPayment)
		assert.Equal(t, cumulativePayment.Bytes(), protoCumulativePayment)
	})

	t.Run("protobuf edge cases", func(t *testing.T) {
		// Test with empty byte arrays
		ledger, err := meterer.NewLocalAccountLedgerFromProtobuf(
			make(map[uint32]*disperser_v2.QuorumReservation),
			make(map[uint32]*disperser_v2.PeriodRecords),
			[]byte{},
			[]byte{},
		)
		require.NoError(t, err)

		reservations, periodRecords, onchainPayment, cumulativePayment := ledger.GetAccountStateProtobuf()
		assert.Empty(t, reservations)
		assert.Empty(t, periodRecords)
		assert.Equal(t, big.NewInt(0).Bytes(), onchainPayment)
		assert.Equal(t, big.NewInt(0).Bytes(), cumulativePayment)

		// Test with all nil maps
		ledger, err = meterer.NewLocalAccountLedgerFromProtobuf(
			nil,
			nil,
			big.NewInt(100).Bytes(),
			big.NewInt(50).Bytes(),
		)
		require.NoError(t, err)

		reservations, periodRecords, onchainPayment, cumulativePayment = ledger.GetAccountStateProtobuf()
		assert.Empty(t, reservations)
		assert.Empty(t, periodRecords)
		assert.Equal(t, big.NewInt(100).Bytes(), onchainPayment)
		assert.Equal(t, big.NewInt(50).Bytes(), cumulativePayment)

		// Test with empty records list
		emptyRecords := map[uint32]*disperser_v2.PeriodRecords{
			1: {Records: []*disperser_v2.PeriodRecord{}}, // Empty records
		}

		ledger, err = meterer.NewLocalAccountLedgerFromProtobuf(
			make(map[uint32]*disperser_v2.QuorumReservation),
			emptyRecords,
			[]byte{},
			[]byte{},
		)
		require.NoError(t, err)

		_, periodRecords, _, _ = ledger.GetAccountStateProtobuf()
		assert.Empty(t, periodRecords) // Empty records should not be included
	})
}

// TestLocalAccountLedger_CoreBehavior tests the main business logic: reservation → overflow → on-demand fallback
func TestLocalAccountLedger_CoreBehavior(t *testing.T) {
	accountID := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")
	ctx := context.Background()
	now := time.Now()

	// Create ledger with reservation: 50 symbols/sec * 2 sec window = 100 symbol limit
	reservations := map[uint32]*disperser_v2.QuorumReservation{
		0: {
			SymbolsPerSecond: 50,
			StartTimestamp:   uint32(now.Add(-time.Hour).Unix()),
			EndTimestamp:     uint32(now.Add(time.Hour).Unix()),
		},
		1: {
			SymbolsPerSecond: 50,
			StartTimestamp:   uint32(now.Add(-time.Hour).Unix()),
			EndTimestamp:     uint32(now.Add(time.Hour).Unix()),
		},
	}

	ledger, err := meterer.NewLocalAccountLedgerFromProtobuf(
		reservations,
		make(map[uint32]*disperser_v2.PeriodRecords),
		big.NewInt(1000).Bytes(), // Sufficient on-demand balance
		big.NewInt(0).Bytes(),
	)
	require.NoError(t, err)

	params := &meterer.PaymentVaultParams{
		QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
			0: {MinNumSymbols: 1, ReservationRateLimitWindow: 2},
			1: {MinNumSymbols: 1, ReservationRateLimitWindow: 2},
		},
		QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
			0: {OnDemandPricePerSymbol: 1},
			1: {OnDemandPricePerSymbol: 1},
		},
		OnDemandQuorumNumbers: []core.QuorumID{0, 1},
	}

	t.Run("reservation path", func(t *testing.T) {
		// 80 symbols within 100 limit - should use reservation
		payment, err := ledger.Debit(ctx, accountID, now.UnixNano(), 80, []core.QuorumID{0, 1}, params)
		assert.NoError(t, err)
		assert.Nil(t, payment) // reservation returns nil
	})

	t.Run("overflow path", func(t *testing.T) {
		// 30 more symbols (80+30=110) - should overflow to overflow bin
		payment, err := ledger.Debit(ctx, accountID, now.UnixNano(), 30, []core.QuorumID{0, 1}, params)
		assert.NoError(t, err)
		assert.Nil(t, payment) // still reservation (overflow bin)
	})

	t.Run("on-demand fallback path", func(t *testing.T) {
		// Another request should fall back to on-demand
		payment, err := ledger.Debit(ctx, accountID, now.UnixNano(), 20, []core.QuorumID{0, 1}, params)
		assert.NoError(t, err)
		assert.NotNil(t, payment) // on-demand payment
		assert.Equal(t, big.NewInt(20), payment)
	})

	t.Run("insufficient balance error", func(t *testing.T) {
		// Exceed available balance
		payment, err := ledger.Debit(ctx, accountID, now.UnixNano(), 2000, []core.QuorumID{0, 1}, params)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient ondemand payment")
		assert.Nil(t, payment)
	})

	t.Run("config errors", func(t *testing.T) {
		// Missing quorum config
		badParams := &meterer.PaymentVaultParams{
			QuorumProtocolConfigs: make(map[core.QuorumID]*core.PaymentQuorumProtocolConfig),
			QuorumPaymentConfigs:  make(map[core.QuorumID]*core.PaymentQuorumConfig),
			OnDemandQuorumNumbers: []core.QuorumID{},
		}

		payment, err := ledger.Debit(ctx, accountID, now.UnixNano(), 50, []core.QuorumID{0}, badParams)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no payment method available")
		assert.Nil(t, payment)

		// Empty quorums
		payment, err = ledger.Debit(ctx, accountID, now.UnixNano(), 50, []core.QuorumID{}, params)
		assert.Error(t, err)
		assert.Nil(t, payment)

		// Quorum mismatch
		mismatchParams := &meterer.PaymentVaultParams{
			QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
				0: {MinNumSymbols: 1, ReservationRateLimitWindow: 2},
			},
			QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
				0: {OnDemandPricePerSymbol: 1},
			},
			OnDemandQuorumNumbers: []core.QuorumID{0}, // Only quorum 0 enabled
		}

		payment, err = ledger.Debit(ctx, accountID, now.UnixNano(), 50, []core.QuorumID{99}, mismatchParams) // Request non-existent quorum
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no payment method available")
		assert.Nil(t, payment)
	})

	t.Run("edge cases", func(t *testing.T) {
		// Create fresh ledger for edge case testing (previous tests exhausted the reservation)
		freshLedger, err := meterer.NewLocalAccountLedgerFromProtobuf(
			reservations,
			make(map[uint32]*disperser_v2.PeriodRecords),
			big.NewInt(1000).Bytes(), // Sufficient on-demand balance
			big.NewInt(0).Bytes(),
		)
		require.NoError(t, err)

		// Zero symbols (LocalAccountLedger applies MinNumSymbols and uses reservation)
		payment, err := freshLedger.Debit(ctx, accountID, now.UnixNano(), 0, []core.QuorumID{0}, params)
		assert.NoError(t, err) // LocalAccountLedger accepts 0 symbols (applies MinNumSymbols)
		assert.Nil(t, payment) // Should use reservation

		// Very large symbols (should overflow to on-demand then fail)
		payment, err = ledger.Debit(ctx, accountID, now.UnixNano(), 999999, []core.QuorumID{0, 1}, params)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient ondemand payment")
		assert.Nil(t, payment)
	})

	t.Run("mixed reservation states", func(t *testing.T) {
		// Test with expired, future, and active reservations
		mixedReservations := map[uint32]*disperser_v2.QuorumReservation{
			0: {
				SymbolsPerSecond: 50,
				StartTimestamp:   uint32(now.Add(time.Hour).Unix()), // Future
				EndTimestamp:     uint32(now.Add(2 * time.Hour).Unix()),
			},
			1: {
				SymbolsPerSecond: 50,
				StartTimestamp:   uint32(now.Add(-2 * time.Hour).Unix()), // Expired
				EndTimestamp:     uint32(now.Add(-time.Hour).Unix()),
			},
			2: {
				SymbolsPerSecond: 50,
				StartTimestamp:   uint32(now.Add(-time.Hour).Unix()), // Active
				EndTimestamp:     uint32(now.Add(time.Hour).Unix()),
			},
		}

		mixedLedger, err := meterer.NewLocalAccountLedgerFromProtobuf(
			mixedReservations,
			make(map[uint32]*disperser_v2.PeriodRecords),
			big.NewInt(1000).Bytes(),
			big.NewInt(0).Bytes(),
		)
		require.NoError(t, err)

		mixedParams := &meterer.PaymentVaultParams{
			QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
				0: {MinNumSymbols: 1, ReservationRateLimitWindow: 2},
				1: {MinNumSymbols: 1, ReservationRateLimitWindow: 2},
				2: {MinNumSymbols: 1, ReservationRateLimitWindow: 2},
			},
			QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
				0: {OnDemandPricePerSymbol: 1},
				1: {OnDemandPricePerSymbol: 1},
				2: {OnDemandPricePerSymbol: 1},
			},
			OnDemandQuorumNumbers: []core.QuorumID{0, 1}, // Quorum 2 not enabled for on-demand
		}

		// Should fail when trying all three quorums (quorum 2 can't fall back to on-demand)
		payment, err := mixedLedger.Debit(ctx, accountID, now.UnixNano(), 50, []core.QuorumID{0, 1, 2}, mixedParams)
		assert.Error(t, err)
		assert.Nil(t, payment)

		// Should succeed with only active reservation quorum
		payment, err = mixedLedger.Debit(ctx, accountID, now.UnixNano(), 50, []core.QuorumID{2}, mixedParams)
		assert.NoError(t, err)
		assert.Nil(t, payment) // Uses reservation

		// Should succeed with on-demand enabled quorums
		payment, err = mixedLedger.Debit(ctx, accountID, now.UnixNano(), 50, []core.QuorumID{0, 1}, mixedParams)
		assert.NoError(t, err)
		assert.NotNil(t, payment) // Uses on-demand
	})
}

// TestLocalAccountLedger_RevertDebit tests the unique revert functionality (not in Accountant)
func TestLocalAccountLedger_RevertDebit(t *testing.T) {
	accountID := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")
	ctx := context.Background()
	now := time.Now()

	t.Run("revert reservation usage", func(t *testing.T) {
		// Setup ledger with reservation
		reservations := map[uint32]*disperser_v2.QuorumReservation{
			0: {
				SymbolsPerSecond: 100,
				StartTimestamp:   uint32(now.Unix()),
				EndTimestamp:     uint32(now.Add(time.Hour).Unix()),
			},
		}

		// Calculate the correct period for the current timestamp
		currentPeriod := payment_logic.GetReservationPeriodByNanosecond(now.UnixNano(), 10)

		// Create ledger with existing usage
		periodRecords := map[uint32]*disperser_v2.PeriodRecords{
			0: {
				Records: []*disperser_v2.PeriodRecord{
					{Index: uint32(currentPeriod), Usage: 50}, // Pre-existing usage for correct period
				},
			},
		}

		ledger, err := meterer.NewLocalAccountLedgerFromProtobuf(
			reservations,
			periodRecords,
			big.NewInt(0).Bytes(),
			big.NewInt(0).Bytes(),
		)
		require.NoError(t, err)

		params := &meterer.PaymentVaultParams{
			QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
				0: {MinNumSymbols: 1, ReservationRateLimitWindow: 10},
			},
			QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
				0: {OnDemandPricePerSymbol: 1},
			},
		}

		// Revert 30 symbols (should succeed)
		err = ledger.RevertDebit(ctx, accountID, now.UnixNano(), 30, []core.QuorumID{0}, params, nil)
		assert.NoError(t, err)

		// Try to revert more than available (should fail)
		err = ledger.RevertDebit(ctx, accountID, now.UnixNano(), 50, []core.QuorumID{0}, params, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient usage")
	})

	t.Run("revert on-demand usage", func(t *testing.T) {
		ledger, err := meterer.NewLocalAccountLedgerFromProtobuf(
			make(map[uint32]*disperser_v2.QuorumReservation),
			make(map[uint32]*disperser_v2.PeriodRecords),
			big.NewInt(0).Bytes(),
			big.NewInt(100).Bytes(), // Current cumulative payment
		)
		require.NoError(t, err)

		// Revert 50 payment (should succeed)
		err = ledger.RevertDebit(ctx, accountID, now.UnixNano(), 0, []core.QuorumID{0}, nil, big.NewInt(50))
		assert.NoError(t, err)

		// Try to revert more than available (should fail)
		err = ledger.RevertDebit(ctx, accountID, now.UnixNano(), 0, []core.QuorumID{0}, nil, big.NewInt(100))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient cumulative payment")
	})

	t.Run("revert edge cases", func(t *testing.T) {
		ledger := meterer.NewLocalAccountLedger()

		// Revert with invalid payment amount (zero or negative)
		err := ledger.RevertDebit(ctx, accountID, now.UnixNano(), 0, []core.QuorumID{0}, nil, big.NewInt(0))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid payment amount")

		err = ledger.RevertDebit(ctx, accountID, now.UnixNano(), 0, []core.QuorumID{0}, nil, big.NewInt(-10))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid payment amount")

		// Revert reservation usage when no reservation exists
		params := &meterer.PaymentVaultParams{
			QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
				0: {MinNumSymbols: 1, ReservationRateLimitWindow: 10},
			},
		}

		err = ledger.RevertDebit(ctx, accountID, now.UnixNano(), 50, []core.QuorumID{0}, params, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot revert reservation usage")

		// Revert with empty quorum list
		err = ledger.RevertDebit(ctx, accountID, now.UnixNano(), 50, []core.QuorumID{}, params, nil)
		assert.Error(t, err)
	})
}

// TestLocalAccountLedger_ConcurrentAccess tests thread safety (critical for production)
func TestLocalAccountLedger_ConcurrentAccess(t *testing.T) {
	ledger, err := meterer.NewLocalAccountLedgerFromProtobuf(
		make(map[uint32]*disperser_v2.QuorumReservation),
		make(map[uint32]*disperser_v2.PeriodRecords),
		big.NewInt(10000).Bytes(), // Large balance for concurrent operations
		big.NewInt(0).Bytes(),
	)
	require.NoError(t, err)

	params := &meterer.PaymentVaultParams{
		QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
			0: {OnDemandPricePerSymbol: 1},
		},
		QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
			0: {MinNumSymbols: 1},
		},
		OnDemandQuorumNumbers: []core.QuorumID{0},
	}

	accountID := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")
	ctx := context.Background()

	const numGoroutines = 50
	const operationsPerGoroutine = 10

	var wg sync.WaitGroup
	var mu sync.Mutex
	payments := make([]*big.Int, 0)

	// Launch concurrent operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				payment, err := ledger.Debit(ctx, accountID, time.Now().UnixNano(), 1, []core.QuorumID{0}, params)
				if err == nil && payment != nil {
					mu.Lock()
					payments = append(payments, payment)
					mu.Unlock()
				}

				// Also test concurrent reads
				ledger.GetAccountStateProtobuf()
			}
		}()
	}

	wg.Wait()

	// Verify no race conditions occurred
	assert.True(t, len(payments) > 0, "Some payments should have succeeded")

	// Verify payments are reasonable (basic sanity check for concurrent safety)
	for _, payment := range payments {
		assert.True(t, payment.Cmp(big.NewInt(0)) > 0, "Payment should be positive")
		assert.True(t, payment.Cmp(big.NewInt(1000)) <= 0, "Payment should be reasonable size")
	}
}

// TestLocalAccountLedger_OverflowRollback tests complex overflow scenarios and rollback behavior
func TestLocalAccountLedger_OverflowRollback(t *testing.T) {
	accountID := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")
	ctx := context.Background()
	now := time.Now()

	// Setup ledger with small limits to force overflow scenarios
	reservations := map[uint32]*disperser_v2.QuorumReservation{
		0: {
			SymbolsPerSecond: 50, // 50 * 2 = 100 symbol limit
			StartTimestamp:   uint32(now.Unix()),
			EndTimestamp:     uint32(now.Add(time.Hour).Unix()),
		},
		1: {
			SymbolsPerSecond: 50,
			StartTimestamp:   uint32(now.Unix()),
			EndTimestamp:     uint32(now.Add(time.Hour).Unix()),
		},
	}

	ledger, err := meterer.NewLocalAccountLedgerFromProtobuf(
		reservations,
		make(map[uint32]*disperser_v2.PeriodRecords),
		big.NewInt(0).Bytes(), // No on-demand balance to force specific errors
		big.NewInt(0).Bytes(),
	)
	require.NoError(t, err)

	params := &meterer.PaymentVaultParams{
		QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
			0: {MinNumSymbols: 1, ReservationRateLimitWindow: 2},
			1: {MinNumSymbols: 1, ReservationRateLimitWindow: 2},
		},
		QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
			0: {OnDemandPricePerSymbol: 1},
			1: {OnDemandPricePerSymbol: 1},
		},
		OnDemandQuorumNumbers: []core.QuorumID{0, 1},
	}

	t.Run("reservation rollback on multi-quorum failure", func(t *testing.T) {
		nowNano := now.UnixNano()

		// Fill up one quorum to its limit
		payment, err := ledger.Debit(ctx, accountID, nowNano, 50, []core.QuorumID{1}, params)
		assert.NoError(t, err)
		assert.Nil(t, payment) // reservation

		// Use both quorums, this should cause quorum 1 to overflow
		payment, err = ledger.Debit(ctx, accountID, nowNano, 60, []core.QuorumID{0, 1}, params)
		assert.NoError(t, err)
		assert.Nil(t, payment) // should still use reservation with overflow

		// Another request on both quorums should fail (overflow bin occupied)
		// and this should rollback properly without partial updates
		payment, err = ledger.Debit(ctx, accountID, nowNano, 60, []core.QuorumID{0, 1}, params)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "reservation limit exceeded")
		assert.Nil(t, payment)

		// Verify that quorum 0 still has available capacity after the failed rollback
		payment, err = ledger.Debit(ctx, accountID, nowNano, 40, []core.QuorumID{0}, params)
		assert.NoError(t, err)
		assert.Nil(t, payment) // Should succeed with reservation
	})

	t.Run("rollback on config errors", func(t *testing.T) {
		// Create fresh ledger for config error testing
		freshLedger, err := meterer.NewLocalAccountLedgerFromProtobuf(
			reservations,
			make(map[uint32]*disperser_v2.PeriodRecords),
			big.NewInt(0).Bytes(), // No on-demand balance to force config error
			big.NewInt(0).Bytes(),
		)
		require.NoError(t, err)

		// Create params with missing config for one quorum
		badParams := &meterer.PaymentVaultParams{
			QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
				0: {MinNumSymbols: 1, ReservationRateLimitWindow: 2},
				// Missing config for quorum 1
			},
			QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
				0: {OnDemandPricePerSymbol: 1},
			},
			OnDemandQuorumNumbers: []core.QuorumID{},
		}

		// This should fail due to missing config and not make any partial updates
		payment, err := freshLedger.Debit(ctx, accountID, now.UnixNano(), 10, []core.QuorumID{0, 1}, badParams)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "quorum config not found")
		assert.Nil(t, payment)

		// Verify quorum 0 is still accessible with proper config
		payment, err = freshLedger.Debit(ctx, accountID, now.UnixNano(), 10, []core.QuorumID{0}, params)
		assert.NoError(t, err)
		assert.Nil(t, payment)
	})

	t.Run("overflow bin timing edge case", func(t *testing.T) {
		// Create a stable timestamp for consistent timing
		testTime := time.Now()

		// Create reservations that are definitely active for testing
		timingReservations := map[uint32]*disperser_v2.QuorumReservation{
			0: {
				SymbolsPerSecond: 50,
				StartTimestamp:   uint32(testTime.Add(-time.Minute).Unix()), // Start 1 minute ago
				EndTimestamp:     uint32(testTime.Add(time.Hour).Unix()),    // End 1 hour from now
			},
		}

		// Create fresh ledger for timing testing
		freshLedger, err := meterer.NewLocalAccountLedgerFromProtobuf(
			timingReservations,
			make(map[uint32]*disperser_v2.PeriodRecords),
			big.NewInt(100).Bytes(), // Small on-demand balance for fallback testing
			big.NewInt(0).Bytes(),
		)
		require.NoError(t, err)

		// Test that current timestamp works (within reservation window)
		payment, err := freshLedger.Debit(ctx, accountID, testTime.UnixNano(), 20, []core.QuorumID{0}, params)
		assert.NoError(t, err)
		assert.Nil(t, payment) // Should use reservation

		// Test with timestamp before reservation starts (falls back to on-demand)
		beforeStartTime := testTime.Add(-2 * time.Minute).UnixNano()
		payment, err = freshLedger.Debit(ctx, accountID, beforeStartTime, 20, []core.QuorumID{0}, params)
		assert.NoError(t, err)    // Falls back to on-demand when reservation is not available
		assert.NotNil(t, payment) // Returns on-demand payment
		assert.Equal(t, big.NewInt(20), payment)
	})
}
