package verifier_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover/v2"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier/v2"
)

func TestVerify(t *testing.T) {
	harness := getTestHarness()

	params := encoding.ParamsFromSysPar(harness.numSys, harness.numPar, uint64(len(harness.paddedGettysburgAddressBytes)))

	proverGroup, err := prover.NewProver(harness.kzgConfig, nil)
	require.Nil(t, err)
	encoder, err := proverGroup.GetKzgEncoder(params)
	require.Nil(t, err)

	verifierGroup, err := verifier.NewVerifier(harness.verifierV2KzgConfig, nil)
	require.Nil(t, err)
	verifier, err := verifierGroup.GetKzgVerifier(params)
	require.Nil(t, err)

	commit, _, _, frames, _, err := encoder.EncodeBytes(harness.paddedGettysburgAddressBytes)
	require.Nil(t, err)
	require.NotNil(t, commit)

	err = verifier.VerifyFrame(&frames[0], 0, commit, params.NumChunks)
	require.Nil(t, err)
}
