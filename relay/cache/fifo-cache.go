package cache

import (
	"github.com/emirpasic/gods/queues"
	"github.com/emirpasic/gods/queues/linkedlistqueue"
)

var _ Cache[string, string] = &FIFOCache[string, string]{}

// FIFOCache is a cache that evicts the least recently added item when the cache is full. Useful for situations
// where time of addition is a better predictor of future access than time of most recent access.
type FIFOCache[K comparable, V any] struct {
	weightCalculator WeightCalculator[K, V]

	currentWeight   uint64
	maxWeight       uint64
	data            map[K]V
	expirationQueue queues.Queue
}

// NewFIFOCache creates a new FIFOCache. If the calculator is nil, the weight of each key-value pair will be 1.
func NewFIFOCache[K comparable, V any](
	maxWeight uint64,
	calculator WeightCalculator[K, V]) Cache[K, V] {

	if calculator == nil {
		calculator = func(K, V) uint64 { return 1 }
	}

	return &FIFOCache[K, V]{
		maxWeight:        maxWeight,
		data:             make(map[K]V),
		weightCalculator: calculator,
		expirationQueue:  linkedlistqueue.New(),
	}
}

func (f *FIFOCache[K, V]) Get(key K) (V, bool) {
	val, ok := f.data[key]
	return val, ok
}

func (f *FIFOCache[K, V]) Put(key K, value V) {
	weight := f.weightCalculator(key, value)
	if weight > f.maxWeight {
		// this item won't fit in the cache no matter what we evict
		return
	}

	old, ok := f.data[key]
	f.currentWeight += weight
	f.data[key] = value
	if ok {
		oldWeight := f.weightCalculator(key, old)
		f.currentWeight -= oldWeight
	} else {
		f.expirationQueue.Enqueue(key)
	}

	if f.currentWeight < f.maxWeight {
		// no need to evict anything
		return
	}

	for f.currentWeight > f.maxWeight {
		val, _ := f.expirationQueue.Dequeue()
		keyToEvict := val.(K)
		weightToEvict := f.weightCalculator(keyToEvict, f.data[keyToEvict])
		delete(f.data, keyToEvict)
		f.currentWeight -= weightToEvict
	}
}

func (f *FIFOCache[K, V]) Size() int {
	return len(f.data)
}

func (f *FIFOCache[K, V]) Weight() uint64 {
	return f.currentWeight
}
