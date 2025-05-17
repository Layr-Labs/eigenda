package blobstore_test

import (
	"context"
	"testing"
	"time"

	corev2 "github.com/Layr-Labs/eigenda/core/v2"

	"github.com/stretchr/testify/assert"
)

func TestStoreGetBlob(t *testing.T) {
	// Create a test-specific environment
	env := setupForTest(t)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Generate a deterministic blob key
	var testBlobKey corev2.BlobKey
	env.rng.Read(testBlobKey[:])
	err := env.blobStore.StoreBlob(ctx, testBlobKey, []byte("testBlobData"))
	assert.NoError(t, err)

	data, err := env.blobStore.GetBlob(ctx, testBlobKey)
	assert.NoError(t, err)
	assert.Equal(t, []byte("testBlobData"), data)
}

func TestGetBlobNotFound(t *testing.T) {
	// Create a test-specific environment
	env := setupForTest(t)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Generate a deterministic blob key
	var testBlobKey corev2.BlobKey
	env.rng.Read(testBlobKey[:])
	data, err := env.blobStore.GetBlob(ctx, testBlobKey)
	assert.Error(t, err)
	assert.Nil(t, data)
}
