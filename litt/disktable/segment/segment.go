package segment

import (
	"fmt"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"sync"
	"sync/atomic"
	"time"
)

// Segment is a chunk of data stored on disk. All data in a particular data segment is expired at the same time.
//
// This struct is not thread safe, access control must be handled by the caller.
type Segment struct { // TODO do we need to export this type?
	// The logger for the segment.
	logger logging.Logger

	// The index of the data segment. The first data segment ever created has index 0, the next has index 1, and so on.
	index uint32

	// This file contains metadata about the segment.
	metadata *metadataFile

	// This file contains the keys for the data segment, and is used for performing garbage collection on the key index.
	keys *keyFile

	// This file contains the values for the data segment.
	values *valueFile

	// The target size for value files.
	targetFileSize uint32

	// lock controls access to the segment.
	lock sync.RWMutex

	// reservationCount is the number of reservations on this segment. The segment will not be deleted until this count
	// reaches zero.
	reservationCount atomic.Int32
}

// NewSegment creates a new data segment.
func NewSegment(
	logger logging.Logger,
	index uint32,
	parentDirectory string,
	targetFileSize uint32) (*Segment, error) {

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

	segment := &Segment{
		logger:         logger,
		index:          index,
		metadata:       metadata,
		keys:           keys,
		values:         values,
		targetFileSize: targetFileSize,
	}

	return segment, nil
}

// Write records a key-value pair in the data segment, returning the resulting address of the data.
// If this file is full, ok will be false and the data will not have been written.
//
// This method does not ensure that the key-value pair is actually written to disk, only that it is recorded
// in the data segment. Flush must be called to ensure that all data previously passed to Put is written to disk.
func (s *Segment) Write(key []byte, value []byte) (address Address, ok bool, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	resultingSize := s.values.currentSize + uint64(len(value)) + 4
	if resultingSize > uint64(s.targetFileSize) {
		// segment is full
		return 0, false, nil
	}

	err = s.keys.write(key)
	if err != nil {
		return 0, false, fmt.Errorf("failed to write key: %v", err)
	}

	address, err = s.values.write(value)
	if err != nil {
		return 0, false, fmt.Errorf("failed to write value: %v", err)
	}

	return address, true, nil
}

// CurrentSize returns the current size of the data segment.
func (s *Segment) CurrentSize() uint64 {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.values.currentSize
}

// Read fetches the data for a key from the data segment.
func (s *Segment) Read(dataAddress Address) ([]byte, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	value, err := s.values.read(dataAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to read value: %v", err)
	}
	return value, nil
}

// GetKeys returns all keys in the data segment. Only permitted to be called after the segment has been sealed.
func (s *Segment) GetKeys() ([][]byte, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	keys, err := s.keys.readKeys()
	if err != nil {
		return nil, fmt.Errorf("failed to read keys: %v", err)
	}
	return keys, nil
}

// Flush writes the data segment to disk.
func (s *Segment) Flush() error {
	s.lock.Lock()
	defer s.lock.Unlock()

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

// Seal flushes all data to disk and finalizes the metadata. After this method is called, no more data can be written
// to the data segment.
func (s *Segment) Seal(now time.Time) error {
	s.lock.Lock()
	defer s.lock.Unlock()

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

// IsSealed returns true if the segment is sealed, and false otherwise.
func (s *Segment) IsSealed() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.metadata.sealed
}

// GetSealTime returns the time at which the segment was sealed. If the file is not sealed, this method will return
// the zero time.
func (s *Segment) GetSealTime() time.Time {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return time.Unix(0, int64(s.metadata.timestamp))
}

// Reserve reserves the segment, preventing it from being deleted. Returns true if the reservation was successful, and
// false otherwise.
func (s *Segment) Reserve() bool {
	for {
		reservations := s.reservationCount.Load()
		if reservations <= 0 {
			return false
		}

		if s.reservationCount.CompareAndSwap(reservations, reservations+1) {
			return true
		}
	}
}

// Release releases the segment, allowing it to be deleted. Deletion happens inside this method if the reservation count
// reaches zero as a result of this call.
func (s *Segment) Release() {
	reservations := s.reservationCount.Add(-1)
	if reservations > 0 {
		return
	}

	go func() {
		err := s.delete()
		if err != nil {
			s.logger.Errorf("failed to delete segment: %v", err)
		}
	}()
}

// delete deletes the segment from disk.
func (s *Segment) delete() error {
	s.lock.Lock()
	defer s.lock.Unlock()

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
