package verification

import (
	"runtime"
	"testing"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg"
	"github.com/Layr-Labs/eigenda/encoding/v2/rs"
	"github.com/Layr-Labs/eigenda/test/random"
	"github.com/stretchr/testify/require"
)

const (
	g1Path = "../../../../resources/srs/g1.point"
)

func randomBlob(t *testing.T, r *random.TestRandom, payloadSize int) *coretypes.Blob {
	blob, err := coretypes.Payload(r.Bytes(payloadSize)).ToBlob(codecs.PolynomialFormCoeff)
	require.NoError(t, err)
	return blob
}

func TestComputeAndCompareKzgCommitmentSuccess(t *testing.T) {
	testRandom := random.NewTestRandom()
	blob := randomBlob(t, testRandom, 100+testRandom.Intn(1000))

	g1Srs, err := kzg.ReadG1Points(g1Path, uint64(blob.LenSymbols()), uint64(runtime.GOMAXPROCS(0)))
	require.NotNil(t, g1Srs)
	require.NoError(t, err)

	commitment, err := GenerateBlobCommitment(g1Srs, blob.GetCoefficients())
	require.NotNil(t, commitment)
	require.NoError(t, err)

	// make sure the commitment verifies correctly
	result, err := GenerateAndCompareBlobCommitment(g1Srs, blob, commitment)
	require.True(t, result)
	require.NoError(t, err)
}

func TestComputeAndCompareKzgCommitmentFailure(t *testing.T) {
	testRandom := random.NewTestRandom()
	blob1 := randomBlob(t, testRandom, 100+testRandom.Intn(1000))

	g1Srs, err := kzg.ReadG1Points(g1Path, 1024, uint64(runtime.GOMAXPROCS(0)))
	require.NotNil(t, g1Srs)
	require.NoError(t, err)

	commitment, err := GenerateBlobCommitment(g1Srs, blob1.GetCoefficients())
	require.NotNil(t, commitment)
	require.NoError(t, err)

	// create a different blob and verify the commitment doesn't match
	blob2 := randomBlob(t, testRandom, 100+testRandom.Intn(1000))
	result, err := GenerateAndCompareBlobCommitment(g1Srs, blob2, commitment)
	require.False(t, result)
	require.NoError(t, err)
}

func TestGenerateBlobCommitmentEquality(t *testing.T) {
	testRandom := random.NewTestRandom()
	blob := randomBlob(t, testRandom, 100+testRandom.Intn(1000))
	coefficients := blob.GetCoefficients()

	g1Srs, err := kzg.ReadG1Points(g1Path, 1024, uint64(runtime.GOMAXPROCS(0)))
	require.NotNil(t, g1Srs)
	require.NoError(t, err)

	// generate two identical commitments
	commitment1, err := GenerateBlobCommitment(g1Srs, coefficients)
	require.NotNil(t, commitment1)
	require.NoError(t, err)
	commitment2, err := GenerateBlobCommitment(g1Srs, coefficients)
	require.NotNil(t, commitment2)
	require.NoError(t, err)

	// commitments to identical coefficients should be equal
	require.Equal(t, commitment1, commitment2)

	// create a different blob
	blob2 := randomBlob(t, testRandom, 100+testRandom.Intn(1000))
	commitmentA, err := GenerateBlobCommitment(g1Srs, blob2.GetCoefficients())
	require.NotNil(t, commitmentA)
	require.NoError(t, err)

	// commitments to different coefficients should not be equal
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
	almostTooLongCoeffs, err := rs.ToFrArray(almostTooLongBytes)
	require.NoError(t, err)
	commitment1, err := GenerateBlobCommitment(g1Srs, almostTooLongCoeffs)
	require.NotNil(t, commitment1)
	require.NoError(t, err)

	// but 1 more byte is more than we can handle
	tooLongBytes := make([]byte, almostTooLongByteCount+1)
	tooLongCoeffs, err := rs.ToFrArray(tooLongBytes)
	require.NoError(t, err)
	commitment2, err := GenerateBlobCommitment(g1Srs, tooLongCoeffs)
	require.Nil(t, commitment2)
	require.NotNil(t, err)
}
