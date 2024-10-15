package tablestore

import (
	"context"
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"sync"
	"time"
)

var _ kvstore.TableStore = &tableStore{}

// The maximum number of keys to delete in a single batch.
var maxDeletionBatchSize uint32 = 1024

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
	base kvstore.Store

	// A map from table names to tables.
	tableMap map[string]kvstore.Table

	// A map containing expiration times. Keys in this table are made up of a timestamp prepended to a key.
	// The value is an empty byte slice. Iterating over this table will return keys in order of expiration time.
	expirationTable kvstore.Table
}

// wrapper wraps the given Store to create a TableStore.
//
// WARNING: it is not safe to access the wrapped store directly while the TableStore is in use. The TableStore uses
// special key formatting, and direct access to the wrapped store may violate the TableStore's invariants, resulting
// in undefined behavior.
func newTableStore(
	logger logging.Logger,
	base kvstore.Store,
	tableIDMap map[uint32]string,
	expirationTable kvstore.Table) kvstore.TableStore {

	ctx, cancel := context.WithCancel(context.Background())
	waitGroup := &sync.WaitGroup{}

	store := &tableStore{
		ctx:             ctx,
		cancel:          cancel,
		waitGroup:       waitGroup,
		logger:          logger,
		base:            base,
		tableMap:        make(map[string]kvstore.Table),
		expirationTable: expirationTable,
	}

	for prefix, name := range tableIDMap {
		table := newTableView(base, store, name, prefix)
		store.tableMap[name] = table
	}

	store.expireKeysInBackground(time.Minute)

	return store
}

// GetTable gets the table with the given name. If the table does not exist, it is first created.
func (t *tableStore) GetTable(name string) (kvstore.Table, error) {
	table, ok := t.tableMap[name]
	if !ok {
		return nil, kvstore.ErrTableNotFound
	}

	return table, nil
}

// GetTables returns a list of all tables in the store in no particular order.
func (t *tableStore) GetTables() []kvstore.Table {
	tables := make([]kvstore.Table, 0, len(t.tableMap))
	for _, table := range t.tableMap {
		tables = append(tables, table)
	}

	return tables
}

// NewBatch creates a new batch for writing to the store.
func (t *tableStore) NewBatch() kvstore.TableStoreBatch {
	return &tableStoreBatch{
		batch:           t.base.NewBatch(),
		expirationTable: t.expirationTable,
	}
}

// ExpireKeysInBackground spawns a background goroutine that periodically checks for expired keys and deletes them.
func (t *tableStore) expireKeysInBackground(gcPeriod time.Duration) {
	ticker := time.NewTicker(gcPeriod)
	go func() {
		defer t.waitGroup.Done()
		for {
			select {
			case now := <-ticker.C:
				err := t.expireKeys(now)
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
func (t *tableStore) expireKeys(now time.Time) error {
	it, err := t.expirationTable.NewIterator(nil)
	if err != nil {
		return err
	}
	defer it.Release()

	batch := t.NewBatch()

	for it.Next() {
		expiryKey := it.Key()
		expiryTimestamp, baseKey := parsePrependedTimestamp(expiryKey)

		if expiryTimestamp.After(now) {
			// No more values to expire
			break
		}

		batch.Delete(baseKey)
		batch.Delete(expiryKey)

		if batch.Size() >= maxDeletionBatchSize {
			err = batch.Apply()
			if err != nil {
				return err
			}
			batch = t.NewBatch()
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
