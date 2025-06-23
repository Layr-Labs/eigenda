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

const (
	// Time related constants
	reservationInterval   = 300 // Changed from 3600 to 300 seconds (5 minutes)
	farFutureTimestamp    = 9999999999
	defaultNumBins        = 3
	defaultMinSymbols     = 32
	defaultUpdateTimeout  = 1 * time.Second
	defaultUpdateInterval = 1 * time.Minute

	// Test data constants
	testBlobLength      = 32
	testUsageAmount     = 100
	testUsageAmount2    = 200
	testLargeUsage      = math.MaxUint32
	testBatchRootValue  = 1
	testSignatureValue1 = 1
	testSignatureValue2 = 2
	testSignatureValue3 = 3

	// Test account and quorum constants
	testAccount1Hex    = "0x1"
	testAccount2Hex    = "0x2"
	testQuorum0        = 0
	testQuorum1        = 1
	testReferenceBlock = 100
	testBlobVersion    = 0
	testRelayKey0      = 0
	testRelayKey1      = 1

	// Test period constants
	testPeriodOffset     = reservationInterval * 2
	testPeriodMultiplier = 2

	// Test usage constants (after rounding)
	testRoundedUsage1         = 4 * defaultMinSymbols
	testRoundedUsage2         = 7 * defaultMinSymbols
	testRoundedUsageSum       = testRoundedUsage1 + testRoundedUsage1
	testBlobLengthForRounding = 97 // Will round up to testRoundedUsage1
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
		MinNumSymbols:     defaultMinSymbols,
		ReservationWindow: reservationInterval,
		NumBins:           defaultNumBins,
		ChainReadTimeout:  defaultUpdateTimeout,
		UpdateInterval:    defaultUpdateInterval,
	}
}

// CreateTestReservation creates a test reservation with the given parameters
func CreateTestReservation(symbolsPerSecond uint64, startTime, endTime time.Time) *core.ReservedPayment {
	return &core.ReservedPayment{
		SymbolsPerSecond: symbolsPerSecond,
		StartTimestamp:   uint64(startTime.Unix()),
		EndTimestamp:     uint64(endTime.Unix()),
	}
}

// InitializePeriodRecord initializes a period record for the given account and quorum
func InitializePeriodRecord(batchMeterer *BatchMeterer, accountID common.Address, quorumID core.QuorumID, period uint64, reservationWindow uint64) {
	accountUsage := batchMeterer.getOrCreateAccountUsage(accountID)
	accountUsage.Lock.Lock()
	defer accountUsage.Lock.Unlock()

	if accountUsage.PeriodRecords[quorumID] == nil {
		accountUsage.PeriodRecords[quorumID] = make([]*pb.PeriodRecord, meterer.MinNumBins)
	}
	relativeIndex := uint32((period / reservationWindow) % uint64(meterer.MinNumBins))
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
	blobLength        uint64 // Required blob length for the test
}

