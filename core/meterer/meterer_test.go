package meterer_test

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

var (
	dockertestPool           *dockertest.Pool
	dockertestResource       *dockertest.Resource
	accountID1               gethcommon.Address
	account1Reservations     *core.ReservedPayment
	account1OnDemandPayments *core.OnDemandPayment
	accountID2               gethcommon.Address
	account2Reservations     *core.ReservedPayment
	account2OnDemandPayments *core.OnDemandPayment
	accountID3               gethcommon.Address
	account3Reservations     *core.ReservedPayment
	mt                       *meterer.Meterer
	store                    meterer.OffchainStore

	deployLocalStack  bool
	localStackPort    = "4566"
	paymentChainState = &mock.MockOnchainPaymentState{}
)

func TestMain(m *testing.M) {
	setup(m)
	code := m.Run()
	teardown()
	os.Exit(code)
}

func createTestStore(t *testing.T) meterer.OffchainStore {
	logger := testutils.GetLogger()
	tmpDir := os.TempDir()
	dbPath := filepath.Join(tmpDir, fmt.Sprintf("test_db_%d", time.Now().UnixNano()))
	store, err := meterer.NewOffchainStore(dbPath, logger)
	if err != nil {
		t.Fatalf("Failed to create offchain store: %v", err)
	}
	return store
}

func setup(m *testing.M) {
	// Create test accounts
	privateKey1, err := crypto.GenerateKey()
	if err != nil {
		panic("failed to generate key 1")
	}
	privateKey2, err := crypto.GenerateKey()
	if err != nil {
		panic("failed to generate key 2")
	}
	privateKey3, err := crypto.GenerateKey()
	if err != nil {
		panic("failed to generate key 3")
	}

	now := uint64(time.Now().Unix())
	accountID1 = crypto.PubkeyToAddress(privateKey1.PublicKey)
	accountID2 = crypto.PubkeyToAddress(privateKey2.PublicKey)
	accountID3 = crypto.PubkeyToAddress(privateKey3.PublicKey)
	account1Reservations = &core.ReservedPayment{SymbolsPerSecond: 20, StartTimestamp: now - 120, EndTimestamp: now + 180, QuorumSplits: []byte{50, 50}, QuorumNumbers: []uint8{0, 1}}
	account2Reservations = &core.ReservedPayment{SymbolsPerSecond: 40, StartTimestamp: now - 120, EndTimestamp: now + 180, QuorumSplits: []byte{30, 70}, QuorumNumbers: []uint8{0, 1}}
	account3Reservations = &core.ReservedPayment{SymbolsPerSecond: 40, StartTimestamp: now + 120, EndTimestamp: now + 180, QuorumSplits: []byte{30, 70}, QuorumNumbers: []uint8{0, 1}}
	account1OnDemandPayments = &core.OnDemandPayment{CumulativePayment: big.NewInt(3864)}
	account2OnDemandPayments = &core.OnDemandPayment{CumulativePayment: big.NewInt(2000)}

	paymentChainState.On("RefreshOnchainPaymentState", testifymock.Anything).Return(nil).Maybe()
	if err := paymentChainState.RefreshOnchainPaymentState(context.Background()); err != nil {
		panic("failed to make initial query to the on-chain state")
	}
}

func teardown() {
	if deployLocalStack {
		deploy.PurgeDockertestResources(dockertestPool, dockertestResource)
	}
}

