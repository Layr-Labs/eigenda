package tablestore

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/Layr-Labs/eigensdk-go/logging"
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
const metadataSchemaVersionKey = "schema_version"

// The metadataDeletionKey holds the value of the table name that is being deleted if one is currently being deleted.
// This is required for atomic deletion of tables.
const metadataDeletionKey = "deletion"

// The current schema version of the metadata table.
const currentSchemaVersion uint64 = 0

var _ kvstore.TableStore = &tableStore{}

// ERR_TABLE_LIMIT_EXCEEDED is returned when the maximum number of tables has been reached.
var ERR_TABLE_LIMIT_EXCEEDED = errors.New("table limit exceeded")

// tableStore is an implementation of TableStore that wraps a Store.
type tableStore struct {
	logger logging.Logger

	// A base store implementation that this TableStore wraps.
	base kvstore.Store

	// A map from table names to table IDs.
	tableIDMap map[string]uint32

	// A map from table IDs to tables.
	tableMap map[uint32]kvstore.Table

	// The highest ID of all user tables in the store. Is -1 if there are no user tables.
	highestTableID int64

	// Builds keys for the metadata table.
	metadataTable kvstore.Table

	// Builds keys for the namespace table.
	namespaceTable kvstore.Table
}

// TableStoreWrapper wraps the given Store to create a TableStore.
//
// WARNING: it is not safe to access the wrapped store directly while the TableStore is in use. The TableStore uses
// special key formatting, and direct access to the wrapped store may violate the TableStore's invariants, resulting
// in undefined behavior.
func TableStoreWrapper(logger logging.Logger, base kvstore.Store) (kvstore.TableStore, error) {

	tableIDMap := make(map[string]uint32)
	tableIdSet := make(map[uint32]bool)
	tableMap := make(map[uint32]kvstore.Table)

	highestTableID := int64(-1)

	metadataTable := newTableView(base, "metadata", metadataTableID)
	schemaKey := []byte(metadataSchemaVersionKey)
	onDiskSchemaBytes, err := metadataTable.Get(schemaKey)

	namespaceTable := newTableView(base, "namespace", namespaceTableID)

	if errors.Is(err, kvstore.ErrNotFound) {
		// This store is new, no on disk schema version exists.
		onDiskSchemaBytes = make([]byte, 8)
		binary.BigEndian.PutUint64(onDiskSchemaBytes, currentSchemaVersion)

		err = metadataTable.Put(schemaKey, onDiskSchemaBytes)
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

		highestTableID, err = loadNamespaceTable(base, namespaceTable, tableIDMap, tableMap, tableIdSet)
		if err != nil {
			return nil, fmt.Errorf("error loading namespace table: %w", err)
		}
	}

	store := &tableStore{
		logger:         logger,
		base:           base,
		tableIDMap:     tableIDMap,
		tableMap:       tableMap,
		metadataTable:  metadataTable,
		namespaceTable: namespaceTable,
		highestTableID: highestTableID,
	}

	err = store.handleIncompleteDeletion()
	if err != nil {
		return nil, fmt.Errorf("error handling incomplete deletion: %w", err)
	}

	return store, nil
}

// loadNamespaceTable loads the namespace table from disk into the given map. Returns the highest table ID found.
func loadNamespaceTable(
	base kvstore.Store,
	namespaceTable kvstore.Table,
	tableIDMap map[string]uint32,
	tableMap map[uint32]kvstore.Table,
	tableIdSet map[uint32]bool) (int64, error) {

	highestTableID := int64(-1)

	it, err := namespaceTable.NewIterator(nil)
	if err != nil {
		return -1, fmt.Errorf("error creating namespace table iterator: %w", err)
	}
	defer it.Release()

	for it.Next() {
		keyBytes := it.Key()
		valueBytes := it.Value()

		tableID := binary.BigEndian.Uint32(keyBytes)
		tableName := string(valueBytes)
		tableIdSet[tableID] = true
		tableIDMap[tableName] = tableID
		tableMap[tableID] = newTableView(base, tableName, tableID)

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
	deletionTableNameBytes, err := t.metadataTable.Get([]byte(metadataDeletionKey))
	if errors.Is(err, kvstore.ErrNotFound) {
		// No deletion in progress, nothing to do.
		return nil
	}

	deletionTableName := string(deletionTableNameBytes)
	t.logger.Errorf("found incomplete deletion of table %s, completing deletion", deletionTableName)

	return t.DropTable(deletionTableName)
}

// Get the prefix for the given table ID and prefix length.
func getPrefix(tableID uint32) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, tableID)
	return bytes
}

