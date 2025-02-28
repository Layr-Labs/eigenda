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

// The table ID reserved for the expiration table. The expiration table is used to store expiration times for keys.
// The keys in the expiration table are created by prepending the expiration time to the key of the data to be expired.
// The value of the key in this table is an empty byte slice. By iterating over this table in lexicographical order,
// keys are encountered in order of expiration time.
const expirationTableID uint32 = math.MaxUint32 - 2

// This key is used to store the schema version in the metadata table.
const metadataSchemaVersionKey = "schema_version"

// The metadataDeletionKey holds the value of the table name that is being deleted if one is currently being deleted.
// This is required for atomic deletion of tables.
const metadataDeletionKey = "deletion"

// When a new table is created, the ID used for that table is stored in the metadata table under this key.
const nextTableIDKey = "next_table_id"

// The current schema version of the metadata table.
const currentSchemaVersion uint64 = 0

// Start creates a new TableStore. This method can be used to instantiate a new store or to load an existing store.
func Start(logger logging.Logger, config *Config) (kvstore.TableStore, error) {
	if config == nil {
		return nil, errors.New("config is required")
	}

	base, err := buildBaseStore(config.Type, logger, config.Path)
	if err != nil {
		return nil, fmt.Errorf("error building base store: %w", err)
	}

	return start(logger, base, config)
}

// Future work: if we ever decide to permit third parties to provide custom store implementations not in this module,
// we will need to make this method public.

// start creates a new TableStore instance with the given base store and table names. If modifySchema is true, the
// tables in the store are made to match the given list of tables by adding and removing tables as needed. If
// modifySchema is false, the tables are loaded as is, and any tables in the provided list are ignored.
func start(
	logger logging.Logger,
	base kvstore.Store[[]byte],
	config *Config) (kvstore.TableStore, error) {

	metadataKeyBuilder := newKeyBuilder("metadata", metadataTableID)
	namespaceKeyBuilder := newKeyBuilder("namespace", namespaceTableID)
	expirationKeyBuilder := newKeyBuilder("expiration", expirationTableID)

	err := validateSchema(base, metadataKeyBuilder)
	if err != nil {
		return nil, fmt.Errorf("error validating schema: %w", err)
	}

	err = handleIncompleteDeletion(logger, base, metadataKeyBuilder, namespaceKeyBuilder)
	if err != nil {
		return nil, fmt.Errorf("error handling incomplete deletion: %w", err)
	}

	tableIDMap, err := loadNamespaceTable(base, namespaceKeyBuilder)
	if err != nil {
		return nil, fmt.Errorf("error loading namespace table: %w", err)
	}

	if config.Schema != nil {
		err = addAndRemoveTables(
			base,
			metadataKeyBuilder,
			namespaceKeyBuilder,
			tableIDMap,
			config.Schema)
		if err != nil {
			return nil, fmt.Errorf("error adding and removing tables: %w", err)
		}
	}

	store := newTableStore(
		logger,
		base,
		tableIDMap,
		expirationKeyBuilder,
		config.GarbageCollectionEnabled,
		config.GarbageCollectionInterval,
		config.GarbageCollectionBatchSize)

	return store, nil
}

// buildBaseStore creates a new base store of the given type.
func buildBaseStore(
	storeType StoreType,
	logger logging.Logger,
	path *string) (kvstore.Store[[]byte], error) {

	switch storeType {
	case LevelDB:
		if path == nil {
			return nil, errors.New("path is required for LevelDB store")
		}
		return leveldb.NewStore(logger, *path)
	case MapStore:
		return mapstore.NewStore(), nil
	default:
		return nil, fmt.Errorf("unknown store type: %d", storeType)
	}
}

// validateSchema loads/initiates the schema version in the metadata table.
func validateSchema(
	base kvstore.Store[[]byte],
	metadataKeyBuilder kvstore.KeyBuilder) error {

	schemaKey := metadataKeyBuilder.Key([]byte(metadataSchemaVersionKey))
	onDiskSchemaBytes, err := base.Get(schemaKey.Raw())

	if err != nil {
		if !errors.Is(err, kvstore.ErrNotFound) {
			return err
		}

		// This store is new, no on disk schema version exists.
		onDiskSchemaBytes = make([]byte, 8)
		binary.BigEndian.PutUint64(onDiskSchemaBytes, currentSchemaVersion)

		err = base.Put(schemaKey.Raw(), onDiskSchemaBytes)
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
	base kvstore.Store[[]byte],
	metadataKeyBuilder kvstore.KeyBuilder,
	namespaceKeyBuilder kvstore.KeyBuilder,
	tableIDMap map[uint32]string,
	currentTables []string) error {

	tablesToAdd, tablesToDrop := computeSchemaChange(tableIDMap, currentTables)

	err := dropTables(base, metadataKeyBuilder, namespaceKeyBuilder, tableIDMap, tablesToDrop)
	if err != nil {
		return fmt.Errorf("error dropping tables: %w", err)
	}

	err = addTables(base, metadataKeyBuilder, namespaceKeyBuilder, tableIDMap, tablesToAdd)
	if err != nil {
		return fmt.Errorf("error adding tables: %w", err)
	}

	err = sanityCheckNamespaceTable(base, namespaceKeyBuilder, tableIDMap, currentTables)
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
	base kvstore.Store[[]byte],
	metadataKeyBuilder kvstore.KeyBuilder,
	namespaceKeyBuilder kvstore.KeyBuilder,
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
		err := dropTable(base, metadataKeyBuilder, namespaceKeyBuilder, tableName, reverseTableIDMap[tableName])
		if err != nil {
			return fmt.Errorf("error dropping table %s: %w", tableName, err)
		}
		delete(tableIDMap, reverseTableIDMap[tableName])
	}

	return nil
}

