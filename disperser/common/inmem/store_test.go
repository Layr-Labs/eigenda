package inmem_test

import (
	"context"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/common/inmem"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestBlobStore(t *testing.T) {
	bs := inmem.NewBlobStore()
	numBlobs := 10
	requestedAt := uint64(time.Now().UnixNano())
	securityParams := []*core.SecurityParam{}

	ctx := context.Background()
	keys := make([]disperser.BlobKey, numBlobs)
	for i := 0; i < numBlobs; i++ {
		blobKey, err := bs.StoreBlob(ctx, &core.Blob{
			RequestHeader: core.BlobRequestHeader{
				SecurityParams: []*core.SecurityParam{},
			},
			Data: []byte{byte(i)},
		}, requestedAt)
		assert.Nil(t, err)
		keys[i] = blobKey
	}

	metas, err := bs.GetBlobMetadataByStatus(ctx, disperser.Processing)
	assert.Nil(t, err)
	assert.Len(t, metas, numBlobs)

	data, err := bs.GetBlobContent(ctx, keys[1].BlobHash)
	assert.Nil(t, err)
	assert.Equal(t, data, []byte{byte(1)})

	metadatas, err := bs.GetBlobMetadataByStatus(ctx, disperser.Processing)
	assert.Nil(t, err)
	assert.Len(t, metadatas, numBlobs)

	blobs, err := bs.GetBlobsByMetadata(ctx, []*disperser.BlobMetadata{metadatas[2], metadatas[5]})
	assert.Nil(t, err)
	assert.Len(t, blobs, 2)
	blobKey1 := metadatas[2].GetBlobKey()
	blobKey2 := metadatas[5].GetBlobKey()
	assert.Len(t, blobs[blobKey1].Data, 1)
	assert.Len(t, blobs[blobKey2].Data, 1)

	meta1, err := bs.GetBlobMetadata(ctx, blobKey1)
	assert.Nil(t, err)
	assert.Equal(t, meta1.BlobStatus, disperser.Processing)
	meta2, err := bs.GetBlobMetadata(ctx, blobKey2)
	assert.Nil(t, err)
	assert.Equal(t, meta2.BlobStatus, disperser.Processing)

	batchHeaderHash := [32]byte{1, 2, 3}
	blobIndex := uint32(0)
	sigRecordHash := [32]byte{0}
	inclusionProof := []byte{1, 2, 3, 4, 5}

	confirmationInfo := &disperser.ConfirmationInfo{
		BatchHeaderHash:         batchHeaderHash,
		BlobIndex:               blobIndex,
		BlobCount:               uint32(numBlobs),
		SignatoryRecordHash:     sigRecordHash,
		ReferenceBlockNumber:    132,
		BatchRoot:               []byte("hello"),
		BlobInclusionProof:      inclusionProof,
		BlobCommitment:          &encoding.BlobCommitments{},
		BatchID:                 99,
		ConfirmationTxnHash:     common.HexToHash("0x123"),
		ConfirmationBlockNumber: uint32(150),
		Fee:                     []byte{0},
	}
	metadata := &disperser.BlobMetadata{
		BlobHash:     meta2.BlobHash,
		MetadataHash: meta2.MetadataHash,
		BlobStatus:   disperser.Processing,
		Expiry:       0,
		NumRetries:   0,
		RequestMetadata: &disperser.RequestMetadata{
			BlobRequestHeader: core.BlobRequestHeader{
				SecurityParams: securityParams,
			},
			RequestedAt: requestedAt,
			BlobSize:    1,
		},
	}
	updated, err := bs.MarkBlobConfirmed(ctx, metadata, confirmationInfo)
	assert.Nil(t, err)
	assert.Equal(t, disperser.Confirmed, updated.BlobStatus)

	meta2, err = bs.GetBlobMetadata(ctx, blobKey2)
	assert.Nil(t, err)
	assert.Equal(t, meta2.BlobStatus, disperser.Confirmed)
	meta1, err = bs.GetBlobMetadata(ctx, blobKey1)
	assert.Nil(t, err)
	assert.Equal(t, meta1.BlobStatus, disperser.Processing)

	err = bs.MarkBlobFailed(ctx, blobKey1)
	assert.Nil(t, err)

	meta1, err = bs.GetBlobMetadata(ctx, blobKey1)
	assert.Nil(t, err)
	assert.Equal(t, meta1.BlobStatus, disperser.Failed)

	allMeta, err := bs.GetAllBlobMetadataByBatch(ctx, batchHeaderHash)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(allMeta))
	assert.Equal(t, allMeta[0].BlobStatus, disperser.Confirmed)
}
