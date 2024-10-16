//go:build gpu
// +build gpu

package gpu

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding/utils/gpu_utils"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/core"
	icicle_bn254 "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
	ecntt "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254/ecntt"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/runtime"
)

func (c *KzgGpuProofDevice) ECNttToGnarkOnDevice(batchPoints core.DeviceSlice, isInverse bool, totalSize int) (core.DeviceSlice, error) {
	output, err := c.ECNttOnDevice(batchPoints, isInverse, totalSize)
	if err != nil {
		return output, err
	}

	return output, nil
}

func (c *KzgGpuProofDevice) ECNttOnDevice(batchPoints core.DeviceSlice, isInverse bool, totalSize int) (core.DeviceSlice, error) {
	var p icicle_bn254.Projective
	var out core.DeviceSlice
	output, err := out.MallocAsync(totalSize*p.Size(), totalSize, *c.Stream)
	// output, err := out.Malloc(totalSize*p.Size(), p.Size())

	if err != runtime.Success {
		return out, fmt.Errorf("%v", "Allocating bytes on device for Projective results failed")
	}

	if isInverse {
		err := ecntt.ECNtt(batchPoints, core.KInverse, &c.NttCfg, output)
		if err != runtime.Success {
			return out, fmt.Errorf("inverse ecntt failed")
		}
	} else {
		err := ecntt.ECNtt(batchPoints, core.KForward, &c.NttCfg, output)
		if err != runtime.Success {
			return out, fmt.Errorf("forward ecntt failed")
		}
	}

	return output, nil
}

func (c *KzgGpuProofDevice) ECNttToGnark(batchPoints core.HostOrDeviceSlice, isInverse bool, totalSize int) ([]bn254.G1Affine, error) {
	output, err := c.ECNtt(batchPoints, isInverse, totalSize)
	if err != nil {
		return nil, err
	}

	// convert icicle projective to gnark affine
	gpuFFTBatch := gpu_utils.HostSliceIcicleProjectiveToGnarkAffine(output, int(c.NumWorker))

	return gpuFFTBatch, nil
}

func (c *KzgGpuProofDevice) ECNtt(batchPoints core.HostOrDeviceSlice, isInverse bool, totalSize int) (core.HostSlice[icicle_bn254.Projective], error) {
	output := make(core.HostSlice[icicle_bn254.Projective], totalSize)

	if isInverse {
		err := ecntt.ECNtt(batchPoints, core.KInverse, &c.NttCfg, output)
		if err != runtime.Success {
			return nil, fmt.Errorf("inverse ecntt failed")
		}
	} else {
		err := ecntt.ECNtt(batchPoints, core.KForward, &c.NttCfg, output)
		if err != runtime.Success {
			return nil, fmt.Errorf("forward ecntt failed")
		}
	}
	return output, nil
}
