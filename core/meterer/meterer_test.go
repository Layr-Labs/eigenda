package meterer_test

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"testing"
	"time"

	commonaws "github.com/Layr-Labs/eigenda/common/aws"
	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/Layr-Labs/eigenda/testbed"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	testifymock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	logger                   = test.GetLogger()
	localstackContainer      *testbed.LocalStackContainer
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
	localstackPort             = "4575"
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
		localstackPort = os.Getenv("LOCALSTACK_PORT")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if deployLocalStack {
		var err error
		localstackContainer, err = testbed.NewLocalStackContainerWithOptions(ctx, testbed.LocalStackOptions{
			ExposeHostPort: true,
			HostPort:       localstackPort,
			Services:       []string{"dynamodb"},
			Logger:         logger,
		})
		if err != nil {
			teardown()
			logger.Fatal("Failed to start localstack container:", err)
		}
	}

	clientConfig = commonaws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     fmt.Sprintf("http://0.0.0.0:%s", localstackPort),
	}

	var err error
	dynamoClient, err = commondynamodb.NewClient(clientConfig, logger)
	if err != nil {
		teardown()
		logger.Fatal("Failed to create dynamodb client:", err)
	}

	privateKey1, err := crypto.GenerateKey()
	if err != nil {
		teardown()
		logger.Fatal("Failed to generate private key:", err)
	}
	privateKey2, err := crypto.GenerateKey()
	if err != nil {
		teardown()
		logger.Fatal("Failed to generate private key:", err)
	}
	privateKey3, err := crypto.GenerateKey()
	if err != nil {
		teardown()
		logger.Fatal("Failed to generate private key:", err)
	}

	logger = test.GetLogger()
	config := meterer.Config{
		ChainReadTimeout: 3 * time.Second,
		UpdateInterval:   1 * time.Second,
	}

	err = meterer.CreateReservationTable(clientConfig, reservationTableName)
	if err != nil {
		teardown()
		logger.Fatal("Failed to create reservation table:", err)
	}
	err = meterer.CreateOnDemandTable(clientConfig, ondemandTableName)
	if err != nil {
		teardown()
		logger.Fatal("Failed to create ondemand table:", err)
	}
	err = meterer.CreateGlobalReservationTable(clientConfig, globalReservationTableName)
	if err != nil {
		teardown()
		logger.Fatal("Failed to create global reservation table:", err)
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

	store, err := meterer.NewDynamoDBMeteringStore(
		clientConfig,
		reservationTableName,
		ondemandTableName,
		globalReservationTableName,
		logger,
	)

	if err != nil {
		teardown()
		logger.Fatal("Failed to create metering store:", err)
	}

	paymentChainState.On("RefreshOnchainPaymentState", testifymock.Anything).Return(nil).Maybe()
	if err := paymentChainState.RefreshOnchainPaymentState(ctx); err != nil {
		logger.Fatal("Failed to make initial query to the on-chain state:", err)
	}

	// add some default sensible configs
	mt = meterer.NewMeterer(
		config,
		paymentChainState,
		store,
		logger,
		// metrics.NewNoopMetrics(),
	)

	mt.Start(ctx)
}

func teardown() {
	if deployLocalStack {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_ = localstackContainer.Terminate(ctx)
	}
}

