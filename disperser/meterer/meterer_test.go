package meterer_test

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	commonaws "github.com/Layr-Labs/eigenda/common/aws"
	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	test_utils "github.com/Layr-Labs/eigenda/common/aws/dynamodb/utils"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser/meterer"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ory/dockertest/v3"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/stretchr/testify/assert"
)

var (
	dockertestPool     *dockertest.Pool
	dockertestResource *dockertest.Resource
	dynamoClient       *commondynamodb.Client
	clientConfig       commonaws.ClientConfig
	privateKey1        *ecdsa.PrivateKey
	privateKey2        *ecdsa.PrivateKey
	signer             *meterer.EIP712Signer
	mt                 *meterer.Meterer

	deployLocalStack bool
	localStackPort   = "4567"
)

func TestMain(m *testing.M) {
	setup(m)
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup(m *testing.M) {

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

	chainID := big.NewInt(17000)
	verifyingContract := gethcommon.HexToAddress("0x1234000000000000000000000000000000000000")
	signer = meterer.NewEIP712Signer(chainID, verifyingContract)

	privateKey1, err = crypto.GenerateKey()
	privateKey2, err = crypto.GenerateKey()

	logger = logging.NewNoopLogger()
	config := meterer.Config{
		PricePerByte:         1,
		GlobalBytesPerSecond: 1000,
		ReservationWindow:    time.Minute,
	}
	// metrics := metrics.NewNoopMetrics()

	paymentChainState := meterer.NewMockedOnchainPaymentState()

	paymentChainState.InitializeMockData(privateKey1, privateKey2)

	clientConfig := commonaws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     fmt.Sprintf("http://0.0.0.0:4566"),
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
		meterer.TimeoutConfig{},
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

// func setupTestMeterer(t *testing.T) (*meterer.Meterer, *meterer.EIP712Signer, *ecdsa.PrivateKey, *ecdsa.PrivateKey, error) {
// 	chainID := big.NewInt(17000)
// 	verifyingContract := gethcommon.HexToAddress("0x1234000000000000000000000000000000000000")
// 	signer := meterer.NewEIP712Signer(chainID, verifyingContract)

// 	privateKey1, err := crypto.GenerateKey()
// 	privateKey2, err := crypto.GenerateKey()

// 	logger := logging.NewNoopLogger()
// 	config := meterer.Config{
// 		PricePerByte:         1,
// 		GlobalBytesPerSecond: 1000,
// 		ReservationWindow:    time.Minute,
// 	}
// 	// metrics := metrics.NewNoopMetrics()

// 	paymentChainState := meterer.NewMockedOnchainPaymentState()

// 	paymentChainState.InitializeMockData(privateKey1, privateKey2)

// 	clientConfig := commonaws.ClientConfig{
// 		Region:          "us-east-1",
// 		AccessKey:       "localstack",
// 		SecretAccessKey: "localstack",
// 		EndpointURL:     fmt.Sprintf("http://0.0.0.0:4566"),
// 	}

// 	store, err := meterer.NewOffchainStore(
// 		clientConfig,
// 		"reservations",
// 		"ondemand",
// 		"global",
// 		logger,
// 	)
// 	if err != nil {
// 		return nil, nil, nil, nil, fmt.Errorf("failed to create offchain store: %w", err)
// 	}
// 	CreateReservationTable(t, "reservations")
// 	CreateOnDemandTable(t, "ondemand")
// 	CreateGlobalReservationTable(t, "global")

// 	// add some default sensible configs
// 	m, err := meterer.NewMeterer(
// 		config,
// 		meterer.TimeoutConfig{},
// 		paymentChainState,
// 		store,
// 		logging.NewNoopLogger(),
// 		// metrics.NewNoopMetrics(),
// 	)
// 	if err != nil {
// 		return nil, nil, nil, nil, err
// 	}

// 	require.NoError(t, err)
// 	return m, signer, privateKey1, privateKey2, nil

// }

func ConstantCommitment() core.G1Point {
	commitment := core.NewG1Point(big.NewInt(123), big.NewInt(456))
	return *commitment
}

// func TestReservationMetering(t *testing.T) {

// 	m, _, _, err := setupTestMeterer(t)
// 	assert.NoError(t, err)

// 	ctx, cancel := context.WithTimeout(context.Background(), m.TimeoutConfig.ChainReadTimeout)
// 	defer cancel()

// 	// retreiverID := "testRetriever"

// 	// params := []common.RequestParams{
// 	// 	{
// 	// 		RequesterID: retreiverID,
// 	// 		BlobSize:    10,
// 	// 		Rate:        100,
// 	// 	},
// 	// }
// 	header := meterer.BlobHeader{
// 		AccountID: "account1",
// 		BlobSize:  10,
// 		BinIndex:  0,
// 		Nonce:     0,
// 		Signature: []byte{},
// 		Version:   0,
// 	}
// 	activeReservation := meterer.ActiveReservation{
// 		StartEpoch: 0,
// 		EndEpoch:   10,
// 	}

// 	for i := 0; i < 10; i++ {
// 		err := m.ServeReservationRequest(ctx, header, &activeReservation)
// 		assert.NoError(t, err)
// 	}

// 	//TODO: overflown request
// 	err = m.ServeReservationRequest(ctx, header, &activeReservation)
// 	assert.Error(t, err)
// }

func TestMetererAcceptValidHeader(t *testing.T) {
	ctx := context.Background()
	CreateReservationTable(t, "reservations")
	// CreateOnDemandTable(t, "ondemand")
	// CreateGlobalReservationTable(t, "global")
	index := meterer.GetCurrentBinIndex()
	// get 32 raw bytes

	commitment := core.NewG1Point(big.NewInt(0), big.NewInt(1))

	fmt.Println("test invalid index")
	header, err := meterer.ConstructBlobHeader(signer, 1, 1, index, 0, *commitment, 2000, []meterer.BlobQuorumParam{}, privateKey1)
	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *header)
	assert.Error(t, err, "invalid bin index for reservation")

	header, err = meterer.ConstructBlobHeader(signer, 1, 1, index-1, 0, *commitment, 1000, []meterer.BlobQuorumParam{}, privateKey2)
	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *header)

	assert.Error(t, err, "invalid bin index for reservation")

	// test bin overflow
	fmt.Println("test bin usage")
	accountID := crypto.PubkeyToAddress(privateKey2.PublicKey).Hex()
	for i := 0; i < 9; i++ {
		blobSize := 20
		header, err = meterer.ConstructBlobHeader(signer, 1, 1, index, 0, *commitment, uint32(blobSize), []meterer.BlobQuorumParam{}, privateKey2)
		assert.NoError(t, err)
		err = mt.MeterRequest(ctx, *header)
		assert.NoError(t, err)
		item, err := dynamoClient.GetItem(ctx, "reservations", commondynamodb.Key{
			"AccountID": &types.AttributeValueMemberS{Value: accountID},
			"BinIndex":  &types.AttributeValueMemberN{Value: strconv.Itoa(int(index))},
		})
		assert.NoError(t, err)
		assert.Equal(t, accountID, item["AccountID"].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, strconv.Itoa(int(index)), item["BinIndex"].(*types.AttributeValueMemberN).Value)
		assert.Equal(t, strconv.Itoa(int((i+1)*blobSize)), item["BinUsage"].(*types.AttributeValueMemberN).Value)

	}
	// frist over flow is allowed
	header, err = meterer.ConstructBlobHeader(signer, 1, 1, index, 0, *commitment, 25, []meterer.BlobQuorumParam{}, privateKey2)
	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *header)
	assert.NoError(t, err)
	overflowedIndex := index + 2
	item, err := dynamoClient.GetItem(ctx, "reservations", commondynamodb.Key{
		"AccountID": &types.AttributeValueMemberS{Value: accountID},
		"BinIndex":  &types.AttributeValueMemberN{Value: strconv.Itoa(int(overflowedIndex))},
	})
	assert.NoError(t, err)
	assert.Equal(t, accountID, item["AccountID"].(*types.AttributeValueMemberS).Value)
	assert.Equal(t, strconv.Itoa(int(overflowedIndex)), item["BinIndex"].(*types.AttributeValueMemberN).Value)
	assert.Equal(t, strconv.Itoa(int(5)), item["BinUsage"].(*types.AttributeValueMemberN).Value)

	// second over flow
	header, err = meterer.ConstructBlobHeader(signer, 1, 1, index, 0, *commitment, 1, []meterer.BlobQuorumParam{}, privateKey2)
	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *header)
	assert.Error(t, err, "Bin has already been overflowed")

	// overwhelming bin overflow
	fmt.Println("test bin overflow")
	header, err = meterer.ConstructBlobHeader(signer, 1, 1, index-1, 0, *commitment, 1000, []meterer.BlobQuorumParam{}, privateKey2)
	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *header)
	assert.Error(t, err, "Overflow usage exceeds bin limit")
}

