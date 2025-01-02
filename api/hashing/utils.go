package hashing

import (
	"encoding/binary"
	"hash"
)

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
