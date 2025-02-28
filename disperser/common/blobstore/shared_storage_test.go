package blobstore_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"

	"github.com/ethereum/go-ethereum/common"
)

func TestSharedBlobStore(t *testing.T) {
	requestedAt := uint64(time.Now().UnixNano())
	ctx := context.Background()
	blobKey, err := sharedStorage.StoreBlob(ctx, blob, requestedAt)
	assert.Nil(t, err)
	assert.Equal(t, blobHash, blobKey.BlobHash)

	metadatas, err := sharedStorage.GetBlobMetadataByStatus(ctx, disperser.Processing)
	assert.Nil(t, err)
	assert.Len(t, metadatas, 1)
	assertMetadata(t, blobKey, blobSize, requestedAt, disperser.Processing, metadatas[0])

	blobs, err := sharedStorage.GetBlobsByMetadata(ctx, metadatas)
	assert.Nil(t, err)
	assert.Len(t, blobs, 1)
	assertBlob(t, blobs[blobKey])

	data, err := sharedStorage.GetBlobContent(ctx, blobKey.BlobHash)
	assert.Nil(t, err)
	assert.Equal(t, blob.Data, data)

	err = sharedStorage.MarkBlobFailed(ctx, blobKey)
	assert.Nil(t, err)

	metadata1, err := sharedStorage.GetBlobMetadata(ctx, blobKey)
	assert.Nil(t, err)
	assertMetadata(t, blobKey, blobSize, requestedAt, disperser.Failed, metadata1)

	err = sharedStorage.MarkBlobProcessing(ctx, blobKey)
	assert.Nil(t, err)

	metadata1, err = sharedStorage.GetBlobMetadata(ctx, blobKey)
	assert.Nil(t, err)
	assertMetadata(t, blobKey, blobSize, requestedAt, disperser.Processing, metadata1)

	err = sharedStorage.IncrementBlobRetryCount(ctx, metadata1)
	assert.Nil(t, err)
	metadata1, err = sharedStorage.GetBlobMetadata(ctx, blobKey)
	assert.Nil(t, err)
	assert.Equal(t, uint(1), metadata1.NumRetries)

	err = sharedStorage.IncrementBlobRetryCount(ctx, metadata1)
	assert.Nil(t, err)
	metadata1, err = sharedStorage.GetBlobMetadata(ctx, blobKey)
	assert.Nil(t, err)
	assert.Equal(t, uint(2), metadata1.NumRetries)

	batchHeaderHash := [32]byte{1, 2, 3}
	blobIndex := uint32(0)
	confirmationInfo := &disperser.ConfirmationInfo{
		BatchHeaderHash:         batchHeaderHash,
		BlobIndex:               blobIndex,
		BlobCount:               2,
		SignatoryRecordHash:     [32]byte{0},
		ReferenceBlockNumber:    132,
		BatchRoot:               []byte("hello"),
		BlobCommitment:          &encoding.BlobCommitments{},
		BatchID:                 99,
		ConfirmationTxnHash:     common.HexToHash("0x123"),
		ConfirmationBlockNumber: 150,
		Fee:                     []byte{0},
	}
	metadata := &disperser.BlobMetadata{
		BlobHash:     blobKey.BlobHash,
		MetadataHash: blobKey.MetadataHash,
		BlobStatus:   disperser.Processing,
		Expiry:       0,
		NumRetries:   0,
		RequestMetadata: &disperser.RequestMetadata{
			BlobRequestHeader: core.BlobRequestHeader{
				SecurityParams: securityParams,
			},
			RequestedAt: requestedAt,
			BlobSize:    blobSize,
		},
	}
	updatedMetadata, err := sharedStorage.MarkBlobConfirmed(ctx, metadata, confirmationInfo)
	assert.Nil(t, err)
	assert.Equal(t, disperser.Confirmed, updatedMetadata.BlobStatus)

	metadata1, err = sharedStorage.GetBlobMetadata(ctx, blobKey)
	assert.Nil(t, err)
	assertMetadata(t, blobKey, blobSize, requestedAt, disperser.Confirmed, metadata1)

	err = sharedStorage.UpdateConfirmationBlockNumber(ctx, metadata1, 151)
	assert.Nil(t, err)
	metadata1, err = sharedStorage.GetBlobMetadata(ctx, blobKey)
	assert.Nil(t, err)
	assert.Equal(t, uint32(151), metadata1.ConfirmationInfo.ConfirmationBlockNumber)

	err = sharedStorage.MarkBlobFinalized(ctx, blobKey)
	assert.Nil(t, err)
	metadata1, err = sharedStorage.GetBlobMetadata(ctx, blobKey)
	assert.Nil(t, err)
	assert.Equal(t, disperser.Finalized, metadata1.BlobStatus)

	metadata1, err = sharedStorage.GetBlobMetadata(ctx, blobKey)
	assert.Nil(t, err)
	assertMetadata(t, blobKey, blobSize, requestedAt, disperser.Finalized, metadata1)

	allMetadata, err := sharedStorage.GetAllBlobMetadataByBatch(ctx, batchHeaderHash)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(allMetadata))
	assertMetadata(t, blobKey, blobSize, requestedAt, disperser.Finalized, allMetadata[0])

	// Store the second blob and then check the metadata.
	blob.Data = []byte("foo")
	blobSize2 := uint(len(blob.Data))
	blobKey2, err := sharedStorage.StoreBlob(ctx, blob, requestedAt)
	assert.Nil(t, err)
	assert.NotEqual(t, blobKey, blobKey2)
	confirmationInfo = &disperser.ConfirmationInfo{
		BatchHeaderHash:         batchHeaderHash,
		BlobIndex:               uint32(1),
		BlobCount:               2,
		SignatoryRecordHash:     [32]byte{0},
		ReferenceBlockNumber:    132,
		BatchRoot:               []byte("hello"),
		BlobCommitment:          &encoding.BlobCommitments{},
		BatchID:                 99,
		ConfirmationBlockNumber: 150,
		Fee:                     []byte{0},
	}
	metadata = &disperser.BlobMetadata{
		BlobHash:     blobKey2.BlobHash,
		MetadataHash: blobKey2.MetadataHash,
		BlobStatus:   disperser.Processing,
		Expiry:       0,
		NumRetries:   0,
		RequestMetadata: &disperser.RequestMetadata{
			BlobRequestHeader: core.BlobRequestHeader{
				SecurityParams: securityParams,
			},
			RequestedAt: requestedAt,
			BlobSize:    blobSize2,
		},
	}
	updatedMetadata, err = sharedStorage.MarkBlobInsufficientSignatures(ctx, metadata, confirmationInfo)
	assert.Nil(t, err)
	assert.Equal(t, disperser.InsufficientSignatures, updatedMetadata.BlobStatus)

	allMetadata, err = sharedStorage.GetAllBlobMetadataByBatch(ctx, batchHeaderHash)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(allMetadata))
	var blob1Metadata, blob2Metadata *disperser.BlobMetadata
	for i, metadata := range allMetadata {
		if metadata.BlobHash == metadata1.BlobHash {
			blob1Metadata = allMetadata[i]
		} else if metadata.BlobHash == updatedMetadata.BlobHash {
			blob2Metadata = allMetadata[i]
		}
	}
	assert.NotNil(t, blob1Metadata)
	assert.NotNil(t, blob2Metadata)
	assertMetadata(t, blobKey, blobSize, requestedAt, disperser.Finalized, blob1Metadata)
	assertMetadata(t, blobKey2, blobSize2, requestedAt, disperser.InsufficientSignatures, blob2Metadata)

	// Cleanup: Delete test items
	t.Cleanup(func() {
		deleteItems(t, []commondynamodb.Key{
			{
				"MetadataHash": &types.AttributeValueMemberS{Value: blobKey.MetadataHash},
				"BlobHash":     &types.AttributeValueMemberS{Value: blobKey.BlobHash},
			},
			{
				"MetadataHash": &types.AttributeValueMemberS{Value: blobKey2.MetadataHash},
				"BlobHash":     &types.AttributeValueMemberS{Value: blobKey2.BlobHash},
			},
		})
	})
}

