//go:build icicle

package icicle

import (
	"fmt"
	"sync"

	"github.com/Layr-Labs/eigenda/encoding/icicle"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/core"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254/ntt"
	icicle_runtime "github.com/ingonyama-zk/icicle/v3/wrappers/golang/runtime"
	"golang.org/x/sync/semaphore"
)

const (
	defaultNTTSize = 25 // Used for NTT setup in Icicle backend
)

type RSBackend struct {
	Device icicle_runtime.Device
	// request-weighted semaphore.
	// See [encoding.Config.GPUConcurrentFrameGenerationDangerous] for more details.
	GpuSemaphore *semaphore.Weighted
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
		Device:       icicleDevice.Device,
		GpuSemaphore: semaphore.NewWeighted(gpuConcurrentEncodings),
	}, nil
}

// Encoding Reed Solomon using FFT
func (g *RSBackend) ExtendPolyEval(coeffs []fr.Element) ([]fr.Element, error) {
	// We lock the GPU here to avoid concurrent NTT calls.
	// This is WAY too conservative, and we could achieve way better throughput by submitting multiple
	// NTT jobs in parallel to the GPU, but the issue is that icicle doesn't have nice backpressure,
	// and the GPU kernel just panics if too many jobs are submitted at once.
	// TODO(samlaf): add some kind of job queue, backpressure, exponential backoff,
	// or whatever to allow maximizing GPU utilization.
	// See https://github.com/Layr-Labs/eigenda/pull/2214 for more details.
	g.GpuLock.Lock()
	defer g.GpuLock.Unlock()

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

		// Create a new stream for this operation to allow concurrent GPU operations
		// without interference. Each stream can execute independently.
		stream, err := icicle_runtime.CreateStream()
		if err != icicle_runtime.Success {
			icicleErr = fmt.Errorf("failed to create stream: %v", err.AsString())
			return
		}
		defer func() {
			// Synchronize stream to ensure all GPU operations complete before cleanup
			syncErr := icicle_runtime.SynchronizeStream(stream)
			if syncErr != icicle_runtime.Success && icicleErr == nil {
				icicleErr = fmt.Errorf("stream synchronization failed: %v", syncErr.AsString())
			}
			icicle_runtime.DestroyStream(stream)
		}()

		// Create NTT config for this operation
		cfg := ntt.GetDefaultNttConfig()
		cfg.IsAsync = true
		cfg.StreamHandle = stream
		nttErr := ntt.Ntt(coeffsSlice, core.KForward, &cfg, outputEvals)
		if nttErr != icicle_runtime.Success {
			icicleErr = fmt.Errorf("NTT operation failed: %v", nttErr.AsString())
			return
		}
	})
	wg.Wait()

	// Check if there was a panic
	if icicleErr != nil {
		return nil, icicleErr
	}
	return outputEvals, nil
}
