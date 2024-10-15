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
)

// The table ID reserved for the metadata table. The metadata table is used for internal bookkeeping.
// The following data is currently stored in the metadata table:
//   - Schema version (in case we ever need to do a schema migration)
//   - Table deletion markers used to detect a crash during table deletion and to complete
//     the deletion when the store is next started.
const metadataTableID uint32 = math.MaxUint32

// The table ID reserved for the namespace table. Keys in the namespace table are table IDs (uint32)
// and values are table names (string). Although no two tables are permitted to have share a name or a table ID,
// once a table is dropped, future tables may be instantiated with the same name and the table ID may be reused.
const namespaceTableID uint32 = math.MaxUint32 - 1

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

// Start creates a new TableStore instance of the given type. The store will be created at the given path.
// This method can be used to instantiate a new store or to load an existing store.
// This method will set up a table for each table name provided, and will drop all tables not in the list.
// Dropping a table is irreversible and will delete all data in the table, so be very careful not to call
// this method with table names omitted by mistake.
func (t StoreType) Start(logger logging.Logger, path string, tables ...string) (kvstore.TableStore, error) {
	base, err := buildBaseStore(t, logger, path)
	if err != nil {
		return nil, fmt.Errorf("error building base store: %w", err)
	}

	return start(logger, base, true, tables...)
}

// Load loads a table store from disk without modifying the table schema. If there is no existing store at the given
// path, this method will create one and return a store without any tables.
func (t StoreType) Load(logger logging.Logger, path string) (kvstore.TableStore, error) {
	base, err := buildBaseStore(t, logger, path)
	if err != nil {
		return nil, fmt.Errorf("error building base store: %w", err)
	}

	return start(logger, base, false)
}

// Future work: if we ever decide to permit third parties to provide custom store implementations not in this module,
// we will need to make this method public.

// start creates a new TableStore instance with the given base store and table names. If modifySchema is true, the
// tables in the store are made to match the given list of tables by adding and removing tables as needed. If
// modifySchema is false, the tables are loaded as is, and any tables in the provided list are ignored.
func start(
	logger logging.Logger,
	base kvstore.Store,
	modifySchema bool,
	tables ...string) (kvstore.TableStore, error) {

	metadataTable := newTableView(base, "metadata", metadataTableID)
	namespaceTable := newTableView(base, "namespace", namespaceTableID)

	err := validateSchema(metadataTable)
	if err != nil {
		return nil, fmt.Errorf("error validating schema: %w", err)
	}

	err = handleIncompleteDeletion(logger, base, metadataTable, namespaceTable)
	if err != nil {
		return nil, fmt.Errorf("error handling incomplete deletion: %w", err)
	}

	tableIDMap, err := loadNamespaceTable(namespaceTable)
	if err != nil {
		return nil, fmt.Errorf("error loading namespace table: %w", err)
	}

	if modifySchema {
		err = addAndRemoveTables(base, metadataTable, namespaceTable, tableIDMap, tables)
		if err != nil {
			return nil, fmt.Errorf("error adding and removing tables: %w", err)
		}
	}

	tableMap := make(map[string]kvstore.Table, len(tableIDMap))
	for tableID, tableName := range tableIDMap {
		tableMap[tableName] = newTableView(base, tableName, tableID)
	}
	return newTableStore(logger, base, tableMap), nil
}

// buildBaseStore creates a new base store of the given type.
func buildBaseStore(storeType StoreType, logger logging.Logger, path string) (kvstore.Store, error) {
	switch storeType {
	case LevelDB:
		return leveldb.NewStore(logger, path)
	case MapStore:
		return mapstore.NewStore(), nil
	default:
		return nil, fmt.Errorf("unknown store type: %d", storeType)
	}
}

// validateSchema loads/initiates the schema version in the metadata table.
func validateSchema(metadataTable kvstore.Table) error {

	schemaKey := []byte(metadataSchemaVersionKey)
	onDiskSchemaBytes, err := metadataTable.Get(schemaKey)

	if err != nil {
		if !errors.Is(err, kvstore.ErrNotFound) {
			return err
		}

		// This store is new, no on disk schema version exists.
		onDiskSchemaBytes = make([]byte, 8)
		binary.BigEndian.PutUint64(onDiskSchemaBytes, currentSchemaVersion)

		err = metadataTable.Put(schemaKey, onDiskSchemaBytes)
		if err != nil {
			return fmt.Errorf("error setting schema version in metadata table: %w", err)
		}

		return nil
	}

	// Verify schema version.
	onDiskSchema := binary.BigEndian.Uint64(onDiskSchemaBytes)
	if onDiskSchema != currentSchemaVersion {
		// In the future if we change schema versions, we may need to write migration code here.
		return fmt.Errorf(
			"incompatible schema version: code is at version %d, data on disk is at version %d",
			currentSchemaVersion, onDiskSchema)
	}

	return nil
}

