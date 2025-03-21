package segment

import (
	"encoding/binary"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/Layr-Labs/eigenda/litt/util"
)

const (
	// The current serialization version. If we ever change how we serialize data, bump this version.
	currentSerializationVersion = uint32(0)

	// MetadataFileExtension is the file extension for the metadata file. This file contains metadata about the data
	// segment, such as serialization version and expiration time.
	MetadataFileExtension = ".metadata"

	// MetadataSwapExtension is the file extension for the metadata swap file. This file is used to atomically update
	// the metadata file by doing an atomic rename of the swap file to the metadata file. If this file is ever
	// present when the database first starts, it is an artifact of a crash during a metadata update, and should be
	// deleted.
	MetadataSwapExtension = ".metadata.swap"

	// The size of the metadata file in bytes. This is a constant, so it's convenient to have it here.
	// - 4 bytes for version
	// - 4 bytes for the sharding factor
	// - 4 bytes for salt
	// - 8 bytes for timestamp
	// - and 1 byte for sealed.
	metadataSize = 21
)

// metadataFile contains metadata about a segment.
type metadataFile struct {
	// The segment index. This value is encoded in the file name.
	index uint32

	// The serialization version for this segment, used to permit smooth data migrations.
	// This value is encoded in the file.
	serializationVersion uint32

	// The sharding factor for this segment. This value is encoded in the file.
	shardingFactor uint32

	// A random number, used to make the sharding hash function hard for an attacker to predict.
	// This value is encoded in the file.
	salt uint32

	// The time when the last value was written into the segment, in nanoseconds since the epoch. A segment can
	// only be deleted when all values within it are expired, and so we only need to keep track of the timestamp of
	// the last value (which always expires last). This value is irrelevant if the segment is not yet sealed.
	// This value is encoded in the file.
	timestamp uint64

	// If true, the segment is sealed and no more data can be written to it. If false, then data can still be written to
	// this segment. This value is encoded in the file.
	sealed bool

	// The parent directory containing this file. This value is not encoded in file, and is stored here
	// for bookkeeping purposes.
	parentDirectory string
}

// newMetadataFile creates a new metadata file. When this method returns, the metadata file will
// be durably written to disk.
//
// Note that shardingFactor and salt parameters are ignored if this is not a new metadata file. Metadata files
// loaded from disk always use their original sharding factor and salt values.
func newMetadataFile(
	index uint32,
	shardingFactor uint32,
	salt uint32,
	parentDirectory string) (*metadataFile, error) {

	file := &metadataFile{
		index:           index,
		parentDirectory: parentDirectory,
	}

	filePath := file.path()
	exists, _, err := util.VerifyFilePermissions(filePath)
	if err != nil {
		return nil, fmt.Errorf("file %s has incorrect permissions: %v", filePath, err)
	}

	if exists {
		// File exists. Load it.
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read metadata file %s: %v", filePath, err)
		}
		err = file.deserialize(data)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize metadata file %s: %v", filePath, err)
		}
	} else {
		// File does not exist. Create it.
		file.serializationVersion = currentSerializationVersion
		file.shardingFactor = shardingFactor
		file.salt = salt
		err = file.write()
		if err != nil {
			return nil, fmt.Errorf("failed to write metadata file: %v", err)
		}
	}

	return file, nil
}

// Name returns the file name for this metadata file.
func (m *metadataFile) name() string {
	return fmt.Sprintf("%d%s", m.index, MetadataFileExtension)
}

// Path returns the full path to this metadata file.
func (m *metadataFile) path() string {
	return path.Join(m.parentDirectory, m.name())
}

// SwapName returns the file name for the swap file for this metadata file.
func (m *metadataFile) swapName() string {
	return fmt.Sprintf("%d%s", m.index, MetadataSwapExtension)
}

// SwapPath returns the full path to the swap file for this metadata file.
func (m *metadataFile) swapPath() string {
	return path.Join(m.parentDirectory, m.swapName())
}

// Seal seals the segment. This action will atomically write the metadata file to disk one final time,
// and should only be performed when all data that will be written to the key/value files has been made durable.
func (m *metadataFile) seal(now time.Time) error {
	m.sealed = true
	m.timestamp = uint64(now.UnixNano())
	err := m.write()
	if err != nil {
		return fmt.Errorf("failed to write sealed metadata file: %v", err)
	}
	return nil
}

// serialize serializes the metadata file to a byte array.
func (m *metadataFile) serialize() []byte {
	// 4 bytes for version, 8 bytes for timestamp, 1 byte for sealed
	data := make([]byte, metadataSize)

	// Write the version
	binary.BigEndian.PutUint32(data[0:4], m.serializationVersion)

	// Write the sharding factor
	binary.BigEndian.PutUint32(data[4:8], m.shardingFactor)

	// Write the salt
	binary.BigEndian.PutUint32(data[8:12], m.salt)

	// Write the timestamp
	binary.BigEndian.PutUint64(data[12:20], m.timestamp)

	// Write the sealed flag
	if m.sealed {
		data[20] = 1
	} else {
		data[20] = 0
	}

	return data
}

// deserialize deserializes the metadata file from a byte array.
func (m *metadataFile) deserialize(data []byte) error {
	if len(data) != metadataSize {
		return fmt.Errorf("metadata file is not the correct size: %d", len(data))
	}

	m.serializationVersion = binary.BigEndian.Uint32(data[0:4])
	if m.serializationVersion != currentSerializationVersion {
		return fmt.Errorf("unsupported serialization version: %d", m.serializationVersion)
	}

	m.shardingFactor = binary.BigEndian.Uint32(data[4:8])
	m.salt = binary.BigEndian.Uint32(data[8:12])
	m.timestamp = binary.BigEndian.Uint64(data[12:20])
	m.sealed = data[20] == 1

	return nil
}

// write atomically writes the metadata file to disk.
func (m *metadataFile) write() error {
	bytes := m.serialize()
	swapPath := m.swapPath()
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

	metadataPath := m.path()
	err = os.Rename(swapPath, metadataPath)
	if err != nil {
		return fmt.Errorf("failed to rename swap file %s to metadata file %s: %v", swapPath, metadataPath, err)
	}

	return nil
}

// delete deletes the metadata file from disk.
func (m *metadataFile) delete() error {
	err := os.Remove(m.path())
	if err != nil {
		return fmt.Errorf("failed to remove metadata file %s: %v", m.path(), err)
	}

	return nil
}
