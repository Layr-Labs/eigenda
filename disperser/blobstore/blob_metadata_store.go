package blobstore

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	statusIndexName = "StatusIndex"
	batchIndexName  = "BatchIndex"
)

// BlobMetadataStore is a blob metadata storage backed by DynamoDB
// The blob metadata is stored in a single table and replicated in several indexes.
// - Metadata: (Partition Key: BlobKey, Sort Key: MetadataHash) -> Metadata
// - Indexes
//   - StatusIndex: (Partition Key: Status, Sort Key: RequestedAt) -> Metadata
//   - BatchIndex: (Partition Key: BatchHeaderHash, Sort Key: BlobIndex) -> Metadata
type BlobMetadataStore struct {
	dynamoDBClient *commondynamodb.Client
	logger         common.Logger
	tableName      string
	ttl            time.Duration
}

func NewBlobMetadataStore(dynamoDBClient *commondynamodb.Client, logger common.Logger, tableName string, ttl time.Duration) *BlobMetadataStore {
	logger.Debugf("creating blob metadata store with table %s with TTL: %s", tableName, ttl)
	return &BlobMetadataStore{
		dynamoDBClient: dynamoDBClient,
		logger:         logger,
		tableName:      tableName,
		ttl:            ttl,
	}
}

func (s *BlobMetadataStore) QueueNewBlobMetadata(ctx context.Context, blobMetadata *disperser.BlobMetadata) error {
	item, err := MarshalBlobMetadata(blobMetadata)
	if err != nil {
		return err
	}

	return s.dynamoDBClient.PutItem(ctx, s.tableName, item)
}

func (s *BlobMetadataStore) GetBlobMetadata(ctx context.Context, metadataKey disperser.BlobKey) (*disperser.BlobMetadata, error) {
	item, err := s.dynamoDBClient.GetItem(ctx, s.tableName, map[string]types.AttributeValue{
		"BlobHash": &types.AttributeValueMemberS{
			Value: metadataKey.BlobHash,
		},
		"MetadataHash": &types.AttributeValueMemberS{
			Value: metadataKey.MetadataHash,
		},
	})
	if err != nil {
		return nil, err
	}

	metadata, err := UnmarshalBlobMetadata(item)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

// GetBlobMetadataByStatus returns all the metadata with the given status
// Because this function scans the entire index, it should only be used for status with a limited number of items.
// It should only be used to filter "Processing" status. To support other status, a streaming version should be implemented.
func (s *BlobMetadataStore) GetBlobMetadataByStatus(ctx context.Context, status disperser.BlobStatus) ([]*disperser.BlobMetadata, error) {
	items, err := s.dynamoDBClient.QueryIndex(ctx, s.tableName, statusIndexName, "BlobStatus = :status", commondynamodb.ExpresseionValues{
		":status": &types.AttributeValueMemberN{
			Value: strconv.Itoa(int(status)),
		}})
	if err != nil {
		return nil, err
	}

	metadata := make([]*disperser.BlobMetadata, len(items))
	for i, item := range items {
		metadata[i], err = UnmarshalBlobMetadata(item)
		if err != nil {
			return nil, err
		}
	}

	return metadata, nil
}

func (s *BlobMetadataStore) GetAllBlobMetadataByBatch(ctx context.Context, batchHeaderHash [32]byte) ([]*disperser.BlobMetadata, error) {
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

	metadatas := make([]*disperser.BlobMetadata, len(items))
	for i, item := range items {
		metadatas[i], err = UnmarshalBlobMetadata(item)
		if err != nil {
			return nil, err
		}
	}

	return metadatas, nil
}

func (s *BlobMetadataStore) GetBlobMetadataInBatch(ctx context.Context, batchHeaderHash [32]byte, blobIndex uint32) (*disperser.BlobMetadata, error) {
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
		return nil, fmt.Errorf("there is no metadata for batch %s and blob index %d", batchHeaderHash, blobIndex)
	}

	if len(items) > 1 {
		s.logger.Error("there are multiple metadata for batch %s and blob index %d", batchHeaderHash, blobIndex)
	}

	metadata, err := UnmarshalBlobMetadata(items[0])
	if err != nil {
		return nil, err
	}
	return metadata, nil
}

func (s *BlobMetadataStore) IncrementNumRetries(ctx context.Context, existingMetadata *disperser.BlobMetadata) error {
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

func (s *BlobMetadataStore) UpdateBlobMetadata(ctx context.Context, metadataKey disperser.BlobKey, updated *disperser.BlobMetadata) error {
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

func GenerateTableSchema(metadataTableName string, readCapacityUnits int64, writeCapacityUnits int64) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("BlobHash"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("MetadataHash"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("BlobStatus"),
				AttributeType: types.ScalarAttributeTypeN,
			},
			{
				AttributeName: aws.String("RequestedAt"),
				AttributeType: types.ScalarAttributeTypeN,
			},
			{
				AttributeName: aws.String("BatchHeaderHash"),
				AttributeType: types.ScalarAttributeTypeB,
			},
			{
				AttributeName: aws.String("BlobIndex"),
				AttributeType: types.ScalarAttributeTypeN,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("BlobHash"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("MetadataHash"),
				KeyType:       types.KeyTypeRange,
			},
		},
		TableName: aws.String(metadataTableName),
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String(statusIndexName),
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
				IndexName: aws.String(batchIndexName),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("BatchHeaderHash"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("BlobIndex"),
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

func MarshalBlobMetadata(metadata *disperser.BlobMetadata) (commondynamodb.Item, error) {
	basicFields, err := attributevalue.MarshalMap(metadata)
	if err != nil {
		return nil, err
	}

	if metadata.RequestMetadata == nil {
		return basicFields, nil
	}

	requestMetadata, err := attributevalue.MarshalMap(metadata.RequestMetadata)
	if err != nil {
		return nil, err
	}

	// Flatten the request metadata
	for k, v := range requestMetadata {
		basicFields[k] = v
	}

	if metadata.ConfirmationInfo == nil {
		return basicFields, nil
	}

	confirmationInfo, err := attributevalue.MarshalMap(metadata.ConfirmationInfo)
	if err != nil {
		return nil, err
	}

	// Flatten the confirmation info
	for k, v := range confirmationInfo {
		basicFields[k] = v
	}

	return basicFields, nil
}

func UnmarshalBlobMetadata(item commondynamodb.Item) (*disperser.BlobMetadata, error) {
	metadata := disperser.BlobMetadata{}
	err := attributevalue.UnmarshalMap(item, &metadata)
	if err != nil {
		return nil, err
	}

	requestMetadata := disperser.RequestMetadata{}
	err = attributevalue.UnmarshalMap(item, &requestMetadata)
	if err != nil {
		return nil, err
	}
	metadata.RequestMetadata = &requestMetadata
	if metadata.BlobStatus != disperser.Confirmed && metadata.BlobStatus != disperser.Finalized {
		return &metadata, nil
	}

	confirmationInfo := disperser.ConfirmationInfo{}
	err = attributevalue.UnmarshalMap(item, &confirmationInfo)
	if err != nil {
		return nil, err
	}
	metadata.ConfirmationInfo = &confirmationInfo

	return &metadata, nil
}
