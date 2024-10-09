package meterer_test

import (
	"context"
	"fmt"
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
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"

	"github.com/Layr-Labs/eigensdk-go/logging"
)

var (
	dockertestPool     *dockertest.Pool
	dockertestResource *dockertest.Resource
	dynamoClient       *commondynamodb.Client
	clientConfig       commonaws.ClientConfig
	accountID1         string
	accountID2         string
	mt                 *meterer.Meterer

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
	accountID1 = crypto.PubkeyToAddress(privateKey1.PublicKey).Hex()
	privateKey2, err := crypto.GenerateKey()
	if err != nil {
		teardown()
		panic("failed to generate private key")
	}
	accountID2 = crypto.PubkeyToAddress(privateKey2.PublicKey).Hex()

	logger = logging.NewNoopLogger()
	config := meterer.Config{
		PricePerSymbol:         1,
		MinNumSymbols:          1,
		GlobalSymbolsPerSecond: 1000,
		ReservationWindow:      1,
		ChainReadTimeout:       3 * time.Second,
	}

	err = meterer.CreateReservationTable(clientConfig, "reservations")
	if err != nil {
		teardown()
		panic("failed to create reservation table")
	}
	err = meterer.CreateOnDemandTable(clientConfig, "ondemand")
	if err != nil {
		teardown()
		panic("failed to create ondemand table")
	}
	err = meterer.CreateGlobalReservationTable(clientConfig, "global")
	if err != nil {
		teardown()
		panic("failed to create global reservation table")
	}

	store, err := meterer.NewOffchainStore(
		clientConfig,
		"reservations",
		"ondemand",
		"global",
		logger,
	)

	if err != nil {
		teardown()
		panic("failed to create offchain store")
	}

	// add some default sensible configs
	mt, err = meterer.NewMeterer(
		config,
		paymentChainState,
		store,
		logging.NewNoopLogger(),
		// metrics.NewNoopMetrics(),
	)

	if err != nil {
		teardown()
		panic("failed to create meterer")
	}
}

func teardown() {
	if deployLocalStack {
		deploy.PurgeDockertestResources(dockertestPool, dockertestResource)
	}
}

func TestMetererReservations(t *testing.T) {
	ctx := context.Background()
	meterer.CreateReservationTable(clientConfig, "reservations")
	binIndex := meterer.GetBinIndex(uint64(time.Now().Unix()), mt.ReservationWindow)
	quoromNumbers := []uint8{0, 1}
	// test invalid quorom ID
	blob, header := createMetererInput(1, 0, 1000, []uint8{0, 1, 2}, accountID1)
	err := mt.MeterRequest(ctx, *blob, *header)
	assert.ErrorContains(t, err, "quorum number mismatch")

	// test non-existent account
	unregisteredUser, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}
	blob, header = createMetererInput(1, 0, 1000, []uint8{0, 1, 2}, crypto.PubkeyToAddress(unregisteredUser.PublicKey).Hex())
	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *blob, *header)
	assert.ErrorContains(t, err, "failed to get active reservation by account: reservation not found")

	// test invalid bin index
	blob, header = createMetererInput(binIndex, 0, 2000, quoromNumbers, accountID1)
	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *blob, *header)
	assert.ErrorContains(t, err, "invalid bin index for reservation")

	// test bin usage metering
	dataLength := uint(20)
	for i := 0; i < 9; i++ {
		blob, header = createMetererInput(binIndex, 0, dataLength, quoromNumbers, accountID2)
		assert.NoError(t, err)
		err = mt.MeterRequest(ctx, *blob, *header)
		assert.NoError(t, err)
		item, err := dynamoClient.GetItem(ctx, "reservations", commondynamodb.Key{
			"AccountID": &types.AttributeValueMemberS{Value: accountID2},
			"BinIndex":  &types.AttributeValueMemberN{Value: strconv.Itoa(int(binIndex))},
		})
		assert.NoError(t, err)
		assert.Equal(t, accountID2, item["AccountID"].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, strconv.Itoa(int(binIndex)), item["BinIndex"].(*types.AttributeValueMemberN).Value)
		assert.Equal(t, strconv.Itoa((i+1)*int(dataLength)), item["BinUsage"].(*types.AttributeValueMemberN).Value)

	}
	// frist over flow is allowed
	blob, header = createMetererInput(binIndex, 0, 25, quoromNumbers, accountID2)
	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *blob, *header)
	assert.NoError(t, err)
	overflowedBinIndex := binIndex + 2
	item, err := dynamoClient.GetItem(ctx, "reservations", commondynamodb.Key{
		"AccountID": &types.AttributeValueMemberS{Value: accountID2},
		"BinIndex":  &types.AttributeValueMemberN{Value: strconv.Itoa(int(overflowedBinIndex))},
	})
	assert.NoError(t, err)
	assert.Equal(t, accountID2, item["AccountID"].(*types.AttributeValueMemberS).Value)
	assert.Equal(t, strconv.Itoa(int(overflowedBinIndex)), item["BinIndex"].(*types.AttributeValueMemberN).Value)
	assert.Equal(t, strconv.Itoa(int(5)), item["BinUsage"].(*types.AttributeValueMemberN).Value)

	// second over flow
	blob, header = createMetererInput(binIndex, 0, 1, quoromNumbers, accountID2)
	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *blob, *header)
	assert.ErrorContains(t, err, "bin has already been filled")

	// overwhelming bin overflow for empty bins (assuming all previous requests happened within 1 reservation window)
	blob, header = createMetererInput(binIndex-1, 0, 1000, quoromNumbers, accountID2)
	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *blob, *header)
	assert.ErrorContains(t, err, "overflow usage exceeds bin limit")
}

