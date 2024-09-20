package tablestore

import (
	"encoding/binary"
	"github.com/Layr-Labs/eigenda/common/kvstore"
)

var _ kvstore.Key = &key{}

// key is a key in a TableStore.
type key struct {
	// the prefix for the table
	prefix []byte
	// the key within the table
	key []byte
}

// GetKeyBytes returns the key within the table, interpreted as a byte slice.
func (k *key) GetKeyBytes() []byte {
	return k.key
}

// GetKeyString returns the key within the table, interpreted as a string. Calling this
// method on keys that do not represent a string may return odd results.
func (k *key) GetKeyString() string {
	return string(k.key)
}

// GetKeyUint32 returns the key within the table, interpreted as a uint32. Calling this
// method on keys that do not represent a uint32 may return odd results.
func (k *key) GetKeyUint32() uint32 {
	if len(k.key) == 4 {
		return binary.BigEndian.Uint32(k.key)
	} else if len(k.key) == 0 {
		return 0
	} else if len(k.key) < 4 {
		slice := make([]byte, 4)
		copy(slice[4-len(k.key):], k.key)
		return binary.BigEndian.Uint32(slice)
	} else {
		return binary.BigEndian.Uint32(k.key[:4])
	}
}

// GetKeyUint64 returns the key within the table, interpreted as a uint64. Calling this
// method on keys that do not represent a uint64 may return odd results.
func (k *key) GetKeyUint64() uint64 {
	if len(k.key) == 8 {
		return binary.BigEndian.Uint64(k.key)
	} else if len(k.key) == 0 {
		return 0
	} else if len(k.key) < 8 {
		slice := make([]byte, 8)
		copy(slice[8-len(k.key):], k.key)
		return binary.BigEndian.Uint64(slice)
	} else {
		return binary.BigEndian.Uint64(k.key[:8])
	}
}

// GetRawBytes gets the representation of the key as used internally by the store.
func (k *key) GetRawBytes() []byte {
	return append(k.prefix, k.key...)
}
