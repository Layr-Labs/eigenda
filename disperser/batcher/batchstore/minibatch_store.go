package batchstore

import (
	"context"
	"fmt"
	"time"

	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/disperser/batcher"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

const (
	batchKey             = "BATCH#"
	minibatchKey         = "MINIBATCH#"
	dispersalRequestKey  = "DISPERSAL_REQUEST#"
	dispersalResponseKey = "DISPERSAL_RESPONSE#"
)

type MinibatchStore struct {
	dynamoDBClient *commondynamodb.Client
	tableName      string
	logger         logging.Logger
	ttl            time.Duration
}

func NewMinibatchStore(dynamoDBClient *commondynamodb.Client, logger logging.Logger, tableName string, ttl time.Duration) *MinibatchStore {
	logger.Debugf("creating minibatch store with table %s with TTL: %s", tableName, ttl)
	return &MinibatchStore{
		dynamoDBClient: dynamoDBClient,
		logger:         logger.With("component", "MinibatchStore"),
		tableName:      tableName,
		ttl:            ttl,
	}
}

func GenerateTableSchema(tableName string, readCapacityUnits int64, writeCapacityUnits int64) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("PK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("SK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("PK"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("SK"),
				KeyType:       types.KeyTypeRange,
			},
		},
		TableName:              aws.String(tableName),
		GlobalSecondaryIndexes: nil,
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(readCapacityUnits),
			WriteCapacityUnits: aws.Int64(writeCapacityUnits),
		},
	}
}

