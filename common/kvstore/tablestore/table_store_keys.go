package tablestore

import (
	"encoding/binary"
	"github.com/Layr-Labs/eigenda/common/kvstore"
)

const prefixLength = 4

var _ kvstore.Key = (*key)(nil)

type key struct {
	data []byte
}

// RawKey returns the raw byte slice that represents the key.
func (k *key) Raw() []byte {
	return k.data
}

// AsString interprets the key as a string and returns it.
func (k *key) AsString() string {
	return string(k.data[prefixLength:])
}

// AsBytes interprets the key as a byte slice and returns it.
func (k *key) AsBytes() []byte {
	return k.data[prefixLength:]
}

// AsUint64 interprets the key as a uint64 and returns it.
func (k *key) AsUint64() (uint64, error) {
	if len(k.data) != 12 {
		return 0, kvstore.ErrInvalidKey
	}
	return binary.BigEndian.Uint64(k.data[prefixLength:]), nil
}

// AsInt64 interprets the key as an int64 and returns it.
func (k *key) AsInt64() (int64, error) {
	if len(k.data) != 12 {
		return 0, kvstore.ErrInvalidKey
	}
	return int64(binary.BigEndian.Uint64(k.data[prefixLength:])), nil
}

// AsUint32 interprets the key as a uint32 and returns it.
func (k *key) AsUint32() (uint32, error) {
	if len(k.data) != 8 {
		return 0, kvstore.ErrInvalidKey
	}
	return binary.BigEndian.Uint32(k.data[prefixLength:]), nil
}

// AsInt32 interprets the key as an int32 and returns it.
func (k *key) AsInt32() (int32, error) {
	if len(k.data) != 8 {
		return 0, kvstore.ErrInvalidKey
	}
	return int32(binary.BigEndian.Uint32(k.data[prefixLength:])), nil
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
	return &key{data: result}
}

// StringKey creates a key from a string. Equivalent to Key([]byte(key)).
func (k *keyBuilder) StringKey(data string) kvstore.Key {
	return k.Key([]byte(data))
}

// Uint64Key creates a key from a uint64. Resulting key is an 8-byte big-endian representation of the uint64.
func (k *keyBuilder) Uint64Key(data uint64) kvstore.Key {
	result := make([]byte, 12)
	copy(result, k.prefix)
	binary.BigEndian.PutUint64(result[prefixLength:], data)
	return &key{data: result}
}

// Int64Key creates a key from an int64. Resulting key is an 8-byte big-endian representation of the int64.
func (k *keyBuilder) Int64Key(data int64) kvstore.Key {
	result := make([]byte, 12)
	copy(result, k.prefix)
	binary.BigEndian.PutUint64(result[prefixLength:], uint64(data))
	return &key{data: result}
}

// Uint32Key creates a key from a uint32. Resulting key is a 4-byte big-endian representation of the uint32.
func (k *keyBuilder) Uint32Key(data uint32) kvstore.Key {
	result := make([]byte, 8)
	copy(result, k.prefix)
	binary.BigEndian.PutUint32(result[prefixLength:], data)
	return &key{data: result}
}

// Int32Key creates a key from an int32. Resulting key is a 4-byte big-endian representation of the int32.
func (k *keyBuilder) Int32Key(data int32) kvstore.Key {
	result := make([]byte, 8)
	copy(result, k.prefix)
	binary.BigEndian.PutUint32(result[prefixLength:], uint32(data))
	return &key{data: result}
}
