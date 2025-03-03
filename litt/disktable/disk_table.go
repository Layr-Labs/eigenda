package disktable

import (
	"context"
	"fmt"
	"os"
	"path"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Layr-Labs/eigenda/litt"
	"github.com/Layr-Labs/eigenda/litt/disktable/keymap"
	"github.com/Layr-Labs/eigenda/litt/disktable/segment"
	"github.com/Layr-Labs/eigenda/litt/types"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

var _ litt.ManagedTable = (*DiskTable)(nil)

// segmentDirectory is the directory where segment files are stored, relative to the root directory.
const segmentDirectory = "segments"

// DiskTable manages a table's Segments.
type DiskTable struct {
	// The context for the disk table.
	ctx context.Context

	// The logger for the disk table.
	logger logging.Logger

	// The root directories for the disk table.
	roots []string

	// The directories where segment files are stored.
	segmentDirectories []string

	// The table's name.
	name string

	// The table's metadata.
	metadata *tableMetadata

	// Random data to make the sharding hash function hard for an attacker to predict.
	salt uint32

	// A map of keys to their addresses.
	keyMap keymap.KeyMap

	// unflushedDataCache is a map of keys to their values that may not have been flushed to disk yet. This is used as a
	// lookup table when data is requested from the table before it has been flushed to disk.
	unflushedDataCache sync.Map

	// The index of the lowest numbered segment. After initial creation, only the garbage collection
	// thread is permitted to read/write this value  for the sake of thread safety.
	lowestSegmentIndex uint32

	// The index of the highest numbered segment. All writes are applied to this segment.
	highestSegmentIndex uint32

	// All segments currently in use.
	segments map[uint32]*segment.Segment

	// The target size for value files.
	targetFileSize uint32

	// segmentLock protects access to the segments map and highestSegmentIndex.
	// Does not protect the segments themselves.
	segmentLock sync.RWMutex

	// alive is an atomic boolean that is true if the disk table is alive, or false if it has been shut down.
	alive atomic.Bool

	// This channel can be used to block until the disk table has been stopped. The channel has a capacity of 1, and
	// there is an element in the channel up until the disk table has been stopped.
	stopChannel chan struct{}

	// controllerChan is the channel for messages sent to controller goroutine. No data managed by the DiskTable
	// may be mutated by anything other than the controller, with the exception of the
	// segmentLock and reference counting.
	controllerChan chan any

	// garbageCollectionPeriod is the period at which garbage collection is run.
	garbageCollectionPeriod time.Duration

	// timeSource is the time source used by the disk table.
	timeSource func() time.Time

	// Set to true when there is a fatal error on a goroutine that doesn't have a way to return the error.
	// This is used during testing.
	fatalError atomic.Bool
}

// NewDiskTable creates a new DiskTable.
func NewDiskTable(
	ctx context.Context,
	logger logging.Logger,
	timeSource func() time.Time,
	name string,
	keyMap keymap.KeyMap,
	roots []string,
	targetFileSize uint32,
	controlChannelSize int,
	shardingFactor uint32,
	salt uint32,
	gcPeriod time.Duration) (litt.ManagedTable, error) {

	if gcPeriod <= 0 {
		return nil, fmt.Errorf("garbage collection period must be greater than 0")
	}

	for _, root := range roots {
		_, err := os.Stat(root)
		if os.IsNotExist(err) {
			err := os.MkdirAll(root, 0755)
			if err != nil {
				return nil, fmt.Errorf("failed to create root directory: %v", err)
			}
		} else if err != nil {
			return nil, fmt.Errorf("failed to stat root directory: %v", err)
		}
	}

	segDirs := make([]string, 0, len(roots))
	for _, root := range roots {
		segDir := path.Join(root, segmentDirectory)
		segDirs = append(segDirs, segDir)
		_, err := os.Stat(segDir)
		if os.IsNotExist(err) {
			err := os.MkdirAll(segDir, 0755)
			if err != nil {
				return nil, fmt.Errorf("failed to create root directory: %v", err)
			}
		} else if err != nil {
			return nil, fmt.Errorf("failed to stat root directory: %v", err)
		}
	}

	var metadataFilePath string
	var metadata *tableMetadata

	for _, root := range roots {
		possibleMetadataPath := metadataPath(root)
		_, err := os.Stat(possibleMetadataPath)
		if err == nil {
			// We've found an existing metadata file. Use it.
			metadataFilePath = possibleMetadataPath
		} else if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to stat metadata file: %v", err)
		}
	}
	if metadataFilePath == "" {
		// No metadata file exists yet. Create a new one in the first root.
		var err error
		metadataDir := roots[0]
		metadata, err = newTableMetadata(logger, metadataDir, 0, shardingFactor)
		if err != nil {
			return nil, fmt.Errorf("failed to create table metadata: %v", err)
		}
	} else {
		// Metadata file exists, so we need to load it.
		var err error
		metadataDir := path.Dir(metadataFilePath)
		metadata, err = loadTableMetadata(logger, metadataDir)
		if err != nil {
			return nil, fmt.Errorf("failed to load table metadata: %v", err)
		}
	}

	table := &DiskTable{
		ctx:                     ctx,
		logger:                  logger,
		timeSource:              timeSource,
		roots:                   roots,
		segmentDirectories:      segDirs,
		name:                    name,
		metadata:                metadata,
		salt:                    salt,
		keyMap:                  keyMap,
		targetFileSize:          targetFileSize,
		segments:                make(map[uint32]*segment.Segment),
		controllerChan:          make(chan any, controlChannelSize),
		stopChannel:             make(chan struct{}, 1),
		garbageCollectionPeriod: gcPeriod,
	}
	table.alive.Store(true)
	table.stopChannel <- struct{}{}

	var err error
	table.lowestSegmentIndex, table.highestSegmentIndex, table.segments, err =
		segment.GatherSegmentFiles(
			ctx,
			logger,
			table.segmentDirectories,
			timeSource(),
			shardingFactor,
			salt,
			true)
	if err != nil {
		return nil, fmt.Errorf("failed to gather segment files: %v", err)
	}

	err = table.keyMap.LoadFromSegments(table.segments, table.lowestSegmentIndex, table.highestSegmentIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to load key map from segments: %v", err)
	}

	return table, nil
}

