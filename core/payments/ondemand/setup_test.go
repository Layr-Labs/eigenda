package ondemand_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	commonaws "github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/testbed"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/require"
)

const (
	// used if DEPLOY_LOCALSTACK != "false"
	defaultLocalstackPort = "4573"
)

var (
	dynamoClient *dynamodb.Client
)

// TestMain sets up Localstack/Dynamo for all tests in the ondemand package and tears down after.
func TestMain(m *testing.M) {
	localstackPort := defaultLocalstackPort

	var localstackContainer *testbed.LocalStackContainer
	var deployLocalStack bool

	if os.Getenv("DEPLOY_LOCALSTACK") != "false" {
		deployLocalStack = true
		var err error
		localstackContainer, err = testbed.NewLocalStackContainerWithOptions(context.Background(), testbed.LocalStackOptions{
			ExposeHostPort: true,
			HostPort:       localstackPort,
			Services:       []string{"dynamodb"},
		})
		if err != nil {
			_ = localstackContainer.Terminate(context.Background())
			panic("failed to start localstack container: " + err.Error())
		}
	} else {
		// localstack is already deployed
		localstackPort = os.Getenv("LOCALSTACK_PORT")
	}

	clientConfig := commonaws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     fmt.Sprintf("http://0.0.0.0:%s", localstackPort),
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
	if deployLocalStack {
		_ = localstackContainer.Terminate(context.Background())
	}
	os.Exit(code)
}

// createPaymentTable creates a DynamoDB table for on-demand payment testing
// Uses the existing CreateOnDemandTable function from meterer package to ensure
// our test table schema exactly matches the production schema.
// Appends a random suffix to the table name to prevent collisions between tests.
func createPaymentTable(t *testing.T, tableName string) string {
	testRandom := random.NewTestRandom()
	randomSuffix := testRandom.Intn(999999)
	fullTableName := fmt.Sprintf("%s_%d", tableName, randomSuffix)

	// Create local client config for table creation
	localstackPort := defaultLocalstackPort
	if os.Getenv("DEPLOY_LOCALSTACK") == "false" {
		localstackPort = os.Getenv("LOCALSTACK_PORT")
	}

	clientConfig := commonaws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     fmt.Sprintf("http://0.0.0.0:%s", localstackPort),
	}

	err := meterer.CreateOnDemandTable(clientConfig, fullTableName)
	require.NoError(t, err)

	return fullTableName
}

// deleteTable deletes a DynamoDB table used in testing
func deleteTable(t *testing.T, tableName string) {
	ctx := context.Background()
	_, err := dynamoClient.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	})
	require.NoError(t, err)
}
