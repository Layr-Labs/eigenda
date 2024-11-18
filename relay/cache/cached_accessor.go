package cache

import (
	lru "github.com/hashicorp/golang-lru/v2"
	"sync"
)

// CachedAccessor is an interface for accessing a resource that is cached. It assumes that cache misses
// are expensive, and prevents multiple concurrent cache misses for the same key.
type CachedAccessor[K comparable, V any] interface {
	// Get returns the value for the given key. If the value is not in the cache, it will be fetched using the Accessor.
	Get(key K) (V, error)
}

// Accessor is function capable of fetching a value from a resource. Used by CachedAccessor when there is a cache miss.
type Accessor[K comparable, V any] func(key K) (V, error)

// accessResult is a struct that holds the result of an Accessor call.
type accessResult[V any] struct {
	// wg.Wait() will block until the value is fetched.
	wg sync.WaitGroup
	// value is the value fetched by the Accessor, or nil if there was an error.
	value V
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

	// lookupsInProgress has an entry for each key that is currently being looked up via the accessor. The value
	// is written into the channel when it is eventually fetched. If a key is requested more than once while a
	// lookup in progress, the second (and following) requests will wait for the result of the first lookup
	// to be written into the channel.
	lookupsInProgress map[K]*accessResult[V]

	// cache is the LRU cache used to store values fetched by the accessor.
	cache *lru.Cache[K, V]

	// concurrencyLimiter is a channel used to limit the number of concurrent lookups that can be in progress.
	concurrencyLimiter chan struct{}

	// lock is used to protect the cache and lookupsInProgress map.
	cacheLock sync.Mutex

	// accessor is the function used to fetch values that are not in the cache.
	accessor Accessor[K, V]
}

// NewCachedAccessor creates a new CachedAccessor. The cacheSize parameter specifies the maximum number of items
// that can be stored in the cache. The concurrencyLimit parameter specifies the maximum number of concurrent
// lookups that can be in progress at any given time. If a greater number of lookups are requested, the excess
// lookups will block until a lookup completes. If concurrencyLimit is zero, then no limits are imposed. The accessor
// parameter is the function used to fetch values that are not in the cache.
func NewCachedAccessor[K comparable, V any](
	cacheSize int,
	concurrencyLimit int,
	accessor Accessor[K, V]) (CachedAccessor[K, V], error) {

	cache, err := lru.New[K, V](cacheSize)
	if err != nil {
		return nil, err
	}

	lookupsInProgress := make(map[K]*accessResult[V])

	var concurrencyLimiter chan struct{}
	if concurrencyLimit > 0 {
		concurrencyLimiter = make(chan struct{}, concurrencyLimit)
	}

	return &cachedAccessor[K, V]{
		cache:              cache,
		concurrencyLimiter: concurrencyLimiter,
		accessor:           accessor,
		lookupsInProgress:  lookupsInProgress,
	}, nil
}

func newAccessResult[V any]() *accessResult[V] {
	result := &accessResult[V]{
		wg: sync.WaitGroup{},
	}
	result.wg.Add(1)
	return result
}

func (c *cachedAccessor[K, V]) Get(key K) (V, error) {

	c.cacheLock.Lock()

	// first, attempt to get the value from the cache
	v, ok := c.cache.Get(key)
	if ok {
		c.cacheLock.Unlock()
		return v, nil
	}

	// if that fails, check if a lookup is already in progress. If not, start a new one.
	result, alreadyLoading := c.lookupsInProgress[key]
	if !alreadyLoading {
		result = newAccessResult[V]()
		c.lookupsInProgress[key] = result
	}

	c.cacheLock.Unlock()

	if alreadyLoading {
		// The result is being fetched on another goroutine. Wait for it to finish.
		result.wg.Wait()
		return result.value, result.err
	} else {
		// We are the first goroutine to request this key.

		if c.concurrencyLimiter != nil {
			c.concurrencyLimiter <- struct{}{}
		}

		value, err := c.accessor(key)

		if c.concurrencyLimiter != nil {
			<-c.concurrencyLimiter
		}

		c.cacheLock.Lock()

		// Update the cache if the fetch was successful.
		if err == nil {
			c.cache.Add(key, value)
		}

		// Provide the result to all other goroutines that may be waiting for it.
		result.err = err
		result.value = value
		result.wg.Done()

		// Clean up the lookupInProgress map.
		delete(c.lookupsInProgress, key)

		c.cacheLock.Unlock()

		return value, err
	}
}
