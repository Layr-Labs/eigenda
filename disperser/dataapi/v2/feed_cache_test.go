package v2

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test item type with timestamp
type testItem struct {
	ts   time.Time
	data string
}

// Test fetcher with instrumentation to track fetch count
type testFetcher struct {
	fetchCount atomic.Int64
	baseTime   time.Time
}

func newTestFetcher(baseTime time.Time) *testFetcher {
	return &testFetcher{baseTime: baseTime}
}

// Implement fetch method matching the interface expected by FeedCache
func (tf *testFetcher) fetch(start, end time.Time, order SortOrder, limit int) ([]*testItem, error) {
	tf.fetchCount.Add(1)
	var items []*testItem

	count := 0

	if order == Ascending {
		// Generate items every minute within the range [start, end) in ascending order
		for t := start; t.Before(end); t = t.Add(time.Minute) {
			if limit > 0 && count >= limit {
				break
			}

			items = append(items, &testItem{
				ts:   t,
				data: t.Format(time.RFC3339),
			})
			count++
		}
	} else { // Descending order
		// Generate items every minute within the range [start, end) in descending order
		// Start from (end - 1 minute) and go backwards to start
		for t := end.Add(-time.Minute); !t.Before(start); t = t.Add(-time.Minute) {
			if limit > 0 && count >= limit {
				break
			}

			items = append(items, &testItem{
				ts:   t,
				data: t.Format(time.RFC3339),
			})
			count++
		}
	}

	return items, nil
}

func (tf *testFetcher) getFetchCount() int64 {
	return tf.fetchCount.Load()
}

// Setup helper for tests
func setupTestCache(maxItems int) (*FeedCache[testItem], *testFetcher, time.Time) {
	baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	fetcher := newTestFetcher(baseTime)

	timestampFn := func(item *testItem) time.Time {
		return item.ts
	}

	cache := NewFeedCache[testItem](
		maxItems,
		fetcher.fetch,
		timestampFn,
	)

	return cache, fetcher, baseTime
}

