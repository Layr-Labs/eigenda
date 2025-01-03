package verification

import (
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/stretchr/testify/require"
	"math"
	"runtime"
	"testing"
)

const g1Path = "../../../../inabox/resources/kzg/g1.point"

// computeSrsNumber computes the number of SRS elements that need to be loaded for a message of given byte count
func computeSrsNumber(byteCount int) uint64 {
	return uint64(math.Ceil(float64(byteCount) / 32))
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
	randomBytes := getRandomPaddedBytes(testRandom, 100+testRandom.Intn(1000))

	srsNumberToLoad := computeSrsNumber(len(randomBytes))

	g1Srs, err := kzg.ReadG1Points(g1Path, srsNumberToLoad, uint64(runtime.GOMAXPROCS(0)))
	require.NotNil(t, g1Srs)
	require.NoError(t, err)

	commitment, err := GenerateBlobCommitment(g1Srs, randomBytes)
	require.NotNil(t, commitment)
	require.NoError(t, err)

	// make sure the commitment verifies correctly
	err = GenerateAndCompareBlobCommitment(
		g1Srs,
		randomBytes,
		commitment)
	require.NoError(t, err)
}

func TestComputeAndCompareKzgCommitmentFailure(t *testing.T) {
	testRandom := random.NewTestRandom(t)
	randomBytes := getRandomPaddedBytes(testRandom, 100+testRandom.Intn(1000))

	srsNumberToLoad := computeSrsNumber(len(randomBytes))

	g1Srs, err := kzg.ReadG1Points(g1Path, srsNumberToLoad, uint64(runtime.GOMAXPROCS(0)))
	require.NotNil(t, g1Srs)
	require.NoError(t, err)

	commitment, err := GenerateBlobCommitment(g1Srs, randomBytes)
	require.NotNil(t, commitment)
	require.NoError(t, err)

	// randomly modify the bytes, and make sure the commitment verification fails
	randomlyModifyBytes(testRandom, randomBytes)
	err = GenerateAndCompareBlobCommitment(
		g1Srs,
		randomBytes,
		commitment)
	require.NotNil(t, err)
}

func TestGenerateBlobCommitmentEquality(t *testing.T) {
	testRandom := random.NewTestRandom(t)
	randomBytes := getRandomPaddedBytes(testRandom, 100+testRandom.Intn(1000))

	srsNumberToLoad := computeSrsNumber(len(randomBytes))

	g1Srs, err := kzg.ReadG1Points(g1Path, srsNumberToLoad, uint64(runtime.GOMAXPROCS(0)))
	require.NotNil(t, g1Srs)
	require.NoError(t, err)

	// generate two identical commitments
	commitment1, err := GenerateBlobCommitment(g1Srs, randomBytes)
	require.NotNil(t, commitment1)
	require.NoError(t, err)
	commitment2, err := GenerateBlobCommitment(g1Srs, randomBytes)
	require.NotNil(t, commitment2)
	require.NoError(t, err)

	// commitments to identical bytes should be equal
	require.Equal(t, commitment1, commitment2)

	// randomly modify a byte
	randomlyModifyBytes(testRandom, randomBytes)
	commitmentA, err := GenerateBlobCommitment(g1Srs, randomBytes)
	require.NotNil(t, commitmentA)
	require.NoError(t, err)

	// commitments to non-identical bytes should not be equal
	require.NotEqual(t, commitment1, commitmentA)
}

func TestGenerateBlobCommitmentTooLong(t *testing.T) {
	srsNumberToLoad := uint64(500)

	g1Srs, err := kzg.ReadG1Points(g1Path, srsNumberToLoad, uint64(runtime.GOMAXPROCS(0)))
	require.NotNil(t, g1Srs)
	require.NoError(t, err)

	// this is the absolute maximum number of bytes we can handle, given how the verifier was configured
	almostTooLongByteCount := srsNumberToLoad * 32

	// an array of exactly this size should be fine
	almostTooLongBytes := make([]byte, almostTooLongByteCount)
	commitment1, err := GenerateBlobCommitment(g1Srs, almostTooLongBytes)
	require.NotNil(t, commitment1)
	require.NoError(t, err)

	// but 1 more byte is more than we can handle
	tooLongBytes := make([]byte, almostTooLongByteCount+1)
	commitment2, err := GenerateBlobCommitment(g1Srs, tooLongBytes)
	require.Nil(t, commitment2)
	require.NotNil(t, err)
}
