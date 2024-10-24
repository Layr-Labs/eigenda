package blobstore_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	core "github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

func TestBlobMetadataStoreOperations(t *testing.T) {
	ctx := context.Background()
	blobHeader1 := &corev2.BlobHeader{
		BlobVersion:     0,
		QuorumNumbers:   []core.QuorumID{0},
		BlobCommitments: mockCommitment,
		PaymentMetadata: core.PaymentMetadata{
			AccountID:         "0x123",
			BinIndex:          0,
			CumulativePayment: big.NewInt(532),
		},
	}
	blobKey1, err := blobHeader1.BlobKey()
	assert.NoError(t, err)
	blobHeader2 := &corev2.BlobHeader{
		BlobVersion:     0,
		QuorumNumbers:   []core.QuorumID{1},
		BlobCommitments: mockCommitment,
		PaymentMetadata: core.PaymentMetadata{
			AccountID:         "0x456",
			BinIndex:          2,
			CumulativePayment: big.NewInt(999),
		},
	}
	blobKey2, err := blobHeader2.BlobKey()
	assert.NoError(t, err)

	now := time.Now()
	metadata1 := &v2.BlobMetadata{
		BlobHeader: *blobHeader1,
		BlobStatus: v2.Queued,
		Expiry:     uint64(now.Add(time.Hour).Unix()),
		NumRetries: 0,
		UpdatedAt:  uint64(now.UnixNano()),
	}
	metadata2 := &v2.BlobMetadata{
		BlobHeader: *blobHeader2,
		BlobStatus: v2.Certified,
		Expiry:     uint64(now.Add(time.Hour).Unix()),
		NumRetries: 0,
		UpdatedAt:  uint64(now.UnixNano()),
	}
	err = blobMetadataStore.PutBlobMetadata(ctx, metadata1)
	assert.NoError(t, err)
	err = blobMetadataStore.PutBlobMetadata(ctx, metadata2)
	assert.NoError(t, err)

	fetchedMetadata, err := blobMetadataStore.GetBlobMetadata(ctx, blobKey1)
	assert.NoError(t, err)
	assert.Equal(t, metadata1, fetchedMetadata)
	fetchedMetadata, err = blobMetadataStore.GetBlobMetadata(ctx, blobKey2)
	assert.NoError(t, err)
	assert.Equal(t, metadata2, fetchedMetadata)

	queued, err := blobMetadataStore.GetBlobMetadataByStatus(ctx, v2.Queued, 0)
	assert.NoError(t, err)
	assert.Len(t, queued, 1)
	assert.Equal(t, metadata1, queued[0])
	certified, err := blobMetadataStore.GetBlobMetadataByStatus(ctx, v2.Certified, 0)
	assert.NoError(t, err)
	assert.Len(t, certified, 1)
	assert.Equal(t, metadata2, certified[0])

	queuedCount, err := blobMetadataStore.GetBlobMetadataCountByStatus(ctx, v2.Queued)
	assert.NoError(t, err)
	assert.Equal(t, int32(1), queuedCount)

	deleteItems(t, []commondynamodb.Key{
		{
			"PK": &types.AttributeValueMemberS{Value: "BlobKey#" + blobKey1.Hex()},
			"SK": &types.AttributeValueMemberS{Value: "BlobMetadata"},
		},
		{
			"PK": &types.AttributeValueMemberS{Value: "BlobKey#" + blobKey2.Hex()},
			"SK": &types.AttributeValueMemberS{Value: "BlobMetadata"},
		},
	})
}

func deleteItems(t *testing.T, keys []commondynamodb.Key) {
	failed, err := dynamoClient.DeleteItems(context.Background(), metadataTableName, keys)
	assert.NoError(t, err)
	assert.Len(t, failed, 0)
}
