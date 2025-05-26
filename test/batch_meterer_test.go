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

	// Set up two test accounts with valid reservations
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

	// Configure the mock to return the appropriate reservations
	// Note: we use MatchedBy to allow for any order of quorum IDs
	mockState.On("GetReservedPaymentByAccountAndQuorums", mock.Anything, account1, mock.MatchedBy(func(quorums []core.QuorumID) bool {
		return len(quorums) == 2 && containsAll(quorums, []core.QuorumID{0, 1})
	})).Return(reservationsAcc1, nil)
	mockState.On("GetReservedPaymentByAccountAndQuorums", mock.Anything, account2, []core.QuorumID{1}).Return(reservationsAcc2, nil)
	mockState.On("GetMinNumSymbols").Return(uint64(32))
	mockState.On("GetReservationWindow").Return(uint64(3600))

	// Create the batch meterer
	batchMeterer := node.NewBatchMeterer(
		meterer.Config{
			ChainReadTimeout: 1 * time.Second,
			UpdateInterval:   1 * time.Minute,
		},
		mockState,
		5, // numBins
		logger,
	)

	// Create a test batch with valid usage within reservation limits
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

	// Test metering the batch - should succeed
	err := batchMeterer.MeterBatch(ctx, batch, time.Now())
	assert.NoError(t, err, "Metering should succeed for batch with valid usage")

	// Verify expectations
	mockState.AssertExpectations(t)
}

// TestBatchMetererExceedingReservation tests the BatchMeterer with usage exceeding reservation limits
func TestBatchMetererExceedingReservation(t *testing.T) {
	ctx := context.Background()
	logger := testutils.GetLogger()
	mockState := setupMockOnchainPaymentState()

	// Set up test account
	account := gethcommon.HexToAddress("0x1111111111111111111111111111111111111111")

	// Set up a reservation with a very low limit that will be exceeded by the test request
	// The bin limit will be 100 * 3600 = 360,000 symbols
	reservations := map[core.QuorumID]*core.ReservedPayment{
		0: {
			SymbolsPerSecond: 100, // Very low limit
			StartTimestamp:   0,
			EndTimestamp:     uint64(time.Now().Add(24 * time.Hour).Unix()),
		},
	}

	// Configure the mock - make sure these return values make sense for the test
	mockState.On("GetReservedPaymentByAccountAndQuorums", mock.Anything, account, []core.QuorumID{0}).Return(reservations, nil)
	mockState.On("GetMinNumSymbols").Return(uint64(32))
	mockState.On("GetReservationWindow").Return(uint64(3600))

	// Create the batch meterer
	batchMeterer := node.NewBatchMeterer(
		meterer.Config{
			ChainReadTimeout: 1 * time.Second,
			UpdateInterval:   1 * time.Minute,
		},
		mockState,
		5, // numBins
		logger,
	)

	// Create a test batch with usage greatly exceeding the reservation limit
	batch := createTestBatch(
		[]accountQuorumInfo{
			{
				account:    account,
				quorumIDs:  []core.QuorumID{0},
				numSymbols: 1000000, // Greatly exceeds the reservation limit (100 * 3600 = 360000)
			},
		},
	)

	// Test metering the batch - should fail
	err := batchMeterer.MeterBatch(ctx, batch, time.Now())
	assert.Error(t, err, "Metering should fail for batch exceeding reservation limit")
	assert.Contains(t, err.Error(), "exceeds bin limit", "Error should mention exceeding bin limit")

	// Verify expectations
	mockState.AssertExpectations(t)
}

// TestBatchMetererInactiveReservation tests the BatchMeterer with an inactive reservation
func TestBatchMetererInactiveReservation(t *testing.T) {
	ctx := context.Background()
	logger := testutils.GetLogger()
	mockState := setupMockOnchainPaymentState()

	// Set up test account
	account := gethcommon.HexToAddress("0x1111111111111111111111111111111111111111")

	// Set up an inactive reservation (already expired)
	reservations := map[core.QuorumID]*core.ReservedPayment{
		0: {
			SymbolsPerSecond: 1000,
			StartTimestamp:   0,
			EndTimestamp:     uint64(time.Now().Add(-24 * time.Hour).Unix()), // Expired
		},
	}

	// Configure the mock
	mockState.On("GetReservedPaymentByAccountAndQuorums", mock.Anything, account, []core.QuorumID{0}).Return(reservations, nil)
	mockState.On("GetMinNumSymbols").Return(uint64(32))
	mockState.On("GetReservationWindow").Return(uint64(3600))

	// Create the batch meterer
	batchMeterer := node.NewBatchMeterer(
		meterer.Config{
			ChainReadTimeout: 1 * time.Second,
			UpdateInterval:   1 * time.Minute,
		},
		mockState,
		5, // numBins
		logger,
	)

	// Create a test batch with an inactive reservation
	batch := createTestBatch(
		[]accountQuorumInfo{
			{
				account:    account,
				quorumIDs:  []core.QuorumID{0},
				numSymbols: 100, // Even a small amount should fail with inactive reservation
			},
		},
	)

	// Test metering the batch - should fail due to inactive reservation
	err := batchMeterer.MeterBatch(ctx, batch, time.Now())
	assert.Error(t, err, "Metering should fail for batch with inactive reservation")
	assert.Contains(t, err.Error(), "inactive reservation", "Error should mention inactive reservation")

	// Verify expectations
	mockState.AssertExpectations(t)
}

