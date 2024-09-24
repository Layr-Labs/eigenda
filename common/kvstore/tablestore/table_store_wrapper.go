package tablestore

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"math"
	"sort"
)

// Table ID 0 is reserved for use internal use by the metadata table.
const metadataTableID uint32 = math.MaxUint32

// Table ID 1 is reserved for use by the namespace table. This stores a mapping between IDs and table names.
const namespaceTableID uint32 = math.MaxUint32 - 1

// The number of tables reserved for internal use.
const reservedTableCount uint32 = 2

// The first table ID that can be used by user tables.
const maxUserTableCount = math.MaxUint32 - reservedTableCount

// This key is used to store the schema version in the metadata table.
const schemaVersionKey = "schema_version"

// The deletionKey holds the value of the table name that is being deleted if one is currently being deleted.
// This is required for atomic deletion of tables.
const deletionKey = "deletion"

// The current schema version of the metadata table.
const currentSchemaVersion uint64 = 0

var _ kvstore.TableStore = &tableStore{}

// ERR_TABLE_LIMIT_EXCEEDED is returned when the maximum number of tables has been reached.
var ERR_TABLE_LIMIT_EXCEEDED = errors.New("table limit exceeded")

// tableStore is an implementation of TableStore that wraps a Store.
type tableStore struct {
	// A base store implementation that this TableStore wraps.
	base kvstore.Store

	// A map from table names to table IDs.
	tableMap map[string]uint32

	// The highest ID of all user tables in the store. Is -1 if there are no user tables.
	highestTableID int64

	// Builds keys for the metadata table.
	metadataTable kvstore.KeyBuilder

	// Builds keys for the namespace table.
	namespaceTable kvstore.KeyBuilder
}

// TableStoreWrapper wraps the given Store to create a TableStore.
//
// WARNING: it is not safe to access the wrapped store directly while the TableStore is in use. The TableStore uses
// special key formatting, and direct access to the wrapped store may violate the TableStore's invariants, resulting
// in undefined behavior.
func TableStoreWrapper(base kvstore.Store) (kvstore.TableStore, error) {

	tableMap := make(map[string]uint32)

	// Setup tables for internal use.
	metadataTable := &keyBuilder{prefix: getPrefix(metadataTableID)}
	namespaceTable := &keyBuilder{prefix: getPrefix(namespaceTableID)}

	highestTableID := int64(-1)

	schemaKey := metadataTable.StringKey(schemaVersionKey).GetRawBytes()
	onDiskSchemaBytes, err := base.Get(schemaKey)

	if errors.Is(err, kvstore.ErrNotFound) {
		// This store is new, no on disk schema version exists.
		onDiskSchemaBytes = make([]byte, 8)
		binary.BigEndian.PutUint64(onDiskSchemaBytes, currentSchemaVersion)

		err = base.Put(schemaKey, onDiskSchemaBytes)
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

		highestTableID, err = loadNamespaceTable(base, tableMap)
		if err != nil {
			return nil, fmt.Errorf("error loading namespace table: %w", err)
		}
	}

	store := &tableStore{
		base:           base,
		tableMap:       tableMap,
		highestTableID: highestTableID,
		metadataTable:  metadataTable,
		namespaceTable: namespaceTable,
	}

	err = store.handleIncompleteDeletion()
	if err != nil {
		return nil, fmt.Errorf("error handling incomplete deletion: %w", err)
	}

	return store, nil
}

// loadNamespaceTable loads the namespace table from disk into the given map. Returns the highest table ID found.
func loadNamespaceTable(base kvstore.Store, tableMap map[string]uint32) (int64, error) {
	highestTableID := int64(-1)

	it, err := base.NewIterator(getPrefix(namespaceTableID))
	if err != nil {
		return -1, fmt.Errorf("error creating namespace table iterator: %w", err)
	}
	defer it.Release()

	for it.Next() {
		keyBytes := it.Key()
		valueBytes := it.Value()

		_, tableIDBytes := parseKeyBytes(keyBytes)
		tableName := string(valueBytes)
		tableID := binary.BigEndian.Uint32(tableIDBytes)
		tableMap[tableName] = tableID

		if int64(tableID) > highestTableID {
			highestTableID = int64(tableID)
		}
	}

	return highestTableID, nil
}

