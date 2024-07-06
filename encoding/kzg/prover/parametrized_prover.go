package prover

import (
	"fmt"
	"log"
	"time"

	"github.com/Layr-Labs/eigenda/encoding"

	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

type ParametrizedProver struct {
	*rs.Encoder

	*kzg.KzgConfig
	Ks *kzg.KZGSettings

	Fs         *fft.FFTSettings
	Ks         *kzg.KZGSettings
	SFs        *fft.FFTSettings   // fft used for submatrix product helper
	FFTPointsT [][]bn254.G1Affine // transpose of FFTPoints

	UseGpu   bool
	Computer ProofComputer
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

	startTime := time.Now()
	// compute chunks
	poly, frames, indices, err := g.Encoder.Encode(inputFr)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	rsEncodeDone := time.Now()

	// compute commit for the full poly
	commit, err := g.Computer.ComputeCommitment(poly.Coeffs)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	commitDone := time.Now()

	lengthCommitment, err := g.Computer.ComputeLengthCommitment(poly.Coeffs)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	lengthCommitDone := time.Now()

	lengthProof, err := g.Computer.ComputeLengthProof(poly.Coeffs)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	lengthProofDone := time.Now()

	// compute proofs
	paddedCoeffs := make([]fr.Element, g.NumEvaluations())
	// polyCoeffs has less points than paddedCoeffs in general due to erasure redundancy
	copy(paddedCoeffs, poly.Coeffs)
	proofs, err := g.Computer.ComputeMultiFrameProof(paddedCoeffs, g.NumChunks, g.ChunkLength, g.NumWorker)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("could not generate proofs: %v", err)
	}
	multiProofDone := time.Now()

	if g.Verbose {
		log.Printf("    RS encode     takes  %v\n", rsEncodeDone.Sub(startTime))
		log.Printf("    Commiting     takes  %v\n", commitDone.Sub(rsEncodeDone))
		log.Printf("    LengthCommit  takes  %v\n", lengthCommitDone.Sub(commitDone))
		log.Printf("    lengthProof   takes  %v\n", lengthProofDone.Sub(lengthCommitDone))
		log.Printf("    multiProof    takes  %v\n", multiProofDone.Sub(lengthProofDone))
		log.Printf("Meta infro. order %v. shift %v\n", len(g.Srs.G2), g.SRSOrder-uint64(len(inputFr)))
	}

	// assemble frames
	kzgFrames := make([]encoding.Frame, len(frames))
	for i, index := range indices {
		kzgFrames[i] = encoding.Frame{
			Proof:  proofs[index],
			Coeffs: frames[i].Coeffs,
		}
	}

	if g.Verbose {
		log.Printf("Total encoding took      %v\n", time.Since(startTime))
	}
	return commit, lengthCommitment, lengthProof, kzgFrames, indices, nil
}
