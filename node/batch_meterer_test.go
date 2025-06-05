package node

import (
	"context"
	"errors"
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
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBatchMetererConfig holds configuration for the test batch meterer
type TestBatchMetererConfig struct {
	MinNumSymbols     uint64
	ReservationWindow uint64
	NumBins           uint32
	ChainReadTimeout  time.Duration
	UpdateInterval    time.Duration
}

// DefaultTestConfig returns a default test configuration
func DefaultTestConfig() TestBatchMetererConfig {
	return TestBatchMetererConfig{
		MinNumSymbols:     32,
		ReservationWindow: 3600,
		NumBins:           3,
		ChainReadTimeout:  1 * time.Second,
		UpdateInterval:    1 * time.Minute,
	}
}

// TestReservation holds test reservation data
type TestReservation struct {
	SymbolsPerSecond uint64
	StartTime        time.Time
	EndTime          time.Time
}

// CreateTestReservation creates a test reservation with the given parameters
func CreateTestReservation(symbolsPerSecond uint64, startTime, endTime time.Time) *core.ReservedPayment {
	return &core.ReservedPayment{
		SymbolsPerSecond: symbolsPerSecond,
		StartTimestamp:   uint64(startTime.Unix()),
		EndTimestamp:     uint64(endTime.Unix()),
	}
}

// CreateTestUpdateRecord creates a test update record
func CreateTestUpdateRecord(accountID common.Address, quorumID core.QuorumID, period uint64, usage uint64) updateRecord {
	return updateRecord{
		accountID: accountID,
		quorumID:  quorumID,
		period:    period,
		usage:     usage,
	}
}

// InitializePeriodRecord initializes a period record for the given account and quorum
func InitializePeriodRecord(batchMeterer *BatchMeterer, accountID common.Address, quorumID core.QuorumID, period uint64) {
	accountUsage := batchMeterer.getOrCreateAccountUsage(accountID)
	accountUsage.Lock.Lock()
	defer accountUsage.Lock.Unlock()

	if accountUsage.PeriodRecords[quorumID] == nil {
		accountUsage.PeriodRecords[quorumID] = make([]*pb.PeriodRecord, 3)
	}
	relativeIndex := uint32((period / 3600) % 3)
	accountUsage.PeriodRecords[quorumID][relativeIndex] = &pb.PeriodRecord{
		Index: uint32(period),
		Usage: 0,
	}
}

// accountQuorumInfo holds information about an account's quorum usage
type accountQuorumInfo struct {
	account           common.Address
	quorumID          core.QuorumID
	timestamp         uint64
	cumulativePayment uint64
}

// createTestBatch creates a test batch with the given account quorum information
func createTestBatch(accountQuorumInfo []accountQuorumInfo) *corev2.Batch {
	batch := &corev2.Batch{
		BatchHeader: &corev2.BatchHeader{
			BatchRoot:            [32]byte{1, 1, 1},
			ReferenceBlockNumber: 100,
		},
		BlobCertificates: make([]*corev2.BlobCertificate, len(accountQuorumInfo)),
	}

	for i, info := range accountQuorumInfo {
		batch.BlobCertificates[i] = &corev2.BlobCertificate{
			BlobHeader: &corev2.BlobHeader{
				BlobVersion:     0,
				BlobCommitments: encoding.BlobCommitments{},
				QuorumNumbers:   []core.QuorumID{info.quorumID},
				PaymentMetadata: core.PaymentMetadata{
					AccountID:         info.account,
					Timestamp:         int64(info.timestamp),
					CumulativePayment: big.NewInt(int64(info.cumulativePayment)),
				},
			},
			Signature: []byte{1, 2, 3},
			RelayKeys: []corev2.RelayKey{0, 1},
		}
	}

	return batch
}

// testFixtures contains common test data and configurations
type testFixtures struct {
	ctx               context.Context
	logger            logging.Logger
	mockState         *coremock.MockOnchainPaymentState
	config            TestBatchMetererConfig
	batchMeterer      *BatchMeterer
	account1          common.Address
	account2          common.Address
	quorum0           core.QuorumID
	quorum1           core.QuorumID
	reservationWindow uint64
}

// setupTestFixtures creates and returns a new testFixtures instance
func setupTestFixtures(t *testing.T) *testFixtures {
	ctx := context.Background()
	logger := testutils.GetLogger()
	mockState := new(coremock.MockOnchainPaymentState)
	config := DefaultTestConfig()

	mockState.On("GetMinNumSymbols").Return(config.MinNumSymbols)
	mockState.On("GetReservationWindow").Return(config.ReservationWindow)

	batchMeterer := NewBatchMeterer(meterer.Config{
		ChainReadTimeout: config.ChainReadTimeout,
		UpdateInterval:   config.UpdateInterval,
	}, mockState, config.NumBins, logger)

	return &testFixtures{
		ctx:               ctx,
		logger:            logger,
		mockState:         mockState,
		config:            config,
		batchMeterer:      batchMeterer,
		account1:          common.HexToAddress("0x1"),
		account2:          common.HexToAddress("0x2"),
		quorum0:           0,
		quorum1:           1,
		reservationWindow: config.ReservationWindow,
	}
}

// resetBatchMeterer clears the batchMeterer's cache
func (f *testFixtures) resetBatchMeterer() {
	f.batchMeterer.AccountUsages = sync.Map{}
}

// createTestReservationMap creates a map of reservations for testing
func (f *testFixtures) createTestReservationMap(symbolsPerSecond uint64, startTime, endTime time.Time, quorumIDs ...core.QuorumID) map[core.QuorumID]*core.ReservedPayment {
	reservations := make(map[core.QuorumID]*core.ReservedPayment)
	for _, quorumID := range quorumIDs {
		reservations[quorumID] = CreateTestReservation(symbolsPerSecond, startTime, endTime)
	}
	return reservations
}

// setupPeriodRecord initializes a period record for testing
func (f *testFixtures) setupPeriodRecord(accountID common.Address, quorumID core.QuorumID, period uint64) {
	InitializePeriodRecord(f.batchMeterer, accountID, quorumID, period)
}

// getRelativeIndex calculates the relative index for a period
func (f *testFixtures) getRelativeIndex(period uint64) uint32 {
	return uint32((period / f.reservationWindow) % uint64(f.config.NumBins))
}

// verifyUsage verifies the usage for a given account, quorum, and period
func (f *testFixtures) verifyUsage(t *testing.T, accountID common.Address, quorumID core.QuorumID, period uint64, expectedUsage uint64) {
	accountUsage := f.batchMeterer.getOrCreateAccountUsage(accountID)
	accountUsage.Lock.RLock()
	defer accountUsage.Lock.RUnlock()

	relativeIndex := f.getRelativeIndex(period)
	require.NotNil(t, accountUsage.PeriodRecords[quorumID], "Period records for quorum %d should exist", quorumID)
	require.NotNil(t, accountUsage.PeriodRecords[quorumID][relativeIndex], "Period record should exist")
	assert.Equal(t, expectedUsage, accountUsage.PeriodRecords[quorumID][relativeIndex].Usage)
}

// TestBatchMeterRequest tests the BatchMeterRequest function
func TestBatchMeterRequest(t *testing.T) {
	f := setupTestFixtures(t)

	t.Run("valid reservation", func(t *testing.T) {
		f.resetBatchMeterer()

		// Create a valid reservation
		now := time.Now()
		reservations := f.createTestReservationMap(100, time.Unix(0, 0), time.Unix(9999999999, 0), f.quorum0)

		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil).Once()

		currentPeriod := meterer.GetReservationPeriod(now.Unix(), f.reservationWindow)
		f.setupPeriodRecord(f.account1, f.quorum0, currentPeriod)

		updates := []updateRecord{
			CreateTestUpdateRecord(f.account1, f.quorum0, currentPeriod, 100),
		}

		err := f.batchMeterer.BatchMeterRequest(f.ctx, updates, now)
		require.NoError(t, err)

		f.verifyUsage(t, f.account1, f.quorum0, currentPeriod, 100)
	})

	t.Run("period transition", func(t *testing.T) {
		f.resetBatchMeterer()

		// Create a valid reservation
		reservations := f.createTestReservationMap(100, time.Unix(0, 0), time.Unix(9999999999, 0), f.quorum0)

		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil).Times(2)

		// Create usage map for current period
		now := time.Now()
		currentPeriod := meterer.GetReservationPeriod(now.Unix(), f.reservationWindow)

		// Initialize period records
		f.setupPeriodRecord(f.account1, f.quorum0, currentPeriod)

		// Create updates array for current period
		updates := []updateRecord{
			{
				accountID: f.account1,
				quorumID:  f.quorum0,
				period:    currentPeriod,
				usage:     100,
			},
		}

		// Process the request
		err := f.batchMeterer.BatchMeterRequest(f.ctx, updates, now)
		require.NoError(t, err)

		// Create usage map for next period
		nextNow := now.Add(time.Hour)
		nextPeriod := meterer.GetReservationPeriod(nextNow.Unix(), f.reservationWindow)

		// Initialize period records for next period
		f.setupPeriodRecord(f.account1, f.quorum0, nextPeriod)

		// Create updates array for next period
		updatesNext := []updateRecord{
			{
				accountID: f.account1,
				quorumID:  f.quorum0,
				period:    nextPeriod, // Use next period
				usage:     100,
			},
		}

		// Process the request for next period
		err = f.batchMeterer.BatchMeterRequest(f.ctx, updatesNext, nextNow)
		require.NoError(t, err)

		// Verify both periods were updated correctly
		f.verifyUsage(t, f.account1, f.quorum0, currentPeriod, 100)
		f.verifyUsage(t, f.account1, f.quorum0, nextPeriod, 100)
	})

	t.Run("multiple periods in same request", func(t *testing.T) {
		f.resetBatchMeterer()

		// Create a valid reservation
		reservations := f.createTestReservationMap(100, time.Unix(0, 0), time.Unix(9999999999, 0), f.quorum0)

		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil).Once() // Only need one call since we process both periods together

		// Create usage map
		now := time.Now()
		currentPeriod := meterer.GetReservationPeriod(now.Unix(), f.reservationWindow)
		prevPeriod := currentPeriod - 3600 // Previous period

		// Initialize period records
		f.setupPeriodRecord(f.account1, f.quorum0, currentPeriod)
		f.setupPeriodRecord(f.account1, f.quorum0, prevPeriod)

		// Create updates array with both periods
		updates := []updateRecord{
			{
				accountID: f.account1,
				quorumID:  f.quorum0,
				period:    currentPeriod,
				usage:     50,
			},
			{
				accountID: f.account1,
				quorumID:  f.quorum0,
				period:    prevPeriod,
				usage:     50,
			},
		}

		// Process both periods in a single request
		err := f.batchMeterer.BatchMeterRequest(f.ctx, updates, now)
		require.NoError(t, err)

		// Verify both periods were updated correctly
		f.verifyUsage(t, f.account1, f.quorum0, currentPeriod, 50)
		f.verifyUsage(t, f.account1, f.quorum0, prevPeriod, 50)
	})

	t.Run("circular buffer wrapping", func(t *testing.T) {
		// Clear the batchMeterer cache to ensure a clean state
		f.resetBatchMeterer()

		// Create a reservation
		reservations := f.createTestReservationMap(100, time.Unix(0, 0), time.Unix(9999999999, 0), f.quorum0)

		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil).Maybe()

		// Make requests across 4 periods to test buffer wrapping
		now := time.Now()
		for i := 0; i < 4; i++ {
			periodTime := now.Add(time.Duration(i*3600) * time.Second)
			periodIndex := meterer.GetReservationPeriod(periodTime.Unix(), f.reservationWindow)
			updates := []updateRecord{
				{
					accountID: f.account1,
					quorumID:  f.quorum0,
					period:    periodIndex,
					usage:     100,
				},
			}
			err := f.batchMeterer.BatchMeterRequest(f.ctx, updates, periodTime)
			require.NoError(t, err)
		}

		// Verify the oldest record was overwritten
		accountUsage := f.batchMeterer.getOrCreateAccountUsage(f.account1)
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()

		// Should only have 3 records (buffer size)
		records := accountUsage.PeriodRecords[f.quorum0]
		assert.Equal(t, 3, len(records))
	})

	t.Run("just started reservation", func(t *testing.T) {
		// Clear existing expectations
		f.mockState.ExpectedCalls = nil
		f.mockState.Calls = nil

		// Create a reservation that just started
		now := time.Now()
		reservations := f.createTestReservationMap(100, now, now.Add(24*time.Hour), f.quorum0)

		f.mockState.On("GetMinNumSymbols").Return(uint64(32)).Maybe()
		f.mockState.On("GetReservationWindow").Return(uint64(3600)).Maybe()
		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil).Once()

		currentPeriod := meterer.GetReservationPeriod(now.Unix(), f.reservationWindow)
		f.setupPeriodRecord(f.account1, f.quorum0, currentPeriod)
		updates := []updateRecord{
			{
				accountID: f.account1,
				quorumID:  f.quorum0,
				period:    currentPeriod,
				usage:     100,
			},
		}

		err := f.batchMeterer.BatchMeterRequest(f.ctx, updates, now)
		require.NoError(t, err)
	})

	t.Run("zero symbols per second", func(t *testing.T) {
		// Clear existing expectations
		f.mockState.ExpectedCalls = nil
		f.mockState.Calls = nil

		// Create a reservation with zero symbols per second
		now := time.Now()
		reservations := f.createTestReservationMap(0, now, now.Add(24*time.Hour), f.quorum0)

		f.mockState.On("GetMinNumSymbols").Return(uint64(32)).Maybe()
		f.mockState.On("GetReservationWindow").Return(uint64(3600)).Maybe()
		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil).Once()

		currentPeriod := meterer.GetReservationPeriod(now.Unix(), f.reservationWindow)
		f.setupPeriodRecord(f.account1, f.quorum0, currentPeriod)
		updates := []updateRecord{
			{
				accountID: f.account1,
				quorumID:  f.quorum0,
				period:    currentPeriod,
				usage:     100,
			},
		}

		err := f.batchMeterer.BatchMeterRequest(f.ctx, updates, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed usage validation for quorum")
	})

	t.Run("invalid timestamps", func(t *testing.T) {
		// Create a reservation with invalid timestamps
		now := time.Now()
		reservations := f.createTestReservationMap(100, now.Add(-24*time.Hour), now, f.quorum0)

		f.mockState.On("GetMinNumSymbols").Return(uint64(32)).Maybe()
		f.mockState.On("GetReservationWindow").Return(uint64(3600)).Maybe()
		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil).Twice()

		currentPeriod := meterer.GetReservationPeriod(now.Unix(), f.reservationWindow)
		f.setupPeriodRecord(f.account1, f.quorum0, currentPeriod)
		updates := []updateRecord{
			{
				accountID: f.account1,
				quorumID:  f.quorum0,
				period:    currentPeriod,
				usage:     100,
			},
		}

		err := f.batchMeterer.BatchMeterRequest(f.ctx, updates, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid reservation period for quorum")
	})
}

