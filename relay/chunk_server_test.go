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

	chunkReader, chunkWriter := buildChunkStore(t, logger, nil)

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

	server, err := newChunkServer(context.Background(), logger, chunkReader, 10, 32)
	require.NoError(t, err)

	// Read it back.
	for key, frames := range expectedFrames {

		mMap := make(metadataMap)
		fragmentInfo := fragmentInfoMap[key]
		mMap[key] = &blobMetadata{
			totalChunkSizeBytes: fragmentInfo.TotalChunkSizeBytes,
			fragmentSizeBytes:   fragmentInfo.FragmentSizeBytes,
		}

		fMap, err := server.GetFrames(context.Background(), &mMap)
		require.NoError(t, err)

		require.Equal(t, 1, len(*fMap))
		readFrames := (*fMap)[key]
		require.NotNil(t, readFrames)

		// TODO make sure this equality actually works
		require.Equal(t, frames, readFrames)
	}
}
