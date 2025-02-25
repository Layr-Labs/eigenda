package v2_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	v2 "github.com/Layr-Labs/eigenda/disperser/dataapi/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testItem struct {
	ts   time.Time
	data string
}

type testFetcher struct {
	fetchCount atomic.Int64
	baseTime   time.Time
}

func newTestFetcher(baseTime time.Time) *testFetcher {
	return &testFetcher{baseTime: baseTime}
}

func (tf *testFetcher) fetch(start, end time.Time) ([]*testItem, error) {
	tf.fetchCount.Add(1)
	var items []*testItem
	// Generate items every minute within the range [start, end)
	for t := start; t.Before(end); t = t.Add(time.Minute) {
		items = append(items, &testItem{
			ts:   t,
			data: t.Format(time.RFC3339),
		})
	}
	return items, nil
}

func (tf *testFetcher) getFetchCount() int64 {
	return tf.fetchCount.Load()
}

func setupTestCache(maxItems int) (*v2.FeedCache[testItem], *testFetcher, time.Time) {
	baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	fetcher := newTestFetcher(baseTime)

	timestampFn := func(item *testItem) time.Time {
		return item.ts
	}

	cache := v2.NewFeedCache(maxItems, fetcher.fetch, timestampFn, v2.Ascending)
	return cache, fetcher, baseTime
}

func TestFullCacheHit(t *testing.T) {
	cache, fetcher, baseTime := setupTestCache(100)

	// Initial fetch
	start := baseTime
	end := baseTime.Add(5 * time.Minute)
	_, err := cache.Get(start, end)
	require.NoError(t, err)
	assert.Equal(t, int64(1), fetcher.getFetchCount())

	// Query subset of cached range - should be a cache hit
	subStart := baseTime.Add(1 * time.Minute)
	subEnd := baseTime.Add(3 * time.Minute)
	items, err := cache.Get(subStart, subEnd)
	require.NoError(t, err)
	require.Len(t, items, 2)

	// Verify no additional DB fetch
	assert.Equal(t, int64(1), fetcher.getFetchCount())

	// Verify items are correct
	for i, item := range items {
		expectedTime := subStart.Add(time.Duration(i) * time.Minute)
		assert.Equal(t, expectedTime, item.ts)
	}
}

func TestPartialOverlap_NewerRange(t *testing.T) {
	cache, fetcher, baseTime := setupTestCache(100)

	// Initial fetch [0:00, 5:00)
	start := baseTime
	end := baseTime.Add(5 * time.Minute)
	_, err := cache.Get(start, end)
	require.NoError(t, err)
	assert.Equal(t, int64(1), fetcher.getFetchCount())

	// Query overlapping range [3:00, 8:00)
	newStart := baseTime.Add(3 * time.Minute)
	newEnd := baseTime.Add(8 * time.Minute)
	items, err := cache.Get(newStart, newEnd)
	require.NoError(t, err)
	require.Len(t, items, 5)

	// Should have one more fetch for the gap
	assert.Equal(t, int64(2), fetcher.getFetchCount())

	// Verify items are correct and in order
	for i, item := range items {
		expectedTime := newStart.Add(time.Duration(i) * time.Minute)
		assert.Equal(t, expectedTime, item.ts)
	}

	// Query within the extended range - should be a cache hit
	subStart := baseTime
	subEnd := baseTime.Add(8 * time.Minute)
	_, err = cache.Get(subStart, subEnd)
	require.NoError(t, err)
	assert.Equal(t, int64(2), fetcher.getFetchCount())
}

func TestPartialOverlap_OlderRange(t *testing.T) {
	cache, fetcher, baseTime := setupTestCache(100)

	// Initial fetch [5:00, 10:00)
	start := baseTime.Add(5 * time.Minute)
	end := baseTime.Add(10 * time.Minute)
	items, err := cache.Get(start, end)
	require.NoError(t, err)
	require.Len(t, items, 5)
	assert.Equal(t, int64(1), fetcher.getFetchCount())

	// Query older overlapping range [3:00, 6:00)
	oldStart := baseTime.Add(3 * time.Minute)
	oldEnd := baseTime.Add(6 * time.Minute)
	items, err = cache.Get(oldStart, oldEnd)
	require.NoError(t, err)
	require.Len(t, items, 3)

	// Should have one more fetch for the gap
	assert.Equal(t, int64(2), fetcher.getFetchCount())

	// Verify items are correct and in order
	for i, item := range items {
		expectedTime := oldStart.Add(time.Duration(i) * time.Minute)
		assert.Equal(t, expectedTime, item.ts)
	}

	// Query within the extended range - should be a cache hit
	subStart := baseTime.Add(3 * time.Minute)
	subEnd := baseTime.Add(10 * time.Minute)
	_, err = cache.Get(subStart, subEnd)
	require.NoError(t, err)
	assert.Equal(t, int64(2), fetcher.getFetchCount())
}