// This method adds and removes tables as needed to match the given list of tables. The table ID map is updated
// to reflect the new state of the tables.
func addAndRemoveTables(
	base kvstore.Store,
	metadataTable kvstore.Table,
	namespaceTable kvstore.Table,
	tableIDMap map[uint32]string,
	currentTables []string) error {

	tablesToAdd, tablesToDrop := computeSchemaChange(tableIDMap, currentTables)

	err := dropTables(base, metadataTable, namespaceTable, tableIDMap, tablesToDrop)
	if err != nil {
		return fmt.Errorf("error dropping tables: %w", err)
	}

	err = addTables(namespaceTable, tableIDMap, tablesToAdd)
	if err != nil {
		return fmt.Errorf("error adding tables: %w", err)
	}

	err = sanityCheckNamespaceTable(namespaceTable, tableIDMap, currentTables)
	if err != nil {
		return fmt.Errorf("error sanity checking namespace table: %w", err)
	}

	return nil
}

// Compute the tables that need to be added and dropped to match the given list of tables.
func computeSchemaChange(
	tableIDMap map[uint32]string,
	currentTables []string) (tablesToAdd []string, tablesToDrop []string) {

	tablesToAdd = make([]string, 0)
	tablesToDrop = make([]string, 0)

	originalTablesSet := make(map[string]struct{})
	for _, table := range tableIDMap {
		originalTablesSet[table] = struct{}{}
	}
	newTablesSet := make(map[string]bool)
	for _, table := range currentTables {
		newTablesSet[table] = true
	}
	for table := range originalTablesSet {
		if !newTablesSet[table] {
			tablesToDrop = append(tablesToDrop, table)
		}
	}
	for table := range newTablesSet {
		if _, exists := originalTablesSet[table]; !exists {
			tablesToAdd = append(tablesToAdd, table)
		}
	}

	return tablesToAdd, tablesToDrop
}

// Drop a list of tables. Updates the table ID map as well as data within the store.
func dropTables(
	base kvstore.Store,
	metadataTable kvstore.Table,
	namespaceTable kvstore.Table,
	tableIDMap map[uint32]string,
	tablesToDrop []string) error {

	if len(tablesToDrop) == 0 {
		// bail out early
		return nil
	}

	reverseTableIDMap := make(map[string]uint32)
	for tableName, tableID := range tableIDMap {
		reverseTableIDMap[tableID] = tableName
	}
	for _, tableName := range tablesToDrop {
		err := dropTable(base, metadataTable, namespaceTable, tableName, reverseTableIDMap[tableName])
		if err != nil {
			return fmt.Errorf("error dropping table %s: %w", tableName, err)
		}
		delete(tableIDMap, reverseTableIDMap[tableName])
	}

	return nil
}

// Add tables to the store. Updates the table ID map as well as data within the store.
func addTables(
	namespaceTable kvstore.Table,
	tableIDMap map[uint32]string,
	tablesToAdd []string) error {

	if len(tablesToAdd) == 0 {
		// bail out early
		return nil
	}

	// Determine the table IDs for the new tables to be added.
	// We want to fill gaps prior to assigning new IDs.
	newTableIDs := make([]uint32, 0, len(tablesToAdd))
	nextID := uint32(0)
	for len(newTableIDs) < len(tablesToAdd) {
		if _, alreadyUsed := tableIDMap[nextID]; !alreadyUsed {
			newTableIDs = append(newTableIDs, nextID)
		}
		nextID++
	}

	// Sort tables to add. Ensures deterministic table IDs given the same input.
	sort.Strings(tablesToAdd)

	// Add new tables.
	for i, tableName := range tablesToAdd {
		tableID := newTableIDs[i]
		err := createTable(namespaceTable, tableName, tableID)
		if err != nil {
			return fmt.Errorf("error creating table %s: %w", tableName, err)
		}
		tableIDMap[tableID] = tableName
	}

	return nil
}

