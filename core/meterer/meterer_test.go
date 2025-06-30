package meterer_test

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	commonaws "github.com/Layr-Labs/eigenda/common/aws"
	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

var (
	dockertestPool           *dockertest.Pool
	dockertestResource       *dockertest.Resource
	dynamoClient             commondynamodb.Client
	clientConfig             commonaws.ClientConfig
	accountID1               gethcommon.Address
	account1Reservations     *core.ReservedPayment
	account1OnDemandPayments *core.OnDemandPayment
	accountID2               gethcommon.Address
	account2Reservations     *core.ReservedPayment
	account2OnDemandPayments *core.OnDemandPayment
	accountID3               gethcommon.Address
	account3Reservations     *core.ReservedPayment
	mt                       *meterer.Meterer

	deployLocalStack           bool
	localStackPort             = "4566"
	paymentChainState          = &mock.MockOnchainPaymentState{}
	ondemandTableName          = "ondemand_meterer"
	reservationTableName       = "reservations_meterer"
	globalReservationTableName = "global_reservation_meterer"
)

func TestMain(m *testing.M) {
	setup(m)
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup(_ *testing.M) {
	deployLocalStack = (os.Getenv("DEPLOY_LOCALSTACK") != "false")
	if !deployLocalStack {
		localStackPort = os.Getenv("LOCALSTACK_PORT")
	}

	if deployLocalStack {
		var err error
		dockertestPool, dockertestResource, err = deploy.StartDockertestWithLocalstackContainer(localStackPort)
		if err != nil {
			teardown()
			panic("failed to start localstack container")
		}
	}

	loggerConfig := common.DefaultLoggerConfig()
	logger, err := common.NewLogger(loggerConfig)
	if err != nil {
		teardown()
		panic("failed to create logger")
	}

	clientConfig = commonaws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     fmt.Sprintf("http://0.0.0.0:%s", localStackPort),
	}

	dynamoClient, err = commondynamodb.NewClient(clientConfig, logger)
	if err != nil {
		teardown()
		panic("failed to create dynamodb client")
	}

	privateKey1, err := crypto.GenerateKey()
	if err != nil {
		teardown()
		panic("failed to generate private key")
	}
	privateKey2, err := crypto.GenerateKey()
	if err != nil {
		teardown()
		panic("failed to generate private key")
	}
	privateKey3, err := crypto.GenerateKey()
	if err != nil {
		teardown()
		panic("failed to generate private key")
	}

	logger = testutils.GetLogger()
	config := meterer.Config{
		ChainReadTimeout: 3 * time.Second,
		UpdateInterval:   1 * time.Second,
	}

	err = meterer.CreateReservationTable(clientConfig, reservationTableName)
	if err != nil {
		teardown()
		panic("failed to create reservation table")
	}
	err = meterer.CreateOnDemandTable(clientConfig, ondemandTableName)
	if err != nil {
		teardown()
		panic("failed to create ondemand table")
	}
	err = meterer.CreateGlobalReservationTable(clientConfig, globalReservationTableName)
	if err != nil {
		teardown()
		panic("failed to create global reservation table")
	}

	now := time.Now()
	accountID1 = crypto.PubkeyToAddress(privateKey1.PublicKey)
	accountID2 = crypto.PubkeyToAddress(privateKey2.PublicKey)
	accountID3 = crypto.PubkeyToAddress(privateKey3.PublicKey)
	account1Reservations = &core.ReservedPayment{SymbolsPerSecond: 20, StartTimestamp: uint64(now.Add(-2 * time.Minute).Unix()), EndTimestamp: uint64(now.Add(3 * time.Minute).Unix())}
	account2Reservations = &core.ReservedPayment{SymbolsPerSecond: 40, StartTimestamp: uint64(now.Add(-2 * time.Minute).Unix()), EndTimestamp: uint64(now.Add(3 * time.Minute).Unix())}
	account3Reservations = &core.ReservedPayment{SymbolsPerSecond: 40, StartTimestamp: uint64(now.Add(2 * time.Minute).Unix()), EndTimestamp: uint64(now.Add(3 * time.Minute).Unix())}
	account1OnDemandPayments = &core.OnDemandPayment{CumulativePayment: big.NewInt(3864)}
	account2OnDemandPayments = &core.OnDemandPayment{CumulativePayment: big.NewInt(2000)}

	store, err := meterer.NewDynamoDBMeteringStore(
		clientConfig,
		reservationTableName,
		ondemandTableName,
		globalReservationTableName,
		logger,
	)

	if err != nil {
		teardown()
		panic("failed to create metering store")
	}

	paymentChainState.On("RefreshOnchainPaymentState", testifymock.Anything).Return(nil).Maybe()

	// add some default sensible configs
	mt = meterer.NewMeterer(
		config,
		paymentChainState,
		store,
		logger,
	)

	mt.Start(context.Background())
}