// Add tables to the store. Updates the table ID map as well as data within the store.
func addTables(
	base kvstore.Store[[]byte],
	metadataKeyBuilder kvstore.KeyBuilder,
	namespaceKeyBuilder kvstore.KeyBuilder,
	tableIDMap map[uint32]string,
	tablesToAdd []string) error {

	if len(tablesToAdd) == 0 {
		// bail out early
		return nil
	}

	// Sort tables to add. Ensures deterministic table IDs given the same input.
	sort.Strings(tablesToAdd)

	// Add new tables.
	for _, tableName := range tablesToAdd {
		tableID, err := createTable(base, metadataKeyBuilder, namespaceKeyBuilder, tableName)
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
	base kvstore.Store[[]byte],
	namespaceKeyBuilder kvstore.KeyBuilder,
	tableIDMap map[uint32]string,
	currentTableList []string) error {

	// TODO also check that no table has ID greater than next table ID

	parsedNamespaceTable, err := loadNamespaceTable(base, namespaceKeyBuilder)
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
	base kvstore.Store[[]byte],
	metadataKeyBuilder kvstore.KeyBuilder,
	namespaceKeyBuilder kvstore.KeyBuilder) error {

	deletionKey := metadataKeyBuilder.Key([]byte(metadataDeletionKey))
	deletionValue, err := base.Get(deletionKey.Raw())
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
	return dropTable(base, metadataKeyBuilder, namespaceKeyBuilder, deletionTableName, deletionTableID)
}

// Get the prefix for the given table ID and prefix length.
func getPrefix(tableID uint32) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, tableID)
	return bytes
}

// loadNamespaceTable loads the namespace table from disk into the given map.
// Returns a map from table IDs to table names.
func loadNamespaceTable(
	base kvstore.Store[[]byte],
	namespaceKeyBuilder kvstore.KeyBuilder) (map[uint32]string, error) {

	tableIDMap := make(map[uint32]string)

	it, err := base.NewIterator(namespaceKeyBuilder.Key([]byte{}).Raw())
	if err != nil {
		return nil, fmt.Errorf("error creating namespace table iterator: %w", err)
	}
	defer it.Release()

	for it.Next() {
		keyBytes := it.Key()
		valueBytes := it.Value()

		tableID := binary.BigEndian.Uint32(keyBytes[prefixLength:])
		tableName := string(valueBytes)
		tableIDMap[tableID] = tableName
	}
	return tableIDMap, nil
}

// createTable creates a new table with the given name.
// Returns ErrTableLimitExceeded if the maximum number of tables has been reached.
func createTable(
	base kvstore.Store[[]byte],
	metadataKeyBuilder kvstore.KeyBuilder,
	namespaceKeyBuilder kvstore.KeyBuilder,
	name string) (uint32, error) {

	batch := base.NewBatch()

	var tableID uint32
	tableIDBytes, err := base.Get(metadataKeyBuilder.Key([]byte(nextTableIDKey)).Raw())
	if err != nil {
		if errors.Is(err, kvstore.ErrNotFound) {
			tableIDBytes = make([]byte, 4)
		} else {
			return 0, fmt.Errorf("error reading next table ID: %w", err)
		}
	}
	tableID = binary.BigEndian.Uint32(tableIDBytes)

	if tableID == expirationTableID {
		return 0, errors.New("table limit exceeded")
	}

	nextTableID := tableID + 1
	nextTableIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(nextTableIDBytes, nextTableID)

	keyBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(keyBytes, tableID)
	batch.Put(namespaceKeyBuilder.Key(keyBytes).Raw(), []byte(name))
	batch.Put(metadataKeyBuilder.Key([]byte(nextTableIDKey)).Raw(), nextTableIDBytes)

	err = batch.Apply()
	if err != nil {
		return 0, fmt.Errorf("error updating namespace table: %w", err)
	}

	return tableID, nil
}

// DropTable deletes the table with the given name. This is a no-op if the table does not exist.
func dropTable(
	base kvstore.Store[[]byte],
	metadataKeyBuilder kvstore.KeyBuilder,
	namespaceKeyBuilder kvstore.KeyBuilder,
	name string,
	tableID uint32) error {

	// This single atomic operation ensures that the table is deleted completely, even if we crash
	// in the middle of the operation. When next starting up, if an entry is observed in this location,
	// then the interrupted deletion can be completed.
	deletionKey := metadataKeyBuilder.Key([]byte(metadataDeletionKey))

	deletionValue := make([]byte, 4+len(name))
	binary.BigEndian.PutUint32(deletionValue, tableID)
	copy(deletionValue[4:], name)

	err := base.Put(deletionKey.Raw(), deletionValue)
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
	keyBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(keyBytes, tableID)
	tableKey := namespaceKeyBuilder.Key(keyBytes)
	err = base.Delete(tableKey.Raw())
	if err != nil {
		return fmt.Errorf("error deleting from namespace table: %w", err)
	}

	// Finally, remove the deletion key from the metadata table.
	return base.Delete(deletionKey.Raw())
}
