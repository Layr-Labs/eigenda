package leveldb

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/syndtr/goleveldb/leveldb"
)

// MetricsConfig holds configuration for the metrics collector
type MetricsConfig struct {
	CollectionInterval   time.Duration
	DegradationThreshold time.Duration
	Name                 string // Identifier for this LevelDB instance
}

// DefaultMetricsConfig provides sensible defaults
var DefaultMetricsConfig = MetricsConfig{
	CollectionInterval:   3 * time.Second,
	DegradationThreshold: time.Minute,
	Name:                 "default",
}

// MetricsCollector manages LevelDB metrics collection
type MetricsCollector struct {
	db     *leveldb.DB
	logger logging.Logger
	config MetricsConfig

	// Synchronization
	mu       sync.RWMutex
	stopChan chan struct{}
	stopped  bool

	// State tracking
	lastStats      leveldb.DBStats
	lastCollection time.Time
	lastWarning    time.Time
}

// Metrics definitions
var (
	// Compaction metrics
	compactionLatency    *prometheus.HistogramVec
	compactionThroughput *prometheus.GaugeVec
	totalCompactionTime  *prometheus.GaugeVec
	compactionCount      *prometheus.GaugeVec

	// Resource utilization metrics
	openTableCount *prometheus.GaugeVec
	blockCacheSize *prometheus.GaugeVec

	// Performance metrics
	readThroughput  *prometheus.GaugeVec
	writeThroughput *prometheus.GaugeVec
	writePaused     *prometheus.GaugeVec

	// Level-specific metrics
	levelTableCount *prometheus.GaugeVec
	levelSize       *prometheus.GaugeVec
	levelReadBytes  *prometheus.GaugeVec
	levelWriteBytes *prometheus.GaugeVec
)

