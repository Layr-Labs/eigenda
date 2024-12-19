package verification

import (
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/stretchr/testify/require"
	"runtime"
	"testing"
)

func getKzgConfig() *kzg.KzgConfig {
	return &kzg.KzgConfig{
		G1Path:          "../../../../inabox/resources/kzg/g1.point",
		G2Path:          "../../../../inabox/resources/kzg/g2.point",
		G2PowerOf2Path:  "../../../../inabox/resources/kzg/g2.point.powerOf2",
		CacheDir:        "../../../../inabox/resources/kzg/SRSTables",
		SRSOrder:        3000,
		SRSNumberToLoad: 2900,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
		LoadG2Points:    false,
	}
}

// randomlyModifyBytes picks a random byte from the input array, and increments it
func randomlyModifyBytes(testRandom *random.TestRandom, inputBytes []byte) {
	indexToModify := testRandom.Intn(len(inputBytes))
	inputBytes[indexToModify] = inputBytes[indexToModify] + 1
}

func getRandomPaddedBytes(testRandom *random.TestRandom, count int) []byte {
	return codec.ConvertByPaddingEmptyByte(testRandom.Bytes(count))
}

func TestComputeAndCompareKzgCommitmentSuccess(t *testing.T) {
	testRandom := random.NewTestRandom(t)
	randomBytes := getRandomPaddedBytes(testRandom, 1000)

	kzgVerifier, err := verifier.NewVerifier(getKzgConfig(), nil)
	require.NotNil(t, kzgVerifier)
	require.NoError(t, err)

	commitment, err := GenerateBlobCommitment(kzgVerifier, randomBytes)
	require.NotNil(t, commitment)
	require.NoError(t, err)

	// make sure the commitment verifies correctly
	err = GenerateAndCompareBlobCommitment(
		kzgVerifier,
		commitment,
		randomBytes)
	require.NoError(t, err)
}

func TestComputeAndCompareKzgCommitmentFailure(t *testing.T) {
	testRandom := random.NewTestRandom(t)
	randomBytes := getRandomPaddedBytes(testRandom, 1000)

	kzgVerifier, err := verifier.NewVerifier(getKzgConfig(), nil)
	require.NotNil(t, kzgVerifier)
	require.NoError(t, err)

	commitment, err := GenerateBlobCommitment(kzgVerifier, randomBytes)
	require.NotNil(t, commitment)
	require.NoError(t, err)

	// randomly modify the bytes, and make sure the commitment verification fails
	randomlyModifyBytes(testRandom, randomBytes)
	err = GenerateAndCompareBlobCommitment(
		kzgVerifier,
		commitment,
		randomBytes)
	require.NotNil(t, err)
}

func TestGenerateBlobCommitmentEquality(t *testing.T) {
	testRandom := random.NewTestRandom(t)
	randomBytes := getRandomPaddedBytes(testRandom, 1000)

	kzgVerifier, err := verifier.NewVerifier(getKzgConfig(), nil)
	require.NotNil(t, kzgVerifier)
	require.NoError(t, err)

	// generate two identical commitments
	commitment1, err := GenerateBlobCommitment(kzgVerifier, randomBytes)
	require.NotNil(t, commitment1)
	require.NoError(t, err)
	commitment2, err := GenerateBlobCommitment(kzgVerifier, randomBytes)
	require.NotNil(t, commitment2)
	require.NoError(t, err)

	// commitments to identical bytes should be equal
	require.Equal(t, commitment1, commitment2)

	// randomly modify a byte
	randomlyModifyBytes(testRandom, randomBytes)
	commitmentA, err := GenerateBlobCommitment(kzgVerifier, randomBytes)
	require.NotNil(t, commitmentA)
	require.NoError(t, err)

	// commitments to non-identical bytes should not be equal
	require.NotEqual(t, commitment1, commitmentA)
}

func TestGenerateBlobCommitmentTooLong(t *testing.T) {
	kzgVerifier, err := verifier.NewVerifier(getKzgConfig(), nil)
	require.NotNil(t, kzgVerifier)
	require.NoError(t, err)

	// this is the absolute maximum number of bytes we can handle, given how the verifier was configured
	almostTooLongByteCount := 2900 * 32

	// an array of exactly this size should be fine
	almostTooLongBytes := make([]byte, almostTooLongByteCount)
	commitment1, err := GenerateBlobCommitment(kzgVerifier, almostTooLongBytes)
	require.NotNil(t, commitment1)
	require.NoError(t, err)

	// but 1 more byte is more than we can handle
	tooLongBytes := make([]byte, almostTooLongByteCount+1)
	commitment2, err := GenerateBlobCommitment(kzgVerifier, tooLongBytes)
	require.Nil(t, commitment2)
	require.NotNil(t, err)
}
