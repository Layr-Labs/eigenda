package batchstore

import (
	"context"
	"fmt"
	"strconv"
	"time"

	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/batcher"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

const (
	batchStatusIndexName          = "BatchStatusIndex"
	blobMinibatchMappingIndexName = "BlobMinibatchMappingIndex"
	batchSKPrefix                 = "BATCH#"
	dispersalSKPrefix             = "DISPERSAL#"
	blobMinibatchMappingSKPrefix  = "BLOB_MINIBATCH_MAPPING#"
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
			{
				AttributeName: aws.String("BlobHash"),
				AttributeType: types.ScalarAttributeTypeS,
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
			{
				IndexName: aws.String(blobMinibatchMappingIndexName),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("BlobHash"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("SK"),
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

func MarshalDispersal(response *batcher.MinibatchDispersal) (map[string]types.AttributeValue, error) {
	fields, err := attributevalue.MarshalMap(*response)
	if err != nil {
		return nil, err
	}
	fields["BatchID"] = &types.AttributeValueMemberS{Value: response.BatchID.String()}
	fields["SK"] = &types.AttributeValueMemberS{Value: dispersalSKPrefix + fmt.Sprintf("%d#%s", response.MinibatchIndex, response.OperatorID.Hex())}
	fields["OperatorID"] = &types.AttributeValueMemberS{Value: response.OperatorID.Hex()}
	fields["RespondedAt"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", response.RespondedAt.UTC().Unix())}
	fields["RequestedAt"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", response.RequestedAt.UTC().Unix())}
	return fields, nil
}

func MarshalDispersalResponse(response *batcher.DispersalResponse) (map[string]types.AttributeValue, error) {
	fields, err := attributevalue.MarshalMap(*response)
	if err != nil {
		return nil, err
	}
	return fields, nil
}

func MarshalBlobMinibatchMapping(blobMinibatchMapping *batcher.BlobMinibatchMapping) (map[string]types.AttributeValue, error) {
	fields, err := attributevalue.MarshalMap(*blobMinibatchMapping)
	if err != nil {
		return nil, err
	}
	fields["BatchID"] = &types.AttributeValueMemberS{Value: blobMinibatchMapping.BatchID.String()}
	fields["SK"] = &types.AttributeValueMemberS{Value: blobMinibatchMappingSKPrefix + fmt.Sprintf("%s#%s#%d", blobMinibatchMapping.BlobKey.MetadataHash, blobMinibatchMapping.BatchID, blobMinibatchMapping.BlobIndex)}
	fields["BlobHash"] = &types.AttributeValueMemberS{Value: blobMinibatchMapping.BlobKey.BlobHash}
	fields["MetadataHash"] = &types.AttributeValueMemberS{Value: blobMinibatchMapping.BlobKey.MetadataHash}
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

func UnmarshalBlobKey(item commondynamodb.Item) (*disperser.BlobKey, error) {
	blobKey := disperser.BlobKey{}
	err := attributevalue.UnmarshalMap(item, &blobKey)
	if err != nil {
		return nil, err
	}

	return &blobKey, nil
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

func UnmarshalDispersal(item commondynamodb.Item) (*batcher.MinibatchDispersal, error) {
	response := batcher.MinibatchDispersal{}
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
	response.RequestedAt = response.RequestedAt.UTC()
	return &response, nil
}

func UnmarshalBlobMinibatchMapping(item commondynamodb.Item) (*batcher.BlobMinibatchMapping, error) {
	blobMinibatchMapping := batcher.BlobMinibatchMapping{}
	err := attributevalue.UnmarshalMap(item, &blobMinibatchMapping)
	if err != nil {
		return nil, err
	}

	batchID, err := UnmarshalBatchID(item)
	if err != nil {
		return nil, err
	}
	blobMinibatchMapping.BatchID = *batchID

	blobKey, err := UnmarshalBlobKey(item)
	if err != nil {
		return nil, err
	}
	blobMinibatchMapping.BlobKey = blobKey

	return &blobMinibatchMapping, nil
}

func (m *MinibatchStore) PutBatch(ctx context.Context, batch *batcher.BatchRecord) error {
	item, err := MarshalBatchRecord(batch)
	if err != nil {
		return err
	}
	constraint := "attribute_not_exists(BatchID) AND attribute_not_exists(SK)"
	return m.dynamoDBClient.PutItemWithCondition(ctx, m.tableName, item, constraint)
}

func (m *MinibatchStore) PutDispersal(ctx context.Context, response *batcher.MinibatchDispersal) error {
	item, err := MarshalDispersal(response)
	if err != nil {
		return err
	}

	return m.dynamoDBClient.PutItem(ctx, m.tableName, item)
}

func (m *MinibatchStore) UpdateDispersalResponse(ctx context.Context, dispersal *batcher.MinibatchDispersal, response *batcher.DispersalResponse) error {
	marshaledResponse, err := MarshalDispersalResponse(response)
	if err != nil {
		return err
	}
	_, err = m.dynamoDBClient.UpdateItem(ctx, m.tableName, map[string]types.AttributeValue{
		"BatchID": &types.AttributeValueMemberS{
			Value: dispersal.BatchID.String(),
		},
		"SK": &types.AttributeValueMemberS{
			Value: dispersalSKPrefix + fmt.Sprintf("%d#%s", dispersal.MinibatchIndex, dispersal.OperatorID.Hex()),
		},
	}, marshaledResponse)

	return err
}

func (m *MinibatchStore) PutBlobMinibatchMappings(ctx context.Context, blobMinibatchMappings []*batcher.BlobMinibatchMapping) error {
	items := make([]map[string]types.AttributeValue, len(blobMinibatchMappings))
	var err error
	for i, blobMinibatchMapping := range blobMinibatchMappings {
		items[i], err = MarshalBlobMinibatchMapping(blobMinibatchMapping)
		if err != nil {
			return err
		}
	}

	failedItems, err := m.dynamoDBClient.PutItems(ctx, m.tableName, items)
	if err != nil {
		return err
	}
	if len(failedItems) > 0 {
		return fmt.Errorf("failed to put blob minibatch mappings: %v", failedItems)
	}

	return nil
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

func (m *MinibatchStore) BatchDispersed(ctx context.Context, batchID uuid.UUID, numMinibatches uint) (bool, error) {
	dispersals, err := m.GetDispersalsByBatchID(ctx, batchID)
	if err != nil {
		return false, fmt.Errorf("failed to get dispersal responses for batch %s - %v", batchID.String(), err)
	}
	if len(dispersals) == 0 {
		m.logger.Info("no dispersals found", "batchID", batchID)
		return false, nil
	}

	minibatchIndices := make(map[uint]struct{})
	for _, dispersal := range dispersals {
		minibatchIndices[dispersal.MinibatchIndex] = struct{}{}
		if dispersal.RespondedAt.IsZero() || dispersal.Error != nil {
			m.logger.Info("response pending", "batchID", batchID, "minibatchIndex", dispersal.MinibatchIndex, "operatorID", dispersal.OperatorID.Hex())
			return false, nil
		}
	}
	if len(minibatchIndices) != int(numMinibatches) {
		m.logger.Info("number of minibatches does not match", "batchID", batchID, "numMinibatches", numMinibatches, "minibatchIndices", len(minibatchIndices))
		return false, nil
	}
	for i := uint(0); i < numMinibatches; i++ {
		if _, ok := minibatchIndices[i]; !ok {
			m.logger.Info("number of minibatches does not match", "batchID", batchID, "minibatchIndex", i, "numMinibatches", numMinibatches)
			return false, nil
		}
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

func (m *MinibatchStore) MarkBatchFormed(ctx context.Context, batchID uuid.UUID, numMinibatches uint) error {
	_, err := m.dynamoDBClient.UpdateItem(ctx, m.tableName, map[string]types.AttributeValue{
		"BatchID": &types.AttributeValueMemberS{Value: batchID.String()},
		"SK":      &types.AttributeValueMemberS{Value: batchSKPrefix + batchID.String()},
	}, commondynamodb.Item{
		"NumMinibatches": &types.AttributeValueMemberN{
			Value: strconv.Itoa(int(numMinibatches)),
		},
		"BatchStatus": &types.AttributeValueMemberN{
			Value: strconv.Itoa(int(batcher.BatchStatusFormed)),
		},
	})
	return err
}

func (m *MinibatchStore) GetLatestFormedBatch(ctx context.Context) (batch *batcher.BatchRecord, err error) {
	formed, err := m.GetBatchesByStatus(ctx, batcher.BatchStatusFormed)
	if err != nil {
		return nil, err
	}
	if len(formed) == 0 {
		return nil, nil
	}
	return formed[len(formed)-1], nil
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

func (m *MinibatchStore) GetDispersal(ctx context.Context, batchID uuid.UUID, minibatchIndex uint, opID core.OperatorID) (*batcher.MinibatchDispersal, error) {
	item, err := m.dynamoDBClient.GetItem(ctx, m.tableName, map[string]types.AttributeValue{
		"BatchID": &types.AttributeValueMemberS{
			Value: batchID.String(),
		},
		"SK": &types.AttributeValueMemberS{
			Value: dispersalSKPrefix + fmt.Sprintf("%d#%s", minibatchIndex, opID.Hex()),
		},
	})
	if err != nil {
		m.logger.Errorf("failed to get dispersal response from DynamoDB: %v", err)
		return nil, err
	}

	if item == nil {
		return nil, nil
	}

	response, err := UnmarshalDispersal(item)
	if err != nil {
		m.logger.Errorf("failed to unmarshal dispersal response from DynamoDB: %v", err)
		return nil, err
	}
	return response, nil
}

func (m *MinibatchStore) GetDispersalsByBatchID(ctx context.Context, batchID uuid.UUID) ([]*batcher.MinibatchDispersal, error) {
	items, err := m.dynamoDBClient.Query(ctx, m.tableName, "BatchID = :batchID AND begins_with(SK, :prefix)", commondynamodb.ExpresseionValues{
		":batchID": &types.AttributeValueMemberS{
			Value: batchID.String(),
		},
		":prefix": &types.AttributeValueMemberS{
			Value: dispersalSKPrefix,
		},
	})
	if err != nil {
		return nil, err
	}

	responses := make([]*batcher.MinibatchDispersal, len(items))
	for i, item := range items {
		responses[i], err = UnmarshalDispersal(item)
		if err != nil {
			m.logger.Errorf("failed to unmarshal dispersal response at index %d: %v", i, err)
			return nil, err
		}
	}

	return responses, nil
}

func (m *MinibatchStore) GetDispersalsByMinibatch(ctx context.Context, batchID uuid.UUID, minibatchIndex uint) ([]*batcher.MinibatchDispersal, error) {
	items, err := m.dynamoDBClient.Query(ctx, m.tableName, "BatchID = :batchID AND SK = :sk", commondynamodb.ExpresseionValues{
		":batchID": &types.AttributeValueMemberS{
			Value: batchID.String(),
		},
		":sk": &types.AttributeValueMemberS{
			Value: dispersalSKPrefix + fmt.Sprintf("%s#%d", batchID.String(), minibatchIndex),
		},
	})
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("no dispersal responses found for BatchID %s MinibatchIndex %d", batchID, minibatchIndex)
	}

	responses := make([]*batcher.MinibatchDispersal, len(items))
	for i, item := range items {
		responses[i], err = UnmarshalDispersal(item)
		if err != nil {
			m.logger.Errorf("failed to unmarshal dispersal response at index %d: %v", i, err)
			return nil, err
		}
	}

	return responses, nil
}

func (m *MinibatchStore) GetBlobMinibatchMappings(ctx context.Context, blobKey disperser.BlobKey) ([]*batcher.BlobMinibatchMapping, error) {
	items, err := m.dynamoDBClient.QueryIndex(ctx, m.tableName, blobMinibatchMappingIndexName, "BlobHash = :blobHash AND begins_with(SK, :prefix)", commondynamodb.ExpresseionValues{
		":blobHash": &types.AttributeValueMemberS{
			Value: blobKey.BlobHash,
		},
		":prefix": &types.AttributeValueMemberS{
			Value: blobMinibatchMappingSKPrefix + blobKey.MetadataHash,
		},
	})
	if err != nil {
		return nil, err
	}

	blobMinibatchMappings := make([]*batcher.BlobMinibatchMapping, len(items))
	for i, item := range items {
		blobMinibatchMappings[i], err = UnmarshalBlobMinibatchMapping(item)
		if err != nil {
			m.logger.Errorf("failed to unmarshal blob minibatch mapping at index %d: %v", i, err)
			return nil, err
		}
	}

	return blobMinibatchMappings, nil
}

func (m *MinibatchStore) GetBlobMinibatchMappingsByBatchID(ctx context.Context, batchID uuid.UUID) ([]*batcher.BlobMinibatchMapping, error) {
	items, err := m.dynamoDBClient.Query(ctx, m.tableName, "BatchID = :batchID", commondynamodb.ExpresseionValues{
		":batchID": &types.AttributeValueMemberS{
			Value: batchID.String(),
		},
	})
	if err != nil {
		return nil, err
	}

	blobMinibatchMappings := make([]*batcher.BlobMinibatchMapping, len(items))
	for i, item := range items {
		blobMinibatchMappings[i], err = UnmarshalBlobMinibatchMapping(item)
		if err != nil {
			m.logger.Errorf("failed to unmarshal blob minibatch mapping at index %d: %v", i, err)
			return nil, err
		}
	}

	return blobMinibatchMappings, nil
}
