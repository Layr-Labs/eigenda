package ondemand_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core/payments/ondemand"
	"github.com/Layr-Labs/eigenda/core/payments/vault"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestNewOnDemandPaymentValidator(t *testing.T) {
	ctx := context.Background()
	tableName := "test-table"
	maxLedgers := 100
	updateInterval := time.Second

	t.Run("nil paymentVault", func(t *testing.T) {
		validator, err := ondemand.NewOnDemandPaymentValidator(
			ctx,
			testutils.GetLogger(),
			maxLedgers,
			nil,
			dynamoClient,
			tableName,
			updateInterval,
		)
		require.Error(t, err)
		require.Nil(t, validator)
	})

	t.Run("nil dynamoClient", func(t *testing.T) {
		testVault := vault.NewTestPaymentVault()

		validator, err := ondemand.NewOnDemandPaymentValidator(
			ctx,
			testutils.GetLogger(),
			maxLedgers,
			testVault,
			nil,
			tableName,
			updateInterval,
		)
		require.Error(t, err)
		require.Nil(t, validator)
	})

	t.Run("zero update interval", func(t *testing.T) {
		testVault := vault.NewTestPaymentVault()

		validator, err := ondemand.NewOnDemandPaymentValidator(
			ctx,
			testutils.GetLogger(),
			maxLedgers,
			testVault,
			dynamoClient,
			tableName,
			0,
		)
		require.Error(t, err)
		require.Nil(t, validator)
	})

	t.Run("negative update interval", func(t *testing.T) {
		testVault := vault.NewTestPaymentVault()

		validator, err := ondemand.NewOnDemandPaymentValidator(
			ctx,
			testutils.GetLogger(),
			maxLedgers,
			testVault,
			dynamoClient,
			tableName,
			-time.Second,
		)
		require.Error(t, err)
		require.Nil(t, validator)
	})
}

func TestDebitMultipleAccounts(t *testing.T) {
	ctx := context.Background()
	tableName := createPaymentTable(t, "TestDebitMultipleAccounts")
	defer deleteTable(t, tableName)

	accountA := gethcommon.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	accountB := gethcommon.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")

	testVault := vault.NewTestPaymentVault()
	testVault.SetDeposit(accountA, big.NewInt(10000))
	testVault.SetDeposit(accountB, big.NewInt(20000))

	paymentValidator, err := ondemand.NewOnDemandPaymentValidator(
		ctx,
		testutils.GetLogger(),
		10,
		testVault,
		dynamoClient,
		tableName,
		time.Second,
	)
	require.NoError(t, err)
	require.NotNil(t, paymentValidator)
	defer paymentValidator.Stop()

	// debit from account A
	err = paymentValidator.Debit(ctx, accountA, uint32(50), []uint8{0})
	require.NoError(t, err, "first debit from account A should succeed")

	// debit from account B
	err = paymentValidator.Debit(ctx, accountB, uint32(75), []uint8{0, 1})
	require.NoError(t, err, "first debit from account B should succeed")

	// debit from account A (should reuse cached ledger)
	err = paymentValidator.Debit(ctx, accountA, uint32(25), []uint8{0})
	require.NoError(t, err, "second debit from account A should succeed")
}

func TestDebitInsufficientFunds(t *testing.T) {
	ctx := context.Background()
	tableName := createPaymentTable(t, "TestDebitInsufficientFunds")
	defer deleteTable(t, tableName)

	accountID := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")

	testVault := vault.NewTestPaymentVault()
	testVault.SetPricePerSymbol(1000)
	testVault.SetDeposit(accountID, big.NewInt(5000))

	paymentValidator, err := ondemand.NewOnDemandPaymentValidator(
		ctx,
		testutils.GetLogger(),
		10,
		testVault,
		dynamoClient,
		tableName,
		time.Second,
	)
	require.NoError(t, err)
	defer paymentValidator.Stop()

	// Try to debit more than available funds (5000 wei / 1000 wei per symbol = 5 symbols max)
	err = paymentValidator.Debit(ctx, accountID, uint32(10), []uint8{0})
	require.Error(t, err, "debit should fail when insufficient funds")
	var insufficientFundsErr *ondemand.InsufficientFundsError
	require.ErrorAs(t, err, &insufficientFundsErr, "error should be InsufficientFundsError")
}

func TestLRUCacheEvictionAndReload(t *testing.T) {
	ctx := context.Background()
	tableName := createPaymentTable(t, "TestLRUCacheEvictionAndReload")
	defer deleteTable(t, tableName)

	accountA := gethcommon.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	accountB := gethcommon.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	accountC := gethcommon.HexToAddress("0xcccccccccccccccccccccccccccccccccccccccc")

	testVault := vault.NewTestPaymentVault()
	testVault.SetPricePerSymbol(1000)
	// Account A has 8000 wei total deposits (can afford 8 symbols at 1000 wei each)
	testVault.SetDeposit(accountA, big.NewInt(8000))
	testVault.SetDeposit(accountB, big.NewInt(5000))
	testVault.SetDeposit(accountC, big.NewInt(3000))

	// Create paymentValidator with small LRU cache size to force eviction
	paymentValidator, err := ondemand.NewOnDemandPaymentValidator(
		ctx,
		testutils.GetLogger(),
		2,
		testVault,
		dynamoClient,
		tableName,
		time.Second,
	)
	require.NoError(t, err)
	require.NotNil(t, paymentValidator)
	defer paymentValidator.Stop()

	// Make a dispersal from account A using 3/4 of total deposits (6 symbols = 6000 wei)
	err = paymentValidator.Debit(ctx, accountA, uint32(6), []uint8{0})
	require.NoError(t, err, "first debit from account A should succeed")

	// Add dispersals from accounts B and C to evict account A from cache
	err = paymentValidator.Debit(ctx, accountB, uint32(3), []uint8{0})
	require.NoError(t, err, "debit from account B should succeed")
	err = paymentValidator.Debit(ctx, accountC, uint32(2), []uint8{0})
	require.NoError(t, err, "debit from account C should succeed")

	// At this point, account A should have been evicted from the LRU cache
	// Cache now contains accounts B and C only

	// Attempt another dispersal from account A that should fail if it was instantiated correctly
	// Account A had 8000 wei total, spent 6000 wei, has 2000 wei left
	// Trying to spend 3000 wei (3 symbols) should fail
	err = paymentValidator.Debit(ctx, accountA, uint32(3), []uint8{0})
	require.Error(t, err, "second debit from account A should fail due to insufficient funds")
	var insufficientFundsErr *ondemand.InsufficientFundsError
	require.ErrorAs(t, err, &insufficientFundsErr, "error should be InsufficientFundsError")
}