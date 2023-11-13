package leveldb

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type keyValueReader struct {
	getFn func(key []byte) ([]byte, error)
}

type keyValueWriter struct {
	putFn func(key, value []byte) error
}

func (r keyValueReader) get(key []byte, value any) error {
	data, err := r.getFn(key)
	if err != nil {
		return err
	}
	return decode(data, value)
}

func (w keyValueWriter) put(key []byte, value any) error {
	data, err := encode(value)
	if err != nil {
		return err
	}
	return w.putFn(key, data)
}

type opener func(path string) (*leveldb.DB, error)

type levelDB struct {
	keyValueReader
	keyValueWriter

	db   *leveldb.DB
	Path string
}

func newLevelDB(path string, opener ...opener) (*levelDB, error) {
	var (
		ldb *leveldb.DB
		err error
	)

	if len(opener) > 0 {
		ldb, err = opener[0](path)
	} else {
		ldb, err = leveldb.OpenFile(path, &opt.Options{Filter: filter.NewBloomFilter(10)})
	}
	if err != nil {
		return nil, err
	}

	db := &levelDB{
		keyValueReader: keyValueReader{
			getFn: func(key []byte) ([]byte, error) {
				return ldb.Get(key, nil)
			},
		},
		keyValueWriter: keyValueWriter{
			putFn: func(key, value []byte) error {
				return ldb.Put(key, value, nil)
			},
		},
		db:   ldb,
		Path: path,
	}
	return db, nil
}

func (l *levelDB) Close() {
	_ = l.db.Close()
}

func (l *levelDB) Get(key []byte, value any) error {
	data, err := l.db.Get(key, nil)
	if err != nil {
		return err
	}
	return decode(data, value)
}

func (l *levelDB) Has(key []byte) bool {
	ok, _ := l.db.Has(key, nil)
	return ok
}

func (l *levelDB) Put(key []byte, value any) error {
	return l.put(key, value)
}

func (l *levelDB) Iter(prefix []byte) *iter {
	it := l.db.NewIterator(util.BytesPrefix(prefix), nil)
	return &iter{it: it}
}

func (l *levelDB) Tx() (*transaction, error) {
	return newTransaction(l)
}

type iter struct {
	it iterator.Iterator
}

func (i *iter) First() bool {
	return i.it.First()
}

func (i *iter) Next() bool {
	return i.it.Next()
}

func (i *iter) Value(v any) error {
	return decode(i.it.Value(), v)
}

func (i *iter) Release() {
	i.it.Release()
}

type transaction struct {
	keyValueReader
	keyValueWriter

	b   *leveldb.Batch
	sn  *leveldb.Snapshot
	db  *leveldb.DB
	err error
}

func newTransaction(l *levelDB) (*transaction, error) {
	sn, err := l.db.GetSnapshot()
	if err != nil {
		return nil, err
	}
	b := new(leveldb.Batch)

	tx := &transaction{
		keyValueReader: keyValueReader{
			getFn: func(key []byte) ([]byte, error) {
				return sn.Get(key, nil)
			},
		},
		keyValueWriter: keyValueWriter{
			putFn: func(key, value []byte) error {
				b.Put(key, value)
				return nil
			},
		},
		b:  b,
		sn: sn,
		db: l.db,
	}
	return tx, nil
}

func (t *transaction) Empty() bool {
	it := t.sn.NewIterator(nil, nil)
	defer it.Release()
	return !it.First()
}

func (t *transaction) Has(key []byte) (bool, error) {
	return t.sn.Has(key, nil)
}

func (t *transaction) Get(key []byte, value any) error {
	if t.err != nil {
		return t.err
	}
	return t.get(key, value)
}

func (t *transaction) Put(key []byte, value any) {
	if t.err != nil {
		return
	}
	t.err = t.put(key, value)
}

func (t *transaction) Iter(prefix []byte) *iter {
	it := t.sn.NewIterator(util.BytesPrefix(prefix), nil)
	return &iter{it: it}
}

func (t *transaction) Delete(key []byte) {
	if t.err != nil {
		return
	}
	t.b.Delete(key)
}

func (t *transaction) Commit() error {
	if t.err != nil {
		return t.err
	}
	t.err = t.db.Write(t.b, nil)
	return t.err
}

func (t *transaction) Discard() {
	t.b.Reset()
	t.err = nil
}

func (t *transaction) SetErr(err error) {
	if t.err == nil {
		t.err = err
	}
}