func TestSharedBlobStoreBlobMetadataStoreOperationsWithPagination(t *testing.T) {
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

	// Setup: Queue new blob metadata
	err := blobMetadataStore.QueueNewBlobMetadata(ctx, metadata1)
	assert.NoError(t, err)
	err = blobMetadataStore.QueueNewBlobMetadata(ctx, metadata2)
	assert.NoError(t, err)

	// Test: Fetch individual blob metadata
	fetchedMetadata, err := sharedStorage.GetBlobMetadata(ctx, blobKey1)
	assert.NoError(t, err)
	assert.Equal(t, metadata1, fetchedMetadata)
	fetchedMetadata, err = sharedStorage.GetBlobMetadata(ctx, blobKey2)
	assert.NoError(t, err)
	assert.Equal(t, metadata2, fetchedMetadata)

	// Test: Fetch blob metadata by status with pagination
	t.Run("Fetch Processing Blobs", func(t *testing.T) {
		processing, lastEvaluatedKey, err := sharedStorage.GetBlobMetadataByStatusWithPagination(ctx, disperser.Processing, 1, nil)
		assert.NoError(t, err)
		assert.Len(t, processing, 1)
		assert.Equal(t, metadata1, processing[0])
		assert.NotNil(t, lastEvaluatedKey)

		// Fetch next page (should be empty)
		nextProcessing, nextLastEvaluatedKey, err := sharedStorage.GetBlobMetadataByStatusWithPagination(ctx, disperser.Processing, 1, lastEvaluatedKey)
		assert.NoError(t, err)
		assert.Len(t, nextProcessing, 0)
		assert.Nil(t, nextLastEvaluatedKey)
	})

	t.Run("Fetch Finalized Blobs", func(t *testing.T) {
		finalized, lastEvaluatedKey, err := sharedStorage.GetBlobMetadataByStatusWithPagination(ctx, disperser.Finalized, 1, nil)
		assert.NoError(t, err)
		assert.Len(t, finalized, 1)
		assert.Equal(t, metadata2, finalized[0])
		assert.NotNil(t, lastEvaluatedKey)

		// Fetch next page (should be empty)
		nextFinalized, nextLastEvaluatedKey, err := sharedStorage.GetBlobMetadataByStatusWithPagination(ctx, disperser.Finalized, 1, lastEvaluatedKey)
		assert.NoError(t, err)
		assert.Len(t, nextFinalized, 0)
		assert.Nil(t, nextLastEvaluatedKey)
	})

	// Cleanup: Delete test items
	t.Cleanup(func() {
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
	})
}

