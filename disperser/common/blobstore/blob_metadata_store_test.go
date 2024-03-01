package blobstore_test

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"

	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/common/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDynamoDBClient is a mock of DynamoDB interface
type MockDynamoDBClient struct {
	mock.Mock
}

func TestBlobMetadataStoreOperations(t *testing.T) {
	ctx := context.Background()
	blobKey1 := disperser.BlobKey{
		BlobHash:     blobHash,
		MetadataHash: "hash",
	}
	metadata1 := &disperser.BlobMetadata{
		MetadataHash: blobKey1.MetadataHash,
		BlobHash:     blobHash,
		BlobStatus:   disperser.Processing,
		Expiry:       0,
		NumRetries:   0,
		RequestMetadata: &disperser.RequestMetadata{
			BlobRequestHeader: blob.RequestHeader,
			BlobSize:          blobSize,
			RequestedAt:       123,
		},
	}
	blobKey2 := disperser.BlobKey{
		BlobHash:     "blob2",
		MetadataHash: "hash2",
	}
	metadata2 := &disperser.BlobMetadata{
		MetadataHash: blobKey2.MetadataHash,
		BlobHash:     blobKey2.BlobHash,
		BlobStatus:   disperser.Finalized,
		Expiry:       0,
		NumRetries:   0,
		RequestMetadata: &disperser.RequestMetadata{
			BlobRequestHeader: blob.RequestHeader,
			BlobSize:          blobSize,
			RequestedAt:       123,
		},
		ConfirmationInfo: &disperser.ConfirmationInfo{},
	}
	err := blobMetadataStore.QueueNewBlobMetadata(ctx, metadata1)
	assert.NoError(t, err)
	err = blobMetadataStore.QueueNewBlobMetadata(ctx, metadata2)
	assert.NoError(t, err)

	fetchedMetadata, err := blobMetadataStore.GetBlobMetadata(ctx, blobKey1)
	assert.NoError(t, err)
	assert.Equal(t, metadata1, fetchedMetadata)
	fetchedMetadata, err = blobMetadataStore.GetBlobMetadata(ctx, blobKey2)
	assert.NoError(t, err)
	assert.Equal(t, metadata2, fetchedMetadata)

	processing, err := blobMetadataStore.GetBlobMetadataByStatus(ctx, disperser.Processing)
	assert.NoError(t, err)
	assert.Len(t, processing, 1)
	assert.Equal(t, metadata1, processing[0])

	processingCount, err := blobMetadataStore.GetBlobMetadataByStatusCount(ctx, disperser.Processing)
	assert.NoError(t, err)
	assert.Equal(t, int32(1), processingCount)

	err = blobMetadataStore.IncrementNumRetries(ctx, metadata1)
	assert.NoError(t, err)
	fetchedMetadata, err = blobMetadataStore.GetBlobMetadata(ctx, blobKey1)
	assert.NoError(t, err)
	metadata1.NumRetries = 1
	assert.Equal(t, metadata1, fetchedMetadata)

	finalized, err := blobMetadataStore.GetBlobMetadataByStatus(ctx, disperser.Finalized)
	assert.NoError(t, err)
	assert.Len(t, finalized, 1)
	assert.Equal(t, metadata2, finalized[0])

	finalizedCount, err := blobMetadataStore.GetBlobMetadataByStatusCount(ctx, disperser.Finalized)
	assert.NoError(t, err)
	assert.Equal(t, int32(1), finalizedCount)

	confirmedMetadata := getConfirmedMetadata(t, blobKey1)
	err = blobMetadataStore.UpdateBlobMetadata(ctx, blobKey1, confirmedMetadata)
	assert.NoError(t, err)

	metadata, err := blobMetadataStore.GetBlobMetadataInBatch(ctx, confirmedMetadata.ConfirmationInfo.BatchHeaderHash, confirmedMetadata.ConfirmationInfo.BlobIndex)
	assert.NoError(t, err)
	assert.Equal(t, metadata, confirmedMetadata)

	confirmedCount, err := blobMetadataStore.GetBlobMetadataByStatusCount(ctx, disperser.Confirmed)
	assert.NoError(t, err)
	assert.Equal(t, int32(1), confirmedCount)

	deleteItems(t, []commondynamodb.Key{
		{
			"MetadataHash": &types.AttributeValueMemberS{Value: blobKey1.MetadataHash},
			"BlobHash":     &types.AttributeValueMemberS{Value: blobKey1.BlobHash},
		},
		{
			"MetadataHash": &types.AttributeValueMemberS{Value: blobKey2.MetadataHash},
			"BlobHash":     &types.AttributeValueMemberS{Value: blobKey2.BlobHash},
		},
	})
}

