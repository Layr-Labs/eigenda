package ondemand_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/core/payments/ondemand"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewOnDemandPaymentValidator(t *testing.T) {
	mockOnChainState := &coremock.MockOnchainPaymentState{}
	tableName := "test-table"
	maxLedgers := 100

	t.Run("nil onChainState", func(t *testing.T) {
		validator, err := ondemand.NewOnDemandPaymentValidator(
			testutils.GetLogger(),
			maxLedgers,
			nil,
			dynamoClient,
			tableName,
		)
		require.Error(t, err)
		require.Nil(t, validator)
	})

	t.Run("nil dynamoClient", func(t *testing.T) {
		validator, err := ondemand.NewOnDemandPaymentValidator(
			testutils.GetLogger(),
			maxLedgers,
			mockOnChainState,
			nil,
			tableName,
		)
		require.Error(t, err)
		require.Nil(t, validator)
	})

	t.Run("empty table name", func(t *testing.T) {
		validator, err := ondemand.NewOnDemandPaymentValidator(
			testutils.GetLogger(),
			maxLedgers,
			mockOnChainState,
			dynamoClient,
			"",
		)
		require.Error(t, err)
		require.Nil(t, validator)
	})

	t.Run("zero max ledgers", func(t *testing.T) {
		validator, err := ondemand.NewOnDemandPaymentValidator(
			testutils.GetLogger(),
			0,
			mockOnChainState,
			dynamoClient,
			tableName,
		)
		require.Error(t, err)
		require.Nil(t, validator)
	})

	t.Run("negative max ledgers", func(t *testing.T) {
		validator, err := ondemand.NewOnDemandPaymentValidator(
			testutils.GetLogger(),
			-1,
			mockOnChainState,
			dynamoClient,
			tableName,
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

	mockOnChainState := &coremock.MockOnchainPaymentState{}
	mockOnChainState.On("GetPricePerSymbol").Return(uint64(100))
	mockOnChainState.On("GetMinNumSymbols").Return(uint64(1))

	mockOnChainState.On("GetOnDemandPaymentByAccount", mock.Anything, accountA).Return(
		&core.OnDemandPayment{CumulativePayment: big.NewInt(10000)}, nil)
	mockOnChainState.On("GetOnDemandPaymentByAccount", mock.Anything, accountB).Return(
		&core.OnDemandPayment{CumulativePayment: big.NewInt(20000)}, nil)

	paymentValidator, err := ondemand.NewOnDemandPaymentValidator(
		testutils.GetLogger(),
		10,
		mockOnChainState,
		dynamoClient,
		tableName,
	)
	require.NoError(t, err)
	require.NotNil(t, paymentValidator)

	// debit from account A
	err = paymentValidator.Debit(ctx, accountA, uint32(50), []uint8{0})
	require.NoError(t, err, "first debit from account A should succeed")

	// debit from account B
	err = paymentValidator.Debit(ctx, accountB, uint32(75), []uint8{0, 1})
	require.NoError(t, err, "first debit from account B should succeed")

	// debit from account A (should reuse cached ledger)
	err = paymentValidator.Debit(ctx, accountA, uint32(25), []uint8{0})
	require.NoError(t, err, "second debit from account A should succeed")

	// Each account should only trigger GetOnDemandPaymentByAccount once (on first access)
	mockOnChainState.AssertNumberOfCalls(t, "GetOnDemandPaymentByAccount", 2)
	mockOnChainState.AssertCalled(t, "GetPricePerSymbol")
	mockOnChainState.AssertCalled(t, "GetMinNumSymbols")
}

func TestDebitInsufficientFunds(t *testing.T) {
	ctx := context.Background()
	tableName := createPaymentTable(t, "TestDebitInsufficientFunds")
	defer deleteTable(t, tableName)

	accountID := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")

	mockOnChainState := &coremock.MockOnchainPaymentState{}
	mockOnChainState.On("GetPricePerSymbol").Return(uint64(1000))
	mockOnChainState.On("GetMinNumSymbols").Return(uint64(1))
	mockOnChainState.On("GetOnDemandPaymentByAccount", mock.Anything, accountID).Return(
		&core.OnDemandPayment{CumulativePayment: big.NewInt(5000)}, nil)

	paymentValidator, err := ondemand.NewOnDemandPaymentValidator(
		testutils.GetLogger(),
		10,
		mockOnChainState,
		dynamoClient,
		tableName,
	)
	require.NoError(t, err)

	// Try to debit more than available funds (5000 wei / 1000 wei per symbol = 5 symbols max)
	err = paymentValidator.Debit(ctx, accountID, uint32(10), []uint8{0})
	require.Error(t, err, "debit should fail when insufficient funds")
	var insufficientFundsErr *ondemand.InsufficientFundsError
	require.ErrorAs(t, err, &insufficientFundsErr, "error should be InsufficientFundsError")

	updates := []ondemand.TotalDepositUpdate{
		{
			// Update total deposits to 15000 wei (enough for 15 symbols at 1000 wei each)
			AccountAddress:  accountID,
			NewTotalDeposit: big.NewInt(15000),
		},
		{
			// Also include an untracked account that should be skipped, to exercise that logic
			AccountAddress:  gethcommon.HexToAddress("0xcccccccccccccccccccccccccccccccccccccccc"),
			NewTotalDeposit: big.NewInt(50000),
		},
	}
	err = paymentValidator.UpdateTotalDeposits(updates)
	require.NoError(t, err, "updating total deposits should succeed")

	// Retry the same debit that previously failed - should now succeed
	err = paymentValidator.Debit(ctx, accountID, uint32(10), []uint8{0})
	require.NoError(t, err, "debit should now succeed after increasing deposits")
}

func TestLRUCacheEvictionAndReload(t *testing.T) {
	ctx := context.Background()
	tableName := createPaymentTable(t, "TestLRUCacheEvictionAndReload")
	defer deleteTable(t, tableName)

	accountA := gethcommon.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	accountB := gethcommon.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	accountC := gethcommon.HexToAddress("0xcccccccccccccccccccccccccccccccccccccccc")

	mockOnChainState := &coremock.MockOnchainPaymentState{}
	mockOnChainState.On("GetPricePerSymbol").Return(uint64(1000))
	mockOnChainState.On("GetMinNumSymbols").Return(uint64(1))

	// Account A has 8000 wei total deposits (can afford 8 symbols at 1000 wei each)
	mockOnChainState.On("GetOnDemandPaymentByAccount", mock.Anything, accountA).Return(
		&core.OnDemandPayment{CumulativePayment: big.NewInt(8000)}, nil)

	mockOnChainState.On("GetOnDemandPaymentByAccount", mock.Anything, accountB).Return(
		&core.OnDemandPayment{CumulativePayment: big.NewInt(5000)}, nil)
	mockOnChainState.On("GetOnDemandPaymentByAccount", mock.Anything, accountC).Return(
		&core.OnDemandPayment{CumulativePayment: big.NewInt(3000)}, nil)

	// Create paymentValidator with small LRU cache size to force eviction
	paymentValidator, err := ondemand.NewOnDemandPaymentValidator(
		testutils.GetLogger(),
		2,
		mockOnChainState,
		dynamoClient,
		tableName,
	)
	require.NoError(t, err)
	require.NotNil(t, paymentValidator)

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

	// Verify that GetOnDemandPaymentByAccount was called exactly 4 times:
	// 1. Initial call for account A
	// 2. Initial call for account B
	// 3. Initial call for account C
	// 4. Second call for account A after it was evicted and accessed again
	mockOnChainState.AssertNumberOfCalls(t, "GetOnDemandPaymentByAccount", 4)
}
