package test

import (
	"context"
	"fmt"
	"os"
	"sync"

	commonaws "github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/test/testbed"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const (
	// used if DEPLOY_LOCALSTACK != "false"
	DefaultLocalstackPort = "4573"
)

var (
	logger        = GetLogger()
	dynamoClientX *dynamodb.Client
	lock          sync.Mutex
)

// GetOrDeployLocalstack deploys a Localstack DynamoDB instance for testing,
// returning existing client if already deployed. Returns a function that should be
// called to clean up resources after tests complete.
func GetOrDeployLocalstack() (*dynamodb.Client, func()) {
	lock.Lock()
	defer lock.Unlock()

	if dynamoClientX != nil {
		return dynamoClientX, func() {
			// If already deployed somewhere else, no cleanup needed here.
		}
	}

	localstackPort := DefaultLocalstackPort

	var localstackContainer *testbed.LocalStackContainer

	ctx := context.Background()

	if os.Getenv("DEPLOY_LOCALSTACK") != "false" {
		var err error
		localstackContainer, err = testbed.NewLocalStackContainerWithOptions(ctx, testbed.LocalStackOptions{
			ExposeHostPort: true,
			HostPort:       localstackPort,
			Services:       []string{"dynamodb"},
			Logger:         logger,
		})
		if err != nil {
			_ = localstackContainer.Terminate(ctx)
			logger.Fatal("Failed to start localstack container:", err)
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
	dynamoClientX = dynamodb.NewFromConfig(awsConfig)

	return dynamoClientX, func() {
		lock.Lock()
		defer lock.Unlock()
		_ = localstackContainer.Terminate(ctx)
		dynamoClientX = nil
	}
}
