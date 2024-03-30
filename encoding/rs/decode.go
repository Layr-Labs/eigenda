package rs

import (
	"errors"
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// Decode data when some chunks from systematic nodes are lost. It first uses FFT to recover
// the whole polynomial. Then it extracts only the systematic chunks.
// It takes a list of available frame, and return the original encoded data
// storing the evaluation points, since it is where RS is applied. The input frame contains
// the coefficient of the interpolating polynomina, hence interpolation is needed before
// recovery.
// maxInputSize is the upper bound of the original data size. This is needed because
// the frames and indices don't encode the length of the original data. If maxInputSize
// is smaller than the original input size, decoded data will be trimmed to fit the maxInputSize.
//
// When asEval is false, the above description is the behavior. When asEval is true, the program proforms
// an additional FFT to transform back to evaluation representation. Under which case, maxInputSize
// must equal to the number of bytes after taking the IFFT, which has to be power of 2
func (g *Encoder) Decode(frames []Frame, indices []uint64, maxInputSize uint64, asEval bool) ([]byte, error) {
	if asEval {
		return g.DecodeAsEval(frames, indices, maxInputSize)
	} else {
		return g.DecodeAsCoeff(frames, indices, maxInputSize)
	}
}

func (g *Encoder) DecodeAsCoeff(frames []Frame, indices []uint64, maxInputSize uint64) ([]byte, error) {
	numSys := encoding.GetNumSys(maxInputSize, g.ChunkLength)

	if uint64(len(frames)) < numSys {
		return nil, errors.New("number of frame must be sufficient")
	}

	samples := make([]*fr.Element, g.NumEvaluations())
	// copy evals based on frame coeffs into samples
	for i, d := range indices {
		f := frames[i]
		e, err := GetLeadingCosetIndex(d, g.NumChunks)
		if err != nil {
			return nil, err
		}

		evals, err := g.GetInterpolationPolyEval(f.Coeffs, uint32(e))
		if err != nil {
			return nil, err
		}

		// Some pattern i butterfly swap. Find the leading coset, then increment by number of coset
		for j := uint64(0); j < g.ChunkLength; j++ {
			p := j*g.NumChunks + uint64(e)
			samples[p] = new(fr.Element)

			samples[p].Set(&evals[j])
		}
	}

	reconstructedData := make([]fr.Element, g.NumEvaluations())
	missingIndices := false
	for i, s := range samples {
		if s == nil {
			missingIndices = true
			break
		}
		reconstructedData[i] = *s
	}

	if missingIndices {
		var err error
		reconstructedData, err = g.Fs.RecoverPolyFromSamples(
			samples,
			g.Fs.ZeroPolyViaMultiplication,
		)
		if err != nil {
			return nil, err
		}
	}

	reconstructedPoly, err := g.Fs.FFT(reconstructedData, true)
	if err != nil {
		return nil, err
	}

	data := ToByteArray(reconstructedPoly, maxInputSize)

	return data, nil
}

// DecodeAsEval assumes the input has been IFFT transformed before sending to batcher to process
// this function should be used to reconstruct the data when api server from disperser contains IFFT
func (g *Encoder) DecodeAsEval(frames []Frame, indices []uint64, maxInputSize uint64) ([]byte, error) {
	if len(frames) == 0 {
		return nil, errors.New("number of frame must be greater than 1")
	}

	// get length for Fr, anything wit eval, 32 is the unit
	paddedInputLength := maxInputSize / 32

	fmt.Println("maxInputSize", maxInputSize, "g.ChunkLength", g.ChunkLength)

	// if number of data is less than padded data length
	if uint64(len(frames))*g.ChunkLength < paddedInputLength {
		// if the number of bytes is less than input size, then the number of points is insufficient
		// in the else case, the there is sufficient number of points, we can still recover the data
		if uint64(len(frames))*g.ChunkLength*encoding.NUMBER_FR_SECURITY_BYTES < maxInputSize {
			return nil, errors.New("number of frame must be sufficient")
		}
	}

	if len(frames) != len(indices) {
		return nil, fmt.Errorf("inconsistent number of frames and indices %d %d", len(frames), len(indices))
	}

	samples := make([]*fr.Element, g.NumEvaluations())
	// copy evals based on frame coeffs into samples
	for i, d := range indices {
		f := frames[i]
		e, err := GetLeadingCosetIndex(d, g.NumChunks)
		if err != nil {
			return nil, err
		}

		evals, err := g.GetInterpolationPolyEval(f.Coeffs, uint32(e))
		if err != nil {
			return nil, err
		}

		// Some pattern i butterfly swap. Find the leading coset, then increment by number of coset
		for j := uint64(0); j < g.ChunkLength; j++ {
			p := j*g.NumChunks + uint64(e)
			samples[p] = new(fr.Element)

			samples[p].Set(&evals[j])
		}
	}

	reconstructedData := make([]fr.Element, g.NumEvaluations())
	missingIndices := false
	for i, s := range samples {
		if s == nil {
			missingIndices = true
			break
		}
		reconstructedData[i] = *s
	}

	if missingIndices {
		var err error
		reconstructedData, err = g.Fs.RecoverPolyFromSamples(
			samples,
			g.Fs.ZeroPolyViaMultiplication,
		)
		if err != nil {
			return nil, err
		}
	}

	reconstructedPoly, err := g.Fs.FFT(reconstructedData, true)
	if err != nil {
		return nil, err
	}

	evalsFr, err := g.Fs.ConvertCoeffsToEvals(reconstructedPoly[:paddedInputLength])
	if err != nil {
		return nil, err
	}

	fmt.Println("maxInputSize", maxInputSize)

	data := ToByteArray(evalsFr, maxInputSize)

	return data, nil
}
