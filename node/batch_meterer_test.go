package node

import (
	"context"
	"math"
	"math/big"
	"sync"
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

// Helper to create a fresh BatchMeterer and mock for each test
func newTestBatchMeterer(t *testing.T) (*BatchMeterer, *coremock.MockOnchainPaymentState, context.Context) {
	ctx := context.Background()
	logger := testutils.GetLogger()
	mockState := new(coremock.MockOnchainPaymentState)
	mockState.On("GetMinNumSymbols").Return(uint64(32)).Maybe()
	mockState.On("GetReservationWindow").Return(uint64(3600)).Maybe()

	config := meterer.Config{
		ChainReadTimeout: 1 * time.Second,
		UpdateInterval:   1 * time.Minute,
	}

	batchMeterer := NewBatchMeterer(config, mockState, 3, logger)
	return batchMeterer, mockState, ctx
}

// TestBatchMeterRequest tests the BatchMeterRequest function
func TestBatchMeterRequest(t *testing.T) {
	t.Run("valid reservation", func(t *testing.T) {
		batchMeterer, mockState, ctx := newTestBatchMeterer(t)
		// Create a valid reservation
		reservations := map[core.QuorumID]*core.ReservedPayment{
			0: {
				SymbolsPerSecond: 100,
				StartTimestamp:   0,
				EndTimestamp:     9999999999,
			},
		}
		mockState.On("GetReservedPaymentByAccountAndQuorums",
			ctx,
			common.HexToAddress("0x1"),
			[]core.QuorumID{0},
		).Return(reservations, nil).Once()

		now := time.Now()
		currentPeriod := meterer.GetReservationPeriod(now.Unix(), 3600)
		usageMap := map[common.Address]map[core.QuorumID]map[uint64]uint64{
			common.HexToAddress("0x1"): {
				0: {currentPeriod: 100},
			},
		}

		err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)
		require.NoError(t, err)

		// Verify usage was recorded
		accountUsage := batchMeterer.getOrCreateAccountUsage(common.HexToAddress("0x1"))
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()

		relativeIndex := uint32((currentPeriod / 3600) % 3)
		assert.Equal(t, uint64(100), accountUsage.PeriodRecords[0][relativeIndex].Usage)
	})

	t.Run("reservation edge cases", func(t *testing.T) {
		t.Run("just started reservation", func(t *testing.T) {
			batchMeterer, mockState, ctx := newTestBatchMeterer(t)
			// Create a reservation that just started
			now := time.Now()
			reservations := map[core.QuorumID]*core.ReservedPayment{
				0: {
					SymbolsPerSecond: 100,
					StartTimestamp:   uint64(now.Unix()),
					EndTimestamp:     uint64(now.Add(24 * time.Hour).Unix()),
				},
			}
			mockState.On("GetReservedPaymentByAccountAndQuorums",
				ctx,
				common.HexToAddress("0x1"),
				[]core.QuorumID{0},
			).Return(reservations, nil).Once()

			currentPeriod := meterer.GetReservationPeriod(now.Unix(), 3600)
			usageMap := map[common.Address]map[core.QuorumID]map[uint64]uint64{
				common.HexToAddress("0x1"): {
					0: {currentPeriod: 100},
				},
			}

			err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)
			require.NoError(t, err)
		})

		t.Run("just expired reservation", func(t *testing.T) {
			batchMeterer, mockState, ctx := newTestBatchMeterer(t)
			// Create an expired reservation
			now := time.Now()
			reservations := map[core.QuorumID]*core.ReservedPayment{
				0: {
					SymbolsPerSecond: 100,
					StartTimestamp:   uint64(now.Add(-24 * time.Hour).Unix()),
					EndTimestamp:     uint64(now.Add(-1 * time.Second).Unix()),
				},
			}
			mockState.On("GetReservedPaymentByAccountAndQuorums",
				ctx,
				common.HexToAddress("0x1"),
				[]core.QuorumID{0},
			).Return(reservations, nil).Once()

			currentPeriod := meterer.GetReservationPeriod(now.Unix(), 3600)
			usageMap := map[common.Address]map[core.QuorumID]map[uint64]uint64{
				common.HexToAddress("0x1"): {
					0: {currentPeriod: 100},
				},
			}

			err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "inactive reservation")
		})

		t.Run("zero symbols per second", func(t *testing.T) {
			batchMeterer, mockState, ctx := newTestBatchMeterer(t)
			// Create a reservation with zero symbols per second
			now := time.Now()
			reservations := map[core.QuorumID]*core.ReservedPayment{
				0: {
					SymbolsPerSecond: 0,
					StartTimestamp:   uint64(now.Unix()),
					EndTimestamp:     uint64(now.Add(24 * time.Hour).Unix()),
				},
			}
			mockState.On("GetReservedPaymentByAccountAndQuorums",
				ctx,
				common.HexToAddress("0x1"),
				[]core.QuorumID{0},
			).Return(reservations, nil).Once()

			currentPeriod := meterer.GetReservationPeriod(now.Unix(), 3600)
			usageMap := map[common.Address]map[core.QuorumID]map[uint64]uint64{
				common.HexToAddress("0x1"): {
					0: {currentPeriod: 100},
				},
			}

			err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "failed usage validation for quorum")
		})

		t.Run("invalid timestamps", func(t *testing.T) {
			batchMeterer, mockState, ctx := newTestBatchMeterer(t)
			// Create a reservation with invalid timestamps
			now := time.Now()
			reservations := map[core.QuorumID]*core.ReservedPayment{
				0: {
					SymbolsPerSecond: 100,
					StartTimestamp:   uint64(now.Add(-24 * time.Hour).Unix()),
					EndTimestamp:     uint64(now.Unix()),
				},
			}
			mockState.On("GetReservedPaymentByAccountAndQuorums",
				ctx,
				common.HexToAddress("0x1"),
				[]core.QuorumID{0},
			).Return(reservations, nil).Once()

			currentPeriod := meterer.GetReservationPeriod(now.Unix(), 3600)
			usageMap := map[common.Address]map[core.QuorumID]map[uint64]uint64{
				common.HexToAddress("0x1"): {
					0: {currentPeriod: 100},
				},
			}

			err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid reservation period for quorum")
		})
	})

	t.Run("period record management", func(t *testing.T) {
		t.Run("period transition", func(t *testing.T) {
			batchMeterer, mockState, ctx := newTestBatchMeterer(t)
			// Create a reservation that spans multiple periods
			reservations := map[core.QuorumID]*core.ReservedPayment{
				0: {
					SymbolsPerSecond: 100,
					StartTimestamp:   0,
					EndTimestamp:     9999999999,
				},
			}
			mockState.On("GetReservedPaymentByAccountAndQuorums",
				ctx,
				common.HexToAddress("0x1"),
				[]core.QuorumID{0},
			).Return(reservations, nil)

			// First request in current period
			now := time.Now()
			currentPeriod := meterer.GetReservationPeriod(now.Unix(), 3600)
			usageMap := map[common.Address]map[core.QuorumID]map[uint64]uint64{
				common.HexToAddress("0x1"): {
					0: {currentPeriod: 100},
				},
			}
			err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)
			require.NoError(t, err)

			// Second request in next period
			nextPeriodTime := now.Add(time.Duration(3600) * time.Second)
			nextPeriod := meterer.GetReservationPeriod(nextPeriodTime.Unix(), 3600)
			usageMapNext := map[common.Address]map[core.QuorumID]map[uint64]uint64{
				common.HexToAddress("0x1"): {
					0: {nextPeriod: 100},
				},
			}
			err = batchMeterer.BatchMeterRequest(ctx, usageMapNext, nextPeriodTime)
			require.NoError(t, err)

			// Verify both periods have correct usage
			accountUsage := batchMeterer.getOrCreateAccountUsage(common.HexToAddress("0x1"))
			accountUsage.Lock.RLock()
			defer accountUsage.Lock.RUnlock()

			currentRelativeIndex := uint32((currentPeriod / 3600) % 3)
			nextRelativeIndex := uint32((nextPeriod / 3600) % 3)
			assert.Equal(t, uint64(100), accountUsage.PeriodRecords[0][currentRelativeIndex].Usage)
			assert.Equal(t, uint64(100), accountUsage.PeriodRecords[0][nextRelativeIndex].Usage)
		})

		t.Run("circular buffer wrapping", func(t *testing.T) {
			// Create a reservation
			batchMeterer, mockState, ctx := newTestBatchMeterer(t)
			reservations := map[core.QuorumID]*core.ReservedPayment{
				0: {
					SymbolsPerSecond: 100,
					StartTimestamp:   0,
					EndTimestamp:     9999999999,
				},
			}
			mockState.On("GetReservedPaymentByAccountAndQuorums",
				ctx,
				common.HexToAddress("0x1"),
				[]core.QuorumID{0},
			).Return(reservations, nil).Maybe()

			// Make requests across 4 periods to test buffer wrapping
			now := time.Now()
			for i := 0; i < 4; i++ {
				periodTime := now.Add(time.Duration(i*3600) * time.Second)
				periodIndex := meterer.GetReservationPeriod(periodTime.Unix(), 3600)
				usageMap := map[common.Address]map[core.QuorumID]map[uint64]uint64{
					common.HexToAddress("0x1"): {
						0: {periodIndex: 100},
					},
				}
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
	})

	t.Run("reservation at capacity", func(t *testing.T) {
		batchMeterer, mockState, ctx := newTestBatchMeterer(t)
		// Create a reservation at capacity
		reservations := map[core.QuorumID]*core.ReservedPayment{
			0: {
				SymbolsPerSecond: 1, // Only 1 symbol per second
				StartTimestamp:   0,
				EndTimestamp:     9999999999,
			},
		}
		mockState.On("GetReservedPaymentByAccountAndQuorums",
			ctx,
			common.HexToAddress("0x1"),
			[]core.QuorumID{0},
		).Return(reservations, nil).Once()

		// Setup a usage record for this account/quorum that already has some usage
		accountUsage := batchMeterer.getOrCreateAccountUsage(common.HexToAddress("0x1"))
		now := time.Now()
		currentPeriod := meterer.GetReservationPeriod(now.Unix(), 3600)
		accountUsage.Lock.Lock()
		accountUsage.PeriodRecords[0] = make([]*pb.PeriodRecord, 3)
		relativeIndex := uint32((currentPeriod / 3600) % 3)
		accountUsage.PeriodRecords[0][relativeIndex] = &pb.PeriodRecord{
			Index: uint32(currentPeriod),
			Usage: 3600, // Set usage to capacity (1 symbol per second for 1 hour)
		}
		// Initialize other elements to avoid nil pointer dereference
		for i := 0; i < 3; i++ {
			if i != int(relativeIndex) {
				accountUsage.PeriodRecords[0][i] = &pb.PeriodRecord{
					Index: 0,
					Usage: 0,
				}
			}
		}
		accountUsage.Lock.Unlock()

		usageMap := map[common.Address]map[core.QuorumID]map[uint64]uint64{
			common.HexToAddress("0x1"): {
				0: {currentPeriod: 1},
			},
		}

		err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed usage validation for quorum")
	})

	t.Run("multiple periods in same request", func(t *testing.T) {
		batchMeterer, mockState, ctx := newTestBatchMeterer(t)
		// Create a valid reservation
		reservations := map[core.QuorumID]*core.ReservedPayment{
			0: {
				SymbolsPerSecond: 100,
				StartTimestamp:   0,
				EndTimestamp:     9999999999,
			},
		}
		mockState.On("GetReservedPaymentByAccountAndQuorums",
			ctx,
			common.HexToAddress("0x1"),
			[]core.QuorumID{0},
		).Return(reservations, nil).Once()

		// Meter usage for current and previous period in a single call
		now := time.Now()
		currentPeriod := meterer.GetReservationPeriod(now.Unix(), 3600)
		previousPeriod := currentPeriod - 3600

		// Initialize period records to ensure clean state
		accountUsage := batchMeterer.getOrCreateAccountUsage(common.HexToAddress("0x1"))
		accountUsage.Lock.Lock()
		accountUsage.PeriodRecords[0] = make([]*pb.PeriodRecord, 3)
		for i := range accountUsage.PeriodRecords[0] {
			accountUsage.PeriodRecords[0][i] = &pb.PeriodRecord{
				Index: 0,
				Usage: 0,
			}
		}
		accountUsage.Lock.Unlock()

		usageMap := map[common.Address]map[core.QuorumID]map[uint64]uint64{
			common.HexToAddress("0x1"): {
				0: {
					currentPeriod:  50,
					previousPeriod: 50,
				},
			},
		}
		err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)
		require.NoError(t, err)

		// Verify usage was recorded for both periods
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()

		currentRelativeIndex := uint32((currentPeriod / 3600) % 3)
		previousRelativeIndex := uint32((previousPeriod / 3600) % 3)
		assert.Equal(t, uint64(50), accountUsage.PeriodRecords[0][currentRelativeIndex].Usage)
		assert.Equal(t, uint64(50), accountUsage.PeriodRecords[0][previousRelativeIndex].Usage)
	})

	t.Run("invalid period index", func(t *testing.T) {
		batchMeterer, mockState, ctx := newTestBatchMeterer(t)
		// Create a valid reservation
		reservations := map[core.QuorumID]*core.ReservedPayment{
			0: {
				SymbolsPerSecond: 100,
				StartTimestamp:   0,
				EndTimestamp:     9999999999,
			},
		}
		mockState.On("GetReservedPaymentByAccountAndQuorums",
			ctx,
			common.HexToAddress("0x1"),
			[]core.QuorumID{0},
		).Return(reservations, nil).Once()

		now := time.Now()
		currentPeriod := meterer.GetReservationPeriod(now.Unix(), 3600)
		invalidPeriod := currentPeriod - 3600*3 // Make it clearly invalid by going back 3 periods

		// Initialize the period records
		accountUsage := batchMeterer.getOrCreateAccountUsage(common.HexToAddress("0x1"))
		accountUsage.Lock.Lock()
		accountUsage.PeriodRecords[0] = make([]*pb.PeriodRecord, 3)
		for i := range accountUsage.PeriodRecords[0] {
			accountUsage.PeriodRecords[0][i] = &pb.PeriodRecord{
				Index: 0,
				Usage: 0,
			}
		}
		accountUsage.Lock.Unlock()

		usageMap := map[common.Address]map[core.QuorumID]map[uint64]uint64{
			common.HexToAddress("0x1"): {
				0: {invalidPeriod: 100},
			},
		}

		err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid reservation period for quorum")
	})

	t.Run("period index overflow", func(t *testing.T) {
		batchMeterer, mockState, ctx := newTestBatchMeterer(t)
		// Create a valid reservation
		reservations := map[core.QuorumID]*core.ReservedPayment{
			0: {
				SymbolsPerSecond: 100,
				StartTimestamp:   0,
				EndTimestamp:     9999999999,
			},
		}
		mockState.On("GetReservedPaymentByAccountAndQuorums",
			ctx,
			common.HexToAddress("0x1"),
			[]core.QuorumID{0},
		).Return(reservations, nil).Once()

		now := time.Now()
		currentPeriod := meterer.GetReservationPeriod(now.Unix(), 3600)
		farFuturePeriod := currentPeriod + 3600*4 // Beyond the 3-bin window
		usageMap := map[common.Address]map[core.QuorumID]map[uint64]uint64{
			common.HexToAddress("0x1"): {
				0: {farFuturePeriod: 100},
			},
		}

		err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid reservation period for quorum")
	})
}

