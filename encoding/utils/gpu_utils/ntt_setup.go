//go:build gpu
// +build gpu

package gpu_utils

import (
	"log"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr/fft"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/core"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
	bn254_icicle "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254/ntt"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/runtime"
)

// batchSize is number of batches
func SetupNTT() (core.NTTConfig[[bn254_icicle.SCALAR_LIMBS]uint32], runtime.EIcicleError) {
	log.Println("Getting default NTT domain")
	cfg := core.GetDefaultNTTInitDomainConfig()

	// maximally possible
	exp := 25
	initDomain(exp, cfg)

	log.Println("Getting default NTT")
	cfgBn254 := ntt.GetDefaultNttConfig()

	cfgBn254.IsAsync = true
	cfgBn254.Ordering = core.KNN

	log.Println("Creating stream")
	streamBn254, err := runtime.CreateStream()
	if err != runtime.Success {
		return cfgBn254, err
	}

	cfgBn254.StreamHandle = streamBn254

	return cfgBn254, runtime.Success
}

func initDomain(largestTestSize int, cfg core.NTTInitDomainConfig) runtime.EIcicleError {
	rouMont, _ := fft.Generator(uint64(1 << largestTestSize))
	rou := rouMont.Bits()
	rouIcicle := bn254.ScalarField{}
	limbs := core.ConvertUint64ArrToUint32Arr(rou[:])

	rouIcicle.FromLimbs(limbs)
	e := ntt.InitDomain(rouIcicle, cfg)
	return e
}
