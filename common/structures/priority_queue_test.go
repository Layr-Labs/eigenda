package structures

import (
	"math/rand"
	"slices"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

// Note: I can't use the normal test random utility in this file due to a circular dependency

func TestInsertThenRemove(t *testing.T) {
	count := 1024

	values := make([]int, count)
	pq := NewPriorityQueue[int](func(a, b int) bool {
		return a < b
	})

	for i := 0; i < count; i++ {
		next := rand.Intn(10000)
		values[i] = next
		pq.Push(next)

		require.Equal(t, i+1, pq.Size())
	}

	// sort the values into the order we expect to see them come out of the priority queue
	slices.Sort(values)

	previous := -1
	for i := 0; i < count; i++ {
		require.Equal(t, values[i], pq.Peek())

		value, ok := pq.TryPeek()
		require.True(t, ok)
		require.Equal(t, values[i], value)

		require.Equal(t, count-i, pq.Size())

		if i%2 == 0 {
			value = pq.Pop()
			require.Equal(t, values[i], value)
		} else {
			var ok bool
			value, ok = pq.TryPop()
			require.True(t, ok)
			require.Equal(t, values[i], value)
		}
		require.GreaterOrEqual(t, value, previous)
		previous = value

		require.Equal(t, count-i-1, pq.Size())
	}

	_, ok := pq.TryPop()
	require.False(t, ok)
}

func TestIteration(t *testing.T) {
	count := 1024

	values := make([]int, count)
	pq := NewPriorityQueue[int](func(a, b int) bool {
		return a < b
	})

	for i := 0; i < count; i++ {
		next := rand.Intn(10000)
		values[i] = next
		pq.Push(next)
	}

	// sort the values into the order we expect to see them come out of the priority queue
	slices.Sort(values)

	index := 0
	for item := range pq.PopIterator() {
		require.Equal(t, values[index], item)
		index++
	}
	require.Equal(t, count, index)
	require.Equal(t, 0, pq.Size())

}

func TestRandomOperations(t *testing.T) {
	count := 256

	values := make([]int, 0, count)
	pq := NewPriorityQueue[int](func(a, b int) bool {
		return a < b
	})

	for i := 0; i < count; i++ {

		choice := rand.Float64()

		if choice < 0.6 || len(values) == 0 {
			// insert
			next := rand.Intn(10000)
			values = append(values, next)
			sort.Ints(values)
			pq.Push(next)
		} else {
			// remove
			expected := values[0]
			values = values[1:]

			value, ok := pq.TryPop()
			require.True(t, ok)
			require.Equal(t, expected, value)
		}
	}

	// pop remaining items
	for i := 0; i < len(values); i++ {
		expected := values[i]
		value, ok := pq.TryPop()
		require.True(t, ok)
		require.Equal(t, expected, value)
	}

	require.Equal(t, 0, pq.Size())
}