// TestBatchMeterMeterBatch tests the MeterBatch function
func TestBatchMeterMeterBatch(t *testing.T) {
	f := setupTestFixtures(t)

	t.Run("successful batch metering", func(t *testing.T) {
		// Create reservations for accounts
		reservationsAcc1 := f.createTestReservationMap(100, time.Unix(0, 0), time.Unix(9999999999, 0), f.quorum0, f.quorum1)
		reservationsAcc2 := f.createTestReservationMap(200, time.Unix(0, 0), time.Unix(9999999999, 0), f.quorum1)

		// Set up the mock for retrieving reservations
		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0, f.quorum1},
		).Return(reservationsAcc1, nil)
		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservationsAcc1, nil)

		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account2,
			[]core.QuorumID{f.quorum1},
		).Return(reservationsAcc2, nil)

		// Create a sample batch with valid timestamps and blob lengths
		now := time.Now()
		timestampNano := uint64(now.UnixNano())
		periodFromTimestamp := meterer.GetReservationPeriodByNanosecond(int64(timestampNano), 3600)

		// Initialize period records for both accounts using periodFromTimestamp
		f.setupPeriodRecord(f.account1, f.quorum0, periodFromTimestamp)
		f.setupPeriodRecord(f.account1, f.quorum1, periodFromTimestamp)
		f.setupPeriodRecord(f.account2, f.quorum1, periodFromTimestamp)

		batch := createTestBatch([]accountQuorumInfo{
			{
				account:           f.account1,
				quorumID:          f.quorum0,
				timestamp:         timestampNano,
				cumulativePayment: 100,
			},
			{
				account:           f.account2,
				quorumID:          f.quorum1,
				timestamp:         timestampNano,
				cumulativePayment: 200,
			},
		})

		// Set blob lengths to ensure they're counted
		for _, cert := range batch.BlobCertificates {
			cert.BlobHeader.BlobCommitments.Length = 32
		}

		// Call the function
		err := f.batchMeterer.MeterBatch(f.ctx, batch, now)
		require.NoError(t, err)

		// Verify usage was updated for both accounts
		f.verifyUsage(t, f.account1, f.quorum0, periodFromTimestamp, 32)
		f.verifyUsage(t, f.account2, f.quorum1, periodFromTimestamp, 32)
	})

	t.Run("batch with invalid account", func(t *testing.T) {
		// Clear existing expectations
		f.mockState.ExpectedCalls = nil
		f.mockState.Calls = nil

		// Set up basic mock expectations
		f.mockState.On("GetMinNumSymbols").Return(uint64(32)).Maybe()
		f.mockState.On("GetReservationWindow").Return(uint64(3600)).Maybe()

		// Create valid reservations for account 1
		reservationsAcc1 := f.createTestReservationMap(100, time.Unix(0, 0), time.Unix(9999999999, 0), f.quorum0)

		// Set up the mock for retrieving reservations
		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservationsAcc1, nil)

		// Account 2 has no reservation - return empty map
		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account2,
			[]core.QuorumID{f.quorum1},
		).Return(map[core.QuorumID]*core.ReservedPayment{}, nil)

		// Initialize period records for both accounts
		now := time.Now()
		currentPeriod := meterer.GetReservationPeriod(now.Unix(), f.reservationWindow)
		f.setupPeriodRecord(f.account1, f.quorum0, currentPeriod)
		f.setupPeriodRecord(f.account2, f.quorum1, currentPeriod)

		batch := createTestBatch([]accountQuorumInfo{
			{
				account:           f.account1,
				quorumID:          f.quorum0,
				timestamp:         uint64(now.UnixNano()),
				cumulativePayment: 100,
			},
			{
				account:           f.account2,
				quorumID:          f.quorum1,
				timestamp:         uint64(now.UnixNano()),
				cumulativePayment: 200,
			},
		})

		// Set blob lengths in the batch
		for _, cert := range batch.BlobCertificates {
			cert.BlobHeader.BlobCommitments.Length = 32 // Set a valid length
		}

		// Call the function
		err := f.batchMeterer.MeterBatch(f.ctx, batch, now)

		// Assert the results - should fail due to account 2 having no reservation
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no reservation for quorum")
	})

	t.Run("nil batch", func(t *testing.T) {
		// Call with nil batch
		now := time.Now()
		err := f.batchMeterer.MeterBatch(f.ctx, nil, now)

		// Assert the results
		require.Error(t, err)
		assert.Contains(t, err.Error(), "batch is nil")
	})
}

