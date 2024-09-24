package tablestore

import (
	"encoding/binary"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/Layr-Labs/eigenda/common/kvstore/leveldb"
	"github.com/Layr-Labs/eigenda/common/kvstore/mapstore"
	"github.com/stretchr/testify/assert"
	"math"
	"os"
	"sort"
	"testing"
)

var dbPath = "test-store"

func deleteDBDirectory(t *testing.T) {
	err := os.RemoveAll(dbPath)
	assert.NoError(t, err)
}

func verifyDBIsDeleted(t *testing.T) {
	_, err := os.Stat(dbPath)
	assert.True(t, os.IsNotExist(err))
}

func TestTableCount(t *testing.T) {
	base := mapstore.NewStore()
	store, err := TableStoreWrapper(base)
	assert.NoError(t, err)

	// table count needs to fit into 32 bytes, and two tables are reserved for internal use
	maxTableCount := uint32(math.MaxUint32 - 2)
	assert.Equal(t, maxTableCount, store.GetMaxTableCount())
	assert.Equal(t, uint32(0), store.GetTableCount())

	_, _, err = store.GetOrCreateTable("table1")
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), store.GetTableCount())

	_, _, err = store.GetOrCreateTable("table2")
	assert.NoError(t, err)
	assert.Equal(t, uint32(2), store.GetTableCount())

	_, _, err = store.GetOrCreateTable("table3")
	assert.NoError(t, err)
	assert.Equal(t, uint32(3), store.GetTableCount())

	err = store.Destroy()
	assert.NoError(t, err)
}

func TestTableList(t *testing.T) {
	deleteDBDirectory(t)

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	base, err := leveldb.NewStore(logger, dbPath)
	assert.NoError(t, err)
	store, err := TableStoreWrapper(base)
	assert.NoError(t, err)

	tables := store.GetTables()
	assert.Equal(t, 0, len(tables))

	// Add some tables

	_, _, err = store.GetOrCreateTable("table1")

	tables = store.GetTables()
	assert.Equal(t, 1, len(tables))
	assert.Equal(t, "table1", tables[0])

	_, _, err = store.GetOrCreateTable("table2")
	assert.NoError(t, err)
	tables = store.GetTables()
	assert.Equal(t, 2, len(tables))
	sort.Strings(tables)
	assert.Equal(t, "table1", tables[0])
	assert.Equal(t, "table2", tables[1])

	_, _, err = store.GetOrCreateTable("table3")
	assert.NoError(t, err)
	tables = store.GetTables()
	assert.Equal(t, 3, len(tables))
	sort.Strings(tables)
	assert.Equal(t, "table1", tables[0])
	assert.Equal(t, "table2", tables[1])
	assert.Equal(t, "table3", tables[2])

	// Duplicate table additions should be no-ops
	_, _, err = store.GetOrCreateTable("table1")
	assert.NoError(t, err)
	tables = store.GetTables()
	assert.Equal(t, 3, len(tables))
	sort.Strings(tables)
	assert.Equal(t, "table1", tables[0])
	assert.Equal(t, "table2", tables[1])
	assert.Equal(t, "table3", tables[2])

	// Tables should survive a restart
	err = store.Shutdown()
	assert.NoError(t, err)

	base, err = leveldb.NewStore(logger, dbPath)
	assert.NoError(t, err)
	store, err = TableStoreWrapper(base)
	assert.NoError(t, err)

	tables = store.GetTables()
	assert.Equal(t, 3, len(tables))
	sort.Strings(tables)
	assert.Equal(t, "table1", tables[0])
	assert.Equal(t, "table2", tables[1])
	assert.Equal(t, "table3", tables[2])

	// Delete a table

	err = store.DropTable("table2")
	assert.NoError(t, err)
	tables = store.GetTables()
	assert.Equal(t, 2, len(tables))
	sort.Strings(tables)
	assert.Equal(t, "table1", tables[0])

	// Table should be deleted after a restart
	err = store.Shutdown()
	assert.NoError(t, err)

	base, err = leveldb.NewStore(logger, dbPath)
	assert.NoError(t, err)
	store, err = TableStoreWrapper(base)
	assert.NoError(t, err)

	tables = store.GetTables()
	assert.Equal(t, 2, len(tables))
	sort.Strings(tables)
	assert.Equal(t, "table1", tables[0])
	assert.Equal(t, "table3", tables[1])

	// Add a table back in (this uses a different code path)
	_, _, err = store.GetOrCreateTable("table4")
	assert.NoError(t, err)
	tables = store.GetTables()
	assert.Equal(t, 3, len(tables))
	sort.Strings(tables)
	assert.Equal(t, "table1", tables[0])
	assert.Equal(t, "table3", tables[1])
	assert.Equal(t, "table4", tables[2])

	// Delete the rest of the tables
	err = store.DropTable("table1")
	assert.NoError(t, err)
	err = store.DropTable("table3")
	assert.NoError(t, err)
	err = store.DropTable("table4")
	assert.NoError(t, err)

	tables = store.GetTables()
	assert.Equal(t, 0, len(tables))

	err = store.Destroy()
	assert.NoError(t, err)
	verifyDBIsDeleted(t)
}

