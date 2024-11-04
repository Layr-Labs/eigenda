package relay

import (
	lru "github.com/hashicorp/golang-lru/v2"
	"sync"
)

// CachedAccessor is an interface for accessing a resource that is cached. It assumes that cache misses
// are expensive, and prevents multiple concurrent cache misses for the same key.
type CachedAccessor[K comparable, V any] interface {
	// Get returns the value for the given key. If the value is not in the cache, it will be fetched using the Accessor.
	Get(key K) (*V, error)
}

// Accessor is function capable of fetching a value from a resource. Used by CachedAccessor when there is a cache miss.
type Accessor[K comparable, V any] func(key K) (V, error)

// accessResult is a struct that holds the result of an Accessor call.
type accessResult[V any] struct {
	// wg.Wait() will block until the value is fetched.
	wg sync.WaitGroup
	// value is the value fetched by the Accessor, or nil if there was an error.
	value *V
	// err is the error returned by the Accessor, or nil if the fetch was successful.
	err error
}

var _ CachedAccessor[string, string] = &cachedAccessor[string, string]{}

// Future work: the cache used in this implementation is suboptimal when storing items that have a large
// variance in size. The current implementation uses a fixed size cache, which requires the cached to be
// sized to the largest item that will be stored. This cache should be replaced with an implementation
// whose size can be specified by memory footprint in bytes.

// cachedAccessor is an implementation of CachedAccessor.
type cachedAccessor[K comparable, V any] struct {
	// used to synchronize access to the cache and lookupsInProgress
	lock sync.Mutex

	// lookupsInProgress has an entry for each key that is currently being looked up via the accessor. The value
	// is written into the channel when it is eventually fetched. If a key is requested more than once while a
	// lookup in progress, the second (and following) requests will wait for the result of the first lookup
	// to be written into the channel.
	lookupsInProgress map[K]*accessResult[V]

	// cache is the LRU cache used to store values fetched by the accessor.
	cache *lru.Cache[K, *V]

	// accessor is the function used to fetch values that are not in the cache.
	accessor Accessor[K, *V]
}

// NewCachedAccessor creates a new CachedAccessor.
func NewCachedAccessor[K comparable, V any](cacheSize int, accessor Accessor[K, *V]) (CachedAccessor[K, V], error) {

	cache, err := lru.New[K, *V](cacheSize)
	if err != nil {
		return nil, err
	}

	return &cachedAccessor[K, V]{
		lock:              sync.Mutex{},
		lookupsInProgress: make(map[K]*accessResult[V]),
		cache:             cache,
		accessor:          accessor,
	}, nil
}

func (c *cachedAccessor[K, V]) newAccessResult() *accessResult[V] {
	result := &accessResult[V]{
		wg: sync.WaitGroup{},
	}
	result.wg.Add(1)
	return result
}

func (c *cachedAccessor[K, V]) Get(key K) (*V, error) {
	// first, attempt to get the value from the cache
	v, ok := c.cache.Get(key)
	if ok {
		return v, nil
	}

	// if that fails, check if a lookup is already in progress. If not, start a new one.
	c.lock.Lock()
	result, lookupInProgress := c.lookupsInProgress[key]
	if !lookupInProgress {
		// No lookup in progress, start a new one.
		result = c.newAccessResult()
		c.lookupsInProgress[key] = result
	}
	c.lock.Unlock()

	if lookupInProgress {
		// Wait for the value to be fetched by the other goroutine.
		result.wg.Wait()
		return result.value, result.err
	} else {
		// We are responsible for fetching the value.
		value, err := c.accessor(key)

		// Update the cache if the fetch was successful.
		if err == nil {
			c.cache.Add(key, value)
		}

		result.err = err
		result.value = value
		result.wg.Done()

		// Clean up the lookupInProgress map.
		c.lock.Lock()
		delete(c.lookupsInProgress, key)
		c.lock.Unlock()

		return value, err
	}
}
