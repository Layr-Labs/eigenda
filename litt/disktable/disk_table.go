package disktable

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Layr-Labs/eigenda/litt"
	"github.com/Layr-Labs/eigenda/litt/disktable/keymap"
	"github.com/Layr-Labs/eigenda/litt/disktable/segment"
	"github.com/Layr-Labs/eigenda/litt/metrics"
	"github.com/Layr-Labs/eigenda/litt/types"
	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

var _ litt.ManagedTable = (*DiskTable)(nil)

// segmentDirectory is the directory where segment files are stored, relative to the root directory.
const segmentDirectory = "segments"

// keymapReloadBatchSize is the size of the batch used for reloading keys from segments into the keymap.
const keymapReloadBatchSize = 1024

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
	keymap keymap.Keymap

	// The path to the keymap directory.
	keymapPath string

	// The type file for the keymap.
	keymapTypeFile *keymap.KeymapTypeFile

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

	// The maximum number of keys in a segment.
	maxKeyCount uint64

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

	// whether fsync mode is enabled.
	fsync bool

	// The number of bytes contained within the immutable segments. This tracks the number of bytes that are
	// on disk, not bytes in memory. For thread safety, this variable may only be read/written in the constructor
	// and in the control loop.
	immutableSegmentSize uint64

	// The number of bytes contained within all segments, including the mutable segment. This tracks the number of
	// bytes that are on disk, not bytes in memory.
	size atomic.Uint64

	// The number of keys in the table.
	keyCount atomic.Int64

	// Encapsulates metrics for the database.
	metrics *metrics.LittDBMetrics
}

// NewDiskTable creates a new DiskTable.
func NewDiskTable(
	ctx context.Context,
	logger logging.Logger,
	timeSource func() time.Time,
	name string,
	keymap keymap.Keymap,
	keymapPath string,
	keymapTypeFile *keymap.KeymapTypeFile,
	roots []string,
	targetFileSize uint32,
	controlChannelSize int,
	shardingFactor uint32,
	saltShaker *rand.Rand,
	ttl time.Duration,
	gcPeriod time.Duration,
	reloadKeymap bool,
	fsync bool,
	metrics *metrics.LittDBMetrics) (litt.ManagedTable, error) {

	if gcPeriod <= 0 {
		return nil, fmt.Errorf("garbage collection period must be greater than 0")
	}

	// If the root directories don't exist, create them.
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

	// For each root directory, create a segment directory if it doesn't exist.
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

	// Find the table metadata file or create a new one.
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
		keymap:                  keymap,
		keymapPath:              keymapPath,
		keymapTypeFile:          keymapTypeFile,
		targetFileSize:          targetFileSize,
		segments:                make(map[uint32]*segment.Segment),
		controllerChannel:       make(chan any, controlChannelSize),
		flushChannel:            make(chan any, tableFlushChannelCapacity),
		garbageCollectionPeriod: gcPeriod,
		fsync:                   fsync,
		metrics:                 metrics,
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
			fsync)
	if err != nil {
		return nil, fmt.Errorf("failed to gather segment files: %v", err)
	}

	for _, seg := range table.segments {
		table.immutableSegmentSize += seg.Size()
	}

	// Create the mutable segment
	creatingFirstSegment := len(table.segments) == 0

	var nextSegmentIndex uint32
	if creatingFirstSegment {
		nextSegmentIndex = 0
	} else {
		nextSegmentIndex = table.highestSegmentIndex + 1
	}
	mutableSegment, err := segment.NewSegment(
		logger,
		dbPanic,
		nextSegmentIndex,
		segDirs,
		timeSource(),
		metadata.GetShardingFactor(),
		saltShaker.Uint32(),
		false,
		fsync)
	if err != nil {
		return nil, fmt.Errorf("failed to create mutable segment: %v", err)
	}
	if !creatingFirstSegment {
		table.segments[table.highestSegmentIndex].SetNextSegment(mutableSegment)
		table.highestSegmentIndex++
	}
	table.segments[nextSegmentIndex] = mutableSegment

	table.updateCurrentSize()

	if reloadKeymap {
		logger.Infof("reloading keymap from segments")
		err = table.reloadKeymap()
		if err != nil {
			return nil, fmt.Errorf("failed to load keymap from segments: %v", err)
		}
	}

	go table.controlLoop()
	go table.flushLoop()

	return table, nil
}

func (d *DiskTable) KeyCount() uint64 {
	return uint64(d.keyCount.Load())
}

func (d *DiskTable) Size() uint64 {
	return d.size.Load()
}

