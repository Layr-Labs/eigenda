package disktable

import (
	"encoding/binary"
	"fmt"
	"os"
	"path"
	"sync/atomic"
	"time"

	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

const tableMetadataSerializationVersion = 0

const tableMetadataFileExtension = ".metadata"
const tableMetadataFileName = "table" + tableMetadataFileExtension
const tableMetadataSwapFileExtension = ".mswap"
const tableMetadataSwapFileName = "table" + tableMetadataSwapFileExtension

const tableMetadataSize = 16

// tableMetadata contains table data that is preserved across restarts.
type tableMetadata struct {
	logger logging.Logger

	tableDirectory string

	// the table's TTL, accessed/modified by concurrent goroutines
	ttl atomic.Pointer[time.Duration]

	// the table's sharding factor, accessed/modified by concurrent goroutines
	shardingFactor atomic.Uint32
}

// newTableMetadata creates a new table metadata object.
func newTableMetadata(
	logger logging.Logger,
	tableDirectory string,
	ttl time.Duration,
	shardingFactor uint32) (*tableMetadata, error) {

	metadata := &tableMetadata{
		logger:         logger,
		tableDirectory: tableDirectory,
	}
	metadata.ttl.Store(&ttl)
	metadata.shardingFactor.Store(shardingFactor)

	err := metadata.deleteOrphanedSwapFile()
	if err != nil {
		return nil, fmt.Errorf("failed to delete orphaned swap file: %v", err)
	}

	err = metadata.write()
	if err != nil {
		return nil, fmt.Errorf("failed to write table metadata: %v", err)
	}

	return metadata, nil
}

// loadTableMetadata loads the table metadata from disk.
func loadTableMetadata(logger logging.Logger, tableDirectory string) (*tableMetadata, error) {
	mPath := metadataPath(tableDirectory)

	exists, err := util.Exists(mPath)
	if err != nil {
		return nil, fmt.Errorf("failed to check if table metadata file exists: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("table metadata file does not exist: %s", mPath)
	}

	data, err := os.ReadFile(mPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read table metadata file %s: %v", mPath, err)
	}

	metadata, err := deserialize(data)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize table metadata: %v", err)
	}
	metadata.logger = logger
	metadata.tableDirectory = tableDirectory

	err = metadata.deleteOrphanedSwapFile()
	if err != nil {
		return nil, fmt.Errorf("failed to delete orphaned swap file: %v", err)
	}

	return metadata, nil
}

// Size returns the size of the table metadata file in bytes.
func (t *tableMetadata) Size() uint64 {
	return tableMetadataSize
}

// GetTTL returns the time-to-live for the table.
func (t *tableMetadata) GetTTL() time.Duration {
	return *t.ttl.Load()
}

// SetTTL sets the time-to-live for the table.
func (t *tableMetadata) SetTTL(ttl time.Duration) error {
	t.ttl.Store(&ttl)
	err := t.write()
	if err != nil {
		return fmt.Errorf("failed to update table metadata: %v", err)
	}
	return nil
}

// GetShardingFactor returns the sharding factor for the table.
func (t *tableMetadata) GetShardingFactor() uint32 {
	return t.shardingFactor.Load()
}

// SetShardingFactor sets the sharding factor for the table.
func (t *tableMetadata) SetShardingFactor(shardingFactor uint32) error {
	t.shardingFactor.Store(shardingFactor)
	err := t.write()
	if err != nil {
		return fmt.Errorf("failed to update table metadata: %v", err)
	}
	return nil
}

// Store atomically stores the table metadata to disk.
func (t *tableMetadata) write() error {
	bytes := t.serialize()
	swapPath := t.swapPath()
	swapFile, err := os.Create(swapPath)
	if err != nil {
		return fmt.Errorf("failed to create swap file %s: %v", swapPath, err)
	}

	_, err = swapFile.Write(bytes)
	if err != nil {
		return fmt.Errorf("failed to write to swap file %s: %v", swapPath, err)
	}

	err = swapFile.Close()
	if err != nil {
		return fmt.Errorf("failed to close swap file %s: %v", swapPath, err)
	}

	mPath := metadataPath(t.tableDirectory)
	err = os.Rename(swapPath, mPath)
	if err != nil {
		return fmt.Errorf("failed to rename swap file %s to metadata file %s: %v", swapPath, mPath, err)
	}

	return nil
}

// serialize serializes the table metadata to a byte slice.
func (t *tableMetadata) serialize() []byte {
	// 4 bytes for version
	// 8 bytes for TTL
	// 4 bytes for sharding factor
	data := make([]byte, tableMetadataSize)

	// Write the version
	binary.BigEndian.PutUint32(data[0:4], tableMetadataSerializationVersion)

	// Write the TTL
	ttlNanoseconds := t.GetTTL().Nanoseconds()
	binary.BigEndian.PutUint64(data[4:12], uint64(ttlNanoseconds))

	// Write the sharding factor
	binary.BigEndian.PutUint32(data[12:16], t.GetShardingFactor())

	return data
}

// deserialize deserializes the table metadata from a byte slice.
func deserialize(data []byte) (*tableMetadata, error) {
	// 4 bytes for version
	// 8 bytes for TTL
	// 4 bytes for sharding factor
	if len(data) != tableMetadataSize {
		return nil, fmt.Errorf("metadata file is not the correct size, expected 16 bytes, got %d", len(data))
	}

	serializationVersion := binary.BigEndian.Uint32(data[0:4])
	if serializationVersion != tableMetadataSerializationVersion {
		return nil, fmt.Errorf("unsupported serialization version: %d", serializationVersion)
	}

	ttl := time.Duration(binary.BigEndian.Uint64(data[4:12]))
	shardingFactor := binary.BigEndian.Uint32(data[12:16])

	metadata := &tableMetadata{}
	metadata.ttl.Store(&ttl)
	metadata.shardingFactor.Store(shardingFactor)

	return metadata, nil
}

// delete deletes the table metadata from disk.
func (t *tableMetadata) delete() error {
	metadataPath := path.Join(t.tableDirectory, tableMetadataFileName)
	err := os.Remove(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to delete table metadata file %s: %v", metadataPath, err)
	}
	return nil
}

// path returns the path to the table metadata file.
func metadataPath(tableDirectory string) string {
	return path.Join(tableDirectory, tableMetadataFileName)
}

// swapPath returns the path to the table metadata swap file.
func (t *tableMetadata) swapPath() string {
	return path.Join(t.tableDirectory, tableMetadataSwapFileName)
}

// deleteOrphanedSwapFile deletes the orphaned swap file if it exists.
// This can happen if the process crashes while writing the metadata file (recoverable).
func (t *tableMetadata) deleteOrphanedSwapFile() error {
	swapPath := t.swapPath()

	exists, err := util.Exists(swapPath)
	if err != nil {
		return fmt.Errorf("failed to check if swap file exists: %v", err)
	}

	if exists {
		t.logger.Warnf("Found orphaned table metadata swap file %s, deleting", swapPath)

		// delete orphaned swap file
		err = os.Remove(swapPath)
		if err != nil {
			return fmt.Errorf("failed to remove file %s: %v", swapPath, err)
		}
	}

	return nil
}
