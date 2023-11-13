package toeplitz_test

import (
	"testing"

	"github.com/Layr-Labs/eigenda/pkg/encoding/utils/toeplitz"
	kzg "github.com/Layr-Labs/eigenda/pkg/kzg"
	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
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
	v := make([]bls.Fr, 7)
	v[0] = bls.ToFr("7")
	v[1] = bls.ToFr("11")
	v[2] = bls.ToFr("5")
	v[3] = bls.ToFr("6")
	v[4] = bls.ToFr("3")
	v[5] = bls.ToFr("8")
	v[6] = bls.ToFr("1")
	fs := kzg.NewFFTSettings(4)

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
	v := make([]bls.Fr, 2)
	v[0] = bls.ToFr("4")
	v[1] = bls.ToFr("2")
	fs := kzg.NewFFTSettings(4)

	_, err := toeplitz.NewToeplitz(v, fs)
	assert.EqualError(t, err, "num diagonal vector must be odd")
}

// Expand toeplitz matrix into circular matrix
// the outcome is a also concise representation
// if   V is (v_0, v_1, v_2, v_3, v_4, v_5, v_6)
// then E is (v_0, v_6, v_5, v_4, 0,   v_3, v_2, v_1)
func TestExtendCircularVec(t *testing.T) {
	v := make([]bls.Fr, 7)
	v[0] = bls.ToFr("7")
	v[1] = bls.ToFr("11")
	v[2] = bls.ToFr("5")
	v[3] = bls.ToFr("6")
	v[4] = bls.ToFr("3")
	v[5] = bls.ToFr("8")
	v[6] = bls.ToFr("1")

	fs := kzg.NewFFTSettings(4)
	c, err := toeplitz.NewToeplitz(v, fs)
	require.Nil(t, err)

	cVec := c.ExtendCircularVec()
	assert.Equal(t, cVec[0], v[0])
	assert.Equal(t, cVec[1], v[6])
	assert.Equal(t, cVec[2], v[5])
	assert.Equal(t, cVec[3], v[4])
	assert.Equal(t, cVec[4], bls.ZERO)
	assert.Equal(t, cVec[5], v[3])
	assert.Equal(t, cVec[6], v[2])
	assert.Equal(t, cVec[7], v[1])
}

// if   col Vector is [v_0, v_1, v_2, v_3, 0, v_4, v_5, v_6]
// then row Vector is [v_0, v_6, v_5, v_4, 0, v_3, v_2, v_1]
// this operation is involutory. i.e. f(f(v)) = v
func TestFromColVToRowV(t *testing.T) {
	v := make([]bls.Fr, 7)
	v[0] = bls.ToFr("7")
	v[1] = bls.ToFr("11")
	v[2] = bls.ToFr("5")
	v[3] = bls.ToFr("6")
	v[4] = bls.ToFr("3")
	v[5] = bls.ToFr("8")
	v[6] = bls.ToFr("1")

	fs := kzg.NewFFTSettings(4)
	c, err := toeplitz.NewToeplitz(v, fs)
	require.Nil(t, err)

	cVec := c.ExtendCircularVec()
	rVec := c.FromColVToRowV(cVec)

	assert.Equal(t, rVec[0], v[0])
	assert.Equal(t, rVec[1], v[1])
	assert.Equal(t, rVec[2], v[2])
	assert.Equal(t, rVec[3], v[3])
	assert.Equal(t, rVec[4], bls.ZERO)
	assert.Equal(t, rVec[5], v[4])
	assert.Equal(t, rVec[6], v[5])
	assert.Equal(t, rVec[7], v[6])

	// involutory
	cVec = c.FromColVToRowV(rVec)
	assert.Equal(t, cVec[0], v[0])
	assert.Equal(t, cVec[1], v[6])
	assert.Equal(t, cVec[2], v[5])
	assert.Equal(t, cVec[3], v[4])
	assert.Equal(t, cVec[4], bls.ZERO)
	assert.Equal(t, cVec[5], v[3])
	assert.Equal(t, cVec[6], v[2])
	assert.Equal(t, cVec[7], v[1])
}

func TestMultiplyToeplitz(t *testing.T) {
	v := make([]bls.Fr, 7)
	v[0] = bls.ToFr("7")
	v[1] = bls.ToFr("11")
	v[2] = bls.ToFr("5")
	v[3] = bls.ToFr("6")
	v[4] = bls.ToFr("3")
	v[5] = bls.ToFr("8")
	v[6] = bls.ToFr("1")

	fs := kzg.NewFFTSettings(4)
	toe, err := toeplitz.NewToeplitz(v, fs)

	require.Nil(t, err)

	x := make([]bls.Fr, 4)
	x[0] = bls.ToFr("1")
	x[1] = bls.ToFr("2")
	x[2] = bls.ToFr("3")
	x[3] = bls.ToFr("4")

	b, err := toe.Multiply(x)
	require.Nil(t, err)

	assert.Equal(t, b[0], bls.ToFr("68"))
	assert.Equal(t, b[1], bls.ToFr("68"))
	assert.Equal(t, b[2], bls.ToFr("75"))
	assert.Equal(t, b[3], bls.ToFr("50"))

	// Assert with direct multiplication
	b2 := toe.DirectMultiply(x)
	assert.Equal(t, b, b2)
}

func TestMultiplyPointsToeplitz(t *testing.T) {
	v := make([]bls.Fr, 7)
	v[0] = bls.ToFr("7")
	v[1] = bls.ToFr("11")
	v[2] = bls.ToFr("5")
	v[3] = bls.ToFr("6")
	v[4] = bls.ToFr("3")
	v[5] = bls.ToFr("8")
	v[6] = bls.ToFr("1")
	fs := kzg.NewFFTSettings(4)
	toe, err := toeplitz.NewToeplitz(v, fs)
	require.Nil(t, err)

	x := make([]bls.G1Point, 8)
	x[0] = bls.GenG1
	x[1] = bls.GenG1
	x[2] = bls.GenG1
	x[3] = bls.GenG1
	x[4] = bls.GenG1
	x[5] = bls.GenG1
	x[6] = bls.GenG1
	x[7] = bls.GenG1

	b1, err := toe.MultiplyPoints(x, false, true)
	require.Nil(t, err)

	//b2, err := toe.MultiplyPoints(x, false, false)
	require.Nil(t, err)

	sum := bls.LinCombG1(x[:7], toe.V)
	assert.Equal(t, &b1[0], sum)

	// TODO: Calculate inverse
	// b2, err := toe.MultiplyPoints(x, true, true)
	// require.Nil(t, err)

	// res, err := toe.MultiplyPoints(b2, false, true)
	// require.Nil(t, err)
	// assert.Equal(t, res[0].X.String(), x[0].X.String())
}