func TestMetererOnDemand(t *testing.T) {
	ctx := context.Background()
	meterer.CreateOnDemandTable(clientConfig, "ondemand")
	meterer.CreateGlobalReservationTable(clientConfig, "global")
	quorumNumbers := []uint8{0, 1}
	binIndex := uint32(0) // this field doesn't matter for on-demand payments wrt global rate limit

	// test unregistered account
	unregisteredUser, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}
	blob, header := createMetererInput(1, 1, 1000, quorumNumbers, crypto.PubkeyToAddress(unregisteredUser.PublicKey).Hex())
	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *blob, *header)
	assert.ErrorContains(t, err, "failed to get on-demand payment by account: payment not found")

	// test invalid quorom ID
	blob, header = createMetererInput(1, 1, 1000, []uint8{0, 1, 2}, accountID1)
	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *blob, *header)
	assert.ErrorContains(t, err, "invalid quorum for On-Demand Request")

	// test insufficient cumulative payment
	blob, header = createMetererInput(0, 1, 2000, quorumNumbers, accountID1)

	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *blob, *header)
	assert.ErrorContains(t, err, "insufficient cumulative payment increment")
	// No rollback after meter request
	result, err := dynamoClient.QueryIndex(ctx, "ondemand", "AccountIDIndex", "AccountID = :account", commondynamodb.ExpressionValues{
		":account": &types.AttributeValueMemberS{
			Value: accountID1,
		}})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))

	// test duplicated cumulative payments
	blob, header = createMetererInput(binIndex, uint64(100), 100, quorumNumbers, accountID2)
	err = mt.MeterRequest(ctx, *blob, *header)
	assert.NoError(t, err)
	blob, header = createMetererInput(binIndex, uint64(100), 100, quorumNumbers, accountID2)
	err = mt.MeterRequest(ctx, *blob, *header)
	assert.ErrorContains(t, err, "exact payment already exists")

	// test valid payments
	for i := 1; i < 9; i++ {
		blob, header = createMetererInput(binIndex, uint64(100*(i+1)), 100, quorumNumbers, accountID2)
		err = mt.MeterRequest(ctx, *blob, *header)
		assert.NoError(t, err)
	}

	// test cumulative payment on-chain constraint
	blob, header = createMetererInput(binIndex, 1001, 1, quorumNumbers, accountID2)
	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *blob, *header)
	assert.ErrorContains(t, err, "invalid on-demand payment: request claims a cumulative payment greater than the on-chain deposit")

	// test insufficient increment in cumulative payment
	blob, header = createMetererInput(binIndex, 901, 2, quorumNumbers, accountID2)
	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *blob, *header)
	assert.ErrorContains(t, err, "invalid on-demand payment: insufficient cumulative payment increment")

	// test cannot insert cumulative payment in out of order
	blob, header = createMetererInput(binIndex, 50, 50, quorumNumbers, accountID2)
	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *blob, *header)
	assert.ErrorContains(t, err, "invalid on-demand payment: breaking cumulative payment invariants")

	numPrevRecords := 12
	result, err = dynamoClient.QueryIndex(ctx, "ondemand", "AccountIDIndex", "AccountID = :account", commondynamodb.ExpressionValues{
		":account": &types.AttributeValueMemberS{
			Value: accountID2,
		}})
	assert.NoError(t, err)
	assert.Equal(t, numPrevRecords, len(result))
	// test failed global rate limit
	blob, header = createMetererInput(binIndex, 1002, 1001, quorumNumbers, accountID1)
	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *blob, *header)
	assert.ErrorContains(t, err, "failed global rate limiting")
	// Correct rollback
	result, err = dynamoClient.QueryIndex(ctx, "ondemand", "AccountIDIndex", "AccountID = :account", commondynamodb.ExpressionValues{
		":account": &types.AttributeValueMemberS{
			Value: accountID2,
		}})
	assert.NoError(t, err)
	assert.Equal(t, numPrevRecords, len(result))
}

