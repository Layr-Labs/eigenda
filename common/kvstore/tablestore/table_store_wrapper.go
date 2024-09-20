package tablestore

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"math"
)

// Table ID 0 is reserved for use internal use by the metadata table.
const metadataTableID = 0

// Table ID 1 is reserved for use by the namespace table. This stores a mapping between table names and table IDs.
const namespaceTableID = 1

// The first table ID that can be used for user-created tables.
const firstTableID = 2

// This key is used to store the maximum table count in the metadata table.
const maxTableCountKey = "max_table_count"

// This key is used to store the schema version in the metadata table.
const schemaVersionKey = "schema_version"

// The current schema version of the metadata table.
const currentSchemaVersion uint64 = 0

var _ kvstore.TableStore = &tableStore{}
var _ kvstore.Key = &key{}
var _ kvstore.KeyBuilder = &keyBuilder{}

var ERR_TABLE_LIMIT_EXCEEDED = errors.New("table limit exceeded")

// tableStore is an implementation of TableStore that wraps a Store.
type tableStore struct {
	// A base store implementation that this TableStore wraps.
	base kvstore.Store

	// A map from table names to table IDs.
	tableMap map[string]uint64

	// A map from table IDs to table prefixes.
	prefixMap map[uint64][]byte

	// The highest ID of all tables in the store.
	highestTableID uint64

	// The maximum number of tables that can be created in the store without resizing. This count includes
	// the two reserved tables for metadata and namespace.
	maxTableCount uint64

	// Builds keys for the metadata table.
	metadataTable kvstore.KeyBuilder

	// Builds keys for the namespace table.
	namespaceTable kvstore.KeyBuilder
}

// key is a key in a TableStore.
type key struct {
	// the prefix for the table
	prefix []byte
	// the key within the table
	key []byte
}

// GetKeyBytes returns the key within the table, interpreted as a byte slice.
func (k *key) GetKeyBytes() []byte {
	return k.key
}

// GetKeyString returns the key within the table, interpreted as a string. Calling this
// method on keys that do not represent a string may return odd results.
func (k *key) GetKeyString() string {
	return string(k.key)
}

// GetKeyUint64 returns the key within the table, interpreted as a uint64. Calling this
// method on keys that do not represent a uint64 may return odd results.
func (k *key) GetKeyUint64() uint64 {
	if len(k.key) != 8 {
		return binary.BigEndian.Uint64(k.key)
	} else if len(k.key) == 0 {
		return 0
	} else if len(k.key) < 8 {
		slice := make([]byte, 8)
		copy(slice[8-len(k.key):], k.key)
		return binary.BigEndian.Uint64(slice)
	} else {
		return binary.BigEndian.Uint64(k.key[:8])
	}
}

// GetRawBytes gets the representation of the key as used internally by the store.
func (k *key) GetRawBytes() []byte {
	return append(k.prefix, k.key...)
}

// keyBuilder is used to create new keys in a specific table.
type keyBuilder struct {
	// the prefix for the table
	prefix []byte
}

// Key creates a new key in a specific table using the given key bytes.
func (k *keyBuilder) Key(keyBytes []byte) kvstore.Key {
	return &key{
		prefix: k.prefix,
		key:    keyBytes,
	}
}

// StringKey creates a new key in a specific table using the given key string.
func (k *keyBuilder) StringKey(keyString string) kvstore.Key {
	return &key{
		prefix: k.prefix,
		key:    []byte(keyString),
	}
}

// Uint64Key creates a new key in a specific table using the given uint64 as a key.
func (k *keyBuilder) Uint64Key(uKey uint64) kvstore.Key {
	keyBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(keyBytes, uKey)
	return &key{
		prefix: k.prefix,
		key:    keyBytes,
	}
}

// TODO consider what happens if we crash during metadata operations. Probably needs to be batched.

