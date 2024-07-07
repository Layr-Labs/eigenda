package gpu

import (
	"errors"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/utils/gpu_utils"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ingonyama-zk/icicle/v2/wrappers/golang/core"
	bn254_icicle "github.com/ingonyama-zk/icicle/v2/wrappers/golang/curves/bn254"
	"github.com/ingonyama-zk/icicle/v2/wrappers/golang/curves/bn254/ntt"
)

type GpuComputer struct {
	Fs *fft.FFTSettings
	encoding.EncodingParams
	NttCfg core.NTTConfig[[bn254_icicle.SCALAR_LIMBS]uint32]
}

// Encoding Reed Solomon using FFT
func (g *GpuComputer) ExtendPolyEval(coeffs []fr.Element) ([]fr.Element, error) {

	if len(coeffs) > int(g.NumEvaluations()) {
		return nil, errors.New("the provided encoding parameters are not sufficient for the size of the data input")
	}

	pdCoeffs := make([]fr.Element, g.NumEvaluations())
	for i := 0; i < len(coeffs); i++ {
		pdCoeffs[i].Set(&coeffs[i])
	}
	for i := len(coeffs); i < len(pdCoeffs); i++ {
		pdCoeffs[i].SetZero()
	}

	scalarsSF := gpu_utils.ConvertFrToScalarFieldsBytes(pdCoeffs)

	scalars := core.HostSliceFromElements[bn254_icicle.ScalarField](scalarsSF)

	outputDevice := make(core.HostSlice[bn254_icicle.ScalarField], len(pdCoeffs))

	ntt.Ntt(scalars, core.KForward, &g.NttCfg, outputDevice)

	evals := gpu_utils.ConvertScalarFieldsToFrBytes(outputDevice)

	return evals, nil
}