// updateCurrentSize updates the size of the table.
func (d *DiskTable) updateCurrentSize() {
	size := d.immutableSegmentSize + d.segments[d.highestSegmentIndex].Size() + d.metadata.Size()
	d.size.Store(size)
}

// reloadKeymap reloads the keymap from the segments. This is necessary when the keymap is lost, the keymap doesn't
// save its data on disk, or we are migrating from one keymap type to another.
func (d *DiskTable) reloadKeymap() error {

	start := d.timeSource()
	defer func() {
		d.logger.Infof("spent %v reloading keymap", d.timeSource().Sub(start))
	}()

	// It's possible that some of the data written near the end of the previous session was corrupted.
	// Read data from the end until the first valid key/value pair is found.
	isValid := false

	batch := make([]*types.KAPair, 0, keymapReloadBatchSize)

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
				if len(batch) == keymapReloadBatchSize {
					err = d.keymap.Put(batch)
					if err != nil {
						return fmt.Errorf("failed to put keys: %v", err)
					}
					batch = make([]*types.KAPair, 0, keymapReloadBatchSize)
				}
			}
		}

	}

	if len(batch) > 0 {
		err := d.keymap.Put(batch)
		if err != nil {
			return fmt.Errorf("failed to put keys: %v", err)
		}
	}

	return nil
}

