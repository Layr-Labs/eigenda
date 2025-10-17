package rs_test

import (
	"encoding/binary"
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/codec"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigenda/test/random"
	"github.com/stretchr/testify/require"
)

func TestFrameCoeffsSliceSerialization(t *testing.T) {
	rand := random.NewTestRandom()

	payload := rand.Bytes(1024 + rand.Intn(1024))
	paddedPayload := codec.ConvertByPaddingEmptyByte(payload)

	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(paddedPayload)))
	cfg := encoding.DefaultConfig()
	enc := rs.NewEncoder(common.TestLogger(t), cfg)

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
	rand := random.NewTestRandom()

	payload := rand.Bytes(1024 + rand.Intn(1024))
	paddedPayload := codec.ConvertByPaddingEmptyByte(payload)

	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(paddedPayload)))
	cfg := encoding.DefaultConfig()
	enc := rs.NewEncoder(common.TestLogger(t), cfg)

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
