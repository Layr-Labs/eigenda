package verifier_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Layr-Labs/eigenda/encoding"

	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
)

func TestVerify(t *testing.T) {
	group, err := prover.NewProver(kzgConfig, nil)
	require.Nil(t, err)

	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(gettysburgAddressBytes)))

	enc, err := group.GetKzgEncoder(params)
	require.Nil(t, err)
	require.NotNil(t, enc)

	commit, _, _, frames, _, err := enc.EncodeBytes(gettysburgAddressBytes)
	require.Nil(t, err)
	require.NotNil(t, commit)
	require.NotNil(t, frames)

	n := uint8(math.Log2(float64(params.NumEvaluations())))
	fs := fft.NewFFTSettings(n)
	require.NotNil(t, fs)

	lc := fs.ExpandedRootsOfUnity[uint64(0)]
	require.NotNil(t, lc)

	g2Atn, err := kzg.ReadG2Point(uint64(len(frames[0].Coeffs)), kzgConfig.SRSOrder, kzgConfig.G2Path)
	require.Nil(t, err)
	assert.Nil(t, verifier.VerifyFrame(&frames[0], enc.Ks, commit, &lc, &g2Atn))
}
