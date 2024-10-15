package tablestore

import "github.com/Layr-Labs/eigenda/common/kvstore"

// tableStoreBatch is a batch for writing to a table store.
type tableStoreBatch struct {
	store *tableStore
	batch kvstore.StoreBatch
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
