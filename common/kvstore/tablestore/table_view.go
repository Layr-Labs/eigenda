package tablestore

import (
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

var _ kvstore.Store = &tableView{}

// tableView allows a single table in a TableStore to be accessed as if it were the only table in the store.
type tableView struct {
	base       kvstore.Store
	keyBuilder kvstore.KeyBuilder
}

// NewTableView creates a new view into a table in a TableStore.
func NewTableView(base kvstore.Store, keyBuilder kvstore.KeyBuilder) kvstore.Store {
	return &tableView{
		base:       base,
		keyBuilder: keyBuilder,
	}
}

func (t *tableView) Put(key []byte, value []byte) error {
	return t.base.Put(t.keyBuilder.Key(key).GetRawBytes(), value)
}

func (t *tableView) Get(key []byte) ([]byte, error) {
	return t.base.Get(t.keyBuilder.Key(key).GetRawBytes())
}

func (t *tableView) Delete(key []byte) error {
	return t.base.Delete(t.keyBuilder.Key(key).GetRawBytes())
}

func (t *tableView) DeleteBatch(keys [][]byte) error {
	prefixedKeys := make([][]byte, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = t.keyBuilder.Key(key).GetRawBytes()
	}
	return t.base.DeleteBatch(prefixedKeys)
}

func (t *tableView) WriteBatch(keys [][]byte, values [][]byte) error {
	prefixedKeys := make([][]byte, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = t.keyBuilder.Key(key).GetRawBytes()
	}
	return t.base.WriteBatch(prefixedKeys, values)
}

func (t *tableView) NewIterator(prefix []byte) (iterator.Iterator, error) {
	return t.base.NewIterator(t.keyBuilder.Key(prefix).GetRawBytes())
}

func (t *tableView) Shutdown() error {
	// Intentional no-op, this is just a view into the table
	return nil
}

func (t *tableView) Destroy() error {
	// Intentional no-op, this is just a view into the table
	return nil
}
