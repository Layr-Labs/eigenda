package chunkgroup

import (
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/lightnode"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
	"testing"
	"time"
)

func randomAssignment(nextShuffleTime time.Time) *assignment {
	id := rand.Uint64()
	seed := rand.Uint64()
	registrationTime := time.Unix(int64(rand.Uint32()), 0)

	registration := lightnode.NewRegistration(id, seed, registrationTime)

	return &assignment{
		registration: registration,
		endOfEpoch:   nextShuffleTime,
		chunkGroup:   rand.Uint64(),
	}
}

func TestEmptyQueue(t *testing.T) {
	queue := newAssignmentQueue()
	assert.Equal(t, uint(0), queue.Size())
	assert.Nil(t, queue.Pop())
	assert.Nil(t, queue.Peek())
	assert.Equal(t, uint(0), queue.Size())
}

func TestInOrderInsertion(t *testing.T) {
	tu.InitializeRandom()

	queue := newAssignmentQueue()
	assert.Equal(t, uint(0), queue.Size())

	startTime := time.Unix(int64(rand.Uint32()), 0)
	numberOfElements := uint(100)
	expectedOrder := make([]*assignment, 0, numberOfElements)

	// Insert elements in order.
	for i := uint(0); i < numberOfElements; i++ {
		registration := randomAssignment(startTime.Add(time.Second * time.Duration(i)))
		expectedOrder = append(expectedOrder, registration)
		queue.Push(registration)

		// Pushing more than once should be a no-op.
		pushCount := rand.Intn(3)
		for j := 0; j < pushCount; j++ {
			queue.Push(registration)
		}

		assert.Equal(t, i+1, queue.Size())
	}

	// Pop elements in order.
	for i := uint(0); i < numberOfElements; i++ {
		preview := queue.Peek()
		assert.Equal(t, expectedOrder[i], preview)
		assert.Equal(t, numberOfElements-i, queue.Size())

		next := queue.Pop()
		assert.Equal(t, expectedOrder[i], next)
		assert.Equal(t, numberOfElements-i-1, queue.Size())
	}
}

func TestReverseOrderInsertion(t *testing.T) {
	tu.InitializeRandom()

	queue := newAssignmentQueue()
	assert.Equal(t, uint(0), queue.Size())

	startTime := time.Unix(int64(rand.Uint32()), 0)
	numberOfElements := uint(100)
	expectedOrder := make([]*assignment, 0, numberOfElements)

	// Generate the elements that will eventually be inserted.
	for i := uint(0); i < numberOfElements; i++ {
		assignment := randomAssignment(startTime.Add(time.Second * time.Duration(i)))
		expectedOrder = append(expectedOrder, assignment)
	}

	// Insert elements in reverse order.
	for i := int(numberOfElements - 1); i >= 0; i-- {
		queue.Push(expectedOrder[i])

		// Pushing more than once should be a no-op.
		pushCount := rand.Intn(3)
		for j := 0; j < pushCount; j++ {
			queue.Push(expectedOrder[i])
		}

		assert.Equal(t, numberOfElements-uint(i), queue.Size())
	}

	// Pop elements in order.
	for i := uint(0); i < numberOfElements; i++ {
		preview := queue.Peek()
		assert.Equal(t, expectedOrder[i], preview)
		assert.Equal(t, numberOfElements-i, queue.Size())

		next := queue.Pop()
		assert.Equal(t, expectedOrder[i], next)
		assert.Equal(t, numberOfElements-i-1, queue.Size())
	}
}

