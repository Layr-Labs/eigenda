package dynamodb

import (
	"context"
	"errors"
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
	// DynamoBatchWriteLimit is the maximum number of items that can be written in a single batch
	// Reference: https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_TransactWriteItems.html
	DynamoBatchWriteLimit = 25
	// DynamoBatchReadLimit is the maximum number of items that can be read in a single batch
	DynamoBatchReadLimit = 100
)

type batchOperation uint

const (
	update batchOperation = iota
	delete
)

var (
	once               sync.Once
	clientRef          *client
	ErrConditionFailed = errors.New("condition failed")
)

type Item = map[string]types.AttributeValue
type Key = map[string]types.AttributeValue
type ExpressionValues = map[string]types.AttributeValue

type QueryResult struct {
	Items            []Item
	LastEvaluatedKey Key
}

// TransactAddOp defines an operation for TransactAddBy
// Value can be positive (increment) or negative (decrement)
type TransactAddOp struct {
	Key   Key
	Attr  string
	Value float64
}

type Client interface {
	DeleteTable(ctx context.Context, tableName string) error
	PutItem(ctx context.Context, tableName string, item Item) error
	PutItemWithCondition(ctx context.Context, tableName string, item Item, condition string, expressionAttributeNames map[string]string, expressionAttributeValues map[string]types.AttributeValue) error
	PutItemWithConditionAndReturn(ctx context.Context, tableName string, item Item, condition string, expressionAttributeNames map[string]string, expressionAttributeValues map[string]types.AttributeValue) (Item, error)
	PutItems(ctx context.Context, tableName string, items []Item) ([]Item, error)
	UpdateItem(ctx context.Context, tableName string, key Key, item Item) (Item, error)
	UpdateItemWithCondition(ctx context.Context, tableName string, key Key, item Item, condition expression.ConditionBuilder) (Item, error)
	IncrementBy(ctx context.Context, tableName string, key Key, attr string, value uint64) (Item, error)
	GetItem(ctx context.Context, tableName string, key Key) (Item, error)
	GetItemWithInput(ctx context.Context, input *dynamodb.GetItemInput) (Item, error)
	GetItems(ctx context.Context, tableName string, keys []Key, consistentRead bool) ([]Item, error)
	QueryIndex(ctx context.Context, tableName string, indexName string, keyCondition string, expAttributeValues ExpressionValues) ([]Item, error)
	Query(ctx context.Context, tableName string, keyCondition string, expAttributeValues ExpressionValues) ([]Item, error)
	QueryWithInput(ctx context.Context, input *dynamodb.QueryInput) ([]Item, error)
	QueryIndexCount(ctx context.Context, tableName string, indexName string, keyCondition string, expAttributeValues ExpressionValues) (int32, error)
	QueryIndexWithPagination(ctx context.Context, tableName string, indexName string, keyCondition string, expAttributeValues ExpressionValues, limit int32, exclusiveStartKey map[string]types.AttributeValue, ascending bool) (QueryResult, error)
	DeleteItem(ctx context.Context, tableName string, key Key) error
	DeleteItems(ctx context.Context, tableName string, keys []Key) ([]Key, error)
	TableExists(ctx context.Context, name string) error
	TransactAddBy(ctx context.Context, tableName string, ops []TransactAddOp) error
}

type client struct {
	dynamoClient *dynamodb.Client
	logger       logging.Logger
}

var _ Client = (*client)(nil)

func NewClient(cfg commonaws.ClientConfig, logger logging.Logger) (*client, error) {
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
		clientRef = &client{dynamoClient: dynamoClient, logger: logger.With("component", "DynamodbClient")}
	})
	return clientRef, err
}

func (c *client) DeleteTable(ctx context.Context, tableName string) error {
	_, err := c.dynamoClient.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName)})
	if err != nil {
		return fmt.Errorf("failed to delete table %s: %w", tableName, err)
	}
	return nil
}

func (c *client) PutItem(ctx context.Context, tableName string, item Item) (err error) {
	_, err = c.dynamoClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName), Item: item,
	})
	if err != nil {
		return fmt.Errorf("failed to put item in table %s: %w", tableName, err)
	}
	return nil
}

func (c *client) PutItemWithCondition(
	ctx context.Context,
	tableName string,
	item Item,
	condition string,
	expressionAttributeNames map[string]string,
	expressionAttributeValues map[string]types.AttributeValue,
) (err error) {
	_, err = c.dynamoClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName), Item: item,
		ConditionExpression:       aws.String(condition),
		ExpressionAttributeNames:  expressionAttributeNames,
		ExpressionAttributeValues: expressionAttributeValues,
	})
	var ccfe *types.ConditionalCheckFailedException
	if errors.As(err, &ccfe) {
		return ErrConditionFailed
	}
	if err != nil {
		return fmt.Errorf("failed to put item in table %s: %w", tableName, err)
	}
	return nil
}

