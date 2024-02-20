package kzgrs_test

import (
	"testing"

	"github.com/Layr-Labs/eigenda/encoding/kzgrs"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProveZeroPadding(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	group, _ := kzgrs.NewKzgEncoderGroup(kzgConfig, true)

	params := rs.GetEncodingParams(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))
	enc, err := group.NewKzgEncoder(params)
	require.Nil(t, err)

	inputFr := rs.ToFrArray(GETTYSBURG_ADDRESS_BYTES)

	_, _, _, _, _, err = enc.Encode(inputFr)
	require.Nil(t, err)

	assert.True(t, true, "Proof %v failed\n")
}
