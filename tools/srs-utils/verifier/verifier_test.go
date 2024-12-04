package verifier_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Layr-Labs/eigenda/tools/srs-utils/verifier"
	"github.com/consensys/gnark-crypto/ecc/bn254"
)

func GetGeneratorPoints(n uint64) ([]bn254.G1Affine, []bn254.G2Affine) {

	secret := new(big.Int)
	secret.SetString("10", 10)

	g1SRS := make([]bn254.G1Affine, n)
	g2SRS := make([]bn254.G2Affine, n)

	multiplier := new(big.Int)
	multiplier.SetString("1", 10)

	_, _, _, g2Gen := bn254.Generators()

	for i := uint64(0); i < n; i++ {
		var s1Out bn254.G1Affine
		var s2Out bn254.G2Affine
		s1Out.ScalarMultiplicationBase(multiplier)
		s2Out.ScalarMultiplication(&g2Gen, multiplier)
		g1SRS[i] = s1Out
		g2SRS[i] = s2Out

		multiplier = multiplier.Mul(multiplier, secret)
	}

	return g1SRS, g2SRS
}

func TestCheckG1(t *testing.T) {
	numSRS := uint64(10)
	g1SRS, g2SRS := GetGeneratorPoints(numSRS)
	numWorker := 1
	results := make(chan error, numWorker)
	go verifier.G1CheckWorker(g1SRS, g2SRS, &g2SRS[0], &g2SRS[1], 0, 9, results)
	for i := 0; i < numWorker; i++ {
		err := <-results
		require.Nil(t, err)
	}
	close(results)

	results = make(chan error, numWorker)
	// corrupt a point
	g1SRS[numSRS/2] = g1SRS[numSRS/2-1]
	go verifier.G1CheckWorker(g1SRS, g2SRS, &g2SRS[0], &g2SRS[1], 0, 9, results)
	for i := 0; i < numWorker; i++ {
		err := <-results
		require.NotNil(t, err)
	}
	close(results)

}

func TestCheckG2(t *testing.T) {
	numSRS := uint64(10)
	g1SRS, g2SRS := GetGeneratorPoints(numSRS)

	numWorker := 1
	results := make(chan error, numWorker)
	go verifier.G2CheckWorker(g1SRS, g2SRS, &g1SRS[0], &g2SRS[0], 0, 10, results)
	for i := 0; i < numWorker; i++ {
		err := <-results
		require.Nil(t, err)
	}
	close(results)

	results = make(chan error, numWorker)
	// corrupt a point
	g1SRS[numSRS/2] = g1SRS[numSRS/2-1]
	go verifier.G2CheckWorker(g1SRS, g2SRS, &g1SRS[0], &g2SRS[0], 0, 10, results)
	for i := 0; i < numWorker; i++ {
		err := <-results
		require.NotNil(t, err)
	}
	close(results)

}
