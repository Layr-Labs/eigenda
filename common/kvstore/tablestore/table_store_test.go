package tablestore

import (
	"encoding/binary"
	"fmt"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/stretchr/testify/assert"
	"math"
	"os"
	"sort"
	"strings"
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
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	builder, err := MapStore.Builder(logger, "")
	assert.NoError(t, err)

	// table count needs to fit into 32 bytes, and two tables are reserved for internal use
	maxTableCount := uint32(math.MaxUint32 - 2)
	assert.Equal(t, maxTableCount, builder.GetMaxTableCount())
	assert.Equal(t, uint32(0), builder.GetTableCount())

	err = builder.CreateTable("table1")
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), builder.GetTableCount())

	err = builder.CreateTable("table2")
	assert.NoError(t, err)
	assert.Equal(t, uint32(2), builder.GetTableCount())

	err = builder.CreateTable("table3")
	assert.NoError(t, err)
	assert.Equal(t, uint32(3), builder.GetTableCount())

	store, err := builder.Build()
	assert.NoError(t, err)

	table1, err := store.GetTable("table1")
	assert.NoError(t, err)
	assert.Equal(t, "table1", table1.Name())

	table2, err := store.GetTable("table2")
	assert.NoError(t, err)
	assert.Equal(t, "table2", table2.Name())

	table3, err := store.GetTable("table3")
	assert.NoError(t, err)
	assert.Equal(t, "table3", table3.Name())

	err = store.Destroy()
	assert.NoError(t, err)
}

func TestTableList(t *testing.T) {
	deleteDBDirectory(t)

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	builder, err := LevelDB.Builder(logger, dbPath)
	assert.NoError(t, err)

	tableNames := builder.GetTableNames()
	assert.Equal(t, 0, len(tableNames))

	// Add some tables

	err = builder.CreateTable("table1")
	assert.NoError(t, err)

	tableNames = builder.GetTableNames()
	assert.Equal(t, 1, len(tableNames))
	assert.Equal(t, "table1", tableNames[0])

	err = builder.CreateTable("table2")
	assert.NoError(t, err)
	tableNames = builder.GetTableNames()
	assert.Equal(t, 2, len(tableNames))
	sort.SliceStable(tableNames, func(i, j int) bool {
		return tableNames[i] < tableNames[j]
	})
	assert.Equal(t, "table1", tableNames[0])
	assert.Equal(t, "table2", tableNames[1])

	err = builder.CreateTable("table3")
	assert.NoError(t, err)
	tableNames = builder.GetTableNames()
	assert.Equal(t, 3, len(tableNames))
	sort.SliceStable(tableNames, func(i, j int) bool {
		return tableNames[i] < tableNames[j]
	})
	assert.Equal(t, "table1", tableNames[0])
	assert.Equal(t, "table2", tableNames[1])
	assert.Equal(t, "table3", tableNames[2])

	// Duplicate table additions should be no-ops
	err = builder.CreateTable("table1")
	assert.NoError(t, err)
	tableNames = builder.GetTableNames()
	assert.Equal(t, 3, len(tableNames))
	sort.SliceStable(tableNames, func(i, j int) bool {
		return tableNames[i] < tableNames[j]
	})
	assert.Equal(t, "table1", tableNames[0])
	assert.Equal(t, "table2", tableNames[1])
	assert.Equal(t, "table3", tableNames[2])

	// Build the store and check the list of tables.
	store, err := builder.Build()
	assert.NoError(t, err)

	tables := store.GetTables()
	assert.Equal(t, 3, len(tables))
	sort.SliceStable(tables, func(i, j int) bool {
		return strings.Compare(tables[i].Name(), tables[j].Name()) < 0
	})
	assert.Equal(t, "table1", tables[0].Name())
	assert.Equal(t, "table2", tables[1].Name())
	assert.Equal(t, "table3", tables[2].Name())

	// Tables should survive a restart
	err = store.Shutdown()
	assert.NoError(t, err)

	builder, err = LevelDB.Builder(logger, dbPath)
	assert.NoError(t, err)

	tableNames = builder.GetTableNames()
	assert.Equal(t, 3, len(tableNames))
	sort.SliceStable(tableNames, func(i, j int) bool {
		return tableNames[i] < tableNames[j]
	})
	assert.Equal(t, "table1", tableNames[0])
	assert.Equal(t, "table2", tableNames[1])
	assert.Equal(t, "table3", tableNames[2])

	// Delete a table

	err = builder.DropTable("table2")
	assert.NoError(t, err)
	tableNames = builder.GetTableNames()
	assert.Equal(t, 2, len(tableNames))
	sort.SliceStable(tableNames, func(i, j int) bool {
		return tableNames[i] < tableNames[j]
	})
	assert.Equal(t, "table1", tableNames[0])

	// Build the store and check the list of tables.

	store, err = builder.Build()
	assert.NoError(t, err)

	tables = store.GetTables()
	assert.Equal(t, 2, len(tables))
	sort.SliceStable(tables, func(i, j int) bool {
		return strings.Compare(tables[i].Name(), tables[j].Name()) < 0
	})
	assert.Equal(t, "table1", tables[0].Name())
	assert.Equal(t, "table3", tables[1].Name())

	// Table should be deleted after a restart
	err = store.Shutdown()
	assert.NoError(t, err)

	builder, err = LevelDB.Builder(logger, dbPath)
	assert.NoError(t, err)

	tableNames = builder.GetTableNames()
	assert.Equal(t, 2, len(tableNames))
	sort.SliceStable(tableNames, func(i, j int) bool {
		return tableNames[i] < tableNames[j]
	})
	assert.Equal(t, "table1", tableNames[0])
	assert.Equal(t, "table3", tableNames[1])

	// Add a table back in (this uses a different code path)
	err = builder.CreateTable("table4")
	assert.NoError(t, err)
	tableNames = builder.GetTableNames()
	assert.Equal(t, 3, len(tableNames))
	sort.SliceStable(tableNames, func(i, j int) bool {
		return tableNames[i] < tableNames[j]
	})
	assert.Equal(t, "table1", tableNames[0])
	assert.Equal(t, "table3", tableNames[1])
	assert.Equal(t, "table4", tableNames[2])

	// Delete the rest of the tables
	err = builder.DropTable("table1")
	assert.NoError(t, err)
	err = builder.DropTable("table3")
	assert.NoError(t, err)
	err = builder.DropTable("table4")
	assert.NoError(t, err)

	tableNames = builder.GetTableNames()
	assert.Equal(t, 0, len(tableNames))

	err = builder.Destroy()
	assert.NoError(t, err)
	verifyDBIsDeleted(t)
}

