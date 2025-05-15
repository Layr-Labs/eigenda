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
	quorumNumber := uint8(0)

	binUsage, err := tc.store.UpdateReservationBin(tc.ctx, accountID, reservationPeriod, size, quorumNumber)
	require.NoError(t, err)
	assert.Equal(t, size, binUsage)

	// Get the bin directly from DynamoDB to verify
	item, err := dynamoClient.GetItem(tc.ctx, tc.reservationTable, commondynamodb.Key{
		"AccountID":         &types.AttributeValueMemberS{Value: accountID.Hex()},
		"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(reservationPeriod, 10)},
		"QuorumNumber":      &types.AttributeValueMemberN{Value: strconv.FormatUint(uint64(quorumNumber), 10)},
	})
	require.NoError(t, err)
	binUsageStr := item["BinUsage"].(*types.AttributeValueMemberN).Value
	binUsageVal, err := strconv.ParseUint(binUsageStr, 10, 64)
	require.NoError(t, err)
	assert.Equal(t, size, binUsageVal)

	// Test updating existing bin
	additionalSize := uint64(500)
	binUsage, err = tc.store.UpdateReservationBin(tc.ctx, accountID, reservationPeriod, additionalSize, quorumNumber)
	require.NoError(t, err)
	assert.Equal(t, size+additionalSize, binUsage)

	// Verify updated bin
	item, err = dynamoClient.GetItem(tc.ctx, tc.reservationTable, commondynamodb.Key{
		"AccountID":         &types.AttributeValueMemberS{Value: accountID.Hex()},
		"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(reservationPeriod, 10)},
		"QuorumNumber":      &types.AttributeValueMemberN{Value: strconv.FormatUint(uint64(quorumNumber), 10)},
	})
	require.NoError(t, err)
	binUsageStr = item["BinUsage"].(*types.AttributeValueMemberN).Value
	binUsageVal, err = strconv.ParseUint(binUsageStr, 10, 64)
	require.NoError(t, err)
	assert.Equal(t, size+additionalSize, binUsageVal)
	
	// Test with a different quorum
	quorumNumber2 := uint8(1)
	size2 := uint64(2000)
	binUsage, err = tc.store.UpdateReservationBin(tc.ctx, accountID, reservationPeriod, size2, quorumNumber2)
	require.NoError(t, err)
	assert.Equal(t, size2, binUsage)
	
	// Verify second quorum bin
	item, err = dynamoClient.GetItem(tc.ctx, tc.reservationTable, commondynamodb.Key{
		"AccountID":         &types.AttributeValueMemberS{Value: accountID.Hex()},
		"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(reservationPeriod, 10)},
		"QuorumNumber":      &types.AttributeValueMemberN{Value: strconv.FormatUint(uint64(quorumNumber2), 10)},
	})
	require.NoError(t, err)
	binUsageStr = item["BinUsage"].(*types.AttributeValueMemberN).Value
	binUsageVal, err = strconv.ParseUint(binUsageStr, 10, 64)
	require.NoError(t, err)
	assert.Equal(t, size2, binUsageVal)
	
	// First quorum should remain unchanged
	item, err = dynamoClient.GetItem(tc.ctx, tc.reservationTable, commondynamodb.Key{
		"AccountID":         &types.AttributeValueMemberS{Value: accountID.Hex()},
		"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(reservationPeriod, 10)},
		"QuorumNumber":      &types.AttributeValueMemberN{Value: strconv.FormatUint(uint64(quorumNumber), 10)},
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
	payment1 := core.PaymentMetadata{
		AccountID:         accountID,
		Timestamp:         time.Now().Unix(),
		CumulativePayment: big.NewInt(100),
	}
	charge1 := big.NewInt(100)

	// Add the payment
	oldPayment, err := tc.store.AddOnDemandPayment(tc.ctx, payment1, charge1)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(0), oldPayment, "Old payment should be 0 for first payment")

	// Verify the payment was added with the correct structure
	item, err := dynamoClient.GetItem(tc.ctx, tc.onDemandTable, commondynamodb.Key{
		"AccountID": &types.AttributeValueMemberS{Value: accountID.Hex()},
	})
	require.NoError(t, err)
	require.NotNil(t, item, "Item should exist in the table")

	// Verify the CumulativePayment field
	cumulativePaymentStr := item["CumulativePayment"].(*types.AttributeValueMemberN).Value
	cumulativePaymentVal, err := strconv.ParseInt(cumulativePaymentStr, 10, 64)
	require.NoError(t, err)
	assert.Equal(t, payment1.CumulativePayment.Int64(), cumulativePaymentVal)

	// Test case: Add a larger payment with sufficient increment
	payment2 := core.PaymentMetadata{
		AccountID:         accountID,
		Timestamp:         time.Now().Unix(),
		CumulativePayment: big.NewInt(200),
	}
	charge2 := big.NewInt(100) // The same charge is fine because 200-100=100 >= 100

	oldPayment, err = tc.store.AddOnDemandPayment(tc.ctx, payment2, charge2)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(100), oldPayment, "Old payment should be 100")

	// Verify the payment was updated
	item, err = dynamoClient.GetItem(tc.ctx, tc.onDemandTable, commondynamodb.Key{
		"AccountID": &types.AttributeValueMemberS{Value: accountID.Hex()},
	})
	require.NoError(t, err)
	cumulativePaymentStr = item["CumulativePayment"].(*types.AttributeValueMemberN).Value
	cumulativePaymentVal, err = strconv.ParseInt(cumulativePaymentStr, 10, 64)
	require.NoError(t, err)
	assert.Equal(t, payment2.CumulativePayment.Int64(), cumulativePaymentVal)

	// Test case: Add a larger payment but with insufficient increment
	payment3 := core.PaymentMetadata{
		AccountID:         accountID,
		Timestamp:         time.Now().Unix(),
		CumulativePayment: big.NewInt(250), // Only 50 more than previous 200
	}
	charge3 := big.NewInt(100) // But we need a minimum increment of 100

	oldPayment, err = tc.store.AddOnDemandPayment(tc.ctx, payment3, charge3)
	require.Error(t, err) // Should fail due to insufficient increment
	assert.Contains(t, err.Error(), "insufficient cumulative payment increment")
	require.Nil(t, oldPayment, "Old payment should be nil on error")

	// Verify the payment wasn't updated
	item, err = dynamoClient.GetItem(tc.ctx, tc.onDemandTable, commondynamodb.Key{
		"AccountID": &types.AttributeValueMemberS{Value: accountID.Hex()},
	})
	require.NoError(t, err)
	cumulativePaymentStr = item["CumulativePayment"].(*types.AttributeValueMemberN).Value
	cumulativePaymentVal, err = strconv.ParseInt(cumulativePaymentStr, 10, 64)
	require.NoError(t, err)
	assert.Equal(t, payment2.CumulativePayment.Int64(), cumulativePaymentVal, "Payment should not have been updated")

	// Test case: Add a smaller payment (should fail)
	payment4 := core.PaymentMetadata{
		AccountID:         accountID,
		Timestamp:         time.Now().Unix(),
		CumulativePayment: big.NewInt(150),
	}
	charge4 := big.NewInt(50)

	oldPayment, err = tc.store.AddOnDemandPayment(tc.ctx, payment4, charge4)
	require.Error(t, err) // Should fail since payment is smaller than current
	assert.Contains(t, err.Error(), "insufficient cumulative payment increment")
	require.Nil(t, oldPayment, "Old payment should be nil on error")

	// Verify the payment wasn't updated
	item, err = dynamoClient.GetItem(tc.ctx, tc.onDemandTable, commondynamodb.Key{
		"AccountID": &types.AttributeValueMemberS{Value: accountID.Hex()},
	})
	require.NoError(t, err)
	cumulativePaymentStr = item["CumulativePayment"].(*types.AttributeValueMemberN).Value
	cumulativePaymentVal, err = strconv.ParseInt(cumulativePaymentStr, 10, 64)
	require.NoError(t, err)
	assert.Equal(t, payment2.CumulativePayment.Int64(), cumulativePaymentVal, "Payment should not have been updated")
}