func (d *DiskTable) Name() string {
	return d.name
}

// Start starts the disk table.
func (d *DiskTable) Start() error {
	if !d.alive.Load() {
		return fmt.Errorf("DB is shut down")
	}
	go d.controlLoop()
	return nil
}

// Stop stops the disk table. Flushes all data out to disk.
func (d *DiskTable) Stop() error {
	alive := d.alive.Swap(false)
	if alive {
		flushReq := &flushRequest{
			shutdown:     true,
			responseChan: make(chan error, 1),
		}
		d.controllerChan <- flushReq
		err := <-flushReq.responseChan

		if err != nil {
			return fmt.Errorf("failed to flush: %v", err)
		}
	}

	// Wait for the control loop to stop.
	d.stopChannel <- struct{}{}
	<-d.stopChannel

	return nil
}

// Destroy stops the disk table and delete all files.
func (d *DiskTable) Destroy() error {
	err := d.Stop()
	if err != nil {
		return fmt.Errorf("failed to stop: %v", err)
	}

	d.logger.Infof("deleting disk table at path(s): %v", d.roots)

	for _, seg := range d.segments {
		seg.Release()
	}
	for _, seg := range d.segments {
		seg.BlockUntilFullyDeleted()
	}

	for _, segDir := range d.segmentDirectories {
		err = os.Remove(segDir)
		if err != nil {
			return fmt.Errorf("failed to remove root directory: %v", err)
		}
	}

	err = d.keyMap.Destroy()
	if err != nil {
		return fmt.Errorf("failed to destroy key map: %v", err)
	}

	err = d.metadata.delete()
	if err != nil {
		return fmt.Errorf("failed to delete metadata: %v", err)
	}

	for _, root := range d.roots {
		err = os.Remove(root)
		if err != nil {
			return fmt.Errorf("failed to remove root directory: %v", err)
		}
	}

	return nil
}

