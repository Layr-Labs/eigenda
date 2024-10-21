package tablestore

import (
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var _ iterator.Iterator = &tableStoreIterator{}

type tableStoreIterator struct {
	baseIterator iterator.Iterator
	keyBuilder   kvstore.KeyBuilder
}

// newTableStoreIterator creates a new table store iterator that iterates over a table.
func newTableStoreIterator(base kvstore.Store[[]byte], k kvstore.Key) (*tableStoreIterator, error) {

	baseIterator, err := base.NewIterator(k.Raw())
	if err != nil {
		return nil, err
	}

	return &tableStoreIterator{
		baseIterator: baseIterator,
		keyBuilder:   k.Builder(),
	}, nil
}

func (t *tableStoreIterator) First() bool {
	return t.baseIterator.First()
}

func (t *tableStoreIterator) Last() bool {
	return t.baseIterator.Last()
}

func (t *tableStoreIterator) Seek(key []byte) bool {
	return t.baseIterator.Seek(t.keyBuilder.Key(key).Raw())
}

func (t *tableStoreIterator) Next() bool {
	return t.baseIterator.Next()
}

func (t *tableStoreIterator) Prev() bool {
	return t.baseIterator.Prev()
}

func (t *tableStoreIterator) Release() {
	t.baseIterator.Release()
}

func (t *tableStoreIterator) SetReleaser(releaser util.Releaser) {
	t.baseIterator.SetReleaser(releaser)
}

func (t *tableStoreIterator) Valid() bool {
	return t.baseIterator.Valid()
}

func (t *tableStoreIterator) Error() error {
	return t.baseIterator.Error()
}

func (t *tableStoreIterator) Key() []byte {
	baseKey := t.baseIterator.Key()
	return baseKey[prefixLength:]
}

func (t *tableStoreIterator) Value() []byte {
	return t.baseIterator.Value()
}