// PutItemWithConditionAndReturn puts an item in the table with a condition and returns the old item if it exists
func (c *client) PutItemWithConditionAndReturn(
	ctx context.Context,
	tableName string,
	item Item,
	condition string,
	expressionAttributeNames map[string]string,
	expressionAttributeValues map[string]types.AttributeValue,
) (Item, error) {
	result, err := c.dynamoClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName), Item: item,
		ConditionExpression:       aws.String(condition),
		ExpressionAttributeNames:  expressionAttributeNames,
		ExpressionAttributeValues: expressionAttributeValues,
		ReturnValues:              types.ReturnValueAllOld,
	})
	var ccfe *types.ConditionalCheckFailedException
	if errors.As(err, &ccfe) {
		return nil, ErrConditionFailed
	}
	if err != nil {
		return nil, fmt.Errorf("failed to put item in table %s: %w", tableName, err)
	}

	return result.Attributes, nil
}

// PutItems puts items in batches of 25 items (which is a limit DynamoDB imposes)
// It returns the items that failed to be put.
func (c *client) PutItems(ctx context.Context, tableName string, items []Item) ([]Item, error) {
	return c.writeItems(ctx, tableName, items, update)
}

func (c *client) UpdateItem(ctx context.Context, tableName string, key Key, item Item) (Item, error) {
	update := expression.UpdateBuilder{}
	for itemKey, itemValue := range item {
		// Ignore primary key updates
		if _, ok := key[itemKey]; ok {
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

func (c *client) UpdateItemWithCondition(
	ctx context.Context,
	tableName string,
	key Key,
	item Item,
	condition expression.ConditionBuilder,
) (Item, error) {
	update := expression.UpdateBuilder{}
	for itemKey, itemValue := range item {
		// Ignore primary key updates
		if _, ok := key[itemKey]; ok {
			continue
		}
		update = update.Set(expression.Name(itemKey), expression.Value(itemValue))
	}

	expr, err := expression.NewBuilder().WithUpdate(update).WithCondition(condition).Build()
	if err != nil {
		return nil, err
	}

	resp, err := c.dynamoClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       key,
		ConditionExpression:       expr.Condition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		ReturnValues:              types.ReturnValueUpdatedNew,
	})

	var ccfe *types.ConditionalCheckFailedException
	if errors.As(err, &ccfe) {
		return nil, ErrConditionFailed
	}

	if err != nil {
		return nil, err
	}

	return resp.Attributes, err
}

// IncrementBy increments the attribute by the value for item that matches with the key
func (c *client) IncrementBy(ctx context.Context, tableName string, key Key, attr string, value uint64) (Item, error) {
	// ADD numeric values; small amounts of precision loss if the uint64 value is large and cannot be representing as a float64.
	// We don't expect such a large value to be incremented as it is used in units of dispersed symbols.
	update := expression.UpdateBuilder{}
	update = update.Add(expression.Name(attr), expression.Value(aws.Float64(float64(value))))
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

	return resp.Attributes, nil
}

func (c *client) GetItem(ctx context.Context, tableName string, key Key) (Item, error) {
	resp, err := c.dynamoClient.GetItem(ctx, &dynamodb.GetItemInput{Key: key, TableName: aws.String(tableName)})
	if err != nil {
		return nil, err
	}

	return resp.Item, nil
}

// GetItemWithInput is a wrapper for the GetItem function that allows for a custom GetItemInput
func (c *client) GetItemWithInput(ctx context.Context, input *dynamodb.GetItemInput) (Item, error) {
	resp, err := c.dynamoClient.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	return resp.Item, nil
}

// GetItems returns the items for the given keys
// Note: ordering of items is not guaranteed
func (c *client) GetItems(ctx context.Context, tableName string, keys []Key, consistentRead bool) ([]Item, error) {
	items, err := c.readItems(ctx, tableName, keys, consistentRead)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// QueryIndex returns all items in the index that match the given key
func (c *client) QueryIndex(ctx context.Context, tableName string, indexName string, keyCondition string, expAttributeValues ExpressionValues) ([]Item, error) {
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
func (c *client) Query(ctx context.Context, tableName string, keyCondition string, expAttributeValues ExpressionValues) ([]Item, error) {
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

// QueryWithInput is a wrapper for the Query function that allows for a custom query input
func (c *client) QueryWithInput(ctx context.Context, input *dynamodb.QueryInput) ([]Item, error) {
	response, err := c.dynamoClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}
	return response.Items, nil
}

// QueryIndexCount returns the count of the items in the index that match the given key
func (c *client) QueryIndexCount(ctx context.Context, tableName string, indexName string, keyCondition string, expAttributeValues ExpressionValues) (int32, error) {
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
func (c *client) QueryIndexWithPagination(ctx context.Context, tableName string, indexName string, keyCondition string, expAttributeValues ExpressionValues, limit int32, exclusiveStartKey map[string]types.AttributeValue, ascending bool) (QueryResult, error) {
	var queryInput *dynamodb.QueryInput

	// Fetch all items if limit is 0
	if limit > 0 {
		queryInput = &dynamodb.QueryInput{
			TableName:                 aws.String(tableName),
			IndexName:                 aws.String(indexName),
			KeyConditionExpression:    aws.String(keyCondition),
			ExpressionAttributeValues: expAttributeValues,
			Limit:                     &limit,
			ScanIndexForward:          aws.Bool(ascending),
		}
	} else {
		queryInput = &dynamodb.QueryInput{
			TableName:                 aws.String(tableName),
			IndexName:                 aws.String(indexName),
			KeyConditionExpression:    aws.String(keyCondition),
			ExpressionAttributeValues: expAttributeValues,
			ScanIndexForward:          aws.Bool(ascending),
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

func (c *client) DeleteItem(ctx context.Context, tableName string, key Key) error {
	_, err := c.dynamoClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{Key: key, TableName: aws.String(tableName)})
	if err != nil {
		return err
	}

	return nil
}

// DeleteItems deletes items in batches of 25 items (which is a limit DynamoDB imposes)
// It returns the items that failed to be deleted.
func (c *client) DeleteItems(ctx context.Context, tableName string, keys []Key) ([]Key, error) {
	return c.writeItems(ctx, tableName, keys, delete)
}

// writeItems writes items in batches of 25 items (which is a limit DynamoDB imposes)
// update and delete operations are supported.
// For update operation, requestItems is []Item.
// For delete operation, requestItems is []Key.
func (c *client) writeItems(ctx context.Context, tableName string, requestItems []map[string]types.AttributeValue, operation batchOperation) ([]map[string]types.AttributeValue, error) {
	startIndex := 0
	failedItems := make([]map[string]types.AttributeValue, 0)
	for startIndex < len(requestItems) {
		remainingNumKeys := float64(len(requestItems) - startIndex)
		batchSize := int(math.Min(float64(DynamoBatchWriteLimit), remainingNumKeys))
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

		startIndex += DynamoBatchWriteLimit
	}

	return failedItems, nil
}

func (c *client) readItems(
	ctx context.Context,
	tableName string,
	keys []Key,
	consistentRead bool,
) ([]Item, error) {
	startIndex := 0
	items := make([]Item, 0)
	for startIndex < len(keys) {
		remainingNumKeys := float64(len(keys) - startIndex)
		batchSize := int(math.Min(float64(DynamoBatchReadLimit), remainingNumKeys))
		keysBatch := keys[startIndex : startIndex+batchSize]
		output, err := c.dynamoClient.BatchGetItem(ctx, &dynamodb.BatchGetItemInput{
			RequestItems: map[string]types.KeysAndAttributes{
				tableName: {
					Keys:           keysBatch,
					ConsistentRead: aws.Bool(consistentRead),
				},
			},
		})
		if err != nil {
			return nil, err
		}

		if len(output.Responses) > 0 {
			for _, resp := range output.Responses {
				items = append(items, resp...)
			}
		}

		if output.UnprocessedKeys != nil {
			keys = append(keys, output.UnprocessedKeys[tableName].Keys...)
		}

		startIndex += batchSize
	}

	return items, nil
}

// TableExists checks if a table exists and can be described
func (c *client) TableExists(ctx context.Context, name string) error {
	if name == "" {
		return errors.New("table name is empty")
	}
	_, err := c.dynamoClient.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(name),
	})
	if err != nil {
		return err
	}
	return nil
}

// TransactAddBy performs atomic add (increment or decrement) on multiple items using DynamoDB's TransactWriteItems API.
// Each operation is specified by a key, attribute name, and value (positive for increment, negative for decrement).
// All operations are performed atomically; if any fail, none are applied.
// Uses TransactAddOp struct.
func (c *client) TransactAddBy(ctx context.Context, tableName string, ops []TransactAddOp) error {
	if len(ops) == 0 {
		return nil
	}
	if len(ops) > DynamoBatchWriteLimit {
		return fmt.Errorf("DynamoDB TransactWriteItems limit is %d operations per transaction", DynamoBatchWriteLimit)
	}

	transactItems := make([]types.TransactWriteItem, len(ops))
	for i, op := range ops {
		update := expression.UpdateBuilder{}
		update = update.Add(expression.Name(op.Attr), expression.Value(aws.Float64(op.Value)))
		expr, err := expression.NewBuilder().WithUpdate(update).Build()
		if err != nil {
			return fmt.Errorf("failed to build update expression: %w", err)
		}
		transactItems[i] = types.TransactWriteItem{
			Update: &types.Update{
				TableName:                           aws.String(tableName),
				Key:                                 op.Key,
				UpdateExpression:                    expr.Update(),
				ExpressionAttributeNames:            expr.Names(),
				ExpressionAttributeValues:           expr.Values(),
				ReturnValuesOnConditionCheckFailure: types.ReturnValuesOnConditionCheckFailureAllOld,
			},
		}
	}

	_, err := c.dynamoClient.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	})
	if err != nil {
		return fmt.Errorf("TransactWriteItems failed: %w", err)
	}
	return nil
}
