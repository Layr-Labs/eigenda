package controller_test

import (
	"testing"
	"time"

	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/disperser/controller"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/stretchr/testify/require"
)

const numObjects = 12

func TestRecoverState(t *testing.T) {
	ctx := t.Context()
	logger := test.GetLogger()
	keys := make([]corev2.BlobKey, numObjects)
	metadatas := make([]*v2.BlobMetadata, numObjects)
	for i := 0; i < numObjects; i++ {
		key, header := newBlob(t, []uint8{0, 1})
		keys[i] = key
		now := time.Now()
		metadatas[i] = &v2.BlobMetadata{
			BlobHeader: header,
			BlobStatus: v2.GatheringSignatures,
			Expiry:     uint64(now.Add(time.Hour).Unix()),
			NumRetries: 0,
			UpdatedAt:  uint64(now.UnixNano()) - uint64(i),
		}
		err := blobMetadataStore.PutBlobMetadata(ctx, metadatas[i])
		require.NoError(t, err)
	}
	err := controller.RecoverState(ctx, blobMetadataStore, logger)
	require.NoError(t, err)

	// check that all blobs are in Failed state
	for i := 0; i < numObjects; i++ {
		metadata, err := blobMetadataStore.GetBlobMetadata(ctx, keys[i])
		require.NoError(t, err)
		require.Equal(t, v2.Failed, metadata.BlobStatus)
	}

	deleteBlobs(t, blobMetadataStore, keys, nil)
}
