package tablestore

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"math"
)

// TODO instead use high numbered tables!

// Table ID 0 is reserved for use internal use by the metadata table.
const metadataTableID uint32 = math.MaxUint32

// Table ID 1 is reserved for use by the namespace table. This stores a mapping between IDs and table names.
const namespaceTableID uint32 = math.MaxUint32 - 1

// Table ID 2 is reserved for use by the deletion table. This is used to make table deletion atomic.
const deletionTableID uint32 = math.MaxUint32 - 2

// The number of tables reserved for internal use.
const reservedTableCount uint32 = 3

// The first table ID that can be used by user tables.
const maxUserTableCount = math.MaxUint32 - reservedTableCount

// This key is used to store the schema version in the metadata table.
const schemaVersionKey = "schema_version"

// The current schema version of the metadata table.
const currentSchemaVersion uint64 = 0

var _ kvstore.TableStore = &tableStore{}

var ERR_TABLE_LIMIT_EXCEEDED = errors.New("table limit exceeded")

// tableStore is an implementation of TableStore that wraps a Store.
type tableStore struct {
	// A base store implementation that this TableStore wraps.
	base kvstore.Store

	// A map from table names to table IDs.
	tableMap map[string]uint32

	// A map from table IDs to table prefixes.
	prefixMap map[uint32][]byte

	// The highest ID of all user tables in the store. Is -1 if there are no user tables.
	highestTableID int64

	// Builds keys for the metadata table.
	metadataTable kvstore.KeyBuilder

	// Builds keys for the namespace table.
	namespaceTable kvstore.KeyBuilder

	// Builds keys for the deletion table.
	deletionTable kvstore.KeyBuilder
}

// TableStoreWrapper wraps the given Store to create a TableStore.
//
// Note that the max table count cannot be changed once the TableStore is created. Use the value 0 to use the default
// table size, which is (2^8 - 3) for new tables, or equal to the previous maximum table size if the TableStore is
// reloaded from disk. If the need arises, we may need to write migration code to resize the maximum number of tables,
// but this feature is not currently supported.
func TableStoreWrapper(base kvstore.Store) (kvstore.TableStore, error) {

	prefixMap := make(map[uint32][]byte)
	tableMap := make(map[string]uint32)

	// Setup tables for internal use.
	metadataPrefix := getPrefix(metadataTableID)
	metadataTable := &keyBuilder{prefix: metadataPrefix}
	namespacePrefix := getPrefix(namespaceTableID)
	namespaceTable := &keyBuilder{prefix: namespacePrefix}
	deletionPrefix := getPrefix(deletionTableID)
	deletionTable := &keyBuilder{prefix: deletionPrefix}

	highestTableID := int64(-1)

	schemaKey := metadataTable.StringKey(schemaVersionKey).GetRawBytes()
	onDiskSchemaBytes, err := base.Get(schemaKey)

	if errors.Is(err, kvstore.ErrNotFound) {
		// This store is new, no on disk schema version exists.
		err = base.Put(schemaKey, []byte{byte(currentSchemaVersion)})
		if err != nil {
			return nil, fmt.Errorf("error setting schema version in metadata table: %w", err)
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

		// Load namespace table.
		it, err := base.NewIterator(namespacePrefix)
		if err != nil {
			return nil, fmt.Errorf("error creating namespace table iterator: %w", err)
		}
		defer it.Release()

		// TODO verify this is the correct pattern, i.e. we aren't skipping the first entry
		for it.Next() {
			keyBytes := it.Key()
			valueBytes := it.Value()

			_, tableIDBytes := parseKeyBytes(keyBytes)
			tableName := string(valueBytes)
			tableID := binary.BigEndian.Uint32(tableIDBytes)
			tableMap[tableName] = tableID

			prefixMap[tableID] = getPrefix(tableID)

			if int64(tableID) > highestTableID {
				highestTableID = int64(tableID)
			}
		}
	}

	// TODO handle deletion that was started but not finished

	store := &tableStore{
		base:           base,
		tableMap:       tableMap,
		prefixMap:      prefixMap,
		highestTableID: highestTableID,
		metadataTable:  metadataTable,
		namespaceTable: namespaceTable,
		deletionTable:  deletionTable,
	}

	return store, nil
}

// Get the prefix for the given table ID and prefix length.
func getPrefix(tableID uint32) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, tableID)
	return bytes
}

