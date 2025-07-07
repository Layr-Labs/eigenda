package meterer

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"testing"
	"time"

	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer/payment_logic"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	testReservationWindow = 3600 // 1 hour
	testUpdateTimeout     = 1 * time.Second
	testUpdateInterval    = 1 * time.Minute
	testMinSymbols        = 32
	testSymbolsPerSecond  = 100
	testFarFutureTime     = 9999999999
)

// MockOnchainPaymentState provides a mock implementation of OnchainPayment
type MockOnchainPaymentState struct {
	mock.Mock
}

func (m *MockOnchainPaymentState) RefreshOnchainPaymentState(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockOnchainPaymentState) GetPaymentGlobalParams() (*PaymentVaultParams, error) {
	args := m.Called()
	return args.Get(0).(*PaymentVaultParams), args.Error(1)
}

func (m *MockOnchainPaymentState) GetReservedPaymentByAccountAndQuorums(ctx context.Context, accountID common.Address, quorumNumbers []core.QuorumID) (map[core.QuorumID]*core.ReservedPayment, error) {
	args := m.Called(ctx, accountID, quorumNumbers)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[core.QuorumID]*core.ReservedPayment), args.Error(1)
}

func (m *MockOnchainPaymentState) GetOnDemandPaymentByAccount(ctx context.Context, accountID common.Address) (*core.OnDemandPayment, error) {
	args := m.Called(ctx, accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*core.OnDemandPayment), args.Error(1)
}

// batchLedgerTestFixtures contains common test data and configurations
type batchLedgerTestFixtures struct {
	ctx              context.Context
	logger           logging.Logger
	mockPaymentState *MockOnchainPaymentState
	batchLedger      *BatchLedger
	config           Config
	params           *PaymentVaultParams
	account1         common.Address
	account2         common.Address
	quorum0          core.QuorumID
	quorum1          core.QuorumID
}

// setupBatchLedgerTestFixtures creates and returns a new batchLedgerTestFixtures instance
func setupBatchLedgerTestFixtures(_ *testing.T) *batchLedgerTestFixtures {
	ctx := context.Background()
	logger := testutils.GetLogger()
	mockPaymentState := new(MockOnchainPaymentState)

	config := Config{
		ChainReadTimeout: testUpdateTimeout,
		UpdateInterval:   testUpdateInterval,
	}

	// Set up payment vault params - only reservation quorums (no on-demand)
	params := &PaymentVaultParams{
		QuorumProtocolConfigs: make(map[core.QuorumID]*core.PaymentQuorumProtocolConfig),
		QuorumPaymentConfigs:  make(map[core.QuorumID]*core.PaymentQuorumConfig),
	}
	for i := uint8(0); i < 2; i++ {
		params.QuorumProtocolConfigs[core.QuorumID(i)] = &core.PaymentQuorumProtocolConfig{
			MinNumSymbols:              testMinSymbols,
			ReservationRateLimitWindow: testReservationWindow,
		}
		params.QuorumPaymentConfigs[core.QuorumID(i)] = &core.PaymentQuorumConfig{
			OnDemandSymbolsPerSecond: testSymbolsPerSecond,
			OnDemandPricePerSymbol:   uint64(1),
		}
	}

	mockPaymentState.On("GetPaymentGlobalParams").Return(params, nil)

	// Create batch ledger (no server ledger needed)
	batchLedger := NewBatchLedger(mockPaymentState, config, logger)

	return &batchLedgerTestFixtures{
		ctx:              ctx,
		logger:           logger,
		mockPaymentState: mockPaymentState,
		batchLedger:      batchLedger,
		config:           config,
		params:           params,
		account1:         common.HexToAddress("0x1"),
		account2:         common.HexToAddress("0x2"),
		quorum0:          core.QuorumID(0),
		quorum1:          core.QuorumID(1),
	}
}

// createTestReservation creates a test reservation
func createTestReservation(symbolsPerSecond uint64, startTime, endTime time.Time) *core.ReservedPayment {
	return &core.ReservedPayment{
		SymbolsPerSecond: symbolsPerSecond,
		StartTimestamp:   uint64(startTime.Unix()),
		EndTimestamp:     uint64(endTime.Unix()),
	}
}