func (d *DiskTable) Name() string {
	return d.name
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

	err := util.Send(d.panic, d.controllerChannel, request)
	if err != nil {
		return fmt.Errorf("failed to send shutdown request: %v", err)
	}

	_, err = util.Await(d.panic, shutdownCompleteChan)
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

	// release all segments
	for _, seg := range d.segments {
		seg.Release()
	}
	// wait for segments to delete themselves
	for _, seg := range d.segments {
		err = seg.BlockUntilFullyDeleted()
		if err != nil {
			return fmt.Errorf("failed to delete segment: %v", err)
		}
	}

	// delete all segment directories
	for _, segDir := range d.segmentDirectories {
		err = os.Remove(segDir)
		if err != nil {
			return fmt.Errorf("failed to remove root directory: %v", err)
		}
	}

	// destroy the keymap
	err = d.keymap.Destroy()
	if err != nil {
		return fmt.Errorf("failed to destroy keymap: %v", err)
	}
	err = d.keymapTypeFile.Delete()
	if err != nil {
		return fmt.Errorf("failed to delete keymap type file: %v", err)
	}
	_, err = os.Stat(d.keymapPath)
	if err == nil {
		err = os.Remove(d.keymapPath)
		if err != nil {
			return fmt.Errorf("failed to remove keymap directory: %v", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat keymap directory: %v", err)
	}

	// delete the metadata file
	err = d.metadata.delete()
	if err != nil {
		return fmt.Errorf("failed to delete metadata: %v", err)
	}

	// delete the root directories for the table
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
	err := util.Send(d.panic, d.controllerChannel, request)
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

	var cacheHit bool
	var dataSize uint64
	if d.metrics != nil {
		start := d.timeSource()
		defer func() {
			end := d.timeSource()
			delta := end.Sub(start)
			d.metrics.ReportReadOperation(d.name, delta, dataSize, cacheHit)
		}()
	}

	// First, check if the key is in the unflushed data map.
	// If so, return it from there.
	if value, ok := d.unflushedDataCache.Load(string(key)); ok {
		bytes := value.([]byte)
		cacheHit = true
		dataSize = uint64(len(bytes))
		return bytes, true, nil
	}

	// Look up the address of the data.
	address, ok, err := d.keymap.Get(key)
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

	dataSize += uint64(len(data))

	return data, true, nil
}

func (d *DiskTable) Put(key []byte, value []byte) error {
	if ok, err := d.panic.IsOk(); !ok {
		return fmt.Errorf("Cannot process Put() request: %v", err)
	}

	if d.metrics != nil {
		start := d.timeSource()
		defer func() {
			end := d.timeSource()
			delta := end.Sub(start)
			d.metrics.ReportWriteOperation(d.name, delta, 1, uint64(len(value)))
		}()
	}

	d.unflushedDataCache.Store(string(key), value)

	writeReq := &controlLoopWriteRequest{
		values: make([]*types.KVPair, 1),
	}
	writeReq.values[0] = &types.KVPair{
		Key:   key,
		Value: value,
	}

	err := util.Send(d.panic, d.controllerChannel, writeReq)
	if err != nil {
		return fmt.Errorf("failed to send write request: %v", err)
	}

	d.keyCount.Add(1)

	return nil
}

func (d *DiskTable) PutBatch(batch []*types.KVPair) error {
	if ok, err := d.panic.IsOk(); !ok {
		return fmt.Errorf("Cannot process PutBatch() request: %v", err)
	}

	if d.metrics != nil {
		start := d.timeSource()
		totalSize := uint64(0)
		for _, kv := range batch {
			totalSize += uint64(len(kv.Value))
		}
		defer func() {
			end := d.timeSource()
			delta := end.Sub(start)
			d.metrics.ReportWriteOperation(d.name, delta, uint64(len(batch)), totalSize)
		}()
	}

	for _, kv := range batch {
		d.unflushedDataCache.Store(string(kv.Key), kv.Value)
	}

	request := &controlLoopWriteRequest{
		values: batch,
	}
	err := util.Send(d.panic, d.controllerChannel, request)
	if err != nil {
		return fmt.Errorf("failed to send write request: %v", err)
	}

	d.keyCount.Add(int64(len(batch)))

	return nil
}

// handleControlLoopWriteRequest handles a controlLoopWriteRequest control message.
func (d *DiskTable) handleControlLoopWriteRequest(req *controlLoopWriteRequest) {
	for _, kv := range req.values {
		// Do the write.
		seg := d.segments[d.highestSegmentIndex]
		keyCount, shardSize, err := seg.Write(kv)
		if err != nil {
			d.panic.Panic(fmt.Errorf("failed to write to segment %d: %v", d.highestSegmentIndex, err))
			return
		}

		// Check to see if the write caused the mutable segment to become full.
		if shardSize > uint64(d.targetFileSize) || keyCount >= d.maxKeyCount {
			// Mutable segment is full. Before continuing, we need to expand the segments.
			err = d.expandSegments()
			if err != nil {
				d.panic.Panic(fmt.Errorf("failed to expand segments: %v", err))
				return
			}
		}
	}

	d.updateCurrentSize()
}

func (d *DiskTable) Exists(key []byte) (bool, error) {
	_, ok := d.unflushedDataCache.Load(string(key))
	if ok {
		return true, nil
	}

	_, ok, err := d.keymap.Get(key)
	if err != nil {
		return false, fmt.Errorf("failed to get address: %v", err)
	}

	return ok, nil
}

// Flush flushes all data to disk. Blocks until all data previously submitted to Put has been written to disk.
func (d *DiskTable) Flush() error {
	if ok, err := d.panic.IsOk(); !ok {
		return fmt.Errorf("Cannot process Flush() request: %v", err)
	}

	if d.metrics != nil {
		start := d.timeSource()
		defer func() {
			end := d.timeSource()
			delta := end.Sub(start)
			d.metrics.ReportFlushOperation(d.name, delta)
		}()
	}

	flushReq := &controlLoopFlushRequest{
		responseChan: make(chan struct{}, 1),
	}
	err := util.Send(d.panic, d.controllerChannel, flushReq)
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
	err := util.Send(d.panic, d.flushChannel, request)
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
		return
	}

	// Flush the keys that are now durable in the segment.
	err = d.writeKeysToKeymap(durableKeys)
	if err != nil {
		d.panic.Panic(fmt.Errorf("failed to flush keys: %v", err))
		return
	}

	// Stop the keymap
	err = d.keymap.Stop()
	if err != nil {
		d.panic.Panic(fmt.Errorf("failed to stop keymap: %v", err))
		return
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

	if d.metrics != nil {
		start := d.timeSource()
		defer func() {
			end := d.timeSource()
			delta := end.Sub(start)
			d.metrics.ReportGarbageCollectionLatency(d.name, delta)
		}()
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

		err = d.keymap.Delete(keys)
		if err != nil {
			d.panic.Panic(fmt.Errorf("failed to delete keys: %v", err))
			return
		}

		d.immutableSegmentSize -= seg.Size()
		d.keyCount.Add(-1 * int64(len(keys)))

		// Deletion of segment files will happen when the segment is released by all reservation holders.
		seg.Release()
		d.segmentLock.Lock()
		delete(d.segments, index)
		d.segmentLock.Unlock()

		d.lowestSegmentIndex++
	}

	d.updateCurrentSize()
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

	d.immutableSegmentSize += d.segments[d.highestSegmentIndex].Size()

	// Seal the previous segment.
	flushLoopResponseChan := make(chan struct{}, 1)
	request := &flushLoopSealRequest{
		now:          now,
		responseChan: flushLoopResponseChan,
	}
	err := util.Send(d.panic, d.flushChannel, request)
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
		false,
		d.fsync)
	if err != nil {
		return err
	}
	d.segments[d.highestSegmentIndex].SetNextSegment(newSegment)
	d.highestSegmentIndex++

	d.segmentLock.Lock()
	d.segments[d.highestSegmentIndex] = newSegment
	d.segmentLock.Unlock()

	d.updateCurrentSize()

	return nil
}

func (d *DiskTable) ScheduleImmediateGC() error {
	request := &controlLoopGCRequest{
		completionChan: make(chan struct{}, 1),
	}

	err := util.Send(d.panic, d.controllerChannel, request)
	if err != nil {
		return fmt.Errorf("failed to send GC request: %v", err)
	}

	_, err = util.Await(d.panic, request.completionChan)
	if err != nil {
		return fmt.Errorf("failed to await GC completion: %v", err)
	}

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
		d.panic.Panic(fmt.Errorf("failed to seal segment %d: %v", d.highestSegmentIndex, err))
		return
	}

	// Flush the keys that are now durable in the segment.
	err = d.writeKeysToKeymap(durableKeys)
	if err != nil {
		d.panic.Panic(fmt.Errorf("failed to flush keys: %v", err))
		return
	}

	req.responseChan <- struct{}{}
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
		return
	}

	// The flush loop is responsible for the remaining parts of the flush.
	request := &flushLoopFlushRequest{
		flushWaitFunction: flushWaitFunction,
		responseChan:      req.responseChan,
	}
	err = util.Send(d.panic, d.flushChannel, request)
	if err != nil {
		d.logger.Errorf("failed to send flush request to flush loop: %v", err)
	}
}