func TestUniqueKeySpace(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	builder, err := MapStore.Builder(logger, dbPath)
	assert.NoError(t, err)

	err = builder.CreateTable("table1")
	assert.NoError(t, err)
	err = builder.CreateTable("table2")
	assert.NoError(t, err)

	store, err := builder.Build()
	assert.NoError(t, err)

	table1, err := store.GetTable("table1")
	assert.NoError(t, err)

	table2, err := store.GetTable("table2")
	assert.NoError(t, err)

	// Write to the tables

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

	// Delete a key from one table but not the other

	err = table1.Delete([]byte("key1"))
	assert.NoError(t, err)

	_, err = table1.Get([]byte("key1"))
	assert.Equal(t, kvstore.ErrNotFound, err)

	value, err = table2.Get([]byte("key1"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("value2"), value)
}

func TestBatchOperations(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	builder, err := MapStore.Builder(logger, dbPath)
	assert.NoError(t, err)

	err = builder.CreateTable("table1")
	assert.NoError(t, err)
	err = builder.CreateTable("table2")
	assert.NoError(t, err)
	err = builder.CreateTable("table3")
	assert.NoError(t, err)

	store, err := builder.Build()
	assert.NoError(t, err)

	table1, err := store.GetTable("table1")
	assert.NoError(t, err)

	table2, err := store.GetTable("table2")
	assert.NoError(t, err)

	table3, err := store.GetTable("table3")
	assert.NoError(t, err)

	// Test a batch with just puts

	batch := store.NewBatch()
	for i := 0; i < 10; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		v := make([]byte, 8)
		binary.BigEndian.PutUint64(v, uint64(i))
		batch.Put(table1.TableKey(k), v)

		v = make([]byte, 8)
		binary.BigEndian.PutUint64(v, uint64(i+10))
		batch.Put(table2.TableKey(k), v)

		v = make([]byte, 8)
		binary.BigEndian.PutUint64(v, uint64(i+20))
		batch.Put(table3.TableKey(k), v)
	}

	err = batch.Apply()
	assert.NoError(t, err)

	for i := 0; i < 10; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		value, err := table1.Get(k)
		assert.NoError(t, err)
		assert.Equal(t, uint64(i), binary.BigEndian.Uint64(value))

		value, err = table2.Get(k)
		assert.NoError(t, err)
		assert.Equal(t, uint64(i+10), binary.BigEndian.Uint64(value))

		value, err = table3.Get(k)
		assert.NoError(t, err)
		assert.Equal(t, uint64(i+20), binary.BigEndian.Uint64(value))
	}

	// Test a batch with just deletes

	// Delete odd keys
	batch = store.NewBatch()
	for i := 0; i < 10; i++ {
		if i%2 == 1 {
			k := make([]byte, 8)
			binary.BigEndian.PutUint64(k, uint64(i))

			batch.Delete(table1.TableKey(k))
			batch.Delete(table2.TableKey(k))
			batch.Delete(table3.TableKey(k))
		}
	}

	err = batch.Apply()
	assert.NoError(t, err)

	for i := 0; i < 10; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		if i%2 == 1 {
			_, err = table1.Get(k)
			assert.Equal(t, kvstore.ErrNotFound, err)

			_, err = table2.Get(k)
			assert.Equal(t, kvstore.ErrNotFound, err)

			_, err = table3.Get(k)
			assert.Equal(t, kvstore.ErrNotFound, err)
		} else {
			value, err := table1.Get(k)
			assert.NoError(t, err)
			assert.Equal(t, uint64(i), binary.BigEndian.Uint64(value))

			value, err = table2.Get(k)
			assert.NoError(t, err)
			assert.Equal(t, uint64(i+10), binary.BigEndian.Uint64(value))

			value, err = table3.Get(k)
			assert.NoError(t, err)
			assert.Equal(t, uint64(i+20), binary.BigEndian.Uint64(value))
		}
	}

	// Test a batch with a mix of puts and deletes

	// Write back in odd numbers, but delete multiples of 4
	batch = store.NewBatch()
	for i := 0; i < 10; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))
		if i%4 == 0 {
			batch.Delete(table1.TableKey(k))
			batch.Delete(table2.TableKey(k))
			batch.Delete(table3.TableKey(k))
		} else if i%2 == 1 {
			v := make([]byte, 8)
			binary.BigEndian.PutUint64(v, uint64(2*i))
			batch.Put(table1.TableKey(k), v)

			v = make([]byte, 8)
			binary.BigEndian.PutUint64(v, uint64(2*i+10))
			batch.Put(table2.TableKey(k), v)

			v = make([]byte, 8)
			binary.BigEndian.PutUint64(v, uint64(2*i+20))
			batch.Put(table3.TableKey(k), v)
		}
	}

	err = batch.Apply()
	assert.NoError(t, err)

	for i := 0; i < 10; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		if i%4 == 0 {
			_, err = table1.Get(k)
			assert.Equal(t, kvstore.ErrNotFound, err)

			_, err = table2.Get(k)
			assert.Equal(t, kvstore.ErrNotFound, err)

			_, err = table3.Get(k)
			assert.Equal(t, kvstore.ErrNotFound, err)
		} else if i%2 == 1 {
			val, err := table1.Get(k)
			assert.NoError(t, err)
			assert.Equal(t, uint64(2*i), binary.BigEndian.Uint64(val))

			val, err = table2.Get(k)
			assert.NoError(t, err)
			assert.Equal(t, uint64(2*i+10), binary.BigEndian.Uint64(val))

			val, err = table3.Get(k)
			assert.NoError(t, err)
			assert.Equal(t, uint64(2*i+20), binary.BigEndian.Uint64(val))
		} else {
			val, err := table1.Get(k)
			assert.NoError(t, err)
			assert.Equal(t, uint64(i), binary.BigEndian.Uint64(val))

			val, err = table2.Get(k)
			assert.NoError(t, err)
			assert.Equal(t, uint64(i+10), binary.BigEndian.Uint64(val))

			val, err = table3.Get(k)
			assert.NoError(t, err)
			assert.Equal(t, uint64(i+20), binary.BigEndian.Uint64(val))
		}
	}

	err = store.Destroy()
	assert.NoError(t, err)
}

