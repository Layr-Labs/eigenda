package tablestore

import (
	"encoding/binary"
	"github.com/Layr-Labs/eigenda/common/kvstore"
)

var _ kvstore.Table = &keyBuilder{}

// keyBuilder is used to create new keys in a specific table.
type keyBuilder struct {
	// the prefix for the table
	prefix []byte
}

// Key creates a new key in a specific table using the given key bytes.
func (k *keyBuilder) Key(keyBytes []byte) kvstore.Key {
	return &key{
		prefix: k.prefix,
		key:    keyBytes,
	}
}

// StringKey creates a new key in a specific table using the given key string.
func (k *keyBuilder) StringKey(keyString string) kvstore.Key {
	return &key{
		prefix: k.prefix,
		key:    []byte(keyString),
	}
}

// Uint32Key creates a new key in a specific table using the given uint32 as a key.
func (k *keyBuilder) Uint32Key(uKey uint32) kvstore.Key {
	keyBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(keyBytes, uKey)
	return &key{
		prefix: k.prefix,
		key:    keyBytes,
	}
}

// Uint64Key creates a new key in a specific table using the given uint64 as a key.
func (k *keyBuilder) Uint64Key(uKey uint64) kvstore.Key {
	keyBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(keyBytes, uKey)
	return &key{
		prefix: k.prefix,
		key:    keyBytes,
	}
}
