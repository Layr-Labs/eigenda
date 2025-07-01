package meterer_test

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"testing"
	"time"

	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/core/payment"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testContext struct {
	ctx              context.Context
	store            payment.PaymentOffchainState
	reservationTable string
	onDemandTable    string
	globalBinTable   string
}

// setupTest creates a test context with tables created
func setupTest(t *testing.T) *testContext {
	tc := &testContext{
		ctx:              context.Background(),
		reservationTable: fmt.Sprintf("reservation_test_%d", rand.Int()),
		onDemandTable:    fmt.Sprintf("ondemand_test_%d", rand.Int()),
		globalBinTable:   fmt.Sprintf("global_bin_test_%d", rand.Int()),
	}

	err := meterer.CreateReservationTable(clientConfig, tc.reservationTable)
	require.NoError(t, err)

	err = meterer.CreateOnDemandTable(clientConfig, tc.onDemandTable)
	require.NoError(t, err)

	err = meterer.CreateGlobalReservationTable(clientConfig, tc.globalBinTable)
	require.NoError(t, err)

	tc.store, err = meterer.NewDynamoDBPaymentOffchainState(
		clientConfig,
		tc.reservationTable,
		tc.onDemandTable,
		tc.globalBinTable,
		nil, // Logger not needed for test
	)
	require.NoError(t, err)

	return tc
}

// cleanupTest removes all tables created for a test
func cleanupTest(t *testing.T, tc *testContext) {
	_ = dynamoClient.DeleteTable(tc.ctx, tc.reservationTable)
	_ = dynamoClient.DeleteTable(tc.ctx, tc.onDemandTable)
	_ = dynamoClient.DeleteTable(tc.ctx, tc.globalBinTable)
}

func withTestContext(t *testing.T, testFn func(t *testing.T, tc *testContext)) {
	tc := setupTest(t)
	t.Cleanup(func() { cleanupTest(t, tc) })
	testFn(t, tc)
}

// TestIncrementBinUsages_EdgeCases tests the IncrementBinUsages function with edge cases
func TestIncrementBinUsages_EdgeCases(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		withTestContext(t, func(t *testing.T, tc *testContext) {
			binUsages, errs := tc.store.IncrementBinUsages(tc.ctx, gethcommon.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"), []core.QuorumID{}, map[core.QuorumID]uint64{}, map[core.QuorumID]uint64{})
			assert.Empty(t, binUsages)
			assert.Empty(t, errs)
		})
	})

	t.Run("exceed transaction limit", func(t *testing.T) {
		withTestContext(t, func(t *testing.T, tc *testContext) {
			accountID := gethcommon.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd")
			reservationPeriod := uint64(42)
			size := uint64(100)
			var quorums []core.QuorumID
			periods := make(map[core.QuorumID]uint64)
			for i := 0; i < commondynamodb.DynamoBatchWriteLimit+1; i++ {
				quorums = append(quorums, core.QuorumID(i))
				periods[core.QuorumID(i)] = reservationPeriod
			}
			binUsages, err := tc.store.IncrementBinUsages(tc.ctx, accountID, quorums, periods, map[core.QuorumID]uint64{core.QuorumID(0): size})
			assert.Empty(t, binUsages)
			assert.Error(t, err)
		})
	})

	t.Run("existing bin (increment)", func(t *testing.T) {
		withTestContext(t, func(t *testing.T, tc *testContext) {
			accountID := gethcommon.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd")
			reservationPeriod := uint64(42)
			size := uint64(100)
			quorum := core.QuorumID(1)
			periods := map[core.QuorumID]uint64{quorum: reservationPeriod}
			// First increment
			_, err := tc.store.IncrementBinUsages(tc.ctx, accountID, []core.QuorumID{quorum}, periods, map[core.QuorumID]uint64{quorum: size})
			assert.NoError(t, err)
			// Second increment
			binUsages, err := tc.store.IncrementBinUsages(tc.ctx, accountID, []core.QuorumID{quorum}, periods, map[core.QuorumID]uint64{quorum: size})
			assert.NoError(t, err)
			assert.Equal(t, size*2, binUsages[quorum])
		})
	})

	t.Run("nonexistent bin (first write)", func(t *testing.T) {
		withTestContext(t, func(t *testing.T, tc *testContext) {
			size := uint64(100)
			quorums := []core.QuorumID{10, 11}
			periods := map[core.QuorumID]uint64{10: 42, 11: 42}
			binUsages, err := tc.store.IncrementBinUsages(tc.ctx, gethcommon.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"), quorums, periods, map[core.QuorumID]uint64{10: size, 11: size})
			for _, quorum := range quorums {
				assert.NoError(t, err)
				assert.Equal(t, size, binUsages[quorum])
			}
		})
	})

	t.Run("exceed transaction limit", func(t *testing.T) {
		withTestContext(t, func(t *testing.T, tc *testContext) {
			accountID := gethcommon.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd")
			reservationPeriod := uint64(42)
			sizes := make(map[core.QuorumID]uint64)
			var quorums []core.QuorumID
			periods := make(map[core.QuorumID]uint64)
			for i := 0; i < commondynamodb.DynamoBatchWriteLimit+1; i++ { // 26 > DynamoDB batch limit
				quorums = append(quorums, core.QuorumID(i))
				periods[core.QuorumID(i)] = reservationPeriod
				sizes[core.QuorumID(i)] = uint64(i)
			}
			_, err := tc.store.IncrementBinUsages(tc.ctx, accountID, quorums, periods, sizes)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), fmt.Sprintf("limit is %d", commondynamodb.DynamoBatchWriteLimit))
		})
	})
}