func TestDropTable(t *testing.T) {
	deleteDBDirectory(t)

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	builder, err := LevelDB.Builder(logger, dbPath)
	assert.NoError(t, err)

	err = builder.CreateTable("table1")
	assert.NoError(t, err)
	err = builder.CreateTable("table2")
	assert.NoError(t, err)
	err = builder.CreateTable("table3")
	assert.NoError(t, err)

	store, err := builder.Build()
	assert.NoError(t, err)

	table1, err := store.GetTable("table1")
	assert.NoError(t, err)

	table2, err := store.GetTable("table2")
	assert.NoError(t, err)

	table3, err := store.GetTable("table3")
	assert.NoError(t, err)

	// Insert some data into the tables
	for i := 0; i < 100; i++ {
		value := make([]byte, 8)
		binary.BigEndian.PutUint64(value, uint64(i))

		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		err = table1.Put(k, value)
		assert.NoError(t, err)

		err = table2.Put(k, value)
		assert.NoError(t, err)

		err = table3.Put(k, value)
		assert.NoError(t, err)
	}

	// Verify the data is there
	for i := 0; i < 100; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		expectedValue := make([]byte, 8)
		binary.BigEndian.PutUint64(expectedValue, uint64(i))

		value, err := table1.Get(k)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)

		value, err = table2.Get(k)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)

		_, err = table3.Get(k)
		assert.NoError(t, err)
	}

	// In order to drop a table, we will need to close the store and reopen it with a builder.
	// Only the builder can drop tables.
	err = store.Shutdown()
	assert.NoError(t, err)

	builder, err = LevelDB.Builder(logger, dbPath)
	assert.NoError(t, err)

	err = builder.DropTable("table2")
	assert.NoError(t, err)

	store, err = builder.Build()
	assert.NoError(t, err)

	table1, err = store.GetTable("table1")
	assert.NoError(t, err)

	_, err = store.GetTable("table2")
	assert.Equal(t, kvstore.ErrTableNotFound, err)

	table3, err = store.GetTable("table3")
	assert.NoError(t, err)

	for i := 0; i < 100; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		expectedValue := make([]byte, 8)
		binary.BigEndian.PutUint64(expectedValue, uint64(i))

		value, err := table1.Get(k)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)

		value, err = table3.Get(k)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)
	}

	// Restart the store so that we can drop another table.
	err = store.Shutdown()
	assert.NoError(t, err)

	builder, err = LevelDB.Builder(logger, dbPath)
	assert.NoError(t, err)

	err = builder.DropTable("table1")
	assert.NoError(t, err)

	store, err = builder.Build()
	assert.NoError(t, err)

	_, err = store.GetTable("table1")
	assert.Equal(t, kvstore.ErrTableNotFound, err)

	_, err = store.GetTable("table2")
	assert.Equal(t, kvstore.ErrTableNotFound, err)

	table3, err = store.GetTable("table3")
	assert.NoError(t, err)

	for i := 0; i < 100; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		expectedValue := make([]byte, 8)
		binary.BigEndian.PutUint64(expectedValue, uint64(i))

		value, err := table3.Get(k)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)
	}

	// Restart the store so that we can drop the last table.
	err = store.Shutdown()
	assert.NoError(t, err)

	builder, err = LevelDB.Builder(logger, dbPath)
	assert.NoError(t, err)

	err = builder.DropTable("table3")
	assert.NoError(t, err)

	store, err = builder.Build()
	assert.NoError(t, err)

	_, err = store.GetTable("table1")
	assert.Equal(t, kvstore.ErrTableNotFound, err)

	_, err = store.GetTable("table2")
	assert.Equal(t, kvstore.ErrTableNotFound, err)

	_, err = store.GetTable("table3")
	assert.Equal(t, kvstore.ErrTableNotFound, err)

	err = store.Destroy()
	assert.NoError(t, err)

	verifyDBIsDeleted(t)
}