// SetTTL sets the TTL for the disk table. If set to 0, no TTL is enforced. This setting effects both new
// data and data already written.
func (d *DiskTable) SetTTL(ttl time.Duration) error {
	if !d.alive.Load() {
		return fmt.Errorf("DB is shut down")
	}

	err := d.metadata.SetTTL(ttl)
	if err != nil {
		return fmt.Errorf("failed to set TTL: %v", err)
	}
	return nil
}

func (d *DiskTable) SetShardingFactor(shardingFactor uint32) error {
	if !d.alive.Load() {
		return fmt.Errorf("DB is shut down")
	}

	err := d.metadata.SetShardingFactor(shardingFactor)
	if err != nil {
		return fmt.Errorf("failed to set sharding factor: %v", err)
	}

	return nil
}

// getSegment returns the segment with the given index. Segment is reserved, and it is the caller's responsibility to
// release the reservation when done.
func (d *DiskTable) getReservedSegment(index uint32) (*segment.Segment, bool) {
	d.segmentLock.RLock()
	defer d.segmentLock.RUnlock()

	seg, ok := d.segments[index]
	if !ok {
		return nil, false
	}

	ok = seg.Reserve()
	if !ok {
		// segmented was deleted out from under us
		return nil, false
	}

	return seg, true
}

func (d *DiskTable) Get(key []byte) ([]byte, bool, error) {
	if !d.alive.Load() {
		return nil, false, fmt.Errorf("DB is shut down")
	}

	// First, check if the key is in the unflushed data map.
	// If so, return it from there.
	if value, ok := d.unflushedDataCache.Load(string(key)); ok {
		return value.([]byte), true, nil
	}

	// Look up the address of the data.
	address, ok, err := d.keyMap.Get(key)
	if err != nil {
		return nil, false, fmt.Errorf("failed to get address: %v", err)
	}
	if !ok {
		return nil, false, nil
	}

	// Reserve the segment that contains the data.
	seg, ok := d.getReservedSegment(address.Index())
	if !ok {
		return nil, false, nil
	}
	defer seg.Release()

	// Read the data from disk.
	data, err := seg.Read(key, address)

	if err != nil {
		return nil, false, fmt.Errorf("failed to read data: %v", err)
	}

	return data, true, nil
}

func (d *DiskTable) Put(key []byte, value []byte) error {
	if !d.alive.Load() {
		return fmt.Errorf("DB is shut down")
	}

	d.unflushedDataCache.Store(string(key), value)

	writeReq := &writeRequest{
		values: make([]*types.KVPair, 1),
	}
	writeReq.values[0] = &types.KVPair{
		Key:   key,
		Value: value,
	}

	d.controllerChan <- writeReq

	return nil
}

func (d *DiskTable) PutBatch(batch []*types.KVPair) error {
	if !d.alive.Load() {
		return fmt.Errorf("DB is shut down")
	}

	for _, kv := range batch {
		d.unflushedDataCache.Store(string(kv.Key), kv.Value)
	}

	request := &writeRequest{
		values: batch,
	}
	d.controllerChan <- request
	return nil
}

// writeRequest is a request to write a key-value pair.
type writeRequest struct {
	values []*types.KVPair
}