// TestUpdateGlobalBin tests the UpdateGlobalBin function
func TestUpdateGlobalBin(t *testing.T) {
	withTestContext(t, func(t *testing.T, tc *testContext) {
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
	})
}

// TestAddOnDemandPayment tests the AddOnDemandPayment function
func TestAddOnDemandPayment(t *testing.T) {
	withTestContext(t, func(t *testing.T, tc *testContext) {
		accountID := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")
		payment1 := payment.PaymentMetadata{
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
		payment2 := payment.PaymentMetadata{
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
		payment3 := payment.PaymentMetadata{
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
		payment4 := payment.PaymentMetadata{
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
	})
}

// TestRollbackOnDemandPayment tests the RollbackOnDemandPayment function
func TestRollbackOnDemandPayment(t *testing.T) {
	withTestContext(t, func(t *testing.T, tc *testContext) {
		// Create and add a payment
		accountID := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")
		cumulativePayment := big.NewInt(1000)
		paymentCharged := big.NewInt(500)

		paymentMetadata := payment.PaymentMetadata{
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
		newPaymentMetadata := payment.PaymentMetadata{
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
	})
}

// TestGetLargestCumulativePayment tests the GetLargestCumulativePayment function
func TestGetLargestCumulativePayment(t *testing.T) {
	withTestContext(t, func(t *testing.T, tc *testContext) {
		// Create an account to test with
		accountID := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")

		// Test case 1: No payment exists yet
		largest, err := tc.store.GetLargestCumulativePayment(tc.ctx, accountID)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(0), largest, "Initial largest payment should be 0")

		// Test case 2: Add first payment of 100 with charge of 100
		payment1 := payment.PaymentMetadata{
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
		payment2 := payment.PaymentMetadata{
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
		payment3 := payment.PaymentMetadata{
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
		payment4 := payment.PaymentMetadata{
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
		payment5 := payment.PaymentMetadata{
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
	})
}

// TestGetPeriodRecords tests the GetPeriodRecords function
func TestGetPeriodRecords(t *testing.T) {
	accountID := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")
	quorum1 := core.QuorumID(1)
	quorum2 := core.QuorumID(2)
	period1 := uint64(42)
	period2 := uint64(43)
	size1 := uint64(100)
	size2 := uint64(200)

	t.Run("empty input", func(t *testing.T) {
		withTestContext(t, func(t *testing.T, tc *testContext) {
			records, err := tc.store.GetPeriodRecords(tc.ctx, accountID, []core.QuorumID{}, []uint64{}, 5)
			require.NoError(t, err)
			assert.Empty(t, records)
		})
	})

	t.Run("single quorum single period", func(t *testing.T) {
		withTestContext(t, func(t *testing.T, tc *testContext) {
			// Setup
			_, err := tc.store.IncrementBinUsages(tc.ctx, accountID, []core.QuorumID{quorum1}, map[core.QuorumID]uint64{quorum1: period1}, map[core.QuorumID]uint64{quorum1: size1})
			require.NoError(t, err)

			// Get the records
			records, err := tc.store.GetPeriodRecords(tc.ctx, accountID, []core.QuorumID{quorum1}, []uint64{period1}, 1)
			require.NoError(t, err)
			require.Len(t, records, 1)
			require.Contains(t, records, quorum1)
			require.Len(t, records[quorum1].Records, 1)
			assert.Equal(t, period1, uint64(records[quorum1].Records[0].Index))
			assert.Equal(t, size1, records[quorum1].Records[0].Usage)
		})
	})

	t.Run("multiple quorums multiple periods", func(t *testing.T) {
		withTestContext(t, func(t *testing.T, tc *testContext) {
			// Setup
			_, err := tc.store.IncrementBinUsages(tc.ctx, accountID, []core.QuorumID{quorum1, quorum2},
				map[core.QuorumID]uint64{quorum1: period1, quorum2: period2},
				map[core.QuorumID]uint64{quorum1: size1, quorum2: size2})
			require.NoError(t, err)

			// Get the records
			records, err := tc.store.GetPeriodRecords(tc.ctx, accountID, []core.QuorumID{quorum1, quorum2}, []uint64{period1, period2}, 1)
			require.NoError(t, err)
			require.Len(t, records, 2)

			// Check quorum1 records
			require.Contains(t, records, quorum1)
			require.Len(t, records[quorum1].Records, 1)
			assert.Equal(t, period1, uint64(records[quorum1].Records[0].Index))
			assert.Equal(t, size1, records[quorum1].Records[0].Usage)

			// Check quorum2 records
			require.Contains(t, records, quorum2)
			require.Len(t, records[quorum2].Records, 1)
			assert.Equal(t, period2, uint64(records[quorum2].Records[0].Index))
			assert.Equal(t, size2, records[quorum2].Records[0].Usage)
		})
	})

	t.Run("multiple periods for same quorum", func(t *testing.T) {
		withTestContext(t, func(t *testing.T, tc *testContext) {
			// Setup
			_, err := tc.store.IncrementBinUsages(tc.ctx, accountID, []core.QuorumID{quorum1},
				map[core.QuorumID]uint64{quorum1: period1},
				map[core.QuorumID]uint64{quorum1: size1})
			require.NoError(t, err)

			_, err = tc.store.IncrementBinUsages(tc.ctx, accountID, []core.QuorumID{quorum1},
				map[core.QuorumID]uint64{quorum1: period2},
				map[core.QuorumID]uint64{quorum1: size2})
			require.NoError(t, err)

			// Get the records
			records, err := tc.store.GetPeriodRecords(tc.ctx, accountID, []core.QuorumID{quorum1}, []uint64{period1}, 2)
			require.NoError(t, err)
			require.Len(t, records, 1)
			require.Contains(t, records, quorum1)
			require.Len(t, records[quorum1].Records, 2)

			// Find and verify period1 record
			var period1Record *pb.PeriodRecord
			for _, record := range records[quorum1].Records {
				if record.Index == uint32(period1) {
					period1Record = record
					break
				}
			}
			require.NotNil(t, period1Record, "Period1 record not found")
			assert.Equal(t, size1, period1Record.Usage)

			// Find and verify period2 record
			var period2Record *pb.PeriodRecord
			for _, record := range records[quorum1].Records {
				if record.Index == uint32(period2) {
					period2Record = record
					break
				}
			}
			require.NotNil(t, period2Record, "Period2 record not found")
			assert.Equal(t, size2, period2Record.Usage)

			// Verify no other periods exist
			for _, record := range records[quorum1].Records {
				assert.True(t, record.Index == uint32(period1) || record.Index == uint32(period2),
					"Unexpected period index found: %d", record.Index)
			}
		})
	})

	t.Run("numBins larger than existing periods", func(t *testing.T) {
		withTestContext(t, func(t *testing.T, tc *testContext) {
			// Setup
			_, err := tc.store.IncrementBinUsages(tc.ctx, accountID, []core.QuorumID{quorum1},
				map[core.QuorumID]uint64{quorum1: 1},
				map[core.QuorumID]uint64{quorum1: size1})
			require.NoError(t, err)

			_, err = tc.store.IncrementBinUsages(tc.ctx, accountID, []core.QuorumID{quorum1},
				map[core.QuorumID]uint64{quorum1: 2},
				map[core.QuorumID]uint64{quorum1: size2})
			require.NoError(t, err)

			// Request 5 bins but only 2 exist
			records, err := tc.store.GetPeriodRecords(tc.ctx, accountID, []core.QuorumID{quorum1}, []uint64{1}, 5)
			require.NoError(t, err)
			require.Len(t, records, 1)
			require.Contains(t, records, quorum1)
			require.Len(t, records[quorum1].Records, 2, "Should only return existing periods")

			// Verify period 1
			var period1Record *pb.PeriodRecord
			for _, record := range records[quorum1].Records {
				if record.Index == 1 {
					period1Record = record
					break
				}
			}
			require.NotNil(t, period1Record, "Period 1 record not found")
			assert.Equal(t, size1, period1Record.Usage)

			// Verify period 2
			var period2Record *pb.PeriodRecord
			for _, record := range records[quorum1].Records {
				if record.Index == 2 {
					period2Record = record
					break
				}
			}
			require.NotNil(t, period2Record, "Period 2 record not found")
			assert.Equal(t, size2, period2Record.Usage)

			// Verify no other periods exist
			for _, record := range records[quorum1].Records {
				assert.True(t, record.Index == 1 || record.Index == 2,
					"Unexpected period index found: %d", record.Index)
			}
		})
	})

	t.Run("nonexistent periods", func(t *testing.T) {
		withTestContext(t, func(t *testing.T, tc *testContext) {
			// Try to get records for periods that don't exist
			records, err := tc.store.GetPeriodRecords(tc.ctx, accountID, []core.QuorumID{quorum1}, []uint64{999}, 1)
			require.NoError(t, err)
			assert.Empty(t, records)
		})
	})

	t.Run("mismatched quorum and period lengths", func(t *testing.T) {
		withTestContext(t, func(t *testing.T, tc *testContext) {
			// Test with mismatched lengths
			records, err := tc.store.GetPeriodRecords(tc.ctx, accountID, []core.QuorumID{quorum1}, []uint64{period1, period2}, 1)
			require.NoError(t, err)
			assert.Empty(t, records)
		})
	})
}
