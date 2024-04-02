package fft

import (
	"fmt"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// Convert evaluation of a polynomial to coefficient of polynomial by taking Inverse FFT
func (fs *FFTSettings) ConvertEvalsToCoeffs(evals []fr.Element) ([]fr.Element, error) {
	if fs.MaxWidth < uint64(len(evals)) {
		return nil, fmt.Errorf("size of fft domain %d is insufficient to convert %d evaluations", fs.MaxWidth, len(evals))
	}

	coeffs, err := fs.FFT(evals, true)
	if err != nil {
		return nil, err
	}

	return coeffs, nil
}

func (fs *FFTSettings) ConvertCoeffsToEvals(coeffs []fr.Element) ([]fr.Element, error) {
	if fs.MaxWidth < uint64(len(coeffs)) {
		return nil, fmt.Errorf("size of fft domain %d is insufficient to convert %d coefficients", fs.MaxWidth, len(coeffs))
	}

	evals, err := fs.FFT(coeffs, false)
	if err != nil {
		return nil, err
	}

	return evals, nil
}
