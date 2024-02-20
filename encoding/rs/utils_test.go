package encoder_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	rs "github.com/Layr-Labs/eigenda/encoding/rs"
)

func TestGetEncodingParams(t *testing.T) {
	params := rs.GetEncodingParams(1, 4, 1000)

	require.NotNil(t, params)
	assert.Equal(t, params.ChunkLen, uint64(64))
	// assert.Equal(t, params.DataLen, uint64(1000))
	assert.Equal(t, params.NumChunks, uint64(8))
	assert.Equal(t, params.NumEvaluations(), uint64(512))
}

func TestGetLeadingCoset(t *testing.T) {
	a, err := rs.GetLeadingCosetIndex(0, 10)
	require.Nil(t, err, "err not nil")
	assert.Equal(t, a, uint32(0))
}

func TestGetNumElement(t *testing.T) {
	numEle := rs.GetNumElement(1000, BYTES_PER_COEFFICIENT)
	assert.Equal(t, numEle, uint64(33))
}

func TestToFrArrayAndToByteArray_AreInverses(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	numEle := rs.GetNumElement(1000, BYTES_PER_COEFFICIENT)
	assert.Equal(t, numEle, uint64(33))

	params := rs.GetEncodingParams(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))
	enc, _ := rs.NewEncoder(params, true)
	require.NotNil(t, enc)

	dataFr := rs.ToFrArray(GETTYSBURG_ADDRESS_BYTES)
	require.NotNil(t, dataFr)

	assert.Equal(t, rs.ToByteArray(dataFr, uint64(len(GETTYSBURG_ADDRESS_BYTES))), GETTYSBURG_ADDRESS_BYTES)
}

func TestRoundUpDivision(t *testing.T) {
	a := rs.RoundUpDivision(1, 5)
	b := rs.RoundUpDivision(5, 1)

	assert.Equal(t, a, uint64(1))
	assert.Equal(t, b, uint64(5))
}
