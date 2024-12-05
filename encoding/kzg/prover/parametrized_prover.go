package prover

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/hashicorp/go-multierror"

	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

type ParametrizedProver struct {
	encoding.EncodingParams
	*rs.Encoder

	KzgConfig *kzg.KzgConfig
	Ks        *kzg.KZGSettings

	KzgMultiProofBackend  KzgMultiProofsBackend
	KzgCommitmentsBackend KzgCommitmentsBackend
}

type rsEncodeResult struct {
	Frames   []rs.Frame
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

type commitmentsResult struct {
	commitment       *bn254.G1Affine
	lengthCommitment *bn254.G2Affine
	lengthProof      *bn254.G2Affine
	Error            error
}

// just a wrapper to take bytes not Fr Element
func (g *ParametrizedProver) EncodeBytes(inputBytes []byte) (*bn254.G1Affine, *bn254.G2Affine, *bn254.G2Affine, []encoding.Frame, []uint32, error) {
	inputFr, err := rs.ToFrArray(inputBytes)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("cannot convert bytes to field elements, %w", err)
	}

	return g.Encode(inputFr)
}

func (g *ParametrizedProver) Encode(inputFr []fr.Element) (*bn254.G1Affine, *bn254.G2Affine, *bn254.G2Affine, []encoding.Frame, []uint32, error) {
	if err := g.validateInput(inputFr); err != nil {
		return nil, nil, nil, nil, nil, err
	}

	encodeStart := time.Now()

	commitmentsChan := make(chan commitmentsResult, 1)

	// inputFr is untouched
	// compute chunks
	go func() {
		commitment, lengthCommitment, lengthProof, err := g.GetCommitments(inputFr, uint64(len(inputFr)))

		commitmentsChan <- commitmentsResult{
			commitment:       commitment,
			lengthCommitment: lengthCommitment,
			lengthProof:      lengthProof,
			Error:            err,
		}
	}()

	frames, indices, err := g.GetFrames(inputFr)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	commitmentResult := <-commitmentsChan
	if commitmentResult.Error != nil {
		return nil, nil, nil, nil, nil, commitmentResult.Error
	}

	slog.Info("Encoding process details",
		"Input_size_bytes", len(inputFr)*encoding.BYTES_PER_SYMBOL,
		"Num_chunks", g.NumChunks,
		"Chunk_length", g.ChunkLength,
		"Total_duration", time.Since(encodeStart),
	)

	return commitmentResult.commitment, commitmentResult.lengthCommitment, commitmentResult.lengthProof, frames, indices, nil
}

func (g *ParametrizedProver) GetCommitments(inputFr []fr.Element, length uint64) (*bn254.G1Affine, *bn254.G2Affine, *bn254.G2Affine, error) {
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
		commit, err := g.KzgCommitmentsBackend.ComputeCommitment(inputFr)
		commitmentChan <- commitmentResult{
			Commitment: commit,
			Err:        err,
			Duration:   time.Since(start),
		}
	}()

	go func() {
		start := time.Now()
		lengthCommitment, err := g.KzgCommitmentsBackend.ComputeLengthCommitment(inputFr)
		lengthCommitmentChan <- lengthCommitmentResult{
			LengthCommitment: lengthCommitment,
			Err:              err,
			Duration:         time.Since(start),
		}
	}()

	go func() {
		start := time.Now()
		lengthProof, err := g.KzgCommitmentsBackend.ComputeLengthProofForLength(inputFr, length)
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
		"Commiting_duration", commitmentResult.Duration,
		"LengthCommit_duration", lengthCommitmentResult.Duration,
		"lengthProof_duration", lengthProofResult.Duration,
		"SRSOrder", g.KzgConfig.SRSOrder,
		"SRSOrder_shift", g.KzgConfig.SRSOrder-uint64(len(inputFr)),
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

		frames, indices, err := g.Encoder.Encode(inputFr, g.EncodingParams)
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
		paddedCoeffs := make([]fr.Element, g.NumEvaluations())
		// polyCoeffs has less points than paddedCoeffs in general due to erasure redundancy
		copy(paddedCoeffs, inputFr)

		numBlob := 1
		flatpaddedCoeffs := make([]fr.Element, 0, numBlob*len(paddedCoeffs))
		for i := 0; i < numBlob; i++ {
			flatpaddedCoeffs = append(flatpaddedCoeffs, paddedCoeffs...)
		}

		proofs, err := g.KzgMultiProofBackend.ComputeMultiFrameProof(flatpaddedCoeffs, g.NumChunks, g.ChunkLength, g.KzgConfig.NumWorker)
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
		"Num_chunks", g.NumChunks,
		"Chunk_length", g.ChunkLength,
		"Total_duration", totalProcessingTime,
		"RS_encode_duration", rsResult.Duration,
		"multiProof_duration", proofsResult.Duration,
		"SRSOrder", g.KzgConfig.SRSOrder,
		"SRSOrder_shift", g.KzgConfig.SRSOrder-uint64(len(inputFr)),
	)

	// assemble frames
	kzgFrames := make([]encoding.Frame, len(rsResult.Frames))
	for i, index := range rsResult.Indices {
		kzgFrames[i] = encoding.Frame{
			Proof:  proofsResult.Proofs[index],
			Coeffs: rsResult.Frames[i].Coeffs,
		}
	}

	return kzgFrames, rsResult.Indices, nil
}

func (g *ParametrizedProver) GetMultiFrameProofs(inputFr []fr.Element) ([]encoding.Proof, error) {
	if err := g.validateInput(inputFr); err != nil {
		return nil, err
	}

	start := time.Now()

	// Pad the input polynomial to the number of evaluations
	paddingStart := time.Now()
	paddedCoeffs := make([]fr.Element, g.NumEvaluations())
	copy(paddedCoeffs, inputFr)
	paddingEnd := time.Since(paddingStart)

	proofs, err := g.KzgMultiProofBackend.ComputeMultiFrameProof(paddedCoeffs, g.NumChunks, g.ChunkLength, g.KzgConfig.NumWorker)

	end := time.Since(start)

	slog.Info("ComputeMultiFrameProofs process details",
		"Input_size_bytes", len(inputFr)*encoding.BYTES_PER_SYMBOL,
		"Num_chunks", g.NumChunks,
		"Chunk_length", g.ChunkLength,
		"Total_duration", end,
		"Padding_duration", paddingEnd,
		"SRSOrder", g.KzgConfig.SRSOrder,
		"SRSOrder_shift", g.KzgConfig.SRSOrder-uint64(len(inputFr)),
	)

	return proofs, err
}

func (g *ParametrizedProver) validateInput(inputFr []fr.Element) error {
	if len(inputFr) > int(g.KzgConfig.SRSNumberToLoad) {
		return fmt.Errorf("poly Coeff length %v is greater than Loaded SRS points %v", len(inputFr), int(g.KzgConfig.SRSNumberToLoad))
	}

	return nil
}
