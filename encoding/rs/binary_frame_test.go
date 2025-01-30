package rs_test

import (
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParsingBinaryFrame(t *testing.T) {
	rand := random.NewTestRandom(t)
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	payload := rand.Bytes(1024 + rand.Intn(1024))
	paddedPayload := codec.ConvertByPaddingEmptyByte(payload)

	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(paddedPayload)))
	cfg := encoding.DefaultConfig()
	enc, err := rs.NewEncoder(cfg)
	require.Nil(t, err)

	coeffs, _, err := enc.EncodeBytes(paddedPayload, params)
	require.Nil(t, err)
	require.NotNil(t, coeffs, err)

	coeff := coeffs[0]
	coeffBytes := make([]byte, rs.CoeffsSize(coeff))
	rs.SerializeFrameCoeffs(coeff, coeffBytes)

	g1, err := randomG1()
	require.NoError(t, err)
	proof := g1.G1Affine
	proofBytes := make([]byte, rs.SerializedProofLength)
	err = rs.SerializeFrameProof(proof, proofBytes)
	require.NoError(t, err)

	binaryFrame := rs.BuildBinaryFrame(proofBytes, coeffBytes)

	frame, err := rs.DeserializeBinaryFrame(binaryFrame)
	require.NoError(t, err)

	require.True(t, proof.Equal(&frame.Proof))
	require.Equal(t, ([]encoding.Symbol)(coeff), frame.Coeffs)
}

func TestParsingBinaryFrames(t *testing.T) {
	rand := random.NewTestRandom(t)
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	payload := rand.Bytes(1024 + rand.Intn(1024))
	paddedPayload := codec.ConvertByPaddingEmptyByte(payload)

	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(paddedPayload)))
	cfg := encoding.DefaultConfig()
	enc, err := rs.NewEncoder(cfg)
	require.Nil(t, err)

	coeffs, _, err := enc.EncodeBytes(paddedPayload, params)
	require.Nil(t, err)
	require.NotNil(t, coeffs, err)

	proofs := make([]*encoding.Proof, len(coeffs))
	for i := 0; i < len(coeffs); i++ {
		g1, err := randomG1()
		require.NoError(t, err)
		proofs[i] = g1.G1Affine
	}

	binaryFrames := make([][]byte, len(coeffs))

	for i := 0; i < len(coeffs); i++ {
		coeffBytes := make([]byte, rs.CoeffsSize(coeffs[i]))
		rs.SerializeFrameCoeffs(coeffs[i], coeffBytes)

		proofBytes := make([]byte, rs.SerializedProofLength)
		err = rs.SerializeFrameProof(proofs[i], proofBytes)
		require.NoError(t, err)

		binaryFrames[i] = rs.BuildBinaryFrame(proofBytes, coeffBytes)
	}

	combinedBinaryFrames := rs.CombineBinaryFrames(binaryFrames)
	splitBinaryFrames, err := rs.SplitBinaryFrames(combinedBinaryFrames)
	require.NoError(t, err)
	frames, err := rs.DeserializeBinaryFrames(combinedBinaryFrames)
	require.NoError(t, err)

	for i := 0; i < len(coeffs); i++ {
		// sanity check split frame
		binaryFrame := splitBinaryFrames[i]
		require.Equal(t, binaryFrames[i], binaryFrame)

		frame := frames[i]
		require.True(t, proofs[i].Equal(&frame.Proof))
		require.Equal(t, ([]encoding.Symbol)(coeffs[i]), frame.Coeffs)
	}
}