// TestBatchRequestUsageEdgeCases tests edge cases in batch request usage calculation
func TestBatchMeterRequestUsageEdgeCases(t *testing.T) {
	f := setupTestFixtures(t)

	t.Run("empty batch", func(t *testing.T) {
		batch := &corev2.Batch{
			BatchHeader:      &corev2.BatchHeader{},
			BlobCertificates: []*corev2.BlobCertificate{},
		}
		_, err := f.batchMeterer.BatchRequestUsage(batch)
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
		_, err := f.batchMeterer.BatchRequestUsage(batch)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "blob certificate has nil header")
	})

	t.Run("zero length blob", func(t *testing.T) {
		batch := createTestBatch([]accountQuorumInfo{
			{
				account:           f.account1,
				quorumID:          f.quorum0,
				timestamp:         0,
				cumulativePayment: 0,
			},
		})
		// Set blob length to 0 explicitly
		batch.BlobCertificates[0].BlobHeader.BlobCommitments.Length = 0
		updates, err := f.batchMeterer.BatchRequestUsage(batch)
		require.NoError(t, err)
		currentPeriod := meterer.GetReservationPeriodByNanosecond(0, 3600)
		var found bool
		for _, update := range updates {
			if update.accountID == f.account1 && update.quorumID == f.quorum0 && update.period == currentPeriod {
				assert.Equal(t, uint64(32), update.usage) // MinNumSymbols is 32
				found = true
			}
		}
		assert.True(t, found, "expected updateRecord not found")
	})

	t.Run("duplicate account quorum", func(t *testing.T) {
		batch := createTestBatch([]accountQuorumInfo{
			{
				account:           f.account1,
				quorumID:          f.quorum0,
				timestamp:         0,
				cumulativePayment: 32,
			},
			{
				account:           f.account1,
				quorumID:          f.quorum0,
				timestamp:         0,
				cumulativePayment: 32,
			},
		})
		// Set blob lengths to ensure they're counted
		for _, cert := range batch.BlobCertificates {
			cert.BlobHeader.BlobCommitments.Length = 32
		}
		updates, err := f.batchMeterer.BatchRequestUsage(batch)
		require.NoError(t, err)
		currentPeriod := meterer.GetReservationPeriodByNanosecond(0, 3600)
		var found bool
		for _, update := range updates {
			if update.accountID == f.account1 && update.quorumID == f.quorum0 && update.period == currentPeriod {
				assert.Equal(t, uint64(64), update.usage) // 32 + 32 = 64
				found = true
			}
		}
		assert.True(t, found, "expected updateRecord not found")
	})

	t.Run("very large usage", func(t *testing.T) {
		batch := createTestBatch([]accountQuorumInfo{
			{
				account:           f.account1,
				quorumID:          f.quorum0,
				timestamp:         0,
				cumulativePayment: 0,
			},
		})
		// Set blob length to a very large value
		batch.BlobCertificates[0].BlobHeader.BlobCommitments.Length = math.MaxUint32
		updates, err := f.batchMeterer.BatchRequestUsage(batch)
		require.NoError(t, err)
		currentPeriod := meterer.GetReservationPeriodByNanosecond(0, 3600)
		var found bool
		for _, update := range updates {
			if update.accountID == f.account1 && update.quorumID == f.quorum0 && update.period == currentPeriod {
				assert.Greater(t, update.usage, uint64(0))
				// Check that usage is reasonable (less than MaxUint64/2 to avoid overflow)
				assert.Less(t, update.usage, uint64(math.MaxUint64/2))
				found = true
			}
		}
		assert.True(t, found, "expected updateRecord not found")
	})
}

