package tablestore

import (
	"context"
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"sync"
	"time"
)

var _ kvstore.TableStore = &tableStore{}

// tableStore is an implementation of TableStore that wraps a Store.
type tableStore struct {
	// the context for the store
	ctx context.Context

	// the cancel function for the store
	cancel context.CancelFunc

	// this wait group is completed when the garbage collection goroutine is done
	waitGroup *sync.WaitGroup

	// the logger for the store
	logger logging.Logger

	// A base store implementation that this TableStore wraps.
	base kvstore.Store[[]byte]

	// A map from table names to key builders.
	keyBuilderMap map[string]kvstore.KeyBuilder

	// A key builder for the expiration table. Keys in this table are made up of a timestamp prepended to a key.
	// The value is an empty byte slice. Iterating over this table will return keys in order of expiration time.
	expirationKeyBuilder kvstore.KeyBuilder
}

// wrapper wraps the given Store to create a TableStore.
//
// WARNING: it is not safe to access the wrapped store directly while the TableStore is in use. The TableStore uses
// special key formatting, and direct access to the wrapped store may violate the TableStore's invariants, resulting
// in undefined behavior.
func newTableStore(
	logger logging.Logger,
	base kvstore.Store[[]byte],
	tableIDMap map[uint32]string,
	expirationKeyBuilder kvstore.KeyBuilder,
	gcEnabled bool,
	gcPeriod time.Duration,
	gcBatchSize uint32) kvstore.TableStore {

	ctx, cancel := context.WithCancel(context.Background())
	waitGroup := &sync.WaitGroup{}

	store := &tableStore{
		ctx:                  ctx,
		cancel:               cancel,
		waitGroup:            waitGroup,
		logger:               logger,
		base:                 base,
		keyBuilderMap:        make(map[string]kvstore.KeyBuilder),
		expirationKeyBuilder: expirationKeyBuilder,
	}

	for prefix, name := range tableIDMap {
		store.keyBuilderMap[name] = newKeyBuilder(name, prefix)
	}

	if gcEnabled {
		store.expireKeysInBackground(gcPeriod, gcBatchSize)
	}

	return store
}

// GetKeyBuilder gets the key builder for a particular table. Returns ErrTableNotFound if the table does not exist.
func (t *tableStore) GetKeyBuilder(name string) (kvstore.KeyBuilder, error) {
	table, ok := t.keyBuilderMap[name]
	if !ok {
		return nil, kvstore.ErrTableNotFound
	}

	return table, nil
}

// GetKeyBuilders returns a list of all tables in the store in no particular order.
func (t *tableStore) GetKeyBuilders() []kvstore.KeyBuilder {
	tables := make([]kvstore.KeyBuilder, 0, len(t.keyBuilderMap))
	for _, kb := range t.keyBuilderMap {
		tables = append(tables, kb)
	}

	return tables
}

// GetTables returns a list of all tables in the store in no particular order.
func (t *tableStore) GetTables() []string {
	tables := make([]string, 0, len(t.keyBuilderMap))
	for _, kb := range t.keyBuilderMap {
		tables = append(tables, kb.TableName())
	}

	return tables
}

// NewBatch creates a new batch for writing to the store.
func (t *tableStore) NewBatch() kvstore.Batch[kvstore.Key] {
	return newTableStoreBatch(t.base, t.expirationKeyBuilder)
}

// NewTTLBatch creates a new batch for writing to the store with TTLs.
func (t *tableStore) NewTTLBatch() kvstore.TTLBatch[kvstore.Key] {
	return newTableStoreBatch(t.base, t.expirationKeyBuilder)
}

// Put adds a key-value pair to the store.
func (t *tableStore) Put(k kvstore.Key, value []byte) error {
	return t.base.Put(k.Raw(), value)
}

// Get retrieves the value for a key from the store.
func (t *tableStore) Get(k kvstore.Key) ([]byte, error) {
	return t.base.Get(k.Raw())
}

// Delete removes a key from the store.
func (t *tableStore) Delete(k kvstore.Key) error {
	return t.base.Delete(k.Raw())
}

// PutWithTTL adds a key-value pair to the store that expires after a specified duration.
func (t *tableStore) PutWithTTL(key kvstore.Key, value []byte, ttl time.Duration) error {
	batch := t.NewTTLBatch()
	batch.PutWithTTL(key, value, ttl)
	return batch.Apply()
}

// PutWithExpiration adds a key-value pair to the store that expires at a specified time.
func (t *tableStore) PutWithExpiration(key kvstore.Key, value []byte, expiryTime time.Time) error {
	batch := t.NewTTLBatch()
	batch.PutWithExpiration(key, value, expiryTime)
	return batch.Apply()
}

// NewIterator returns an iterator that can be used to iterate over a subset of the keys in the store.
func (t *tableStore) NewIterator(prefix kvstore.Key) (iterator.Iterator, error) {
	return newTableStoreIterator(t.base, prefix)
}

// NewTableIterator returns an iterator that can be used to iterate over all keys in a table.
func (t *tableStore) NewTableIterator(builder kvstore.KeyBuilder) (iterator.Iterator, error) {
	return newTableStoreIterator(t.base, builder.Key([]byte{}))
}

// ExpireKeysInBackground spawns a background goroutine that periodically checks for expired keys and deletes them.
func (t *tableStore) expireKeysInBackground(gcPeriod time.Duration, gcBatchSize uint32) {
	t.waitGroup.Add(1)
	ticker := time.NewTicker(gcPeriod)
	go func() {
		defer t.waitGroup.Done()
		for {
			select {
			case now := <-ticker.C:
				err := t.expireKeys(now, gcBatchSize)
				if err != nil {
					t.logger.Error("Error expiring keys", err)
					continue
				}
			case <-t.ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

// Delete all keys with a TTL that has expired.
func (t *tableStore) expireKeys(now time.Time, gcBatchSize uint32) error {
	it, err := t.NewTableIterator(t.expirationKeyBuilder)
	if err != nil {
		return err
	}
	defer it.Release()

	batch := t.base.NewBatch()

	for it.Next() {
		expiryKey := it.Key()
		expiryTimestamp, baseKey := parsePrependedTimestamp(expiryKey)

		if expiryTimestamp.After(now) {
			// No more values to expire
			break
		}

		batch.Delete(baseKey)
		batch.Delete(expiryKey)

		if batch.Size() >= gcBatchSize {
			err = batch.Apply()
			if err != nil {
				return err
			}
			batch = t.base.NewBatch()
		}
	}

	if batch.Size() > 0 {
		return batch.Apply()
	}
	return nil
}

// Shutdown shuts down the store, flushing any remaining cached data to disk.
func (t *tableStore) Shutdown() error {
	t.cancel()
	t.waitGroup.Wait()
	return t.base.Shutdown()
}

// Destroy shuts down and permanently deletes all data in the store.
func (t *tableStore) Destroy() error {
	t.cancel()
	t.waitGroup.Wait()
	return t.base.Destroy()
}
