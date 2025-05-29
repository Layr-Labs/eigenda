package node

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestBatchToRequestInfos tests the BatchToRequestInfos function
func TestBatchToRequestInfos(t *testing.T) {
	// Create a mock OnchainPayment
	mockChainPayment := new(coremock.MockOnchainPaymentState)
	// Set up required method calls
	mockChainPayment.On("GetReservedPaymentByAccount", mock.Anything, mock.Anything).Return(map[core.QuorumID]*core.ReservedPayment{}, nil).Maybe()

	logger := testutils.GetLogger()

	config := meterer.Config{
		ChainReadTimeout: 1 * time.Second,
		UpdateInterval:   1 * time.Minute,
	}

	batchMeterer := NewBatchMeterer(config, mockChainPayment, 3, logger)

	t.Run("success case", func(t *testing.T) {
		// Create a sample batch
		batch := createSampleBatch(t)

		// Call the function
		requestInfos, err := batchMeterer.BatchToRequestInfos(batch)

		// Assert the results
		require.NoError(t, err)
		require.Len(t, requestInfos, 2)

		// Check the first request
		assert.Equal(t, common.HexToAddress("0x1"), requestInfos[0].AccountID)
		assert.Equal(t, []core.QuorumID{0, 1}, requestInfos[0].QuorumIDs)
		assert.Equal(t, uint64(100), requestInfos[0].NumSymbols)

		// Check the second request
		assert.Equal(t, common.HexToAddress("0x2"), requestInfos[1].AccountID)
		assert.Equal(t, []core.QuorumID{1}, requestInfos[1].QuorumIDs)
		assert.Equal(t, uint64(200), requestInfos[1].NumSymbols)
	})

	t.Run("nil batch", func(t *testing.T) {
		// Call with nil batch
		_, err := batchMeterer.BatchToRequestInfos(nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "batch is nil")
	})

	t.Run("empty blob certificates", func(t *testing.T) {
		// Create a batch with no blob certificates
		batch := &corev2.Batch{
			BatchHeader:      &corev2.BatchHeader{},
			BlobCertificates: []*corev2.BlobCertificate{},
		}

		// Call the function
		_, err := batchMeterer.BatchToRequestInfos(batch)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "batch has no blob certificates")
	})

	t.Run("nil blob header", func(t *testing.T) {
		// Create a batch with a blob certificate that has a nil header
		batch := &corev2.Batch{
			BatchHeader: &corev2.BatchHeader{},
			BlobCertificates: []*corev2.BlobCertificate{
				{
					BlobHeader: nil,
				},
			},
		}

		// Call the function
		_, err := batchMeterer.BatchToRequestInfos(batch)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "blob certificate has nil header")
	})
}

