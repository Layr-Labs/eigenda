package cache

import (
	"time"

	"github.com/emirpasic/gods/queues"
	"github.com/emirpasic/gods/queues/linkedlistqueue"
)

var _ Cache[string, string] = &FIFOCache[string, string]{}

// FIFOCache is a cache that evicts the least recently added item when the cache is full. Useful for situations
// where time of addition is a better predictor of future access than time of most recent access.
type FIFOCache[K comparable, V any] struct {
	// A function that calculates the weight of a key-value pair. If nil, the weight of each key-value pair will be 1.
	weightCalculator WeightCalculator[K, V]

	// The sum of all weights of the items in the cache.
	currentWeight uint64

	// The maximum weight of the cache. If the current weight exceeds this, items will be evicted until the current
	// weight is less than or equal to this value.
	maxWeight uint64

	// The data stored in the cache.
	data map[K]*wrappedValue[V]

	// Keeps track of the order of insertion into the cache, is used to evict the least recently added item when the
	// cache is full.
	evictionQueue queues.Queue

	// encapsulates metrics for the cache
	metrics *CacheMetrics

	// nextSerialNumber is used to assign a unique serial number to each insertion record.
	nextSerialNumber uint64
}

// insertionRecord is a record of when a key was inserted into the cache, and is used to decide when it should be
// evicted.
type insertionRecord struct {
	// The key that was added to the cache.
	key any

	// Each entry into the cache is assigned a serial number. If an element is inserted twice, this serial number
	// prevents the eviction queue from being confused by the fact that there are multiple entries in the queue
	// for a single key in the cache.
	serialNumber uint64
}

// Contains the value with additional metadata used for eviction and metrics tracking.
type wrappedValue[V any] struct {

	// The value stored in the cache.
	value V

	// The serial number of the insertion.
	serialNumber uint64

	// The time at which the key was added to the cache.
	timestamp time.Time

	// The weight of the key-value pair. Important to keep this around instead of recomputing in case the base
	// type doesn't always return the same weight for the same key-value pair (e.g. if it's a weak pointer).
	weight uint64
}

// NewFIFOCache creates a new FIFOCache. If the calculator is nil, the weight of each key-value pair will be 1.
func NewFIFOCache[K comparable, V any](
	maxWeight uint64,
	calculator WeightCalculator[K, V],
	metrics *CacheMetrics) Cache[K, V] {

	if calculator == nil {
		calculator = func(K, V) uint64 { return 1 }
	}

	return &FIFOCache[K, V]{
		maxWeight:        maxWeight,
		data:             make(map[K]*wrappedValue[V]),
		weightCalculator: calculator,
		evictionQueue:    linkedlistqueue.New(),
		metrics:          metrics,
		nextSerialNumber: 0,
	}
}

func (f *FIFOCache[K, V]) Get(key K) (V, bool) {
	val, ok := f.data[key]

	if ok {
		return val.value, true
	}

	var zero V
	return zero, false
}

func (f *FIFOCache[K, V]) Put(key K, value V) {
	weight := f.weightCalculator(key, value)
	if weight > f.maxWeight {
		// this item won't fit in the cache no matter what we evict
		return
	}

	serialNumber := f.nextSerialNumber
	f.nextSerialNumber++

	old, alreadyPresent := f.data[key]
	if alreadyPresent {
		f.currentWeight -= f.weightCalculator(key, old.value)
	}

	f.currentWeight += weight

	f.data[key] = &wrappedValue[V]{
		value:        value,
		serialNumber: serialNumber,
		timestamp:    time.Now(),
		weight:       weight,
	}

	f.evictionQueue.Enqueue(&insertionRecord{
		key:          key,
		serialNumber: serialNumber,
	})

	if f.currentWeight > f.maxWeight {
		f.evict()
	}

	f.metrics.reportInsertion(weight)
	f.metrics.reportCurrentSize(len(f.data), f.currentWeight)
}

func (f *FIFOCache[K, V]) evict() {
	now := time.Now()

	for f.currentWeight > f.maxWeight {
		next, _ := f.evictionQueue.Dequeue()
		record := next.(*insertionRecord)

		keyToEvict := record.key.(K)

		current := f.data[keyToEvict]
		if current.serialNumber == record.serialNumber {
			// The record matches the value currently in the cache.

			delete(f.data, keyToEvict)
			f.currentWeight -= current.weight
			f.metrics.reportEviction(now.Sub(current.timestamp))
		}
	}
}

func (f *FIFOCache[K, V]) Size() int {
	return len(f.data)
}

func (f *FIFOCache[K, V]) Weight() uint64 {
	return f.currentWeight
}

func (f *FIFOCache[K, V]) SetMaxWeight(capacity uint64) {
	f.maxWeight = capacity
	f.evict()
}
