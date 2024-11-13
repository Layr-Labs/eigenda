package rs

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	rb "github.com/Layr-Labs/eigenda/encoding/utils/reverseBits"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

type ParametrizedEncoder struct {
	*encoding.Config
	encoding.EncodingParams
	Fs                *fft.FFTSettings
	RSEncoderComputer EncoderDevice
}

// PadPolyEval pads the input polynomial coefficients to match the number of evaluations
// required by the encoder.
func (g *ParametrizedEncoder) PadPolyEval(coeffs []fr.Element) ([]fr.Element, error) {
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
func (g *ParametrizedEncoder) MakeFrames(
	polyEvals []fr.Element,
) ([]Frame, []uint32, error) {
	// reverse dataFr making easier to sample points
	err := rb.ReverseBitOrderFr(polyEvals)
	if err != nil {
		return nil, nil, err
	}

	indices := make([]uint32, 0)
	frames := make([]Frame, g.NumChunks)

	numWorker := uint64(g.Config.NumWorker)
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

func (g *ParametrizedEncoder) interpolyWorker(
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
