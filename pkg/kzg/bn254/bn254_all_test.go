package bn254

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPointG1Marshalling(t *testing.T) {
	var x Fr
	SetFr(&x, "44689111813071777962210527909085028157792767057343609826799812096627770269092")
	var point G1Point
	MulG1(&point, &GenG1, &x)

	bytes := point.MarshalText()

	var anotherPoint G1Point
	err := anotherPoint.UnmarshalText(bytes)
	require.Nil(t, err)

	assert.True(t, EqualG1(&point, &anotherPoint), "G1 points did not match\n%s\n%s", StrG1(&point), StrG1(&anotherPoint))
}

func TestPointG1Marshalling_InvalidG1(t *testing.T) {
	var g1 *G1Point
	err := g1.UnmarshalText([]byte(""))
	assert.EqualError(t, err, "cannot decode into nil G1Point")

	g1 = new(G1Point)
	err = g1.UnmarshalText([]byte("G"))
	assert.EqualError(t, err, "encoding/hex: invalid byte: U+0047 'G'")

	err = g1.UnmarshalText([]byte("8000000000000000000000000000000000000000000000000000000000000099"))
	assert.EqualError(t, err, "invalid compressed coordinate: square root doesn't exist")
}

func TestPointG2Marshalling(t *testing.T) {
	var x Fr
	SetFr(&x, "44689111813071777962210527909085028157792767057343609826799812096627770269092")
	var point G2Point
	MulG2(&point, &GenG2, &x)

	bytes := point.MarshalText()

	var anotherPoint G2Point
	err := anotherPoint.UnmarshalText(bytes)
	require.Nil(t, err)

	assert.True(t, EqualG2(&point, &anotherPoint), "G2 points did not match:\n%s\n%s", StrG2(&point), StrG2(&anotherPoint))
}

func TestPointG2Marshalling_InvalidG2(t *testing.T) {
	var g2 *G2Point
	err := g2.UnmarshalText([]byte(""))
	assert.EqualError(t, err, "cannot decode into nil G2Point")

	g2 = new(G2Point)
	err = g2.UnmarshalText([]byte("G"))
	assert.EqualError(t, err, "encoding/hex: invalid byte: U+0047 'G'")

	err = g2.UnmarshalText([]byte("898e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed"))
	assert.EqualError(t, err, "invalid point: subgroup check failed")

	err = g2.UnmarshalText([]byte("998e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992ffff"))
	assert.EqualError(t, err, "invalid compressed coordinate: square root doesn't exist")

}