func TestBlobMetadataStoreOperationsWithPagination(t *testing.T) {
	ctx := context.Background()
	blobKey1 := disperser.BlobKey{
		BlobHash:     blobHash,
		MetadataHash: "hash",
	}
	metadata1 := &disperser.BlobMetadata{
		MetadataHash: blobKey1.MetadataHash,
		BlobHash:     blobHash,
		BlobStatus:   disperser.Processing,
		Expiry:       0,
		NumRetries:   0,
		RequestMetadata: &disperser.RequestMetadata{
			BlobRequestHeader: blob.RequestHeader,
			BlobSize:          blobSize,
			RequestedAt:       123,
		},
	}
	blobKey2 := disperser.BlobKey{
		BlobHash:     "blob2",
		MetadataHash: "hash2",
	}
	metadata2 := &disperser.BlobMetadata{
		MetadataHash: blobKey2.MetadataHash,
		BlobHash:     blobKey2.BlobHash,
		BlobStatus:   disperser.Finalized,
		Expiry:       0,
		NumRetries:   0,
		RequestMetadata: &disperser.RequestMetadata{
			BlobRequestHeader: blob.RequestHeader,
			BlobSize:          blobSize,
			RequestedAt:       123,
		},
		ConfirmationInfo: &disperser.ConfirmationInfo{},
	}
	err := blobMetadataStore.QueueNewBlobMetadata(ctx, metadata1)
	assert.NoError(t, err)
	err = blobMetadataStore.QueueNewBlobMetadata(ctx, metadata2)
	assert.NoError(t, err)

	fetchedMetadata, err := blobMetadataStore.GetBlobMetadata(ctx, blobKey1)
	assert.NoError(t, err)
	assert.Equal(t, metadata1, fetchedMetadata)
	fetchedMetadata, err = blobMetadataStore.GetBlobMetadata(ctx, blobKey2)
	assert.NoError(t, err)
	assert.Equal(t, metadata2, fetchedMetadata)

	processing, lastEvaluatedKey, err := blobMetadataStore.GetBlobMetadataByStatusWithPagination(ctx, disperser.Processing, 1, nil)
	assert.NoError(t, err)
	assert.Len(t, processing, 1)
	assert.Equal(t, metadata1, processing[0])
	assert.NotNil(t, lastEvaluatedKey)

	finalized, lastEvaluatedKey, err := blobMetadataStore.GetBlobMetadataByStatusWithPagination(ctx, disperser.Finalized, 1, nil)
	assert.NoError(t, err)
	assert.Len(t, finalized, 1)
	assert.Equal(t, metadata2, finalized[0])
	assert.NotNil(t, lastEvaluatedKey)

	finalized, lastEvaluatedKey, err = blobMetadataStore.GetBlobMetadataByStatusWithPagination(ctx, disperser.Finalized, 1, lastEvaluatedKey)
	assert.NoError(t, err)
	assert.Len(t, finalized, 0)
	assert.Nil(t, lastEvaluatedKey)

	deleteItems(t, []commondynamodb.Key{
		{
			"MetadataHash": &types.AttributeValueMemberS{Value: blobKey1.MetadataHash},
			"BlobHash":     &types.AttributeValueMemberS{Value: blobKey1.BlobHash},
		},
		{
			"MetadataHash": &types.AttributeValueMemberS{Value: blobKey2.MetadataHash},
			"BlobHash":     &types.AttributeValueMemberS{Value: blobKey2.BlobHash},
		},
	})
}

