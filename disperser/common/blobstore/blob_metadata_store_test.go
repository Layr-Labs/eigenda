package blobstore_test

import (
	"context"
	"testing"
	"time"

	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
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

	confirmedMetadata := getConfirmedMetadata(t, blobKey1)
	err = blobMetadataStore.UpdateBlobMetadata(ctx, blobKey1, confirmedMetadata)
	assert.NoError(t, err)

	metadata, err := blobMetadataStore.GetBlobMetadataInBatch(ctx, confirmedMetadata.ConfirmationInfo.BatchHeaderHash, confirmedMetadata.ConfirmationInfo.BlobIndex)
	assert.NoError(t, err)
	assert.Equal(t, metadata, confirmedMetadata)

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
	assert.Nil(t, lastEvaluatedKey)

	finalized, lastEvaluatedKey, err := blobMetadataStore.GetBlobMetadataByStatusWithPagination(ctx, disperser.Finalized, 1, nil)
	assert.NoError(t, err)
	assert.Len(t, finalized, 1)
	assert.Equal(t, metadata2, finalized[0])
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
	commitment := &core.Commitment{
		G1Point: &bn254.G1Point{
			X: commitX,
			Y: commitY,
		},
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
			BlobCommitment: &core.BlobCommitments{
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
