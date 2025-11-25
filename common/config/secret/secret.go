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
func (s *Secret) Get() string {
	s.lock.Lock()
	defer s.lock.Unlock()
	value := <-s.vault
	s.vault <- value
	return value
}

// Set updates the secret value, returning the old value.
func (s *Secret) Set(value string) string {
	s.lock.Lock()
	defer s.lock.Unlock()
	oldValue := <-s.vault
	s.vault <- value
	return oldValue
}

func (s *Secret) String() string {
	return "****"
}

func (s *Secret) GoString() string {
	return "****"
}
