//go:build icicle

package icicle

import (
	"fmt"
	"sync"

	"github.com/Layr-Labs/eigenda/encoding/icicle"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/core"
	iciclebn254 "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254/ntt"
	icicle_runtime "github.com/ingonyama-zk/icicle/v3/wrappers/golang/runtime"
)

const (
	defaultNTTSize = 25 // Used for NTT setup in Icicle backend
)

type RSBackend struct {
	NttCfg  core.NTTConfig[[iciclebn254.SCALAR_LIMBS]uint32]
	Device  icicle_runtime.Device
	GpuLock sync.Mutex
}

func BuildRSBackend(logger logging.Logger, enableGPU bool) (*RSBackend, error) {
	icicleDevice, err := icicle.NewIcicleDevice(icicle.IcicleDeviceConfig{
		Logger:    logger,
		GPUEnable: enableGPU,
		NTTSize:   defaultNTTSize,
		// No MSM setup needed for encoder
	})
	if err != nil {
		return nil, fmt.Errorf("setup icicle device: %w", err)
	}
	return &RSBackend{
		NttCfg: icicleDevice.NttCfg,
		Device: icicleDevice.Device,
	}, nil
}

// Encoding Reed Solomon using FFT
func (g *RSBackend) ExtendPolyEval(coeffs []fr.Element) ([]fr.Element, error) {
	// coeffs will be moved to device memory inside Ntt function,
	// and the result copied back into outputEvals.
	coeffsSlice := core.HostSliceFromElements(coeffs)
	outputEvals := make(core.HostSlice[fr.Element], len(coeffs))
	var icicleErr error

	wg := sync.WaitGroup{}
	wg.Add(1)
	icicle_runtime.RunOnDevice(&g.Device, func(args ...any) {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				icicleErr = fmt.Errorf("GPU operation panic: %v", r)
			}
		}()

		cfg := g.NttCfg
		// We just perform the NTT synchronously here; we have nothing to do while waiting.
		cfg.IsAsync = false
		cfg.BatchSize = int32(1)
		ntt.Ntt(coeffsSlice, core.KForward, &cfg, outputEvals)
	})
	wg.Wait()

	// Check if there was a panic
	if icicleErr != nil {
		return nil, icicleErr
	}
	return outputEvals, nil
}
