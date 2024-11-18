package rs_test

import (
	"fmt"
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeDecodeFrame_AreInverses(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))
	cfg := encoding.DefaultConfig()
	enc, err := rs.NewEncoder(cfg)
	assert.Nil(t, err)

	frames, _, err := enc.EncodeBytes(GETTYSBURG_ADDRESS_BYTES, params)
	require.Nil(t, err)
	require.NotNil(t, frames, err)

	b, err := frames[0].Encode()
	require.Nil(t, err)
	require.NotNil(t, b)

	frame, err := rs.Decode(b)
	require.Nil(t, err)
	require.NotNil(t, frame)

	assert.Equal(t, frame, frames[0])
}

func TestGnarkEncodeDecodeFrame_AreInverses(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))
	cfg := encoding.DefaultConfig()
	enc, err := rs.NewEncoder(cfg)
	assert.Nil(t, err)

	frames, _, err := enc.EncodeBytes(GETTYSBURG_ADDRESS_BYTES, params)
	require.Nil(t, err)
	require.NotNil(t, frames, err)

	serializedSize := rs.GnarkFrameSize(&frames[0]) + 4
	bytes := make([]byte, serializedSize)
	rs.GnarkEncodeFrame(&frames[0], bytes)

	fmt.Printf("\n\n\n")

	deserializedFrame, bytesRead, err := rs.GnarkDecodeFrame(bytes)
	assert.NoError(t, err)
	assert.Equal(t, bytesRead, serializedSize)
	assert.Equal(t, &frames[0], deserializedFrame)
}

func TestGnarkEncodeDecodeFrames_AreInverses(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))
	cfg := encoding.DefaultConfig()
	enc, err := rs.NewEncoder(cfg)
	assert.Nil(t, err)

	frames, _, err := enc.EncodeBytes(GETTYSBURG_ADDRESS_BYTES, params)
	assert.NoError(t, err)

	framesPointers := make([]*rs.Frame, len(frames))
	for i, frame := range frames {
		framesPointers[i] = &frame
	}

	encodedFrames, err := rs.GnarkEncodeFrames(framesPointers)
	assert.NoError(t, err)

	decodedFrames, err := rs.GnarkDecodeFrames(encodedFrames)
	assert.NoError(t, err)

	assert.Equal(t, len(framesPointers), len(decodedFrames))
	for i := range framesPointers {
		assert.Equal(t, *framesPointers[i], *decodedFrames[i])
	}
}
