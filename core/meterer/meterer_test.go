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

	"github.com/Layr-Labs/eigensdk-go/logging"
)

var (
	dockertestPool           *dockertest.Pool
	dockertestResource       *dockertest.Resource
	dynamoClient             commondynamodb.Client
	clientConfig             commonaws.ClientConfig
	accountID1               gethcommon.Address
	account1Reservations     *core.ActiveReservation
	account1OnDemandPayments *core.OnDemandPayment
	accountID2               gethcommon.Address
	account2Reservations     *core.ActiveReservation
	account2OnDemandPayments *core.OnDemandPayment
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

	deployLocalStack = !(os.Getenv("DEPLOY_LOCALSTACK") == "false")
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

	logger = logging.NewNoopLogger()
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

	now := uint64(time.Now().Unix())
	accountID1 = crypto.PubkeyToAddress(privateKey1.PublicKey)
	accountID2 = crypto.PubkeyToAddress(privateKey2.PublicKey)
	account1Reservations = &core.ActiveReservation{SymbolsPerSecond: 100, StartTimestamp: now + 1200, EndTimestamp: now + 1800, QuorumSplits: []byte{50, 50}, QuorumNumbers: []uint8{0, 1}}
	account2Reservations = &core.ActiveReservation{SymbolsPerSecond: 200, StartTimestamp: now - 120, EndTimestamp: now + 180, QuorumSplits: []byte{30, 70}, QuorumNumbers: []uint8{0, 1}}
	account1OnDemandPayments = &core.OnDemandPayment{CumulativePayment: big.NewInt(3864)}
	account2OnDemandPayments = &core.OnDemandPayment{CumulativePayment: big.NewInt(2000)}

	store, err := meterer.NewOffchainStore(
		clientConfig,
		reservationTableName,
		ondemandTableName,
		globalReservationTableName,
		logger,
	)

	if err != nil {
		teardown()
		panic("failed to create offchain store")
	}

	paymentChainState.On("RefreshOnchainPaymentState", testifymock.Anything).Return(nil).Maybe()
	if err := paymentChainState.RefreshOnchainPaymentState(context.Background(), nil); err != nil {
		panic("failed to make initial query to the on-chain state")
	}

	// add some default sensible configs
	mt = meterer.NewMeterer(
		config,
		paymentChainState,
		store,
		logger,
		// metrics.NewNoopMetrics(),
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
	paymentChainState.On("GetReservationWindow", testifymock.Anything).Return(uint32(1), nil)
	paymentChainState.On("GetGlobalSymbolsPerSecond", testifymock.Anything).Return(uint64(1009), nil)
	paymentChainState.On("GetGlobalRateBinInterval", testifymock.Anything).Return(uint32(1), nil)
	paymentChainState.On("GetMinNumSymbols", testifymock.Anything).Return(uint32(3), nil)

	reservationPeriod := meterer.GetReservationPeriod(uint64(time.Now().Unix()), mt.ChainPaymentState.GetReservationWindow())
	quoromNumbers := []uint8{0, 1}

	paymentChainState.On("GetActiveReservationByAccount", testifymock.Anything, testifymock.MatchedBy(func(account gethcommon.Address) bool {
		return account == accountID1
	})).Return(account1Reservations, nil)
	paymentChainState.On("GetActiveReservationByAccount", testifymock.Anything, testifymock.MatchedBy(func(account gethcommon.Address) bool {
		return account == accountID2
	})).Return(account2Reservations, nil)
	paymentChainState.On("GetActiveReservationByAccount", testifymock.Anything, testifymock.Anything).Return(&core.ActiveReservation{}, fmt.Errorf("reservation not found"))

	// test invalid quorom ID
	header := createPaymentHeader(1, big.NewInt(0), accountID1)
	err := mt.MeterRequest(ctx, *header, 1000, []uint8{0, 1, 2})
	assert.ErrorContains(t, err, "quorum number mismatch")

	// overwhelming bin overflow for empty bins
	header = createPaymentHeader(reservationPeriod-1, big.NewInt(0), accountID2)
	err = mt.MeterRequest(ctx, *header, 10, quoromNumbers)
	assert.NoError(t, err)
	// overwhelming bin overflow for empty bins
	header = createPaymentHeader(reservationPeriod-1, big.NewInt(0), accountID2)
	err = mt.MeterRequest(ctx, *header, 1000, quoromNumbers)
	assert.ErrorContains(t, err, "overflow usage exceeds bin limit")

	// test non-existent account
	unregisteredUser, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}
	header = createPaymentHeader(1, big.NewInt(0), crypto.PubkeyToAddress(unregisteredUser.PublicKey))
	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *header, 1000, []uint8{0, 1, 2})
	assert.ErrorContains(t, err, "failed to get active reservation by account: reservation not found")

	// test invalid bin index
	header = createPaymentHeader(reservationPeriod, big.NewInt(0), accountID1)
	err = mt.MeterRequest(ctx, *header, 2000, quoromNumbers)
	assert.ErrorContains(t, err, "invalid bin index for reservation")

	// test bin usage metering
	symbolLength := uint(20)
	requiredLength := uint(21) // 21 should be charged for length of 20 since minNumSymbols is 3
	for i := 0; i < 9; i++ {
		header = createPaymentHeader(reservationPeriod, big.NewInt(0), accountID2)
		err = mt.MeterRequest(ctx, *header, symbolLength, quoromNumbers)
		assert.NoError(t, err)
		item, err := dynamoClient.GetItem(ctx, reservationTableName, commondynamodb.Key{
			"AccountID":         &types.AttributeValueMemberS{Value: accountID2.Hex()},
			"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.Itoa(int(reservationPeriod))},
		})
		assert.NoError(t, err)
		assert.Equal(t, accountID2.Hex(), item["AccountID"].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, strconv.Itoa(int(reservationPeriod)), item["ReservationPeriod"].(*types.AttributeValueMemberN).Value)
		assert.Equal(t, strconv.Itoa((i+1)*int(requiredLength)), item["BinUsage"].(*types.AttributeValueMemberN).Value)

	}
	// first over flow is allowed
	header = createPaymentHeader(reservationPeriod, big.NewInt(0), accountID2)
	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *header, 25, quoromNumbers)
	assert.NoError(t, err)
	overflowedReservationPeriod := reservationPeriod + 2
	item, err := dynamoClient.GetItem(ctx, reservationTableName, commondynamodb.Key{
		"AccountID":         &types.AttributeValueMemberS{Value: accountID2.Hex()},
		"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.Itoa(int(overflowedReservationPeriod))},
	})
	assert.NoError(t, err)
	assert.Equal(t, accountID2.Hex(), item["AccountID"].(*types.AttributeValueMemberS).Value)
	assert.Equal(t, strconv.Itoa(int(overflowedReservationPeriod)), item["ReservationPeriod"].(*types.AttributeValueMemberN).Value)
	// 25 rounded up to the nearest multiple of minNumSymbols - (200-21*9) = 16
	assert.Equal(t, strconv.Itoa(int(16)), item["BinUsage"].(*types.AttributeValueMemberN).Value)

	// second over flow
	header = createPaymentHeader(reservationPeriod, big.NewInt(0), accountID2)
	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *header, 1, quoromNumbers)
	assert.ErrorContains(t, err, "bin has already been filled")
}

