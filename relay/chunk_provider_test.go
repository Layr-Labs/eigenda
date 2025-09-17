package relay

import (
	"fmt"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigenda/test/random"
	"github.com/stretchr/testify/require"

	"github.com/Layr-Labs/eigenda/crypto/ecc/bn254"
	"github.com/Layr-Labs/eigenda/encoding/rs"
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
	random.InitializeRandom()

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
	random.InitializeRandom()

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

func TestParsingBundle(t *testing.T) {
	rand := random.NewTestRandom()
	numNode, numSys := uint64(4), uint64(3)
	numPar := numNode - numSys

	payload := rand.Bytes(1024 + rand.Intn(1024))
	paddedPayload := codec.ConvertByPaddingEmptyByte(payload)

	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(paddedPayload)))
	cfg := encoding.DefaultConfig()
	enc, err := rs.NewEncoder(cfg)
	require.Nil(t, err)

	// Build some random coefficients
	coeffs, _, err := enc.EncodeBytes(paddedPayload, params)
	require.Nil(t, err)
	require.NotNil(t, coeffs, err)
	serializedCoeffs, err := rs.SerializeFrameCoeffsSlice(coeffs)
	require.NoError(t, err)
	elementCount, splitSerializedCoeffs, err := rs.SplitSerializedFrameCoeffs(serializedCoeffs)
	require.NoError(t, err)
	require.Equal(t, uint32(len(coeffs[0])), elementCount)
	require.Equal(t, len(coeffs), len(splitSerializedCoeffs))

	// Build some random proofs
	proofs := make([]*encoding.Proof, len(coeffs))
	for i := 0; i < len(coeffs); i++ {
		g1, err := randomG1()
		require.NoError(t, err)
		proof := g1.G1Affine
		proofs[i] = proof
	}
	serializedProofs, err := rs.SerializeFrameProofs(proofs)
	require.NoError(t, err)
	splitProofs, err := rs.SplitSerializedFrameProofs(serializedProofs)
	require.NoError(t, err)
	require.Equal(t, len(proofs), len(splitProofs))

	// Build binary Frames
	binaryFrames, err := buildChunksData(splitProofs, int(elementCount), splitSerializedCoeffs)
	require.NoError(t, err)

	// convert binary Frames into a serialized bundle
	serializedBundle, err := binaryFrames.FlattenToBundle()
	require.NoError(t, err)

	// construct a standard core.Bundle, serialize it, and compare bytes.
	// Should produce the exact same bytes through the new and old paths.
	bundle := make(core.Bundle, len(proofs))
	for i := 0; i < len(proofs); i++ {
		bundle[i] = &encoding.Frame{
			Proof:  *proofs[i],
			Coeffs: coeffs[i],
		}
	}
	canonicalSerializedBundle, err := bundle.Serialize()
	require.NoError(t, err)
	require.Equal(t, canonicalSerializedBundle, serializedBundle)

	// parse back to proofs and coefficients
	deserializedBundle := core.Bundle{}
	deserializedBundle, err = deserializedBundle.Deserialize(serializedBundle)
	require.NoError(t, err)

	for i := 0; i < len(proofs); i++ {
		expectedProof := proofs[i]
		deserializedProof := &deserializedBundle[i].Proof
		require.True(t, expectedProof.Equal(deserializedProof))

		expectedCoeffs := coeffs[i]
		deserializedCoeffs := (rs.FrameCoeffs)(deserializedBundle[i].Coeffs)
		require.Equal(t, expectedCoeffs, deserializedCoeffs)
	}
}

// randomG1 generates a random G1 point. There is no direct way to generate a random G1 point in the bn254 library,
// but we can generate a random BLS key and steal the public key.
func randomG1() (*bn254.G1Point, error) {
	key, err := bn254.GenRandomBlsKeys()
	if err != nil {
		return nil, fmt.Errorf("failed to generate random BLS keys: %w", err)
	}
	return key.PubKey, nil
}
