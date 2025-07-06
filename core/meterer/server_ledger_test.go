package meterer_test

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/core/mock"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

// TestServerLedgerAccountManagement tests the account ledger caching functionality
func TestServerLedgerAccountManagement(t *testing.T) {
	ctx := context.Background()
	logger := testutils.GetLogger()
	config := meterer.Config{
		ChainReadTimeout: 3 * time.Second,
		UpdateInterval:   1 * time.Second,
	}

	// Setup DynamoDB store for testing
	store, err := meterer.NewDynamoDBMeteringStore(
		clientConfig,
		reservationTableName,
		ondemandTableName,
		globalReservationTableName,
		logger,
	)
	assert.NoError(t, err)

	// Create ServerLedger
	serverLedger := meterer.NewServerLedger(config, store, logger)
	assert.NotNil(t, serverLedger)
	assert.Equal(t, 0, serverLedger.GetAccountCount())

	// Setup mock chain payment state
	chainPaymentState := &mock.MockOnchainPaymentState{}
	chainPaymentState.On("RefreshOnchainPaymentState", testifymock.Anything).Return(nil).Maybe()
	chainPaymentState.On("GetPaymentGlobalParams").Return(&meterer.PaymentVaultParams{
		QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
			0: {OnDemandSymbolsPerSecond: 1000, OnDemandPricePerSymbol: 2},
		},
		QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
			0: {MinNumSymbols: 3, OnDemandRateLimitWindow: 1},
		},
		OnDemandQuorumNumbers: []core.QuorumID{0},
	}, nil)

	// Test account creation
	privateKey1, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountID1 := crypto.PubkeyToAddress(privateKey1.PublicKey)

	privateKey2, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountID2 := crypto.PubkeyToAddress(privateKey2.PublicKey)

	// Setup chain payment state responses for test accounts
	chainPaymentState.On("GetReservedPaymentByAccountAndQuorums", testifymock.Anything, accountID1, testifymock.Anything).Return(
		map[core.QuorumID]*core.ReservedPayment{}, nil,
	)
	chainPaymentState.On("GetOnDemandPaymentByAccount", testifymock.Anything, accountID1).Return(
		&core.OnDemandPayment{CumulativePayment: big.NewInt(1000)}, nil,
	)
	chainPaymentState.On("GetReservedPaymentByAccountAndQuorums", testifymock.Anything, accountID2, testifymock.Anything).Return(
		map[core.QuorumID]*core.ReservedPayment{}, nil,
	)
	chainPaymentState.On("GetOnDemandPaymentByAccount", testifymock.Anything, accountID2).Return(
		&core.OnDemandPayment{CumulativePayment: big.NewInt(2000)}, nil,
	)

	// Test 1: Create first account ledger
	ledger1, err := serverLedger.GetAccountLedger(ctx, accountID1, chainPaymentState)
	assert.NoError(t, err)
	assert.NotNil(t, ledger1)
	assert.Equal(t, 1, serverLedger.GetAccountCount())

	// Test 2: Get same account ledger again (should return cached)
	ledger1_cached, err := serverLedger.GetAccountLedger(ctx, accountID1, chainPaymentState)
	assert.NoError(t, err)
	assert.Same(t, ledger1, ledger1_cached) // Should be exact same instance
	assert.Equal(t, 1, serverLedger.GetAccountCount())

	// Test 3: Create second account ledger
	ledger2, err := serverLedger.GetAccountLedger(ctx, accountID2, chainPaymentState)
	assert.NoError(t, err)
	assert.NotNil(t, ledger2)
	assert.NotSame(t, ledger1, ledger2) // Should be different instances
	assert.Equal(t, 2, serverLedger.GetAccountCount())

	// Test 4: Remove first account
	serverLedger.RemoveAccount(accountID1)
	assert.Equal(t, 1, serverLedger.GetAccountCount())

	// Test 5: Get removed account again (should create new instance)
	ledger1_new, err := serverLedger.GetAccountLedger(ctx, accountID1, chainPaymentState)
	assert.NoError(t, err)
	assert.NotNil(t, ledger1_new)
	assert.NotSame(t, ledger1, ledger1_new) // Should be new instance
	assert.Equal(t, 2, serverLedger.GetAccountCount())

	// Test 6: Remove all accounts
	serverLedger.RemoveAccount(accountID1)
	serverLedger.RemoveAccount(accountID2)
	assert.Equal(t, 0, serverLedger.GetAccountCount())
}

