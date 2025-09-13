package reservation_test

import (
	"context"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/testutils"
	bindings "github.com/Layr-Labs/eigenda/contracts/bindings/v2/PaymentVault"
	"github.com/Layr-Labs/eigenda/core/payments/reservation"
	"github.com/Layr-Labs/eigenda/core/payments/vault"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestNewReservationLedgerCacheInvalidParams(t *testing.T) {
	testTime := time.Date(1971, 8, 15, 0, 0, 0, 0, time.UTC)

	t.Run("nil payment vault", func(t *testing.T) {
		cache, err := reservation.NewReservationLedgerCache(
			context.Background(),
			testutils.GetLogger(),
			10,
			nil, // nil payment vault
			func() time.Time { return testTime },
			reservation.OverfillOncePermitted,
			10*time.Second,
			time.Second,
		)
		require.Error(t, err)
		require.Nil(t, cache)
	})

	t.Run("nil time source", func(t *testing.T) {
		cache, err := reservation.NewReservationLedgerCache(
			context.Background(),
			testutils.GetLogger(),
			10,
			vault.NewTestPaymentVault(),
			nil, // nil time source
			reservation.OverfillOncePermitted,
			10*time.Second,
			time.Second,
		)
		require.Error(t, err)
		require.Nil(t, cache)
	})

	t.Run("invalid capacity duration", func(t *testing.T) {
		cache, err := reservation.NewReservationLedgerCache(
			context.Background(),
			testutils.GetLogger(),
			10,
			vault.NewTestPaymentVault(),
			func() time.Time { return testTime },
			reservation.OverfillOncePermitted,
			0, // invalid capacity duration (zero)
			time.Second,
		)
		require.Error(t, err)
		require.Nil(t, cache)
	})
}

func TestLRUCacheEvictionAndReload(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	accountA := gethcommon.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	accountB := gethcommon.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	accountC := gethcommon.HexToAddress("0xcccccccccccccccccccccccccccccccccccccccc")

	testTime := time.Date(1971, 8, 15, 0, 0, 0, 0, time.UTC)
	timeSource := func() time.Time { return testTime }

	testVault := vault.NewTestPaymentVault()
	testVault.SetReservation(accountA, &bindings.IPaymentVaultReservation{
		SymbolsPerSecond: 8,
		StartTimestamp:   uint64(testTime.Unix() - 3600), // started 1 hour ago
		EndTimestamp:     uint64(testTime.Unix() + 3600), // ends in 1 hour
		QuorumNumbers:    []byte{0},
		QuorumSplits:     []byte{100},
	})
	testVault.SetReservation(accountB, &bindings.IPaymentVaultReservation{
		SymbolsPerSecond: 5,
		StartTimestamp:   uint64(testTime.Unix() - 3600),
		EndTimestamp:     uint64(testTime.Unix() + 3600),
		QuorumNumbers:    []byte{0},
		QuorumSplits:     []byte{100},
	})
	testVault.SetReservation(accountC, &bindings.IPaymentVaultReservation{
		SymbolsPerSecond: 3,
		StartTimestamp:   uint64(testTime.Unix() - 3600),
		EndTimestamp:     uint64(testTime.Unix() + 3600),
		QuorumNumbers:    []byte{0},
		QuorumSplits:     []byte{100},
	})

	ledgerCache, err := reservation.NewReservationLedgerCache(
		ctx,
		testutils.GetLogger(),
		2, // Small cache size to force eviction
		testVault,
		timeSource,
		reservation.OverfillOncePermitted,
		time.Second,
		time.Millisecond,
	)
	require.NoError(t, err)
	require.NotNil(t, ledgerCache)

	// Get ledger for account A and perform a debit
	ledgerA, err := ledgerCache.GetOrCreate(ctx, accountA)
	require.NoError(t, err)
	success, _, err := ledgerA.Debit(testTime, testTime, uint32(9), []uint8{0})
	require.NoError(t, err, "first debit from account A should succeed")
	require.True(t, success, "first debit from account A should succeed")

	// Add accounts B and C to cache, evicting account A
	ledgerB, err := ledgerCache.GetOrCreate(ctx, accountB)
	require.NoError(t, err)
	success, _, err = ledgerB.Debit(testTime, testTime, uint32(3), []uint8{0})
	require.NoError(t, err, "debit from account B should succeed")
	require.True(t, success, "debit from account B should succeed")
	ledgerC, err := ledgerCache.GetOrCreate(ctx, accountC)
	require.NoError(t, err)
	success, _, err = ledgerC.Debit(testTime, testTime, uint32(2), []uint8{0})
	require.NoError(t, err, "debit from account C should succeed")
	require.True(t, success, "debit from account C should succeed")

	// At this point, account A should have been evicted from the LRU cache
	// Cache now contains accounts B and C only

	// Get account A again - should reload from vault with fresh (empty) state
	ledgerAReloaded, err := ledgerCache.GetOrCreate(ctx, accountA)
	require.NoError(t, err)

	// Account A starts fresh with an empty bucket
	// Since bucket capacity is 1 second and rate is 8 symbols/sec, it can hold 8 symbols total
	// The fresh bucket should allow the full capacity
	success, _, err = ledgerAReloaded.Debit(testTime, testTime, uint32(8), []uint8{0})
	require.NoError(t, err, "second debit from reloaded account A should not error")
	require.True(t, success, "second debit from reloaded account A should succeed with fresh bucket")

	// Now trying to add 1 more symbol should fail on capacity since bucket is full
	success, _, err = ledgerAReloaded.Debit(testTime, testTime, uint32(1), []uint8{0})
	require.NoError(t, err, "third debit from account A should not error")
	require.False(t, success, "third debit from account A should fail due to insufficient capacity")

	// simulate a new reservation update for account A with higher capacity
	testVault.SetReservation(accountA, &bindings.IPaymentVaultReservation{
		SymbolsPerSecond: 12, // increased capacity
		StartTimestamp:   uint64(testTime.Unix() - 3600),
		EndTimestamp:     uint64(testTime.Unix() + 3600),
		QuorumNumbers:    []byte{0},
		QuorumSplits:     []byte{100},
	})

	// sleep for long enough for the update to be picked up by the monitor
	time.Sleep(time.Millisecond * 10)

	// try adding more symbols again - should work due to increased capacity
	success, _, err = ledgerAReloaded.Debit(testTime, testTime, uint32(4), []uint8{0})
	require.NoError(t, err)
	require.True(t, success)
}
