package prover_test

import (
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/v1/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/v1/kzg/verifier"
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

	commit, _, _, frames, _, err := enc.Encode(inputFr)
	require.Nil(t, err)

	verifierGroup, err := verifier.NewVerifier(kzgConfig, nil)
	require.Nil(t, err)
	verifier, err := verifierGroup.GetKzgVerifier(params)
	require.Nil(t, err)

	for i, frame := range frames {
		err = verifier.VerifyFrame(&frame, uint64(i), commit, params.NumChunks)
		require.Nil(t, err)
	}
}
