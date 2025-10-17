package committer

import (
	"testing"

	"github.com/Layr-Labs/eigenda/test/random"
	"github.com/stretchr/testify/require"
)

func BenchmarkCommitter_Commit(b *testing.B) {
	blobLen := uint64(1 << 19) // 2^19 = 524,288 field elements = 16 MiB
	config := Config{
		SRSNumberToLoad:   blobLen,
		G1SRSPath:         "../../../resources/srs/g1.point",
		G2SRSPath:         "../../../resources/srs/g2.point",
		G2TrailingSRSPath: "../../../resources/srs/g2.trailing.point",
	}
	committer, err := NewFromConfig(config)
	require.NoError(b, err)

	rand := random.NewTestRandom()
	blob := rand.FrElements(blobLen)

	// G1 MSM
	b.Run("blob commitment", func(b *testing.B) {
		for b.Loop() {
			_, err := committer.computeCommitmentV2(blob)
			require.NoError(b, err)
		}
	})

	// G2 MSM
	b.Run("blob length commitment", func(b *testing.B) {
		for b.Loop() {
			_, err := committer.computeLengthCommitmentV2(blob)
			require.NoError(b, err)
		}
	})

	// G2 MSM
	b.Run("blob length proof", func(b *testing.B) {
		for b.Loop() {
			_, err := committer.computeLengthProofV2(blob)
			require.NoError(b, err)
		}
	})

	b.Run("all 3", func(b *testing.B) {
		for b.Loop() {
			_, _, _, err := committer.GetCommitments(blob)
			require.NoError(b, err)
		}
	})
}
