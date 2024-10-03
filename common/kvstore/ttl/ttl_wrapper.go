package ttl

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
	"time"
)

var _ kvstore.Store = &ttlStore{}

// ttlStore adds a time-to-live (TTL) capability to the store.
//
// This store utilizes the properties of store iteration. Namely, that the keys are returned in lexicographical order,
// as well as the ability to filter keys by prefix. "Regular" keys are stored in the store with a prefix "k", while
// special expiry keys are stored with a prefix "e". The expiry key also contains the expiry time in hexadecimal format,
// such that when iterating over expiry keys in lexicographical order, the keys are ordered by expiry time. The value
// each expiry key points to is the regular key that is to be deleted when the expiry time is reached. In order to
// efficiently delete expired keys, the expiry keys must be iterated over periodically to find and delete expired keys.
type ttlStore struct {
	store  kvstore.Store
	ctx    context.Context
	cancel context.CancelFunc

	logger logging.Logger
}

// TTLWrapper extends the given store with TTL capabilities. Periodically checks for expired keys and deletes them
// with a period of gcPeriod. If gcPeriod is 0, no background goroutine is spawned to check for expired keys.
//
// Note: it is unsafe to access the wrapped store directly while the TTLStore is in use. The TTLStore uses special
// key formatting, and direct access to the wrapped store may violate the TTLStore's invariants, resulting in
// undefined behavior.
func TTLWrapper(
	ctx context.Context,
	logger logging.Logger,
	store kvstore.Store,
	gcPeriod time.Duration) kvstore.TTLStore {

	ctx, cancel := context.WithCancel(ctx)

	ttlStore := &ttlStore{
		store:  store,
		ctx:    ctx,
		cancel: cancel,
		logger: logger,
	}
	if gcPeriod > 0 {
		ttlStore.expireKeysInBackground(gcPeriod)
	}
	return ttlStore
}

var keyPrefix = []byte("k")
var expiryPrefix = []byte("e")
var maxDeletionBatchSize = 1000

// PutWithTTL adds a key-value pair to the store that expires after a specified time-to-live (TTL).
// Key is eventually deleted after the TTL elapses.
func (store *ttlStore) PutWithTTL(key []byte, value []byte, ttl time.Duration) error {
	expiryTime := time.Now().Add(ttl)
	return store.PutWithExpiration(key, value, expiryTime)
}

// PutBatchWithTTL adds multiple key-value pairs to the store that expire after a specified time-to-live (TTL).
func (store *ttlStore) PutBatchWithTTL(keys [][]byte, values [][]byte, ttl time.Duration) error {
	expiryTime := time.Now().Add(ttl)
	return store.PutBatchWithExpiration(keys, values, expiryTime)
}

// buildExpiryKey creates an expiry key from the given expiry time.
// The expiry key is composed of the following 3 components appended one after the other:
// - a one byte "e" prefix
// - the expiry time in hexadecimal format (8 bytes)
// - and the base key.
func buildExpiryKey(
	baseKey []byte,
	expiryTime time.Time) []byte {

	expiryKeyLength := 1 /* prefix */ + 8 /* expiry timestamp */ + len(baseKey)
	expiryKey := make([]byte, expiryKeyLength)

	expiryKey[0] = 'e'
	expiryUnixNano := expiryTime.UnixNano()
	binary.BigEndian.PutUint64(expiryKey[1:], uint64(expiryUnixNano))

	copy(expiryKey[9:], baseKey)

	return expiryKey
}

// parseExpiryKey extracts the expiry time and base key from the given expiry key.
func parseExpiryKey(expiryKey []byte) (baseKey []byte, expiryTime time.Time) {
	expiryUnixNano := int64(binary.BigEndian.Uint64(expiryKey[1:]))
	expiryTime = time.Unix(0, expiryUnixNano)

	baseKey = expiryKey[9:]
	return
}

// PutWithExpiration adds a key-value pair to the store that expires at a specified time.
// Key is eventually deleted after the expiry time.
func (store *ttlStore) PutWithExpiration(key []byte, value []byte, expiryTime time.Time) error {
	prefixedKey := append(keyPrefix, key...)

	batchKeys := make([][]byte, 2)
	batchValues := make([][]byte, 2)

	batchKeys[0] = prefixedKey
	batchValues[0] = value

	batchKeys[1] = buildExpiryKey(key, expiryTime)
	batchValues[1] = make([]byte, 0)

	return store.store.WriteBatch(batchKeys, batchValues)
}

