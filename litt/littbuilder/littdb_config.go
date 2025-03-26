package littbuilder

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/cache"
	"github.com/Layr-Labs/eigenda/litt"
	tablecache "github.com/Layr-Labs/eigenda/litt/cache"
	"github.com/Layr-Labs/eigenda/litt/disktable"
	"github.com/Layr-Labs/eigenda/litt/disktable/keymap"
	"github.com/Layr-Labs/eigenda/litt/metrics"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// keymapBuilders contains builders for all supported keymap types.
var keymapBuilders = map[keymap.KeymapType]keymap.KeymapBuilder{
	keymap.MemKeymapType:     keymap.NewMemKeymapBuilder(),
	keymap.LevelDBKeymapType: keymap.NewLevelDBKeymapBuilder(),
}

// LittDBConfig is configuration for a litt.DB.
type LittDBConfig struct {
	// The context for the database. If nil, context.Background() is used.
	CTX context.Context

	// The paths where the database will store its files. If the path does not exist, it will be created.
	// If more than one path is provided, then the database will do its best to spread out the data across
	// the paths. If the database is restarted, it will attempt to load data from all paths. Note: the number
	// of paths should not exceed the sharding factor, or else data may not be split across all paths.
	Paths []string

	// The logger for the database. If nil, a logger is built using the LoggerConfig.
	Logger logging.Logger

	// The logger configuration for the database. Ignored if Logger is not nil.
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

	// If enabled, the database will return an error if a key is written but that key is already present in
	// the database. Updating existing keys is illegal and may result in unexpected behavior, and so this check
	// acts as a safety mechanism against this sort of illegal operation. Unfortunately, if using a keymap other
	// than keymap.MemKeymapType, performing this check may be very expensive. By default, this is false.
	DoubleWriteProtection bool

	// If enabled, collect DB metrics and export them to prometheus. By default, this is false.
	MetricsEnabled bool

	// The namespace to use for metrics. If empty, the default namespace "litt" is used.
	MetricsNamespace string

	// The prometheus registry to use for metrics. If nil and metrics are enabled, a new registry is created.
	MetricsRegistry *prometheus.Registry

	// The port to use for the metrics server. Ignored if MetricsEnabled is false or MetricsRegistry is not nil.
	MetricsPort int

	// The interval at which various DB metrics are updated. The default is 1 second.
	MetricsUpdateInterval time.Duration
}

// DefaultConfig returns a Config with default values.
func DefaultConfig(paths ...string) (*LittDBConfig, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("at least one path must be provided")
	}

	seed := time.Now().UnixNano()
	saltShaker := rand.New(rand.NewSource(seed))

	return &LittDBConfig{
		CTX:                   context.Background(),
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
		DoubleWriteProtection: false,
		MetricsEnabled:        false,
		MetricsNamespace:      "litt",
		MetricsPort:           8080,
		MetricsUpdateInterval: time.Second,
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
	kmap, requiresReload, err = builderForConfiguredType.Build(logger, keymapDataDirectory, c.DoubleWriteProtection)
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
	ttl time.Duration,
	metrics *metrics.LittDBMetrics) (litt.ManagedTable, error) {

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
		c.Fsync,
		metrics)

	if err != nil {
		return nil, fmt.Errorf("error creating table: %w", err)
	}

	tableCache := cache.NewFIFOCache[string, []byte](c.CacheSize, cacheWeight)
	tableCache = cache.NewThreadSafeCache(tableCache)
	cachedTable := tablecache.NewCachedTable(table, tableCache)

	return cachedTable, nil
}

// buildLogger creates a new logger based on the configuration.
func (c *LittDBConfig) buildLogger() (logging.Logger, error) {
	if c.Logger != nil {
		return c.Logger, nil
	}

	return common.NewLogger(c.LoggerConfig)
}

// buildMetrics creates a new metrics object based on the configuration. If the returned server is not nil,
// then it is the responsibility of the caller to eventually call server.Shutdown().
func (c *LittDBConfig) buildMetrics(logger logging.Logger) (*metrics.LittDBMetrics, *http.Server) {
	if !c.MetricsEnabled {
		return nil, nil
	}

	var registry *prometheus.Registry
	var server *http.Server

	if c.MetricsRegistry != nil {
		registry = prometheus.NewRegistry()
		registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
		registry.MustRegister(collectors.NewGoCollector())

		logger.Infof("Starting metrics server at port %d", c.MetricsPort)
		addr := fmt.Sprintf(":%d", c.MetricsPort)
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.HandlerFor(
			registry,
			promhttp.HandlerOpts{},
		))
		server = &http.Server{
			Addr:    addr,
			Handler: mux,
		}

		go func() {
			err := server.ListenAndServe()
			if err != nil && !strings.Contains(err.Error(), "http: Server closed") {
				logger.Errorf("metrics server error: %v", err)
			}
		}()
	}

	return metrics.NewLittDBMetrics(registry, c.MetricsNamespace), server
}

// Build creates a new litt.DB from the configuration. After calling Build, the configuration should not be
// modified. This method is syntactic equivalent to NewDB(config).
func (c *LittDBConfig) Build() (litt.DB, error) {
	return NewDB(c)
}