// func TestMetererRejectInvalidSignature(t *testing.T) {
// 	m, signer, privateKey, err := setupTestMeterer(t)
// 	assert.NoError(t, err)
// 	ctx := context.Background()

// 	header, err := meterer.ConstructBlobHeader(signer, 1, 1, 0, 1000, ConstantCommitment(), 1000, []meterer.BlobQuorumParam{}, privateKey)
// 	assert.NoError(t, err)
// 	header.Signature = []byte("invalid signature")
// 	err = m.MeterRequest(ctx, *header)
// 	assert.Error(t, err)
// }

// func TestMetererRejectInsufficientPayment(t *testing.T) {
// 	m, signer, privateKey, err := setupTestMeterer(t)
// 	assert.NoError(t, err)
// 	ctx := context.Background()

// 	header, err := meterer.ConstructBlobHeader(signer, 1, 1, 0, 500, core.G1Point{}, 1000, []meterer.BlobQuorumParam{}, privateKey)
// 	assert.NoError(t, err)
// 	err = m.MeterRequest(ctx, *header)
// 	assert.Error(t, err)
// }

// func TestMetererRejectInvalidBinIndex(t *testing.T) {
// 	m, signer, privateKey, err := setupTestMeterer(t)
// 	assert.NoError(t, err)
// 	ctx := context.Background()

