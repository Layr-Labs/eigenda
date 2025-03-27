package memtable

import (
	"fmt"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/litt"
	"github.com/Layr-Labs/eigenda/litt/types"
	"github.com/emirpasic/gods/queues"
	"github.com/emirpasic/gods/queues/linkedlistqueue"
)

var _ litt.ManagedTable = &memTable{}

// expirationRecord is a record of when a key was inserted into the table, and for when it should be deleted.
type expirationRecord struct {
	// The time at which the key was inserted into the table.
	creationTime time.Time
	// A stringified version of the key.
	key string
}

// memTable is a simple implementation of a Table that stores its data in memory.
type memTable struct {
	// A function that returns the current time.
	timeSource func() time.Time

	// The name of the table.
	name string

	// The time-to-live for data in this table.
	ttl time.Duration

	// The actual data store.
	data map[string][]byte

	// Keeps track of when data should be deleted.
	expirationQueue queues.Queue

	// Protects access to data and expirationQueue.
	//
	// This implementation could be made with smaller granularity locks to improve multithreaded performance,
	// at the cost of code complexity. But since this implementation is primary intended for use in tests,
	// such optimization is not necessary.
	lock sync.RWMutex
}

// NewMemTable creates a new in-memory table.
func NewMemTable(timeSource func() time.Time, name string, ttl time.Duration) litt.ManagedTable {
	return &memTable{
		timeSource:      timeSource,
		name:            name,
		ttl:             ttl,
		data:            make(map[string][]byte),
		expirationQueue: linkedlistqueue.New(),
	}
}

func (m *memTable) Size() uint64 {
	// Technically speaking, this table stores zero bytes on disk, and this method
	// is contractually obligated to return only the size of the data on disk.
	return 0
}

func (m *memTable) Name() string {
	return m.name
}

func (m *memTable) KeyCount() uint64 {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return uint64(len(m.data))
}

func (m *memTable) Put(key []byte, value []byte) error {
	stringKey := string(key)
	expiration := &expirationRecord{
		creationTime: m.timeSource(),
		key:          stringKey,
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	_, ok := m.data[stringKey]
	if ok {
		return fmt.Errorf("key %x already exists", key)
	}
	m.data[stringKey] = value
	m.expirationQueue.Enqueue(expiration)

	return nil
}

func (m *memTable) PutBatch(batch []*types.KVPair) error {
	for _, kv := range batch {
		err := m.Put(kv.Key, kv.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *memTable) Get(key []byte) ([]byte, bool, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	value, ok := m.data[string(key)]
	if !ok {
		return nil, false, nil
	}

	return value, true, nil
}

func (m *memTable) Flush() error {
	// This is a no-op for a memory table. Memory tables are ephemeral by nature.
	return nil
}

func (m *memTable) SetTTL(ttl time.Duration) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.ttl = ttl
	return nil
}

func (m *memTable) doGarbageCollection() error { // TODO this is never called
	m.lock.Lock()
	defer m.lock.Unlock()

	now := m.timeSource()
	earliestPermittedCreationTime := now.Add(-m.ttl)

	for {
		item, ok := m.expirationQueue.Peek()
		if !ok {
			break
		}
		expiration := item.(*expirationRecord)
		if expiration.creationTime.After(earliestPermittedCreationTime) {
			break
		}
		m.expirationQueue.Dequeue()
		delete(m.data, expiration.key)
	}

	return nil
}

func (m *memTable) Destroy() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.data = make(map[string][]byte)
	m.expirationQueue.Clear()

	return nil
}

func (m *memTable) Stop() error {
	// no-op
	return nil
}

func (m *memTable) SetCacheSize(size uint64) error {
	// The memory table doesn't have a cache... it's already one giant cache.
	return nil
}

func (m *memTable) SetShardingFactor(shardingFactor uint32) error {
	// the memory table has no concept of sharding
	return nil
}

func (m *memTable) ScheduleImmediateGC() error {
	return m.doGarbageCollection()
}
