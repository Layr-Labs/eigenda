package traffic

import (
	"fmt"
	"math/rand"
	"sync"
)

// BlobTable tracks blobs written by the traffic generator. This is a thread safe data structure.
type BlobTable struct {

	// blobs contains all blobs currently tracked by the requiredReads.
	blobs []*BlobMetadata

	// size describes the total number of blobs currently tracked by the requiredReads.
	// size may be smaller than the capacity of the blobs slice.
	size uint32

	// lock is used to synchronize access to the requiredReads.
	lock sync.Mutex
}

// NewBlobTable creates a new BlobTable instance.
func NewBlobTable() BlobTable {
	return BlobTable{
		blobs: make([]*BlobMetadata, 1024),
		size:  0,
	}
}

// Size returns the total number of blobs currently tracked by the requiredReads.
func (table *BlobTable) Size() uint32 {
	table.lock.Lock()
	defer table.lock.Unlock()

	return table.size
}

// Add a blob to the requiredReads.
func (table *BlobTable) Add(blob *BlobMetadata) {
	table.lock.Lock()
	defer table.lock.Unlock()

	// TODO this calculation is probably a little wrong
	if table.size == uint32(len(table.blobs)) {
		panic(fmt.Sprintf("blob requiredReads is full, cannot add blob %x", blob.Key))
	}

	blob.index = table.size
	table.blobs[table.size] = blob
	table.size++
}

// AddOrReplace adds a blob to the requiredReads if there is capacity or replaces an existing blob at random
// if the requiredReads is full. This method is a no-op if maximumCapacity is 0.
func (table *BlobTable) AddOrReplace(blob *BlobMetadata, maximumCapacity uint32) {
	if maximumCapacity == 0 {
		return
	}

	table.lock.Lock()
	defer table.lock.Unlock()

	if table.size >= maximumCapacity {
		// replace random existing blob
		index := rand.Int31n(int32(table.size))
		table.blobs[index] = blob
		blob.index = uint32(index)
	} else {
		// add new blob
		blob.index = table.size
		table.blobs[table.size] = blob
		table.size++
	}
}

// GetRandom returns a random blob currently tracked by the requiredReads. Returns nil if the requiredReads is empty.
// Optionally decrements the read  permits of the blob if decrement is true. If the number of read permits
// reaches 0, the blob is removed  from the requiredReads. Returns the blob metadata (if there is at least one blob
// in the table) and a boolean indicating whether the blob was removed from the table as a result of this operation.
func (table *BlobTable) GetRandom(decrement bool) (*BlobMetadata, bool) {
	table.lock.Lock()
	defer table.lock.Unlock()

	if table.size == 0 {
		return nil, false
	}

	blob := table.blobs[rand.Int31n(int32(table.size))] // TODO make sure we can get items if we overflow an int32

	removed := false
	if decrement && blob.remainingReadPermits != -1 {
		blob.remainingReadPermits--
		if blob.remainingReadPermits == 0 {
			table.remove(blob)
			removed = true
		}
	}

	return blob, removed
}

// remove a blob from the requiredReads.
func (table *BlobTable) remove(blob *BlobMetadata) {
	if table.blobs[blob.index] != blob {
		panic(fmt.Sprintf("blob %x is not not present in the requiredReads at index %d", blob.Key, blob.index))
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