// 	header, err := meterer.ConstructBlobHeader(signer, 1, 1, 0, 1000000, core.G1Point{}, 1000, []meterer.BlobQuorumParam{}, privateKey) // Invalid bin index
// 	assert.NoError(t, err)
// 	err = m.MeterRequest(ctx, *header)
// 	assert.Error(t, err)
// }

// func TestMetererAcceptMultipleValidHeaders(t *testing.T) {
// 	m, signer, privateKey, err := setupTestMeterer(t)
// 	assert.NoError(t, err)
// 	ctx := context.Background()

// 	header1, err := meterer.ConstructBlobHeader(signer, 1, 1, 0, 1000, core.G1Point{}, 1000, []meterer.BlobQuorumParam{}, privateKey)
// 	assert.NoError(t, err)
// 	err = m.MeterRequest(ctx, *header1)
// 	assert.NoError(t, err)

// 	header2, err := meterer.ConstructBlobHeader(signer, 1, 2, 1, 2000, core.G1Point{}, 1000, []meterer.BlobQuorumParam{}, privateKey)
// 	assert.NoError(t, err)
// 	err = m.MeterRequest(ctx, *header2)
// 	assert.NoError(t, err)
// }

// func TestMetererRejectDuplicateNonce(t *testing.T) {
// 	m, signer, privateKey, err := setupTestMeterer(t)
// 	assert.NoError(t, err)
// 	ctx := context.Background()

