package gpu

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding/utils/gpu_utils"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ingonyama-zk/icicle/v2/wrappers/golang/core"
	bn254_icicle "github.com/ingonyama-zk/icicle/v2/wrappers/golang/curves/bn254"
	bn254_icicle_ntt "github.com/ingonyama-zk/icicle/v2/wrappers/golang/curves/bn254/ntt"
)

func (c *GpuComputeDevice) NTT(batchFr [][]fr.Element) ([][]fr.Element, error) {
	if len(batchFr) == 0 {
		return nil, fmt.Errorf("input to NTT contains no blob")
	}

	numSymbol := len(batchFr[0])
	batchSize := len(batchFr)

	totalSize := numSymbol * batchSize

	// prepare scaler fields
	flattenBatchFr := make([]fr.Element, 0)
	for i := 0; i < len(batchFr); i++ {
		flattenBatchFr = append(flattenBatchFr, batchFr[i]...)
	}
	flattenBatchSf := gpu_utils.ConvertFrToScalarFieldsBytes(flattenBatchFr)
	scalarsCopy := core.HostSliceFromElements[bn254_icicle.ScalarField](flattenBatchSf)

	// run ntt
	output := make(core.HostSlice[bn254_icicle.ScalarField], totalSize)
	bn254_icicle_ntt.Ntt(scalarsCopy, core.KForward, &c.NttCfg, output)
	flattenBatchFrOutput := gpu_utils.ConvertScalarFieldsToFrBytes(output)

	// convert ntt output from icicle to gnark
	nttOutput := make([][]fr.Element, len(batchFr))
	for i := 0; i < len(batchFr); i++ {
		nttOutput[i] = flattenBatchFrOutput[i*numSymbol : (i+1)*numSymbol]
	}

	return nttOutput, nil
}
