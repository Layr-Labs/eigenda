//go:build !bignum_pure && !bignum_hol256
// +build !bignum_pure,!bignum_hol256

package bn254

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPointCompression(t *testing.T) {
	var x Fr
	SetFr(&x, "44689111813071777962210527909085028157792767057343609826799812096627770269092")
	var point G1Point
	MulG1(&point, &GenG1, &x)
	got := ToCompressedG1(&point)

	expected, err := hex.DecodeString("cb6a192085f84b7f34e29b5e5056927581f573201abb1ce1b8ee132c71fae1e2")
	if !bytes.Equal(expected, got) || err != nil {
		t.Fatalf("Invalid compression result, %x != %x", got, expected)
	}
}

func TestZeroG1(t *testing.T) {
	zero := ZeroG1
	expected := "0"
	if zero.X != zero.Y || zero.X.String() != expected {
		t.Fatalf("ZeroG1 is not zero! expected %s, but got: %s", expected, zero.X.String())
	}
}

func TestGenG1(t *testing.T) {
	one := GenG1
	expectedX := "1"
	expectedY := "2"

	if one.X.String() != expectedX {
		t.Fatalf("expected %s, but got: %s", expectedX, one.X.String())
	}

	if one.Y.String() != expectedY {
		t.Fatalf("expected %s, but got: %s", expectedY, one.Y.String())
	}
}

func TestClearG1(t *testing.T) {
	clear := GenG1
	ClearG1(&clear)

	if clear != ZERO_G1 {
		t.Fatalf("ClearG1 failed %v", clear.X.String())
	}
}

func TestCopyG1(t *testing.T) {
	one := GenG1
	var dst G1Point
	CopyG1(&dst, &one)

	if one != dst {
		t.Fatalf("CopyG1 failed %v", dst.X.String())
	}
}

func TestMulG1(t *testing.T) {
	g1 := GenG1
	two := ToFr("2")
	var prod G1Point
	MulG1(&prod, &g1, &two)

	var add G1Point
	AddG1(&add, &g1, &g1)

	if add != prod {
		t.Fatalf("MulG1 failed")
	}
}

func TestAddG1(t *testing.T) {
	zero := ZERO_G1
	one := GenG1
	var dst G1Point
	AddG1(&dst, &zero, &one)

	if dst != one {
		t.Fatalf("AddG1 failed %v", dst.X.String())
	}
}

func TestSubG1(t *testing.T) {
	one := GenG1
	var dst G1Point
	SubG1(&dst, &one, &one)

	if dst != ZERO_G1 {
		t.Fatalf("SubG1 failed %v", dst.X.String())
	}
}

func TestStrG1(t *testing.T) {
	g1 := GenG1
	g1Str := StrG1(&g1)
	expected := "E([1,2])\n"
	if g1Str != expected {
		t.Fatalf("StrG1 failed, %v %v", g1Str, expected)
	}
}

func TestStringG1(t *testing.T) {
	g1 := GenG1
	g1Str := g1.String()
	expected := "E([1,2])\n"
	if g1Str != expected {
		t.Fatalf("StrG1 failed, %v %v", g1Str, expected)
	}
}

func TestNegG1(t *testing.T) {
	negG1 := GenG1
	NegG1(&negG1)

	// Add
	var res G1Point
	AddG1(&res, &GenG1, &negG1)
	if res != ZERO_G1 {
		t.Fatal("NegG1 failed")
	}
}

func TestLinCombG1(t *testing.T) {
	// TODO: use random poly and g1 points
	poly := []Fr{
		ToFr("1"), ToFr("2"), ToFr("3"), ToFr("4"), ToFr("5"),
	}
	one := GenG1
	val := []G1Point{
		one, one, one, one, one,
	}
	lin := LinCombG1(val, poly)

	var scalar Fr
	AsFr(&scalar, uint64(15))
	var product G1Point
	MulG1(&product, &one, &scalar)
	if *lin != product {
		t.Fatal("Linear combination != product!")
	}
}

func TestZeroG2(t *testing.T) {
	zero := ZeroG2
	expected := "0+0*u"

	if zero.X != zero.Y || zero.X.String() != expected {
		t.Fatalf("ZeroG2 is not zero! expected %s, but got: %s", expected, zero.X.String())
	}
}

