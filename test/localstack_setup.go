package test

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Layr-Labs/eigenda/common"
	commonaws "github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/test/testbed"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const LocalstackPort = uint16(4573)

// DeployDynamoLocalstack deploys a Localstack DynamoDB instance for testing.
func DeployDynamoLocalstack(ctx context.Context) (func(), error) {

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	localstackContainer, err := testbed.NewLocalStackContainerWithOptions(ctx, testbed.LocalStackOptions{
		ExposeHostPort: true,
		HostPort:       fmt.Sprintf("%d", LocalstackPort),
		Services:       []string{"dynamodb"},
		Logger:         logger,
	})
	if err != nil {
		if strings.Contains(err.Error(), "port is already allocated") {
			// Assume localstack is already deployed
			logger.Warnf("Localstack port %d is already allocated, assuming Localstack is already running",
				LocalstackPort)
			return func() {}, nil
		} else {
			return nil, fmt.Errorf("failed to start localstack container: %w", err)
		}
	}

	return func() {
		if os.Getenv("CI") != "" {
			// Special case: in CI environments, never tear down localstack.
			return
		}

		_ = localstackContainer.Terminate(ctx)
	}, nil
}

// GetDynamoClient returns a DynamoDB client connected to Localstack for testing.
func GetDynamoClient() (*dynamodb.Client, error) {
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
