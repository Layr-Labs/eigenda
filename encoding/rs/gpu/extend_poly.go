//go:build gpu
// +build gpu

package gpu

import (
	"sync"

	"github.com/Layr-Labs/eigenda/encoding/utils/gpu_utils"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/core"
	icicle_bn254 "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254/ntt"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/runtime"
)

type RsGpuComputeDevice struct {
	NttCfg  core.NTTConfig[[icicle_bn254.SCALAR_LIMBS]uint32]
	GpuLock *sync.Mutex
	Device  runtime.Device
}

// Encoding Reed Solomon using FFT
func (g *RsGpuComputeDevice) ExtendPolyEval(coeffs []fr.Element) ([]fr.Element, error) {
	g.GpuLock.Lock()
	defer g.GpuLock.Unlock()

	wg := sync.WaitGroup{}
	wg.Add(1)

	var evals []fr.Element
	g.NttCfg.BatchSize = int32(1)
	runtime.RunOnDevice(&g.Device, func(args ...any) {
		defer wg.Done()
		scalarsSF := gpu_utils.ConvertFrToScalarFieldsBytes(coeffs)
		scalars := core.HostSliceFromElements[icicle_bn254.ScalarField](scalarsSF)
		outputDevice := make(core.HostSlice[icicle_bn254.ScalarField], len(coeffs))
		ntt.Ntt(scalars, core.KForward, &g.NttCfg, outputDevice)
		evals = gpu_utils.ConvertScalarFieldsToFrBytes(outputDevice)
	})

	wg.Wait()

	return evals, nil
}
