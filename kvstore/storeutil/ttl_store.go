package storeutil

import (
	"fmt"
	"github.com/Layr-Labs/eigenda/kvstore"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
	"strconv"
	"time"
)

var _ kvstore.Store = &TTLStore{}

// TTLStore adds a time-to-live (TTL) capability to the store.
type TTLStore struct {
	store kvstore.Store
}

// TTLWrapper extends the given store with TTL capabilities.
func TTLWrapper(store kvstore.Store) *TTLStore {
	return &TTLStore{
		store: store,
	}
}

var keyPrefix = []byte("k")
var expiryPrefix = []byte("e")

var maxInt64 = int64(^uint64(0) >> 1)
var maxIntLengthBase16 = len(fmt.Sprintf("%x", maxInt64))
var expiryKeyPadding = fmt.Sprintf("%%0%dx", maxIntLengthBase16)

// PutWithTTL adds a key-value pair to the store that expires after a specified time-to-live (TTL).
// Key is eventually deleted after the TTL elapses.
func (store *TTLStore) PutWithTTL(key []byte, value []byte, ttl time.Duration) error {
	expiryTime := time.Now().Add(ttl)
	return store.PutWithExpiration(key, value, expiryTime)
}

// PutWithExpiration adds a key-value pair to the store that expires at a specified time.
// Key is eventually deleted after the expiry time.
func (store *TTLStore) PutWithExpiration(key []byte, value []byte, expiryTime time.Time) error {
	prefixedKey := append(keyPrefix, key...)

	batchKeys := make([][]byte, 2)
	batchValues := make([][]byte, 2)

	batchKeys[0] = prefixedKey
	batchValues[0] = value

	// The expiry key takes the form of the string "e<expiry timestamp in hexadecimal>".
	// The expiry timestamp is padded with zeros to ensure that the expiry key is lexicographically
	// ordered by time of expiry.
	// TODO verify padding
	expiryKey := append(expiryPrefix, []byte(fmt.Sprintf(expiryKeyPadding, expiryTime))...)

	batchKeys[1] = expiryKey
	batchValues[1] = prefixedKey

	return store.store.WriteBatch(batchKeys, batchValues)
}

// Delete all keys with a TTL that has expired.
func (store *TTLStore) expireKeys(now time.Time) error {

	// TODO add really strong documentation

	it, err := store.store.NewIterator(expiryPrefix)
	if err != nil {
		return err
	}
	defer it.Release()

	keysToDelete := make([][]byte, 0)

	for it.Next() {
		expiryKey := it.Key()
		expiryHex := string(expiryKey[len(expiryPrefix):])
		expiryValue, err := strconv.ParseUint(expiryHex, 16, 64)
		if err != nil {
			return err
		}
		expiryTime := time.Unix(0, int64(expiryValue))

		if expiryTime.After(now) {
			// No more values to expire
			return nil
		}

		keysToDelete = append(keysToDelete, it.Value())
	}

	return store.DeleteBatch(keysToDelete)
}

func (store *TTLStore) Put(key []byte, value []byte) error {
	prefixedKey := append(keyPrefix, key...)
	return store.store.Put(prefixedKey, value)
}

func (store *TTLStore) Get(key []byte) ([]byte, error) {
	prefixedKey := append(keyPrefix, key...)
	return store.store.Get(prefixedKey)
}

func (store *TTLStore) Delete(key []byte) error {
	prefixedKey := append(keyPrefix, key...)
	return store.store.Delete(prefixedKey)
}

func (store *TTLStore) DeleteBatch(keys [][]byte) error {
	prefixedKeys := make([][]byte, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = append(keyPrefix, key...)
	}
	return store.store.DeleteBatch(prefixedKeys)
}

func (store *TTLStore) WriteBatch(keys, values [][]byte) error {
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

func (store *TTLStore) NewIterator(prefix []byte) (iterator.Iterator, error) {
	prefixedPrefix := append(keyPrefix, prefix...)
	baseIterator, err := store.store.NewIterator(prefixedPrefix)
	if err != nil {
		return nil, err
	}

	return &ttlIterator{
		baseIterator: baseIterator,
	}, nil
}

func (store *TTLStore) Shutdown() error {
	return store.store.Shutdown()
}

func (store *TTLStore) Destroy() error {
	return store.store.Destroy()
}