// TestRollbackOnDemandPayment tests the RollbackOnDemandPayment function
func TestRollbackOnDemandPayment(t *testing.T) {
	tc := setupTest(t)

	// Create and add a payment
	accountID := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")
	cumulativePayment := big.NewInt(1000)
	paymentCharged := big.NewInt(500)

	paymentMetadata := core.PaymentMetadata{
		AccountID:         accountID,
		Timestamp:         time.Now().Unix(),
		CumulativePayment: cumulativePayment,
	}

	oldPayment, err := tc.store.AddOnDemandPayment(tc.ctx, paymentMetadata, paymentCharged)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(0), oldPayment, "Old payment should be 0 for first payment")

	// Verify the payment was added
	item, err := dynamoClient.GetItem(tc.ctx, tc.onDemandTable, commondynamodb.Key{
		"AccountID": &types.AttributeValueMemberS{Value: accountID.Hex()},
	})
	require.NoError(t, err)
	require.NotNil(t, item, "Item should exist in the table")

	// Add another payment
	newCumulativePayment := big.NewInt(2000)
	newPaymentMetadata := core.PaymentMetadata{
		AccountID:         accountID,
		Timestamp:         time.Now().Unix(),
		CumulativePayment: newCumulativePayment,
	}
	newPaymentCharged := big.NewInt(1000)

	oldPayment, err = tc.store.AddOnDemandPayment(tc.ctx, newPaymentMetadata, newPaymentCharged)
	require.NoError(t, err)
	require.Equal(t, cumulativePayment, oldPayment, "Old payment should be 1000 for second payment")

	// Test case 1: Rollback to previous payment
	err = tc.store.RollbackOnDemandPayment(tc.ctx, accountID, newCumulativePayment, oldPayment)
	require.NoError(t, err)

	// Verify the payment was rolled back
	item, err = dynamoClient.GetItem(tc.ctx, tc.onDemandTable, commondynamodb.Key{
		"AccountID": &types.AttributeValueMemberS{Value: accountID.Hex()},
	})
	require.NoError(t, err)
	require.NotNil(t, item, "Item should still exist in the table")

	cumulativePaymentStr := item["CumulativePayment"].(*types.AttributeValueMemberN).Value
	cumulativePaymentVal, err := strconv.ParseInt(cumulativePaymentStr, 10, 64)
	require.NoError(t, err)
	assert.Equal(t, oldPayment.Int64(), cumulativePaymentVal, "Payment should be rolled back to 1000")

	// Test case 2: Rollback to a different value directly
	// The value will be updated regardless of what the current value is
	err = tc.store.RollbackOnDemandPayment(tc.ctx, accountID, big.NewInt(1000), big.NewInt(500))
	require.NoError(t, err)

	// Verify the payment was updated to the new value
	item, err = dynamoClient.GetItem(tc.ctx, tc.onDemandTable, commondynamodb.Key{
		"AccountID": &types.AttributeValueMemberS{Value: accountID.Hex()},
	})
	require.NoError(t, err)
	require.NotNil(t, item, "Item should still exist in the table")

	cumulativePaymentStr = item["CumulativePayment"].(*types.AttributeValueMemberN).Value
	cumulativePaymentVal, err = strconv.ParseInt(cumulativePaymentStr, 10, 64)
	require.NoError(t, err)
	assert.Equal(t, int64(500), cumulativePaymentVal, "Payment should be set to 500 regardless of current value")

	// Test case 3: Rollback to zero (should delete the record)
	err = tc.store.RollbackOnDemandPayment(tc.ctx, accountID, big.NewInt(500), big.NewInt(0))
	require.NoError(t, err)

	// payment is set back to 0
	largest, err := tc.store.GetLargestCumulativePayment(tc.ctx, accountID)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(0), largest, "Payment should be set to 0")

	// Test case 4: Trying to rollback non-matching payment should not cause an error
	err = tc.store.RollbackOnDemandPayment(tc.ctx, accountID, big.NewInt(9999), big.NewInt(500))
	require.NoError(t, err)
}