// 	header1, err := meterer.ConstructBlobHeader(signer, 1, 1, 0, 1000, core.G1Point{}, 1000, []meterer.BlobQuorumParam{}, privateKey)
// 	assert.NoError(t, err)
// 	err = m.MeterRequest(ctx, *header1)
// 	assert.NoError(t, err)

// 	header2, err := meterer.ConstructBlobHeader(signer, 1, 1, 1, 2000, core.G1Point{}, 1000, []meterer.BlobQuorumParam{}, privateKey) // Same nonce
// 	assert.NoError(t, err)
// 	err = m.MeterRequest(ctx, *header2)
// 	assert.Error(t, err)
// }

// func TestMetererRejectGlobalRateLimit(t *testing.T) {
// 	m, signer, privateKey, err := setupTestMeterer(t)
// 	assert.NoError(t, err)
// 	ctx := context.Background()

// 	// Submit requests up to the global rate limit
// 	for i := 0; i < 1000; i++ {
// 		header, err := meterer.ConstructBlobHeader(signer, 1, uint32(i), 0, uint64((i+1)*1000), core.G1Point{}, 1000, []meterer.BlobQuorumParam{}, privateKey)
// 		assert.NoError(t, err)
// 		err = m.MeterRequest(ctx, *header)
// 		assert.NoError(t, err)
// 	}

// 	// This request should be rejected due to exceeding the global rate limit
// 	header, err := meterer.ConstructBlobHeader(signer, 1, 1000, 0, 1001000, core.G1Point{}, 1000, []meterer.BlobQuorumParam{}, privateKey)
// 	assert.NoError(t, err)
// 	err = m.MeterRequest(ctx, *header)
// 	assert.Error(t, err)
// }

func CreateReservationTable(t *testing.T, tableName string) {
	ctx := context.Background()
	tableDescription, err := test_utils.CreateTable(ctx, clientConfig, tableName, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("AccountID"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("BinIndex"),
				AttributeType: types.ScalarAttributeTypeN,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("AccountID"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("BinIndex"),
				KeyType:       types.KeyTypeRange,
			},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("AccountIDIndex"),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("AccountID"),
						KeyType:       types.KeyTypeHash,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll, // ProjectionTypeAll means all attributes are projected into the index
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(10),
					WriteCapacityUnits: aws.Int64(10),
				},
			},
		},
		TableName: aws.String(tableName),
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, tableDescription)
}

func CreateGlobalReservationTable(t *testing.T, tableName string) {
	ctx := context.Background()
	tableDescription, err := test_utils.CreateTable(ctx, clientConfig, tableName, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("BinIndex"),
				AttributeType: types.ScalarAttributeTypeN,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("BinIndex"),
				KeyType:       types.KeyTypeHash,
			},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("BinIndexIndex"),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("BinIndex"),
						KeyType:       types.KeyTypeHash,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll, // ProjectionTypeAll means all attributes are projected into the index
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(10),
					WriteCapacityUnits: aws.Int64(10),
				},
			},
		},
		TableName: aws.String(tableName),
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, tableDescription)
}

func CreateOnDemandTable(t *testing.T, tableName string) {
	ctx := context.Background()
	tableDescription, err := test_utils.CreateTable(ctx, clientConfig, tableName, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("AccountID"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("CumulativePayments"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("AccountID"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("CumulativePayments"),
				KeyType:       types.KeyTypeRange,
			},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("AccountIDIndex"),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("AccountID"),
						KeyType:       types.KeyTypeHash,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll, // ProjectionTypeAll means all attributes are projected into the index
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(10),
					WriteCapacityUnits: aws.Int64(10),
				},
			},
		},
		TableName: aws.String(tableName),
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, tableDescription)
}
