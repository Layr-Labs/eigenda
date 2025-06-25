package cache

import (
	"weak"
)

var _ Cache[string, string] = &weakCache[string, string]{}

// TODO unit test

// A weakCache wraps another cache. It uses weak pointers to hold values, making it so that the garbage collector
// can reclaim memory used by values if there is extreme memory pressure.
// Works with any type V - internally creates pointers for weak reference management.
type weakCache[K comparable, V any] struct {
	cache Cache[K, weak.Pointer[V]]
}

// Create a new weak cache by wrapping an existing cache.
func NewWeakCacheWrapper[K comparable, V any](cache Cache[K, weak.Pointer[V]]) Cache[K, V] {
	return &weakCache[K, V]{
		cache: cache,
	}
}

// Create a new weak cache. The base cache type is a FIFO cache.
func NewWeakFIFOCache[K comparable, V any](
	maxWeight uint64,
	calculator WeightCalculator[K, V],
	metrics *CacheMetrics,
) Cache[K, V] {

	var wrappedCalculator WeightCalculator[K, weak.Pointer[V]]
	if calculator != nil {
		wrappedCalculator = func(key K, value weak.Pointer[V]) uint64 {
			v := value.Value()
			if v == nil {
				// If the value has been garbage collected, we treat it as having zero weight.
				// In practice, we should never trigger this, since the inner cache always computes cache
				// weight before the object becomes eligible for garbage collection (the outer context will
				// be holding a strong reference to the value at the moment it is put into the cache).
				return 0
			}

			return calculator(key, *v)
		}
	}

	baseCache := NewFIFOCache[K, weak.Pointer[V]](
		maxWeight,
		wrappedCalculator,
		metrics)

	return NewWeakCacheWrapper[K, V](baseCache)
}

func (w *weakCache[K, V]) Get(key K) (V, bool) {
	pointer, ok := w.cache.Get(key)
	if !ok {
		// The value is not in the cache.
		var zero V
		return zero, false
	}

	value := pointer.Value()
	if value == nil {
		// The value has been garbage collected, pretend like the value doesn't exist.
		var zero V
		return zero, false
	}

	return *value, true
}

func (w *weakCache[K, V]) Put(key K, value V) {
	// Create a copy of the value on the heap so we can take its address
	valueCopy := value
	w.cache.Put(key, weak.Make(&valueCopy))
}

func (w *weakCache[K, V]) Size() int {
	return w.cache.Size()
}

func (w *weakCache[K, V]) Weight() uint64 {
	return w.cache.Weight()
}

func (w *weakCache[K, V]) SetMaxWeight(capacity uint64) {
	w.cache.SetMaxWeight(capacity)
}
