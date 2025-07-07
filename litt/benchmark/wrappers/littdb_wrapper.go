package wrappers

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/litt"
	"github.com/Layr-Labs/eigenda/litt/benchmark/config"
	"github.com/Layr-Labs/eigenda/litt/littbuilder"
)

var _ DatabaseWrapper = (*LittDBWrapper)(nil)

// LittDB doesn't need a new instance per thread, it's perfectly happy to be used by multiple threads at once.
var _ ThreadLocalDatabaseWrapper = (*LittDBWrapper)(nil)

// LittDBWrapper is a wrapper around a LittDB database, enabling it to be used by the benchmark engine.
type LittDBWrapper struct {
	// The database to be benchmarked.
	db litt.DB

	// The table in the database where data is stored.
	table litt.Table
}

// Instantiate a new LittDBWrapper with the given configuration.
func NewLittDBWrapper(cfg *config.BenchmarkConfig) (*LittDBWrapper, error) {
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

	return &LittDBWrapper{
		db:    db,
		table: table,
	}, nil
}

func (w *LittDBWrapper) BuildThreadLocalWrapper() (ThreadLocalDatabaseWrapper, error) {
	return w, nil
}

func (w *LittDBWrapper) Close() error {
	return w.db.Close()
}

func (w *LittDBWrapper) Put(key, value []byte) error {
	return w.table.Put(key, value)
}

func (w *LittDBWrapper) Get(key []byte) (value []byte, exists bool, err error) {
	return w.table.Get(key)
}

func (w *LittDBWrapper) Flush() error {
	return w.table.Flush()
}