func TestSharedBlobStoreGetAllBlobMetadataByBatchWithPagination(t *testing.T) {
	ctx := context.Background()
	batchHeaderHash := [32]byte{1, 2, 3}

	// Create and store multiple blob metadata for the same batch
	numBlobs := 5
	blobKeys := make([]disperser.BlobKey, numBlobs)
	for i := 0; i < numBlobs; i++ {
		blobKey := disperser.BlobKey{
			BlobHash:     fmt.Sprintf("blob%d", i),
			MetadataHash: fmt.Sprintf("hash%d", i),
		}
		blobKeys[i] = blobKey

		metadata := &disperser.BlobMetadata{
			BlobHash:     blobKey.BlobHash,
			MetadataHash: blobKey.MetadataHash,
			BlobStatus:   disperser.Confirmed,
			RequestMetadata: &disperser.RequestMetadata{
				BlobRequestHeader: blob.RequestHeader,
				BlobSize:          blobSize,
				RequestedAt:       uint64(time.Now().UnixNano()),
			},
			ConfirmationInfo: &disperser.ConfirmationInfo{
				BatchHeaderHash: batchHeaderHash,
				BlobIndex:       uint32(i),
			},
		}

		err := blobMetadataStore.QueueNewBlobMetadata(ctx, metadata)
		assert.NoError(t, err)
	}

	// Test pagination with a page size of 2
	t.Run("Fetch All Blobs with Pagination", func(t *testing.T) {
		var allFetchedMetadata []*disperser.BlobMetadata
		var lastEvaluatedKey *disperser.BatchIndexExclusiveStartKey
		pageSize := int32(2)

		for {
			fetchedMetadata, newLastEvaluatedKey, err := sharedStorage.GetAllBlobMetadataByBatchWithPagination(ctx, batchHeaderHash, pageSize, lastEvaluatedKey)
			assert.NoError(t, err)

			allFetchedMetadata = append(allFetchedMetadata, fetchedMetadata...)

			if newLastEvaluatedKey == nil {
				assert.Len(t, fetchedMetadata, numBlobs%int(pageSize))
				break
			} else {
				assert.Len(t, fetchedMetadata, int(pageSize))
			}
			lastEvaluatedKey = newLastEvaluatedKey
		}

		assert.Len(t, allFetchedMetadata, numBlobs)

		// Verify that all blob metadata is fetched and in the correct order
		for i, metadata := range allFetchedMetadata {
			assert.Equal(t, fmt.Sprintf("blob%d", i), metadata.BlobHash)
			assert.Equal(t, fmt.Sprintf("hash%d", i), metadata.MetadataHash)
			assert.Equal(t, uint32(i), metadata.ConfirmationInfo.BlobIndex)
		}
	})

	// Test pagination with a page size of 10
	t.Run("Fetch All Blobs with Pagination (Page Size > Num Blobs)", func(t *testing.T) {
		var allFetchedMetadata []*disperser.BlobMetadata
		var lastEvaluatedKey *disperser.BatchIndexExclusiveStartKey
		pageSize := int32(10)

		for {
			fetchedMetadata, newLastEvaluatedKey, err := sharedStorage.GetAllBlobMetadataByBatchWithPagination(ctx, batchHeaderHash, pageSize, lastEvaluatedKey)
			assert.NoError(t, err)

			allFetchedMetadata = append(allFetchedMetadata, fetchedMetadata...)

			if newLastEvaluatedKey == nil {
				assert.Len(t, fetchedMetadata, numBlobs)
				break
			} else {
				assert.Len(t, fetchedMetadata, int(pageSize))
			}

			lastEvaluatedKey = newLastEvaluatedKey
		}

		assert.Len(t, allFetchedMetadata, numBlobs)

		// Verify that all blob metadata is fetched and in the correct order
		for i, metadata := range allFetchedMetadata {
			assert.Equal(t, fmt.Sprintf("blob%d", i), metadata.BlobHash)
			assert.Equal(t, fmt.Sprintf("hash%d", i), metadata.MetadataHash)
			assert.Equal(t, uint32(i), metadata.ConfirmationInfo.BlobIndex)
		}
	})

	// Test invalid batch header hash
	t.Run("Fetch All Blobs with Invalid Batch Header Hash", func(t *testing.T) {
		invalidBatchHeaderHash := [32]byte{4, 5, 6}
		allFetchedMetadata, lastEvaluatedKey, err := sharedStorage.GetAllBlobMetadataByBatchWithPagination(ctx, invalidBatchHeaderHash, 10, nil)
		assert.NoError(t, err)
		assert.Len(t, allFetchedMetadata, 0)
		assert.Nil(t, lastEvaluatedKey)
	})

	// Cleanup: Delete test items
	t.Cleanup(func() {
		var keys []commondynamodb.Key
		for _, blobKey := range blobKeys {
			keys = append(keys, commondynamodb.Key{
				"MetadataHash": &types.AttributeValueMemberS{Value: blobKey.MetadataHash},
				"BlobHash":     &types.AttributeValueMemberS{Value: blobKey.BlobHash},
			})
		}
		deleteItems(t, keys)
	})
}

