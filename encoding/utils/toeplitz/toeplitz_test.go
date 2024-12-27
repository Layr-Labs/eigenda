package toeplitz_test

import (
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/utils/toeplitz"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// V is ordered as (v_0, .., v_6), so it creates a
// matrix below. Slice must be odd
// v_0 v_1 v_2 v_3
// v_6 v_0 v_1 v_2
// v_5 v_6 v_0 v_1
// v_4 v_5 v_6 v_0
func TestNewToeplitz(t *testing.T) {
	v := make([]fr.Element, 7)
	v[0].SetInt64(int64(7))
	v[1].SetInt64(int64(11))
	v[2].SetInt64(int64(5))
	v[3].SetInt64(int64(6))
	v[4].SetInt64(int64(3))
	v[5].SetInt64(int64(8))
	v[6].SetInt64(int64(1))
	fs := fft.NewFFTSettings(4)

	toe, err := toeplitz.NewToeplitz(v, fs)
	require.Nil(t, err)

	assert.Equal(t, v[0], toe.V[0])
	assert.Equal(t, v[1], toe.V[1])
	assert.Equal(t, v[2], toe.V[2])
	assert.Equal(t, v[3], toe.V[3])
	assert.Equal(t, v[4], toe.V[4])
	assert.Equal(t, v[5], toe.V[5])
}

func TestNewToeplitz_InvalidSize(t *testing.T) {
	v := make([]fr.Element, 2)
	v[0].SetInt64(int64(4))
	v[1].SetInt64(int64(2))
	fs := fft.NewFFTSettings(4)

	_, err := toeplitz.NewToeplitz(v, fs)
	assert.EqualError(t, err, "num diagonal vector must be odd")
}

// Expand toeplitz matrix into circular matrix
// the outcome is also a concise representation
// if   V is (v_0, v_1, v_2, v_3, v_4, v_5, v_6)
// then E is (v_0, v_6, v_5, v_4, 0,   v_3, v_2, v_1)
func TestExtendCircularVec(t *testing.T) {
	v := make([]fr.Element, 7)
	v[0].SetInt64(int64(7))
	v[1].SetInt64(int64(11))
	v[2].SetInt64(int64(5))
	v[3].SetInt64(int64(6))
	v[4].SetInt64(int64(3))
	v[5].SetInt64(int64(8))
	v[6].SetInt64(int64(1))

	fs := fft.NewFFTSettings(4)
	c, err := toeplitz.NewToeplitz(v, fs)
	require.Nil(t, err)

	cVec := c.ExtendCircularVec()
	assert.Equal(t, cVec[0], v[0])
	assert.Equal(t, cVec[1], v[6])
	assert.Equal(t, cVec[2], v[5])
	assert.Equal(t, cVec[3], v[4])
	assert.Equal(t, cVec[4], encoding.ZERO)
	assert.Equal(t, cVec[5], v[3])
	assert.Equal(t, cVec[6], v[2])
	assert.Equal(t, cVec[7], v[1])
}

// if   col Vector is [v_0, v_1, v_2, v_3, 0, v_4, v_5, v_6]
// then row Vector is [v_0, v_6, v_5, v_4, 0, v_3, v_2, v_1]
// this operation is involutory. i.e. f(f(v)) = v
func TestFromColVToRowV(t *testing.T) {
	v := make([]fr.Element, 7)
	v[0].SetInt64(int64(7))
	v[1].SetInt64(int64(11))
	v[2].SetInt64(int64(5))
	v[3].SetInt64(int64(6))
	v[4].SetInt64(int64(3))
	v[5].SetInt64(int64(8))
	v[6].SetInt64(int64(1))

	fs := fft.NewFFTSettings(4)
	c, err := toeplitz.NewToeplitz(v, fs)
	require.Nil(t, err)

	cVec := c.ExtendCircularVec()
	rVec := c.FromColVToRowV(cVec)

	assert.Equal(t, rVec[0], v[0])
	assert.Equal(t, rVec[1], v[1])
	assert.Equal(t, rVec[2], v[2])
	assert.Equal(t, rVec[3], v[3])
	assert.Equal(t, rVec[4], encoding.ZERO)
	assert.Equal(t, rVec[5], v[4])
	assert.Equal(t, rVec[6], v[5])
	assert.Equal(t, rVec[7], v[6])

	// involutory
	cVec = c.FromColVToRowV(rVec)
	assert.Equal(t, cVec[0], v[0])
	assert.Equal(t, cVec[1], v[6])
	assert.Equal(t, cVec[2], v[5])
	assert.Equal(t, cVec[3], v[4])
	assert.Equal(t, cVec[4], encoding.ZERO)
	assert.Equal(t, cVec[5], v[3])
	assert.Equal(t, cVec[6], v[2])
	assert.Equal(t, cVec[7], v[1])
}

func TestMultiplyToeplitz(t *testing.T) {
	v := make([]fr.Element, 7)
	v[0].SetInt64(int64(7))
	v[1].SetInt64(int64(11))
	v[2].SetInt64(int64(5))
	v[3].SetInt64(int64(6))
	v[4].SetInt64(int64(3))
	v[5].SetInt64(int64(8))
	v[6].SetInt64(int64(1))

	fs := fft.NewFFTSettings(4)
	toe, err := toeplitz.NewToeplitz(v, fs)

	require.Nil(t, err)

	x := make([]fr.Element, 4)
	x[0].SetInt64(int64(1))
	x[1].SetInt64(int64(2))
	x[2].SetInt64(int64(3))
	x[3].SetInt64(int64(4))

	b, err := toe.Multiply(x)
	require.Nil(t, err)

	p := make([]fr.Element, 4)
	p[0].SetInt64(int64(68))
	p[1].SetInt64(int64(68))
	p[2].SetInt64(int64(75))
	p[3].SetInt64(int64(50))

	assert.Equal(t, b[0], p[0])
	assert.Equal(t, b[1], p[1])
	assert.Equal(t, b[2], p[2])
	assert.Equal(t, b[3], p[3])

	// Assert with direct multiplication
	b2 := toe.DirectMultiply(x)
	assert.Equal(t, b, b2)
}
