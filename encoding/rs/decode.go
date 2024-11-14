package rs

import (
	"errors"

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
func (e *Encoder) Decode(frames []Frame, indices []uint64, maxInputSize uint64, params encoding.EncodingParams) ([]byte, error) {
	// Get encoder
	g, err := e.GetRsEncoder(params)
	if err != nil {
		return nil, err
	}

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