func assertMetadata(t *testing.T, blobKey disperser.BlobKey, expectedBlobSize uint, expectedRequestedAt uint64, expectedStatus disperser.BlobStatus, actualMetadata *disperser.BlobMetadata) {
	assert.NotNil(t, actualMetadata)
	assert.Equal(t, expectedStatus, actualMetadata.BlobStatus)
	assert.Equal(t, blob.RequestHeader, actualMetadata.RequestMetadata.BlobRequestHeader)
	assert.Equal(t, blobKey.BlobHash, actualMetadata.BlobHash)
	assert.Equal(t, blobKey.MetadataHash, actualMetadata.MetadataHash)
	assert.Equal(t, expectedBlobSize, actualMetadata.RequestMetadata.BlobSize)
	assert.Equal(t, expectedRequestedAt, actualMetadata.RequestMetadata.RequestedAt)
	metadataSuffix, err := metadataSuffix(actualMetadata.RequestMetadata.RequestedAt, actualMetadata.RequestMetadata.SecurityParams)
	assert.Nil(t, err)
	assert.Equal(t, metadataSuffix, actualMetadata.MetadataHash)
}

func assertBlob(t *testing.T, blob *core.Blob) {
	assert.NotNil(t, blob)
	assert.Equal(t, blob.Data, blob.Data)
	assert.Equal(t, blob.RequestHeader.SecurityParams, blob.RequestHeader.SecurityParams)
}

func metadataSuffix(requestedAt uint64, securityParams []*core.SecurityParam) (string, error) {
	var str string
	str = fmt.Sprintf("%d/", requestedAt)
	for _, param := range securityParams {
		appendStr := fmt.Sprintf("%d/%d/", param.QuorumID, param.AdversaryThreshold)
		// Append String incase of multiple securityParams
		str = str + appendStr
	}
	bytes := []byte(str)
	return hex.EncodeToString(sha256.New().Sum(bytes)), nil
}
