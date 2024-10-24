package kvstore

import "errors"

// ErrInvalidKey is returned when a key cannot be interpreted as the requested type.
var ErrInvalidKey = errors.New("invalid key")

// Key represents a key in a TableStore. Each key is scoped to a specific table.
type Key interface {
	// AsString interprets the key as a string and returns it.
	AsString() string
	// AsBytes interprets the key as a byte slice and returns it.
	AsBytes() []byte
	// AsUint64 interprets the key as a uint64 and returns it.
	// Returns an error if the data cannot be interpreted as a uint64.
	AsUint64() (uint64, error)
	// AsInt64 interprets the key as an int64 and returns it.
	// Returns an error if the data cannot be interpreted as an int64.
	AsInt64() (int64, error)
	// AsUint32 interprets the key as a uint32 and returns it.
	// Returns an error if the data cannot be interpreted as a uint32.
	AsUint32() (uint32, error)
	// AsInt32 interprets the key as an int32 and returns it.
	// Returns an error if the data cannot be interpreted as an int32.
	AsInt32() (int32, error)
	// Raw returns the raw byte slice that represents the key. This value
	// may not be equal to the byte slice that was used to create the key, and
	// should be treated as an opaque value.
	Raw() []byte
	// Builder returns the KeyBuilder that created this key.
	Builder() KeyBuilder
}

// KeyBuilder is used to create keys for a TableStore. Each KeyBuilder is scoped to a particular table,
// and can be used to create keys that are within that table.
type KeyBuilder interface {
	// TableName returns the name of the table that this KeyBuilder is scoped to.
	TableName() string
	// Key creates a key from a byte slice.
	Key(key []byte) Key
	// StringKey creates a key from a string. Equivalent to Key([]byte(key)).
	StringKey(key string) Key
	// Uint64Key creates a key from a uint64. Resulting key is an 8-byte big-endian representation of the uint64.
	Uint64Key(key uint64) Key
	// Int64Key creates a key from an int64. Resulting key is an 8-byte big-endian representation of the int64.
	Int64Key(key int64) Key
	// Uint32Key creates a key from a uint32. Resulting key is a 4-byte big-endian representation of the uint32.
	Uint32Key(key uint32) Key
	// Int32Key creates a key from an int32. Resulting key is a 4-byte big-endian representation of the int32.
	Int32Key(key int32) Key
}
