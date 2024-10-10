package tablestore

import (
	"encoding/binary"
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var _ kvstore.Table = &tableView{}

// tableView allows data in a table to be accessed as if it were a Store.
type tableView struct {
	// base is the underlying store.
	base kvstore.Store
	// name is the name of the table.
	name string
	// prefix is the prefix for all keys in the table.
	prefix uint32
}

// NewTableView creates a new view into a table in a New.
func newTableView(
	base kvstore.Store,
	name string,
	prefix uint32) kvstore.Table {

	return &tableView{
		base:   base,
		name:   name,
		prefix: prefix,
	}
}

// Name returns the name of the table.
func (t *tableView) Name() string {
	return t.name
}

// Put adds a key-value pair to the table.
func (t *tableView) Put(key []byte, value []byte) error {
	if value == nil {
		value = []byte{}
	}

	k := t.TableKey(key)
	return t.base.Put(k, value)
}

// Get retrieves a value from the table.
func (t *tableView) Get(key []byte) ([]byte, error) {
	k := t.TableKey(key)
	return t.base.Get(k)
}

// Delete removes a key-value pair from the table.
func (t *tableView) Delete(key []byte) error {
	k := t.TableKey(key)
	return t.base.Delete(k)
}

// iteratorWrapper wraps the base iterator to iterate only over keys in the table.
type iteratorWrapper struct {
	base iterator.Iterator
}

func (i *iteratorWrapper) First() bool {
	return i.base.First()
}

func (i *iteratorWrapper) Last() bool {
	return i.base.Last()
}

func (i *iteratorWrapper) Seek(key []byte) bool {
	return i.base.Seek(key)
}

func (i *iteratorWrapper) Next() bool {
	return i.base.Next()
}

func (i *iteratorWrapper) Prev() bool {
	return i.base.Prev()
}

func (i *iteratorWrapper) Release() {
	i.base.Release()
}

func (i *iteratorWrapper) SetReleaser(releaser util.Releaser) {
	i.base.SetReleaser(releaser)
}

func (i *iteratorWrapper) Valid() bool {
	return i.base.Valid()
}

func (i *iteratorWrapper) Error() error {
	return i.base.Error()
}

func (i *iteratorWrapper) Key() []byte {
	baseKey := i.base.Key()
	return baseKey[4:]
}

func (i *iteratorWrapper) Value() []byte {
	return i.base.Value()
}

// NewIterator creates a new iterator. Only keys prefixed with the given prefix will be iterated.
func (t *tableView) NewIterator(prefix []byte) (iterator.Iterator, error) {
	if prefix == nil {
		prefix = []byte{}
	}

	it, err := t.base.NewIterator(t.TableKey(prefix))
	if err != nil {
		return nil, err
	}

	return &iteratorWrapper{
		base: it,
	}, nil
}

// Shutdown shuts down the table.
func (t *tableView) Shutdown() error {
	return t.base.Shutdown()
}

// Destroy shuts down a table and deletes all data in it.
func (t *tableView) Destroy() error {
	return t.base.Destroy()
}

// tableBatch is a batch for a table in a New.
type tableBatch struct {
	table kvstore.Table
	batch kvstore.StoreBatch
}

// Put schedules a key-value pair to be added to the table.
func (t *tableBatch) Put(key []byte, value []byte) {
	if value == nil {
		value = []byte{}
	}
	k := t.table.TableKey(key)
	t.batch.Put(k, value)
}

// Delete schedules a key-value pair to be removed from the table.
func (t *tableBatch) Delete(key []byte) {
	k := t.table.TableKey(key)
	t.batch.Delete(k)
}

// Apply applies the batch to the table.
func (t *tableBatch) Apply() error {
	return t.batch.Apply()
}

// Size returns the number of operations in the batch.
func (t *tableBatch) Size() uint32 {
	return t.batch.Size()
}

// NewBatch creates a new batch for the table.
func (t *tableView) NewBatch() kvstore.StoreBatch {
	return &tableBatch{
		table: t,
		batch: t.base.NewBatch(),
	}
}

// TableKey creates a key scoped to this table.
func (t *tableView) TableKey(key []byte) kvstore.TableKey {
	result := make([]byte, 4+len(key))
	binary.BigEndian.PutUint32(result, t.prefix)
	copy(result[4:], key)
	return result
}
