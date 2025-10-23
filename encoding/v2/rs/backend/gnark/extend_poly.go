package gnark

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding/v2/fft"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

type RSBackend struct {
	Fs *fft.FFTSettings
}

func NewRSBackend(fs *fft.FFTSettings) *RSBackend {
	return &RSBackend{
		Fs: fs,
	}
}

// Encoding Reed Solomon using FFT
func (g *RSBackend) ExtendPolyEval(coeffs []fr.Element) ([]fr.Element, error) {
	evals, err := g.Fs.FFT(coeffs, false)
	if err != nil {
		return nil, fmt.Errorf("fft: %w", err)
	}

	return evals, nil
}
