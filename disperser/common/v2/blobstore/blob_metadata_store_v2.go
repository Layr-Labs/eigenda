package blobstore

import (
	"context"
	"fmt"
	"strconv"
	"time"

	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/disperser"
	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	StatusIndexName            = "StatusIndex"
	OperatorDispersalIndexName = "OperatorDispersalIndex"
	OperatorResponseIndexName  = "OperatorResponseIndex"
)

// BlobMetadataStore is a blob metadata storage backed by DynamoDB
type BlobMetadataStore struct {
	dynamoDBClient *commondynamodb.Client
	logger         logging.Logger
	tableName      string
	ttl            time.Duration
}

func NewBlobMetadataStore(dynamoDBClient *commondynamodb.Client, logger logging.Logger, tableName string, ttl time.Duration) *BlobMetadataStore {
	logger.Debugf("creating blob metadata store v2 with table %s with TTL: %s", tableName, ttl)
	return &BlobMetadataStore{
		dynamoDBClient: dynamoDBClient,
		logger:         logger.With("component", "BlobMetadataStoreV2"),
		tableName:      tableName,
		ttl:            ttl,
	}
}

func (s *BlobMetadataStore) CreateBlobMetadata(ctx context.Context, blobMetadata *v2.BlobMetadata) error {
	item, err := MarshalBlobMetadata(blobMetadata)
	if err != nil {
		return err
	}

	return s.dynamoDBClient.PutItem(ctx, s.tableName, item)
}

func (s *BlobMetadataStore) GetBlobMetadata(ctx context.Context, blobKey disperser.BlobKey) (*v2.BlobMetadata, error) {
	item, err := s.dynamoDBClient.GetItem(ctx, s.tableName, map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{
			Value: blobKey.BlobHash,
		},
		"SK": &types.AttributeValueMemberS{
			Value: blobKey.MetadataHash,
		},
	})

	if item == nil {
		return nil, fmt.Errorf("%w: metadata not found for key %s", disperser.ErrMetadataNotFound, blobKey)
	}

	if err != nil {
		return nil, err
	}

	metadata, err := UnmarshalBlobMetadata(item)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

// GetBulkBlobMetadata returns the metadata for the given blob keys
// Note: ordering of items is not guaranteed
func (s *BlobMetadataStore) GetBulkBlobMetadata(ctx context.Context, blobKeys []disperser.BlobKey) ([]*v2.BlobMetadata, error) {
	keys := make([]map[string]types.AttributeValue, len(blobKeys))
	for i := 0; i < len(blobKeys); i += 1 {
		keys[i] = map[string]types.AttributeValue{
			"BlobHash":     &types.AttributeValueMemberS{Value: blobKeys[i].BlobHash},
			"MetadataHash": &types.AttributeValueMemberS{Value: blobKeys[i].MetadataHash},
		}
	}
	items, err := s.dynamoDBClient.GetItems(ctx, s.tableName, keys)
	if err != nil {
		return nil, err
	}

	metadata := make([]*v2.BlobMetadata, len(items))
	for i, item := range items {
		metadata[i], err = UnmarshalBlobMetadata(item)
		if err != nil {
			return nil, err
		}
	}

	return metadata, nil
}

// GetBlobMetadataByStatus returns all the metadata with the given status
// Because this function scans the entire index, it should only be used for status with a limited number of items.
// It should only be used to filter "Processing" status. To support other status, a streaming version should be implemented.
func (s *BlobMetadataStore) GetBlobMetadataByStatus(ctx context.Context, status disperser.BlobStatus) ([]*v2.BlobMetadata, error) {
	items, err := s.dynamoDBClient.QueryIndex(ctx, s.tableName, expiryIndexName, "BlobStatus = :status AND Expiry > :expiry", commondynamodb.ExpresseionValues{
		":status": &types.AttributeValueMemberN{
			Value: strconv.Itoa(int(status)),
		},
		":expiry": &types.AttributeValueMemberN{
			Value: strconv.FormatInt(time.Now().Unix(), 10),
		}})
	if err != nil {
		return nil, err
	}

	metadata := make([]*v2.BlobMetadata, len(items))
	for i, item := range items {
		metadata[i], err = UnmarshalBlobMetadata(item)
		if err != nil {
			return nil, err
		}
	}

	return metadata, nil
}

