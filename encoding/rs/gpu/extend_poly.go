package gpu

import (
	"sync"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/utils/gpu_utils"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ingonyama-zk/icicle/v2/wrappers/golang/core"
	bn254_icicle "github.com/ingonyama-zk/icicle/v2/wrappers/golang/curves/bn254"
	"github.com/ingonyama-zk/icicle/v2/wrappers/golang/curves/bn254/ntt"
)

type GpuComputeDevice struct {
	Fs *fft.FFTSettings
	encoding.EncodingParams
	NttCfg  core.NTTConfig[[bn254_icicle.SCALAR_LIMBS]uint32]
	GpuLock *sync.Mutex
}

// Encoding Reed Solomon using FFT
func (g *GpuComputeDevice) ExtendPolyEval(coeffs []fr.Element) ([]fr.Element, error) {

	g.GpuLock.Lock()
	defer g.GpuLock.Unlock()

	scalarsSF := gpu_utils.ConvertFrToScalarFieldsBytes(coeffs)

	scalars := core.HostSliceFromElements[bn254_icicle.ScalarField](scalarsSF)

	outputDevice := make(core.HostSlice[bn254_icicle.ScalarField], len(coeffs))

	ntt.Ntt(scalars, core.KForward, &g.NttCfg, outputDevice)

	evals := gpu_utils.ConvertScalarFieldsToFrBytes(outputDevice)

	return evals, nil
}
