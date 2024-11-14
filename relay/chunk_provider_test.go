package relay

import (
	"context"
	"github.com/Layr-Labs/eigenda/common"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFetchingIndividualBlobs(t *testing.T) {
	tu.InitializeRandom()

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	setup(t)
	defer teardown()

	chunkReader, chunkWriter := buildChunkStore(t, logger)

	expectedFrames := make(map[v2.BlobKey][]*encoding.Frame)
	fragmentInfoMap := make(map[v2.BlobKey]*encoding.FragmentInfo)

	// Write some data.
	blobCount := 10
	for i := 0; i < blobCount; i++ {

		header, _, frames := randomBlobChunks(t)
		blobKey, err := header.BlobKey()
		require.NoError(t, err)

		rsFrames, proofs := disassembleFrames(frames)

		err = chunkWriter.PutChunkProofs(context.Background(), blobKey, proofs)
		require.NoError(t, err)

		fragmentInfo, err := chunkWriter.PutChunkCoefficients(context.Background(), blobKey, rsFrames)
		require.NoError(t, err)

		expectedFrames[blobKey] = frames
		fragmentInfoMap[blobKey] = fragmentInfo
	}

	server, err := newChunkProvider(context.Background(), logger, chunkReader, 10, 32)
	require.NoError(t, err)

	// Read it back.
	for key, frames := range expectedFrames {

		mMap := make(metadataMap)
		fragmentInfo := fragmentInfoMap[key]
		mMap[key] = &blobMetadata{
			totalChunkSizeBytes: fragmentInfo.TotalChunkSizeBytes,
			fragmentSizeBytes:   fragmentInfo.FragmentSizeBytes,
		}

		fMap, err := server.GetFrames(context.Background(), mMap)
		require.NoError(t, err)

		require.Equal(t, 1, len(fMap))
		readFrames := (fMap)[key]
		require.NotNil(t, readFrames)

		// TODO: when I inspect this data using a debugger, the proofs are all made up of 0s... something
		//  is wrong with the way the data is generated in the test.
		require.Equal(t, frames, readFrames)
	}

	// Read it back again to test caching.
	for key, frames := range expectedFrames {

		mMap := make(metadataMap)
		fragmentInfo := fragmentInfoMap[key]
		mMap[key] = &blobMetadata{
			totalChunkSizeBytes: fragmentInfo.TotalChunkSizeBytes,
			fragmentSizeBytes:   fragmentInfo.FragmentSizeBytes,
		}

		fMap, err := server.GetFrames(context.Background(), mMap)
		require.NoError(t, err)

		require.Equal(t, 1, len(fMap))
		readFrames := (fMap)[key]
		require.NotNil(t, readFrames)

		require.Equal(t, frames, readFrames)
	}
}

func TestFetchingBatchedBlobs(t *testing.T) {
	tu.InitializeRandom()

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	setup(t)
	defer teardown()

	chunkReader, chunkWriter := buildChunkStore(t, logger)

	expectedFrames := make(map[v2.BlobKey][]*encoding.Frame)
	fragmentInfoMap := make(map[v2.BlobKey]*encoding.FragmentInfo)

	// Write some data.
	blobCount := 10
	for i := 0; i < blobCount; i++ {

		header, _, frames := randomBlobChunks(t)
		blobKey, err := header.BlobKey()
		require.NoError(t, err)

		rsFrames, proofs := disassembleFrames(frames)

		err = chunkWriter.PutChunkProofs(context.Background(), blobKey, proofs)
		require.NoError(t, err)

		fragmentInfo, err := chunkWriter.PutChunkCoefficients(context.Background(), blobKey, rsFrames)
		require.NoError(t, err)

		expectedFrames[blobKey] = frames
		fragmentInfoMap[blobKey] = fragmentInfo
	}

	server, err := newChunkProvider(context.Background(), logger, chunkReader, 10, 32)
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

		fMap, err := server.GetFrames(context.Background(), mMap)
		require.NoError(t, err)

		require.Equal(t, batchSize, len(fMap))
		for key := range mMap {

			readFrames := (fMap)[key]
			require.NotNil(t, readFrames)

			expectedFramesForBlob := expectedFrames[key]
			require.Equal(t, expectedFramesForBlob, readFrames)
		}
	}
}
