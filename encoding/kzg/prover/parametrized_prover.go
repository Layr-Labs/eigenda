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

	Computer ProofDevice
}

type rsEncodeResult struct {
	Frames   []rs.Frame
	Indices  []uint32
	Duration time.Duration
	Err      error
}
type lengthCommitmentResult struct {
	LengthCommitment bn254.G2Affine
	Duration         time.Duration
	Err              error
}
type lengthProofResult struct {
	LengthProof bn254.G2Affine
	Duration    time.Duration
	Err         error
}
type commitmentResult struct {
	Commitment bn254.G1Affine
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

	if len(inputFr) > int(g.KzgConfig.SRSNumberToLoad) {
		return nil, nil, nil, nil, nil, fmt.Errorf("poly Coeff length %v is greater than Loaded SRS points %v", len(inputFr), int(g.KzgConfig.SRSNumberToLoad))
	}

	encodeStart := time.Now()

	commitmentsChan := make(chan commitmentsResult, 1)

	// inputFr is untouched
	// compute chunks
	go func() {
		commitment, lengthCommitment, lengthProof, err := g.GetCommitments(inputFr)

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

	totalProcessingTime := time.Since(encodeStart)

	if g.Verbose {
		log.Printf("Total encoding took      %v\n", totalProcessingTime)
	}
	return commitmentResult.commitment, commitmentResult.lengthCommitment, commitmentResult.lengthProof, frames, indices, nil
}

func (g *ParametrizedProver) GetCommitments(inputFr []fr.Element) (*bn254.G1Affine, *bn254.G2Affine, *bn254.G2Affine, error) {

	if len(inputFr) > int(g.KzgConfig.SRSNumberToLoad) {
		return nil, nil, nil, fmt.Errorf("poly Coeff length %v is greater than Loaded SRS points %v", len(inputFr), int(g.KzgConfig.SRSNumberToLoad))
	}

	encodeStart := time.Now()

	lengthCommitmentChan := make(chan lengthCommitmentResult, 1)
	lengthProofChan := make(chan lengthProofResult, 1)
	commitmentChan := make(chan commitmentResult, 1)

	// compute commit for the full poly
	go func() {
		start := time.Now()
		commit, err := g.Computer.ComputeCommitment(inputFr)
		commitmentChan <- commitmentResult{
			Commitment: *commit,
			Err:        err,
			Duration:   time.Since(start),
		}
	}()

	go func() {
		start := time.Now()
		lengthCommitment, err := g.Computer.ComputeLengthCommitment(inputFr)
		lengthCommitmentChan <- lengthCommitmentResult{
			LengthCommitment: *lengthCommitment,
			Err:              err,
			Duration:         time.Since(start),
		}
	}()

	go func() {
		start := time.Now()
		lengthProof, err := g.Computer.ComputeLengthProof(inputFr)
		lengthProofChan <- lengthProofResult{
			LengthProof: *lengthProof,
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

	log.Printf("\n\t\tCommiting     %-v\n\t\tLengthCommit  %-v\n\t\tlengthProof   %-v\n\t\tMetaInfo. order  %-v shift %v\n",
		commitmentResult.Duration,
		lengthCommitmentResult.Duration,
		lengthProofResult.Duration,
		g.SRSOrder,
		g.SRSOrder-uint64(len(inputFr)),
	)

	if g.Verbose {
		log.Printf("Total encoding took      %v\n", totalProcessingTime)
	}
	return &commitmentResult.Commitment, &lengthCommitmentResult.LengthCommitment, &lengthProofResult.LengthProof, nil
}

func (g *ParametrizedProver) GetFrames(inputFr []fr.Element) ([]encoding.Frame, []uint32, error) {

	if len(inputFr) > int(g.KzgConfig.SRSNumberToLoad) {
		return nil, nil, fmt.Errorf("poly Coeff length %v is greater than Loaded SRS points %v", len(inputFr), int(g.KzgConfig.SRSNumberToLoad))
	}

	proofChan := make(chan proofsResult, 1)
	rsChan := make(chan rsEncodeResult, 1)

	// inputFr is untouched
	// compute chunks
	go func() {
		start := time.Now()
		frames, indices, err := g.Encoder.Encode(inputFr)
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

		proofs, err := g.Computer.ComputeMultiFrameProof(flatpaddedCoeffs, g.NumChunks, g.ChunkLength, g.NumWorker)
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

	log.Printf("\n\t\tRS encode     %-v\n\t\tmultiProof    %-v\n\t\tMetaInfo. order  %-v shift %v\n",
		rsResult.Duration,
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

	return kzgFrames, rsResult.Indices, nil

}
