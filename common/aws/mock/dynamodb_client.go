package mock

import (
	"context"

	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/mock"
)

type MockDynamoDBClient struct {
	mock.Mock
}

var _ dynamodb.Client = (*MockDynamoDBClient)(nil)

func (c *MockDynamoDBClient) DeleteTable(ctx context.Context, tableName string) error {
	args := c.Called()
	return args.Error(0)
}

func (c *MockDynamoDBClient) PutItem(ctx context.Context, tableName string, item dynamodb.Item) error {
	args := c.Called()
	return args.Error(0)
}

func (c *MockDynamoDBClient) PutItemWithCondition(ctx context.Context, tableName string, item dynamodb.Item, condition string, expressionAttributeNames map[string]string, expressionAttributeValues map[string]types.AttributeValue) error {
	args := c.Called()
	return args.Error(0)
}

func (c *MockDynamoDBClient) PutItems(ctx context.Context, tableName string, items []dynamodb.Item) ([]dynamodb.Item, error) {
	args := c.Called(ctx, tableName, items)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dynamodb.Item), args.Error(1)
}

func (c *MockDynamoDBClient) UpdateItem(ctx context.Context, tableName string, key dynamodb.Key, item dynamodb.Item) (dynamodb.Item, error) {
	args := c.Called()
	return args.Get(0).(dynamodb.Item), args.Error(1)
}

func (c *MockDynamoDBClient) UpdateItemWithCondition(ctx context.Context, tableName string, key dynamodb.Key, item dynamodb.Item, condition expression.ConditionBuilder) (dynamodb.Item, error) {
	args := c.Called()
	return args.Get(0).(dynamodb.Item), args.Error(1)
}

func (c *MockDynamoDBClient) IncrementBy(ctx context.Context, tableName string, key dynamodb.Key, attr string, value uint64) (dynamodb.Item, error) {
	args := c.Called()
	return args.Get(0).(dynamodb.Item), args.Error(1)
}

func (c *MockDynamoDBClient) GetItem(ctx context.Context, tableName string, key dynamodb.Key) (dynamodb.Item, error) {
	args := c.Called()
	return args.Get(0).(dynamodb.Item), args.Error(1)
}

func (c *MockDynamoDBClient) GetItems(ctx context.Context, tableName string, keys []dynamodb.Key) ([]dynamodb.Item, error) {
	args := c.Called()
	return args.Get(0).([]dynamodb.Item), args.Error(1)
}

func (c *MockDynamoDBClient) QueryIndex(ctx context.Context, tableName string, indexName string, keyCondition string, expAttributeValues dynamodb.ExpressionValues) ([]dynamodb.Item, error) {
	args := c.Called()
	return args.Get(0).([]dynamodb.Item), args.Error(1)
}

func (c *MockDynamoDBClient) Query(ctx context.Context, tableName string, keyCondition string, expAttributeValues dynamodb.ExpressionValues) ([]dynamodb.Item, error) {
	args := c.Called()
	return args.Get(0).([]dynamodb.Item), args.Error(1)
}

func (c *MockDynamoDBClient) QueryWithInput(ctx context.Context, input *awsdynamodb.QueryInput) ([]dynamodb.Item, error) {
	args := c.Called()
	return args.Get(0).([]dynamodb.Item), args.Error(1)
}

func (c *MockDynamoDBClient) QueryIndexCount(ctx context.Context, tableName string, indexName string, keyCondition string, expAttributeValues dynamodb.ExpressionValues) (int32, error) {
	args := c.Called()
	return args.Get(0).(int32), args.Error(1)
}

func (c *MockDynamoDBClient) QueryIndexWithPagination(ctx context.Context, tableName string, indexName string, keyCondition string, expAttributeValues dynamodb.ExpressionValues, limit int32, exclusiveStartKey map[string]types.AttributeValue) (dynamodb.QueryResult, error) {
	args := c.Called()
	return args.Get(0).(dynamodb.QueryResult), args.Error(1)
}

func (c *MockDynamoDBClient) DeleteItem(ctx context.Context, tableName string, key dynamodb.Key) error {
	args := c.Called()
	return args.Error(0)
}

func (c *MockDynamoDBClient) DeleteItems(ctx context.Context, tableName string, keys []dynamodb.Key) ([]dynamodb.Key, error) {
	args := c.Called()
	return args.Get(0).([]dynamodb.Key), args.Error(1)
}

func (c *MockDynamoDBClient) TableExists(ctx context.Context, name string) error {
	args := c.Called()
	return args.Error(0)
}
