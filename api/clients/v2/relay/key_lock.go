package relay

import (
	"sync"
)

// KeyLock is a utility that provides a way to lock access to a given key of type T
//
// This utility is useful in situations where you want to synchronize operations for something that doesn't exist
// in a concrete form. For example, perhaps you only want to create connections with a given peer on a single
// thread of execution, but the new peer could appear simultaneously in concurrent operations. This utility allows
// the first thread which encounters the new peer to perform necessary initialization tasks, and store generated
// artifacts in a central location for subsequent callers to access.
type KeyLock[T comparable] struct {
	// Map from key T to a mutex that corresponds to that key
	keyMutexMap map[T]*sync.Mutex
	// Used to lock access to the keyMutexMap, so that only a single mutex is created for each key
	globalMutex sync.Mutex
}

// NewKeyLock constructs a KeyLock utility
func NewKeyLock[T comparable]() *KeyLock[T] {
	return &KeyLock[T]{
		keyMutexMap: make(map[T]*sync.Mutex),
	}
}

// AcquireKeyLock acquires an exclusive lock on a conceptual key, and returns a function to release the lock
//
// The caller MUST eventually invoke the returned unlock function, or all future calls with the same key will block
// indefinitely
func (kl *KeyLock[T]) AcquireKeyLock(key T) func() {
	// we must globally synchronize access to the mutex map, so that only a single mutex will be created for a given key
	kl.globalMutex.Lock()
	keyMutex, valueAlreadyExists := kl.keyMutexMap[key]
	if !valueAlreadyExists {
		keyMutex = &sync.Mutex{}
		kl.keyMutexMap[key] = keyMutex
	}
	kl.globalMutex.Unlock()

	keyMutex.Lock()
	return keyMutex.Unlock
}
