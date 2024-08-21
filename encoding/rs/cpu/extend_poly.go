package cpu

import (
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

type CpuComputeDevice struct {
	Fs *fft.FFTSettings

	encoding.EncodingParams
}

// Encoding Reed Solomon using FFT
func (g *CpuComputeDevice) ExtendPolyEval(coeffs []fr.Element) ([]fr.Element, error) {
	evals, err := g.Fs.FFT(coeffs, false)
	if err != nil {
		return nil, err
	}

	return evals, nil
}
