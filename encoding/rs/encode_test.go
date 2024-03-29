package rs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
)

func TestEncodeDecode_InvertsWhenSamplingAllFrames(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))

	enc, _ := rs.NewEncoder(params, true)
	require.NotNil(t, enc)

	inputFr := rs.ToPaddedFrArray(GETTYSBURG_ADDRESS_BYTES)
	_, frames, _, err := enc.Encode(inputFr)
	assert.Nil(t, err)

	// sample some frames
	samples, indices := sampleFrames(frames, uint64(len(frames)))
	data, err := enc.DecodeBytes(samples, indices, uint64(len(GETTYSBURG_ADDRESS_BYTES)))

	require.Nil(t, err)
	require.NotNil(t, data)

	assert.Equal(t, data, GETTYSBURG_ADDRESS_BYTES)
}

func TestEncodeDecode_InvertsWhenSamplingMissingFrame(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))
	enc, _ := rs.NewEncoder(params, true)
	require.NotNil(t, enc)

	inputFr := rs.ToPaddedFrArray(GETTYSBURG_ADDRESS_BYTES)
	_, frames, _, err := enc.Encode(inputFr)
	assert.Nil(t, err)

	// sample some frames
	samples, indices := sampleFrames(frames, uint64(numSys))

	data, err := enc.DecodeBytes(samples, indices, uint64(len(GETTYSBURG_ADDRESS_BYTES)))

	require.Nil(t, err)
	require.NotNil(t, data)

	assert.Equal(t, data, GETTYSBURG_ADDRESS_BYTES)
}

func TestEncodeDecode_ErrorsWhenNotEnoughSampledFrames(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))
	enc, _ := rs.NewEncoder(params, true)
	require.NotNil(t, enc)
	inputFr := rs.ToPaddedFrArray(GETTYSBURG_ADDRESS_BYTES)
	_, frames, _, err := enc.Encode(inputFr)
	assert.Nil(t, err)

	// sample some frames
	samples, indices := sampleFrames(frames, uint64(numSys-1))
	data, err := enc.DecodeBytes(samples, indices, uint64(len(GETTYSBURG_ADDRESS_BYTES)))
	require.Nil(t, data)
	require.NotNil(t, err)

	assert.EqualError(t, err, "number of frame must be sufficient")
}

func TestEncodeDecodeAsEval_InvertsWhenSamplingAllFrames(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))

	enc, _ := rs.NewEncoder(params, true)
	require.NotNil(t, enc)

	inputFr := rs.ToPaddedFrArray(GETTYSBURG_ADDRESS_BYTES)

	_, frames, _, err := enc.EncodeAsEval(inputFr)
	assert.Nil(t, err)

	// sample some frames
	//samples, indices := sampleFrames(frames, uint64(len(frames)))
	//data, err := enc.DecodeBytesAsEval(samples, indices, uint64(len(GETTYSBURG_ADDRESS_BYTES)))
	indices := make([]uint64, len(frames))
	for i := 0; i < len(indices); i++ {
		indices[i] = uint64(i)
	}
	data, err := enc.DecodeBytesAsEval(frames, indices, uint64(len(GETTYSBURG_ADDRESS_BYTES)))

	require.Nil(t, err)
	require.NotNil(t, data)

	assert.Equal(t, data, GETTYSBURG_ADDRESS_BYTES)
}
