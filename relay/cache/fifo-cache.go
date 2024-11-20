package cache

import (
	"errors"
	"github.com/Layr-Labs/eigenda/common/queue"
)

var _ Cache[string, string] = &FIFOCache[string, string]{}

// FIFOCache is a cache that evicts the least recently added item when the cache is full. Useful for situations
// where time of addition is a better predictor of future access than time of most recent access.
type FIFOCache[K comparable, V any] struct {
	weightCalculator WeightCalculator[K, V]

	currentWeight   uint64
	maxWeight       uint64
	data            map[K]V
	expirationQueue queue.Queue[K]
}

// NewFIFOCache creates a new FIFOCache.
func NewFIFOCache[K comparable, V any](maxWeight uint64) *FIFOCache[K, V] {
	defaultWeightCalculator := func(key K, value V) uint64 {
		return uint64(1)
	}

	return &FIFOCache[K, V]{
		maxWeight:        maxWeight,
		data:             make(map[K]V),
		weightCalculator: defaultWeightCalculator,
		expirationQueue:  &queue.LinkedQueue[K]{},
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
		f.expirationQueue.Push(key)
	}

	if f.currentWeight < f.maxWeight {
		// no need to evict anything
		return
	}

	for f.currentWeight > f.maxWeight {
		keyToEvict, _ := f.expirationQueue.Pop()
		weightToEvict := f.weightCalculator(keyToEvict, f.data[keyToEvict])
		delete(f.data, keyToEvict)
		f.currentWeight -= weightToEvict
	}
}

func (f *FIFOCache[K, V]) WithWeightCalculator(weightCalculator WeightCalculator[K, V]) error {
	if f.Size() > 0 {
		return errors.New("cannot set weight calculator on non-empty cache")
	}
	f.weightCalculator = weightCalculator
	return nil
}

func (f *FIFOCache[K, V]) Size() int {
	return len(f.data)
}

func (f *FIFOCache[K, V]) Weight() uint64 {
	return f.currentWeight
}