func TestIteration(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	builder, err := MapStore.Builder(logger, dbPath)
	assert.NoError(t, err)

	err = builder.CreateTable("table1")
	assert.NoError(t, err)

	err = builder.CreateTable("table2")
	assert.NoError(t, err)

	store, err := builder.Build()
	assert.NoError(t, err)

	table1, err := store.GetTable("table1")
	assert.NoError(t, err)

	table2, err := store.GetTable("table2")
	assert.NoError(t, err)

	// Prefix "qwer"
	for i := 0; i < 100; i++ {
		k := []byte(fmt.Sprintf("qwer%3d", i))

		value := make([]byte, 8)
		binary.BigEndian.PutUint64(value, uint64(i))
		err = table1.Put(k, value)
		assert.NoError(t, err)

		value = make([]byte, 8)
		binary.BigEndian.PutUint64(value, uint64(2*i))
		err = table2.Put(k, value)
		assert.NoError(t, err)
	}

	// Prefix "asdf"
	for i := 0; i < 100; i++ {
		k := []byte(fmt.Sprintf("asdf%3d", i))

		value := make([]byte, 8)
		binary.BigEndian.PutUint64(value, uint64(i))
		err = table1.Put(k, value)
		assert.NoError(t, err)

		value = make([]byte, 8)
		binary.BigEndian.PutUint64(value, uint64(2*i))
		err = table2.Put(k, value)
		assert.NoError(t, err)
	}

	// Iterate table 1 with no prefix filter
	it, err := table1.NewIterator(nil)
	assert.NoError(t, err)

	count := 0
	for it.Next() {
		k := it.Key()
		v := it.Value()

		if count < 100 {
			// First we should see the keys with prefix "asdf" since they are lexicographically first
			expectedKey := []byte(fmt.Sprintf("asdf%3d", count))
			expectedValue := make([]byte, 8)
			binary.BigEndian.PutUint64(expectedValue, uint64(count))

			assert.Equal(t, expectedKey, k)
			assert.Equal(t, expectedValue, v)
		} else {
			// Then we should see the keys with prefix "qwer"
			adjustedCount := count - 100
			expectedKey := []byte(fmt.Sprintf("qwer%3d", adjustedCount))
			expectedValue := make([]byte, 8)
			binary.BigEndian.PutUint64(expectedValue, uint64(adjustedCount))

			assert.Equal(t, expectedKey, k)
			assert.Equal(t, expectedValue, v)
		}

		count++
	}
	it.Release()

	// Iterate table 2 with no prefix filter
	it, err = table2.NewIterator(nil)
	assert.NoError(t, err)

	count = 0
	for it.Next() {
		k := it.Key()
		v := it.Value()

		if count < 100 {
			// First we should see the keys with prefix "asdf" since they are lexicographically first
			expectedKey := []byte(fmt.Sprintf("asdf%3d", count))
			expectedValue := make([]byte, 8)
			binary.BigEndian.PutUint64(expectedValue, uint64(2*count))

			assert.Equal(t, expectedKey, k)
			assert.Equal(t, expectedValue, v)
		} else {
			// Then we should see the keys with prefix "qwer"
			adjustedCount := count - 100
			expectedKey := []byte(fmt.Sprintf("qwer%3d", adjustedCount))
			expectedValue := make([]byte, 8)
			binary.BigEndian.PutUint64(expectedValue, uint64(2*adjustedCount))

			assert.Equal(t, expectedKey, k)
			assert.Equal(t, expectedValue, v)
		}

		count++
	}
	it.Release()

	// Iterate over the "qwer" keys from table 1
	it, err = table1.NewIterator([]byte("qwer"))
	assert.NoError(t, err)

	count = 0
	for it.Next() {
		k := it.Key()
		v := it.Value()

		expectedKey := []byte(fmt.Sprintf("qwer%3d", count))
		expectedValue := make([]byte, 8)
		binary.BigEndian.PutUint64(expectedValue, uint64(count))

		assert.Equal(t, expectedKey, k)
		assert.Equal(t, expectedValue, v)

		count++
	}
	it.Release()

	// Iterate over the "asdf" keys from table 2
	it, err = table2.NewIterator([]byte("asdf"))
	assert.NoError(t, err)

	count = 0
	for it.Next() {
		k := it.Key()
		v := it.Value()

		expectedKey := []byte(fmt.Sprintf("asdf%3d", count))
		expectedValue := make([]byte, 8)
		binary.BigEndian.PutUint64(expectedValue, uint64(2*count))

		assert.Equal(t, expectedKey, k)
		assert.Equal(t, expectedValue, v)

		count++
	}
	it.Release()

	err = store.Destroy()
	assert.NoError(t, err)
}

