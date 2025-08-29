package ondemand_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core/payments/mock"
	"github.com/Layr-Labs/eigenda/core/payments/ondemand"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
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
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockVault := mock.NewMockPaymentVault(ctrl)
		// No expectations needed - the constructor should fail before using the vault

		validator, err := ondemand.NewOnDemandPaymentValidator(
			ctx,
			testutils.GetLogger(),
			maxLedgers,
			mockVault,
			nil,
			tableName,
			updateInterval,
		)
		require.Error(t, err)
		require.Nil(t, validator)
	})

	t.Run("zero update interval", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockVault := mock.NewMockPaymentVault(ctrl)

		validator, err := ondemand.NewOnDemandPaymentValidator(
			ctx,
			testutils.GetLogger(),
			maxLedgers,
			mockVault,
			dynamoClient,
			tableName,
			0,
		)
		require.Error(t, err)
		require.Nil(t, validator)
	})

	t.Run("negative update interval", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockVault := mock.NewMockPaymentVault(ctrl)

		validator, err := ondemand.NewOnDemandPaymentValidator(
			ctx,
			testutils.GetLogger(),
			maxLedgers,
			mockVault,
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockVault := mock.NewMockPaymentVault(ctrl)
	
	// Setup expectations for constructor
	mockVault.EXPECT().GetPricePerSymbol(gomock.Any()).Return(uint64(100), nil)
	mockVault.EXPECT().GetMinNumSymbols(gomock.Any()).Return(uint64(1), nil)
	
	// Setup expectations for GetTotalDeposit calls during GetOrCreate
	mockVault.EXPECT().GetTotalDeposit(gomock.Any(), accountA).Return(big.NewInt(10000), nil)
	mockVault.EXPECT().GetTotalDeposit(gomock.Any(), accountB).Return(big.NewInt(20000), nil)

	paymentValidator, err := ondemand.NewOnDemandPaymentValidator(
		ctx,
		testutils.GetLogger(),
		10,
		mockVault,
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockVault := mock.NewMockPaymentVault(ctrl)
	
	// Setup expectations for constructor
	mockVault.EXPECT().GetPricePerSymbol(gomock.Any()).Return(uint64(1000), nil)
	mockVault.EXPECT().GetMinNumSymbols(gomock.Any()).Return(uint64(1), nil)
	
	// Setup expectations for GetTotalDeposit
	mockVault.EXPECT().GetTotalDeposit(gomock.Any(), accountID).Return(big.NewInt(5000), nil)

	paymentValidator, err := ondemand.NewOnDemandPaymentValidator(
		ctx,
		testutils.GetLogger(),
		10,
		mockVault,
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockVault := mock.NewMockPaymentVault(ctrl)
	
	// Setup expectations for constructor
	mockVault.EXPECT().GetPricePerSymbol(gomock.Any()).Return(uint64(1000), nil)
	mockVault.EXPECT().GetMinNumSymbols(gomock.Any()).Return(uint64(1), nil)
	
	// Account A has 8000 wei total deposits (can afford 8 symbols at 1000 wei each)
	// Account A will be called twice due to cache eviction
	mockVault.EXPECT().GetTotalDeposit(gomock.Any(), accountA).Return(big.NewInt(8000), nil).Times(2)
	mockVault.EXPECT().GetTotalDeposit(gomock.Any(), accountB).Return(big.NewInt(5000), nil)
	mockVault.EXPECT().GetTotalDeposit(gomock.Any(), accountC).Return(big.NewInt(3000), nil)

	// Create paymentValidator with small LRU cache size to force eviction
	paymentValidator, err := ondemand.NewOnDemandPaymentValidator(
		ctx,
		testutils.GetLogger(),
		2,
		mockVault,
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