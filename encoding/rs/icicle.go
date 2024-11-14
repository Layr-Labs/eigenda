//go:build icicle

package rs

import (
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/icicle"
	rs_icicle "github.com/Layr-Labs/eigenda/encoding/rs/icicle"
)

const (
	defaultNTTSize = 25 // Used for NTT setup in Icicle backend
)

func CreateIcicleBackendEncoder(e *Encoder, params encoding.EncodingParams, fs *fft.FFTSettings) (*ParametrizedEncoder, error) {
	icicleDevice, err := icicle.NewIcicleDevice(icicle.IcicleDeviceConfig{
		EnableGPU: e.Config.EnableGPU,
		NTTSize:   defaultNTTSize,
		// No MSM setup needed for encoder
	})
	if err != nil {
		return nil, err
	}

	return &ParametrizedEncoder{
		Config:         e.Config,
		EncodingParams: params,
		Fs:             fs,
		RSEncoderComputer: &rs_icicle.RsIcicleComputeDevice{
			NttCfg: icicleDevice.NttCfg,
			Device: icicleDevice.Device,
		},
	}, nil
}
