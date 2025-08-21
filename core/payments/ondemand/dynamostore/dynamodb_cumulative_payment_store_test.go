package dynamostore_test

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"testing"

	commonaws "github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/core/payments/ondemand"
	"github.com/Layr-Labs/eigenda/core/payments/ondemand/dynamostore"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
)

var (
	dockertestPool     *dockertest.Pool
	dockertestResource *dockertest.Resource
	dynamoClient       *dynamodb.Client
	clientConfig       commonaws.ClientConfig

	// whether these test must deploy localstack
	deployLocalStack bool
	// if this test deploys localstack, do so with this port. otherwise, use the value in env var LOCALSTACK_PORT
	localStackPort = "4566"

	accountID = gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_LOCALSTACK") != "false" {
		// deploy localstack
		deployLocalStack = true
		var err error
		dockertestPool, dockertestResource, err = deploy.StartDockertestWithLocalstackContainer(localStackPort)
		if err != nil {
			teardown()
			panic("failed to start localstack container: " + err.Error())
		}
	} else {
		// localstack is already deployed
		localStackPort = os.Getenv("LOCALSTACK_PORT")
	}

	clientConfig = commonaws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     fmt.Sprintf("http://0.0.0.0:%s", localStackPort),
	}

	awsConfig := aws.Config{
		Region: clientConfig.Region,
		Credentials: aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     clientConfig.AccessKey,
				SecretAccessKey: clientConfig.SecretAccessKey,
			}, nil
		}),
		EndpointResolverWithOptions: aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				if clientConfig.EndpointURL != "" {
					return aws.Endpoint{
						PartitionID:   "aws",
						URL:           clientConfig.EndpointURL,
						SigningRegion: clientConfig.Region,
					}, nil
				}
				return aws.Endpoint{}, &aws.EndpointNotFoundError{}
			}),
	}
	dynamoClient = dynamodb.NewFromConfig(awsConfig)

	code := m.Run()
	teardown()
	os.Exit(code)
}

func teardown() {
	if deployLocalStack {
		deploy.PurgeDockertestResources(dockertestPool, dockertestResource)
	}
}

func createPaymentTable(t *testing.T, tableName string) {
	// Use the existing CreateOnDemandTable function from meterer package
	// This ensures our test table schema exactly matches the production schema
	err := meterer.CreateOnDemandTable(clientConfig, tableName)
	require.NoError(t, err)
}

func deleteTable(t *testing.T, tableName string) {
	ctx := context.Background()
	_, err := dynamoClient.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	})
	require.NoError(t, err)
}

func TestConstructor(t *testing.T) {
	tableName := "TestConstructor"

	store, err := dynamostore.NewDynamoDBCumulativePaymentStore(nil, tableName, accountID)
	require.Error(t, err, "nil client should error")
	require.Nil(t, store)

	store, err = dynamostore.NewDynamoDBCumulativePaymentStore(dynamoClient, "", accountID)
	require.Error(t, err, "empty table name should error")
	require.Nil(t, store)

	store, err = dynamostore.NewDynamoDBCumulativePaymentStore(dynamoClient, tableName, gethcommon.Address{})
	require.Error(t, err, "zero address should error")
	require.Nil(t, store)
}

