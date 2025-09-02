package common

import (
	"fmt"
	"math"

	"github.com/Layr-Labs/eigenda/common/enforce"
)

// The minimum initial capacity of a RandomAccessDeque.
const minimumInitialCapacity = 32

// A double-ended queue (deque) that supports O(1) lookup by index.
//
// - Insertion time: O(1) average, O(n) worst-case (when resizing is needed)
// - Deletion time: O(1) average, array space is not reclaimed
// - Lookup time by index: O(1)
// - Iteration: O(1) to build iterator, O(1) per step
//
// This data structure is not thread safe.
type RandomAccessDeque[T any] struct {
	// The current number of elements in the deque.
	size uint64
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

	if initialCapacity < minimumInitialCapacity {
		initialCapacity = minimumInitialCapacity
	}

	return &RandomAccessDeque[T]{
		data:            make([]T, initialCapacity),
		startIndex:      0,
		endIndex:        0,
		initialCapacity: initialCapacity,
	}
}

// Get the number of elements in the deque.
//
// O(1)
func (s *RandomAccessDeque[T]) Size() uint64 {
	return s.size
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
		s.startIndex--
	}

	s.data[s.startIndex] = value
	s.size++
}

// Return the value at the front of the deque without removing it. If the deque is empty, returns an error.
//
// O(1)
func (s *RandomAccessDeque[T]) PeekFront() (value T, err error) {
	if s.size == 0 {
		var zero T
		return zero, fmt.Errorf("cannot peek front: deque is empty")
	}

	value, err = s.Get(0)
	enforce.NilError(err, "Get failed, this should never happen if size check passes")

	return value, nil
}

// Remove and return the value at the front of the deque. If the deque is empty, returns an error.
//
// O(1)
func (s *RandomAccessDeque[T]) PopFront() (value T, err error) {
	if s.size == 0 {
		var zero T
		return zero, fmt.Errorf("cannot pop front: deque is empty")
	}

	value = s.data[s.startIndex]

	var zero T
	s.data[s.startIndex] = zero

	if s.startIndex == uint64(len(s.data)-1) {
		// wrap around
		s.startIndex = 0
	} else {
		s.startIndex++
	}

	s.size--

	return value, nil
}

// Insert a value at the back of the deque. This value will have index Size()-1 after insertion.
//
// O(1) average, O(n) worst-case (when resizing is needed)
func (s *RandomAccessDeque[T]) PushBack(value T) {
	s.resizeForInsertion()

	s.data[s.endIndex] = value

	if s.endIndex == uint64(len(s.data)-1) {
		// wrap around
		s.endIndex = 0
	} else {
		s.endIndex++
	}

	s.size++
}

// Return the value at the back of the deque without removing it. If the deque is empty, returns an error.
//
// O(1)
func (s *RandomAccessDeque[T]) PeekBack() (value T, err error) {
	if s.size == 0 {
		var zero T
		return zero, fmt.Errorf("cannot peek back: deque is empty")
	}

	value, err = s.Get(s.size - 1)
	enforce.NilError(err, "Get failed, this should never happen if size check passes")

	return value, nil
}

// Remove and return the value at the back of the deque. If the deque is empty, returns an error.
//
// O(1)
func (s *RandomAccessDeque[T]) PopBack() (value T, err error) {
	if s.size == 0 {
		var zero T
		return zero, fmt.Errorf("cannot pop back: deque is empty")
	}

	var backIndex uint64
	if s.endIndex == 0 {
		backIndex = uint64(len(s.data)) - 1
	} else {
		backIndex = s.endIndex - 1
	}

	value = s.data[backIndex]

	var zero T
	s.data[backIndex] = zero

	s.endIndex = backIndex

	s.size--

	return value, nil
}

// Get the value at the specified index. If the index is out of bounds returns an error.
func (s *RandomAccessDeque[T]) Get(index uint64) (value T, err error) {
	if index >= s.size {
		var zero T
		return zero, fmt.Errorf("index %d out of bounds (size %d)", index, s.size)
	}

	realIndex := (s.startIndex + index) % uint64(len(s.data))
	return s.data[realIndex], nil
}

// Get an element indexed from the last thing in the deque. Equivalent to Get(Size() - 1 - index).
// If the index is out of bounds returns an error.
func (s *RandomAccessDeque[T]) GetFromBack(index uint64) (value T, err error) {
	if index >= s.size {
		var zero T
		return zero, fmt.Errorf("index %d out of bounds (size %d)", index, s.size)
	}

	value, err = s.Get(s.size - 1 - index)
	enforce.NilError(err, "Get failed, this should never happen if size check passes")

	return value, nil
}

// Set the value at the specified index, replacing the existing value, which is returned.
// If the index is out of bounds returns an error.
func (s *RandomAccessDeque[T]) Set(index uint64, value T) (previousValue T, err error) {
	if index >= s.size {
		var zero T
		return zero, fmt.Errorf("index %d out of bounds (size %d)", index, s.size)
	}

	realIndex := (s.startIndex + index) % uint64(len(s.data))
	previousValue = s.data[realIndex]
	s.data[realIndex] = value
	return previousValue, nil
}

// Set an element indexed from the last thing in the deque, replacing the existing value, which is returned.
// Equivalent to Set(Size() - 1 - index, value). // TODO test this
func (s *RandomAccessDeque[T]) SetFromBack(index uint64, value T) (previousValue T, err error) {
	if index >= s.size {
		var zero T
		return zero, fmt.Errorf("index %d out of bounds (size %d)", index, s.size)
	}

	previousValue, err = s.Set(s.size-1-index, value)
	enforce.NilError(err, "Set failed, this should never happen if size check passes")

	return previousValue, nil
}

