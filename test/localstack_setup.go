package test

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"sync"

	commonaws "github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/test/testbed"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const (
	LocalstackPort = uint16(4573)
)

var (
	logger = GetLogger()
	lock   sync.Mutex
)

// DeployDynamoLocalstack deploys a Localstack DynamoDB instance for testing.
func DeployDynamoLocalstack(ctx context.Context) (func(), error) {
	lock.Lock()
	defer lock.Unlock()

	var localstackContainer *testbed.LocalStackContainer
	shouldTearDown := true

	// Unfortunately, we have to use environment variables to control this, because
	// tests across different packages do not share state. This may cause race conditions,
	// but that can't be helped with the current state of the testing framework.

	if os.Getenv("DEPLOY_LOCALSTACK") != "false" {
		var err error
		localstackContainer, err = testbed.NewLocalStackContainerWithOptions(ctx, testbed.LocalStackOptions{
			ExposeHostPort: true,
			HostPort:       fmt.Sprintf("%d", LocalstackPort),
			Services:       []string{"dynamodb"},
			Logger:         logger,
		})
		if err != nil {
			_ = localstackContainer.Terminate(ctx)
			logger.Fatal("Failed to start localstack container:", err)
		}
	} else {
		// assume localstack is already deployed
		shouldTearDown = false
	}

	if os.Getenv("CI") != "false" {
		// Special case: in CI environments, never tear down localstack.
		shouldTearDown = false
		_ = os.Setenv("DEPLOY_LOCALSTACK", "false")
	}

	return func() {
		if !shouldTearDown {
			// If localstack was not deployed here, do not tear down here either
			return
		}

		lock.Lock()
		defer lock.Unlock()

		// TODO: temporary thing to debug tests in CI
		fmt.Printf("calling cleanup function in GetOrDeployLocalstack\n")
		debug.PrintStack()

		_ = localstackContainer.Terminate(ctx)
	}, nil
}

// GetDynamoClient returns a DynamoDB client connected to Localstack for testing.
func GetDynamoClient() (*dynamodb.Client, error) {
	lock.Lock()
	defer lock.Unlock()

	clientConfig := commonaws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     fmt.Sprintf("http://0.0.0.0:%d", LocalstackPort),
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
	return dynamodb.NewFromConfig(awsConfig), nil
}
