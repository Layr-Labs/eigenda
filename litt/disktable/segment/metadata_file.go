package segment

import (
	"encoding/binary"
	"fmt"
	"os"
	"path"
	"strconv"
	"time"
)

const (

	// MetadataFileExtension is the file extension for the metadata file.
	MetadataFileExtension = ".metadata"

	// MetadataSwapExtension is the file extension for the metadata swap file. This file is used to atomically update
	// the metadata file by doing an atomic rename of the swap file to the metadata file. If this file is ever
	// present when the database first starts, it is an artifact of a crash during a metadata update, and should be
	// deleted.
	MetadataSwapExtension = ".metadata.swap"

	// OldMetadataSize is the size of the format 0 metadata file in bytes.
	// This is a constant, so it's convenient to have it here.
	// - 4 bytes for version
	// - 4 bytes for the sharding factor
	// - 4 bytes for salt
	// - 8 bytes for lastValueTimestamp
	// - and 1 byte for sealed.
	OldMetadataSize = 21

	// MetadataSize is the size of the metadata file in bytes. This is a constant, so it's convenient to have it here.
	// - 4 bytes for version
	// - 4 bytes for the sharding factor
	// - 16 bytes for salt
	// - 8 bytes for lastValueTimestamp
	// - and 1 byte for sealed.
	MetadataSize = 33
)

// metadataFile contains metadata about a segment. This file contains metadata about the data segment, such as
// serialization version and the lastValueTimestamp when the file was sealed.
type metadataFile struct {
	// The segment index. This value is encoded in the file name.
	index uint32

	// The serialization version for this segment, used to permit smooth data migrations.
	// This value is encoded in the file.
	segmentVersion SegmentVersion

	// The sharding factor for this segment. This value is encoded in the file.
	shardingFactor uint32

	// A random number, used to make the sharding hash function hard for an attacker to predict.
	// This value is encoded in the file. Note: after the hash function change, this value is
	// only used for data written with the old hash function.
	legacySalt uint32

	// A random byte array, used to make the sharding hash function hard for an attacker to predict.
	// This value is encoded in the file.
	salt [16]byte

	// The time when the last value was written into the segment, in nanoseconds since the epoch. A segment can
	// only be deleted when all values within it are expired, and so we only need to keep track of the lastValueTimestamp of
	// the last value (which always expires last). This value is irrelevant if the segment is not yet sealed.
	// This value is encoded in the file.
	lastValueTimestamp uint64

	// If true, the segment is sealed and no more data can be written to it. If false, then data can still be written to
	// this segment. This value is encoded in the file.
	sealed bool

	// The parent directory containing this file. This value is not encoded in file, and is stored here
	// for bookkeeping purposes.
	parentDirectory string
}

// createMetadataFile creates a new metadata file. When this method returns, the metadata file will
// be durably written to disk.
func createMetadataFile(
	index uint32,
	shardingFactor uint32,
	salt [16]byte,
	parentDirectory string) (*metadataFile, error) {

	file := &metadataFile{
		index:           index,
		parentDirectory: parentDirectory,
	}

	file.segmentVersion = LatestSegmentVersion
	file.shardingFactor = shardingFactor
	file.salt = salt
	err := file.write()
	if err != nil {
		return nil, fmt.Errorf("failed to write metadata file: %v", err)
	}

	return file, nil
}

// loadMetadataFile loads the metadata file from disk, looking in the given parent directories until it finds the file.
// If the file is not found, it returns an error.
func loadMetadataFile(index uint32, parentDirectories []string) (*metadataFile, error) {
	metadataFileName := fmt.Sprintf("%d%s", index, MetadataFileExtension)
	metadataPath, err := lookForFile(parentDirectories, metadataFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to find metadata file: %w", err)
	}
	if metadataPath == "" {
		return nil, fmt.Errorf("failed to find metadata file %s", metadataFileName)
	}
	parentDirectory := path.Dir(metadataPath)

	file := &metadataFile{
		index:           index,
		parentDirectory: parentDirectory,
	}

	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata file %s: %v", metadataPath, err)
	}
	err = file.deserialize(data)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize metadata file %s: %v", metadataPath, err)
	}

	return file, nil
}

