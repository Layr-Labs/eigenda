package wrappers

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/Layr-Labs/eigenda/common/kvstore/tablestore"
	"github.com/Layr-Labs/eigenda/litt/benchmark/config"
	"github.com/Layr-Labs/eigenda/litt/util"
)

var _ DatabaseWrapper = (*LevelDBWrapper)(nil)

type LevelDBWrapper struct {
	store      kvstore.TableStore
	keyBuilder kvstore.KeyBuilder
	ttl        time.Duration
}

var _ ThreadLocalDatabaseWrapper = (*ThreadLocalLevelDBWrapper)(nil)

type ThreadLocalLevelDBWrapper struct {
	store      kvstore.TableStore
	keyBuilder kvstore.KeyBuilder
	batch      kvstore.TTLBatch[kvstore.Key]
	ttl        time.Duration
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
	storeConfig.Schema = []string{"benchmark"}
	storeConfig.Path = &path

	store, err := tablestore.Start(logger, storeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to start tablestore: %w", err)
	}

	keyBuilder, err := store.GetKeyBuilder("benchmark")
	if err != nil {
		return nil, fmt.Errorf("failed to get key builder: %w", err)
	}

	ttl := time.Duration(cfg.TTLHours * float64(time.Hour))

	return &LevelDBWrapper{
		store:      store,
		keyBuilder: keyBuilder,
		ttl:        ttl,
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
	t.batch.PutWithTTL(t.keyBuilder.Key(key), value, t.ttl)
	return nil
}

func (t *ThreadLocalLevelDBWrapper) Get(key []byte) (value []byte, exists bool, err error) {
	// Note that this actually has slightly wrong semantics. This method SHOULD be returning any value previously
	// written by Put, but since this is a simple implementation, it won't actually return values currently waiting
	// inside a batch. Not important to implement this correctly, since the main point of this benchmark suite
	// is to exercise LittDB, not the LevelDB wrapper.

	value, err = t.store.Get(t.keyBuilder.Key(key))
	if err != nil {
		if errors.Is(err, kvstore.ErrNotFound) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("failed to get key %s: %w", key, err)
	}
	return value, true, nil
}

func (t *ThreadLocalLevelDBWrapper) Flush() error {
	if t.batch.Size() == 0 {
		// No changes to flush
		return nil
	}
	err := t.batch.Apply()
	if err != nil {
		return fmt.Errorf("failed to apply batch: %w", err)
	}
	return nil
}
