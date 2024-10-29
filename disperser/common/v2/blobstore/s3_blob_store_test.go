package blobstore_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStoreGetBlob(t *testing.T) {
	err := blobStore.StoreBlob(context.Background(), "testBlobKey", []byte("testBlobData"))
	assert.NoError(t, err)
	data, err := blobStore.GetBlob(context.Background(), "testBlobKey")
	assert.NoError(t, err)
	assert.Equal(t, []byte("testBlobData"), data)
}

func TestGetBlobNotFound(t *testing.T) {
	data, err := blobStore.GetBlob(context.Background(), "nonExistentBlobKey")
	assert.Error(t, err)
	assert.Nil(t, data)
}