// TestServerLedgerConcurrentAccess tests thread safety of account ledger management
func TestServerLedgerConcurrentAccess(t *testing.T) {
	ctx := context.Background()
	logger := testutils.GetLogger()
	config := meterer.Config{
		ChainReadTimeout: 3 * time.Second,
		UpdateInterval:   1 * time.Second,
	}

	// Setup DynamoDB store for testing
	store, err := meterer.NewDynamoDBMeteringStore(
		clientConfig,
		reservationTableName,
		ondemandTableName,
		globalReservationTableName,
		logger,
	)
	assert.NoError(t, err)

	// Create ServerLedger
	serverLedger := meterer.NewServerLedger(config, store, logger)

	// Setup mock chain payment state
	chainPaymentState := &mock.MockOnchainPaymentState{}
	chainPaymentState.On("RefreshOnchainPaymentState", testifymock.Anything).Return(nil).Maybe()
	chainPaymentState.On("GetPaymentGlobalParams").Return(&meterer.PaymentVaultParams{
		QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
			0: {OnDemandSymbolsPerSecond: 1000, OnDemandPricePerSymbol: 2},
		},
		QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
			0: {MinNumSymbols: 3, OnDemandRateLimitWindow: 1},
		},
		OnDemandQuorumNumbers: []core.QuorumID{0},
	}, nil)

	// Create test account
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountID := crypto.PubkeyToAddress(privateKey.PublicKey)

	chainPaymentState.On("GetReservedPaymentByAccountAndQuorums", testifymock.Anything, accountID, testifymock.Anything).Return(
		map[core.QuorumID]*core.ReservedPayment{}, nil,
	)
	chainPaymentState.On("GetOnDemandPaymentByAccount", testifymock.Anything, accountID).Return(
		&core.OnDemandPayment{CumulativePayment: big.NewInt(1000)}, nil,
	)

	// Test concurrent access to same account
	numGoroutines := 10
	ledgers := make([]*meterer.ServerAccountLedger, numGoroutines)
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			ledger, err := serverLedger.GetAccountLedger(ctx, accountID, chainPaymentState)
			assert.NoError(t, err)
			ledgers[index] = ledger
		}(i)
	}

	wg.Wait()

	// All goroutines should get the same ledger instance
	for i := 1; i < numGoroutines; i++ {
		assert.Same(t, ledgers[0], ledgers[i])
	}

	// Should only have one account in cache
	assert.Equal(t, 1, serverLedger.GetAccountCount())
}

// TestServerLedgerGlobalRateLimit tests the global on-demand rate limiting functionality
func TestServerLedgerGlobalRateLimit(t *testing.T) {
	ctx := context.Background()
	logger := testutils.GetLogger()
	config := meterer.Config{
		ChainReadTimeout: 3 * time.Second,
		UpdateInterval:   1 * time.Second,
	}

	// Setup DynamoDB store for testing
	store, err := meterer.NewDynamoDBMeteringStore(
		clientConfig,
		reservationTableName,
		ondemandTableName,
		globalReservationTableName,
		logger,
	)
	assert.NoError(t, err)

	// Note: We skip clearing existing global data since each test uses unique timestamps/periods
	// and DynamoDB operations in tests are isolated by design

	// Create ServerLedger
	serverLedger := meterer.NewServerLedger(config, store, logger)

	// Setup payment vault params with specific rate limits
	params := &meterer.PaymentVaultParams{
		QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
			0: {OnDemandSymbolsPerSecond: 1000, OnDemandPricePerSymbol: 2}, // 1000 symbols/sec = 3,600,000 per hour
		},
		QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
			0: {MinNumSymbols: 3, OnDemandRateLimitWindow: 1},
		},
		OnDemandQuorumNumbers: []core.QuorumID{0},
	}

	now := time.Now()
	timestampNs := now.UnixNano()

	// Test 1: Normal usage within limit
	usage1, err := serverLedger.ValidateGlobalOnDemandUsage(ctx, timestampNs, 1000000, params) // 1M symbols
	assert.NoError(t, err)
	assert.Equal(t, uint64(1000000), usage1)

	// Test 2: Additional usage still within limit
	usage2, err := serverLedger.ValidateGlobalOnDemandUsage(ctx, timestampNs, 1000000, params) // Another 1M symbols
	assert.NoError(t, err)
	assert.Equal(t, uint64(2000000), usage2) // Cumulative: 2M symbols

	// Test 3: Usage approaching limit
	usage3, err := serverLedger.ValidateGlobalOnDemandUsage(ctx, timestampNs, 1000000, params) // Another 1M symbols
	assert.NoError(t, err)
	assert.Equal(t, uint64(3000000), usage3) // Cumulative: 3M symbols

	// Test 4: Usage exceeding limit (global limit is 3,600,000 per hour)
	_, err = serverLedger.ValidateGlobalOnDemandUsage(ctx, timestampNs, 1000000, params) // This would make 4M total
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "global rate limit exceeded")

	// Test 5: Usage in different time period should work
	nextHour := timestampNs + int64(time.Hour)
	usage4, err := serverLedger.ValidateGlobalOnDemandUsage(ctx, nextHour, 1000000, params)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1000000), usage4) // Fresh period, so starts from 1M
}

