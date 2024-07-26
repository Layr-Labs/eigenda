package dynamodb

import (
	"context"
	"fmt"
	"math"
	"sync"

	commonaws "github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	// dynamoBatchLimit is the maximum number of items that can be written in a single batch
	dynamoBatchLimit = 25
)

type batchOperation uint

const (
	update batchOperation = iota
	delete
)

var (
	once      sync.Once
	clientRef *Client
)

type Item = map[string]types.AttributeValue
type Key = map[string]types.AttributeValue
type ExpresseionValues = map[string]types.AttributeValue

type QueryResult struct {
	Items            []Item
	LastEvaluatedKey Key
}

type Client struct {
	dynamoClient *dynamodb.Client
	logger       logging.Logger
}

func NewClient(cfg commonaws.ClientConfig, logger logging.Logger) (*Client, error) {
	var err error
	once.Do(func() {
		createClient := func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if cfg.EndpointURL != "" {
				return aws.Endpoint{
					PartitionID:   "aws",
					URL:           cfg.EndpointURL,
					SigningRegion: cfg.Region,
				}, nil
			}

			// returning EndpointNotFoundError will allow the service to fallback to its default resolution
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		}
		customResolver := aws.EndpointResolverWithOptionsFunc(createClient)

		options := [](func(*config.LoadOptions) error){
			config.WithRegion(cfg.Region),
			config.WithEndpointResolverWithOptions(customResolver),
			config.WithRetryMode(aws.RetryModeStandard),
		}
		// If access key and secret access key are not provided, use the default credential provider
		if len(cfg.AccessKey) > 0 && len(cfg.SecretAccessKey) > 0 {
			options = append(options, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretAccessKey, "")))
		}
		awsConfig, errCfg := config.LoadDefaultConfig(context.Background(), options...)

		if errCfg != nil {
			err = errCfg
			return
		}
		dynamoClient := dynamodb.NewFromConfig(awsConfig)
		clientRef = &Client{dynamoClient: dynamoClient, logger: logger.With("component", "DynamodbClient")}
	})
	return clientRef, err
}

func (c *Client) DeleteTable(ctx context.Context, tableName string) error {
	_, err := c.dynamoClient.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName)})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) PutItem(ctx context.Context, tableName string, item Item) (err error) {
	_, err = c.dynamoClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName), Item: item,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) PutItemWithCondition(ctx context.Context, tableName string, item Item, condition string) (err error) {
	_, err = c.dynamoClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName), Item: item,
		ConditionExpression: aws.String(condition),
	})
	if err != nil {
		return err
	}

	return nil
}

// PutItems puts items in batches of 25 items (which is a limit DynamoDB imposes)
// It returns the items that failed to be put.
func (c *Client) PutItems(ctx context.Context, tableName string, items []Item) ([]Item, error) {
	return c.writeItems(ctx, tableName, items, update)
}

func (c *Client) UpdateItem(ctx context.Context, tableName string, key Key, item Item) (Item, error) {
	update := expression.UpdateBuilder{}
	for itemKey, itemValue := range item {
		if _, ok := key[itemKey]; ok {
			// Cannot update the key
			continue
		}
		update = update.Set(expression.Name(itemKey), expression.Value(itemValue))
	}

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return nil, err
	}

	resp, err := c.dynamoClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       key,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		ReturnValues:              types.ReturnValueUpdatedNew,
	})

	if err != nil {
		return nil, err
	}

	return resp.Attributes, err
}

func (c *Client) GetItem(ctx context.Context, tableName string, key Key) (Item, error) {
	resp, err := c.dynamoClient.GetItem(ctx, &dynamodb.GetItemInput{Key: key, TableName: aws.String(tableName)})
	if err != nil {
		return nil, err
	}

	return resp.Item, nil
}

// QueryIndex returns all items in the index that match the given key
func (c *Client) QueryIndex(ctx context.Context, tableName string, indexName string, keyCondition string, expAttributeValues ExpresseionValues) ([]Item, error) {
	response, err := c.dynamoClient.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		IndexName:                 aws.String(indexName),
		KeyConditionExpression:    aws.String(keyCondition),
		ExpressionAttributeValues: expAttributeValues,
	})
	if err != nil {
		return nil, err
	}

	return response.Items, nil
}

// Query returns all items in the primary index that match the given expression
func (c *Client) Query(ctx context.Context, tableName string, keyCondition string, expAttributeValues ExpresseionValues) ([]Item, error) {
	response, err := c.dynamoClient.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		KeyConditionExpression:    aws.String(keyCondition),
		ExpressionAttributeValues: expAttributeValues,
	})
	if err != nil {
		return nil, err
	}

	return response.Items, nil
}