// MetadataFileExtension is the file extension for the metadata file. Metadata file names have the form "X.metadata",
// where X is the segment index.
func getMetadataFileIndex(fileName string) (uint32, error) {
	indexString := path.Base(fileName)[:len(fileName)-len(MetadataFileExtension)]
	index, err := strconv.Atoi(indexString)
	if err != nil {
		return 0, fmt.Errorf("failed to parse index from file name %s: %v", fileName, err)
	}

	return uint32(index), nil
}

// Size returns the size of the metadata file in bytes.
func (m *metadataFile) Size() uint64 {
	if m.segmentVersion == OldHashFunctionSerializationVersion {
		return OldMetadataSize
	} else {
		return MetadataSize
	}
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
	m.lastValueTimestamp = uint64(now.UnixNano())
	err := m.write()
	if err != nil {
		return fmt.Errorf("failed to write sealed metadata file: %v", err)
	}
	return nil
}

func (m *metadataFile) serializeLegacy() []byte {
	data := make([]byte, OldMetadataSize)

	// Write the version
	binary.BigEndian.PutUint32(data[0:4], uint32(m.segmentVersion))

	// Write the sharding factor
	binary.BigEndian.PutUint32(data[4:8], m.shardingFactor)

	// Write the salt
	binary.BigEndian.PutUint32(data[8:12], m.legacySalt)

	// Write the lastValueTimestamp
	binary.BigEndian.PutUint64(data[12:20], m.lastValueTimestamp)

	// Write the sealed flag
	if m.sealed {
		data[20] = 1
	} else {
		data[20] = 0
	}

	return data
}

// serialize serializes the metadata file to a byte array.
func (m *metadataFile) serialize() []byte {
	if m.segmentVersion == OldHashFunctionSerializationVersion {
		return m.serializeLegacy()
	}

	data := make([]byte, MetadataSize)

	// Write the version
	binary.BigEndian.PutUint32(data[0:4], uint32(m.segmentVersion))

	// Write the sharding factor
	binary.BigEndian.PutUint32(data[4:8], m.shardingFactor)

	// Write the salt
	copy(data[8:24], m.salt[:])

	// Write the lastValueTimestamp
	binary.BigEndian.PutUint64(data[24:32], m.lastValueTimestamp)

	// Write the sealed flag
	if m.sealed {
		data[32] = 1
	} else {
		data[32] = 0
	}

	return data
}

// deserialize deserializes the metadata file from a byte array.
func (m *metadataFile) deserialize(data []byte) error {
	if len(data) < 4 {
		return fmt.Errorf("metadata file is not the correct size, expected at least 4 bytes, got %d", len(data))
	}

	m.segmentVersion = SegmentVersion(binary.BigEndian.Uint32(data[0:4]))
	if m.segmentVersion > LatestSegmentVersion {
		return fmt.Errorf("unsupported serialization version: %d", m.segmentVersion)
	}

	if m.segmentVersion == OldHashFunctionSerializationVersion {
		if len(data) != OldMetadataSize {
			return fmt.Errorf("metadata file is not the correct size, expected %d, got %d",
				OldMetadataSize, len(data))
		}

		// TODO (cody.littley): delete this after all data is migrated to the new hash function.
		m.shardingFactor = binary.BigEndian.Uint32(data[4:8])
		m.legacySalt = binary.BigEndian.Uint32(data[8:12])
		m.lastValueTimestamp = binary.BigEndian.Uint64(data[12:20])
		m.sealed = data[20] == 1
		return nil
	}

	if len(data) != MetadataSize {
		return fmt.Errorf("metadata file is not the correct size, expected %d, got %d",
			MetadataSize, len(data))
	}

	m.shardingFactor = binary.BigEndian.Uint32(data[4:8])
	m.salt = [16]byte(data[8:24])
	m.lastValueTimestamp = binary.BigEndian.Uint64(data[24:32])
	m.sealed = data[32] == 1

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
