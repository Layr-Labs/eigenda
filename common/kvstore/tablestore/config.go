package tablestore

import "time"

// StoreType describes the underlying store implementation.
type StoreType int

const (
	// LevelDB is a LevelDB-backed store.
	LevelDB StoreType = iota
	// MapStore is an in-memory store. This store does not preserve data across restarts.
	MapStore
)

// Config is the configuration for a TableStore.
type Config struct {
	// The type of the base store. Default is LevelDB.
	Type StoreType
	// The path to the file system directory where the store will write its data. Default is nil.
	// Some store implementations may ignore this field (e.g. MapStore). Other store implementations may require
	// this field to be set (e.g. LevelDB).
	Path *string
	// If true, the store will perform garbage collection on a background goroutine. Default is true.
	GarbageCollectionEnabled bool
	// If garbage collection is enabled, this is the interval at which it will run. Default is 5 minutes.
	GarbageCollectionInterval time.Duration
	// If garbage collection is enabled, this is the maximum number of entries to delete in a single batch during
	// garbage collection. Default is 1024.
	GarbageCollectionBatchSize uint32
	// The list of tables to create on startup. Any pre-existing table not in this list will be deleted. If
	// this list is nil, the previous schema will be carried forward with no modifications. Default is nil.
	Schema []string
}

// DefaultConfig returns a Config with default values.
func DefaultConfig() *Config {
	return &Config{
		Type:                       LevelDB,
		Path:                       nil,
		GarbageCollectionEnabled:   true,
		GarbageCollectionInterval:  5 * time.Minute,
		GarbageCollectionBatchSize: 1024,
		Schema:                     nil,
	}
}

// DefaultLevelDBConfig returns a Config with default values for a LevelDB store.
func DefaultLevelDBConfig(path string) *Config {
	config := DefaultConfig()
	config.Type = LevelDB
	config.Path = &path
	return config
}

// DefaultMapStoreConfig returns a Config with default values for a MapStore.
func DefaultMapStoreConfig() *Config {
	config := DefaultConfig()
	config.Type = MapStore
	return config
}
