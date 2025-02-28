package litt

import (
	"time"

	"github.com/Layr-Labs/eigenda/litt/types"
)

// Table is a key-value store with a namespace that does not overlap with other tables.
// Values may be written to the table, but once written, they may not be changed or deleted (except via TTL).
type Table interface {
	// Name returns the name of the table.
	Name() string

	// Put stores a value in the database. May not be used to overwrite an existing value.
	// Note that when this method returns, data written may not be crash durable on disk
	// (although the write does have atomicity). In order to ensure crash durability, call Flush().
	Put(key []byte, value []byte) error

	// PutBatch stores multiple values in the database. Similar to Put, but allows for multiple values to be written
	// at once. This may improve performance, but is otherwise has identical properties to a sequence of Put calls
	// (i.e. this method does not atomically write the entire batch).
	PutBatch(batch []*types.KVPair) error

	// Get retrieves a value from the database. The returned boolean indicates whether the key exists in the database
	// (returns false if the key does not exist).
	//
	// For the sake of performance, the returned data is NOT safe to mutate. If you need to modify the data,
	// make a copy of it first. Better to avoid a copy if it's not necessary, though.
	Get(key []byte) ([]byte, bool, error)

	// Flush ensures that all data written to the database is crash durable on disk. When this method returns,
	// all data written by Put() operations is guaranteed to be crash durable. Put() operations called synchronously
	// with this method may not be crash durable after this method returns.
	//
	// Note that data flushed at the same time is not atomic. If the process crashes mid-flush, some data
	// being flushed may become persistent, while some may not. Each individual key-value pair is atomic
	// in the event of a crash, though.
	Flush() error

	// SetTTL sets the time to live for data in this table. This TTL is immediately applied to data already in
	// the table. Note that deletion is lazy. That is, when the data expires, it may not be deleted immediately.
	SetTTL(ttl time.Duration) error

	// SetShardingFactor sets the number of write shards used. Increasing this value increases the number of parallel
	// writes that can be performed.
	SetShardingFactor(shardingFactor uint32) error

	// SetCacheSize sets the cache size, in bytes, for the table. For tables without a cache, this method does nothing.
	// If the cache size is set to 0 (default), the cache is disabled. The size of each cache entry is equal to the sum
	// the key length and the value length. Note that the actual in-memory footprint if the cache will be slightly
	// larger than the cache size due to implementation overhead (e.g. pointers, slice headers, map entries, etc.).
	SetCacheSize(size uint64)
}

// ManagedTable is a Table that can perform garbage collection on its data. This type should not be directly used
// by clients, and is a type that is used internally by the database.
type ManagedTable interface {
	Table

	// Start starts the table. Should be called prior to any other methods.
	Start() error

	// Stop shuts down the table, flushing data to disk.
	Stop() error

	// Destroy cleans up resources used by the table. Table will not be usable after this method is called.
	Destroy() error
}
