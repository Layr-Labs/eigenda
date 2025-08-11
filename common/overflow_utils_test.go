package common

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSafeAddInt64(t *testing.T) {
	t.Run("normal addition", func(t *testing.T) {
		result, err := SafeAddInt64(10, 20)
		require.NoError(t, err)
		assert.Equal(t, int64(30), result)
	})

	t.Run("positive overflow", func(t *testing.T) {
		_, err := SafeAddInt64(math.MaxInt64, 1)
		require.Error(t, err)
	})

	t.Run("negative overflow", func(t *testing.T) {
		_, err := SafeAddInt64(math.MinInt64, -1)
		require.Error(t, err)
	})
}

func TestSafeSubtractInt64(t *testing.T) {
	t.Run("normal subtraction", func(t *testing.T) {
		result, err := SafeSubtractInt64(30, 20)
		require.NoError(t, err)
		assert.Equal(t, int64(10), result)
	})

	t.Run("positive overflow", func(t *testing.T) {
		_, err := SafeSubtractInt64(math.MaxInt64, -1)
		require.Error(t, err)
	})

	t.Run("negative overflow", func(t *testing.T) {
		_, err := SafeSubtractInt64(math.MinInt64, 1)
		require.Error(t, err)
	})
}

func TestSafeMultiplyInt64(t *testing.T) {
	t.Run("normal multiplication", func(t *testing.T) {
		result, err := SafeMultiplyInt64(10, 20)
		require.NoError(t, err)
		assert.Equal(t, int64(200), result)
	})

	t.Run("zero multiplication", func(t *testing.T) {
		result, err := SafeMultiplyInt64(0, 100)
		require.NoError(t, err)
		assert.Equal(t, int64(0), result)

		result, err = SafeMultiplyInt64(100, 0)
		require.NoError(t, err)
		assert.Equal(t, int64(0), result)
	})

	t.Run("positive overflow", func(t *testing.T) {
		_, err := SafeMultiplyInt64(math.MaxInt64, 2)
		require.Error(t, err)

		_, err = SafeMultiplyInt64(-math.MaxInt64, -2)
		require.Error(t, err)
	})

	t.Run("negative overflow", func(t *testing.T) {
		_, err := SafeMultiplyInt64(2, math.MinInt64)
		require.Error(t, err)

		_, err = SafeMultiplyInt64(math.MinInt64, 2)
		require.Error(t, err)
	})
}
