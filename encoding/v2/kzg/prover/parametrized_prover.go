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

func (g *ParametrizedProver) GetProofs(inputFr []fr.Element) ([]encoding.Proof, error) {
	if err := g.validateInput(inputFr); err != nil {
		return nil, err
	}

	// pad inputFr to NumEvaluations(), which encodes the RS redundancy
	paddedCoeffs := make([]fr.Element, g.encodingParams.NumEvaluations())
	copy(paddedCoeffs, inputFr)

	proofs, err := g.kzgMultiProofBackend.ComputeMultiFrameProofV2(
		paddedCoeffs, g.encodingParams.NumChunks, g.encodingParams.ChunkLength, g.computeMultiproofNumWorker)
	if err != nil {
		return nil, fmt.Errorf("compute multi frame proof: %w", err)
	}
	return proofs, nil
}

func (g *ParametrizedProver) validateInput(inputFr []fr.Element) error {
	if len(inputFr) > int(g.srsNumberToLoad) {
		return fmt.Errorf("poly Coeff length %v is greater than Loaded SRS points %v",
			len(inputFr), int(g.srsNumberToLoad))
	}

	return nil
}