func TestMetererReservations(t *testing.T) {
	store = createTestStore(t)
	mt = meterer.NewMeterer(
		meterer.Config{
			ChainReadTimeout: 5 * time.Second,
			UpdateInterval:   1 * time.Second,
		},
		paymentChainState,
		store,
		testutils.GetLogger(),
	)
	mt.Start(context.Background())
	ctx := context.Background()
	paymentChainState.On("GetReservationWindow", testifymock.Anything).Return(uint64(5), nil)
	paymentChainState.On("GetGlobalSymbolsPerSecond", testifymock.Anything).Return(uint64(1009), nil)
	paymentChainState.On("GetGlobalRatePeriodInterval", testifymock.Anything).Return(uint64(1), nil)
	paymentChainState.On("GetMinNumSymbols", testifymock.Anything).Return(uint64(3), nil)

	now := time.Now()
	reservationPeriod := meterer.GetReservationPeriodByNanosecond(now.UnixNano(), mt.ChainPaymentState.GetReservationWindow())
	quoromNumbers := []uint8{0, 1}

	paymentChainState.On("GetReservedPaymentByAccount", testifymock.Anything, testifymock.MatchedBy(func(account gethcommon.Address) bool {
		return account == accountID1
	})).Return(account1Reservations, nil)
	paymentChainState.On("GetReservedPaymentByAccount", testifymock.Anything, testifymock.MatchedBy(func(account gethcommon.Address) bool {
		return account == accountID2
	})).Return(account2Reservations, nil)
	paymentChainState.On("GetReservedPaymentByAccount", testifymock.Anything, testifymock.MatchedBy(func(account gethcommon.Address) bool {
		return account == accountID3
	})).Return(account3Reservations, nil)
	paymentChainState.On("GetReservedPaymentByAccount", testifymock.Anything, testifymock.Anything).Return(&core.ReservedPayment{}, fmt.Errorf("reservation not found"))

	// test not active reservation
	header := createPaymentHeader(1, big.NewInt(0), accountID1)
	_, err := mt.MeterRequest(ctx, *header, 1000, []uint8{0, 1, 2}, now)
	assert.ErrorContains(t, err, "reservation not active")

	// test invalid quorom ID
	header = createPaymentHeader(now.UnixNano(), big.NewInt(0), accountID1)
	_, err = mt.MeterRequest(ctx, *header, 1000, []uint8{0, 1, 2}, now)
	assert.ErrorContains(t, err, "invalid quorum for reservation")

	// small bin overflow for empty bin
	header = createPaymentHeader(now.UnixNano()-int64(mt.ChainPaymentState.GetReservationWindow())*1e9, big.NewInt(0), accountID2)
	_, err = mt.MeterRequest(ctx, *header, 10, quoromNumbers, now)
	assert.NoError(t, err)
	// overwhelming bin overflow for empty bins
	header = createPaymentHeader(now.UnixNano()-int64(mt.ChainPaymentState.GetReservationWindow())*1e9, big.NewInt(0), accountID2)
	_, err = mt.MeterRequest(ctx, *header, 1000, quoromNumbers, now)
	assert.ErrorContains(t, err, "overflow usage exceeds bin limit")

	// test non-existent account
	unregisteredUser, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}
	header = createPaymentHeader(1, big.NewInt(0), crypto.PubkeyToAddress(unregisteredUser.PublicKey))
	assert.NoError(t, err)
	_, err = mt.MeterRequest(ctx, *header, 1000, []uint8{0, 1, 2}, time.Now())
	assert.ErrorContains(t, err, "failed to get active reservation by account: reservation not found")

	// test inactive reservation
	header = createPaymentHeader(now.UnixNano(), big.NewInt(0), accountID3)
	_, err = mt.MeterRequest(ctx, *header, 1000, []uint8{0}, now)
	assert.ErrorContains(t, err, "reservation not active")

	// test invalid reservation period
	header = createPaymentHeader(now.UnixNano()-2*int64(mt.ChainPaymentState.GetReservationWindow())*1e9, big.NewInt(0), accountID1)
	_, err = mt.MeterRequest(ctx, *header, 2000, quoromNumbers, now)
	assert.ErrorContains(t, err, "invalid reservation period for reservation")

	// test bin usage metering
	symbolLength := uint64(20)
	requiredLength := uint(21) // 21 should be charged for length of 20 since minNumSymbols is 3
	for i := 0; i < 9; i++ {
		reservationPeriod = meterer.GetReservationPeriodByNanosecond(now.UnixNano(), mt.ChainPaymentState.GetReservationWindow())
		header = createPaymentHeader(now.UnixNano(), big.NewInt(0), accountID2)
		symbolsCharged, err := mt.MeterRequest(ctx, *header, symbolLength, quoromNumbers, now)
		assert.NoError(t, err)
		// Verify reservation bin usage
		binUsage, err := store.UpdateReservationBin(ctx, accountID2, reservationPeriod, 0)
		assert.NoError(t, err)
		assert.Equal(t, uint64(requiredLength)*uint64(i+1), binUsage)
		assert.Equal(t, uint64(requiredLength), symbolsCharged)
	}
	// first over flow is allowed
	header = createPaymentHeader(now.UnixNano(), big.NewInt(0), accountID2)
	symbolsCharged, err := mt.MeterRequest(ctx, *header, 25, quoromNumbers, now)
	assert.NoError(t, err)
	assert.Equal(t, uint64(27), symbolsCharged)
	overflowedReservationPeriod := reservationPeriod + 2
	// Verify overflow bin usage
	binUsage, err := store.UpdateReservationBin(ctx, accountID2, overflowedReservationPeriod, 0)
	assert.NoError(t, err)
	assert.Equal(t, uint64(16), binUsage) // 25 rounded up to the nearest multiple of minNumSymbols - (200-21*9) = 16

	// second over flow
	header = createPaymentHeader(now.UnixNano(), big.NewInt(0), accountID2)
	assert.NoError(t, err)
	_, err = mt.MeterRequest(ctx, *header, 1, quoromNumbers, now)
	assert.ErrorContains(t, err, "bin has already been filled")
}

