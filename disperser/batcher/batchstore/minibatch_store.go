package batchstore

import (
	"context"
	"fmt"
	"strconv"
	"time"

	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser/batcher"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

const (
	batchStatusIndexName      = "BatchStatusIndex"
	batchSKPrefix             = "BATCH#"
	minibatchSKPrefix         = "MINIBATCH#"
	dispersalRequestSKPrefix  = "DISPERSAL_REQUEST#"
	dispersalResponseSKPrefix = "DISPERSAL_RESPONSE#"
)

type MinibatchStore struct {
	dynamoDBClient *commondynamodb.Client
	tableName      string
	logger         logging.Logger
	ttl            time.Duration
}

var _ batcher.MinibatchStore = (*MinibatchStore)(nil)

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
				AttributeName: aws.String("BatchID"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("SK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("BatchStatus"),
				AttributeType: types.ScalarAttributeTypeN,
			},
			{
				AttributeName: aws.String("CreatedAt"),
				AttributeType: types.ScalarAttributeTypeN,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("BatchID"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("SK"),
				KeyType:       types.KeyTypeRange,
			},
		},
		TableName: aws.String(tableName),
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String(batchStatusIndexName),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("BatchStatus"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("CreatedAt"),
						KeyType:       types.KeyTypeRange,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(readCapacityUnits),
					WriteCapacityUnits: aws.Int64(writeCapacityUnits),
				},
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(readCapacityUnits),
			WriteCapacityUnits: aws.Int64(writeCapacityUnits),
		},
	}
}