// TestGetPeriodRecords tests the GetPeriodRecords function with quorum-specific retrieval
func TestGetPeriodRecords(t *testing.T) {
	tc := setupTest(t)

	// Setup: Create several period records with different quorums
	accountID := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")
	reservationPeriod1 := uint64(100)
	reservationPeriod2 := uint64(101)
	reservationPeriod3 := uint64(102)
	
	// Create data for quorum 0
	quorum0 := uint8(0)
	size1Quorum0 := uint64(1000)
	size2Quorum0 := uint64(1500)
	size3Quorum0 := uint64(2000)
	
	// Create data for quorum 1
	quorum1 := uint8(1)
	size1Quorum1 := uint64(3000)
	size2Quorum1 := uint64(3500)
	size3Quorum1 := uint64(4000)
	
	// Add records for quorum 0
	_, err := tc.store.UpdateReservationBin(tc.ctx, accountID, reservationPeriod1, size1Quorum0, quorum0)
	require.NoError(t, err)
	_, err = tc.store.UpdateReservationBin(tc.ctx, accountID, reservationPeriod2, size2Quorum0, quorum0)
	require.NoError(t, err)
	_, err = tc.store.UpdateReservationBin(tc.ctx, accountID, reservationPeriod3, size3Quorum0, quorum0)
	require.NoError(t, err)
	
	// Add records for quorum 1
	_, err = tc.store.UpdateReservationBin(tc.ctx, accountID, reservationPeriod1, size1Quorum1, quorum1)
	require.NoError(t, err)
	_, err = tc.store.UpdateReservationBin(tc.ctx, accountID, reservationPeriod2, size2Quorum1, quorum1)
	require.NoError(t, err)
	_, err = tc.store.UpdateReservationBin(tc.ctx, accountID, reservationPeriod3, size3Quorum1, quorum1)
	require.NoError(t, err)
	
	// Test single quorum retrieval - quorum 0
	records, err := tc.store.GetPeriodRecords(tc.ctx, accountID, reservationPeriod1, quorum0)
	require.NoError(t, err)
	require.Equal(t, 3, len(records))
	
	// Verify records for quorum 0
	assert.Equal(t, uint32(reservationPeriod1), records[0].Index)
	assert.Equal(t, size1Quorum0, records[0].Usage)
	assert.Equal(t, uint32(quorum0), records[0].QuorumNumber)
	
	assert.Equal(t, uint32(reservationPeriod2), records[1].Index)
	assert.Equal(t, size2Quorum0, records[1].Usage)
	assert.Equal(t, uint32(quorum0), records[1].QuorumNumber)
	
	assert.Equal(t, uint32(reservationPeriod3), records[2].Index)
	assert.Equal(t, size3Quorum0, records[2].Usage)
	assert.Equal(t, uint32(quorum0), records[2].QuorumNumber)
	
	// Test single quorum retrieval - quorum 1
	records, err = tc.store.GetPeriodRecords(tc.ctx, accountID, reservationPeriod1, quorum1)
	require.NoError(t, err)
	require.Equal(t, 3, len(records))
	
	// Verify records for quorum 1
	assert.Equal(t, uint32(reservationPeriod1), records[0].Index)
	assert.Equal(t, size1Quorum1, records[0].Usage)
	assert.Equal(t, uint32(quorum1), records[0].QuorumNumber)
	
	assert.Equal(t, uint32(reservationPeriod2), records[1].Index)
	assert.Equal(t, size2Quorum1, records[1].Usage)
	assert.Equal(t, uint32(quorum1), records[1].QuorumNumber)
	
	assert.Equal(t, uint32(reservationPeriod3), records[2].Index)
	assert.Equal(t, size3Quorum1, records[2].Usage)
	assert.Equal(t, uint32(quorum1), records[2].QuorumNumber)
}

