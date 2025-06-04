package integration_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/node"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestBatchMetererValidReservation tests the BatchMeterer with valid reservation parameters
func TestBatchMetererValidReservation(t *testing.T) {
	ctx := context.Background()
	logger := testutils.GetLogger()
	mockState := setupMockOnchainPaymentState()

	account1 := gethcommon.HexToAddress("0x1111111111111111111111111111111111111111")
	account2 := gethcommon.HexToAddress("0x2222222222222222222222222222222222222222")

	// Set up the reservations for account1 - both quorums 0 and 1
	reservationsAcc1 := map[core.QuorumID]*core.ReservedPayment{
		0: {
			SymbolsPerSecond: 1000,
			StartTimestamp:   0,
			EndTimestamp:     uint64(time.Now().Add(24 * time.Hour).Unix()),
		},
		1: {
			SymbolsPerSecond: 500,
			StartTimestamp:   0,
			EndTimestamp:     uint64(time.Now().Add(24 * time.Hour).Unix()),
		},
	}

	// Set up the reservations for account2 - only quorum 1
	reservationsAcc2 := map[core.QuorumID]*core.ReservedPayment{
		1: {
			SymbolsPerSecond: 800,
			StartTimestamp:   0,
			EndTimestamp:     uint64(time.Now().Add(24 * time.Hour).Unix()),
		},
	}

	mockState.On("GetReservedPaymentByAccountAndQuorums", mock.Anything, account1, mock.MatchedBy(func(quorums []core.QuorumID) bool {
		return len(quorums) == 2 && containsAll(quorums, []core.QuorumID{0, 1})
	})).Return(reservationsAcc1, nil)
	mockState.On("GetReservedPaymentByAccountAndQuorums", mock.Anything, account2, []core.QuorumID{1}).Return(reservationsAcc2, nil)
	mockState.On("GetMinNumSymbols").Return(uint64(32))
	mockState.On("GetReservationWindow").Return(uint64(3600))

	batchMeterer := node.NewBatchMeterer(
		meterer.Config{
			ChainReadTimeout: 1 * time.Second,
			UpdateInterval:   1 * time.Minute,
		},
		mockState,
		5, // numBins
		logger,
	)

	batch := createTestBatch(
		[]accountQuorumInfo{
			{
				account:    account1,
				quorumIDs:  []core.QuorumID{0, 1},
				numSymbols: 500, // Well under the reservation limit
			},
			{
				account:    account2,
				quorumIDs:  []core.QuorumID{1},
				numSymbols: 400, // Well under the reservation limit
			},
		},
	)

	err := batchMeterer.MeterBatch(ctx, batch, time.Now())
	assert.NoError(t, err, "Metering should succeed for batch with valid usage")
	mockState.AssertExpectations(t)
}

// TestBatchMetererExceedingReservation tests the BatchMeterer with usage exceeding reservation limits
func TestBatchMetererExceedingReservation(t *testing.T) {
	ctx := context.Background()
	logger := testutils.GetLogger()
	mockState := setupMockOnchainPaymentState()

	account := gethcommon.HexToAddress("0x1111111111111111111111111111111111111111")
	reservations := map[core.QuorumID]*core.ReservedPayment{
		0: {
			SymbolsPerSecond: 100, // Very low limit
			StartTimestamp:   0,
			EndTimestamp:     uint64(time.Now().Add(24 * time.Hour).Unix()),
		},
	}
	mockState.On("GetReservedPaymentByAccountAndQuorums", mock.Anything, account, []core.QuorumID{0}).Return(reservations, nil)
	mockState.On("GetMinNumSymbols").Return(uint64(32))
	mockState.On("GetReservationWindow").Return(uint64(3600))

	batchMeterer := node.NewBatchMeterer(
		meterer.Config{
			ChainReadTimeout: 1 * time.Second,
			UpdateInterval:   1 * time.Minute,
		},
		mockState,
		5, // numBins
		logger,
	)

	batch := createTestBatch(
		[]accountQuorumInfo{
			{
				account:    account,
				quorumIDs:  []core.QuorumID{0},
				numSymbols: 1000000, // Greatly exceeds the reservation limit (100 * 3600 = 360000)
			},
		},
	)

	err := batchMeterer.MeterBatch(ctx, batch, time.Now())
	assert.Error(t, err, "Metering should fail for batch exceeding reservation limit")
	assert.Contains(t, err.Error(), "exceeds bin limit", "Error should mention exceeding bin limit")
	mockState.AssertExpectations(t)
}

