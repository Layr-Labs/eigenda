package node

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"testing"
	"time"

	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBatchMeterRequest tests the BatchMeterRequest function
func TestBatchMeterRequest(t *testing.T) {
	ctx := context.Background()

	t.Run("all valid reservations", func(t *testing.T) {
		// Create a mock OnchainPayment
		mockChainPayment := new(coremock.MockOnchainPaymentState)

		// Set up the mock expectations
		mockChainPayment.On("GetMinNumSymbols").Return(uint64(10))
		mockChainPayment.On("GetReservationWindow").Return(uint64(3600))

		// Create reservations for account 1
		reservationsAcc1 := map[core.QuorumID]*core.ReservedPayment{
			0: {
				SymbolsPerSecond: 100,
				StartTimestamp:   0,
				EndTimestamp:     9999999999,
			},
			1: {
				SymbolsPerSecond: 100,
				StartTimestamp:   0,
				EndTimestamp:     9999999999,
			},
		}

		// Create reservations for account 2
		reservationsAcc2 := map[core.QuorumID]*core.ReservedPayment{
			1: {
				SymbolsPerSecond: 200,
				StartTimestamp:   0,
				EndTimestamp:     9999999999,
			},
		}

		// Set up the mock for retrieving reservations
		mockChainPayment.On("GetReservedPaymentByAccountAndQuorums", ctx, common.HexToAddress("0x1"), []core.QuorumID{0, 1}).Return(reservationsAcc1, nil)
		mockChainPayment.On("GetReservedPaymentByAccountAndQuorums", ctx, common.HexToAddress("0x2"), []core.QuorumID{1}).Return(reservationsAcc2, nil)

		logger := testutils.GetLogger()

		config := meterer.Config{
			ChainReadTimeout: 1 * time.Second,
			UpdateInterval:   1 * time.Minute,
		}

		batchMeterer := NewBatchMeterer(config, mockChainPayment, 3, logger)

		// Create usage map
		usageMap := map[common.Address]map[core.QuorumID]uint64{
			common.HexToAddress("0x1"): {
				0: 100,
				1: 100,
			},
			common.HexToAddress("0x2"): {
				1: 200,
			},
		}

		// Call the function
		now := time.Now()
		err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)

		// Assert the results
		require.NoError(t, err)

		// Verify expectations
		mockChainPayment.AssertExpectations(t)
	})

	t.Run("expired reservation", func(t *testing.T) {
		// Create a mock OnchainPayment
		mockChainPayment := new(coremock.MockOnchainPaymentState)

		// Set up the mock expectations
		mockChainPayment.On("GetMinNumSymbols").Return(uint64(10))
		mockChainPayment.On("GetReservationWindow").Return(uint64(3600))

		// Create an expired reservation
		reservations := map[core.QuorumID]*core.ReservedPayment{
			0: {
				SymbolsPerSecond: 100,
				StartTimestamp:   0,
				EndTimestamp:     1000, // Past timestamp
			},
		}

		// Set up the mock for retrieving reservations
		mockChainPayment.On("GetReservedPaymentByAccountAndQuorums", ctx, common.HexToAddress("0x1"), []core.QuorumID{0}).Return(reservations, nil)

		logger := testutils.GetLogger()

		config := meterer.Config{
			ChainReadTimeout: 1 * time.Second,
			UpdateInterval:   1 * time.Minute,
		}

		batchMeterer := NewBatchMeterer(config, mockChainPayment, 3, logger)

		// Create usage map
		usageMap := map[common.Address]map[core.QuorumID]uint64{
			common.HexToAddress("0x1"): {
				0: 100,
			},
		}

		// Current time is after the reservation end time
		now := time.Unix(2000, 0)
		err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)

		// Assert the results
		require.Error(t, err)
		assert.Contains(t, err.Error(), "inactive reservation")

		// Verify expectations
		mockChainPayment.AssertExpectations(t)
	})

	t.Run("reservation at capacity", func(t *testing.T) {
		// Create a mock OnchainPayment
		mockChainPayment := new(coremock.MockOnchainPaymentState)

		// Set up the mock expectations
		mockChainPayment.On("GetMinNumSymbols").Return(uint64(10))
		mockChainPayment.On("GetReservationWindow").Return(uint64(3600))

		// Create a reservation with limited capacity
		reservations := map[core.QuorumID]*core.ReservedPayment{
			0: {
				SymbolsPerSecond: 1, // Only 3600 symbols per window
				StartTimestamp:   0,
				EndTimestamp:     9999999999,
			},
		}

		// Set up the mock for retrieving reservations
		mockChainPayment.On("GetReservedPaymentByAccountAndQuorums", ctx, common.HexToAddress("0x1"), []core.QuorumID{0}).Return(reservations, nil)

		logger := testutils.GetLogger()

		config := meterer.Config{
			ChainReadTimeout: 1 * time.Second,
			UpdateInterval:   1 * time.Minute,
		}

		batchMeterer := NewBatchMeterer(config, mockChainPayment, 3, logger)

		// Create a usage record for this account/quorum that already has some usage
		accountUsage := batchMeterer.getOrCreateAccountUsage(common.HexToAddress("0x1"))
		reservationPeriod := meterer.GetReservationPeriod(time.Now().Unix(), 3600)

		// Create a record with usage already at capacity (3600 symbols per window)
		accountUsage.PeriodRecords[0] = []*pb.PeriodRecord{
			{
				Index: uint32(reservationPeriod),
				Usage: 3600, // Already at full capacity
			},
		}

		// Create usage map
		usageMap := map[common.Address]map[core.QuorumID]uint64{
			common.HexToAddress("0x1"): {
				0: 10, // Attempt to add more usage (any amount)
			},
		}

		now := time.Now()
		err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)

		// Assert the results
		require.Error(t, err)
		assert.Contains(t, err.Error(), "bin has already been filled")

		// Verify expectations
		mockChainPayment.AssertExpectations(t)
	})

	t.Run("reservation with overflow", func(t *testing.T) {
		// Create a mock OnchainPayment
		mockChainPayment := new(coremock.MockOnchainPaymentState)

		// Set up the mock expectations
		mockChainPayment.On("GetMinNumSymbols").Return(uint64(10))
		mockChainPayment.On("GetReservationWindow").Return(uint64(3600))

		// Create a reservation with enough capacity for overflow
		reservations := map[core.QuorumID]*core.ReservedPayment{
			0: {
				SymbolsPerSecond: 10, // 36000 symbols per window
				StartTimestamp:   0,
				EndTimestamp:     9999999999,
			},
		}

		// Set up the mock for retrieving reservations
		mockChainPayment.On("GetReservedPaymentByAccountAndQuorums", ctx, common.HexToAddress("0x1"), []core.QuorumID{0}).Return(reservations, nil)

		logger := testutils.GetLogger()

		config := meterer.Config{
			ChainReadTimeout: 1 * time.Second,
			UpdateInterval:   1 * time.Minute,
		}

		batchMeterer := NewBatchMeterer(config, mockChainPayment, 3, logger)

		// Create usage map with usage that would cause overflow
		usageMap := map[common.Address]map[core.QuorumID]uint64{
			common.HexToAddress("0x1"): {
				0: 37000, // Slightly more than the bin limit, but less than 2x
			},
		}

		now := time.Now()
		err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)

		// Assert the results
		require.NoError(t, err)

		// Verify expectations
		mockChainPayment.AssertExpectations(t)
	})

	t.Run("missing reservation for quorum", func(t *testing.T) {
		// Create a mock OnchainPayment
		mockChainPayment := new(coremock.MockOnchainPaymentState)

		// Set up the mock expectations
		mockChainPayment.On("GetMinNumSymbols").Return(uint64(10))
		mockChainPayment.On("GetReservationWindow").Return(uint64(3600))

		// Create a reservation for quorum 0 but not for quorum 1
		reservations := map[core.QuorumID]*core.ReservedPayment{
			0: {
				SymbolsPerSecond: 100,
				StartTimestamp:   0,
				EndTimestamp:     9999999999,
			},
		}

		// Set up the mock for retrieving reservations
		mockChainPayment.On("GetReservedPaymentByAccountAndQuorums", ctx, common.HexToAddress("0x1"), []core.QuorumID{0, 1}).Return(reservations, nil)

		logger := testutils.GetLogger()

		config := meterer.Config{
			ChainReadTimeout: 1 * time.Second,
			UpdateInterval:   1 * time.Minute,
		}

		batchMeterer := NewBatchMeterer(config, mockChainPayment, 3, logger)

		// Create usage map
		usageMap := map[common.Address]map[core.QuorumID]uint64{
			common.HexToAddress("0x1"): {
				0: 100,
				1: 100,
			},
		}

		now := time.Now()
		err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)

		// Assert the results
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no reservation for quorum")

		// Verify expectations
		mockChainPayment.AssertExpectations(t)
	})

	t.Run("error retrieving reservation from chain", func(t *testing.T) {
		// Create a mock OnchainPayment
		mockChainPayment := new(coremock.MockOnchainPaymentState)

		// Set up the mock expectations
		mockChainPayment.On("GetMinNumSymbols").Return(uint64(10))
		mockChainPayment.On("GetReservationWindow").Return(uint64(3600))

		// Set up the mock to return an error
		mockChainPayment.On("GetReservedPaymentByAccountAndQuorums", ctx, common.HexToAddress("0x1"), []core.QuorumID{0}).Return(map[core.QuorumID]*core.ReservedPayment{}, fmt.Errorf("chain error"))

		logger := testutils.GetLogger()

		config := meterer.Config{
			ChainReadTimeout: 1 * time.Second,
			UpdateInterval:   1 * time.Minute,
		}

		batchMeterer := NewBatchMeterer(config, mockChainPayment, 3, logger)

		// Create usage map
		usageMap := map[common.Address]map[core.QuorumID]uint64{
			common.HexToAddress("0x1"): {
				0: 100,
			},
		}

		now := time.Now()
		err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)

		// Assert the results
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get reservations")

		// Verify expectations
		mockChainPayment.AssertExpectations(t)
	})
}