// TestGetPeriodRecordsMultiQuorum tests the GetPeriodRecordsMultiQuorum function
func TestGetPeriodRecordsMultiQuorum(t *testing.T) {
	tc := setupTest(t)

	// Setup: Create several period records with different quorums
	accountID := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")
	reservationPeriod1 := uint64(100)
	reservationPeriod2 := uint64(101)
	
	// Create data for quorum 0, 1, and 2
	quorum0 := uint8(0)
	quorum1 := uint8(1)
	quorum2 := uint8(2)
	
	// Add records
	_, err := tc.store.UpdateReservationBin(tc.ctx, accountID, reservationPeriod1, 1000, quorum0)
	require.NoError(t, err)
	_, err = tc.store.UpdateReservationBin(tc.ctx, accountID, reservationPeriod1, 2000, quorum1)
	require.NoError(t, err)
	_, err = tc.store.UpdateReservationBin(tc.ctx, accountID, reservationPeriod1, 3000, quorum2)
	require.NoError(t, err)
	_, err = tc.store.UpdateReservationBin(tc.ctx, accountID, reservationPeriod2, 1500, quorum0)
	require.NoError(t, err)
	_, err = tc.store.UpdateReservationBin(tc.ctx, accountID, reservationPeriod2, 2500, quorum1)
	require.NoError(t, err)
	_, err = tc.store.UpdateReservationBin(tc.ctx, accountID, reservationPeriod2, 3500, quorum2)
	require.NoError(t, err)
	
	// Test multi-quorum retrieval
	records, err := tc.store.GetPeriodRecordsMultiQuorum(tc.ctx, accountID, reservationPeriod1, []uint8{quorum0, quorum1})
	require.NoError(t, err)
	require.Equal(t, 4, len(records), "Should retrieve 4 records (2 periods × 2 quorums)")
	
	// Create a map to verify all expected records are present
	recordMap := make(map[string]bool)
	for _, record := range records {
		key := fmt.Sprintf("%d-%d", record.Index, record.QuorumNumber)
		recordMap[key] = true
		
		// Verify the record values
		if record.Index == uint32(reservationPeriod1) && record.QuorumNumber == uint32(quorum0) {
			assert.Equal(t, uint64(1000), record.Usage)
		} else if record.Index == uint32(reservationPeriod1) && record.QuorumNumber == uint32(quorum1) {
			assert.Equal(t, uint64(2000), record.Usage)
		} else if record.Index == uint32(reservationPeriod2) && record.QuorumNumber == uint32(quorum0) {
			assert.Equal(t, uint64(1500), record.Usage)
		} else if record.Index == uint32(reservationPeriod2) && record.QuorumNumber == uint32(quorum1) {
			assert.Equal(t, uint64(2500), record.Usage)
		}
	}
	
	// Check that all expected records are present
	assert.True(t, recordMap[fmt.Sprintf("%d-%d", reservationPeriod1, quorum0)])
	assert.True(t, recordMap[fmt.Sprintf("%d-%d", reservationPeriod1, quorum1)])
	assert.True(t, recordMap[fmt.Sprintf("%d-%d", reservationPeriod2, quorum0)])
	assert.True(t, recordMap[fmt.Sprintf("%d-%d", reservationPeriod2, quorum1)])
	
	// Quorum2 records should not be present
	assert.False(t, recordMap[fmt.Sprintf("%d-%d", reservationPeriod1, quorum2)])
	assert.False(t, recordMap[fmt.Sprintf("%d-%d", reservationPeriod2, quorum2)])
	
	// Test with single quorum using multi-quorum API
	records, err = tc.store.GetPeriodRecordsMultiQuorum(tc.ctx, accountID, reservationPeriod1, []uint8{quorum2})
	require.NoError(t, err)
	require.Equal(t, 2, len(records), "Should retrieve 2 records (2 periods × 1 quorum)")
	
	// Verify the records
	for _, record := range records {
		assert.Equal(t, uint32(quorum2), record.QuorumNumber)
		if record.Index == uint32(reservationPeriod1) {
			assert.Equal(t, uint64(3000), record.Usage)
		} else if record.Index == uint32(reservationPeriod2) {
			assert.Equal(t, uint64(3500), record.Usage)
		}
	}
	
	// Test with empty quorum list
	records, err = tc.store.GetPeriodRecordsMultiQuorum(tc.ctx, accountID, reservationPeriod1, []uint8{})
	require.NoError(t, err)
	require.Equal(t, 0, len(records), "Should retrieve 0 records with empty quorum list")
}