// GetBlobMetadataCountByStatus returns the count of all the metadata with the given status
// Because this function scans the entire index, it should only be used for status with a limited number of items.
// It should only be used to filter "Processing" status. To support other status, a streaming version should be implemented.
func (s *BlobMetadataStore) GetBlobMetadataCountByStatus(ctx context.Context, status disperser.BlobStatus) (int32, error) {
	count, err := s.dynamoDBClient.QueryIndexCount(ctx, s.tableName, expiryIndexName, "BlobStatus = :status AND Expiry > :expiry", commondynamodb.ExpresseionValues{
		":status": &types.AttributeValueMemberN{
			Value: strconv.Itoa(int(status)),
		},
		":expiry": &types.AttributeValueMemberN{
			Value: strconv.FormatInt(time.Now().Unix(), 10),
		},
	})
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetBlobMetadataByStatusWithPagination returns all the metadata with the given status upto the specified limit
// along with items, also returns a pagination token that can be used to fetch the next set of items
//
// Note that this may not return all the metadata for the batch if dynamodb query limit is reached.
// e.g 1mb limit for a single query
func (s *BlobMetadataStore) GetBlobMetadataByStatusWithPagination(ctx context.Context, status disperser.BlobStatus, limit int32, exclusiveStartKey *disperser.BlobStoreExclusiveStartKey) ([]*v2.BlobMetadata, *disperser.BlobStoreExclusiveStartKey, error) {

	var attributeMap map[string]types.AttributeValue
	var err error

	// Convert the exclusive start key to a map of AttributeValue
	if exclusiveStartKey != nil {
		attributeMap, err = convertToAttribMap(exclusiveStartKey)
		if err != nil {
			return nil, nil, err
		}
	}

	queryResult, err := s.dynamoDBClient.QueryIndexWithPagination(ctx, s.tableName, expiryIndexName, "BlobStatus = :status AND Expiry > :expiry", commondynamodb.ExpresseionValues{
		":status": &types.AttributeValueMemberN{
			Value: strconv.Itoa(int(status)),
		},
		":expiry": &types.AttributeValueMemberN{
			Value: strconv.FormatInt(time.Now().Unix(), 10),
		},
	}, limit, attributeMap)

	if err != nil {
		return nil, nil, err
	}

	// When no more results to fetch, the LastEvaluatedKey is nil
	if queryResult.Items == nil && queryResult.LastEvaluatedKey == nil {
		return nil, nil, nil
	}

	metadata := make([]*v2.BlobMetadata, len(queryResult.Items))
	for i, item := range queryResult.Items {
		metadata[i], err = UnmarshalBlobMetadata(item)
		if err != nil {
			return nil, nil, err
		}
	}

	lastEvaluatedKey := queryResult.LastEvaluatedKey
	if lastEvaluatedKey == nil {
		return metadata, nil, nil
	}

	// Convert the last evaluated key to a disperser.BlobStoreExclusiveStartKey
	exclusiveStartKey, err = convertToExclusiveStartKey(lastEvaluatedKey)
	if err != nil {
		return nil, nil, err
	}
	return metadata, exclusiveStartKey, nil
}

func (s *BlobMetadataStore) GetAllBlobMetadataByBatch(ctx context.Context, batchHeaderHash [32]byte) ([]*v2.BlobMetadata, error) {
	items, err := s.dynamoDBClient.QueryIndex(ctx, s.tableName, batchIndexName, "BatchHeaderHash = :batch_header_hash", commondynamodb.ExpresseionValues{
		":batch_header_hash": &types.AttributeValueMemberB{
			Value: batchHeaderHash[:],
		},
	})
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("there is no metadata for batch %x", batchHeaderHash)
	}

	metadatas := make([]*v2.BlobMetadata, len(items))
	for i, item := range items {
		metadatas[i], err = UnmarshalBlobMetadata(item)
		if err != nil {
			return nil, err
		}
	}

	return metadatas, nil
}

