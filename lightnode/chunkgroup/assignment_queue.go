package chunkgroup

import (
	"container/heap"
	"fmt"
)

// assignmentHeap implements the heap.Interface for assignment objects, used to create a priority queue.
type assignmentHeap struct {
	data []*assignment
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
	h.data = append(h.data, x.(*assignment))
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

	// A set of node IDs in the queue. This is used to do efficient removals.
	// A true value indicates that the node is in the queue. A false value indicates
	// that the node was removed from the queue but has not yet been fully deleted.
	nodeIdSet map[uint64]bool

	// The number of elements in the queue. Tracked separately since the heap and NodeIdSet
	// may contain removed nodes that have not yet been fully garbage collected.
	size uint
}

// newAssignmentQueue creates a new priority queue.
func newAssignmentQueue() *assignmentQueue {
	return &assignmentQueue{
		heap: &assignmentHeap{
			data: make([]*assignment, 0),
		},
		nodeIdSet: make(map[uint64]bool),
	}
}

// Size returns the number of elements in the priority queue.
func (queue *assignmentQueue) Size() uint {
	return queue.size
}

// Push adds an assignment to the priority queue. This is a no-op if the assignment is already in the queue.
func (queue *assignmentQueue) Push(assignment *assignment) {
	notRemoved, present := queue.nodeIdSet[assignment.registration.ID()]
	if present && notRemoved {
		return
	}

	queue.size++

	if !present {
		heap.Push(queue.heap, assignment)
	}

	queue.nodeIdSet[assignment.registration.ID()] = true
}

// Pop removes and returns the assignment with the earliest endOfEpoch.
func (queue *assignmentQueue) Pop() *assignment {
	queue.collectGarbage()
	if queue.size == 0 {
		return nil
	}
	assignment := heap.Pop(queue.heap).(*assignment)
	delete(queue.nodeIdSet, assignment.registration.ID())
	queue.size--
	return assignment
}

// Peek returns the assignment with the earliest endOfEpoch without removing it from the queue. Returns
// nil if the queue is empty.
func (queue *assignmentQueue) Peek() *assignment {
	queue.collectGarbage()
	if queue.size == 0 {
		return nil
	}
	return queue.heap.data[0]
}

// Remove removes the light node with the given ID from the priority queue.
// This is a no-op if the light node is not in the queue.
func (queue *assignmentQueue) Remove(lightNodeId uint64) {
	// Deletion is lazy. The node is fully removed when it reaches the top of the heap.

	notRemoved, present := queue.nodeIdSet[lightNodeId]
	if !present || !notRemoved {
		// Element is either not in the queue or has already been marked for removal.
		return
	}

	queue.size--

	queue.nodeIdSet[lightNodeId] = false
}

// collectGarbage removes all nodes that have been removed from the queue but have not yet been fully deleted.
// This is done by popping elements from the heap until the first element is not marked for deletion.
func (queue *assignmentQueue) collectGarbage() {
	if len(queue.heap.data) == 0 {
		return
	}

	// sanity check to prevent infinite loops
	maxIterations := len(queue.heap.data)

	for {
		maxIterations--
		if maxIterations < 0 {
			panic("garbage collection did not terminate")
		}

		next := queue.heap.data[0]

		notRemoved, present := queue.nodeIdSet[next.registration.ID()]
		if !present {
			panic(fmt.Sprintf("node %d is not in the nodeIdSet", next.registration.ID()))
		}

		if notRemoved {
			// Once we find the first element that is not marked for deletion, we can stop.
			return
		}

		heap.Pop(queue.heap)
	}
}
