//go:build !icicle

package icicle

import (
	"context"
	"errors"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

type RSBackend struct{}

func (g *RSBackend) ExtendPolyEvalV2(_ context.Context, coeffs []fr.Element) ([]fr.Element, error) {
	// Not supported
	return nil, errors.New("icicle backend called without icicle build tag")
}

func BuildRSBackend(
	logger logging.Logger, enableGPU bool, gpuConcurrentTasks int64) (*RSBackend, error) {
	// Not supported
	return nil, errors.New("icicle backend called without icicle build tag")
}
