package kvstore

import "errors"

// ErrInvalidKey is returned when a key cannot be interpreted as the requested type.
var ErrInvalidKey = errors.New("invalid key")

// Key represents a key in a TableStore. Each key is scoped to a specific table.
type Key interface {
	// Bytes returns the key as a byte slice. Does not include internal metadata (i.e. the table).
	Bytes() []byte
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
}
