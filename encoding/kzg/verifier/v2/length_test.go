package verifier_test

import (
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover/v2"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier/v2"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLengthProof(t *testing.T) {
	harness := getTestHarness()

	group, err := prover.NewProver(harness.kzgConfig, nil)
	require.Nil(t, err)

	v, err := verifier.NewVerifier(harness.verifierV2KzgConfig, nil)
	require.Nil(t, err)

	params := encoding.ParamsFromSysPar(harness.numSys, harness.numPar, uint64(len(harness.paddedGettysburgAddressBytes)))
	enc, err := group.GetKzgEncoder(params)
	require.Nil(t, err)

	numBlob := 5
	for z := 0; z < numBlob; z++ {
		extra := make([]byte, z*32*2)
		inputBytes := append(harness.paddedGettysburgAddressBytes, extra...)
		inputFr, err := rs.ToFrArray(inputBytes)
		require.Nil(t, err)

		_, lengthCommitment, lengthProof, _, _, err := enc.Encode(inputFr)
		require.Nil(t, err)

		length := len(inputFr)
		assert.NoError(t, v.VerifyLengthProof(lengthCommitment, lengthProof, uint64(length)),
			"low degree verification failed\n")

		length = len(inputFr) - 10
		assert.Error(t, v.VerifyLengthProof(lengthCommitment, lengthProof, uint64(length)),
			"low degree verification failed\n")
	}
}
