package tablestore

import (
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

// Table ID 0 is reserved for use internal use by the metadata table.
const metadataTableID = 0

// Table ID 1 is reserved for use by the namespace table. This stores a mapping between table names and table IDs. TODO
const namespaceTableID = 1

// This key is used to store the maximum table count in the metadata table.
const maxTableCountKey = "max_table_count"

var _ kvstore.TableStore = &tableStore{}
var _ kvstore.Key = &key{}
var _ kvstore.KeyBuilder = &keyBuilder{}

// tableStore is an implementation of TableStore that wraps a Store.
type tableStore struct {
	// A base store implementation that this TableStore wraps.
	base kvstore.Store

	// A map from table names to table IDs.
	tableMap map[string]uint64
}

// key is a key in a TableStore.
type key struct {
	tableID uint64
	key     []byte
}

// GetKey returns the key within the table.
func (k *key) GetKey() []byte {
	return k.key
}

// GetInternalRepresentation gets the representation of the key as used internally by the store.
func (k *key) GetInternalRepresentation() []byte {
	//TODO implement me
	panic("implement me")
}

// keyBuilder is used to create new keys in a specific table.
type keyBuilder struct {
	tableID uint64
}

// Key creates a new key in a specific table using the given key bytes.
func (k *keyBuilder) Key(keyBytes []byte) kvstore.Key {
	return &key{
		tableID: k.tableID,
		key:     keyBytes,
	}
}

// TableStoreWrapper wraps the given Store to create a TableStore.
//
// Note that the max table count cannot be changed once the TableStore is created. Use the value 0 to use the default
// table size, which is (2^8 - 2) for new tables, or equal to the previous maximum table size if the TableStore is
// reloaded from disk. If the need arises, we may need to write migration code to resize the maximum number of tables,
// but this feature is not currently supported.
func TableStoreWrapper(base kvstore.Store, maxTableCount uint64) kvstore.TableStore {
	return &tableStore{base: base}
}

func (t *tableStore) GetOrCreateTable(name string) (kvstore.KeyBuilder, error) {
	//TODO implement me
	panic("implement me")
}

func (t *tableStore) DropTable(name string) error {
	//TODO implement me
	panic("implement me")
}

func (t *tableStore) ResizeMaxTables(newCount uint64) error {
	//TODO implement me
	panic("implement me")
}

func (t *tableStore) GetKeyBuilder(name string) (kvstore.KeyBuilder, error) {
	//TODO implement me
	panic("implement me")
}

func (t *tableStore) GetMaxTableCount() uint64 {
	//TODO implement me
	panic("implement me")
}

func (t *tableStore) GetCurrentTableCount() uint64 {
	//TODO implement me
	panic("implement me")
}

func (t *tableStore) GetTables() ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (t *tableStore) Put(key kvstore.Key, value []byte) error {
	//TODO implement me
	panic("implement me")
}

func (t *tableStore) Get(key kvstore.Key) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (t *tableStore) Delete(key kvstore.Key) error {
	//TODO implement me
	panic("implement me")
}

func (t *tableStore) DeleteBatch(keys []kvstore.Key) error {
	//TODO implement me
	panic("implement me")
}

func (t *tableStore) WriteBatch(keys []kvstore.Key, values [][]byte) error {
	//TODO implement me
	panic("implement me")
}

func (t *tableStore) NewIterator(prefix kvstore.Key) (iterator.Iterator, error) {
	//TODO implement me
	panic("implement me")
}

func (t *tableStore) Shutdown() error {
	//TODO implement me
	panic("implement me")
}

func (t *tableStore) Destroy() error {
	//TODO implement me
	panic("implement me")
}
