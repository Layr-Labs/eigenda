package disktable

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/litt"
	"github.com/Layr-Labs/eigenda/litt/disktable/keymap"
	"github.com/Layr-Labs/eigenda/litt/disktable/segment"
	"github.com/Layr-Labs/eigenda/litt/types"
	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

var _ litt.ManagedTable = (*DiskTable)(nil)

// segmentDirectory is the directory where segment files are stored, relative to the root directory.
const segmentDirectory = "segments"

// keyMapReloadBatchSize is the size of the batch used for reloading keys from segments into the key map.
const keyMapReloadBatchSize = 1024

const tableFlushChannelCapacity = 8

// DiskTable manages a table's Segments.
type DiskTable struct {
	// The logger for the disk table.
	logger logging.Logger

	// panic is a struct that permits the DB to "panic". There are many goroutines that function under the hood, and
	// many of these threads could, in theory, encounter errors which are unrecoverable. In such situations, the
	// desirable outcome is for the DB to report the error and then refuse to do additional work. If the DB is in a
	// broken state, it is much better to refuse to do work than to continue to do work and potentially corrupt data.
	panic *util.DBPanic

	// The root directories for the disk table.
	roots []string

	// The directories where segment files are stored.
	segmentDirectories []string

	// The table's name.
	name string

	// The table's metadata.
	metadata *tableMetadata

	// A source of randomness used for generating sharding salt.
	saltShaker *rand.Rand

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

	// controllerChannel is the channel for messages sent to controller goroutine. No data managed by the DiskTable
	// may be mutated by anything other than the controller, with the exception of the
	// segmentLock and reference counting.
	controllerChannel chan any

	// flushChannel is a channel used to enqueue work on the flush loop goroutine.
	flushChannel chan any

	// garbageCollectionPeriod is the period at which garbage collection is run.
	garbageCollectionPeriod time.Duration

	// timeSource is the time source used by the disk table.
	timeSource func() time.Time
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
	saltShaker *rand.Rand,
	ttl time.Duration,
	gcPeriod time.Duration,
	reloadKeyMap bool) (litt.ManagedTable, error) {

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
		metadata, err = newTableMetadata(logger, metadataDir, ttl, shardingFactor)
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

	dbPanic := util.NewDBPanic(ctx, logger)

	table := &DiskTable{
		logger:                  logger,
		panic:                   dbPanic,
		timeSource:              timeSource,
		roots:                   roots,
		segmentDirectories:      segDirs,
		name:                    name,
		metadata:                metadata,
		saltShaker:              saltShaker,
		keyMap:                  keyMap,
		targetFileSize:          targetFileSize,
		segments:                make(map[uint32]*segment.Segment),
		controllerChannel:       make(chan any, controlChannelSize),
		flushChannel:            make(chan any, tableFlushChannelCapacity),
		garbageCollectionPeriod: gcPeriod,
	}

	var err error
	table.lowestSegmentIndex, table.highestSegmentIndex, table.segments, err =
		segment.GatherSegmentFiles(
			logger,
			dbPanic,
			table.segmentDirectories,
			timeSource(),
			shardingFactor,
			saltShaker.Uint32(),
			true)
	if err != nil {
		return nil, fmt.Errorf("failed to gather segment files: %v", err)
	}

	if reloadKeyMap {
		logger.Infof("reloading key map from segments")
		err = table.reloadKeyMap()
		if err != nil {
			return nil, fmt.Errorf("failed to load key map from segments: %v", err)
		}
	}

	return table, nil
}

// reloadKeyMap reloads the key map from the segments. This is necessary when the key map is lost, the key map doesn't
// save its data on disk, or we are migrating from one key map type to another.
func (d *DiskTable) reloadKeyMap() error {

	start := d.timeSource()
	defer func() {
		d.logger.Infof("spent %v reloading key map", d.timeSource().Sub(start))
	}()

	// It's possible that some of the data written near the end of the previous session was corrupted.
	// Read data from the end until the first valid key/value pair is found.
	isValid := false

	batch := make([]*types.KAPair, 0, keyMapReloadBatchSize)

	for i := d.highestSegmentIndex; i >= d.lowestSegmentIndex && i+1 != 0; i-- {
		if !d.segments[i].IsSealed() {
			// ignore unsealed segment, this will have been created in the current session and will not
			// yet contain any data.
			continue
		}

		keys, err := d.segments[i].GetKeys()
		if err != nil {
			return fmt.Errorf("failed to get keys from segment: %v", err)
		}
		for keyIndex := len(keys) - 1; keyIndex >= 0; keyIndex-- {
			key := keys[keyIndex]

			if !isValid {
				_, err = d.segments[i].Read(key.Key, key.Address)
				if err == nil {
					// we found a valid key/value pair. All subsequent keys are valid.
					isValid = true
				} else {
					// This is not cause for alarm (probably).
					// This can happen when the database is not cleanly shut down,
					// and just means that some data near the end was not fully committed.
					d.logger.Infof("truncated value for key %s with address %s", key.Key, key.Address)
				}
			}

			if isValid {
				batch = append(batch, key)
				if len(batch) == keyMapReloadBatchSize {
					err = d.keyMap.Put(batch)
					if err != nil {
						return fmt.Errorf("failed to put keys: %v", err)
					}
					batch = make([]*types.KAPair, 0, keyMapReloadBatchSize)
				}
			}
		}

		if len(batch) > 0 {
			err = d.keyMap.Put(batch)
			if err != nil {
				return fmt.Errorf("failed to put keys: %v", err)
			}
		}
	}

	return nil
}