// TableStoreWrapper wraps the given Store to create a TableStore.
//
// Note that the max table count cannot be changed once the TableStore is created. Use the value 0 to use the default
// table size, which is (2^8 - 2) for new tables, or equal to the previous maximum table size if the TableStore is
// reloaded from disk. If the need arises, we may need to write migration code to resize the maximum number of tables,
// but this feature is not currently supported.
func TableStoreWrapper(
	base kvstore.Store,
	maxTableCount uint64) (kvstore.TableStore, error) {

	prefixMap := make(map[uint64][]byte, maxTableCount)
	// TODO we don't know the prefix length without reading from disk... circular dependency?
	prefixLength := getPrefixLength(maxTableCount)

	// Set up metadata table.
	metadataPrefix, err := getPrefix(metadataTableID, prefixLength)
	if err != nil {
		return nil, fmt.Errorf("error getting prefix for metadata table: %w", err)
	}
	prefixMap[metadataTableID] = metadataPrefix
	metadataTable := &keyBuilder{
		prefix: metadataPrefix,
	}

	// Set up namespace table.
	namespacePrefix, err := getPrefix(namespaceTableID, prefixLength)
	if err != nil {
		return nil, fmt.Errorf("error getting prefix for namespace table: %w", err)
	}
	prefixMap[namespaceTableID] = namespacePrefix
	namespaceTable := &keyBuilder{
		prefix: namespacePrefix,
	}

	schemaKey := metadataTable.StringKey(schemaVersionKey).GetRawBytes()
	countKey := metadataTable.StringKey(maxTableCountKey).GetRawBytes()

	onDiskSchemaBytes, err := base.Get(schemaKey)
	if errors.Is(err, kvstore.ErrNotFound) {
		// This store is new, no on disk schema version exists.
		err = base.Put(schemaKey, []byte{byte(currentSchemaVersion)})
		if err != nil {
			return nil, fmt.Errorf("error setting schema version in metadata table: %w", err)
		}

		if maxTableCount == 0 {
			// Use the default value.
			maxTableCount = uint64(math.Pow(2, 8))
		}

		err = base.Put(countKey, []byte{byte(maxTableCount)})
		if err != nil {
			return nil, fmt.Errorf("error setting max table count in metadata table: %w", err)
		}
	} else {
		// This store is not new. Load data from disk.

		// Verify schema version.
		onDiskSchema := binary.BigEndian.Uint64(onDiskSchemaBytes)
		if onDiskSchema != currentSchemaVersion {
			// In the future if we change schema versions, we may need to write migration code here.
			return nil, fmt.Errorf(
				"incompatible schema version: code is at version %d, data on disk is at version %d",
				currentSchemaVersion, onDiskSchema)
		}

		// Verify max table count.
		onDiskMaxTableCountBytes, err := base.Get(metadataTable.StringKey(maxTableCountKey).GetRawBytes())
		if err != nil {
			return nil, fmt.Errorf("error reading max table count from metadata table: %w", err)
		}

		onDiskMaxTableCount := binary.BigEndian.Uint64(onDiskMaxTableCountBytes)
		if maxTableCount == 0 {
			// Use the on disk value.
			maxTableCount = onDiskMaxTableCount
		} else if maxTableCount != onDiskMaxTableCount {
			return nil, fmt.Errorf(
				"max table count mismatch: code expects %d, data on disk is %d",
				maxTableCount, onDiskMaxTableCount)
		}

		// Load namespace table.
		iterator, err := base.NewIterator(namespacePrefix)
		if err != nil {
			return nil, fmt.Errorf("error creating namespace table iterator: %w", err)
		}
		defer iterator.Release()

	}

	store := &tableStore{
		base:           base,
		tableMap:       make(map[string]uint64, maxTableCount),
		prefixMap:      prefixMap,
		highestTableID: firstTableID - 1,
		maxTableCount:  maxTableCount + 2, /* 2 tables are reserved for internal use */
		metadataTable:  metadataTable,
		namespaceTable: namespaceTable,
	}

	// TODO other stuff

	return store, nil
}

// Get the length of the prefix required to support the given number of tables. This method does not
// give special treatment to the metadata table or the namespace table.
func getPrefixLength(maxTableCount uint64) uint64 {
	log2 := math.Log2(float64(maxTableCount))

	// Round up to the nearest integer, this is the number of bits we need to represent the table ID.
	bits := math.Ceil(log2)

	// Round up to the nearest multiple of 8, this is the number of bytes we need to represent the table ID.
	return uint64(math.Ceil(bits / 8))
}

// Get the prefix for the given table ID and prefix length.
func getPrefix(tableID uint64, prefixLength uint64) ([]byte, error) {
	if prefixLength > 8 {
		return nil, errors.New("prefix length greater than 8 is not supported")
	}

	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, tableID)
	return bytes[8-prefixLength:], nil
}

// GetOrCreateTable creates a new table with the given name if one does not exist.
func (t *tableStore) GetOrCreateTable(name string) (kvstore.KeyBuilder, error) {
	tableID, ok := t.tableMap[name]
	if ok {
		// Table already exists.
		prefix := t.prefixMap[tableID]
		return &keyBuilder{
			prefix: prefix,
		}, nil
	}

	currentTableCount := uint64(len(t.tableMap))
	if currentTableCount == t.maxTableCount {
		return nil, ERR_TABLE_LIMIT_EXCEEDED
	}

	if t.highestTableID-1 == currentTableCount {
		// There are no gaps in the table IDs, so we can just use the next available ID.
		tableID = t.highestTableID + 1
		t.highestTableID = tableID
	} else {
		// Find the first unused table ID.
		for i := uint64(firstTableID); i < t.maxTableCount; i++ {
			_, ok = t.prefixMap[i]
			if !ok {
				// We've found an unused table ID.
				tableID = i
				break
			}
		}
	}

	t.tableMap[name] = tableID

	prefix, err := getPrefix(tableID, getPrefixLength(t.maxTableCount))
	if err != nil {
		return nil, fmt.Errorf("error getting prefix for table ID %d: %w", tableID, err)
	}
	t.prefixMap[tableID] = prefix

	err = t.base.Put(t.namespaceTable.Uint64Key(tableID).GetRawBytes(), []byte(name))
	if err != nil {
		return nil, fmt.Errorf("error updating namespace table: %w", err)
	}

	return &keyBuilder{
		prefix: prefix,
	}, nil
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
