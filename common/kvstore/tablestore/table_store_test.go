package tablestore

import (
	"encoding/binary"
	"fmt"
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
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	base := mapstore.NewStore()
	store, err := TableStoreWrapper(logger, base)
	assert.NoError(t, err)

	// table count needs to fit into 32 bytes, and two tables are reserved for internal use
	maxTableCount := uint32(math.MaxUint32 - 2)
	assert.Equal(t, maxTableCount, store.GetMaxTableCount())
	assert.Equal(t, uint32(0), store.GetTableCount())

	table1, err := store.GetTable("table1")
	assert.NoError(t, err)
	assert.Equal(t, "table1", table1.Name())
	assert.Equal(t, uint32(1), store.GetTableCount())

	table2, err := store.GetTable("table2")
	assert.NoError(t, err)
	assert.Equal(t, "table2", table2.Name())
	assert.Equal(t, uint32(2), store.GetTableCount())

	table3, err := store.GetTable("table3")
	assert.NoError(t, err)
	assert.Equal(t, "table3", table3.Name())
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
	store, err := TableStoreWrapper(logger, base)
	assert.NoError(t, err)

	tables := store.GetTables()
	assert.Equal(t, 0, len(tables))

	// Add some tables

	_, err = store.GetTable("table1")
	assert.NoError(t, err)

	tables = store.GetTables()
	assert.Equal(t, 1, len(tables))
	assert.Equal(t, "table1", tables[0].Name())

	_, err = store.GetTable("table2")
	assert.NoError(t, err)
	tables = store.GetTables()
	assert.Equal(t, 2, len(tables))
	sort.SliceStable(tables, func(i, j int) bool {
		return tables[i].Name() < tables[j].Name()
	})
	assert.Equal(t, "table1", tables[0].Name())
	assert.Equal(t, "table2", tables[1].Name())

	_, err = store.GetTable("table3")
	assert.NoError(t, err)
	tables = store.GetTables()
	assert.Equal(t, 3, len(tables))
	sort.SliceStable(tables, func(i, j int) bool {
		return tables[i].Name() < tables[j].Name()
	})
	assert.Equal(t, "table1", tables[0].Name())
	assert.Equal(t, "table2", tables[1].Name())
	assert.Equal(t, "table3", tables[2].Name())

	// Duplicate table additions should be no-ops
	_, err = store.GetTable("table1")
	assert.NoError(t, err)
	tables = store.GetTables()
	assert.Equal(t, 3, len(tables))
	sort.SliceStable(tables, func(i, j int) bool {
		return tables[i].Name() < tables[j].Name()
	})
	assert.Equal(t, "table1", tables[0].Name())
	assert.Equal(t, "table2", tables[1].Name())
	assert.Equal(t, "table3", tables[2].Name())

	// Tables should survive a restart
	err = store.Shutdown()
	assert.NoError(t, err)

	base, err = leveldb.NewStore(logger, dbPath)
	assert.NoError(t, err)
	store, err = TableStoreWrapper(logger, base)
	assert.NoError(t, err)

	tables = store.GetTables()
	assert.Equal(t, 3, len(tables))
	sort.SliceStable(tables, func(i, j int) bool {
		return tables[i].Name() < tables[j].Name()
	})
	assert.Equal(t, "table1", tables[0].Name())
	assert.Equal(t, "table2", tables[1].Name())
	assert.Equal(t, "table3", tables[2].Name())

	// Delete a table

	err = store.DropTable("table2")
	assert.NoError(t, err)
	tables = store.GetTables()
	assert.Equal(t, 2, len(tables))
	sort.SliceStable(tables, func(i, j int) bool {
		return tables[i].Name() < tables[j].Name()
	})
	assert.Equal(t, "table1", tables[0].Name())

	// Table should be deleted after a restart
	err = store.Shutdown()
	assert.NoError(t, err)

	base, err = leveldb.NewStore(logger, dbPath)
	assert.NoError(t, err)
	store, err = TableStoreWrapper(logger, base)
	assert.NoError(t, err)

	tables = store.GetTables()
	assert.Equal(t, 2, len(tables))
	sort.SliceStable(tables, func(i, j int) bool {
		return tables[i].Name() < tables[j].Name()
	})
	assert.Equal(t, "table1", tables[0].Name())
	assert.Equal(t, "table3", tables[1].Name())

	// Add a table back in (this uses a different code path)
	_, err = store.GetTable("table4")
	assert.NoError(t, err)
	tables = store.GetTables()
	assert.Equal(t, 3, len(tables))
	sort.SliceStable(tables, func(i, j int) bool {
		return tables[i].Name() < tables[j].Name()
	})
	assert.Equal(t, "table1", tables[0].Name())
	assert.Equal(t, "table3", tables[1].Name())
	assert.Equal(t, "table4", tables[2].Name())

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
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	base := mapstore.NewStore()
	store, err := TableStoreWrapper(logger, base)
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

// TODO test that mixes put and delete
func TestBatchOperations(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	base := mapstore.NewStore()
	store, err := TableStoreWrapper(logger, base)
	assert.NoError(t, err)

	table1, err := store.GetTable("table1")
	assert.NoError(t, err)

	table2, err := store.GetTable("table2")
	assert.NoError(t, err)

	table3, err := store.GetTable("table3")
	assert.NoError(t, err)

	batch := store.NewBatch()

	for i := 0; i < 10; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))
		v := make([]byte, 8)
		binary.BigEndian.PutUint64(v, uint64(i))
		batch.Put(table1.TableKey(k), v)
	}
	for i := 0; i < 10; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))
		v := make([]byte, 8)
		binary.BigEndian.PutUint64(v, uint64(i+10))
		batch.Put(table2.TableKey(k), v)
	}
	for i := 0; i < 10; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))
		v := make([]byte, 8)
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

	err = store.Destroy()
	assert.NoError(t, err)
}

