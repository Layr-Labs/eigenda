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
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser/meterer"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
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

	chainID := big.NewInt(17000)
	verifyingContract := gethcommon.HexToAddress("0x1234000000000000000000000000000000000000")
	signer = meterer.NewEIP712Signer(chainID, verifyingContract)

	privateKey1, err = crypto.GenerateKey()
	privateKey2, err = crypto.GenerateKey()

	logger = logging.NewNoopLogger()
	config := meterer.Config{
		PricePerChargeable:   1,
		MinChargeableSize:    1,
		GlobalBytesPerSecond: 1000,
		ReservationWindow:    time.Minute,
	}

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

func TestMetererReservations(t *testing.T) {
	ctx := context.Background()
	meterer.CreateReservationTable(clientConfig, "reservations")
	index := meterer.GetCurrentBinIndex()
	commitment := core.NewG1Point(big.NewInt(0), big.NewInt(1))

	// test invalid signature
	invalidHeader := &meterer.BlobHeader{
		Version:           1,
		AccountID:         crypto.PubkeyToAddress(privateKey1.PublicKey).Hex(),
		Nonce:             1,
		BinIndex:          uint64(time.Now().Unix()),
		CumulativePayment: 0,
		Commitment:        *commitment,
		BlobSize:          2000,
		BlobQuorumParams:  []meterer.BlobQuorumParam{},
		Signature:         []byte{78, 212, 55, 45, 156, 217, 21, 240, 47, 141, 18, 213, 226, 196, 4, 51, 245, 110, 20, 106, 244, 142, 142, 49, 213, 21, 34, 151, 118, 254, 46, 89, 48, 84, 250, 46, 179, 228, 46, 51, 106, 164, 122, 11, 26, 101, 10, 10, 243, 2, 30, 46, 95, 125, 189, 237, 236, 91, 130, 224, 240, 151, 106, 204, 1},
	}
	err := mt.MeterRequest(ctx, *invalidHeader)
	assert.Error(t, err, "invalid signature: recovered address * does not match account ID *")

	// test non-existent account
	unregisteredUser, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}
	header, err := meterer.ConstructBlobHeader(signer, 1, 1, 1, 0, *commitment, 1000, []meterer.BlobQuorumParam{}, unregisteredUser)
	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *header)
	assert.Error(t, err, "failed to get on-demand payment by account: reservation not found")

	// test invalid index
	header, err = meterer.ConstructBlobHeader(signer, 1, 1, index, 0, *commitment, 2000, []meterer.BlobQuorumParam{}, privateKey1)
	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *header)
	assert.Error(t, err, "invalid bin index for reservation")

	header, err = meterer.ConstructBlobHeader(signer, 1, 1, index-1, 0, *commitment, 1000, []meterer.BlobQuorumParam{}, privateKey2)
	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *header)
	assert.Error(t, err, "invalid bin index for reservation")

	// test bin usage
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
	header, err = meterer.ConstructBlobHeader(signer, 1, 1, index-1, 0, *commitment, 1000, []meterer.BlobQuorumParam{}, privateKey2)
	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *header)
	assert.Error(t, err, "Overflow usage exceeds bin limit")
}

