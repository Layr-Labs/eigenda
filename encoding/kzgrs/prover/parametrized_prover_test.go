package prover_test

import (
	"fmt"
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzgrs"
	"github.com/Layr-Labs/eigenda/encoding/kzgrs/prover"
	"github.com/Layr-Labs/eigenda/encoding/kzgrs/verifier"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProveAllCosetThreads(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	group, _ := prover.NewProver(kzgConfig, true)

	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))
	enc, err := group.GetKzgEncoder(params)
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

		g2Atn, err := kzgrs.ReadG2Point(uint64(len(f.Coeffs)), kzgConfig)
		require.Nil(t, err)
		assert.True(t, verifier.VerifyFrame(&f, enc.Ks, commit, &lc, &g2Atn), "Proof %v failed\n", i)
	}
}