// GetBlobMetadataByStatusWithPagination returns all the metadata with the given status upto the specified limit
// along with items, also returns a pagination token that can be used to fetch the next set of items
//
// Note that this may not return all the metadata for the batch if dynamodb query limit is reached.
// e.g 1mb limit for a single query
func (s *BlobMetadataStore) GetAllBlobMetadataByBatchWithPagination(
	ctx context.Context,
	batchHeaderHash [32]byte,
	limit int32,
	exclusiveStartKey *disperser.BatchIndexExclusiveStartKey,
) ([]*v2.BlobMetadata, *disperser.BatchIndexExclusiveStartKey, error) {
	var attributeMap map[string]types.AttributeValue
	var err error

	// Convert the exclusive start key to a map of AttributeValue
	if exclusiveStartKey != nil {
		attributeMap, err = convertToAttribMapBatchIndex(exclusiveStartKey)
		if err != nil {
			return nil, nil, err
		}
	}

	queryResult, err := s.dynamoDBClient.QueryIndexWithPagination(
		ctx,
		s.tableName,
		batchIndexName,
		"BatchHeaderHash = :batch_header_hash",
		commondynamodb.ExpresseionValues{
			":batch_header_hash": &types.AttributeValueMemberB{
				Value: batchHeaderHash[:],
			},
		},
		limit,
		attributeMap,
	)
	if err != nil {
		return nil, nil, err
	}

	s.logger.Info("Query result", "items", len(queryResult.Items), "lastEvaluatedKey", queryResult.LastEvaluatedKey)
	// When no more results to fetch, the LastEvaluatedKey is nil
	if queryResult.Items == nil && queryResult.LastEvaluatedKey == nil {
		return nil, nil, nil
	}

	metadata := make([]*v2.BlobMetadata, len(queryResult.Items))
	for i, item := range queryResult.Items {
		metadata[i], err = UnmarshalBlobMetadata(item)
		if err != nil {
			return nil, nil, err
		}
	}

	lastEvaluatedKey := queryResult.LastEvaluatedKey
	if lastEvaluatedKey == nil {
		return metadata, nil, nil
	}

	// Convert the last evaluated key to a disperser.BatchIndexExclusiveStartKey
	exclusiveStartKey, err = convertToExclusiveStartKeyBatchIndex(lastEvaluatedKey)
	if err != nil {
		return nil, nil, err
	}
	return metadata, exclusiveStartKey, nil
}

func (s *BlobMetadataStore) GetBlobMetadataInBatch(ctx context.Context, batchHeaderHash [32]byte, blobIndex uint32) (*v2.BlobMetadata, error) {
	items, err := s.dynamoDBClient.QueryIndex(ctx, s.tableName, batchIndexName, "BatchHeaderHash = :batch_header_hash AND BlobIndex = :blob_index", commondynamodb.ExpresseionValues{
		":batch_header_hash": &types.AttributeValueMemberB{
			Value: batchHeaderHash[:],
		},
		":blob_index": &types.AttributeValueMemberN{
			Value: strconv.Itoa(int(blobIndex)),
		}})
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("%w: there is no metadata for batch %s and blob index %d", disperser.ErrMetadataNotFound, hexutil.Encode(batchHeaderHash[:]), blobIndex)
	}

	if len(items) > 1 {
		s.logger.Error("there are multiple metadata for batch %s and blob index %d", hexutil.Encode(batchHeaderHash[:]), blobIndex)
	}

	metadata, err := UnmarshalBlobMetadata(items[0])
	if err != nil {
		return nil, err
	}
	return metadata, nil
}

func (s *BlobMetadataStore) IncrementNumRetries(ctx context.Context, existingMetadata *v2.BlobMetadata) error {
	_, err := s.dynamoDBClient.UpdateItem(ctx, s.tableName, map[string]types.AttributeValue{
		"BlobHash": &types.AttributeValueMemberS{
			Value: existingMetadata.BlobHash,
		},
		"MetadataHash": &types.AttributeValueMemberS{
			Value: existingMetadata.MetadataHash,
		},
	}, commondynamodb.Item{
		"NumRetries": &types.AttributeValueMemberN{
			Value: strconv.Itoa(int(existingMetadata.NumRetries + 1)),
		},
	})

	return err
}

func (s *BlobMetadataStore) UpdateConfirmationBlockNumber(ctx context.Context, existingMetadata *v2.BlobMetadata, confirmationBlockNumber uint32) error {
	updated := *existingMetadata
	if updated.ConfirmationInfo == nil {
		return fmt.Errorf("failed to update confirmation block number because confirmation info is missing for blob key %s", existingMetadata.GetBlobKey().String())
	}

	updated.ConfirmationInfo.ConfirmationBlockNumber = confirmationBlockNumber
	item, err := MarshalBlobMetadata(&updated)
	if err != nil {
		return err
	}

	_, err = s.dynamoDBClient.UpdateItem(ctx, s.tableName, map[string]types.AttributeValue{
		"BlobHash": &types.AttributeValueMemberS{
			Value: existingMetadata.BlobHash,
		},
		"MetadataHash": &types.AttributeValueMemberS{
			Value: existingMetadata.MetadataHash,
		},
	}, item)

	return err
}