// GetTable gets the table with the given name. If the table does not exist, it is first created.
func (t *tableStore) GetTable(name string) (kvstore.Table, error) {

	tableID, ok := t.tableIDMap[name]
	if ok {
		return t.tableMap[tableID], nil
	}

	currentTableCount := uint32(len(t.tableIDMap))
	if currentTableCount == math.MaxUint32-reservedTableCount {
		return nil, ERR_TABLE_LIMIT_EXCEEDED
	}

	if uint32(t.highestTableID+1) == currentTableCount {
		// There are no gaps in the table IDs, so we can just use the next available ID.
		tableID = uint32(t.highestTableID + 1)
		t.highestTableID = int64(tableID)
	} else {
		// Find the first unused table ID. This may not be efficient for a large number of table deletions
		// followed by a large number of table creations, but let's cross that bridge when we get to it.
		sortedTableIDs := make([]uint32, 0, currentTableCount)
		for _, id := range t.tableIDMap {
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

	t.tableIDMap[name] = tableID
	table := newTableView(t.base, name, tableID)
	t.tableMap[tableID] = table

	tableKey := make([]byte, 4)
	binary.BigEndian.PutUint32(tableKey, tableID)
	err := t.namespaceTable.Put(tableKey, []byte(name))

	if err != nil {
		return nil, fmt.Errorf("error updating namespace table: %w", err)
	}

	return table, nil
}

// DropTable deletes the table with the given name. This is a no-op if the table does not exist.
func (t *tableStore) DropTable(name string) error {
	tableID, ok := t.tableIDMap[name]
	if !ok {
		// Table does not exist, nothing to do.
		return nil
	}

	// This single atomic operation ensures that the table is deleted completely, even if we crash
	// in the middle of the operation. When next starting up, if an entry is observed in this location,
	// then the interrupted deletion can be completed.
	deletionKey := []byte(metadataDeletionKey)
	err := t.metadataTable.Put(deletionKey, []byte(name))
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
	tableKey := make([]byte, 4)
	binary.BigEndian.PutUint32(tableKey, tableID)
	err = t.namespaceTable.Delete(tableKey)
	if err != nil {
		return fmt.Errorf("error deleting from namespace table: %w", err)
	}
	delete(t.tableIDMap, name)
	delete(t.tableMap, tableID)

	// Update highestTableID as needed.
	for ; t.highestTableID >= 0; t.highestTableID-- {
		if _, ok := t.tableMap[uint32(t.highestTableID)]; ok {
			break
		}
	}

	// Finally, remove the deletion key from the metadata table.
	return t.metadataTable.Delete(deletionKey)
}

// GetMaxTableCount returns the maximum number of tables that can be created in the store.
func (t *tableStore) GetMaxTableCount() uint32 {
	return maxUserTableCount
}

// GetTableCount returns the current number of tables in the store.
func (t *tableStore) GetTableCount() uint32 {
	return uint32(len(t.tableIDMap))
}

// GetTables returns a list of all tables in the store in no particular order.
func (t *tableStore) GetTables() []kvstore.Table {
	tables := make([]kvstore.Table, 0, len(t.tableIDMap))
	for _, table := range t.tableMap {
		tables = append(tables, table)
	}

	return tables
}

// tableStoreBatch is a batch for writing to a table store.
type tableStoreBatch struct {
	store *tableStore
	batch kvstore.Batch[[]byte]
}

// Put adds a key-value pair to the batch.
func (t *tableStoreBatch) Put(key kvstore.TableKey, value []byte) {
	t.batch.Put(key, value)
}

// Delete removes a key-value pair from the batch.
func (t *tableStoreBatch) Delete(key kvstore.TableKey) {
	t.batch.Delete(key)
}

// Apply applies the batch to the store.
func (t *tableStoreBatch) Apply() error {
	return t.batch.Apply()
}

// Size returns the number of operations in the batch.
func (t *tableStoreBatch) Size() uint32 {
	return t.batch.Size()
}

// NewBatch creates a new batch for writing to the store.
func (t *tableStore) NewBatch() kvstore.Batch[kvstore.TableKey] {
	return &tableStoreBatch{
		store: t,
		batch: t.base.NewBatch(),
	}
}

// Shutdown shuts down the store, flushing any remaining cached data to disk.
func (t *tableStore) Shutdown() error {
	return t.base.Shutdown()
}

// Destroy shuts down and permanently deletes all data in the store.
func (t *tableStore) Destroy() error {
	return t.base.Destroy()
}
