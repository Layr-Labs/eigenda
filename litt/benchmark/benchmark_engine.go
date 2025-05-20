package benchmark

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Layr-Labs/eigenda/litt"
	"github.com/Layr-Labs/eigenda/litt/littbuilder"
)

// BenchmarkEngine is a tool for benchmarking LittDB performance.
type BenchmarkEngine struct {
	ctx    context.Context
	cancel context.CancelFunc

	// The configuration for the benchmark.
	config *BenchmarkConfig

	// The database to be benchmarked.
	db litt.DB

	// The table in the database where data is stored.
	table litt.Table
}

// NewBenchmarkEngine creates a new BenchmarkEngine with the given configuration.
func NewBenchmarkEngine(configPath string) (*BenchmarkEngine, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file %s: %w", configPath, err)
	}

	db, err := littbuilder.NewDB(config.LittConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create db: %w", err)
	}

	table, err := db.GetTable("benchmark")
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &BenchmarkEngine{
		ctx:    ctx,
		cancel: cancel,
		config: config,
		db:     db,
		table:  table,
	}, nil
}

// Run executes the benchmark. This method blocks forever, or until the benchmark is stopped via control-C or
// encounters an error.
func (b *BenchmarkEngine) Run() error {

	go b.writer()

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

}
