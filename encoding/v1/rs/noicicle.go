//go:build !icicle

package rs

import (
	"errors"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
)

func CreateIcicleBackendEncoder(p *Encoder, params encoding.EncodingParams, fs *fft.FFTSettings) (*ParametrizedEncoder, error) {
	// Not supported
	return nil, errors.New("icicle backend called without icicle build tag")
}