func TestPartialOverlap_NewerAndOlderRange(t *testing.T) {
	cache, fetcher, baseTime := setupTestCache(100)

	// Initial fetch [5:00, 10:00)
	start := baseTime.Add(5 * time.Minute)
	end := baseTime.Add(10 * time.Minute)
	items, err := cache.Get(start, end)
	require.NoError(t, err)
	require.Len(t, items, 5)
	assert.Equal(t, int64(1), fetcher.getFetchCount())

	// Query a larger range [3:00, 12:00)
	extendedStart := baseTime.Add(3 * time.Minute)
	extendedEnd := baseTime.Add(12 * time.Minute)
	items, err = cache.Get(extendedStart, extendedEnd)
	require.NoError(t, err)
	require.Len(t, items, 9)

	// Should have two more fetches (two gaps)
	assert.Equal(t, int64(3), fetcher.getFetchCount())

	// Verify items are correct and in order
	for i, item := range items {
		expectedTime := extendedStart.Add(time.Duration(i) * time.Minute)
		assert.Equal(t, expectedTime, item.ts)
	}

	// Query within the extended range - should be a cache hit
	_, err = cache.Get(extendedStart, extendedEnd)
	require.NoError(t, err)
	assert.Equal(t, int64(3), fetcher.getFetchCount())
}

func TestNoOverlap_NewerRange(t *testing.T) {
	cache, fetcher, baseTime := setupTestCache(100)

	// Initial fetch [0:00, 5:00)
	start := baseTime
	end := baseTime.Add(5 * time.Minute)
	_, err := cache.Get(start, end)
	require.NoError(t, err)
	assert.Equal(t, int64(1), fetcher.getFetchCount())

	// Query non-overlapping but newer range [10:00, 15:00)
	newStart := baseTime.Add(10 * time.Minute)
	newEnd := baseTime.Add(15 * time.Minute)
	items, err := cache.Get(newStart, newEnd)
	require.NoError(t, err)
	require.Len(t, items, 5)

	// Should have one more fetch for the new range
	assert.Equal(t, int64(2), fetcher.getFetchCount())

	// Query the old range again - should be a new fetch since cache was updated
	_, err = cache.Get(start, end)
	require.NoError(t, err)
	assert.Equal(t, int64(3), fetcher.getFetchCount())

	// Query the same range again - should hit the cache
	_, err = cache.Get(newStart, newEnd)
	require.NoError(t, err)
	assert.Equal(t, int64(3), fetcher.getFetchCount())
}

func TestNoOverlap_OlderRange(t *testing.T) {
	cache, fetcher, baseTime := setupTestCache(100)

	// Initial fetch for [5:00, 10:00)
	start := baseTime.Add(5 * time.Minute)
	end := baseTime.Add(10 * time.Minute)
	items, err := cache.Get(start, end)
	require.NoError(t, err)
	require.Len(t, items, 5)
	assert.Equal(t, int64(1), fetcher.getFetchCount())

	// Query older range [0:00, 3:00)
	oldStart := baseTime
	oldEnd := baseTime.Add(3 * time.Minute)
	items, err = cache.Get(oldStart, oldEnd)
	require.NoError(t, err)
	require.Len(t, items, 3)
	assert.Equal(t, int64(2), fetcher.getFetchCount()) // New fetch needed

	// Query cached range again - should still hit original cache
	items, err = cache.Get(start, end)
	require.NoError(t, err)
	require.Len(t, items, 5)
	assert.Equal(t, int64(2), fetcher.getFetchCount()) // No new fetch

	// Verify first item is from the original cache start time
	firstItem := items[0]
	assert.Equal(t, start, firstItem.ts)

	// Query that overlaps both ranges - should use cache for newer part
	overlapStart := baseTime.Add(2 * time.Minute)
	overlapEnd := baseTime.Add(7 * time.Minute)
	_, err = cache.Get(overlapStart, overlapEnd)
	require.NoError(t, err)
	assert.Equal(t, int64(3), fetcher.getFetchCount()) // One new fetch for older part
}

func TestEviction(t *testing.T) {
	cache, fetcher, baseTime := setupTestCache(3)

	// Fetch 5 minutes worth of data
	start := baseTime
	end := baseTime.Add(5 * time.Minute)
	items, err := cache.Get(start, end)
	require.NoError(t, err)
	require.Len(t, items, 5) // Got all items from DB
	assert.Equal(t, int64(1), fetcher.getFetchCount())

	// Get the items again, have partial cache hit but will still need a fetch
	items, err = cache.Get(start, end)
	require.NoError(t, err)
	require.Len(t, items, 5)
	assert.Equal(t, int64(2), fetcher.getFetchCount())

	expectedStart := baseTime
	for i, item := range items {
		expectedTime := expectedStart.Add(time.Duration(i) * time.Minute)
		assert.Equal(t, expectedTime, item.ts)
	}

	// Try to get the evicted items - should cause a new fetch
	oldStart := baseTime
	oldEnd := baseTime.Add(2 * time.Minute)
	items, err = cache.Get(oldStart, oldEnd)
	require.NoError(t, err)
	require.Len(t, items, 2)
	assert.Equal(t, int64(3), fetcher.getFetchCount()) // New fetch needed

	hitStart := baseTime.Add(2 * time.Minute)
	hitEnd := baseTime.Add(5 * time.Minute)
	items, err = cache.Get(hitStart, hitEnd)
	require.NoError(t, err)
	require.Len(t, items, 3)
	assert.Equal(t, int64(3), fetcher.getFetchCount())
}

func TestConcurrentAccess(t *testing.T) {
	cache, _, baseTime := setupTestCache(100)

	var wg sync.WaitGroup
	concurrentRequests := 10

	// Launch multiple goroutines to access cache concurrently
	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go func(offset int) {
			defer wg.Done()
			start := baseTime.Add(time.Duration(offset) * time.Minute)
			end := start.Add(5 * time.Minute)
			items, err := cache.Get(start, end)
			assert.NoError(t, err)
			assert.NotNil(t, items)
		}(i)
	}

	wg.Wait()
}