// TestMeterBatch tests the MeterBatch function
func TestMeterBatch(t *testing.T) {
	ctx := context.Background()

	t.Run("successful batch metering", func(t *testing.T) {
		// Create a mock OnchainPayment
		mockChainPayment := new(coremock.MockOnchainPaymentState)

		// Set up the mock expectations
		mockChainPayment.On("GetMinNumSymbols").Return(uint64(10))
		mockChainPayment.On("GetReservationWindow").Return(uint64(3600))

		// Create reservations for accounts
		reservationsAcc1 := map[core.QuorumID]*core.ReservedPayment{
			0: {
				SymbolsPerSecond: 100,
				StartTimestamp:   0,
				EndTimestamp:     9999999999,
			},
			1: {
				SymbolsPerSecond: 100,
				StartTimestamp:   0,
				EndTimestamp:     9999999999,
			},
		}

		reservationsAcc2 := map[core.QuorumID]*core.ReservedPayment{
			1: {
				SymbolsPerSecond: 200,
				StartTimestamp:   0,
				EndTimestamp:     9999999999,
			},
		}

		// Set up the mock for retrieving reservations
		mockChainPayment.On("GetReservedPaymentByAccountAndQuorums", ctx, common.HexToAddress("0x1"), []core.QuorumID{0, 1}).Return(reservationsAcc1, nil)
		mockChainPayment.On("GetReservedPaymentByAccountAndQuorums", ctx, common.HexToAddress("0x2"), []core.QuorumID{1}).Return(reservationsAcc2, nil)

		logger := testutils.GetLogger()

		config := meterer.Config{
			ChainReadTimeout: 1 * time.Second,
			UpdateInterval:   1 * time.Minute,
		}

		batchMeterer := NewBatchMeterer(config, mockChainPayment, 3, logger)

		// Create a sample batch
		batch := createSampleBatch(t)

		// Call the function
		now := time.Now()
		err := batchMeterer.MeterBatch(ctx, batch, now)

		// Assert the results
		require.NoError(t, err)

		// Verify expectations
		mockChainPayment.AssertExpectations(t)
	})

	t.Run("batch with invalid account", func(t *testing.T) {
		// Create a mock OnchainPayment
		mockChainPayment := new(coremock.MockOnchainPaymentState)

		// Set up the mock expectations
		mockChainPayment.On("GetMinNumSymbols").Return(uint64(10))
		mockChainPayment.On("GetReservationWindow").Return(uint64(3600))

		// Create valid reservations for account 1
		reservationsAcc1 := map[core.QuorumID]*core.ReservedPayment{
			0: {
				SymbolsPerSecond: 100,
				StartTimestamp:   0,
				EndTimestamp:     9999999999,
			},
			1: {
				SymbolsPerSecond: 100,
				StartTimestamp:   0,
				EndTimestamp:     9999999999,
			},
		}

		// Account 2 has no reservation for quorum 1
		reservationsAcc2 := map[core.QuorumID]*core.ReservedPayment{}

		// Set up the mock for retrieving reservations
		mockChainPayment.On("GetReservedPaymentByAccountAndQuorums", ctx, common.HexToAddress("0x1"), []core.QuorumID{0, 1}).Return(reservationsAcc1, nil)
		mockChainPayment.On("GetReservedPaymentByAccountAndQuorums", ctx, common.HexToAddress("0x2"), []core.QuorumID{1}).Return(reservationsAcc2, nil)

		logger := testutils.GetLogger()

		config := meterer.Config{
			ChainReadTimeout: 1 * time.Second,
			UpdateInterval:   1 * time.Minute,
		}

		batchMeterer := NewBatchMeterer(config, mockChainPayment, 3, logger)

		// Create a sample batch
		batch := createSampleBatch(t)

		// Call the function
		now := time.Now()
		err := batchMeterer.MeterBatch(ctx, batch, now)

		// Assert the results - the batch should fail due to account 2
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no reservation for quorum")

		// Verify expectations
		mockChainPayment.AssertExpectations(t)
	})

	t.Run("nil batch", func(t *testing.T) {
		// Create a mock OnchainPayment
		mockChainPayment := new(coremock.MockOnchainPaymentState)

		logger := testutils.GetLogger()

		config := meterer.Config{
			ChainReadTimeout: 1 * time.Second,
			UpdateInterval:   1 * time.Minute,
		}

		batchMeterer := NewBatchMeterer(config, mockChainPayment, 3, logger)

		// Call with nil batch
		now := time.Now()
		err := batchMeterer.MeterBatch(ctx, nil, now)

		// Assert the results
		require.Error(t, err)
		assert.Contains(t, err.Error(), "batch is nil")
	})
}

