package kzgEncoder_test

import (
	"testing"

	rs "github.com/Layr-Labs/eigenda/pkg/encoding/encoder"
	kzgRs "github.com/Layr-Labs/eigenda/pkg/encoding/kzgEncoder"
	"github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBatchEquivalence(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	group, _ := kzgRs.NewKzgEncoderGroup(kzgConfig, true)
	params := rs.GetEncodingParams(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))
	enc, err := group.NewKzgEncoder(params)
	require.Nil(t, err)

	numBlob := 5
	commitPairs := make([]kzgRs.CommitmentPair, numBlob)
	for z := 0; z < numBlob; z++ {
		inputFr := rs.ToFrArray(GETTYSBURG_ADDRESS_BYTES)

		commit, g2commit, _, _, _, err := enc.Encode(inputFr)
		require.Nil(t, err)

		commitPairs[z] = kzgRs.CommitmentPair{
			Commitment:       *commit,
			LengthCommitment: *g2commit,
		}
	}

	assert.True(t, group.BatchVerifyCommitEquivalence(commitPairs) == nil, "batch equivalence test failed\n")

	for z := 0; z < numBlob; z++ {
		inputFr := rs.ToFrArray(GETTYSBURG_ADDRESS_BYTES)

		commit, g2commit, _, _, _, err := enc.Encode(inputFr)
		require.Nil(t, err)

		var modifiedCommit bn254.G1Point
		bn254.AddG1(&modifiedCommit, commit, commit)
		commitPairs[z] = kzgRs.CommitmentPair{
			Commitment:       modifiedCommit,
			LengthCommitment: *g2commit,
		}
	}

	assert.False(t, group.BatchVerifyCommitEquivalence(commitPairs) == nil, "batch equivalence negative test failed\n")

	for z := 0; z < numBlob; z++ {
		inputFr := rs.ToFrArray(GETTYSBURG_ADDRESS_BYTES)

		commit, g2commit, _, _, _, err := enc.Encode(inputFr)
		require.Nil(t, err)

		commitPairs[z] = kzgRs.CommitmentPair{
			Commitment:       *commit,
			LengthCommitment: *g2commit,
		}
	}

	bn254.AddG1(&commitPairs[numBlob/2].Commitment, &commitPairs[numBlob/2].Commitment, &commitPairs[numBlob/2].Commitment)

	assert.False(t, group.BatchVerifyCommitEquivalence(commitPairs) == nil, "batch equivalence negative test failed in outer loop\n")
}
