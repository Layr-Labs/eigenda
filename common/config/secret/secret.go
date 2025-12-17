package secret

import (
	"fmt"
	"sync"
)

var _ fmt.Stringer = &Secret{}
var _ fmt.GoStringer = &Secret{}

// Secret holds a string that should be kept secret. It is intentionally designed in a way that makes it very hard
// to accidentally expose the secret value, even if you print structs that contain it or use reflection.
type Secret struct {
	lock sync.Mutex
	// The secret lives in this channel, which cannot be introspected or automatically printed using reflection.
	// Doesn't protect against deep magic (e.g. direct inspection of memory), but any golang library that uses
	// reflection to print struct fields won't be able to see inside this.
	vault chan string
}

// Create a new secret.
func NewSecret(value string) *Secret {
	s := &Secret{
		vault: make(chan string, 1),
	}
	s.vault <- value
	return s
}

// Get returns the secret value.
//
// Safe to call on a nil *Secret, in which case it returns an empty string.
func (s *Secret) Get() string {
	if s == nil {
		return ""
	}

	s.lock.Lock()
	defer s.lock.Unlock()
	value := <-s.vault
	s.vault <- value
	return value
}

// Set updates the secret value, returning the old value.
//
// Not safe to call on a nil *Secret (will panic).
func (s *Secret) Set(value string) string {
	if s == nil {
		panic("cannot set value on nil Secret")
	}

	s.lock.Lock()
	defer s.lock.Unlock()
	oldValue := <-s.vault
	s.vault <- value
	return oldValue
}

func (s *Secret) String() string {
	if s == nil {
		return ""
	}

	return "****"
}

func (s *Secret) GoString() string {
	if s == nil {
		return ""
	}

	return "****"
}