// Helper function to create a sample batch for testing
func createSampleBatch(t *testing.T) *corev2.Batch {
	// Create a batch header
	batchHeader := &corev2.BatchHeader{
		BatchRoot:            [32]byte{1, 2, 3},
		ReferenceBlockNumber: 100,
	}

	// Create blob certificates
	blobCerts := []*corev2.BlobCertificate{
		{
			BlobHeader: &corev2.BlobHeader{
				BlobVersion: 1,
				BlobCommitments: encoding.BlobCommitments{
					Length: 100,
				},
				QuorumNumbers: []core.QuorumID{0, 1},
				PaymentMetadata: core.PaymentMetadata{
					AccountID:         common.HexToAddress("0x1"),
					Timestamp:         time.Now().UnixNano(),
					CumulativePayment: big.NewInt(0),
				},
			},
			Signature: []byte{1, 2, 3},
			RelayKeys: []corev2.RelayKey{1, 2},
		},
		{
			BlobHeader: &corev2.BlobHeader{
				BlobVersion: 1,
				BlobCommitments: encoding.BlobCommitments{
					Length: 200,
				},
				QuorumNumbers: []core.QuorumID{1},
				PaymentMetadata: core.PaymentMetadata{
					AccountID:         common.HexToAddress("0x2"),
					Timestamp:         time.Now().UnixNano(),
					CumulativePayment: big.NewInt(0),
				},
			},
			Signature: []byte{4, 5, 6},
			RelayKeys: []corev2.RelayKey{3},
		},
	}

	return &corev2.Batch{
		BatchHeader:      batchHeader,
		BlobCertificates: blobCerts,
	}
}

