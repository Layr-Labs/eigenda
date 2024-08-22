//go:build gpu
// +build gpu

package gpu_utils

import (
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/ingonyama-zk/icicle/v2/wrappers/golang/core"
	bn254_icicle "github.com/ingonyama-zk/icicle/v2/wrappers/golang/curves/bn254"
	icicle_bn254_msm "github.com/ingonyama-zk/icicle/v2/wrappers/golang/curves/bn254/msm"
)

func SetupMsm(rowsG1 [][]bn254.G1Affine, srsG1 []bn254.G1Affine) ([]bn254_icicle.Affine, []bn254_icicle.Affine, core.MSMConfig) {
	rowsG1Icicle := make([]bn254_icicle.Affine, 0)

	for _, row := range rowsG1 {
		rowsG1Icicle = append(rowsG1Icicle, BatchConvertGnarkAffineToIcicleAffine(row)...)
	}

	srsG1Icicle := BatchConvertGnarkAffineToIcicleAffine(srsG1)

	msmCfg := icicle_bn254_msm.GetDefaultMSMConfig()

	return rowsG1Icicle, srsG1Icicle, msmCfg
}
