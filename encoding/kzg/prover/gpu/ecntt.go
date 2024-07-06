package gpu

import (
	"fmt"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	cr "github.com/ingonyama-zk/icicle/v2/wrappers/golang/cuda_runtime"
	icicle_bn254 "github.com/ingonyama-zk/icicle/v2/wrappers/golang/curves/bn254"
	ecntt "github.com/ingonyama-zk/icicle/v2/wrappers/golang/curves/bn254/ecntt"

	"github.com/ingonyama-zk/icicle/v2/wrappers/golang/core"
)

func (c *GpuComputer) ECNtt(batchPoints []bn254.G1Affine, isInverse bool) ([]bn254.G1Affine, error) {
	totalNumSym := len(batchPoints)

	// convert gnark affine to icicle projective on slice
	pointsIcileProjective := BatchConvertGnarkAffineToIcicleProjective(batchPoints)
	pointsCopy := core.HostSliceFromElements[icicle_bn254.Projective](pointsIcileProjective)

	output := make(core.HostSlice[icicle_bn254.Projective], int(totalNumSym))

	// compute
	if isInverse {
		err := ecntt.ECNtt(pointsCopy, core.KInverse, &c.cfg, output)
		if err.CudaErrorCode != cr.CudaSuccess || err.IcicleErrorCode != core.IcicleSuccess {
			return nil, fmt.Errorf("inverse ecntt failed")
		}
	} else {
		err := ecntt.ECNtt(pointsCopy, core.KForward, &c.cfg, output)
		if err.CudaErrorCode != cr.CudaSuccess || err.IcicleErrorCode != core.IcicleSuccess {
			return nil, fmt.Errorf("forward ecntt failed")
		}
	}

	// convert icicle projective to gnark affine
	gpuFFTBatch := make([]bn254.G1Affine, len(batchPoints))
	for j := 0; j < totalNumSym; j++ {
		gpuFFTBatch[j] = IcicleProjectiveToGnarkAffine(output[j])
	}

	return gpuFFTBatch, nil
}
