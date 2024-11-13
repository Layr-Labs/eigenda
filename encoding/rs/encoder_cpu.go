//go:build !gpu
// +build !gpu

package rs

import (
	"math"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/rs/cpu"

	_ "go.uber.org/automaxprocs"
)

func (g *Encoder) newEncoder(params encoding.EncodingParams) (*ParametrizedEncoder, error) {
	err := params.Validate()
	if err != nil {
		return nil, err
	}

	n := uint8(math.Log2(float64(params.NumEvaluations())))
	if params.ChunkLength == 1 {
		n = uint8(math.Log2(float64(2 * params.NumChunks)))
	}
	fs := fft.NewFFTSettings(n)

	// Set RS CPU computer
	RsComputeDevice := &cpu.RsCpuComputeDevice{
		Fs: fs,
	}

	return &ParametrizedEncoder{
		Config:            g.Config,
		EncodingParams:    params,
		Fs:                fs,
		RSEncoderComputer: RsComputeDevice,
	}, nil
}
