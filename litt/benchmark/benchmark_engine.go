package benchmark

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/litt"
	"github.com/Layr-Labs/eigenda/litt/benchmark/config"
	"github.com/Layr-Labs/eigenda/litt/littbuilder"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/docker/go-units"
	"golang.org/x/time/rate"
)

// BenchmarkEngine is a tool for benchmarking LittDB performance.
type BenchmarkEngine struct {
	ctx    context.Context
	cancel context.CancelFunc
	logger logging.Logger

	// The configuration for the benchmark.
	config *config.BenchmarkConfig

	// The database to be benchmarked.
	db litt.DB

	// The table in the database where data is stored.
	table litt.Table

	// Keeps track of data to read and write.
	dataTracker *DataTracker

	// The maximum write throughput in bytes per second for each worker thread.
	writeBytesPerSecondPerThread uint64
}

// NewBenchmarkEngine creates a new BenchmarkEngine with the given configuration.
func NewBenchmarkEngine(configPath string) (*BenchmarkEngine, error) {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file %s: %w", configPath, err)
	}

	cfg.LittConfig.Logger, err = common.NewLogger(cfg.LittConfig.LoggerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	cfg.LittConfig.ShardingFactor = uint32(len(cfg.LittConfig.Paths))

	db, err := littbuilder.NewDB(cfg.LittConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create db: %w", err)
	}

	table, err := db.GetTable("benchmark")
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	ttl := time.Duration(cfg.TTLHours * float64(time.Hour))
	err = table.SetTTL(ttl)
	if err != nil {
		return nil, fmt.Errorf("failed to set TTL for table: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	dataTracker, err := NewDataTracker(ctx, cfg)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create data tracker: %w", err)
	}

	writeBytesPerSecond := uint64(cfg.MaximumWriteThroughputMB * float64(units.MiB))
	writeBytesPerSecondPerThread := writeBytesPerSecond / uint64(cfg.WriterParallelism)

	return &BenchmarkEngine{
		ctx:                          ctx,
		cancel:                       cancel,
		logger:                       cfg.LittConfig.Logger,
		config:                       cfg,
		db:                           db,
		table:                        table,
		dataTracker:                  dataTracker,
		writeBytesPerSecondPerThread: writeBytesPerSecondPerThread,
	}, nil
}

// Logger returns the logger used by the benchmark engine.
func (b *BenchmarkEngine) Logger() logging.Logger {
	return b.logger
}

// Run executes the benchmark. This method blocks forever, or until the benchmark is stopped via control-C or
// encounters an error.
func (b *BenchmarkEngine) Run() error {

	for i := 0; i < b.config.WriterParallelism; i++ {
		go b.writer()
	}

	// Create a channel to listen for OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for signal
	<-sigChan

	// Cancel the context when signal is received
	b.cancel()

	return nil
}

// writer runs on a goroutine and writes data to the database.
func (b *BenchmarkEngine) writer() {

	maxBatchSize := uint64(b.config.BatchSizeMB * float64(units.MiB))
	throttle := rate.NewLimiter(rate.Limit(b.writeBytesPerSecondPerThread), int(b.writeBytesPerSecondPerThread))

	for {
		select {
		case <-b.ctx.Done():
			break
		default:
			batchSize := uint64(0)

			for batchSize < maxBatchSize {
				writeInfo := b.dataTracker.GetWriteInfo()
				batchSize += uint64(len(writeInfo.Value))

				reservation := throttle.ReserveN(time.Now(), len(writeInfo.Value))
				if !reservation.OK() {
					continue
				}
				if reservation.Delay() > 0 {
					time.Sleep(reservation.Delay())
				}

				err := b.table.Put(writeInfo.Key, writeInfo.Value)
				if err != nil {
					panic(fmt.Sprintf("failed to write data: %v", err)) // TODO not clean
				}

				b.dataTracker.ReportWrite(writeInfo.KeyIndex)
			}

			err := b.table.Flush()
			if err != nil {
				panic(fmt.Sprintf("failed to flush data: %v", err)) // TODO not clean
			}
		}
	}
}
