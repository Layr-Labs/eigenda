package blobstore_test

import (
	"context"
	"testing"

	tu "github.com/Layr-Labs/eigenda/common/testutils"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"

	"github.com/stretchr/testify/assert"
)

func TestStoreGetBlob(t *testing.T) {
	testBlobKey := corev2.BlobKey(tu.RandomBytes(32))
	err := blobStore.StoreBlob(context.Background(), testBlobKey, []byte("testBlobData"))
	assert.NoError(t, err)
	data, err := blobStore.GetBlob(context.Background(), testBlobKey)
	assert.NoError(t, err)
	assert.Equal(t, []byte("testBlobData"), data)
}

func TestGetBlobNotFound(t *testing.T) {
	testBlobKey := corev2.BlobKey(tu.RandomBytes(32))
	data, err := blobStore.GetBlob(context.Background(), testBlobKey)
	assert.Error(t, err)
	assert.Nil(t, data)
}
