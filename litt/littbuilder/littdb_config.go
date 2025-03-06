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
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// keyMapBuilders contains a list of the supported key map builders.
var keyMapBuilders = []keymap.KeyMapBuilder{
	keymap.NewMemKeyMapBuilder(),
	keymap.NewLevelDBKeyMapBuilder(),
}

// LittDBConfig is configuration for a litt.DB.
type LittDBConfig struct {
	// The paths where the database will store its files. If the path does not exist, it will be created.
	// If more than one path is provided, then the database will do its best to spread out the data across
	// the paths. If the database is restarted, it will attempt to load data from all paths. Note: the number
	// of paths should not exceed the sharding factor, or else data may not be split across all paths.
	Paths []string

	// The logger configuration for the database.
	loggerConfig common.LoggerConfig

	// The type of the key map. Choices are keymap.MemKeyMapType and keymap.LevelDBKeyMapType.
	// Default is keymap.LevelDBKeyMapType.
	KeyMapType keymap.KeyMapType

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

	// The random number generator used for generating sharding salts. The default is math/rand. TODO
	SaltShaker rand.Rand

	// The salt used for sharding. Chosen randomly at boot time by default.
	// This doesn't need to be cryptographically secure, but it should be kept private.
	// In theory, an attacker who knows the salt could craft keys that all hash to the same shard.
	Salt uint32 // TODO regenerate this for each segment

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
		KeyMapType:            keymap.LevelDBKeyMapType,
		ControlChannelSize:    64,
		TargetSegmentFileSize: math.MaxUint32,
	}, nil
}

// cacheWeight is a function that calculates the weight of a cache entry.
func cacheWeight(key string, value []byte) uint64 {
	return uint64(len(key) + len(value))
}

// TODO restrict the names that can be used for tables

// buildKeyMap creates a new key map based on the configuration.
func (c *LittDBConfig) buildKeyMap(logger logging.Logger, tableName string) (keymap.KeyMap, bool, error) {

	keymapDirectories := make([]string, len(c.Paths))
	for i, p := range c.Paths {
		keymapDirectories[i] = path.Join(p, tableName)
	}

	// For each KeyMap type we are not using, delete any files associated with it if they exist.
	var chosenBuilder keymap.KeyMapBuilder
	for _, builder := range keyMapBuilders {
		if builder.Type() == c.KeyMapType {
			chosenBuilder = builder
		} else {
			err := builder.DeleteFiles(logger, keymapDirectories)
			if err != nil {
				return nil, false, fmt.Errorf("error deleting key map files: %w", err)
			}
		}
	}
	if chosenBuilder == nil {
		return nil, false, fmt.Errorf("unsupported key map type: %v", c.KeyMapType)
	}

	keyMap, requiresReload, err := chosenBuilder.Build(logger, keymapDirectories)
	if err != nil {
		return nil, false, fmt.Errorf("error building key map: %w", err)
	}

	return keyMap, requiresReload, nil
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

	keyMap, requiresReload, err := c.buildKeyMap(logger, name)
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
		timeSource,
		name,
		keyMap,
		tableRoots,
		c.TargetSegmentFileSize,
		c.ControlChannelSize,
		c.ShardingFactor,
		c.Salt,
		ttl,
		c.GCPeriod,
		requiresReload)

	if err != nil {
		return nil, fmt.Errorf("error creating table: %w", err)
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

	return NewDB(ctx, logger, c.TimeSource, c.TTL, c.GCPeriod, c.buildTable), nil
}
