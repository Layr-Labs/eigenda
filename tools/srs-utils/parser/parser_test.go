package parser_test

import (
	"testing"

	"github.com/Layr-Labs/eigenda/tools/srs-utils/parser"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestG1GeneratorPointsFromChallengeFile(t *testing.T) {
	// this file a truncated files from the original challenge_0085 file
	// this file contains only metadata and 4 g1 points, starting from
	// bn254 g1 generator
	filePath := "../resources/challenge_0085_with_4_g1_points"

	p := parser.Params{
		NumPoint:         4,
		NumTotalG1Points: 4,
		G1Size:           64,
		G2Size:           128,
	}

	p.SetG1StartBytePos(0)

	_, _, g1AffGen, _ := bn254.Generators()

	g1points, err := parser.ParseG1PointSection(filePath, p, 1)
	require.Nil(t, err)
	assert.Equal(t, len(g1points), 4)
	assert.Equal(t, g1points[0], g1AffGen)
}
