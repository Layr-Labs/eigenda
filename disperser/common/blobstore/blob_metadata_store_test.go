package blobstore_test

import (
	"context"
	"testing"
	"time"

	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestBlobMetadataStoreOperations(t *testing.T) {
	ctx := context.Background()
	blobKey1 := disperser.BlobKey{
		BlobHash:     blobHash,
		MetadataHash: "hash",
	}
	now := time.Now()
	metadata1 := &disperser.BlobMetadata{
		MetadataHash: blobKey1.MetadataHash,
		BlobHash:     blobHash,
		BlobStatus:   disperser.Processing,
		Expiry:       uint64(now.Add(time.Hour).Unix()),
		NumRetries:   0,
		RequestMetadata: &disperser.RequestMetadata{
			BlobRequestHeader: blob.RequestHeader,
			BlobSize:          blobSize,
			RequestedAt:       uint64(now.Unix()),
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
		Expiry:       uint64(now.Add(time.Hour).Unix()),
		NumRetries:   0,
		RequestMetadata: &disperser.RequestMetadata{
			BlobRequestHeader: blob.RequestHeader,
			BlobSize:          blobSize,
			RequestedAt:       uint64(now.Unix()),
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

	fetchBulk, err := blobMetadataStore.GetBulkBlobMetadata(ctx, []disperser.BlobKey{blobKey1, blobKey2})
	assert.NoError(t, err)
	assert.Equal(t, metadata1, fetchBulk[0])
	assert.Equal(t, metadata2, fetchBulk[1])

	processing, err := blobMetadataStore.GetBlobMetadataByStatus(ctx, disperser.Processing)
	assert.NoError(t, err)
	assert.Len(t, processing, 1)
	assert.Equal(t, metadata1, processing[0])

	processingCount, err := blobMetadataStore.GetBlobMetadataCountByStatus(ctx, disperser.Processing)
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

	finalizedCount, err := blobMetadataStore.GetBlobMetadataCountByStatus(ctx, disperser.Finalized)
	assert.NoError(t, err)
	assert.Equal(t, int32(1), finalizedCount)

	confirmedMetadata := getConfirmedMetadata(t, metadata1, 1)
	err = blobMetadataStore.UpdateBlobMetadata(ctx, blobKey1, confirmedMetadata)
	assert.NoError(t, err)

	metadata, err := blobMetadataStore.GetBlobMetadataInBatch(ctx, confirmedMetadata.ConfirmationInfo.BatchHeaderHash, confirmedMetadata.ConfirmationInfo.BlobIndex)
	assert.NoError(t, err)
	assert.Equal(t, metadata, confirmedMetadata)

	confirmedCount, err := blobMetadataStore.GetBlobMetadataCountByStatus(ctx, disperser.Confirmed)
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
	now := time.Now()
	metadata1 := &disperser.BlobMetadata{
		MetadataHash: blobKey1.MetadataHash,
		BlobHash:     blobHash,
		BlobStatus:   disperser.Processing,
		Expiry:       uint64(now.Add(time.Hour).Unix()),
		NumRetries:   0,
		RequestMetadata: &disperser.RequestMetadata{
			BlobRequestHeader: blob.RequestHeader,
			BlobSize:          blobSize,
			RequestedAt:       uint64(now.Unix()),
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
		Expiry:       uint64(now.Add(time.Hour).Unix()),
		NumRetries:   0,
		RequestMetadata: &disperser.RequestMetadata{
			BlobRequestHeader: blob.RequestHeader,
			BlobSize:          blobSize,
			RequestedAt:       uint64(now.Unix()),
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

func TestGetAllBlobMetadataByBatchWithPagination(t *testing.T) {
	ctx := context.Background()
	blobKey1 := disperser.BlobKey{
		BlobHash:     blobHash,
		MetadataHash: "hash",
	}
	expiry := uint64(time.Now().Add(time.Hour).Unix())
	metadata1 := &disperser.BlobMetadata{
		MetadataHash: blobKey1.MetadataHash,
		BlobHash:     blobHash,
		BlobStatus:   disperser.Processing,
		Expiry:       expiry,
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
		Expiry:       expiry,
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

	confirmedMetadata1 := getConfirmedMetadata(t, metadata1, 1)
	err = blobMetadataStore.UpdateBlobMetadata(ctx, blobKey1, confirmedMetadata1)
	assert.NoError(t, err)

	confirmedMetadata2 := getConfirmedMetadata(t, metadata2, 2)
	err = blobMetadataStore.UpdateBlobMetadata(ctx, blobKey2, confirmedMetadata2)
	assert.NoError(t, err)

	// Fetch the blob metadata with limit 1
	metadata, exclusiveStartKey, err := blobMetadataStore.GetAllBlobMetadataByBatchWithPagination(ctx, confirmedMetadata1.ConfirmationInfo.BatchHeaderHash, 1, nil)
	assert.NoError(t, err)
	assert.Equal(t, metadata[0], confirmedMetadata1)
	assert.NotNil(t, exclusiveStartKey)
	assert.Equal(t, confirmedMetadata1.ConfirmationInfo.BlobIndex, exclusiveStartKey.BlobIndex)

	// Get the next blob metadata with limit 1 and the exclusive start key
	metadata, exclusiveStartKey, err = blobMetadataStore.GetAllBlobMetadataByBatchWithPagination(ctx, confirmedMetadata1.ConfirmationInfo.BatchHeaderHash, 1, exclusiveStartKey)
	assert.NoError(t, err)
	assert.Equal(t, metadata[0], confirmedMetadata2)
	assert.Equal(t, confirmedMetadata2.ConfirmationInfo.BlobIndex, exclusiveStartKey.BlobIndex)

	// Fetching the next blob metadata should return an empty list
	metadata, exclusiveStartKey, err = blobMetadataStore.GetAllBlobMetadataByBatchWithPagination(ctx, confirmedMetadata1.ConfirmationInfo.BatchHeaderHash, 1, exclusiveStartKey)
	assert.NoError(t, err)
	assert.Len(t, metadata, 0)
	assert.Nil(t, exclusiveStartKey)

	// Fetch the blob metadata with limit 2
	metadata, exclusiveStartKey, err = blobMetadataStore.GetAllBlobMetadataByBatchWithPagination(ctx, confirmedMetadata1.ConfirmationInfo.BatchHeaderHash, 2, nil)
	assert.NoError(t, err)
	assert.Len(t, metadata, 2)
	assert.Equal(t, metadata[0], confirmedMetadata1)
	assert.Equal(t, metadata[1], confirmedMetadata2)
	assert.NotNil(t, exclusiveStartKey)
	assert.Equal(t, confirmedMetadata2.ConfirmationInfo.BlobIndex, exclusiveStartKey.BlobIndex)

	// Fetch the blob metadata with limit 3 should return only 2 items
	metadata, exclusiveStartKey, err = blobMetadataStore.GetAllBlobMetadataByBatchWithPagination(ctx, confirmedMetadata1.ConfirmationInfo.BatchHeaderHash, 3, nil)
	assert.NoError(t, err)
	assert.Len(t, metadata, 2)
	assert.Equal(t, metadata[0], confirmedMetadata1)
	assert.Equal(t, metadata[1], confirmedMetadata2)
	assert.Nil(t, exclusiveStartKey)

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

func TestFilterOutExpiredBlobMetadata(t *testing.T) {
	ctx := context.Background()

	blobKey := disperser.BlobKey{
		BlobHash:     "blob1",
		MetadataHash: "hash1",
	}
	now := time.Now()
	metadata := &disperser.BlobMetadata{
		MetadataHash: blobKey.MetadataHash,
		BlobHash:     blobKey.BlobHash,
		BlobStatus:   disperser.Processing,
		Expiry:       uint64(now.Add(-1).Unix()),
		NumRetries:   0,
		RequestMetadata: &disperser.RequestMetadata{
			BlobRequestHeader: blob.RequestHeader,
			BlobSize:          blobSize,
			RequestedAt:       uint64(now.Add(-1000).Unix()),
		},
		ConfirmationInfo: &disperser.ConfirmationInfo{},
	}

	err := blobMetadataStore.QueueNewBlobMetadata(ctx, metadata)
	assert.NoError(t, err)

	processing, err := blobMetadataStore.GetBlobMetadataByStatus(ctx, disperser.Processing)
	assert.NoError(t, err)
	assert.Len(t, processing, 0)

	processingCount, err := blobMetadataStore.GetBlobMetadataCountByStatus(ctx, disperser.Processing)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), processingCount)

	processing, _, err = blobMetadataStore.GetBlobMetadataByStatusWithPagination(ctx, disperser.Processing, 10, nil)
	assert.NoError(t, err)
	assert.Len(t, processing, 0)

	deleteItems(t, []commondynamodb.Key{
		{
			"MetadataHash": &types.AttributeValueMemberS{Value: blobKey.MetadataHash},
			"BlobHash":     &types.AttributeValueMemberS{Value: blobKey.BlobHash},
		},
	})
}

func deleteItems(t *testing.T, keys []commondynamodb.Key) {
	_, err := dynamoClient.DeleteItems(context.Background(), metadataTableName, keys)
	assert.NoError(t, err)
}

func getConfirmedMetadata(t *testing.T, metadata *disperser.BlobMetadata, blobIndex uint32) *disperser.BlobMetadata {
	batchHeaderHash := [32]byte{1, 2, 3}
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
	confirmationInfo := &disperser.ConfirmationInfo{
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
	}
	metadata.BlobStatus = disperser.Confirmed
	metadata.ConfirmationInfo = confirmationInfo
	return metadata
}
