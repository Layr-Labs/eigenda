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

func TestDebitMultipleAccounts(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	tableName := createPaymentTable(t, "TestDebitMultipleAccounts")
	defer deleteTable(t, tableName)

	accountA := gethcommon.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	accountB := gethcommon.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")

	testVault := vault.NewTestPaymentVault()
	testVault.SetTotalDeposit(accountA, big.NewInt(10000))
	testVault.SetTotalDeposit(accountB, big.NewInt(20000))

	config, err := ondemandvalidation.NewOnDemandLedgerCacheConfig(
		10,
		tableName,
		time.Second,
	)
	require.NoError(t, err)

	paymentValidator, err := ondemandvalidation.NewOnDemandPaymentValidator(
		ctx,
		test.GetLogger(),
		config,
		testVault,
		dynamoClient,
		nil,
		nil,
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
}

func TestDebitInsufficientFunds(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	tableName := createPaymentTable(t, "TestDebitInsufficientFunds")
	defer deleteTable(t, tableName)

	accountID := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")

	testVault := vault.NewTestPaymentVault()
	testVault.SetPricePerSymbol(1000)
	testVault.SetTotalDeposit(accountID, big.NewInt(5000))

	config, err := ondemandvalidation.NewOnDemandLedgerCacheConfig(
		10,
		tableName,
		time.Second,
	)
	require.NoError(t, err)

	paymentValidator, err := ondemandvalidation.NewOnDemandPaymentValidator(
		ctx,
		test.GetLogger(),
		config,
		testVault,
		dynamoClient,
		nil,
		nil,
	)
	require.NoError(t, err)

	// Try to debit more than available funds (5000 wei / 1000 wei per symbol = 5 symbols max)
	err = paymentValidator.Debit(ctx, accountID, uint32(10), []uint8{0})
	require.Error(t, err, "debit should fail when insufficient funds")
	var insufficientFundsErr *ondemand.InsufficientFundsError
	require.ErrorAs(t, err, &insufficientFundsErr, "error should be InsufficientFundsError")
}
