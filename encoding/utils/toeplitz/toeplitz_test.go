package toeplitz

import (
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"

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

	toe, err := NewToeplitz(v, fs)
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

	_, err := NewToeplitz(v, fs)
	assert.EqualError(t, err, "num diagonal vector must be odd")
}

// Expand toeplitz matrix into circular matrix
// the outcome is a also concise representation
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
	toep, err := NewToeplitz(v, fs)
	require.Nil(t, err)

	cVec := toep.extendCirculantVec()
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
	toep, err := NewToeplitz(v, fs)
	require.Nil(t, err)

	cVec := toep.extendCirculantVec()
	rVec := toep.fromColVToRowV(cVec)
	assert.Equal(t, rVec[0], v[0])
	assert.Equal(t, rVec[1], v[1])
	assert.Equal(t, rVec[2], v[2])
	assert.Equal(t, rVec[3], v[3])
	assert.Equal(t, rVec[4], encoding.ZERO)
	assert.Equal(t, rVec[5], v[4])
	assert.Equal(t, rVec[6], v[5])
	assert.Equal(t, rVec[7], v[6])

	// involutory
	cVec = toep.fromColVToRowV(rVec)
	assert.Equal(t, cVec[0], v[0])
	assert.Equal(t, cVec[1], v[6])
	assert.Equal(t, cVec[2], v[5])
	assert.Equal(t, cVec[3], v[4])
	assert.Equal(t, cVec[4], encoding.ZERO)
	assert.Equal(t, cVec[5], v[3])
	assert.Equal(t, cVec[6], v[2])
	assert.Equal(t, cVec[7], v[1])
}