// createTestBatch creates a test batch with the given account information
// Replicates the structure from BatchMeterer tests
func createTestBatch(accountID common.Address, quorumID core.QuorumID, timestampNs int64, numSymbols uint64) *corev2.Batch {
	return &corev2.Batch{
		BatchHeader: &corev2.BatchHeader{
			BatchRoot:            [32]byte{1, 2, 3},
			ReferenceBlockNumber: 100,
		},
		BlobCertificates: []*corev2.BlobCertificate{
			{
				BlobHeader: &corev2.BlobHeader{
					BlobVersion:     0,
					BlobCommitments: encoding.BlobCommitments{Length: uint(numSymbols)},
					QuorumNumbers:   []core.QuorumID{quorumID},
					PaymentMetadata: core.PaymentMetadata{
						AccountID:         accountID,
						Timestamp:         timestampNs,
						CumulativePayment: big.NewInt(0),
					},
				},
				Signature: []byte{1, 2, 3},
				RelayKeys: []corev2.RelayKey{0, 1},
			},
		},
	}
}

// createTestBatchMultiple creates a test batch with multiple blob certificates
func createTestBatchMultiple(infos []struct {
	accountID   common.Address
	quorumID    core.QuorumID
	timestampNs int64
	numSymbols  uint64
}) *corev2.Batch {
	batch := &corev2.Batch{
		BatchHeader: &corev2.BatchHeader{
			BatchRoot:            [32]byte{1, 2, 3},
			ReferenceBlockNumber: 100,
		},
		BlobCertificates: make([]*corev2.BlobCertificate, len(infos)),
	}

	for i, info := range infos {
		batch.BlobCertificates[i] = &corev2.BlobCertificate{
			BlobHeader: &corev2.BlobHeader{
				BlobVersion:     0,
				BlobCommitments: encoding.BlobCommitments{Length: uint(info.numSymbols)},
				QuorumNumbers:   []core.QuorumID{info.quorumID},
				PaymentMetadata: core.PaymentMetadata{
					AccountID:         info.accountID,
					Timestamp:         info.timestampNs,
					CumulativePayment: big.NewInt(0),
				},
			},
			Signature: []byte{1, 2, 3},
			RelayKeys: []corev2.RelayKey{0, 1},
		}
	}

	return batch
}

