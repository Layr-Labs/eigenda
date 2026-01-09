package common

import (
	"container/heap"
	"fmt"
	"iter"
)

// A standard priority queue implementation using golang's container/heap package under the hood.
//
// By design, this implementation does not attempt to reclaim memory if the heap is large and then shrinks.
// As a general rule of thumb, if there are X items in the queue at one moment in time, it's likely that there will be
// on the order of X items in the queue at other times as well.
//
// This implementation is not thread safe.
type PriorityQueue[T any] struct {
	// Implementation of the heap interface.
	heap *heapImpl[T]
}

// Create a new priority queue that orders elements of type T according to the provided lessThan function.
func NewPriorityQueue[T any](
	// A function that returns true if a is less than b (i.e., it should show up earlier in the priority queue).
	lessThan func(a T, b T) bool,
) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		heap: &heapImpl[T]{
			items:      make([]T, 0),
			lessThan:   lessThan,
			rightIndex: -1,
		},
	}
}

// Size returns the number of items in the priority queue.
func (pq *PriorityQueue[T]) Size() int {
	return pq.heap.Len()
}

// Push adds an item to the priority queue.
func (pq *PriorityQueue[T]) Push(item T) {
	heap.Push(pq.heap, item)
}

// Pop removes and returns the highest-priority item from the priority queue.
//
// This method will panic if the priority queue is empty.
func (pq *PriorityQueue[T]) Pop() T {
	return heap.Pop(pq.heap).(T)
}

// TryPop attempts to remove and return the highest-priority item from the priority queue. If that is not possible
// (because the queue is empty), it returns false and a zero-value item.
func (pq *PriorityQueue[T]) TryPop() (value T, ok bool) {
	if pq.Size() == 0 {
		var zero T
		return zero, false
	}
	return pq.Pop(), true
}

// Peek returns the highest-priority item from the priority queue without removing it.
//
// This method will panic if the priority queue is empty.
func (pq *PriorityQueue[T]) Peek() T {
	return pq.heap.items[0]
}

// TryPeek attempts to return the highest-priority item from the priority queue without removing it. If that is not
// possible (because the queue is empty), it returns false and a zero-value item.
func (pq *PriorityQueue[T]) TryPeek() (value T, ok bool) {
	if pq.Size() == 0 {
		var zero T
		return zero, false
	}
	return pq.Peek(), true
}

// Build an iterator that pops all items from the priority queue in order.
func (pq *PriorityQueue[T]) PopIterator() iter.Seq[T] {
	return func(yield func(T) bool) {
		for pq.Size() > 0 {
			next := pq.Pop()
			if !yield(next) {
				return
			}
		}
	}
}

var _ heap.Interface = (*heapImpl[any])(nil)

// Implements the heap.Interface for PriorityQueue. This is a non-exported type, since we don't want to expose the
// ugly heap methods to users of PriorityQueue.
type heapImpl[T any] struct {
	// The items in the priority queue. May be longer than the number of items currently in the heap.
	// Intentionally does not shrink the slice when items are popped for efficiency.
	items []T

	// The index of the last valid item in the items slice. Will be -1 if the heap is empty.
	rightIndex int

	// Function to compare two items of type T. Should return true if a has higher priority than b.
	lessThan func(a T, b T) bool
}

func (h *heapImpl[T]) Len() int {
	return h.rightIndex + 1
}

func (h *heapImpl[T]) Less(i int, j int) bool {
	if i < 0 || i > h.rightIndex || j < 0 || j > h.rightIndex {
		panic(fmt.Sprintf("index out of range: i=%d, j=%d, rightIndex=%d", i, j, h.rightIndex))
	}
	return h.lessThan(h.items[i], h.items[j])
}

func (h *heapImpl[T]) Pop() any {
	if h.rightIndex < 0 {
		panic("pop from empty priority queue")
	}

	value := h.items[h.rightIndex]

	var zero T
	h.items[h.rightIndex] = zero

	h.rightIndex--
	return value
}

func (h *heapImpl[T]) Push(x any) {
	if len(h.items) > h.rightIndex+1 {
		h.items[h.rightIndex+1] = x.(T)
	} else {
		h.items = append(h.items, x.(T))
	}
	h.rightIndex++
}

func (h *heapImpl[T]) Swap(i int, j int) {
	if i < 0 || i > h.rightIndex || j < 0 || j > h.rightIndex {
		panic(fmt.Sprintf("index out of range: i=%d, j=%d, rightIndex=%d", i, j, h.rightIndex))
	}

	h.items[i], h.items[j] = h.items[j], h.items[i]
}
