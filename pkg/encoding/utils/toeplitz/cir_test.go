package toeplitz_test

import (
	"testing"

	"github.com/Layr-Labs/eigenda/pkg/encoding/utils/toeplitz"
	kzg "github.com/Layr-Labs/eigenda/pkg/kzg"
	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCircular(t *testing.T) {
	v := make([]bls.Fr, 4)
	v[0] = bls.ToFr("7")
	v[1] = bls.ToFr("6")
	v[2] = bls.ToFr("5")
	v[3] = bls.ToFr("11")
	fs := kzg.NewFFTSettings(4)

	c := toeplitz.NewCircular(v, fs)

	assert.Equal(t, v[0], c.V[0])
	assert.Equal(t, v[1], c.V[1])
	assert.Equal(t, v[2], c.V[2])
	assert.Equal(t, v[3], c.V[3])
}

func TestMultiplyCircular(t *testing.T) {
	v := make([]bls.Fr, 4)
	v[0] = bls.ToFr("7")
	v[1] = bls.ToFr("11")
	v[2] = bls.ToFr("5")
	v[3] = bls.ToFr("6")
	fs := kzg.NewFFTSettings(4)

	c := toeplitz.NewCircular(v, fs)

	x := make([]bls.Fr, 4)
	x[0] = bls.ToFr("1")
	x[1] = bls.ToFr("2")
	x[2] = bls.ToFr("3")
	x[3] = bls.ToFr("4")
	b, err := c.Multiply(x)
	require.Nil(t, err)

	assert.Equal(t, b[0], bls.ToFr("68"))
	assert.Equal(t, b[1], bls.ToFr("73"))
	assert.Equal(t, b[2], bls.ToFr("82"))
	assert.Equal(t, b[3], bls.ToFr("67"))

	// Assert with direct multiplication
	b2 := c.DirectMultiply(x)
	assert.Equal(t, b, b2)
}

func TestMultiplyCircular_InvalidDimensions(t *testing.T) {
	v := make([]bls.Fr, 2)
	v[0] = bls.ToFr("7")
	v[1] = bls.ToFr("11")
	fs := kzg.NewFFTSettings(2)

	c := toeplitz.NewCircular(v, fs)

	x := make([]bls.Fr, 4)
	x[0] = bls.ToFr("1")
	x[1] = bls.ToFr("2")
	x[2] = bls.ToFr("3")
	x[3] = bls.ToFr("4")
	_, err := c.Multiply(x)
	assert.EqualError(t, err, "dimension inconsistent")
}

func TestMultiplyPointsCircular(t *testing.T) {
	v := make([]bls.Fr, 4)
	v[0] = bls.ToFr("7")
	v[1] = bls.ToFr("6")
	v[2] = bls.ToFr("5")
	v[3] = bls.ToFr("11")
	fs := kzg.NewFFTSettings(4)

	c := toeplitz.NewCircular(v, fs)

	x := make([]bls.G1Point, 4)
	x[0] = bls.GenG1
	x[1] = bls.GenG1
	x[2] = bls.GenG1
	x[3] = bls.GenG1

	b, err := c.MultiplyPoints(x, false, true)
	require.Nil(t, err)

	sum := bls.LinCombG1(x, v)
	assert.Equal(t, &b[0], sum)
}

func TestMultiplyPointsCircular_InvalidDimension(t *testing.T) {
	v := make([]bls.Fr, 2)
	v[0] = bls.ToFr("7")
	v[1] = bls.ToFr("6")
	fs := kzg.NewFFTSettings(2)

	c := toeplitz.NewCircular(v, fs)

	x := make([]bls.G1Point, 4)
	x[0] = bls.GenG1
	x[1] = bls.GenG1
	x[2] = bls.GenG1
	x[3] = bls.GenG1

	_, err := c.MultiplyPoints(x, false, true)
	assert.EqualError(t, err, "dimension inconsistent. Input != vector")
}

func TestHadamard_InvalidDimension(t *testing.T) {
	a := make([]bls.Fr, 2)
	a[0] = bls.ToFr("1")
	a[1] = bls.ToFr("2")

	b := make([]bls.Fr, 1)
	b[0] = bls.ToFr("3")

	c := make([]bls.Fr, 3)
	err := toeplitz.Hadamard(a, b, c)
	assert.EqualError(t, err, "dimension inconsistent. Cannot do Hadamard Product on Fr")

	// TODO: This causes a panic because there are no checks on the size of c
	// b = make([]bls.Fr, 2)
	// b[0] = bls.ToFr("3")
	// b[1] = bls.ToFr("4")

	// c = make([]bls.Fr, 1)
	// fmt.Println(len(a), len(b), len(c))
	// err = kzgRs.Hadamard(a, b, c)
	// require.Nil(t, err)
}

func TestHadamardPoint_InvalidDimension(t *testing.T) {
	a := make([]bls.G1Point, 2)
	a[0] = bls.GenG1
	a[1] = bls.GenG1

	b := make([]bls.Fr, 1)
	b[0] = bls.ToFr("1")

	c := make([]bls.G1Point, 3)
	err := toeplitz.HadamardPoints(a, b, c)
	assert.EqualError(t, err, "dimension inconsistent. Cannot do Hadamard Product on Points")

	// TODO: This causes a panic because there are no checks on the size of c
	// b = make([]bls.Fr, 2)
	// b[0] = bls.ToFr("3")
	// b[1] = bls.ToFr("4")

	// c = make([]bls.Fr, 1)
	// fmt.Println(len(a), len(b), len(c))
	// err = kzgRs.Hadamard(a, b, c)
	// require.Nil(t, err)
}
