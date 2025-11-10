package test

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"sync"

	commonaws "github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/test/testbed"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const (
	// used if DEPLOY_LOCALSTACK != "false"
	DefaultLocalstackPort = uint16(4573)
)

var (
	logger         = GetLogger()
	lock           sync.Mutex
	localstackPort uint16
	deployed       bool
)

// DeployDynamoLocalstack deploys a Localstack DynamoDB instance for testing.
func DeployDynamoLocalstack() (func(), error) {
	lock.Lock()
	defer lock.Unlock()

	if deployed {
		// Already deployed somewhere else
		return func() {

		}, nil
	}

	localstackPort = DefaultLocalstackPort

	var localstackContainer *testbed.LocalStackContainer

	ctx := context.Background()

	if os.Getenv("DEPLOY_LOCALSTACK") != "false" {
		var err error
		localstackContainer, err = testbed.NewLocalStackContainerWithOptions(ctx, testbed.LocalStackOptions{
			ExposeHostPort: true,
			HostPort:       fmt.Sprintf("%d", localstackPort),
			Services:       []string{"dynamodb"},
			Logger:         logger,
		})
		if err != nil {
			_ = localstackContainer.Terminate(ctx)
			logger.Fatal("Failed to start localstack container:", err)
		}
	} else {
		// assume localstack is already deployed
		port, err := strconv.ParseUint(os.Getenv("LOCALSTACK_PORT"), 10, 16)
		if err != nil {
			logger.Fatal("Failed to parse LOCALSTACK_PORT:", err)
		}
		localstackPort = uint16(port)
	}

	deployed = true

	return func() {
		lock.Lock()
		defer lock.Unlock()

		// TODO: temporary thing to debug tests in CI
		fmt.Printf("calling cleanup function in GetOrDeployLocalstack\n")
		debug.PrintStack()

		_ = localstackContainer.Terminate(ctx)
		deployed = false
	}, nil
}

// GetDynamoClient returns a DynamoDB client connected to Localstack for testing.
func GetDynamoClient() (*dynamodb.Client, error) {
	lock.Lock()
	defer lock.Unlock()

	if !deployed {
		return nil, fmt.Errorf("localstack not deployed; call DeployDynamoLocalstack first")
	}

	clientConfig := commonaws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     fmt.Sprintf("http://0.0.0.0:%d", localstackPort),
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

// GetLocalstackPort returns the port number where Localstack is running.
func GetLocalstackPort() uint16 {
	lock.Lock()
	defer lock.Unlock()
	return localstackPort
}