func TestBlobMetadataStoreOperationsWithPaginationNoStoredBlob(t *testing.T) {
	ctx := context.Background()
	// Query BlobMetadataStore for a blob that does not exist
	// This should return nil for both the blob and lastEvaluatedKey
	processing, lastEvaluatedKey, err := blobMetadataStore.GetBlobMetadataByStatusWithPagination(ctx, disperser.Processing, 1, nil)
	assert.NoError(t, err)
	assert.Nil(t, processing)
	assert.Nil(t, lastEvaluatedKey)
}

func deleteItems(t *testing.T, keys []commondynamodb.Key) {
	_, err := dynamoClient.DeleteItems(context.Background(), metadataTableName, keys)
	assert.NoError(t, err)
}

func getConfirmedMetadata(t *testing.T, metadataKey disperser.BlobKey) *disperser.BlobMetadata {
	batchHeaderHash := [32]byte{1, 2, 3}
	blobIndex := uint32(1)
	requestedAt := uint64(time.Now().Nanosecond())
	var commitX, commitY fp.Element
	_, err := commitX.SetString("21661178944771197726808973281966770251114553549453983978976194544185382599016")
	assert.NoError(t, err)
	_, err = commitY.SetString("9207254729396071334325696286939045899948985698134704137261649190717970615186")
	assert.NoError(t, err)
	commitment := &encoding.G1Commitment{
		X: commitX,
		Y: commitY,
	}
	dataLength := 32
	batchID := uint32(99)
	batchRoot := []byte("hello")
	referenceBlockNumber := uint32(132)
	confirmationBlockNumber := uint32(150)
	sigRecordHash := [32]byte{0}
	fee := []byte{0}
	inclusionProof := []byte{1, 2, 3, 4, 5}
	return &disperser.BlobMetadata{
		BlobHash:     metadataKey.BlobHash,
		MetadataHash: metadataKey.MetadataHash,
		BlobStatus:   disperser.Confirmed,
		Expiry:       0,
		NumRetries:   0,
		RequestMetadata: &disperser.RequestMetadata{
			BlobRequestHeader: core.BlobRequestHeader{
				SecurityParams: securityParams,
			},
			RequestedAt: requestedAt,
			BlobSize:    blobSize,
		},
		ConfirmationInfo: &disperser.ConfirmationInfo{
			BatchHeaderHash:      batchHeaderHash,
			BlobIndex:            blobIndex,
			SignatoryRecordHash:  sigRecordHash,
			ReferenceBlockNumber: referenceBlockNumber,
			BatchRoot:            batchRoot,
			BlobInclusionProof:   inclusionProof,
			BlobCommitment: &encoding.BlobCommitments{
				Commitment: commitment,
				Length:     uint(dataLength),
			},
			BatchID:                 batchID,
			ConfirmationTxnHash:     common.HexToHash("0x123"),
			ConfirmationBlockNumber: confirmationBlockNumber,
			Fee:                     fee,
		},
	}
}

func TestGetBlobMetadataByStatusWithPagination_ErrorQueryingDynamoDB(t *testing.T) {
	ctx := context.Background()

	// Use your setup for mocking the DynamoDB client
	mockDynamoDBClient := new(MockDynamoDBClient)
	mockDynamoDBClient.On("QueryIndexWithPagination", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(commondynamodb.QueryResult{}, errors.New("dynamodb query error"))
	blobMetadataStore = blobstore.NewBlobMetadataStore(mockDynamoDBClient, logger, metadataTableName, time.Hour)

	metadatas, key, err := blobMetadataStore.GetBlobMetadataByStatusWithPagination(ctx, disperser.Confirmed, 1, nil)
	assert.Equal(t, 0, len(metadatas))
	assert.Nil(t, key)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "dynamodb query error")
}