// TestBatchRequestUsageEdgeCases tests edge cases in batch request usage calculation
func TestBatchRequestUsageEdgeCases(t *testing.T) {
	logger := testutils.GetLogger()
	mockState := new(coremock.MockOnchainPaymentState)
	mockState.On("GetMinNumSymbols").Return(uint64(32))
	mockState.On("GetReservationWindow").Return(uint64(3600))

	config := meterer.Config{
		ChainReadTimeout: 1 * time.Second,
		UpdateInterval:   1 * time.Minute,
	}

	batchMeterer := NewBatchMeterer(config, mockState, 3, logger)

	t.Run("empty batch", func(t *testing.T) {
		batch := &corev2.Batch{
			BatchHeader:      &corev2.BatchHeader{},
			BlobCertificates: []*corev2.BlobCertificate{},
		}
		_, err := batchMeterer.BatchRequestUsage(batch)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "batch is nil or empty")
	})

	t.Run("nil blob header", func(t *testing.T) {
		batch := &corev2.Batch{
			BatchHeader: &corev2.BatchHeader{},
			BlobCertificates: []*corev2.BlobCertificate{
				{
					BlobHeader: nil,
				},
			},
		}
		_, err := batchMeterer.BatchRequestUsage(batch)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "blob certificate has nil header")
	})

	t.Run("zero length blob", func(t *testing.T) {
		batch := createTestBatch([]accountQuorumInfo{
			{
				account:    common.HexToAddress("0x1"),
				quorumIDs:  []core.QuorumID{0},
				numSymbols: 0,
			},
		})
		usageMap, err := batchMeterer.BatchRequestUsage(batch)
		require.NoError(t, err)
		// Should use minimum symbols
		assert.Equal(t, uint64(32), usageMap[common.HexToAddress("0x1")][0])
	})

	t.Run("duplicate account quorum", func(t *testing.T) {
		// Create a batch with two blobs for the same account and quorum
		// Each blob should get minimum 32 symbols
		batch := createTestBatch([]accountQuorumInfo{
			{
				account:    common.HexToAddress("0x1"),
				quorumIDs:  []core.QuorumID{0},
				numSymbols: 32, // Minimum symbols
			},
			{
				account:    common.HexToAddress("0x1"),
				quorumIDs:  []core.QuorumID{0},
				numSymbols: 32, // Minimum symbols
			},
		})
		usageMap, err := batchMeterer.BatchRequestUsage(batch)
		require.NoError(t, err)
		// Should aggregate usage, each blob gets minimum 32 symbols
		assert.Equal(t, uint64(64), usageMap[common.HexToAddress("0x1")][0])
	})

	t.Run("very large usage", func(t *testing.T) {
		batch := createTestBatch([]accountQuorumInfo{
			{
				account:    common.HexToAddress("0x1"),
				quorumIDs:  []core.QuorumID{0},
				numSymbols: math.MaxUint64 - 100,
			},
		})
		usageMap, err := batchMeterer.BatchRequestUsage(batch)
		require.NoError(t, err)
		// Should handle large numbers without overflow
		assert.Greater(t, usageMap[common.HexToAddress("0x1")][0], uint64(0))
	})
}

