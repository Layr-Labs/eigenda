package rs

import (
	"fmt"
	"log/slog"
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
	paddingDuration := time.Since(intermediate)

	intermediate = time.Now()

	polyEvals, err := encoder.RSEncoderComputer.ExtendPolyEval(pdCoeffs)
	if err != nil {
		return nil, nil, err
	}
	extensionDuration := time.Since(intermediate)

	intermediate = time.Now()

	// create frames to group relevant info
	frames, indices, err := encoder.MakeFrames(polyEvals)
	if err != nil {
		return nil, nil, err
	}

	framesDuration := time.Since(intermediate)

	slog.Info("RSEncode details",
		"input_size_bytes", len(inputFr)*encoding.BYTES_PER_SYMBOL,
		"num_chunks", encoder.NumChunks,
		"chunk_length", encoder.ChunkLength,
		"padding_duration", paddingDuration,
		"extension_duration", extensionDuration,
		"frames_duration", framesDuration,
		"total_duration", time.Since(start))

	return frames, indices, nil
}
