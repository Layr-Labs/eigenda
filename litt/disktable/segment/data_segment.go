package segment

import (
	"fmt"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"time"
)

// Segment is a chunk of data stored on disk. All data in a particular data segment is expired at the same time.
//
// This struct is not thread safe, access control must be handled by the caller.
type Segment struct {
	// The index of the data segment. The first data segment ever created has index 0, the next has index 1, and so on.
	index uint32

	// This file contains metadata about the segment.
	metadata *metadataFile

	// This file contains the keys for the data segment, and is used for performing garbage collection on the key index.
	keys *keyFile

	// This file contains the values for the data segment.
	values *valueFile
}

// NewSegment creates a new data segment.
func NewSegment(
	logger logging.Logger,
	index uint32,
	parentDirectory string) (*Segment, error) {

	metadata, err := newMetadataFile(index, parentDirectory)
	if err != nil {
		return nil, fmt.Errorf("failed to open metadata file: %v", err)
	}

	keys, err := newKeyFile(logger, index, parentDirectory, metadata.sealed)
	if err != nil {
		return nil, fmt.Errorf("failed to open key file: %v", err)
	}

	values, err := newValueFile(logger, index, parentDirectory, metadata.sealed)
	if err != nil {
		return nil, fmt.Errorf("failed to open value file: %v", err)
	}

	return &Segment{
		index:    index,
		metadata: metadata,
		keys:     keys,
		values:   values,
	}, nil
}

// Index returns the index of the data segment.
func (s *Segment) Index() uint32 {
	return s.index
}

// Put records a key-value pair in the data segment, returning the resulting address of the data.
// This method does not ensure that the key-value pair is actually written to disk, only that it is recorded
// in the data segment. Flush must be called to ensure that all data previously passed to Put is written to disk.
func (s *Segment) Put(key []byte, value []byte) (Address, error) {
	err := s.keys.write(key)
	if err != nil {
		return 0, fmt.Errorf("failed to write key: %v", err)
	}

	address, err := s.values.write(value)
	if err != nil {
		return 0, fmt.Errorf("failed to write value: %v", err)
	}

	return address, nil
}

// Get fetches the data for a key from the data segment.
func (s *Segment) Get(dataAddress Address) ([]byte, error) {
	value, err := s.values.read(dataAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to read value: %v", err)
	}
	return value, nil
}

// GetKeys returns all keys in the data segment. Only permitted to be called after the segment has been sealed.
func (s *Segment) GetKeys() ([][]byte, error) {
	keys, err := s.keys.readKeys()
	if err != nil {
		return nil, fmt.Errorf("failed to read keys: %v", err)
	}
	return keys, nil
}

// Flush writes the data segment to disk.
func (s *Segment) Flush() error {
	err := s.keys.flush()
	if err != nil {
		return fmt.Errorf("failed to flush key file: %v", err)
	}

	err = s.values.flush()
	if err != nil {
		return fmt.Errorf("failed to flush value file: %v", err)
	}

	return nil
}

// Delete deletes the data segment from disk.
func (s *Segment) Delete() error {
	err := s.keys.delete()
	if err != nil {
		return fmt.Errorf("failed to delete key file: %v", err)
	}
	err = s.values.delete()
	if err != nil {
		return fmt.Errorf("failed to delete value file: %v", err)
	}
	err = s.metadata.delete()
	if err != nil {
		return fmt.Errorf("failed to delete metadata file: %v", err)
	}
	return nil
}

// Seal flushes all data to disk and finalizes the metadata. After this method is called, no more data can be written
// to the data segment.
func (s *Segment) Seal(now time.Time) error {
	err := s.keys.seal()
	if err != nil {
		return fmt.Errorf("failed to seal key file: %v", err)
	}

	err = s.values.seal()
	if err != nil {
		return fmt.Errorf("failed to seal value file: %v", err)
	}

	err = s.metadata.seal(now)
	if err != nil {
		return fmt.Errorf("failed to seal metadata file: %v", err)
	}

	return nil
}