// TestOverflowEdgeCases tests edge cases in overflow handling
func TestBatchMeterOverflowEdgeCases(t *testing.T) {
	f := setupTestFixtures(t)

	t.Run("exact bin limit", func(t *testing.T) {
		// Create a reservation
		binLimit := uint64(3600)
		reservations := f.createTestReservationMap(1, time.Unix(0, 0), time.Unix(9999999999, 0), f.quorum0)

		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil).Once()

		// Initialize the period record
		f.setupPeriodRecord(f.account1, f.quorum0, meterer.GetReservationPeriod(time.Now().Unix(), 3600))
		now := time.Now()
		currentPeriod := meterer.GetReservationPeriod(now.Unix(), 3600)
		accountUsage := f.batchMeterer.getOrCreateAccountUsage(f.account1)
		accountUsage.Lock.Lock()
		relativeIndex := f.getRelativeIndex(currentPeriod)
		accountUsage.PeriodRecords[f.quorum0] = make([]*pb.PeriodRecord, 3)
		accountUsage.PeriodRecords[f.quorum0][relativeIndex] = &pb.PeriodRecord{
			Index: uint32(currentPeriod),
			Usage: 0,
		}
		accountUsage.Lock.Unlock()

		// Try to use exactly bin limit
		updates := []updateRecord{
			{
				accountID: f.account1,
				quorumID:  f.quorum0,
				period:    currentPeriod,
				usage:     binLimit,
			},
		}

		err := f.batchMeterer.BatchMeterRequest(f.ctx, updates, now)
		require.NoError(t, err)

		// Verify usage is exactly at bin limit
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()
		assert.Equal(t, binLimit, accountUsage.PeriodRecords[f.quorum0][relativeIndex].Usage)
	})

	t.Run("just over bin limit", func(t *testing.T) {
		// Create a reservation
		binLimit := uint64(3600)
		reservations := f.createTestReservationMap(1, time.Unix(0, 0), time.Unix(9999999999, 0), f.quorum0)

		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil).Once()

		// Initialize the period record with some usage
		f.setupPeriodRecord(f.account1, f.quorum0, meterer.GetReservationPeriod(time.Now().Unix(), 3600))
		now := time.Now()
		currentPeriod := meterer.GetReservationPeriod(now.Unix(), 3600)
		accountUsage := f.batchMeterer.getOrCreateAccountUsage(f.account1)
		accountUsage.Lock.Lock()
		relativeIndex := f.getRelativeIndex(currentPeriod)
		accountUsage.PeriodRecords[f.quorum0] = make([]*pb.PeriodRecord, 3)
		accountUsage.PeriodRecords[f.quorum0][relativeIndex] = &pb.PeriodRecord{
			Index: uint32(currentPeriod),
			Usage: binLimit - 1,
		}
		accountUsage.Lock.Unlock()

		// Try to use just over bin limit
		updates := []updateRecord{
			{
				accountID: f.account1,
				quorumID:  f.quorum0,
				period:    currentPeriod,
				usage:     2, // Add 2 more to exceed limit
			},
		}

		err := f.batchMeterer.BatchMeterRequest(f.ctx, updates, now)
		require.NoError(t, err)

		// Verify overflow was handled
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()
		assert.Equal(t, binLimit, accountUsage.PeriodRecords[f.quorum0][relativeIndex].Usage)
	})

	t.Run("exactly 2x bin limit", func(t *testing.T) {
		// Create a reservation
		binLimit := uint64(3600)
		reservations := f.createTestReservationMap(1, time.Unix(0, 0), time.Unix(9999999999, 0), f.quorum0)

		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil).Once()

		// Initialize the period record
		f.setupPeriodRecord(f.account1, f.quorum0, meterer.GetReservationPeriod(time.Now().Unix(), 3600))
		now := time.Now()
		currentPeriod := meterer.GetReservationPeriod(now.Unix(), 3600)
		accountUsage := f.batchMeterer.getOrCreateAccountUsage(f.account1)
		accountUsage.Lock.Lock()
		relativeIndex := f.getRelativeIndex(currentPeriod)
		accountUsage.PeriodRecords[f.quorum0] = make([]*pb.PeriodRecord, 3)
		accountUsage.PeriodRecords[f.quorum0][relativeIndex] = &pb.PeriodRecord{
			Index: uint32(currentPeriod),
			Usage: 0,
		}
		accountUsage.Lock.Unlock()

		// Try to use exactly 2x bin limit
		updates := []updateRecord{
			{
				accountID: f.account1,
				quorumID:  f.quorum0,
				period:    currentPeriod,
				usage:     2 * binLimit,
			},
		}

		err := f.batchMeterer.BatchMeterRequest(f.ctx, updates, now)
		require.NoError(t, err)

		// Verify overflow was handled
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()
		assert.Equal(t, binLimit, accountUsage.PeriodRecords[f.quorum0][relativeIndex].Usage)
		overflowPeriod := currentPeriod + 3600
		overflowRelativeIndex := f.getRelativeIndex(overflowPeriod)
		assert.Equal(t, binLimit, accountUsage.PeriodRecords[f.quorum0][overflowRelativeIndex].Usage)
	})

	t.Run("over 2x bin limit", func(t *testing.T) {
		// Create a reservation
		binLimit := uint64(3600)
		reservations := f.createTestReservationMap(1, time.Unix(0, 0), time.Unix(9999999999, 0), f.quorum0)

		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil).Once()

		// Initialize the period record
		f.setupPeriodRecord(f.account1, f.quorum0, meterer.GetReservationPeriod(time.Now().Unix(), 3600))
		accountUsage := f.batchMeterer.getOrCreateAccountUsage(f.account1)
		accountUsage.Lock.Lock()
		accountUsage.PeriodRecords[f.quorum0] = make([]*pb.PeriodRecord, 3)
		for i := range accountUsage.PeriodRecords[f.quorum0] {
			accountUsage.PeriodRecords[f.quorum0][i] = &pb.PeriodRecord{
				Index: 0,
				Usage: 0,
			}
		}
		accountUsage.Lock.Unlock()

		// Try to use over 2x bin limit
		now := time.Now()
		currentPeriod := meterer.GetReservationPeriod(now.Unix(), 3600)
		updates := []updateRecord{
			{
				accountID: f.account1,
				quorumID:  f.quorum0,
				period:    currentPeriod,
				usage:     2*binLimit + 1,
			},
		}

		err := f.batchMeterer.BatchMeterRequest(f.ctx, updates, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "overflow usage exceeds bin limit")
	})
}

