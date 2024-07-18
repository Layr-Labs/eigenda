package traffic

import (
	"fmt"
	"math/rand"
	"sync"
)

// BlobTable tracks blobs written by the traffic generator. This is a thread safe data structure.
type BlobTable struct {

	// blobs contains all blobs currently tracked by the table.
	blobs []*BlobMetadata

	// size describes the total number of blobs currently tracked by the table.
	// size may be smaller than the capacity of the blobs slice.
	size uint32

	// lock is used to synchronize access to the table.
	lock sync.Mutex
}

// NewBlobTable creates a new BlobTable instance.
func NewBlobTable() BlobTable {
	return BlobTable{
		blobs: make([]*BlobMetadata, 1024),
		size:  0,
	}
}

// Size returns the total number of blobs currently tracked by the table.
func (table *BlobTable) Size() uint32 {
	table.lock.Lock()
	defer table.lock.Unlock()

	return table.size
}

// Add a blob to the table.
func (table *BlobTable) Add(blob *BlobMetadata) {
	table.lock.Lock()
	defer table.lock.Unlock()

	if table.size == uint32(len(table.blobs)) {
		panic(fmt.Sprintf("blob table is full, cannot add blob %x", blob.Key))
	}

	blob.index = table.size
	table.blobs[table.size] = blob
	table.size++
}

// GetRandom returns a random blob currently tracked by the table. Returns nil if the table is empty.
// Optionally decrements the read  permits of the blob if decrement is true. If the number of read permits
// reaches 0, the blob is removed  from the table.
func (table *BlobTable) GetRandom(decrement bool) *BlobMetadata {
	table.lock.Lock()
	defer table.lock.Unlock()

	if table.size == 0 {
		return nil
	}

	blob := table.blobs[rand.Int31n(int32(table.size))] // TODO make sure we can get items if we overflow an int32

	if decrement && blob.remainingReadPermits != -1 {
		blob.remainingReadPermits--
		if blob.remainingReadPermits == 0 {
			table.remove(blob)
		}
	}

	return blob
}

// remove a blob from the table.
func (table *BlobTable) remove(blob *BlobMetadata) {
	if table.blobs[blob.index] != blob {
		panic(fmt.Sprintf("blob %x is not not present in the table at index %d", blob.Key, blob.index))
	}

	if table.size == 1 {
		table.blobs[0] = nil
	} else {
		// Move the last blob to the position of the blob being removed.
		table.blobs[blob.index] = table.blobs[table.size-1]
		table.blobs[blob.index].index = blob.index
		table.blobs[table.size-1] = nil
	}
	table.size--
}