// TestBatchMetererInactiveReservation tests the BatchMeterer with an inactive reservation
func TestBatchMetererInactiveReservation(t *testing.T) {
	ctx := context.Background()
	logger := testutils.GetLogger()
	mockState := setupMockOnchainPaymentState()

	account := gethcommon.HexToAddress("0x1111111111111111111111111111111111111111")
	reservations := map[core.QuorumID]*core.ReservedPayment{
		0: {
			SymbolsPerSecond: 1000,
			StartTimestamp:   0,
			EndTimestamp:     uint64(time.Now().Add(-24 * time.Hour).Unix()), // Expired
		},
	}
	mockState.On("GetReservedPaymentByAccountAndQuorums", mock.Anything, account, []core.QuorumID{0}).Return(reservations, nil)
	mockState.On("GetMinNumSymbols").Return(uint64(32))
	mockState.On("GetReservationWindow").Return(uint64(3600))

	batchMeterer := node.NewBatchMeterer(
		meterer.Config{
			ChainReadTimeout: 1 * time.Second,
			UpdateInterval:   1 * time.Minute,
		},
		mockState,
		5, // numBins
		logger,
	)
	batch := createTestBatch(
		[]accountQuorumInfo{
			{
				account:    account,
				quorumIDs:  []core.QuorumID{0},
				numSymbols: 100, // Even a small amount should fail with inactive reservation
			},
		},
	)
	err := batchMeterer.MeterBatch(ctx, batch, time.Now())
	assert.Error(t, err, "Metering should fail for batch with inactive reservation")
	assert.Contains(t, err.Error(), "inactive reservation", "Error should mention inactive reservation")

	mockState.AssertExpectations(t)
}

// TestBatchMetererInvalidReservationPeriod tests the BatchMeterer with an invalid reservation period
func TestBatchMetererInvalidReservationPeriod(t *testing.T) {
	ctx := context.Background()
	logger := testutils.GetLogger()
	mockState := setupMockOnchainPaymentState()

	account := gethcommon.HexToAddress("0x1111111111111111111111111111111111111111")
	reservationWindow := uint64(3600) // 1 hour in seconds
	now := time.Now()

	// Create a reservation with valid start/end timestamps
	reservations := map[core.QuorumID]*core.ReservedPayment{
		0: {
			SymbolsPerSecond: 1000,
			StartTimestamp:   uint64(now.Add(-24 * time.Hour).Unix()), // Started 24 hours ago
			EndTimestamp:     uint64(now.Add(24 * time.Hour).Unix()),  // Ends 24 hours in the future
		},
	}

	// Clear existing expectations and call history
	mockState.ExpectedCalls = nil
	mockState.Calls = nil

	// Set up expectations for all three test cases (valid timestamp, before start, after end)
	mockState.On("GetReservedPaymentByAccountAndQuorums", mock.Anything, account, []core.QuorumID{0}).Return(reservations, nil).Times(3)
	mockState.On("GetMinNumSymbols").Return(uint64(32)).Maybe()
	mockState.On("GetReservationWindow").Return(reservationWindow).Maybe()

	batchMeterer := node.NewBatchMeterer(
		meterer.Config{
			ChainReadTimeout: 1 * time.Second,
			UpdateInterval:   1 * time.Minute,
		},
		mockState,
		5, // numBins
		logger,
	)

	// Create a test batch
	batch := createTestBatch(
		[]accountQuorumInfo{
			{
				account:    account,
				quorumIDs:  []core.QuorumID{0},
				numSymbols: 100,
			},
		},
	)

	// Test with a valid timestamp in the middle of the reservation period
	validTimestamp := time.Unix(int64(reservations[0].StartTimestamp)+3600, 0) // 1 hour after reservation start
	err := batchMeterer.MeterBatch(ctx, batch, validTimestamp)
	assert.NoError(t, err, "Metering should succeed with valid timestamp")

	// Test with a timestamp that's way too old (outside the start of the reservation)
	oldTimestamp := time.Unix(int64(reservations[0].StartTimestamp)-7200, 0) // 2 hours before reservation start
	err = batchMeterer.MeterBatch(ctx, batch, oldTimestamp)
	assert.Error(t, err, "Metering should fail with timestamp before reservation start")
	assert.Contains(t, err.Error(), "inactive reservation", "Error should mention inactive reservation")

	// Test with a timestamp that's after the end of the reservation
	futureTimestamp := time.Unix(int64(reservations[0].EndTimestamp)+7200, 0) // 2 hours after reservation end
	err = batchMeterer.MeterBatch(ctx, batch, futureTimestamp)
	assert.Error(t, err, "Metering should fail with timestamp after reservation end")
	assert.Contains(t, err.Error(), "inactive reservation", "Error should mention inactive reservation")

	mockState.AssertExpectations(t)
}

