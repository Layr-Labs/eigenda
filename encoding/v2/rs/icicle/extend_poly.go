//go:build icicle

package icicle

import (
	"fmt"
	"sync"

	"github.com/Layr-Labs/eigenda/encoding/icicle"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/core"
	iciclebn254 "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254/ntt"
	icicleRuntime "github.com/ingonyama-zk/icicle/v3/wrappers/golang/runtime"
)

type RsIcicleBackend struct {
	NttCfg  core.NTTConfig[[iciclebn254.SCALAR_LIMBS]uint32]
	Device  icicleRuntime.Device
	GpuLock sync.Mutex
}

// Encoding Reed Solomon using FFT
func (g *RsIcicleBackend) ExtendPolyEval(coeffs []fr.Element) ([]fr.Element, error) {
	// Lock the GPU for operations
	g.GpuLock.Lock()
	defer g.GpuLock.Unlock()

	// Convert and prepare data
	g.NttCfg.BatchSize = int32(1)
	scalarsSF := icicle.ConvertFrToScalarFieldsBytes(coeffs)
	scalars := core.HostSliceFromElements[iciclebn254.ScalarField](scalarsSF)
	outputDevice := make(core.HostSlice[iciclebn254.ScalarField], len(coeffs))

	// Set device
	err := icicleRuntime.SetDevice(&g.Device)
	if err != icicleRuntime.Success {
		return nil, fmt.Errorf("failed to set device: %v", err.AsString())
	}

	// Perform NTT
	var icicleErr error
	wg := sync.WaitGroup{}
	wg.Add(1)
	icicleRuntime.RunOnDevice(&g.Device, func(args ...any) {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				icicleErr = fmt.Errorf("GPU operation panic: %v", r)
			}
		}()

		ntt.Ntt(scalars, core.KForward, &g.NttCfg, outputDevice)
	})

	wg.Wait()

	// Check if there was a panic
	if icicleErr != nil {
		return nil, icicleErr
	}

	evals := icicle.ConvertScalarFieldsToFrBytes(outputDevice)
	return evals, nil
}