func teardown() {
	if deployLocalStack {
		deploy.PurgeDockertestResources(dockertestPool, dockertestResource)
	}
}

func TestMetererReservations(t *testing.T) {
	ctx := context.Background()

	// Create mock payment vault params
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
	paymentChainState.On("GetPaymentGlobalParams").Return(mockParams, nil)

	now := time.Now()
	quoromNumbers := []uint8{0, 1}
	reservationPeriods := make([]uint64, len(quoromNumbers))
	for i, quorumNumber := range quoromNumbers {
		reservationPeriods[i] = meterer.GetReservationPeriodByNanosecond(now.UnixNano(), mockParams.QuorumProtocolConfigs[core.QuorumID(quorumNumber)].ReservationRateLimitWindow)
	}

	// Update mocks for GetReservedPaymentByAccountAndQuorums
	paymentChainState.On("GetReservedPaymentByAccountAndQuorums", testifymock.Anything, accountID1, testifymock.Anything).Return(
		map[core.QuorumID]*core.ReservedPayment{0: account1Reservations, 1: account1Reservations},
	)
	paymentChainState.On("GetReservedPaymentByAccountAndQuorums", testifymock.Anything, accountID2, testifymock.Anything).Return(
		map[core.QuorumID]*core.ReservedPayment{0: account2Reservations, 1: account2Reservations},
	)
	paymentChainState.On("GetReservedPaymentByAccountAndQuorums", testifymock.Anything, accountID3, testifymock.Anything).Return(
		map[core.QuorumID]*core.ReservedPayment{0: account3Reservations, 1: account3Reservations},
	)
	paymentChainState.On("GetReservedPaymentByAccountAndQuorums", testifymock.Anything, testifymock.Anything, testifymock.Anything).Return(
		func(ctx context.Context, account gethcommon.Address, quorums []core.QuorumID) map[core.QuorumID]*core.ReservedPayment {
			return map[core.QuorumID]*core.ReservedPayment{}
		},
		fmt.Errorf("reservation not found"),
	)

	// test not active reservation
	request := createDebitSlip(1, big.NewInt(0), accountID1, 1000, []uint8{0, 1})
	request.ReceivedAt = now
	_, err := mt.MeterRequest(ctx, request)
	assert.ErrorContains(t, err, "reservation not active")

	// test invalid quorom ID
	request = createDebitSlip(now.UnixNano(), big.NewInt(0), accountID1, 1000, []uint8{0, 1, 2})
	request.ReceivedAt = now
	_, err = mt.MeterRequest(ctx, request)
	assert.ErrorContains(t, err, "quorum number mismatch")

	// small bin overflow for empty bin (using one quorum for protocol parameters for now)
	reservationWindow := mockParams.QuorumProtocolConfigs[meterer.OnDemandQuorumID].ReservationRateLimitWindow
	request = createDebitSlip(now.UnixNano()-int64(reservationWindow)*1e9, big.NewInt(0), accountID2, 10, quoromNumbers)
	request.ReceivedAt = now
	_, err = mt.MeterRequest(ctx, request)
	assert.NoError(t, err)
	// overwhelming bin overflow for empty bins
	request = createDebitSlip(now.UnixNano()-int64(reservationWindow)*1e9, big.NewInt(0), accountID2, 1000, quoromNumbers)
	request.ReceivedAt = now
	_, err = mt.MeterRequest(ctx, request)
	assert.ErrorContains(t, err, "overflow usage exceeds bin limit")

	// test non-existent account
	unregisteredUser, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}
	request = createDebitSlip(1, big.NewInt(0), crypto.PubkeyToAddress(unregisteredUser.PublicKey), 1000, []uint8{0, 1, 2})
	request.ReceivedAt = time.Now()
	assert.NoError(t, err)
	_, err = mt.MeterRequest(ctx, request)
	assert.ErrorContains(t, err, "failed to get active reservation by account: reservation not found")

	// test inactive reservation
	request = createDebitSlip(now.UnixNano(), big.NewInt(0), accountID3, 1000, []uint8{0})
	request.ReceivedAt = now
	_, err = mt.MeterRequest(ctx, request)
	assert.ErrorContains(t, err, "reservation not active")

	// test invalid reservation period
	request = createDebitSlip(now.UnixNano()-2*int64(reservationWindow)*1e9, big.NewInt(0), accountID1, 2000, quoromNumbers)
	request.ReceivedAt = now
	_, err = mt.MeterRequest(ctx, request)
	assert.ErrorContains(t, err, "invalid reservation period for reservation")
	// test bin usage metering
	symbolLength := uint64(20)
	requiredLength := uint(21) // 21 should be charged for length of 20 since minNumSymbols is 3
	accountAndQuorums := []string{}
	for _, quorum := range quoromNumbers {
		accountAndQuorums = append(accountAndQuorums, fmt.Sprintf("%s:%d", accountID2.Hex(), quorum))
	}
	for i := 0; i < 9; i++ {
		reservationPeriod := meterer.GetReservationPeriodByNanosecond(now.UnixNano(), reservationWindow)
		request = createDebitSlip(now.UnixNano(), big.NewInt(0), accountID2, symbolLength, quoromNumbers)
		request.ReceivedAt = now
		symbolsCharged, err := mt.MeterRequest(ctx, request)
		assert.NoError(t, err)
		for _, accountAndQuorum := range accountAndQuorums {
			item, err := dynamoClient.GetItem(ctx, reservationTableName, commondynamodb.Key{
				"AccountID":         &types.AttributeValueMemberS{Value: accountAndQuorum},
				"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(reservationPeriod, 10)},
			})
			assert.NotNil(t, item)
			assert.NoError(t, err)
			assert.Equal(t, uint64(requiredLength), symbolsCharged)
			assert.Equal(t, accountAndQuorum, item["AccountID"].(*types.AttributeValueMemberS).Value)
			assert.Equal(t, strconv.Itoa(int(reservationPeriod)), item["ReservationPeriod"].(*types.AttributeValueMemberN).Value)
			assert.Equal(t, strconv.Itoa((i+1)*int(requiredLength)), item["BinUsage"].(*types.AttributeValueMemberN).Value)
		}
	}
	// first over flow is allowed
	request = createDebitSlip(now.UnixNano(), big.NewInt(0), accountID2, 25, quoromNumbers)
	request.ReceivedAt = now
	symbolsCharged, err := mt.MeterRequest(ctx, request)
	assert.NoError(t, err)
	assert.Equal(t, uint64(27), symbolsCharged)

	for _, accountAndQuorum := range accountAndQuorums {
		reservationPeriod := meterer.GetReservationPeriodByNanosecond(now.UnixNano(), mockParams.QuorumProtocolConfigs[meterer.OnDemandQuorumID].ReservationRateLimitWindow)
		overflowedReservationPeriod := meterer.GetOverflowPeriod(reservationPeriod, mockParams.QuorumProtocolConfigs[meterer.OnDemandQuorumID].ReservationRateLimitWindow)
		item, err := dynamoClient.GetItem(ctx, reservationTableName, commondynamodb.Key{
			"AccountID":         &types.AttributeValueMemberS{Value: accountAndQuorum},
			"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(overflowedReservationPeriod, 10)},
		})
		assert.NotNil(t, item)
		assert.NoError(t, err)
		assert.Equal(t, accountAndQuorum, item["AccountID"].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, strconv.Itoa(int(overflowedReservationPeriod)), item["ReservationPeriod"].(*types.AttributeValueMemberN).Value)
		assert.Equal(t, strconv.Itoa(int(16)), item["BinUsage"].(*types.AttributeValueMemberN).Value)
	}

	// second over flow
	request = createDebitSlip(now.UnixNano(), big.NewInt(0), accountID2, 1, quoromNumbers)
	request.ReceivedAt = now
	_, err = mt.MeterRequest(ctx, request)
	assert.ErrorContains(t, err, "bin has already been filled")

	// Test quorum-specific behavior - one quorum succeeds, one fails
	// First, reset the bin data for quorum 1 used in the previous test
	accountAndQuorum1 := fmt.Sprintf("%s_%d", accountID2.Hex(), uint8(1))
	err = dynamoClient.DeleteItem(ctx, reservationTableName, commondynamodb.Key{
		"AccountID":         &types.AttributeValueMemberS{Value: accountAndQuorum1},
		"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.Itoa(int(reservationPeriods[0]))},
	})
	assert.NoError(t, err)

	// Now try a request that should succeed for quorum 1 but fail for quorum 0
	request = createDebitSlip(now.UnixNano(), big.NewInt(0), accountID2, 50, quoromNumbers)
	request.ReceivedAt = now
	_, err = mt.MeterRequest(ctx, request)
	assert.ErrorContains(t, err, "bin has already been filled")

	// Verify quorum 1 was not updated (because the operation should be atomic)
	item, err := dynamoClient.GetItem(ctx, reservationTableName, commondynamodb.Key{
		"AccountID":         &types.AttributeValueMemberS{Value: accountAndQuorum1},
		"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.Itoa(int(reservationPeriods[0]))},
	})
	assert.NoError(t, err)
	// The item should not exist or have zero usage since the batched update failed
	if len(item) > 0 {
		if binUsage, ok := item["BinUsage"]; ok {
			binUsageStr := binUsage.(*types.AttributeValueMemberN).Value
			binUsageVal, _ := strconv.ParseUint(binUsageStr, 10, 64)
			assert.Zero(t, binUsageVal, "Bin usage for quorum 1 should be zero since the batched update should have failed atomically")
		}
	}

}

