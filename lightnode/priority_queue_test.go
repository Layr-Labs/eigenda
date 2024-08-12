package lightnode

import (
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
	"testing"
	"time"
)

func RandomRegistration(nextShuffleTime time.Time) *Registration {
	id := rand.Uint64()
	seed := rand.Uint64()
	registrationTime := time.Unix(int64(rand.Uint32()), 0)

	registration := NewRegistration(id, seed, registrationTime)
	registration.nextShuffleTime = nextShuffleTime

	return registration
}

func TestEmptyQueue(t *testing.T) {
	queue := NewPriorityQueue()
	assert.Equal(t, 0, queue.Len())
	assert.Nil(t, queue.Pop())
	assert.Nil(t, queue.Peek())
	assert.Equal(t, 0, queue.Len())
}

func TestInOrderInsertion(t *testing.T) {
	tu.InitializeRandom()

	queue := NewPriorityQueue()
	assert.Equal(t, 0, queue.Len())

	startTime := time.Unix(int64(rand.Uint32()), 0)
	numberOfElements := 100
	expectedOrder := make([]*Registration, 0, numberOfElements)

	// Insert elements in order.
	for i := 0; i < numberOfElements; i++ {
		registration := RandomRegistration(startTime.Add(time.Second * time.Duration(i)))
		expectedOrder = append(expectedOrder, registration)
		queue.Push(registration)
		assert.Equal(t, i+1, queue.Len())
	}

	// Pop elements in order.
	for i := 0; i < numberOfElements; i++ {
		preview := queue.Peek()
		assert.Equal(t, expectedOrder[i], preview)
		assert.Equal(t, numberOfElements-i, queue.Len())

		next := queue.Pop()
		assert.Equal(t, expectedOrder[i], next)
		assert.Equal(t, numberOfElements-i-1, queue.Len())
	}
}

func TestReverseOrderInsertion(t *testing.T) {
	tu.InitializeRandom()

	queue := NewPriorityQueue()
	assert.Equal(t, 0, queue.Len())

	startTime := time.Unix(int64(rand.Uint32()), 0)
	numberOfElements := 100
	expectedOrder := make([]*Registration, 0, numberOfElements)

	// Generate the elements that will eventually be inserted.
	for i := 0; i < numberOfElements; i++ {
		registration := RandomRegistration(startTime.Add(time.Second * time.Duration(i)))
		expectedOrder = append(expectedOrder, registration)
	}

	// Insert elements in reverse order.
	for i := numberOfElements - 1; i >= 0; i-- {
		queue.Push(expectedOrder[i])
		assert.Equal(t, numberOfElements-i, queue.Len())
	}

	// Pop elements in order.
	for i := 0; i < numberOfElements; i++ {
		preview := queue.Peek()
		assert.Equal(t, expectedOrder[i], preview)
		assert.Equal(t, numberOfElements-i, queue.Len())

		next := queue.Pop()
		assert.Equal(t, expectedOrder[i], next)
		assert.Equal(t, numberOfElements-i-1, queue.Len())
	}
}

func TestRandomInsertion(t *testing.T) {
	tu.InitializeRandom()

	queue := NewPriorityQueue()
	assert.Equal(t, 0, queue.Len())

	startTime := time.Unix(int64(rand.Uint32()), 0)
	numberOfElements := 100
	expectedOrder := make([]*Registration, 0, numberOfElements)

	// Generate the elements that will eventually be inserted.
	for i := 0; i < numberOfElements; i++ {
		registration := RandomRegistration(startTime.Add(time.Second * time.Duration(i)))
		expectedOrder = append(expectedOrder, registration)
	}

	// Insert elements in random order.
	perm := rand.Perm(numberOfElements)
	for i := 0; i < numberOfElements; i++ {
		queue.Push(expectedOrder[perm[i]])
		assert.Equal(t, i+1, queue.Len())
	}

	// Pop elements in order.
	for i := 0; i < numberOfElements; i++ {
		preview := queue.Peek()
		assert.Equal(t, expectedOrder[i], preview)
		assert.Equal(t, numberOfElements-i, queue.Len())

		next := queue.Pop()
		assert.Equal(t, expectedOrder[i], next)
		assert.Equal(t, numberOfElements-i-1, queue.Len())
	}
}

// TODO test removal once implemented