// handleWriteRequest handles a writeRequest control message.
func (d *DiskTable) handleWriteRequest(req *writeRequest) {
	for _, kv := range req.values {
		// Do the write.
		seg := d.segments[d.highestSegmentIndex]
		shardSize, err := seg.Write(kv)
		if err != nil {
			d.panic(fmt.Errorf("failed to write to segment %d: %v", d.highestSegmentIndex, err))
		}

		// Check to see if the write caused the mutable segment to become full.
		if shardSize > uint64(d.targetFileSize) {
			// Mutable segment is full. Before continuing, we need to expand the segments.
			err = d.expandSegments()
			if err != nil {
				d.panic(fmt.Errorf("failed to expand segments: %v", err))
			}
		}
	}
}

// panic! Something just went very wrong. (╯°□°)╯︵ ┻━┻
func (d *DiskTable) panic(err error) {
	d.fatalError.Store(true)
	d.logger.Fatalf("unrecoverable DB error, database is shutting down: %v", err)
	err = d.Stop()
	if err != nil {
		d.logger.Fatalf("failed to stop DB: %v", err)
	}
}

// expandSegments checks if the highest segment is full, and if so, creates a new segment.
func (d *DiskTable) expandSegments() error {
	// Seal the previous segment.
	now := d.timeSource()
	durableKeys, err := d.segments[d.highestSegmentIndex].Seal(now)
	if err != nil {
		return fmt.Errorf("failed to seal segment %d: %v", d.highestSegmentIndex, err)
	}

	// Flush the keys that are now durable in the segment.
	err = d.flushKeys(durableKeys)
	if err != nil {
		return fmt.Errorf("failed to flush keys: %v", err)
	}

	// Create a new segment.
	newSegment, err := segment.NewSegment(
		d.ctx,
		d.logger,
		d.highestSegmentIndex+1,
		d.segmentDirectories,
		now,
		d.metadata.GetShardingFactor(),
		d.salt,
		false)
	if err != nil {
		d.segmentLock.Unlock()
		return fmt.Errorf("failed to create new segment: %v", err)
	}
	d.segments[d.highestSegmentIndex].SetNextSegment(newSegment)
	d.highestSegmentIndex++

	d.segmentLock.Lock()
	d.segments[d.highestSegmentIndex] = newSegment
	d.segmentLock.Unlock()

	return nil
}

// Flush flushes all data to disk. Blocks until all data previously submitted to Put has been written to disk.
func (d *DiskTable) Flush() error {
	if !d.alive.Load() {
		return fmt.Errorf("DB is shut down")
	}

	flushReq := &flushRequest{
		responseChan: make(chan error, 1),
	}
	d.controllerChan <- flushReq
	err := <-flushReq.responseChan
	if err != nil {
		return fmt.Errorf("failed to flush: %v", err)
	}

	return nil
}

// flushRequest is a request to flush the writer.
type flushRequest struct {
	shutdown     bool
	responseChan chan error
}

// handleFlushRequest handles a flushRequest control message.
func (d *DiskTable) handleFlushRequest(req *flushRequest) {

	var durableKeys []*types.KAPair
	var err error
	if req.shutdown {
		durableKeys, err = d.segments[d.highestSegmentIndex].Seal(d.timeSource())
	} else {
		durableKeys, err = d.segments[d.highestSegmentIndex].Flush()
	}

	if err != nil {
		err = fmt.Errorf("failed to flush mutable segment: %v", err)
		req.responseChan <- err
		return
	}

	err = d.flushKeys(durableKeys)
	if err != nil {
		err = fmt.Errorf("failed to flush keys: %v", err)
		d.panic(err)
		req.responseChan <- err
		return
	}

	if req.shutdown {
		err = d.keyMap.Stop()
		if err != nil {
			err = fmt.Errorf("failed to stop key map: %v", err)
			d.panic(err)
			req.responseChan <- err
			return
		}
	}

	req.responseChan <- nil
}