// Clear all elements from the deque. Reclaims space in the underlying array.
//
// O(1)
func (s *RandomAccessDeque[T]) Clear() {
	s.startIndex = 0
	s.endIndex = 0
	s.size = 0
	// Reset the underlying array to allow garbage collection of contained elements.
	s.data = make([]T, s.initialCapacity)
}

// Get an iterator over the elements in the deque, from front to back. It is not safe to get an iterator,
// modify the deque, and then use the iterator again.
//
// O(1) to call this method, O(1) per iteration step.
func (s *RandomAccessDeque[T]) Iterator() func(yield func(uint64, T) bool) {
	if s.size == 0 {
		return func(yield func(uint64, T) bool) {
			// no-op
		}
	}

	iterator, err := s.IteratorFrom(0)
	enforce.NilError(err, "IteratorFrom failed, this should never happen")

	return iterator
}

// Get an iterator over the elements in the deque, from the specified index to back. It is not safe to get an iterator,
// modify the deque, and then use the iterator again.
//
// O(1) to call this method, O(1) per iteration step.
func (s *RandomAccessDeque[T]) IteratorFrom(index uint64) (func(yield func(uint64, T) bool), error) {

	if index >= s.size {
		return nil, fmt.Errorf("index %d out of bounds (size %d)", index, s.size)
	}

	return func(yield func(uint64, T) bool) {
		for i := index; i < s.size; i++ {
			value, err := s.Get(i)
			enforce.NilError(err, "Get failed, did you modify the deque while iterating?!?")

			yield(i, value)
		}
	}, nil
}

// Get an iterator over the elements in the deque, from back to front. It is not safe to get an iterator,
// modify the deque, and then use the iterator again.
//
// // O(1) to call this method, O(1) per iteration step.
func (s *RandomAccessDeque[T]) ReverseIterator() func(yield func(uint64, T) bool) {
	if s.size == 0 {
		return func(yield func(uint64, T) bool) {
			// no-op
		}
	}

	iterator, err := s.ReverseIteratorFrom(s.size - 1)
	enforce.NilError(err, "ReverseIteratorFrom failed, this should never happen")

	return iterator
}

// Get an iterator over the elements in the deque, from the specified index to front. It is not safe to get an iterator,
// modify the deque, and then use the iterator again.
//
// O(1) to call this method, O(1) per iteration step.
func (s *RandomAccessDeque[T]) ReverseIteratorFrom(index uint64) (func(yield func(uint64, T) bool), error) {

	if index >= s.size {
		return nil, fmt.Errorf("index %d out of bounds (size %d)", index, s.size)
	}

	return func(yield func(uint64, T) bool) {
		for i := index; i != math.MaxUint64; i-- {
			value, err := s.Get(i)
			enforce.NilError(err, "Get failed, did you modify the deque while iterating?!?")

			yield(i, value)
		}
	}, nil
}

// Resize the underlying array to accommodate at least one more insertion. Preserves existing elements.
// If no resizing is needed, this is a no-op.
func (s *RandomAccessDeque[T]) resizeForInsertion() {
	remainingCapacity := uint64(len(s.data)) - s.size

	if remainingCapacity > 0 {
		return
	}

	newData := make([]T, len(s.data)*2)

	for index, value := range s.Iterator() {
		newData[index] = value
	}

	s.data = newData
	s.startIndex = 0
	s.endIndex = s.size
}

// Perform a binary search in the deque for an element matching the compare function. Assumes that
// the deque is sorted according to the same compare function. If an exact match can't be found,
// returns the index of the location where the value would be inserted if it were inserted in the proper location.
//
// The compare function `compare(a V, b T) int` should return:
//   - negative value if a < b
//   - zero if a == b
//   - positive value if a > b
//
// If the deque is not sorted or if the ordering is not a total ordering, the return value is undefined. This function
// is not defined as a method on RandomAccessDeque due to this fact. Not all RandomAccessDeque instances will be sorted,
// and so this function is not always valid to call.
func BinarySearchInOrderedDeque[V any, T any](
	deque *RandomAccessDeque[T],
	value V,
	compare func(a V, b T) int) (index uint64, exact bool) {

	if deque.size == 0 {
		return 0, false
	}

	// Index is the external index in the deque, from 0 to size-1, not indices as they
	// appear in the underlying array.
	left := uint64(0)
	right := deque.size - 1
	var targetIndex uint64

	for left < right {
		targetIndex = left + (right-left)/2
		target, err := deque.Get(targetIndex)
		enforce.NilError(err, "Get failed, this should never happen with valid indices")

		cmp := compare(value, target)

		if cmp == 0 {
			// We've found an exact match.
			return targetIndex, true
		} else if cmp < 0 {
			// value < target, search left half
			//
			//      value is here
			//  |-----------------------|-----------------------|
			// left                   target                  right
			right = targetIndex - 1
		} else {
			// value > target, search right half
			//
			//                               value is here
			//  |-----------------------|-----------------------|
			// left                   target                  right
			left = targetIndex + 1
		}
	}

	element, err := deque.Get(left)
	enforce.NilError(err, "Get failed, this should never happen with valid indices")
	cmp := compare(value, element)
	if cmp == 0 {
		// We've found an exact match.
		return left, true
	} else if cmp < 0 {
		// value < element, so missing value should go to the left of it
		return left, false
	}
	// value > element, so missing value should go to the right of it
	return left + 1, false
}
