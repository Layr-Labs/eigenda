package litt

// dataSegment is a chunk of data stored on disk. All data in a particular data segment is expired at the same time.
type dataSegment struct {
	// The index of the data segment. The first data segment ever created has index 0, the next has index 1, and so on.
	index uint32
}

// HasCapacityFor returns true if the data segment has capacity for the given key-value pair.
// Data files are not permitted to grow beyond 2^32 bytes in size.
func (d *dataSegment) HasCapacityFor(key []byte, value []byte) bool {
	//TODO implement me
	panic("implement me")
}

// Put records a key-value pair in the data segment, returning the resulting address of the data.
// This method does not ensure that the key-value pair is actually written to disk, only that it is recorded
// in the data segment. Flush must be called to ensure that all data previously passed to Put is written to disk.
func (d *dataSegment) Put(key []byte, value []byte) (address, error) {
	return 0, nil
}

// Get fetches the data for a key from the data segment.
func (d *dataSegment) Get(dataAddress address) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

// Flush writes the data segment to disk.
func (d *dataSegment) Flush() error {
	//TODO implement me
	panic("implement me")
}

// Delete deletes the data segment from disk.
func (d *dataSegment) Delete() error {
	//TODO implement me
	panic("implement me")
}