func TestBatchLedger_MeterBatch_Reservations(t *testing.T) {
	f := setupBatchLedgerTestFixtures(t)

	t.Run("successful reservation batch", func(t *testing.T) {
		// Set up reservations for both accounts
		reservations1 := map[core.QuorumID]*core.ReservedPayment{
			f.quorum0: createTestReservation(testSymbolsPerSecond, time.Unix(0, 0), time.Unix(testFarFutureTime, 0)),
		}
		reservations2 := map[core.QuorumID]*core.ReservedPayment{
			f.quorum1: createTestReservation(testSymbolsPerSecond, time.Unix(0, 0), time.Unix(testFarFutureTime, 0)),
		}

		f.mockPaymentState.On("GetReservedPaymentByAccountAndQuorums", f.ctx, f.account1, mock.Anything).Return(reservations1, nil)
		f.mockPaymentState.On("GetReservedPaymentByAccountAndQuorums", f.ctx, f.account2, mock.Anything).Return(reservations2, nil)

		// Create batch with multiple blob certificates
		now := time.Now()
		batch := createTestBatchMultiple([]struct {
			accountID   common.Address
			quorumID    core.QuorumID
			timestampNs int64
			numSymbols  uint64
		}{
			{
				accountID:   f.account1,
				quorumID:    f.quorum0,
				timestampNs: now.UnixNano(),
				numSymbols:  64,
			},
			{
				accountID:   f.account2,
				quorumID:    f.quorum1,
				timestampNs: now.UnixNano(),
				numSymbols:  96,
			},
		})

		// Process batch
		err := f.batchLedger.MeterBatch(f.ctx, batch, now)
		require.NoError(t, err)

		// Verify usage was tracked for both accounts
		accountUsage1 := f.batchLedger.getAccount(f.account1)
		assert.NotNil(t, accountUsage1)
		accountUsage2 := f.batchLedger.getAccount(f.account2)
		assert.NotNil(t, accountUsage2)
	})

	t.Run("atomic failure prevents all changes", func(t *testing.T) {
		// Clear previous state
		f.batchLedger.accounts = sync.Map{}
		f.mockPaymentState.ExpectedCalls = nil

		// Set up reservations - account1 has reservation, account2 doesn't
		reservations1 := map[core.QuorumID]*core.ReservedPayment{
			f.quorum0: createTestReservation(testSymbolsPerSecond, time.Unix(0, 0), time.Unix(testFarFutureTime, 0)),
		}
		emptyReservations := make(map[core.QuorumID]*core.ReservedPayment)

		f.mockPaymentState.On("GetPaymentGlobalParams").Return(f.params, nil)
		f.mockPaymentState.On("GetReservedPaymentByAccountAndQuorums", f.ctx, f.account1, mock.Anything).Return(reservations1, nil)
		f.mockPaymentState.On("GetReservedPaymentByAccountAndQuorums", f.ctx, f.account2, mock.Anything).Return(emptyReservations, nil)

		// Create batch with multiple blob certificates - should fail on second account
		now := time.Now()
		batch := createTestBatchMultiple([]struct {
			accountID   common.Address
			quorumID    core.QuorumID
			timestampNs int64
			numSymbols  uint64
		}{
			{
				accountID:   f.account1,
				quorumID:    f.quorum0,
				timestampNs: now.UnixNano(),
				numSymbols:  64,
			},
			{
				accountID:   f.account2,
				quorumID:    f.quorum1,
				timestampNs: now.UnixNano(),
				numSymbols:  96,
			},
		})

		// Process batch - should fail on second account
		err := f.batchLedger.MeterBatch(f.ctx, batch, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no reservation")
	})
}

func TestBatchLedger_ProcessBatch_EdgeCases(t *testing.T) {
	f := setupBatchLedgerTestFixtures(t)

	t.Run("empty batch", func(t *testing.T) {
		now := time.Now()
		batch := &corev2.Batch{}

		err := f.batchLedger.MeterBatch(f.ctx, batch, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "batch is nil or empty")
	})

	t.Run("nil batch", func(t *testing.T) {
		now := time.Now()
		var batch *corev2.Batch

		err := f.batchLedger.MeterBatch(f.ctx, batch, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "batch is nil or empty")
	})

	t.Run("reservation lookup failure", func(t *testing.T) {
		// Clear previous state
		f.batchLedger.accounts = sync.Map{}
		f.mockPaymentState.ExpectedCalls = nil

		// Mock failure in getting reservations
		f.mockPaymentState.On("GetPaymentGlobalParams").Return(f.params, nil)
		f.mockPaymentState.On("GetReservedPaymentByAccountAndQuorums", f.ctx, f.account1, mock.Anything).Return(nil, errors.New("chain error"))

		now := time.Now()
		batch := createTestBatch(f.account1, f.quorum0, now.UnixNano(), 64)

		err := f.batchLedger.MeterBatch(f.ctx, batch, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get reservations")
	})
}

func TestBatchLedger_UsageTracking(t *testing.T) {
	f := setupBatchLedgerTestFixtures(t)

	// Set up reservation
	reservations := map[core.QuorumID]*core.ReservedPayment{
		f.quorum0: createTestReservation(testSymbolsPerSecond, time.Unix(0, 0), time.Unix(testFarFutureTime, 0)),
	}

	f.mockPaymentState.On("GetReservedPaymentByAccountAndQuorums", f.ctx, f.account1, mock.Anything).Return(reservations, nil)

	now := time.Now()
	batch := createTestBatch(f.account1, f.quorum0, now.UnixNano(), 64)

	// Process batch
	err := f.batchLedger.MeterBatch(f.ctx, batch, now)
	require.NoError(t, err)

	// Verify usage tracking
	accountUsage := f.batchLedger.getAccount(f.account1)
	require.NotNil(t, accountUsage)

	// Verify period records were created
	assert.NotNil(t, accountUsage.PeriodRecords[f.quorum0])
}

// TestBatchLedger_OverflowEdgeCases tests overflow handling like BatchMeterer
func TestBatchLedger_OverflowEdgeCases(t *testing.T) {
	f := setupBatchLedgerTestFixtures(t)

	t.Run("exact bin limit", func(t *testing.T) {
		// Set up reservation with specific rate
		symbolsPerSecond := uint64(100)
		reservations := map[core.QuorumID]*core.ReservedPayment{
			f.quorum0: createTestReservation(symbolsPerSecond, time.Unix(0, 0), time.Unix(testFarFutureTime, 0)),
		}
		f.mockPaymentState.On("GetReservedPaymentByAccountAndQuorums", f.ctx, f.account1, mock.Anything).Return(reservations, nil)

		// Calculate exact bin limit
		binLimit := symbolsPerSecond * testReservationWindow
		now := time.Now()
		batch := createTestBatch(f.account1, f.quorum0, now.UnixNano(), binLimit)

		err := f.batchLedger.MeterBatch(f.ctx, batch, now)
		require.NoError(t, err)

		// Verify usage equals bin limit
		accountUsage := f.batchLedger.getAccount(f.account1)
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()

		// Find the period record that was actually used
		var periodRecord *pb.PeriodRecord
		for _, record := range accountUsage.PeriodRecords[f.quorum0] {
			if record != nil && record.Usage > 0 {
				periodRecord = record
				break
			}
		}
		require.NotNil(t, periodRecord, "should have found a period record with usage")
		assert.Equal(t, binLimit, periodRecord.Usage)
	})

	t.Run("just over bin limit", func(t *testing.T) {
		f.batchLedger.accounts = sync.Map{}
		f.mockPaymentState.ExpectedCalls = nil

		symbolsPerSecond := uint64(100)
		reservations := map[core.QuorumID]*core.ReservedPayment{
			f.quorum0: createTestReservation(symbolsPerSecond, time.Unix(0, 0), time.Unix(testFarFutureTime, 0)),
		}
		f.mockPaymentState.On("GetPaymentGlobalParams").Return(f.params, nil)
		f.mockPaymentState.On("GetReservedPaymentByAccountAndQuorums", f.ctx, f.account1, mock.Anything).Return(reservations, nil)

		// Try usage just over bin limit (should overflow to next period)
		binLimit := symbolsPerSecond * testReservationWindow
		slightlyOver := binLimit + 10
		now := time.Now()
		batch := createTestBatch(f.account1, f.quorum0, now.UnixNano(), slightlyOver)

		err := f.batchLedger.MeterBatch(f.ctx, batch, now)
		require.NoError(t, err)

		// Verify current period capped at bin limit and overflow period has remaining usage
		accountUsage := f.batchLedger.getAccount(f.account1)
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()

		// Find the period records that were actually used
		usedRecords := make([]*pb.PeriodRecord, 0)
		for _, record := range accountUsage.PeriodRecords[f.quorum0] {
			if record != nil && record.Usage > 0 {
				usedRecords = append(usedRecords, record)
			}
		}
		require.Len(t, usedRecords, 2, "should have 2 period records with usage (current + overflow)")

		// Calculate the actual charged amount (rounds up to nearest multiple of minSymbols)
		actualCharged := payment_logic.SymbolsCharged(slightlyOver, testMinSymbols)
		expectedOverflow := actualCharged - binLimit

		// Verify one period is capped at bin limit and the other has the overflow
		foundBinLimit := false
		foundOverflow := false
		for _, record := range usedRecords {
			if record.Usage == binLimit {
				foundBinLimit = true
			} else if record.Usage == expectedOverflow {
				foundOverflow = true
			}
		}
		assert.True(t, foundBinLimit, "should have one period at bin limit")
		assert.True(t, foundOverflow, "should have one period with overflow amount")
	})

	t.Run("exactly 2x bin limit", func(t *testing.T) {
		f.batchLedger.accounts = sync.Map{}
		f.mockPaymentState.ExpectedCalls = nil

		symbolsPerSecond := uint64(100)
		reservations := map[core.QuorumID]*core.ReservedPayment{
			f.quorum0: createTestReservation(symbolsPerSecond, time.Unix(0, 0), time.Unix(testFarFutureTime, 0)),
		}
		f.mockPaymentState.On("GetPaymentGlobalParams").Return(f.params, nil)
		f.mockPaymentState.On("GetReservedPaymentByAccountAndQuorums", f.ctx, f.account1, mock.Anything).Return(reservations, nil)

		// Try usage exactly 2x bin limit
		binLimit := symbolsPerSecond * testReservationWindow
		exactly2x := 2 * binLimit
		now := time.Now()
		batch := createTestBatch(f.account1, f.quorum0, now.UnixNano(), exactly2x)

		err := f.batchLedger.MeterBatch(f.ctx, batch, now)
		require.NoError(t, err)

		// Verify both current and overflow periods are at limit
		accountUsage := f.batchLedger.getAccount(f.account1)
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()

		// Find the period records that were actually used
		usedRecords := make([]*pb.PeriodRecord, 0)
		for _, record := range accountUsage.PeriodRecords[f.quorum0] {
			if record != nil && record.Usage > 0 {
				usedRecords = append(usedRecords, record)
			}
		}
		require.Len(t, usedRecords, 2, "should have 2 period records with usage (current + overflow)")

		// Verify both periods are at bin limit (2x bin limit split across 2 periods)
		for _, record := range usedRecords {
			assert.Equal(t, binLimit, record.Usage, "each period should be at bin limit")
		}

		// Verify total usage equals exactly 2x
		totalUsage := uint64(0)
		for _, record := range usedRecords {
			totalUsage += record.Usage
		}
		assert.Equal(t, exactly2x, totalUsage)
	})

	t.Run("over 2x bin limit fails", func(t *testing.T) {
		f.batchLedger.accounts = sync.Map{}
		f.mockPaymentState.ExpectedCalls = nil

		symbolsPerSecond := uint64(100)
		reservations := map[core.QuorumID]*core.ReservedPayment{
			f.quorum0: createTestReservation(symbolsPerSecond, time.Unix(0, 0), time.Unix(testFarFutureTime, 0)),
		}
		f.mockPaymentState.On("GetPaymentGlobalParams").Return(f.params, nil)
		f.mockPaymentState.On("GetReservedPaymentByAccountAndQuorums", f.ctx, f.account1, mock.Anything).Return(reservations, nil)

		// Try usage over 2x bin limit
		binLimit := symbolsPerSecond * testReservationWindow
		excessiveUsage := 2*binLimit + 1
		now := time.Now()
		batch := createTestBatch(f.account1, f.quorum0, now.UnixNano(), excessiveUsage)

		err := f.batchLedger.MeterBatch(f.ctx, batch, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "overflow usage exceeds bin limit")
	})
}

// TestBatchLedger_PeriodRecordEdgeCases tests period record management like BatchMeterer
func TestBatchLedger_PeriodRecordEdgeCases(t *testing.T) {
	f := setupBatchLedgerTestFixtures(t)

	t.Run("period record initialization", func(t *testing.T) {
		reservations := map[core.QuorumID]*core.ReservedPayment{
			f.quorum0: createTestReservation(testSymbolsPerSecond, time.Unix(0, 0), time.Unix(testFarFutureTime, 0)),
		}
		f.mockPaymentState.On("GetReservedPaymentByAccountAndQuorums", f.ctx, f.account1, mock.Anything).Return(reservations, nil)

		now := time.Now()
		batch := createTestBatch(f.account1, f.quorum0, now.UnixNano(), 64)

		err := f.batchLedger.MeterBatch(f.ctx, batch, now)
		require.NoError(t, err)

		// Verify period record was created correctly
		accountUsage := f.batchLedger.getAccount(f.account1)
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()
		assert.NotNil(t, accountUsage.PeriodRecords[f.quorum0])

		// Find the period record that was actually used
		var periodRecord *pb.PeriodRecord
		for _, record := range accountUsage.PeriodRecords[f.quorum0] {
			if record != nil && record.Usage > 0 {
				periodRecord = record
				break
			}
		}
		require.NotNil(t, periodRecord, "should have found a period record with usage")
		assert.Equal(t, uint64(64), periodRecord.Usage)
	})

	t.Run("circular buffer wrapping", func(t *testing.T) {
		f.batchLedger.accounts = sync.Map{}
		f.mockPaymentState.ExpectedCalls = nil

		reservations := map[core.QuorumID]*core.ReservedPayment{
			f.quorum0: createTestReservation(testSymbolsPerSecond, time.Unix(0, 0), time.Unix(testFarFutureTime, 0)),
		}
		f.mockPaymentState.On("GetPaymentGlobalParams").Return(f.params, nil).Times(4)
		f.mockPaymentState.On("GetReservedPaymentByAccountAndQuorums", f.ctx, f.account1, mock.Anything).Return(reservations, nil).Times(8) // 2 calls per MeterBatch (validation + commit)

		// Process 4 periods to test circular buffer behavior (MinNumBins = 3)
		baseTime := time.Now()
		for i := 0; i < 4; i++ {
			timestamp := baseTime.Add(time.Duration(i*int(testReservationWindow)) * time.Second)
			batch := createTestBatch(f.account1, f.quorum0, timestamp.UnixNano(), 32)
			err := f.batchLedger.MeterBatch(f.ctx, batch, timestamp)
			require.NoError(t, err)
		}

		// Verify only 3 records remain (circular buffer size)
		accountUsage := f.batchLedger.getAccount(f.account1)
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()
		nonNilCount := 0
		for _, record := range accountUsage.PeriodRecords[f.quorum0] {
			if record != nil {
				nonNilCount++
			}
		}
		assert.Equal(t, 3, nonNilCount) // MinNumBins
	})
}

// TestBatchLedger_ReservationEdgeCases tests reservation validation like BatchMeterer
func TestBatchLedger_ReservationEdgeCases(t *testing.T) {
	f := setupBatchLedgerTestFixtures(t)

	t.Run("inactive reservation", func(t *testing.T) {
		// Create reservation that's not active yet
		futureStart := time.Now().Add(24 * time.Hour)
		reservations := map[core.QuorumID]*core.ReservedPayment{
			f.quorum0: createTestReservation(testSymbolsPerSecond, futureStart, time.Unix(testFarFutureTime, 0)),
		}
		f.mockPaymentState.On("GetPaymentGlobalParams").Return(f.params, nil)
		f.mockPaymentState.On("GetReservedPaymentByAccountAndQuorums", f.ctx, f.account1, mock.Anything).Return(reservations, nil)

		now := time.Now()
		batch := createTestBatch(f.account1, f.quorum0, now.UnixNano(), 64)

		err := f.batchLedger.MeterBatch(f.ctx, batch, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "inactive reservation")
	})

	t.Run("expired reservation", func(t *testing.T) {
		f.batchLedger.accounts = sync.Map{}
		f.mockPaymentState.ExpectedCalls = nil

		// Create reservation that has expired
		pastEnd := time.Now().Add(-24 * time.Hour)
		reservations := map[core.QuorumID]*core.ReservedPayment{
			f.quorum0: createTestReservation(testSymbolsPerSecond, time.Unix(0, 0), pastEnd),
		}
		f.mockPaymentState.On("GetPaymentGlobalParams").Return(f.params, nil)
		f.mockPaymentState.On("GetReservedPaymentByAccountAndQuorums", f.ctx, f.account1, mock.Anything).Return(reservations, nil)

		now := time.Now()
		batch := createTestBatch(f.account1, f.quorum0, now.UnixNano(), 64)

		err := f.batchLedger.MeterBatch(f.ctx, batch, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "inactive reservation")
	})

	t.Run("zero symbols per second", func(t *testing.T) {
		f.batchLedger.accounts = sync.Map{}
		f.mockPaymentState.ExpectedCalls = nil

		// Create reservation with zero symbols per second
		reservations := map[core.QuorumID]*core.ReservedPayment{
			f.quorum0: createTestReservation(0, time.Unix(0, 0), time.Unix(testFarFutureTime, 0)),
		}
		f.mockPaymentState.On("GetPaymentGlobalParams").Return(f.params, nil)
		f.mockPaymentState.On("GetReservedPaymentByAccountAndQuorums", f.ctx, f.account1, mock.Anything).Return(reservations, nil)

		now := time.Now()
		batch := createTestBatch(f.account1, f.quorum0, now.UnixNano(), 64)

		err := f.batchLedger.MeterBatch(f.ctx, batch, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("invalid reservation period", func(t *testing.T) {
		f.batchLedger.accounts = sync.Map{}
		f.mockPaymentState.ExpectedCalls = nil

		// Create reservation that's active for current time
		now := time.Now()
		reservations := map[core.QuorumID]*core.ReservedPayment{
			f.quorum0: createTestReservation(testSymbolsPerSecond, now.Add(-1*time.Hour), now.Add(1*time.Hour)),
		}
		f.mockPaymentState.On("GetPaymentGlobalParams").Return(f.params, nil)
		f.mockPaymentState.On("GetReservedPaymentByAccountAndQuorums", f.ctx, f.account1, mock.Anything).Return(reservations, nil)

		// Create a timestamp that's too far in the past to be a valid period
		// This should create a period that's more than 1 reservation window ago
		invalidPastTime := now.Add(-time.Duration(testReservationWindow*2) * time.Second)
		batch := createTestBatch(f.account1, f.quorum0, invalidPastTime.UnixNano(), 64)

		// Use current time as batchReceivedAt so reservation is active
		err := f.batchLedger.MeterBatch(f.ctx, batch, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid reservation period")
	})

	t.Run("just started reservation", func(t *testing.T) {
		f.batchLedger.accounts = sync.Map{}
		f.mockPaymentState.ExpectedCalls = nil

		// Create reservation that just started
		now := time.Now()
		reservations := map[core.QuorumID]*core.ReservedPayment{
			f.quorum0: createTestReservation(testSymbolsPerSecond, now.Add(-1*time.Second), time.Unix(testFarFutureTime, 0)),
		}
		f.mockPaymentState.On("GetPaymentGlobalParams").Return(f.params, nil)
		f.mockPaymentState.On("GetReservedPaymentByAccountAndQuorums", f.ctx, f.account1, mock.Anything).Return(reservations, nil)

		batch := createTestBatch(f.account1, f.quorum0, now.UnixNano(), 64)

		err := f.batchLedger.MeterBatch(f.ctx, batch, now)
		require.NoError(t, err)

		// Verify successful processing
		accountUsage := f.batchLedger.getAccount(f.account1)
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()

		// Find the period record that was actually used (non-nil with usage > 0)
		var periodRecord *pb.PeriodRecord
		for _, record := range accountUsage.PeriodRecords[f.quorum0] {
			if record != nil && record.Usage > 0 {
				periodRecord = record
				break
			}
		}
		require.NotNil(t, periodRecord, "should have found a period record with usage")
		assert.Equal(t, uint64(64), periodRecord.Usage)
	})
}

// TestBatchLedger_BatchUsageEdgeCases tests batch processing edge cases like BatchMeterer
func TestBatchLedger_BatchUsageEdgeCases(t *testing.T) {
	f := setupBatchLedgerTestFixtures(t)

	t.Run("nil blob header", func(t *testing.T) {
		batch := &corev2.Batch{
			BatchHeader: &corev2.BatchHeader{
				BatchRoot:            [32]byte{1, 2, 3},
				ReferenceBlockNumber: 100,
			},
			BlobCertificates: []*corev2.BlobCertificate{
				{
					BlobHeader: nil, // Nil header
					Signature:  []byte{1, 2, 3},
				},
			},
		}

		now := time.Now()
		err := f.batchLedger.MeterBatch(f.ctx, batch, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "blob certificate has nil header")
	})

	t.Run("zero length blob", func(t *testing.T) {
		reservations := map[core.QuorumID]*core.ReservedPayment{
			f.quorum0: createTestReservation(testSymbolsPerSecond, time.Unix(0, 0), time.Unix(testFarFutureTime, 0)),
		}
		f.mockPaymentState.On("GetReservedPaymentByAccountAndQuorums", f.ctx, f.account1, mock.Anything).Return(reservations, nil)

		now := time.Now()
		batch := createTestBatch(f.account1, f.quorum0, now.UnixNano(), 0) // Zero length

		err := f.batchLedger.MeterBatch(f.ctx, batch, now)
		require.NoError(t, err)

		// Verify usage rounds up to MinNumSymbols
		accountUsage := f.batchLedger.getAccount(f.account1)
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()

		// Find the period record that was actually used
		var periodRecord *pb.PeriodRecord
		for _, record := range accountUsage.PeriodRecords[f.quorum0] {
			if record != nil && record.Usage > 0 {
				periodRecord = record
				break
			}
		}
		require.NotNil(t, periodRecord, "should have found a period record with usage")
		assert.Equal(t, uint64(testMinSymbols), periodRecord.Usage) // Should round up to 32
	})

	t.Run("duplicate account quorum in batch", func(t *testing.T) {
		f.batchLedger.accounts = sync.Map{}
		f.mockPaymentState.ExpectedCalls = nil

		reservations := map[core.QuorumID]*core.ReservedPayment{
			f.quorum0: createTestReservation(testSymbolsPerSecond, time.Unix(0, 0), time.Unix(testFarFutureTime, 0)),
		}
		f.mockPaymentState.On("GetPaymentGlobalParams").Return(f.params, nil)
		f.mockPaymentState.On("GetReservedPaymentByAccountAndQuorums", f.ctx, f.account1, mock.Anything).Return(reservations, nil)

		// Create batch with same account/quorum twice
		now := time.Now()
		batch := createTestBatchMultiple([]struct {
			accountID   common.Address
			quorumID    core.QuorumID
			timestampNs int64
			numSymbols  uint64
		}{
			{
				accountID:   f.account1,
				quorumID:    f.quorum0,
				timestampNs: now.UnixNano(),
				numSymbols:  32,
			},
			{
				accountID:   f.account1,
				quorumID:    f.quorum0,
				timestampNs: now.UnixNano(),
				numSymbols:  32,
			},
		})

		err := f.batchLedger.MeterBatch(f.ctx, batch, now)
		require.NoError(t, err)

		// Verify usage gets summed (32 + 32 = 64)
		accountUsage := f.batchLedger.getAccount(f.account1)
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()

		// Find the period record that was actually used
		var periodRecord *pb.PeriodRecord
		for _, record := range accountUsage.PeriodRecords[f.quorum0] {
			if record != nil && record.Usage > 0 {
				periodRecord = record
				break
			}
		}
		require.NotNil(t, periodRecord, "should have found a period record with usage")
		assert.Equal(t, uint64(64), periodRecord.Usage)
	})

	t.Run("very large usage", func(t *testing.T) {
		f.batchLedger.accounts = sync.Map{}
		f.mockPaymentState.ExpectedCalls = nil

		// Create large reservation to handle big usage
		largeSymbolsPerSecond := uint64(1000000)
		reservations := map[core.QuorumID]*core.ReservedPayment{
			f.quorum0: createTestReservation(largeSymbolsPerSecond, time.Unix(0, 0), time.Unix(testFarFutureTime, 0)),
		}
		f.mockPaymentState.On("GetPaymentGlobalParams").Return(f.params, nil)
		f.mockPaymentState.On("GetReservedPaymentByAccountAndQuorums", f.ctx, f.account1, mock.Anything).Return(reservations, nil)

		// Try very large usage (but within bounds)
		largeUsage := uint64(4294967295) // MaxUint32
		now := time.Now()
		batch := createTestBatch(f.account1, f.quorum0, now.UnixNano(), largeUsage)

		err := f.batchLedger.MeterBatch(f.ctx, batch, now)
		require.NoError(t, err)

		// Verify usage was recorded (may be rounded)
		accountUsage := f.batchLedger.getAccount(f.account1)
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()

		// Find the period record that was actually used
		var periodRecord *pb.PeriodRecord
		for _, record := range accountUsage.PeriodRecords[f.quorum0] {
			if record != nil && record.Usage > 0 {
				periodRecord = record
				break
			}
		}
		require.NotNil(t, periodRecord, "should have found a period record with usage")
		assert.True(t, periodRecord.Usage > 0)
		assert.True(t, periodRecord.Usage <= largeSymbolsPerSecond*testReservationWindow)
	})
}

// TestBatchLedger_PeriodTransition tests period transition scenarios like BatchMeterer
func TestBatchLedger_PeriodTransition(t *testing.T) {
	f := setupBatchLedgerTestFixtures(t)

	t.Run("multiple periods in same batch", func(t *testing.T) {
		reservations := map[core.QuorumID]*core.ReservedPayment{
			f.quorum0: createTestReservation(testSymbolsPerSecond, time.Unix(0, 0), time.Unix(testFarFutureTime, 0)),
		}
		f.mockPaymentState.On("GetPaymentGlobalParams").Return(f.params, nil)
		f.mockPaymentState.On("GetReservedPaymentByAccountAndQuorums", f.ctx, f.account1, mock.Anything).Return(reservations, nil)

		// Create batch with timestamps from different but valid periods
		// Use current period and previous period (both are valid)
		baseTime := time.Now()
		currentTime := baseTime
		previousPeriodTime := baseTime.Add(-time.Duration(testReservationWindow) * time.Second)

		batch := createTestBatchMultiple([]struct {
			accountID   common.Address
			quorumID    core.QuorumID
			timestampNs int64
			numSymbols  uint64
		}{
			{
				accountID:   f.account1,
				quorumID:    f.quorum0,
				timestampNs: currentTime.UnixNano(),
				numSymbols:  32,
			},
			{
				accountID:   f.account1,
				quorumID:    f.quorum0,
				timestampNs: previousPeriodTime.UnixNano(),
				numSymbols:  32,
			},
		})

		err := f.batchLedger.MeterBatch(f.ctx, batch, baseTime)
		require.NoError(t, err)

		// Verify usage tracked in both periods
		accountUsage := f.batchLedger.getAccount(f.account1)
		accountUsage.Lock.RLock()
		defer accountUsage.Lock.RUnlock()

		// Should have records for both periods
		assert.NotNil(t, accountUsage.PeriodRecords[f.quorum0])

		// Count non-nil period records
		nonNilCount := 0
		for _, record := range accountUsage.PeriodRecords[f.quorum0] {
			if record != nil && record.Usage > 0 {
				nonNilCount++
			}
		}
		assert.GreaterOrEqual(t, nonNilCount, 1) // At least one period should have usage
	})
}

// TestBatchLedger_AtomicityValidation tests the validate-first atomicity behavior
func TestBatchLedger_AtomicityValidation(t *testing.T) {
	f := setupBatchLedgerTestFixtures(t)

	t.Run("validation failure prevents all changes", func(t *testing.T) {
		// Set up one valid and one invalid account
		reservations1 := map[core.QuorumID]*core.ReservedPayment{
			f.quorum0: createTestReservation(testSymbolsPerSecond, time.Unix(0, 0), time.Unix(testFarFutureTime, 0)),
		}
		emptyReservations := make(map[core.QuorumID]*core.ReservedPayment)

		f.mockPaymentState.On("GetPaymentGlobalParams").Return(f.params, nil)
		f.mockPaymentState.On("GetReservedPaymentByAccountAndQuorums", f.ctx, f.account1, mock.Anything).Return(reservations1, nil)
		f.mockPaymentState.On("GetReservedPaymentByAccountAndQuorums", f.ctx, f.account2, mock.Anything).Return(emptyReservations, nil)

		// Create batch with valid and invalid accounts
		now := time.Now()
		batch := createTestBatchMultiple([]struct {
			accountID   common.Address
			quorumID    core.QuorumID
			timestampNs int64
			numSymbols  uint64
		}{
			{
				accountID:   f.account1,
				quorumID:    f.quorum0,
				timestampNs: now.UnixNano(),
				numSymbols:  64,
			},
			{
				accountID:   f.account2, // No reservation
				quorumID:    f.quorum1,
				timestampNs: now.UnixNano(),
				numSymbols:  96,
			},
		})

		// Should fail due to account2
		err := f.batchLedger.MeterBatch(f.ctx, batch, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no reservation")

		// Verify NO changes were made to account1 (atomicity)
		accountUsage1 := f.batchLedger.getAccount(f.account1)
		accountUsage1.Lock.RLock()
		defer accountUsage1.Lock.RUnlock()

		// No usage should have been applied due to atomicity (validation prevents all changes)
		if accountUsage1.PeriodRecords[f.quorum0] != nil {
			for _, record := range accountUsage1.PeriodRecords[f.quorum0] {
				if record != nil {
					assert.Equal(t, uint64(0), record.Usage, "all period records should have zero usage due to validation failure")
				}
			}
		}
	})
}