func TestMetererReservations(t *testing.T) {
	ctx := t.Context()

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
	header := createPaymentHeader(t, 1, big.NewInt(0), accountID1)
	_, err := mt.MeterRequest(ctx, *header, 1000, []uint8{0, 1, 2}, now)
	require.ErrorContains(t, err, "reservation not active", "should error when reservation timestamp is not active")

	// test invalid quorom ID
	header = createPaymentHeader(t, now.UnixNano(), big.NewInt(0), accountID1)
	_, err = mt.MeterRequest(ctx, *header, 1000, []uint8{0, 1, 2}, now)
	require.ErrorContains(t, err, "invalid quorum for reservation",
		"should error when quorum IDs are invalid for reservation")

	// small bin overflow for empty bin
	header = createPaymentHeader(t,
		now.UnixNano()-int64(mt.ChainPaymentState.GetReservationWindow())*1e9, big.NewInt(0), accountID2)
	_, err = mt.MeterRequest(ctx, *header, 10, quoromNumbers, now)
	require.NoError(t, err, "small bin overflow should succeed")
	// overwhelming bin overflow for empty bins
	header = createPaymentHeader(t,
		now.UnixNano()-int64(mt.ChainPaymentState.GetReservationWindow())*1e9, big.NewInt(0), accountID2)
	_, err = mt.MeterRequest(ctx, *header, 1000, quoromNumbers, now)
	require.ErrorContains(t, err, "overflow usage exceeds bin limit", "overwhelming bin overflow should fail")

	// test non-existent account
	unregisteredUser, err := crypto.GenerateKey()
	require.NoError(t, err, "failed to generate key for unregistered user")
	header = createPaymentHeader(t, 1, big.NewInt(0), crypto.PubkeyToAddress(unregisteredUser.PublicKey))
	require.NoError(t, err, "key generation should succeed")
	_, err = mt.MeterRequest(ctx, *header, 1000, []uint8{0, 1, 2}, time.Now())
	require.ErrorContains(t, err, "failed to get active reservation by account: reservation not found", "unregistered user should fail reservation lookup")

	// test inactive reservation
	header = createPaymentHeader(t, now.UnixNano(), big.NewInt(0), accountID3)
	_, err = mt.MeterRequest(ctx, *header, 1000, []uint8{0}, now)
	require.ErrorContains(t, err, "reservation not active", "inactive reservation should fail")

	// test invalid reservation period
	header = createPaymentHeader(t,
		now.UnixNano()-2*int64(mt.ChainPaymentState.GetReservationWindow())*1e9, big.NewInt(0), accountID1)
	_, err = mt.MeterRequest(ctx, *header, 2000, quoromNumbers, now)
	require.ErrorContains(t, err, "invalid reservation period for reservation", "invalid reservation period should fail")

	// test bin usage metering
	symbolLength := uint64(20)
	requiredLength := uint(21) // 21 should be charged for length of 20 since minNumSymbols is 3
	for i := 0; i < 9; i++ {
		reservationPeriod = meterer.GetReservationPeriodByNanosecond(now.UnixNano(), mt.ChainPaymentState.GetReservationWindow())
		header = createPaymentHeader(t, now.UnixNano(), big.NewInt(0), accountID2)
		symbolsCharged, err := mt.MeterRequest(ctx, *header, symbolLength, quoromNumbers, now)
		require.NoError(t, err, "valid reservation request should succeed")
		item, err := dynamoClient.GetItem(ctx, reservationTableName, commondynamodb.Key{
			"AccountID":         &types.AttributeValueMemberS{Value: accountID2.Hex()},
			"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.Itoa(int(reservationPeriod))},
		})
		require.NotNil(t, item, "reservation record should exist in database")
		require.NoError(t, err, "database query should succeed")
		require.Equal(t, uint64(requiredLength), symbolsCharged)
		require.Equal(t, accountID2.Hex(), item["AccountID"].(*types.AttributeValueMemberS).Value)
		require.Equal(t, strconv.Itoa(int(reservationPeriod)), item["ReservationPeriod"].(*types.AttributeValueMemberN).Value)
		require.Equal(t, strconv.Itoa((i+1)*int(requiredLength)), item["BinUsage"].(*types.AttributeValueMemberN).Value)
	}
	// first over flow is allowed
	header = createPaymentHeader(t, now.UnixNano(), big.NewInt(0), accountID2)
	symbolsCharged, err := mt.MeterRequest(ctx, *header, 25, quoromNumbers, now)
	require.NoError(t, err, "first overflow should be allowed")
	require.Equal(t, uint64(27), symbolsCharged)
	overflowedReservationPeriod := reservationPeriod + 2
	item, err := dynamoClient.GetItem(ctx, reservationTableName, commondynamodb.Key{
		"AccountID":         &types.AttributeValueMemberS{Value: accountID2.Hex()},
		"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.Itoa(int(overflowedReservationPeriod))},
	})
	require.NoError(t, err)
	require.Equal(t, accountID2.Hex(), item["AccountID"].(*types.AttributeValueMemberS).Value)
	require.Equal(t, strconv.Itoa(int(overflowedReservationPeriod)),
		item["ReservationPeriod"].(*types.AttributeValueMemberN).Value)
	// 25 rounded up to the nearest multiple of minNumSymbols - (200-21*9) = 16
	require.Equal(t, strconv.Itoa(int(16)), item["BinUsage"].(*types.AttributeValueMemberN).Value)

	// second over flow
	header = createPaymentHeader(t, now.UnixNano(), big.NewInt(0), accountID2)
	require.NoError(t, err)
	_, err = mt.MeterRequest(ctx, *header, 1, quoromNumbers, now)
	require.ErrorContains(t, err, "bin has already been filled")
}