// handleFlushLoopFlushRequest handles the part of the flush that is performed on the flush loop.
func (d *DiskTable) handleFlushLoopFlushRequest(req *flushLoopFlushRequest) {

	var segmentFlushStart time.Time
	if d.metrics != nil {
		segmentFlushStart = d.timeSource()
	}

	durableKeys, err := req.flushWaitFunction()
	if err != nil {
		d.panic.Panic(fmt.Errorf("failed to flush mutable segment: %v", err))
		return
	}

	if d.metrics != nil {
		segmentFlushEnd := d.timeSource()
		delta := segmentFlushEnd.Sub(segmentFlushStart)
		d.metrics.ReportSegmentFlushLatency(d.name, delta)
	}

	err = d.writeKeysToKeymap(durableKeys)
	if err != nil {
		d.panic.Panic(fmt.Errorf("failed to flush keys: %v", err))
		return
	}

	req.responseChan <- struct{}{}
}

// writeKeysToKeymap flushes all keys to the keymap. Once they are flushed, it also removes the keys from the
// unflushedDataCache.
func (d *DiskTable) writeKeysToKeymap(keys []*types.KAPair) error {
	if len(keys) == 0 {
		// Nothing to flush.
		return nil
	}

	if d.metrics != nil {
		start := d.timeSource()
		defer func() {
			end := d.timeSource()
			delta := end.Sub(start)
			d.metrics.ReportKeymapFlushLatency(d.name, delta)
		}()
	}

	err := d.keymap.Put(keys)
	if err != nil {
		return fmt.Errorf("failed to flush keys: %v", err)
	}

	// Keys are now durably written to both the segment and the keymap. It is therefore safe to remove them from the
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
		return
	}

	// This seals the current mutable segment and creates a new one. The new segment will have the new sharding factor.
	err = d.expandSegments()
	if err != nil {
		d.panic.Panic(fmt.Errorf("failed to expand segments: %v", err))
		return
	}
}

// controlLoopFlushRequest is a request to flush the writer that is sent to the control loop.
type controlLoopFlushRequest struct {
	// responseChan produces a value when the flush is complete.
	responseChan chan struct{}
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

// controlLoopGCRequest is a request to run garbage collection that is sent to the control loop.
type controlLoopGCRequest struct {
	completionChan chan struct{}
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
				d.handleControlLoopWriteRequest(req)
			} else if req, ok := message.(*controlLoopFlushRequest); ok {
				d.handleControlLoopFlushRequest(req)
			} else if req, ok := message.(*controlLoopSetShardingFactorRequest); ok {
				d.handleControlLoopSetShardingFactorRequest(req)
			} else if req, ok := message.(*controlLoopShutdownRequest); ok {
				d.handleControlLoopShutdownRequest(req)
				return
			} else if req, ok := message.(*controlLoopGCRequest); ok {
				d.doGarbageCollection()
				req.completionChan <- struct{}{}
			} else {
				d.panic.Panic(fmt.Errorf("Unknown control message type %T", message))
				return
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

	// responseChan sends an object when the flush is complete.
	responseChan chan struct{}
}

// flushLoopSealRequest is a request to seal the mutable segment that is sent to the flush loop.
type flushLoopSealRequest struct {
	// the time when the segment is sealed
	now time.Time
	// responseChan sends an object when the seal is complete.
	responseChan chan struct{}
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
				return
			}
		}
	}
}