// QueryIndexCount returns the count of the items in the index that match the given key
func (c *Client) QueryIndexCount(ctx context.Context, tableName string, indexName string, keyCondition string, expAttributeValues ExpresseionValues) (int32, error) {
	response, err := c.dynamoClient.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		IndexName:                 aws.String(indexName),
		KeyConditionExpression:    aws.String(keyCondition),
		ExpressionAttributeValues: expAttributeValues,
		Select:                    types.SelectCount,
	})
	if err != nil {
		return 0, err
	}

	return response.Count, nil
}

// QueryIndexWithPagination returns all items in the index that match the given key
// Results are limited to the given limit and the pagination token is returned
// When limit is 0, all items are returned
func (c *Client) QueryIndexWithPagination(ctx context.Context, tableName string, indexName string, keyCondition string, expAttributeValues ExpresseionValues, limit int32, exclusiveStartKey map[string]types.AttributeValue) (QueryResult, error) {
	var queryInput *dynamodb.QueryInput

	// Fetch all items if limit is 0
	if limit > 0 {
		queryInput = &dynamodb.QueryInput{
			TableName:                 aws.String(tableName),
			IndexName:                 aws.String(indexName),
			KeyConditionExpression:    aws.String(keyCondition),
			ExpressionAttributeValues: expAttributeValues,
			Limit:                     &limit,
		}
	} else {
		queryInput = &dynamodb.QueryInput{
			TableName:                 aws.String(tableName),
			IndexName:                 aws.String(indexName),
			KeyConditionExpression:    aws.String(keyCondition),
			ExpressionAttributeValues: expAttributeValues,
		}
	}

	// If a pagination token was provided, set it as the ExclusiveStartKey
	if exclusiveStartKey != nil {
		queryInput.ExclusiveStartKey = exclusiveStartKey
	}

	response, err := c.dynamoClient.Query(ctx, queryInput)
	if err != nil {
		return QueryResult{}, err
	}

	if len(response.Items) == 0 {
		return QueryResult{Items: nil, LastEvaluatedKey: nil}, nil
	}

	// Return the items and the pagination token
	return QueryResult{
		Items:            response.Items,
		LastEvaluatedKey: response.LastEvaluatedKey,
	}, nil
}

func (c *Client) DeleteItem(ctx context.Context, tableName string, key Key) error {
	_, err := c.dynamoClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{Key: key, TableName: aws.String(tableName)})
	if err != nil {
		return err
	}

	return nil
}

// DeleteItems deletes items in batches of 25 items (which is a limit DynamoDB imposes)
// It returns the items that failed to be deleted.
func (c *Client) DeleteItems(ctx context.Context, tableName string, keys []Key) ([]Key, error) {
	return c.writeItems(ctx, tableName, keys, delete)
}

// writeItems writes items in batches of 25 items (which is a limit DynamoDB imposes)
// update and delete operations are supported.
// For update operation, requestItems is []Item.
// For delete operation, requestItems is []Key.
func (c *Client) writeItems(ctx context.Context, tableName string, requestItems []map[string]types.AttributeValue, operation batchOperation) ([]map[string]types.AttributeValue, error) {
	startIndex := 0
	failedItems := make([]map[string]types.AttributeValue, 0)
	for startIndex < len(requestItems) {
		remainingNumKeys := float64(len(requestItems) - startIndex)
		batchSize := int(math.Min(float64(dynamoBatchLimit), remainingNumKeys))
		writeRequests := make([]types.WriteRequest, batchSize)
		for i := 0; i < batchSize; i += 1 {
			item := requestItems[startIndex+i]
			if operation == update {
				writeRequests[i] = types.WriteRequest{PutRequest: &types.PutRequest{Item: item}}
			} else if operation == delete {
				writeRequests[i] = types.WriteRequest{DeleteRequest: &types.DeleteRequest{Key: item}}
			} else {
				return nil, fmt.Errorf("unknown batch operation: %d", operation)
			}
		}
		// write batch
		output, err := c.dynamoClient.BatchWriteItem(
			ctx,
			&dynamodb.BatchWriteItemInput{
				RequestItems: map[string][]types.WriteRequest{tableName: writeRequests},
			},
		)
		if err != nil {
			return nil, err
		}

		// check for unprocessed items
		if len(output.UnprocessedItems) > 0 {
			for _, req := range output.UnprocessedItems[tableName] {
				failedItems = append(failedItems, req.DeleteRequest.Key)
			}
		}

		startIndex += dynamoBatchLimit
	}

	return failedItems, nil
}
