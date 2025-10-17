//go:build icicle

package icicle

import (
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/core"
	iciclebn254 "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/runtime"
)

// SetupMsmG1 initializes the MSM configuration for G1 points.
func SetupMsmG1(rowsG1 [][]bn254.G1Affine, srsG1 []bn254.G1Affine) ([]iciclebn254.Affine, []iciclebn254.Affine, core.MSMConfig, runtime.EIcicleError) {
	// Calculate total length needed for rowsG1Icicle
	totalLen := 0
	for _, row := range rowsG1 {
		totalLen += len(row)
	}

	// Pre-allocate slice with exact capacity needed
	rowsG1Icicle := make([]iciclebn254.Affine, totalLen)

	currentIdx := 0
	for _, row := range rowsG1 {
		converted := BatchConvertGnarkAffineToIcicleAffine(row)
		copy(rowsG1Icicle[currentIdx:], converted)
		currentIdx += len(row)
	}

	srsG1Icicle := BatchConvertGnarkAffineToIcicleAffine(srsG1)
	cfgBn254 := core.GetDefaultMSMConfig()
	cfgBn254.IsAsync = true

	streamBn254, err := runtime.CreateStream()
	if err != runtime.Success {
		return nil, nil, cfgBn254, err
	}

	cfgBn254.StreamHandle = streamBn254
	return rowsG1Icicle, srsG1Icicle, cfgBn254, runtime.Success
}
