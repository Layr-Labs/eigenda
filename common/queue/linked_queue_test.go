package queue

import (
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
	"testing"
)

func TestEmptyQueue(t *testing.T) {
	var q LinkedQueue[int]
	require.Equal(t, 0, q.Size())

	next, ok := q.Peek()
	require.False(t, ok)
	require.Equal(t, 0, next)

	next, ok = q.Pop()
	require.False(t, ok)
	require.Equal(t, 0, next)

	require.Equal(t, 0, q.Size())
}

func TestRandomOperations(t *testing.T) {
	tu.InitializeRandom()

	var q LinkedQueue[int]
	expectedValues := make([]int, 0)

	for i := 0; i < 1000; i++ {
		if rand.Int()%2 == 0 || len(expectedValues) == 0 {
			// push an item
			itemToPush := rand.Int()
			q.Push(itemToPush)
			expectedValues = append(expectedValues, itemToPush)
		} else {
			// pop an item

			next, ok := q.Pop()
			expectedNext := expectedValues[0]
			expectedValues = expectedValues[1:]

			require.True(t, ok)
			require.Equal(t, expectedNext, next)
		}

		require.Equal(t, len(expectedValues), q.Size())

		next, ok := q.Peek()
		if len(expectedValues) == 0 {
			require.False(t, ok)
		} else {
			require.True(t, ok)
			require.Equal(t, expectedValues[0], next)
		}
	}
}
