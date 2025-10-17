//go:build icicle

package icicle

import (
	"fmt"

	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/core"
	iciclebn254 "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254/msm"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/runtime"
)

// MsmBatchOnDevice function supports batch across blobs.
// totalSize is the number of output points, which equals to numPoly * 2 * dimE , dimE is number of chunks
func (c *KzgMultiProofIcicleBackend) MsmBatchOnDevice(rowsFrIcicleCopy core.DeviceSlice, rowsG1Icicle []iciclebn254.Affine, totalSize int) (core.DeviceSlice, error) {
	rowsG1IcicleCopy := core.HostSliceFromElements[iciclebn254.Affine](rowsG1Icicle)

	var p iciclebn254.Projective
	var out core.DeviceSlice

	_, err := out.Malloc(p.Size(), totalSize)
	if err != runtime.Success {
		return out, fmt.Errorf("allocating bytes on device failed: %v", err.AsString())
	}

	err = msm.Msm(rowsFrIcicleCopy, rowsG1IcicleCopy, &c.MsmCfg, out)
	if err != runtime.Success {
		return out, fmt.Errorf("msm error: %v", err.AsString())
	}

	return out, nil
}
