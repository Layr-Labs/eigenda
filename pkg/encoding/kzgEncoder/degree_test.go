package kzgEncoder_test

import (
	"testing"

	rs "github.com/Layr-Labs/eigenda/pkg/encoding/encoder"
	kzgRs "github.com/Layr-Labs/eigenda/pkg/encoding/kzgEncoder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLengthProof(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	group, _ := kzgRs.NewKzgEncoderGroup(kzgConfig, true)
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
		assert.True(t, group.VerifyCommit(lowDegreeCommitment, lowDegreeProof, uint64(length)) == nil, "low degree verification failed\n")

		length = len(inputFr) - 10
		assert.False(t, group.VerifyCommit(lowDegreeCommitment, lowDegreeProof, uint64(length)) == nil, "low degree verification failed\n")
	}
}
