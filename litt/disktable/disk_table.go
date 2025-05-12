package disktable

import (
	"errors"
	"fmt"
	"math"
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

	// fatalErrorHandler is a struct that permits the DB to "panic". There are many goroutines that function under the
	// hood, and many of these threads could, in theory, encounter errors which are unrecoverable. In such situations,
	// the desirable outcome is for the DB to report the error and then refuse to do additional work. If the DB is in a
	// broken state, it is much better to refuse to do work than to continue to do work and potentially corrupt data.
	fatalErrorHandler *util.FatalErrorHandler

	// The root directories for the disk table.
	roots []string

	// The directories where segment files are stored.
	segmentDirectories []string

	// The table's name.
	name string

	// The table's metadata.
	metadata *tableMetadata

	// A map of keys to their addresses.
	keymap keymap.Keymap

	// The path to the keymap directory.
	keymapPath string

	// The type file for the keymap.
	keymapTypeFile *keymap.KeymapTypeFile

	// unflushedDataCache is a map of keys to their values that may not have been flushed to disk yet. This is used as a
	// lookup table when data is requested from the table before it has been flushed to disk.
	unflushedDataCache sync.Map

	// clock is the time source used by the disk table.
	clock func() time.Time

	// The number of bytes contained within all segments, including the mutable segment. This tracks the number of
	// bytes that are on disk, not bytes in memory.
	size atomic.Uint64

	// The number of keys in the table.
	keyCount atomic.Int64

	// The control loop is a goroutine responsible for scheduling operations that mutate the table.
	controlLoop *controlLoop

	// The flush loop is a goroutine responsible for blocking on flush operations.
	flushLoop *flushLoop

	// Encapsulates metrics for the database.
	metrics *metrics.LittDBMetrics
}

