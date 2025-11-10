package ondemand_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core/payments/ondemand"
	"github.com/Layr-Labs/eigenda/core/payments/ondemand/ondemandvalidation"
	"github.com/Layr-Labs/eigenda/core/payments/vault"
	"github.com/Layr-Labs/eigenda/test"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestNewOnDemandLedgerCacheInvalidParams(t *testing.T) {
	ctx := t.Context()

	t.Run("nil payment vault", func(t *testing.T) {
		config, err := ondemandvalidation.NewOnDemandLedgerCacheConfig(
			10,
			"tableName",
			time.Second,
		)
		require.NoError(t, err)

		cleanup, err := test.DeployDynamoLocalstack(t.Context())
		require.NoError(t, err)
		defer cleanup()

		dynamoClient, err := test.GetDynamoClient()
		require.NoError(t, err)

		cache, err := ondemandvalidation.NewOnDemandLedgerCache(
			ctx,
			test.GetLogger(),
			config,
			nil, // nil payment vault
			dynamoClient,
			nil,
		)
		require.Error(t, err)
		require.Nil(t, cache)
	})

	t.Run("nil dynamo client", func(t *testing.T) {
		config, err := ondemandvalidation.NewOnDemandLedgerCacheConfig(
			10,
			"tableName",
			time.Second,
		)
		require.NoError(t, err)

		cache, err := ondemandvalidation.NewOnDemandLedgerCache(
			ctx,
			test.GetLogger(),
			config,
			vault.NewTestPaymentVault(),
			nil, // nil dynamo client
			nil,
		)
		require.Error(t, err)
		require.Nil(t, cache)
	})
}

func TestLRUCacheEvictionAndReload(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	tableName := createPaymentTable(t, "TestLRUCacheEvictionAndReload")
	defer deleteTable(t, tableName)

	accountA := gethcommon.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	accountB := gethcommon.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	accountC := gethcommon.HexToAddress("0xcccccccccccccccccccccccccccccccccccccccc")

	testVault := vault.NewTestPaymentVault()
	testVault.SetPricePerSymbol(1000)
	// Account A has 8000 wei total deposits (can afford 8 symbols at 1000 wei each)
	testVault.SetTotalDeposit(accountA, big.NewInt(8000))
	testVault.SetTotalDeposit(accountB, big.NewInt(5000))
	testVault.SetTotalDeposit(accountC, big.NewInt(3000))

	config, err := ondemandvalidation.NewOnDemandLedgerCacheConfig(
		2, // Small cache size to force eviction
		tableName,
		time.Millisecond, // update frequently
	)
	require.NoError(t, err)

	cleanup, err := test.DeployDynamoLocalstack(t.Context())
	require.NoError(t, err)
	defer cleanup()

	dynamoClient, err := test.GetDynamoClient()
	require.NoError(t, err)

	ledgerCache, err := ondemandvalidation.NewOnDemandLedgerCache(
		ctx,
		test.GetLogger(),
		config,
		testVault,
		dynamoClient,
		nil,
	)
	require.NoError(t, err)
	require.NotNil(t, ledgerCache)

	// Get ledger for account A and perform a debit
	ledgerA, err := ledgerCache.GetOrCreate(ctx, accountA)
	require.NoError(t, err)
	_, err = ledgerA.Debit(ctx, uint32(6), []uint8{0}) // 6 symbols = 6000 wei
	require.NoError(t, err, "first debit from account A should succeed")

	// Add accounts B and C to cache, evicting account A
	ledgerB, err := ledgerCache.GetOrCreate(ctx, accountB)
	require.NoError(t, err)
	_, err = ledgerB.Debit(ctx, uint32(3), []uint8{0})
	require.NoError(t, err, "debit from account B should succeed")
	ledgerC, err := ledgerCache.GetOrCreate(ctx, accountC)
	require.NoError(t, err)
	_, err = ledgerC.Debit(ctx, uint32(2), []uint8{0})
	require.NoError(t, err, "debit from account C should succeed")

	// At this point, account A should have been evicted from the LRU cache
	// Cache now contains accounts B and C only

	// Get account A again - should reload from DynamoDB with persisted state
	ledgerAReloaded, err := ledgerCache.GetOrCreate(ctx, accountA)
	require.NoError(t, err)

	// Account A had 8000 wei total, spent 6000 wei, has 2000 wei left
	// Trying to spend 3000 wei (3 symbols) should fail
	_, err = ledgerAReloaded.Debit(ctx, uint32(3), []uint8{0})
	require.Error(t, err, "second debit from account A should fail due to insufficient funds")
	var insufficientFundsErr *ondemand.InsufficientFundsError
	require.ErrorAs(t, err, &insufficientFundsErr, "error should be InsufficientFundsError")

	// simulate a new deposit by account A
	testVault.SetTotalDeposit(accountA, big.NewInt(10000))

	// wait for the monitor to pick up the deposit update
	test.AssertEventuallyTrue(t, func() bool {
		_, err := ledgerAReloaded.Debit(ctx, uint32(3), []uint8{0})
		return err == nil
	}, time.Second)
}
