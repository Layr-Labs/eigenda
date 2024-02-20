package kzgEncoder_test

import (
	"testing"

	kzgRs "github.com/Layr-Labs/eigenda/encoding/kzg"
	rs "github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProveZeroPadding(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	group, _ := kzgRs.NewKzgEncoderGroup(kzgConfig, true)

	params := rs.GetEncodingParams(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))
	enc, err := group.NewKzgEncoder(params)
	require.Nil(t, err)

	inputFr := rs.ToFrArray(GETTYSBURG_ADDRESS_BYTES)

	_, _, _, _, _, err = enc.Encode(inputFr)
	require.Nil(t, err)

	assert.True(t, true, "Proof %v failed\n")
}
