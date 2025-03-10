package segment

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"path"

	"github.com/Layr-Labs/eigensdk-go/logging"
)

// ValuesFileExtension is the file extension for the values file. This file contains the values for the data
// segment. Value files are written in the form "X-Y.values", where X is the segment index and Y is the shard number.
const ValuesFileExtension = ".values"

// valueFile represents a file that stores values.
type valueFile struct {
	// The logger for the value file.
	logger logging.Logger

	// The segment index.
	index uint32

	// The shard number of this value file.
	shard uint32

	// The parent directory containing this file.
	parentDirectory string

	// The writer for the file. If the file is sealed, this value is nil.
	writer *bufio.Writer

	// The current size of the file in bytes. Includes both flushed and unflushed data.
	currentSize uint64
}

// newValueFile creates a new value file.
func newValueFile(
	logger logging.Logger,
	index uint32,
	shard uint32,
	parentDirectory string,
	sealed bool) (*valueFile, error) {

	values := &valueFile{
		logger:          logger,
		index:           index,
		shard:           shard,
		parentDirectory: parentDirectory,
	}

	filePath := values.path()
	exists, size, err := verifyFilePermissions(filePath)
	if err != nil {
		return nil, fmt.Errorf("file %s has incorrect permissions: %v", filePath, err)
	}

	values.currentSize = uint64(size)

	if sealed {
		if !exists {
			return nil, fmt.Errorf("value file %s does not exist", filePath)
		}

	} else {
		if exists {
			return nil, fmt.Errorf("value file %s already exists", filePath)
		}

		// Open the file for writing.
		file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open value file %s: %v", filePath, err)
		}

		values.writer = bufio.NewWriter(file)
	}

	return values, nil
}

// name returns the name of the value file.
func (v *valueFile) name() string {
	return fmt.Sprintf("%d-%d%s", v.index, v.shard, ValuesFileExtension)
}

// path returns the path to the value file.
func (v *valueFile) path() string {
	return path.Join(v.parentDirectory, v.name())
}

// read reads a value from the value file.
func (v *valueFile) read(firstByteIndex uint32) ([]byte, error) {
	if uint64(firstByteIndex) >= v.currentSize {
		return nil, fmt.Errorf("index %d is out of bounds (current size is %d)", firstByteIndex, v.currentSize)
	}

	file, err := os.OpenFile(v.path(), os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open value file: %v", err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			v.logger.Errorf("failed to close value file: %v", err)
		}
	}()

	_, err = file.Seek(int64(firstByteIndex), 0)
	reader := bufio.NewReader(file)

	// Read the length of the value.
	var length uint32
	err = binary.Read(reader, binary.BigEndian, &length)
	if err != nil {
		return nil, fmt.Errorf("failed to read value length from value file: %v", err)
	}

	// Read the value itself.
	value := make([]byte, length)
	bytesRead, err := io.ReadFull(reader, value)
	if err != nil {
		return nil, fmt.Errorf("failed to read value from value file: %v", err)
	}

	if uint32(bytesRead) != length {
		return nil, fmt.Errorf("failed to read value from value file: read %d bytes, expected %d", bytesRead, length)
	}

	return value, nil
}

// write writes a value to the value file, returning the index of the first byte written.
func (v *valueFile) write(value []byte) (uint32, error) {
	if v.writer == nil {
		return 0, fmt.Errorf("value file is sealed")
	}

	if v.currentSize > math.MaxUint32 {
		// No matter what, we can't start a new value if its first byte would be beyond position 2^32.
		// This is because we only have 32 bits in an address to store the position of a value's first byte.
		return 0, fmt.Errorf("value file already contains %d bytes, cannot add a new value", v.currentSize)
	}

	firstByteIndex := uint32(v.currentSize)

	// First, write the length of the value.
	err := binary.Write(v.writer, binary.BigEndian, uint32(len(value)))
	if err != nil {
		return 0, fmt.Errorf("failed to write value length to value file: %v", err)
	}

	// Then, write the value itself.
	_, err = v.writer.Write(value)
	if err != nil {
		return 0, fmt.Errorf("failed to write value to value file: %v", err)
	}

	v.currentSize += uint64(len(value)) + 4

	return firstByteIndex, nil
}

// flush writes all unflushed data to disk.
func (v *valueFile) flush() error {
	if v.writer == nil {
		return fmt.Errorf("value file is sealed")
	}

	err := v.writer.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush value file: %v", err)
	}

	return nil
}

// seal seals the value file.
func (v *valueFile) seal() error {
	if v.writer == nil {
		return fmt.Errorf("value file is already sealed")
	}

	err := v.flush()
	if err != nil {
		return fmt.Errorf("failed to flush value file: %v", err)
	}

	v.writer = nil
	return nil
}

// delete deletes the value file.
func (v *valueFile) delete() error {
	if v.writer != nil {
		return fmt.Errorf("value file is not sealed")
	}

	err := os.Remove(v.path())
	if err != nil {
		return fmt.Errorf("failed to delete value file: %v", err)
	}

	return nil
}
