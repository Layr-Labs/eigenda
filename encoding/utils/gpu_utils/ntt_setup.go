package gpu_utils

import (
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
