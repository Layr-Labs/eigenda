package common

import "sync"

// IndexLock is similar to a sync.Mutex, but it allows for different indices to be locked independently. There
// is a probability that any two indices' locks interfere with each other, but this can be made arbitrarily small
// by configuration.
//
// Internally, an IndexLock keeps an array of mutexes. Each index is mapped onto one of these mutexes in the array,
// such that the same index always maps to the same mutex. Note that due to this mapping, otherwise unrelated indices
// may end up using the same mutex. Increasing the number of locks will decrease the probability of unrelated indices
// contending for the same lock, but will also increase memory usage.
type IndexLock struct {
	locks []sync.Mutex
}

// NewIndexLock creates a new IndexLock.
func NewIndexLock(numLocks uint32) *IndexLock {
	locks := make([]sync.Mutex, numLocks)
	return &IndexLock{locks: locks}
}

// Lock locks the given index. Two calls to Lock with the same index will attempt to acquire the same lock.
// Two calls to Lock with different indices may or may not acquire the same lock. After calling lock,
// the caller must eventually also call Unlock.
func (i *IndexLock) Lock(index uint64) {
	lockIndex := index % uint64(len(i.locks))
	i.locks[lockIndex].Lock()
}

// Unlock unlocks the given index. It is an error to call Unlock with an index that has not been locked.
func (i *IndexLock) Unlock(index uint64) {
	lockIndex := index % uint64(len(i.locks))
	i.locks[lockIndex].Unlock()
}