// NewDiskTable creates a new DiskTable.
func NewDiskTable(
	config *litt.Config,
	name string,
	keymap keymap.Keymap,
	keymapPath string,
	keymapTypeFile *keymap.KeymapTypeFile,
	roots []string,
	reloadKeymap bool,
	metrics *metrics.LittDBMetrics) (litt.ManagedTable, error) {

	if config.GCPeriod <= 0 {
		return nil, errors.New("garbage collection period must be greater than 0")
	}

	// If the root directories don't exist, create them.
	for _, root := range roots {
		exists, err := util.Exists(root)
		if err != nil {
			return nil, fmt.Errorf("failed to check if root directory exists: %w", err)
		}
		if !exists {
			err = os.MkdirAll(root, 0755)
			if err != nil {
				return nil, fmt.Errorf("failed to create root directory: %w", err)
			}
		}
	}

	// For each root directory, create a segment directory if it doesn't exist.
	segDirs := make([]string, 0, len(roots))
	for _, root := range roots {
		segDir := path.Join(root, segmentDirectory)
		segDirs = append(segDirs, segDir)

		exists, err := util.Exists(segDir)
		if err != nil {
			return nil, fmt.Errorf("failed to check if segment directory exists: %w", err)
		}

		if !exists {
			err := os.MkdirAll(segDir, 0755)
			if err != nil {
				return nil, fmt.Errorf("failed to create segment directory: %w", err)
			}
		}
	}

	var metadataFilePath string
	var metadata *tableMetadata

	// Find the table metadata file or create a new one.
	for _, root := range roots {
		possibleMetadataPath := metadataPath(root)
		exists, err := util.Exists(possibleMetadataPath)
		if err != nil {
			return nil, fmt.Errorf("failed to check if metadata file exists: %w", err)
		}
		if exists {
			if metadataFilePath != "" {
				return nil, fmt.Errorf("multiple metadata files found: %s and %s",
					metadataFilePath, possibleMetadataPath)
			}

			// We've found an existing metadata file. Use it.
			metadataFilePath = possibleMetadataPath
		}
	}
	if metadataFilePath == "" {
		// No metadata file exists yet. Create a new one in the first root.
		var err error
		metadataDir := roots[0]
		metadata, err = newTableMetadata(config.Logger, metadataDir, config.TTL, config.ShardingFactor)
		if err != nil {
			return nil, fmt.Errorf("failed to create table metadata: %w", err)
		}
	} else {
		// Metadata file exists, so we need to load it.
		var err error
		metadataDir := path.Dir(metadataFilePath)
		metadata, err = loadTableMetadata(config.Logger, metadataDir)
		if err != nil {
			return nil, fmt.Errorf("failed to load table metadata: %w", err)
		}
	}

	fatalErrorHandler := util.NewFatalErrorHandler(config.CTX, config.Logger, config.FatalErrorCallback)

	table := &DiskTable{
		logger:             config.Logger,
		fatalErrorHandler:  fatalErrorHandler,
		clock:              config.Clock,
		roots:              roots,
		segmentDirectories: segDirs,
		name:               name,
		metadata:           metadata,
		keymap:             keymap,
		keymapPath:         keymapPath,
		keymapTypeFile:     keymapTypeFile,
		metrics:            metrics,
	}

	lowestSegmentIndex, highestSegmentIndex, segments, err :=
		segment.GatherSegmentFiles(
			config.Logger,
			fatalErrorHandler,
			table.segmentDirectories,
			config.Clock())
	if err != nil {
		return nil, fmt.Errorf("failed to gather segment files: %w", err)
	}

	immutableSegmentSize := uint64(0)
	for _, seg := range segments {
		immutableSegmentSize += seg.Size()
	}

	// Create the mutable segment
	creatingFirstSegment := len(segments) == 0

	var nextSegmentIndex uint32
	if creatingFirstSegment {
		nextSegmentIndex = 0
	} else {
		nextSegmentIndex = highestSegmentIndex + 1
	}
	salt := [16]byte{}
	_, err = config.SaltShaker.Read(salt[:])
	if err != nil {
		return nil, fmt.Errorf("failed to read salt: %w", err)
	}
	mutableSegment, err := segment.CreateSegment(
		config.Logger,
		fatalErrorHandler,
		nextSegmentIndex,
		segDirs,
		metadata.GetShardingFactor(),
		salt,
		config.Fsync)
	if err != nil {
		return nil, fmt.Errorf("failed to create mutable segment: %w", err)
	}
	if !creatingFirstSegment {
		segments[highestSegmentIndex].SetNextSegment(mutableSegment)
		highestSegmentIndex++
	}
	segments[nextSegmentIndex] = mutableSegment

	if reloadKeymap {
		config.Logger.Infof("reloading keymap from segments")
		err = table.reloadKeymap(segments, lowestSegmentIndex, highestSegmentIndex)
		if err != nil {
			return nil, fmt.Errorf("failed to load keymap from segments: %w", err)
		}
	}

	tableSaltShaker := rand.New(rand.NewSource(config.SaltShaker.Int63()))

	// Start the flush loop.
	fLoop := &flushLoop{
		logger:            config.Logger,
		diskTable:         table,
		fatalErrorHandler: fatalErrorHandler,
		flushChannel:      make(chan any, tableFlushChannelCapacity),
		metrics:           metrics,
		clock:             config.Clock,
		name:              name,
	}
	table.flushLoop = fLoop
	go fLoop.run()

	// Start the control loop.
	cLoop := &controlLoop{
		logger:                  config.Logger,
		diskTable:               table,
		fatalErrorHandler:       fatalErrorHandler,
		controllerChannel:       make(chan any, config.ControlChannelSize),
		lowestSegmentIndex:      lowestSegmentIndex,
		highestSegmentIndex:     highestSegmentIndex,
		segments:                segments,
		size:                    &table.size,
		keyCount:                &table.keyCount,
		targetFileSize:          config.TargetSegmentFileSize,
		targetKeyFileSize:       config.TargetSegmentKeyFileSize,
		maxKeyCount:             config.MaxSegmentKeyCount,
		clock:                   config.Clock,
		segmentDirectories:      segDirs,
		saltShaker:              tableSaltShaker,
		metadata:                metadata,
		fsync:                   config.Fsync,
		metrics:                 metrics,
		name:                    name,
		gcBatchSize:             config.GCBatchSize,
		keymap:                  keymap,
		flushLoop:               fLoop,
		garbageCollectionPeriod: config.GCPeriod,
	}
	cLoop.threadsafeHighestSegmentIndex.Store(highestSegmentIndex)
	table.controlLoop = cLoop
	cLoop.updateCurrentSize()
	go cLoop.run()

	return table, nil
}

