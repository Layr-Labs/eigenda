package rs

import (
	"fmt"
	"log"
	"time"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

type GlobalPoly struct {
	Coeffs []fr.Element
	Values []fr.Element
}

// just a wrapper to take bytes not Fr Element
func (g *Encoder) EncodeBytes(inputBytes []byte, params encoding.EncodingParams) ([]Frame, []uint32, error) {
	inputFr, err := ToFrArray(inputBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot convert bytes to field elements, %w", err)
	}
	return g.Encode(inputFr, params)
}

// Encode function takes input in unit of Fr Element, creates a kzg commit and a list of frames
// which contains a list of multireveal interpolating polynomial coefficients, a G1 proof and a
// low degree proof corresponding to the interpolating polynomial. Each frame is an independent
// group of data verifiable to the kzg commitment. The encoding functions ensures that in each
// frame, the multireveal interpolating coefficients are identical to the part of input bytes
// in the form of field element. The extra returned integer list corresponds to which leading
// coset root of unity, the frame is proving against, which can be deduced from a frame's index
func (g *Encoder) Encode(inputFr []fr.Element, params encoding.EncodingParams) ([]Frame, []uint32, error) {
	start := time.Now()
	intermediate := time.Now()

	// Get RS encoder from params
	encoder, err := g.GetRsEncoder(params)
	if err != nil {
		return nil, nil, err
	}

	pdCoeffs, err := encoder.PadPolyEval(inputFr)
	if err != nil {
		return nil, nil, err
	}

	if g.verbose {
		log.Printf("    Padding takes  %v\n", time.Since(intermediate))
	}

	polyEvals, err := encoder.RSEncoderComputer.ExtendPolyEval(pdCoeffs)
	if err != nil {
		return nil, nil, err
	}

	if g.verbose {
		log.Printf("    Extending evaluation takes  %v\n", time.Since(intermediate))
	}

	// create frames to group relevant info
	frames, indices, err := encoder.MakeFrames(polyEvals)
	if err != nil {
		return nil, nil, err
	}

	log.Printf("  SUMMARY: RSEncode %v byte among %v numChunks with chunkLength %v takes %v\n",
		len(inputFr)*encoding.BYTES_PER_SYMBOL, encoder.NumChunks, encoder.ChunkLength, time.Since(start))

	return frames, indices, nil
}
