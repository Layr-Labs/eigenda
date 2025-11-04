//go:build icicle

package icicle

import (
	"context"
	"fmt"
	"sync"

	_ "github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/icicle"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/core"
	iciclebn254 "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254/ntt"
	icicle_runtime "github.com/ingonyama-zk/icicle/v3/wrappers/golang/runtime"
	"golang.org/x/sync/semaphore"
)

const (
	defaultNTTSize = 25 // Used for NTT setup in Icicle backend
)

type RSBackend struct {
	NttCfg core.NTTConfig[[iciclebn254.SCALAR_LIMBS]uint32]
	Device icicle_runtime.Device
	// request-weighted semaphore.
	// See [encoding.Config.GPUConcurrentFrameGenerationDangerous] for more details.
	GpuSemaphore *semaphore.Weighted
}

func BuildRSBackend(logger logging.Logger, enableGPU bool, gpuConcurrentEncodings int64) (*RSBackend, error) {
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
		NttCfg:       icicleDevice.NttCfg,
		Device:       icicleDevice.Device,
		GpuSemaphore: semaphore.NewWeighted(gpuConcurrentEncodings),
	}, nil
}

// Encoding Reed Solomon using FFT
func (g *RSBackend) ExtendPolyEvalV2(ctx context.Context, coeffs []fr.Element) ([]fr.Element, error) {
	// We acquire a semaphore here to avoid too many concurrent NTT calls.
	// This is a very unideal and coarse grain solution, but unfortunately
	// icicle doesn't have nice backpressure, and the GPU kernel just panics if RAM is exhausted.
	// In its current implementation, icicle's NTT kernel takes RAM = input+output size.
	// We could use a finer-grained semaphore that calculates the RAM usage per request,
	// but this would feel very hardcoded and hardware dependent (although we can request RAM available on the device
	// dynamically using icicle APIs). For now opting to keep this simple.
	// TODO(samlaf): rethink this approach.
	err := g.GpuSemaphore.Acquire(ctx, 1)
	if err != nil {
		return nil, fmt.Errorf("acquiring GPU semaphore: %w", err)
	}
	defer g.GpuSemaphore.Release(1)

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