func (d *DiskTable) Name() string {
	return d.name
}

// Start starts the disk table.
func (d *DiskTable) Start() error {
	if ok, err := d.panic.IsOk(); !ok {
		return fmt.Errorf("Cannot process Start() request: %v", err)
	}

	go d.controlLoop()
	go d.flushLoop()
	return nil
}

// Stop stops the disk table. Flushes all data out to disk.
func (d *DiskTable) Stop() error {
	if ok, err := d.panic.IsOk(); !ok {
		return fmt.Errorf("Cannot process Stop() request: %v", err)
	}

	d.panic.Shutdown()

	shutdownCompleteChan := make(chan struct{}, 1)
	request := &controlLoopShutdownRequest{
		shutdownCompleteChan: shutdownCompleteChan,
	}
	err := util.SendAny(d.panic, d.controllerChannel, request)
	if err != nil {
		return fmt.Errorf("failed to send shutdown request: %v", err)
	}

	_, err = util.Await(d.panic, shutdownCompleteChan) // TODO don't return an error on the channel
	if err != nil {
		return fmt.Errorf("failed to shutdown: %v", err)
	}

	return nil
}

// Destroy stops the disk table and delete all files.
func (d *DiskTable) Destroy() error {
	if ok, err := d.panic.IsOk(); !ok {
		return fmt.Errorf("Cannot process Destroy() request: %v", err)
	}

	err := d.Stop()
	if err != nil {
		return fmt.Errorf("failed to stop: %v", err)
	}

	d.logger.Infof("deleting disk table at path(s): %v", d.roots)

	for _, seg := range d.segments {
		seg.Release()
	}
	for _, seg := range d.segments {
		err = seg.BlockUntilFullyDeleted()
		if err != nil {
			return fmt.Errorf("failed to delete segment: %v", err)
		}
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
	if ok, err := d.panic.IsOk(); !ok {
		return fmt.Errorf("Cannot process SetTTL() request: %v", err)
	}

	err := d.metadata.SetTTL(ttl)
	if err != nil {
		return fmt.Errorf("failed to set TTL: %v", err)
	}
	return nil
}

// TODO test this on a live table

func (d *DiskTable) SetShardingFactor(shardingFactor uint32) error {
	if ok, err := d.panic.IsOk(); !ok {
		return fmt.Errorf("Cannot process SetShardingFactor() request: %v", err)
	}

	if shardingFactor == 0 {
		return fmt.Errorf("sharding factor must be greater than 0")
	}

	request := &controlLoopSetShardingFactorRequest{
		shardingFactor: shardingFactor,
	}
	err := util.SendAny(d.panic, d.controllerChannel, request)
	if err != nil {
		return fmt.Errorf("failed to send sharding factor request: %v", err)
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
	if ok, err := d.panic.IsOk(); !ok {
		return nil, false, fmt.Errorf("Cannot process Get() request: %v", err)
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
	if ok, err := d.panic.IsOk(); !ok {
		return fmt.Errorf("Cannot process Put() request: %v", err)
	}

	d.unflushedDataCache.Store(string(key), value)

	writeReq := &controlLoopWriteRequest{
		values: make([]*types.KVPair, 1),
	}
	writeReq.values[0] = &types.KVPair{
		Key:   key,
		Value: value,
	}

	err := util.SendAny(d.panic, d.controllerChannel, writeReq)
	if err != nil {
		return fmt.Errorf("failed to send write request: %v", err)
	}

	return nil
}

func (d *DiskTable) PutBatch(batch []*types.KVPair) error {
	if ok, err := d.panic.IsOk(); !ok {
		return fmt.Errorf("Cannot process PutBatch() request: %v", err)
	}

	for _, kv := range batch {
		d.unflushedDataCache.Store(string(kv.Key), kv.Value)
	}

	request := &controlLoopWriteRequest{
		values: batch,
	}
	err := util.SendAny(d.panic, d.controllerChannel, request)
	if err != nil {
		return fmt.Errorf("failed to send write request: %v", err)
	}
	return nil
}

// handleWriteRequest handles a controlLoopWriteRequest control message.
func (d *DiskTable) handleWriteRequest(req *controlLoopWriteRequest) {
	for _, kv := range req.values {
		// Do the write.
		seg := d.segments[d.highestSegmentIndex]
		shardSize, err := seg.Write(kv)
		if err != nil {
			d.panic.Panic(fmt.Errorf("failed to write to segment %d: %v", d.highestSegmentIndex, err))
		}

		// Check to see if the write caused the mutable segment to become full.
		if shardSize > uint64(d.targetFileSize) {
			// Mutable segment is full. Before continuing, we need to expand the segments.
			err = d.expandSegments()
			if err != nil {
				d.panic.Panic(fmt.Errorf("failed to expand segments: %v", err))
			}
		}
	}
}

// TODO ensure that if we panic, this doesn't block forever

// Flush flushes all data to disk. Blocks until all data previously submitted to Put has been written to disk.
func (d *DiskTable) Flush() error {
	if ok, err := d.panic.IsOk(); !ok {
		return fmt.Errorf("Cannot process Flush() request: %v", err)
	}

	flushReq := &controlLoopFlushRequest{
		responseChan: make(chan error, 1),
	}
	err := util.SendAny(d.panic, d.controllerChannel, flushReq)
	if err != nil {
		return fmt.Errorf("failed to send flush request: %v", err)
	}

	_, err = util.Await(d.panic, flushReq.responseChan)
	if err != nil {
		return fmt.Errorf("failed to flush: %v", err)
	}

	return nil
}

// handleControlLoopShutdownRequest performs tasks necessary to cleanly shut down the disk table.
func (d *DiskTable) handleControlLoopShutdownRequest(req *controlLoopShutdownRequest) {
	// Instruct the flush loop to stop.
	shutdownCompleteChan := make(chan struct{})
	request := &flushLoopShutdownRequest{
		shutdownCompleteChan: shutdownCompleteChan,
	}
	err := util.SendAny(d.panic, d.flushChannel, request)
	if err != nil {
		d.logger.Errorf("failed to send shutdown request to flush loop: %v", err)
		return
	}

	_, err = util.Await(d.panic, shutdownCompleteChan)
	if err != nil {
		d.logger.Errorf("failed to shutdown flush loop: %v", err)
		return
	}

	// Seal the mutable segment
	durableKeys, err := d.segments[d.highestSegmentIndex].Seal(d.timeSource())
	if err != nil {
		d.panic.Panic(fmt.Errorf("failed to seal mutable segment: %v", err))
	}

	// Flush the keys that are now durable in the segment.
	err = d.writeKeysToKeyMap(durableKeys)
	if err != nil {
		d.panic.Panic(fmt.Errorf("failed to flush keys: %v", err))
	}

	// Stop the key map
	err = d.keyMap.Stop()
	if err != nil {
		d.panic.Panic(fmt.Errorf("failed to stop key map: %v", err))
	}

	req.shutdownCompleteChan <- struct{}{}
}

// doGarbageCollection performs garbage collection on all segments, deleting old ones as necessary.
func (d *DiskTable) doGarbageCollection() {
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
			// Segment is not old enough to be deleted.
			return
		}

		// Segment is old enough to be deleted.
		keys, err := seg.GetKeys()
		if err != nil {
			d.panic.Panic(fmt.Errorf("failed to get keys: %v", err))
			return
		}

		err = d.keyMap.Delete(keys)
		if err != nil {
			d.panic.Panic(fmt.Errorf("failed to delete keys: %v", err))
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

func (d *DiskTable) SetCacheSize(_ uint64) error {
	if ok, err := d.panic.IsOk(); !ok {
		return fmt.Errorf("Cannot process SetCacheSize() request: %v", err)
	}

	// this implementation does not provide a cache, if a cache is needed then it must be provided by a wrapper
	return nil
}

// expandSegments checks if the highest segment is full, and if so, creates a new segment.
func (d *DiskTable) expandSegments() error {
	now := d.timeSource()

	// Seal the previous segment.
	flushLoopResponseChan := make(chan error, 1)
	request := &flushLoopSealRequest{
		now:          now,
		responseChan: flushLoopResponseChan,
	}
	err := util.SendAny(d.panic, d.flushChannel, request)
	if err != nil {
		return fmt.Errorf("failed to send seal request: %v", err)
	}

	// Unfortunately, it is necessary to block until the sealing has been completed. Although this may result
	// in a brief interruption in new write work being sent to the segment, expanding the number of segments is
	// infrequent, even for very high throughput workloads.
	_, err = util.Await(d.panic, flushLoopResponseChan)
	if err != nil {
		return fmt.Errorf("failed to seal segment: %v", err)
	}

	// Create a new segment.
	newSegment, err := segment.NewSegment(
		d.logger,
		d.panic,
		d.highestSegmentIndex+1,
		d.segmentDirectories,
		now,
		d.metadata.GetShardingFactor(),
		d.saltShaker.Uint32(),
		false)
	if err != nil {
		return err
	}
	d.segments[d.highestSegmentIndex].SetNextSegment(newSegment)
	d.highestSegmentIndex++

	d.segmentLock.Lock()
	d.segments[d.highestSegmentIndex] = newSegment
	d.segmentLock.Unlock()

	return nil
}

// handleFlushLoopSealRequest handles the part of the seal operation that is performed on the flush loop.
// We don't want to send a flush request to a segment that has already been sealed. By performing the sealing
// on the flush loop, we ensure that this can never happen. Any previously scheduled flush requests against the
// segment that is being sealed will be processed prior to this request being processed due to the FIFO nature
// of the flush loop channel. When a sealing operation begins, the control loop blocks, and does not unblock until
// the seal is finished and a new mutable segment has been created. This means that no future flush requests will be
// sent to the segment that is being sealed, since only the control loop can schedule work for the flush loop.
func (d *DiskTable) handleFlushLoopSealRequest(req *flushLoopSealRequest) {
	durableKeys, err := d.segments[d.highestSegmentIndex].Seal(req.now)
	if err != nil {
		req.responseChan <- fmt.Errorf("failed to seal segment %d: %v", d.highestSegmentIndex, err)
		return
	}

	// Flush the keys that are now durable in the segment.
	err = d.writeKeysToKeyMap(durableKeys)
	if err != nil {
		req.responseChan <- fmt.Errorf("failed to flush keys: %v", err)
		return
	}

	req.responseChan <- nil
}

// handleControlLoopFlushRequest handles the part of the flush that is performed on the control loop.
// The control loop is responsible for enqueuing the flush request in the segment's work queue (thus
// ensuring a serial ordering with respect to other operations on the control loop), but not for
// waiting for the segment to finish the flush.
func (d *DiskTable) handleControlLoopFlushRequest(req *controlLoopFlushRequest) {
	// This method will enqueue a flush operation within the segment. Once that is done,
	// it becomes the responsibility of the flush loop to wait for the flush to complete.
	flushWaitFunction, err := d.segments[d.highestSegmentIndex].Flush()
	if err != nil {
		d.panic.Panic(fmt.Errorf("failed to flush segment %d: %v", d.highestSegmentIndex, err))
	}

	// The flush loop is responsible for the remaining parts of the flush.
	request := &flushLoopFlushRequest{
		flushWaitFunction: flushWaitFunction,
		responseChan:      req.responseChan,
	}
	err = util.SendAny(d.panic, d.flushChannel, request)
	if err != nil {
		d.logger.Errorf("failed to send flush request to flush loop: %v", err)
	}
}

// handleFlushLoopFlushRequest handles the part of the flush that is performed on the flush loop.
func (d *DiskTable) handleFlushLoopFlushRequest(req *flushLoopFlushRequest) {
	durableKeys, err := req.flushWaitFunction()
	if err != nil {
		err = fmt.Errorf("failed to flush mutable segment: %v", err)
		req.responseChan <- fmt.Errorf("failed to flush mutable segment: %v", err)
		d.panic.Panic(err)
	}

	err = d.writeKeysToKeyMap(durableKeys)
	if err != nil {
		err = fmt.Errorf("failed to flush keys: %v", err)
		req.responseChan <- err
		d.panic.Panic(err)
	}

	req.responseChan <- nil
}

// writeKeysToKeyMap flushes all keys to the key map. Once they are flushed, it also removes the keys from the
// unflushedDataCache.
func (d *DiskTable) writeKeysToKeyMap(keys []*types.KAPair) error {
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

// handleControlLoopSetShardingFactorRequest updates the sharding factor of the disk table. If the requested
// sharding factor is the same as before, no action is taken. If it is different, the sharding factor is updated,
// the current mutable segment is sealed, and a new mutable segment is created.
func (d *DiskTable) handleControlLoopSetShardingFactorRequest(req *controlLoopSetShardingFactorRequest) {

	if req.shardingFactor == d.metadata.GetShardingFactor() {
		// No action necessary.
		return
	}
	err := d.metadata.SetShardingFactor(req.shardingFactor)
	if err != nil {
		d.panic.Panic(fmt.Errorf("failed to set sharding factor: %v", err))
	}

	// This seals the current mutable segment and creates a new one. The new segment will have the new sharding factor.
	err = d.expandSegments()
	if err != nil {
		d.panic.Panic(fmt.Errorf("failed to expand segments: %v", err))
	}
}

// controlLoopFlushRequest is a request to flush the writer that is sent to the control loop.
type controlLoopFlushRequest struct {
	// responseChan will produce a nil if the flush was successful, or an error if it was not.
	responseChan chan error
}

// controlLoopWriteRequest is a request to write a key-value pair that is sent to the control loop.
type controlLoopWriteRequest struct {
	// values is a slice of key-value pairs to write.
	values []*types.KVPair
}

// controlLoopSetShardingFactorRequest is a request to set the sharding factor that is sent to the control loop.
type controlLoopSetShardingFactorRequest struct {
	// shardingFactor is the new sharding factor to set.
	shardingFactor uint32
}

// controlLoopShutdownRequest is a request to shut down the table that is sent to the control loop.
type controlLoopShutdownRequest struct {
	// responseChan will produce a single struct{} when the control loop has stopped
	// (i.e. when the handleControlLoopShutdownRequest is complete).
	shutdownCompleteChan chan struct{}
}

// controlLoop is the main loop for the disk table. It has sole responsibility for mutating data managed by the
// disk table (this vastly simplifies locking and thread safety).
func (d *DiskTable) controlLoop() {
	ticker := time.NewTicker(d.garbageCollectionPeriod)

	for {
		select {
		case <-d.panic.ImmediateShutdownRequired():
			d.logger.Infof("context done, shutting down disk table control loop")
		case message := <-d.controllerChannel:
			if req, ok := message.(*controlLoopWriteRequest); ok {
				d.handleWriteRequest(req)
			} else if req, ok := message.(*controlLoopFlushRequest); ok {
				d.handleControlLoopFlushRequest(req)
			} else if req, ok := message.(*controlLoopSetShardingFactorRequest); ok {
				d.handleControlLoopSetShardingFactorRequest(req)
			} else if req, ok := message.(*controlLoopShutdownRequest); ok {
				d.handleControlLoopShutdownRequest(req)
				return
			} else {
				d.panic.Panic(fmt.Errorf("Unknown control message type %T", message))
			}
		case <-ticker.C:
			d.doGarbageCollection()
		}
	}
}

// flushLoopFlushRequest is a request to flush the writer that is sent to the flush loop.
type flushLoopFlushRequest struct {
	// flushWaitFunction is the function that will wait for the flush to complete.
	flushWaitFunction segment.FlushWaitFunction

	// responseChan the flush loop sends a nil if the flush was successfully completed, or an error if it was not.
	responseChan chan error
}

// flushLoopSealRequest is a request to seal the mutable segment that is sent to the flush loop.
type flushLoopSealRequest struct {
	// the time when the segment is sealed
	now time.Time
	// responseChan will produce a nil if the seal was successful, or an error if it was not.
	responseChan chan error
}

// flushLoopShutdownRequest is a request to shut down the flush loop.
type flushLoopShutdownRequest struct {
	// responseChan will produce a single struct{} when the flush loop has stopped.
	shutdownCompleteChan chan struct{}
}

// flushLoop is responsible for handling operations that flush data (i.e. calls to Flush() and when the mutable segment
// is sealed). In theory, this work could be done on the main control loop, but doing so would block new writes while
// a flush is in progress. In order to keep the writing threads busy, it is critical that flush do not block the
// control loop.
func (d *DiskTable) flushLoop() {
	for {
		select {
		case <-d.panic.ImmediateShutdownRequired():
			d.logger.Infof("context done, shutting down disk table flush loop")
		case message := <-d.flushChannel:
			if req, ok := message.(*flushLoopFlushRequest); ok {
				d.handleFlushLoopFlushRequest(req)
			} else if req, ok := message.(*flushLoopSealRequest); ok {
				d.handleFlushLoopSealRequest(req)
			} else if req, ok := message.(*flushLoopShutdownRequest); ok {
				req.shutdownCompleteChan <- struct{}{}
				return
			} else {
				d.panic.Panic(fmt.Errorf("Unknown flush message type %T", message))
			}
		}
	}
}
