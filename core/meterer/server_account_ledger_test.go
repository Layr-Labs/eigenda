package meterer_test

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/core/meterer/payment_logic"
	"github.com/Layr-Labs/eigenda/core/mock"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

// Simple OnchainPayment implementation for testing
type SimpleOnchainPaymentState struct {
	Params           *meterer.PaymentVaultParams
	OnDemandPayments map[gethcommon.Address]*core.OnDemandPayment
	ReservedPayments map[gethcommon.Address]map[core.QuorumID]*core.ReservedPayment
}

func (s *SimpleOnchainPaymentState) RefreshOnchainPaymentState(ctx context.Context) error {
	return nil
}

func (s *SimpleOnchainPaymentState) GetReservedPaymentByAccountAndQuorums(ctx context.Context, accountID gethcommon.Address, quorumNumbers []core.QuorumID) (map[core.QuorumID]*core.ReservedPayment, error) {
	if reservations, exists := s.ReservedPayments[accountID]; exists {
		return reservations, nil
	}
	return map[core.QuorumID]*core.ReservedPayment{}, nil
}

func (s *SimpleOnchainPaymentState) GetOnDemandPaymentByAccount(ctx context.Context, accountID gethcommon.Address) (*core.OnDemandPayment, error) {
	if payment, exists := s.OnDemandPayments[accountID]; exists {
		return payment, nil
	}
	return &core.OnDemandPayment{CumulativePayment: big.NewInt(0)}, nil
}

func (s *SimpleOnchainPaymentState) GetPaymentGlobalParams() (*meterer.PaymentVaultParams, error) {
	return s.Params, nil
}