func TestGenG2(t *testing.T) {
	two := GenG2
	expectedTwoX := "10857046999023057135944570762232829481370756359578518086990519993285655852781+11559732032986387107991004021392285783925812861821192530917403151452391805634*u"
	expectedTwoY := "8495653923123431417604973247489272438418190587263600148770280649306958101930+4082367875863433681332203403145435568316851327593401208105741076214120093531*u"

	if two.X.String() != expectedTwoX {
		t.Fatalf("expected %s, but got: %s", expectedTwoX, two.X.String())
	}

	if two.Y.String() != expectedTwoY {
		t.Fatalf("expected %s, but got: %s", expectedTwoY, two.Y.String())
	}
}

func TestClearG2(t *testing.T) {
	clear := GenG2
	ClearG2(&clear)

	if clear != ZERO_G2 {
		t.Fatalf("ClearG2 failed %v", clear.X.String())
	}
}

func TestCopyG2(t *testing.T) {
	two := GenG2
	var dst G2Point
	CopyG2(&dst, &two)

	if two != dst {
		t.Fatalf("CopyG2 failed %v", dst.X.String())
	}
}

func TestMulG2(t *testing.T) {
	g2 := GenG2
	two := ToFr("2")
	var prod G2Point
	MulG2(&prod, &g2, &two)

	var add G2Point
	AddG2(&add, &g2, &g2)

	if add != prod {
		t.Fatalf("MulG2 failed")
	}
}

func TestAddG2(t *testing.T) {
	zero := ZERO_G2
	one := GenG2
	var dst G2Point
	AddG2(&dst, &zero, &one)

	if dst != one {
		t.Fatalf("AddG2 failed %v", dst.X.String())
	}
}

func TestSubG2(t *testing.T) {
	one := GenG2
	var dst G2Point
	SubG2(&dst, &one, &one)

	if dst != ZERO_G2 {
		t.Fatalf("SubG1 failed %v", dst.X.String())
	}
}

func TestStrG2(t *testing.T) {
	g2 := GenG2
	g2Str := StrG2(&g2)
	expected := "E([10857046999023057135944570762232829481370756359578518086990519993285655852781+11559732032986387107991004021392285783925812861821192530917403151452391805634*u,8495653923123431417604973247489272438418190587263600148770280649306958101930+4082367875863433681332203403145435568316851327593401208105741076214120093531*u])"
	if g2Str != expected {
		t.Fatalf("StrG2 failed, %v %v", g2Str, expected)
	}
}

func TestStringG2(t *testing.T) {
	g2 := GenG2
	g2Str := g2.String()
	expected := "E([10857046999023057135944570762232829481370756359578518086990519993285655852781+11559732032986387107991004021392285783925812861821192530917403151452391805634*u,8495653923123431417604973247489272438418190587263600148770280649306958101930+4082367875863433681332203403145435568316851327593401208105741076214120093531*u])"
	if g2Str != expected {
		t.Fatalf("StrG2 failed, %v %v", g2Str, expected)
	}
}

func TestNegG2(t *testing.T) {
	negG2 := GenG2
	NegG2(&negG2)

	// Add
	var res G2Point
	AddG2(&res, &GenG2, &negG2)
	if res != ZERO_G2 {
		t.Fatal("NegG2 failed")
	}
}

func TestLinCombG2(t *testing.T) {
	// TODO: use random poly and g1 points
	poly := []Fr{
		ToFr("1"), ToFr("2"), ToFr("3"), ToFr("4"), ToFr("5"),
	}
	one := GenG2
	val := []G2Point{
		one, one, one, one, one,
	}
	lin := LinCombG2(val, poly)

	var scalar Fr
	AsFr(&scalar, uint64(15))
	var product G2Point
	MulG2(&product, &one, &scalar)
	if *lin != product {
		t.Fatal("Linear combination != product!")
	}
}

func TestGenerators(t *testing.T) {
	g1, g2 := Generators()
	if g1 != GenG1 && g2 != GenG2 {
		t.Fatal("Setup failed")
	}
}

func TestPairingsVerify(t *testing.T) {
	g1, g2 := Generators()

	// Invalid pairing
	val := PairingsVerify(&g1, &g2, &ZeroG1, &ZeroG2)
	assert.False(t, val)

	// TODO: test valid pairing
}
