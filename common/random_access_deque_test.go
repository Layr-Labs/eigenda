package common_test

import (
	"os"
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/stretchr/testify/require"

	"github.com/emirpasic/gods/lists/doublylinkedlist"
)

func TestRandomDequeOperations(t *testing.T) {
	rand := random.NewTestRandom()

	// The gods implementation of a doubly linked has a fmt.Println()... sigh. Squelch standard out during this test.
	null, _ := os.Open(os.DevNull)
	old := os.Stdout
	os.Stdout = null
	defer func() {
		_ = null.Close()
		os.Stdout = old
	}()

	initialSize := rand.Uint64Range(0, 8)

	deque := common.NewRandomAccessDeque[int](initialSize)

	// Use a linked list library we trust to verify correctness. The linked list can't do O(1) index access, but we can
	// work around that in the test code.
	expectedData := doublylinkedlist.New()
	expectedSize := uint64(0)

	operationCount := 10_000
	for i := 0; i < operationCount; i++ {

		// Do a random mutation.
		choice := rand.Float64()

		// nolint:nestif
		if choice < 0.001 {
			// ~1 time per 1000 operations
			// clear

			deque.Clear()
			expectedData.Clear()
			expectedSize = 0

		} else if choice < 0.25 {
			// ~25% chance
			// Add to the front

			value := rand.Int()
			deque.PushFront(value)
			expectedData.Insert(0, value)

			expectedSize++
		} else if choice < 0.5 {
			// ~25% chance
			// Add to the back

			value := rand.Int()
			deque.PushBack(value)
			expectedData.Add(value)

			expectedSize++
		} else if choice < 0.7 {
			// ~20% chance
			// Remove from the front

			if expectedSize == 0 {
				_, err := deque.PopFront()
				require.Error(t, err)
			} else {
				value, err := deque.PopFront()
				require.NoError(t, err)

				expectedValue, ok := expectedData.Get(0)
				require.True(t, ok)
				expectedData.Remove(0)

				require.Equal(t, expectedValue, value)

				expectedSize--
			}
		} else if choice < 0.9 {
			// ~20% chance
			// remove from the back

			if expectedSize == 0 {
				_, err := deque.PopBack()
				require.Error(t, err)
			} else {
				value, err := deque.PopBack()
				require.NoError(t, err)

				expectedValue, ok := expectedData.Get(expectedData.Size() - 1)
				require.True(t, ok)
				expectedData.Remove(expectedData.Size() - 1)

				require.Equal(t, expectedValue, value)

				expectedSize--
			}
		} else {
			// ~10% chance
			// set a random index

			if expectedSize == 0 {
				_, err := deque.Set(0, rand.Int())
				require.Error(t, err)
				_, err = deque.Set(rand.Uint64(), rand.Int())
				require.Error(t, err)
			} else {
				index := 0
				if expectedSize > 2 {
					index = rand.Intn(int(expectedSize - 1))
				}

				newValue := rand.Int()

				// This is O(2 * n)... hard to test this efficiently without a trusted reference implementation
				// that supports O(1) index access. ;(
				expectedOldValue, ok := expectedData.Get(index)
				require.True(t, ok)
				expectedData.Set(index, newValue)

				oldValue, err := deque.Set(uint64(index), newValue)
				require.NoError(t, err)

				require.Equal(t, expectedOldValue, oldValue)
			}
		}

		// Always check things that are fast to check.
		require.Equal(t, expectedSize, deque.Size(), "size mismatch after %d operations", i)
		if expectedSize == 0 {
			_, err := deque.PeekFront()
			require.Error(t, err)
			_, err = deque.PeekBack()
			require.Error(t, err)
			_, err = deque.PopFront()
			require.Error(t, err)
			_, err = deque.PopBack()
			require.Error(t, err)
			_, err = deque.Get(0)
			require.Error(t, err)
			_, err = deque.Get(rand.Uint64())
			require.Error(t, err)
			_, err = deque.Set(0, rand.Int())
			require.Error(t, err)
			_, err = deque.Set(rand.Uint64(), rand.Int())
			require.Error(t, err)
		} else {
			expected, ok := expectedData.Get(0)
			require.True(t, ok)
			actual, err := deque.PeekFront()
			require.NoError(t, err)
			require.Equal(t, expected, actual)

			expected, ok = expectedData.Get(expectedData.Size() - 1)
			require.True(t, ok)
			actual, err = deque.PeekBack()
			require.NoError(t, err)
			require.Equal(t, expected, actual)
		}

		if i%1000 == 0 {
			// Once in a while, verify the entire contents of the deque. It's expensive to do this in every iteration.

			// Create a copy of the expected data for efficient verification.
			expectedArray := make([]int, 0, expectedData.Size())
			expectedData.Each(func(index int, value interface{}) {
				expectedArray = append(expectedArray, value.(int))
			})

			if expectedData.Size() > 0 {
				// Verify a random index.
				index := 0
				if expectedData.Size() > 2 {
					index = rand.Intn(expectedData.Size() - 1)
				}
				value, err := deque.Get(uint64(index))
				require.NoError(t, err)
				require.Equal(t, expectedArray[index], value)
			}

			// Iterate forwards
			expectedIndex := 0
			for index, value := range deque.Iterator() {
				require.Equal(t, expectedIndex, index)
				expectedIndex++

				require.True(t, index < expectedData.Size())

				require.Equal(t, expectedArray[index], value)
			}
			require.Equal(t, expectedData.Size(), expectedIndex, "forward iteration count mismatch")

			// Iterate backwards
			expectedIndex = expectedData.Size() - 1
			for index, value := range deque.ReverseIterator() {
				require.Equal(t, expectedIndex, index)
				expectedIndex--

				require.True(t, index >= 0)

				require.Equal(t, expectedArray[index], value)
			}
			require.Equal(t, -1, expectedIndex, "backward iteration count mismatch")
		}
	}
}
