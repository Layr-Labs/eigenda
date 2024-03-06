package verifier_test

import (
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLengthProof(t *testing.T) {

	group, _ := prover.NewProver(kzgConfig, true)
	v, _ := verifier.NewVerifier(kzgConfig, true)
	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(gettysburgAddressBytes)))
	enc, err := group.GetKzgEncoder(params)
	require.Nil(t, err)

	numBlob := 5
	for z := 0; z < numBlob; z++ {
		extra := make([]byte, z*31*2)
		inputBytes := append(gettysburgAddressBytes, extra...)
		inputFr := rs.ToFrArray(inputBytes)

		_, lowDegreeCommitment, lowDegreeProof, _, _, err := enc.Encode(inputFr)
		require.Nil(t, err)

		length := len(inputFr)
		assert.NoError(t, v.VerifyCommit(lowDegreeCommitment, lowDegreeProof, uint64(length)), "low degree verification failed\n")

		length = len(inputFr) - 10
		assert.Error(t, v.VerifyCommit(lowDegreeCommitment, lowDegreeProof, uint64(length)), "low degree verification failed\n")
	}
}
