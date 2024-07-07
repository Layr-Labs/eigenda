package gpu

import (
	"fmt"

	"github.com/ingonyama-zk/icicle/v2/wrappers/golang/core"
	cr "github.com/ingonyama-zk/icicle/v2/wrappers/golang/cuda_runtime"
	bn254_icicle "github.com/ingonyama-zk/icicle/v2/wrappers/golang/curves/bn254"
	"github.com/ingonyama-zk/icicle/v2/wrappers/golang/curves/bn254/vecOps"
)

// numRow and numCol describes input dimension
func Transpose(coeffStoreFFT core.HostSlice[bn254_icicle.ScalarField], l, numPoly, numCol int) (core.HostSlice[bn254_icicle.ScalarField], error) {
	totalSize := l * numPoly * numCol
	ctx, err := cr.GetDefaultDeviceContext()
	if err != cr.CudaSuccess {
		return nil, fmt.Errorf("allocating bytes on device for projective results failed")
	}

	transposedNTTOutput := make(core.HostSlice[bn254_icicle.ScalarField], totalSize)

	for i := 0; i < numPoly; i++ {
		vecOps.TransposeMatrix(coeffStoreFFT[i*l*numCol:(i+1)*l*numCol], transposedNTTOutput[i*l*numCol:(i+1)*l*numCol], l, numCol, ctx, false, false)
	}

	return transposedNTTOutput, nil
}