func MarshalBatchRecord(batch *batcher.BatchRecord) (map[string]types.AttributeValue, error) {
	fields, err := attributevalue.MarshalMap(batch)
	if err != nil {
		return nil, err
	}
	fields["PK"] = &types.AttributeValueMemberS{Value: batchKey + batch.ID.String()}
	fields["SK"] = &types.AttributeValueMemberS{Value: batchKey + batch.ID.String()}
	fields["CreatedAt"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", batch.CreatedAt.UTC().Unix())}
	return fields, nil
}

func MarshalMinibatchRecord(minibatch *batcher.MinibatchRecord) (map[string]types.AttributeValue, error) {
	fields, err := attributevalue.MarshalMap(minibatch)
	if err != nil {
		return nil, err
	}
	fields["PK"] = &types.AttributeValueMemberS{Value: batchKey + minibatch.BatchID.String()}
	fields["SK"] = &types.AttributeValueMemberS{Value: minibatchKey + fmt.Sprintf("%d", minibatch.MinibatchIndex)}
	return fields, nil
}

func MarshalDispersalRequest(request *batcher.DispersalRequest) (map[string]types.AttributeValue, error) {
	fields, err := attributevalue.MarshalMap(request)
	if err != nil {
		return nil, err
	}
	fields["PK"] = &types.AttributeValueMemberS{Value: batchKey + request.BatchID.String()}
	fields["SK"] = &types.AttributeValueMemberS{Value: dispersalRequestKey + fmt.Sprintf("%d", request.MinibatchIndex)}
	fields["RequestedAt"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", request.RequestedAt.UTC().Unix())}
	return fields, nil
}

func MarshalDispersalResponse(response *batcher.DispersalResponse) (map[string]types.AttributeValue, error) {
	fields, err := attributevalue.MarshalMap(response)
	if err != nil {
		return nil, err
	}
	fields["PK"] = &types.AttributeValueMemberS{Value: batchKey + response.BatchID.String()}
	fields["SK"] = &types.AttributeValueMemberS{Value: dispersalResponseKey + fmt.Sprintf("%d", response.MinibatchIndex)}
	fields["RespondedAt"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", response.RequestedAt.UTC().Unix())}
	return fields, nil
}
func UnmarshalBatchRecord(item commondynamodb.Item) (*batcher.BatchRecord, error) {
	batch := batcher.BatchRecord{}
	err := attributevalue.UnmarshalMap(item, &batch)
	if err != nil {
		return nil, err
	}
	batch.CreatedAt = batch.CreatedAt.UTC()
	return &batch, nil
}

func UnmarshalMinibatchRecord(item commondynamodb.Item) (*batcher.MinibatchRecord, error) {
	minibatch := batcher.MinibatchRecord{}
	err := attributevalue.UnmarshalMap(item, &minibatch)
	if err != nil {
		return nil, err
	}
	return &minibatch, nil
}

func UnmarshalDispersalRequest(item commondynamodb.Item) (*batcher.DispersalRequest, error) {
	request := batcher.DispersalRequest{}
	err := attributevalue.UnmarshalMap(item, &request)
	if err != nil {
		return nil, err
	}
	request.RequestedAt = request.RequestedAt.UTC()
	return &request, nil
}

func UnmarshalDispersalResponse(item commondynamodb.Item) (*batcher.DispersalResponse, error) {
	response := batcher.DispersalResponse{}
	err := attributevalue.UnmarshalMap(item, &response)
	if err != nil {
		return nil, err
	}
	response.RespondedAt = response.RespondedAt.UTC()
	return &response, nil
}

func (m *MinibatchStore) PutBatch(ctx context.Context, batch *batcher.BatchRecord) error {
	item, err := MarshalBatchRecord(batch)
	if err != nil {
		return err
	}

	return m.dynamoDBClient.PutItem(ctx, m.tableName, item)
}

func (m *MinibatchStore) PutMinibatch(ctx context.Context, minibatch *batcher.MinibatchRecord) error {
	item, err := MarshalMinibatchRecord(minibatch)
	if err != nil {
		return err
	}

	return m.dynamoDBClient.PutItem(ctx, m.tableName, item)
}

func (m *MinibatchStore) PutDispersalRequest(ctx context.Context, request *batcher.DispersalRequest) error {
	item, err := MarshalDispersalRequest(request)
	if err != nil {
		return err
	}

	fmt.Printf("%v", item)
	return m.dynamoDBClient.PutItem(ctx, m.tableName, item)
}

func (m *MinibatchStore) PutDispersalResponse(ctx context.Context, response *batcher.DispersalResponse) error {
	item, err := MarshalDispersalResponse(response)
	if err != nil {
		return err
	}

	return m.dynamoDBClient.PutItem(ctx, m.tableName, item)
}

func (m *MinibatchStore) GetBatch(ctx context.Context, batchID uuid.UUID) (*batcher.BatchRecord, error) {
	item, err := m.dynamoDBClient.GetItem(ctx, m.tableName, map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{
			Value: batchKey + batchID.String(),
		},
		"SK": &types.AttributeValueMemberS{
			Value: batchKey + batchID.String(),
		},
	})
	if err != nil {
		m.logger.Errorf("failed to get batch from DynamoDB: %v", err)
		return nil, err
	}

	if item == nil {
		return nil, nil
	}

	batch, err := UnmarshalBatchRecord(item)
	if err != nil {
		m.logger.Errorf("failed to unmarshal batch record from DynamoDB: %v", err)
		return nil, err
	}
	return batch, nil
}

func (m *MinibatchStore) GetMinibatch(ctx context.Context, batchID uuid.UUID, minibatchIndex uint) (*batcher.MinibatchRecord, error) {
	item, err := m.dynamoDBClient.GetItem(ctx, m.tableName, map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{
			Value: batchKey + batchID.String(),
		},
		"SK": &types.AttributeValueMemberS{
			Value: minibatchKey + fmt.Sprintf("%d", minibatchIndex),
		},
	})
	if err != nil {
		m.logger.Errorf("failed to get minibatch from DynamoDB: %v", err)
		return nil, err
	}

	if item == nil {
		return nil, nil
	}

	minibatch, err := UnmarshalMinibatchRecord(item)
	if err != nil {
		m.logger.Errorf("failed to unmarshal minibatch record from DynamoDB: %v", err)
		return nil, err
	}
	return minibatch, nil
}

func (m *MinibatchStore) GetDispersalRequest(ctx context.Context, batchID uuid.UUID, minibatchIndex uint) (*batcher.DispersalRequest, error) {
	item, err := m.dynamoDBClient.GetItem(ctx, m.tableName, map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{
			Value: batchKey + batchID.String(),
		},
		"SK": &types.AttributeValueMemberS{
			Value: dispersalRequestKey + fmt.Sprintf("%d", minibatchIndex),
		},
	})
	if err != nil {
		m.logger.Errorf("failed to get dispersal request from DynamoDB: %v", err)
		return nil, err
	}

	if item == nil {
		return nil, nil
	}

	request, err := UnmarshalDispersalRequest(item)
	if err != nil {
		m.logger.Errorf("failed to unmarshal dispersal request from DynamoDB: %v", err)
		return nil, err
	}
	return request, nil
}

func (m *MinibatchStore) GetDispersalResponse(ctx context.Context, batchID uuid.UUID, minibatchIndex uint) (*batcher.DispersalResponse, error) {
	item, err := m.dynamoDBClient.GetItem(ctx, m.tableName, map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{
			Value: batchKey + batchID.String(),
		},
		"SK": &types.AttributeValueMemberS{
			Value: dispersalResponseKey + fmt.Sprintf("%d", minibatchIndex),
		},
	})
	if err != nil {
		m.logger.Errorf("failed to get dispersal response from DynamoDB: %v", err)
		return nil, err
	}

	if item == nil {
		return nil, nil
	}

	response, err := UnmarshalDispersalResponse(item)
	if err != nil {
		m.logger.Errorf("failed to unmarshal dispersal response from DynamoDB: %v", err)
		return nil, err
	}
	return response, nil
}
