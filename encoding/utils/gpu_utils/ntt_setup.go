//go:build gpu
// +build gpu

package gpu_utils

import (
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/fft"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/core"
	bn254_icicle "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
	bn254_icicle_ntt "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254/ntt"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/runtime"
)

// batchSize is number of batches
func SetupNTT() (core.NTTConfig[[bn254_icicle.SCALAR_LIMBS]uint32], runtime.EIcicleError) {
	cfgBn254 := bn254_icicle_ntt.GetDefaultNttConfig()
	cfgBn254.IsAsync = true
	cfgBn254.Ordering = core.KNN

	streamBn254, err := runtime.CreateStream()
	if err != runtime.Success {
		return cfgBn254, err
	}

	cfgBn254.StreamHandle = streamBn254

	cfg := core.GetDefaultNTTInitDomainConfig()

	// maximally possible
	exp := 28

	rouMont, _ := fft.Generator(uint64(1 << exp))
	rou := rouMont.Bits()
	rouIcicle := bn254_icicle.ScalarField{}
	limbs := core.ConvertUint64ArrToUint32Arr(rou[:])
	rouIcicle.FromLimbs(limbs)
	bn254_icicle_ntt.InitDomain(rouIcicle, cfg)

	return cfgBn254, runtime.Success
}
