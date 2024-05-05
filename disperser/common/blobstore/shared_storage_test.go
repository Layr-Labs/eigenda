package blobstore_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/encoding"
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
	fmt.Println("Num Retries", metadata1.NumRetries)
	assert.Nil(t, err)
	assert.Equal(t, uint(1), metadata1.NumRetries)

	err = sharedStorage.IncrementBlobRetryCount(ctx, metadata1)
	assert.Nil(t, err)
	metadata1, err = sharedStorage.GetBlobMetadata(ctx, blobKey)
	fmt.Println("Num Retries", metadata1.NumRetries)
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
		AccountID:    "test",
		RequestMetadata: &disperser.RequestMetadata{
			BlobRequestHeader: core.BlobRequestHeader{
				SecurityParams: securityParams,
				BlobAuthHeader: blob.RequestHeader.BlobAuthHeader,
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

	err = sharedStorage.MarkBlobFinalized(ctx, blobKey)
	assert.Nil(t, err)

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
		AccountID:    "test",
		RequestMetadata: &disperser.RequestMetadata{
			BlobRequestHeader: core.BlobRequestHeader{
				SecurityParams: securityParams,
				BlobAuthHeader: blob.RequestHeader.BlobAuthHeader,
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