// TestMeterBatch tests the MeterBatch function
func TestMeterBatch(t *testing.T) {
	ctx := context.Background()

	t.Run("successful batch metering", func(t *testing.T) {
		// Create a mock OnchainPayment
		mockChainPayment := new(coremock.MockOnchainPaymentState)

		// Set up the mock expectations
		mockChainPayment.On("GetMinNumSymbols").Return(uint64(10)).Maybe()
		mockChainPayment.On("GetReservationWindow").Return(uint64(3600)).Maybe()

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

		// Set up the mock for retrieving reservations with exact argument matching
		mockChainPayment.On("GetReservedPaymentByAccountAndQuorums",
			ctx,
			common.HexToAddress("0x1"),
			[]core.QuorumID{0, 1},
		).Return(reservationsAcc1, nil).Once()

		mockChainPayment.On("GetReservedPaymentByAccountAndQuorums",
			ctx,
			common.HexToAddress("0x2"),
			[]core.QuorumID{1},
		).Return(reservationsAcc2, nil).Once()

		logger := testutils.GetLogger()

		config := meterer.Config{
			ChainReadTimeout: 1 * time.Second,
			UpdateInterval:   1 * time.Minute,
		}

		batchMeterer := NewBatchMeterer(config, mockChainPayment, 3, logger)

		// Create a sample batch
		batch := createTestBatch([]accountQuorumInfo{
			{
				account:    common.HexToAddress("0x1"),
				quorumIDs:  []core.QuorumID{0, 1},
				numSymbols: 100,
			},
			{
				account:    common.HexToAddress("0x2"),
				quorumIDs:  []core.QuorumID{1},
				numSymbols: 200,
			},
		})

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
		mockChainPayment.On("GetMinNumSymbols").Return(uint64(10)).Maybe()
		mockChainPayment.On("GetReservationWindow").Return(uint64(3600)).Maybe()

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

		// Set up the mock for retrieving reservations with exact argument matching
		mockChainPayment.On("GetReservedPaymentByAccountAndQuorums",
			ctx,
			common.HexToAddress("0x1"),
			[]core.QuorumID{0, 1},
		).Return(reservationsAcc1, nil)

		mockChainPayment.On("GetReservedPaymentByAccountAndQuorums",
			ctx,
			common.HexToAddress("0x2"),
			[]core.QuorumID{1},
		).Return(reservationsAcc2, nil)

		logger := testutils.GetLogger()

		config := meterer.Config{
			ChainReadTimeout: 1 * time.Second,
			UpdateInterval:   1 * time.Minute,
		}

		batchMeterer := NewBatchMeterer(config, mockChainPayment, 3, logger)

		// Create a sample batch
		batch := createTestBatch([]accountQuorumInfo{
			{
				account:    common.HexToAddress("0x1"),
				quorumIDs:  []core.QuorumID{0, 1},
				numSymbols: 100,
			},
			{
				account:    common.HexToAddress("0x2"),
				quorumIDs:  []core.QuorumID{1},
				numSymbols: 200,
			},
		})

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
		currentPeriod := meterer.GetReservationPeriod(time.Now().Unix(), 3600)
		assert.Equal(t, uint64(32), usageMap[common.HexToAddress("0x1")][0][currentPeriod])
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
		currentPeriod := meterer.GetReservationPeriod(time.Now().Unix(), 3600)
		assert.Equal(t, uint64(64), usageMap[common.HexToAddress("0x1")][0][currentPeriod])
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
		currentPeriod := meterer.GetReservationPeriod(time.Now().Unix(), 3600)
		assert.Greater(t, usageMap[common.HexToAddress("0x1")][0][currentPeriod], uint64(0))
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
		mockState.On("GetMinNumSymbols").Return(uint64(32)).Maybe()
		mockState.On("GetReservationWindow").Return(uint64(3600)).Maybe()
		mockState.On("GetReservedPaymentByAccountAndQuorums",
			ctx,
			common.HexToAddress("0x1"),
			[]core.QuorumID{0},
		).Return(reservations, nil).Once()

		now := time.Now()
		currentPeriod := meterer.GetReservationPeriod(now.Unix(), 3600)
		usageMap := map[common.Address]map[core.QuorumID]map[uint64]uint64{
			common.HexToAddress("0x1"): {
				0: {currentPeriod: binLimit},
			},
		}

		err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)
		require.NoError(t, err)

		// Verify usage is exactly at bin limit
		accountUsage := batchMeterer.getOrCreateAccountUsage(common.HexToAddress("0x1"))
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()

		relativeIndex := uint32((currentPeriod / 3600) % 3)
		assert.Equal(t, binLimit, accountUsage.PeriodRecords[0][relativeIndex].Usage)
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
		mockState.On("GetMinNumSymbols").Return(uint64(32)).Maybe()
		mockState.On("GetReservationWindow").Return(uint64(3600)).Maybe()
		mockState.On("GetReservedPaymentByAccountAndQuorums",
			ctx,
			common.HexToAddress("0x1"),
			[]core.QuorumID{0},
		).Return(reservations, nil).Once()

		// Initialize the period record with some usage
		accountUsage := batchMeterer.getOrCreateAccountUsage(common.HexToAddress("0x1"))
		now := time.Now()
		currentPeriod := meterer.GetReservationPeriod(now.Unix(), 3600)
		accountUsage.Lock.Lock()
		reservationWindow := uint64(3600)
		numBins := uint32(3)
		relativeIndex := uint32((currentPeriod / reservationWindow) % uint64(numBins))
		accountUsage.PeriodRecords[0] = make([]*pb.PeriodRecord, numBins)
		for i := range accountUsage.PeriodRecords[0] {
			accountUsage.PeriodRecords[0][i] = &pb.PeriodRecord{
				Index: 0,
				Usage: 0,
			}
		}
		accountUsage.PeriodRecords[0][relativeIndex].Index = uint32(currentPeriod)
		accountUsage.PeriodRecords[0][relativeIndex].Usage = binLimit - 1
		accountUsage.Lock.Unlock()

		// Try to use just over bin limit
		usageMap := map[common.Address]map[core.QuorumID]map[uint64]uint64{
			common.HexToAddress("0x1"): {
				0: {currentPeriod: 2}, // Add 2 more to exceed limit
			},
		}

		err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)
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
		mockState.On("GetMinNumSymbols").Return(uint64(32)).Maybe()
		mockState.On("GetReservationWindow").Return(uint64(3600)).Maybe()
		mockState.On("GetReservedPaymentByAccountAndQuorums",
			ctx,
			common.HexToAddress("0x1"),
			[]core.QuorumID{0},
		).Return(reservations, nil).Once()

		// Initialize the period record with some usage
		accountUsage := batchMeterer.getOrCreateAccountUsage(common.HexToAddress("0x1"))
		now := time.Now()
		currentPeriod := meterer.GetReservationPeriod(now.Unix(), 3600)
		accountUsage.Lock.Lock()
		reservationWindow := uint64(3600)
		numBins := uint32(3)
		relativeIndex := uint32((currentPeriod / reservationWindow) % uint64(numBins))
		accountUsage.PeriodRecords[0] = make([]*pb.PeriodRecord, numBins)
		for i := range accountUsage.PeriodRecords[0] {
			accountUsage.PeriodRecords[0][i] = &pb.PeriodRecord{
				Index: 0,
				Usage: 0,
			}
		}
		accountUsage.PeriodRecords[0][relativeIndex].Index = uint32(currentPeriod)
		accountUsage.PeriodRecords[0][relativeIndex].Usage = 0
		accountUsage.Lock.Unlock()

		// Try to use exactly 2x bin limit
		usageMap := map[common.Address]map[core.QuorumID]map[uint64]uint64{
			common.HexToAddress("0x1"): {
				0: {currentPeriod: 2 * binLimit},
			},
		}

		err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)
		require.NoError(t, err)

		// Verify overflow was handled
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()

		assert.Equal(t, binLimit, accountUsage.PeriodRecords[0][relativeIndex].Usage)
		overflowPeriod := currentPeriod + reservationWindow
		overflowRelativeIndex := uint32((overflowPeriod / reservationWindow) % uint64(numBins))
		assert.Equal(t, binLimit, accountUsage.PeriodRecords[0][overflowRelativeIndex].Usage)
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
		mockState.On("GetMinNumSymbols").Return(uint64(32)).Maybe()
		mockState.On("GetReservationWindow").Return(uint64(3600)).Maybe()
		mockState.On("GetReservedPaymentByAccountAndQuorums",
			ctx,
			common.HexToAddress("0x1"),
			[]core.QuorumID{0},
		).Return(reservations, nil).Once()

		// Initialize the period record with some usage
		accountUsage := batchMeterer.getOrCreateAccountUsage(common.HexToAddress("0x1"))
		accountUsage.Lock.Lock()
		accountUsage.PeriodRecords[0] = make([]*pb.PeriodRecord, 3)
		for i := range accountUsage.PeriodRecords[0] {
			accountUsage.PeriodRecords[0][i] = &pb.PeriodRecord{
				Index: 0,
				Usage: 0,
			}
		}
		accountUsage.Lock.Unlock()

		// Try to use over 2x bin limit
		now := time.Now()
		currentPeriod := meterer.GetReservationPeriod(now.Unix(), 3600)
		usageMap := map[common.Address]map[core.QuorumID]map[uint64]uint64{
			common.HexToAddress("0x1"): {
				0: {currentPeriod: 2*binLimit + 1},
			},
		}

		err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)
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
		// Clear the batchMeterer cache to ensure a clean state
		batchMeterer.AccountUsages = sync.Map{}

		// Create a reservation that spans multiple periods
		reservations := map[core.QuorumID]*core.ReservedPayment{
			0: {
				SymbolsPerSecond: 100,
				StartTimestamp:   0,
				EndTimestamp:     9999999999,
			},
		}
		mockState.On("GetMinNumSymbols").Return(uint64(32)).Maybe()
		mockState.On("GetReservationWindow").Return(uint64(3600)).Maybe()
		mockState.On("GetReservedPaymentByAccountAndQuorums",
			ctx,
			common.HexToAddress("0x1"),
			[]core.QuorumID{0},
		).Return(reservations, nil)

		// First request in current period
		now := time.Now()
		currentPeriod := meterer.GetReservationPeriod(now.Unix(), 3600)
		usageMap := map[common.Address]map[core.QuorumID]map[uint64]uint64{
			common.HexToAddress("0x1"): {
				0: {currentPeriod: 100},
			},
		}
		err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)
		require.NoError(t, err)

		// Second request in next period
		nextPeriodTime := now.Add(time.Duration(3600) * time.Second)
		nextPeriod := meterer.GetReservationPeriod(nextPeriodTime.Unix(), 3600)
		usageMapNext := map[common.Address]map[core.QuorumID]map[uint64]uint64{
			common.HexToAddress("0x1"): {
				0: {nextPeriod: 100},
			},
		}
		err = batchMeterer.BatchMeterRequest(ctx, usageMapNext, nextPeriodTime)
		require.NoError(t, err)

		// Verify both periods have correct usage
		accountUsage := batchMeterer.getOrCreateAccountUsage(common.HexToAddress("0x1"))
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()

		currentRelativeIndex := uint32((currentPeriod / 3600) % 3)
		nextRelativeIndex := uint32((nextPeriod / 3600) % 3)
		assert.Equal(t, uint64(100), accountUsage.PeriodRecords[0][currentRelativeIndex].Usage)
		assert.Equal(t, uint64(100), accountUsage.PeriodRecords[0][nextRelativeIndex].Usage)
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
		mockState.On("GetReservedPaymentByAccountAndQuorums",
			ctx,
			common.HexToAddress("0x1"),
			[]core.QuorumID{0},
		).Return(reservations, nil).Maybe()

		// Make requests across 4 periods to test buffer wrapping
		now := time.Now()
		for i := 0; i < 4; i++ {
			periodTime := now.Add(time.Duration(i*3600) * time.Second)
			periodIndex := meterer.GetReservationPeriod(periodTime.Unix(), 3600)
			usageMap := map[common.Address]map[core.QuorumID]map[uint64]uint64{
				common.HexToAddress("0x1"): {
					0: {periodIndex: 100},
				},
			}
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
		mockState.On("GetMinNumSymbols").Return(uint64(32)).Maybe()
		mockState.On("GetReservationWindow").Return(uint64(3600)).Maybe()
		mockState.On("GetReservedPaymentByAccountAndQuorums",
			ctx,
			common.HexToAddress("0x1"),
			[]core.QuorumID{0},
		).Return(reservations, nil).Once()

		currentPeriod := meterer.GetReservationPeriod(now.Unix(), 3600)
		usageMap := map[common.Address]map[core.QuorumID]map[uint64]uint64{
			common.HexToAddress("0x1"): {
				0: {currentPeriod: 100},
			},
		}

		err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)
		require.NoError(t, err)
	})

	t.Run("zero symbols per second", func(t *testing.T) {
		// Clear existing expectations
		mockState.ExpectedCalls = nil
		mockState.Calls = nil

		// Create a reservation with zero symbols per second
		now := time.Now()
		reservations := map[core.QuorumID]*core.ReservedPayment{
			0: {
				SymbolsPerSecond: 0,
				StartTimestamp:   uint64(now.Unix()),
				EndTimestamp:     uint64(now.Add(24 * time.Hour).Unix()),
			},
		}
		mockState.On("GetMinNumSymbols").Return(uint64(32)).Maybe()
		mockState.On("GetReservationWindow").Return(uint64(3600)).Maybe()
		mockState.On("GetReservedPaymentByAccountAndQuorums",
			ctx,
			common.HexToAddress("0x1"),
			[]core.QuorumID{0},
		).Return(reservations, nil).Once()

		currentPeriod := meterer.GetReservationPeriod(now.Unix(), 3600)
		usageMap := map[common.Address]map[core.QuorumID]map[uint64]uint64{
			common.HexToAddress("0x1"): {
				0: {currentPeriod: 100},
			},
		}

		err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed usage validation for quorum")
	})

	t.Run("invalid timestamps", func(t *testing.T) {
		// Create a reservation with invalid timestamps
		now := time.Now()
		reservations := map[core.QuorumID]*core.ReservedPayment{
			0: {
				SymbolsPerSecond: 100,
				StartTimestamp:   uint64(now.Add(-24 * time.Hour).Unix()),
				EndTimestamp:     uint64(now.Unix()), // End time in the past
			},
		}
		mockState.On("GetMinNumSymbols").Return(uint64(32)).Maybe()
		mockState.On("GetReservationWindow").Return(uint64(3600)).Maybe()
		mockState.On("GetReservedPaymentByAccountAndQuorums",
			ctx,
			common.HexToAddress("0x1"),
			[]core.QuorumID{0},
		).Return(reservations, nil).Once()

		currentPeriod := meterer.GetReservationPeriod(now.Unix(), 3600)
		usageMap := map[common.Address]map[core.QuorumID]map[uint64]uint64{
			common.HexToAddress("0x1"): {
				0: {currentPeriod: 100},
			},
		}

		err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid reservation period for quorum")
	})
}

