package prover

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/Layr-Labs/eigenda/common/math"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/hashicorp/go-multierror"

	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// ParametrizedProver is a prover that is configured for a specific encoding configuration.
// It contains a specific FFT setup and pre-transformed SRS points for that specific encoding config.
// Note that commitments are not dependent on the FFT setup.
// TODO(samlaf): move the commitment functionality back to the prover, not parametrizedProver.
type ParametrizedProver struct {
	srsNumberToLoad uint64

	encodingParams encoding.EncodingParams
	encoder        *rs.Encoder

	computeMultiproofNumWorker uint64
	kzgMultiProofBackend       KzgMultiProofsBackendV2
}

type rsEncodeResult struct {
	Frames   []rs.FrameCoeffs
	Indices  []uint32
	Duration time.Duration
	Err      error
}

type proofsResult struct {
	Proofs   []bn254.G1Affine
	Duration time.Duration
	Err      error
}

func (g *ParametrizedProver) GetFrames(inputFr []fr.Element) ([]encoding.Frame, []uint32, error) {
	if err := g.validateInput(inputFr); err != nil {
		return nil, nil, err
	}

	encodeStart := time.Now()

	proofChan := make(chan proofsResult, 1)
	rsChan := make(chan rsEncodeResult, 1)

	// inputFr is untouched
	// compute chunks
	go func() {
		start := time.Now()

		frames, indices, err := g.encoder.Encode(inputFr, g.encodingParams)
		rsChan <- rsEncodeResult{
			Frames:   frames,
			Indices:  indices,
			Err:      err,
			Duration: time.Since(start),
		}
	}()

	go func() {
		start := time.Now()
		// compute proofs to a polynomial that contains power of 2 symbols
		// padding are 0 representing coefficients at higher monomials
		paddedCoeffs := make([]fr.Element, math.NextPowOf2u64(uint64(len(inputFr))))
		copy(paddedCoeffs, inputFr)

		proofs, err := g.kzgMultiProofBackend.ComputeMultiFrameProofV2(
			paddedCoeffs, g.encodingParams.NumChunks, g.encodingParams.ChunkLength, g.computeMultiproofNumWorker)
		proofChan <- proofsResult{
			Proofs:   proofs,
			Err:      err,
			Duration: time.Since(start),
		}
	}()

	rsResult := <-rsChan
	proofsResult := <-proofChan

	if rsResult.Err != nil || proofsResult.Err != nil {
		return nil, nil, multierror.Append(rsResult.Err, proofsResult.Err)
	}

	totalProcessingTime := time.Since(encodeStart)
	slog.Info("Frame process details",
		"Input_size_bytes", len(inputFr)*encoding.BYTES_PER_SYMBOL,
		"Num_chunks", g.encodingParams.NumChunks,
		"Chunk_length", g.encodingParams.ChunkLength,
		"Total_duration", totalProcessingTime,
		"RS_encode_duration", rsResult.Duration,
		"multiProof_duration", proofsResult.Duration,
		"SRSOrder", encoding.SRSOrder,
		// TODO(samlaf): should we take NextPowerOf2(len(inputFr)) instead?
		"SRSOrder_shift", encoding.SRSOrder-uint64(len(inputFr)),
	)

	// assemble frames
	kzgFrames := make([]encoding.Frame, len(rsResult.Frames))
	for i, index := range rsResult.Indices {
		kzgFrames[i] = encoding.Frame{
			Proof:  proofsResult.Proofs[index],
			Coeffs: rsResult.Frames[i],
		}
	}

	return kzgFrames, rsResult.Indices, nil
}

func (g *ParametrizedProver) validateInput(inputFr []fr.Element) error {
	if len(inputFr) > int(g.srsNumberToLoad) {
		return fmt.Errorf("poly Coeff length %v is greater than Loaded SRS points %v",
			len(inputFr), int(g.srsNumberToLoad))
	}

	return nil
}
