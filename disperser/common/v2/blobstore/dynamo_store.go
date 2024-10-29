package blobstore

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser"
	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	StatusIndexName            = "StatusIndex"
	OperatorDispersalIndexName = "OperatorDispersalIndex"
	OperatorResponseIndexName  = "OperatorResponseIndex"

	blobKeyPrefix  = "BlobKey#"
	blobMetadataSK = "BlobMetadata"
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
		logger:         logger.With("component", "blobMetadataStoreV2"),
		tableName:      tableName,
		ttl:            ttl,
	}
}

func (s *BlobMetadataStore) PutBlobMetadata(ctx context.Context, blobMetadata *v2.BlobMetadata) error {
	item, err := MarshalBlobMetadata(blobMetadata)
	if err != nil {
		return err
	}

	return s.dynamoDBClient.PutItem(ctx, s.tableName, item)
}

func (s *BlobMetadataStore) GetBlobMetadata(ctx context.Context, blobKey core.BlobKey) (*v2.BlobMetadata, error) {
	item, err := s.dynamoDBClient.GetItem(ctx, s.tableName, map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{
			Value: blobKeyPrefix + blobKey.Hex(),
		},
		"SK": &types.AttributeValueMemberS{
			Value: blobMetadataSK,
		},
	})

	if item == nil {
		return nil, fmt.Errorf("%w: metadata not found for key %s", disperser.ErrMetadataNotFound, blobKey.Hex())
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

// GetBlobMetadataByStatus returns all the metadata with the given status
// Because this function scans the entire index, it should only be used for status with a limited number of items.
func (s *BlobMetadataStore) GetBlobMetadataByStatus(ctx context.Context, status v2.BlobStatus, lastUpdatedAt uint64) ([]*v2.BlobMetadata, error) {
	items, err := s.dynamoDBClient.QueryIndex(ctx, s.tableName, StatusIndexName, "BlobStatus = :status AND UpdatedAt > :updatedAt", commondynamodb.ExpressionValues{
		":status": &types.AttributeValueMemberN{
			Value: strconv.Itoa(int(status)),
		},
		":updatedAt": &types.AttributeValueMemberN{
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
func (s *BlobMetadataStore) GetBlobMetadataCountByStatus(ctx context.Context, status v2.BlobStatus) (int32, error) {
	count, err := s.dynamoDBClient.QueryIndexCount(ctx, s.tableName, StatusIndexName, "BlobStatus = :status", commondynamodb.ExpressionValues{
		":status": &types.AttributeValueMemberN{
			Value: strconv.Itoa(int(status)),
		},
	})
	if err != nil {
		return 0, err
	}

	return count, nil
}

func GenerateTableSchema(tableName string, readCapacityUnits int64, writeCapacityUnits int64) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			// PK is the composite partition key
			{
				AttributeName: aws.String("PK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			// SK is the composite sort key
			{
				AttributeName: aws.String("SK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("BlobStatus"),
				AttributeType: types.ScalarAttributeTypeN,
			},
			{
				AttributeName: aws.String("UpdatedAt"),
				AttributeType: types.ScalarAttributeTypeN,
			},
			{
				AttributeName: aws.String("OperatorID"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("DispersedAt"),
				AttributeType: types.ScalarAttributeTypeN,
			},
			{
				AttributeName: aws.String("RespondedAt"),
				AttributeType: types.ScalarAttributeTypeN,
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
						AttributeName: aws.String("UpdatedAt"),
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
	fields, err := attributevalue.MarshalMap(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal blob metadata: %w", err)
	}

	// Add PK and SK fields
	blobKey, err := metadata.BlobHeader.BlobKey()
	if err != nil {
		return nil, err
	}
	fields["PK"] = &types.AttributeValueMemberS{Value: blobKeyPrefix + blobKey.Hex()}
	fields["SK"] = &types.AttributeValueMemberS{Value: blobMetadataSK}

	return fields, nil
}

func UnmarshalBlobKey(item commondynamodb.Item) (core.BlobKey, error) {
	type Blob struct {
		PK string
	}

	blob := Blob{}
	err := attributevalue.UnmarshalMap(item, &blob)
	if err != nil {
		return core.BlobKey{}, err
	}

	bk := strings.TrimPrefix(blob.PK, blobKeyPrefix)
	return core.HexToBlobKey(bk)
}

func UnmarshalBlobMetadata(item commondynamodb.Item) (*v2.BlobMetadata, error) {
	metadata := v2.BlobMetadata{}
	err := attributevalue.UnmarshalMap(item, &metadata)
	if err != nil {
		return nil, err
	}
	return &metadata, nil
}
