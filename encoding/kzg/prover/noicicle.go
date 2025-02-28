//go:build !icicle

package prover

import (
	"errors"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
)

func CreateIcicleBackendProver(p *Prover, params encoding.EncodingParams, fs *fft.FFTSettings, ks *kzg.KZGSettings) (*ParametrizedProver, error) {
	// Not supported
	return nil, errors.New("icicle backend called without icicle build tag")
}