func TestUniqueKeySpace(t *testing.T) {
	base := mapstore.NewStore()
	store, err := TableStoreWrapper(base)
	assert.NoError(t, err)

	kb1, table1, err := store.GetOrCreateTable("table1")
	assert.NoError(t, err)
	kb2, table2, err := store.GetOrCreateTable("table2")
	assert.NoError(t, err)

	// Write to keys with the same name in different tables using table store views

	err = table1.Put([]byte("key1"), []byte("value1"))
	assert.NoError(t, err)
	err = table2.Put([]byte("key1"), []byte("value2"))
	assert.NoError(t, err)

	value, err := table1.Get([]byte("key1"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("value1"), value)

	value, err = table2.Get([]byte("key1"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("value2"), value)

	// Write to keys with the same name in different tables using key builders

	err = store.Put(kb1.StringKey("key2"), []byte("value3"))
	assert.NoError(t, err)
	err = store.Put(kb2.StringKey("key"), []byte("value4"))
	assert.NoError(t, err)

	value, err = store.Get(kb1.StringKey("key2"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("value3"), value)

	value, err = store.Get(kb2.StringKey("key"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("value4"), value)

	// Delete a key from one table but not the other using table store views

	err = table1.Delete([]byte("key1"))
	assert.NoError(t, err)

	value, err = table1.Get([]byte("key1"))
	assert.Equal(t, kvstore.ErrNotFound, err)

	value, err = table2.Get([]byte("key1"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("value2"), value)

	// Delete a key from one table but not the other using key builders

	err = store.Delete(kb1.StringKey("key2"))
	assert.NoError(t, err)

	value, err = store.Get(kb1.StringKey("key2"))
	assert.Equal(t, kvstore.ErrNotFound, err)

	value, err = store.Get(kb2.StringKey("key"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("value4"), value)
}

func TestBatchOperations(t *testing.T) {
	base := mapstore.NewStore()
	store, err := TableStoreWrapper(base)
	assert.NoError(t, err)

	kb1, _, err := store.GetOrCreateTable("table1")
	assert.NoError(t, err)

	kb2, _, err := store.GetOrCreateTable("table2")
	assert.NoError(t, err)

	kb3, _, err := store.GetOrCreateTable("table3")
	assert.NoError(t, err)

	keys := make([]kvstore.Key, 0)

	for i := 0; i < 10; i++ {
		keys = append(keys, kb1.Uint64Key(uint64(i)))
	}
	for i := 0; i < 10; i++ {
		keys = append(keys, kb2.Uint64Key(uint64(i)))
	}
	for i := 0; i < 10; i++ {
		keys = append(keys, kb3.Uint64Key(uint64(i)))
	}

	values := make([][]byte, 0)
	for i := 0; i < 30; i++ {
		valueBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(valueBytes, uint64(i))
		values = append(values, valueBytes)
	}

	err = store.WriteBatch(keys, values)
	assert.NoError(t, err)

	for i := 0; i < 10; i++ {
		value, err := store.Get(kb1.Uint64Key(uint64(i)))
		assert.NoError(t, err)
		assert.Equal(t, uint64(i), binary.BigEndian.Uint64(value))
	}
	for i := 0; i < 10; i++ {
		value, err := store.Get(kb2.Uint64Key(uint64(i)))
		assert.NoError(t, err)
		assert.Equal(t, uint64(i+10), binary.BigEndian.Uint64(value))
	}
	for i := 0; i < 10; i++ {
		value, err := store.Get(kb3.Uint64Key(uint64(i)))
		assert.NoError(t, err)
		assert.Equal(t, uint64(i+20), binary.BigEndian.Uint64(value))
	}

	// Delete odd keys
	keys = make([]kvstore.Key, 0)
	for i := 0; i < 10; i++ {
		if i%2 == 1 {
			keys = append(keys, kb1.Uint64Key(uint64(i)))
		}
	}
	for i := 0; i < 10; i++ {
		if i%2 == 1 {
			keys = append(keys, kb2.Uint64Key(uint64(i)))
		}
	}
	for i := 0; i < 10; i++ {
		if i%2 == 1 {
			keys = append(keys, kb3.Uint64Key(uint64(i)))
		}
	}

	err = store.DeleteBatch(keys)
	assert.NoError(t, err)

	for i := 0; i < 10; i++ {
		if i%2 == 1 {
			_, err := store.Get(kb1.Uint64Key(uint64(i)))
			assert.Equal(t, kvstore.ErrNotFound, err)
		} else {
			value, err := store.Get(kb1.Uint64Key(uint64(i)))
			assert.NoError(t, err)
			assert.Equal(t, uint64(i), binary.BigEndian.Uint64(value))
		}
	}
	for i := 0; i < 10; i++ {
		if i%2 == 1 {
			_, err := store.Get(kb2.Uint64Key(uint64(i)))
			assert.Equal(t, kvstore.ErrNotFound, err)
		} else {
			value, err := store.Get(kb2.Uint64Key(uint64(i)))
			assert.NoError(t, err)
			assert.Equal(t, uint64(i+10), binary.BigEndian.Uint64(value))
		}
	}
	for i := 0; i < 10; i++ {
		if i%2 == 1 {
			_, err := store.Get(kb3.Uint64Key(uint64(i)))
			assert.Equal(t, kvstore.ErrNotFound, err)
		} else {
			value, err := store.Get(kb3.Uint64Key(uint64(i)))
			assert.NoError(t, err)
			assert.Equal(t, uint64(i+20), binary.BigEndian.Uint64(value))
		}
	}

	err = store.Destroy()
	assert.NoError(t, err)
}

// TODO
//  - drop table
//  - interrupted drop table
//  - create table inside gap, look inside to make sure things happen as expected
//  - testing the different methods on key builders and keys
//  - iteration
//  - random operations test