func (d *DiskTable) KeyCount() uint64 {
	return uint64(d.keyCount.Load())
}

func (d *DiskTable) Size() uint64 {
	return d.size.Load()
}

// reloadKeymap reloads the keymap from the segments. This is necessary when the keymap is lost, the keymap doesn't
// save its data on disk, or we are migrating from one keymap type to another.
func (d *DiskTable) reloadKeymap(
	segments map[uint32]*segment.Segment,
	lowestSegmentIndex uint32,
	highestSegmentIndex uint32) error {

	start := d.clock()
	defer func() {
		d.logger.Infof("spent %v reloading keymap", d.clock().Sub(start))
	}()

	// It's possible that some of the data written near the end of the previous session was corrupted.
	// Read data from the end until the first valid key/value pair is found.
	isValid := false

	batch := make([]*types.KAPair, 0, keymapReloadBatchSize)

	for i := highestSegmentIndex; i >= lowestSegmentIndex && i+1 != 0; i-- {
		if !segments[i].IsSealed() {
			// ignore unsealed segment, this will have been created in the current session and will not
			// yet contain any data.
			continue
		}

		keys, err := segments[i].GetKeys()
		if err != nil {
			return fmt.Errorf("failed to get keys from segment %d: %w", i, err)
		}
		for keyIndex := len(keys) - 1; keyIndex >= 0; keyIndex-- {
			key := keys[keyIndex]

			if !isValid {
				_, err = segments[i].Read(key.Key, key.Address)
				if err == nil {
					// we found a valid key/value pair. All subsequent keys are valid.
					isValid = true
				} else {
					// This is not cause for alarm (probably).
					// This can happen when the database is not cleanly shut down,
					// and just means that some data near the end was not fully committed.
					d.logger.Infof("truncated value for key %s with address %s for segment %d",
						key.Key, key.Address, i)
				}
			}

			if isValid {
				batch = append(batch, key)
				if len(batch) == keymapReloadBatchSize {
					err = d.keymap.Put(batch)
					if err != nil {
						return fmt.Errorf("failed to put keys for segment %d: %w", i, err)
					}
					batch = make([]*types.KAPair, 0, keymapReloadBatchSize)
				}
			}
		}

	}

	if len(batch) > 0 {
		err := d.keymap.Put(batch)
		if err != nil {
			return fmt.Errorf("failed to put keys: %w", err)
		}
	}

	// Now that the keymap is loaded, write the marker file that indicates that the keymap is fully loaded.
	// If we crash prior to writing this file, the keymap will reload from the segments again.
	keymapInitializedFile := path.Join(d.keymapPath, keymap.KeymapInitializedFileName)
	err := os.MkdirAll(d.keymapPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create keymap directory: %w", err)
	}

	f, err := os.Create(keymapInitializedFile)
	if err != nil {
		return fmt.Errorf("failed to create keymap initialized file after reload: %w", err)
	}
	err = f.Close()
	if err != nil {
		return fmt.Errorf("failed to close keymap initialized file after reload: %w", err)
	}

	return nil
}

func (d *DiskTable) Name() string {
	return d.name
}

// Close stops the disk table. Flushes all data out to disk.
func (d *DiskTable) Close() error {
	if ok, err := d.fatalErrorHandler.IsOk(); !ok {
		return fmt.Errorf("cannot process Stop() request, DB is in panicked state due to error: %w", err)
	}

	d.fatalErrorHandler.Shutdown()

	shutdownCompleteChan := make(chan struct{}, 1)
	request := &controlLoopShutdownRequest{
		shutdownCompleteChan: shutdownCompleteChan,
	}

	err := d.controlLoop.enqueue(request)
	if err != nil {
		return fmt.Errorf("failed to send shutdown request: %w", err)
	}

	_, err = util.AwaitIfNotFatal(d.fatalErrorHandler, shutdownCompleteChan)
	if err != nil {
		return fmt.Errorf("failed to shutdown: %w", err)
	}

	return nil
}

