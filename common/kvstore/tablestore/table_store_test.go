package tablestore

import (
	"encoding/binary"
	"fmt"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/Layr-Labs/eigenda/common/kvstore/leveldb"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"math/rand"
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

func TestTableList(t *testing.T) {
	deleteDBDirectory(t)

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	tStore, err := LevelDB.Start(logger, dbPath)
	assert.NoError(t, err)

	tables := tStore.GetTables()
	assert.Equal(t, 0, len(tables))

	err = tStore.Shutdown()
	assert.NoError(t, err)

	// Add some tables

	tStore, err = LevelDB.Start(logger, dbPath, "table1")
	assert.NoError(t, err)

	tables = tStore.GetTables()
	assert.Equal(t, 1, len(tables))
	assert.Equal(t, "table1", tables[0].Name())

	err = tStore.Shutdown()
	assert.NoError(t, err)

	tStore, err = LevelDB.Start(logger, dbPath, "table1", "table2")
	assert.NoError(t, err)
	tables = tStore.GetTables()
	assert.Equal(t, 2, len(tables))
	sort.SliceStable(tables, func(i, j int) bool {
		return tables[i].Name() < tables[j].Name()
	})
	assert.Equal(t, "table1", tables[0].Name())
	assert.Equal(t, "table2", tables[1].Name())

	err = tStore.Shutdown()
	assert.NoError(t, err)

	tStore, err = LevelDB.Start(logger, dbPath, "table1", "table2", "table3")
	assert.NoError(t, err)
	tables = tStore.GetTables()
	assert.Equal(t, 3, len(tables))
	sort.SliceStable(tables, func(i, j int) bool {
		return tables[i].Name() < tables[j].Name()
	})
	assert.Equal(t, "table1", tables[0].Name())
	assert.Equal(t, "table2", tables[1].Name())
	assert.Equal(t, "table3", tables[2].Name())

	err = tStore.Shutdown()
	assert.NoError(t, err)

	// Restarting with the same tables should work.
	tStore, err = LevelDB.Start(logger, dbPath, "table1", "table2", "table3")
	assert.NoError(t, err)
	tables = tStore.GetTables()
	assert.Equal(t, 3, len(tables))
	sort.SliceStable(tables, func(i, j int) bool {
		return tables[i].Name() < tables[j].Name()
	})
	assert.Equal(t, "table1", tables[0].Name())
	assert.Equal(t, "table2", tables[1].Name())
	assert.Equal(t, "table3", tables[2].Name())

	err = tStore.Shutdown()
	assert.NoError(t, err)

	// Delete a table

	tStore, err = LevelDB.Start(logger, dbPath, "table1", "table3")
	assert.NoError(t, err)
	tables = tStore.GetTables()
	assert.Equal(t, 2, len(tables))
	sort.SliceStable(tables, func(i, j int) bool {
		return tables[i].Name() < tables[j].Name()
	})
	assert.Equal(t, "table1", tables[0].Name())

	err = tStore.Shutdown()
	assert.NoError(t, err)

	// Add a table back in (this uses a different code path)
	tStore, err = LevelDB.Start(logger, dbPath, "table1", "table3", "table4")
	assert.NoError(t, err)
	tables = tStore.GetTables()
	assert.Equal(t, 3, len(tables))
	sort.SliceStable(tables, func(i, j int) bool {
		return tables[i].Name() < tables[j].Name()
	})
	assert.Equal(t, "table1", tables[0].Name())
	assert.Equal(t, "table3", tables[1].Name())
	assert.Equal(t, "table4", tables[2].Name())

	err = tStore.Shutdown()
	assert.NoError(t, err)

	// Delete the rest of the tables
	tStore, err = LevelDB.Start(logger, dbPath)
	assert.NoError(t, err)

	tables = tStore.GetTables()
	assert.Equal(t, 0, len(tables))

	err = tStore.Destroy()
	assert.NoError(t, err)
	verifyDBIsDeleted(t)
}

func TestUniqueKeySpace(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	store, err := MapStore.Start(logger, dbPath, "table1", "table2")
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

	store, err := MapStore.Start(logger, dbPath, "table1", "table2", "table3")
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

	store, err := LevelDB.Start(logger, dbPath, "table1", "table2", "table3")
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

	// In order to drop a table, we will need to close the store and reopen it.
	err = store.Shutdown()
	assert.NoError(t, err)

	store, err = LevelDB.Start(logger, dbPath, "table1", "table3")
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

	store, err = LevelDB.Start(logger, dbPath, "table3")
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

	store, err = LevelDB.Start(logger, dbPath)
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

func TestSimultaneousAddAndDrop(t *testing.T) {
	deleteDBDirectory(t)

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	store, err := LevelDB.Start(logger, dbPath, "table1", "table2", "table3", "table4", "table5")
	assert.NoError(t, err)

	table1, err := store.GetTable("table1")
	assert.NoError(t, err)

	table2, err := store.GetTable("table2")
	assert.NoError(t, err)

	table3, err := store.GetTable("table3")
	assert.NoError(t, err)

	table4, err := store.GetTable("table4")
	assert.NoError(t, err)

	table5, err := store.GetTable("table5")
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

		err = table4.Put(k, value)
		assert.NoError(t, err)

		err = table5.Put(k, value)
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

		value, err = table3.Get(k)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)

		value, err = table4.Get(k)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)

		value, err = table5.Get(k)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)
	}

	// In order to drop a table, we will need to close the store and reopen it.
	err = store.Shutdown()
	assert.NoError(t, err)

	store, err = LevelDB.Start(logger, dbPath, "table1", "table5", "table6", "table7", "table8", "table9")
	assert.NoError(t, err)

	table1, err = store.GetTable("table1")
	assert.NoError(t, err)

	_, err = store.GetTable("table2")
	assert.Equal(t, kvstore.ErrTableNotFound, err)

	_, err = store.GetTable("table3")
	assert.Equal(t, kvstore.ErrTableNotFound, err)

	_, err = store.GetTable("table4")
	assert.Equal(t, kvstore.ErrTableNotFound, err)

	table5, err = store.GetTable("table5")
	assert.NoError(t, err)

	table6, err := store.GetTable("table6")
	assert.NoError(t, err)

	table7, err := store.GetTable("table7")
	assert.NoError(t, err)

	table8, err := store.GetTable("table8")
	assert.NoError(t, err)

	table9, err := store.GetTable("table9")
	assert.NoError(t, err)

	// Check data in the tables that were not dropped
	for i := 0; i < 100; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		expectedValue := make([]byte, 8)
		binary.BigEndian.PutUint64(expectedValue, uint64(i))

		value, err := table1.Get(k)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)

		value, err = table5.Get(k)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)
	}

	// Verify the table IDs. There should be no gaps, and the tables that did not get dropped should have the same IDs
	// as before.
	tableView1 := table1.(*tableView)
	assert.Equal(t, uint32(0), tableView1.prefix)

	tableView5 := table5.(*tableView)
	assert.Equal(t, uint32(4), tableView5.prefix)

	tableView6 := table6.(*tableView)
	assert.Equal(t, uint32(1), tableView6.prefix)

	tableView7 := table7.(*tableView)
	assert.Equal(t, uint32(2), tableView7.prefix)

	tableView8 := table8.(*tableView)
	assert.Equal(t, uint32(3), tableView8.prefix)

	tableView9 := table9.(*tableView)
	assert.Equal(t, uint32(5), tableView9.prefix)

	err = store.Destroy()
	assert.NoError(t, err)

	verifyDBIsDeleted(t)
}

