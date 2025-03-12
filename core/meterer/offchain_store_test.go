package meterer_test

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"testing"
	"time"

	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testContext struct {
	ctx              context.Context
	store            meterer.OffchainStore
	reservationTable string
	onDemandTable    string
	globalBinTable   string
}

// setupTest creates a test context with tables created and cleaned up after the test
func setupTest(t *testing.T) *testContext {
	tc := &testContext{
		ctx:              context.Background(),
		reservationTable: fmt.Sprintf("reservation_test_%d", rand.Int()),
		onDemandTable:    fmt.Sprintf("ondemand_test_%d", rand.Int()),
		globalBinTable:   fmt.Sprintf("global_bin_test_%d", rand.Int()),
	}

	var err error

	// Create the tables
	err = meterer.CreateReservationTable(clientConfig, tc.reservationTable)
	require.NoError(t, err)

	err = meterer.CreateOnDemandTable(clientConfig, tc.onDemandTable)
	require.NoError(t, err)

	err = meterer.CreateGlobalReservationTable(clientConfig, tc.globalBinTable)
	require.NoError(t, err)

	// Register cleanup to remove tables after test completes
	t.Cleanup(func() {
		cleanupTables(tc)
	})

	// Create the OffchainStore
	tc.store, err = meterer.NewOffchainStore(
		clientConfig,
		tc.reservationTable,
		tc.onDemandTable,
		tc.globalBinTable,
		nil, // Logger not needed for test
	)
	require.NoError(t, err)

	return tc
}

// cleanupTables removes all tables created for a test
func cleanupTables(tc *testContext) {
	_ = dynamoClient.DeleteTable(tc.ctx, tc.reservationTable)
	_ = dynamoClient.DeleteTable(tc.ctx, tc.onDemandTable)
	_ = dynamoClient.DeleteTable(tc.ctx, tc.globalBinTable)
}

// TestUpdateReservationBin tests the UpdateReservationBin function
func TestUpdateReservationBin(t *testing.T) {
	tc := setupTest(t)

	// Test updating bin that doesn't exist yet (should create it)
	accountID := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")
	reservationPeriod := uint64(1)
	size := uint64(1000)

	binUsage, err := tc.store.UpdateReservationBin(tc.ctx, accountID, reservationPeriod, size)
	require.NoError(t, err)
	assert.Equal(t, size, binUsage)

	// Get the bin directly from DynamoDB to verify
	item, err := dynamoClient.GetItem(tc.ctx, tc.reservationTable, commondynamodb.Key{
		"AccountID":         &types.AttributeValueMemberS{Value: accountID.Hex()},
		"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(reservationPeriod, 10)},
	})
	require.NoError(t, err)
	binUsageStr := item["BinUsage"].(*types.AttributeValueMemberN).Value
	binUsageVal, err := strconv.ParseUint(binUsageStr, 10, 64)
	require.NoError(t, err)
	assert.Equal(t, size, binUsageVal)

	// Test updating existing bin
	additionalSize := uint64(500)
	binUsage, err = tc.store.UpdateReservationBin(tc.ctx, accountID, reservationPeriod, additionalSize)
	require.NoError(t, err)
	assert.Equal(t, size+additionalSize, binUsage)

	// Verify updated bin
	item, err = dynamoClient.GetItem(tc.ctx, tc.reservationTable, commondynamodb.Key{
		"AccountID":         &types.AttributeValueMemberS{Value: accountID.Hex()},
		"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(reservationPeriod, 10)},
	})
	require.NoError(t, err)
	binUsageStr = item["BinUsage"].(*types.AttributeValueMemberN).Value
	binUsageVal, err = strconv.ParseUint(binUsageStr, 10, 64)
	require.NoError(t, err)
	assert.Equal(t, size+additionalSize, binUsageVal)
}

// TestUpdateGlobalBin tests the UpdateGlobalBin function
func TestUpdateGlobalBin(t *testing.T) {
	tc := setupTest(t)

	// Test updating global bin that doesn't exist yet (should create it)
	reservationPeriod := uint64(1)
	size := uint64(2000)

	binUsage, err := tc.store.UpdateGlobalBin(tc.ctx, reservationPeriod, size)
	require.NoError(t, err)
	assert.Equal(t, size, binUsage)

	// Get the bin directly from DynamoDB to verify
	item, err := dynamoClient.GetItem(tc.ctx, tc.globalBinTable, commondynamodb.Key{
		"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(reservationPeriod, 10)},
	})
	require.NoError(t, err)
	binUsageStr := item["BinUsage"].(*types.AttributeValueMemberN).Value
	binUsageVal, err := strconv.ParseUint(binUsageStr, 10, 64)
	require.NoError(t, err)
	assert.Equal(t, size, binUsageVal)

	// Test updating existing bin
	additionalSize := uint64(1000)
	binUsage, err = tc.store.UpdateGlobalBin(tc.ctx, reservationPeriod, additionalSize)
	require.NoError(t, err)
	assert.Equal(t, size+additionalSize, binUsage)

	// Verify updated bin
	item, err = dynamoClient.GetItem(tc.ctx, tc.globalBinTable, commondynamodb.Key{
		"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(reservationPeriod, 10)},
	})
	require.NoError(t, err)
	binUsageStr = item["BinUsage"].(*types.AttributeValueMemberN).Value
	binUsageVal, err = strconv.ParseUint(binUsageStr, 10, 64)
	require.NoError(t, err)
	assert.Equal(t, size+additionalSize, binUsageVal)
}

