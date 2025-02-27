package segment

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Layr-Labs/eigenda/litt/types"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// Segment is a chunk of data stored on disk. All data in a particular data segment is expired at the same time.
//
// This struct is not thread safe, access control must be handled by the caller.
type Segment struct {
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

	// lock controls access to the segment.
	lock sync.RWMutex

	// deletionMutex permits a caller to block until this segment is fully deleted. The channel has a capacity of 1, and
	// there is an element in the channel up until the segment is fully deleted.
	deletionChannel chan struct{}

	// reservationCount is the number of reservations on this segment. The segment will not be deleted until this count
	// reaches zero.
	reservationCount atomic.Int32

	// nextSegment is the next segment in the chain (i.e. the segment with index+1). Each segment takes a reservation
	// on the next segment in the sequence. This reservation is released when the segment is fully deleted. This
	// ensures that segments are always deleted strictly in sequence. This makes it impossible for a crash to cause
	// segment X to be missing while segment X-1 is present.
	nextSegment *Segment
}

// NewSegment creates a new data segment.
//
// Note that shardingFactor and salt parameters are ignored if this is not a new segment. Segments loaded from
// disk always use their original sharding factor and salt values
func NewSegment(
	logger logging.Logger,
	index uint32,
	parentDirectory string,
	now time.Time,
	shardingFactor uint32,
	salt uint32,
	sealIfUnsealed bool) (*Segment, error) {

	metadata, err := newMetadataFile(index, shardingFactor, salt, parentDirectory)
	if err != nil {
		return nil, fmt.Errorf("failed to open metadata file: %v", err)
	}

	if sealIfUnsealed && !metadata.sealed {
		err = metadata.seal(now)
		if err != nil {
			return nil, fmt.Errorf("failed to seal segment: %v", err)
		}
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
		logger:          logger,
		index:           index,
		metadata:        metadata,
		keys:            keys,
		values:          values,
		deletionChannel: make(chan struct{}, 1),
	}

	// This element is removed from the channel when the segment is fully deleted.
	segment.deletionChannel <- struct{}{}

	// Segments are returned with an initial reference count of 1, as the caller of the constructor is considered to
	// have a reference to the segment.
	segment.reservationCount.Store(1)

	return segment, nil
}

// SetNextSegment sets the next segment in the chain.
func (s *Segment) SetNextSegment(nextSegment *Segment) {
	nextSegment.Reserve()
	s.nextSegment = nextSegment
}

// Write records a key-value pair in the data segment, returning the resulting address of the data.
//
// This method does not ensure that the key-value pair is actually written to disk, only that it is recorded
// in the data segment. Flush must be called to ensure that all data previously passed to Put is written to disk.
func (s *Segment) Write(key []byte, value []byte) (address types.Address, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	address, err = s.values.write(value)
	if err != nil {
		return 0, fmt.Errorf("failed to write value: %v", err)
	}

	err = s.keys.write(key, address)
	if err != nil {
		return 0, fmt.Errorf("failed to write key: %v", err)
	}

	return address, nil
}

// CurrentSize returns the current size of the data segment.
func (s *Segment) CurrentSize() uint64 {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.values.currentSize
}

// Read fetches the data for a key from the data segment.
func (s *Segment) Read(dataAddress types.Address) ([]byte, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	value, err := s.values.read(dataAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to read value: %v", err)
	}
	return value, nil
}

// GetKeys returns all keys in the data segment. Only permitted to be called after the segment has been sealed.
func (s *Segment) GetKeys() ([]*types.KAPair, error) {
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

	if reservations < 0 {
		// This should be impossible.
		s.logger.Errorf("segment %d has negative reservation count: %d", s.index, reservations)
	}

	go func() {
		err := s.delete()
		if err != nil {
			s.logger.Errorf("failed to delete segment: %v", err)
		}
	}()
}

// BlockUntilFullyDeleted blocks until the segment is fully deleted. If the segment is not yet fully released,
// this method will block until it is.
func (s *Segment) BlockUntilFullyDeleted() {
	s.deletionChannel <- struct{}{}
	<-s.deletionChannel
}

// delete deletes the segment from disk.
func (s *Segment) delete() error {
	s.lock.Lock()
	defer func() {
		s.lock.Unlock()
		<-s.deletionChannel
	}()

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

	// The next segment is now eligible for deletion once it is fully released by other reservation holders.
	if s.nextSegment != nil {
		s.nextSegment.Release()
	}

	return nil
}

func (s *Segment) String() string {
	var sealedString string
	if s.metadata.sealed {
		sealedString = "sealed"
	} else {
		sealedString = "unsealed"
	}

	return fmt.Sprintf("[seg %d - %s]", s.index, sealedString)
}
