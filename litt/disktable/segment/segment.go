package segment

import (
	"errors"
	"fmt"
	"math"
	"path"
	"sync/atomic"
	"time"

	"github.com/Layr-Labs/eigenda/litt/types"
	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// unflushedKeysInitialCapacity is the initial capacity of the unflushedKeys slice. This slice is used to store keys
// that have been written to the segment but have not yet been flushed to disk.
const unflushedKeysInitialCapacity = 128

// shardControlChannelCapacity is the capacity of the channel used to send messages to the shard control loop.
const shardControlChannelCapacity = 32

// Segment is a chunk of data stored on disk. All data in a particular data segment is expired at the same time.
//
// This struct is not safe for operations that mutate the segment, access control must be handled by the caller.
type Segment struct {
	// The logger for the segment.
	logger logging.Logger

	// Used to signal an unrecoverable error in the segment. If fatalErrorHandler.Panic() is called, the entire DB
	// enters a "panic" state and will refuse to do additional work.
	fatalErrorHandler *util.FatalErrorHandler

	// The index of the data segment. The first data segment ever created has index 0, the next has index 1, and so on.
	index uint32

	// This file contains metadata about the segment.
	metadata *metadataFile

	// This file contains the keys for the data segment, and is used for performing garbage collection on the key index.
	keys *keyFile

	// The value files, one for each shard in the segment. Indexed by shard number.
	shards []*valueFile

	// shardSizes is a list of the current sizes of each shard in the segment. Indexed by shard number. This
	// value is only tracked for mutable segments (i.e. the unsealed segment), meaning that if this segment was loaded
	// from disk, the values in this slice will be zero.
	shardSizes []uint64

	// The current size of the key file in bytes. This is only tracked for mutable segments, meaning that if this
	// segment was loaded from disk, this value will be zero.
	keyFileSize uint64

	// The maximum size of all shards in this segment.
	maxShardSize uint64

	// The number of keys written to this segment.
	keyCount uint64

	// shardChannels is a list of channels used to send messages to the goroutine responsible for writing to
	// each shard. Indexed by shard number.
	shardChannels []chan any

	// keyFileChannel is a channel used to send messages to the goroutine responsible for writing to the key file.
	keyFileChannel chan any

	// deletionChannel permits a caller to block until this segment is fully deleted. An element is inserted into
	// the channel when the segment is fully deleted.
	deletionChannel chan struct{}

	// reservationCount is the number of reservations on this segment. The segment will not be deleted until this count
	// reaches zero.
	reservationCount atomic.Int32

	// nextSegment is the next segment in the chain (i.e. the segment with index+1). Each segment takes a reservation
	// on the next segment in the sequence. This reservation is released when the segment is fully deleted. This
	// ensures that segments are always deleted strictly in sequence. This makes it impossible for a crash to cause
	// segment X to be missing while segment X-1 is present.
	nextSegment *Segment

	// Used as a sanity checker. For each value written to the segment, the segment must eventually return
	// a key to be written to the keymap. This value tracks the number of values that have been written to the
	// segment but have not yet been flushed to the keymap. When the segment is eventually sealed, the code
	// asserts that this value is zero. This check should never fail, but is a nice safety net.
	unflushedKeyCount atomic.Int64
}

// CreateSegment creates a new data segment.
//
// Note that shardingFactor and salt parameters are ignored if this is not a new segment. Segments loaded from
// disk always use their original sharding factor and salt values
func CreateSegment(
	logger logging.Logger,
	fatalErrorHandler *util.FatalErrorHandler,
	index uint32,
	parentDirectories []string,
	now time.Time,
	shardingFactor uint32,
	salt uint32,
	fsync bool) (*Segment, error) {

	if len(parentDirectories) == 0 {
		return nil, errors.New("no parent directories provided")
	}

	// look for the metadata file
	metadataPath, err := lookForFile(parentDirectories, fmt.Sprintf("%d%s", index, MetadataFileExtension))
	var metadataDir string
	if err != nil {
		return nil, fmt.Errorf("failed to find metadata file: %v", err)
	}
	if metadataPath != "" {
		metadataDir = path.Dir(metadataPath)
	} else {
		// By default, put the metadata file in the first parent directory.
		metadataDir = parentDirectories[0]
	}
	metadata, err := createMetadataFile(index, shardingFactor, salt, metadataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to open metadata file: %v", err)
	}

	keysPath, err := lookForFile(parentDirectories, fmt.Sprintf("%d%s", index, KeyFileExtension))
	var keysDirectory string
	if err != nil {
		return nil, fmt.Errorf("failed to find key file: %v", err)
	}
	if keysPath != "" {
		keysDirectory = path.Dir(keysPath)
	} else {
		// By default, put the key file in the first parent directory.
		keysDirectory = parentDirectories[0]
	}
	keys, err := createKeyFile(logger, index, keysDirectory)
	if err != nil {
		return nil, fmt.Errorf("failed to open key file: %v", err)
	}

	keyFileSize := keys.Size()

	shards := make([]*valueFile, metadata.shardingFactor)
	for shard := uint32(0); shard < metadata.shardingFactor; shard++ {
		// Assign value files to parent directories in a round-robin fashion.
		// Assign the first shard to the directory at index 1. The first directory
		// is used by the keymap, so if we have enough directories we don't want to
		// use it for value files too.
		parentDirectory := parentDirectories[int(shard+1)%len(parentDirectories)]

		values, err := createValueFile(logger, index, shard, parentDirectory, fsync)
		if err != nil {
			return nil, fmt.Errorf("failed to open value file: %v", err)
		}
		shards[shard] = values
	}

	shardSizes := make([]uint64, metadata.shardingFactor)

	shardChannels := make([]chan any, metadata.shardingFactor)
	for shard := uint32(0); shard < metadata.shardingFactor; shard++ {
		shardChannels[shard] = make(chan any, shardControlChannelCapacity)
	}

	// If at all possible, we want to size this channel so that the goroutines writing data to the sharded value files
	// do not block on insertion into this channel. Scale the size of this channel by the number of shards, as more
	// shards mean there may be a higher rate of writes to this channel.
	keyFileChannel := make(chan any, shardControlChannelCapacity*metadata.shardingFactor)

	segment := &Segment{
		logger:            logger,
		fatalErrorHandler: fatalErrorHandler,
		index:             index,
		metadata:          metadata,
		keys:              keys,
		shards:            shards,
		shardSizes:        shardSizes,
		keyFileSize:       keyFileSize,
		shardChannels:     shardChannels,
		keyFileChannel:    keyFileChannel,
		deletionChannel:   make(chan struct{}, 1),
	}

	// Segments are returned with an initial reference count of 1, as the caller of the constructor is considered to
	// have a reference to the segment.
	segment.reservationCount.Store(1)

	// Start up the control loops.
	for shard := uint32(0); shard < metadata.shardingFactor; shard++ {
		go segment.shardControlLoop(shard)
	}

	go segment.keyFileControlLoop()

	return segment, nil
}

// LoadSegment loads an existing segment from disk. If that segment is unsealed, this method will seal it.
func LoadSegment(logger logging.Logger,
	fatalErrorHandler *util.FatalErrorHandler,
	index uint32,
	parentDirectories []string,
	now time.Time) (*Segment, error) {

	if len(parentDirectories) == 0 {
		return nil, errors.New("no parent directories provided")
	}

	// Look for the metadata file.
	metadata, err := loadMetadataFile(index, parentDirectories)
	if err != nil {
		return nil, fmt.Errorf("failed to open metadata file: %w", err)
	}

	// seal the segment if it is unsealed
	if !metadata.sealed {
		err = metadata.seal(now)
		if err != nil {
			return nil, fmt.Errorf("failed to seal segment: %w", err)
		}
	}

	// Look for the key file.
	keys, err := loadKeyFile(logger, index, parentDirectories)
	if err != nil {
		return nil, fmt.Errorf("failed to open key file: %v", err)
	}
	keyFileSize := keys.Size()

	// Look for the value files. There should be one for each shard.
	shards := make([]*valueFile, metadata.shardingFactor)
	for shard := uint32(0); shard < metadata.shardingFactor; shard++ {
		values, err := loadValueFile(logger, index, shard, parentDirectories)
		if err != nil {
			return nil, fmt.Errorf("failed to open value file: %v", err)
		}
		shards[shard] = values
	}

	segment := &Segment{
		logger:            logger,
		fatalErrorHandler: fatalErrorHandler,
		index:             index,
		metadata:          metadata,
		keys:              keys,
		shards:            shards,
		keyFileSize:       keyFileSize,
		deletionChannel:   make(chan struct{}, 1),
	}

	// Segments are returned with an initial reference count of 1, as the caller of the constructor is considered to
	// have a reference to the segment.
	segment.reservationCount.Store(1)

	return segment, nil

}

// Size returns the size of the segment in bytes. Counts bytes that are on disk or that will eventually end up on disk.
// This method is not thread safe, and should not be called concurrently with methods that modify the segment.
func (s *Segment) Size() uint64 {
	size := s.metadata.Size()

	if s.IsSealed() {
		// This segment is immutable, so it's thread safe to query the files directly.
		size += s.keys.Size()
		for _, shard := range s.shards {
			size += shard.Size()
		}
	} else {
		// This segment is mutable. We must use our local reckoning of the sizes of the files.
		size += s.keyFileSize
		for _, shardSize := range s.shardSizes {
			size += shardSize
		}
	}

	return size
}

// lookForFile looks for a file in a list of directories. It returns an error if the file appears
// in more than one directory, and an empty string if the file is not found. If the file is found and
// there are no errors, this method returns the path to the file.
func lookForFile(directories []string, fileName string) (string, error) {
	locations := make([]string, 0, 1)
	for _, directory := range directories {
		potentialLocation := path.Join(directory, fileName)
		exists, err := util.Exists(potentialLocation)
		if err != nil {
			return "", fmt.Errorf("failed to check if file %s exists: %v", potentialLocation, err)
		}
		if exists {
			locations = append(locations, potentialLocation)
		}
	}

	if len(locations) > 1 {
		return "", fmt.Errorf("file %s found in multiple directories: %v", fileName, locations)
	}

	if len(locations) == 0 {
		return "", nil
	}
	return locations[0], nil
}

// SetNextSegment sets the next segment in the chain.
func (s *Segment) SetNextSegment(nextSegment *Segment) {
	nextSegment.Reserve()
	s.nextSegment = nextSegment
}

// GetShard returns the shard number for a key.
func (s *Segment) GetShard(key []byte) uint32 {
	if s.metadata.shardingFactor == 1 {
		// Shortcut: if we have one shard, we don't need to hash the key to figure out the mapping.
		return 0
	}

	return util.HashKey(key, s.metadata.salt) % s.metadata.shardingFactor
}

// Write records a key-value pair in the data segment, returning the maximum size of all shards within this segment.
//
// This method does not ensure that the key-value pair is actually written to disk, only that it will eventually be
// written to disk. Flush must be called to ensure that all data previously passed to Write is written to disk.
func (s *Segment) Write(data *types.KVPair) (keyCount uint64, keyFileSize uint64, err error) {
	if s.metadata.sealed {
		return 0, 0, fmt.Errorf("segment is sealed, cannot write data")
	}

	shard := s.GetShard(data.Key)
	currentSize := s.shardSizes[shard]

	if currentSize > math.MaxUint32 {
		// No matter the configuration, we absolutely cannot permit a value to be written if the first byte of the
		// value would be beyond position 2^32. This is because we only have 32 bits in an address to store the
		// position of a value's first byte.
		return 0, 0,
			fmt.Errorf("value file already contains %d bytes, cannot add a new value", currentSize)
	}
	s.unflushedKeyCount.Add(1)
	firstByteIndex := uint32(currentSize)

	s.shardSizes[shard] += uint64(len(data.Value)) + 4 /* uint32 length */
	if s.shardSizes[shard] > s.maxShardSize {
		s.maxShardSize = s.shardSizes[shard]
	}
	s.keyCount++
	s.keyFileSize += uint64(len(data.Key)) + 4 /* uint32 length */ + 8 /* uint64 Address */

	// Forward the value to the shard control loop, which asynchronously writes it to the value file.
	shardRequest := &valueToWrite{
		value:                  data.Value,
		expectedFirstByteIndex: firstByteIndex,
	}
	err = util.SendIfNotFatal(s.fatalErrorHandler, s.shardChannels[shard], shardRequest)
	if err != nil {
		return 0, 0,
			fmt.Errorf("failed to send value to shard control loop: %v", err)
	}

	// Forward the value to the key and its address file control loop, which asynchronously writes it to the key file.
	keyRequest := &types.KAPair{
		Key:     data.Key,
		Address: types.NewAddress(s.index, firstByteIndex),
	}
	err = util.SendIfNotFatal(s.fatalErrorHandler, s.keyFileChannel, keyRequest)
	if err != nil {
		return 0, 0,
			fmt.Errorf("failed to send key to key file control loop: %v", err)
	}

	return s.keyCount, s.keyFileSize, nil
}

// GetMaxShardSize returns the maximum size of all shards in this segment.
func (s *Segment) GetMaxShardSize() uint64 {
	return s.maxShardSize
}

// Read fetches the data for a key from the data segment.
//
// It is only thread safe to read from a segment if the key being read has previously been flushed to disk.
func (s *Segment) Read(key []byte, dataAddress types.Address) ([]byte, error) {
	shard := s.GetShard(key)
	values := s.shards[shard]

	value, err := values.read(dataAddress.Offset())
	if err != nil {
		return nil, fmt.Errorf("failed to read value: %v", err)
	}
	return value, nil
}

// GetKeys returns all keys in the data segment. Only permitted to be called after the segment has been sealed.
func (s *Segment) GetKeys() ([]*types.KAPair, error) {
	if !s.metadata.sealed {
		return nil, fmt.Errorf("segment is not sealed, cannot read keys")
	}

	keys, err := s.keys.readKeys()
	if err != nil {
		return nil, fmt.Errorf("failed to read keys: %v", err)
	}
	return keys, nil
}

// FlushWaitFunction is a function that waits for a flush operation to complete. It returns the addresses of the data
// that was flushed, or an error if the flush operation failed.
type FlushWaitFunction func() ([]*types.KAPair, error)

// Flush schedules a flush operation. Flush operations are performed serially in the order they are scheduled.
// This method returns a function that, when called, will block until the flush operation is complete. The function
// returns the addresses of the data that was flushed, or an error if the flush operation failed.
func (s *Segment) Flush() (FlushWaitFunction, error) {
	return s.flush(false)
}

func (s *Segment) flush(seal bool) (FlushWaitFunction, error) {
	// Schedule a flush for all shards.
	shardResponseChannels := make([]chan struct{}, s.metadata.shardingFactor)
	for shard, shardChannel := range s.shardChannels {
		shardResponseChannels[shard] = make(chan struct{}, 1)
		request := &shardFlushRequest{
			seal:              seal,
			completionChannel: shardResponseChannels[shard],
		}
		err := util.SendIfNotFatal(s.fatalErrorHandler, shardChannel, request)
		if err != nil {
			return nil, fmt.Errorf("failed to send flush request to shard %d: %v", shard, err)
		}
	}

	// Schedule a flush for the key channel.
	// Now that all shards have sent their key/address pairs to the key file, flush the key file.
	keyResponseChannel := make(chan *keyFileFlushResponse, 1)
	request := &keyFileFlushRequest{
		seal:              seal,
		completionChannel: keyResponseChannel,
	}
	err := util.SendIfNotFatal(s.fatalErrorHandler, s.keyFileChannel, request)
	if err != nil {
		return nil, fmt.Errorf("failed to send flush request to key file: %v", err)
	}

	return func() ([]*types.KAPair, error) {
		// Wait for each shard to finish flushing.
		for i := range s.shardChannels {
			_, err := util.AwaitIfNotFatal(s.fatalErrorHandler, shardResponseChannels[i])
			if err != nil {
				return nil, fmt.Errorf("failed to flush shard %d: %v", i, err)
			}
		}

		keyFlushResponse, err := util.AwaitIfNotFatal(s.fatalErrorHandler, keyResponseChannel)
		if err != nil {
			return nil, fmt.Errorf("failed to flush key file: %v", err)
		}

		s.unflushedKeyCount.Add(-int64(len(keyFlushResponse.addresses)))
		return keyFlushResponse.addresses, nil
	}, nil
}

// Seal flushes all data to disk and finalizes the metadata. Returns addresses that became durable as a result of
// this method call. After this method is called, no more data can be written to this segment.
func (s *Segment) Seal(now time.Time) ([]*types.KAPair, error) {
	flushWaitFunction, err := s.flush(true)
	if err != nil {
		return nil, fmt.Errorf("failed to flush segment: %v", err)
	}
	addresses, err := flushWaitFunction()
	if err != nil {
		return nil, fmt.Errorf("failed to flush segment: %v", err)
	}

	// Seal the metadata file.
	err = s.metadata.seal(now)
	if err != nil {
		return nil, fmt.Errorf("failed to seal metadata file: %v", err)
	}

	unflushedKeyCount := s.unflushedKeyCount.Load()
	if s.unflushedKeyCount.Load() != 0 {
		return nil, fmt.Errorf("segment %d has %d unflushedKeyCount keys", s.index, unflushedKeyCount)
	}

	return addresses, nil
}

// IsSealed returns true if the segment is sealed, and false otherwise.
func (s *Segment) IsSealed() bool {
	return s.metadata.sealed
}

// GetSealTime returns the time at which the segment was sealed. If the file is not sealed, this method will return
// the zero time.
func (s *Segment) GetSealTime() time.Time {
	return time.Unix(0, int64(s.metadata.lastValueTimestamp))
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

// Release releases a reservation held on this segment. A segment cannot be deleted until all reservations on it
// have been released. The last call to Release() that releases the final reservation schedules the segment for
// asynchronous deletion.
func (s *Segment) Release() {
	reservations := s.reservationCount.Add(-1)

	if reservations > 0 {
		return
	}

	if reservations < 0 {
		// This should be impossible.
		s.fatalErrorHandler.Panic(
			fmt.Errorf("segment %d has negative reservation count: %d", s.index, reservations))
	}

	go func() {
		err := s.delete()
		if err != nil {
			s.fatalErrorHandler.Panic(fmt.Errorf("failed to delete segment: %v", err))
		}
	}()
}

// BlockUntilFullyDeleted blocks until the segment is fully deleted. If the segment is not yet fully released,
// this method will block until it is. This method should only be called once per segment (the second call
// will block forever!).
func (s *Segment) BlockUntilFullyDeleted() error {
	_, err := util.AwaitIfNotFatal(s.fatalErrorHandler, s.deletionChannel)
	if err != nil {
		return fmt.Errorf("failed to await segment deletion: %v", err)
	}
	return nil
}

// delete deletes the segment from disk.
func (s *Segment) delete() error {
	defer func() {
		s.deletionChannel <- struct{}{}
	}()

	err := s.keys.delete()
	if err != nil {
		return fmt.Errorf("failed to delete key file, segment %d: %v", s.index, err)
	}
	for shardIndex, shard := range s.shards {
		err = shard.delete()
		if err != nil {
			return fmt.Errorf("failed to delete value file, segment %d, shard %d: %v", s.index, shardIndex, err)
		}
	}
	err = s.metadata.delete()
	if err != nil {
		return fmt.Errorf("failed to delete metadata file, segment %d: %v", s.index, err)
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

// handleShardFlushRequest handles a request to flush a shard to disk.
func (s *Segment) handleShardFlushRequest(shard uint32, request *shardFlushRequest) {
	if request.seal {
		err := s.shards[shard].seal()
		if err != nil {
			s.fatalErrorHandler.Panic(fmt.Errorf("failed to seal value file: %v", err))
		}
	} else {
		err := s.shards[shard].flush()
		if err != nil {
			s.fatalErrorHandler.Panic(fmt.Errorf("failed to flush value file: %v", err))
		}
	}
	request.completionChannel <- struct{}{}
}

// handleShardWrite applies a single write operation to a shard.
func (s *Segment) handleShardWrite(shard uint32, data *valueToWrite) {
	firstByteIndex, err := s.shards[shard].write(data.value)
	if err != nil {
		s.fatalErrorHandler.Panic(fmt.Errorf("failed to write value to value file: %v", err))
	}

	if firstByteIndex != data.expectedFirstByteIndex {
		// This should never happen. But it's a good sanity check.
		if firstByteIndex != data.expectedFirstByteIndex {
			s.fatalErrorHandler.Panic(
				fmt.Errorf("expected first byte index %d, got %d", data.expectedFirstByteIndex, firstByteIndex))
		}
	}
}

// handleKeyFileWrite writes a key to the key file.
func (s *Segment) handleKeyFileWrite(data *types.KAPair) {
	err := s.keys.write(data.Key, data.Address)
	if err != nil {
		s.fatalErrorHandler.Panic(fmt.Errorf("failed to write key to key file: %v", err))
	}
}

// handleKeyFileFlushRequest handles a request to flush the key file to disk.
func (s *Segment) handleKeyFileFlushRequest(request *keyFileFlushRequest, unflushedKeys []*types.KAPair) {
	if request.seal {
		err := s.keys.seal()
		if err != nil {
			s.fatalErrorHandler.Panic(fmt.Errorf("failed to seal key file: %v", err))
		}
	} else {
		err := s.keys.flush()
		if err != nil {
			s.fatalErrorHandler.Panic(fmt.Errorf("failed to flush key file: %v", err))
		}
	}

	request.completionChannel <- &keyFileFlushResponse{
		addresses: unflushedKeys,
	}
}

// shardFlushRequest is a message sent to shard control loops to request that they flush their data to disk.
type shardFlushRequest struct {
	// If true, seal the shard after flushing. If false, do not seal the shard.
	seal bool

	// As each shard finishes its flush it will send an object to this channel.
	completionChannel chan struct{}
}

// valueToWrite is a message sent to the shard control loop to request that it write a value to the value file.
type valueToWrite struct {
	value                  []byte
	expectedFirstByteIndex uint32
}

// shardControlLoop is the main loop for performing modifications to a particular shard. Each shard is managed
// by its own goroutine, which is running this function.
func (s *Segment) shardControlLoop(shard uint32) {
	for {
		select {
		case <-s.fatalErrorHandler.ImmediateShutdownRequired():
			s.logger.Infof("segment %d shard %d control loop exiting, context cancelled", s.index, shard)
			return
		case operation := <-s.shardChannels[shard]:
			if flushRequest, ok := operation.(*shardFlushRequest); ok {
				s.handleShardFlushRequest(shard, flushRequest)
				if flushRequest.seal {
					// After sealing, we can exit the control loop.
					return
				}
			} else if data, ok := operation.(*valueToWrite); ok {
				s.handleShardWrite(shard, data)
				continue
			} else {
				s.fatalErrorHandler.Panic(
					fmt.Errorf("unknown operation type in shard control loop: %T", operation))
			}
		}
	}
}

// keyFileFlushRequest is a message sent to the key file control loop to request that it flush its data to disk.
type keyFileFlushRequest struct {
	// If true, seal the key file after flushing. If false, do not seal the key file.
	seal bool

	// As the key file finishes its flush, it will either send an error if something went wrong, or nil if the flush was
	// successful.
	completionChannel chan *keyFileFlushResponse
}

// keyFileFlushResponse is a message sent from the key file control loop to the caller of Flush to indicate that the
// key file has been flushed.
type keyFileFlushResponse struct {
	addresses []*types.KAPair
}

// keyFileControlLoop is the main loop for performing modifications to the key file. This goroutine is responsible
// for writing key-address pairs to the key file.
func (s *Segment) keyFileControlLoop() {
	unflushedKeys := make([]*types.KAPair, 0, unflushedKeysInitialCapacity)

	for {
		select {
		case <-s.fatalErrorHandler.ImmediateShutdownRequired():
			s.logger.Infof("segment %d key file control loop exiting, context cancelled", s.index)
			return
		case operation := <-s.keyFileChannel:

			if flushRequest, ok := operation.(*keyFileFlushRequest); ok {
				s.handleKeyFileFlushRequest(flushRequest, unflushedKeys)
				unflushedKeys = make([]*types.KAPair, 0, unflushedKeysInitialCapacity)

				if flushRequest.seal {
					// After sealing, we can exit the control loop.
					return
				}

			} else if data, ok := operation.(*types.KAPair); ok {
				s.handleKeyFileWrite(data)
				unflushedKeys = append(unflushedKeys, data)

			} else {
				s.fatalErrorHandler.Panic(
					fmt.Errorf("unknown operation type in key file control loop: %T", operation))
			}
		}
	}
}