// This method handles cleanup of incomplete deletions. Since deletion of a table requires multiple operations that
// are not atomic in aggregate, it is possible that a table deletion may have been started without being completed.
// This method makes sure that any such incomplete deletions are completed.
func (t *tableStore) handleIncompleteDeletion() error {
	deletionTableNameBytes, err := t.base.Get(t.metadataTable.StringKey(deletionKey).GetRawBytes())
	if errors.Is(err, kvstore.ErrNotFound) {
		// No deletion in progress, nothing to do.
		return nil
	}

	deletionTableName := string(deletionTableNameBytes)
	// TODO log something

	return t.DropTable(deletionTableName)
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
func (t *tableStore) GetOrCreateTable(name string) (kvstore.KeyBuilder, kvstore.Store, error) {
	tableID, ok := t.tableMap[name]
	if ok {
		// Table already exists.
		kb := &keyBuilder{
			prefix: getPrefix(tableID),
		}
		return kb, newTableView(t.base, kb), nil
	}

	currentTableCount := uint32(len(t.tableMap))
	if currentTableCount == math.MaxUint32-reservedTableCount {
		return nil, nil, ERR_TABLE_LIMIT_EXCEEDED
	}

	if uint32(t.highestTableID+1) == currentTableCount {
		// There are no gaps in the table IDs, so we can just use the next available ID.
		tableID = uint32(t.highestTableID + 1)
		t.highestTableID = int64(tableID)
	} else {
		// TODO write a unit test for this case specifically
		// Find the first unused table ID. This may not be efficient for a large number of table deletions
		// followed by a large number of table creations, but let's cross that bridge when we get to it.
		sortedTableIDs := make([]uint32, 0, currentTableCount)
		for _, id := range t.tableMap {
			sortedTableIDs = append(sortedTableIDs, id)
		}
		sort.Slice(sortedTableIDs, func(i, j int) bool {
			return sortedTableIDs[i] < sortedTableIDs[j]
		})
		next := uint32(0)
		for _, id := range sortedTableIDs {
			if id != next {
				tableID = next
				break
			}
			next++
		}
	}

	t.tableMap[name] = tableID

	err := t.base.Put(t.namespaceTable.Uint32Key(tableID).GetRawBytes(), []byte(name))
	if err != nil {
		return nil, nil, fmt.Errorf("error updating namespace table: %w", err)
	}

	kb := &keyBuilder{
		prefix: getPrefix(tableID),
	}
	return kb, newTableView(t.base, kb), nil
}

// DropTable deletes the table with the given name. This is a no-op if the table does not exist.
func (t *tableStore) DropTable(name string) error {

	fmt.Printf("tableStore.DropTable %s\n", name)     // TODO
	fmt.Printf("Table index: %d\n", t.tableMap[name]) // TODO

	tableID, ok := t.tableMap[name]
	if !ok {
		// Table does not exist, nothing to do.
		return nil
	}

	// This single atomic operation ensures that the table is deleted completely, even if we crash
	// in the middle of the operation. When next starting up, if an entry is observed in this location,
	// then the interrupted deletion can be completed.
	err := t.base.Put(t.metadataTable.StringKey(deletionKey).GetRawBytes(), []byte(name))
	if err != nil {
		return fmt.Errorf("error updating metadata table for deletion: %w", err)
	}

	it, err := t.base.NewIterator(getPrefix(tableID))
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

	// Finally, remove the deletion key from the metadata table.
	return t.base.Delete(t.metadataTable.StringKey(deletionKey).GetRawBytes())
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

// GetTableCount returns the current number of tables in the store.
func (t *tableStore) GetTableCount() uint32 {
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
	it, err := t.base.NewIterator(iteratorPrefix)
	if err != nil {
		return nil, err
	}

	return newTableIterator(it, prefix.GetKeyBuilder()), nil
}

// Shutdown shuts down the store, flushing any remaining cached data to disk.
func (t *tableStore) Shutdown() error {
	return t.base.Shutdown()
}

// Destroy shuts down and permanently deletes all data in the store.
func (t *tableStore) Destroy() error {
	return t.base.Destroy()
}
