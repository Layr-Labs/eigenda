package litt

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/litt/disktable/keymap"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/docker/go-units"
	"github.com/prometheus/client_golang/prometheus"
)

// Config is configuration for a litt.DB.
type Config struct {
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
	LoggerConfig *common.LoggerConfig

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

	// The maximum number of keys in a segment. The default is 32,000. For workloads with moderately large values
	// (i.e. in the kb+ range), this threshold is unlikely to be relevant. For workloads with very small values,
	// this constant prevents a segment from accumulating too many keys. A segment with too many keys may have
	// undesirable properties such as a very large key file and very slow garbage collection (since no kv-pair in
	// a segment can be deleted until the entire segment is deleted).
	MaxSegmentKeyCount uint64

	// The desired maximum size for a key file. The default is 1 MB. When a key file exceeds this size, the segment
	// will close the current segment and begin writing to a new one. For workloads with moderately large values,
	// this threshold is unlikely to be relevant. For workloads with very small values, this constant prevents a key
	// file from growing too large. A key file with too many keys may have undesirable properties such as very slow
	// garbage collection (since no kv-pair in a segment can be deleted until the entire segment is deleted).
	TargetSegmentKeyFileSize uint64

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
func DefaultConfig(paths ...string) (*Config, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("at least one path must be provided")
	}

	seed := time.Now().UnixNano()
	saltShaker := rand.New(rand.NewSource(seed))

	loggerConfig := common.DefaultLoggerConfig()

	return &Config{
		CTX:                      context.Background(),
		Paths:                    paths,
		LoggerConfig:             &loggerConfig,
		TimeSource:               time.Now,
		GCPeriod:                 5 * time.Minute,
		ShardingFactor:           8,
		SaltShaker:               saltShaker,
		KeymapType:               keymap.LevelDBKeymapType,
		ControlChannelSize:       64,
		TargetSegmentFileSize:    math.MaxUint32,
		MaxSegmentKeyCount:       32_000,
		TargetSegmentKeyFileSize: units.MiB,
		Fsync:                    true,
		DoubleWriteProtection:    false,
		MetricsEnabled:           false,
		MetricsNamespace:         "litt",
		MetricsPort:              8080,
		MetricsUpdateInterval:    time.Second,
	}, nil
}

// SanityCheck performs a sanity check on the configuration, returning an error if any of the configuration
// settings are invalid. The config returned by DefaultConfig() is guaranteed to pass this check if unmodified.
func (c *Config) SanityCheck() error {
	if c.CTX == nil {
		return fmt.Errorf("context cannot be nil")
	}
	if c.Paths == nil || len(c.Paths) == 0 {
		return fmt.Errorf("at least one path must be provided")
	}
	if c.Logger == nil && c.LoggerConfig == nil {
		return fmt.Errorf("logger or logger config must be provided")
	}
	if c.TimeSource == nil {
		return fmt.Errorf("time source cannot be nil")
	}
	if c.ShardingFactor == 0 {
		return fmt.Errorf("sharding factor must be at least 1")
	}
	if c.ControlChannelSize == 0 {
		return fmt.Errorf("control channel size must be at least 1")
	}
	if c.TargetSegmentFileSize == 0 {
		return fmt.Errorf("target segment file size must be at least 1")
	}
	if c.MaxSegmentKeyCount == 0 {
		return fmt.Errorf("max segment key count must be at least 1")
	}
	if c.TargetSegmentKeyFileSize == 0 {
		return fmt.Errorf("target segment key file size must be at least 1")
	}
	if c.GCPeriod == 0 {
		return fmt.Errorf("gc period must be at least 1")
	}
	if c.SaltShaker == nil {
		return fmt.Errorf("salt shaker cannot be nil")
	}
	if (c.MetricsEnabled || c.MetricsRegistry != nil) && c.MetricsUpdateInterval == 0 {
		return fmt.Errorf("metrics update interval must be at least 1 if metrics are enabled")
	}

	return nil
}
