package cache

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common/cache"
	"github.com/Layr-Labs/eigenda/litt"
	"github.com/Layr-Labs/eigenda/litt/types"
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
		c.cache.Put(string(kv.Key), kv.Value)
	}
	return nil
}

func (c *cachedTable) Get(key []byte) ([]byte, bool, error) {

	value, ok := c.cache.Get(string(key))
	if ok {
		return value, true, nil
	}

	value, ok, err := c.base.Get(key)
	if err != nil {
		return nil, false, err
	}

	if ok {
		c.cache.Put(string(key), value)
	}

	return value, ok, nil
}

func (c *cachedTable) Flush() error {
	return c.base.Flush()
}

func (c *cachedTable) SetTTL(ttl time.Duration) error {
	return c.base.SetTTL(ttl)
}

func (c *cachedTable) SetCacheSize(size uint64) error {
	c.cache.SetCapacity(size)
	err := c.base.SetCacheSize(size)
	if err != nil {
		return fmt.Errorf("failed to set base table cache size: %w", err)
	}
	return nil
}

func (c *cachedTable) Stop() error {
	return c.base.Stop()
}

func (c *cachedTable) Destroy() error {
	return c.base.Destroy()
}

func (c *cachedTable) SetShardingFactor(shardingFactor uint32) error {
	return c.base.SetShardingFactor(shardingFactor)
}
