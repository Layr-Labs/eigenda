package kzgEncoder_test

import (
	"fmt"
	"testing"

	rs "github.com/Layr-Labs/eigenda/pkg/encoding/encoder"
	kzgRs "github.com/Layr-Labs/eigenda/pkg/encoding/kzgEncoder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProveAllCosetThreads(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	group, _ := kzgRs.NewKzgEncoderGroup(kzgConfig, true)

	params := rs.GetEncodingParams(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))
	enc, err := group.NewKzgEncoder(params)
	require.Nil(t, err)

	inputFr := rs.ToFrArray(GETTYSBURG_ADDRESS_BYTES)

	commit, _, _, frames, fIndices, err := enc.Encode(inputFr)
	require.Nil(t, err)

	for i := 0; i < len(frames); i++ {
		f := frames[i]
		j := fIndices[i]

		q, err := rs.GetLeadingCosetIndex(uint64(i), numSys+numPar)
		require.Nil(t, err)

		assert.Equal(t, j, q, "leading coset inconsistency")

		fmt.Printf("frame %v leading coset %v\n", i, j)
		lc := enc.Fs.ExpandedRootsOfUnity[uint64(j)]

		assert.True(t, f.Verify(enc.Ks, commit, &lc), "Proof %v failed\n", i)
	}
}
