package rs_test

import (
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParsingBundle(t *testing.T) {
	rand := random.NewTestRandom(t)
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

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
	binaryFrames, err := rs.BuildChunksData(splitProofs, int(elementCount), splitSerializedCoeffs)
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