func (s *BlobMetadataStore) UpdateBlobMetadata(ctx context.Context, metadataKey disperser.BlobKey, updated *v2.BlobMetadata) error {
	item, err := MarshalBlobMetadata(updated)
	if err != nil {
		return err
	}

	_, err = s.dynamoDBClient.UpdateItem(ctx, s.tableName, map[string]types.AttributeValue{
		"BlobHash": &types.AttributeValueMemberS{
			Value: metadataKey.BlobHash,
		},
		"MetadataHash": &types.AttributeValueMemberS{
			Value: metadataKey.MetadataHash,
		},
	}, item)

	return err
}

func (s *BlobMetadataStore) SetBlobStatus(ctx context.Context, metadataKey disperser.BlobKey, status disperser.BlobStatus) error {
	_, err := s.dynamoDBClient.UpdateItem(ctx, s.tableName, map[string]types.AttributeValue{
		"BlobHash": &types.AttributeValueMemberS{
			Value: metadataKey.BlobHash,
		},
		"MetadataHash": &types.AttributeValueMemberS{
			Value: metadataKey.MetadataHash,
		},
	}, commondynamodb.Item{
		"BlobStatus": &types.AttributeValueMemberN{
			Value: strconv.Itoa(int(status)),
		},
	})

	return err
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
		TableName: aws.String(tableName),
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String(StatusIndexName),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("BlobStatus"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("RequestedAt"),
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
				IndexName: aws.String(OperatorDispersalIndexName),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("OperatorID"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("DispersedAt"),
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
				IndexName: aws.String(OperatorResponseIndexName),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("OperatorID"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("RespondedAt"),
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

func MarshalBlobMetadata(metadata *v2.BlobMetadata) (commondynamodb.Item, error) {
	return attributevalue.MarshalMap(metadata)
}

func UnmarshalBlobMetadata(item commondynamodb.Item) (*v2.BlobMetadata, error) {
	metadata := v2.BlobMetadata{}
	err := attributevalue.UnmarshalMap(item, &metadata)
	if err != nil {
		return nil, err
	}

	return &metadata, nil
}

func convertToExclusiveStartKey(exclusiveStartKeyMap map[string]types.AttributeValue) (*disperser.BlobStoreExclusiveStartKey, error) {
	blobStoreExclusiveStartKey := disperser.BlobStoreExclusiveStartKey{}
	err := attributevalue.UnmarshalMap(exclusiveStartKeyMap, &blobStoreExclusiveStartKey)
	if err != nil {
		return nil, err
	}

	return &blobStoreExclusiveStartKey, nil
}

func convertToExclusiveStartKeyBatchIndex(exclusiveStartKeyMap map[string]types.AttributeValue) (*disperser.BatchIndexExclusiveStartKey, error) {
	blobStoreExclusiveStartKey := disperser.BatchIndexExclusiveStartKey{}
	err := attributevalue.UnmarshalMap(exclusiveStartKeyMap, &blobStoreExclusiveStartKey)
	if err != nil {
		return nil, err
	}

	return &blobStoreExclusiveStartKey, nil
}

func convertToAttribMap(blobStoreExclusiveStartKey *disperser.BlobStoreExclusiveStartKey) (map[string]types.AttributeValue, error) {
	if blobStoreExclusiveStartKey == nil {
		// Return an empty map or nil
		return nil, nil
	}

	avMap, err := attributevalue.MarshalMap(blobStoreExclusiveStartKey)
	if err != nil {
		return nil, err
	}
	return avMap, nil
}

func convertToAttribMapBatchIndex(blobStoreExclusiveStartKey *disperser.BatchIndexExclusiveStartKey) (map[string]types.AttributeValue, error) {
	if blobStoreExclusiveStartKey == nil {
		// Return an empty map or nil
		return nil, nil
	}

	avMap, err := attributevalue.MarshalMap(blobStoreExclusiveStartKey)
	if err != nil {
		return nil, err
	}
	return avMap, nil
}