// TestBatchUpdateReservationBins tests the BatchUpdateReservationBins function
func TestBatchUpdateReservationBins(t *testing.T) {
	tc := setupTest(t)
	
	// Create test data
	accountID := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")
	updates := []meterer.ReservationBinUpdate{
		{
			AccountID:         accountID,
			ReservationPeriod: 100,
			Size:              1000,
			QuorumNumber:      0,
		},
		{
			AccountID:         accountID,
			ReservationPeriod: 101,
			Size:              2000, 
			QuorumNumber:      0,
		},
		{
			AccountID:         accountID,
			ReservationPeriod: 100,
			Size:              3000,
			QuorumNumber:      1,
		},
	}
	
	// Perform batch update
	results, errors := tc.store.BatchUpdateReservationBins(tc.ctx, updates)
	require.Equal(t, 0, len(errors), "Should have no errors")
	require.Equal(t, 3, len(results), "Should have 3 results")
	
	// Verify results match expected sizes
	assert.Equal(t, uint64(1000), results[0])
	assert.Equal(t, uint64(2000), results[1])
	assert.Equal(t, uint64(3000), results[2])
	
	// Verify records in database
	item, err := dynamoClient.GetItem(tc.ctx, tc.reservationTable, commondynamodb.Key{
		"AccountID":         &types.AttributeValueMemberS{Value: accountID.Hex()},
		"ReservationPeriod": &types.AttributeValueMemberN{Value: "100"},
		"QuorumNumber":      &types.AttributeValueMemberN{Value: "0"},
	})
	require.NoError(t, err)
	binUsageStr := item["BinUsage"].(*types.AttributeValueMemberN).Value
	binUsageVal, err := strconv.ParseUint(binUsageStr, 10, 64)
	require.NoError(t, err)
	assert.Equal(t, uint64(1000), binUsageVal)
	
	item, err = dynamoClient.GetItem(tc.ctx, tc.reservationTable, commondynamodb.Key{
		"AccountID":         &types.AttributeValueMemberS{Value: accountID.Hex()},
		"ReservationPeriod": &types.AttributeValueMemberN{Value: "101"},
		"QuorumNumber":      &types.AttributeValueMemberN{Value: "0"},
	})
	require.NoError(t, err)
	binUsageStr = item["BinUsage"].(*types.AttributeValueMemberN).Value
	binUsageVal, err = strconv.ParseUint(binUsageStr, 10, 64)
	require.NoError(t, err)
	assert.Equal(t, uint64(2000), binUsageVal)
	
	item, err = dynamoClient.GetItem(tc.ctx, tc.reservationTable, commondynamodb.Key{
		"AccountID":         &types.AttributeValueMemberS{Value: accountID.Hex()},
		"ReservationPeriod": &types.AttributeValueMemberN{Value: "100"},
		"QuorumNumber":      &types.AttributeValueMemberN{Value: "1"},
	})
	require.NoError(t, err)
	binUsageStr = item["BinUsage"].(*types.AttributeValueMemberN).Value
	binUsageVal, err = strconv.ParseUint(binUsageStr, 10, 64)
	require.NoError(t, err)
	assert.Equal(t, uint64(3000), binUsageVal)
	
	// Test updating existing bins
	updates = []meterer.ReservationBinUpdate{
		{
			AccountID:         accountID,
			ReservationPeriod: 100,
			Size:              500,
			QuorumNumber:      0,
		},
		{
			AccountID:         accountID,
			ReservationPeriod: 100,
			Size:              700,
			QuorumNumber:      1,
		},
	}
	
	results, errors = tc.store.BatchUpdateReservationBins(tc.ctx, updates)
	require.Equal(t, 0, len(errors), "Should have no errors")
	require.Equal(t, 2, len(results), "Should have 2 results")
	
	// Verify updated results include previous values
	assert.Equal(t, uint64(1500), results[0]) // 1000 + 500
	assert.Equal(t, uint64(3700), results[1]) // 3000 + 700
	
	// Verify records in database were updated
	item, err = dynamoClient.GetItem(tc.ctx, tc.reservationTable, commondynamodb.Key{
		"AccountID":         &types.AttributeValueMemberS{Value: accountID.Hex()},
		"ReservationPeriod": &types.AttributeValueMemberN{Value: "100"},
		"QuorumNumber":      &types.AttributeValueMemberN{Value: "0"},
	})
	require.NoError(t, err)
	binUsageStr = item["BinUsage"].(*types.AttributeValueMemberN).Value
	binUsageVal, err = strconv.ParseUint(binUsageStr, 10, 64)
	require.NoError(t, err)
	assert.Equal(t, uint64(1500), binUsageVal)
	
	item, err = dynamoClient.GetItem(tc.ctx, tc.reservationTable, commondynamodb.Key{
		"AccountID":         &types.AttributeValueMemberS{Value: accountID.Hex()},
		"ReservationPeriod": &types.AttributeValueMemberN{Value: "100"},
		"QuorumNumber":      &types.AttributeValueMemberN{Value: "1"},
	})
	require.NoError(t, err)
	binUsageStr = item["BinUsage"].(*types.AttributeValueMemberN).Value
	binUsageVal, err = strconv.ParseUint(binUsageStr, 10, 64)
	require.NoError(t, err)
	assert.Equal(t, uint64(3700), binUsageVal)
}