// ServerAccountLedgerTest mirrors the meterer tests but for ServerAccountLedger
func TestServerAccountLedgerReservations(t *testing.T) {
	ctx := context.Background()

	// Setup equivalent to meterer test
	logger := testutils.GetLogger()
	config := meterer.Config{
		ChainReadTimeout: 3 * time.Second,
		UpdateInterval:   1 * time.Second,
	}

	// Use the same DynamoDB setup as meterer tests
	store, err := meterer.NewDynamoDBMeteringStore(
		clientConfig,
		reservationTableName,
		ondemandTableName,
		globalReservationTableName,
		logger,
	)
	assert.NoError(t, err)

	// Create mock payment vault params identical to meterer test
	mockParams := &meterer.PaymentVaultParams{
		QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
			0: {
				OnDemandSymbolsPerSecond: 1009,
				OnDemandPricePerSymbol:   2,
			},
			1: {
				OnDemandSymbolsPerSecond: 1009,
				OnDemandPricePerSymbol:   2,
			},
		},
		QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
			0: {
				MinNumSymbols:              3,
				ReservationRateLimitWindow: 5,
				OnDemandRateLimitWindow:    1,
			},
			1: {
				MinNumSymbols:              3,
				ReservationRateLimitWindow: 5,
				OnDemandRateLimitWindow:    1,
			},
		},
		OnDemandQuorumNumbers: []core.QuorumID{0, 1},
	}

	// Setup chain payment state mock identical to meterer test
	chainPaymentState := &mock.MockOnchainPaymentState{}
	chainPaymentState.On("RefreshOnchainPaymentState", testifymock.Anything).Return(nil).Maybe()
	chainPaymentState.On("GetPaymentGlobalParams").Return(mockParams, nil)

	// Setup reservation responses identical to meterer test
	chainPaymentState.On("GetReservedPaymentByAccountAndQuorums", testifymock.Anything, accountID1, testifymock.Anything).Return(
		map[core.QuorumID]*core.ReservedPayment{0: account1Reservations, 1: account1Reservations}, nil,
	)
	chainPaymentState.On("GetReservedPaymentByAccountAndQuorums", testifymock.Anything, accountID2, testifymock.Anything).Return(
		map[core.QuorumID]*core.ReservedPayment{0: account2Reservations, 1: account2Reservations}, nil,
	)
	chainPaymentState.On("GetReservedPaymentByAccountAndQuorums", testifymock.Anything, accountID3, testifymock.Anything).Return(
		map[core.QuorumID]*core.ReservedPayment{0: account3Reservations, 1: account3Reservations}, nil,
	)
	chainPaymentState.On("GetReservedPaymentByAccountAndQuorums", testifymock.Anything, testifymock.Anything, testifymock.Anything).Return(
		func(ctx context.Context, account gethcommon.Address, quorums []core.QuorumID) map[core.QuorumID]*core.ReservedPayment {
			return map[core.QuorumID]*core.ReservedPayment{}
		},
		fmt.Errorf("reservation not found"),
	)

	// Setup on-demand payment responses for constructor calls
	chainPaymentState.On("GetOnDemandPaymentByAccount", testifymock.Anything, accountID1).Return(account1OnDemandPayments, nil)
	chainPaymentState.On("GetOnDemandPaymentByAccount", testifymock.Anything, accountID2).Return(account2OnDemandPayments, nil)
	chainPaymentState.On("GetOnDemandPaymentByAccount", testifymock.Anything, accountID3).Return(&core.OnDemandPayment{CumulativePayment: big.NewInt(0)}, nil)
	chainPaymentState.On("GetOnDemandPaymentByAccount", testifymock.Anything, testifymock.Anything).Return(&core.OnDemandPayment{}, fmt.Errorf("payment not found"))

	now := time.Now()
	quorumNumbers := []uint8{0, 1}

	// Test 1: Not active reservation - should match meterer behavior
	sal1, err := meterer.NewServerAccountLedger(ctx, accountID1, chainPaymentState, store, config, logger)
	assert.NoError(t, err)

	header := createPaymentHeader(1, big.NewInt(0), accountID1)
	_, err = sal1.Debit(ctx, *header, 1000, []uint8{0, 1}, mockParams, now)
	assert.ErrorContains(t, err, "reservation not active")

	// Test 2: Invalid quorum ID - should match meterer behavior
	header = createPaymentHeader(now.UnixNano(), big.NewInt(0), accountID1)
	_, err = sal1.Debit(ctx, *header, 1000, []uint8{0, 1, 2}, mockParams, now)
	assert.ErrorContains(t, err, "quorum number mismatch")

	// Test 3: Small bin overflow for empty bin - should match meterer behavior
	// Clear any existing data first to ensure clean state
	reservationWindow := mockParams.QuorumProtocolConfigs[meterer.OnDemandQuorumID].ReservationRateLimitWindow
	testTimestamp := now.UnixNano() - int64(reservationWindow)*1e9
	reservationPeriod := payment_logic.GetReservationPeriodByNanosecond(testTimestamp, reservationWindow)

	// Clear previous test data for account2 in this reservation period
	for _, quorum := range quorumNumbers {
		accountAndQuorum := fmt.Sprintf("%s:%d", accountID2.Hex(), quorum)
		_ = dynamoClient.DeleteItem(ctx, reservationTableName, commondynamodb.Key{
			"AccountID":         &types.AttributeValueMemberS{Value: accountAndQuorum},
			"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(reservationPeriod, 10)},
		})
	}

	sal2, err := meterer.NewServerAccountLedger(ctx, accountID2, chainPaymentState, store, config, logger)
	assert.NoError(t, err)

	header = createPaymentHeader(testTimestamp, big.NewInt(0), accountID2)
	_, err = sal2.Debit(ctx, *header, 10, quorumNumbers, mockParams, now)
	assert.NoError(t, err)

	// Test 4: Overwhelming bin overflow - should match meterer behavior
	header = createPaymentHeader(testTimestamp, big.NewInt(0), accountID2)
	_, err = sal2.Debit(ctx, *header, 1000, quorumNumbers, mockParams, now)
	assert.ErrorContains(t, err, "overflow usage exceeds bin limit")

	// Test 5: Non-existent account - should match meterer behavior
	// Note: ServerAccountLedger constructor will fail if account doesn't exist
	unregisteredUser, err := crypto.GenerateKey()
	assert.NoError(t, err)
	unregisteredAccountID := crypto.PubkeyToAddress(unregisteredUser.PublicKey)

	// This should fail during construction like meterer fails during request
	_, err = meterer.NewServerAccountLedger(ctx, unregisteredAccountID, chainPaymentState, store, config, logger)
	assert.ErrorContains(t, err, "failed to get reservations: reservation not found")

	// Test 6: Inactive reservation - should match meterer behavior
	sal3, err := meterer.NewServerAccountLedger(ctx, accountID3, chainPaymentState, store, config, logger)
	assert.NoError(t, err)

	header = createPaymentHeader(now.UnixNano(), big.NewInt(0), accountID3)
	_, err = sal3.Debit(ctx, *header, 1000, []uint8{0}, mockParams, now)
	assert.ErrorContains(t, err, "reservation not active")

	// Test 7: Invalid reservation period - should match meterer behavior
	header = createPaymentHeader(now.UnixNano()-2*int64(reservationWindow)*1e9, big.NewInt(0), accountID1)
	_, err = sal1.Debit(ctx, *header, 2000, quorumNumbers, mockParams, now)
	assert.ErrorContains(t, err, "invalid reservation period for reservation")

	// Test 8: Bin usage metering - should match meterer behavior exactly
	// This test follows the same sequence as meterer_test.go lines 268-317
	symbolLength := uint64(20)
	accountAndQuorums := []string{}
	for _, quorum := range quorumNumbers {
		accountAndQuorums = append(accountAndQuorums, fmt.Sprintf("%s:%d", accountID2.Hex(), quorum))
	}

	// Clear any previous data for clean test (using current time not past time)
	currentReservationPeriod := payment_logic.GetReservationPeriodByNanosecond(now.UnixNano(), reservationWindow)
	overflowReservationPeriod := payment_logic.GetOverflowPeriod(currentReservationPeriod, reservationWindow)

	// Also clear periods from earlier tests that used different timestamps
	testTimestampFromEarlierTest := now.UnixNano() - int64(reservationWindow)*1e9
	pastReservationPeriod := payment_logic.GetReservationPeriodByNanosecond(testTimestampFromEarlierTest, reservationWindow)
	pastOverflowPeriod := payment_logic.GetOverflowPeriod(pastReservationPeriod, reservationWindow)

	for _, accountAndQuorum := range accountAndQuorums {
		// Clear current period data
		_ = dynamoClient.DeleteItem(ctx, reservationTableName, commondynamodb.Key{
			"AccountID":         &types.AttributeValueMemberS{Value: accountAndQuorum},
			"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(currentReservationPeriod, 10)},
		})
		// Clear overflow period data to prevent pollution
		_ = dynamoClient.DeleteItem(ctx, reservationTableName, commondynamodb.Key{
			"AccountID":         &types.AttributeValueMemberS{Value: accountAndQuorum},
			"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(overflowReservationPeriod, 10)},
		})
		// Clear past period data from earlier tests
		_ = dynamoClient.DeleteItem(ctx, reservationTableName, commondynamodb.Key{
			"AccountID":         &types.AttributeValueMemberS{Value: accountAndQuorum},
			"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(pastReservationPeriod, 10)},
		})
		// Clear past overflow period data from earlier tests
		_ = dynamoClient.DeleteItem(ctx, reservationTableName, commondynamodb.Key{
			"AccountID":         &types.AttributeValueMemberS{Value: accountAndQuorum},
			"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(pastOverflowPeriod, 10)},
		})
	}

	// Create fresh ledger for account2 for clean testing
	sal2_fresh, err := meterer.NewServerAccountLedger(ctx, accountID2, chainPaymentState, store, config, logger)
	assert.NoError(t, err)

	// Fill the bin with 9 requests of 20 symbols each (189 total, under 200 limit)
	for i := 0; i < 9; i++ {
		header = createPaymentHeader(now.UnixNano(), big.NewInt(0), accountID2)
		payment, err := sal2_fresh.Debit(ctx, *header, symbolLength, quorumNumbers, mockParams, now)
		assert.NoError(t, err)
		assert.Nil(t, payment) // Should be nil for reservation requests

		// Verify database state matches meterer behavior exactly
		for _, accountAndQuorum := range accountAndQuorums {
			item, err := dynamoClient.GetItem(ctx, reservationTableName, commondynamodb.Key{
				"AccountID":         &types.AttributeValueMemberS{Value: accountAndQuorum},
				"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(currentReservationPeriod, 10)},
			})
			assert.NotNil(t, item)
			assert.NoError(t, err)
			// Verify symbols charged calculation matches meterer: 20 symbols -> 21 charged (rounded up to multiple of 3)
			expectedSymbolsCharged := payment_logic.SymbolsCharged(symbolLength, mockParams.QuorumProtocolConfigs[meterer.OnDemandQuorumID].MinNumSymbols)
			assert.Equal(t, uint64(21), expectedSymbolsCharged) // 20 -> 21 with minSymbols=3
			assert.Equal(t, accountAndQuorum, item["AccountID"].(*types.AttributeValueMemberS).Value)
			assert.Equal(t, strconv.Itoa(int(currentReservationPeriod)), item["ReservationPeriod"].(*types.AttributeValueMemberN).Value)
			assert.Equal(t, strconv.Itoa((i+1)*int(expectedSymbolsCharged)), item["BinUsage"].(*types.AttributeValueMemberN).Value)
		}
	}

	// Test 9: First overflow is allowed - should match meterer behavior
	// 189 + 27 = 216 total, overflow = 216-200 = 16 symbols
	header = createPaymentHeader(now.UnixNano(), big.NewInt(0), accountID2)
	payment, err := sal2_fresh.Debit(ctx, *header, 25, quorumNumbers, mockParams, now)
	assert.NoError(t, err)
	assert.Nil(t, payment) // Should be nil for reservation requests

	// Verify symbols charged matches meterer: 25 symbols -> 27 charged (rounded up to multiple of 3)
	overflowSymbolsCharged := payment_logic.SymbolsCharged(25, mockParams.QuorumProtocolConfigs[meterer.OnDemandQuorumID].MinNumSymbols)
	assert.Equal(t, uint64(27), overflowSymbolsCharged) // 25 -> 27 with minSymbols=3

	// Verify overflow period usage matches meterer behavior
	overflowedReservationPeriod := payment_logic.GetOverflowPeriod(currentReservationPeriod, mockParams.QuorumProtocolConfigs[meterer.OnDemandQuorumID].ReservationRateLimitWindow)

	// Check each quorum separately with detailed debugging
	for i, accountAndQuorum := range accountAndQuorums {
		item, err := dynamoClient.GetItem(ctx, reservationTableName, commondynamodb.Key{
			"AccountID":         &types.AttributeValueMemberS{Value: accountAndQuorum},
			"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(overflowedReservationPeriod, 10)},
		})
		assert.NotNil(t, item, "Overflow record should exist for %s", accountAndQuorum)
		assert.NoError(t, err)
		assert.Equal(t, accountAndQuorum, item["AccountID"].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, strconv.Itoa(int(overflowedReservationPeriod)), item["ReservationPeriod"].(*types.AttributeValueMemberN).Value)

		binUsageStr := item["BinUsage"].(*types.AttributeValueMemberN).Value
		assert.Equal(t, strconv.Itoa(int(16)), binUsageStr, "Quorum %d should have 16 overflow symbols, got %s", i, binUsageStr) // 216-200=16
	}

	// Test 10: Second overflow fails - should match meterer behavior
	header = createPaymentHeader(now.UnixNano(), big.NewInt(0), accountID2)
	_, err = sal2_fresh.Debit(ctx, *header, 1, quorumNumbers, mockParams, now)
	assert.ErrorContains(t, err, "bin has already been filled")
}

