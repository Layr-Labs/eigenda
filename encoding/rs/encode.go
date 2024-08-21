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
func (g *Encoder) EncodeBytes(inputBytes []byte) ([]Frame, []uint32, error) {
	inputFr, err := ToFrArray(inputBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot convert bytes to field elements, %w", err)
	}
	return g.Encode(inputFr)
}

// Encode function takes input in unit of Fr Element, creates a kzg commit and a list of frames
// which contains a list of multireveal interpolating polynomial coefficients, a G1 proof and a
// low degree proof corresponding to the interpolating polynomial. Each frame is an independent
// group of data verifiable to the kzg commitment. The encoding functions ensures that in each
// frame, the multireveal interpolating coefficients are identical to the part of input bytes
// in the form of field element. The extra returned integer list corresponds to which leading
// coset root of unity, the frame is proving against, which can be deduced from a frame's index
func (g *Encoder) Encode(inputFr []fr.Element) ([]Frame, []uint32, error) {
	start := time.Now()
	intermediate := time.Now()

	pdCoeffs, err := g.PadPolyEval(inputFr)
	if err != nil {
		return nil, nil, err
	}

	polyEvals, err := g.Computer.ExtendPolyEval(pdCoeffs)
	if err != nil {
		return nil, nil, err
	}

	if g.verbose {
		log.Printf("    Extending evaluation takes  %v\n", time.Since(intermediate))
	}

	// create frames to group relevant info
	frames, indices, err := g.MakeFrames(polyEvals)
	if err != nil {
		return nil, nil, err
	}

	log.Printf("  SUMMARY: RSEncode %v byte among %v numChunks with chunkLength %v takes %v\n",
		len(inputFr)*encoding.BYTES_PER_SYMBOL, g.NumChunks, g.ChunkLength, time.Since(start))

	return frames, indices, nil
}

// PadPolyEval pads the input polynomial coefficients to match the number of evaluations
// required by the encoder.
func (g *Encoder) PadPolyEval(coeffs []fr.Element) ([]fr.Element, error) {
	numEval := int(g.NumEvaluations())
	if len(coeffs) > numEval {
		return nil, fmt.Errorf("the provided encoding parameters are not sufficient for the size of the data input")
	}

	pdCoeffs := make([]fr.Element, numEval)
	copy(pdCoeffs, coeffs)

	// Pad the remaining elements with zeroes
	for i := len(coeffs); i < numEval; i++ {
		pdCoeffs[i].SetZero()
	}

	return pdCoeffs, nil
}

// MakeFrames function takes extended evaluation data and bundles relevant information into Frame.
// Every frame is verifiable to the commitment.
func (g *Encoder) MakeFrames(
	polyEvals []fr.Element,
) ([]Frame, []uint32, error) {
	// reverse dataFr making easier to sample points
	err := rb.ReverseBitOrderFr(polyEvals)
	if err != nil {
		return nil, nil, err
	}

	indices := make([]uint32, 0)
	frames := make([]Frame, g.NumChunks)

	numWorker := uint64(g.NumRSWorker)

	if numWorker > g.NumChunks {
		numWorker = g.NumChunks
	}

	jobChan := make(chan JobRequest, numWorker)
	results := make(chan error, numWorker)

	for w := uint64(0); w < numWorker; w++ {
		go g.interpolyWorker(
			polyEvals,
			jobChan,
			results,
			frames,
		)
	}

	for i := uint64(0); i < g.NumChunks; i++ {
		j := rb.ReverseBitsLimited(uint32(g.NumChunks), uint32(i))
		jr := JobRequest{
			Index: i,
		}
		jobChan <- jr
		indices = append(indices, j)
	}
	close(jobChan)

	for w := uint64(0); w < numWorker; w++ {
		interPolyErr := <-results
		if interPolyErr != nil {
			err = interPolyErr
		}
	}

	if err != nil {
		return nil, nil, fmt.Errorf("proof worker error: %v", err)
	}

	return frames, indices, nil
}

type JobRequest struct {
	Index uint64
}

func (g *Encoder) interpolyWorker(
	polyEvals []fr.Element,
	jobChan <-chan JobRequest,
	results chan<- error,
	frames []Frame,
) {

	for jr := range jobChan {
		i := jr.Index
		j := rb.ReverseBitsLimited(uint32(g.NumChunks), uint32(i))
		ys := polyEvals[g.ChunkLength*i : g.ChunkLength*(i+1)]
		err := rb.ReverseBitOrderFr(ys)
		if err != nil {
			results <- err
			continue
		}
		coeffs, err := g.GetInterpolationPolyCoeff(ys, uint32(j))
		if err != nil {
			results <- err
			continue
		}

		frames[i].Coeffs = coeffs
	}

	results <- nil

}
