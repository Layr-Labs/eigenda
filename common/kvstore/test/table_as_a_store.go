package test

import (
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

var _ kvstore.Store[[]byte] = &tableAsAStore{}

// tableAsAStore allows a TableStore to masquerade as a Store. Useful for reusing unit tests that operate on a Store
// against a TableStore. Base TableStore is required to have a table called "test".
type tableAsAStore struct {
	tableStore kvstore.TableStore
	kb         kvstore.KeyBuilder
}

func NewTableAsAStore(tableStore kvstore.TableStore) (kvstore.Store[[]byte], error) {
	kb, err := tableStore.GetKeyBuilder("test")
	if err != nil {
		return nil, err
	}

	return &tableAsAStore{
		tableStore: tableStore,
		kb:         kb,
	}, nil
}

func (t *tableAsAStore) Put(k []byte, value []byte) error {
	return t.tableStore.Put(t.kb.Key(k), value)
}

func (t *tableAsAStore) Get(k []byte) ([]byte, error) {
	return t.tableStore.Get(t.kb.Key(k))
}

func (t *tableAsAStore) Delete(k []byte) error {
	return t.tableStore.Delete(t.kb.Key(k))
}

func (t *tableAsAStore) NewBatch() kvstore.Batch[[]byte] {
	return &wrappedTableStoreBatch{
		tableStoreBatch: t.tableStore.NewBatch(),
		kb:              t.kb,
	}
}

func (t *tableAsAStore) NewIterator(prefix []byte) (iterator.Iterator, error) {
	return t.tableStore.NewIterator(t.kb.Key(prefix))
}

func (t *tableAsAStore) Shutdown() error {
	return t.tableStore.Shutdown()
}

func (t *tableAsAStore) Destroy() error {
	return t.tableStore.Destroy()
}

var _ kvstore.Batch[[]byte] = &wrappedTableStoreBatch{}

type wrappedTableStoreBatch struct {
	tableStoreBatch kvstore.Batch[kvstore.Key]
	kb              kvstore.KeyBuilder
}

func (w *wrappedTableStoreBatch) Put(k []byte, value []byte) {
	w.tableStoreBatch.Put(w.kb.Key(k), value)
}

func (w *wrappedTableStoreBatch) Delete(k []byte) {
	w.tableStoreBatch.Delete(w.kb.Key(k))
}

func (w *wrappedTableStoreBatch) Apply() error {
	return w.tableStoreBatch.Apply()
}

func (w *wrappedTableStoreBatch) Size() uint32 {
	return w.tableStoreBatch.Size()
}