func TestServerAccountLedgerOnDemand(t *testing.T) {
	ctx := context.Background()
	quorumNumbers := []uint8{0, 1}

	// Setup equivalent to meterer test
	logger := testutils.GetLogger()
	config := meterer.Config{
		ChainReadTimeout: 3 * time.Second,
		UpdateInterval:   1 * time.Second,
	}

	// Use the same DynamoDB setup as meterer tests
	store, err := meterer.NewDynamoDBMeteringStore(
		clientConfig,
		reservationTableName,
		ondemandTableName,
		globalReservationTableName,
		logger,
	)
	assert.NoError(t, err)

	// Create mock payment vault params for on-demand test - identical to meterer test
	mockParams := &meterer.PaymentVaultParams{
		QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
			meterer.OnDemandQuorumID: {
				OnDemandSymbolsPerSecond: 1009,
				OnDemandPricePerSymbol:   2,
			},
		},
		QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
			meterer.OnDemandQuorumID: {
				MinNumSymbols:           3,
				OnDemandRateLimitWindow: 1,
			},
		},
		OnDemandQuorumNumbers: []core.QuorumID{0, 1},
	}

	// Setup chain payment state mock identical to meterer test
	chainPaymentState := &mock.MockOnchainPaymentState{}
	chainPaymentState.On("RefreshOnchainPaymentState", testifymock.Anything).Return(nil).Maybe()
	chainPaymentState.On("GetPaymentGlobalParams").Return(mockParams, nil)

	// Setup ondemand payment responses - need to be more specific due to constructor calls
	// Create local variables to ensure consistent test values
	testAccount1OnDemandPayments := &core.OnDemandPayment{CumulativePayment: big.NewInt(3864)}
	testAccount2OnDemandPayments := &core.OnDemandPayment{CumulativePayment: big.NewInt(2000)}

	chainPaymentState.On("GetOnDemandPaymentByAccount", testifymock.Anything, accountID1).Return(testAccount1OnDemandPayments, nil)
	chainPaymentState.On("GetOnDemandPaymentByAccount", testifymock.Anything, accountID2).Return(testAccount2OnDemandPayments, nil)
	chainPaymentState.On("GetOnDemandPaymentByAccount", testifymock.Anything, testifymock.Anything).Return(&core.OnDemandPayment{}, fmt.Errorf("payment not found"))

	// Setup reservation responses for constructor calls
	chainPaymentState.On("GetReservedPaymentByAccountAndQuorums", testifymock.Anything, accountID1, testifymock.Anything).Return(
		map[core.QuorumID]*core.ReservedPayment{}, nil,
	)
	chainPaymentState.On("GetReservedPaymentByAccountAndQuorums", testifymock.Anything, accountID2, testifymock.Anything).Return(
		map[core.QuorumID]*core.ReservedPayment{}, nil,
	)
	chainPaymentState.On("GetReservedPaymentByAccountAndQuorums", testifymock.Anything, testifymock.Anything, testifymock.Anything).Return(
		map[core.QuorumID]*core.ReservedPayment{}, fmt.Errorf("reservation not found"),
	)

	now := time.Now()

	// Clear any previous on-demand payment records from earlier tests
	result1, err := dynamoClient.Query(ctx, ondemandTableName, "AccountID = :account", commondynamodb.ExpressionValues{
		":account": &types.AttributeValueMemberS{
			Value: accountID1.Hex(),
		}})
	assert.NoError(t, err)
	for _, item := range result1 {
		_ = dynamoClient.DeleteItem(ctx, ondemandTableName, commondynamodb.Key{
			"AccountID": item["AccountID"],
		})
	}

	result2, err := dynamoClient.Query(ctx, ondemandTableName, "AccountID = :account", commondynamodb.ExpressionValues{
		":account": &types.AttributeValueMemberS{
			Value: accountID2.Hex(),
		}})
	assert.NoError(t, err)
	for _, item := range result2 {
		_ = dynamoClient.DeleteItem(ctx, ondemandTableName, commondynamodb.Key{
			"AccountID": item["AccountID"],
		})
	}

	// Test 1: Unregistered account - should match meterer behavior
	// Note: ServerAccountLedger constructor will fail if account doesn't exist
	unregisteredUser, err := crypto.GenerateKey()
	assert.NoError(t, err)
	unregisteredAccountID := crypto.PubkeyToAddress(unregisteredUser.PublicKey)

	// This should fail during construction - ServerAccountLedger tries reservations first
	_, err = meterer.NewServerAccountLedger(ctx, unregisteredAccountID, chainPaymentState, store, config, logger)
	assert.ErrorContains(t, err, "failed to get reservations: reservation not found")

	// Test 2: Invalid quorum ID - should match meterer behavior
	sal1, err := meterer.NewServerAccountLedger(ctx, accountID1, chainPaymentState, store, config, logger)
	assert.NoError(t, err)

	header := createPaymentHeader(now.UnixNano(), big.NewInt(2), accountID1)
	_, err = sal1.Debit(ctx, *header, 1000, []uint8{0, 1, 2}, mockParams, now)
	assert.ErrorContains(t, err, "invalid quorum for On-Demand Request")

	// Test 3: Insufficient cumulative payment - should match meterer behavior
	header = createPaymentHeader(now.UnixNano(), big.NewInt(1), accountID1)
	_, err = sal1.Debit(ctx, *header, 1000, quorumNumbers, mockParams, now)
	// Should fail with payment validation from DB store
	assert.ErrorContains(t, err, "payment validation failed: payment charged is greater than cumulative payment")

	// Verify no record for invalid payment - should match meterer behavior
	result, err := dynamoClient.Query(ctx, ondemandTableName, "AccountID = :account", commondynamodb.ExpressionValues{
		":account": &types.AttributeValueMemberS{
			Value: accountID1.Hex(),
		}})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(result))

	// Test 4: Valid payments sequence - should match meterer behavior exactly
	sal2, err := meterer.NewServerAccountLedger(ctx, accountID2, chainPaymentState, store, config, logger)
	assert.NoError(t, err)

	symbolLength := uint64(100)
	minSymbols := mockParams.QuorumProtocolConfigs[meterer.OnDemandQuorumID].MinNumSymbols
	pricePerSymbol := mockParams.QuorumPaymentConfigs[meterer.OnDemandQuorumID].OnDemandPricePerSymbol
	symbolsCharged := payment_logic.SymbolsCharged(symbolLength, minSymbols)
	priceCharged := payment_logic.PaymentCharged(symbolsCharged, pricePerSymbol)
	assert.Equal(t, big.NewInt(int64(102*pricePerSymbol)), priceCharged)

	header = createPaymentHeader(now.UnixNano(), priceCharged, accountID2)
	payment, err := sal2.Debit(ctx, *header, symbolLength, quorumNumbers, mockParams, now)
	assert.NoError(t, err)
	assert.Equal(t, priceCharged, payment) // Should return payment charged for on-demand

	// Verify symbols charged calculation matches meterer: 100 symbols -> 102 charged (rounded up to multiple of 3)
	assert.Equal(t, uint64(102), symbolsCharged)

	// Test 5: Duplicated cumulative payment - should match meterer behavior
	header = createPaymentHeader(now.UnixNano(), priceCharged, accountID2)
	_, err = sal2.Debit(ctx, *header, symbolLength, quorumNumbers, mockParams, now)
	assert.ErrorContains(t, err, "insufficient cumulative payment increment")

	// Test 6: Valid payment sequence - should match meterer behavior
	for i := 1; i < 9; i++ {
		cumulativePayment := new(big.Int).Mul(priceCharged, big.NewInt(int64(i+1)))
		// Check if this would exceed on-chain deposit (2000)
		if cumulativePayment.Cmp(big.NewInt(2000)) > 0 {
			break // Stop before exceeding deposit
		}
		header = createPaymentHeader(now.UnixNano(), cumulativePayment, accountID2)
		payment, err := sal2.Debit(ctx, *header, symbolLength, quorumNumbers, mockParams, now)
		assert.NoError(t, err)
		assert.Equal(t, priceCharged, payment)
		// Verify that ServerAccountLedger calculates the same payment amount as meterer
		// Each iteration should charge the same amount (102 symbols * 2 price = 204 payment)
		assert.Equal(t, big.NewInt(int64(102*pricePerSymbol)), payment)
	}

	// Test 7: Cumulative payment on-chain constraint - should match meterer behavior
	header = createPaymentHeader(now.UnixNano(), big.NewInt(2023), accountID2)
	_, err = sal2.Debit(ctx, *header, 1, quorumNumbers, mockParams, now)
	assert.ErrorContains(t, err, "request claims a cumulative payment greater than the on-chain deposit")

	// Test 8: Insufficient increment in cumulative payment - should match meterer behavior
	previousCumulativePayment := priceCharged.Mul(priceCharged, big.NewInt(9))
	symbolLength = uint64(2)
	symbolsCharged = payment_logic.SymbolsCharged(symbolLength, minSymbols)
	priceCharged = payment_logic.PaymentCharged(symbolsCharged, pricePerSymbol)
	header = createPaymentHeader(now.UnixNano(), big.NewInt(0).Add(previousCumulativePayment, big.NewInt(0).Sub(priceCharged, big.NewInt(1))), accountID2)
	_, err = sal2.Debit(ctx, *header, symbolLength, quorumNumbers, mockParams, now)
	assert.ErrorContains(t, err, "insufficient cumulative payment increment")

	// Test 9: Cannot insert cumulative payment out of order - should match meterer behavior
	symbolsCharged = payment_logic.SymbolsCharged(uint64(50), minSymbols)
	header = createPaymentHeader(now.UnixNano(), payment_logic.PaymentCharged(symbolsCharged, pricePerSymbol), accountID2)
	_, err = sal2.Debit(ctx, *header, 50, quorumNumbers, mockParams, now)
	assert.ErrorContains(t, err, "insufficient cumulative payment increment")

	// Verify database state matches meterer behavior
	result, err = dynamoClient.Query(ctx, ondemandTableName, "AccountID = :account", commondynamodb.ExpressionValues{
		":account": &types.AttributeValueMemberS{
			Value: accountID2.Hex(),
		}})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))
}

