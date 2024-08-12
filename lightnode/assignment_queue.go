package lightnode

import "container/heap"

// assignmentHeap implements the heap.Interface for chunkGroupAssignment objects, used to create a priority queue.
type assignmentHeap struct {
	data []*chunkGroupAssignment
}

// Len returns the number of elements in the priority queue.
func (h *assignmentHeap) Len() int {
	return len(h.data)
}

// Less returns whether the element with index i should sort before the element with index j.
// This assignmentHeap sorts based on the nextShuffleTime of the light nodes.
func (h *assignmentHeap) Less(i int, j int) bool {
	ii := h.data[i]
	jj := h.data[j]
	return ii.nextShuffleTime.Before(jj.nextShuffleTime)
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

// assignmentQueue is a priority queue that sorts light nodes based on their nextShuffleTime.
type assignmentQueue struct {
	heap *assignmentHeap
}

// newAssignmentQueue creates a new priority queue.
func newAssignmentQueue() *assignmentQueue {
	return &assignmentQueue{
		heap: &assignmentHeap{
			data: make([]*chunkGroupAssignment, 0),
		},
	}
}

// Len returns the number of elements in the priority queue.
func (queue *assignmentQueue) Len() int {
	return queue.heap.Len()
}

// Push adds an assignment to the priority queue.
func (queue *assignmentQueue) Push(assignment *chunkGroupAssignment) {
	heap.Push(queue.heap, assignment)
}

// Pop removes and returns the assignment with the earliest nextShuffleTime.
func (queue *assignmentQueue) Pop() *chunkGroupAssignment {
	if queue.Len() == 0 {
		return nil
	}
	return heap.Pop(queue.heap).(*chunkGroupAssignment)
}

// Peek returns the assignment with the earliest nextShuffleTime without removing it from the queue. Returns
// nil if the queue is empty.
func (queue *assignmentQueue) Peek() *chunkGroupAssignment {
	if queue.Len() == 0 {
		return nil
	}
	return queue.heap.data[0]
}

// Remove removes the light node with the given ID from the priority queue.
func (queue *assignmentQueue) Remove(lightNodeId uint64) {
	//TODO implement me
	panic("implement")
}
