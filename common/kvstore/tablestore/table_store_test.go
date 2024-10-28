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

	config := DefaultLevelDBConfig(dbPath)
	tStore, err := Start(logger, config)
	assert.NoError(t, err)

	tables := tStore.GetTables()
	assert.Equal(t, 0, len(tables))

	err = tStore.Shutdown()
	assert.NoError(t, err)

	// Add some tables

	config = DefaultLevelDBConfig(dbPath)
	config.Schema = []string{"table1"}
	tStore, err = Start(logger, config)
	assert.NoError(t, err)

	tables = tStore.GetTables()
	assert.Equal(t, 1, len(tables))
	assert.Equal(t, "table1", tables[0])

	err = tStore.Shutdown()
	assert.NoError(t, err)

	config = DefaultLevelDBConfig(dbPath)
	config.Schema = []string{"table1", "table2"}
	tStore, err = Start(logger, config)
	assert.NoError(t, err)
	tables = tStore.GetTables()
	assert.Equal(t, 2, len(tables))
	sort.SliceStable(tables, func(i, j int) bool {
		return tables[i] < tables[j]
	})
	assert.Equal(t, "table1", tables[0])
	assert.Equal(t, "table2", tables[1])

	err = tStore.Shutdown()
	assert.NoError(t, err)

	config = DefaultLevelDBConfig(dbPath)
	config.Schema = []string{"table1", "table2", "table3"}
	tStore, err = Start(logger, config)
	assert.NoError(t, err)
	tables = tStore.GetTables()
	assert.Equal(t, 3, len(tables))
	sort.SliceStable(tables, func(i, j int) bool {
		return tables[i] < tables[j]
	})
	assert.Equal(t, "table1", tables[0])
	assert.Equal(t, "table2", tables[1])
	assert.Equal(t, "table3", tables[2])

	err = tStore.Shutdown()
	assert.NoError(t, err)

	// Restarting with the same tables should work.
	config = DefaultLevelDBConfig(dbPath)
	config.Schema = []string{"table1", "table2", "table3"}
	tStore, err = Start(logger, config)
	assert.NoError(t, err)
	tables = tStore.GetTables()
	assert.Equal(t, 3, len(tables))
	sort.SliceStable(tables, func(i, j int) bool {
		return tables[i] < tables[j]
	})
	assert.Equal(t, "table1", tables[0])
	assert.Equal(t, "table2", tables[1])
	assert.Equal(t, "table3", tables[2])

	err = tStore.Shutdown()
	assert.NoError(t, err)

	// Delete a table
	config = DefaultLevelDBConfig(dbPath)
	config.Schema = []string{"table1", "table3"}
	tStore, err = Start(logger, config)
	assert.NoError(t, err)
	tables = tStore.GetTables()
	assert.Equal(t, 2, len(tables))
	sort.SliceStable(tables, func(i, j int) bool {
		return tables[i] < tables[j]
	})
	assert.Equal(t, "table1", tables[0])

	err = tStore.Shutdown()
	assert.NoError(t, err)

	// Add a table back in (this uses a different code path)
	config = DefaultLevelDBConfig(dbPath)
	config.Schema = []string{"table1", "table3", "table4"}
	tStore, err = Start(logger, config)
	assert.NoError(t, err)
	tables = tStore.GetTables()
	assert.Equal(t, 3, len(tables))
	sort.SliceStable(tables, func(i, j int) bool {
		return tables[i] < tables[j]
	})
	assert.Equal(t, "table1", tables[0])
	assert.Equal(t, "table3", tables[1])
	assert.Equal(t, "table4", tables[2])

	err = tStore.Shutdown()
	assert.NoError(t, err)

	// Delete the rest of the tables
	config = DefaultLevelDBConfig(dbPath)
	config.Schema = []string{}
	tStore, err = Start(logger, config)
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

	config := DefaultMapStoreConfig()
	config.Schema = []string{"table1", "table2"}
	store, err := Start(logger, config)
	assert.NoError(t, err)

	kb1, err := store.GetKeyBuilder("table1")
	assert.NoError(t, err)

	kb2, err := store.GetKeyBuilder("table2")
	assert.NoError(t, err)

	// Write to the tables

	err = store.Put(kb1.Key([]byte("key1")), []byte("value1"))
	assert.NoError(t, err)
	err = store.Put(kb2.Key([]byte("key1")), []byte("value2"))
	assert.NoError(t, err)

	value, err := store.Get(kb1.Key([]byte("key1")))
	assert.NoError(t, err)
	assert.Equal(t, []byte("value1"), value)

	value, err = store.Get(kb2.Key([]byte("key1")))
	assert.NoError(t, err)
	assert.Equal(t, []byte("value2"), value)

	// Delete a key from one table but not the other

	err = store.Delete(kb1.Key([]byte("key1")))
	assert.NoError(t, err)

	_, err = store.Get(kb1.Key([]byte("key1")))
	assert.Equal(t, kvstore.ErrNotFound, err)

	value, err = store.Get(kb2.Key([]byte("key1")))
	assert.NoError(t, err)
	assert.Equal(t, []byte("value2"), value)

	err = store.Destroy()
	assert.NoError(t, err)
}