func TestRandomInsertion(t *testing.T) {
	tu.InitializeRandom()

	queue := newAssignmentQueue()
	assert.Equal(t, uint(0), queue.Size())

	startTime := time.Unix(int64(rand.Uint32()), 0)
	numberOfElements := uint(100)
	expectedOrder := make([]*assignment, 0, numberOfElements)

	// Generate the elements that will eventually be inserted.
	for i := uint(0); i < numberOfElements; i++ {
		assignment := randomAssignment(startTime.Add(time.Second * time.Duration(i)))
		expectedOrder = append(expectedOrder, assignment)
	}

	// Insert elements in random order.
	perm := rand.Perm(int(numberOfElements))
	for i := uint(0); i < numberOfElements; i++ {

		queue.Push(expectedOrder[perm[i]])

		// Pushing more than once should be a no-op.
		pushCount := rand.Intn(3)
		for j := 0; j < pushCount; j++ {
			queue.Push(expectedOrder[perm[i]])
		}

		assert.Equal(t, i+1, queue.Size())
	}

	// Pop elements in order.
	for i := uint(0); i < numberOfElements; i++ {
		preview := queue.Peek()
		assert.Equal(t, expectedOrder[i], preview)
		assert.Equal(t, numberOfElements-i, queue.Size())

		next := queue.Pop()
		assert.Equal(t, expectedOrder[i], next)
		assert.Equal(t, numberOfElements-i-1, queue.Size())
	}
}

func TestPeriodicRemoval(t *testing.T) {
	tu.InitializeRandom()

	queue := newAssignmentQueue()
	assert.Equal(t, uint(0), queue.Size())

	startTime := time.Unix(int64(rand.Uint32()), 0)
	numberOfElements := uint(100)
	expectedOrder := make([]*assignment, 0, numberOfElements)
	expectedOrderWithRemovals := make([]*assignment, 0, numberOfElements)

	// Generate the elements that will eventually be inserted.
	for i := uint(0); i < numberOfElements; i++ {
		assignment := randomAssignment(startTime.Add(time.Second * time.Duration(i)))
		expectedOrder = append(expectedOrder, assignment)

		// We will remove every 7th element.
		if i%7 != 0 {
			expectedOrderWithRemovals = append(expectedOrderWithRemovals, assignment)
		}
	}

	// Insert elements in random order.
	perm := rand.Perm(int(numberOfElements))
	for i := uint(0); i < numberOfElements; i++ {

		queue.Push(expectedOrder[perm[i]])

		// Pushing more than once should be a no-op.
		pushCount := rand.Intn(3)
		for j := 0; j < pushCount; j++ {
			queue.Push(expectedOrder[perm[i]])
		}

		assert.Equal(t, i+1, queue.Size())
	}

	removalCount := uint(0)

	// Remove every 7th element.
	for i := uint(0); i < numberOfElements; i++ {
		if i%7 == 0 {
			queue.Remove(expectedOrder[i].registration.ID())

			// Removing more than once should be a no-op.
			removeCount := rand.Intn(3)
			for j := 0; j < removeCount; j++ {
				queue.Remove(expectedOrder[i].registration.ID())
			}

			removalCount++
			assert.Equal(t, numberOfElements-removalCount, queue.Size())
		}
	}

	// Pop elements in order. We shouldn't see the removed elements.
	for i := uint(0); i < (numberOfElements - removalCount); i++ {
		preview := queue.Peek()
		assert.Equal(t, expectedOrderWithRemovals[i], preview)
		assert.Equal(t, numberOfElements-i-removalCount, queue.Size())

		next := queue.Pop()
		assert.Equal(t, expectedOrderWithRemovals[i], next)
		assert.Equal(t, numberOfElements-i-removalCount-1, queue.Size())
	}
}

