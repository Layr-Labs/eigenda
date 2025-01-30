package rs_test

import (
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/stretchr/testify/require"
)

// TODO find and replace "frame" terminology

func TestFrameCoeffsSerialization(t *testing.T) {
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

	serializedSize := rs.CoeffsSize(coeffs[0])
	bytes := make([]byte, serializedSize)
	rs.SerializeFrameCoeffs(coeffs[0], bytes)

	deserializedCoeffs, bytesRead, err := rs.DeserializeFrameCoeffs(bytes)
	require.NoError(t, err)
	require.Equal(t, bytesRead, serializedSize)
	require.Equal(t, coeffs[0], deserializedCoeffs)
}

func TestFrameCoeffsSliceSerialization(t *testing.T) {
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
	require.NoError(t, err)

	encodedCoeffs, err := rs.SerializeFrameCoeffsSlice(coeffs)
	require.NoError(t, err)

	decodedCoeffs, err := rs.DeserializeFrameCoeffsSlice(encodedCoeffs)
	require.NoError(t, err)

	require.Equal(t, len(coeffs), len(decodedCoeffs))
	for i := range coeffs {
		require.Equal(t, coeffs[i], decodedCoeffs[i])
	}
}

func TestSplitSerializedFrameCoeffs(t *testing.T) {
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
	require.NoError(t, err)

	encodedCoeffs, err := rs.SerializeFrameCoeffsSlice(coeffs)
	require.NoError(t, err)

	splitCoeffBytes, err := rs.SplitSerializedFrameCoeffs(encodedCoeffs)
	require.NoError(t, err)

	// The length of the split coeffs should be equal to the length of the serialized coeffs minus 4 (the frame count)
	totalLength := 0
	for _, coeffBytes := range splitCoeffBytes {
		totalLength += len(coeffBytes)
	}
	require.Equal(t, len(encodedCoeffs)-4, totalLength)

	// deserializing each FrameCoeffs individually should yield the same coeffs as the original
	for i, coeffsBytes := range splitCoeffBytes {
		deserializedFromCoeffBytes, length, err := rs.DeserializeFrameCoeffs(coeffsBytes)
		require.NoError(t, err)
		require.Equal(t, uint32(len(coeffsBytes)), length)
		require.Equal(t, coeffs[i], deserializedFromCoeffBytes)
	}

	// recombining the split coeffs should yield the original serialized coeffs
	combinedCoeffs := rs.CombineSerializedFrameCoeffs(splitCoeffBytes)
	require.Equal(t, encodedCoeffs, combinedCoeffs)

	// finally, parse the combined coeffs (for the sake of sanity)
	decodedCoeffs, err := rs.DeserializeFrameCoeffsSlice(combinedCoeffs)
	require.NoError(t, err)
	for i := range coeffs {
		require.Equal(t, coeffs[i], decodedCoeffs[i])
	}
}
