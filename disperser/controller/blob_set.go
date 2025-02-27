package controller

import (
	"sync"

	v2 "github.com/Layr-Labs/eigenda/core/v2"
)

// BlobSet is a set of blob keys. This can be used to track a set of blobs.
type BlobSet interface {
	AddBlob(blobKey v2.BlobKey)
	RemoveBlob(blobKey v2.BlobKey)
	Size() int
	Contains(blobKey v2.BlobKey) bool
}

type blobSet struct {
	mu    sync.RWMutex
	blobs map[v2.BlobKey]struct{}
}

func NewBlobSet() BlobSet {
	return &blobSet{
		blobs: make(map[v2.BlobKey]struct{}),
	}
}

func (q *blobSet) AddBlob(blobKey v2.BlobKey) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.blobs[blobKey] = struct{}{}
}

func (q *blobSet) RemoveBlob(blobKey v2.BlobKey) {
	q.mu.Lock()
	defer q.mu.Unlock()

	delete(q.blobs, blobKey)
}

func (q *blobSet) Size() int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return len(q.blobs)
}

func (q *blobSet) Contains(blobKey v2.BlobKey) bool {
	q.mu.RLock()
	defer q.mu.RUnlock()

	_, ok := q.blobs[blobKey]
	return ok
}
