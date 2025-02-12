package v2

import (
	"sort"
	"sync"
	"time"
)

const (
	bucketMinutes = 1
)

// FeedBucket represents a cached bucket of feed items
type FeedBucket[T any] struct {
	StartTime time.Time
	Items     []T
}

// FeedCache implements a time-based caching system for feed items
type FeedCache[T any] struct {
	cache        map[int64]FeedBucket[T]
	mutex        sync.RWMutex
	bucketSize   time.Duration
	maxBuckets   int
	getTimestamp func(T) time.Time
}

// NewFeedCache creates a new cache with specified bucket size and maximum number of buckets
func NewFeedCache[T any](
	bucketSizeMinutes,
	maxBuckets int,
	getTimestamp func(T) time.Time,
) *FeedCache[T] {
	return &FeedCache[T]{
		cache:        make(map[int64]FeedBucket[T]),
		bucketSize:   time.Duration(bucketSizeMinutes) * time.Minute,
		maxBuckets:   maxBuckets,
		getTimestamp: getTimestamp,
	}
}

// getBucketKey returns the cache key for a given timestamp
func (c *FeedCache[T]) getBucketKey(timestamp time.Time) int64 {
	return timestamp.Truncate(c.bucketSize).Unix()
}

// GetItemsInRange retrieves items within the specified time range [startTime, endTime).
// The interval is inclusive of startTime and exclusive of endTime.
// fetchFromDB is a function that retrieves items from the database for a specific time bucket.
// The items returned from fetchFromDB are guaranteed to be sorted by timestamp in ascending order.
func (c *FeedCache[T]) GetItemsInRange(
	startTime, endTime time.Time,
	fetchFromDB func(start, end time.Time) ([]T, error),
) ([]T, error) {
	var results []T

	// Calculate all bucket keys needed
	var bucketKeys []int64
	current := startTime.Truncate(c.bucketSize)
	for current.Before(endTime) {
		bucketKeys = append(bucketKeys, c.getBucketKey(current))
		current = current.Add(c.bucketSize)
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Process each bucket
	for _, bucketKey := range bucketKeys {
		bucket, exists := c.cache[bucketKey]
		if !exists {
			// Convert key back to time
			bucketStart := time.Unix(bucketKey, 0)
			bucketEnd := bucketStart.Add(c.bucketSize)

			// Fetch items from database
			items, err := fetchFromDB(bucketStart, bucketEnd)
			if err != nil {
				return nil, err
			}

			// Cache the bucket
			bucket = FeedBucket[T]{
				StartTime: bucketStart,
				Items:     items,
			}
			c.cache[bucketKey] = bucket

			// Cleanup old buckets if needed
			c.cleanupOldBuckets()
		}

		// Filter items within the requested range
		for _, item := range bucket.Items {
			ts := c.getTimestamp(item)
			if (ts.Equal(startTime) || ts.After(startTime)) &&
				ts.Before(endTime) {
				results = append(results, item)
			}
		}
	}

	// Results are already sorted by timestamp since each bucket is pre-sorted from fetchFromDB
	// and we process buckets in chronological order
	return results, nil
}

// cleanupOldBuckets removes the oldest buckets when cache size exceeds maxBuckets
func (c *FeedCache[T]) cleanupOldBuckets() {
	if len(c.cache) <= c.maxBuckets {
		return
	}

	// Convert keys to times and sort
	var times []int64
	for key := range c.cache {
		times = append(times, key)
	}
	sort.Slice(times, func(i, j int) bool {
		return times[i] < times[j]
	})

	// Remove oldest buckets
	for i := 0; i < len(times)-c.maxBuckets; i++ {
		delete(c.cache, times[i])
	}
}
