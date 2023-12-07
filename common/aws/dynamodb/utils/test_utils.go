package test_utils

import (
	"context"
	"errors"
	"time"

	commonaws "github.com/Layr-Labs/eigenda/common/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	// waiterDuration is the duration to wait for a table to be created
	waiterDuration = 15 * time.Second
)

func CreateTable(ctx context.Context, cfg commonaws.ClientConfig, name string, input *dynamodb.CreateTableInput) (*types.TableDescription, error) {
	c, err := getClient(cfg)
	if err != nil {
		return nil, err
	}
	table, err := c.CreateTable(ctx, input)
	if err != nil {
		return nil, err
	}

	waiter := dynamodb.NewTableExistsWaiter(c)
	err = waiter.Wait(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(name),
	}, waiterDuration)
	if err != nil {
		return nil, err
	}

	return table.TableDescription, nil
}

func CreateTableIfNotExists(ctx context.Context, cfg commonaws.ClientConfig, name string, input *dynamodb.CreateTableInput) (*types.TableDescription, error) {
	c, err := getClient(cfg)
	if err != nil {
		return nil, err
	}

	// Check if the table already exists
	_, err = c.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(name),
	})

	// If the table does not exist, create it
	if err != nil {
		var aerr *types.ResourceNotFoundException
		if errors.As(err, &aerr) {
			// Table does not exist, so create it
			table, err := c.CreateTable(ctx, input)
			if err != nil {
				return nil, err
			}

			// Wait for the table to be created
			waiter := dynamodb.NewTableExistsWaiter(c)
			err = waiter.Wait(ctx, &dynamodb.DescribeTableInput{
				TableName: aws.String(name),
			}, waiterDuration)
			if err != nil {
				return nil, err
			}

			return table.TableDescription, nil
		} else {
			// Some other error occurred
			return nil, err
		}
	}

	// If the table exists, return its description
	tableDescription, err := c.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(name),
	})
	if err != nil {
		return nil, err
	}

	return tableDescription.Table, nil
}

func getClient(clientConfig commonaws.ClientConfig) (*dynamodb.Client, error) {
	createClient := func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if clientConfig.EndpointURL != "" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           clientConfig.EndpointURL,
				SigningRegion: clientConfig.Region,
			}, nil
		}

		// returning EndpointNotFoundError will allow the service to fallback to its default resolution
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	}
	customResolver := aws.EndpointResolverWithOptionsFunc(createClient)

	cfg, errCfg := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(clientConfig.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(clientConfig.AccessKey, clientConfig.SecretAccessKey, "")),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithRetryMode(aws.RetryModeStandard),
	)
	if errCfg != nil {
		return nil, errCfg
	}
	return dynamodb.NewFromConfig(cfg), nil
}
