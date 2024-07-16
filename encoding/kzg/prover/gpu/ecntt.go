package gpu

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding/utils/gpu_utils"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	cr "github.com/ingonyama-zk/icicle/v2/wrappers/golang/cuda_runtime"
	icicle_bn254 "github.com/ingonyama-zk/icicle/v2/wrappers/golang/curves/bn254"
	ecntt "github.com/ingonyama-zk/icicle/v2/wrappers/golang/curves/bn254/ecntt"

	"github.com/ingonyama-zk/icicle/v2/wrappers/golang/core"
)

func (c *GpuComputeDevice) ECNttToGnark(batchPoints core.HostOrDeviceSlice, isInverse bool, totalSize int) ([]bn254.G1Affine, error) {
	output, err := c.ECNtt(batchPoints, isInverse, totalSize)
	if err != nil {
		return nil, err
	}

	// convert icicle projective to gnark affine
	gpuFFTBatch := gpu_utils.HostSliceIcicleProjectiveToGnarkAffine(output, int(c.NumWorker))

	return gpuFFTBatch, nil
}

func (c *GpuComputeDevice) ECNtt(batchPoints core.HostOrDeviceSlice, isInverse bool, totalSize int) (core.HostSlice[icicle_bn254.Projective], error) {
	output := make(core.HostSlice[icicle_bn254.Projective], totalSize)

	if isInverse {
		err := ecntt.ECNtt(batchPoints, core.KInverse, &c.NttCfg, output)
		if err.CudaErrorCode != cr.CudaSuccess || err.IcicleErrorCode != core.IcicleSuccess {
			return nil, fmt.Errorf("inverse ecntt failed")
		}
	} else {
		err := ecntt.ECNtt(batchPoints, core.KForward, &c.NttCfg, output)
		if err.CudaErrorCode != cr.CudaSuccess || err.IcicleErrorCode != core.IcicleSuccess {
			return nil, fmt.Errorf("forward ecntt failed")
		}
	}
	return output, nil
}
