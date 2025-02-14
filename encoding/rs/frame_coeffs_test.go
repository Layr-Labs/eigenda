package rs_test

import (
	"encoding/binary"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/stretchr/testify/require"
)

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

	elementCount, splitCoeffBytes, err := rs.SplitSerializedFrameCoeffs(encodedCoeffs)
	require.NoError(t, err)
	require.Equal(t, elementCount, uint32(len(coeffs[0])))

	// recombining the split coeffs should yield the original serialized coeffs
	combinedCoeffs := make([]byte, len(encodedCoeffs))
	binary.BigEndian.PutUint32(combinedCoeffs, elementCount)
	for i, splitCoeff := range splitCoeffBytes {
		copy(combinedCoeffs[4+i*len(splitCoeff):], splitCoeff)
	}

	require.Equal(t, encodedCoeffs, combinedCoeffs)
}