func TestRestart(t *testing.T) {
	deleteDBDirectory(t)

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	builder, err := LevelDB.Builder(logger, dbPath)
	assert.NoError(t, err)

	err = builder.CreateTable("table1")
	assert.NoError(t, err)

	err = builder.CreateTable("table2")
	assert.NoError(t, err)

	store, err := builder.Build()
	assert.NoError(t, err)

	table1, err := store.GetTable("table1")
	assert.NoError(t, err)

	table2, err := store.GetTable("table2")
	assert.NoError(t, err)

	for i := 0; i < 100; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		value1 := make([]byte, 8)
		binary.BigEndian.PutUint64(value1, uint64(i))

		value2 := make([]byte, 8)
		binary.BigEndian.PutUint64(value2, uint64(i*2))

		err = table1.Put(k, value1)
		assert.NoError(t, err)

		err = table2.Put(k, value2)
		assert.NoError(t, err)
	}

	err = store.Shutdown()
	assert.NoError(t, err)

	builder, err = LevelDB.Builder(logger, dbPath)
	assert.NoError(t, err)

	store, err = builder.Build()
	assert.NoError(t, err)

	table1, err = store.GetTable("table1")
	assert.NoError(t, err)
	table2, err = store.GetTable("table2")
	assert.NoError(t, err)

	for i := 0; i < 100; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		expectedValue1 := make([]byte, 8)
		binary.BigEndian.PutUint64(expectedValue1, uint64(i))

		expectedValue2 := make([]byte, 8)
		binary.BigEndian.PutUint64(expectedValue2, uint64(i*2))

		value1, err := table1.Get(k)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue1, value1)

		value2, err := table2.Get(k)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue2, value2)
	}

	err = store.Destroy()
	assert.NoError(t, err)

	verifyDBIsDeleted(t)
}