func TestGetBlobMetadataByStatusWithPagination_WithyResult(t *testing.T) {
	ctx := context.Background()
	blobKey1 := disperser.BlobKey{
		BlobHash:     blobHash,
		MetadataHash: "hash",
	}

	metadata1 := &disperser.BlobMetadata{
		MetadataHash: blobKey1.MetadataHash,
		BlobHash:     blobHash,
		BlobStatus:   disperser.Processing,
		Expiry:       0,
		NumRetries:   0,
		RequestMetadata: &disperser.RequestMetadata{
			BlobRequestHeader: blob.RequestHeader,
			BlobSize:          blobSize,
			RequestedAt:       123,
		},
	}
	item, err := attributevalue.MarshalMap(metadata1)
	if err != nil {
		// Handle the error, perhaps log it or return it
		log.Fatalf("failed to marshal metadata to Item: %v", err)
	}

	var result commondynamodb.QueryResult
	result.Items = append(result.Items, item)
	result.LastEvaluatedKey = map[string]types.AttributeValue{
		"BlobHash":     &types.AttributeValueMemberS{Value: blobKey1.BlobHash},
		"MetadataHash": &types.AttributeValueMemberS{Value: blobKey1.MetadataHash},
	}
	// Use your setup for mocking the DynamoDB client
	mockDynamoDBClient := new(MockDynamoDBClient)
	mockDynamoDBClient.On("QueryIndexWithPagination", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(result, nil)
	blobMetadataStore = blobstore.NewBlobMetadataStore(mockDynamoDBClient, logger, metadataTableName, time.Hour)

	metadatas, key, err := blobMetadataStore.GetBlobMetadataByStatusWithPagination(ctx, disperser.Confirmed, 1, nil)
	assert.Equal(t, 1, len(metadatas))
	assert.NotNil(t, key)
	assert.NoError(t, err)
}

func (m *MockDynamoDBClient) PutItem(ctx context.Context, tableName string, item map[string]types.AttributeValue) error {
	args := m.Called(ctx, tableName, item)
	return args.Error(0)
}

func (m *MockDynamoDBClient) PutItems(ctx context.Context, tableName string, items []commondynamodb.Item) ([]commondynamodb.Item, error) {
	args := m.Called(ctx, tableName, items)
	return args.Get(0).([]commondynamodb.Item), args.Error(1)
}

func (m *MockDynamoDBClient) GetItem(ctx context.Context, tableName string, key map[string]types.AttributeValue) (map[string]types.AttributeValue, error) {
	args := m.Called(ctx, tableName, key)
	return args.Get(0).(map[string]types.AttributeValue), args.Error(1)
}

func (m *MockDynamoDBClient) QueryIndex(ctx context.Context, tableName string, indexName string, keyCondition string, expressionValues commondynamodb.ExpresseionValues) ([]commondynamodb.Item, error) {
	args := m.Called(ctx, tableName, indexName, keyCondition, expressionValues)
	return args.Get(0).([]commondynamodb.Item), args.Error(1)
}

func (m *MockDynamoDBClient) QueryIndexWithPagination(ctx context.Context, tableName string, indexName string, keyCondition string, expressionValues commondynamodb.ExpresseionValues, limit int32, exclusiveStartKey map[string]types.AttributeValue) (commondynamodb.QueryResult, error) {
	args := m.Called(ctx, tableName, indexName, keyCondition, expressionValues, limit, exclusiveStartKey)
	return args.Get(0).(commondynamodb.QueryResult), args.Error(1)
}

func (m *MockDynamoDBClient) QueryIndexCount(ctx context.Context, tableName string, indexName string, keyCondition string, expressionValues commondynamodb.ExpresseionValues) (int32, error) {
	args := m.Called(ctx, tableName, indexName, keyCondition, expressionValues)
	return int32(args.Int(0)), args.Error(1)
}

func (m *MockDynamoDBClient) UpdateItem(ctx context.Context, tableName string, key commondynamodb.Key, update commondynamodb.Item) (commondynamodb.Item, error) {
	args := m.Called(ctx, tableName, key, update)
	return args.Get(0).(commondynamodb.Item), args.Error(1)
}

func (m *MockDynamoDBClient) DeleteItem(ctx context.Context, tableName string, key commondynamodb.Key) error {
	args := m.Called(ctx, tableName, key)
	return args.Error(0)
}

func (m *MockDynamoDBClient) DeleteItems(ctx context.Context, tableName string, keys []commondynamodb.Key) ([]commondynamodb.Key, error) {
	args := m.Called(ctx, tableName, keys)
	return args.Get(0).([]commondynamodb.Key), args.Error(1)
}

func (m *MockDynamoDBClient) DeleteTable(ctx context.Context, tableName string) error {
	args := m.Called(ctx, tableName)
	return args.Error(0)
}
