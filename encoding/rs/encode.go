package rs

import (
	"fmt"
	"log"
	"time"

	"github.com/Layr-Labs/eigenda/encoding"
	rb "github.com/Layr-Labs/eigenda/encoding/utils/reverseBits"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

type GlobalPoly struct {
	Coeffs []fr.Element
	Values []fr.Element
}

// just a wrapper to take bytes not Fr Element
func (g *Encoder) EncodeBytes(inputBytes []byte) (*GlobalPoly, []Frame, []uint32, error) {
	inputFr := ToFrArray(inputBytes)
	return g.Encode(inputFr)
}

// Encode function takes input in unit of Fr Element, creates a kzg commit and a list of frames
// which contains a list of multireveal interpolating polynomial coefficients, a G1 proof and a
// low degree proof corresponding to the interpolating polynomial. Each frame is an independent
// group of data verifiable to the kzg commitment. The encoding functions ensures that in each
// frame, the multireveal interpolating coefficients are identical to the part of input bytes
// in the form of field element. The extra returned integer list corresponds to which leading
// coset root of unity, the frame is proving against, which can be deduced from a frame's index
func (g *Encoder) Encode(inputFr []fr.Element) (*GlobalPoly, []Frame, []uint32, error) {
	start := time.Now()
	intermediate := time.Now()

	polyCoeffs := inputFr

	// extend data based on Sys, Par ratio. The returned fullCoeffsPoly is padded with 0 to ease proof
	polyEvals, _, err := g.ExtendPolyEval(polyCoeffs)
	if err != nil {
		return nil, nil, nil, err
	}

	poly := &GlobalPoly{
		Values: polyEvals,
		Coeffs: polyCoeffs,
	}

	if g.verbose {
		log.Printf("    Extending evaluation takes  %v\n", time.Since(intermediate))
	}

	// create frames to group relevant info
	frames, indices, err := g.MakeFrames(polyEvals)
	if err != nil {
		return nil, nil, nil, err
	}

	log.Printf("  SUMMARY: Encode %v byte among %v numNode takes %v\n",
		len(inputFr)*encoding.BYTES_PER_COEFFICIENT, g.NumChunks, time.Since(start))

	return poly, frames, indices, nil
}

// This Function takes extended evaluation data and bundles relevant information into Frame.
// Every frame is verifiable to the commitment.
func (g *Encoder) MakeFrames(
	polyEvals []fr.Element,
) ([]Frame, []uint32, error) {
	// reverse dataFr making easier to sample points
	err := rb.ReverseBitOrderFr(polyEvals)
	if err != nil {
		return nil, nil, err
	}
	k := uint64(0)

	indices := make([]uint32, 0)
	frames := make([]Frame, g.NumChunks)

	for i := uint64(0); i < uint64(g.NumChunks); i++ {

		// finds out which coset leader i-th node is having
		j := rb.ReverseBitsLimited(uint32(g.NumChunks), uint32(i))

		// mutltiprover return proof in butterfly order
		frame := Frame{}
		indices = append(indices, j)

		ys := polyEvals[g.ChunkLength*i : g.ChunkLength*(i+1)]
		err := rb.ReverseBitOrderFr(ys)
		if err != nil {
			return nil, nil, err
		}
		coeffs, err := g.GetInterpolationPolyCoeff(ys, uint32(j))
		if err != nil {
			return nil, nil, err
		}

		frame.Coeffs = coeffs

		frames[k] = frame
		k++
	}

	return frames, indices, nil
}

// Encoding Reed Solomon using FFT
func (g *Encoder) ExtendPolyEval(coeffs []fr.Element) ([]fr.Element, []fr.Element, error) {

	if len(coeffs) > int(g.NumEvaluations()) {
		return nil, nil, fmt.Errorf("the provided encoding parameters are not sufficient for the size of the data input")
	}

	pdCoeffs := make([]fr.Element, g.NumEvaluations())
	for i := 0; i < len(coeffs); i++ {
		//bls.CopyFr(&pdCoeffs[i], &coeffs[i])
		pdCoeffs[i].Set(&coeffs[i])
	}
	for i := len(coeffs); i < len(pdCoeffs); i++ {
		pdCoeffs[i].SetZero()
		//bls.CopyFr(&pdCoeffs[i], &bls.ZERO)
	}

	evals, err := g.Fs.FFT(pdCoeffs, false)
	if err != nil {
		return nil, nil, err
	}

	return evals, pdCoeffs, nil
}
