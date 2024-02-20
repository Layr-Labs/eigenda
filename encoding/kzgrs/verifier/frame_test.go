package kzgrs_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Layr-Labs/eigenda/encoding/kzgrs"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	kzg "github.com/Layr-Labs/eigenda/pkg/kzg"
)

func TestEncodeDecodeFrame_AreInverses(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	group, _ := kzgrs.NewKzgEncoderGroup(kzgConfig, true)

	params := rs.GetEncodingParams(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))

	enc, err := group.NewKzgEncoder(params)

	require.Nil(t, err)
	require.NotNil(t, enc)

	_, _, _, frames, _, err := enc.EncodeBytes(GETTYSBURG_ADDRESS_BYTES)
	require.Nil(t, err)
	require.NotNil(t, frames, err)

	b, err := frames[0].Encode()
	require.Nil(t, err)
	require.NotNil(t, b)

	frame, err := kzgrs.Decode(b)
	require.Nil(t, err)
	require.NotNil(t, frame)

	assert.Equal(t, frame, frames[0])
}

func TestVerify(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	group, _ := kzgrs.NewKzgEncoderGroup(kzgConfig, true)

	params := rs.GetEncodingParams(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))

	enc, err := group.NewKzgEncoder(params)
	require.Nil(t, err)
	require.NotNil(t, enc)

	commit, _, _, frames, _, err := enc.EncodeBytes(GETTYSBURG_ADDRESS_BYTES)
	require.Nil(t, err)
	require.NotNil(t, commit)
	require.NotNil(t, frames)

	n := uint8(math.Log2(float64(params.NumEvaluations())))
	fs := kzg.NewFFTSettings(n)
	require.NotNil(t, fs)

	lc := enc.Fs.ExpandedRootsOfUnity[uint64(0)]
	require.NotNil(t, lc)

	g2Atn, err := kzgrs.ReadG2Point(uint64(len(frames[0].Coeffs)), kzgConfig)
	require.Nil(t, err)
	assert.True(t, frames[0].Verify(enc.Ks, commit, &lc, &g2Atn))
}