func TestMetererOnDemand(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	store := createTestStore(t)
	quorumNumbers := []uint8{0, 1}
	paymentChainState := &mock.MockOnchainPaymentState{}
	paymentChainState.On("GetPricePerSymbol").Return(uint64(2))
	paymentChainState.On("GetMinNumSymbols").Return(uint64(1))
	paymentChainState.On("GetGlobalRatePeriodInterval").Return(uint64(300))
	paymentChainState.On("GetGlobalSymbolsPerSecond").Return(uint64(10000))
	paymentChainState.On("GetOnDemandQuorumNumbers", testifymock.Anything).Return(quorumNumbers, nil)
	paymentChainState.On("GetOnDemandPaymentByAccount", testifymock.Anything, testifymock.MatchedBy(func(account gethcommon.Address) bool {
		return account == accountID1
	})).Return(&core.OnDemandPayment{
		CumulativePayment: big.NewInt(4000),
	}, nil)
	paymentChainState.On("GetOnDemandPaymentByAccount", testifymock.Anything, testifymock.MatchedBy(func(account gethcommon.Address) bool {
		return account == accountID2
	})).Return(&core.OnDemandPayment{
		CumulativePayment: big.NewInt(2040),
	}, nil)
	paymentChainState.On("GetOnDemandPaymentByAccount", testifymock.Anything, testifymock.Anything).Return(nil, fmt.Errorf("payment not found"))

	mt := meterer.NewMeterer(
		meterer.Config{
			ChainReadTimeout: 1 * time.Second,
			UpdateInterval:   1 * time.Second,
		},
		paymentChainState,
		store,
		testutils.GetLogger(),
	)
	mt.Start(ctx)

	// test unregistered account
	unregisteredUser, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}
	header := createPaymentHeader(now.UnixNano(), big.NewInt(2), crypto.PubkeyToAddress(unregisteredUser.PublicKey))
	assert.NoError(t, err)
	_, err = mt.MeterRequest(ctx, *header, 1000, quorumNumbers, now)
	assert.ErrorContains(t, err, "failed to get on-demand payment by account: payment not found")

	// test invalid quorom ID
	header = createPaymentHeader(now.UnixNano(), big.NewInt(2), accountID1)
	_, err = mt.MeterRequest(ctx, *header, 1000, []uint8{0, 1, 2}, now)
	assert.ErrorContains(t, err, "invalid on-demand request: invalid quorum for On-Demand Request")

	// test insufficient cumulative payment
	header = createPaymentHeader(now.UnixNano(), big.NewInt(4001), accountID1)
	_, err = mt.MeterRequest(ctx, *header, 1000, quorumNumbers, now)
	assert.ErrorContains(t, err, "invalid on-demand request: request claims a cumulative payment greater than the on-chain deposit")
	// Verify no record for invalid payment
	oldPayment, err := store.GetLargestCumulativePayment(ctx, accountID1)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(0), oldPayment)

	// test successful payment
	header = createPaymentHeader(now.UnixNano(), big.NewInt(204), accountID2)
	_, err = mt.MeterRequest(ctx, *header, 100, quorumNumbers, now)
	assert.NoError(t, err)
	oldPayment, err = store.GetLargestCumulativePayment(ctx, accountID2)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(204), oldPayment)

	// test insufficient cumulative payment increment
	header = createPaymentHeader(now.UnixNano(), big.NewInt(204), accountID2)
	_, err = mt.MeterRequest(ctx, *header, 100, quorumNumbers, now)
	assert.ErrorContains(t, err, "invalid on-demand request: failed to update cumulative payment: insufficient cumulative payment increment")
	oldPayment, err = store.GetLargestCumulativePayment(ctx, accountID2)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(204), oldPayment)

	// test successful payment
	for i := 0; i < 9; i++ {
		header = createPaymentHeader(now.UnixNano(), big.NewInt(int64((i+2)*204)), accountID2)
		_, err = mt.MeterRequest(ctx, *header, 100, quorumNumbers, now)
		assert.NoError(t, err)
		oldPayment, err = store.GetLargestCumulativePayment(ctx, accountID2)
		assert.NoError(t, err)
		assert.Equal(t, big.NewInt(int64((i+2)*204)), oldPayment)
	}

	// test failed global rate limit
	header = createPaymentHeader(now.UnixNano(), big.NewInt(2023), accountID2)
	_, err = mt.MeterRequest(ctx, *header, 1, quorumNumbers, now)
	assert.ErrorContains(t, err, "invalid on-demand request: failed to update cumulative payment: insufficient cumulative payment increment")
	oldPayment, err = store.GetLargestCumulativePayment(ctx, accountID2)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(2040), oldPayment)

	// test insufficient cumulative payment increment
	header = createPaymentHeader(now.UnixNano(), big.NewInt(1841), accountID2)
	_, err = mt.MeterRequest(ctx, *header, 2, quorumNumbers, now)
	assert.ErrorContains(t, err, "invalid on-demand request: failed to update cumulative payment: insufficient cumulative payment increment")
	oldPayment, err = store.GetLargestCumulativePayment(ctx, accountID2)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(2040), oldPayment)

	// test insufficient cumulative payment increment
	header = createPaymentHeader(now.UnixNano(), big.NewInt(102), accountID2)
	_, err = mt.MeterRequest(ctx, *header, 50, quorumNumbers, now)
	assert.ErrorContains(t, err, "invalid on-demand request: failed to update cumulative payment: insufficient cumulative payment increment")
	oldPayment, err = store.GetLargestCumulativePayment(ctx, accountID2)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(2040), oldPayment)

	// test insufficient cumulative payment increment
	header = createPaymentHeader(now.UnixNano(), big.NewInt(60), accountID2)
	_, err = mt.MeterRequest(ctx, *header, 30, quorumNumbers, now)
	assert.ErrorContains(t, err, "invalid on-demand request: failed to update cumulative payment: insufficient cumulative payment increment")
	oldPayment, err = store.GetLargestCumulativePayment(ctx, accountID2)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(2040), oldPayment)

	// test successful payment
	header = createPaymentHeader(now.UnixNano(), big.NewInt(3862), accountID1)
	_, err = mt.MeterRequest(ctx, *header, 1010, quorumNumbers, now)
	assert.NoError(t, err)
	oldPayment, err = store.GetLargestCumulativePayment(ctx, accountID1)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(3862), oldPayment)
}

