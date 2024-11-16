//go:build icicle

package icicle

import (
	"sync"

	"github.com/Layr-Labs/eigenda/encoding/icicle"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/core"
	iciclebn254 "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254/ntt"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/runtime"
)

type RsIcicleComputeDevice struct {
	NttCfg core.NTTConfig[[iciclebn254.SCALAR_LIMBS]uint32]
	Device runtime.Device
}

// Encoding Reed Solomon using FFT
func (g *RsIcicleComputeDevice) ExtendPolyEval(coeffs []fr.Element) ([]fr.Element, error) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	var evals []fr.Element
	g.NttCfg.BatchSize = int32(1)
	runtime.RunOnDevice(&g.Device, func(args ...any) {
		defer wg.Done()
		scalarsSF := icicle.ConvertFrToScalarFieldsBytes(coeffs)
		scalars := core.HostSliceFromElements[iciclebn254.ScalarField](scalarsSF)
		outputDevice := make(core.HostSlice[iciclebn254.ScalarField], len(coeffs))
		ntt.Ntt(scalars, core.KForward, &g.NttCfg, outputDevice)
		evals = icicle.ConvertScalarFieldsToFrBytes(outputDevice)
	})

	wg.Wait()

	return evals, nil
}
