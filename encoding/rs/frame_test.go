package rs_test

import (
	"fmt"
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/stretchr/testify/require"
)

func TestEncodeDecodeFrame_AreInverses(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))
	cfg := encoding.DefaultConfig()
	enc, err := rs.NewEncoder(cfg)
	require.Nil(t, err)

	frames, _, err := enc.EncodeBytes(GETTYSBURG_ADDRESS_BYTES, params)
	require.Nil(t, err)
	require.NotNil(t, frames, err)

	b, err := frames[0].Encode()
	require.Nil(t, err)
	require.NotNil(t, b)

	frame, err := rs.Decode(b)
	require.Nil(t, err)
	require.NotNil(t, frame)

	require.Equal(t, frame, frames[0])
}

func TestGnarkEncodeDecodeFrame_AreInverses(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))
	cfg := encoding.DefaultConfig()
	enc, err := rs.NewEncoder(cfg)
	require.Nil(t, err)

	frames, _, err := enc.EncodeBytes(GETTYSBURG_ADDRESS_BYTES, params)
	require.Nil(t, err)
	require.NotNil(t, frames, err)

	serializedSize := rs.GnarkFrameSize(&frames[0]) + 4
	bytes := make([]byte, serializedSize)
	rs.GnarkEncodeFrame(&frames[0], bytes)

	fmt.Printf("\n\n\n")

	deserializedFrame, bytesRead, err := rs.GnarkDecodeFrame(bytes)
	require.NoError(t, err)
	require.Equal(t, bytesRead, serializedSize)
	require.Equal(t, &frames[0], deserializedFrame)
}

func TestGnarkEncodeDecodeFrames_AreInverses(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))
	cfg := encoding.DefaultConfig()
	enc, err := rs.NewEncoder(cfg)
	require.Nil(t, err)

	frames, _, err := enc.EncodeBytes(GETTYSBURG_ADDRESS_BYTES, params)
	require.NoError(t, err)

	framesPointers := make([]*rs.Frame, len(frames))
	for i, frame := range frames {
		framesPointers[i] = &frame
	}

	encodedFrames, err := rs.GnarkEncodeFrames(framesPointers)
	require.NoError(t, err)

	decodedFrames, err := rs.GnarkDecodeFrames(encodedFrames)
	require.NoError(t, err)

	require.Equal(t, len(framesPointers), len(decodedFrames))
	for i := range framesPointers {
		require.Equal(t, *framesPointers[i], *decodedFrames[i])
	}
}

func TestGnarkSplitBinaryFrames(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))
	cfg := encoding.DefaultConfig()
	enc, err := rs.NewEncoder(cfg)
	require.Nil(t, err)

	frames, _, err := enc.EncodeBytes(GETTYSBURG_ADDRESS_BYTES, params)
	require.NoError(t, err)

	framesPointers := make([]*rs.Frame, len(frames))
	for i, frame := range frames {
		framesPointers[i] = &frame
	}

	encodedFrames, err := rs.GnarkEncodeFrames(framesPointers)
	require.NoError(t, err)

	splitFrameBytes, err := rs.GnarkSplitBinaryFrames(encodedFrames)
	require.NoError(t, err)

	// The length of the split frames should be equal to the length of the serialized frames minus 4 (the frame count)
	totalLength := 0
	for _, frameBytes := range splitFrameBytes {
		totalLength += len(frameBytes)
	}
	require.Equal(t, len(encodedFrames)-4, totalLength)

	// deserializing each frame individually should yield the same frame as the original
	for i, frameBytes := range splitFrameBytes {
		deserializedFromFrameBytes, length, err := rs.GnarkDecodeFrame(frameBytes)
		require.NoError(t, err)
		require.Equal(t, uint32(len(frameBytes)), length)
		require.Equal(t, *framesPointers[i], *deserializedFromFrameBytes)
	}

	// recombining the split frames should yield the original serialized frames
	combinedFrames := rs.CombineBinaryFrames(splitFrameBytes)
	require.Equal(t, encodedFrames, combinedFrames)

	// finally, parse the combined frames (for the sake of sanity)
	decodedFrames, err := rs.GnarkDecodeFrames(combinedFrames)
	require.NoError(t, err)
	for i := range framesPointers {
		require.Equal(t, *framesPointers[i], *decodedFrames[i])
	}
}
