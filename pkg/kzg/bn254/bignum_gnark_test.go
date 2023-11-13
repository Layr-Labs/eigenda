package bn254

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	two := ToFr("2")
	expected := "2"

	assert.Equal(t, two.String(), expected)
}

func TestFrToBytes(t *testing.T) {
	two := ToFr("2")
	b := FrToBytes(&two)
	expected := [32]byte{}
	expected[31] = 2

	assert.Equal(t, b, expected)
}

func TestFrSetBytes(t *testing.T) {
	two := make([]byte, 32)
	two[31] = 2
	var dst Fr
	FrSetBytes(&dst, two)
	expected := ToFr("2")

	assert.Equal(t, dst, expected)
}

func TestFrFrom32(t *testing.T) {
	one := [32]byte{}
	one[31] = 1
	var dst Fr
	FrFrom32(&dst, one)
	expected := ToFr("1")

	assert.Equal(t, dst, expected)
}

func TestFrTo32(t *testing.T) {
	one := ToFr("1")
	b := FrTo32(&one)
	expected := [32]byte{}
	expected[31] = 1

	assert.Equal(t, b, expected)
}

func TestCopyFrStr(t *testing.T) {
	one := ToFr("1")
	var dst Fr
	CopyFr(&dst, &one)

	assert.Equal(t, dst, one)
}

func TestFrStr(t *testing.T) {
	s := FrStr(nil)
	expected := "<nil>"

	assert.Equal(t, s, expected)

	one := ToFr("1")
	s = FrStr(&one)
	expected = "1"

	assert.Equal(t, s, expected)
}

func TestEqualZero(t *testing.T) {
	zero := ToFr("0")
	eq := EqualZero(&zero)

	assert.True(t, eq)
}

func TestEqualOne(t *testing.T) {
	one := ToFr("1")
	eq := EqualOne(&one)

	assert.True(t, eq)
}

func TestEqualFr(t *testing.T) {
	one := ToFr("1")
	two := ToFr("2")

	eq := EqualFr(&one, &two)
	assert.False(t, eq)

	eq = EqualFr(&one, &one)
	assert.True(t, eq)
}

func TestAddFr(t *testing.T) {
	one := ToFr("1")
	two := ToFr("2")
	three := ToFr("3")

	var res Fr
	AsFr(&res, 0)
	AddModFr(&res, &one, &two)

	assert.Equal(t, res, three)
}

func TestSubFr(t *testing.T) {
	one := ToFr("1")
	two := ToFr("2")
	three := ToFr("3")

	var res Fr
	AsFr(&res, 0)
	SubModFr(&res, &three, &two)

	assert.Equal(t, res, one)
}

func TestDivModFr(t *testing.T) {
	two := ToFr("2")
	three := ToFr("3")
	six := ToFr("6")

	var res Fr
	AsFr(&res, 0)
	DivModFr(&res, &six, &three)

	assert.Equal(t, res, two)
}

func TestMulModFr(t *testing.T) {
	two := ToFr("2")
	three := ToFr("3")
	six := ToFr("6")

	var res Fr
	AsFr(&res, 0)
	MulModFr(&res, &two, &three)

	assert.Equal(t, res, six)
}

func TestInvModFr(t *testing.T) {
	two := ToFr("2")
	invTwo := ToFr("10944121435919637611123202872628637544274182200208017171849102093287904247809")

	var res Fr
	AsFr(&res, 0)
	InvModFr(&res, &two)

	assert.Equal(t, res, invTwo)
}

func TestEvalPolyAt(t *testing.T) {
	one := ToFr("1")
	var dst Fr
	p := []Fr{
		ToFr("1"),
		ToFr("2"),
		ToFr("3"),
	}
	EvalPolyAt(&dst, p, &one)

	assert.Equal(t, dst, ToFr("6"))
}