// createTestBatch creates a test batch with the given account quorum information
func createTestBatch(accountQuorumInfo []accountQuorumInfo) *corev2.Batch {
	batch := &corev2.Batch{
		BatchHeader: &corev2.BatchHeader{
			BatchRoot:            [32]byte{testBatchRootValue, testBatchRootValue, testBatchRootValue},
			ReferenceBlockNumber: testReferenceBlock,
		},
		BlobCertificates: make([]*corev2.BlobCertificate, len(accountQuorumInfo)),
	}

	for i, info := range accountQuorumInfo {
		batch.BlobCertificates[i] = &corev2.BlobCertificate{
			BlobHeader: &corev2.BlobHeader{
				BlobVersion:     testBlobVersion,
				BlobCommitments: encoding.BlobCommitments{},
				QuorumNumbers:   []core.QuorumID{info.quorumID},
				PaymentMetadata: core.PaymentMetadata{
					AccountID:         info.account,
					Timestamp:         int64(info.timestamp),
					CumulativePayment: big.NewInt(int64(info.cumulativePayment)),
				},
			},
			Signature: []byte{testSignatureValue1, testSignatureValue2, testSignatureValue3},
			RelayKeys: []corev2.RelayKey{testRelayKey0, testRelayKey1},
		}
		batch.BlobCertificates[i].BlobHeader.BlobCommitments.Length = uint(info.blobLength)
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

	// Set up default mock for GetPaymentGlobalParams
	mockParams := &meterer.PaymentVaultParams{
		QuorumProtocolConfigs: make(map[core.QuorumID]*core.PaymentQuorumProtocolConfig),
	}
	for i := uint8(0); i < 2; i++ {
		mockParams.QuorumProtocolConfigs[core.QuorumID(i)] = &core.PaymentQuorumProtocolConfig{
			MinNumSymbols:              config.MinNumSymbols,
			ReservationRateLimitWindow: config.ReservationWindow,
		}
	}
	mockState.On("GetPaymentGlobalParams").Return(mockParams, nil)

	batchMeterer := NewBatchMeterer(meterer.Config{
		ChainReadTimeout: config.ChainReadTimeout,
		UpdateInterval:   config.UpdateInterval,
	}, mockState, logger)

	return &testFixtures{
		ctx:               ctx,
		logger:            logger,
		mockState:         mockState,
		config:            config,
		batchMeterer:      batchMeterer,
		account1:          common.HexToAddress(testAccount1Hex),
		account2:          common.HexToAddress(testAccount2Hex),
		quorum0:           testQuorum0,
		quorum1:           testQuorum1,
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
	InitializePeriodRecord(f.batchMeterer, accountID, quorumID, period, f.reservationWindow)
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

// TestBatchMeterProcessBatch tests the processBatch function
func TestBatchMeterProcessBatch(t *testing.T) {
	f := setupTestFixtures(t)

	t.Run("valid reservation", func(t *testing.T) {
		f.resetBatchMeterer()

		// Create a valid reservation
		now := time.Now()
		reservations := f.createTestReservationMap(testUsageAmount, time.Unix(0, 0), time.Unix(farFutureTimestamp, 0), f.quorum0)

		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil).Once()

		currentPeriod := meterer.GetReservationPeriod(now.Unix(), f.reservationWindow)
		f.setupPeriodRecord(f.account1, f.quorum0, currentPeriod)

		updates := []*UpdateRecord{
			newUpdateRecord(f.account1, f.quorum0, currentPeriod, testUsageAmount),
		}

		params, _ := f.mockState.GetPaymentGlobalParams()
		err := f.batchMeterer.processBatch(f.ctx, params, updates, now)
		require.NoError(t, err)

		f.verifyUsage(t, f.account1, f.quorum0, currentPeriod, testUsageAmount)
	})

	t.Run("period transition", func(t *testing.T) {
		f.resetBatchMeterer()

		// Create a valid reservation
		reservations := f.createTestReservationMap(testUsageAmount, time.Unix(0, 0), time.Unix(farFutureTimestamp, 0), f.quorum0)

		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil).Times(testPeriodMultiplier)

		// Create usage map for current period
		now := time.Now()
		currentPeriod := meterer.GetReservationPeriod(now.Unix(), f.reservationWindow)

		// Initialize period records for both current and next period
		f.setupPeriodRecord(f.account1, f.quorum0, currentPeriod)
		nextPeriod := currentPeriod + reservationInterval
		f.setupPeriodRecord(f.account1, f.quorum0, nextPeriod)

		// Create updates array for current period
		updates := []*UpdateRecord{
			newUpdateRecord(f.account1, f.quorum0, currentPeriod, testUsageAmount),
		}

		// Process the request
		params, _ := f.mockState.GetPaymentGlobalParams()
		err := f.batchMeterer.processBatch(f.ctx, params, updates, now)
		require.NoError(t, err)

		// Create usage map for next period
		nextNow := now.Add(time.Duration(reservationInterval) * time.Second)

		// Create updates array for next period
		updatesNext := []*UpdateRecord{
			newUpdateRecord(f.account1, f.quorum0, nextPeriod, testUsageAmount),
		}

		// Process the request for next period
		params, _ = f.mockState.GetPaymentGlobalParams()
		err = f.batchMeterer.processBatch(f.ctx, params, updatesNext, nextNow)
		require.NoError(t, err)

		// Verify both periods were updated correctly
		f.verifyUsage(t, f.account1, f.quorum0, currentPeriod, testUsageAmount)
		f.verifyUsage(t, f.account1, f.quorum0, nextPeriod, testUsageAmount)
	})

	t.Run("multiple periods in same request", func(t *testing.T) {
		f.resetBatchMeterer()

		// Create a valid reservation
		reservations := f.createTestReservationMap(testUsageAmount, time.Unix(0, 0), time.Unix(farFutureTimestamp, 0), f.quorum0)

		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil).Once() // Only need one call since we process both periods together

		// Create usage map
		now := time.Now()
		currentPeriod := meterer.GetReservationPeriod(now.Unix(), f.reservationWindow)
		prevPeriod := currentPeriod - reservationInterval // Previous period

		// Initialize period records
		f.setupPeriodRecord(f.account1, f.quorum0, currentPeriod)
		f.setupPeriodRecord(f.account1, f.quorum0, prevPeriod)

		// Create updates array with both periods
		updates := []*UpdateRecord{
			newUpdateRecord(f.account1, f.quorum0, currentPeriod, testUsageAmount),
			newUpdateRecord(f.account1, f.quorum0, prevPeriod, testUsageAmount),
		}

		// Process both periods in a single request
		params, _ := f.mockState.GetPaymentGlobalParams()
		err := f.batchMeterer.processBatch(f.ctx, params, updates, now)
		require.NoError(t, err)

		// Verify both periods were updated correctly
		f.verifyUsage(t, f.account1, f.quorum0, currentPeriod, testUsageAmount)
		f.verifyUsage(t, f.account1, f.quorum0, prevPeriod, testUsageAmount)
	})

	t.Run("circular buffer wrapping", func(t *testing.T) {
		// Clear the batchMeterer cache to ensure a clean state
		f.resetBatchMeterer()

		// Create a reservation
		reservations := f.createTestReservationMap(testUsageAmount, time.Unix(0, 0), time.Unix(farFutureTimestamp, 0), f.quorum0)

		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil).Maybe()

		// Make requests across 4 periods to test buffer wrapping
		now := time.Now()
		for i := 0; i < 4; i++ {
			periodTime := now.Add(time.Duration(i*reservationInterval) * time.Second)
			periodIndex := meterer.GetReservationPeriod(periodTime.Unix(), f.reservationWindow)
			updates := []*UpdateRecord{
				newUpdateRecord(f.account1, f.quorum0, periodIndex, testUsageAmount),
			}
			params, _ := f.mockState.GetPaymentGlobalParams()
			err := f.batchMeterer.processBatch(f.ctx, params, updates, periodTime)
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
}

// TestBatchMeterMeterBatch tests the MeterBatch function
func TestBatchMeterMeterBatch(t *testing.T) {
	f := setupTestFixtures(t)

	t.Run("successful batch metering", func(t *testing.T) {
		// Create reservations for accounts
		reservationsAcc1 := f.createTestReservationMap(testUsageAmount, time.Unix(0, 0), time.Unix(farFutureTimestamp, 0), f.quorum0, f.quorum1)
		reservationsAcc2 := f.createTestReservationMap(testUsageAmount2, time.Unix(0, 0), time.Unix(farFutureTimestamp, 0), f.quorum1)

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
		periodFromTimestamp := meterer.GetReservationPeriodByNanosecond(int64(timestampNano), reservationInterval)

		// Initialize period records for both accounts using periodFromTimestamp
		f.setupPeriodRecord(f.account1, f.quorum0, periodFromTimestamp)
		f.setupPeriodRecord(f.account1, f.quorum1, periodFromTimestamp)
		f.setupPeriodRecord(f.account2, f.quorum1, periodFromTimestamp)

		batch := createTestBatch([]accountQuorumInfo{
			{
				account:           f.account1,
				quorumID:          f.quorum0,
				timestamp:         timestampNano,
				cumulativePayment: testUsageAmount,
				blobLength:        testBlobLengthForRounding,
			},
			{
				account:           f.account2,
				quorumID:          f.quorum1,
				timestamp:         timestampNano,
				cumulativePayment: testUsageAmount2,
				blobLength:        testBlobLengthForRounding + 96, // Will round up to testRoundedUsage2
			},
		})

		err := f.batchMeterer.MeterBatch(f.ctx, batch, now)
		require.NoError(t, err)
		f.verifyUsage(t, f.account1, f.quorum0, periodFromTimestamp, testRoundedUsage1)
		f.verifyUsage(t, f.account2, f.quorum1, periodFromTimestamp, testRoundedUsage2)
	})

	t.Run("batch with invalid account", func(t *testing.T) {
		// Clear existing expectations
		f.mockState.ExpectedCalls = nil
		f.mockState.Calls = nil

		// Set up basic mock expectations
		mockParams := &meterer.PaymentVaultParams{
			QuorumProtocolConfigs: make(map[core.QuorumID]*core.PaymentQuorumProtocolConfig),
		}
		for i := uint8(0); i < 2; i++ {
			mockParams.QuorumProtocolConfigs[core.QuorumID(i)] = &core.PaymentQuorumProtocolConfig{
				MinNumSymbols:              defaultMinSymbols,
				ReservationRateLimitWindow: reservationInterval,
			}
		}
		f.mockState.On("GetPaymentGlobalParams").Return(mockParams, nil)

		// Create valid reservations for account 1
		reservationsAcc1 := f.createTestReservationMap(testUsageAmount, time.Unix(0, 0), time.Unix(farFutureTimestamp, 0), f.quorum0)

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
				cumulativePayment: testUsageAmount,
				blobLength:        testBlobLength,
			},
			{
				account:           f.account2,
				quorumID:          f.quorum1,
				timestamp:         uint64(now.UnixNano()),
				cumulativePayment: testUsageAmount2,
				blobLength:        testBlobLength,
			},
		})

		// Set blob lengths in the batch
		for _, cert := range batch.BlobCertificates {
			cert.BlobHeader.BlobCommitments.Length = testBlobLength // Set a valid length
		}

		// Call the function
		err := f.batchMeterer.MeterBatch(f.ctx, batch, now)

		// Assert the results - should fail due to account 2 having no reservation
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no reservation")
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

// TestProcessBatchUsageEdgeCases tests edge cases in batch request usage calculation
func TestProcessBatchUsageEdgeCases(t *testing.T) {
	f := setupTestFixtures(t)

	t.Run("empty batch", func(t *testing.T) {
		batch := &corev2.Batch{
			BatchHeader:      &corev2.BatchHeader{},
			BlobCertificates: []*corev2.BlobCertificate{},
		}
		params, _ := f.mockState.GetPaymentGlobalParams()
		_, err := f.batchMeterer.batchRequestUsage(params, batch)
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
		params, _ := f.mockState.GetPaymentGlobalParams()
		_, err := f.batchMeterer.batchRequestUsage(params, batch)
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
				blobLength:        0,
			},
		})
		// Set blob length to 0 explicitly
		batch.BlobCertificates[0].BlobHeader.BlobCommitments.Length = 0
		params, _ := f.mockState.GetPaymentGlobalParams()
		updates, err := f.batchMeterer.batchRequestUsage(params, batch)
		require.NoError(t, err)
		currentPeriod := meterer.GetReservationPeriodByNanosecond(0, reservationInterval)
		var found bool
		for _, update := range updates {
			if update.accountID == f.account1 && update.quorumID == f.quorum0 && update.period == currentPeriod {
				assert.Equal(t, uint64(defaultMinSymbols), update.usage) // MinNumSymbols is 32
				found = true
			}
		}
		assert.True(t, found, "expected UpdateRecord not found")
	})

	t.Run("duplicate account quorum", func(t *testing.T) {
		batch := createTestBatch([]accountQuorumInfo{
			{
				account:           f.account1,
				quorumID:          f.quorum0,
				timestamp:         0,
				cumulativePayment: testUsageAmount,
				blobLength:        testBlobLengthForRounding,
			},
			{
				account:           f.account1,
				quorumID:          f.quorum0,
				timestamp:         0,
				cumulativePayment: testUsageAmount,
				blobLength:        testBlobLengthForRounding,
			},
		})
		params, _ := f.mockState.GetPaymentGlobalParams()
		updates, err := f.batchMeterer.batchRequestUsage(params, batch)
		require.NoError(t, err)
		currentPeriod := meterer.GetReservationPeriodByNanosecond(0, reservationInterval)
		var found bool
		for _, update := range updates {
			if update.accountID == f.account1 && update.quorumID == f.quorum0 && update.period == currentPeriod {
				assert.Equal(t, uint64(testRoundedUsageSum), update.usage)
				found = true
			}
		}
		assert.True(t, found, "expected UpdateRecord not found")
	})

	t.Run("very large usage", func(t *testing.T) {
		batch := createTestBatch([]accountQuorumInfo{
			{
				account:           f.account1,
				quorumID:          f.quorum0,
				timestamp:         0,
				cumulativePayment: 0,
				blobLength:        0,
			},
		})
		// Set blob length to a very large value
		batch.BlobCertificates[0].BlobHeader.BlobCommitments.Length = testLargeUsage
		params, _ := f.mockState.GetPaymentGlobalParams()
		updates, err := f.batchMeterer.batchRequestUsage(params, batch)
		require.NoError(t, err)
		currentPeriod := meterer.GetReservationPeriodByNanosecond(0, reservationInterval)
		var found bool
		for _, update := range updates {
			if update.accountID == f.account1 && update.quorumID == f.quorum0 && update.period == currentPeriod {
				assert.Greater(t, update.usage, uint64(0))
				// Check that usage is reasonable (less than MaxUint64/2 to avoid overflow)
				assert.Less(t, update.usage, uint64(math.MaxUint64/2))
				found = true
			}
		}
		assert.True(t, found, "expected UpdateRecord not found")
	})
}