func TestMetererOnDemand(t *testing.T) {
	ctx := context.Background()
	meterer.CreateOnDemandTable(clientConfig, "ondemand")
	meterer.CreateGlobalReservationTable(clientConfig, "global")
	commitment := core.NewG1Point(big.NewInt(0), big.NewInt(1))

	// test invalid signature
	invalidHeader := &meterer.BlobHeader{
		Version:           1,
		AccountID:         crypto.PubkeyToAddress(privateKey1.PublicKey).Hex(),
		Nonce:             1,
		BinIndex:          uint64(time.Now().Unix()),
		CumulativePayment: 1,
		Commitment:        *commitment,
		BlobSize:          2000,
		BlobQuorumParams:  []meterer.BlobQuorumParam{},
		Signature:         []byte{78, 212, 55, 45, 156, 217, 21, 240, 47, 141, 18, 213, 226, 196, 4, 51, 245, 110, 20, 106, 244, 142, 142, 49, 213, 21, 34, 151, 118, 254, 46, 89, 48, 84, 250, 46, 179, 228, 46, 51, 106, 164, 122, 11, 26, 101, 10, 10, 243, 2, 30, 46, 95, 125, 189, 237, 236, 91, 130, 224, 240, 151, 106, 204, 1},
	}
	err := mt.MeterRequest(ctx, *invalidHeader)
	assert.Error(t, err, "invalid signature: recovered address * does not match account ID *")

	// test unregistered account
	unregisteredUser, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}
	header, err := meterer.ConstructBlobHeader(signer, 1, 1, 1, 1, *commitment, 1000, []meterer.BlobQuorumParam{}, unregisteredUser)
	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *header)
	assert.Error(t, err, "failed to get on-demand payment by account: payment not found")

	// test insufficient cumulative payment
	header, err = meterer.ConstructBlobHeader(signer, 1, 1, 0, 1, *commitment, 2000, []meterer.BlobQuorumParam{}, privateKey1)
	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *header)
	assert.Error(t, err, "insufficient cumulative payment increment")
	// Correct rollback (TODO: discuss if there should be a rollback for invalid payments or penalize client)
	result, err := dynamoClient.QueryIndex(ctx, "ondemand", "AccountIDIndex", "AccountID = :account", commondynamodb.ExpresseionValues{
		":account": &types.AttributeValueMemberS{
			Value: crypto.PubkeyToAddress(privateKey1.PublicKey).Hex(),
		}})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(result))

	// test failed global bin index
	header, err = meterer.ConstructBlobHeader(signer, 1, 1, 0, 1, *commitment, 1, []meterer.BlobQuorumParam{}, privateKey1)
	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *header)
	assert.Error(t, err, "failed global bin index")
	// Correct rollback
	result, err = dynamoClient.QueryIndex(ctx, "ondemand", "AccountIDIndex", "AccountID = :account", commondynamodb.ExpresseionValues{
		":account": &types.AttributeValueMemberS{
			Value: crypto.PubkeyToAddress(privateKey1.PublicKey).Hex(),
		}})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(result))

	header, err = meterer.ConstructBlobHeader(signer, 1, 1, uint64(time.Now().Unix()), uint64(100), *commitment, 100, []meterer.BlobQuorumParam{}, privateKey1)
	err = mt.MeterRequest(ctx, *header)
	assert.NoError(t, err)
	header, err = meterer.ConstructBlobHeader(signer, 1, 1, uint64(time.Now().Unix()), uint64(100), *commitment, 100, []meterer.BlobQuorumParam{}, privateKey1)
	err = mt.MeterRequest(ctx, *header)
	assert.Error(t, err, "exact payment already exists")

	// test valid payments
	for i := 1; i < 10; i++ {
		header, err = meterer.ConstructBlobHeader(signer, 1, 1, uint64(time.Now().Unix()), uint64(100*(i+1)), *commitment, 100, []meterer.BlobQuorumParam{}, privateKey1)
		err = mt.MeterRequest(ctx, *header)
		assert.NoError(t, err)
	}

	// test insufficient remaining balance from cumulative payment
	header, err = meterer.ConstructBlobHeader(signer, 1, 1, uint64(time.Now().Unix()), 1, *commitment, 1, []meterer.BlobQuorumParam{}, privateKey1)
	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *header)
	assert.Error(t, err, "insufficient cumulative payment increment")

	// test cannot insert cumulative payment in out of order
	header, err = meterer.ConstructBlobHeader(signer, 1, 1, uint64(time.Now().Unix()), 0, *commitment, 50, []meterer.BlobQuorumParam{}, privateKey2)
	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *header)
	assert.Error(t, err, "cannot insert cumulative payment in out of order")

	// test failed global rate limit
	header, err = meterer.ConstructBlobHeader(signer, 1, 1, uint64(time.Now().Unix()), 1001, *commitment, 1001, []meterer.BlobQuorumParam{}, privateKey2)
	assert.NoError(t, err)
	err = mt.MeterRequest(ctx, *header)
	assert.Error(t, err, "failed global bin index")
	// Correct rollback
	result, err = dynamoClient.QueryIndex(ctx, "ondemand", "AccountIDIndex", "AccountID = :account", commondynamodb.ExpresseionValues{
		":account": &types.AttributeValueMemberS{
			Value: crypto.PubkeyToAddress(privateKey2.PublicKey).Hex(),
		}})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(result))
}
