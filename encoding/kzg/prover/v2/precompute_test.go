package prover_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/consensys/gnark-crypto/ecc/bn254"
)

func TestNewSRSTable_PreComputeWorks(t *testing.T) {
	harness := getTestHarness()

	kzgConfig := harness.proverV2KzgConfig
	kzgConfig.CacheDir = "./data/SRSTable"
	params := encoding.ParamsFromSysPar(harness.numSys, harness.numPar, uint64(len(harness.paddedGettysburgAddressBytes)))
	require.NotNil(t, params)

	s1, err := kzg.ReadG1Points(kzgConfig.G1Path, kzgConfig.SRSNumberToLoad, kzgConfig.NumWorker)
	require.Nil(t, err)
	require.NotNil(t, s1)

	_, err = kzg.ReadG2Points(kzgConfig.G2Path, kzgConfig.SRSNumberToLoad, kzgConfig.NumWorker)
	require.Nil(t, err)

	subTable1, err := prover.NewSRSTable(kzgConfig.CacheDir, s1, kzgConfig.NumWorker)
	require.Nil(t, err)
	require.NotNil(t, subTable1)

	fftPoints1, err := subTable1.GetSubTables(params.NumChunks, params.ChunkLength)
	require.Nil(t, err)
	require.NotNil(t, fftPoints1)

	subTable2, err := prover.NewSRSTable(kzgConfig.CacheDir, s1, kzgConfig.NumWorker)
	require.Nil(t, err)
	require.NotNil(t, subTable2)

	fftPoints2, err := subTable2.GetSubTables(params.NumChunks, params.ChunkLength)
	require.Nil(t, err)
	require.NotNil(t, fftPoints2)

	// Result of non precomputed GetSubTables should equal precomputed GetSubTables
	assert.Equal(t, fftPoints1, fftPoints2)
}

// This test reproduces the scenario where SRS_LOAD=2097152 and computing a subtable
// with the parameters (DimE=4, CosetSize=2097152) would cause a panic.
// The issue: m = numChunks*chunkLen - 1 = 4*2097152 - 1 = 8388607
// When j=0, k starts at m - cosetSize = 8388607 - 2097152 = 6291455
// Since 6291455 >= 2097152 (the length of our SRS), we get:
// panic: runtime error: index out of range [6291455] with length 2097152
func TestSRSTable_InsufficientSRSPoints_NoPanic(t *testing.T) {
	// Create a limited SRS with only 2097152 points
	limitedSRSSize := uint64(2097152)
	limitedSRS := make([]bn254.G1Affine, limitedSRSSize)

	// Initialize with some dummy points (doesn't matter what they are for this test)
	var generator bn254.G1Affine
	_, err := generator.X.SetString("1")
	require.NoError(t, err)
	_, err = generator.Y.SetString("2")
	require.NoError(t, err)
	for i := range limitedSRS {
		limitedSRS[i] = generator
	}

	// Create SRSTable with limited SRS points
	tempDir := t.TempDir()
	srsTable, err := prover.NewSRSTable(tempDir, limitedSRS, 1)
	require.NoError(t, err)

	// Try to create subtables with the following parameters
	numChunks := uint64(4)
	chunkLen := uint64(2097152)

	// This should return an error instead of panicking
	fftPoints, err := srsTable.GetSubTables(numChunks, chunkLen)

	assert.Error(t, err)
	assert.Nil(t, fftPoints)
	assert.Contains(t, err.Error(), "insufficient SRS points")
}