// TestOverflowEdgeCases tests edge cases in overflow handling
func TestBatchMeterOverflowEdgeCases(t *testing.T) {
	f := setupTestFixtures(t)

	t.Run("exact bin limit", func(t *testing.T) {
		// Create a reservation
		binLimit := uint64(reservationInterval)
		reservations := f.createTestReservationMap(1, time.Unix(0, 0), time.Unix(farFutureTimestamp, 0), f.quorum0)

		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil).Once()

		// Initialize the period record
		f.setupPeriodRecord(f.account1, f.quorum0, meterer.GetReservationPeriod(time.Now().Unix(), reservationInterval))
		now := time.Now()
		currentPeriod := meterer.GetReservationPeriod(now.Unix(), reservationInterval)
		accountUsage := f.batchMeterer.getOrCreateAccountUsage(f.account1)
		accountUsage.Lock.Lock()
		relativeIndex := f.getRelativeIndex(currentPeriod)
		accountUsage.PeriodRecords[f.quorum0] = make([]*pb.PeriodRecord, defaultNumBins)
		accountUsage.PeriodRecords[f.quorum0][relativeIndex] = &pb.PeriodRecord{
			Index: uint32(currentPeriod),
			Usage: 0,
		}
		accountUsage.Lock.Unlock()

		// Try to use exactly bin limit
		updates := []*UpdateRecord{
			newUpdateRecord(f.account1, f.quorum0, currentPeriod, binLimit),
		}

		params, _ := f.mockState.GetPaymentGlobalParams()
		err := f.batchMeterer.processBatch(f.ctx, params, updates, now)
		require.NoError(t, err)

		// Verify usage is exactly at bin limit
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()
		assert.Equal(t, binLimit, accountUsage.PeriodRecords[f.quorum0][relativeIndex].Usage)
	})

	t.Run("just over bin limit", func(t *testing.T) {
		// Create a reservation
		binLimit := uint64(reservationInterval)
		reservations := f.createTestReservationMap(1, time.Unix(0, 0), time.Unix(farFutureTimestamp, 0), f.quorum0)

		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil).Once()

		// Initialize the period record with some usage
		f.setupPeriodRecord(f.account1, f.quorum0, meterer.GetReservationPeriod(time.Now().Unix(), reservationInterval))
		now := time.Now()
		currentPeriod := meterer.GetReservationPeriod(now.Unix(), reservationInterval)
		accountUsage := f.batchMeterer.getOrCreateAccountUsage(f.account1)
		accountUsage.Lock.Lock()
		relativeIndex := f.getRelativeIndex(currentPeriod)
		accountUsage.PeriodRecords[f.quorum0] = make([]*pb.PeriodRecord, defaultNumBins)
		accountUsage.PeriodRecords[f.quorum0][relativeIndex] = &pb.PeriodRecord{
			Index: uint32(currentPeriod),
			Usage: binLimit - 1,
		}
		accountUsage.Lock.Unlock()

		// Try to use just over bin limit
		updates := []*UpdateRecord{
			newUpdateRecord(f.account1, f.quorum0, currentPeriod, 2), // Add 2 more to exceed limit
		}

		params, _ := f.mockState.GetPaymentGlobalParams()
		err := f.batchMeterer.processBatch(f.ctx, params, updates, now)
		require.NoError(t, err)

		// Verify overflow was handled
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()
		assert.Equal(t, binLimit, accountUsage.PeriodRecords[f.quorum0][relativeIndex].Usage)
	})

	t.Run("exactly 2x bin limit", func(t *testing.T) {
		// Create a reservation
		binLimit := uint64(reservationInterval)
		reservations := f.createTestReservationMap(1, time.Unix(0, 0), time.Unix(farFutureTimestamp, 0), f.quorum0)

		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil).Once()

		// Initialize the period record
		f.setupPeriodRecord(f.account1, f.quorum0, meterer.GetReservationPeriod(time.Now().Unix(), reservationInterval))
		now := time.Now()
		currentPeriod := meterer.GetReservationPeriod(now.Unix(), reservationInterval)
		accountUsage := f.batchMeterer.getOrCreateAccountUsage(f.account1)
		accountUsage.Lock.Lock()
		relativeIndex := f.getRelativeIndex(currentPeriod)
		accountUsage.PeriodRecords[f.quorum0] = make([]*pb.PeriodRecord, defaultNumBins)
		accountUsage.PeriodRecords[f.quorum0][relativeIndex] = &pb.PeriodRecord{
			Index: uint32(currentPeriod),
			Usage: 0,
		}
		accountUsage.Lock.Unlock()

		// Try to use exactly 2x bin limit
		updates := []*UpdateRecord{
			newUpdateRecord(f.account1, f.quorum0, currentPeriod, 2*binLimit),
		}

		params, _ := f.mockState.GetPaymentGlobalParams()
		err := f.batchMeterer.processBatch(f.ctx, params, updates, now)
		require.NoError(t, err)

		// Verify overflow was handled
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()
		assert.Equal(t, binLimit, accountUsage.PeriodRecords[f.quorum0][relativeIndex].Usage)
		overflowPeriod := currentPeriod + reservationInterval
		overflowRelativeIndex := f.getRelativeIndex(overflowPeriod)
		assert.Equal(t, binLimit, accountUsage.PeriodRecords[f.quorum0][overflowRelativeIndex].Usage)
	})

	t.Run("over 2x bin limit", func(t *testing.T) {
		// Create a reservation
		binLimit := uint64(reservationInterval)
		reservations := f.createTestReservationMap(1, time.Unix(0, 0), time.Unix(farFutureTimestamp, 0), f.quorum0)

		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil).Once()

		// Initialize the period record
		f.setupPeriodRecord(f.account1, f.quorum0, meterer.GetReservationPeriod(time.Now().Unix(), reservationInterval))
		accountUsage := f.batchMeterer.getOrCreateAccountUsage(f.account1)
		accountUsage.Lock.Lock()
		accountUsage.PeriodRecords[f.quorum0] = make([]*pb.PeriodRecord, defaultNumBins)
		for i := range accountUsage.PeriodRecords[f.quorum0] {
			accountUsage.PeriodRecords[f.quorum0][i] = &pb.PeriodRecord{
				Index: 0,
				Usage: 0,
			}
		}
		accountUsage.Lock.Unlock()

		// Try to use over 2x bin limit
		now := time.Now()
		currentPeriod := meterer.GetReservationPeriod(now.Unix(), reservationInterval)
		updates := []*UpdateRecord{
			newUpdateRecord(f.account1, f.quorum0, currentPeriod, 2*binLimit+1),
		}

		params, _ := f.mockState.GetPaymentGlobalParams()
		err := f.batchMeterer.processBatch(f.ctx, params, updates, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "usage exceeds bin limit")
	})
}

