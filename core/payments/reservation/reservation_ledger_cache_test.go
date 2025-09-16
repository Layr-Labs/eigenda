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
			t.Context(),
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
			t.Context(),
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
			t.Context(),
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

func TestLRUCacheNormalEviction(t *testing.T) {
	ctx := t.Context()

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

	// Get ledger for account A without performing a debit (bucket remains empty)
	ledgerA, err := ledgerCache.GetOrCreate(ctx, accountA)
	require.NoError(t, err)
	require.NotNil(t, ledgerA)

	// Add accounts B and C to cache
	// This should evict A normally since its bucket is empty
	ledgerB, err := ledgerCache.GetOrCreate(ctx, accountB)
	require.NoError(t, err)
	require.NotNil(t, ledgerB)

	ledgerC, err := ledgerCache.GetOrCreate(ctx, accountC)
	require.NoError(t, err)
	require.NotNil(t, ledgerC)

	// Get account A again - it should be a new instance since it was evicted
	ledgerAReloaded, err := ledgerCache.GetOrCreate(ctx, accountA)
	require.NoError(t, err)
	require.NotNil(t, ledgerAReloaded)

	// The pointers should NOT be the same - this is a new ledger instance
	require.NotSame(t, ledgerA, ledgerAReloaded, "ledger A should have been evicted and recreated, different objects")
}

func TestLRUCachePrematureEviction(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
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
	require.NoError(t, err)
	require.True(t, success, "first debit from account A should succeed")

	// Add accounts B and C to cache
	// This should result in the cache being resized, since A will be evicted prematurely
	ledgerB, err := ledgerCache.GetOrCreate(ctx, accountB)
	require.NoError(t, err)
	require.NotNil(t, ledgerB)

	ledgerC, err := ledgerCache.GetOrCreate(ctx, accountC)
	require.NoError(t, err)
	require.NotNil(t, ledgerC)

	// the LRU cache will have attempted to evict account A, but A's bucket wasn't empty! therefore the cache will have
	// been resized, and the original ledger A should still be present
	ledgerAReloaded, err := ledgerCache.GetOrCreate(ctx, accountA)
	require.NoError(t, err)

	// The pointers should be the same - ledger A should still be in cache
	require.Same(t, ledgerA, ledgerAReloaded, "ledger A should not have been evicted, same object should be returned")

	// Account A should still have its previous debit of 9 symbols
	success, _, err = ledgerAReloaded.Debit(testTime, testTime, uint32(1), []uint8{0})
	require.NoError(t, err)
	require.False(t, success, "second debit from account A should fail - it is over capacity")

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
