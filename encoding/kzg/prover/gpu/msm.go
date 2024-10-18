//go:build gpu
// +build gpu

package gpu

import (
	"fmt"

	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/core"
	icicle_bn254 "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254/msm"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/runtime"
)

// MsmBatchOnDevice function supports batch across blobs.
// totalSize is the number of output points, which equals to numPoly * 2 * dimE , dimE is number of chunks
func (c *KzgGpuProofDevice) MsmBatchOnDevice(rowsFrIcicleCopy core.DeviceSlice, rowsG1Icicle []icicle_bn254.Affine, totalSize int, device runtime.Device) (core.DeviceSlice, error) {
	rowsG1IcicleCopy := core.HostSliceFromElements[icicle_bn254.Affine](rowsG1Icicle)

	var p icicle_bn254.Projective
	var out core.DeviceSlice

	_, err := out.Malloc(p.Size(), totalSize)

	if err != runtime.Success {
		fmt.Println(err)
		return out, fmt.Errorf("%v", "Allocating bytes on device for Projective results failed")
	}

	err = msm.Msm(rowsFrIcicleCopy, rowsG1IcicleCopy, &c.MsmCfg, out)
	if err != runtime.Success {
		return out, fmt.Errorf("%v", "Msm failed")
	}

	return out, nil
}
