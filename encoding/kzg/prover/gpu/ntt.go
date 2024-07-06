package gpu

import (
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/fft"
	"github.com/ingonyama-zk/icicle/v2/wrappers/golang/core"
	bn254_icicle "github.com/ingonyama-zk/icicle/v2/wrappers/golang/curves/bn254"
	bn254_icicle_ntt "github.com/ingonyama-zk/icicle/v2/wrappers/golang/curves/bn254/ntt"
)

// batchSize is number of batches
func SetupNTT() core.NTTConfig[[bn254_icicle.SCALAR_LIMBS]uint32] {
	cfg := bn254_icicle_ntt.GetDefaultNttConfig()

	cfg.Ordering = core.KNN
	cfg.NttAlgorithm = core.Radix2

	// batchSize will change later when used
	cfg.BatchSize = int32(1)
	cfg.NttAlgorithm = core.Radix2

	// maximally possible
	exp := 28

	rouMont, _ := fft.Generator(uint64(1 << exp))
	rou := rouMont.Bits()
	rouIcicle := bn254_icicle.ScalarField{}
	limbs := core.ConvertUint64ArrToUint32Arr(rou[:])
	rouIcicle.FromLimbs(limbs)
	bn254_icicle_ntt.InitDomain(rouIcicle, cfg.Ctx, false)

	return cfg
}

func (c *GpuComputer) NTT(batchFr [][]fr.Element) ([][]fr.Element, error) {
	numSymbol := len(batchFr[0])
	batchSize := len(batchFr)

	totalSize := numSymbol * batchSize

	// prepare scaler fields
	flattenBatchFr := make([]fr.Element, 0)
	for i := 0; i < len(batchFr); i++ {
		flattenBatchFr = append(flattenBatchFr, batchFr[i]...)
	}
	flattenBatchSf := ConvertFrToScalarFieldsBytes(flattenBatchFr)
	scalarsCopy := core.HostSliceFromElements[bn254_icicle.ScalarField](flattenBatchSf)

	// run ntt
	output := make(core.HostSlice[bn254_icicle.ScalarField], totalSize)
	bn254_icicle_ntt.Ntt(scalarsCopy, core.KForward, &c.cfg, output)
	flattenBatchFrOutput := ConvertScalarFieldsToFrBytes(output)

	// convert ntt output from icicle to gnark
	nttOutput := make([][]fr.Element, len(batchFr))
	for i := 0; i < len(batchFr); i++ {
		nttOutput[i] = flattenBatchFrOutput[i*numSymbol : (i+1)*numSymbol]
	}

	return nttOutput, nil
}
