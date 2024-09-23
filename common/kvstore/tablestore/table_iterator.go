package tablestore

import (
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var _ iterator.Iterator = &tableIterator{}

// tableIterator is an iterator that iterates over a single table in a TableStore.
type tableIterator struct {
	base       iterator.Iterator
	keyBuilder kvstore.KeyBuilder
}

// NewTableIterator creates a new iterator that iterates over a single table in a TableStore.
func newTableIterator(base iterator.Iterator, keyBuilder kvstore.KeyBuilder) iterator.Iterator {
	return &tableIterator{
		base:       base,
		keyBuilder: keyBuilder,
	}
}

func (t *tableIterator) First() bool {
	return t.base.First()
}

func (t *tableIterator) Last() bool {
	return t.base.Last()
}

func (t *tableIterator) Seek(key []byte) bool {
	return t.base.Seek(t.keyBuilder.Key(key).GetRawBytes())
}

func (t *tableIterator) Next() bool {
	return t.base.Next()
}

func (t *tableIterator) Prev() bool {
	return t.base.Prev()
}

func (t *tableIterator) Release() {
	t.base.Release()
}

func (t *tableIterator) SetReleaser(releaser util.Releaser) {
	t.base.SetReleaser(releaser)
}

func (t *tableIterator) Valid() bool {
	return t.base.Valid()
}

func (t *tableIterator) Error() error {
	return t.base.Error()
}

func (t *tableIterator) Key() []byte {
	rawKey := t.base.Key()
	return rawKey[4:]
}

func (t *tableIterator) Value() []byte {
	return t.base.Value()
}