func TestIteration(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	store, err := MapStore.Start(logger, dbPath, "table1", "table2")
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

	store, err := LevelDB.Start(logger, dbPath, "table1", "table2")
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

	store, err = LevelDB.Start(logger, dbPath, "table1", "table2")
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

// getTableNameList returns a list of table names from a table map.
func getTableNameList(tableMap map[string]kvstore.Table) []string {
	names := make([]string, 0, len(tableMap))
	for name := range tableMap {
		names = append(names, name)
	}
	return names
}

func TestRandomOperations(t *testing.T) {
	tu.InitializeRandom()

	deleteDBDirectory(t)

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	store, err := LevelDB.Start(logger, dbPath)
	assert.NoError(t, err)

	tables := make(map[string]kvstore.Table)
	expectedData := make(expectedStoreData)

	for i := 0; i < 10000; i++ {

		choice := rand.Float64()

		if choice < 0.01 {
			// restart the store
			err = store.Shutdown()
			assert.NoError(t, err)

			store, err = LevelDB.Start(logger, dbPath, getTableNameList(tables)...)
			assert.NoError(t, err)

			for tableName := range tables {
				table, err := store.GetTable(tableName)
				assert.NoError(t, err)
				tables[tableName] = table
			}
		} else if len(tables) == 0 || choice < 0.02 {
			// Create a new table. Requires the store to be restarted.

			err = store.Shutdown()
			assert.NoError(t, err)

			tableNames := getTableNameList(tables)
			name := tu.RandomString(8)
			tableNames = append(tableNames, name)

			store, err = LevelDB.Start(logger, dbPath, tableNames...)
			assert.NoError(t, err)

			expectedData[name] = make(expectedTableData)
			assert.NoError(t, err)
			tables[name] = nil

			for tableName := range tables {
				table, err := store.GetTable(tableName)
				assert.NoError(t, err)
				tables[tableName] = table
			}
		} else if choice < 0.025 {
			// Drop a table. Requires the store to be restarted.

			err = store.Shutdown()
			assert.NoError(t, err)

			var name string
			for n := range tables {
				name = n
				break
			}
			delete(tables, name)

			store, err = LevelDB.Start(logger, dbPath, getTableNameList(tables)...)
			assert.NoError(t, err)

			// Delete all expected data for the table
			delete(expectedData, name)

			for tableName := range tables {
				table, err := store.GetTable(tableName)
				assert.NoError(t, err)
				tables[tableName] = table
			}
		} else if choice < 0.9 || len(expectedData) == 0 {
			// Write a value

			var tableName string
			for n := range tables {
				tableName = n
				break
			}
			table := tables[tableName]

			k := []byte(tu.RandomString(32))
			v := tu.RandomBytes(32)

			expectedData[tableName][string(k)] = v

			err = table.Put(k, v)
			assert.NoError(t, err)
		} else {
			// Delete a value

			var tableName string
			for n := range tables {
				tableName = n
				break
			}
			table := tables[tableName]

			if len(expectedData[tableName]) == 0 {
				// no data in this table, skip
				continue
			}

			var k string
			for k = range expectedData[tableName] {
				break
			}

			delete(expectedData[tableName], k)
			err = table.Delete([]byte(k))
			assert.NoError(t, err)
		}

		// Every once in a while, check that the store matches the expected data
		if i%100 == 0 {
			// Every so often, check that the store matches the expected data.
			for tableName, tableData := range expectedData {
				table := tables[tableName]

				for k := range tableData {
					expectedValue := tableData[k]
					value, err := table.Get([]byte(k))
					assert.NoError(t, err)
					assert.Equal(t, expectedValue, value)
				}
			}
		}
	}

	err = store.Destroy()
	assert.NoError(t, err)
}

var _ kvstore.Store = &explodingStore{}

// explodingStore is a store that returns an error after a certain number of operations.
// Used to intentionally crash table deletion to exercise table deletion recovery.
type explodingStore struct {
	base               kvstore.Store
	deletionsRemaining int
}

func (e *explodingStore) Put(key []byte, value []byte) error {
	return e.base.Put(key, value)
}

func (e *explodingStore) Get(key []byte) ([]byte, error) {
	return e.base.Get(key)
}

func (e *explodingStore) Delete(key []byte) error {
	if e.deletionsRemaining == 0 {
		return fmt.Errorf("intentional error")
	}
	e.deletionsRemaining--
	return e.base.Delete(key)
}

func (e *explodingStore) NewBatch() kvstore.StoreBatch {
	panic("not used")
}

func (e *explodingStore) NewIterator(prefix []byte) (iterator.Iterator, error) {
	return e.base.NewIterator(prefix)
}

func (e *explodingStore) Shutdown() error {
	return e.base.Shutdown()
}

func (e *explodingStore) Destroy() error {
	return e.base.Destroy()
}

func TestInterruptedTableDeletion(t *testing.T) {
	deleteDBDirectory(t)

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	store, err := LevelDB.Start(logger, dbPath, "table1", "table2")
	assert.NoError(t, err)

	table1, err := store.GetTable("table1")
	assert.NoError(t, err)

	table2, err := store.GetTable("table2")
	assert.NoError(t, err)

	// Write some data to the tables
	for i := 0; i < 100; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		value := make([]byte, 8)
		binary.BigEndian.PutUint64(value, uint64(i))

		err = table1.Put(k, value)
		assert.NoError(t, err)

		err = table2.Put(k, value)
		assert.NoError(t, err)
	}

	// Drop one of the tables (requires restart). Use a store that causes the drop operation to fail partway through.
	err = store.Shutdown()
	assert.NoError(t, err)

	base, err := leveldb.NewStore(logger, dbPath)
	assert.NoError(t, err)

	explodingBase := &explodingStore{
		base:               base,
		deletionsRemaining: 50,
	}

	_, err = start(logger, explodingBase, true, "table2")
	assert.Error(t, err)

	err = explodingBase.Shutdown()
	assert.NoError(t, err)

	// Restart the store. The table should be gone by the time the method returns.
	store, err = LevelDB.Start(logger, dbPath, "table2")
	assert.NoError(t, err)

	tables := store.GetTables()
	assert.Equal(t, 1, len(tables))
	assert.Equal(t, "table2", tables[0].Name())
	table2, err = store.GetTable("table2")
	assert.NoError(t, err)

	// Check that the data in the remaining table is still there. We shouldn't see any data from the deleted table.
	for i := 0; i < 100; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))
		value, err := table2.Get(k)
		assert.NoError(t, err)
		assert.Equal(t, uint64(i), binary.BigEndian.Uint64(value))
	}

	err = store.Destroy()
	assert.NoError(t, err)

	verifyDBIsDeleted(t)
}

func TestLoadWithoutModifiedSchema(t *testing.T) {
	deleteDBDirectory(t)

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	store, err := LevelDB.Start(logger, dbPath, "table1", "table2")
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

	store, err = LevelDB.Load(logger, dbPath)
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
