package ondemand_test

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"testing"

	commonaws "github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/core/payments/ondemand"
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

	dynamoClient *dynamodb.Client
	clientConfig commonaws.ClientConfig

	// whether these tests must deploy localstack
	deployLocalStack bool
	// if this test deploys localstack, do so with this port. otherwise, use the value in env var LOCALSTACK_PORT
	localStackPort = "4566"

	accountID = gethcommon.HexToAddress("0x1234567890123456789012345678901234567890")
)

// TestMain sets up Localstack/Dynamo for tests and tears down after.
func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_LOCALSTACK") != "false" {
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

	store, err := ondemand.NewCumulativePaymentStore(nil, tableName, accountID)
	require.Error(t, err, "nil client should error")
	require.Nil(t, store)

	store, err = ondemand.NewCumulativePaymentStore(dynamoClient, "", accountID)
	require.Error(t, err, "empty table name should error")
	require.Nil(t, store)

	store, err = ondemand.NewCumulativePaymentStore(dynamoClient, tableName, gethcommon.Address{})
	require.Error(t, err, "zero address should error")
	require.Nil(t, store)
}

func TestStoreCumulativePaymentInputValidation(t *testing.T) {
	tableName := "StoreInputValidation"
	createPaymentTable(t, tableName)
	defer deleteTable(t, tableName)

	store, err := ondemand.NewCumulativePaymentStore(dynamoClient, tableName, accountID)
	require.NoError(t, err)

	ctx := context.Background()

	err = store.StoreCumulativePayment(ctx, nil)
	require.Error(t, err, "nil amount should error")

	err = store.StoreCumulativePayment(ctx, big.NewInt(-100))
	require.Error(t, err, "negative amount should error")
}

func TestStoreThenGet(t *testing.T) {
	tableName := "StoreThenGet"
	createPaymentTable(t, tableName)
	defer deleteTable(t, tableName)

	store, err := ondemand.NewCumulativePaymentStore(dynamoClient, tableName, accountID)
	require.NoError(t, err)
	ctx := context.Background()

	value, err := store.GetCumulativePayment(ctx)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(0), value, "get when missing should return 0")

	require.NoError(t, store.StoreCumulativePayment(ctx, big.NewInt(100)))
	value, err = store.GetCumulativePayment(ctx)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(100), value)

	require.NoError(t, store.StoreCumulativePayment(ctx, big.NewInt(200)))
	value, err = store.GetCumulativePayment(ctx)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(200), value)

	require.NoError(t, store.StoreCumulativePayment(ctx, big.NewInt(50)))
	value, err = store.GetCumulativePayment(ctx)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(50), value)

}

func TestDifferentAddresses(t *testing.T) {
	tableName := "DifferentAddresses"
	createPaymentTable(t, tableName)
	defer deleteTable(t, tableName)

	accountA := gethcommon.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	accountB := gethcommon.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")

	storeA, err := ondemand.NewCumulativePaymentStore(dynamoClient, tableName, accountA)
	require.NoError(t, err)
	storeB, err := ondemand.NewCumulativePaymentStore(dynamoClient, tableName, accountB)
	require.NoError(t, err)

	ctx := context.Background()
	require.NoError(t, storeA.StoreCumulativePayment(ctx, big.NewInt(100)))
	require.NoError(t, storeB.StoreCumulativePayment(ctx, big.NewInt(300)))

	valueA, err := storeA.GetCumulativePayment(ctx)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(100), valueA)

	valueB, err := storeB.GetCumulativePayment(ctx)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(300), valueB)
}