func TestDropTable(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	base := mapstore.NewStore()
	store, err := TableStoreWrapper(logger, base)
	assert.NoError(t, err)

	table1, err := store.GetTable("table1")
	assert.NoError(t, err)

	table2, err := store.GetTable("table2")
	assert.NoError(t, err)

	table3, err := store.GetTable("table3")
	assert.NoError(t, err)

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

	err = store.DropTable("table2")
	assert.NoError(t, err)

	for i := 0; i < 100; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		expectedValue := make([]byte, 8)
		binary.BigEndian.PutUint64(expectedValue, uint64(i))

		value, err := table1.Get(k)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)

		_, err = table2.Get(k)
		assert.Equal(t, kvstore.ErrNotFound, err)

		value, err = table3.Get(k)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)
	}

	err = store.DropTable("table1")
	assert.NoError(t, err)

	for i := 0; i < 100; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		expectedValue := make([]byte, 8)
		binary.BigEndian.PutUint64(expectedValue, uint64(i))

		_, err := table1.Get(k)
		assert.Equal(t, kvstore.ErrNotFound, err)

		_, err = table2.Get(k)
		assert.Equal(t, kvstore.ErrNotFound, err)

		value, err := table3.Get(k)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)
	}

	err = store.DropTable("table3")
	assert.NoError(t, err)

	for i := 0; i < 100; i++ {
		k := make([]byte, 8)
		binary.BigEndian.PutUint64(k, uint64(i))

		_, err := table1.Get(k)
		assert.Equal(t, kvstore.ErrNotFound, err)

		_, err = table2.Get(k)
		assert.Equal(t, kvstore.ErrNotFound, err)

		_, err = table3.Get(k)
		assert.Equal(t, kvstore.ErrNotFound, err)
	}

	err = store.Destroy()
	assert.NoError(t, err)
}

func TestIteration(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	base := mapstore.NewStore()
	store, err := TableStoreWrapper(logger, base)
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

		err = table2.Put(k, value)
		assert.NoError(t, err)
	}

	// TODO put different values in table 1 and 2, iterate over each

	// Iterate with no prefix filter
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

	// Iterate over the "qwer" keys
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
}

