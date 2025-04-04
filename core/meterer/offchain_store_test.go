package meterer_test

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testContext struct {
	ctx     context.Context
	dbPath  string
	store   meterer.OffchainStore
	cleanup func()
}

func setupTestContext(t *testing.T) *testContext {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, fmt.Sprintf("test_db_%d", rand.Int()))

	tc := &testContext{
		ctx:    context.Background(),
		dbPath: dbPath,
	}

	// Create the OffchainStore
	logger := testutils.GetLogger()
	var err error
	tc.store, err = meterer.NewOffchainStore(
		dbPath,
		logger,
	)
	require.NoError(t, err)

	// Set up cleanup function after store is created
	tc.cleanup = func() {
		t.Logf("Cleaning up store")
		// First destroy the store to ensure all files are closed
		if err := tc.store.Destroy(); err != nil {
			t.Logf("Failed to destroy store: %v", err)
		}

		// Give a small delay to ensure all file handles are released
		time.Sleep(100 * time.Millisecond)

		// Remove all files in the test directory
		t.Logf("Removing test directory: %s", tmpDir)
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("Failed to remove test directory: %v", err)
		}

		// Double check that the directory is gone
		if _, err := os.Stat(tmpDir); !os.IsNotExist(err) {
			t.Logf("Test directory still exists after cleanup: %v", err)
			// Try to list remaining files for debugging
			if files, err := os.ReadDir(tmpDir); err == nil {
				for _, file := range files {
					t.Logf("Remaining file: %s", file.Name())
				}
			}
		}
	}

	// Register cleanup with testing.T
	t.Cleanup(tc.cleanup)

	return tc
}

// TestUpdateReservationBin tests the UpdateReservationBin function
func TestUpdateReservationBin(t *testing.T) {
	tc := setupTestContext(t)

	// Test updating bin that doesn't exist yet (should create it)
	accountID := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")
	reservationPeriod := uint64(1)
	size := uint64(1000)

	// First verify it doesn't exist
	binUsage, err := tc.store.GetReservationBin(tc.ctx, accountID, reservationPeriod)
	require.NoError(t, err)
	assert.Equal(t, uint64(0), binUsage)

	// Update the bin
	binUsage, err = tc.store.UpdateReservationBin(tc.ctx, accountID, reservationPeriod, size)
	require.NoError(t, err)
	assert.Equal(t, size, binUsage)

	// Verify the update
	binUsage, err = tc.store.GetReservationBin(tc.ctx, accountID, reservationPeriod)
	require.NoError(t, err)
	assert.Equal(t, size, binUsage)

	// Test updating existing bin
	additionalSize := uint64(500)
	binUsage, err = tc.store.UpdateReservationBin(tc.ctx, accountID, reservationPeriod, additionalSize)
	require.NoError(t, err)
	assert.Equal(t, size+additionalSize, binUsage)

	// Verify the update
	binUsage, err = tc.store.GetReservationBin(tc.ctx, accountID, reservationPeriod)
	require.NoError(t, err)
	assert.Equal(t, size+additionalSize, binUsage)
}

// TestUpdateGlobalBin tests the UpdateGlobalBin function
func TestUpdateGlobalBin(t *testing.T) {
	tc := setupTestContext(t)

	// Test updating global bin that doesn't exist yet (should create it)
	reservationPeriod := uint64(1)
	size := uint64(2000)

	// First verify it doesn't exist
	binUsage, err := tc.store.GetGlobalBin(tc.ctx, reservationPeriod)
	require.NoError(t, err)
	assert.Equal(t, uint64(0), binUsage)

	// Update the bin
	binUsage, err = tc.store.UpdateGlobalBin(tc.ctx, reservationPeriod, size)
	require.NoError(t, err)
	assert.Equal(t, size, binUsage)

	// Verify the update
	binUsage, err = tc.store.GetGlobalBin(tc.ctx, reservationPeriod)
	require.NoError(t, err)
	assert.Equal(t, size, binUsage)

	// Test updating existing bin
	additionalSize := uint64(1000)
	binUsage, err = tc.store.UpdateGlobalBin(tc.ctx, reservationPeriod, additionalSize)
	require.NoError(t, err)
	assert.Equal(t, size+additionalSize, binUsage)

	// Verify the update
	binUsage, err = tc.store.GetGlobalBin(tc.ctx, reservationPeriod)
	require.NoError(t, err)
	assert.Equal(t, size+additionalSize, binUsage)
}

