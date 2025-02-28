package rs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
)

func TestGetEncodingParams(t *testing.T) {
	params := encoding.ParamsFromSysPar(1, 4, 1000)

	require.NotNil(t, params)
	assert.Equal(t, params.ChunkLength, uint64(32)) // 1000/32/1 => 32
	// assert.Equal(t, params.DataLen, uint64(1000))
	assert.Equal(t, params.NumChunks, uint64(8))
	assert.Equal(t, params.NumEvaluations(), uint64(256))
}

func TestGetLeadingCoset(t *testing.T) {
	a, err := rs.GetLeadingCosetIndex(0, 10)
	require.Nil(t, err, "err not nil")
	assert.Equal(t, a, uint32(0))
}

func TestGetNumElement(t *testing.T) {
	numEle := rs.GetNumElement(1000, encoding.BYTES_PER_SYMBOL)
	assert.Equal(t, numEle, uint64(32))
}

func TestToFrArrayAndToByteArray_AreInverses(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	numEle := rs.GetNumElement(1000, encoding.BYTES_PER_SYMBOL)
	assert.Equal(t, numEle, uint64(32))

	cfg := encoding.DefaultConfig()
	enc, err := rs.NewEncoder(cfg)
	assert.Nil(t, err)
	require.NotNil(t, enc)

	dataFr, err := rs.ToFrArray(GETTYSBURG_ADDRESS_BYTES)
	require.Nil(t, err)
	require.NotNil(t, dataFr)

	assert.Equal(t, rs.ToByteArray(dataFr, uint64(len(GETTYSBURG_ADDRESS_BYTES))), GETTYSBURG_ADDRESS_BYTES)
}

func TestRoundUpDivision(t *testing.T) {
	a := rs.RoundUpDivision(1, 5)
	b := rs.RoundUpDivision(5, 1)

	assert.Equal(t, a, uint64(1))
	assert.Equal(t, b, uint64(5))
}
