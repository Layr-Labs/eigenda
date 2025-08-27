package common

import "fmt"

// A double-ended queue (deque) that supports O(1) lookup by index.
//
// - Insertion time: O(1) average, O(n) worst-case (when resizing is needed)
// - Deletion time: O(1) average, array space is not reclaimed
// - Lookup time by index: O(1)
//
// This data structure is not thread safe.
type RandomAccessDeque[T any] struct {
	// Underlying data storage
	data []T
	// The index in data that corresponds to the logical start of the deque.
	startIndex uint64
	// The index in data that corresponds to the logical end of the deque (one past the last element).
	endIndex uint64
	// The initial capacity of the deque. Used when calling Clear().
	initialCapacity uint64
}

// Create a new RandomAccessDeque with the specified initial capacity. Queue can grow beyond this capacity if needed.
func NewRandomAccessDeque[T any](initialCapacity uint64) *RandomAccessDeque[T] {
	return &RandomAccessDeque[T]{
		data:            make([]T, initialCapacity),
		startIndex:      0,
		endIndex:        1,
		initialCapacity: initialCapacity,
	}
}

// Get the number of elements in the deque.
//
// O(1)
func (s *RandomAccessDeque[T]) Size() uint64 {
	if s.endIndex >= s.startIndex {
		return s.endIndex - s.startIndex
	}
	return uint64(len(s.data)) - s.startIndex + s.endIndex
}

// Insert a value at the front of the deque. This value will have index 0 after insertion, and all other values will
// have their indices increased by 1.
//
// O(1) average, O(n) worst-case (when resizing is needed)
func (s *RandomAccessDeque[T]) PushFront(value T) {
	s.resizeForInsertion()

	if s.startIndex == 0 {
		// wrap around
		s.startIndex = uint64(len(s.data)) - 1
	} else {
		s.startIndex -= 1
	}

	err := s.Set(0, value)
	enforce.NoError(err, "failed to push value")

	// TODO
}

// Return the value at the front of the deque without removing it. If the deque is empty, returns an error.
//
// O(1)
func (s *RandomAccessDeque[T]) PeekFront() (value T, err error) {
	size := s.Size()
	if size == 0 {
		var zero T
		return zero, fmt.Errorf("cannot peek front: deque is empty")
	}

	value, err = s.Get(0)

	if err != nil {
		var zero T
		return zero, fmt.Errorf("cannot peek front: %w", err)
	}

	return value, nil
}

// Remove and return the value at the front of the deque. If the deque is empty, returns an error.
//
// O(1)
func (s *RandomAccessDeque[T]) PopFront() (value T, err error) {
	size := s.Size()
	if size == 0 {
		var zero T
		return zero, fmt.Errorf("cannot pop front: deque is empty")
	}

	var zero T
	value, err = s.Set(0, zero)

	if err != nil {
		var zero T
		return zero, fmt.Errorf("cannot pop front: %w", err)
	}

	if s.startIndex == uint64(len(s.data)-1) {
		// wrap around
		s.startIndex = 0
	} else {
		s.startIndex += 1
	}

	return value, nil
}

// Insert a value at the back of the deque. This value will have index Size()-1 after insertion.
//
// O(1) average, O(n) worst-case (when resizing is needed)
func (s *RandomAccessDeque[T]) PushBack(value T) {
	// TODO
}

// Return the value at the back of the deque without removing it. If the deque is empty, returns an error.
//
// O(1)
func (s *RandomAccessDeque[T]) PeekBack() (value T, err error) {
	size := s.Size()
	if size == 0 {
		var zero T
		return zero, fmt.Errorf("cannot peek back: deque is empty")
	}

	value, err = s.Get(size - 1)

	if err != nil {
		var zero T
		return zero, fmt.Errorf("cannot peek back: %w", err)
	}

	return value, nil
}

// Remove and return the value at the back of the deque. If the deque is empty, returns an error.
//
// O(1)
func (s *RandomAccessDeque[T]) PopBack() (value T, err error) {
	size := s.Size()
	if size == 0 {
		var zero T
		return zero, fmt.Errorf("cannot pop back: deque is empty")
	}

	var zero T
	value, err = s.Set(size-1, zero)

	if err != nil {
		var zero T
		return zero, fmt.Errorf("cannot pop back: %w", err)
	}

	if s.endIndex == 0 {
		// wrap around
		s.endIndex = uint64(len(s.data)) - 1
	} else {
		s.endIndex -= 1
	}

	return value, nil
}

// Get the value at the specified index. If the index is out of bounds returns an error.
func (s *RandomAccessDeque[T]) Get(index uint64) (value T, err error) {
	size := s.Size()
	if index >= size {
		var zero T
		return zero, fmt.Errorf("index %d out of bounds (size %d)", index, size)
	}

	realIndex := (s.startIndex + index) % uint64(len(s.data))
	return s.data[realIndex], nil
}

// Set the value at the specified index, replacing the existing value, which is returned.
// If the index is out of bounds returns an error.
func (s *RandomAccessDeque[T]) Set(index uint64, value T) (previousValue T, err error) {
	size := s.Size()
	if index >= size {
		var zero T
		return zero, fmt.Errorf("index %d out of bounds (size %d)", index, size)
	}

	realIndex := (s.startIndex + index) % uint64(len(s.data))
	previousValue = s.data[realIndex]
	s.data[realIndex] = value
	return previousValue, nil
}

// Clear all elements from the deque. Reclaims space in the underlying array.
//
// O(1)
func (s *RandomAccessDeque[T]) Clear() {
	s.startIndex = 0
	s.endIndex = 0
	// Reset the underlying array to allow garbage collection of contained elements.
	s.data = make([]T, s.initialCapacity)
}

// Get an iterator over the elements in the deque, from front to back. It is not safe to get an iterator,
// modify the deque, and then use the iterator again.
//
// O(1) to call this method, O(1) per iteration step.
func (s *RandomAccessDeque[T]) Iterator() func(yield func(int, T) bool) {
	// TODO
	return nil
}

// Get an iterator over the elements in the deque, from back to front. It is not safe to get an iterator,
// modify the deque, and then use the iterator again.
//
// // O(1) to call this method, O(1) per iteration step.
func (s *RandomAccessDeque[T]) ReverseIterator() func(yield func(int, T) bool) {
	// TODO
	return nil
}

// Resize the underlying array to accommodate at least one more insertion. Preserves existing elements.
// If no resizing is needed, this is a no-op.
func (s *RandomAccessDeque[T]) resizeForInsertion() {
	size := s.Size()
	remainingCapacity := uint64(len(s.data)) - size

	if remainingCapacity > 0 {
		return
	}

	newData := make([]T, len(s.data)*2)

	for index, value := range s.Iterator() {
		newData[index] = value
	}

	s.data = newData
	s.startIndex = 0
	s.endIndex = size
}