// TestOverflowEdgeCases tests edge cases in overflow handling
func TestOverflowEdgeCases(t *testing.T) {
	ctx := context.Background()
	logger := testutils.GetLogger()
	mockState := new(coremock.MockOnchainPaymentState)
	mockState.On("GetMinNumSymbols").Return(uint64(32))
	mockState.On("GetReservationWindow").Return(uint64(3600))

	config := meterer.Config{
		ChainReadTimeout: 1 * time.Second,
		UpdateInterval:   1 * time.Minute,
	}

	batchMeterer := NewBatchMeterer(config, mockState, 3, logger)

	t.Run("exact bin limit", func(t *testing.T) {
		// Create a reservation with exact bin limit
		binLimit := uint64(3600) // 1 symbol per second
		reservations := map[core.QuorumID]*core.ReservedPayment{
			0: {
				SymbolsPerSecond: 1,
				StartTimestamp:   0,
				EndTimestamp:     9999999999,
			},
		}
		mockState.On("GetReservedPaymentByAccountAndQuorums", ctx, common.HexToAddress("0x1"), []core.QuorumID{0}).Return(reservations, nil)

		usageMap := map[common.Address]map[core.QuorumID]uint64{
			common.HexToAddress("0x1"): {0: binLimit},
		}

		err := batchMeterer.BatchMeterRequest(ctx, usageMap, time.Now())
		require.NoError(t, err)

		// Verify usage is exactly at bin limit
		accountUsage := batchMeterer.getOrCreateAccountUsage(common.HexToAddress("0x1"))
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()

		currentPeriod := meterer.GetReservationPeriod(time.Now().Unix(), 3600)
		assert.Equal(t, binLimit, accountUsage.PeriodRecords[0][currentPeriod%3].Usage)
	})

	t.Run("just over bin limit", func(t *testing.T) {
		// Create a reservation
		binLimit := uint64(3600)
		reservations := map[core.QuorumID]*core.ReservedPayment{
			0: {
				SymbolsPerSecond: 1,
				StartTimestamp:   0,
				EndTimestamp:     9999999999,
			},
		}
		mockState.On("GetReservedPaymentByAccountAndQuorums", ctx, common.HexToAddress("0x1"), []core.QuorumID{0}).Return(reservations, nil)

		// Initialize the period record with some usage
		accountUsage := batchMeterer.getOrCreateAccountUsage(common.HexToAddress("0x1"))
		currentPeriod := meterer.GetReservationPeriod(time.Now().Unix(), 3600)
		accountUsage.Lock.Lock()
		reservationWindow := uint64(3600)
		numBins := uint32(3)
		relativeIndex := uint32((currentPeriod / reservationWindow) % uint64(numBins))
		accountUsage.PeriodRecords[0] = make([]*pb.PeriodRecord, numBins)
		for i := range accountUsage.PeriodRecords[0] {
			accountUsage.PeriodRecords[0][i] = &pb.PeriodRecord{}
		}
		accountUsage.PeriodRecords[0][relativeIndex].Index = uint32(currentPeriod)
		accountUsage.PeriodRecords[0][relativeIndex].Usage = binLimit - 1
		accountUsage.Lock.Unlock()

		// Try to use just over bin limit
		usageMap := map[common.Address]map[core.QuorumID]uint64{
			common.HexToAddress("0x1"): {0: 2}, // Add 2 more to exceed limit
		}

		err := batchMeterer.BatchMeterRequest(ctx, usageMap, time.Now())
		require.NoError(t, err)

		// Verify overflow was handled
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()

		relativeIndex = uint32((currentPeriod / reservationWindow) % uint64(numBins))
		assert.Equal(t, binLimit, accountUsage.PeriodRecords[0][relativeIndex].Usage)
		overflowPeriod := currentPeriod + reservationWindow // Add reservation window
		overflowRelativeIndex := uint32((overflowPeriod / reservationWindow) % uint64(numBins))
		assert.Equal(t, uint64(1), accountUsage.PeriodRecords[0][overflowRelativeIndex].Usage)
	})

	t.Run("exactly 2x bin limit", func(t *testing.T) {
		// Create a reservation
		binLimit := uint64(3600)
		reservations := map[core.QuorumID]*core.ReservedPayment{
			0: {
				SymbolsPerSecond: 1,
				StartTimestamp:   0,
				EndTimestamp:     9999999999,
			},
		}
		mockState.On("GetReservedPaymentByAccountAndQuorums", ctx, common.HexToAddress("0x1"), []core.QuorumID{0}).Return(reservations, nil)

		// Initialize the period record with some usage
		accountUsage := batchMeterer.getOrCreateAccountUsage(common.HexToAddress("0x1"))
		currentPeriod := meterer.GetReservationPeriod(time.Now().Unix(), 3600)
		accountUsage.Lock.Lock()
		accountUsage.PeriodRecords[0] = make([]*pb.PeriodRecord, 3)
		accountUsage.PeriodRecords[0][currentPeriod%3] = &pb.PeriodRecord{
			Index: uint32(currentPeriod),
			Usage: 0,
		}
		accountUsage.Lock.Unlock()

		// Try to use exactly 2x bin limit
		usageMap := map[common.Address]map[core.QuorumID]uint64{
			common.HexToAddress("0x1"): {0: 2 * binLimit},
		}

		err := batchMeterer.BatchMeterRequest(ctx, usageMap, time.Now())
		require.NoError(t, err)

		// Verify overflow was handled
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()

		assert.Equal(t, binLimit, accountUsage.PeriodRecords[0][currentPeriod%3].Usage)
		overflowPeriod := currentPeriod + 3600 // Add reservation window
		assert.Equal(t, binLimit, accountUsage.PeriodRecords[0][overflowPeriod%3].Usage)
	})

	t.Run("over 2x bin limit", func(t *testing.T) {
		// Create a reservation
		binLimit := uint64(3600)
		reservations := map[core.QuorumID]*core.ReservedPayment{
			0: {
				SymbolsPerSecond: 1,
				StartTimestamp:   0,
				EndTimestamp:     9999999999,
			},
		}
		mockState.On("GetReservedPaymentByAccountAndQuorums", ctx, common.HexToAddress("0x1"), []core.QuorumID{0}).Return(reservations, nil)

		// Initialize the period record with some usage
		accountUsage := batchMeterer.getOrCreateAccountUsage(common.HexToAddress("0x1"))
		currentPeriod := meterer.GetReservationPeriod(time.Now().Unix(), 3600)
		accountUsage.Lock.Lock()
		accountUsage.PeriodRecords[0] = make([]*pb.PeriodRecord, 3)
		accountUsage.PeriodRecords[0][currentPeriod%3] = &pb.PeriodRecord{
			Index: uint32(currentPeriod),
			Usage: 0,
		}
		accountUsage.Lock.Unlock()

		// Try to use over 2x bin limit
		usageMap := map[common.Address]map[core.QuorumID]uint64{
			common.HexToAddress("0x1"): {0: 2*binLimit + 1},
		}

		err := batchMeterer.BatchMeterRequest(ctx, usageMap, time.Now())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "overflow usage exceeds bin limit")
	})
}

