package controller_test

import (
	"context"
	"testing"

	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/stretchr/testify/require"
)

func deleteBlobs(t *testing.T, blobMetadataStore *blobstore.BlobMetadataStore, keys []corev2.BlobKey, batchHeaderHashes [][32]byte) {
	ctx := context.Background()
	for _, key := range keys {
		err := blobMetadataStore.DeleteBlobMetadata(ctx, key)
		require.NoError(t, err)
		err = blobMetadataStore.DeleteBlobCertificate(ctx, key)
		require.NoError(t, err)
	}

	for _, bhh := range batchHeaderHashes {
		err := blobMetadataStore.DeleteBatchHeader(ctx, bhh)
		require.NoError(t, err)
	}
}
