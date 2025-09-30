package verifier_test

import (
	"strconv"
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg/committer"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier/v2"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigenda/test/random"
	"github.com/stretchr/testify/require"
)

func TestLengthProof(t *testing.T) {
	harness := getTestHarness()
	testRand := random.NewTestRandom(134)
	maxNumSymbols := uint64(1 << 19) // our stored G1 and G2 files only contain this many pts
	harness.proverV2KzgConfig.SRSNumberToLoad = maxNumSymbols

	committer, err := committer.NewFromConfig(*harness.committerConfig)
	require.Nil(t, err)

	v, err := verifier.NewVerifier(harness.verifierV2KzgConfig, nil)
	require.Nil(t, err)

	for numSymbols := uint64(1); numSymbols < maxNumSymbols; numSymbols *= 2 {
		t.Run("numSymbols="+strconv.Itoa(int(numSymbols)), func(t *testing.T) {
			inputBytes := testRand.Bytes(int(numSymbols) * encoding.BYTES_PER_SYMBOL)
			for i := range numSymbols {
				inputBytes[i*encoding.BYTES_PER_SYMBOL] = 0
			}
			inputFr, err := rs.ToFrArray(inputBytes)
			require.Nil(t, err)
			require.Equal(t, uint64(len(inputFr)), numSymbols)

			_, lengthCommitment, lengthProof, err := committer.GetCommitments(inputFr)
			require.Nil(t, err)

			require.NoError(t, v.VerifyLengthProof(lengthCommitment, lengthProof, numSymbols),
				"low degree verification failed\n")

			require.Error(t, v.VerifyLengthProof(lengthCommitment, lengthProof, numSymbols*2),
				"low degree verification failed\n")
		})
	}
}
