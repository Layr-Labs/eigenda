package ondemand_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/core/payments/ondemand"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestConstructor(t *testing.T) {
	tableName := "TestConstructor"
	accountID := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")

	store, err := ondemand.NewCumulativePaymentStore(nil, tableName, accountID)
	require.Error(t, err, "nil client should error")
	require.Nil(t, store)

	store, err = ondemand.NewCumulativePaymentStore(dynamoClient, "", accountID)
	require.Error(t, err, "empty table name should error")
	require.Nil(t, store)

	store, err = ondemand.NewCumulativePaymentStore(dynamoClient, tableName, gethcommon.Address{})
	require.Error(t, err, "zero address should error")
	require.Nil(t, store)
}

func TestStoreCumulativePaymentInputValidation(t *testing.T) {
	tableName := createPaymentTable(t, "StoreInputValidation")
	defer deleteTable(t, tableName)
	
	accountID := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")
	store, err := ondemand.NewCumulativePaymentStore(dynamoClient, tableName, accountID)
	require.NoError(t, err)

	ctx := context.Background()

	err = store.StoreCumulativePayment(ctx, nil)
	require.Error(t, err, "nil amount should error")

	err = store.StoreCumulativePayment(ctx, big.NewInt(-100))
	require.Error(t, err, "negative amount should error")
}

func TestStoreThenGet(t *testing.T) {
	tableName := createPaymentTable(t, "StoreThenGet")
	defer deleteTable(t, tableName)

	accountID := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")
	store, err := ondemand.NewCumulativePaymentStore(dynamoClient, tableName, accountID)
	require.NoError(t, err)
	ctx := context.Background()

	value, err := store.GetCumulativePayment(ctx)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(0), value, "get when missing should return 0")

	require.NoError(t, store.StoreCumulativePayment(ctx, big.NewInt(100)))
	value, err = store.GetCumulativePayment(ctx)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(100), value)

	require.NoError(t, store.StoreCumulativePayment(ctx, big.NewInt(200)))
	value, err = store.GetCumulativePayment(ctx)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(200), value)

	require.NoError(t, store.StoreCumulativePayment(ctx, big.NewInt(50)))
	value, err = store.GetCumulativePayment(ctx)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(50), value)

}

func TestDifferentAddresses(t *testing.T) {
	tableName := createPaymentTable(t, "DifferentAddresses")
	defer deleteTable(t, tableName)

	accountA := gethcommon.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	accountB := gethcommon.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")

	storeA, err := ondemand.NewCumulativePaymentStore(dynamoClient, tableName, accountA)
	require.NoError(t, err)
	storeB, err := ondemand.NewCumulativePaymentStore(dynamoClient, tableName, accountB)
	require.NoError(t, err)

	ctx := context.Background()
	require.NoError(t, storeA.StoreCumulativePayment(ctx, big.NewInt(100)))
	require.NoError(t, storeB.StoreCumulativePayment(ctx, big.NewInt(300)))

	valueA, err := storeA.GetCumulativePayment(ctx)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(100), valueA)

	valueB, err := storeB.GetCumulativePayment(ctx)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(300), valueB)
}