func TestMetererOnDemand(t *testing.T) {
	ctx := context.Background()
	quorumNumbers := []uint8{0, 1}
	paymentChainState.On("GetPricePerSymbol", testifymock.Anything).Return(uint32(2), nil)
	paymentChainState.On("GetMinNumSymbols", testifymock.Anything).Return(uint32(3), nil)
	reservationPeriod := uint32(0) // this field doesn't matter for on-demand payments wrt global rate limit

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
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}
	header := createPaymentHeader(reservationPeriod, big.NewInt(2), crypto.PubkeyToAddress(unregisteredUser.PublicKey))
	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *header, 1000, quorumNumbers)
	assert.ErrorContains(t, err, "failed to get on-demand payment by account: payment not found")

	// test invalid quorom ID
	header = createPaymentHeader(reservationPeriod, big.NewInt(2), accountID1)
	err = mt.MeterRequest(ctx, *header, 1000, []uint8{0, 1, 2})
	assert.ErrorContains(t, err, "invalid quorum for On-Demand Request")

	// test insufficient cumulative payment
	header = createPaymentHeader(reservationPeriod, big.NewInt(1), accountID1)
	err = mt.MeterRequest(ctx, *header, 1000, quorumNumbers)
	assert.ErrorContains(t, err, "insufficient cumulative payment increment")
	// No rollback after meter request
	result, err := dynamoClient.Query(ctx, ondemandTableName, "AccountID = :account", commondynamodb.ExpressionValues{
		":account": &types.AttributeValueMemberS{
			Value: accountID1.Hex(),
		}})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))

	// test duplicated cumulative payments
	symbolLength := uint(100)
	priceCharged := mt.PaymentCharged(symbolLength)
	assert.Equal(t, big.NewInt(int64(102*mt.ChainPaymentState.GetPricePerSymbol())), priceCharged)
	header = createPaymentHeader(reservationPeriod, priceCharged, accountID2)
	err = mt.MeterRequest(ctx, *header, symbolLength, quorumNumbers)
	assert.NoError(t, err)
	header = createPaymentHeader(reservationPeriod, priceCharged, accountID2)
	err = mt.MeterRequest(ctx, *header, symbolLength, quorumNumbers)
	assert.ErrorContains(t, err, "exact payment already exists")

	// test valid payments
	for i := 1; i < 9; i++ {
		header = createPaymentHeader(reservationPeriod, new(big.Int).Mul(priceCharged, big.NewInt(int64(i+1))), accountID2)
		err = mt.MeterRequest(ctx, *header, symbolLength, quorumNumbers)
		assert.NoError(t, err)
	}

	// test cumulative payment on-chain constraint
	header = createPaymentHeader(reservationPeriod, big.NewInt(2023), accountID2)
	err = mt.MeterRequest(ctx, *header, 1, quorumNumbers)
	assert.ErrorContains(t, err, "invalid on-demand payment: request claims a cumulative payment greater than the on-chain deposit")

	// test insufficient increment in cumulative payment
	previousCumulativePayment := priceCharged.Mul(priceCharged, big.NewInt(9))
	symbolLength = uint(2)
	priceCharged = mt.PaymentCharged(symbolLength)
	header = createPaymentHeader(reservationPeriod, big.NewInt(0).Add(previousCumulativePayment, big.NewInt(0).Sub(priceCharged, big.NewInt(1))), accountID2)
	err = mt.MeterRequest(ctx, *header, symbolLength, quorumNumbers)
	assert.ErrorContains(t, err, "invalid on-demand payment: insufficient cumulative payment increment")
	previousCumulativePayment = big.NewInt(0).Add(previousCumulativePayment, priceCharged)

	// test cannot insert cumulative payment in out of order
	header = createPaymentHeader(reservationPeriod, mt.PaymentCharged(50), accountID2)
	err = mt.MeterRequest(ctx, *header, 50, quorumNumbers)
	assert.ErrorContains(t, err, "invalid on-demand payment: breaking cumulative payment invariants")

	numPrevRecords := 12
	result, err = dynamoClient.Query(ctx, ondemandTableName, "AccountID = :account", commondynamodb.ExpressionValues{
		":account": &types.AttributeValueMemberS{
			Value: accountID2.Hex(),
		}})
	assert.NoError(t, err)
	assert.Equal(t, numPrevRecords, len(result))
	// test failed global rate limit (previously payment recorded: 2, global limit: 1009)
	header = createPaymentHeader(reservationPeriod, big.NewInt(0).Add(previousCumulativePayment, mt.PaymentCharged(1010)), accountID1)
	err = mt.MeterRequest(ctx, *header, 1010, quorumNumbers)
	assert.ErrorContains(t, err, "failed global rate limiting")
	// Correct rollback
	result, err = dynamoClient.Query(ctx, ondemandTableName, "AccountID = :account", commondynamodb.ExpressionValues{
		":account": &types.AttributeValueMemberS{
			Value: accountID2.Hex(),
		}})
	assert.NoError(t, err)
	assert.Equal(t, numPrevRecords, len(result))
}