func TestMetererOnDemand(t *testing.T) {
	ctx := t.Context()
	quorumNumbers := []uint8{0, 1}
	paymentChainState.On("GetPricePerSymbol", testifymock.Anything, testifymock.Anything).Return(uint64(2), nil)
	paymentChainState.On("GetMinNumSymbols", testifymock.Anything, testifymock.Anything).Return(uint64(3), nil)
	now := time.Now()

	paymentChainState.On("GetOnDemandPaymentByAccount", testifymock.Anything, testifymock.MatchedBy(func(account gethcommon.Address) bool {
		return account == accountID1
	})).Return(account1OnDemandPayments, nil)
	paymentChainState.On("GetOnDemandPaymentByAccount", testifymock.Anything, testifymock.MatchedBy(func(account gethcommon.Address) bool {
		return account == accountID2
	})).Return(account2OnDemandPayments, nil)
	paymentChainState.On("GetOnDemandPaymentByAccount", testifymock.Anything, testifymock.Anything).Return(&core.OnDemandPayment{}, fmt.Errorf("payment not found"))
	paymentChainState.On("GetOnDemandQuorumNumbers", testifymock.Anything).Return(quorumNumbers, nil)

	// test unregistered account
	unregisteredUser, err := crypto.GenerateKey()
	require.NoError(t, err, "failed to generate key for unregistered user")
	header := createPaymentHeader(t, now.UnixNano(), big.NewInt(2), crypto.PubkeyToAddress(unregisteredUser.PublicKey))
	require.NoError(t, err)
	_, err = mt.MeterRequest(ctx, *header, 1000, quorumNumbers, now)
	require.ErrorContains(t, err, "failed to get on-demand payment by account: payment not found")

	// test invalid quorom ID
	header = createPaymentHeader(t, now.UnixNano(), big.NewInt(2), accountID1)
	_, err = mt.MeterRequest(ctx, *header, 1000, []uint8{0, 1, 2}, now)
	require.ErrorContains(t, err, "invalid quorum for On-Demand Request")

	// test insufficient cumulative payment
	header = createPaymentHeader(t, now.UnixNano(), big.NewInt(1), accountID1)
	_, err = mt.MeterRequest(ctx, *header, 1000, quorumNumbers, now)
	require.ErrorContains(t, err, "payment validation failed: payment charged is greater than cumulative payment")
	// No record for invalid payment
	result, err := dynamoClient.Query(ctx, ondemandTableName, "AccountID = :account", commondynamodb.ExpressionValues{
		":account": &types.AttributeValueMemberS{
			Value: accountID1.Hex(),
		}})
	require.NoError(t, err)
	require.Equal(t, 0, len(result))

	// test duplicated cumulative payments
	symbolLength := uint64(100)
	symbolsCharged := mt.SymbolsCharged(symbolLength)
	priceCharged := meterer.PaymentCharged(symbolsCharged, mt.ChainPaymentState.GetPricePerSymbol())
	require.Equal(t, big.NewInt(int64(102*mt.ChainPaymentState.GetPricePerSymbol())), priceCharged)
	header = createPaymentHeader(t, now.UnixNano(), priceCharged, accountID2)
	symbolsCharged, err = mt.MeterRequest(ctx, *header, symbolLength, quorumNumbers, now)
	require.NoError(t, err)
	require.Equal(t, uint64(102), symbolsCharged)
	header = createPaymentHeader(t, now.UnixNano(), priceCharged, accountID2)
	_, err = mt.MeterRequest(ctx, *header, symbolLength, quorumNumbers, now)
	// Doesn't check for exact payment, checks for increment
	require.ErrorContains(t, err, "insufficient cumulative payment increment")

	// test valid payments
	for i := 1; i < 9; i++ {
		header = createPaymentHeader(t, now.UnixNano(), new(big.Int).Mul(priceCharged, big.NewInt(int64(i+1))), accountID2)
		symbolsCharged, err = mt.MeterRequest(ctx, *header, symbolLength, quorumNumbers, now)
		require.NoError(t, err)
		require.Equal(t, uint64(102), symbolsCharged)
	}

	// test cumulative payment on-chain constraint
	header = createPaymentHeader(t, now.UnixNano(), big.NewInt(2023), accountID2)
	_, err = mt.MeterRequest(ctx, *header, 1, quorumNumbers, now)
	require.ErrorContains(t, err,
		"invalid on-demand request: request claims a cumulative payment greater than the on-chain deposit")

	// test insufficient increment in cumulative payment
	previousCumulativePayment := priceCharged.Mul(priceCharged, big.NewInt(9))
	symbolLength = uint64(2)
	symbolsCharged = mt.SymbolsCharged(symbolLength)
	priceCharged = meterer.PaymentCharged(symbolsCharged, mt.ChainPaymentState.GetPricePerSymbol())
	header = createPaymentHeader(t, now.UnixNano(),
		big.NewInt(0).Add(previousCumulativePayment, big.NewInt(0).Sub(priceCharged, big.NewInt(1))), accountID2)
	_, err = mt.MeterRequest(ctx, *header, symbolLength, quorumNumbers, now)
	require.ErrorContains(t, err, "insufficient cumulative payment increment")
	previousCumulativePayment = big.NewInt(0).Add(previousCumulativePayment, priceCharged)

	// test cannot insert cumulative payment in out of order
	symbolsCharged = mt.SymbolsCharged(uint64(50))
	header = createPaymentHeader(t, now.UnixNano(),
		meterer.PaymentCharged(symbolsCharged, mt.ChainPaymentState.GetPricePerSymbol()), accountID2)
	_, err = mt.MeterRequest(ctx, *header, 50, quorumNumbers, now)
	require.ErrorContains(t, err, "insufficient cumulative payment increment")

	result, err = dynamoClient.Query(ctx, ondemandTableName, "AccountID = :account", commondynamodb.ExpressionValues{
		":account": &types.AttributeValueMemberS{
			Value: accountID2.Hex(),
		}})
	require.NoError(t, err)
	require.Equal(t, 1, len(result))

	// with rollback of invalid payments, users cannot cheat by inserting an invalid cumulative payment
	symbolsCharged = mt.SymbolsCharged(uint64(30))
	header = createPaymentHeader(t, now.UnixNano(),
		meterer.PaymentCharged(symbolsCharged, mt.ChainPaymentState.GetPricePerSymbol()), accountID2)
	_, err = mt.MeterRequest(ctx, *header, 30, quorumNumbers, now)
	require.ErrorContains(t, err, "insufficient cumulative payment increment")

	// test failed global rate limit (previously payment recorded: 2, global limit: 1009)
	header = createPaymentHeader(t, now.UnixNano(),
		big.NewInt(0).Add(previousCumulativePayment,
			meterer.PaymentCharged(1010, mt.ChainPaymentState.GetPricePerSymbol())), accountID1)
	_, err = mt.MeterRequest(ctx, *header, 1010, quorumNumbers, now)
	require.ErrorContains(t, err, "failed global rate limiting")
	// Correct rollback
	result, err = dynamoClient.Query(ctx, ondemandTableName, "AccountID = :account", commondynamodb.ExpressionValues{
		":account": &types.AttributeValueMemberS{
			Value: accountID2.Hex(),
		}})
	require.NoError(t, err)
	require.Equal(t, 1, len(result))
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
			require.Equal(t, tt.expected, result)
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
			require.Equal(t, tt.expected, result)
		})
	}
}

func createPaymentHeader(
	t *testing.T, timestamp int64, cumulativePayment *big.Int, accountID gethcommon.Address,
) *core.PaymentMetadata {
	t.Helper()
	return &core.PaymentMetadata{
		AccountID:         accountID,
		Timestamp:         timestamp,
		CumulativePayment: cumulativePayment,
	}
}
