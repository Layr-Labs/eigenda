package tablestore

import (
	"encoding/binary"
	"github.com/Layr-Labs/eigenda/common/kvstore"
)

const prefixLength = 4

var _ kvstore.Key = (*key)(nil)

type key struct {
	keyBuilder kvstore.KeyBuilder
	data       []byte
}

// Builder returns the KeyBuilder that was used to create the key.
func (k *key) Builder() kvstore.KeyBuilder {
	return k.keyBuilder
}

// Raw returns the raw byte slice that represents the key.
func (k *key) Raw() []byte {
	return k.data
}

// Bytes interprets the key as a byte slice and returns it.
func (k *key) Bytes() []byte {
	return k.data[prefixLength:]
}

var _ kvstore.KeyBuilder = (*keyBuilder)(nil)

type keyBuilder struct {
	tableName string
	prefix    []byte
}

// newKeyBuilder creates a new KeyBuilder for the given table name and prefix.
func newKeyBuilder(tableName string, prefix uint32) kvstore.KeyBuilder {
	prefixBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(prefixBytes, prefix)
	return &keyBuilder{
		tableName: tableName,
		prefix:    prefixBytes,
	}
}

// TableName returns the name of the table that this KeyBuilder is scoped to.
func (k *keyBuilder) TableName() string {
	return k.tableName
}

// Key creates a key from a byte slice.
func (k *keyBuilder) Key(data []byte) kvstore.Key {
	result := make([]byte, prefixLength+len(data))
	copy(result, k.prefix)
	copy(result[prefixLength:], data)
	return &key{
		keyBuilder: k,
		data:       result,
	}
}
