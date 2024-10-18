package tablestore

import (
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"time"
)

var _ kvstore.TTLBatch = &tableViewBatch{}

// tableViewBatch is a batch for a table in a TableStore.
type tableViewBatch struct {
	table kvstore.Table
	batch kvstore.TTLBatch
}

// PutWithTTL schedules a key-value pair to be added to the table with a time-to-live (TTL).
func (t *tableViewBatch) PutWithTTL(key []byte, value []byte, ttl time.Duration) {
	tableKey := t.table.TableKey(key)
	t.batch.PutWithTTL(tableKey, value, ttl)
}

// PutWithExpiration schedules a key-value pair to be added to the table with an expiration time.
func (t *tableViewBatch) PutWithExpiration(key []byte, value []byte, expiryTime time.Time) {
	tableKey := t.table.TableKey(key)
	t.batch.PutWithExpiration(tableKey, value, expiryTime)
}

// Put schedules a key-value pair to be added to the table.
func (t *tableViewBatch) Put(key []byte, value []byte) {
	k := t.table.TableKey(key)
	t.batch.Put(k, value)
}

// Delete schedules a key-value pair to be removed from the table.
func (t *tableViewBatch) Delete(key []byte) {
	k := t.table.TableKey(key)
	t.batch.Delete(k)
}

// Apply applies the batch to the table.
func (t *tableViewBatch) Apply() error {
	return t.batch.Apply()
}

// Size returns the number of operations in the batch.
func (t *tableViewBatch) Size() uint32 {
	return t.batch.Size()
}
