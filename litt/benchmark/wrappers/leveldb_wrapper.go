package wrappers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/Layr-Labs/eigenda/common/kvstore/tablestore"
	"github.com/Layr-Labs/eigenda/litt/benchmark/config"
	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var _ DatabaseWrapper = (*LevelDBWrapper)(nil)
var tableName = "benchmark"

type LevelDBWrapper struct {
	store      kvstore.TableStore
	keyBuilder kvstore.KeyBuilder
	ttl        time.Duration
	metrics    *basicWrapperMetrics
}

var _ ThreadLocalDatabaseWrapper = (*ThreadLocalLevelDBWrapper)(nil)

type ThreadLocalLevelDBWrapper struct {
	store      kvstore.TableStore
	keyBuilder kvstore.KeyBuilder
	batch      kvstore.TTLBatch[kvstore.Key]
	ttl        time.Duration
	metrics    *basicWrapperMetrics
}

var _ WrapperFactory = NewLevelDBWrapper

func NewLevelDBWrapper(cfg *config.BenchmarkConfig) (DatabaseWrapper, error) {

	// Level DB can't split data across multiple paths, so we just use the first path in the config.
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

	logger, err := common.NewLogger(common.DefaultConsoleLoggerConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	storeConfig := tablestore.DefaultConfig()
	storeConfig.Schema = []string{tableName}
	storeConfig.Path = &path

	store, err := tablestore.Start(logger, storeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to start tablestore: %w", err)
	}

	keyBuilder, err := store.GetKeyBuilder(tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get key builder: %w", err)
	}

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

		fmt.Printf("Registry: %v\n", registry) // TODO

		metrics = newBasicWrapperMetrics(registry, tableName)
	} else {
		logger.Infof("Metrics are disabled, not starting metrics server, metrics enabled flag = %v", cfg.LittConfig.MetricsEnabled) // TODO
	}

	return &LevelDBWrapper{
		store:      store,
		keyBuilder: keyBuilder,
		ttl:        ttl,
		metrics:    metrics,
	}, nil
}

func (w *LevelDBWrapper) BuildThreadLocalWrapper() (ThreadLocalDatabaseWrapper, error) {
	batch := w.store.NewTTLBatch()

	return &ThreadLocalLevelDBWrapper{
		store:      w.store,
		keyBuilder: w.keyBuilder,
		batch:      batch,
		ttl:        w.ttl,
	}, nil
}

func (t *ThreadLocalLevelDBWrapper) Put(key []byte, value []byte) error {
	start := time.Now()
	t.batch.PutWithTTL(t.keyBuilder.Key(key), value, t.ttl)
	elapsed := time.Since(start)
	t.metrics.RecordWriteLatency(elapsed)
	t.metrics.RecordBytesWritten(uint64(len(key) + len(value)))

	return nil
}

func (t *ThreadLocalLevelDBWrapper) Get(key []byte) (value []byte, exists bool, err error) {
	// Note that this actually has slightly wrong semantics. This method SHOULD be returning any value previously
	// written by Put, but since this is a simple implementation, it won't actually return values currently waiting
	// inside a batch. Not important to implement this correctly, since the main point of this benchmark suite
	// is to exercise LittDB, not the LevelDB wrapper.

	start := time.Now()

	value, err = t.store.Get(t.keyBuilder.Key(key))
	if err != nil {
		if errors.Is(err, kvstore.ErrNotFound) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("failed to get key %s: %w", key, err)
	}

	end := time.Since(start)
	t.metrics.RecordReadLatency(end)
	t.metrics.RecordBytesRead(uint64(len(key) + len(value)))

	return value, true, nil
}

func (t *ThreadLocalLevelDBWrapper) Flush() error {
	if t.batch.Size() == 0 {
		// No changes to flush
		return nil
	}
	start := time.Now()

	err := t.batch.Apply()
	if err != nil {
		return fmt.Errorf("failed to apply batch: %w", err)
	}

	elapsed := time.Since(start)
	t.metrics.RecordFlushLatency(elapsed)

	return nil
}