func TestBatchOperations(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	config := DefaultMapStoreConfig()
	config.Schema = []string{"table1", "table2", "table3"}
	store, err := Start(logger, config)
	assert.NoError(t, err)

	kb1, err := store.GetKeyBuilder("table1")
	assert.NoError(t, err)

	kb2, err := store.GetKeyBuilder("table2")
	assert.NoError(t, err)

	kb3, err := store.GetKeyBuilder("table3")
	assert.NoError(t, err)

	// Test a batch with just puts

	batch := store.NewBatch()
	for i := 0; i < 10; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		v := make([]byte, 8)
		binary.BigEndian.PutUint64(v, uint64(i))
		batch.Put(kb1.Key(k), v)

		v = make([]byte, 8)
		binary.BigEndian.PutUint64(v, uint64(i+10))
		batch.Put(kb2.Key(k), v)

		v = make([]byte, 8)
		binary.BigEndian.PutUint64(v, uint64(i+20))
		batch.Put(kb3.Key(k), v)
	}

	err = batch.Apply()
	assert.NoError(t, err)

	for i := 0; i < 10; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		value, err := store.Get(kb1.Key(k))
		assert.NoError(t, err)
		assert.Equal(t, uint64(i), binary.BigEndian.Uint64(value))

		value, err = store.Get(kb2.Key(k))
		assert.NoError(t, err)
		assert.Equal(t, uint64(i+10), binary.BigEndian.Uint64(value))

		value, err = store.Get(kb3.Key(k))
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

			batch.Delete(kb1.Key(k))
			batch.Delete(kb2.Key(k))
			batch.Delete(kb3.Key(k))
		}
	}

	err = batch.Apply()
	assert.NoError(t, err)

	for i := 0; i < 10; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		if i%2 == 1 {
			_, err = store.Get(kb1.Key(k))
			assert.Equal(t, kvstore.ErrNotFound, err)

			_, err = store.Get(kb2.Key(k))
			assert.Equal(t, kvstore.ErrNotFound, err)

			_, err = store.Get(kb3.Key(k))
			assert.Equal(t, kvstore.ErrNotFound, err)
		} else {
			value, err := store.Get(kb1.Key(k))
			assert.NoError(t, err)
			assert.Equal(t, uint64(i), binary.BigEndian.Uint64(value))

			value, err = store.Get(kb2.Key(k))
			assert.NoError(t, err)
			assert.Equal(t, uint64(i+10), binary.BigEndian.Uint64(value))

			value, err = store.Get(kb3.Key(k))
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
			batch.Delete(kb1.Key(k))
			batch.Delete(kb2.Key(k))
			batch.Delete(kb3.Key(k))
		} else if i%2 == 1 {
			v := make([]byte, 8)
			binary.BigEndian.PutUint64(v, uint64(2*i))
			batch.Put(kb1.Key(k), v)

			v = make([]byte, 8)
			binary.BigEndian.PutUint64(v, uint64(2*i+10))
			batch.Put(kb2.Key(k), v)

			v = make([]byte, 8)
			binary.BigEndian.PutUint64(v, uint64(2*i+20))
			batch.Put(kb3.Key(k), v)
		}
	}

	err = batch.Apply()
	assert.NoError(t, err)

	for i := 0; i < 10; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		if i%4 == 0 {
			_, err = store.Get(kb1.Key(k))
			assert.Equal(t, kvstore.ErrNotFound, err)

			_, err = store.Get(kb2.Key(k))
			assert.Equal(t, kvstore.ErrNotFound, err)

			_, err = store.Get(kb3.Key(k))
			assert.Equal(t, kvstore.ErrNotFound, err)
		} else if i%2 == 1 {
			val, err := store.Get(kb1.Key(k))
			assert.NoError(t, err)
			assert.Equal(t, uint64(2*i), binary.BigEndian.Uint64(val))

			val, err = store.Get(kb2.Key(k))
			assert.NoError(t, err)
			assert.Equal(t, uint64(2*i+10), binary.BigEndian.Uint64(val))

			val, err = store.Get(kb3.Key(k))
			assert.NoError(t, err)
			assert.Equal(t, uint64(2*i+20), binary.BigEndian.Uint64(val))
		} else {
			val, err := store.Get(kb1.Key(k))
			assert.NoError(t, err)
			assert.Equal(t, uint64(i), binary.BigEndian.Uint64(val))

			val, err = store.Get(kb2.Key(k))
			assert.NoError(t, err)
			assert.Equal(t, uint64(i+10), binary.BigEndian.Uint64(val))

			val, err = store.Get(kb3.Key(k))
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

	config := DefaultLevelDBConfig(dbPath)
	config.Schema = []string{"table1", "table2", "table3"}
	store, err := Start(logger, config)
	assert.NoError(t, err)

	kb1, err := store.GetKeyBuilder("table1")
	assert.NoError(t, err)

	kb2, err := store.GetKeyBuilder("table2")
	assert.NoError(t, err)

	kb3, err := store.GetKeyBuilder("table3")
	assert.NoError(t, err)

	// Insert some data into the tables
	for i := 0; i < 100; i++ {
		value := make([]byte, 8)
		binary.BigEndian.PutUint64(value, uint64(i))

		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		err = store.Put(kb1.Key(k), value)
		assert.NoError(t, err)

		err = store.Put(kb2.Key(k), value)
		assert.NoError(t, err)

		err = store.Put(kb3.Key(k), value)
		assert.NoError(t, err)
	}

	// Verify the data is there
	for i := 0; i < 100; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		expectedValue := make([]byte, 8)
		binary.BigEndian.PutUint64(expectedValue, uint64(i))

		value, err := store.Get(kb1.Key(k))
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)

		value, err = store.Get(kb2.Key(k))
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)

		_, err = store.Get(kb3.Key(k))
		assert.NoError(t, err)
	}

	// In order to drop a table, we will need to close the store and reopen it.
	err = store.Shutdown()
	assert.NoError(t, err)

	config = DefaultLevelDBConfig(dbPath)
	config.Schema = []string{"table1", "table3"}
	store, err = Start(logger, config)
	assert.NoError(t, err)

	kb1, err = store.GetKeyBuilder("table1")
	assert.NoError(t, err)

	_, err = store.GetKeyBuilder("table2")
	assert.Equal(t, kvstore.ErrTableNotFound, err)

	kb3, err = store.GetKeyBuilder("table3")
	assert.NoError(t, err)

	for i := 0; i < 100; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		expectedValue := make([]byte, 8)
		binary.BigEndian.PutUint64(expectedValue, uint64(i))

		value, err := store.Get(kb1.Key(k))
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)

		value, err = store.Get(kb3.Key(k))
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)
	}

	// Restart the store so that we can drop another table.
	err = store.Shutdown()
	assert.NoError(t, err)

	config = DefaultLevelDBConfig(dbPath)
	config.Schema = []string{"table3"}
	store, err = Start(logger, config)
	assert.NoError(t, err)

	_, err = store.GetKeyBuilder("table1")
	assert.Equal(t, kvstore.ErrTableNotFound, err)

	_, err = store.GetKeyBuilder("table2")
	assert.Equal(t, kvstore.ErrTableNotFound, err)

	kb3, err = store.GetKeyBuilder("table3")
	assert.NoError(t, err)

	for i := 0; i < 100; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		expectedValue := make([]byte, 8)
		binary.BigEndian.PutUint64(expectedValue, uint64(i))

		value, err := store.Get(kb3.Key(k))
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)
	}

	// Restart the store so that we can drop the last table.
	err = store.Shutdown()
	assert.NoError(t, err)

	config = DefaultLevelDBConfig(dbPath)
	config.Schema = []string{}
	store, err = Start(logger, config)
	assert.NoError(t, err)

	_, err = store.GetKeyBuilder("table1")
	assert.Equal(t, kvstore.ErrTableNotFound, err)

	_, err = store.GetKeyBuilder("table2")
	assert.Equal(t, kvstore.ErrTableNotFound, err)

	_, err = store.GetKeyBuilder("table3")
	assert.Equal(t, kvstore.ErrTableNotFound, err)

	err = store.Destroy()
	assert.NoError(t, err)

	verifyDBIsDeleted(t)
}

