package littbuilder

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
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

// keymapBuilders contains builders for all supported keymap types.
var keymapBuilders = map[keymap.KeymapType]keymap.KeymapBuilder{
	keymap.MemKeymapType:     keymap.NewMemKeymapBuilder(),
	keymap.LevelDBKeymapType: keymap.NewLevelDBKeymapBuilder(),
}

// LittDBConfig is configuration for a litt.DB.
type LittDBConfig struct {
	// The paths where the database will store its files. If the path does not exist, it will be created.
	// If more than one path is provided, then the database will do its best to spread out the data across
	// the paths. If the database is restarted, it will attempt to load data from all paths. Note: the number
	// of paths should not exceed the sharding factor, or else data may not be split across all paths.
	Paths []string

	// The logger configuration for the database.
	LoggerConfig common.LoggerConfig

	// The type of the keymap. Choices are keymap.MemKeymapType and keymap.LevelDBKeymapType.
	// Default is keymap.LevelDBKeymapType.
	KeymapType keymap.KeymapType

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

	// The random number generator used for generating sharding salts. The default is a standard rand.New()
	// seeded by the current time.
	SaltShaker *rand.Rand

	// The size of the cache for tables that have not had their cache size set. The default is 0 (no cache).
	// Cache size is in bytes, and includes the size of both the key and the value. Cache size can be set
	// individually on each table by calling Table.SetCacheSize().
	CacheSize uint64

	// The time source used by the database. This can be substituted for an artificial time source
	// for testing purposes. The default is time.Now.
	TimeSource func() time.Time

	// If true, then flush operations will call fsync on the underlying file to ensure data is flushed out of the
	// operating system's buffer and onto disk. Setting this to false means that even after flushing data,
	// there may be data loss in the advent of an OS/hardware crash.
	//
	// The default is true.
	//
	// Enabling fsync may have performance implications, although this strongly depends on the workload. For large
	// batches that are flushed infrequently, benchmark data suggests that the impact is minimal. For small batches
	// that are flushed frequently, the difference can be severe. For example, when enabled in unit tests that do
	// super tiny and frequent flushes, the difference in performance was an order of magnitude.
	Fsync bool
}

// DefaultConfig returns a Config with default values.
func DefaultConfig(paths ...string) (*LittDBConfig, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("at least one path must be provided")
	}

	seed := time.Now().UnixNano()
	saltShaker := rand.New(rand.NewSource(seed))

	return &LittDBConfig{
		Paths:                 paths,
		LoggerConfig:          common.DefaultLoggerConfig(),
		TimeSource:            time.Now,
		GCPeriod:              5 * time.Minute,
		ShardingFactor:        8,
		SaltShaker:            saltShaker,
		KeymapType:            keymap.LevelDBKeymapType,
		ControlChannelSize:    64,
		TargetSegmentFileSize: math.MaxUint32,
		Fsync:                 true,
	}, nil
}

// cacheWeight is a function that calculates the weight of a cache entry.
func cacheWeight(key string, value []byte) uint64 {
	return uint64(len(key) + len(value))
}

// buildKeymap creates a new keymap based on the configuration.
func (c *LittDBConfig) buildKeymap(
	logger logging.Logger,
	tableName string,
) (kmap keymap.Keymap, keymapPath string, keymapTypeFile *keymap.KeymapTypeFile, requiresReload bool, err error) {

	builderForConfiguredType, ok := keymapBuilders[c.KeymapType]
	if !ok {
		return nil, "", nil, false,
			fmt.Errorf("unsupported keymap type: %v", c.KeymapType)
	}

	keymapDirectories := make([]string, len(c.Paths))
	for i, p := range c.Paths {
		keymapDirectories[i] = path.Join(p, tableName, keymap.KeymapDirectoryName)
	}

	var keymapDirectory string
	for _, directory := range keymapDirectories {
		exists, err := keymap.KeymapFileExists(directory)
		if err != nil {
			return nil, "", nil, false,
				fmt.Errorf("error checking for keymap type file: %w", err)
		}
		if exists {
			keymapDirectory = directory
			keymapTypeFile, err = keymap.LoadKeymapTypeFile(directory)
			if err != nil {
				return nil, "", nil, false,
					fmt.Errorf("error loading keymap type file: %w", err)
			}
			break
		}
	}

	newKeymap := false
	if keymapTypeFile == nil {
		// No previous keymap exists. Either we are starting fresh or the keymap was deleted manually.
		newKeymap = true

		// by convention, always select the first path as the keymap directory
		keymapDirectory = keymapDirectories[0]
		keymapTypeFile = keymap.NewKeymapTypeFile(keymapDirectory, c.KeymapType)

		// create the keymap directory
		err := os.MkdirAll(keymapDirectory, 0755)
		if err != nil {
			return nil, "", nil, false,
				fmt.Errorf("error creating keymap directory: %w", err)
		}

		// write the keymap type file
		err = keymapTypeFile.Write()
		if err != nil {
			return nil, "", nil, false,
				fmt.Errorf("error writing keymap type file: %w", err)
		}

	} else {
		// A previous keymap exists. Check if the keymap type has changed.

		builderForTypeOnDisk, ok := keymapBuilders[keymapTypeFile.Type()]
		if !ok {
			return nil, "", nil, false,
				fmt.Errorf("unsupported keymap type: %v", keymapTypeFile.Type())
		}

		if c.KeymapType != keymapTypeFile.Type() {
			// The previously used keymap type is different from the one in the configuration.

			// delete the old keymap
			err := builderForTypeOnDisk.DeleteFiles(logger, keymapDirectory)
			if err != nil {
				return nil, "", nil, false,
					fmt.Errorf("error deleting keymap files: %w", err)
			}

			// delete the keymap type file
			err = keymapTypeFile.Delete()
			if err != nil {
				return nil, "", nil, false,
					fmt.Errorf("error deleting keymap type file: %w", err)
			}

			// finally, delete the keymap directory
			_, err = os.Stat(keymapDirectory)
			if err == nil {
				err = os.Remove(keymapDirectory)
				if err != nil {
					return nil, "", nil, false,
						fmt.Errorf("error deleting keymap directory: %w", err)
				}
			} else if !os.IsNotExist(err) {
				return nil, "", nil, false,
					fmt.Errorf("error checking for keymap directory: %w", err)
			}
		}
	}

	keymapDataDirectory := path.Join(keymapDirectory, keymap.KeymapDataDirectoryName)
	kmap, requiresReload, err = builderForConfiguredType.Build(logger, keymapDataDirectory)
	if err != nil {
		return nil, "", nil, false,
			fmt.Errorf("error building keymap: %w", err)
	}

	return kmap, keymapDirectory, keymapTypeFile, requiresReload || newKeymap, nil
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

	kmap, keymapDirectory, keymapTypeFile, requiresReload, err := c.buildKeymap(logger, name)
	if err != nil {
		return nil, fmt.Errorf("error creating keymap: %w", err)
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
		kmap,
		keymapDirectory,
		keymapTypeFile,
		tableRoots,
		c.TargetSegmentFileSize,
		c.ControlChannelSize,
		c.ShardingFactor,
		c.SaltShaker,
		ttl,
		c.GCPeriod,
		requiresReload,
		c.Fsync)

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
	logger, err := common.NewLogger(c.LoggerConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating logger: %w", err)
	}

	return NewDB(ctx, logger, c.TimeSource, c.TTL, c.GCPeriod, c.buildTable), nil
}
