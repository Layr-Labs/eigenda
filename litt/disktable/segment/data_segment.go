package segment

import (
	"fmt"
	"path"
)

const (
	// KeysFileExtension is the file extension for the keys file. This file contains the keys for the data segment,
	// and is used for performing garbage collection on the key index.
	KeysFileExtension = ".keys"
	// ValuesFileExtension is the file extension for the values file. This file contains the values for the data
	// segment.
	ValuesFileExtension = ".values"
)

// dataSegment is a chunk of data stored on disk. All data in a particular data segment is expired at the same time.
//
// This struct is not thread safe, access control must be handled by the caller.
type dataSegment struct {
	// The index of the data segment. The first data segment ever created has index 0, the next has index 1, and so on.
	index uint32

	// The directory containing the data segment.
	parentDirectory string

	// The current size of data written to this data segment, in bytes.
	currentSize uint64

	// If true, this file is sealed and no more data can be written to it. If false, then data can still be written to
	// this file.
	sealed bool
}

// KeyFileName returns the name of the keys file for the data segment.
func (d *dataSegment) KeyFileName() string {
	return fmt.Sprintf("%d%s", d.index, KeysFileExtension)
}

// KeyFilePath returns the path to the keys file for the data segment.
func (d *dataSegment) KeyFilePath() string {
	return path.Join(d.parentDirectory, d.KeyFileName())
}

// ValueFileName returns the name of the values file for the data segment.
func (d *dataSegment) ValueFileName() string {
	return fmt.Sprintf("%d%s", d.index, ValuesFileExtension)
}

// ValueFilePath returns the path to the values file for the data segment.
func (d *dataSegment) ValueFilePath() string {
	return path.Join(d.parentDirectory, d.ValueFileName())
}

// Index returns the index of the data segment.
func (d *dataSegment) Index() uint32 {
	return d.index
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

// Seal flushes all data to disk and finalizes the metadata. After this method is called, no more data can be written
// to the data segment.
func (d *dataSegment) Seal() error {
	//TODO implement me
	panic("implement me")
}