func TestSimultaneousAddAndDrop(t *testing.T) {
	deleteDBDirectory(t)

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	config := DefaultLevelDBConfig(dbPath)
	config.Schema = []string{"table1", "table2", "table3", "table4", "table5"}
	store, err := Start(logger, config)
	assert.NoError(t, err)

	kb1, err := store.GetKeyBuilder("table1")
	assert.NoError(t, err)

	kb2, err := store.GetKeyBuilder("table2")
	assert.NoError(t, err)

	kb3, err := store.GetKeyBuilder("table3")
	assert.NoError(t, err)

	kb4, err := store.GetKeyBuilder("table4")
	assert.NoError(t, err)

	kb5, err := store.GetKeyBuilder("table5")
	assert.NoError(t, err)

	// Insert some data into the tables
	for i := 0; i < 100; i++ {
		value := make([]byte, 8)
		binary.BigEndian.PutUint64(value, uint64(i))

		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		err = store.Put(kb1.Key(k), value)
		assert.NoError(t, err)

		err = store.Put(kb2.Key(k), value)
		assert.NoError(t, err)

		err = store.Put(kb3.Key(k), value)
		assert.NoError(t, err)

		err = store.Put(kb4.Key(k), value)
		assert.NoError(t, err)

		err = store.Put(kb5.Key(k), value)
		assert.NoError(t, err)
	}

	// Verify the data is there
	for i := 0; i < 100; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		expectedValue := make([]byte, 8)
		binary.BigEndian.PutUint64(expectedValue, uint64(i))

		value, err := store.Get(kb1.Key(k))
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)

		value, err = store.Get(kb2.Key(k))
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)

		value, err = store.Get(kb3.Key(k))
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)

		value, err = store.Get(kb4.Key(k))
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)

		value, err = store.Get(kb5.Key(k))
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)
	}

	// In order to drop a table, we will need to close the store and reopen it.
	err = store.Shutdown()
	assert.NoError(t, err)

	config = DefaultLevelDBConfig(dbPath)
	config.Schema = []string{"table1", "table5", "table6", "table7", "table8", "table9"}
	store, err = Start(logger, config)
	assert.NoError(t, err)

	kb1, err = store.GetKeyBuilder("table1")
	assert.NoError(t, err)

	_, err = store.GetKeyBuilder("table2")
	assert.Equal(t, kvstore.ErrTableNotFound, err)

	_, err = store.GetKeyBuilder("table3")
	assert.Equal(t, kvstore.ErrTableNotFound, err)

	_, err = store.GetKeyBuilder("table4")
	assert.Equal(t, kvstore.ErrTableNotFound, err)

	kb2, err = store.GetKeyBuilder("table5")
	assert.NoError(t, err)

	table6, err := store.GetKeyBuilder("table6")
	assert.NoError(t, err)

	table7, err := store.GetKeyBuilder("table7")
	assert.NoError(t, err)

	table8, err := store.GetKeyBuilder("table8")
	assert.NoError(t, err)

	table9, err := store.GetKeyBuilder("table9")
	assert.NoError(t, err)

	// Check data in the tables that were not dropped
	for i := 0; i < 100; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		expectedValue := make([]byte, 8)
		binary.BigEndian.PutUint64(expectedValue, uint64(i))

		value, err := store.Get(kb1.Key(k))
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)

		value, err = store.Get(kb5.Key(k))
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)
	}

	// Verify the table IDs.
	assert.Equal(t, uint32(0), getTableID(kb1))
	assert.Equal(t, uint32(4), getTableID(kb2))
	assert.Equal(t, uint32(5), getTableID(table6))
	assert.Equal(t, uint32(6), getTableID(table7))
	assert.Equal(t, uint32(7), getTableID(table8))
	assert.Equal(t, uint32(8), getTableID(table9))

	err = store.Destroy()
	assert.NoError(t, err)

	verifyDBIsDeleted(t)
}