func TestServerAccountLedgerRevertDebit(t *testing.T) {
	ctx := context.Background()

	// Setup equivalent to meterer test
	logger := testutils.GetLogger()
	config := meterer.Config{
		ChainReadTimeout: 3 * time.Second,
		UpdateInterval:   1 * time.Second,
	}

	// Use the same DynamoDB setup as meterer tests
	store, err := meterer.NewDynamoDBMeteringStore(
		clientConfig,
		reservationTableName,
		ondemandTableName,
		globalReservationTableName,
		logger,
	)
	assert.NoError(t, err)

	// Create mock payment vault params - must include both quorums 0 and 1
	mockParams := &meterer.PaymentVaultParams{
		QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
			0: {
				OnDemandSymbolsPerSecond: 1009,
				OnDemandPricePerSymbol:   2,
			},
			1: {
				OnDemandSymbolsPerSecond: 1009,
				OnDemandPricePerSymbol:   2,
			},
		},
		QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
			0: {
				MinNumSymbols:           3,
				OnDemandRateLimitWindow: 1,
			},
			1: {
				MinNumSymbols:           3,
				OnDemandRateLimitWindow: 1,
			},
		},
		OnDemandQuorumNumbers: []core.QuorumID{0, 1},
	}

	// Setup chain payment state mock
	chainPaymentState := &mock.MockOnchainPaymentState{}
	chainPaymentState.On("RefreshOnchainPaymentState", testifymock.Anything).Return(nil).Maybe()
	chainPaymentState.On("GetPaymentGlobalParams").Return(mockParams, nil)

	// Use local test values to ensure clean state
	testAccount2OnDemandPayments := &core.OnDemandPayment{CumulativePayment: big.NewInt(2000)}
	chainPaymentState.On("GetOnDemandPaymentByAccount", testifymock.Anything, accountID2).Return(testAccount2OnDemandPayments, nil)
	chainPaymentState.On("GetReservedPaymentByAccountAndQuorums", testifymock.Anything, accountID2, testifymock.Anything).Return(
		map[core.QuorumID]*core.ReservedPayment{}, nil,
	)

	now := time.Now()
	quorumNumbers := []uint8{0, 1}

	// Clear any previous on-demand payment records from earlier tests with enhanced cleanup
	result, err := dynamoClient.Query(ctx, ondemandTableName, "AccountID = :account", commondynamodb.ExpressionValues{
		":account": &types.AttributeValueMemberS{
			Value: accountID2.Hex(),
		}})
	assert.NoError(t, err)
	for _, item := range result {
		deleteErr := dynamoClient.DeleteItem(ctx, ondemandTableName, commondynamodb.Key{
			"AccountID": item["AccountID"], // Only AccountID - no CumulativePayment needed!
		})
		if deleteErr != nil {
			t.Logf("Failed to delete item: %v", deleteErr)
		}
	}

	// Double-check cleanup worked
	verifyResult, err := dynamoClient.Query(ctx, ondemandTableName, "AccountID = :account", commondynamodb.ExpressionValues{
		":account": &types.AttributeValueMemberS{
			Value: accountID2.Hex(),
		}})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(verifyResult), "Database cleanup failed - records still exist for RevertDebit test")

	// Test RevertDebit with on-demand payment
	sal2, err2 := meterer.NewServerAccountLedger(ctx, accountID2, chainPaymentState, store, config, logger)
	assert.NoError(t, err2)

	symbolLength := uint64(100)
	pricePerSymbol := mockParams.QuorumPaymentConfigs[meterer.OnDemandQuorumID].OnDemandPricePerSymbol
	symbolsCharged := payment_logic.SymbolsCharged(symbolLength, mockParams.QuorumProtocolConfigs[meterer.OnDemandQuorumID].MinNumSymbols)
	priceCharged := payment_logic.PaymentCharged(symbolsCharged, pricePerSymbol)

	// Make a successful payment
	header := createPaymentHeader(now.UnixNano(), priceCharged, accountID2)
	payment, err := sal2.Debit(ctx, *header, symbolLength, quorumNumbers, mockParams, now)
	assert.NoError(t, err)
	assert.Equal(t, priceCharged, payment)

	// Test RevertDebit functionality
	err = sal2.RevertDebit(ctx, *header, symbolLength, quorumNumbers, mockParams, now, payment)
	assert.NoError(t, err)

	// Verify payment was rolled back by checking database state
	result2, err3 := dynamoClient.Query(ctx, ondemandTableName, "AccountID = :account", commondynamodb.ExpressionValues{
		":account": &types.AttributeValueMemberS{
			Value: accountID2.Hex(),
		}})
	assert.NoError(t, err3)
	// Should have reverted to original state
	if len(result2) > 0 {
		cumulativePayment := result2[0]["CumulativePayment"].(*types.AttributeValueMemberN).Value
		assert.Equal(t, "0", cumulativePayment) // Should be rolled back to zero (oldPayment)
	}
}

