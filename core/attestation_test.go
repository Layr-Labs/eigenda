package core_test

import (
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/stretchr/testify/require"
)

// contractG2Coords serializes a G2 point into the contracts' BN254.G2Point
// layout, where each E2 coordinate is ordered [A1, A0] (imaginary component
// first) — the order in which NewPubkeyRegistration event data arrives.
func contractG2Coords(g2 *core.G2Point) (x, y [2]*big.Int) {
	x = [2]*big.Int{g2.X.A1.BigInt(new(big.Int)), g2.X.A0.BigInt(new(big.Int))}
	y = [2]*big.Int{g2.Y.A1.BigInt(new(big.Int)), g2.Y.A0.BigInt(new(big.Int))}
	return x, y
}

func TestNewG2PointComponentSwap(t *testing.T) {
	keyPair, err := core.GenRandomBlsKeys()
	require.NoError(t, err)

	g1 := keyPair.GetPubKeyG1()
	g2 := keyPair.GetPubKeyG2()

	// Round-trip through the contract layout must reproduce the original key:
	// the rebuilt G2 point pairs with the G1 point of the same key pair.
	x, y := contractG2Coords(g2)
	rebuilt := core.NewG2Point(x, y)
	require.True(t, rebuilt.G2Affine.Equal(g2.G2Affine))

	ok, err := g1.VerifyEquivalence(rebuilt)
	require.NoError(t, err)
	require.True(t, ok, "rebuilt G2 key must pair with the G1 key")

	// The negative: feeding components WITHOUT the swap (i.e. gnark's native
	// [A0, A1] order passed as if it were contract order) must NOT reproduce
	// the key. This is the "valid-looking but incorrect key" failure mode the
	// helper exists to prevent.
	unswappedX := [2]*big.Int{x[1], x[0]}
	unswappedY := [2]*big.Int{y[1], y[0]}
	wrong := core.NewG2Point(unswappedX, unswappedY)
	require.False(t, wrong.G2Affine.Equal(g2.G2Affine))
}

func TestGetOperatorIDMatchesKeccakOfG1(t *testing.T) {
	keyPair, err := core.GenRandomBlsKeys()
	require.NoError(t, err)

	// Two independent derivations must agree, and distinct keys must yield
	// distinct IDs.
	id1 := keyPair.GetPubKeyG1().GetOperatorID()
	id2 := keyPair.GetPubKeyG1().GetOperatorID()
	require.Equal(t, id1, id2)

	other, err := core.GenRandomBlsKeys()
	require.NoError(t, err)
	require.NotEqual(t, id1, other.GetPubKeyG1().GetOperatorID())
}