// Destroy stops the disk table and delete all files.
func (d *DiskTable) Destroy() error {
	if ok, err := d.fatalErrorHandler.IsOk(); !ok {
		return fmt.Errorf("Cannot process Destroy() request, DB is in panicked state due to error: %w", err)
	}

	err := d.Close()
	if err != nil {
		return fmt.Errorf("failed to stop: %w", err)
	}

	d.logger.Infof("deleting disk table at path(s): %v", d.roots)

	// release all segments
	segments, err := d.controlLoop.getSegments()
	if err != nil {
		return fmt.Errorf("failed to get segments: %w", err)
	}
	for _, seg := range segments {
		seg.Release()
	}
	// wait for segments to delete themselves
	for _, seg := range segments {
		err = seg.BlockUntilFullyDeleted()
		if err != nil {
			return fmt.Errorf("failed to delete segment: %w", err)
		}
	}

	// delete all segment directories
	for _, segDir := range d.segmentDirectories {
		err = os.Remove(segDir)
		if err != nil {
			return fmt.Errorf("failed to remove segment directory: %w", err)
		}
	}

	// destroy the keymap
	err = d.keymap.Destroy()
	if err != nil {
		return fmt.Errorf("failed to destroy keymap: %w", err)
	}
	err = d.keymapTypeFile.Delete()
	if err != nil {
		return fmt.Errorf("failed to delete keymap type file: %w", err)
	}
	exists, err := util.Exists(d.keymapPath)
	if err != nil {
		return fmt.Errorf("failed to check if keymap directory exists: %w", err)
	}
	if exists {
		err = os.RemoveAll(d.keymapPath)
		if err != nil {
			return fmt.Errorf("failed to remove keymap directory: %w", err)
		}
	}

	// delete the metadata file
	err = d.metadata.delete()
	if err != nil {
		return fmt.Errorf("failed to delete metadata: %w", err)
	}

	// delete the root directories for the table
	for _, root := range d.roots {
		err = os.Remove(root)
		if err != nil {
			return fmt.Errorf("failed to remove root directory: %w", err)
		}
	}

	return nil
}

// SetTTL sets the TTL for the disk table. If set to 0, no TTL is enforced. This setting affects both new
// data and data already written.
func (d *DiskTable) SetTTL(ttl time.Duration) error {
	if ok, err := d.fatalErrorHandler.IsOk(); !ok {
		return fmt.Errorf("Cannot process SetTTL() request, DB is in panicked state due to error: %w", err)
	}

	err := d.metadata.SetTTL(ttl)
	if err != nil {
		return fmt.Errorf("failed to set TTL: %w", err)
	}
	return nil
}

func (d *DiskTable) SetShardingFactor(shardingFactor uint32) error {
	if ok, err := d.fatalErrorHandler.IsOk(); !ok {
		return fmt.Errorf(
			"Cannot process SetShardingFactor() request, DB is in panicked state due to error: %w", err)
	}

	if shardingFactor == 0 {
		return fmt.Errorf("sharding factor must be greater than 0")
	}

	request := &controlLoopSetShardingFactorRequest{
		shardingFactor: shardingFactor,
	}
	err := d.controlLoop.enqueue(request)
	if err != nil {
		return fmt.Errorf("failed to send sharding factor request: %w", err)
	}

	return nil
}

