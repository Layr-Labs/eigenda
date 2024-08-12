package lightnode

import "container/heap"

// registrationHeap implements the heap.Interface, used to create a priority queue.
type registrationHeap struct {
	data []*Registration
}

// Len returns the number of elements in the priority queue.
func (h *registrationHeap) Len() int {
	return len(h.data)
}

// Less returns whether the element with index i should sort before the element with index j.
// This registrationHeap sorts based on the nextShuffleTime of the light nodes.
func (h *registrationHeap) Less(i int, j int) bool {
	ii := h.data[i]
	jj := h.data[j]
	return ii.nextShuffleTime.Before(jj.nextShuffleTime)
}

// Swap swaps the elements with indexes i and j.
func (h *registrationHeap) Swap(i int, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

// Push adds an element to the end of the priority queue.
func (h *registrationHeap) Push(x any) {
	h.data = append(h.data, x.(*Registration))
}

// Pop removes and returns the last element in the priority queue.
func (h *registrationHeap) Pop() any {
	n := len(h.data)
	x := h.data[n-1]
	h.data = h.data[:n-1]
	return x
}

// PriorityQueue is a priority queue that sorts light nodes based on their nextShuffleTime.
type PriorityQueue struct {
	heap *registrationHeap
}

// NewPriorityQueue creates a new priority queue.
func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{
		heap: &registrationHeap{
			data: make([]*Registration, 0),
		},
	}
}

// Len returns the number of elements in the priority queue.
func (queue *PriorityQueue) Len() int {
	return queue.heap.Len()
}

// Push adds a light node to the priority queue.
func (queue *PriorityQueue) Push(registration *Registration) {
	heap.Push(queue.heap, registration)
}

// Pop removes and returns the light node with the earliest nextShuffleTime.
func (queue *PriorityQueue) Pop() *Registration {
	if queue.Len() == 0 {
		return nil
	}
	return heap.Pop(queue.heap).(*Registration)
}

// Peek returns the light node with the earliest nextShuffleTime without removing it from the queue. Returns
// nil if the queue is empty.
func (queue *PriorityQueue) Peek() *Registration {
	if queue.Len() == 0 {
		return nil
	}
	return queue.heap.data[0]
}

// Remove removes the light node with the given ID from the priority queue.
func (queue *PriorityQueue) Remove(lightNodeId uint64) {
	//TODO implement me
	panic("implement")
}
