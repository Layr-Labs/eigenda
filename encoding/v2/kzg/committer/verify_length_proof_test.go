package committer_test

import (
	"crypto/rand"
	"strconv"
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/codec"
	kzgcommitment "github.com/Layr-Labs/eigenda/encoding/v2/kzg/committer"
	"github.com/Layr-Labs/eigenda/encoding/v2/rs"
	"github.com/Layr-Labs/eigenda/test/random"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/stretchr/testify/require"
)

func TestBatchEquivalence(t *testing.T) {
	paddedGettysburgAddressBytes := codec.ConvertByPaddingEmptyByte([]byte("Fourscore and seven years ago our fathers brought forth, on this continent, a new nation, conceived in liberty, and dedicated to the proposition that all men are created equal. Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived, and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting-place for those who here gave their lives, that that nation might live. It is altogether fitting and proper that we should do this. But, in a larger sense, we cannot dedicate, we cannot consecrate—we cannot hallow—this ground. The brave men, living and dead, who struggled here, have consecrated it far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us—that from these honored dead we take increased devotion to that cause for which they here gave the last full measure of devotion—that we here highly resolve that these dead shall not have died in vain—that this nation, under God, shall have a new birth of freedom, and that government of the people, by the people, for the people, shall not perish from the earth."))

	committer, err := kzgcommitment.NewFromConfig(kzgcommitment.Config{
		SRSNumberToLoad:   4096,
		G1SRSPath:         "../../../../resources/srs/g1.point",
		G2SRSPath:         "../../../../resources/srs/g2.point",
		G2TrailingSRSPath: "../../../../resources/srs/g2.trailing.point",
	})
	require.NoError(t, err)

	commitment, err := committer.GetCommitmentsForPaddedLength(paddedGettysburgAddressBytes)
	require.NoError(t, err)

	numBlob := 5
	commitments := make([]encoding.BlobCommitments, numBlob)
	for z := 0; z < numBlob; z++ {
		commitments[z] = commitment
	}

	require.NoError(t, kzgcommitment.VerifyCommitEquivalenceBatch(commitments), "batch equivalence negative test failed\n")

	var modifiedCommit bn254.G1Affine
	modifiedCommit.Add((*bn254.G1Affine)(commitment.Commitment), (*bn254.G1Affine)(commitment.Commitment))

	for z := 0; z < numBlob; z++ {
		commitments[z].Commitment = (*encoding.G1Commitment)(&modifiedCommit)
	}

	require.Error(t, kzgcommitment.VerifyCommitEquivalenceBatch(commitments), "batch equivalence negative test failed\n")

	for z := 0; z < numBlob; z++ {
		commitments[z] = commitment
	}
	commitments[numBlob/2].Commitment = (*encoding.G1Commitment)(&modifiedCommit)

	require.Error(t, kzgcommitment.VerifyCommitEquivalenceBatch(commitments),
		"batch equivalence negative test failed in outer loop\n")
}

func TestLengthProof(t *testing.T) {
	testRand := random.NewTestRandom(134)
	maxNumSymbols := uint64(1 << 19) // our stored G1 and G2 files only contain this many pts

	committer, err := kzgcommitment.NewFromConfig(kzgcommitment.Config{
		SRSNumberToLoad:   maxNumSymbols,
		G1SRSPath:         "../../../../resources/srs/g1.point",
		G2SRSPath:         "../../../../resources/srs/g2.point",
		G2TrailingSRSPath: "../../../../resources/srs/g2.trailing.point",
	})
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

			commitments, err := committer.GetCommitmentsForPaddedLength(inputBytes)
			require.Nil(t, err)

			require.NoError(t, kzgcommitment.VerifyLengthProof(commitments), "low degree verification failed\n")

			commitments.Length *= 2
			require.Error(t, kzgcommitment.VerifyLengthProof(commitments), "low degree verification failed\n")
		})
	}
}

func BenchmarkVerifyBlob(b *testing.B) {
	committer, err := kzgcommitment.NewFromConfig(kzgcommitment.Config{
		SRSNumberToLoad:   4096,
		G1SRSPath:         "../../../../resources/srs/g1.point",
		G2SRSPath:         "../../../../resources/srs/g2.point",
		G2TrailingSRSPath: "../../../../resources/srs/g2.trailing.point",
	})
	require.NoError(b, err)

	blobSize := 8 * 256
	numSamples := 30
	blobs := make([][]byte, numSamples)
	for i := 0; i < numSamples; i++ {
		blob := make([]byte, blobSize)
		_, _ = rand.Read(blob)
		blobs[i] = blob
	}

	commitments, err := committer.GetCommitmentsForPaddedLength(codec.ConvertByPaddingEmptyByte(blobs[0]))
	require.NoError(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err = kzgcommitment.VerifyLengthProof(commitments)
		require.NoError(b, err)
	}
}