// TestAddOnDemandPayment tests the AddOnDemandPayment function
func TestAddOnDemandPayment(t *testing.T) {
	tc := setupTest(t)

	accountID := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")
	payment := core.PaymentMetadata{
		AccountID:         accountID,
		Timestamp:         time.Now().Unix(),
		CumulativePayment: big.NewInt(100),
	}
	charge := big.NewInt(100)

	// Add the payment
	err := tc.store.AddOnDemandPayment(tc.ctx, payment, charge)
	require.NoError(t, err)

	// Verify the payment was added - using the exact key structure expected by DynamoDB
	item, err := dynamoClient.GetItem(tc.ctx, tc.onDemandTable, commondynamodb.Key{
		"AccountID":          &types.AttributeValueMemberS{Value: accountID.Hex()},
		"CumulativePayments": &types.AttributeValueMemberN{Value: payment.CumulativePayment.String()},
	})
	require.NoError(t, err)
	require.NotNil(t, item, "Item should exist in the table")

	// Verify the PaymentCharged field
	paymentChargedStr := item["PaymentCharged"].(*types.AttributeValueMemberN).Value
	paymentChargedVal, err := strconv.ParseInt(paymentChargedStr, 10, 64)
	require.NoError(t, err)
	assert.Equal(t, charge.Int64(), paymentChargedVal)

	// Attempt to add the same payment again, should return an error
	err = tc.store.AddOnDemandPayment(tc.ctx, payment, charge)
	require.Error(t, err)
	assert.Equal(t, "exact payment already exists", err.Error())
}

// TestRemoveOnDemandPayment tests the RemoveOnDemandPayment function
func TestRemoveOnDemandPayment(t *testing.T) {
	tc := setupTest(t)

	// Create and add a payment first
	accountID := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")
	cumulativePayment := big.NewInt(1000)
	paymentCharged := big.NewInt(500)

	paymentMetadata := core.PaymentMetadata{
		AccountID:         accountID,
		Timestamp:         time.Now().Unix(),
		CumulativePayment: cumulativePayment,
	}

	err := tc.store.AddOnDemandPayment(tc.ctx, paymentMetadata, paymentCharged)
	require.NoError(t, err)

	// Verify the payment was added before removal
	item, err := dynamoClient.GetItem(tc.ctx, tc.onDemandTable, commondynamodb.Key{
		"AccountID":          &types.AttributeValueMemberS{Value: accountID.Hex()},
		"CumulativePayments": &types.AttributeValueMemberN{Value: cumulativePayment.String()},
	})
	require.NoError(t, err)
	require.NotNil(t, item, "Item should exist before removal")

	// Test removing the payment
	err = tc.store.RemoveOnDemandPayment(tc.ctx, accountID, cumulativePayment)
	require.NoError(t, err)

	// Verify the payment was removed
	item, err = dynamoClient.GetItem(tc.ctx, tc.onDemandTable, commondynamodb.Key{
		"AccountID":          &types.AttributeValueMemberS{Value: accountID.Hex()},
		"CumulativePayments": &types.AttributeValueMemberN{Value: cumulativePayment.String()},
	})
	require.NoError(t, err)
	assert.Nil(t, item, "Item should be deleted")

	// Removing non-existent payment should work (not error)
	err = tc.store.RemoveOnDemandPayment(tc.ctx, accountID, big.NewInt(9999))
	require.NoError(t, err)
}

// TestGetRelevantOnDemandRecords tests the GetRelevantOnDemandRecords function
func TestGetRelevantOnDemandRecords(t *testing.T) {
	tc := setupTest(t)

	// Create and add multiple payments for the same account but with different cumulative payments
	accountID := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")

	// Helper function to create payment metadata
	createPayment := func(amount int64) (core.PaymentMetadata, *big.Int) {
		return core.PaymentMetadata{
			AccountID:         accountID,
			Timestamp:         time.Now().Unix(),
			CumulativePayment: big.NewInt(amount),
		}, big.NewInt(amount)
	}

	// Add payments with cumulative values 100, 300, 600
	payments := []int64{100, 300, 600}
	for i, amt := range payments {
		payment, charge := createPayment(amt)
		if i > 0 {
			// Each charge is the difference between this payment and the previous one
			charge = big.NewInt(amt - payments[i-1])
		}
		err := tc.store.AddOnDemandPayment(tc.ctx, payment, charge)
		require.NoError(t, err)
	}

	// Verify each payment was added correctly
	for _, amt := range payments {
		item, err := dynamoClient.GetItem(tc.ctx, tc.onDemandTable, commondynamodb.Key{
			"AccountID":          &types.AttributeValueMemberS{Value: accountID.Hex()},
			"CumulativePayments": &types.AttributeValueMemberN{Value: strconv.FormatInt(amt, 10)},
		})
		require.NoError(t, err)
		require.NotNil(t, item, fmt.Sprintf("Payment with cumulative amount %d should exist", amt))
	}

	// Test getting relevant records for different cumulative payment values
	testCases := []struct {
		requestedPayment int64
		expectedPrev     int64
		expectedNext     int64
		expectedCharge   int64
	}{
		{200, 100, 300, 200}, // Between 100 and 300
		{400, 300, 600, 300}, // Between 300 and 600
		{700, 600, 0, 0},     // Larger than any payment
		{50, 0, 100, 100},    // Smaller than any payment
	}

	for _, testCase := range testCases {
		cumulativePayment := big.NewInt(testCase.requestedPayment)
		prevPmt, nextPmt, chargeOfNextPmt, err := tc.store.GetRelevantOnDemandRecords(tc.ctx, accountID, cumulativePayment)
		require.NoError(t, err)

		assert.Equal(t, big.NewInt(testCase.expectedPrev), prevPmt)
		assert.Equal(t, big.NewInt(testCase.expectedNext), nextPmt)
		assert.Equal(t, big.NewInt(testCase.expectedCharge), chargeOfNextPmt)
	}
}
