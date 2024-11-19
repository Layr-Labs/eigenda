package tablestore

import (
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"time"
)

var _ kvstore.TTLBatch[kvstore.Key] = &tableStoreBatch{}

// tableStoreBatch is a batch for writing to a table store.
type tableStoreBatch struct {
	baseBatch            kvstore.Batch[[]byte]
	expirationKeyBuilder kvstore.KeyBuilder
}

// newTableStoreBatch creates a new batch for writing to a table store.
func newTableStoreBatch(
	base kvstore.Store[[]byte],
	expirationKeyBuilder kvstore.KeyBuilder) kvstore.TTLBatch[kvstore.Key] {

	return &tableStoreBatch{
		baseBatch:            base.NewBatch(),
		expirationKeyBuilder: expirationKeyBuilder,
	}
}

// PutWithTTL adds a key-value pair to the batch that expires after a specified duration.
func (t *tableStoreBatch) PutWithTTL(k kvstore.Key, value []byte, ttl time.Duration) {
	expirationTime := time.Now().Add(ttl)
	t.PutWithExpiration(k, value, expirationTime)
}

// PutWithExpiration adds a key-value pair to the batch that expires at a specified time.
func (t *tableStoreBatch) PutWithExpiration(k kvstore.Key, value []byte, expiryTime time.Time) {
	expirationKey := t.expirationKeyBuilder.Key(prependTimestamp(expiryTime, k.Raw()))

	t.baseBatch.Put(k.Raw(), value)
	t.baseBatch.Put(expirationKey.Raw(), []byte{})
}

// Put adds a key-value pair to the batch.
func (t *tableStoreBatch) Put(k kvstore.Key, value []byte) {
	t.baseBatch.Put(k.Raw(), value)
}

// Delete removes a key-value pair from the batch.
func (t *tableStoreBatch) Delete(k kvstore.Key) {
	t.baseBatch.Delete(k.Raw())
}

// Apply applies the batch to the store.
func (t *tableStoreBatch) Apply() error {
	return t.baseBatch.Apply()
}

// Size returns the number of operations in the batch.
func (t *tableStoreBatch) Size() uint32 {
	return t.baseBatch.Size()
}
