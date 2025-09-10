package relay

import (
	"testing"
	"time"

	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/stretchr/testify/require"
)

func deserializeBinaryFrames(t *testing.T, binaryFrames *core.ChunksData) []*encoding.Frame {
	t.Helper()
	bundleBytes, err := binaryFrames.FlattenToBundle()
	require.NoError(t, err)
	bundle := core.Bundle{}
	bundle, err = bundle.Deserialize(bundleBytes)
	require.NoError(t, err)
	return bundle
}

func TestFetchingIndividualBlobs(t *testing.T) {
	ctx := t.Context()
	tu.InitializeRandom()

	setup(t)
	defer teardown(t)

	chunkReader, chunkWriter := buildChunkStore(t, logger)

	expectedFrames := make(map[v2.BlobKey][]*encoding.Frame)
	fragmentInfoMap := make(map[v2.BlobKey]*encoding.FragmentInfo)

	// Write some data.
	blobCount := 10
	for i := 0; i < blobCount; i++ {

		header, _, frames := randomBlobChunks(t)
		blobKey, err := header.BlobKey()
		require.NoError(t, err)

		rsFrames, proofs := disassembleFrames(t, frames)

		err = chunkWriter.PutFrameProofs(ctx, blobKey, proofs)
		require.NoError(t, err)

		fragmentInfo, err := chunkWriter.PutFrameCoefficients(ctx, blobKey, rsFrames)
		require.NoError(t, err)

		expectedFrames[blobKey] = frames
		fragmentInfoMap[blobKey] = fragmentInfo
	}

	server, err := newChunkProvider(
		ctx,
		logger,
		chunkReader,
		1024*1024*32,
		32,
		10*time.Second,
		10*time.Second,
		nil)
	require.NoError(t, err)

	// Read it back.
	for key, frames := range expectedFrames {

		mMap := make(metadataMap)
		fragmentInfo := fragmentInfoMap[key]
		mMap[key] = &blobMetadata{
			totalChunkSizeBytes: fragmentInfo.TotalChunkSizeBytes,
			fragmentSizeBytes:   fragmentInfo.FragmentSizeBytes,
		}

		fMap, err := server.GetFrames(ctx, mMap)
		require.NoError(t, err)

		require.Equal(t, 1, len(fMap))
		readFrames := (fMap)[key]
		require.NotNil(t, readFrames)

		// TODO: when I inspect this data using a debugger, the proofs are all made up of 0s... something
		//  is wrong with the way the data is generated in the test.
		deserializedFrames := deserializeBinaryFrames(t, readFrames)
		require.Equal(t, frames, deserializedFrames)
	}

	// Read it back again to test caching.
	for key, frames := range expectedFrames {
		mMap := make(metadataMap)
		fragmentInfo := fragmentInfoMap[key]
		mMap[key] = &blobMetadata{
			totalChunkSizeBytes: fragmentInfo.TotalChunkSizeBytes,
			fragmentSizeBytes:   fragmentInfo.FragmentSizeBytes,
		}

		fMap, err := server.GetFrames(ctx, mMap)
		require.NoError(t, err)

		require.Equal(t, 1, len(fMap))
		readFrames := (fMap)[key]
		require.NotNil(t, readFrames)

		deserializedFrames := deserializeBinaryFrames(t, readFrames)
		require.Equal(t, frames, deserializedFrames)
	}
}

func TestFetchingBatchedBlobs(t *testing.T) {
	ctx := t.Context()
	tu.InitializeRandom()

	setup(t)
	defer teardown(t)

	chunkReader, chunkWriter := buildChunkStore(t, logger)

	expectedFrames := make(map[v2.BlobKey][]*encoding.Frame)
	fragmentInfoMap := make(map[v2.BlobKey]*encoding.FragmentInfo)

	// Write some data.
	blobCount := 10
	for i := 0; i < blobCount; i++ {

		header, _, frames := randomBlobChunks(t)
		blobKey, err := header.BlobKey()
		require.NoError(t, err)

		rsFrames, proofs := disassembleFrames(t, frames)

		err = chunkWriter.PutFrameProofs(ctx, blobKey, proofs)
		require.NoError(t, err)

		fragmentInfo, err := chunkWriter.PutFrameCoefficients(ctx, blobKey, rsFrames)
		require.NoError(t, err)

		expectedFrames[blobKey] = frames
		fragmentInfoMap[blobKey] = fragmentInfo
	}

	server, err := newChunkProvider(
		ctx,
		logger,
		chunkReader,
		1024*1024*32,
		32,
		10*time.Second,
		10*time.Second,
		nil)
	require.NoError(t, err)

	// Read it back.
	batchSize := 3
	for i := 0; i < 10; i++ {

		mMap := make(metadataMap)
		for key := range expectedFrames {
			mMap[key] = &blobMetadata{
				totalChunkSizeBytes: fragmentInfoMap[key].TotalChunkSizeBytes,
				fragmentSizeBytes:   fragmentInfoMap[key].FragmentSizeBytes,
			}
			if len(mMap) == batchSize {
				break
			}
		}

		fMap, err := server.GetFrames(ctx, mMap)
		require.NoError(t, err)

		require.Equal(t, batchSize, len(fMap))
		for key := range mMap {

			readFrames := (fMap)[key]
			require.NotNil(t, readFrames)

			expectedFramesForBlob := expectedFrames[key]
			deserializedFrames := deserializeBinaryFrames(t, readFrames)
			require.Equal(t, expectedFramesForBlob, deserializedFrames)
		}
	}
}
