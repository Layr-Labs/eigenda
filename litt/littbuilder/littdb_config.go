package littbuilder

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"path"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/cache"
	"github.com/Layr-Labs/eigenda/litt"
	tablecache "github.com/Layr-Labs/eigenda/litt/cache"
	"github.com/Layr-Labs/eigenda/litt/disktable"
	"github.com/Layr-Labs/eigenda/litt/disktable/keymap"
	"github.com/Layr-Labs/eigenda/litt/memtable"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// DBType is an enum for the type of database.
type DBType int

// DiskDB stores data persistently on disk
const DiskDB DBType = 0

// MemDB stores data in memory. This data is not persistent. This database type is useful for testing, but
// probably isn't suitable for production use (i.e. just use a map).
const MemDB DBType = 1

// KeyMapType is an enum for the type of key map. A key map is used to store the mapping between keys and the
// addresses of their corresponding values.
type KeyMapType int

// MemKeyMap stores the key map in memory. It is much faster than LevelDBKeyMap, but it may take a lot longer
// to load from disk, and may have a large memory footprint (depending on the number of keys).
const MemKeyMap KeyMapType = 0

// LevelDBKeyMap stores the key map on disk using LevelDB. This is slower than MemKeyMap, but has a much smaller
// memory footprint, and is much faster to load from disk.
const LevelDBKeyMap KeyMapType = 1

// TODO can these configurations be sorted better?

// LittDBConfig is configuration for a litt.DB.
type LittDBConfig struct {
	// The paths where the database will store its files. If the path does not exist, it will be created.
	// If more than one path is provided, then the database will do its best to spread out the data across
	// the paths. If the database is restarted, it will attempt to load data from all paths. Note: the number
	// of paths should not exceed the sharding factor, or else data may not be split across all paths.
	Paths []string

	// The logger configuration for the database.
	loggerConfig common.LoggerConfig

	// The type of the DB. Choices are DiskDB and MemDB. Default is DiskDB.
	DBType DBType

	// The type of the key map. Choices are MemKeyMap and LevelDBKeyMap. Default is LevelDBKeyMap.
	KeyMapType KeyMapType

	// The default TTL for newly created tables (either ones with data on disk or new tables).
	// The default is 0 (no TTL). TTL can be set individually on each table by calling Table.SetTTL().
	TTL time.Duration

	// The size of the control channel for the segment manager. The default is 64.
	ControlChannelSize int

	// The target size for segments. The default is math.MaxUint32.
	TargetSegmentFileSize uint32

	// The period between garbage collection runs. The default is 5 minutes.
	GCPeriod time.Duration

	// The sharding factor for the database. The default is 8. Must be at least 1.
	ShardingFactor uint32

	// The salt used for sharding. Chosen randomly at boot time by default.
	// This doesn't need to be cryptographically secure, but it should be kept private.
	// In theory, an attacker who knows the salt could craft keys that all hash to the same shard.
	Salt uint32

	// The size of the cache for tables that have not had their cache size set. The default is 0 (no cache).
	// Cache size is in bytes, and includes the size of both the key and the value. Cache size can be set
	// individually on each table by calling Table.SetCacheSize().
	CacheSize uint64

	// The time source used by the database. This can be substituted for an artificial time source
	// for testing purposes. The default is time.Now.
	TimeSource func() time.Time
}

// DefaultConfig returns a Config with default values.
func DefaultConfig(paths ...string) (*LittDBConfig, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("at least one path must be provided")
	}

	return &LittDBConfig{
		Paths:                 paths,
		loggerConfig:          common.DefaultLoggerConfig(),
		TimeSource:            time.Now,
		GCPeriod:              5 * time.Minute,
		ShardingFactor:        8,
		Salt:                  rand.Uint32(),
		DBType:                DiskDB,
		KeyMapType:            LevelDBKeyMap,
		ControlChannelSize:    64,
		TargetSegmentFileSize: math.MaxUint32,
	}, nil
}

// cacheWeight is a function that calculates the weight of a cache entry.
func cacheWeight(key string, value []byte) uint64 {
	return uint64(len(key) + len(value))
}

// buildKeyMap creates a new key map based on the configuration.
func (c *LittDBConfig) buildKeyMap(name string, logger logging.Logger) (keymap.KeyMap, error) {
	var keyMap keymap.KeyMap
	var err error

	switch c.KeyMapType {
	case MemKeyMap:
		keyMap = keymap.NewMemKeyMap(logger)
	case LevelDBKeyMap:
		keymapPath := path.Join(c.Paths[0], name, "keymap")
		keyMap, err = keymap.NewLevelDBKeyMap(logger, keymapPath)
	default:
		return nil, fmt.Errorf("unsupported key map type: %v", c.KeyMapType)
	}

	if err != nil {
		return nil, fmt.Errorf("error creating key map: %w", err)
	}

	return keyMap, nil
}

// buildTable creates a new table based on the configuration.
func (c *LittDBConfig) buildTable(
	ctx context.Context,
	logger logging.Logger,
	timeSource func() time.Time,
	name string,
	ttl time.Duration) (litt.ManagedTable, error) {

	var table litt.ManagedTable

	if c.ShardingFactor < 1 {
		return nil, fmt.Errorf("sharding factor must be at least 1")
	}

	switch c.DBType {
	case DiskDB:
		keyMap, err := c.buildKeyMap(name, logger)
		if err != nil {
			return nil, fmt.Errorf("error creating key map: %w", err)
		}

		tableRoots := make([]string, len(c.Paths))
		for i, p := range c.Paths {
			tableRoots[i] = path.Join(p, name)
		}

		table, err = disktable.NewDiskTable(
			ctx,
			logger,
			time.Now,
			name,
			keyMap,
			tableRoots,
			c.TargetSegmentFileSize,
			c.ControlChannelSize,
			c.ShardingFactor,
			c.Salt,
			c.TTL,
			c.GCPeriod)

		if err != nil {
			return nil, fmt.Errorf("error creating table: %w", err)
		}
	case MemDB:
		table = memtable.NewMemTable(timeSource, name, ttl)
	default:
		return nil, fmt.Errorf("unsupported DB type: %v", c.DBType)
	}

	tableCache := cache.NewFIFOCache[string, []byte](c.CacheSize, cacheWeight)
	tableCache = cache.NewThreadSafeCache(tableCache)
	cachedTable := tablecache.NewCachedTable(table, tableCache)

	return cachedTable, nil

}

// Build creates a new litt.DB from the configuration. After calling Build, the configuration should not be
// modified.
func (c *LittDBConfig) Build(ctx context.Context) (litt.DB, error) {
	logger, err := common.NewLogger(c.loggerConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating logger: %w", err)
	}

	return newDB(ctx, logger, c.TimeSource, c.TTL, c.GCPeriod, c.buildTable), nil
}
