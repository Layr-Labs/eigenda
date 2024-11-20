package queue

var _ Queue[string] = &LinkedQueue[string]{}

// LinkedQueue is a queue that uses a linked list to store values. It is not thread safe.
type LinkedQueue[T any] struct {
	front *node[T]
	back  *node[T]
	size  int
}

// node is a single element in the linked list.
type node[T any] struct {
	value T
	next  *node[T]
}

func (l *LinkedQueue[T]) Push(value T) {
	if l.size == 0 {
		l.front = &node[T]{value: value}
		l.back = l.front
	} else {
		n := &node[T]{value: value}
		l.back.next = n
		l.back = n
	}
	l.size++
}

func (l *LinkedQueue[T]) Pop() (T, bool) {
	if l.size == 0 {
		var zero T
		return zero, false
	}

	value := l.front.value
	l.front = l.front.next
	l.size--
	return value, true
}

func (l *LinkedQueue[T]) Peek() (T, bool) {
	if l.size == 0 {
		var zero T
		return zero, false
	}
	return l.front.value, true
}

func (l *LinkedQueue[T]) Size() int {
	return l.size
}