// Perform sanity checks on the namespace table.
// This method makes potential logic in the namespace errors fail fast and visibly.
func sanityCheckNamespaceTable(
	namespaceTable kvstore.Table,
	tableIDMap map[uint32]string,
	currentTableList []string) error {

	parsedNamespaceTable, err := loadNamespaceTable(namespaceTable)
	if err != nil {
		return fmt.Errorf("error loading namespace table: %w", err)
	}

	if len(parsedNamespaceTable) != len(tableIDMap) {
		return fmt.Errorf("namespace table has %d entries, but expected %d", len(parsedNamespaceTable), len(tableIDMap))
	}

	reverseNamespaceTable := make(map[string]uint32)
	for tableID, tableName := range tableIDMap {
		reverseNamespaceTable[tableName] = tableID
		if parsedNamespaceTable[tableID] != tableName {
			return fmt.Errorf("namespace table has mismatched entry for table %s", tableName)
		}
	}

	for _, tableName := range currentTableList {
		_, exists := parsedNamespaceTable[reverseNamespaceTable[tableName]]
		if !exists {
			return fmt.Errorf("namespace table is missing entry for table %s", tableName)
		}
	}

	return nil
}

// This method handles cleanup of incomplete deletions. Since deletion of a table requires multiple operations that
// are not atomic in aggregate, it is possible that a table deletion may have been started without being completed.
// This method makes sure that any such incomplete deletions are completed.
func handleIncompleteDeletion(
	logger logging.Logger,
	base kvstore.Store,
	metadataTable kvstore.Table,
	namespaceTable kvstore.Table) error {

	deletionValue, err := metadataTable.Get([]byte(metadataDeletionKey))
	if errors.Is(err, kvstore.ErrNotFound) {
		// No deletion in progress, nothing to do.
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading metadata table for deletion: %w", err)
	}

	deletionTableID := binary.BigEndian.Uint32(deletionValue)
	deletionTableName := string(deletionValue[4:])

	logger.Errorf("found incomplete deletion of table %s, completing deletion", deletionTableName)
	return dropTable(base, metadataTable, namespaceTable, deletionTableName, deletionTableID)
}

// Get the prefix for the given table ID and prefix length.
func getPrefix(tableID uint32) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, tableID)
	return bytes
}

// loadNamespaceTable loads the namespace table from disk into the given map.
// Returns a map from table IDs to table names.
func loadNamespaceTable(namespaceTable kvstore.Table) (map[uint32]string, error) {

	tableIDMap := make(map[uint32]string)

	it, err := namespaceTable.NewIterator(nil)
	if err != nil {
		return nil, fmt.Errorf("error creating namespace table iterator: %w", err)
	}
	defer it.Release()

	for it.Next() {
		keyBytes := it.Key()
		valueBytes := it.Value()

		tableID := binary.BigEndian.Uint32(keyBytes)
		tableName := string(valueBytes)
		tableIDMap[tableID] = tableName
	}
	return tableIDMap, nil
}

// CreateTable creates a new table with the given name. If a table with the given name already exists,
// this method becomes a no-op. Returns ErrTableLimitExceeded if the maximum number of tables has been reached.
func createTable(
	namespaceTable kvstore.Table,
	name string,
	tableID uint32) error {

	tableKey := make([]byte, 4)
	binary.BigEndian.PutUint32(tableKey, tableID)
	err := namespaceTable.Put(tableKey, []byte(name))

	if err != nil {
		return fmt.Errorf("error updating namespace table: %w", err)
	}

	return nil
}

// DropTable deletes the table with the given name. This is a no-op if the table does not exist.
func dropTable(
	base kvstore.Store,
	metadataTable kvstore.Table,
	namespaceTable kvstore.Table,
	name string,
	tableID uint32) error {

	// This single atomic operation ensures that the table is deleted completely, even if we crash
	// in the middle of the operation. When next starting up, if an entry is observed in this location,
	// then the interrupted deletion can be completed.
	deletionKey := []byte(metadataDeletionKey)

	deletionValue := make([]byte, 4+len(name))
	binary.BigEndian.PutUint32(deletionValue, tableID)
	copy(deletionValue[4:], name)

	err := metadataTable.Put(deletionKey, deletionValue)
	if err != nil {
		return fmt.Errorf("error updating metadata table for deletion: %w", err)
	}

	it, err := base.NewIterator(getPrefix(tableID))
	if err != nil {
		return fmt.Errorf("error creating iterator for table: %w", err)
	}
	defer it.Release()

	// Future work: if this is a performance bottleneck, use batching.
	for it.Next() {
		err = base.Delete(it.Key())
		if err != nil {
			return fmt.Errorf("error deleting key from table: %w", err)
		}
	}

	// All table entries have been deleted. Now delete the table from the namespace table.
	tableKey := make([]byte, 4)
	binary.BigEndian.PutUint32(tableKey, tableID)
	err = namespaceTable.Delete(tableKey)
	if err != nil {
		return fmt.Errorf("error deleting from namespace table: %w", err)
	}

	// Finally, remove the deletion key from the metadata table.
	return metadataTable.Delete(deletionKey)
}