func TestUpperBound(t *testing.T) {
	tableName := "UpperBound"
	createPaymentTable(t, tableName)
	defer deleteTable(t, tableName)

	store, err := dynamostore.NewDynamoDBCumulativePaymentStore(dynamoClient, tableName, accountID)
	require.NoError(t, err)

	ctx := context.Background()
	maxCumulativePayment := big.NewInt(1000)

	newValue, err := store.AddCumulativePayment(ctx, big.NewInt(200), maxCumulativePayment)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(200), newValue)

	newValue, err = store.AddCumulativePayment(ctx, big.NewInt(801), maxCumulativePayment)
	require.Error(t, err, "exceeding capacity shouldn't be possible")
	require.Nil(t, newValue)
	var insufficientFundsErr *ondemand.InsufficientFundsError
	require.ErrorAs(t, err, &insufficientFundsErr)
	require.Equal(t, big.NewInt(200), insufficientFundsErr.CurrentCumulativePayment)
	require.Equal(t, maxCumulativePayment, insufficientFundsErr.MaxCumulativePayment)
	require.Equal(t, big.NewInt(801), insufficientFundsErr.BlobCost)

	newValue, err = store.AddCumulativePayment(ctx, big.NewInt(800), maxCumulativePayment)
	require.NoError(t, err, "exact maximum shouldn't error")
	require.Equal(t, big.NewInt(1000), newValue)

	newValue, err = store.AddCumulativePayment(ctx, big.NewInt(1), maxCumulativePayment)
	require.Error(t, err, "even one more wei should fail")
	require.Nil(t, newValue)
}

func TestSingleDispersalExceedsMaximum(t *testing.T) {
	tableName := "TestSingleDispersalExceedsMaximum"
	createPaymentTable(t, tableName)
	defer deleteTable(t, tableName)

	store, err := dynamostore.NewDynamoDBCumulativePaymentStore(dynamoClient, tableName, accountID)
	require.NoError(t, err)

	ctx := context.Background()

	_, err = store.AddCumulativePayment(ctx, big.NewInt(2000), big.NewInt(1000))
	require.Error(t, err, "single payment exceeding total deposits should fail")
	var insufficientFundsErr *ondemand.InsufficientFundsError
	require.ErrorAs(t, err, &insufficientFundsErr)
}

func TestAddCumulativePaymentInputValidation(t *testing.T) {
	tableName := "AddCumulativePaymentInputValidation"
	createPaymentTable(t, tableName)
	defer deleteTable(t, tableName)

	store, err := dynamostore.NewDynamoDBCumulativePaymentStore(dynamoClient, tableName, accountID)
	require.NoError(t, err)

	ctx := context.Background()

	_, err = store.AddCumulativePayment(ctx, nil, big.NewInt(1000))
	require.Error(t, err, "nil amount should error")

	_, err = store.AddCumulativePayment(ctx, big.NewInt(0), big.NewInt(1000))
	require.Error(t, err, "zero amount should error")

	_, err = store.AddCumulativePayment(ctx, big.NewInt(-100), big.NewInt(1000))
	require.Error(t, err, "negative amount should error")

	_, err = store.AddCumulativePayment(ctx, big.NewInt(100), nil)
	require.Error(t, err, "nil max should error")

	_, err = store.AddCumulativePayment(ctx, big.NewInt(100), big.NewInt(0))
	require.Error(t, err, "zero max should error")

	_, err = store.AddCumulativePayment(ctx, big.NewInt(100), big.NewInt(-1000))
	require.Error(t, err, "negative max should error")
}

func TestLargeNumbers(t *testing.T) {
	tableName := "TestLargeNumbers"
	createPaymentTable(t, tableName)
	defer deleteTable(t, tableName)

	store, err := dynamostore.NewDynamoDBCumulativePaymentStore(dynamoClient, tableName, accountID)
	require.NoError(t, err)

	ctx := context.Background()

	// Test with very large numbers (typical Wei values)
	amount := new(big.Int)
	amount.SetString("1000000000000000000", 10) // 1 ETH in Wei

	maxCumulativePayment := new(big.Int)
	maxCumulativePayment.SetString("100000000000000000000", 10) // 100 ETH in Wei

	newValue, err := store.AddCumulativePayment(ctx, amount, maxCumulativePayment)
	require.NoError(t, err)
	require.Equal(t, amount, newValue)

	// Add more
	newValue, err = store.AddCumulativePayment(ctx, amount, maxCumulativePayment)
	require.NoError(t, err)
	expected := new(big.Int)
	expected.SetString("2000000000000000000", 10) // 2 ETH in Wei
	require.Equal(t, expected, newValue)
}
