package verifier_test

import (
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/Layr-Labs/eigenda/encoding/rs"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBatchEquivalence(t *testing.T) {

	group, _ := prover.NewProver(kzgConfig, true)
	v, _ := verifier.NewVerifier(kzgConfig, true)
	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(gettysburgAddressBytes)))
	enc, err := group.GetKzgEncoder(params)
	require.Nil(t, err)

	inputFr := rs.ToPaddedFrArray(gettysburgAddressBytes)
	commit, g2commit, _, _, _, _, err := enc.Encode(inputFr)
	require.Nil(t, err)

	numBlob := 5
	commitPairs := make([]verifier.CommitmentPair, numBlob)
	for z := 0; z < numBlob; z++ {
		commitPairs[z] = verifier.CommitmentPair{
			Commitment:       *commit,
			LengthCommitment: *g2commit,
		}
	}

	assert.NoError(t, v.BatchVerifyCommitEquivalence(commitPairs), "batch equivalence negative test failed\n")

	var modifiedCommit bn254.G1Affine
	modifiedCommit.Add(commit, commit)

	for z := 0; z < numBlob; z++ {
		commitPairs[z] = verifier.CommitmentPair{
			Commitment:       modifiedCommit,
			LengthCommitment: *g2commit,
		}
	}

	assert.Error(t, v.BatchVerifyCommitEquivalence(commitPairs), "batch equivalence negative test failed\n")

	for z := 0; z < numBlob; z++ {
		commitPairs[z] = verifier.CommitmentPair{
			Commitment:       *commit,
			LengthCommitment: *g2commit,
		}
	}

	commitPairs[numBlob/2].Commitment.Add(&commitPairs[numBlob/2].Commitment, &commitPairs[numBlob/2].Commitment)

	assert.Error(t, v.BatchVerifyCommitEquivalence(commitPairs), "batch equivalence negative test failed in outer loop\n")
}
