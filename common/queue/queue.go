package queue

// Queue is an interface for a generic queue. It's absurd there isn't an equivalent in the standard golang libraries.
type Queue[T any] interface {
	// Push adds a value to the queue.
	Push(value T)

	// Pop removes and returns the value at the front of the queue.
	// If the queue is empty, the second return value will be false.
	Pop() (T, bool)

	// Peek returns the value at the front of the queue without removing it.
	// If the queue is empty, the second return value will be false.
	Peek() (T, bool)

	// Size returns the number of values in the queue.
	Size() int
}
