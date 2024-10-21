//go:build gpu
// +build gpu

package gpu

import (
	"sync"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/utils/gpu_utils"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/core"
	icicle_bn254 "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254/ntt"
)

type GpuComputeDevice struct {
	NttCfg  core.NTTConfig[[icicle_bn254.SCALAR_LIMBS]uint32]
	GpuLock *sync.Mutex

	encoding.EncodingParams
}

// Encoding Reed Solomon using FFT
func (g *GpuComputeDevice) ExtendPolyEval(coeffs []fr.Element) ([]fr.Element, error) {

	g.GpuLock.Lock()
	defer g.GpuLock.Unlock()

	scalarsSF := gpu_utils.ConvertFrToScalarFieldsBytes(coeffs)

	scalars := core.HostSliceFromElements[icicle_bn254.ScalarField](scalarsSF)

	outputDevice := make(core.HostSlice[icicle_bn254.ScalarField], len(coeffs))

	ntt.Ntt(scalars, core.KForward, &g.NttCfg, outputDevice)

	evals := gpu_utils.ConvertScalarFieldsToFrBytes(outputDevice)

	return evals, nil
}