// Test a full cache hit scenario
func TestFullCacheHit(t *testing.T) {
	cache, fetcher, baseTime := setupTestCache(100)

	// Initial fetch
	start := baseTime
	end := baseTime.Add(5 * time.Minute)
	_, err := cache.Get(start, end, Ascending, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(1), fetcher.getFetchCount())

	// Query subset of cached range - should be a cache hit
	subStart := baseTime.Add(1 * time.Minute)
	subEnd := baseTime.Add(3 * time.Minute)
	items, err := cache.Get(subStart, subEnd, Ascending, 0)
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

// Test partial overlap with newer range
func TestPartialOverlap_NewerRange(t *testing.T) {
	cache, fetcher, baseTime := setupTestCache(100)

	// Initial fetch [0:00, 5:00)
	start := baseTime
	end := baseTime.Add(5 * time.Minute)
	_, err := cache.Get(start, end, Ascending, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(1), fetcher.getFetchCount())

	// Query overlapping range [3:00, 8:00)
	newStart := baseTime.Add(3 * time.Minute)
	newEnd := baseTime.Add(8 * time.Minute)
	items, err := cache.Get(newStart, newEnd, Ascending, 0)
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
	_, err = cache.Get(subStart, subEnd, Ascending, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(2), fetcher.getFetchCount())
}

// Test partial overlap with older range
func TestPartialOverlap_OlderRange(t *testing.T) {
	cache, fetcher, baseTime := setupTestCache(100)

	// Initial fetch [5:00, 10:00)
	start := baseTime.Add(5 * time.Minute)
	end := baseTime.Add(10 * time.Minute)
	items, err := cache.Get(start, end, Ascending, 0)
	require.NoError(t, err)
	require.Len(t, items, 5)
	assert.Equal(t, int64(1), fetcher.getFetchCount())

	// Query older overlapping range [3:00, 6:00)
	oldStart := baseTime.Add(3 * time.Minute)
	oldEnd := baseTime.Add(6 * time.Minute)
	items, err = cache.Get(oldStart, oldEnd, Ascending, 0)
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
	_, err = cache.Get(subStart, subEnd, Ascending, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(2), fetcher.getFetchCount())
}

// Test partial overlap with both newer and older range
func TestPartialOverlap_NewerAndOlderRange(t *testing.T) {
	cache, fetcher, baseTime := setupTestCache(100)

	// Initial fetch [5:00, 10:00)
	start := baseTime.Add(5 * time.Minute)
	end := baseTime.Add(10 * time.Minute)
	items, err := cache.Get(start, end, Ascending, 0)
	require.NoError(t, err)
	require.Len(t, items, 5)
	assert.Equal(t, int64(1), fetcher.getFetchCount())

	// Query a larger range [3:00, 12:00)
	extendedStart := baseTime.Add(3 * time.Minute)
	extendedEnd := baseTime.Add(12 * time.Minute)
	items, err = cache.Get(extendedStart, extendedEnd, Ascending, 0)
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
	_, err = cache.Get(extendedStart, extendedEnd, Ascending, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(3), fetcher.getFetchCount())
}

// Test no overlap with newer range
func TestNoOverlap_NewerRange(t *testing.T) {
	cache, fetcher, baseTime := setupTestCache(100)

	// Initial fetch [0:00, 5:00)
	start := baseTime
	end := baseTime.Add(5 * time.Minute)
	_, err := cache.Get(start, end, Ascending, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(1), fetcher.getFetchCount())

	// Query non-overlapping but newer range [10:00, 15:00)
	newStart := baseTime.Add(10 * time.Minute)
	newEnd := baseTime.Add(15 * time.Minute)
	items, err := cache.Get(newStart, newEnd, Ascending, 0)
	require.NoError(t, err)
	require.Len(t, items, 5)

	// Should have one more fetch for the new range
	assert.Equal(t, int64(2), fetcher.getFetchCount())

	// Query the new range again - should hit the cache
	_, err = cache.Get(newStart, newEnd, Ascending, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(2), fetcher.getFetchCount())
}

// Test no overlap with older range
func TestNoOverlap_OlderRange(t *testing.T) {
	cache, fetcher, baseTime := setupTestCache(100)

	// Initial fetch for [5:00, 10:00)
	start := baseTime.Add(5 * time.Minute)
	end := baseTime.Add(10 * time.Minute)
	items, err := cache.Get(start, end, Ascending, 0)
	require.NoError(t, err)
	require.Len(t, items, 5)
	assert.Equal(t, int64(1), fetcher.getFetchCount())

	// Query older range [0:00, 3:00)
	oldStart := baseTime
	oldEnd := baseTime.Add(3 * time.Minute)
	items, err = cache.Get(oldStart, oldEnd, Ascending, 0)
	require.NoError(t, err)
	require.Len(t, items, 3)
	assert.Equal(t, int64(2), fetcher.getFetchCount()) // New fetch needed

	// Query the new range again - should hit the cache
	_, err = cache.Get(oldStart, oldEnd, Ascending, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(2), fetcher.getFetchCount())
}

// Test cache eviction due to maxItems limit
func TestEviction(t *testing.T) {
	cache, fetcher, baseTime := setupTestCache(3)

	// Fetch 5 minutes worth of data
	start := baseTime
	end := baseTime.Add(5 * time.Minute)
	items, err := cache.Get(start, end, Ascending, 0)
	require.NoError(t, err)
	require.Len(t, items, 5) // Got all items from DB
	assert.Equal(t, int64(1), fetcher.getFetchCount())

	// Query the full range again - should be a partial cache hit
	// Only the most recent 3 items should be in cache due to maxItems
	start2 := baseTime
	end2 := baseTime.Add(5 * time.Minute)
	items2, err := cache.Get(start2, end2, Ascending, 0)
	require.NoError(t, err)
	require.Len(t, items2, 5)
	assert.Equal(t, int64(2), fetcher.getFetchCount()) // Need to fetch older items not in cache

	// Query just the most recent 3 items - should be a cache hit
	recentStart := baseTime.Add(2 * time.Minute)
	recentEnd := baseTime.Add(5 * time.Minute)
	items3, err := cache.Get(recentStart, recentEnd, Ascending, 0)
	require.NoError(t, err)
	require.Len(t, items3, 3)
	assert.Equal(t, int64(2), fetcher.getFetchCount()) // No new fetch needed
}

// Test concurrent access to cache
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
			items, err := cache.Get(start, end, Ascending, 0)
			assert.NoError(t, err)
			assert.NotNil(t, items)
		}(i)
	}

	wg.Wait()
}