// TestServerLedgerAccountCreationFailure tests error handling when account creation fails
func TestServerLedgerAccountCreationFailure(t *testing.T) {
	ctx := context.Background()
	logger := testutils.GetLogger()
	config := meterer.Config{
		ChainReadTimeout: 3 * time.Second,
		UpdateInterval:   1 * time.Second,
	}

	// Setup DynamoDB store for testing
	store, err := meterer.NewDynamoDBMeteringStore(
		clientConfig,
		reservationTableName,
		ondemandTableName,
		globalReservationTableName,
		logger,
	)
	assert.NoError(t, err)

	// Create ServerLedger
	serverLedger := meterer.NewServerLedger(config, store, logger)

	// Setup mock chain payment state that fails for unregistered accounts
	chainPaymentState := &mock.MockOnchainPaymentState{}
	chainPaymentState.On("RefreshOnchainPaymentState", testifymock.Anything).Return(nil).Maybe()
	chainPaymentState.On("GetPaymentGlobalParams").Return(&meterer.PaymentVaultParams{
		QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
			0: {OnDemandSymbolsPerSecond: 1000, OnDemandPricePerSymbol: 2},
		},
		QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
			0: {MinNumSymbols: 3, OnDemandRateLimitWindow: 1},
		},
		OnDemandQuorumNumbers: []core.QuorumID{0},
	}, nil)

	// Create unregistered account
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountID := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Setup chain payment state to fail for this account
	chainPaymentState.On("GetReservedPaymentByAccountAndQuorums", testifymock.Anything, accountID, testifymock.Anything).Return(
		map[core.QuorumID]*core.ReservedPayment{}, fmt.Errorf("reservation not found"),
	)
	chainPaymentState.On("GetOnDemandPaymentByAccount", testifymock.Anything, accountID).Return(
		&core.OnDemandPayment{}, fmt.Errorf("payment not found"),
	)

	// Test that account creation fails
	_, err = serverLedger.GetAccountLedger(ctx, accountID, chainPaymentState)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create account ledger")

	// Account should not be cached
	assert.Equal(t, 0, serverLedger.GetAccountCount())
}

// TestServerLedgerGlobalUsageCalculation tests the global rate limit calculation logic
func TestServerLedgerGlobalUsageCalculation(t *testing.T) {
	ctx := context.Background()
	logger := testutils.GetLogger()
	config := meterer.Config{
		ChainReadTimeout: 3 * time.Second,
		UpdateInterval:   1 * time.Second,
	}

	// Setup DynamoDB store for testing
	store, err := meterer.NewDynamoDBMeteringStore(
		clientConfig,
		reservationTableName,
		ondemandTableName,
		globalReservationTableName,
		logger,
	)
	assert.NoError(t, err)

	// Note: We skip clearing existing global data since each test uses unique timestamps/periods
	// and DynamoDB operations in tests are isolated by design

	// Create ServerLedger
	serverLedger := meterer.NewServerLedger(config, store, logger)

	// Test different rate limit configurations
	testCases := []struct {
		name           string
		symbolsPerSec  uint64
		requestSymbols uint64
		expectedLimit  uint64
		shouldSucceed  bool
	}{
		{
			name:           "Low rate limit",
			symbolsPerSec:  10,    // 10 symbols/sec = 36,000 per hour
			requestSymbols: 20000, // Request 20K symbols
			expectedLimit:  36000, // Limit is 36K
			shouldSucceed:  true,  // 20K < 36K
		},
		{
			name:           "High rate limit",
			symbolsPerSec:  1000,    // 1000 symbols/sec = 3,600,000 per hour
			requestSymbols: 2000000, // Request 2M symbols
			expectedLimit:  3600000, // Limit is 3.6M
			shouldSucceed:  true,    // 2M < 3.6M
		},
		{
			name:           "Exceeded rate limit",
			symbolsPerSec:  100,    // 100 symbols/sec = 360,000 per hour
			requestSymbols: 500000, // Request 500K symbols
			expectedLimit:  360000, // Limit is 360K
			shouldSucceed:  false,  // 500K > 360K
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Use unique timestamp for each test case to ensure different global bin periods
			// Each test case gets a different hour to avoid interference
			baseTime := time.Date(2024, 1, 1, i, 0, 0, 0, time.UTC)
			timestampNs := baseTime.UnixNano()

			params := &meterer.PaymentVaultParams{
				QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
					0: {OnDemandSymbolsPerSecond: tc.symbolsPerSec, OnDemandPricePerSymbol: 2},
				},
				QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
					0: {MinNumSymbols: 3, OnDemandRateLimitWindow: 1},
				},
				OnDemandQuorumNumbers: []core.QuorumID{0},
			}

			if tc.shouldSucceed {
				usage, err := serverLedger.ValidateGlobalOnDemandUsage(ctx, timestampNs, tc.requestSymbols, params)
				assert.NoError(t, err)
				assert.Equal(t, tc.requestSymbols, usage)
				assert.LessOrEqual(t, usage, tc.expectedLimit)
			} else {
				_, err := serverLedger.ValidateGlobalOnDemandUsage(ctx, timestampNs, tc.requestSymbols, params)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "global rate limit exceeded")
			}
		})
	}
}
