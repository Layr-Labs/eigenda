package relay

import (
	"testing"
	"time"

	tu "github.com/Layr-Labs/eigenda/common/testutils"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/stretchr/testify/require"
)

func TestReadWrite(t *testing.T) {
	ctx := t.Context()
	tu.InitializeRandom()

	setup(t)
	defer teardown(t)

	blobStore := buildBlobStore(t, logger)

	expectedData := make(map[v2.BlobKey][]byte)

	blobCount := 10
	for i := 0; i < blobCount; i++ {
		header, data := randomBlob(t)

		blobKey, err := header.BlobKey()
		require.NoError(t, err)
		expectedData[blobKey] = data

		err = blobStore.StoreBlob(ctx, blobKey, data)
		require.NoError(t, err)
	}

	server, err := newBlobProvider(
		ctx,
		logger,
		blobStore,
		1024*1024*32,
		32,
		10*time.Second,
		nil)
	require.NoError(t, err)

	// Read the blobs back.
	for key, data := range expectedData {
		blob, err := server.GetBlob(ctx, key)

		require.NoError(t, err)
		require.Equal(t, data, blob)
	}

	// Read the blobs back again to test caching.
	for key, data := range expectedData {
		blob, err := server.GetBlob(ctx, key)

		require.NoError(t, err)
		require.Equal(t, data, blob)
	}
}

func TestNonExistentBlob(t *testing.T) {
	ctx := t.Context()
	tu.InitializeRandom()

	setup(t)
	defer teardown(t)

	blobStore := buildBlobStore(t, logger)

	server, err := newBlobProvider(
		ctx,
		logger,
		blobStore,
		1024*1024*32,
		32,
		10*time.Second,
		nil)
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		blob, err := server.GetBlob(ctx, v2.BlobKey(tu.RandomBytes(32)))
		require.Error(t, err)
		require.Nil(t, blob)
	}
}