// Maps keys (in string form) to expected values for a particular table.
type expectedTableData map[string][]byte

// Maps table names to expected data for that table.
type expectedStoreData map[string]expectedTableData

//func TestRandomOperations(t *testing.T) {
//	tu.InitializeRandom()
//
//	deleteDBDirectory(t)
//
//	logger, err := common.NewLogger(common.DefaultLoggerConfig())
//	assert.NoError(t, err)
//
//	base, err := leveldb.NewStore(logger, dbPath)
//	assert.NoError(t, err)
//	store, err := wrapper(logger, base)
//	assert.NoError(t, err)
//
//	tables := make(map[string]kvstore.Table)
//	expectedData := make(expectedStoreData)
//
//	for i := 0; i < 10000; i++ {
//
//		choice := rand.Float64()
//
//		if choice < 0.01 {
//			// restart the store
//			err = store.Shutdown()
//			assert.NoError(t, err)
//
//			base, err = leveldb.NewStore(logger, dbPath)
//			assert.NoError(t, err)
//			store, err = wrapper(logger, base)
//			assert.NoError(t, err)
//
//			for tableName := range tables {
//				table, err := store.GetTable(tableName)
//				assert.NoError(t, err)
//				tables[tableName] = table
//			}
//		} else if len(tables) == 0 || choice < 0.1 {
//			// Create a new table.
//			name := tu.RandomString(8)
//			table, err := store.GetTable(name)
//			expectedData[name] = make(expectedTableData)
//			assert.NoError(t, err)
//			tables[name] = table
//		} else if choice < 0.15 {
//			// Drop a table
//
//			var name string
//			for n := range tables {
//				name = n
//				break
//			}
//
//			err := store.DropTable(name)
//			assert.NoError(t, err)
//
//			// Delete all expected data for the table
//			delete(expectedData, name)
//			delete(tables, name)
//		} else if choice < 0.9 || len(expectedData) == 0 {
//			// Write a value
//
//			var tableName string
//			for n := range tables {
//				tableName = n
//				break
//			}
//			table := tables[tableName]
//
//			k := []byte(tu.RandomString(32))
//			v := tu.RandomBytes(32)
//
//			expectedData[tableName][string(k)] = v
//
//			err = table.Put(k, v)
//			assert.NoError(t, err)
//		} else {
//			// Delete a value
//
//			var tableName string
//			for n := range tables {
//				tableName = n
//				break
//			}
//			table := tables[tableName]
//
//			if len(expectedData[tableName]) == 0 {
//				// no data in this table, skip
//				continue
//			}
//
//			var k string
//			for k = range expectedData[tableName] {
//				break
//			}
//
//			delete(expectedData[tableName], k)
//			err = table.Delete([]byte(k))
//			assert.NoError(t, err)
//		}
//
//		// Every once in a while, check that the store matches the expected data
//		if i%100 == 0 {
//			// Every so often, check that the store matches the expected data.
//			for tableName, tableData := range expectedData {
//				table := tables[tableName]
//
//				for k := range tableData {
//					expectedValue := tableData[k]
//					value, err := table.Get([]byte(k))
//					assert.NoError(t, err)
//					assert.Equal(t, expectedValue, value)
//				}
//			}
//		}
//	}
//
//	err = store.Destroy()
//	assert.NoError(t, err)
//}

