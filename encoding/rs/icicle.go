//go:build icicle

package rs

import (
	"sync"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/icicle"
	rsicicle "github.com/Layr-Labs/eigenda/encoding/rs/icicle"
)

const (
	defaultNTTSize = 25 // Used for NTT setup in Icicle backend
)

func CreateIcicleBackendEncoder(e *Encoder, params encoding.EncodingParams, fs *fft.FFTSettings) (*ParametrizedEncoder, error) {
	icicleDevice, err := icicle.NewIcicleDevice(icicle.IcicleDeviceConfig{
		GPUEnable: e.Config.GPUEnable,
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
		RSEncoderComputer: &rsicicle.RsIcicleBackend{
			NttCfg:  icicleDevice.NttCfg,
			Device:  icicleDevice.Device,
			GpuLock: sync.Mutex{},
		},
	}, nil
}