// TestAggregateRequests tests the AggregateRequests function
func TestAggregateRequests(t *testing.T) {
	// Create a mock OnchainPayment
	mockChainPayment := new(coremock.MockOnchainPaymentState)
	mockChainPayment.On("GetMinNumSymbols").Return(uint64(10))
	// Set up required method calls
	mockChainPayment.On("GetReservedPaymentByAccount", mock.Anything, mock.Anything).Return(map[core.QuorumID]*core.ReservedPayment{}, nil).Maybe()

	logger := testutils.GetLogger()

	config := meterer.Config{
		ChainReadTimeout: 1 * time.Second,
		UpdateInterval:   1 * time.Minute,
	}

	batchMeterer := NewBatchMeterer(config, mockChainPayment, 3, logger)

	t.Run("aggregate single account single quorum", func(t *testing.T) {
		// Create a sample request
		requests := []BatchRequestInfo{
			{
				AccountID:  common.HexToAddress("0x1"),
				QuorumIDs:  []core.QuorumID{0},
				NumSymbols: 100,
			},
		}

		// Call the function
		aggregated := batchMeterer.AggregateRequests(requests)

		// Assert the results
		require.Len(t, aggregated, 1)
		assert.Equal(t, uint64(100), aggregated[common.HexToAddress("0x1")][0])
	})

	t.Run("aggregate single account multiple quorums", func(t *testing.T) {
		// Create sample requests
		requests := []BatchRequestInfo{
			{
				AccountID:  common.HexToAddress("0x1"),
				QuorumIDs:  []core.QuorumID{0, 1},
				NumSymbols: 100,
			},
		}

		// Call the function
		aggregated := batchMeterer.AggregateRequests(requests)

		// Assert the results
		require.Len(t, aggregated, 1)
		assert.Equal(t, uint64(100), aggregated[common.HexToAddress("0x1")][0])
		assert.Equal(t, uint64(100), aggregated[common.HexToAddress("0x1")][1])
	})

	t.Run("aggregate multiple accounts multiple quorums", func(t *testing.T) {
		// Create sample requests
		requests := []BatchRequestInfo{
			{
				AccountID:  common.HexToAddress("0x1"),
				QuorumIDs:  []core.QuorumID{0, 1},
				NumSymbols: 100,
			},
			{
				AccountID:  common.HexToAddress("0x2"),
				QuorumIDs:  []core.QuorumID{1},
				NumSymbols: 200,
			},
			{
				AccountID:  common.HexToAddress("0x1"),
				QuorumIDs:  []core.QuorumID{0},
				NumSymbols: 50,
			},
		}

		// Call the function
		aggregated := batchMeterer.AggregateRequests(requests)

		// Assert the results
		require.Len(t, aggregated, 2)
		assert.Equal(t, uint64(150), aggregated[common.HexToAddress("0x1")][0])
		assert.Equal(t, uint64(100), aggregated[common.HexToAddress("0x1")][1])
		assert.Equal(t, uint64(200), aggregated[common.HexToAddress("0x2")][1])
	})

	t.Run("empty requests", func(t *testing.T) {
		// Call with empty requests
		aggregated := batchMeterer.AggregateRequests([]BatchRequestInfo{})
		assert.Empty(t, aggregated)
	})
}