func TestContiguousRemoval(t *testing.T) {
	tu.InitializeRandom()

	queue := newAssignmentQueue()
	assert.Equal(t, uint(0), queue.Size())

	startTime := time.Unix(int64(rand.Uint32()), 0)
	numberOfElements := uint(100)
	expectedOrder := make([]*assignment, 0, numberOfElements)
	expectedOrderWithRemovals := make([]*assignment, 0, numberOfElements)

	// Generate the elements that will eventually be inserted.
	for i := uint(0); i < numberOfElements; i++ {
		assignment := randomAssignment(startTime.Add(time.Second * time.Duration(i)))
		expectedOrder = append(expectedOrder, assignment)

		// We will remove all elements after index 10 and before index 90
		if i <= 10 || i >= 90 {
			expectedOrderWithRemovals = append(expectedOrderWithRemovals, assignment)
		}
	}

	// Insert elements in random order.
	perm := rand.Perm(int(numberOfElements))
	for i := uint(0); i < numberOfElements; i++ {

		queue.Push(expectedOrder[perm[i]])

		// Pushing more than once should be a no-op.
		pushCount := rand.Intn(3)
		for j := 0; j < pushCount; j++ {
			queue.Push(expectedOrder[perm[i]])
		}

		assert.Equal(t, i+1, queue.Size())
	}

	removalCount := uint(0)

	// Remove all elements after index 10 and before index 90
	for i := uint(0); i < numberOfElements; i++ {
		if i > 10 && i < 90 {
			queue.Remove(expectedOrder[i].registration.ID())

			// Removing more than once should be a no-op.
			removeCount := rand.Intn(3)
			for j := 0; j < removeCount; j++ {
				queue.Remove(expectedOrder[i].registration.ID())
			}

			removalCount++
			assert.Equal(t, numberOfElements-removalCount, queue.Size())
		}
	}

	// Pop elements in order. We shouldn't see the removed elements.
	for i := uint(0); i < (numberOfElements - removalCount); i++ {
		preview := queue.Peek()
		assert.Equal(t, expectedOrderWithRemovals[i], preview)
		assert.Equal(t, numberOfElements-i-removalCount, queue.Size())

		next := queue.Pop()
		assert.Equal(t, expectedOrderWithRemovals[i], next)
		assert.Equal(t, numberOfElements-i-removalCount-1, queue.Size())
	}
}

func TestRemoveFollowedByPush(t *testing.T) {
	tu.InitializeRandom()

	queue := newAssignmentQueue()
	assert.Equal(t, uint(0), queue.Size())

	startTime := time.Unix(int64(rand.Uint32()), 0)
	numberOfElements := uint(100)
	expectedOrder := make([]*assignment, 0, numberOfElements)

	// Generate the elements that will eventually be inserted.
	for i := uint(0); i < numberOfElements; i++ {
		assignment := randomAssignment(startTime.Add(time.Second * time.Duration(i)))
		expectedOrder = append(expectedOrder, assignment)
	}

	// Insert elements in random order.
	perm := rand.Perm(int(numberOfElements))
	for i := uint(0); i < numberOfElements; i++ {

		queue.Push(expectedOrder[perm[i]])

		// Pushing more than once should be a no-op.
		pushCount := rand.Intn(3)
		for j := 0; j < pushCount; j++ {
			queue.Push(expectedOrder[perm[i]])
		}

		assert.Equal(t, i+1, queue.Size())
	}

	removalCount := uint(0)

	// Remove every seventh element.
	for i := uint(0); i < numberOfElements; i++ {
		if i%7 == 0 {
			queue.Remove(expectedOrder[i].registration.ID())

			// Removing more than once should be a no-op.
			removeCount := rand.Intn(3)
			for j := 0; j < removeCount; j++ {
				queue.Remove(expectedOrder[i].registration.ID())
			}

			removalCount++
			assert.Equal(t, numberOfElements-removalCount, queue.Size())
		}
	}

	// Push the removed nodes back into the queue.
	for i := uint(0); i < numberOfElements; i++ {
		if i%7 == 0 {
			queue.Push(expectedOrder[i])
			removalCount--
			assert.Equal(t, numberOfElements-removalCount, queue.Size())
		}
	}

	// Pop elements in order.
	for i := uint(0); i < (numberOfElements - removalCount); i++ {
		preview := queue.Peek()
		assert.Equal(t, expectedOrder[i], preview)
		assert.Equal(t, numberOfElements-i, queue.Size())

		next := queue.Pop()
		assert.Equal(t, expectedOrder[i], next)
		assert.Equal(t, numberOfElements-i-1, queue.Size())
	}
}
