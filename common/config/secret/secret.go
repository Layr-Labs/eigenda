package secret

import (
	"fmt"
	"sync"
)

var _ fmt.Stringer = &Secret[any]{}
var _ fmt.GoStringer = &Secret[any]{}

// Secret holds a value that should be kept secret. It is intentionally designed in a way that makes it very hard
// to accidentally expose the secret value, even if you print structs that contain it.
type Secret[T any] struct {
	lock sync.Mutex
	// The secret lives in this channel, which cannot be introspected or automatically printed using reflection.
	// Doesn't protect against deep magic (e.g. direct inspection of memory), but any golang library that uses
	// reflection to print struct fields won't be able to see inside this.
	vault chan T
}

// Create a new secret.
func NewSecret[T any](value T) *Secret[T] {

	s := &Secret[T]{
		vault: make(chan T, 1),
	}
	s.vault <- value
	return s
}

// Get returns the secret value.
func (s *Secret[T]) Get() T {
	s.lock.Lock()
	defer s.lock.Unlock()
	value := <-s.vault
	s.vault <- value
	return value
}

// Set updates the secret value, returning the old value.
func (s *Secret[T]) Set(value T) T {
	s.lock.Lock()
	defer s.lock.Unlock()
	oldValue := <-s.vault
	s.vault <- value
	return oldValue
}

func (s *Secret[T]) String() string {
	return "****"
}

func (s *Secret[T]) GoString() string {
	return "****"
}
