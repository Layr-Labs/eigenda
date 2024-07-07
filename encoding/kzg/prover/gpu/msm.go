package gpu

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding/utils/gpu_utils"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ingonyama-zk/icicle/v2/wrappers/golang/core"
	cr "github.com/ingonyama-zk/icicle/v2/wrappers/golang/cuda_runtime"
	icicle_bn254 "github.com/ingonyama-zk/icicle/v2/wrappers/golang/curves/bn254"
	icicle_bn254_msm "github.com/ingonyama-zk/icicle/v2/wrappers/golang/curves/bn254/msm"
)

// MsmBatch function supports batch across blobs
func (c *GpuComputer) MsmBatch(rowsFr [][]fr.Element, rowsG1 [][]bn254.G1Affine) ([]bn254.G1Affine, error) {
	msmCfg := icicle_bn254_msm.GetDefaultMSMConfig()
	rowsSfIcicle := make([]icicle_bn254.ScalarField, 0)
	rowsAffineIcicle := make([]icicle_bn254.Affine, 0)
	numBatchEle := len(rowsFr)

	// Prepare scalar fields
	for _, row := range rowsFr {
		rowsSfIcicle = append(rowsSfIcicle, gpu_utils.ConvertFrToScalarFieldsBytes(row)...)
	}
	rowsFrIcicleCopy := core.HostSliceFromElements[icicle_bn254.ScalarField](rowsSfIcicle)

	// Prepare icicle g1 affines
	for _, row := range rowsG1 {
		rowsAffineIcicle = append(rowsAffineIcicle, gpu_utils.BatchConvertGnarkAffineToIcicleAffine(row)...)
	}
	rowsG1IcicleCopy := core.HostSliceFromElements[icicle_bn254.Affine](rowsAffineIcicle)

	var p icicle_bn254.Projective
	var out core.DeviceSlice

	// prepare output
	_, err := out.Malloc(numBatchEle*p.Size(), p.Size())
	if err != cr.CudaSuccess {
		return nil, fmt.Errorf("allocating bytes on device for projective results failed")
	}

	err = icicle_bn254_msm.Msm(rowsFrIcicleCopy, rowsG1IcicleCopy, &msmCfg, out)
	if err != cr.CudaSuccess {
		return nil, fmt.Errorf("msm failed")
	}

	// move output out of device
	outHost := make(core.HostSlice[icicle_bn254.Projective], numBatchEle)
	outHost.CopyFromDevice(&out)
	out.Free()

	// convert data back to gnark format
	gnarkOuts := make([]bn254.G1Affine, numBatchEle)
	for i := 0; i < numBatchEle; i++ {
		gnarkOuts[i] = gpu_utils.IcicleProjectiveToGnarkAffine(outHost[i])
	}

	return gnarkOuts, nil
}
