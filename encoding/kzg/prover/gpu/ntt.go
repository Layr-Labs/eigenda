package gpu

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding/utils/gpu_utils"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ingonyama-zk/icicle/v2/wrappers/golang/core"
	bn254_icicle "github.com/ingonyama-zk/icicle/v2/wrappers/golang/curves/bn254"
	bn254_icicle_ntt "github.com/ingonyama-zk/icicle/v2/wrappers/golang/curves/bn254/ntt"
)

func (c *GpuComputeDevice) NTT(batchFr [][]fr.Element) (core.HostSlice[bn254_icicle.ScalarField], error) {
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
	flattenBatchSf := gpu_utils.ConvertFrToScalarFieldsBytesThread(flattenBatchFr, int(c.NumWorker))
	scalarsCopy := core.HostSliceFromElements[bn254_icicle.ScalarField](flattenBatchSf)

	// run ntt
	output := make(core.HostSlice[bn254_icicle.ScalarField], totalSize)
	bn254_icicle_ntt.Ntt(scalarsCopy, core.KForward, &c.NttCfg, output)

	return output, nil
}