// TestPeriodRecordEdgeCases tests edge cases in period record management
func TestPeriodRecordEdgeCases(t *testing.T) {
	ctx := context.Background()
	logger := testutils.GetLogger()
	mockState := new(coremock.MockOnchainPaymentState)
	mockState.On("GetMinNumSymbols").Return(uint64(32))
	mockState.On("GetReservationWindow").Return(uint64(3600))

	config := meterer.Config{
		ChainReadTimeout: 1 * time.Second,
		UpdateInterval:   1 * time.Minute,
	}

	batchMeterer := NewBatchMeterer(config, mockState, 3, logger)

	t.Run("period transition", func(t *testing.T) {
		// Create a reservation that spans multiple periods
		reservations := map[core.QuorumID]*core.ReservedPayment{
			0: {
				SymbolsPerSecond: 100,
				StartTimestamp:   0,
				EndTimestamp:     9999999999,
			},
		}
		mockState.On("GetReservedPaymentByAccountAndQuorums", ctx, common.HexToAddress("0x1"), []core.QuorumID{0}).Return(reservations, nil)

		// First request in current period
		now := time.Now()
		usageMap := map[common.Address]map[core.QuorumID]uint64{
			common.HexToAddress("0x1"): {0: 100},
		}
		err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)
		require.NoError(t, err)

		// Second request in next period
		nextPeriod := now.Add(time.Duration(3600) * time.Second)
		err = batchMeterer.BatchMeterRequest(ctx, usageMap, nextPeriod)
		require.NoError(t, err)

		// Verify both periods have correct usage
		accountUsage := batchMeterer.getOrCreateAccountUsage(common.HexToAddress("0x1"))
		currentPeriod := meterer.GetReservationPeriod(now.Unix(), 3600)
		nextPeriodIndex := meterer.GetReservationPeriod(nextPeriod.Unix(), 3600)

		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()

		assert.Equal(t, uint64(100), accountUsage.PeriodRecords[0][currentPeriod%3].Usage)
		assert.Equal(t, uint64(100), accountUsage.PeriodRecords[0][nextPeriodIndex%3].Usage)
	})

	t.Run("circular buffer wrapping", func(t *testing.T) {
		// Create a reservation
		reservations := map[core.QuorumID]*core.ReservedPayment{
			0: {
				SymbolsPerSecond: 100,
				StartTimestamp:   0,
				EndTimestamp:     9999999999,
			},
		}
		mockState.On("GetReservedPaymentByAccountAndQuorums", ctx, common.HexToAddress("0x1"), []core.QuorumID{0}).Return(reservations, nil)

		// Make requests across 4 periods to test buffer wrapping
		now := time.Now()
		usageMap := map[common.Address]map[core.QuorumID]uint64{
			common.HexToAddress("0x1"): {0: 100},
		}

		for i := 0; i < 4; i++ {
			periodTime := now.Add(time.Duration(i*3600) * time.Second)
			err := batchMeterer.BatchMeterRequest(ctx, usageMap, periodTime)
			require.NoError(t, err)
		}

		// Verify the oldest record was overwritten
		accountUsage := batchMeterer.getOrCreateAccountUsage(common.HexToAddress("0x1"))
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()

		// Should only have 3 records (buffer size)
		records := accountUsage.PeriodRecords[0]
		assert.Equal(t, 3, len(records))
	})
}

