package wrappers

// This file can be used to benchmark BadgerDB performance. Since we don't use BadgerDB outside of this benchmark suite,
// we don't want to add it as a dependency, so by default we comment out this code. To benchmark BadgerDB, do the
// following:
//
// 1. Run "go get github.com/dgraph-io/badger/v4" to add BadgerDB as a dependency. (Don't merge this into master!)
// 2. Uncomment the code in this file.
// 3. Uncomment the "badgerdb" entry in the WrapperFactories map in wrappers/wrapper_factories.go.
// 4. Run the benchmark suite with the "badgerdb" database type in the configuration file.

// Add a single backslash to the line below to uncomment the code in this file.
//*

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/litt/benchmark/config"
	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/dgraph-io/badger/v4"
	"github.com/dgraph-io/badger/v4/options"
	"github.com/docker/go-units"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var _ DatabaseWrapper = (*badgerDBWrapper)(nil)

type badgerDBWrapper struct {
	db      *badger.DB
	metrics *basicWrapperMetrics
	ttl     time.Duration
}

var _ ThreadLocalDatabaseWrapper = (*threadLocalBadgerDBWrapper)(nil)

type threadLocalBadgerDBWrapper struct {
	db          *badger.DB
	transaction *badger.Txn
	metrics     *basicWrapperMetrics
	ttl         time.Duration
}

var _ WrapperFactory = NewBadgerDBWrapper

func NewBadgerDBWrapper(cfg *config.BenchmarkConfig) (DatabaseWrapper, error) {
	// Badger DB can't split data across multiple paths, so we just use the first path in the config.
	path, err := util.SanitizePath(cfg.LittConfig.Paths[0])
	if err != nil {
		return nil, fmt.Errorf("error sanitizing path: %w", err)
	}

	exists, err := util.Exists(path)
	if err != nil {
		return nil, fmt.Errorf("error checking for directory: %w", err)
	}

	if !exists {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return nil, fmt.Errorf("error creating directory: %w", err)
		}
	}

	opts := badger.DefaultOptions(path)

	// The following settings were derived by trial and error. By default, BadgerDB doesn't want to do
	// garbage collection (i.e. compaction) at runtime. If left to its own devices, it will only
	// compact when the database is closed. This is not suitable: we need a database that can run in stable
	// state indefinitely without needing to be closed and reopened. BadgerDB doesn't make this easy.
	// The settings below are how I managed to get it to most reliably do compaction at runtime (although
	// eventually it will become overwhelmed and fail to keep up with compaction, at which point the disk
	// will fill up and the database will stop working).
	opts.Compression = options.None
	opts.Logger = nil
	opts.SyncWrites = true
	opts.ValueThreshold = 0
	opts.BaseTableSize = 10 * units.KiB
	opts.TableSizeMultiplier = 2
	opts.BaseLevelSize = 10 * units.KiB
	opts.LevelSizeMultiplier = 2
	opts.MemTableSize = 10 * units.KiB
	opts.NumMemtables = 1
	opts.NumLevelZeroTables = 1
	opts.MaxLevels = 16

	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("error opening badger database: %w", err)
	}

	logger := cfg.LittConfig.Logger

	// Force the DB to run compaction in the background.
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for {
			<-ticker.C

			logger.Info("Running GC")
			startTime := time.Now()

			err = db.Flatten(8)
			if err != nil {
				logger.Errorf("Error flattening DB: %v", err)
			}

			gcIterations := 0
			for {
				gcIterations++
				err = db.RunValueLogGC(0.125)
				if err != nil {
					if !strings.Contains(err.Error(), "Value log GC attempt didn't result in any cleanup") {
						logger.Errorf("\nError running GC: %v\n", err)
					}
					break
				}
			}

			logger.Infof("\nGC took %v, did %d iterations\n", time.Since(startTime), gcIterations)
		}
	}()

	ttl := time.Duration(cfg.TTLHours * float64(time.Hour))

	var metrics *basicWrapperMetrics
	if cfg.LittConfig.MetricsEnabled {
		registry := prometheus.NewRegistry()
		registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
		registry.MustRegister(collectors.NewGoCollector())

		logger.Infof("Starting metrics server at port %d", cfg.LittConfig.MetricsPort)
		addr := fmt.Sprintf(":%d", cfg.LittConfig.MetricsPort)
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.HandlerFor(
			registry,
			promhttp.HandlerOpts{},
		))
		server := &http.Server{
			Addr:    addr,
			Handler: mux,
		}

		go func() {
			err := server.ListenAndServe()
			if err != nil && !strings.Contains(err.Error(), "http: Server closed") {
				logger.Errorf("metrics server error: %v", err)
			}
		}()

		metrics = newBasicWrapperMetrics(registry, tableName)
	}

	return &badgerDBWrapper{
		db:      db,
		metrics: metrics,
		ttl:     ttl,
	}, nil
}

func (w *badgerDBWrapper) BuildThreadLocalWrapper() (ThreadLocalDatabaseWrapper, error) {
	return &threadLocalBadgerDBWrapper{
		db:          w.db,
		transaction: w.db.NewTransaction(true),
		metrics:     w.metrics,
		ttl:         w.ttl,
	}, nil
}

func (t *threadLocalBadgerDBWrapper) Put(key, value []byte) error {
	start := time.Now()

	entry := badger.NewEntry(key, value).WithTTL(t.ttl)
	err := t.transaction.SetEntry(entry)
	if err != nil {
		if strings.Contains(err.Error(), "Txn is too big to fit into one request") {
			err = t.Flush()
			if err != nil {
				return fmt.Errorf("error flushing transaction: %w", err)
			}
			start = time.Now() // reset the timer so as to not include flush time
			err = t.transaction.SetEntry(entry)
			if err != nil {
				return fmt.Errorf("error setting key entry after flush: %w", err)
			}
		} else {
			return fmt.Errorf("error setting key entry: %w", err)
		}
	}

	elapsed := time.Since(start)
	t.metrics.RecordWriteLatency(elapsed)
	t.metrics.RecordBytesWritten(uint64(len(key) + len(value)))

	return nil
}

func (t *threadLocalBadgerDBWrapper) Get(key []byte) (value []byte, exists bool, err error) {
	// Note that this actually has slightly wrong semantics. This method SHOULD be returning any value previously
	// written by Put, but since this is a simple implementation, it won't actually return values currently waiting
	// inside a batch. Not important to implement this correctly, since the main point of this benchmark suite
	// is to exercise LittDB, not the LevelDB wrapper.

	start := time.Now()

	readTransaction := t.db.NewTransaction(false)
	defer readTransaction.Discard()
	item, err := readTransaction.Get(key)
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("failed to get key %s: %w", string(key), err)
	}

	value, err = item.ValueCopy(nil)
	if err != nil {
		return nil, false, fmt.Errorf("failed to copy value for key %s: %w", string(key), err)
	}

	end := time.Since(start)
	t.metrics.RecordReadLatency(end)
	t.metrics.RecordBytesRead(uint64(len(key) + len(value)))

	return value, true, nil
}

func (t *threadLocalBadgerDBWrapper) Flush() error {
	start := time.Now()

	err := t.transaction.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	elapsed := time.Since(start)
	t.metrics.RecordFlushLatency(elapsed)

	t.transaction = t.db.NewTransaction(true)

	return nil
}

// Leave this here: it allows the entire file to be commented out easily.
//*/
