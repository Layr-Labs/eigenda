package rs_test

import (
	"math"
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	rs_cpu "github.com/Layr-Labs/eigenda/encoding/rs/cpu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeDecodeFrame_AreInverses(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))
	enc, _ := rs.NewEncoder(params, true)

	n := uint8(math.Log2(float64(enc.NumEvaluations())))
	if enc.ChunkLength == 1 {
		n = uint8(math.Log2(float64(2 * enc.NumChunks)))
	}
	fs := fft.NewFFTSettings(n)

	RsComputeDevice := &rs_cpu.RsCpuComputeDevice{
		Fs:             fs,
		EncodingParams: params,
	}

	enc.Computer = RsComputeDevice
	require.NotNil(t, enc)

	frames, _, err := enc.EncodeBytes(GETTYSBURG_ADDRESS_BYTES)
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

func TestEncodeDecodeFrames_AreInverses(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))
	enc, _ := rs.NewEncoder(params, true)

	n := uint8(math.Log2(float64(enc.NumEvaluations())))
	if enc.ChunkLength == 1 {
		n = uint8(math.Log2(float64(2 * enc.NumChunks)))
	}
	fs := fft.NewFFTSettings(n)

	RsComputeDevice := &rs_cpu.RsCpuComputeDevice{
		Fs:             fs,
		EncodingParams: params,
	}

	enc.Computer = RsComputeDevice
	require.NotNil(t, enc)

	frames, _, err := enc.EncodeBytes(GETTYSBURG_ADDRESS_BYTES)
	assert.NoError(t, err)

	framesPointers := make([]*rs.Frame, len(frames))
	for i, frame := range frames {
		framesPointers[i] = &frame
	}

	encodedFrames, err := rs.EncodeFrames(framesPointers)
	assert.NoError(t, err)

	decodedFrames, err := rs.DecodeFrames(encodedFrames)
	assert.NoError(t, err)

	assert.Equal(t, framesPointers, decodedFrames)
}
