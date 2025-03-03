package hashing

import (
	"encoding/binary"
	"fmt"
	"hash"
	"math"
)

// hashLength hashes the length of the given thing.
func hashLength[T any](hasher hash.Hash, thing []T) error {
	if len(thing) > math.MaxUint32 {
		return fmt.Errorf("array is too long: %d", len(thing))
	}

	hashUint32(hasher, uint32(len(thing)))

	return nil
}

// hashByteArray hashes the given byte array.
func hashByteArray(hasher hash.Hash, bytes []byte) error {
	if len(bytes) > math.MaxUint32 {
		return fmt.Errorf("byte array is too long: %d", len(bytes))
	}

	err := hashLength(hasher, bytes)
	if err != nil {
		return fmt.Errorf("failed to hash length: %w", err)
	}
	hasher.Write(bytes)

	return nil
}

// hashUint32Array hashes the given uint32 array.
func hashUint32Array(hasher hash.Hash, values []uint32) error {
	if len(values) > math.MaxUint32 {
		return fmt.Errorf("uint32 array is too long: %d", len(values))
	}

	err := hashLength(hasher, values)
	if err != nil {
		return fmt.Errorf("failed to hash length: %w", err)
	}
	for _, value := range values {
		hashUint32(hasher, value)
	}

	return nil
}

// hashUint32 hashes the given uint32 value.
func hashUint32(hasher hash.Hash, value uint32) {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, value)
	hasher.Write(bytes)
}

// hashUint64 hashes the given uint64 value.
func hashUint64(hasher hash.Hash, value uint64) {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, value)
	hasher.Write(bytes)
}

// hashInt64 hashes the given int64 value.
func hashInt64(hasher hash.Hash, value int64) {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(value))
	hasher.Write(bytes)
}

// hashChar hashes the given byte value.
func hashChar(hasher hash.Hash, value byte) {
	hasher.Write([]byte{value})
}
