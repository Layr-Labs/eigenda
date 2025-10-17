//go:build icicle

package icicle

import (
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/fft"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/core"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254/ntt"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/runtime"
)

// SetupNTT initializes the NTT domain with the domain size of maxScale.
// It returns the NTT configuration and an error if the initialization fails.
func SetupNTT(maxScale uint8) (core.NTTConfig[[bn254.SCALAR_LIMBS]uint32], runtime.EIcicleError) {
	cfg := core.GetDefaultNTTInitDomainConfig()
	cfgBn254 := ntt.GetDefaultNttConfig()
	cfgBn254.IsAsync = true
	cfgBn254.Ordering = core.KNN

	err := initDomain(int(maxScale), cfg)
	if err != runtime.Success {
		return cfgBn254, err
	}

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
