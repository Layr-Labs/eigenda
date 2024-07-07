package gpu_utils

import (
	"github.com/consensys/gnark-crypto/ecc/bn254"
	bn254_icicle "github.com/ingonyama-zk/icicle/v2/wrappers/golang/curves/bn254"
)

func SetupMsm(rowsG1 [][]bn254.G1Affine) []bn254_icicle.Affine {
	rowsG1Icicle := make([]bn254_icicle.Affine, 0)

	for _, row := range rowsG1 {
		rowsG1Icicle = append(rowsG1Icicle, BatchConvertGnarkAffineToIcicleAffine(row)...)
	}
	return rowsG1Icicle
}