// TestPeriodRecordEdgeCases tests edge cases in period record management
func TestBatchMeterPeriodRecordEdgeCases(t *testing.T) {
	f := setupTestFixtures(t)

	t.Run("period transition", func(t *testing.T) {
		// Clear the batchMeterer cache to ensure a clean state
		f.resetBatchMeterer()

		// Create a reservation that spans multiple periods
		reservations := f.createTestReservationMap(100, time.Unix(0, 0), time.Unix(9999999999, 0), f.quorum0)

		f.mockState.On("GetMinNumSymbols").Return(uint64(32)).Maybe()
		f.mockState.On("GetReservationWindow").Return(uint64(3600)).Maybe()
		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil)

		// First request in current period
		now := time.Now()
		currentPeriod := meterer.GetReservationPeriod(now.Unix(), f.reservationWindow)
		f.setupPeriodRecord(f.account1, f.quorum0, currentPeriod)
		updates := []updateRecord{
			{
				accountID: f.account1,
				quorumID:  f.quorum0,
				period:    currentPeriod,
				usage:     100,
			},
		}
		err := f.batchMeterer.BatchMeterRequest(f.ctx, updates, now)
		require.NoError(t, err)

		// Second request in next period
		nextPeriodTime := now.Add(time.Duration(3600) * time.Second)
		nextPeriod := meterer.GetReservationPeriod(nextPeriodTime.Unix(), f.reservationWindow)
		updatesNext := []updateRecord{
			{
				accountID: f.account1,
				quorumID:  f.quorum0,
				period:    nextPeriod, // Use next period instead of current period
				usage:     100,
			},
		}
		err = f.batchMeterer.BatchMeterRequest(f.ctx, updatesNext, nextPeriodTime)
		require.NoError(t, err)

		// Verify both periods have correct usage
		f.verifyUsage(t, f.account1, f.quorum0, currentPeriod, 100)

		f.verifyUsage(t, f.account1, f.quorum0, nextPeriod, 100)
	})

	t.Run("circular buffer wrapping", func(t *testing.T) {
		// Clear the batchMeterer cache to ensure a clean state
		f.resetBatchMeterer()

		// Create a reservation
		reservations := f.createTestReservationMap(100, time.Unix(0, 0), time.Unix(9999999999, 0), f.quorum0)

		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil).Maybe()

		// Make requests across 4 periods to test buffer wrapping
		now := time.Now()
		for i := 0; i < 4; i++ {
			periodTime := now.Add(time.Duration(i*3600) * time.Second)
			periodIndex := meterer.GetReservationPeriod(periodTime.Unix(), f.reservationWindow)
			updates := []updateRecord{
				{
					accountID: f.account1,
					quorumID:  f.quorum0,
					period:    periodIndex,
					usage:     100,
				},
			}
			err := f.batchMeterer.BatchMeterRequest(f.ctx, updates, periodTime)
			require.NoError(t, err)
		}

		// Verify the oldest record was overwritten
		accountUsage := f.batchMeterer.getOrCreateAccountUsage(f.account1)
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()

		// Should only have 3 records (buffer size)
		records := accountUsage.PeriodRecords[f.quorum0]
		assert.Equal(t, 3, len(records))
	})

	t.Run("just started reservation", func(t *testing.T) {
		// Clear existing expectations
		f.mockState.ExpectedCalls = nil
		f.mockState.Calls = nil

		// Create a reservation that just started
		now := time.Now()
		reservations := f.createTestReservationMap(100, now, now.Add(24*time.Hour), f.quorum0)

		f.mockState.On("GetMinNumSymbols").Return(uint64(32)).Maybe()
		f.mockState.On("GetReservationWindow").Return(uint64(3600)).Maybe()
		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil).Once()

		currentPeriod := meterer.GetReservationPeriod(now.Unix(), f.reservationWindow)
		f.setupPeriodRecord(f.account1, f.quorum0, currentPeriod)
		updates := []updateRecord{
			{
				accountID: f.account1,
				quorumID:  f.quorum0,
				period:    currentPeriod,
				usage:     100,
			},
		}

		err := f.batchMeterer.BatchMeterRequest(f.ctx, updates, now)
		require.NoError(t, err)
	})

	t.Run("zero symbols per second", func(t *testing.T) {
		// Clear existing expectations
		f.mockState.ExpectedCalls = nil
		f.mockState.Calls = nil

		// Create a reservation with zero symbols per second
		now := time.Now()
		reservations := f.createTestReservationMap(0, now, now.Add(24*time.Hour), f.quorum0)

		f.mockState.On("GetMinNumSymbols").Return(uint64(32)).Maybe()
		f.mockState.On("GetReservationWindow").Return(uint64(3600)).Maybe()
		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil).Once()

		currentPeriod := meterer.GetReservationPeriod(now.Unix(), f.reservationWindow)
		f.setupPeriodRecord(f.account1, f.quorum0, currentPeriod)
		updates := []updateRecord{
			{
				accountID: f.account1,
				quorumID:  f.quorum0,
				period:    currentPeriod,
				usage:     100,
			},
		}

		err := f.batchMeterer.BatchMeterRequest(f.ctx, updates, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed usage validation for quorum")
	})

	t.Run("invalid timestamps", func(t *testing.T) {
		// Create a reservation with invalid timestamps
		now := time.Now()
		reservations := f.createTestReservationMap(100, now.Add(-24*time.Hour), now, f.quorum0)

		f.mockState.On("GetMinNumSymbols").Return(uint64(32)).Maybe()
		f.mockState.On("GetReservationWindow").Return(uint64(3600)).Maybe()
		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil).Twice()

		currentPeriod := meterer.GetReservationPeriod(now.Unix(), f.reservationWindow)
		f.setupPeriodRecord(f.account1, f.quorum0, currentPeriod)
		updates := []updateRecord{
			{
				accountID: f.account1,
				quorumID:  f.quorum0,
				period:    currentPeriod,
				usage:     100,
			},
		}

		err := f.batchMeterer.BatchMeterRequest(f.ctx, updates, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid reservation period for quorum")
	})
}

