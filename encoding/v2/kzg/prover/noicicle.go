//go:build !icicle

package prover

import (
	"errors"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/v2/fft"
)

func CreateIcicleBackendProver(
	p *Prover, params encoding.EncodingParams, fs *fft.FFTSettings,
) (*ParametrizedProver, error) {
	// Not supported
	return nil, errors.New("icicle backend called without icicle build tag")
}