func TestRestart(t *testing.T) {
	deleteDBDirectory(t)

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	assert.NoError(t, err)

	base, err := leveldb.NewStore(logger, dbPath)
	assert.NoError(t, err)
	store, err := TableStoreWrapper(logger, base)
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

	base, err = leveldb.NewStore(logger, dbPath)
	assert.NoError(t, err)
	store, err = TableStoreWrapper(logger, base)
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

// TODO
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
//	store, err := TableStoreWrapper(logger, base)
//	assert.NoError(t, err)
//
//	tables := make(map[string]kvstore.Table)
//	expectedData := make(map[kvstore.Key][]byte)
//	keysByTable := make(map[string][]kvstore.Key)
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
//			store, err = TableStoreWrapper(logger, base)
//			assert.NoError(t, err)
//		} else if len(tables) == 0 || choice < 0.1 {
//			// Create a new table.
//			name := tu.RandomString(8)
//			kb, _, err := store.GetTable(name)
//			keysByTable[name] = make([]kvstore.Key, 0)
//			assert.NoError(t, err)
//			tables[name] = kb
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
//			for _, k := range keysByTable[name] {
//				delete(expectedData, k)
//			}
//			delete(keysByTable, name)
//			delete(tables, name)
//		} else if choice < 0.9 || len(expectedData) == 0 {
//			// Write a value
//
//			var tableName string
//			for n := range tables {
//				tableName = n
//				break
//			}
//			kb := tables[tableName]
//
//			k := kb.StringKey(tu.RandomString(32))
//			v := tu.RandomBytes(32)
//
//			keysByTable[tableName] = append(keysByTable[tableName], k)
//
//			err := store.Put(k, v)
//			assert.NoError(t, err)
//
//			expectedData[k] = v
//
//		} else {
//			// Drop a value
//			var k kvstore.Key
//			for k = range expectedData {
//				break
//			}
//			delete(expectedData, k)
//			err := store.Delete(k)
//			assert.NoError(t, err)
//		}
//
//		// Every once in a while, check that the store matches the expected data
//		if i%100 == 0 {
//			// Every so often, check that the store matches the expected data.
//			for k, expectedValue := range expectedData {
//				value, err := store.Get(k)
//				assert.NoError(t, err)
//				assert.Equal(t, expectedValue, value)
//			}
//		}
//	}
//
//	err = store.Destroy()
//	assert.NoError(t, err)
//}
//
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
//func (e *explodingStore) DeleteBatch(keys [][]byte) error {
//	return e.base.DeleteBatch(keys)
//}
//
//func (e *explodingStore) WriteBatch(keys, values [][]byte) error {
//	return e.base.WriteBatch(keys, values)
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
//	store, err := TableStoreWrapper(logger, explodingBase)
//	assert.NoError(t, err)
//
//	// Create a few tables
//	kb1, _, err := store.GetTable("table1")
//	assert.NoError(t, err)
//
//	kb2, _, err := store.GetTable("table2")
//	assert.NoError(t, err)
//
//	// Write some data to the tables
//	for i := 0; i < 100; i++ {
//		value := make([]byte, 8)
//		binary.BigEndian.PutUint64(value, uint64(i))
//
//		err = store.Put(kb1.Uint64Key(uint64(i)), value)
//		assert.NoError(t, err)
//
//		err = store.Put(kb2.Uint64Key(uint64(i)), value)
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
//	store, err = TableStoreWrapper(logger, base)
//	assert.NoError(t, err)
//
//	tables := store.GetTables()
//	assert.Equal(t, 1, len(tables))
//	assert.Equal(t, "table2", tables[0])
//
//	// Check that the data in the remaining table is still there. We shouldn't see any data from the deleted table.
//	for i := 0; i < 100; i++ {
//		_, err := store.Get(kb1.Uint64Key(uint64(i)))
//		assert.Equal(t, kvstore.ErrNotFound, err)
//
//		value, err := store.Get(kb2.Uint64Key(uint64(i)))
//		assert.NoError(t, err)
//		assert.Equal(t, uint64(i), binary.BigEndian.Uint64(value))
//	}
//
//	err = store.Destroy()
//	assert.NoError(t, err)
//
//	verifyDBIsDeleted(t)
//}