func TestMetererOnDemand(t *testing.T) {
	ctx := context.Background()
	quorumNumbers := []uint8{0, 1}

	// Create mock payment vault params for on-demand test
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
	paymentChainState.On("GetPaymentGlobalParams").Return(mockParams, nil)

	now := time.Now()

	paymentChainState.On("GetOnDemandPaymentByAccount", testifymock.Anything, testifymock.MatchedBy(func(account gethcommon.Address) bool {
		return account == accountID1
	})).Return(account1OnDemandPayments, nil)
	paymentChainState.On("GetOnDemandPaymentByAccount", testifymock.Anything, testifymock.MatchedBy(func(account gethcommon.Address) bool {
		return account == accountID2
	})).Return(account2OnDemandPayments, nil)
	paymentChainState.On("GetOnDemandPaymentByAccount", testifymock.Anything, testifymock.Anything).Return(&core.OnDemandPayment{}, fmt.Errorf("payment not found"))

	// test unregistered account
	unregisteredUser, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}
	request := createDebitSlip(now.UnixNano(), big.NewInt(2), crypto.PubkeyToAddress(unregisteredUser.PublicKey), 1000, quorumNumbers)
	request.ReceivedAt = now
	_, err = mt.MeterRequest(ctx, request)
	assert.ErrorContains(t, err, "failed to get on-demand payment by account: payment not found")

	// test invalid quorom ID
	request = createDebitSlip(now.UnixNano(), big.NewInt(2), accountID1, 1000, []uint8{0, 1, 2})
	request.ReceivedAt = now
	_, err = mt.MeterRequest(ctx, request)
	assert.ErrorContains(t, err, "invalid quorum for On-Demand Request")

	// test insufficient cumulative payment
	request = createDebitSlip(now.UnixNano(), big.NewInt(1), accountID1, 1000, quorumNumbers)
	request.ReceivedAt = now
	_, err = mt.MeterRequest(ctx, request)
	assert.ErrorContains(t, err, "payment validation failed: payment charged is greater than cumulative payment")
	// No record for invalid payment
	result, err := dynamoClient.Query(ctx, ondemandTableName, "AccountID = :account", commondynamodb.ExpressionValues{
		":account": &types.AttributeValueMemberS{
			Value: accountID1.Hex(),
		}})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(result))

	// test duplicated cumulative payments
	symbolLength := uint64(100)
	minSymbols := mockParams.QuorumProtocolConfigs[meterer.OnDemandQuorumID].MinNumSymbols
	pricePerSymbol := mockParams.QuorumPaymentConfigs[meterer.OnDemandQuorumID].OnDemandPricePerSymbol
	symbolsCharged := meterer.SymbolsCharged(symbolLength, minSymbols)
	priceCharged := meterer.PaymentCharged(symbolsCharged, pricePerSymbol)
	assert.Equal(t, big.NewInt(int64(102*pricePerSymbol)), priceCharged)
	request = createDebitSlip(now.UnixNano(), priceCharged, accountID2, symbolLength, quorumNumbers)
	request.ReceivedAt = now
	symbolsCharged, err = mt.MeterRequest(ctx, request)
	assert.NoError(t, err)
	assert.Equal(t, uint64(102), symbolsCharged)
	request = createDebitSlip(now.UnixNano(), priceCharged, accountID2, symbolLength, quorumNumbers)
	request.ReceivedAt = now
	_, err = mt.MeterRequest(ctx, request)
	// Doesn't check for exact payment, checks for increment
	assert.ErrorContains(t, err, "insufficient cumulative payment increment")

	// test valid payments
	for i := 1; i < 9; i++ {
		request = createDebitSlip(now.UnixNano(), new(big.Int).Mul(priceCharged, big.NewInt(int64(i+1))), accountID2, symbolLength, quorumNumbers)
		request.ReceivedAt = now
		symbolsCharged, err = mt.MeterRequest(ctx, request)
		assert.NoError(t, err)
		assert.Equal(t, uint64(102), symbolsCharged)
	}

	// test cumulative payment on-chain constraint
	request = createDebitSlip(now.UnixNano(), big.NewInt(2023), accountID2, 1, quorumNumbers)
	request.ReceivedAt = now
	_, err = mt.MeterRequest(ctx, request)
	assert.ErrorContains(t, err, "request claims a cumulative payment greater than the on-chain deposit")

	// test insufficient increment in cumulative payment
	previousCumulativePayment := priceCharged.Mul(priceCharged, big.NewInt(9))
	symbolLength = uint64(2)
	symbolsCharged = meterer.SymbolsCharged(symbolLength, minSymbols)
	priceCharged = meterer.PaymentCharged(symbolsCharged, pricePerSymbol)
	request = createDebitSlip(now.UnixNano(), big.NewInt(0).Add(previousCumulativePayment, big.NewInt(0).Sub(priceCharged, big.NewInt(1))), accountID2, symbolLength, quorumNumbers)
	request.ReceivedAt = now
	_, err = mt.MeterRequest(ctx, request)
	assert.ErrorContains(t, err, "insufficient cumulative payment increment")
	previousCumulativePayment = big.NewInt(0).Add(previousCumulativePayment, priceCharged)

	// test cannot insert cumulative payment in out of order
	symbolsCharged = meterer.SymbolsCharged(uint64(50), minSymbols)
	request = createDebitSlip(now.UnixNano(), meterer.PaymentCharged(symbolsCharged, pricePerSymbol), accountID2, 50, quorumNumbers)
	request.ReceivedAt = now
	_, err = mt.MeterRequest(ctx, request)
	assert.ErrorContains(t, err, "insufficient cumulative payment increment")

	result, err = dynamoClient.Query(ctx, ondemandTableName, "AccountID = :account", commondynamodb.ExpressionValues{
		":account": &types.AttributeValueMemberS{
			Value: accountID2.Hex(),
		}})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))

	// with rollback of invalid payments, users cannot cheat by inserting an invalid cumulative payment
	symbolsCharged = meterer.SymbolsCharged(uint64(30), minSymbols)
	request = createDebitSlip(now.UnixNano(), meterer.PaymentCharged(symbolsCharged, pricePerSymbol), accountID2, 30, quorumNumbers)
	request.ReceivedAt = now
	_, err = mt.MeterRequest(ctx, request)
	assert.ErrorContains(t, err, "insufficient cumulative payment increment")

	// test failed global rate limit (previously payment recorded: 2, global limit: 1009)
	request = createDebitSlip(now.UnixNano(), big.NewInt(0).Add(previousCumulativePayment, meterer.PaymentCharged(1010, pricePerSymbol)), accountID1, 1010, quorumNumbers)
	request.ReceivedAt = now
	_, err = mt.MeterRequest(ctx, request)
	assert.ErrorContains(t, err, "failed global rate limiting")
	// Correct rollback
	result, err = dynamoClient.Query(ctx, ondemandTableName, "AccountID = :account", commondynamodb.ExpressionValues{
		":account": &types.AttributeValueMemberS{
			Value: accountID2.Hex(),
		}})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := meterer.SymbolsCharged(tt.symbolLength, tt.minNumSymbols)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func createDebitSlip(timestamp int64, cumulativePayment *big.Int, accountID gethcommon.Address, numSymbols uint64, quorumNumbers []uint8) *meterer.DebitSlip {
	paymentMetadata := core.PaymentMetadata{
		AccountID:         accountID,
		Timestamp:         timestamp,
		CumulativePayment: cumulativePayment,
	}
	return meterer.NewDebitSlip(paymentMetadata, numSymbols, quorumNumbers)
}
