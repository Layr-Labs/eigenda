package prover

import (
	"fmt"
	"log"
	"time"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/hashicorp/go-multierror"

	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

type ParametrizedProver struct {
	*rs.Encoder

	*kzg.KzgConfig
	Ks *kzg.KZGSettings

	Computer ProofComputer
}

type RsEncodeResult struct {
	Frames   []rs.Frame
	Indices  []uint32
	Err      error
	Duration time.Duration
}
type LengthCommitmentResult struct {
	LengthCommitment bn254.G2Affine
	Err              error
	Duration         time.Duration
}
type LengthProofResult struct {
	LengthProof bn254.G2Affine
	Err         error
	Duration    time.Duration
}
type CommitmentResult struct {
	Commitment bn254.G1Affine
	Err        error
	Duration   time.Duration
}
type ProofsResult struct {
	Proofs   []bn254.G1Affine
	Err      error
	Duration time.Duration
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

	if len(inputFr) > int(g.KzgConfig.SRSNumberToLoad) {
		return nil, nil, nil, nil, nil, fmt.Errorf("poly Coeff length %v is greater than Loaded SRS points %v", len(inputFr), int(g.KzgConfig.SRSNumberToLoad))
	}

	encodeStart := time.Now()

	rsChan := make(chan RsEncodeResult, 1)
	lengthCommitmentChan := make(chan LengthCommitmentResult, 1)
	lengthProofChan := make(chan LengthProofResult, 1)
	commitmentChan := make(chan CommitmentResult, 1)
	proofChan := make(chan ProofsResult, 1)

	// inputFr is untouched
	// compute chunks
	go func() {
		start := time.Now()
		frames, indices, err := g.Encoder.Encode(inputFr)
		rsChan <- RsEncodeResult{
			Frames:   frames,
			Indices:  indices,
			Err:      err,
			Duration: time.Since(start),
		}
	}()

	// compute commit for the full poly
	go func() {
		start := time.Now()
		commit, err := g.Computer.ComputeCommitment(inputFr)
		commitmentChan <- CommitmentResult{
			Commitment: *commit,
			Err:        err,
			Duration:   time.Since(start),
		}
	}()

	go func() {
		start := time.Now()
		lengthCommitment, err := g.Computer.ComputeLengthCommitment(inputFr)
		lengthCommitmentChan <- LengthCommitmentResult{
			LengthCommitment: *lengthCommitment,
			Err:              err,
			Duration:         time.Since(start),
		}
	}()

	go func() {
		start := time.Now()
		lengthProof, err := g.Computer.ComputeLengthProof(inputFr)
		lengthProofChan <- LengthProofResult{
			LengthProof: *lengthProof,
			Err:         err,
			Duration:    time.Since(start),
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

		proofs, err := g.Computer.ComputeMultiFrameProof(flatpaddedCoeffs, g.NumChunks, g.ChunkLength, g.NumWorker)
		proofChan <- ProofsResult{
			Proofs:   proofs,
			Err:      err,
			Duration: time.Since(start),
		}
	}()

	lengthProofResult := <-lengthProofChan
	lengthCommitmentResult := <-lengthCommitmentChan
	commitmentResult := <-commitmentChan
	rsResult := <-rsChan
	proofsResult := <-proofChan

	if lengthProofResult.Err != nil || lengthCommitmentResult.Err != nil ||
		commitmentResult.Err != nil || rsResult.Err != nil ||
		proofsResult.Err != nil {
		return nil, nil, nil, nil, nil, multierror.Append(lengthProofResult.Err, lengthCommitmentResult.Err, commitmentResult.Err, rsResult.Err, proofsResult.Err)
	}
	totalProcessingTime := time.Since(encodeStart)

	log.Printf("\n\t\tRS encode     %-v\n\t\tCommiting     %-v\n\t\tLengthCommit  %-v\n\t\tlengthProof   %-v\n\t\tmultiProof    %-v\n\t\tMetaInfo. order  %-v shift %v\n",
		rsResult.Duration,
		commitmentResult.Duration,
		lengthCommitmentResult.Duration,
		lengthProofResult.Duration,
		proofsResult.Duration,
		g.SRSOrder,
		g.SRSOrder-uint64(len(inputFr)),
	)

	// assemble frames
	kzgFrames := make([]encoding.Frame, len(rsResult.Frames))
	for i, index := range rsResult.Indices {
		kzgFrames[i] = encoding.Frame{
			Proof:  proofsResult.Proofs[index],
			Coeffs: rsResult.Frames[i].Coeffs,
		}
	}

	if g.Verbose {
		log.Printf("Total encoding took      %v\n", totalProcessingTime)
	}
	return &commitmentResult.Commitment, &lengthCommitmentResult.LengthCommitment, &lengthProofResult.LengthProof, kzgFrames, rsResult.Indices, nil
}
