package tablestore

import (
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"time"
)

var _ kvstore.TableStoreBatch = &tableStoreBatch{}

// tableStoreBatch is a batch for writing to a table store.
type tableStoreBatch struct {
	batch           kvstore.StoreBatch
	expirationTable kvstore.Table
}

// PutWithTTL adds a key-value pair to the batch that expires after a specified duration.
func (t *tableStoreBatch) PutWithTTL(key kvstore.TableKey, value []byte, ttl time.Duration) {
	expirationTime := time.Now().Add(ttl)
	t.PutWithExpiration(key, value, expirationTime)
}

// PutWithExpiration adds a key-value pair to the batch that expires at a specified time.
func (t *tableStoreBatch) PutWithExpiration(key kvstore.TableKey, value []byte, expiryTime time.Time) {
	expirationKey := t.expirationTable.TableKey(prependTimestamp(expiryTime, key))

	t.Put(key, value)
	t.Put(expirationKey, make([]byte, 0))
}

// Put adds a key-value pair to the batch.
func (t *tableStoreBatch) Put(key kvstore.TableKey, value []byte) {
	if value == nil {
		value = []byte{}
	}
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