// TestGetLargestCumulativePayment tests the GetLargestCumulativePayment function
func TestGetLargestCumulativePayment(t *testing.T) {
	tc := setupTest(t)

	// Create an account to test with
	accountID := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")

	// Test case 1: No payment exists yet
	largest, err := tc.store.GetLargestCumulativePayment(tc.ctx, accountID)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(0), largest, "Initial largest payment should be 0")

	// Test case 2: Add first payment of 100 with charge of 100
	payment1 := core.PaymentMetadata{
		AccountID:         accountID,
		Timestamp:         time.Now().Unix(),
		CumulativePayment: big.NewInt(100),
	}
	oldPayment, err := tc.store.AddOnDemandPayment(tc.ctx, payment1, big.NewInt(100))
	require.NoError(t, err)
	require.Equal(t, big.NewInt(0), oldPayment, "Old payment should be 0 for first payment")

	largest, err = tc.store.GetLargestCumulativePayment(tc.ctx, accountID)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(100), largest, "Largest payment should be 100")

	// Test case 3: Add second payment of 300 with charge of 200 (cumulative)
	payment2 := core.PaymentMetadata{
		AccountID:         accountID,
		Timestamp:         time.Now().Unix(),
		CumulativePayment: big.NewInt(300),
	}
	oldPayment, err = tc.store.AddOnDemandPayment(tc.ctx, payment2, big.NewInt(200))
	require.NoError(t, err)
	require.Equal(t, big.NewInt(100), oldPayment, "Old payment should be 100")

	largest, err = tc.store.GetLargestCumulativePayment(tc.ctx, accountID)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(300), largest, "Largest payment should be 300")

	// Test case 4: Try to add payment of 200 with charge of 100 - should fail since cumulative is less than previous
	payment3 := core.PaymentMetadata{
		AccountID:         accountID,
		Timestamp:         time.Now().Unix(),
		CumulativePayment: big.NewInt(200),
	}
	oldPayment, err = tc.store.AddOnDemandPayment(tc.ctx, payment3, big.NewInt(100))
	require.Error(t, err)
	require.Nil(t, oldPayment, "Old payment should be nil on error")

	largest, err = tc.store.GetLargestCumulativePayment(tc.ctx, accountID)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(300), largest, "Largest payment should still be 300")

	// Test case 5: Add payment of 500 with insufficient charge (250) - should fail
	payment4 := core.PaymentMetadata{
		AccountID:         accountID,
		Timestamp:         time.Now().Unix(),
		CumulativePayment: big.NewInt(500),
	}
	oldPayment, err = tc.store.AddOnDemandPayment(tc.ctx, payment4, big.NewInt(250))
	require.Error(t, err)
	require.Nil(t, oldPayment, "Old payment should be nil on error")

	largest, err = tc.store.GetLargestCumulativePayment(tc.ctx, accountID)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(300), largest, "Largest payment should still be 300")

	// Test case 6: Add valid payment of 500 with sufficient charge (200)
	payment5 := core.PaymentMetadata{
		AccountID:         accountID,
		Timestamp:         time.Now().Unix(),
		CumulativePayment: big.NewInt(500),
	}
	oldPayment, err = tc.store.AddOnDemandPayment(tc.ctx, payment5, big.NewInt(200))
	require.NoError(t, err)
	require.Equal(t, big.NewInt(300), oldPayment, "Old payment should be 300")

	largest, err = tc.store.GetLargestCumulativePayment(tc.ctx, accountID)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(500), largest, "Largest payment should be 500")

	// Test case 7: Roll back the payment
	err = tc.store.RollbackOnDemandPayment(tc.ctx, accountID, big.NewInt(500), big.NewInt(300))
	require.NoError(t, err)

	largest, err = tc.store.GetLargestCumulativePayment(tc.ctx, accountID)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(300), largest, "After rollback, largest payment should be 300")

	// Test case 8: Verify rolling back a non-existent payment has no effect
	err = tc.store.RollbackOnDemandPayment(tc.ctx, accountID, big.NewInt(9999), big.NewInt(500))
	require.NoError(t, err)
}