// Test with limit parameter
func TestWithLimit(t *testing.T) {
	cache, fetcher, baseTime := setupTestCache(100)

	// Fetch with limit
	start := baseTime
	end := baseTime.Add(10 * time.Minute)
	limit := 3

	items, err := cache.Get(start, end, Ascending, limit)
	require.NoError(t, err)
	require.Len(t, items, limit)
	assert.Equal(t, int64(1), fetcher.getFetchCount())

	// Verify only got first 3 items
	for i, item := range items {
		expectedTime := start.Add(time.Duration(i) * time.Minute)
		assert.Equal(t, expectedTime, item.ts)
	}

	// Fetch with descending order and limit
	items, err = cache.Get(start, end, Descending, limit)
	require.NoError(t, err)
	require.Len(t, items, limit)
	assert.Equal(t, int64(2), fetcher.getFetchCount())

	// Verify items are in descending order (most recent first)
	for i, item := range items {
		expectedTime := end.Add(-time.Minute - time.Duration(i)*time.Minute)
		assert.Equal(t, expectedTime, item.ts)
	}
}

// Test how hasMore detection works with limit
func TestHasMoreDetection(t *testing.T) {
	cache, fetcher, baseTime := setupTestCache(100)

	// Fetch with a small limit to trigger hasMore logic
	start := baseTime
	end := baseTime.Add(20 * time.Minute)
	limit := 5

	items, err := cache.Get(start, end, Ascending, limit)
	require.NoError(t, err)
	require.Len(t, items, limit)
	assert.Equal(t, int64(1), fetcher.getFetchCount())

	// The cache should now have a conservative end time based on the last item returned
	// Fetch data that continues just after the cached range
	nextStart := baseTime.Add(5 * time.Minute) // Just after the last item in previous fetch
	nextEnd := baseTime.Add(10 * time.Minute)

	items, err = cache.Get(nextStart, nextEnd, Ascending, 0)
	require.NoError(t, err)
	require.Len(t, items, 5)

	// Should need a new fetch since we reached the limit in the first query
	assert.Equal(t, int64(2), fetcher.getFetchCount())

	// Verify the items are correct
	for i, item := range items {
		expectedTime := nextStart.Add(time.Duration(i) * time.Minute)
		assert.Equal(t, expectedTime, item.ts)
	}
}

