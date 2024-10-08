package tablestore

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/Layr-Labs/eigenda/common/kvstore/leveldb"
	"github.com/Layr-Labs/eigenda/common/kvstore/mapstore"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"math"
	"sort"
	"sync"
)

// The table ID reserved for the metadata table.
const metadataTableID uint32 = math.MaxUint32

// The table ID reserved for the namespace table.
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

// StoreType describes the underlying store implementation.
type StoreType int

const (
	// LevelDB is a LevelDB-backed store.
	LevelDB StoreType = iota
	// MapStore is an in-memory store. This store does not preserve data across restarts.
	MapStore
)

// New creates a new TableStore instance of the given type. The store will be created at the given path.
// This method will set up a table for each table name provided, and will drop all tables not in the list.
// Dropping a table is irreversible and will delete all data in the table, so be very careful not to call
// this method with table names omitted by mistake.
func (t StoreType) New(logger logging.Logger, path string, tables ...string) (kvstore.TableStore, error) {
	builder, err := t.builder(logger, path) // TODO
	if err != nil {
		return nil, fmt.Errorf("error creating store builder: %w", err)
	}

	originalTables := builder.GetTableNames()
	originalTablesSet := make(map[string]bool)
	for _, table := range originalTables {
		originalTablesSet[table] = true
	}

	newTablesSet := make(map[string]bool)
	for _, table := range tables {
		newTablesSet[table] = true
	}

	// Add new tables.
	for _, table := range tables {
		if !originalTablesSet[table] {
			err = builder.CreateTable(table)
			if err != nil {
				return nil, fmt.Errorf("error creating table %s: %w", table, err)
			}

		}
	}

	// Drop tables that are not in the list.
	for _, table := range originalTables {
		if !newTablesSet[table] {
			err = builder.DropTable(table)
			if err != nil {
				return nil, fmt.Errorf("error dropping table %s: %w", table, err)
			}
		}
	}

	store, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("error building store: %w", err)
	}

	return store, nil
}

// Future work: if we ever decide to permit third parties to provide custom store implementations not in this module,
// we will need to provide a way to instantiate a builder using an instantiated base store.

// Builder creates a new TableStoreBuilder of the given type. Any data written to disk by the TableStoreBuilder
// (and the eventual TableStore it builds) will be stored in the given path. Can be used to create a new store
// or to load an existing store from disk.
func (t StoreType) builder(logger logging.Logger, path string) (kvstore.TableStoreBuilder, error) {
	switch t {
	case LevelDB:
		store, err := leveldb.NewStore(logger, path)
		if err != nil {
			return nil, fmt.Errorf("error creating LevelDB store: %w", err)
		}
		return newTableStoreBuilder(logger, store)
	case MapStore:
		store := mapstore.NewStore()
		return newTableStoreBuilder(logger, store)
	default:
		return nil, fmt.Errorf("unknown store type: %d", t)
	}
}

var _ kvstore.TableStoreBuilder = (*tableStoreBuilder)(nil)

type tableStoreBuilder struct {
	logger logging.Logger

	// Used to make this builder thread-safe.
	lock sync.Mutex

	// A base store implementation
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

	// Whether the store has been built.
	built bool
}

// newTableStoreBuilder creates a new TableStoreBuilder instance.
func newTableStoreBuilder(logger logging.Logger, base kvstore.Store) (*tableStoreBuilder, error) {

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

	builder := &tableStoreBuilder{
		logger:         logger,
		base:           base,
		tableIDMap:     tableIDMap,
		tableMap:       tableMap,
		metadataTable:  metadataTable,
		namespaceTable: namespaceTable,
		highestTableID: highestTableID,
	}

	err = builder.handleIncompleteDeletion()
	if err != nil {
		return nil, fmt.Errorf("error handling incomplete deletion: %w", err)
	}

	return builder, nil
}

// This method handles cleanup of incomplete deletions. Since deletion of a table requires multiple operations that
// are not atomic in aggregate, it is possible that a table deletion may have been started without being completed.
// This method makes sure that any such incomplete deletions are completed.
func (t *tableStoreBuilder) handleIncompleteDeletion() error {
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

// CreateTable creates a new table with the given name. If a table with the given name already exists,
// this method becomes a no-op. Returns ErrTableLimitExceeded if the maximum number of tables has been reached.
func (t *tableStoreBuilder) CreateTable(name string) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.built {
		return errors.New("cannot modify a built store")
	}

	tableID, ok := t.tableIDMap[name]
	if ok {
		return nil
	}

	currentTableCount := uint32(len(t.tableIDMap))
	if currentTableCount == math.MaxUint32-reservedTableCount {
		return kvstore.ErrTableLimitExceeded
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
		return fmt.Errorf("error updating namespace table: %w", err)
	}

	return nil
}

// DropTable deletes the table with the given name. This is a no-op if the table does not exist.
func (t *tableStoreBuilder) DropTable(name string) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.built {
		return errors.New("cannot modify a built store")
	}

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

// GetTableCount returns the current number of tables in the store (excluding internal tables utilized by the store).
func (t *tableStoreBuilder) GetTableCount() uint32 {
	t.lock.Lock()
	defer t.lock.Unlock()

	return uint32(len(t.tableIDMap))
}

// GetMaxTableCount returns the maximum number of tables that can be created in the store.
func (t *tableStoreBuilder) GetMaxTableCount() uint32 {
	return maxUserTableCount
}

// GetTableNames returns a list of the names of all tables in the store, in no particular order.
func (t *tableStoreBuilder) GetTableNames() []string {
	t.lock.Lock()
	defer t.lock.Unlock()

	names := make([]string, 0, len(t.tableIDMap))
	for name := range t.tableIDMap {
		names = append(names, name)
	}
	return names
}

// Build creates a new TableStore instance with the specified tables. After this method is called,
// the TableStoreBuilder should not be used again.
func (t *tableStoreBuilder) Build() (kvstore.TableStore, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.built {
		return nil, errors.New("cannot build a store more than once")
	}

	tableMap := make(map[string]kvstore.Table, len(t.tableIDMap))

	for name, tableID := range t.tableIDMap {
		tableMap[name] = t.tableMap[tableID]
	}

	return newTableStore(t.logger, t.base, tableMap), nil
}

// Shutdown shuts down the store, flushing any remaining cached data to disk.
func (t *tableStoreBuilder) Shutdown() error {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.built {
		return errors.New("the builder cannot be used to shut down a store once it has been built")
	}

	return t.base.Shutdown()
}

// Destroy shuts down and permanently deletes all data in the store.
func (t *tableStoreBuilder) Destroy() error {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.built {
		return errors.New("the builder cannot be used to destroy a store once it has been built")
	}

	return t.base.Destroy()
}