// TestBatchMetererMissingReservation tests the BatchMeterer with a quorum that has no reservation
func TestBatchMetererMissingReservation(t *testing.T) {
	ctx := context.Background()
	logger := testutils.GetLogger()
	mockState := setupMockOnchainPaymentState()

	account := gethcommon.HexToAddress("0x1111111111111111111111111111111111111111")
	reservations := map[core.QuorumID]*core.ReservedPayment{
		0: {
			SymbolsPerSecond: 1000,
			StartTimestamp:   0,
			EndTimestamp:     uint64(time.Now().Add(24 * time.Hour).Unix()),
		},
	}

	mockState.On("GetReservedPaymentByAccountAndQuorums", mock.Anything, account, mock.MatchedBy(func(quorums []core.QuorumID) bool {
		return len(quorums) == 2 && containsAll(quorums, []core.QuorumID{0, 1})
	})).Return(reservations, nil)
	mockState.On("GetMinNumSymbols").Return(uint64(32))
	mockState.On("GetReservationWindow").Return(uint64(3600))

	batchMeterer := node.NewBatchMeterer(
		meterer.Config{
			ChainReadTimeout: 1 * time.Second,
			UpdateInterval:   1 * time.Minute,
		},
		mockState,
		5, // numBins
		logger,
	)
	batch := createTestBatch(
		[]accountQuorumInfo{
			{
				account:    account,
				quorumIDs:  []core.QuorumID{0, 1}, // Quorum 1 has no reservation
				numSymbols: 100,
			},
		},
	)

	err := batchMeterer.MeterBatch(ctx, batch, time.Now())
	assert.Error(t, err, "Metering should fail for batch with missing reservation")
	assert.Contains(t, err.Error(), "no reservation for quorum", "Error should mention missing reservation")

	mockState.AssertExpectations(t)
}

// TestBatchMetererMinSymbolsCharge tests that the meterer correctly applies the minimum symbols charge
func TestBatchMetererMinSymbolsCharge(t *testing.T) {
	ctx := context.Background()
	logger := testutils.GetLogger()
	mockState := setupMockOnchainPaymentState()

	account := gethcommon.HexToAddress("0x1111111111111111111111111111111111111111")

	// Set up a reservation with a moderate limit
	reservations := map[core.QuorumID]*core.ReservedPayment{
		0: {
			SymbolsPerSecond: 1000,
			StartTimestamp:   0,
			EndTimestamp:     uint64(time.Now().Add(24 * time.Hour).Unix()),
		},
	}

	minSymbols := uint64(32)

	mockState.On("GetReservedPaymentByAccountAndQuorums", mock.Anything, account, []core.QuorumID{0}).Return(reservations, nil)
	mockState.On("GetMinNumSymbols").Return(minSymbols)
	mockState.On("GetReservationWindow").Return(uint64(3600))

	batchMeterer := node.NewBatchMeterer(
		meterer.Config{
			ChainReadTimeout: 1 * time.Second,
			UpdateInterval:   1 * time.Minute,
		},
		mockState,
		5,
		logger,
	)

	smallBatch := createTestBatch(
		[]accountQuorumInfo{
			{
				account:    account,
				quorumIDs:  []core.QuorumID{0},
				numSymbols: 16, // Less than minimum (32)
			},
		},
	)

	// Test metering the batch - should succeed but charge minimum
	err := batchMeterer.MeterBatch(ctx, smallBatch, time.Now())
	assert.NoError(t, err, "Batch with small symbols should succeed but charge minimum")
	mockState.AssertExpectations(t)
}

