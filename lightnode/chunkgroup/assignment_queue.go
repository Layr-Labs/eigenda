package chunkgroup

import (
	"container/heap"
	"fmt"
)

// assignmentHeap implements the heap.Interface for chunkGroupAssignment objects, used to create a priority queue.
type assignmentHeap struct {
	data []*chunkGroupAssignment
}

// Len returns the number of elements in the priority queue.
func (h *assignmentHeap) Len() int {
	return len(h.data)
}

// Less returns whether the element with index i should sort before the element with index j.
// This assignmentHeap sorts based on the endOfEpoch of the light nodes.
func (h *assignmentHeap) Less(i int, j int) bool {
	ii := h.data[i]
	jj := h.data[j]
	return ii.endOfEpoch.Before(jj.endOfEpoch)
}

// Swap swaps the elements with indexes i and j.
func (h *assignmentHeap) Swap(i int, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

// Push adds an element to the end of the priority queue.
func (h *assignmentHeap) Push(x any) {
	h.data = append(h.data, x.(*chunkGroupAssignment))
}

// Pop removes and returns the last element in the priority queue.
func (h *assignmentHeap) Pop() any {
	n := len(h.data)
	x := h.data[n-1]
	h.data = h.data[:n-1]
	return x
}

// assignmentQueue is a priority queue that sorts light nodes based on their endOfEpoch.
type assignmentQueue struct {
	// The heap that stores the light nodes. Nodes are sorted by their endOfEpoch.
	heap *assignmentHeap

	// A set assignments in the queue. This is used to do efficient removals.
	// A true value indicates that the assignment is in the queue. A false value indicates
	// that the node was removed from the queue but has not yet been fully deleted.
	assignmentSet map[assignmentKey]bool

	// The number of elements in the queue. Tracked separately since the heap and NodeIdSet
	// may contain removed nodes that have not yet been fully garbage collected.
	size uint64
}

// newAssignmentQueue creates a new priority queue.
func newAssignmentQueue() *assignmentQueue {
	return &assignmentQueue{
		heap: &assignmentHeap{
			data: make([]*chunkGroupAssignment, 0),
		},
		assignmentSet: make(map[assignmentKey]bool),
	}
}

// Size returns the number of elements in the priority queue.
func (queue *assignmentQueue) Size() uint64 {
	return queue.size
}

// Push adds an chunkGroupAssignment to the priority queue.
// This is a no-op if the chunkGroupAssignment is already in the queue.
func (queue *assignmentQueue) Push(assignment *chunkGroupAssignment) {
	notRemoved, present := queue.assignmentSet[assignment.key]
	if present && notRemoved {
		return
	}

	queue.size++

	if !present {
		heap.Push(queue.heap, assignment)
	}

	queue.assignmentSet[assignment.key] = true
}

// Pop removes and returns the chunkGroupAssignment with the earliest endOfEpoch.
func (queue *assignmentQueue) Pop() *chunkGroupAssignment {
	queue.collectGarbage()
	if queue.size == 0 {
		return nil
	}
	assignment := heap.Pop(queue.heap).(*chunkGroupAssignment)
	delete(queue.assignmentSet, assignment.key)
	queue.size--
	return assignment
}

// Peek returns the chunkGroupAssignment with the earliest endOfEpoch without removing it from the queue. Returns
// nil if the queue is empty.
func (queue *assignmentQueue) Peek() *chunkGroupAssignment {
	queue.collectGarbage()
	if queue.size == 0 {
		return nil
	}
	return queue.heap.data[0]
}

// Remove removes the assignment with the given key from the queue.
// This is a no-op if the assignment is not in the queue.
func (queue *assignmentQueue) Remove(key assignmentKey) {
	// Deletion is lazy. The assignment is fully removed when it reaches the top of the heap.

	notRemoved, present := queue.assignmentSet[key]
	if !present || !notRemoved {
		// Element is either not in the queue or has already been marked for removal.
		return
	}

	queue.size--

	queue.assignmentSet[key] = false
	queue.collectGarbage()
}

// collectGarbage removes all nodes that have been removed from the queue but have not yet been fully deleted.
// This is done by popping elements from the heap until the first element is not marked for deletion.
func (queue *assignmentQueue) collectGarbage() {
	if len(queue.heap.data) == 0 {
		return
	}

	// sanity check to prevent infinite loops
	maxIterations := len(queue.heap.data) + 1

	for {
		maxIterations--
		if maxIterations < 0 {
			panic("garbage collection did not terminate")
		}

		if len(queue.heap.data) == 0 {
			return
		}

		next := queue.heap.data[0]

		notRemoved, present := queue.assignmentSet[next.key]
		if !present {
			panic(fmt.Sprintf("node %d is not in the assignmentSet", next.registration.ID()))
		}

		if notRemoved {
			// Once we find the first element that is not marked for deletion, we can stop.
			return
		}

		heap.Pop(queue.heap)
		delete(queue.assignmentSet, next.key)
	}
}
