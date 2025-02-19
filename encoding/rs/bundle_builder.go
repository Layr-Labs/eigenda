package rs

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/core"
)

// BuildChunksData creates a binary core.ChunksData object from the given proofs and coefficients.
func BuildChunksData(
	proofs [][]byte,
	chunkLen int,
	coefficients [][]byte) (*core.ChunksData, error) {

	if len(proofs) != len(coefficients) {
		return nil, fmt.Errorf("proofs and coefficients have different lengths (%d vs %d)",
			len(proofs), len(coefficients))
	}

	binaryChunks := make([][]byte, len(proofs))

	for i := 0; i < len(proofs); i++ {
		binaryFrame := make([]byte, len(proofs[i])+len(coefficients[i]))
		copy(binaryFrame, proofs[i])
		copy(binaryFrame[len(proofs[i]):], coefficients[i])
		binaryChunks[i] = binaryFrame
	}

	return &core.ChunksData{
		Chunks:   binaryChunks,
		Format:   core.GnarkChunkEncodingFormat,
		ChunkLen: chunkLen,
	}, nil
}
