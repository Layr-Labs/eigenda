package prover

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding"

	"github.com/Layr-Labs/eigenda/encoding/v2/kzg/prover/backend"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// ParametrizedProver is a prover that is configured for a specific encoding configuration.
// It contains a specific FFT setup and pre-transformed SRS points for that specific encoding config.
type ParametrizedProver struct {
	srsNumberToLoad uint64

	encodingParams encoding.EncodingParams

	computeMultiproofNumWorker uint64
	kzgMultiProofBackend       backend.KzgMultiProofsBackendV2
}

// The inputFr has not been padded to the next power of 2 field of elements. But ComputeMultiFrameProofV2
// requires it.
func (g *ParametrizedProver) GetProofs(inputFr []fr.Element, provingParams ProvingParams) ([]encoding.Proof, error) {
	// pad inputFr to BlobLength if it is not power of 2, which encodes the RS redundancy
	paddedCoeffs := make([]fr.Element, provingParams.BlobLength)
	copy(paddedCoeffs, inputFr)

	proofs, err := g.kzgMultiProofBackend.ComputeMultiFrameProofV2(
		paddedCoeffs, g.encodingParams.NumChunks, g.encodingParams.ChunkLength, g.computeMultiproofNumWorker)
	if err != nil {
		return nil, fmt.Errorf("compute multi frame proof: %w", err)
	}
	return proofs, nil
}