// TestAddOnDemandPayment tests the AddOnDemandPayment function
func TestAddOnDemandPayment(t *testing.T) {
	tc := setupTestContext(t)

	accountID := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")
	payment1 := core.PaymentMetadata{
		AccountID:         accountID,
		Timestamp:         time.Now().Unix(),
		CumulativePayment: big.NewInt(100),
	}
	charge1 := big.NewInt(100)

	// First verify it doesn't exist
	payment, err := tc.store.GetOnDemandPayment(tc.ctx, accountID)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(0), payment)

	// Add the payment
	oldPayment, err := tc.store.AddOnDemandPayment(tc.ctx, payment1, charge1)
	require.NoError(t, err)
	require.Condition(t, func() bool {
		return oldPayment.Cmp(big.NewInt(0)) == 0
	}, "Old payment should be 0 for first payment")

	// Verify the update
	payment, err = tc.store.GetOnDemandPayment(tc.ctx, accountID)
	require.NoError(t, err)
	assert.Equal(t, payment1.CumulativePayment, payment)

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

	// Verify the update
	payment, err = tc.store.GetOnDemandPayment(tc.ctx, accountID)
	require.NoError(t, err)
	assert.Equal(t, payment2.CumulativePayment, payment)

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
	payment, err = tc.store.GetOnDemandPayment(tc.ctx, accountID)
	require.NoError(t, err)
	assert.Equal(t, payment2.CumulativePayment, payment)

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
	payment, err = tc.store.GetOnDemandPayment(tc.ctx, accountID)
	require.NoError(t, err)
	assert.Equal(t, payment2.CumulativePayment, payment)
}

// TestRollbackOnDemandPayment tests the RollbackOnDemandPayment function
func TestRollbackOnDemandPayment(t *testing.T) {
	tc := setupTestContext(t)

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
	require.Condition(t, func() bool {
		return oldPayment.Cmp(big.NewInt(0)) == 0
	}, "Old payment should be 0 for first payment")

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

	// Test case 2: Rollback to a different value directly
	// The value will be updated regardless of what the current value is
	err = tc.store.RollbackOnDemandPayment(tc.ctx, accountID, big.NewInt(1000), big.NewInt(500))
	require.NoError(t, err)

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

// TestGetLargestCumulativePayment tests the GetLargestCumulativePayment function
func TestGetLargestCumulativePayment(t *testing.T) {
	tc := setupTestContext(t)

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
	require.Condition(t, func() bool {
		return oldPayment.Cmp(big.NewInt(0)) == 0
	}, "Old payment should be 0 for first payment")

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

// TestGetPeriodRecords tests the GetPeriodRecords function
func TestGetPeriodRecords(t *testing.T) {
	tc := setupTestContext(t)

	accountID := gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")
	reservationPeriod := uint64(1)

	// Add some test data
	_, err := tc.store.UpdateReservationBin(tc.ctx, accountID, reservationPeriod, 1000)
	require.NoError(t, err)
	_, err = tc.store.UpdateReservationBin(tc.ctx, accountID, reservationPeriod+1, 2000)
	require.NoError(t, err)
	_, err = tc.store.UpdateReservationBin(tc.ctx, accountID, reservationPeriod+2, 3000)
	require.NoError(t, err)

	// Get records
	records, err := tc.store.GetPeriodRecords(tc.ctx, accountID, reservationPeriod)
	require.NoError(t, err)

	// Verify records
	assert.Equal(t, uint32(reservationPeriod), records[0].Index)
	assert.Equal(t, uint64(1000), records[0].Usage)
	assert.Equal(t, uint32(reservationPeriod+1), records[1].Index)
	assert.Equal(t, uint64(2000), records[1].Usage)
	assert.Equal(t, uint32(reservationPeriod+2), records[2].Index)
	assert.Equal(t, uint64(3000), records[2].Usage)
}