func newLevelDBMetrics(reg *prometheus.Registry) error {
	if reg == nil {
		return errors.New("prometheus registry cannot be nil")
	}

	// Compaction metrics
	compactionLatencyMetric := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:      "compaction_duration_seconds",
		Namespace: "eigenda",
		Subsystem: "leveldb",
		Help:      "Duration of compaction operations by type (memory, level0, non-level0)",
		Buckets:   prometheus.ExponentialBuckets(0.001, 2, 15),
	}, []string{"type", "name"})

	if err := reg.Register(compactionLatencyMetric); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			compactionLatency = are.ExistingCollector.(*prometheus.HistogramVec)
		} else {
			return fmt.Errorf("failed to register compaction latency metric: %w", err)
		}
	} else {
		compactionLatency = compactionLatencyMetric
	}

	compactionThroughputMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "compaction_throughput_bytes_per_second",
		Namespace: "eigenda",
		Subsystem: "leveldb",
		Help:      "Rate of data processed during compaction operations (read/write)",
	}, []string{"operation", "name"})

	if err := reg.Register(compactionThroughputMetric); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			compactionThroughput = are.ExistingCollector.(*prometheus.GaugeVec)
		} else {
			return fmt.Errorf("failed to register compaction throughput metric: %w", err)
		}
	} else {
		compactionThroughput = compactionThroughputMetric
	}

	totalCompactionTimeMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "total_compaction_time_seconds",
		Namespace: "eigenda",
		Subsystem: "leveldb",
		Help:      "Total time spent in compaction across all levels",
	}, []string{"name"})

	if err := reg.Register(totalCompactionTimeMetric); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			totalCompactionTime = are.ExistingCollector.(*prometheus.GaugeVec)
		} else {
			return fmt.Errorf("failed to register total compaction time metric: %w", err)
		}
	} else {
		totalCompactionTime = totalCompactionTimeMetric
	}

	// Resource utilization metrics
	readThroughputMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "read_throughput_bytes_per_second",
		Namespace: "eigenda",
		Subsystem: "leveldb",
		Help:      "Rate of bytes read per second",
	}, []string{"name"})

	if err := reg.Register(readThroughputMetric); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			readThroughput = are.ExistingCollector.(*prometheus.GaugeVec)
		} else {
			return fmt.Errorf("failed to register read throughput metric: %w", err)
		}
	} else {
		readThroughput = readThroughputMetric
	}

	writeThroughputMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "write_throughput_bytes_per_second",
		Namespace: "eigenda",
		Subsystem: "leveldb",
		Help:      "Rate of bytes written per second",
	}, []string{"name"})

	if err := reg.Register(writeThroughputMetric); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			writeThroughput = are.ExistingCollector.(*prometheus.GaugeVec)
		} else {
			return fmt.Errorf("failed to register write throughput metric: %w", err)
		}
	} else {
		writeThroughput = writeThroughputMetric
	}

	openTableCountMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "open_tables_total",
		Namespace: "eigenda",
		Subsystem: "leveldb",
		Help:      "Number of currently open tables",
	}, []string{"name"})

	if err := reg.Register(openTableCountMetric); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			openTableCount = are.ExistingCollector.(*prometheus.GaugeVec)
		} else {
			return fmt.Errorf("failed to register open table count metric: %w", err)
		}
	} else {
		openTableCount = openTableCountMetric
	}

	blockCacheSizeMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "block_cache_bytes",
		Namespace: "eigenda",
		Subsystem: "leveldb",
		Help:      "Size of block cache in bytes",
	}, []string{"name"})

	if err := reg.Register(blockCacheSizeMetric); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			blockCacheSize = are.ExistingCollector.(*prometheus.GaugeVec)
		} else {
			return fmt.Errorf("failed to register block cache size metric: %w", err)
		}
	} else {
		blockCacheSize = blockCacheSizeMetric
	}

	// Performance metrics
	compactionCountMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "compactions_total",
		Namespace: "eigenda",
		Subsystem: "leveldb",
		Help:      "Number of compactions by type (memory, level0, nonlevel0, seek)",
	}, []string{"type", "name"})

	if err := reg.Register(compactionCountMetric); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			compactionCount = are.ExistingCollector.(*prometheus.GaugeVec)
		} else {
			return fmt.Errorf("failed to register compaction count metric: %w", err)
		}
	} else {
		compactionCount = compactionCountMetric
	}

	writePausedMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "write_paused",
		Namespace: "eigenda",
		Subsystem: "leveldb",
		Help:      "Whether writes are currently paused (1 for yes, 0 for no)",
	}, []string{"name"})

	if err := reg.Register(writePausedMetric); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			writePaused = are.ExistingCollector.(*prometheus.GaugeVec)
		} else {
			return fmt.Errorf("failed to register write paused metric: %w", err)
		}
	} else {
		writePaused = writePausedMetric
	}

	// Level-specific metrics
	levelTableCountMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "level_tables_total",
		Namespace: "eigenda",
		Subsystem: "leveldb",
		Help:      "Number of tables in each level",
	}, []string{"level", "name"})

	if err := reg.Register(levelTableCountMetric); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			levelTableCount = are.ExistingCollector.(*prometheus.GaugeVec)
		} else {
			return fmt.Errorf("failed to register level table count metric: %w", err)
		}
	} else {
		levelTableCount = levelTableCountMetric
	}

	levelSizeMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "level_size_bytes",
		Namespace: "eigenda",
		Subsystem: "leveldb",
		Help:      "Size of each level in bytes",
	}, []string{"level", "name"})

	if err := reg.Register(levelSizeMetric); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			levelSize = are.ExistingCollector.(*prometheus.GaugeVec)
		} else {
			return fmt.Errorf("failed to register level size metric: %w", err)
		}
	} else {
		levelSize = levelSizeMetric
	}

	levelReadBytesMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "level_read_bytes_total",
		Namespace: "eigenda",
		Subsystem: "leveldb",
		Help:      "Total bytes read from each level",
	}, []string{"level", "name"})

	if err := reg.Register(levelReadBytesMetric); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			levelReadBytes = are.ExistingCollector.(*prometheus.GaugeVec)
		} else {
			return fmt.Errorf("failed to register level read bytes metric: %w", err)
		}
	} else {
		levelReadBytes = levelReadBytesMetric
	}

	levelWriteBytesMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "level_write_bytes_total",
		Namespace: "eigenda",
		Subsystem: "leveldb",
		Help:      "Total bytes written to each level",
	}, []string{"level", "name"})

	if err := reg.Register(levelWriteBytesMetric); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			levelWriteBytes = are.ExistingCollector.(*prometheus.GaugeVec)
		} else {
			return fmt.Errorf("failed to register level write bytes metric: %w", err)
		}
	} else {
		levelWriteBytes = levelWriteBytesMetric
	}

	return nil
}

// NewMetricsCollector creates a new metrics collector with the given configuration
func NewMetricsCollector(db *leveldb.DB, logger logging.Logger, config MetricsConfig, reg *prometheus.Registry) *MetricsCollector {
	if config.CollectionInterval == 0 {
		config.CollectionInterval = DefaultMetricsConfig.CollectionInterval
	}
	if config.DegradationThreshold == 0 {
		config.DegradationThreshold = DefaultMetricsConfig.DegradationThreshold
	}

	if err := newLevelDBMetrics(reg); err != nil {
		logger.Error("Failed to initialize LevelDB metrics", "error", err)
		return nil
	}

	mc := &MetricsCollector{
		db:       db,
		logger:   logger,
		config:   config,
		stopChan: make(chan struct{}),
	}

	go mc.collectionLoop()
	return mc
}

// Stop gracefully stops the metrics collection
func (mc *MetricsCollector) Stop() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if !mc.stopped {
		close(mc.stopChan)
		mc.stopped = true
	}
}

func (mc *MetricsCollector) collectionLoop() {
	ticker := time.NewTicker(mc.config.CollectionInterval)
	defer ticker.Stop()

	for {
		select {
		case <-mc.stopChan:
			return
		case <-ticker.C:
			mc.collectMetrics()
		}
	}
}