// TestReservationEdgeCases tests edge cases in reservation validation
func TestBatchMeterReservationEdgeCases(t *testing.T) {
	f := setupTestFixtures(t)

	t.Run("no reservation", func(t *testing.T) {
		// No reservation for the account
		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(nil, nil).Once()

		// Initialize the period record
		f.setupPeriodRecord(f.account1, f.quorum0, meterer.GetReservationPeriod(time.Now().Unix(), f.reservationWindow))
		now := time.Now()
		currentPeriod := meterer.GetReservationPeriod(now.Unix(), f.reservationWindow)
		accountUsage := f.batchMeterer.getOrCreateAccountUsage(f.account1)
		accountUsage.Lock.Lock()
		relativeIndex := f.getRelativeIndex(currentPeriod)
		accountUsage.PeriodRecords[f.quorum0] = make([]*pb.PeriodRecord, 3)
		accountUsage.PeriodRecords[f.quorum0][relativeIndex] = &pb.PeriodRecord{
			Index: uint32(currentPeriod),
			Usage: 0,
		}
		accountUsage.Lock.Unlock()

		// Try to use without reservation
		updates := []updateRecord{
			{
				accountID: f.account1,
				quorumID:  f.quorum0,
				period:    currentPeriod,
				usage:     100,
			},
		}

		err := f.batchMeterer.BatchMeterRequest(f.ctx, updates, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no reservation")
	})

	t.Run("invalid reservation period", func(t *testing.T) {
		// Create a reservation
		reservations := f.createTestReservationMap(100, time.Unix(0, 0), time.Unix(9999999999, 0), f.quorum0)

		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil).Once()

		// Initialize the period record
		f.setupPeriodRecord(f.account1, f.quorum0, meterer.GetReservationPeriod(time.Now().Unix(), f.reservationWindow))
		now := time.Now()
		currentPeriod := meterer.GetReservationPeriod(now.Unix(), f.reservationWindow)
		accountUsage := f.batchMeterer.getOrCreateAccountUsage(f.account1)
		accountUsage.Lock.Lock()
		relativeIndex := f.getRelativeIndex(currentPeriod)
		accountUsage.PeriodRecords[f.quorum0] = make([]*pb.PeriodRecord, 3)
		accountUsage.PeriodRecords[f.quorum0][relativeIndex] = &pb.PeriodRecord{
			Index: uint32(currentPeriod),
			Usage: 0,
		}
		accountUsage.Lock.Unlock()

		// Try to use with invalid period
		updates := []updateRecord{
			{
				accountID: f.account1,
				quorumID:  f.quorum0,
				period:    currentPeriod + 7200, // Period too far in the future
				usage:     100,
			},
		}

		err := f.batchMeterer.BatchMeterRequest(f.ctx, updates, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid reservation period for quorum")
	})
}

// TestBatchMeterRequestRollback tests the rollback functionality in BatchMeterRequest
func TestBatchMeterRequestRollback(t *testing.T) {
	f := setupTestFixtures(t)

	t.Run("rollback on reservation error", func(t *testing.T) {
		// Clear the batchMeterer cache to ensure a clean state
		f.resetBatchMeterer()

		// Create a reservation
		reservations := f.createTestReservationMap(100, time.Unix(0, 0), time.Unix(9999999999, 0), f.quorum0)

		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil).Once()

		// Initialize period records
		f.setupPeriodRecord(f.account1, f.quorum0, meterer.GetReservationPeriod(time.Now().Unix(), f.reservationWindow))
		now := time.Now()
		currentPeriod := meterer.GetReservationPeriod(now.Unix(), f.reservationWindow)
		accountUsage := f.batchMeterer.getOrCreateAccountUsage(f.account1)
		accountUsage.Lock.Lock()
		relativeIndex := f.getRelativeIndex(currentPeriod)
		accountUsage.PeriodRecords[f.quorum0] = make([]*pb.PeriodRecord, 3)
		accountUsage.PeriodRecords[f.quorum0][relativeIndex] = &pb.PeriodRecord{
			Index: uint32(currentPeriod),
			Usage: 0,
		}
		accountUsage.Lock.Unlock()

		// Create updates
		updates := []updateRecord{
			{
				accountID: f.account1,
				quorumID:  f.quorum0,
				period:    currentPeriod,
				usage:     50,
			},
		}

		// Process the request
		err := f.batchMeterer.BatchMeterRequest(f.ctx, updates, now)
		require.NoError(t, err)

		// Verify usage was updated
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()
		assert.Equal(t, uint64(50), accountUsage.PeriodRecords[f.quorum0][relativeIndex].Usage)

		// Now try to update with an invalid reservation
		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(nil, errors.New("reservation not found")).Once()

		// Create updates for invalid reservation
		updates = []updateRecord{
			{
				accountID: f.account1,
				quorumID:  f.quorum0,
				period:    currentPeriod,
				usage:     100,
			},
		}

		// Process the request
		err = f.batchMeterer.BatchMeterRequest(f.ctx, updates, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "reservation not found")

		// Verify usage was rolled back
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()
		assert.Equal(t, uint64(50), accountUsage.PeriodRecords[f.quorum0][relativeIndex].Usage)
	})

	t.Run("rollback on validation error", func(t *testing.T) {
		// Clear the batchMeterer cache to ensure a clean state
		f.resetBatchMeterer()

		// Create a reservation
		reservations := f.createTestReservationMap(100, time.Unix(0, 0), time.Unix(9999999999, 0), f.quorum0)

		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil).Twice()

		// Initialize period records
		f.setupPeriodRecord(f.account1, f.quorum0, meterer.GetReservationPeriod(time.Now().Unix(), f.reservationWindow))
		now := time.Now()
		currentPeriod := meterer.GetReservationPeriod(now.Unix(), f.reservationWindow)
		accountUsage := f.batchMeterer.getOrCreateAccountUsage(f.account1)
		accountUsage.Lock.Lock()
		relativeIndex := f.getRelativeIndex(currentPeriod)
		accountUsage.PeriodRecords[f.quorum0] = make([]*pb.PeriodRecord, 3)
		accountUsage.PeriodRecords[f.quorum0][relativeIndex] = &pb.PeriodRecord{
			Index: uint32(currentPeriod),
			Usage: 0,
		}
		accountUsage.Lock.Unlock()

		// Create updates
		updates := []updateRecord{
			{
				accountID: f.account1,
				quorumID:  f.quorum0,
				period:    currentPeriod,
				usage:     50,
			},
		}

		// Process the request
		err := f.batchMeterer.BatchMeterRequest(f.ctx, updates, now)
		require.NoError(t, err)

		// Create updates for invalid period (far future period)
		updates = []updateRecord{
			{
				accountID: f.account1,
				quorumID:  f.quorum0,
				period:    currentPeriod + 2*3600, // Far future period, should be invalid
				usage:     100,
			},
		}

		// Process the request
		err = f.batchMeterer.BatchMeterRequest(f.ctx, updates, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid reservation period for quorum")

		// Verify usage was rolled back
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()
		assert.Equal(t, uint64(50), accountUsage.PeriodRecords[f.quorum0][relativeIndex].Usage)
	})
}
