package gnark

import (
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

type RsGnarkBackend struct {
	Fs *fft.FFTSettings
}

// Encoding Reed Solomon using FFT
func (g *RsGnarkBackend) ExtendPolyEval(coeffs []fr.Element) ([]fr.Element, error) {
	evals, err := g.Fs.FFT(coeffs, false)
	if err != nil {
		return nil, err
	}

	return evals, nil
}
