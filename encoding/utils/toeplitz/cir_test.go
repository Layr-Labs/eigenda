package toeplitz_test

import (
	"testing"

	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/utils/toeplitz"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/stretchr/testify/assert"
)

func TestNewCircular(t *testing.T) {
	v := make([]fr.Element, 4)
	v[0].SetInt64(int64(7))
	v[1].SetInt64(int64(6))
	v[2].SetInt64(int64(5))
	v[3].SetInt64(int64(11))
	fs := fft.NewFFTSettings(4)

	c := toeplitz.NewCircular(v, fs)

	assert.Equal(t, v[0], c.V[0])
	assert.Equal(t, v[1], c.V[1])
	assert.Equal(t, v[2], c.V[2])
	assert.Equal(t, v[3], c.V[3])
}

func TestMultiplyCircular_InvalidDimensions(t *testing.T) {
	v := make([]fr.Element, 2)
	v[0].SetInt64(int64(7))
	v[1].SetInt64(int64(11))
	fs := fft.NewFFTSettings(2)

	c := toeplitz.NewCircular(v, fs)

	x := make([]fr.Element, 4)
	x[0].SetInt64(int64(1))
	x[1].SetInt64(int64(2))
	x[2].SetInt64(int64(3))
	x[3].SetInt64(int64(4))
	_, err := c.Multiply(x)
	assert.EqualError(t, err, "dimension inconsistent")
}

func TestHadamard_InvalidDimension(t *testing.T) {
	a := make([]fr.Element, 2)
	a[0].SetInt64(int64(1))
	a[1].SetInt64(int64(2))

	b := make([]fr.Element, 1)
	b[0].SetInt64(int64(3))

	c := make([]fr.Element, 3)
	err := toeplitz.Hadamard(a, b, c)
	assert.EqualError(t, err, "dimension inconsistent. Cannot do Hadamard Product on Fr")

	// TODO: This causes a panic because there are no checks on the size of c
	// b = make([]fr.Element, 2)
	// b[0] = bls.ToFr("3")
	// b[1] = bls.ToFr("4")

	// c = make([]fr.Element, 1)
	// fmt.Println(len(a), len(b), len(c))
	// err = kzgRs.Hadamard(a, b, c)
	// require.Nil(t, err)
}
