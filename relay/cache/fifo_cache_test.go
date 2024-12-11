package cache

import (
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
	"testing"
)

func TestExpirationOrder(t *testing.T) {
	tu.InitializeRandom()

	maxWeight := uint64(10 + rand.Intn(10))
	c := NewFIFOCache[int, int](maxWeight, nil)

	require.Equal(t, uint64(0), c.Weight())
	require.Equal(t, 0, c.Size())

	expectedValues := make(map[int]int)

	// Fill up the cache. Everything should have weight 1.
	for i := 1; i <= int(maxWeight); i++ {

		value := rand.Int()
		expectedValues[i] = value

		// The value shouldn't be present yet
		v, ok := c.Get(i)
		require.False(t, ok)
		require.Equal(t, 0, v)

		c.Put(i, value)

		require.Equal(t, uint64(i), c.Weight())
		require.Equal(t, i, c.Size())
	}

	// Verify that all expected values are present.
	for k, v := range expectedValues {
		value, ok := c.Get(k)
		require.True(t, ok)
		require.Equal(t, v, value)
	}

	// Push the old values out of the queue one at a time.
	for i := 1; i <= int(maxWeight); i++ {
		value := rand.Int()
		expectedValues[-i] = value
		delete(expectedValues, i)

		// The value shouldn't be present yet
		v, ok := c.Get(-i)
		require.False(t, ok)
		require.Equal(t, 0, v)

		c.Put(-i, value)

		require.Equal(t, maxWeight, c.Weight())
		require.Equal(t, int(maxWeight), c.Size())

		// verify that the purged value is specifically not present
		_, ok = c.Get(i)
		require.False(t, ok)

		// verify that only the expected values have been purged. Has the added benefit of randomly
		// reading all the values in the cache, which for a FIFO cache should not influence the order
		// that we purge values.
		for kk, vv := range expectedValues {
			value, ok = c.Get(kk)
			require.True(t, ok)
			require.Equal(t, vv, value)
		}
	}
}

func TestWeightedValues(t *testing.T) {
	tu.InitializeRandom()

	maxWeight := uint64(100 + rand.Intn(100))

	// For this test, weight is simply the key.
	weightCalculator := func(key int, value int) uint64 {
		return uint64(key)
	}

	c := NewFIFOCache[int, int](maxWeight, weightCalculator)

	expectedValues := make(map[int]int)

	require.Equal(t, uint64(0), c.Weight())
	require.Equal(t, 0, c.Size())

	highestUndeletedKey := 0
	expectedWeight := uint64(0)
	for nextKey := 0; nextKey <= int(maxWeight); nextKey++ {
		value := rand.Int()
		c.Put(nextKey, value)
		expectedValues[nextKey] = value
		expectedWeight += uint64(nextKey)

		// simulate the expected removal
		for expectedWeight > maxWeight {
			delete(expectedValues, highestUndeletedKey)
			expectedWeight -= uint64(highestUndeletedKey)
			highestUndeletedKey++
		}

		require.Equal(t, expectedWeight, c.Weight())
		require.Equal(t, len(expectedValues), c.Size())

		// Update a random existing key. Shouldn't affect the weight or removal order.
		for k := range expectedValues {
			value = rand.Int()
			c.Put(k, value)
			expectedValues[k] = value
			break
		}

		// verify that all expected values are present
		for k, v := range expectedValues {
			var ok bool
			value, ok = c.Get(k)
			require.True(t, ok)
			require.Equal(t, v, value)
		}
	}

	// Attempting to insert a value that exceeds the max weight should have no effect.
	c.Put(int(maxWeight)+1, rand.Int())

	for k, v := range expectedValues {
		value, ok := c.Get(k)
		require.True(t, ok)
		require.Equal(t, v, value)
	}
}