// TestBatchMetererOverflowBehavior tests that overflow handling works correctly
// This test verifies that the BatchMeterer correctly:
// 1. Allows usage up to the bin limit for a reservation period
// 2. Allows small overflows to be moved to future bins
// 3. Prevents usage when a bin is already filled
// 4. Allows usage in different reservation periods
func TestBatchMetererOverflowBehavior(t *testing.T) {
	ctx := context.Background()
	logger := testutils.GetLogger()
	mockState := setupMockOnchainPaymentState()

	account := gethcommon.HexToAddress("0x1111111111111111111111111111111111111111")

	// Set up a reservation with a specific symbols per second limit
	// This creates a bin limit of 100 * 3600 = 360,000 symbols
	symbolsPerSecond := uint64(100)
	reservationWindow := uint64(3600)
	binLimit := symbolsPerSecond * reservationWindow

	// Set a long reservation period to allow for overflow testing
	reservations := map[core.QuorumID]*core.ReservedPayment{
		0: {
			SymbolsPerSecond: symbolsPerSecond,
			StartTimestamp:   0,
			EndTimestamp:     uint64(time.Now().Add(30 * 24 * time.Hour).Unix()), // 30 days in the future
		},
	}

	// Set up the required mock behaviors
	mockState.On("GetReservedPaymentByAccountAndQuorums",
		mock.Anything,
		account,
		mock.MatchedBy(func(quorums []core.QuorumID) bool {
			return len(quorums) == 1 && quorums[0] == 0
		}),
	).Return(reservations, nil).Times(4)

	mockState.On("GetMinNumSymbols").Return(uint64(32)).Maybe()
	mockState.On("GetReservationWindow").Return(reservationWindow).Maybe()

	// Create the batch meterer with enough bins to handle overflow
	batchMeterer := node.NewBatchMeterer(
		meterer.Config{
			ChainReadTimeout: 1 * time.Second,
			UpdateInterval:   1 * time.Minute,
		},
		mockState,
		10, // More bins to handle overflow
		logger,
	)

	// Fixed timestamp for deterministic testing
	now := time.Now()

	// Test 1: Create a batch that partially fills the current bin (should succeed)
	partialBatch := createTestBatch(
		[]accountQuorumInfo{
			{
				account:    account,
				quorumIDs:  []core.QuorumID{0},
				numSymbols: binLimit / 2, // Half the bin limit
			},
		},
	)
	err := batchMeterer.MeterBatch(ctx, partialBatch, now)
	assert.NoError(t, err, "Batch with partial bin usage should succeed")

	// Test 2: Create a batch that overflows slightly (should succeed with overflow)
	// The total usage will be (binLimit/2 + binLimit*0.75) which exceeds binLimit
	// The overflow amount will be (binLimit*1.25 - binLimit) = binLimit*0.25
	overflowBatch := createTestBatch(
		[]accountQuorumInfo{
			{
				account:    account,
				quorumIDs:  []core.QuorumID{0},
				numSymbols: binLimit * 3 / 4, // 75% of bin limit, causing overflow
			},
		},
	)
	err = batchMeterer.MeterBatch(ctx, overflowBatch, now)
	assert.NoError(t, err, "Batch with small overflow should succeed by moving excess to future bin")

	// Test 3: Create a batch that attempts to use an already filled bin (should fail)
	// The current bin is now filled from the previous operations
	additionalBatch := createTestBatch(
		[]accountQuorumInfo{
			{
				account:    account,
				quorumIDs:  []core.QuorumID{0},
				numSymbols: binLimit / 4, // Even a small amount should fail
			},
		},
	)
	err = batchMeterer.MeterBatch(ctx, additionalBatch, now)
	assert.Error(t, err, "Batch using already filled bin should fail")
	assert.Contains(t, err.Error(), "bin has already been filled", "Error should mention bin already filled")

	// Test 4: Create a batch that would exceed 2*binLimit (should fail)
	// This should fail because total usage would be > 2*binLimit
	largeOverflowBatch := createTestBatch(
		[]accountQuorumInfo{
			{
				account:    account,
				quorumIDs:  []core.QuorumID{0},
				numSymbols: binLimit * 2, // Would cause total usage to exceed 2*binLimit
			},
		},
	)
	err = batchMeterer.MeterBatch(ctx, largeOverflowBatch, now)
	assert.Error(t, err, "Batch causing large overflow should fail")
	assert.Contains(t, err.Error(), "overflow usage exceeds bin limit", "Error should mention overflow limit exceeded")

	// Test 5: Using a different timestamp (next period) should allow usage again
	// Calculate a timestamp in the next reservation period
	nextPeriodTime := now.Add(time.Duration(reservationWindow) * time.Second)

	// This batch should succeed because it's in a different period
	nextPeriodBatch := createTestBatch(
		[]accountQuorumInfo{
			{
				account:    account,
				quorumIDs:  []core.QuorumID{0},
				numSymbols: binLimit / 2, // Half the bin limit
			},
		},
	)
	err = batchMeterer.MeterBatch(ctx, nextPeriodBatch, nextPeriodTime)
	assert.NoError(t, err, "Batch in a different reservation period should succeed")

	mockState.AssertExpectations(t)
}