//var _ kvstore.Store = &explodingStore{}
//
//// explodingStore is a store that returns an error after a certain number of operations.
//// Used to intentionally crash table deletion to exercise table deletion recovery.
//type explodingStore struct {
//	base               kvstore.Store
//	deletionsRemaining int
//}
//
//func (e *explodingStore) Put(key []byte, value []byte) error {
//	return e.base.Put(key, value)
//}
//
//func (e *explodingStore) Get(key []byte) ([]byte, error) {
//	return e.base.Get(key)
//}
//
//func (e *explodingStore) Delete(key []byte) error {
//	if e.deletionsRemaining == 0 {
//		return fmt.Errorf("intentional error")
//	}
//	e.deletionsRemaining--
//	return e.base.Delete(key)
//}
//
//func (e *explodingStore) NewBatch() kvstore.Batch[[]byte] {
//	panic("not used")
//}
//
//func (e *explodingStore) NewIterator(prefix []byte) (iterator.Iterator, error) {
//	return e.base.NewIterator(prefix)
//}
//
//func (e *explodingStore) Shutdown() error {
//	return e.base.Shutdown()
//}
//
//func (e *explodingStore) Destroy() error {
//	return e.base.Destroy()
//}
//
//func TestInterruptedTableDeletion(t *testing.T) {
//	deleteDBDirectory(t)
//
//	logger, err := common.NewLogger(common.DefaultLoggerConfig())
//	assert.NoError(t, err)
//
//	base, err := leveldb.NewStore(logger, dbPath)
//	assert.NoError(t, err)
//
//	explodingBase := &explodingStore{
//		base:               base,
//		deletionsRemaining: 50,
//	}
//
//	store, err := wrapper(logger, explodingBase)
//	assert.NoError(t, err)
//
//	// Create a few tables
//	table1, err := store.GetTable("table1")
//	assert.NoError(t, err)
//
//	table2, err := store.GetTable("table2")
//	assert.NoError(t, err)
//
//	// Write some data to the tables
//	for i := 0; i < 100; i++ {
//		k := make([]byte, 8)
//		binary.BigEndian.PutUint64(k, uint64(i))
//
//		value := make([]byte, 8)
//		binary.BigEndian.PutUint64(value, uint64(i))
//
//		err = table1.Put(k, value)
//		assert.NoError(t, err)
//
//		err = table2.Put(k, value)
//		assert.NoError(t, err)
//	}
//
//	// Drop one of the tables. This should fail partway through.
//	err = store.DropTable("table1")
//	assert.Error(t, err)
//
//	err = store.Shutdown()
//	assert.NoError(t, err)
//
//	// Restart the store. The table should be gone by the time the method returns.
//	base, err = leveldb.NewStore(logger, dbPath)
//	assert.NoError(t, err)
//	store, err = wrapper(logger, base)
//	assert.NoError(t, err)
//
//	tables := store.GetTables()
//	assert.Equal(t, 1, len(tables))
//	assert.Equal(t, "table2", tables[0].Name())
//	table2, err = store.GetTable("table2")
//	assert.NoError(t, err)
//
//	// Check that the data in the remaining table is still there. We shouldn't see any data from the deleted table.
//	for i := 0; i < 100; i++ {
//		k := make([]byte, 8)
//		binary.BigEndian.PutUint64(k, uint64(i))
//		value, err := table2.Get(k)
//		assert.NoError(t, err)
//		assert.Equal(t, uint64(i), binary.BigEndian.Uint64(value))
//	}
//
//	err = store.Destroy()
//	assert.NoError(t, err)
//
//	verifyDBIsDeleted(t)
//}