// TestReservationEdgeCases tests edge cases in reservation validation
func TestReservationEdgeCases(t *testing.T) {
	ctx := context.Background()
	logger := testutils.GetLogger()
	mockState := new(coremock.MockOnchainPaymentState)
	mockState.On("GetMinNumSymbols").Return(uint64(32))
	mockState.On("GetReservationWindow").Return(uint64(3600))

	config := meterer.Config{
		ChainReadTimeout: 1 * time.Second,
		UpdateInterval:   1 * time.Minute,
	}

	batchMeterer := NewBatchMeterer(config, mockState, 3, logger)

	t.Run("just started reservation", func(t *testing.T) {
		// Create a reservation that just started
		now := time.Now()
		reservations := map[core.QuorumID]*core.ReservedPayment{
			0: {
				SymbolsPerSecond: 100,
				StartTimestamp:   uint64(now.Unix()),
				EndTimestamp:     uint64(now.Add(24 * time.Hour).Unix()),
			},
		}
		mockState.On("GetReservedPaymentByAccountAndQuorums", ctx, common.HexToAddress("0x1"), []core.QuorumID{0}).Return(reservations, nil)

		usageMap := map[common.Address]map[core.QuorumID]uint64{
			common.HexToAddress("0x1"): {0: 100},
		}

		err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)
		require.NoError(t, err)
	})

	t.Run("just expired reservation", func(t *testing.T) {
		// Create a reservation that just expired
		now := time.Now()
		reservations := map[core.QuorumID]*core.ReservedPayment{
			0: {
				SymbolsPerSecond: 100,
				StartTimestamp:   uint64(now.Add(-24 * time.Hour).Unix()),
				EndTimestamp:     uint64(now.Unix()),
			},
		}
		mockState.On("GetReservedPaymentByAccountAndQuorums", ctx, common.HexToAddress("0x1"), []core.QuorumID{0}).Return(reservations, nil)

		usageMap := map[common.Address]map[core.QuorumID]uint64{
			common.HexToAddress("0x1"): {0: 100},
		}

		err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "reservation has expired")
	})

	t.Run("zero symbols per second", func(t *testing.T) {
		// Create a reservation with zero symbols per second
		now := time.Now()
		reservations := map[core.QuorumID]*core.ReservedPayment{
			0: {
				SymbolsPerSecond: 0,
				StartTimestamp:   uint64(now.Unix()),
				EndTimestamp:     uint64(now.Add(24 * time.Hour).Unix()),
			},
		}
		mockState.On("GetReservedPaymentByAccountAndQuorums", ctx, common.HexToAddress("0x1"), []core.QuorumID{0}).Return(reservations, nil)

		usageMap := map[common.Address]map[core.QuorumID]uint64{
			common.HexToAddress("0x1"): {0: 100},
		}

		err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid symbols per second")
	})

	t.Run("invalid timestamps", func(t *testing.T) {
		// Create a reservation with invalid timestamps
		now := time.Now()
		reservations := map[core.QuorumID]*core.ReservedPayment{
			0: {
				SymbolsPerSecond: 100,
				StartTimestamp:   uint64(now.Add(24 * time.Hour).Unix()), // Start time in the future
				EndTimestamp:     uint64(now.Unix()),                     // End time in the past
			},
		}
		mockState.On("GetReservedPaymentByAccountAndQuorums", ctx, common.HexToAddress("0x1"), []core.QuorumID{0}).Return(reservations, nil)

		usageMap := map[common.Address]map[core.QuorumID]uint64{
			common.HexToAddress("0x1"): {0: 100},
		}

		err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid reservation timestamps")
	})
}

// Helper function to create a test batch
func createTestBatch(infos []accountQuorumInfo) *corev2.Batch {
	batchHeader := &corev2.BatchHeader{
		BatchRoot:            [32]byte{1, 2, 3},
		ReferenceBlockNumber: 100,
	}

	blobCerts := make([]*corev2.BlobCertificate, len(infos))
	for i, info := range infos {
		blobCerts[i] = &corev2.BlobCertificate{
			BlobHeader: &corev2.BlobHeader{
				BlobVersion: 1,
				BlobCommitments: encoding.BlobCommitments{
					Length: uint(info.numSymbols),
				},
				QuorumNumbers: info.quorumIDs,
				PaymentMetadata: core.PaymentMetadata{
					AccountID:         info.account,
					Timestamp:         time.Now().UnixNano(),
					CumulativePayment: big.NewInt(0),
				},
			},
			Signature: []byte{1, 2, 3},
			RelayKeys: []corev2.RelayKey{1, 2},
		}
	}

	return &corev2.Batch{
		BatchHeader:      batchHeader,
		BlobCertificates: blobCerts,
	}
}

// Helper struct for creating test batches
type accountQuorumInfo struct {
	account    common.Address
	quorumIDs  []core.QuorumID
	numSymbols uint64
}
