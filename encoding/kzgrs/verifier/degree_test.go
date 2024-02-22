package verifier_test

import (
	"testing"

	"github.com/Layr-Labs/eigenda/encoding/kzgrs/prover"
	"github.com/Layr-Labs/eigenda/encoding/kzgrs/verifier"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLengthProof(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	group, _ := prover.NewProver(kzgConfig, true)
	v, _ := verifier.NewVerifier(kzgConfig, true)
	params := rs.GetEncodingParams(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))
	enc, err := group.NewKzgEncoder(params)
	require.Nil(t, err)

	numBlob := 5
	for z := 0; z < numBlob; z++ {
		extra := make([]byte, z*31*2)
		inputBytes := append(GETTYSBURG_ADDRESS_BYTES, extra...)
		inputFr := rs.ToFrArray(inputBytes)

		_, lowDegreeCommitment, lowDegreeProof, _, _, err := enc.Encode(inputFr)
		require.Nil(t, err)

		length := len(inputFr)
		assert.NoError(t, v.VerifyCommit(lowDegreeCommitment, lowDegreeProof, uint64(length)), "low degree verification failed\n")

		length = len(inputFr) - 10
		assert.Error(t, v.VerifyCommit(lowDegreeCommitment, lowDegreeProof, uint64(length)), "low degree verification failed\n")
	}
}
