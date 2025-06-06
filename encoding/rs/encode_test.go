package rs_test

import (
	"fmt"
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeDecode_InvertsWhenSamplingAllFrames(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))

	cfg := encoding.DefaultConfig()
	enc, err := rs.NewEncoder(cfg)
	assert.Nil(t, err)

	inputFr, err := rs.ToFrArray(GETTYSBURG_ADDRESS_BYTES)
	assert.Nil(t, err)
	frames, _, err := enc.Encode(inputFr, params)
	assert.Nil(t, err)

	// sample some Frames
	samples, indices := sampleFrames(frames, uint64(len(frames)))
	data, err := enc.Decode(samples, indices, uint64(len(GETTYSBURG_ADDRESS_BYTES)), params)

	require.Nil(t, err)
	require.NotNil(t, data)

	assert.Equal(t, data, GETTYSBURG_ADDRESS_BYTES)
}

func TestEncodeDecode_InvertsWhenSamplingMissingFrame(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))

	cfg := encoding.DefaultConfig()
	enc, err := rs.NewEncoder(cfg)
	assert.Nil(t, err)

	inputFr, err := rs.ToFrArray(GETTYSBURG_ADDRESS_BYTES)
	assert.Nil(t, err)
	frames, _, err := enc.Encode(inputFr, params)
	assert.Nil(t, err)

	// sample some Frames
	samples, indices := sampleFrames(frames, uint64(len(frames)-1))
	data, err := enc.Decode(samples, indices, uint64(len(GETTYSBURG_ADDRESS_BYTES)), params)

	require.Nil(t, err)
	require.NotNil(t, data)

	assert.Equal(t, data, GETTYSBURG_ADDRESS_BYTES)
}

func TestEncodeDecode_InvertsWithMissingAndDuplicateFrames(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	numSys := uint64(3)
	numPar := uint64(5)
	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))

	cfg := encoding.DefaultConfig()
	enc, err := rs.NewEncoder(cfg)
	assert.Nil(t, err)

	inputFr, err := rs.ToFrArray(GETTYSBURG_ADDRESS_BYTES)
	assert.Nil(t, err)
	frames, _, err := enc.Encode(inputFr, params)
	assert.Nil(t, err)

	assert.EqualValues(t, len(frames), numSys+numPar)

	// sample some Frames
	samples, indices := sampleFrames(frames, uint64(len(frames))-numPar)

	// duplicate two of the frames
	samples = append(samples, samples[0:2]...)
	indices = append(indices, indices[0:2]...)

	data, err := enc.Decode(samples, indices, uint64(len(GETTYSBURG_ADDRESS_BYTES)), params)

	require.Nil(t, err)
	require.NotNil(t, data)

	assert.Equal(t, data, GETTYSBURG_ADDRESS_BYTES)
}

func TestEncodeDecode_ErrorsWhenNotEnoughSampledFrames(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))
	cfg := encoding.DefaultConfig()
	enc, err := rs.NewEncoder(cfg)
	assert.Nil(t, err)

	fmt.Println("Num Chunks: ", params.NumChunks)

	inputFr, err := rs.ToFrArray(GETTYSBURG_ADDRESS_BYTES)
	assert.Nil(t, err)
	frames, _, err := enc.Encode(inputFr, params)
	assert.Nil(t, err)

	// sample some Frames
	samples, indices := sampleFrames(frames, uint64(len(frames)-2))
	data, err := enc.Decode(samples, indices, uint64(len(GETTYSBURG_ADDRESS_BYTES)), params)

	require.Nil(t, data)
	require.NotNil(t, err)

	assert.EqualError(t, err, "number of frame must be sufficient")
}

func TestEncodeDecode_ErrorsWhenNotEnoughSampledFramesWithDuplicates(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))
	cfg := encoding.DefaultConfig()
	enc, err := rs.NewEncoder(cfg)
	assert.Nil(t, err)

	fmt.Println("Num Chunks: ", params.NumChunks)

	inputFr, err := rs.ToFrArray(GETTYSBURG_ADDRESS_BYTES)
	assert.Nil(t, err)
	frames, _, err := enc.Encode(inputFr, params)
	assert.Nil(t, err)

	// sample some Frames
	samples, indices := sampleFrames(frames, uint64(len(frames)-2))

	// duplicate two of the frames
	samples = append(samples, samples[0:2]...)
	indices = append(indices, indices[0:2]...)

	data, err := enc.Decode(samples, indices, uint64(len(GETTYSBURG_ADDRESS_BYTES)), params)

	require.Nil(t, data)
	require.NotNil(t, err)

	assert.EqualError(t, err, "number of frame must be sufficient")
}