func getTableID(kb kvstore.KeyBuilder) uint32 {
	prefix := kb.(*keyBuilder).prefix
	return binary.BigEndian.Uint32(prefix)
}

func TestIteration(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	config := DefaultMapStoreConfig()
	config.Schema = []string{"table1", "table2"}
	store, err := Start(logger, config)
	assert.NoError(t, err)

	kb1, err := store.GetKeyBuilder("table1")
	assert.NoError(t, err)

	kb2, err := store.GetKeyBuilder("table2")
	assert.NoError(t, err)

	// Prefix "qwer"
	for i := 0; i < 100; i++ {
		k := []byte(fmt.Sprintf("qwer%3d", i))

		value := make([]byte, 8)
		binary.BigEndian.PutUint64(value, uint64(i))
		err = store.Put(kb1.Key(k), value)
		assert.NoError(t, err)

		value = make([]byte, 8)
		binary.BigEndian.PutUint64(value, uint64(2*i))
		err = store.Put(kb2.Key(k), value)
		assert.NoError(t, err)
	}

	// Prefix "asdf"
	for i := 0; i < 100; i++ {
		k := []byte(fmt.Sprintf("asdf%3d", i))

		value := make([]byte, 8)
		binary.BigEndian.PutUint64(value, uint64(i))
		err = store.Put(kb1.Key(k), value)
		assert.NoError(t, err)

		value = make([]byte, 8)
		binary.BigEndian.PutUint64(value, uint64(2*i))
		err = store.Put(kb2.Key(k), value)
		assert.NoError(t, err)
	}

	// Iterate table 1 with no prefix filter
	it, err := store.NewTableIterator(kb1)
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
	it, err = store.NewTableIterator(kb2)
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
	it, err = store.NewIterator(kb1.Key([]byte("qwer")))
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
	it, err = store.NewIterator(kb2.Key([]byte("asdf")))
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

	config := DefaultLevelDBConfig(dbPath)
	config.Schema = []string{"table1", "table2"}
	store, err := Start(logger, config)
	assert.NoError(t, err)

	kb1, err := store.GetKeyBuilder("table1")
	assert.NoError(t, err)

	kb2, err := store.GetKeyBuilder("table2")
	assert.NoError(t, err)

	for i := 0; i < 100; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		value1 := make([]byte, 8)
		binary.BigEndian.PutUint64(value1, uint64(i))

		value2 := make([]byte, 8)
		binary.BigEndian.PutUint64(value2, uint64(i*2))

		err = store.Put(kb1.Key(k), value1)
		assert.NoError(t, err)

		err = store.Put(kb2.Key(k), value2)
		assert.NoError(t, err)
	}

	err = store.Shutdown()
	assert.NoError(t, err)

	// Restart the store
	store, err = Start(logger, config)
	assert.NoError(t, err)

	kb1, err = store.GetKeyBuilder("table1")
	assert.NoError(t, err)
	kb2, err = store.GetKeyBuilder("table2")
	assert.NoError(t, err)

	for i := 0; i < 100; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		expectedValue1 := make([]byte, 8)
		binary.BigEndian.PutUint64(expectedValue1, uint64(i))

		expectedValue2 := make([]byte, 8)
		binary.BigEndian.PutUint64(expectedValue2, uint64(i*2))

		value1, err := store.Get(kb1.Key(k))
		assert.NoError(t, err)
		assert.Equal(t, expectedValue1, value1)

		value2, err := store.Get(kb2.Key(k))
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
func getTableNameList(tableMap map[string]kvstore.KeyBuilder) []string {
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

	config := DefaultLevelDBConfig(dbPath)
	store, err := Start(logger, config)
	assert.NoError(t, err)

	tables := make(map[string]kvstore.KeyBuilder)
	expectedData := make(expectedStoreData)

	for i := 0; i < 10000; i++ {

		choice := rand.Float64()

		if choice < 0.01 {
			// restart the store
			err = store.Shutdown()
			assert.NoError(t, err)

			config = DefaultLevelDBConfig(dbPath)
			config.Schema = getTableNameList(tables)
			store, err = Start(logger, config)
			assert.NoError(t, err)

			for tableName := range tables {
				table, err := store.GetKeyBuilder(tableName)
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

			config = DefaultLevelDBConfig(dbPath)
			config.Schema = tableNames
			store, err = Start(logger, config)
			assert.NoError(t, err)

			expectedData[name] = make(expectedTableData)
			assert.NoError(t, err)
			tables[name] = nil

			for tableName := range tables {
				table, err := store.GetKeyBuilder(tableName)
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

			config = DefaultLevelDBConfig(dbPath)
			config.Schema = getTableNameList(tables)
			store, err = Start(logger, config)
			assert.NoError(t, err)

			// Delete all expected data for the table
			delete(expectedData, name)

			for tableName := range tables {
				table, err := store.GetKeyBuilder(tableName)
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

			err = store.Put(table.Key(k), v)
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
			err = store.Delete(table.Key([]byte(k)))
			assert.NoError(t, err)
		}

		// Every once in a while, check that the store matches the expected data
		if i%100 == 0 {
			// Every so often, check that the store matches the expected data.
			for tableName, tableData := range expectedData {
				table := tables[tableName]

				for k := range tableData {
					expectedValue := tableData[k]
					value, err := store.Get(table.Key([]byte(k)))
					assert.NoError(t, err)
					assert.Equal(t, expectedValue, value)
				}
			}
		}
	}

	err = store.Destroy()
	assert.NoError(t, err)
}

var _ kvstore.Store[[]byte] = &explodingStore{}

// explodingStore is a store that returns an error after a certain number of operations.
// Used to intentionally crash table deletion to exercise table deletion recovery.
type explodingStore struct {
	base               kvstore.Store[[]byte]
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

func (e *explodingStore) NewBatch() kvstore.Batch[[]byte] {
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

	config := DefaultLevelDBConfig(dbPath)
	config.Schema = []string{"table1", "table2"}
	store, err := Start(logger, config)
	assert.NoError(t, err)

	kb1, err := store.GetKeyBuilder("table1")
	assert.NoError(t, err)

	kb2, err := store.GetKeyBuilder("table2")
	assert.NoError(t, err)

	// Write some data to the tables
	for i := 0; i < 100; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		value := make([]byte, 8)
		binary.BigEndian.PutUint64(value, uint64(i))

		err = store.Put(kb1.Key(k), value)
		assert.NoError(t, err)

		err = store.Put(kb2.Key(k), value)
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

	config = DefaultLevelDBConfig(dbPath)
	config.Schema = []string{"table2"}
	_, err = start(logger, explodingBase, config)
	assert.Error(t, err)

	err = explodingBase.Shutdown()
	assert.NoError(t, err)

	// Restart the store. The table should be gone by the time the method returns.
	config = DefaultLevelDBConfig(dbPath)
	config.Schema = []string{"table2"}
	store, err = Start(logger, config)
	assert.NoError(t, err)

	tables := store.GetTables()
	assert.Equal(t, 1, len(tables))
	assert.Equal(t, "table2", tables[0])
	kb2, err = store.GetKeyBuilder("table2")
	assert.NoError(t, err)

	// Check that the data in the remaining table is still there. We shouldn't see any data from the deleted table.
	for i := 0; i < 100; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))
		value, err := store.Get(kb2.Key(k))
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

	config := DefaultLevelDBConfig(dbPath)
	config.Schema = []string{"table1", "table2"}
	store, err := Start(logger, config)
	assert.NoError(t, err)

	kb1, err := store.GetKeyBuilder("table1")
	assert.NoError(t, err)

	kb2, err := store.GetKeyBuilder("table2")
	assert.NoError(t, err)

	for i := 0; i < 100; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		value1 := make([]byte, 8)
		binary.BigEndian.PutUint64(value1, uint64(i))

		value2 := make([]byte, 8)
		binary.BigEndian.PutUint64(value2, uint64(i*2))

		err = store.Put(kb1.Key(k), value1)
		assert.NoError(t, err)

		err = store.Put(kb2.Key(k), value2)
		assert.NoError(t, err)
	}

	err = store.Shutdown()
	assert.NoError(t, err)

	// Load the store without the schema
	config = DefaultLevelDBConfig(dbPath)
	store, err = Start(logger, config)
	assert.NoError(t, err)

	kb1, err = store.GetKeyBuilder("table1")
	assert.NoError(t, err)
	kb2, err = store.GetKeyBuilder("table2")
	assert.NoError(t, err)

	for i := 0; i < 100; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		expectedValue1 := make([]byte, 8)
		binary.BigEndian.PutUint64(expectedValue1, uint64(i))

		expectedValue2 := make([]byte, 8)
		binary.BigEndian.PutUint64(expectedValue2, uint64(i*2))

		value1, err := store.Get(kb1.Key(k))
		assert.NoError(t, err)
		assert.Equal(t, expectedValue1, value1)

		value2, err := store.Get(kb2.Key(k))
		assert.NoError(t, err)
		assert.Equal(t, expectedValue2, value2)
	}

	err = store.Destroy()
	assert.NoError(t, err)

	verifyDBIsDeleted(t)
}