// TestBatchMetererOverflowBinExceedCapacity tests that overflow fails when it would exceed capacity of the target overflow bin
func TestBatchMetererOverflowBinExceedCapacity(t *testing.T) {
	ctx := context.Background()
	logger := testutils.GetLogger()
	mockState := setupMockOnchainPaymentState()

	account := gethcommon.HexToAddress("0x1111111111111111111111111111111111111111")

	// Set up a reservation with a specific symbols per second limit
	symbolsPerSecond := uint64(100)
	reservationWindow := uint64(3600)
	binLimit := symbolsPerSecond * reservationWindow

	// Set a long reservation period to allow for overflow testing
	reservations := map[core.QuorumID]*core.ReservedPayment{
		0: {
			SymbolsPerSecond: symbolsPerSecond,
			StartTimestamp:   0,
			EndTimestamp:     uint64(time.Now().Add(30 * 24 * time.Hour).Unix()), // 30 days in the future
		},
	}

	// Clear all existing expectations and call history
	mockState.ExpectedCalls = nil
	mockState.Calls = nil

	// Configure mock with precise expectations
	// This test has 4 cases that each call GetReservedPaymentByAccountAndQuorums once
	mockState.On("GetReservedPaymentByAccountAndQuorums",
		mock.Anything,
		account,
		mock.MatchedBy(func(quorums []core.QuorumID) bool {
			return len(quorums) == 1 && quorums[0] == 0
		}),
	).Return(reservations, nil).Times(4)

	mockState.On("GetMinNumSymbols").Return(uint64(32)).Maybe()
	mockState.On("GetReservationWindow").Return(reservationWindow).Maybe()

	// Create the batch meterer with only 2 bins to simulate overflow bin capacity issues
	batchMeterer := node.NewBatchMeterer(
		meterer.Config{
			ChainReadTimeout: 1 * time.Second,
			UpdateInterval:   1 * time.Minute,
		},
		mockState,
		2, // Only 2 bins to ensure we hit capacity issues
		logger,
	)

	// Use two different timestamps: current period and overflow period
	now := time.Now()

	// Create batch for overflow period directly to fill it first
	// Calculate a timestamp for the overflow period
	overflowPeriodWindow := 2 * reservationWindow // Same as used in the BatchMeterer implementation
	overflowPeriodTime := now.Add(time.Duration(overflowPeriodWindow) * time.Second)

	// Step 1: Fill the current bin partially (70% capacity)
	currentBatch := createTestBatch(
		[]accountQuorumInfo{
			{
				account:    account,
				quorumIDs:  []core.QuorumID{0},
				numSymbols: binLimit * 7 / 10, // 70% of bin limit
			},
		},
	)
	err := batchMeterer.MeterBatch(ctx, currentBatch, now)
	assert.NoError(t, err, "Batch for current period should succeed")

	// Step 2: Fill the overflow bin partially (70% capacity)
	// This primes the overflow bin with usage so a future overflow might exceed its capacity
	overflowBatch := createTestBatch(
		[]accountQuorumInfo{
			{
				account:    account,
				quorumIDs:  []core.QuorumID{0},
				numSymbols: binLimit * 7 / 10, // 70% of bin limit
			},
		},
	)
	err = batchMeterer.MeterBatch(ctx, overflowBatch, overflowPeriodTime)
	assert.NoError(t, err, "Batch for overflow period should succeed")

	// Step 3: Try to add a batch to the current period that would overflow
	// The overflow would be (70% + 50% - 100%) = 20% of the bin limit
	// But the overflow bin only has 30% capacity left, which is enough
	smallOverflowBatch := createTestBatch(
		[]accountQuorumInfo{
			{
				account:    account,
				quorumIDs:  []core.QuorumID{0},
				numSymbols: binLimit * 5 / 10, // 50% of bin limit
			},
		},
	)
	err = batchMeterer.MeterBatch(ctx, smallOverflowBatch, now)
	assert.NoError(t, err, "Small overflow batch should succeed")

	// Step 4: Try to add another batch that would overflow too much
	// The current bin is already at capacity, and the overflow would exceed the overflow bin's capacity
	largeOverflowBatch := createTestBatch(
		[]accountQuorumInfo{
			{
				account:    account,
				quorumIDs:  []core.QuorumID{0},
				numSymbols: binLimit * 4 / 10, // 40% of bin limit
			},
		},
	)
	err = batchMeterer.MeterBatch(ctx, largeOverflowBatch, now)
	assert.Error(t, err, "Large overflow batch should fail")
	assert.Contains(t, err.Error(), "bin already filled", "Error should mention bin already filled")

	mockState.AssertExpectations(t)
}