func MarshalBatchRecord(batch *batcher.BatchRecord) (map[string]types.AttributeValue, error) {
	fields, err := attributevalue.MarshalMap(*batch)
	if err != nil {
		return nil, err
	}
	fields["BatchID"] = &types.AttributeValueMemberS{Value: batch.ID.String()}
	fields["SK"] = &types.AttributeValueMemberS{Value: batchSKPrefix + batch.ID.String()}
	fields["CreatedAt"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", batch.CreatedAt.UTC().Unix())}
	return fields, nil
}

func MarshalMinibatchRecord(minibatch *batcher.MinibatchRecord) (map[string]types.AttributeValue, error) {
	fields, err := attributevalue.MarshalMap(*minibatch)
	if err != nil {
		return nil, err
	}
	fields["BatchID"] = &types.AttributeValueMemberS{Value: minibatch.BatchID.String()}
	fields["SK"] = &types.AttributeValueMemberS{Value: minibatchSKPrefix + fmt.Sprintf("%d", minibatch.MinibatchIndex)}
	return fields, nil
}

func MarshalDispersalRequest(request *batcher.DispersalRequest) (map[string]types.AttributeValue, error) {
	fields, err := attributevalue.MarshalMap(*request)
	if err != nil {
		return nil, err
	}
	fields["BatchID"] = &types.AttributeValueMemberS{Value: request.BatchID.String()}
	fields["SK"] = &types.AttributeValueMemberS{Value: dispersalRequestSKPrefix + fmt.Sprintf("%d#%s", request.MinibatchIndex, request.OperatorID.Hex())}
	fields["OperatorID"] = &types.AttributeValueMemberS{Value: request.OperatorID.Hex()}
	fields["RequestedAt"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", request.RequestedAt.UTC().Unix())}
	return fields, nil
}

func MarshalDispersalResponse(response *batcher.DispersalResponse) (map[string]types.AttributeValue, error) {
	fields, err := attributevalue.MarshalMap(*response)
	if err != nil {
		return nil, err
	}
	fields["BatchID"] = &types.AttributeValueMemberS{Value: response.BatchID.String()}
	fields["SK"] = &types.AttributeValueMemberS{Value: dispersalResponseSKPrefix + fmt.Sprintf("%d#%s", response.MinibatchIndex, response.OperatorID.Hex())}
	fields["OperatorID"] = &types.AttributeValueMemberS{Value: response.OperatorID.Hex()}
	fields["RespondedAt"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", response.RespondedAt.UTC().Unix())}
	fields["RequestedAt"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", response.DispersalRequest.RequestedAt.UTC().Unix())}
	return fields, nil
}

func UnmarshalBatchID(item commondynamodb.Item) (*uuid.UUID, error) {
	type BatchID struct {
		BatchID string
	}

	batch := BatchID{}
	err := attributevalue.UnmarshalMap(item, &batch)
	if err != nil {
		return nil, err
	}

	batchID, err := uuid.Parse(batch.BatchID)
	if err != nil {
		return nil, err
	}

	return &batchID, nil
}

func UnmarshalOperatorID(item commondynamodb.Item) (*core.OperatorID, error) {
	type OperatorID struct {
		OperatorID string
	}

	dispersal := OperatorID{}
	err := attributevalue.UnmarshalMap(item, &dispersal)
	if err != nil {
		return nil, err
	}

	operatorID, err := core.OperatorIDFromHex(dispersal.OperatorID)
	if err != nil {
		return nil, err
	}

	return &operatorID, nil
}

func UnmarshalBatchRecord(item commondynamodb.Item) (*batcher.BatchRecord, error) {
	batch := batcher.BatchRecord{}
	err := attributevalue.UnmarshalMap(item, &batch)
	if err != nil {
		return nil, err
	}

	batchID, err := UnmarshalBatchID(item)
	if err != nil {
		return nil, err
	}
	batch.ID = *batchID

	batch.CreatedAt = batch.CreatedAt.UTC()
	return &batch, nil
}

func UnmarshalMinibatchRecord(item commondynamodb.Item) (*batcher.MinibatchRecord, error) {
	minibatch := batcher.MinibatchRecord{}
	err := attributevalue.UnmarshalMap(item, &minibatch)
	if err != nil {
		return nil, err
	}

	batchID, err := UnmarshalBatchID(item)
	if err != nil {
		return nil, err
	}
	minibatch.BatchID = *batchID

	return &minibatch, nil
}

func UnmarshalDispersalRequest(item commondynamodb.Item) (*batcher.DispersalRequest, error) {
	request := batcher.DispersalRequest{}
	err := attributevalue.UnmarshalMap(item, &request)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dispersal request from DynamoDB: %v", err)
	}

	batchID, err := UnmarshalBatchID(item)
	if err != nil {
		return nil, err
	}
	request.BatchID = *batchID

	operatorID, err := UnmarshalOperatorID(item)
	if err != nil {
		return nil, err
	}
	request.OperatorID = *operatorID

	request.RequestedAt = request.RequestedAt.UTC()
	return &request, nil
}

func UnmarshalDispersalResponse(item commondynamodb.Item) (*batcher.DispersalResponse, error) {
	response := batcher.DispersalResponse{}
	err := attributevalue.UnmarshalMap(item, &response)
	if err != nil {
		return nil, err
	}

	batchID, err := UnmarshalBatchID(item)
	if err != nil {
		return nil, err
	}
	response.BatchID = *batchID

	operatorID, err := UnmarshalOperatorID(item)
	if err != nil {
		return nil, err
	}
	response.OperatorID = *operatorID

	response.RespondedAt = response.RespondedAt.UTC()
	response.DispersalRequest.RequestedAt = response.DispersalRequest.RequestedAt.UTC()
	return &response, nil
}

func (m *MinibatchStore) PutBatch(ctx context.Context, batch *batcher.BatchRecord) error {
	item, err := MarshalBatchRecord(batch)
	if err != nil {
		return err
	}
	constraint := "attribute_not_exists(BatchID) AND attribute_not_exists(SK)"
	return m.dynamoDBClient.PutItemWithCondition(ctx, m.tableName, item, constraint)
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
		"BatchID": &types.AttributeValueMemberS{
			Value: batchID.String(),
		},
		"SK": &types.AttributeValueMemberS{
			Value: batchSKPrefix + batchID.String(),
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

func (m *MinibatchStore) BatchDispersed(ctx context.Context, batchID uuid.UUID) (bool, error) {
	dispersalRequests, err := m.GetDispersalRequests(ctx, batchID)
	if err != nil {
		return false, fmt.Errorf("failed to get dispersal requests for batch %s - %v", batchID.String(), err)

	}
	dispersalResponses, err := m.GetDispersalResponses(ctx, batchID)
	if err != nil {
		return false, fmt.Errorf("failed to get dispersal responses for batch %s - %v", batchID.String(), err)
	}
	if len(dispersalRequests) != len(dispersalResponses) {
		m.logger.Info("number of minibatch dispersal requests does not match responses", "batchID", batchID, "numRequests", len(dispersalRequests), "numResponses", len(dispersalResponses))
		return false, nil
	}
	if len(dispersalRequests) == 0 || len(dispersalResponses) == 0 {
		m.logger.Info("no dispersal requests or responses found", "batchID", batchID)
		return false, nil
	}
	return true, nil
}

func (m *MinibatchStore) GetBatchesByStatus(ctx context.Context, status batcher.BatchStatus) ([]*batcher.BatchRecord, error) {
	items, err := m.dynamoDBClient.QueryIndex(ctx, m.tableName, batchStatusIndexName, "BatchStatus = :status", commondynamodb.ExpresseionValues{
		":status": &types.AttributeValueMemberN{
			Value: strconv.Itoa(int(status)),
		}})
	if err != nil {
		return nil, err
	}

	batches := make([]*batcher.BatchRecord, len(items))
	for i, item := range items {
		batches[i], err = UnmarshalBatchRecord(item)
		if err != nil {
			m.logger.Errorf("failed to unmarshal batch record at index %d: %v", i, err)
			return nil, err
		}
	}

	return batches, nil
}

func (m *MinibatchStore) GetLatestFormedBatch(ctx context.Context) (batch *batcher.BatchRecord, minibatches []*batcher.MinibatchRecord, err error) {
	formed, err := m.GetBatchesByStatus(ctx, batcher.BatchStatusFormed)
	if err != nil {
		return nil, nil, err
	}
	if len(formed) == 0 {
		return nil, nil, nil
	}
	batch = formed[len(formed)-1]
	minibatches, err = m.GetMinibatches(ctx, batch.ID)
	if err != nil {
		return nil, nil, err
	}
	return batch, minibatches, nil
}

func (m *MinibatchStore) UpdateBatchStatus(ctx context.Context, batchID uuid.UUID, status batcher.BatchStatus) error {
	if status < batcher.BatchStatusFormed || status > batcher.BatchStatusFailed {
		return fmt.Errorf("invalid batch status %v", status)
	}
	_, err := m.dynamoDBClient.UpdateItem(ctx, m.tableName, map[string]types.AttributeValue{
		"BatchID": &types.AttributeValueMemberS{Value: batchID.String()},
		"SK":      &types.AttributeValueMemberS{Value: batchSKPrefix + batchID.String()},
	}, commondynamodb.Item{
		"BatchStatus": &types.AttributeValueMemberN{
			Value: strconv.Itoa(int(status)),
		},
	})

	if err != nil {
		return fmt.Errorf("failed to update batch status: %v", err)
	}

	return nil
}

func (m *MinibatchStore) GetMinibatch(ctx context.Context, batchID uuid.UUID, minibatchIndex uint) (*batcher.MinibatchRecord, error) {
	item, err := m.dynamoDBClient.GetItem(ctx, m.tableName, map[string]types.AttributeValue{
		"BatchID": &types.AttributeValueMemberS{
			Value: batchID.String(),
		},
		"SK": &types.AttributeValueMemberS{
			Value: minibatchSKPrefix + fmt.Sprintf("%d", minibatchIndex),
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

func (m *MinibatchStore) GetMinibatches(ctx context.Context, batchID uuid.UUID) ([]*batcher.MinibatchRecord, error) {
	items, err := m.dynamoDBClient.Query(ctx, m.tableName, "BatchID = :batchID AND begins_with(SK, :prefix)", commondynamodb.ExpresseionValues{
		":batchID": &types.AttributeValueMemberS{
			Value: batchID.String(),
		},
		":prefix": &types.AttributeValueMemberS{
			Value: minibatchSKPrefix,
		},
	})
	if err != nil {
		return nil, err
	}

	minibatches := make([]*batcher.MinibatchRecord, len(items))
	for i, item := range items {
		minibatches[i], err = UnmarshalMinibatchRecord(item)
		if err != nil {
			m.logger.Errorf("failed to unmarshal minibatch record at index %d: %v", i, err)
			return nil, err
		}
	}

	return minibatches, nil
}

func (m *MinibatchStore) GetDispersalRequest(ctx context.Context, batchID uuid.UUID, minibatchIndex uint, opID core.OperatorID) (*batcher.DispersalRequest, error) {
	item, err := m.dynamoDBClient.GetItem(ctx, m.tableName, map[string]types.AttributeValue{
		"BatchID": &types.AttributeValueMemberS{
			Value: batchID.String(),
		},
		"SK": &types.AttributeValueMemberS{
			Value: dispersalRequestSKPrefix + fmt.Sprintf("%d#%s", minibatchIndex, opID.Hex()),
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

func (m *MinibatchStore) GetDispersalRequests(ctx context.Context, batchID uuid.UUID) ([]*batcher.DispersalRequest, error) {
	items, err := m.dynamoDBClient.Query(ctx, m.tableName, "BatchID = :batchID AND begins_with(SK, :prefix)", commondynamodb.ExpresseionValues{
		":batchID": &types.AttributeValueMemberS{
			Value: batchID.String(),
		},
		":prefix": &types.AttributeValueMemberS{
			Value: dispersalRequestSKPrefix,
		},
	})
	if err != nil {
		return nil, err
	}

	requests := make([]*batcher.DispersalRequest, len(items))
	for i, item := range items {
		requests[i], err = UnmarshalDispersalRequest(item)
		if err != nil {
			m.logger.Errorf("failed to unmarshal dispersal requests at index %d: %v", i, err)
			return nil, err
		}
	}

	return requests, nil
}

func (m *MinibatchStore) GetMinibatchDispersalRequests(ctx context.Context, batchID uuid.UUID, minibatchIndex uint) ([]*batcher.DispersalRequest, error) {
	items, err := m.dynamoDBClient.Query(ctx, m.tableName, "BatchID = :batchID AND SK = :sk", commondynamodb.ExpresseionValues{
		":batchID": &types.AttributeValueMemberS{
			Value: batchID.String(),
		},
		":sk": &types.AttributeValueMemberS{
			Value: dispersalRequestSKPrefix + fmt.Sprintf("%s#%d", batchID.String(), minibatchIndex),
		},
	})
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("no dispersal requests found for BatchID %s MinibatchIndex %d", batchID, minibatchIndex)
	}

	requests := make([]*batcher.DispersalRequest, len(items))
	for i, item := range items {
		requests[i], err = UnmarshalDispersalRequest(item)
		if err != nil {
			m.logger.Errorf("failed to unmarshal dispersal requests at index %d: %v", i, err)
			return nil, err
		}
	}

	return requests, nil
}

func (m *MinibatchStore) GetDispersalResponse(ctx context.Context, batchID uuid.UUID, minibatchIndex uint, opID core.OperatorID) (*batcher.DispersalResponse, error) {
	item, err := m.dynamoDBClient.GetItem(ctx, m.tableName, map[string]types.AttributeValue{
		"BatchID": &types.AttributeValueMemberS{
			Value: batchID.String(),
		},
		"SK": &types.AttributeValueMemberS{
			Value: dispersalResponseSKPrefix + fmt.Sprintf("%d#%s", minibatchIndex, opID.Hex()),
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

func (m *MinibatchStore) GetDispersalResponses(ctx context.Context, batchID uuid.UUID) ([]*batcher.DispersalResponse, error) {
	items, err := m.dynamoDBClient.Query(ctx, m.tableName, "BatchID = :batchID AND begins_with(SK, :prefix)", commondynamodb.ExpresseionValues{
		":batchID": &types.AttributeValueMemberS{
			Value: batchID.String(),
		},
		":prefix": &types.AttributeValueMemberS{
			Value: dispersalResponseSKPrefix,
		},
	})
	if err != nil {
		return nil, err
	}

	responses := make([]*batcher.DispersalResponse, len(items))
	for i, item := range items {
		responses[i], err = UnmarshalDispersalResponse(item)
		if err != nil {
			m.logger.Errorf("failed to unmarshal dispersal response at index %d: %v", i, err)
			return nil, err
		}
	}

	return responses, nil
}

func (m *MinibatchStore) GetMinibatchDispersalResponses(ctx context.Context, batchID uuid.UUID, minibatchIndex uint) ([]*batcher.DispersalResponse, error) {
	items, err := m.dynamoDBClient.Query(ctx, m.tableName, "BatchID = :batchID AND SK = :sk", commondynamodb.ExpresseionValues{
		":batchID": &types.AttributeValueMemberS{
			Value: batchID.String(),
		},
		":sk": &types.AttributeValueMemberS{
			Value: dispersalResponseSKPrefix + fmt.Sprintf("%s#%d", batchID.String(), minibatchIndex),
		},
	})
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("no dispersal responses found for BatchID %s MinibatchIndex %d", batchID, minibatchIndex)
	}

	responses := make([]*batcher.DispersalResponse, len(items))
	for i, item := range items {
		responses[i], err = UnmarshalDispersalResponse(item)
		if err != nil {
			m.logger.Errorf("failed to unmarshal dispersal response at index %d: %v", i, err)
			return nil, err
		}
	}

	return responses, nil
}
