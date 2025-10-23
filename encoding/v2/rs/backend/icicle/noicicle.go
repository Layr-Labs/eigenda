//go:build !icicle

package icicle

import (
	"errors"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

type RSBackend struct{}

func (g *RSBackend) ExtendPolyEval(coeffs []fr.Element) ([]fr.Element, error) {
	// Not supported
	return nil, errors.New("icicle backend called without icicle build tag")
}

func BuildRSBackend(
	logger logging.Logger, enableGPU bool) (*RSBackend, error) {
	// Not supported
	return nil, errors.New("icicle backend called without icicle build tag")
}
