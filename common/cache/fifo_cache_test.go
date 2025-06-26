package cache

import (
	"testing"

	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/stretchr/testify/require"
)

// A function that builds a cache
type cacheBuilder func(maxWeight uint64, calculator func(key int, value int) uint64) Cache[int, int]

func expirationOrderTest(t *testing.T, builder cacheBuilder) {
	rand := random.NewTestRandom()

	maxWeight := uint64(10 + rand.Intn(10))
	c := builder(maxWeight, nil)

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

func TestExpirationOrder(t *testing.T) {
	t.Run("FIFO", func(t *testing.T) {
		builder := func(maxWeight uint64, calculator func(key int, value int) uint64) Cache[int, int] {
			return NewFIFOCache[int, int](maxWeight, calculator, nil)
		}
		expirationOrderTest(t, builder)
	})
	t.Run("Thread Safe FIFO", func(t *testing.T) {
		builder := func(maxWeight uint64, calculator func(key int, value int) uint64) Cache[int, int] {
			base := NewFIFOCache[int, int](maxWeight, calculator, nil)
			return NewThreadSafeCache[int, int](base)
		}
		expirationOrderTest(t, builder)
	})
	t.Run("Weak FIFO", func(t *testing.T) {
		// We are using low enough memory that it is unlikely that the weak pointers
		// will be garbage collected during the course of this test.

		builder := func(maxWeight uint64, calculator func(key int, value int) uint64) Cache[int, int] {
			return NewWeakFIFOCache[int, int](maxWeight, calculator, nil)
		}
		expirationOrderTest(t, builder)
	})
}

func weightedValuesTest(t *testing.T, builder cacheBuilder) {
	rand := random.NewTestRandom()

	maxWeight := uint64(100 + rand.Intn(100))

	// For this test, weight is simply the key.
	weightCalculator := func(key int, value int) uint64 {
		return uint64(key)
	}

	c := NewFIFOCache[int, int](maxWeight, weightCalculator, nil)

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

func TestWeightedValues(t *testing.T) {
	t.Run("FIFO", func(t *testing.T) {
		builder := func(maxWeight uint64, calculator func(key int, value int) uint64) Cache[int, int] {
			return NewFIFOCache[int, int](maxWeight, calculator, nil)
		}
		weightedValuesTest(t, builder)
	})
	t.Run("Thread Safe FIFO", func(t *testing.T) {
		builder := func(maxWeight uint64, calculator func(key int, value int) uint64) Cache[int, int] {
			base := NewFIFOCache[int, int](maxWeight, calculator, nil)
			return NewThreadSafeCache[int, int](base)
		}
		weightedValuesTest(t, builder)
	})
	t.Run("Weak FIFO", func(t *testing.T) {
		// We are using low enough memory that it is unlikely that the weak pointers
		// will be garbage collected during the course of this test.

		builder := func(maxWeight uint64, calculator func(key int, value int) uint64) Cache[int, int] {
			return NewWeakFIFOCache[int, int](maxWeight, calculator, nil)
		}
		weightedValuesTest(t, builder)
	})
}

func reinsertionTest(t *testing.T, builder cacheBuilder) {
	rand := random.NewTestRandom()

	maxWeight := uint64(100 + rand.Intn(100))

	// For this test, weight is simply the key.
	weightCalculator := func(key int, value int) uint64 {
		return uint64(key)
	}

	c := NewFIFOCache[int, int](maxWeight, weightCalculator, nil)

	expectedValues := make(map[int]int)

	require.Equal(t, uint64(0), c.Weight())
	require.Equal(t, 0, c.Size())

	highestUndeletedKey := 0
	expectedWeight := uint64(0)
	var nextKey int
	for ; nextKey <= int(maxWeight); nextKey++ {

		expectedWeight += uint64(nextKey)
		if expectedWeight > maxWeight {
			// Don't add enough data to trigger an eviction yet.
			break
		}

		value := rand.Int()
		c.Put(nextKey, value)
		expectedValues[nextKey] = value

		// simulate the expected removal
		for expectedWeight > maxWeight {
			delete(expectedValues, highestUndeletedKey)
			expectedWeight -= uint64(highestUndeletedKey)
			highestUndeletedKey++
		}

		require.Equal(t, expectedWeight, c.Weight())
		require.Equal(t, len(expectedValues), c.Size())

		// verify that all expected values are present
		for k, v := range expectedValues {
			var ok bool
			value, ok = c.Get(k)
			require.True(t, ok)
			require.Equal(t, v, value)
		}
	}

	// Reinsert value 0. It is currently the first value scheduled to be garbage collected, but this should move
	// it to the end of the queue.
	value := rand.Int()
	c.Put(0, value)
	expectedValues[0] = value

	// Insert a value with a weight that will fill up all capacity all by itself. If key 0 is at the front of the GC
	// queue, then we'd expect for it to be a casualty of this operation. If it is correctly at the back of the queue
	// now, then it will not be evicted (since key 0 has a weight of 0).
	bigKey := int(maxWeight)
	value = rand.Int()
	expectedValues[bigKey] = value
	c.Put(bigKey, value)

	for k, v := range expectedValues {
		value, ok := c.Get(k)

		if k == 0 || k == bigKey {
			// There should only be room for the big key and key 0
			require.True(t, ok)
			require.Equal(t, v, value)
		} else {
			// All other keys should have been evicted
			require.False(t, ok)
		}
	}
}

func TestReinsertion(t *testing.T) {
	t.Run("FIFO", func(t *testing.T) {
		builder := func(maxWeight uint64, calculator func(key int, value int) uint64) Cache[int, int] {
			return NewFIFOCache[int, int](maxWeight, calculator, nil)
		}
		reinsertionTest(t, builder)
	})
	t.Run("Thread Safe FIFO", func(t *testing.T) {
		builder := func(maxWeight uint64, calculator func(key int, value int) uint64) Cache[int, int] {
			base := NewFIFOCache[int, int](maxWeight, calculator, nil)
			return NewThreadSafeCache[int, int](base)
		}
		reinsertionTest(t, builder)
	})
	t.Run("Weak FIFO", func(t *testing.T) {
		// We are using low enough memory that it is unlikely that the weak pointers
		// will be garbage collected during the course of this test.

		builder := func(maxWeight uint64, calculator func(key int, value int) uint64) Cache[int, int] {
			return NewWeakFIFOCache[int, int](maxWeight, calculator, nil)
		}
		reinsertionTest(t, builder)
	})
}