// TestBatchMeterPeriodRecordEdgeCases tests edge cases in period record management
func TestBatchMeterPeriodRecordEdgeCases(t *testing.T) {
	f := setupTestFixtures(t)

	t.Run("period record initialization", func(t *testing.T) {
		// Clear the batchMeterer cache to ensure a clean state
		f.resetBatchMeterer()

		// Create a reservation that spans multiple periods
		reservations := f.createTestReservationMap(100, time.Unix(0, 0), time.Unix(farFutureTimestamp, 0), f.quorum0)

		f.mockState.On("GetMinNumSymbols", f.quorum0).Return(uint64(32)).Maybe()
		f.mockState.On("GetReservationWindow", f.quorum0).Return(uint64(reservationInterval)).Maybe()
		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil)

		// Test initialization of period records
		now := time.Now()
		currentPeriod := meterer.GetReservationPeriod(now.Unix(), f.reservationWindow)

		// Create an update to trigger period record initialization
		updates := []*UpdateRecord{
			newUpdateRecord(f.account1, f.quorum0, currentPeriod, 100),
		}

		// Process the request to initialize the period record
		params, _ := f.mockState.GetPaymentGlobalParams()

		err := f.batchMeterer.processBatch(f.ctx, params, updates, now)
		require.NoError(t, err)

		// Verify period record was initialized correctly
		accountUsage := f.batchMeterer.getOrCreateAccountUsage(f.account1)
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()

		relativeIndex := f.getRelativeIndex(currentPeriod)
		require.NotNil(t, accountUsage.PeriodRecords[f.quorum0], "Period records for quorum should exist")
		require.NotNil(t, accountUsage.PeriodRecords[f.quorum0][relativeIndex], "Period record should exist")
		assert.Equal(t, uint32(currentPeriod), accountUsage.PeriodRecords[f.quorum0][relativeIndex].Index)
		assert.Equal(t, uint64(100), accountUsage.PeriodRecords[f.quorum0][relativeIndex].Usage)
	})

	t.Run("period record cleanup", func(t *testing.T) {
		// Clear the batchMeterer cache to ensure a clean state
		f.resetBatchMeterer()

		// Create a reservation
		reservations := f.createTestReservationMap(100, time.Unix(0, 0), time.Unix(farFutureTimestamp, 0), f.quorum0)

		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil).Maybe()

		// Initialize period records for multiple periods
		now := time.Now()
		for i := 0; i < 4; i++ {
			periodTime := now.Add(time.Duration(i*reservationInterval) * time.Second)
			periodIndex := meterer.GetReservationPeriod(periodTime.Unix(), f.reservationWindow)
			updates := []*UpdateRecord{
				newUpdateRecord(f.account1, f.quorum0, periodIndex, 100),
			}
			params, _ := f.mockState.GetPaymentGlobalParams()

			err := f.batchMeterer.processBatch(f.ctx, params, updates, periodTime)
			require.NoError(t, err)
		}

		// Verify old records are cleaned up
		accountUsage := f.batchMeterer.getOrCreateAccountUsage(f.account1)
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()

		records := accountUsage.PeriodRecords[f.quorum0]
		assert.Equal(t, 3, len(records), "Should only have 3 records (buffer size)")
	})
}

