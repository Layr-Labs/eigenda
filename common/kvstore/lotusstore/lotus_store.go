package lotusstore

import (
	"fmt"
	"os"
	"strings"

	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/docker/go-units"
	lotus "github.com/lotusdblabs/lotusdb/v2"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var _ kvstore.Store[[]byte] = &lotusStore{}

type lotusStore struct {
	db           *lotus.DB
	dataDir      string
	batchOptions lotus.BatchOptions
}

func NewStore(dataDir string) (kvstore.Store[[]byte], error) {

	fmt.Printf("starting lotus store at %s\n", dataDir) // TODO

	opts := lotus.DefaultOptions
	opts.DirPath = dataDir
	opts.Sync = true
	opts.MemtableSize = 512 * units.MiB

	db, err := lotus.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("unable to open database: %w", err)
	}

	batchOptions := lotus.BatchOptions{
		WriteOptions: lotus.WriteOptions{
			Sync: true,
		},
	}

	return &lotusStore{
		db:           db,
		dataDir:      dataDir,
		batchOptions: batchOptions,
	}, nil
}

func (l *lotusStore) Put(k []byte, value []byte) error {
	err := l.db.Put(k, value)
	if err != nil {
		return fmt.Errorf("unable to put: %w", err)
	}

	// this is necessary to ensure data durability in the event of a crash
	err = l.db.Sync()
	if err != nil {
		return fmt.Errorf("unable to sync: %w", err)
	}

	return nil
}

func (l *lotusStore) Get(k []byte) ([]byte, error) {
	data, err := l.db.Get(k)
	if err != nil {
		if strings.Contains(err.Error(), "key not found in database") {
			return nil, kvstore.ErrNotFound
		}
	}

	return data, nil
}

func (l *lotusStore) Delete(k []byte) error {
	return l.db.Delete(k)
}

func (l *lotusStore) NewBatch() kvstore.Batch[[]byte] {
	return &lotusBatch{
		batch: l.db.NewBatch(l.batchOptions),
	}
}

func (l *lotusStore) NewIterator(prefix []byte) (iterator.Iterator, error) {
	it, err := l.db.NewIterator(lotus.IteratorOptions{
		Prefix: prefix,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create iterator: %w", err)
	}

	return &lotusIterator{
		iterator: it,
	}, nil
}

func (l *lotusStore) Shutdown() error {
	return l.db.Close()
}

func (l *lotusStore) Destroy() error {
	err := l.Shutdown()
	if err != nil {
		return fmt.Errorf("unable to shutdown: %w", err)
	}

	err = os.RemoveAll(l.dataDir)
	if err != nil {
		return fmt.Errorf("unable to remove data directory %s: %w", l.dataDir, err)
	}

	return nil
}

var _ kvstore.Batch[[]byte] = &lotusBatch{}

type lotusBatch struct {
	batch     *lotus.Batch
	batchSize uint32
}

func (b *lotusBatch) Put(key []byte, value []byte) {
	b.batchSize++
	err := b.batch.Put(key, value)
	if err != nil {
		panic(fmt.Errorf("unable to put: %w", err)) // TODO handle this better if we actually release this
	}
}

func (b *lotusBatch) Delete(key []byte) {
	b.batchSize++
	err := b.batch.Delete(key)
	if err != nil {
		panic(fmt.Errorf("unable to delete: %w", err)) // TODO handle this better if we actually release this
	}
}

func (b *lotusBatch) Apply() error {
	return b.batch.Commit()
}

func (b *lotusBatch) Size() uint32 {
	return b.batchSize
}

var _ iterator.Iterator = &lotusIterator{}

type lotusIterator struct {
	iterator *lotus.Iterator
}

func (i *lotusIterator) First() bool {
	//TODO implement me
	panic("implement me")
}

func (i *lotusIterator) Last() bool {
	//TODO implement me
	panic("implement me")
}

func (i *lotusIterator) Seek(key []byte) bool {
	//TODO implement me
	panic("implement me")
}

func (i *lotusIterator) Next() bool { // TODO
	i.iterator.Next()
	return i.iterator.Valid()
}

func (i *lotusIterator) Prev() bool {
	//TODO implement me
	panic("implement me")
}

func (i *lotusIterator) Release() { // TODO
	err := i.iterator.Close()
	if err != nil {
		panic(fmt.Errorf("unable to close iterator: %w", err)) // TODO handle this better if we actually release this
	}
}

func (i *lotusIterator) SetReleaser(releaser util.Releaser) {
	//TODO implement me
	panic("implement me")
}

func (i *lotusIterator) Valid() bool {
	//TODO implement me
	panic("implement me")
}

func (i *lotusIterator) Error() error {
	//TODO implement me
	panic("implement me")
}

func (i *lotusIterator) Key() []byte { // TODO
	return i.iterator.Key()
}

func (i *lotusIterator) Value() []byte {
	//TODO implement me
	panic("implement me")
}
