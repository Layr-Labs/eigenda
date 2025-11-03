package node

import (
	"testing"

	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/test/random"
	"github.com/stretchr/testify/require"
)

func testIndexToRangeConversion(
	t *testing.T,
	indexProbability float64,
) {

	rand := random.NewTestRandom()
	maxIndex := uint32(1024 * 8)

	indices := make([]uint32, 0)

	// For each possible index, choose whether it will be present based on the given probability.
	// Lower indexProbability values will result in sparse sets of indices, while higher ones will
	// result in denser sets of indices.
	for i := uint32(0); i < maxIndex; i++ {
		if rand.Float64() < indexProbability {
			indices = append(indices, i)
		}
	}

	var blobKey corev2.BlobKey
	chunkRequests := convertIndicesToRangeRequests(blobKey, indices)

	// Iterate over the generated chunk requests and reconstruct the requested indices.
	reconstructedIndices := make([]uint32, 0)
	for _, chunkRequestByRange := range chunkRequests {
		for i := chunkRequestByRange.Start; i < chunkRequestByRange.End; i++ {
			reconstructedIndices = append(reconstructedIndices, i)
		}
	}

	require.Equal(t, indices, reconstructedIndices)

}

func TestIndexToRangeConversion(t *testing.T) {
	t.Run("No Indices", func(t *testing.T) {
		testIndexToRangeConversion(t, 0.0)
	})
	t.Run("Very Sparse Indices", func(t *testing.T) {
		testIndexToRangeConversion(t, 0.01)
	})
	t.Run("Sparse Indices", func(t *testing.T) {
		testIndexToRangeConversion(t, 0.1)
	})
	t.Run("Moderate Indices", func(t *testing.T) {
		testIndexToRangeConversion(t, 0.5)
	})
	t.Run("Dense Indices", func(t *testing.T) {
		testIndexToRangeConversion(t, 0.9)
	})
	t.Run("All Indices", func(t *testing.T) {
		testIndexToRangeConversion(t, 1.0)
	})
}