func TestMeterer_paymentCharged(t *testing.T) {
	tests := []struct {
		name           string
		symbolLength   uint
		pricePerSymbol uint32
		minNumSymbols  uint32
		expected       *big.Int
	}{
		{
			name:           "Data length equal to min chargeable size",
			symbolLength:   1024,
			pricePerSymbol: 1,
			minNumSymbols:  1024,
			expected:       big.NewInt(1024),
		},
		{
			name:           "Data length less than min chargeable size",
			symbolLength:   512,
			pricePerSymbol: 1,
			minNumSymbols:  1024,
			expected:       big.NewInt(1024),
		},
		{
			name:           "Data length greater than min chargeable size",
			symbolLength:   2048,
			pricePerSymbol: 1,
			minNumSymbols:  1024,
			expected:       big.NewInt(2048),
		},
		{
			name:           "Large data length",
			symbolLength:   1 << 20, // 1 MB
			pricePerSymbol: 1,
			minNumSymbols:  1024,
			expected:       big.NewInt(1 << 20),
		},
		{
			name:           "Price not evenly divisible by min chargeable size",
			symbolLength:   1536,
			pricePerSymbol: 1,
			minNumSymbols:  1024,
			expected:       big.NewInt(2048),
		},
	}

	paymentChainState := &mock.MockOnchainPaymentState{}
	for _, tt := range tests {
		paymentChainState.On("GetPricePerSymbol", testifymock.Anything).Return(uint32(tt.pricePerSymbol), nil)
		paymentChainState.On("GetMinNumSymbols", testifymock.Anything).Return(uint32(tt.minNumSymbols), nil)
		t.Run(tt.name, func(t *testing.T) {
			m := &meterer.Meterer{
				ChainPaymentState: paymentChainState,
			}
			result := m.PaymentCharged(tt.symbolLength)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMeterer_symbolsCharged(t *testing.T) {
	tests := []struct {
		name          string
		symbolLength  uint
		minNumSymbols uint32
		expected      uint32
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
		paymentChainState.On("GetMinNumSymbols", testifymock.Anything).Return(uint32(tt.minNumSymbols), nil)
		t.Run(tt.name, func(t *testing.T) {
			m := &meterer.Meterer{
				ChainPaymentState: paymentChainState,
			}
			result := m.SymbolsCharged(tt.symbolLength)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func createPaymentHeader(reservationPeriod uint32, cumulativePayment *big.Int, accountID gethcommon.Address) *core.PaymentMetadata {
	return &core.PaymentMetadata{
		AccountID:         accountID.Hex(),
		ReservationPeriod: reservationPeriod,
		CumulativePayment: cumulativePayment,
	}
}
