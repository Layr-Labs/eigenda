package table

import "sync"

// BlobStore is a thread safe data structure that tracks blobs written by the traffic generator.
type BlobStore struct {

	// blobs contains all blobs currently tracked by the store.
	blobs map[uint64]*BlobMetadata

	// nextKey describes the next key to used for the blobs map.
	nextKey uint64

	lock sync.Mutex
}

// NewBlobStore creates a new BlobStore instance.
func NewBlobStore() *BlobStore {
	return &BlobStore{
		blobs:   make(map[uint64]*BlobMetadata),
		nextKey: 0,
	}
}

// Add a blob to the store.
func (store *BlobStore) Add(blob *BlobMetadata) {
	store.lock.Lock()
	defer store.lock.Unlock()

	store.blobs[store.nextKey] = blob
	store.nextKey++
}

// GetNext returns the next blob in the store. Decrements the blob's read permits, removing it
// from the store if the permits are exhausted. This method makes no guarantees about the order
// in which blobs are returned. Returns nil if no blobs are available.
func (store *BlobStore) GetNext() *BlobMetadata {
	store.lock.Lock()
	defer store.lock.Unlock()

	for key, blob := range store.blobs {
		// Always return the first blob found.

		if blob.RemainingReadPermits > 0 {
			blob.RemainingReadPermits--
			if blob.RemainingReadPermits == 0 {
				delete(store.blobs, key)
			}
		}

		return blob
	}
	return nil
}

// Size returns the number of blobs currently stored.
func (store *BlobStore) Size() uint {
	store.lock.Lock()
	defer store.lock.Unlock()

	return uint(len(store.blobs))
}