func TestPaymentCharged(t *testing.T) {
	tests := []struct {
		name           string
		numSymbols     uint64
		pricePerSymbol uint64
		expected       *big.Int
	}{
		{
			name:           "Simple case: 1024 symbols, price per symbol is 1",
			numSymbols:     1024,
			pricePerSymbol: 1,
			expected:       big.NewInt(1024 * 1),
		},
		{
			name:           "Higher price per symbol",
			numSymbols:     1024,
			pricePerSymbol: 2,
			expected:       big.NewInt(1024 * 2),
		},
		{
			name:           "Zero symbols",
			numSymbols:     0,
			pricePerSymbol: 5,
			expected:       big.NewInt(0),
		},
		{
			name:           "Zero price per symbol",
			numSymbols:     512,
			pricePerSymbol: 0,
			expected:       big.NewInt(0),
		},
		{
			name:           "Large number of symbols",
			numSymbols:     1 << 20, // 1 MB
			pricePerSymbol: 3,
			expected:       new(big.Int).Mul(big.NewInt(1<<20), big.NewInt(3)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := meterer.PaymentCharged(tt.numSymbols, tt.pricePerSymbol)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMeterer_symbolsCharged(t *testing.T) {
	tests := []struct {
		name          string
		symbolLength  uint64
		minNumSymbols uint64
		expected      uint64
	}{
		{
			name:          "Data length equal to min number of symobols",
			symbolLength:  1024,
			minNumSymbols: 1024,
			expected:      1024,
		},
		{
			name:          "Data length less than min number of symbols",
			symbolLength:  512,
			minNumSymbols: 1024,
			expected:      1024,
		},
		{
			name:          "Data length greater than min number of symbols",
			symbolLength:  2048,
			minNumSymbols: 1024,
			expected:      2048,
		},
		{
			name:          "Large data length",
			symbolLength:  1 << 20, // 1 MB
			minNumSymbols: 1024,
			expected:      1 << 20,
		},
		{
			name:          "Very small data length",
			symbolLength:  16,
			minNumSymbols: 1024,
			expected:      1024,
		},
	}

	paymentChainState := &mock.MockOnchainPaymentState{}
	for _, tt := range tests {
		paymentChainState.On("GetMinNumSymbols", testifymock.Anything).Return(uint64(tt.minNumSymbols), nil)
		t.Run(tt.name, func(t *testing.T) {
			m := &meterer.Meterer{
				ChainPaymentState: paymentChainState,
			}
			result := m.SymbolsCharged(tt.symbolLength)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func createPaymentHeader(timestamp int64, cumulativePayment *big.Int, accountID gethcommon.Address) *core.PaymentMetadata {
	return &core.PaymentMetadata{
		AccountID:         accountID,
		Timestamp:         timestamp,
		CumulativePayment: cumulativePayment,
	}
}
