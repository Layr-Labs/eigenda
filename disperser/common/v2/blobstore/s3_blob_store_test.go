package blobstore_test

import (
	"testing"

	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/test/random"
	"github.com/stretchr/testify/assert"
)

func TestStoreGetBlob(t *testing.T) {
	ctx := t.Context()
	testBlobKey := corev2.BlobKey(random.RandomBytes(32))
	err := blobStore.StoreBlob(ctx, testBlobKey, []byte("testBlobData"))
	assert.NoError(t, err)
	data, err := blobStore.GetBlob(ctx, testBlobKey)
	assert.NoError(t, err)
	assert.Equal(t, []byte("testBlobData"), data)
}

func TestGetBlobNotFound(t *testing.T) {
	ctx := t.Context()
	testBlobKey := corev2.BlobKey(random.RandomBytes(32))
	data, err := blobStore.GetBlob(ctx, testBlobKey)
	assert.Error(t, err)
	assert.Nil(t, data)
}
