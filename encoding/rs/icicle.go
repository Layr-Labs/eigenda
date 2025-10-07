//go:build icicle

package rs

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding/icicle"
	rsicicle "github.com/Layr-Labs/eigenda/encoding/rs/icicle"
)

const (
	defaultNTTSize = 25 // Used for NTT setup in Icicle backend
)

func createIcicleBackend(enableGPU bool) (*rsicicle.RsIcicleBackend, error) {
	icicleDevice, err := icicle.NewIcicleDevice(icicle.IcicleDeviceConfig{
		GPUEnable: enableGPU,
		NTTSize:   defaultNTTSize,
		// No MSM setup needed for encoder
	})
	if err != nil {
		return nil, fmt.Errorf("new icicle device: %w", err)
	}

	return &rsicicle.RsIcicleBackend{
		NttCfg: icicleDevice.NttCfg,
		Device: icicleDevice.Device,
	}, nil

}