// flushKeys flushes all keys to the key map. As they are flushed, it also removes them from the unflushedDataCache.
func (d *DiskTable) flushKeys(keys []*types.KAPair) error {
	if len(keys) == 0 {
		// Nothing to flush.
		return nil
	}

	err := d.keyMap.Put(keys)
	if err != nil {
		return fmt.Errorf("failed to flush keys: %v", err)
	}

	// Keys are now durably written to both the segment and the key map. It is therefore safe to remove them from the
	// unflushed data cache.
	for _, ka := range keys {
		d.unflushedDataCache.Delete(string(ka.Key))
	}

	return nil
}

// controlLoop is the main loop for the disk table. It has sole responsibility for mutating data managed by the
// disk table (this vastly simplifies locking and thread safety).
func (d *DiskTable) controlLoop() {
	defer d.shutdownTasks()

	ticker := time.NewTicker(d.garbageCollectionPeriod)

	for {
		select {
		case <-d.ctx.Done():
			d.logger.Infof("context done, shutting down disk table control loop")
		case message := <-d.controllerChan:
			if writeReq, ok := message.(*writeRequest); ok {
				d.handleWriteRequest(writeReq)
			} else if flushReq, ok := message.(*flushRequest); ok {
				d.handleFlushRequest(flushReq)
				if flushReq.shutdown {
					<-d.stopChannel
					return
				}
			} else {
				d.logger.Errorf("Unknown control message type %T", message)
			}
		case <-ticker.C:
			d.doGarbageCollection()
		}
	}
}

// shutdownTasks performs tasks necessary to cleanly shut down the disk table.
func (d *DiskTable) shutdownTasks() {
	// Seal the mutable segment
	durableKeys, err := d.segments[d.highestSegmentIndex].Seal(d.timeSource())
	if err != nil {
		d.logger.Errorf("failed to seal mutable segment: %v", err)
	}

	// Flush the keys that are now durable in the segment.
	err = d.flushKeys(durableKeys)
	if err != nil {
		d.logger.Errorf("failed to flush keys: %v", err)
	}

	// Stop the key map
	err = d.keyMap.Stop()
	if err != nil {
		d.logger.Errorf("failed to stop key map: %v", err)
	}

	// unblock the Stop() method
	<-d.stopChannel
}

// doGarbageCollection performs garbage collection on all segments, deleting old ones as necessary.
func (d *DiskTable) doGarbageCollection() {

	defer func() {
		fmt.Printf("    Finished GC\n") // TODO
	}()

	fmt.Printf("Waking up to do GC, ttl is %v\n", d.metadata.GetTTL()) // TODO

	now := d.timeSource()
	ttl := d.metadata.GetTTL()
	if ttl.Nanoseconds() == 0 {
		// No TTL set, so nothing to do.
		return
	}

	for index := d.lowestSegmentIndex; index <= d.highestSegmentIndex; index++ {
		seg := d.segments[index]
		if !seg.IsSealed() {
			// We can't delete an unsealed segment.
			return
		}

		sealTime := seg.GetSealTime()
		segmentAge := now.Sub(sealTime)
		if segmentAge < ttl {
			fmt.Printf("    Segment %d is not old enough to delete, age is %v\n", index, segmentAge) // TODO
			// Segment is not old enough to be deleted.
			return
		}

		fmt.Printf("    Segment %d is old enough to delete, age is %v\n", index, segmentAge) // TODO

		// Segment is old enough to be deleted.

		keys, err := seg.GetKeys()
		if err != nil {
			d.logger.Errorf("Failed to get keys: %v", err)
			return
		}

		err = d.keyMap.Delete(keys)
		if err != nil {
			d.logger.Errorf("Failed to delete keys: %v", err)
			return
		}

		// Deletion of segment files will happen when the segment is released by all reservation holders.
		seg.Release()
		d.segmentLock.Lock()
		delete(d.segments, index)
		d.segmentLock.Unlock()

		d.lowestSegmentIndex++
	}
}

func (d *DiskTable) SetCacheSize(size uint64) {
	// this implementation does not provide a cache, if a cache is needed then it must be provided by a wrapper
}
