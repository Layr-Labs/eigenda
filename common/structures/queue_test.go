package structures_test

import (
	"testing"

	"github.com/Layr-Labs/eigenda/common/structures"
	"github.com/Layr-Labs/eigenda/test/random"
	"github.com/stretchr/testify/require"
)

// A simple implementation of a queue for testing purposes. It's slow, but easy to reason about.
type simpleQueue[T any] struct {
	data []T
}

func newSimpleQueue[T any]() *simpleQueue[T] {
	return &simpleQueue[T]{
		data: make([]T, 0),
	}
}

func (q *simpleQueue[T]) Push(item T) {
	q.data = append(q.data, item)
}

func (q *simpleQueue[T]) Pop() (T, bool) {
	if len(q.data) == 0 {
		var zero T
		return zero, false
	}
	item := q.data[0]
	q.data = q.data[1:]
	return item, true
}

func (q *simpleQueue[T]) Size() uint64 {
	return uint64(len(q.data))
}

func (q *simpleQueue[T]) Peek() (T, bool) {
	if len(q.data) == 0 {
		var zero T
		return zero, false
	}
	return q.data[0], true
}

func (q *simpleQueue[T]) Clear() {
	q.data = make([]T, 0)
}

func (q *simpleQueue[T]) Get(index int) T {
	if index < 0 || index >= len(q.data) {
		panic("index out of bounds")
	}
	return q.data[index]
}

func (q *simpleQueue[T]) Set(index int, value T) (T, bool) {
	if index < 0 || index >= len(q.data) {
		var zero T
		return zero, false
	}
	old := q.data[index]
	q.data[index] = value
	return old, true
}

func TestRandomQueueOperations(t *testing.T) {
	rand := random.NewTestRandom()

	initialSize := rand.Uint64Range(0, 8)

	queue := structures.NewQueue[int](initialSize)

	// Iterating an empty queue should work as expected
	for range queue.Iterator() {
		t.Fail()
	}

	// Use a simple queue implementation we trust to verify correctness.
	expectedData := newSimpleQueue[int]()
	expectedSize := uint64(0)

	operationCount := 10_000
	for i := 0; i < operationCount; i++ {

		// Do a random mutation.
		choice := rand.Float64()

		// nolint:nestif
		if choice < 0.001 {
			// ~0.1% chance
			// clear

			queue.Clear()
			expectedData.Clear()
			expectedSize = 0

		} else if choice < 0.5 {
			// ~50% chance
			// Push to the queue (enqueue)

			value := rand.Int()
			queue.Push(value)
			expectedData.Push(value)

			expectedSize++
		} else if choice < 0.9 {
			// ~40% chance
			// Pop from the queue (dequeue)

			if expectedSize == 0 {
				_, ok := queue.TryPop()
				require.False(t, ok)
			} else {
				value, ok := queue.TryPop()
				require.True(t, ok)

				expectedValue, expectedOk := expectedData.Peek()
				require.True(t, expectedOk)
				_, _ = expectedData.Pop()

				require.Equal(t, expectedValue, value)

				expectedSize--
			}
		} else {
			// ~10% chance
			// Set a random index

			if expectedSize == 0 {
				// Setting on empty queue should panic
				require.Panics(t, func() {
					queue.Set(0, rand.Int())
				})
				require.Panics(t, func() {
					queue.Set(rand.Uint64(), rand.Int())
				})
			} else {
				index := 0
				if expectedSize > 2 {
					index = rand.Intn(int(expectedSize - 1))
				}

				newValue := rand.Int()

				expectedOldValue := expectedData.Get(index)
				expectedData.Set(index, newValue)

				oldValue := queue.Set(uint64(index), newValue)

				require.Equal(t, expectedOldValue, oldValue)
			}
		}

		// Always check things that are fast to check.
		require.Equal(t, expectedSize, queue.Size(), "size mismatch after %d operations", i)
		require.Equal(t, expectedSize == 0, queue.IsEmpty())

		if expectedSize == 0 {
			_, ok := queue.TryPeek()
			require.False(t, ok)
			_, ok = queue.TryPop()
			require.False(t, ok)

			// Verify panicking operations
			require.Panics(t, func() { queue.Peek() })
			require.Panics(t, func() { queue.Pop() })
			require.Panics(t, func() { queue.Get(0) })
			require.Panics(t, func() { queue.Get(rand.Uint64()) })
		} else {
			expected, ok := expectedData.Peek()
			require.True(t, ok)
			actual, actualOk := queue.TryPeek()
			require.True(t, actualOk)
			require.Equal(t, expected, actual)
		}

		// nolint:nestif
		if i%1000 == 0 {
			// Once in a while, verify the entire contents of the queue. It's expensive to do this in every iteration.

			if expectedData.Size() > 0 {
				// Verify a random index.
				index := 0
				if expectedData.Size() > 2 {
					index = rand.Intn(int(expectedData.Size()) - 1)
				}
				value := queue.Get(uint64(index))
				require.Equal(t, expectedData.Get(index), value)

				// Verify out-of-bounds access panics
				require.Panics(t, func() { queue.Get(expectedSize) })
				require.Panics(t, func() { queue.Get(expectedSize + rand.Uint64Range(1, 100)) })
				require.Panics(t, func() { queue.Set(expectedSize, rand.Int()) })
			}

			// Iterate forwards
			expectedIndex := 0
			for index, value := range queue.Iterator() {
				require.Equal(t, uint64(expectedIndex), index)
				expectedIndex++

				require.True(t, index < expectedData.Size())

				require.Equal(t, expectedData.Get(int(index)), value)
			}
			require.Equal(t, expectedData.Size(), uint64(expectedIndex), "forward iteration count mismatch")
		}
	}
}