// Helper type for test batch creation
type accountQuorumInfo struct {
	account    gethcommon.Address
	quorumIDs  []core.QuorumID
	numSymbols uint64
}

// createTestBatch creates a test batch with the specified accounts, quorums, and symbol counts
func createTestBatch(infos []accountQuorumInfo) *corev2.Batch {
	blobCertificates := make([]*corev2.BlobCertificate, 0, len(infos))

	for _, info := range infos {
		// Create blob header with the account's payment metadata
		blobHeader := &corev2.BlobHeader{
			BlobCommitments: encoding.BlobCommitments{
				Length: uint(info.numSymbols),
			},
			QuorumNumbers: info.quorumIDs,
			PaymentMetadata: core.PaymentMetadata{
				AccountID:         info.account,
				Timestamp:         time.Now().Unix(),
				CumulativePayment: big.NewInt(1000), // Dummy value for testing
			},
		}

		// Create a blob certificate with the header
		blobCert := &corev2.BlobCertificate{
			BlobHeader: blobHeader,
			Signature:  make([]byte, 65), // Dummy signature
		}

		blobCertificates = append(blobCertificates, blobCert)
	}

	// Create the batch with the blob certificates
	return &corev2.Batch{
		BatchHeader: &corev2.BatchHeader{
			ReferenceBlockNumber: 12345,
		},
		BlobCertificates: blobCertificates,
	}
}

// containsAll checks if a slice contains all elements from another slice
func containsAll(slice, contains []core.QuorumID) bool {
	for _, v := range contains {
		found := false
		for _, s := range slice {
			if s == v {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// setupMockOnchainPaymentState creates and configures a mock OnchainPaymentState
func setupMockOnchainPaymentState() *coremock.MockOnchainPaymentState {
	mockState := new(coremock.MockOnchainPaymentState)

	// Set up common mock behavior for methods that might be called implicitly
	mockState.On("RefreshOnchainPaymentState", mock.Anything).Return(nil).Maybe()

	// Add a default stub for GetReservationWindow to avoid nil pointer issues
	mockState.On("GetReservationWindow").Return(uint64(3600)).Maybe()
	mockState.On("GetMinNumSymbols").Return(uint64(32)).Maybe()

	return mockState
}