// TestReservationEdgeCases tests edge cases in reservation validation
func TestReservationEdgeCases(t *testing.T) {
	ctx := context.Background()
	logger := testutils.GetLogger()

	config := meterer.Config{
		ChainReadTimeout: 1 * time.Second,
		UpdateInterval:   1 * time.Minute,
	}

	t.Run("just started reservation", func(t *testing.T) {
		mockState := new(coremock.MockOnchainPaymentState)
		mockState.On("GetMinNumSymbols").Return(uint64(32)).Maybe()
		mockState.On("GetReservationWindow").Return(uint64(3600)).Maybe()

		batchMeterer := NewBatchMeterer(config, mockState, 3, logger)

		// Create a reservation that just started
		now := time.Now()
		reservations := map[core.QuorumID]*core.ReservedPayment{
			0: {
				SymbolsPerSecond: 100,
				StartTimestamp:   uint64(now.Unix()),
				EndTimestamp:     uint64(now.Add(24 * time.Hour).Unix()),
			},
		}
		mockState.On("GetReservedPaymentByAccountAndQuorums", ctx, common.HexToAddress("0x1"), []core.QuorumID{0}).Return(reservations, nil).Once()

		currentPeriod := meterer.GetReservationPeriod(now.Unix(), 3600)
		usageMap := map[common.Address]map[core.QuorumID]map[uint64]uint64{
			common.HexToAddress("0x1"): {
				0: {currentPeriod: 100},
			},
		}

		err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)
		require.NoError(t, err)
	})

	t.Run("just expired reservation", func(t *testing.T) {
		mockState := new(coremock.MockOnchainPaymentState)
		mockState.On("GetMinNumSymbols").Return(uint64(32)).Maybe()
		mockState.On("GetReservationWindow").Return(uint64(3600)).Maybe()

		batchMeterer := NewBatchMeterer(config, mockState, 3, logger)

		// Create an expired reservation
		now := time.Now()
		reservations := map[core.QuorumID]*core.ReservedPayment{
			0: {
				SymbolsPerSecond: 100,
				StartTimestamp:   uint64(now.Add(-24 * time.Hour).Unix()),
				EndTimestamp:     uint64(now.Add(-1 * time.Second).Unix()), // End time is 1 second before now
			},
		}

		mockState.On("GetReservedPaymentByAccountAndQuorums",
			ctx,
			common.HexToAddress("0x1"),
			[]core.QuorumID{0},
		).Return(reservations, nil).Once()

		currentPeriod := meterer.GetReservationPeriod(now.Unix(), 3600)
		usageMap := map[common.Address]map[core.QuorumID]map[uint64]uint64{
			common.HexToAddress("0x1"): {
				0: {currentPeriod: 100},
			},
		}

		err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "inactive reservation")
	})

	t.Run("zero symbols per second", func(t *testing.T) {
		mockState := new(coremock.MockOnchainPaymentState)
		mockState.On("GetMinNumSymbols").Return(uint64(32)).Maybe()
		mockState.On("GetReservationWindow").Return(uint64(3600)).Maybe()

		batchMeterer := NewBatchMeterer(config, mockState, 3, logger)

		// Create a reservation with zero symbols per second
		now := time.Now()
		reservations := map[core.QuorumID]*core.ReservedPayment{
			0: {
				SymbolsPerSecond: 0,
				StartTimestamp:   uint64(now.Unix()),
				EndTimestamp:     uint64(now.Add(24 * time.Hour).Unix()),
			},
		}
		mockState.On("GetReservedPaymentByAccountAndQuorums",
			ctx,
			common.HexToAddress("0x1"),
			[]core.QuorumID{0},
		).Return(reservations, nil).Once()

		currentPeriod := meterer.GetReservationPeriod(now.Unix(), 3600)
		usageMap := map[common.Address]map[core.QuorumID]map[uint64]uint64{
			common.HexToAddress("0x1"): {
				0: {currentPeriod: 100},
			},
		}

		err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed usage validation for quorum")
	})

	t.Run("invalid timestamps", func(t *testing.T) {
		mockState := new(coremock.MockOnchainPaymentState)
		mockState.On("GetMinNumSymbols").Return(uint64(32)).Maybe()
		mockState.On("GetReservationWindow").Return(uint64(3600)).Maybe()

		batchMeterer := NewBatchMeterer(config, mockState, 3, logger)

		// Create a reservation with invalid timestamps
		now := time.Now()
		reservations := map[core.QuorumID]*core.ReservedPayment{
			0: {
				SymbolsPerSecond: 100,
				StartTimestamp:   uint64(now.Add(-24 * time.Hour).Unix()),
				EndTimestamp:     uint64(now.Unix()), // End time in the past
			},
		}
		mockState.On("GetReservedPaymentByAccountAndQuorums",
			ctx,
			common.HexToAddress("0x1"),
			[]core.QuorumID{0},
		).Return(reservations, nil).Once()

		currentPeriod := meterer.GetReservationPeriod(now.Unix(), 3600)
		usageMap := map[common.Address]map[core.QuorumID]map[uint64]uint64{
			common.HexToAddress("0x1"): {
				0: {currentPeriod: 100},
			},
		}

		err := batchMeterer.BatchMeterRequest(ctx, usageMap, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid reservation period for quorum")
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