// TestBatchMeterRequest tests the BatchMeterRequest function
func TestBatchMeterRequest(t *testing.T) {
	ctx := context.Background()

	t.Run("all valid reservations", func(t *testing.T) {
		// Create a mock OnchainPayment
		mockChainPayment := new(coremock.MockOnchainPaymentState)

		// Set up the mock expectations
		mockChainPayment.On("GetMinNumSymbols").Return(uint64(10))
		mockChainPayment.On("GetReservationWindow").Return(uint64(3600))
		// Set up required method calls
		mockChainPayment.On("GetReservedPaymentByAccount", mock.Anything, mock.Anything).Return(map[core.QuorumID]*core.ReservedPayment{}, nil).Maybe()

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

		// Create sample requests
		requests := []BatchRequestInfo{
			{
				AccountID:  common.HexToAddress("0x1"),
				QuorumIDs:  []core.QuorumID{0, 1},
				NumSymbols: 100,
			},
			{
				AccountID:  common.HexToAddress("0x2"),
				QuorumIDs:  []core.QuorumID{1},
				NumSymbols: 200,
			},
		}

		// Call the function
		now := time.Now()
		err := batchMeterer.BatchMeterRequest(ctx, requests, now)

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
		// Set up required method calls
		mockChainPayment.On("GetReservedPaymentByAccount", mock.Anything, mock.Anything).Return(map[core.QuorumID]*core.ReservedPayment{}, nil).Maybe()

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

		// Create sample requests
		requests := []BatchRequestInfo{
			{
				AccountID:  common.HexToAddress("0x1"),
				QuorumIDs:  []core.QuorumID{0},
				NumSymbols: 100,
			},
		}

		// Current time is after the reservation end time
		now := time.Unix(2000, 0)
		err := batchMeterer.BatchMeterRequest(ctx, requests, now)

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
		// Set up required method calls
		mockChainPayment.On("GetReservedPaymentByAccount", mock.Anything, mock.Anything).Return(map[core.QuorumID]*core.ReservedPayment{}, nil).Maybe()

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
		accountUsage.UsageRecords[0] = []UsageRecord{
			{
				Index: reservationPeriod,
				Usage: 3600, // Already at full capacity
			},
		}

		// Create sample requests
		requests := []BatchRequestInfo{
			{
				AccountID:  common.HexToAddress("0x1"),
				QuorumIDs:  []core.QuorumID{0},
				NumSymbols: 10, // Attempt to add more usage (any amount)
			},
		}

		now := time.Now()
		err := batchMeterer.BatchMeterRequest(ctx, requests, now)

		// Assert the results
		require.Error(t, err)
		assert.Contains(t, err.Error(), "bin already filled")

		// Verify expectations
		mockChainPayment.AssertExpectations(t)
	})

	t.Run("reservation with overflow", func(t *testing.T) {
		// Create a mock OnchainPayment
		mockChainPayment := new(coremock.MockOnchainPaymentState)

		// Set up the mock expectations
		mockChainPayment.On("GetMinNumSymbols").Return(uint64(10))
		mockChainPayment.On("GetReservationWindow").Return(uint64(3600))
		// Set up required method calls
		mockChainPayment.On("GetReservedPaymentByAccount", mock.Anything, mock.Anything).Return(map[core.QuorumID]*core.ReservedPayment{}, nil).Maybe()

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

		// Create sample requests
		requests := []BatchRequestInfo{
			{
				AccountID:  common.HexToAddress("0x1"),
				QuorumIDs:  []core.QuorumID{0},
				NumSymbols: 37000, // Slightly more than the bin limit, but less than 2x
			},
		}

		now := time.Now()
		err := batchMeterer.BatchMeterRequest(ctx, requests, now)

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
		// Set up required method calls
		mockChainPayment.On("GetReservedPaymentByAccount", mock.Anything, mock.Anything).Return(map[core.QuorumID]*core.ReservedPayment{}, nil).Maybe()

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

		// Create sample requests
		requests := []BatchRequestInfo{
			{
				AccountID:  common.HexToAddress("0x1"),
				QuorumIDs:  []core.QuorumID{0, 1},
				NumSymbols: 100,
			},
		}

		now := time.Now()
		err := batchMeterer.BatchMeterRequest(ctx, requests, now)

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
		mockChainPayment.On("GetReservationWindow").Return(uint64(3600)).Maybe()
		// Set up required method calls
		mockChainPayment.On("GetReservedPaymentByAccount", mock.Anything, mock.Anything).Return(map[core.QuorumID]*core.ReservedPayment{}, nil).Maybe()

		// Set up the mock to return an error
		mockChainPayment.On("GetReservedPaymentByAccountAndQuorums", ctx, common.HexToAddress("0x1"), []core.QuorumID{0}).Return(map[core.QuorumID]*core.ReservedPayment{}, fmt.Errorf("chain error"))

		logger := testutils.GetLogger()

		config := meterer.Config{
			ChainReadTimeout: 1 * time.Second,
			UpdateInterval:   1 * time.Minute,
		}

		batchMeterer := NewBatchMeterer(config, mockChainPayment, 3, logger)

		// Create sample requests
		requests := []BatchRequestInfo{
			{
				AccountID:  common.HexToAddress("0x1"),
				QuorumIDs:  []core.QuorumID{0},
				NumSymbols: 100,
			},
		}

		now := time.Now()
		err := batchMeterer.BatchMeterRequest(ctx, requests, now)

		// Assert the results
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get reservations")

		// Verify expectations
		mockChainPayment.AssertExpectations(t)
	})

	t.Run("mixed accounts with one invalid", func(t *testing.T) {
		// Create a mock OnchainPayment
		mockChainPayment := new(coremock.MockOnchainPaymentState)

		// Set up the mock expectations
		mockChainPayment.On("GetMinNumSymbols").Return(uint64(10))
		mockChainPayment.On("GetReservationWindow").Return(uint64(3600))
		// Set up required method calls
		mockChainPayment.On("GetReservedPaymentByAccount", mock.Anything, mock.Anything).Return(map[core.QuorumID]*core.ReservedPayment{}, nil).Maybe()

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

		// Account 2 has expired reservation
		reservationsAcc2 := map[core.QuorumID]*core.ReservedPayment{
			1: {
				SymbolsPerSecond: 200,
				StartTimestamp:   0,
				EndTimestamp:     1000, // Expired
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

		// Create sample requests with mixed accounts
		requests := []BatchRequestInfo{
			{
				AccountID:  common.HexToAddress("0x1"),
				QuorumIDs:  []core.QuorumID{0, 1},
				NumSymbols: 100,
			},
			{
				AccountID:  common.HexToAddress("0x2"),
				QuorumIDs:  []core.QuorumID{1},
				NumSymbols: 200,
			},
		}

		// Current time is after account 2's reservation end time
		now := time.Unix(2000, 0)
		err := batchMeterer.BatchMeterRequest(ctx, requests, now)

		// Assert the results - the entire batch should fail due to account 2
		require.Error(t, err)
		assert.Contains(t, err.Error(), "inactive reservation")

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
		// Set up required method calls
		mockChainPayment.On("GetReservedPaymentByAccount", mock.Anything, mock.Anything).Return(map[core.QuorumID]*core.ReservedPayment{}, nil).Maybe()

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

		// Instead of hardcoding the quorum order, use a dynamic approach
		// Create a sample batch to analyze the order of quorums
		batch := createSampleBatch(t)
		quorumIDs1 := make([]core.QuorumID, 0)
		quorumIDs2 := make([]core.QuorumID, 0)

		// Extract the actual quorum orders from the blob certificates
		for _, cert := range batch.BlobCertificates {
			if cert.BlobHeader.PaymentMetadata.AccountID == common.HexToAddress("0x1") {
				quorumIDs1 = cert.BlobHeader.QuorumNumbers
			} else if cert.BlobHeader.PaymentMetadata.AccountID == common.HexToAddress("0x2") {
				quorumIDs2 = cert.BlobHeader.QuorumNumbers
			}
		}

		// Set up the mock for retrieving reservations with the actual quorum orders
		mockChainPayment.On("GetReservedPaymentByAccountAndQuorums", ctx, common.HexToAddress("0x1"), quorumIDs1).Return(reservationsAcc1, nil)
		mockChainPayment.On("GetReservedPaymentByAccountAndQuorums", ctx, common.HexToAddress("0x2"), quorumIDs2).Return(reservationsAcc2, nil)

		logger := testutils.GetLogger()

		config := meterer.Config{
			ChainReadTimeout: 1 * time.Second,
			UpdateInterval:   1 * time.Minute,
		}

		batchMeterer := NewBatchMeterer(config, mockChainPayment, 3, logger)

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
		// Set up required method calls
		mockChainPayment.On("GetReservedPaymentByAccount", mock.Anything, mock.Anything).Return(map[core.QuorumID]*core.ReservedPayment{}, nil).Maybe()

		// Create a sample batch to analyze the order of quorums
		batch := createSampleBatch(t)
		quorumIDs1 := make([]core.QuorumID, 0)
		quorumIDs2 := make([]core.QuorumID, 0)

		// Extract the actual quorum orders from the blob certificates
		for _, cert := range batch.BlobCertificates {
			if cert.BlobHeader.PaymentMetadata.AccountID == common.HexToAddress("0x1") {
				quorumIDs1 = cert.BlobHeader.QuorumNumbers
			} else if cert.BlobHeader.PaymentMetadata.AccountID == common.HexToAddress("0x2") {
				quorumIDs2 = cert.BlobHeader.QuorumNumbers
			}
		}

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

		// Set up the mock for retrieving reservations with the actual quorum orders
		mockChainPayment.On("GetReservedPaymentByAccountAndQuorums", ctx, common.HexToAddress("0x1"), quorumIDs1).Return(reservationsAcc1, nil)
		mockChainPayment.On("GetReservedPaymentByAccountAndQuorums", ctx, common.HexToAddress("0x2"), quorumIDs2).Return(reservationsAcc2, nil)

		logger := testutils.GetLogger()

		config := meterer.Config{
			ChainReadTimeout: 1 * time.Second,
			UpdateInterval:   1 * time.Minute,
		}

		batchMeterer := NewBatchMeterer(config, mockChainPayment, 3, logger)

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
		// Set up required method calls
		mockChainPayment.On("GetReservedPaymentByAccount", mock.Anything, mock.Anything).Return(map[core.QuorumID]*core.ReservedPayment{}, nil).Maybe()

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
