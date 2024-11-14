//go:build icicle

package icicle

import (
	"fmt"

	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/core"
	icicle_bn254 "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
	ecntt "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254/ecntt"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/runtime"
)

func (c *KzgMultiProofIcicleBackend) ECNttToGnarkOnDevice(batchPoints core.DeviceSlice, isInverse bool, totalSize int) (core.DeviceSlice, error) {
	output, err := c.ECNttOnDevice(batchPoints, isInverse, totalSize)
	if err != nil {
		return output, err
	}

	return output, nil
}

func (c *KzgMultiProofIcicleBackend) ECNttOnDevice(batchPoints core.DeviceSlice, isInverse bool, totalSize int) (core.DeviceSlice, error) {
	var p icicle_bn254.Projective
	var out core.DeviceSlice

	output, err := out.Malloc(p.Size(), totalSize)

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
