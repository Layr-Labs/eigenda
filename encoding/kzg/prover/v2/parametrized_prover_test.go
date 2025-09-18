package prover_test

import (
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover/v2"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier/v2"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProveAllCosetThreads(t *testing.T) {
	harness := getTestHarness()

	group, err := prover.NewProver(harness.proverV2KzgConfig, nil)
	require.NoError(t, err)

	params := encoding.ParamsFromSysPar(harness.numSys, harness.numPar, uint64(len(harness.paddedGettysburgAddressBytes)))
	enc, err := group.GetKzgEncoder(params)
	require.Nil(t, err)

	inputFr, err := rs.ToFrArray(harness.paddedGettysburgAddressBytes)
	assert.Nil(t, err)

	commit, _, _, frames, _, err := enc.Encode(inputFr)
	require.Nil(t, err)

	verifierGroup, err := verifier.NewVerifier(harness.verifierV2KzgConfig, nil)
	require.Nil(t, err)
	verifier, err := verifierGroup.GetKzgVerifier(params)
	require.Nil(t, err)

	for i, frame := range frames {
		err = verifier.VerifyFrame(&frame, uint64(i), commit, params.NumChunks)
		require.Nil(t, err)
	}
}
