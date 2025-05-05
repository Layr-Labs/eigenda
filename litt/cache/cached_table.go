package cache

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common/cache"
	"github.com/Layr-Labs/eigenda/litt"
	"github.com/Layr-Labs/eigenda/litt/types"
	"github.com/Layr-Labs/eigenda/litt/util"
)

var _ litt.ManagedTable = &cachedTable{}

// cachedTable wraps a table and adds caching functionality.
type cachedTable struct {
	base  litt.ManagedTable
	cache cache.Cache[string, []byte]
}

// NewCachedTable creates a new table out of a base table and a cache.
func NewCachedTable(base litt.ManagedTable, cache cache.Cache[string, []byte]) litt.ManagedTable {
	return &cachedTable{
		base:  base,
		cache: cache,
	}
}

func (c *cachedTable) KeyCount() uint64 {
	return c.base.KeyCount()
}

func (c *cachedTable) Size() uint64 {
	return c.base.Size()
}

func (c *cachedTable) Name() string {
	return c.base.Name()
}

func (c *cachedTable) Put(key []byte, value []byte) error {
	err := c.base.Put(key, value)
	if err != nil {
		return err
	}
	c.cache.Put(string(key), value)
	return nil
}

func (c *cachedTable) PutBatch(batch []*types.KVPair) error {
	err := c.base.PutBatch(batch)
	if err != nil {
		return err
	}
	for _, kv := range batch {
		c.cache.Put(util.UnsafeBytesToString(kv.Key), kv.Value)
	}
	return nil
}

func (c *cachedTable) Get(key []byte) (value []byte, exists bool, err error) {
	value, exists, _, err = c.CacheAwareGet(key, false)
	return value, exists, err
}

// In theory, there is a race condition here where call to CacheAwareGet() made concurrently with a call to Put()
// might find the data to exist but not to be hot. This is not a problem though, since it will be hard to trigger and
// since it is not a violation of the consistency/correctness guarantees made by LittDB. Caching is inherently a
// "best effort" optimization, and so it's not worth adding extra locking in order to prevent this edge case.
//
// Scenario:
// - Thread A calls Put() on key K
// - Thread B calls CacheAwareGet() on key K with onlyReadFromCache set to true
// - Thread B checks the cache, and finds that the value is not there
// - Thread A finishes the Put() and returns. LittDB flushes the value out to disk.
// - Thread A gets to the part of CacheAwareGet() where it checks the base table for the value. Since the
//   base table has flushed the value out to disk, it says that the value exists but does not fetch it since
//   onlyReadFromCache is true.

func (c *cachedTable) CacheAwareGet(
	key []byte,
	onlyReadFromCache bool,
) (value []byte, exists bool, hot bool, err error) {

	stringKey := util.UnsafeBytesToString(key)

	value, exists = c.cache.Get(stringKey)
	if exists {
		// The value is in the cache
		return value, true, true, nil
	}

	value, exists, hot, err = c.base.CacheAwareGet(key, onlyReadFromCache)
	if err != nil {
		return nil, false, false, err
	}

	if exists {
		c.cache.Put(stringKey, value)
	}

	return value, exists, hot, nil
}

func (c *cachedTable) Exists(key []byte) (exists bool, err error) {
	_, exists = c.cache.Get(util.UnsafeBytesToString(key))
	if exists {
		return true, nil
	}

	return c.base.Exists(key)
}

func (c *cachedTable) Flush() error {
	return c.base.Flush()
}

func (c *cachedTable) SetTTL(ttl time.Duration) error {
	return c.base.SetTTL(ttl)
}

func (c *cachedTable) SetCacheSize(size uint64) error {
	c.cache.SetMaxWeight(size)
	err := c.base.SetCacheSize(size)
	if err != nil {
		return fmt.Errorf("failed to set base table cache size: %w", err)
	}
	return nil
}

func (c *cachedTable) Close() error {
	return c.base.Close()
}

func (c *cachedTable) Destroy() error {
	return c.base.Destroy()
}

func (c *cachedTable) SetShardingFactor(shardingFactor uint32) error {
	return c.base.SetShardingFactor(shardingFactor)
}

func (c *cachedTable) RunGC() error {
	return c.base.RunGC()
}