// Parse the table ID and key bytes from the given key bytes.
func parseKeyBytes(keyBytes []byte) (tableID uint32, key []byte) {
	tableID = binary.BigEndian.Uint32(keyBytes[:4])
	key = keyBytes[4:]
	return
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

	currentTableCount := uint32(len(t.tableMap))
	if currentTableCount == math.MaxUint32-reservedTableCount {
		return nil, ERR_TABLE_LIMIT_EXCEEDED
	}

	if uint32(t.highestTableID+1) == currentTableCount {
		// There are no gaps in the table IDs, so we can just use the next available ID.
		tableID = uint32(t.highestTableID + 1)
		t.highestTableID = int64(tableID)
	} else {
		// Find the first unused table ID.
		for i := uint32(0); i < maxUserTableCount; i++ {
			_, ok = t.prefixMap[i]
			if !ok {
				// We've found an unused table ID.
				tableID = i
				break
			}
		}
	}

	t.tableMap[name] = tableID
	prefix := getPrefix(tableID)
	t.prefixMap[tableID] = prefix

	err := t.base.Put(t.namespaceTable.Uint32Key(tableID).GetRawBytes(), []byte(name))
	if err != nil {
		return nil, fmt.Errorf("error updating namespace table: %w", err)
	}

	return &keyBuilder{
		prefix: prefix,
	}, nil
}

// DropTable deletes the table with the given name. This is a no-op if the table does not exist.
func (t *tableStore) DropTable(name string) error {
	tableID, ok := t.tableMap[name]
	if !ok {
		// Table does not exist, nothing to do.
		return nil
	}

	prefix := t.prefixMap[tableID]
	it, err := t.base.NewIterator(prefix)
	if err != nil {
		return fmt.Errorf("error creating iterator for table: %w", err)
	}
	defer it.Release()

	// Future work: if this is a performance bottleneck, use batching.
	for it.Next() {
		err = t.base.Delete(it.Key())
		if err != nil {
			return fmt.Errorf("error deleting key from table: %w", err)
		}
	}

	// All table entries have been deleted. Now delete the table from the namespace table.
	err = t.base.Delete(t.namespaceTable.Uint32Key(tableID).GetRawBytes())
	if err != nil {
		return fmt.Errorf("error deleting from namespace table: %w", err)
	}
	delete(t.tableMap, name)
	delete(t.prefixMap, tableID)

	return nil
}

// GetKeyBuilder returns a key builder for the table with the given name. Throws an error if the table does not exist.
func (t *tableStore) GetKeyBuilder(name string) (kvstore.KeyBuilder, error) {
	tableID, ok := t.tableMap[name]
	if !ok {
		return nil, fmt.Errorf("table %s does not exist", name)
	}

	return &keyBuilder{
		prefix: getPrefix(tableID),
	}, nil
}

// GetMaxTableCount returns the maximum number of tables that can be created in the store.
func (t *tableStore) GetMaxTableCount() uint32 {
	return maxUserTableCount
}

// GetCurrentTableCount returns the current number of tables in the store.
func (t *tableStore) GetCurrentTableCount() uint32 {
	return uint32(len(t.tableMap))
}

// GetTables returns a list of all tables in the store in no particular order.
func (t *tableStore) GetTables() []string {
	tables := make([]string, 0, len(t.tableMap))
	for name := range t.tableMap {
		tables = append(tables, name)
	}

	return tables
}

// Put stores the given key / value pair in the database, overwriting any existing value for that key.
func (t *tableStore) Put(key kvstore.Key, value []byte) error {
	return t.base.Put(key.GetRawBytes(), value)
}

// Get retrieves the value for the given key from the database.
// Returns a kvstore.ErrNotFound error if the key does not exist.
func (t *tableStore) Get(key kvstore.Key) ([]byte, error) {
	return t.base.Get(key.GetRawBytes())
}

// Delete removes the key from the database. Does not return an error if the key does not exist.
func (t *tableStore) Delete(key kvstore.Key) error {
	return t.base.Delete(key.GetRawBytes())
}

// DeleteBatch atomically removes a list of keys from the database.
func (t *tableStore) DeleteBatch(keys []kvstore.Key) error {
	keyBytes := make([][]byte, len(keys))

	for i, key := range keys {
		keyBytes[i] = key.GetRawBytes()
	}

	return t.base.DeleteBatch(keyBytes)
}

// WriteBatch atomically writes a list of key / value pairs to the database.
func (t *tableStore) WriteBatch(keys []kvstore.Key, values [][]byte) error {
	keyBytes := make([][]byte, len(keys))

	for i, key := range keys {
		keyBytes[i] = key.GetRawBytes()
	}

	return t.base.WriteBatch(keyBytes, values)
}

// NewIterator returns an iterator that can be used to iterate over a subset of the keys in the database.
// Only keys with the given key's table with prefix matching the key will be iterated. The iterator must be closed
// by calling Release() when done. The iterator will return keys in lexicographically sorted order. The iterator
// walks over a consistent snapshot of the database, so it will not see any writes that occur after the iterator
// is created.
func (t *tableStore) NewIterator(prefix kvstore.Key) (iterator.Iterator, error) {
	iteratorPrefix := prefix.GetRawBytes()
	return t.base.NewIterator(iteratorPrefix)
}

// Shutdown shuts down the store, flushing any remaining cached data to disk.
func (t *tableStore) Shutdown() error {
	return t.base.Shutdown()
}

// Destroy shuts down and permanently deletes all data in the store.
func (t *tableStore) Destroy() error {
	return t.base.Destroy()
}
