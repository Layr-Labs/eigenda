package kzgEncoder_test

import (
	"testing"

	rs "github.com/Layr-Labs/eigenda/pkg/encoding/encoder"
	kzgRs "github.com/Layr-Labs/eigenda/pkg/encoding/kzgEncoder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBatchEquivalence(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	group, _ := kzgRs.NewKzgEncoderGroup(kzgConfig)
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
}