// TestBatchMeterReservationEdgeCases tests edge cases in reservation validation
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
		updates := []*UpdateRecord{
			newUpdateRecord(f.account1, f.quorum0, currentPeriod, 100),
		}

		params, _ := f.mockState.GetPaymentGlobalParams()

		err := f.batchMeterer.processBatch(f.ctx, params, updates, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no reservation")
	})

	t.Run("invalid reservation period", func(t *testing.T) {
		// Create a reservation
		reservations := f.createTestReservationMap(100, time.Unix(0, 0), time.Unix(farFutureTimestamp, 0), f.quorum0)

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
		updates := []*UpdateRecord{
			newUpdateRecord(f.account1, f.quorum0, currentPeriod+testPeriodOffset, 100),
		}

		params, _ := f.mockState.GetPaymentGlobalParams()

		err := f.batchMeterer.processBatch(f.ctx, params, updates, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "RESERVATION_PERIOD_INVALID")
	})

	t.Run("just started reservation", func(t *testing.T) {
		// Clear existing expectations
		f.mockState.ExpectedCalls = nil
		f.mockState.Calls = nil

		// Set up basic mock expectations
		mockParams := &meterer.PaymentVaultParams{
			QuorumProtocolConfigs: make(map[core.QuorumID]*core.PaymentQuorumProtocolConfig),
		}
		mockParams.QuorumProtocolConfigs[f.quorum0] = &core.PaymentQuorumProtocolConfig{
			MinNumSymbols:              defaultMinSymbols,
			ReservationRateLimitWindow: reservationInterval,
		}
		f.mockState.On("GetPaymentGlobalParams").Return(mockParams, nil)

		// Create a reservation that just started
		now := time.Now()
		reservations := f.createTestReservationMap(100, now, now.Add(24*time.Hour), f.quorum0)
		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil).Once()

		currentPeriod := meterer.GetReservationPeriod(now.Unix(), f.reservationWindow)
		f.setupPeriodRecord(f.account1, f.quorum0, currentPeriod)
		updates := []*UpdateRecord{
			newUpdateRecord(f.account1, f.quorum0, currentPeriod, 100),
		}

		params, _ := f.mockState.GetPaymentGlobalParams()

		err := f.batchMeterer.processBatch(f.ctx, params, updates, now)
		require.NoError(t, err)
	})

	t.Run("zero symbols per second", func(t *testing.T) {
		// Clear existing expectations
		f.mockState.ExpectedCalls = nil
		f.mockState.Calls = nil

		// Set up basic mock expectations
		mockParams := &meterer.PaymentVaultParams{
			QuorumProtocolConfigs: make(map[core.QuorumID]*core.PaymentQuorumProtocolConfig),
		}
		mockParams.QuorumProtocolConfigs[f.quorum0] = &core.PaymentQuorumProtocolConfig{
			MinNumSymbols:              defaultMinSymbols,
			ReservationRateLimitWindow: reservationInterval,
		}
		f.mockState.On("GetPaymentGlobalParams").Return(mockParams, nil)

		// Create a reservation with zero symbols per second
		now := time.Now()
		reservations := f.createTestReservationMap(0, now, now.Add(24*time.Hour), f.quorum0)
		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil).Once()

		currentPeriod := meterer.GetReservationPeriod(now.Unix(), f.reservationWindow)
		f.setupPeriodRecord(f.account1, f.quorum0, currentPeriod)
		updates := []*UpdateRecord{
			newUpdateRecord(f.account1, f.quorum0, currentPeriod, 100),
		}

		params, _ := f.mockState.GetPaymentGlobalParams()

		err := f.batchMeterer.processBatch(f.ctx, params, updates, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "BIN_ALREADY_FULL")
	})

	t.Run("invalid timestamps", func(t *testing.T) {
		// Create a reservation with invalid timestamps
		now := time.Now()
		reservations := f.createTestReservationMap(100, now.Add(-24*time.Hour), now, f.quorum0)

		f.mockState.On("GetMinNumSymbols", f.quorum0).Return(uint64(32)).Maybe()
		f.mockState.On("GetReservationWindow", f.quorum0).Return(uint64(reservationInterval)).Maybe()
		f.mockState.On("GetReservedPaymentByAccountAndQuorums",
			f.ctx,
			f.account1,
			[]core.QuorumID{f.quorum0},
		).Return(reservations, nil).Twice()

		currentPeriod := meterer.GetReservationPeriod(now.Unix(), f.reservationWindow)
		f.setupPeriodRecord(f.account1, f.quorum0, currentPeriod)
		updates := []*UpdateRecord{
			newUpdateRecord(f.account1, f.quorum0, currentPeriod, 100),
		}

		params, _ := f.mockState.GetPaymentGlobalParams()

		err := f.batchMeterer.processBatch(f.ctx, params, updates, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "RESERVATION_PERIOD_INVALID")
	})
}