// PutBatchWithExpiration adds multiple key-value pairs to the store that expire at a specified time.
func (store *ttlStore) PutBatchWithExpiration(keys [][]byte, values [][]byte, expiryTime time.Time) error {
	if len(keys) != len(values) {
		return fmt.Errorf("keys and values must have the same length (keys: %d, values: %d)", len(keys), len(values))
	}

	batchKeys := make([][]byte, len(keys)*2)
	batchValues := make([][]byte, len(keys)*2)

	for i, key := range keys {
		prefixedKey := append(keyPrefix, key...)

		batchKeys[i*2] = prefixedKey
		batchValues[i*2] = values[i]

		batchKeys[i*2+1] = buildExpiryKey(key, expiryTime)
		batchValues[i*2+1] = make([]byte, 0)
	}

	return store.store.WriteBatch(batchKeys, batchValues)
}

// Spawns a background goroutine that periodically checks for expired keys and deletes them.
func (store *ttlStore) expireKeysInBackground(gcPeriod time.Duration) {
	ticker := time.NewTicker(gcPeriod)
	go func() {
		for {
			select {
			case now := <-ticker.C:
				err := store.expireKeys(now)
				if err != nil {
					store.logger.Error("Error expiring keys", err)
					continue
				}
			case <-store.ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

// Delete all keys with a TTL that has expired.
func (store *ttlStore) expireKeys(now time.Time) error {
	it, err := store.store.NewIterator(expiryPrefix)
	if err != nil {
		return err
	}
	defer it.Release()

	keysToDelete := make([][]byte, 0)

	for it.Next() {
		expiryKey := it.Key()
		baseKey, expiryTimestamp := parseExpiryKey(expiryKey)

		if expiryTimestamp.After(now) {
			// No more values to expire
			break
		}

		prefixedBaseKey := append(keyPrefix, baseKey...)
		keysToDelete = append(keysToDelete, expiryKey)
		keysToDelete = append(keysToDelete, prefixedBaseKey)

		if len(keysToDelete) >= maxDeletionBatchSize {
			err = store.store.DeleteBatch(keysToDelete)
			if err != nil {
				return err
			}
			keysToDelete = make([][]byte, 0)
		}
	}

	if len(keysToDelete) > 0 {
		return store.store.DeleteBatch(keysToDelete)
	}
	return nil
}

func (store *ttlStore) Put(key []byte, value []byte) error {
	prefixedKey := append(keyPrefix, key...)
	return store.store.Put(prefixedKey, value)
}

func (store *ttlStore) Get(key []byte) ([]byte, error) {
	prefixedKey := append(keyPrefix, key...)
	return store.store.Get(prefixedKey)
}

func (store *ttlStore) Delete(key []byte) error {
	prefixedKey := append(keyPrefix, key...)
	return store.store.Delete(prefixedKey)
}

func (store *ttlStore) DeleteBatch(keys [][]byte) error {
	prefixedKeys := make([][]byte, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = append(keyPrefix, key...)
	}
	return store.store.DeleteBatch(prefixedKeys)
}

func (store *ttlStore) WriteBatch(keys, values [][]byte) error {
	prefixedKeys := make([][]byte, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = append(keyPrefix, key...)
	}
	return store.store.WriteBatch(prefixedKeys, values)
}

type ttlIterator struct {
	baseIterator iterator.Iterator
}

func (it *ttlIterator) First() bool {
	return it.baseIterator.First()
}

func (it *ttlIterator) Last() bool {
	return it.baseIterator.Last()
}

func (it *ttlIterator) Seek(key []byte) bool {
	return it.baseIterator.Seek(key)
}

func (it *ttlIterator) Next() bool {
	return it.baseIterator.Next()
}

func (it *ttlIterator) Prev() bool {
	return it.baseIterator.Prev()
}

func (it *ttlIterator) Release() {
	it.baseIterator.Release()
}

func (it *ttlIterator) SetReleaser(releaser util.Releaser) {
	it.baseIterator.SetReleaser(releaser)
}

func (it *ttlIterator) Valid() bool {
	return it.baseIterator.Valid()
}

func (it *ttlIterator) Error() error {
	return it.baseIterator.Error()
}

func (it *ttlIterator) Key() []byte {
	baseKey := it.baseIterator.Key()
	return baseKey[len(keyPrefix):]
}

func (it *ttlIterator) Value() []byte {
	return it.baseIterator.Value()
}

func (store *ttlStore) NewIterator(prefix []byte) (iterator.Iterator, error) {
	prefixedPrefix := append(keyPrefix, prefix...)
	baseIterator, err := store.store.NewIterator(prefixedPrefix)
	if err != nil {
		return nil, err
	}

	return &ttlIterator{
		baseIterator: baseIterator,
	}, nil
}

func (store *ttlStore) Shutdown() error {
	return store.store.Shutdown()
}

func (store *ttlStore) Destroy() error {
	return store.store.Destroy()
}