// TestErrorMessageComparison verifies that both Meterer and ServerAccountLedger fail consistently for unregistered accounts
func TestErrorMessageComparison(t *testing.T) {
	ctx := context.Background()
	logger := testutils.GetLogger()
	config := meterer.Config{
		ChainReadTimeout: 3 * time.Second,
		UpdateInterval:   1 * time.Second,
	}

	// Test account setup
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)
	accountID := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Mock payment vault params
	mockParams := &meterer.PaymentVaultParams{
		QuorumPaymentConfigs: map[core.QuorumID]*core.PaymentQuorumConfig{
			0: {OnDemandSymbolsPerSecond: 1009, OnDemandPricePerSymbol: 2},
		},
		QuorumProtocolConfigs: map[core.QuorumID]*core.PaymentQuorumProtocolConfig{
			0: {MinNumSymbols: 3, ReservationRateLimitWindow: 5, OnDemandRateLimitWindow: 1},
		},
		OnDemandQuorumNumbers: []core.QuorumID{0},
	}

	// Simple chain payment state that fails for unregistered accounts
	chainPaymentState := &SimpleOnchainPaymentState{
		Params:           mockParams,
		OnDemandPayments: map[gethcommon.Address]*core.OnDemandPayment{},
		ReservedPayments: map[gethcommon.Address]map[core.QuorumID]*core.ReservedPayment{},
	}

	// Use a simple DynamoDB store setup for testing
	store, err := meterer.NewDynamoDBMeteringStore(
		clientConfig,
		reservationTableName,
		ondemandTableName,
		globalReservationTableName,
		logger,
	)
	assert.NoError(t, err)

	// Test that both Meterer and ServerAccountLedger fail consistently for unregistered accounts
	metererInstance := meterer.NewMeterer(config, chainPaymentState, store, logger)

	now := time.Now()
	header := &core.PaymentMetadata{
		AccountID:         accountID,
		Timestamp:         now.UnixNano(),
		CumulativePayment: big.NewInt(0),
	}

	// Both should fail with reservation-related errors for unregistered accounts
	_, err1 := metererInstance.MeterRequest(ctx, *header, 1000, []uint8{0}, now)

	// For ServerAccountLedger, test the same operation (Debit) rather than just construction
	sal, err2 := meterer.NewServerAccountLedger(ctx, accountID, chainPaymentState, store, config, logger)
	if err2 == nil {
		// If construction succeeds, try the same operation as meterer
		_, err2 = sal.Debit(ctx, *header, 1000, []uint8{0}, mockParams, now)
	}

	// Both should have similar error patterns (though exact messages may differ due to different execution flows)
	assert.Error(t, err1)
	assert.Error(t, err2)

	// The key is that both fail when trying to handle unregistered accounts, maintaining the same security model
}