func (d *DiskTable) Get(key []byte) (value []byte, exists bool, err error) {
	if ok, err := d.fatalErrorHandler.IsOk(); !ok {
		return nil, false, fmt.Errorf(
			"Cannot process Get() request, DB is in panicked state due to error: %w", err)
	}

	var cacheHit bool
	var dataSize uint64
	if d.metrics != nil {
		start := d.clock()
		defer func() {
			end := d.clock()
			delta := end.Sub(start)
			d.metrics.ReportReadOperation(d.name, delta, dataSize, cacheHit)
		}()
	}

	// First, check if the key is in the unflushed data map.
	// If so, return it from there.
	if value, ok := d.unflushedDataCache.Load(util.UnsafeBytesToString(key)); ok {
		bytes := value.([]byte)
		cacheHit = true
		dataSize = uint64(len(bytes))
		return bytes, true, nil
	}

	// Look up the address of the data.
	address, ok, err := d.keymap.Get(key)
	if err != nil {
		return nil, false, fmt.Errorf("failed to get address: %w", err)
	}
	if !ok {
		return nil, false, nil
	}

	// Reserve the segment that contains the data.
	seg, ok := d.controlLoop.getReservedSegment(address.Index())
	if !ok {
		return nil, false, nil
	}
	defer seg.Release()

	// Read the data from disk.
	data, err := seg.Read(key, address)
	if err != nil {
		return nil, false, fmt.Errorf("failed to read data: %w", err)
	}

	dataSize = uint64(len(data))

	return data, true, nil
}

func (d *DiskTable) CacheAwareGet(
	key []byte,
	onlyReadFromCache bool,
) (value []byte, exists bool, hot bool, err error) {

	if ok, err := d.fatalErrorHandler.IsOk(); !ok {
		return nil, false, false, fmt.Errorf(
			"Cannot process CacheAwareGet() request, DB is in panicked state due to error: %w", err)
	}

	var cacheHit bool
	var dataSize uint64
	if d.metrics != nil {
		start := d.clock()
		defer func() {
			end := d.clock()
			delta := end.Sub(start)
			d.metrics.ReportReadOperation(d.name, delta, dataSize, cacheHit)
		}()
	}

	// First, check if the key is in the unflushed data map. If so, return it from there.
	// Performance wise, this has equivalent semantics to reading the value from
	// a cache, so we'd might as well count it as a cache hit.
	var rawValue any
	if rawValue, exists = d.unflushedDataCache.Load(util.UnsafeBytesToString(key)); exists {
		value = rawValue.([]byte)
		cacheHit = true
		dataSize = uint64(len(value))
		return value, true, true, nil
	}

	// Look up the address of the data.
	var address types.Address
	address, exists, err = d.keymap.Get(key)
	if err != nil {
		return nil, false, false, fmt.Errorf("failed to get address: %w", err)
	}
	if !exists {
		return nil, false, false, nil
	}

	if onlyReadFromCache {
		// The value exists but we are not allowed to read it from disk.
		return nil, true, false, nil
	}

	// Reserve the segment that contains the data.
	seg, ok := d.controlLoop.getReservedSegment(address.Index())
	if !ok {
		// This can happen if there is a race between this thread and the GC thread, i.e.
		// if we start reading a value just as the garbage collector decides to delete it.
		return nil, false, false, nil
	}
	defer seg.Release()

	// Read the data from disk.
	value, err = seg.Read(key, address)
	if err != nil {
		return nil, false, false, fmt.Errorf("failed to read data: %w", err)
	}

	dataSize = uint64(len(value))

	return value, true, false, nil
}

func (d *DiskTable) Put(key []byte, value []byte) error {
	return d.PutBatch([]*types.KVPair{{Key: key, Value: value}})
}

func (d *DiskTable) PutBatch(batch []*types.KVPair) error {
	if ok, err := d.fatalErrorHandler.IsOk(); !ok {
		return fmt.Errorf("Cannot process PutBatch() request, DB is in panicked state due to error: %w", err)
	}

	if d.metrics != nil {
		start := d.clock()
		totalSize := uint64(0)
		for _, kv := range batch {
			totalSize += uint64(len(kv.Value))
		}
		defer func() {
			end := d.clock()
			delta := end.Sub(start)
			d.metrics.ReportWriteOperation(d.name, delta, uint64(len(batch)), totalSize)
		}()
	}

	for _, kv := range batch {
		if len(kv.Key) > math.MaxUint32 {
			return fmt.Errorf("key is too large, length must not exceed 2^32 bytes: %d bytes", len(kv.Key))
		}
		if len(kv.Value) > math.MaxUint32 {
			return fmt.Errorf("value is too large, length must not exceed 2^32 bytes: %d bytes", len(kv.Value))
		}
		if kv.Key == nil {
			return fmt.Errorf("nil keys are not supported")
		}
		if kv.Value == nil {
			return fmt.Errorf("nil values are not supported")
		}

		d.unflushedDataCache.Store(util.UnsafeBytesToString(kv.Key), kv.Value)
	}

	request := &controlLoopWriteRequest{
		values: batch,
	}
	err := d.controlLoop.enqueue(request)
	if err != nil {
		return fmt.Errorf("failed to send write request: %w", err)
	}

	d.keyCount.Add(int64(len(batch)))

	return nil
}

