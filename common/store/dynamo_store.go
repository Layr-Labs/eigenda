package store

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/common"
	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type dynamodbBucketStore[T any] struct {
	client    *commondynamodb.Client
	tableName string
}

func NewDynamoParamStore[T any](client *commondynamodb.Client, tableName string) common.KVStoreVersioned[T] {
	return &dynamodbBucketStore[T]{
		client:    client,
		tableName: tableName,
	}
}

func (s *dynamodbBucketStore[T]) GetItem(ctx context.Context, requesterID string) (*T, error) {

	key := map[string]types.AttributeValue{
		"RequesterID": &types.AttributeValueMemberS{
			Value: requesterID,
		},
	}

	item, err := s.client.GetItem(ctx, s.tableName, key)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("item not found")
	}

	params := new(T)
	err = attributevalue.UnmarshalMap(item, params)
	if err != nil {
		return nil, err
	}

	return params, nil
}

func (s *dynamodbBucketStore[T]) GetItemWithVersion(ctx context.Context, requesterID string) (*T, int, error) {

	key := map[string]types.AttributeValue{
		"RequesterID": &types.AttributeValueMemberS{
			Value: requesterID,
		},
	}

	item, version, err := s.client.GetItemWithVersion(ctx, s.tableName, key)
	if err != nil {
		return nil, 0, err
	}
	if item == nil {
		return nil, 0, fmt.Errorf("item not found")
	}

	params := new(T)
	err = attributevalue.UnmarshalMap(item, params)
	if err != nil {
		return nil, 0, err
	}

	return params, version, nil
}

func (s *dynamodbBucketStore[T]) UpdateItem(ctx context.Context, requesterID string, params *T) error {

	fields, err := attributevalue.MarshalMap(params)
	if err != nil {
		return err
	}

	fields["RequesterID"] = &types.AttributeValueMemberS{
		Value: requesterID,
	}

	return s.client.PutItem(ctx, s.tableName, fields)
}

func (s *dynamodbBucketStore[T]) UpdateItemWithVersion(ctx context.Context, requesterID string, params *T, expectedVersion int) error {

	fields, err := attributevalue.MarshalMap(params)
	if err != nil {
		return err
	}

	fields["RequesterID"] = &types.AttributeValueMemberS{
		Value: requesterID,
	}

	key := map[string]types.AttributeValue{
		"RequesterID": &types.AttributeValueMemberS{
			Value: requesterID,
		},
	}

	_, err = s.client.UpsertItemWithVersion(ctx, s.tableName, key, fields, expectedVersion)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *dynamodbBucketStore[T]) UpdateItemWithExpression(ctx context.Context, requesterID string, customUpdateExpr *expression.UpdateBuilder) error {

	key := map[string]types.AttributeValue{
		"RequesterID": &types.AttributeValueMemberS{
			Value: requesterID,
		},
	}
	_, err := s.client.UpdateWithExpression(ctx, s.tableName, key, customUpdateExpr)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func GenerateTableSchema(readCapacityUnits int64, writeCapacityUnits int64, tableName string) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("RequesterID"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("RequesterID"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName: aws.String(tableName),
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(readCapacityUnits),
			WriteCapacityUnits: aws.Int64(writeCapacityUnits),
		},
	}
}
