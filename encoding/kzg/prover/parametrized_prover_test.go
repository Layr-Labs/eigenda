package prover_test

import (
	"fmt"
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProveAllCosetThreads(t *testing.T) {
	group, err := prover.NewProver(kzgConfig, nil)
	require.NoError(t, err)

	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(gettysburgAddressBytes)))
	enc, err := group.GetKzgEncoder(params)
	require.Nil(t, err)

	inputFr, err := rs.ToFrArray(gettysburgAddressBytes)
	assert.Nil(t, err)

	commit, _, _, frames, fIndices, err := enc.Encode(inputFr)
	require.Nil(t, err)

	for i := 0; i < len(frames); i++ {
		f := frames[i]
		j := fIndices[i]

		q, err := rs.GetLeadingCosetIndex(uint64(i), numSys+numPar)
		require.Nil(t, err)

		assert.Equal(t, j, q, "leading coset inconsistency")

		fmt.Printf("frame %v leading coset %v\n", i, j)
		rs, err := enc.GetRsEncoder(params)
		require.Nil(t, err)

		lc := rs.Fs.ExpandedRootsOfUnity[uint64(j)]

		g2Atn, err := kzg.ReadG2Point(uint64(len(f.Coeffs)), kzgConfig.SRSOrder, kzgConfig.G2Path)
		require.Nil(t, err)
		assert.Nil(t, verifier.VerifyFrame(&f, enc.Ks, commit, &lc, &g2Atn), "Proof %v failed\n", i)
	}
}
