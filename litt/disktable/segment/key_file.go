package segment

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/Layr-Labs/eigenda/litt/types"
	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// KeyFileExtension is the file extension for the keys file. This file contains the keys for the data segment,
// and is used for performing garbage collection on the keymap. It can also be used to rebuild the keymap.
const KeyFileExtension = ".keys"

// keyFile tracks the keys in a segment. It is used to do garbage collection on the keymap.
//
// This struct is NOT goroutine safe. It is unsafe to concurrently call write, flush, or seal on the same key file.
// It is not safe to read a key file until it is sealed. Once sealed, read only operations are goroutine safe.
type keyFile struct {
	// The logger for the key file.
	logger logging.Logger

	// The segment index.
	index uint32

	// The parent directory containing this file.
	parentDirectory string

	// The writer for the file. If the file is sealed, this value is nil.
	writer *bufio.Writer

	// The size of the key file in bytes.
	size uint64
}

// newKeyFile creates a new key file.
func createKeyFile(
	logger logging.Logger,
	index uint32,
	parentDirectory string) (*keyFile, error) {

	keys := &keyFile{
		logger:          logger,
		index:           index,
		parentDirectory: parentDirectory,
	}

	filePath := keys.path()

	exists, _, err := util.VerifyFileProperties(filePath)
	if err != nil {
		return nil, fmt.Errorf("can not write to file: %v", err)
	}

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

	return keys, nil
}

// loadKeyFile loads the key file from disk, looking in the given parent directories until it finds the file.
// If the file is not found, it returns an error.
func loadKeyFile(logger logging.Logger, index uint32, parentDirectories []string) (*keyFile, error) {

	keyFileName := fmt.Sprintf("%d%s", index, KeyFileExtension)
	keysPath, err := lookForFile(parentDirectories, keyFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to find key file: %w", err)
	}
	if keysPath == "" {
		return nil, fmt.Errorf("failed to find key file %s", keyFileName)
	}
	parentDirectory := path.Dir(keysPath)

	keys := &keyFile{
		logger:          logger,
		index:           index,
		parentDirectory: parentDirectory,
	}

	filePath := keys.path()

	exists, size, err := util.VerifyFileProperties(filePath)
	if err != nil {
		return nil, fmt.Errorf("can not write to file: %v", err)
	}

	if exists {
		keys.size = uint64(size)
	}

	if !exists {
		return nil, fmt.Errorf("key file %s does not exist", filePath)
	}

	return keys, nil
}

// Size returns the size of the key file in bytes.
func (k *keyFile) Size() uint64 {
	return k.size
}

// name returns the name of the key file.
func (k *keyFile) name() string {
	return fmt.Sprintf("%d%s", k.index, KeyFileExtension)
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

	k.size += uint64(4 + len(key) + 8)

	return nil
}

// getKeyFileIndex returns the index of the key file from the file name. Key file names have the form "X.keys",
// where X is the segment index.
func getKeyFileIndex(fileName string) (uint32, error) {
	baseName := path.Base(fileName)
	indexString := baseName[:len(baseName)-len(KeyFileExtension)]
	index, err := strconv.Atoi(indexString)
	if err != nil {
		return 0, fmt.Errorf("failed to parse index from file name %s: %v", fileName, err)
	}

	return uint32(index), nil
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
// If there are keys that were only partially written (i.e. keys being written when the process crashed), then
// those keys may not be returned. If a key is returned, it is guaranteed to be "whole" (i.e. a partial key will
// never be returned).
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
		// We need at least 4 bytes to read the length of the key.
		if index+4 > len(keyBytes) {
			// There are fewer than 4 bytes left in the file.
			break
		}
		keyLength := int(binary.BigEndian.Uint32(keyBytes[index : index+4]))
		index += 4

		// We need to read the key, as well as the 8 byte address.
		if index+keyLength+8 > len(keyBytes) {
			// There are insufficient bytes left in the file to read the key and address.
			break
		}

		key := keyBytes[index : index+keyLength]
		index += keyLength

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
		k.logger.Warnf("key file %s has %d partial bytes", k.path(), len(keyBytes)-index)
	}

	return keys, nil
}

// delete deletes the key file.
func (k *keyFile) delete() error {
	return os.Remove(k.path())
}