// Test descending order functionality and mixed query orders
func TestDescendingOrder(t *testing.T) {
	cache, fetcher, baseTime := setupTestCache(100)

	// Initial fetch in descending order
	start := baseTime
	end := baseTime.Add(5 * time.Minute)
	items, err := cache.Get(start, end, Descending, 0)
	require.NoError(t, err)
	require.Len(t, items, 5)
	assert.Equal(t, int64(1), fetcher.getFetchCount())

	// Verify order: [4, 3, 2, 1, 0] minutes from baseTime
	for i, item := range items {
		expectedTime := end.Add(-time.Minute - time.Duration(i)*time.Minute)
		assert.Equal(t, expectedTime, item.ts)
	}

	// Query again in descending order - should use cache
	items2, err := cache.Get(start, end, Descending, 0)
	require.NoError(t, err)
	require.Len(t, items2, 5)
	assert.Equal(t, int64(1), fetcher.getFetchCount()) // No new fetch
	assert.Equal(t, items, items2)                     // Should be identical

	// Query in ascending order - should use cached data but return in different order
	itemsAsc, err := cache.Get(start, end, Ascending, 0)
	require.NoError(t, err)
	require.Len(t, itemsAsc, 5)
	assert.Equal(t, int64(1), fetcher.getFetchCount()) // No new fetch

	// Verify ascending order: [0, 1, 2, 3, 4] minutes from baseTime
	for i, item := range itemsAsc {
		expectedTime := start.Add(time.Duration(i) * time.Minute)
		assert.Equal(t, expectedTime, item.ts)
	}

	// Verify that the internal order is correctly maintained
	// by checking if a small subset is returned correctly in both orderings
	subStart := baseTime.Add(1 * time.Minute)
	subEnd := baseTime.Add(4 * time.Minute)

	// Get subset in ascending order
	subItemsAsc, err := cache.Get(subStart, subEnd, Ascending, 0)
	require.NoError(t, err)
	require.Len(t, subItemsAsc, 3)                     // Minutes 1, 2, 3
	assert.Equal(t, int64(1), fetcher.getFetchCount()) // Still using cache

	// Verify ascending subset
	for i, item := range subItemsAsc {
		expectedTime := subStart.Add(time.Duration(i) * time.Minute)
		assert.Equal(t, expectedTime, item.ts)
	}

	// Get subset in descending order
	subItemsDesc, err := cache.Get(subStart, subEnd, Descending, 0)
	require.NoError(t, err)
	require.Len(t, subItemsDesc, 3)                    // Minutes 3, 2, 1
	assert.Equal(t, int64(1), fetcher.getFetchCount()) // Still using cache

	// Verify descending subset
	for i, item := range subItemsDesc {
		expectedTime := subEnd.Add(-time.Minute - time.Duration(i)*time.Minute)
		assert.Equal(t, expectedTime, item.ts)
	}
}

// Test how cache handles partial overlap with limit
func TestCacheMergeWithLimit(t *testing.T) {
	cache, fetcher, baseTime := setupTestCache(100)

	// First fetch a range with a small limit
	start := baseTime
	end := baseTime.Add(10 * time.Minute)
	initialLimit := 5

	_, err := cache.Get(start, end, Ascending, initialLimit)
	require.NoError(t, err)
	assert.Equal(t, int64(1), fetcher.getFetchCount())

	// Now query a range that overlaps but extends beyond what's cached
	// Should only need to fetch the extension
	extendedStart := baseTime.Add(3 * time.Minute)
	extendedEnd := baseTime.Add(15 * time.Minute)

	items, err := cache.Get(extendedStart, extendedEnd, Ascending, 0)
	require.NoError(t, err)
	require.Len(t, items, 12) // From minute 3 to minute 14

	// Should have made another fetch for the extension
	assert.Equal(t, int64(2), fetcher.getFetchCount())

	// Verify first items come from cache and later items from new fetch
	for i, item := range items {
		expectedTime := extendedStart.Add(time.Duration(i) * time.Minute)
		assert.Equal(t, expectedTime, item.ts)
	}

	// Now query the same range in descending order
	// Should use the cache that we've built up
	itemsDesc, err := cache.Get(extendedStart, extendedEnd, Descending, 0)
	require.NoError(t, err)
	require.Len(t, itemsDesc, 12)
	assert.Equal(t, int64(2), fetcher.getFetchCount()) // No new fetch

	// Verify items are in descending order
	for i, item := range itemsDesc {
		expectedTime := extendedEnd.Add(-time.Minute - time.Duration(i)*time.Minute)
		assert.Equal(t, expectedTime, item.ts)
	}

	// Test with a limit in the opposite order
	newLimit := 6
	limitedItems, err := cache.Get(extendedStart, extendedEnd, Descending, newLimit)
	require.NoError(t, err)
	require.Len(t, limitedItems, newLimit)
	assert.Equal(t, int64(2), fetcher.getFetchCount()) // Still no new fetch

	// Should have the 6 most recent items (minutes 9-14)
	for i, item := range limitedItems {
		expectedTime := extendedEnd.Add(-time.Minute - time.Duration(i)*time.Minute)
		assert.Equal(t, expectedTime, item.ts)
	}
}

// Test invalid parameters
func TestInvalidParameters(t *testing.T) {
	cache, _, baseTime := setupTestCache(100)

	// Test with end before start
	_, err := cache.Get(baseTime.Add(5*time.Minute), baseTime, Ascending, 0)
	assert.Error(t, err)
}