func (d *DiskTable) Exists(key []byte) (bool, error) {
	_, ok := d.unflushedDataCache.Load(util.UnsafeBytesToString(key))
	if ok {
		return true, nil
	}

	_, ok, err := d.keymap.Get(key)
	if err != nil {
		return false, fmt.Errorf("failed to get address: %w", err)
	}

	return ok, nil
}

// Flush flushes all data to disk. Blocks until all data previously submitted to Put has been written to disk.
func (d *DiskTable) Flush() error {
	if ok, err := d.fatalErrorHandler.IsOk(); !ok {
		return fmt.Errorf("Cannot process Flush() request, DB is in panicked state due to error: %w", err)
	}

	if d.metrics != nil {
		start := d.clock()
		defer func() {
			end := d.clock()
			delta := end.Sub(start)
			d.metrics.ReportFlushOperation(d.name, delta)
		}()
	}

	flushReq := &controlLoopFlushRequest{
		responseChan: make(chan struct{}, 1),
	}
	err := d.controlLoop.enqueue(flushReq)
	if err != nil {
		return fmt.Errorf("failed to send flush request: %w", err)
	}

	_, err = util.AwaitIfNotFatal(d.fatalErrorHandler, flushReq.responseChan)
	if err != nil {
		return fmt.Errorf("failed to flush: %w", err)
	}

	return nil
}

func (d *DiskTable) SetWriteCacheSize(size uint64) error {
	if ok, err := d.fatalErrorHandler.IsOk(); !ok {
		return fmt.Errorf(
			"Cannot process SetWriteCacheSize() request, DB is in panicked state due to error: %w", err)
	}

	// this implementation does not provide a cache, if a cache is needed then it must be provided by a wrapper
	return nil
}

func (d *DiskTable) SetReadCacheSize(size uint64) error {
	if ok, err := d.fatalErrorHandler.IsOk(); !ok {
		return fmt.Errorf(
			"Cannot process SetReadCacheSize() request, DB is in panicked state due to error: %w", err)
	}

	// this implementation does not provide a cache, if a cache is needed then it must be provided by a wrapper
	return nil
}

func (d *DiskTable) RunGC() error {
	if ok, err := d.fatalErrorHandler.IsOk(); !ok {
		return fmt.Errorf(
			"Cannot process RunGC() request, DB is in panicked state due to error: %w", err)
	}

	request := &controlLoopGCRequest{
		completionChan: make(chan struct{}, 1),
	}

	err := d.controlLoop.enqueue(request)
	if err != nil {
		return fmt.Errorf("failed to send GC request: %w", err)
	}

	_, err = util.AwaitIfNotFatal(d.fatalErrorHandler, request.completionChan)
	if err != nil {
		return fmt.Errorf("failed to await GC completion: %w", err)
	}

	return nil
}

// writeKeysToKeymap flushes all keys to the keymap. Once they are flushed, it also removes the keys from the
// unflushedDataCache.
func (d *DiskTable) writeKeysToKeymap(keys []*types.KAPair) error {
	if len(keys) == 0 {
		// Nothing to flush.
		return nil
	}

	if d.metrics != nil {
		start := d.clock()
		defer func() {
			end := d.clock()
			delta := end.Sub(start)
			d.metrics.ReportKeymapFlushLatency(d.name, delta)
		}()
	}

	err := d.keymap.Put(keys)
	if err != nil {
		return fmt.Errorf("failed to flush keys: %w", err)
	}

	// Keys are now durably written to both the segment and the keymap. It is therefore safe to remove them from the
	// unflushed data cache.
	for _, ka := range keys {
		d.unflushedDataCache.Delete(util.UnsafeBytesToString(ka.Key))
	}

	return nil
}