func (mc *MetricsCollector) collectMetrics() {
	var stats leveldb.DBStats
	if err := mc.db.Stats(&stats); err != nil {
		mc.logger.Error("Failed to collect database stats", "error", err)
		return
	}

	mc.mu.Lock()
	defer mc.mu.Unlock()

	// Calculate time-based deltas
	timeDelta := time.Since(mc.lastCollection).Seconds()
	if timeDelta == 0 {
		return // Avoid division by zero
	}

	// Process compaction metrics
	mc.processMetrics(&stats, timeDelta)

	// Check for performance degradation
	mc.checkDegradation(&stats)

	// Update state
	mc.lastStats = stats
	mc.lastCollection = time.Now()
}

func (mc *MetricsCollector) processMetrics(stats *leveldb.DBStats, timeDelta float64) {
	// Calculate compaction latencies
	if !mc.lastCollection.IsZero() && len(mc.lastStats.LevelDurations) > 0 {
		for level, duration := range stats.LevelDurations {
			if level < len(mc.lastStats.LevelDurations) {
				prevDuration := mc.lastStats.LevelDurations[level]
				deltaDuration := duration - prevDuration
				if deltaDuration > 0 {
					compactionLatency.WithLabelValues(getLevelName(level), mc.config.Name).Observe(deltaDuration.Seconds())
				}
			}
		}
	}

	// Calculate total compaction time (cumulative)
	var totalDuration time.Duration
	for _, duration := range stats.LevelDurations {
		totalDuration += duration
	}

	totalCompactionTime.WithLabelValues(mc.config.Name).Set(totalDuration.Seconds())

	// Calculate throughput metrics
	if prevStats := mc.lastStats; prevStats.LevelRead != nil {
		readDelta := stats.LevelRead.Sum() - prevStats.LevelRead.Sum()
		writeDelta := stats.LevelWrite.Sum() - prevStats.LevelWrite.Sum()

		compactionThroughput.WithLabelValues("read", mc.config.Name).Set(float64(readDelta) / timeDelta)
		compactionThroughput.WithLabelValues("write", mc.config.Name).Set(float64(writeDelta) / timeDelta)
	}

	// Update compaction counters for each level
	compactionCount.WithLabelValues("memory", mc.config.Name).Set(float64(stats.MemComp))
	compactionCount.WithLabelValues("level0", mc.config.Name).Set(float64(stats.Level0Comp))
	compactionCount.WithLabelValues("nonlevel0", mc.config.Name).Set(float64(stats.NonLevel0Comp))
	compactionCount.WithLabelValues("seek", mc.config.Name).Set(float64(stats.SeekComp))

	// Process resource metrics
	openTableCount.WithLabelValues(mc.config.Name).Set(float64(stats.OpenedTablesCount))
	blockCacheSize.WithLabelValues(mc.config.Name).Set(float64(stats.BlockCacheSize))

	// Process performance metrics
	if !mc.lastCollection.IsZero() {
		readDelta := float64(stats.IORead - mc.lastStats.IORead)
		writeDelta := float64(stats.IOWrite - mc.lastStats.IOWrite)

		readThroughput.WithLabelValues(mc.config.Name).Set(readDelta / timeDelta)
		writeThroughput.WithLabelValues(mc.config.Name).Set(writeDelta / timeDelta)
	}

	// Track write pauses
	if stats.WritePaused {
		writePaused.WithLabelValues(mc.config.Name).Set(1)
	} else {
		writePaused.WithLabelValues(mc.config.Name).Set(0)
	}

	// Process level-specific metrics
	for level := range stats.LevelTablesCounts {
		levelName := getLevelName(level)
		levelTableCount.WithLabelValues(levelName, mc.config.Name).Set(float64(stats.LevelTablesCounts[level]))

		if stats.LevelSizes != nil {
			levelSize.WithLabelValues(levelName, mc.config.Name).Set(float64(stats.LevelSizes[level]))
		}

		if stats.LevelRead != nil {
			levelReadBytes.WithLabelValues(levelName, mc.config.Name).Set(float64(stats.LevelRead[level]))
		}

		if stats.LevelWrite != nil {
			levelWriteBytes.WithLabelValues(levelName, mc.config.Name).Set(float64(stats.LevelWrite[level]))
		}
	}
}

func (mc *MetricsCollector) checkDegradation(stats *leveldb.DBStats) {
	if !stats.WritePaused {
		return
	}

	now := time.Now()
	if now.Sub(mc.lastWarning) < mc.config.DegradationThreshold {
		return
	}

	mc.logger.Warn("Database performance degraded due to compaction")
	mc.lastWarning = now
}

func getLevelName(level int) string {
	if level == 0 {
		return "memory"
	}
	return "level_" + strconv.Itoa(level)
}