func TestMeterer_paymentCharged(t *testing.T) {
	tests := []struct {
		name           string
		dataLength     uint
		pricePerSymbol uint32
		minNumSymbols  uint32
		expected       uint64
	}{
		{
			name:           "Data length equal to min chargeable size",
			dataLength:     1024,
			pricePerSymbol: 100,
			minNumSymbols:  1024,
			expected:       100,
		},
		{
			name:           "Data length less than min chargeable size",
			dataLength:     512,
			pricePerSymbol: 100,
			minNumSymbols:  1024,
			expected:       100,
		},
		{
			name:           "Data length greater than min chargeable size",
			dataLength:     2048,
			pricePerSymbol: 100,
			minNumSymbols:  1024,
			expected:       200,
		},
		{
			name:           "Large data length",
			dataLength:     1 << 20, // 1 MB
			pricePerSymbol: 100,
			minNumSymbols:  1024,
			expected:       102400,
		},
		{
			name:           "Price not evenly divisible by min chargeable size",
			dataLength:     1536,
			pricePerSymbol: 150,
			minNumSymbols:  1024,
			expected:       225,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &meterer.Meterer{
				Config: meterer.Config{
					PricePerSymbol: tt.pricePerSymbol,
					MinNumSymbols:  tt.minNumSymbols,
				},
			}
			result := m.PaymentCharged(tt.dataLength)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMeterer_symbolsCharged(t *testing.T) {
	tests := []struct {
		name          string
		dataLength    uint
		minNumSymbols uint32
		expected      uint32
	}{
		{
			name:          "Data length equal to min chargeable size",
			dataLength:    1024,
			minNumSymbols: 1024,
			expected:      1024,
		},
		{
			name:          "Data length less than min chargeable size",
			dataLength:    512,
			minNumSymbols: 1024,
			expected:      1024,
		},
		{
			name:          "Data length greater than min chargeable size",
			dataLength:    2048,
			minNumSymbols: 1024,
			expected:      2048,
		},
		{
			name:          "Large data length",
			dataLength:    1 << 20, // 1 MB
			minNumSymbols: 1024,
			expected:      1 << 20,
		},
		{
			name:          "Very small data length",
			dataLength:    16,
			minNumSymbols: 1024,
			expected:      1024,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &meterer.Meterer{
				Config: meterer.Config{
					MinNumSymbols: tt.minNumSymbols,
				},
			}
			result := m.SymbolsCharged(tt.dataLength)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func createMetererInput(binIndex uint32, cumulativePayment uint64, dataLength uint, quorumNumbers []uint8, accountID string) (blob *core.Blob, header *core.PaymentMetadata) {
	sp := make([]*core.SecurityParam, len(quorumNumbers))
	for i, quorumID := range quorumNumbers {
		sp[i] = &core.SecurityParam{
			QuorumID: quorumID,
		}
	}
	blob = &core.Blob{
		RequestHeader: core.BlobRequestHeader{
			BlobAuthHeader: core.BlobAuthHeader{
				AccountID: accountID2,
				BlobCommitments: encoding.BlobCommitments{
					Length: dataLength,
				},
			},
			SecurityParams: sp,
		},
	}
	header = &core.PaymentMetadata{
		AccountID:         accountID,
		BinIndex:          binIndex,
		CumulativePayment: cumulativePayment,
	}
	return blob, header
}
