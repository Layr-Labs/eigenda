package table

import "sync"

// BlobTable is a thread safe data structure that tracks blobs written by the traffic generator.
type BlobTable struct {

	// blobs contains all blobs currently tracked by the store.
	blobs map[uint64]*BlobMetadata

	// nextKey describes the next key to used for the blobs map.
	nextKey uint64

	lock sync.Mutex
}

// NewBlobStore creates a new BlobTable instance.
func NewBlobStore() *BlobTable {
	return &BlobTable{
		blobs:   make(map[uint64]*BlobMetadata),
		nextKey: 0,
	}
}

// Add a blob to the store.
func (table *BlobTable) Add(blob *BlobMetadata) {
	table.lock.Lock()
	defer table.lock.Unlock()

	table.blobs[table.nextKey] = blob
	table.nextKey++
}

// GetNext returns the next blob in the store. Decrements the blob's read permits, removing it
// from the store if the permits are exhausted. This method makes no guarantees about the order
// in which blobs are returned. Returns nil if no blobs are available.
func (table *BlobTable) GetNext() *BlobMetadata {
	table.lock.Lock()
	defer table.lock.Unlock()

	for key, blob := range table.blobs {
		// Always return the first blob found.

		if blob.RemainingReadPermits > 0 {
			blob.RemainingReadPermits--
			if blob.RemainingReadPermits == 0 {
				delete(table.blobs, key)
			}
		}

		return blob
	}
	return nil
}

// Size returns the number of blobs currently stored.
func (table *BlobTable) Size() uint {
	table.lock.Lock()
	defer table.lock.Unlock()

	return uint(len(table.blobs))
}

// GetAll returns all blobs currently stored. Intended for test purposes.
func (table *BlobTable) GetAll() []*BlobMetadata {
	table.lock.Lock()
	defer table.lock.Unlock()

	blobs := make([]*BlobMetadata, 0, len(table.blobs))
	for _, blob := range table.blobs {
		blobs = append(blobs, blob)
	}
	return blobs
}
