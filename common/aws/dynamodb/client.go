package dynamodb

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"sync"

	"github.com/Layr-Labs/eigenda/common"
	commonaws "github.com/Layr-Labs/eigenda/common/aws"
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

type Client struct {
	dynamoClient *dynamodb.Client
	logger       common.Logger
}

func NewClient(cfg commonaws.ClientConfig, logger common.Logger) (*Client, error) {
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
		clientRef = &Client{dynamoClient: dynamoClient, logger: logger}
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

// Add a new item to the table if key does not exist
// Add version attribute to item for optimistic locking
func (c *Client) PutItemWithVersion(ctx context.Context, tableName string, item Item) error {
	// Check if the version key exists in the item map, if not set it to 1
	// Set version attribute to 1 if not present
	if _, exists := item["Version"]; !exists {
		item["Version"] = &types.AttributeValueMemberN{Value: "1"}
	}

	_, err := c.dynamoClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
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

func (c *Client) UpsertItemWithExpression(ctx context.Context, tableName string, key Key, item Item, customUpdateExpr *expression.UpdateBuilder) (Item, error) {
	var resp *dynamodb.UpdateItemOutput
	var err error
	// Retry the operation twice if the conditional check fails
	for attempt := 0; attempt < 2; attempt++ {
		currentItem, expectedVersion, err := c.GetItemWithVersion(ctx, tableName, key)
		if err != nil {
			return nil, err
		}

		if currentItem == nil {
			item["Version"] = &types.AttributeValueMemberN{Value: "1"} // Set initial version to 1
			err = c.PutItemWithVersion(ctx, tableName, item)
			if err != nil {
				return nil, err
			}
			return item, nil
		}

		// Prepare the update builder
		var update expression.UpdateBuilder
		if customUpdateExpr != nil {
			// Use the provided custom update expression
			update = *customUpdateExpr
		} else {
			// Build the update expression for other item attributes
			for itemKey, itemValue := range item {
				if _, ok := key[itemKey]; ok {
					continue // Cannot update the key
				}
				if itemKey != "Version" { // Ensure not to update "Version" twice
					update = update.Set(expression.Name(itemKey), expression.Value(itemValue))
				}
			}
		}

		// Update Version attribute
		versionName := expression.Name("Version")
		update = update.Set(versionName, versionName.Plus(expression.Value(1)))

		// Conditional expression to check the version
		condition := expression.Name("Version").Equal(expression.Value(expectedVersion))

		// Build the expression
		expr, err := expression.NewBuilder().WithUpdate(update).WithCondition(condition).Build()
		if err != nil {
			return nil, err
		}

		// Perform the update
		resp, err = c.dynamoClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName:                 aws.String(tableName),
			Key:                       key,
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
			ConditionExpression:       expr.Condition(),
			ReturnValues:              types.ReturnValueUpdatedNew,
		})

		if err != nil {
			continue // If conditional check fails, retry
		}

		break // If successful, break the loop
	}

	if err != nil {
		return nil, err
	}

	return resp.Attributes, nil
}

// UpdateItemWithVersion updates an item with optimistic locking
func (c *Client) UpsertItemWithVersion(ctx context.Context, tableName string, key Key, item Item, expectedVersion int) (Item, error) {
	var resp *dynamodb.UpdateItemOutput
	var err error

	// Retry the operation twice if the conditional check fails
	for attempt := 0; attempt < 2; attempt++ {
		// Get the current version of the item
		fmt.Println("Attempt", attempt)
		fmt.Println("Key", key)
		currentItem, currentVersion, err := c.GetItemWithVersion(ctx, tableName, key)
		if err != nil {
			return nil, err
		}
		fmt.Println("Expected Version Of Item", expectedVersion)

		// If item not found, put a new item with the initial version
		if currentItem == nil {
			fmt.Println("Adding Item with key to store", key)
			item["Version"] = &types.AttributeValueMemberN{Value: "1"} // Set initial version to 1
			err = c.PutItemWithVersion(ctx, tableName, item)
			if err != nil {
				return nil, err
			}
			return item, nil
		}

		// Check CurrentVersion Matches Expected Version
		if expectedVersion != currentVersion {
			fmt.Println("Version mismatch, for key", key)
			// Version mismatch, may retry or handle as needed
			continue
		}

		// Prepare the update builder
		update := expression.UpdateBuilder{}

		// Build the update expression for other item attributes
		for itemKey, itemValue := range item {
			if _, ok := key[itemKey]; ok {
				// Cannot update the key
				continue
			}
			if itemKey != "Version" { // Ensure not to update "Version" twice
				update = update.Set(expression.Name(itemKey), expression.Value(itemValue))
			}
		}

		// Update Version attribute
		versionName := expression.Name("Version")

		// Increment the version number
		update = update.Set(versionName, versionName.Plus(expression.Value(1)))

		// Conditional expression to check the version
		condition := expression.Name("Version").Equal(expression.Value(expectedVersion))

		// Build the expression
		expr, err := expression.NewBuilder().WithUpdate(update).WithCondition(condition).Build()
		if err != nil {
			return nil, err
		}

		// Perform the update
		fmt.Println("Begin Update of Key", key)

		resp, err = c.dynamoClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName:                 aws.String(tableName),
			Key:                       key,
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
			ConditionExpression:       expr.Condition(),
			ReturnValues:              types.ReturnValueUpdatedNew,
		})

		if err != nil {
			fmt.Println("Conditional check failed for key %v, error:", key, err)
			continue
		}

		// If successful, break the loop
		break
	}

	if err != nil {
		return nil, err
	}

	return resp.Attributes, nil
}

// Updates an item with optimistic locking
func (c *Client) UpdateItemWithVersion(ctx context.Context, tableName string, key Key, item Item, expectedVersion int) (Item, error) {
	update := expression.UpdateBuilder{}
	versionName := expression.Name("Version")

	// Increment the version number
	update = update.Set(versionName, versionName.Plus(expression.Value(1)))

	// Build the update expression for other item attributes
	for itemKey, itemValue := range item {
		if _, ok := key[itemKey]; ok {
			// Cannot update the key
			continue
		}
		update = update.Set(expression.Name(itemKey), expression.Value(itemValue))
	}

	// Build the conditional expression to check the version
	condition := expression.Name("Version").Equal(expression.Value(expectedVersion))

	expr, err := expression.NewBuilder().WithUpdate(update).WithCondition(condition).Build()
	if err != nil {
		return nil, err
	}

	resp, err := c.dynamoClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       key,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		ConditionExpression:       expr.Condition(),
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

// GetItemWithVersion returns Item and Version of fetched Item
func (c *Client) GetItemWithVersion(ctx context.Context, tableName string, key Key) (Item, int, error) {
	resp, err := c.dynamoClient.GetItem(ctx, &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, 0, err
	}

	// Extract the version from the item
	version := 0 // Default to 0 if no version is present
	if v, ok := resp.Item["Version"]; ok {
		versionValue, ok := v.(*types.AttributeValueMemberN)
		if !ok {
			return nil, 0, fmt.Errorf("version attribute is not a number")
		}
		version, err = strconv.Atoi(versionValue.Value)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to parse version: %v", err)
		}
	}

	return resp.Item, version, nil
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