// Test empty results
func TestEmptyResults(t *testing.T) {
	cache, fetcher, baseTime := setupTestCache(100)

	// Query for a time range far in the future
	start := baseTime.Add(1000 * time.Hour)
	end := start.Add(5 * time.Minute)

	items, err := cache.Get(start, end, Ascending, 0)
	require.NoError(t, err)
	assert.Empty(t, items)
	assert.Equal(t, int64(1), fetcher.getFetchCount())
}

// Test descending query with overlaps on both sides
func TestDescendingWithOverlap(t *testing.T) {
	cache, fetcher, baseTime := setupTestCache(100)

	// Initial fetch in ascending order (to populate cache)
	start := baseTime.Add(5 * time.Minute)
	end := baseTime.Add(10 * time.Minute)
	_, err := cache.Get(start, end, Ascending, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(1), fetcher.getFetchCount())

	// Now do a descending query that overlaps and extends the range
	extendedStart := baseTime.Add(3 * time.Minute)
	extendedEnd := baseTime.Add(12 * time.Minute)

	items, err := cache.Get(extendedStart, extendedEnd, Descending, 0)
	require.NoError(t, err)
	require.Len(t, items, 9) // 9 items in total (minutes 3-11)

	// Should have made two more fetches (for both extensions)
	assert.Equal(t, int64(3), fetcher.getFetchCount())

	// Verify the items are in descending order
	for i := 0; i < len(items)-1; i++ {
		assert.True(t, items[i].ts.After(items[i+1].ts))
	}

	// First item should be just before extended end
	assert.Equal(t, baseTime.Add(11*time.Minute), items[0].ts)

	// Last item should be extended start
	assert.Equal(t, extendedStart, items[len(items)-1].ts)

	// Now query again in ascending order, should use cached data
	itemsAsc, err := cache.Get(extendedStart, extendedEnd, Ascending, 0)
	require.NoError(t, err)
	require.Len(t, itemsAsc, 9)
	assert.Equal(t, int64(3), fetcher.getFetchCount()) // No new fetch

	// Verify ascending order
	for i := 0; i < len(itemsAsc)-1; i++ {
		assert.True(t, itemsAsc[i].ts.Before(itemsAsc[i+1].ts))
	}
}

// Test changing cache segment with no overlap
func TestCacheSegmentChange(t *testing.T) {
	cache, fetcher, baseTime := setupTestCache(10)

	// Initial fetch
	start1 := baseTime
	end1 := baseTime.Add(5 * time.Minute)
	_, err := cache.Get(start1, end1, Ascending, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(1), fetcher.getFetchCount())

	// Fetch a completely different time range - should replace the cache
	start2 := baseTime.Add(100 * time.Minute)
	end2 := baseTime.Add(105 * time.Minute)
	items2, err := cache.Get(start2, end2, Ascending, 0)
	require.NoError(t, err)
	require.Len(t, items2, 5)
	assert.Equal(t, int64(2), fetcher.getFetchCount())

	// Original range should now cause a cache miss
	_, err = cache.Get(start1, end1, Ascending, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(3), fetcher.getFetchCount())

	// New range should still be cached
	_, err = cache.Get(start2, end2, Ascending, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(3), fetcher.getFetchCount()) // No new fetch
}

// Test that Time.Equal values work correctly
func TestTimeEqualHandling(t *testing.T) {
	cache, fetcher, baseTime := setupTestCache(100)

	// Initial fetch
	start := baseTime
	end := baseTime.Add(5 * time.Minute)
	_, err := cache.Get(start, end, Ascending, 0)
	require.NoError(t, err)

	// Query with exactly the same start and end
	sameStart := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC) // Same as baseTime but new instance
	sameEnd := sameStart.Add(5 * time.Minute)
	_, err = cache.Get(sameStart, sameEnd, Ascending, 0)
	require.NoError(t, err)

	// Should be a cache hit
	assert.Equal(t, int64(1), fetcher.getFetchCount())
}
