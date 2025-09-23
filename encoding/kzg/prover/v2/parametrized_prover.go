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

	computeMultiprooNumWorker uint64
	kzgMultiProofBackend      KzgMultiProofsBackendV2
	kzgCommitmentsBackend     KzgCommitmentsBackendV2
}

type rsEncodeResult struct {
	Frames   []rs.FrameCoeffs
	Indices  []uint32
	Duration time.Duration
	Err      error
}

type lengthCommitmentResult struct {
	LengthCommitment *bn254.G2Affine
	Duration         time.Duration
	Err              error
}

type lengthProofResult struct {
	LengthProof *bn254.G2Affine
	Duration    time.Duration
	Err         error
}

type commitmentResult struct {
	Commitment *bn254.G1Affine
	Duration   time.Duration
	Err        error
}

type proofsResult struct {
	Proofs   []bn254.G1Affine
	Duration time.Duration
	Err      error
}

func (g *ParametrizedProver) GetCommitments(
	inputFr []fr.Element,
) (*bn254.G1Affine, *bn254.G2Affine, *bn254.G2Affine, error) {
	if err := g.validateInput(inputFr); err != nil {
		return nil, nil, nil, err
	}

	encodeStart := time.Now()

	lengthCommitmentChan := make(chan lengthCommitmentResult, 1)
	lengthProofChan := make(chan lengthProofResult, 1)
	commitmentChan := make(chan commitmentResult, 1)

	// compute commit for the full poly
	go func() {
		start := time.Now()
		commit, err := g.kzgCommitmentsBackend.ComputeCommitmentV2(inputFr)
		commitmentChan <- commitmentResult{
			Commitment: commit,
			Err:        err,
			Duration:   time.Since(start),
		}
	}()

	go func() {
		start := time.Now()
		lengthCommitment, err := g.kzgCommitmentsBackend.ComputeLengthCommitmentV2(inputFr)
		lengthCommitmentChan <- lengthCommitmentResult{
			LengthCommitment: lengthCommitment,
			Err:              err,
			Duration:         time.Since(start),
		}
	}()

	go func() {
		start := time.Now()
		// blobLen must always be a power of 2 in V2
		// inputFr is not modified because padding with 0s doesn't change the commitment,
		// but we need to pretend like it was actually padded with 0s to get the correct length proof.
		blobLen := math.NextPowOf2u64(uint64(len(inputFr)))
		lengthProof, err := g.kzgCommitmentsBackend.ComputeLengthProofForLengthV2(inputFr, blobLen)
		lengthProofChan <- lengthProofResult{
			LengthProof: lengthProof,
			Err:         err,
			Duration:    time.Since(start),
		}
	}()

	lengthProofResult := <-lengthProofChan
	lengthCommitmentResult := <-lengthCommitmentChan
	commitmentResult := <-commitmentChan

	if lengthProofResult.Err != nil || lengthCommitmentResult.Err != nil ||
		commitmentResult.Err != nil {
		return nil, nil, nil, multierror.Append(lengthProofResult.Err, lengthCommitmentResult.Err, commitmentResult.Err)
	}
	totalProcessingTime := time.Since(encodeStart)

	slog.Info("Commitment process details",
		"Input_size_bytes", len(inputFr)*encoding.BYTES_PER_SYMBOL,
		"Total_duration", totalProcessingTime,
		"Committing_duration", commitmentResult.Duration,
		"LengthCommit_duration", lengthCommitmentResult.Duration,
		"lengthProof_duration", lengthProofResult.Duration,
		"SRSOrder", encoding.SRSOrder,
		// TODO(samlaf): should we take NextPowerOf2(len(inputFr)) instead?
		"SRSOrder_shift", encoding.SRSOrder-uint64(len(inputFr)),
	)

	return commitmentResult.Commitment, lengthCommitmentResult.LengthCommitment, lengthProofResult.LengthProof, nil
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
		// compute proofs
		paddedCoeffs := make([]fr.Element, g.encodingParams.NumEvaluations())
		// polyCoeffs has less points than paddedCoeffs in general due to erasure redundancy
		copy(paddedCoeffs, inputFr)

		numBlob := 1
		flatpaddedCoeffs := make([]fr.Element, 0, numBlob*len(paddedCoeffs))
		for i := 0; i < numBlob; i++ {
			flatpaddedCoeffs = append(flatpaddedCoeffs, paddedCoeffs...)
		}

		proofs, err := g.kzgMultiProofBackend.ComputeMultiFrameProofV2(
			flatpaddedCoeffs, g.encodingParams.NumChunks, g.encodingParams.ChunkLength, g.computeMultiprooNumWorker)
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