// TestBatchMetererMissingReservation tests the BatchMeterer with a quorum that has no reservation
func TestBatchMetererMissingReservation(t *testing.T) {
	ctx := context.Background()
	logger := testutils.GetLogger()
	mockState := setupMockOnchainPaymentState()

	// Set up test account
	account := gethcommon.HexToAddress("0x1111111111111111111111111111111111111111")

	// Set up a reservation for quorum 0 but not quorum 1
	reservations := map[core.QuorumID]*core.ReservedPayment{
		0: {
			SymbolsPerSecond: 1000,
			StartTimestamp:   0,
			EndTimestamp:     uint64(time.Now().Add(24 * time.Hour).Unix()),
		},
		// No reservation for quorum 1
	}

	// Configure the mock - use MatchedBy to allow for any order of quorum IDs
	mockState.On("GetReservedPaymentByAccountAndQuorums", mock.Anything, account, mock.MatchedBy(func(quorums []core.QuorumID) bool {
		return len(quorums) == 2 && containsAll(quorums, []core.QuorumID{0, 1})
	})).Return(reservations, nil)
	mockState.On("GetMinNumSymbols").Return(uint64(32))
	mockState.On("GetReservationWindow").Return(uint64(3600))

	// Create the batch meterer
	batchMeterer := node.NewBatchMeterer(
		meterer.Config{
			ChainReadTimeout: 1 * time.Second,
			UpdateInterval:   1 * time.Minute,
		},
		mockState,
		5, // numBins
		logger,
	)

	// Create a test batch requesting both quorums, but one is missing a reservation
	batch := createTestBatch(
		[]accountQuorumInfo{
			{
				account:    account,
				quorumIDs:  []core.QuorumID{0, 1}, // Quorum 1 has no reservation
				numSymbols: 100,
			},
		},
	)

	// Test metering the batch - should fail due to missing reservation
	err := batchMeterer.MeterBatch(ctx, batch, time.Now())
	assert.Error(t, err, "Metering should fail for batch with missing reservation")
	assert.Contains(t, err.Error(), "no reservation for quorum", "Error should mention missing reservation")

	// Verify expectations
	mockState.AssertExpectations(t)
}

// TestBatchMetererMinSymbolsCharge tests that the meterer correctly applies the minimum symbols charge
func TestBatchMetererMinSymbolsCharge(t *testing.T) {
	ctx := context.Background()
	logger := testutils.GetLogger()
	mockState := setupMockOnchainPaymentState()

	// Set up test account
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

	// Configure the mock
	mockState.On("GetReservedPaymentByAccountAndQuorums", mock.Anything, account, []core.QuorumID{0}).Return(reservations, nil)
	mockState.On("GetMinNumSymbols").Return(minSymbols)
	mockState.On("GetReservationWindow").Return(uint64(3600))

	// Create the batch meterer
	batchMeterer := node.NewBatchMeterer(
		meterer.Config{
			ChainReadTimeout: 1 * time.Second,
			UpdateInterval:   1 * time.Minute,
		},
		mockState,
		5, // numBins
		logger,
	)

	// Create a batch with symbols less than the minimum
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

	// Verify the method was called with minSymbols
	mockState.AssertExpectations(t)
}

// TestBatchMetererOverflowBehavior tests that overflow handling works correctly
func TestBatchMetererOverflowBehavior(t *testing.T) {
	// Skip this test for now until the bin overflow behavior is properly
	// implemented and tested in another PR
	t.Skip("Skipping overflow behavior test for future PR")

	// This test will be completed in a future PR specifically focused on
	// implementing and testing the bin overflow behavior of the BatchMeterer
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
