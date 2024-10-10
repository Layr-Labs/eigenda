package tablestore

import (
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

var _ kvstore.TableStore = &tableStore{}

// tableStore is an implementation of TableStore that wraps a Store.
type tableStore struct {
	logger logging.Logger

	// A base store implementation that this TableStore wraps.
	base kvstore.Store

	// A map from table names to tables.
	tableMap map[string]kvstore.Table
}

// wrapper wraps the given Store to create a TableStore.
//
// WARNING: it is not safe to access the wrapped store directly while the TableStore is in use. The TableStore uses
// special key formatting, and direct access to the wrapped store may violate the TableStore's invariants, resulting
// in undefined behavior.
func newTableStore(
	logger logging.Logger,
	base kvstore.Store,
	tables map[string]kvstore.Table) kvstore.TableStore {

	return &tableStore{
		logger:   logger,
		base:     base,
		tableMap: tables,
	}
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
func (t *tableStore) NewBatch() kvstore.TableBatch {
	return &tableStoreBatch{
		store: t,
		batch: t.base.NewBatch(),
	}
}

// Shutdown shuts down the store, flushing any remaining cached data to disk.
func (t *tableStore) Shutdown() error {
	return t.base.Shutdown()
}

// Destroy shuts down and permanently deletes all data in the store.
func (t *tableStore) Destroy() error {
	return t.base.Destroy()
}
