package verifier_test

import (
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/stretchr/testify/require"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg/committer"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover/v2"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier/v2"
)

func TestVerify(t *testing.T) {
	harness := getTestHarness()

	params := encoding.ParamsFromSysPar(harness.numSys, harness.numPar, uint64(len(harness.paddedGettysburgAddressBytes)))

	proverGroup, err := prover.NewProver(harness.proverV2KzgConfig, nil)
	require.Nil(t, err)

	committer, err := committer.NewFromConfig(*harness.committerConfig)
	require.Nil(t, err)

	frames, err := proverGroup.GetFrames(harness.paddedGettysburgAddressBytes, params)
	require.Nil(t, err)
	commitments, err := committer.GetCommitmentsForPaddedLength(harness.paddedGettysburgAddressBytes)
	require.Nil(t, err)

	verifierGroup, err := verifier.NewVerifier(harness.verifierV2KzgConfig, nil)
	require.Nil(t, err)
	verifier, err := verifierGroup.GetKzgVerifier(params)
	require.Nil(t, err)

	err = verifier.VerifyFrame(frames[0], 0, (*bn254.G1Affine)(commitments.Commitment), params.NumChunks)
	require.Nil(t, err)
}
