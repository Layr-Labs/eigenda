//go:build icicle

package icicle

import (
	"fmt"
	"runtime"

	"github.com/Layr-Labs/eigenda/encoding/icicle"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/core"
	iciclebn254 "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254/ntt"
	icicleRuntime "github.com/ingonyama-zk/icicle/v3/wrappers/golang/runtime"
)

type RsIcicleBackend struct {
	NttCfg core.NTTConfig[[iciclebn254.SCALAR_LIMBS]uint32]
	Device icicleRuntime.Device
}

// Encoding Reed Solomon using FFT
func (g *RsIcicleBackend) ExtendPolyEval(coeffs []fr.Element) ([]fr.Element, error) {
	// Lock the OS thread for GPU operations
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Set device
	originalDevice, err := icicleRuntime.GetActiveDevice()
	if err != icicleRuntime.Success {
		return nil, fmt.Errorf("failed to get active device: %v", err)
	}
	defer icicleRuntime.SetDevice(originalDevice)

	if err := icicleRuntime.SetDevice(&g.Device); err != icicleRuntime.Success {
		return nil, fmt.Errorf("failed to set device: %v", err)
	}

	// Convert and prepare data
	g.NttCfg.BatchSize = int32(1)
	scalarsSF := icicle.ConvertFrToScalarFieldsBytes(coeffs)
	scalars := core.HostSliceFromElements[iciclebn254.ScalarField](scalarsSF)
	outputDevice := make(core.HostSlice[iciclebn254.ScalarField], len(coeffs))

	// Perform NTT
	ntt.Ntt(scalars, core.KForward, &g.NttCfg, outputDevice)

	// Convert back to fr.Element
	evals := icicle.ConvertScalarFieldsToFrBytes(outputDevice)

	return evals, nil
}
