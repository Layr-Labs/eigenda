package segment

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"path"

	"github.com/Layr-Labs/eigenda/litt/types"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// KeysFileExtension is the file extension for the keys file. This file contains the keys for the data segment,
// and is used for performing garbage collection on the key index.
const KeysFileExtension = ".keys"

// keyFile tracks the keys in a segment. It is used to do garbage collection on the key-to-address map.
type keyFile struct {
	// The logger for the key file.
	logger logging.Logger

	// The segment index.
	index uint32

	// The parent directory containing this file.
	parentDirectory string

	// The writer for the file. If the file is sealed, this value is nil.
	writer *bufio.Writer
}

// newKeyFile creates a new key file.
func newKeyFile(
	logger logging.Logger,
	index uint32,
	parentDirectory string,
	sealed bool) (*keyFile, error) {

	keys := &keyFile{
		logger:          logger,
		index:           index,
		parentDirectory: parentDirectory,
	}

	filePath := keys.path()

	exists, _, err := verifyFilePermissions(filePath)
	if err != nil {
		return nil, fmt.Errorf("file is not writeable: %v", err)
	}

	if sealed {
		if !exists {
			return nil, fmt.Errorf("key file %s does not exist", filePath)
		}
	} else {
		if exists {
			return nil, fmt.Errorf("key file %s already exists", filePath)
		}

		flags := os.O_RDWR | os.O_CREATE
		file, err := os.OpenFile(filePath, flags, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open key file: %v", err)
		}

		writer := bufio.NewWriter(file)
		keys.writer = writer
	}

	return keys, nil
}

// name returns the name of the key file.
func (k *keyFile) name() string {
	return fmt.Sprintf("%d%s", k.index, KeysFileExtension)
}

// path returns the full path to the key file.
func (k *keyFile) path() string {
	return path.Join(k.parentDirectory, k.name())
}

// write writes a key to the key file.
func (k *keyFile) write(key []byte, address types.Address) error {
	if k.writer == nil {
		return fmt.Errorf("key file is sealed")
	}

	// First write the length of the key.
	err := binary.Write(k.writer, binary.BigEndian, uint32(len(key)))
	if err != nil {
		return fmt.Errorf("failed to write key length to key file: %v", err)
	}

	// Next, write the key itself.
	_, err = k.writer.Write(key)
	if err != nil {
		return fmt.Errorf("failed to write key to key file: %v", err)
	}

	// Finally, write the address.
	err = binary.Write(k.writer, binary.BigEndian, address)
	if err != nil {
		return fmt.Errorf("failed to write address to key file: %v", err)
	}

	return nil
}

// flush flushes the key file to disk.
func (k *keyFile) flush() error {
	if k.writer == nil {
		return fmt.Errorf("key file is sealed")
	}

	return k.writer.Flush()
}

// seal seals the key file, preventing further writes.
func (k *keyFile) seal() error {
	if k.writer == nil {
		return fmt.Errorf("key file is already sealed")
	}

	err := k.flush()
	if err != nil {
		return fmt.Errorf("failed to flush key file: %v", err)
	}
	k.writer = nil

	return nil
}

// readKeys reads all keys from the key file. This method returns an error if the key file is not sealed.
func (k *keyFile) readKeys() ([]*types.KAPair, error) {
	if k.writer != nil {
		return nil, fmt.Errorf("key file is not sealed")
	}

	file, err := os.Open(k.path())
	if err != nil {
		return nil, fmt.Errorf("failed to open key file: %v", err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			k.logger.Errorf("failed to close key file: %v", err)
		}
	}()

	// Key files are small as long as key length is sane. Safe to read the whole file into memory.
	keyBytes, err := os.ReadFile(k.path())
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %v", err)
	}
	keys := make([]*types.KAPair, 0)

	index := 0
	for {
		if index+4 >= len(keyBytes) {
			break
		}
		keyLength := binary.BigEndian.Uint32(keyBytes[index : index+4])
		index += 4

		if index+int(keyLength)+8 > len(keyBytes) {
			break
		}

		key := keyBytes[index : index+int(keyLength)]
		index += int(keyLength)

		address := types.Address(binary.BigEndian.Uint64(keyBytes[index : index+8]))
		index += 8

		keys = append(keys,
			&types.KAPair{
				Key:     key,
				Address: address,
			})
	}

	if index != len(keyBytes) {
		// This can happen if there is a crash while we are writing to the key file.
		// Recoverable, but best to note the event in the logs.
		k.logger.Warnf("key file %s has %d corrupted bytes", k.path(), len(keyBytes)-index)
	}

	return keys, nil
}

// delete deletes the key file.
func (k *keyFile) delete() error {
	return os.Remove(k.path())
}