// TestBatchMeterRollback tests the rollback functionality in processBatch
func TestBatchMeterRollback(t *testing.T) {
	f := setupTestFixtures(t)

	t.Run("rollback on reservation error", func(t *testing.T) {
		// Clear the batchMeterer cache to ensure a clean state
		f.resetBatchMeterer()

		// Create a reservation
		reservations := f.createTestReservationMap(100, time.Unix(0, 0), time.Unix(farFutureTimestamp, 0), f.quorum0)

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
		updates := []*UpdateRecord{
			newUpdateRecord(f.account1, f.quorum0, currentPeriod, 50),
		}

		// Process the request
		params, _ := f.mockState.GetPaymentGlobalParams()
		err := f.batchMeterer.processBatch(f.ctx, params, updates, now)
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
		updates = []*UpdateRecord{
			newUpdateRecord(f.account1, f.quorum0, currentPeriod, 100),
		}

		// Process the request
		params, _ = f.mockState.GetPaymentGlobalParams()
		err = f.batchMeterer.processBatch(f.ctx, params, updates, now)
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
		reservations := f.createTestReservationMap(100, time.Unix(0, 0), time.Unix(farFutureTimestamp, 0), f.quorum0)

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
		updates := []*UpdateRecord{
			newUpdateRecord(f.account1, f.quorum0, currentPeriod, 50),
		}

		// Process the request
		params, _ := f.mockState.GetPaymentGlobalParams()
		err := f.batchMeterer.processBatch(f.ctx, params, updates, now)
		require.NoError(t, err)

		// Create updates for invalid period (far future period)
		updates = []*UpdateRecord{
			newUpdateRecord(f.account1, f.quorum0, currentPeriod+testPeriodMultiplier*reservationInterval, 100),
		}

		// Process the request
		params, _ = f.mockState.GetPaymentGlobalParams()
		err = f.batchMeterer.processBatch(f.ctx, params, updates, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "RESERVATION_PERIOD_INVALID")

		// Verify usage was rolled back
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()
		assert.Equal(t, uint64(50), accountUsage.PeriodRecords[f.quorum0][relativeIndex].Usage)
	})
}
