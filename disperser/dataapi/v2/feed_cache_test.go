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

func roundUpToNextMinute(t time.Time) time.Time {
	if t.Truncate(time.Minute) == t {
		return t
	}
	return t.Truncate(time.Minute).Add(time.Minute)
}

// Implement fetch method matching the interface expected by FeedCache
func (tf *testFetcher) fetch(start, end time.Time, order v2.FetchOrder, limit int) ([]*testItem, error) {
	tf.fetchCount.Add(1)
	var items []*testItem

	// Round up next exact minute (i.e. simulating there are only data items at exact minutes)
	start = roundUpToNextMinute(start)
	count := 0

	if order == v2.Ascending {
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
	} else {
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

func (tf *testFetcher) getFetchCount() int {
	return int(tf.fetchCount.Load())
}

// Setup helper for tests
func setupTestCache(maxItems int) (*v2.FeedCache[testItem], *testFetcher, time.Time) {
	baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	fetcher := newTestFetcher(baseTime)

	timestampFn := func(item *testItem) time.Time {
		return item.ts
	}

	cache := v2.NewFeedCache[testItem](
		maxItems,
		fetcher.fetch,
		timestampFn,
	)

	return cache, fetcher, baseTime
}

// Test invalid parameters
func TestInvalidParameters(t *testing.T) {
	cache, _, baseTime := setupTestCache(100)

	// Test with end before start
	_, err := cache.Get(baseTime.Add(5*time.Minute), baseTime, v2.Ascending, 0)
	assert.Error(t, err)
}

// Test a full cache hit scenario
func TestFullCacheHit(t *testing.T) {
	cache, fetcher, baseTime := setupTestCache(100)

	test := func(direction v2.FetchOrder) {
		// Initial fetch with specified direction
		start := baseTime
		end := baseTime.Add(5 * time.Minute)
		_, err := cache.Get(start, end, direction, 0)
		require.NoError(t, err)
		assert.Equal(t, 1, fetcher.getFetchCount())

		subStart := baseTime.Add(1 * time.Minute)
		subEnd := baseTime.Add(3 * time.Minute)

		// Sub range query ascending: full cache hit
		items, err := cache.Get(subStart, subEnd, v2.Ascending, 0)
		require.NoError(t, err)
		require.Len(t, items, 2)
		assert.Equal(t, 1, fetcher.getFetchCount())
		for i, item := range items {
			expectedTime := subStart.Add(time.Duration(i) * time.Minute)
			assert.Equal(t, expectedTime, item.ts)
		}
		// With limit
		items, err = cache.Get(subStart, subEnd, v2.Ascending, 1)
		require.NoError(t, err)
		require.Len(t, items, 1)
		assert.Equal(t, subStart, items[0].ts)

		// Sub range query descending: full cache hit
		items, err = cache.Get(subStart, subEnd, v2.Descending, 0)
		require.NoError(t, err)
		require.Len(t, items, 2)
		assert.Equal(t, 1, fetcher.getFetchCount())
		for i, item := range items {
			expectedTime := subStart.Add(time.Duration(1-i) * time.Minute)
			assert.Equal(t, expectedTime, item.ts)
		}
		// With limit
		items, err = cache.Get(subStart, subEnd, v2.Descending, 1)
		require.NoError(t, err)
		require.Len(t, items, 1)
		assert.Equal(t, subEnd.Add(-time.Minute), items[0].ts)

	}

	t.Run("ascending", func(t *testing.T) {
		test(v2.Ascending)
	})

	t.Run("descending", func(t *testing.T) {
		test(v2.Descending)
	})
}

// Test no overlap with newer range
func TestNoOverlap_NewerRange(t *testing.T) {
	testCases := []struct {
		name                string
		initialDirection    v2.FetchOrder
		newerRangeDirection v2.FetchOrder
		expectedFetchCounts []int // Expected fetch counts after each fetch
	}{
		{
			name:                "Ascending-Ascending",
			initialDirection:    v2.Ascending,
			newerRangeDirection: v2.Ascending,
			expectedFetchCounts: []int{1, 2, 3, 3},
		},
		{
			name:                "Ascending-Descending",
			initialDirection:    v2.Ascending,
			newerRangeDirection: v2.Descending,
			expectedFetchCounts: []int{1, 2, 3, 3},
		},
		{
			name:                "Descending-Ascending",
			initialDirection:    v2.Descending,
			newerRangeDirection: v2.Ascending,
			expectedFetchCounts: []int{1, 2, 3, 3},
		},
		{
			name:                "Descending-Descending",
			initialDirection:    v2.Descending,
			newerRangeDirection: v2.Descending,
			expectedFetchCounts: []int{1, 2, 3, 3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cache, fetcher, baseTime := setupTestCache(100)

			// Initial fetch
			start := baseTime
			end := baseTime.Add(5 * time.Minute)
			_, err := cache.Get(start, end, tc.initialDirection, 0)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedFetchCounts[0], fetcher.getFetchCount())

			// Query non-overlapping but newer range
			newStart := baseTime.Add(10 * time.Minute)
			newEnd := baseTime.Add(15 * time.Minute)
			items, err := cache.Get(newStart, newEnd, tc.newerRangeDirection, 0)
			require.NoError(t, err)
			require.Len(t, items, 5)
			assert.Equal(t, tc.expectedFetchCounts[1], fetcher.getFetchCount())

			// The old cache was dropped
			_, err = cache.Get(start, end, tc.initialDirection, 0)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedFetchCounts[2], fetcher.getFetchCount())

			// Query the new range again - should hit the cache
			_, err = cache.Get(newStart, newEnd, tc.newerRangeDirection, 0)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedFetchCounts[3], fetcher.getFetchCount())
		})
	}
}

// Test no overlap with newer range with limit query param
func TestNoOverlap_NewerRange_WithQueryLimit(t *testing.T) {
	cache, fetcher, baseTime := setupTestCache(100)

	// Initial fetch
	start := baseTime
	end := baseTime.Add(5 * time.Minute)
	_, err := cache.Get(start, end, v2.Ascending, 0)
	require.NoError(t, err)
	assert.Equal(t, int(1), fetcher.getFetchCount())

	// Query non-overlapping but newer range
	// With limit = 2, it'll just fetch 10:00, 11:00
	newStart := baseTime.Add(10 * time.Minute)
	newEnd := baseTime.Add(15 * time.Minute)
	items, err := cache.Get(newStart, newEnd, v2.Ascending, 2)
	require.NoError(t, err)
	require.Len(t, items, 2)
	assert.Equal(t, 2, fetcher.getFetchCount())

	// The old cache was dropped
	_, err = cache.Get(start, end, v2.Ascending, 0)
	require.NoError(t, err)
	assert.Equal(t, 3, fetcher.getFetchCount())

	// Query [10:00, 11:00+1ns) should have full cache hit
	_, err = cache.Get(newStart, newStart.Add(time.Minute).Add(time.Nanosecond), v2.Ascending, 0)
	require.NoError(t, err)
	assert.Equal(t, 3, fetcher.getFetchCount())

	// Query the new range again - should fetch DB
	_, err = cache.Get(newStart, newEnd, v2.Ascending, 0)
	require.NoError(t, err)
	assert.Equal(t, 4, fetcher.getFetchCount())
}

// Test no overlap with older range
func TestNoOverlap_OlderRange(t *testing.T) {
	testCases := []struct {
		name                string
		initialDirection    v2.FetchOrder
		olderRangeDirection v2.FetchOrder
		expectedFetchCounts []int
	}{
		{
			name:                "Ascending-Ascending",
			initialDirection:    v2.Ascending,
			olderRangeDirection: v2.Ascending,
			expectedFetchCounts: []int{1, 2, 2},
		},
		{
			name:                "Ascending-Descending",
			initialDirection:    v2.Ascending,
			olderRangeDirection: v2.Descending,
			expectedFetchCounts: []int{1, 2, 2},
		},
		{
			name:                "Descending-Ascending",
			initialDirection:    v2.Descending,
			olderRangeDirection: v2.Ascending,
			expectedFetchCounts: []int{1, 2, 2},
		},
		{
			name:                "Descending-Descending",
			initialDirection:    v2.Descending,
			olderRangeDirection: v2.Descending,
			expectedFetchCounts: []int{1, 2, 2},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cache, fetcher, baseTime := setupTestCache(100)

			// Initial fetch
			start := baseTime.Add(5 * time.Minute)
			end := baseTime.Add(10 * time.Minute)
			items, err := cache.Get(start, end, tc.initialDirection, 0)
			require.NoError(t, err)
			require.Len(t, items, 5)
			assert.Equal(t, tc.expectedFetchCounts[0], fetcher.getFetchCount())

			// Query older range
			oldStart := baseTime
			oldEnd := baseTime.Add(3 * time.Minute)
			items, err = cache.Get(oldStart, oldEnd, tc.olderRangeDirection, 0)
			require.NoError(t, err)
			require.Len(t, items, 3)
			assert.Equal(t, tc.expectedFetchCounts[1], fetcher.getFetchCount())

			// Query the new range again - should hit the cache
			for limit := 0; limit <= 5; limit++ {
				_, err = cache.Get(start, end, v2.Ascending, limit)
				require.NoError(t, err)
				assert.Equal(t, tc.expectedFetchCounts[2], fetcher.getFetchCount())
				_, err = cache.Get(start, end, v2.Descending, limit)
				require.NoError(t, err)
				assert.Equal(t, tc.expectedFetchCounts[2], fetcher.getFetchCount())
			}
		})
	}
}

// Test with limit parameter
func TestWithLimit(t *testing.T) {
	cache, fetcher, baseTime := setupTestCache(100)

	// Fetch with limit
	start := baseTime
	end := baseTime.Add(10 * time.Minute)
	limit := 3

	// Resulting in cache [0:00, 2:00+ns)
	items, err := cache.Get(start, end, v2.Ascending, limit)
	require.NoError(t, err)
	require.Len(t, items, limit)
	for i, item := range items {
		expectedTime := start.Add(time.Duration(i) * time.Minute)
		assert.Equal(t, expectedTime, item.ts)
	}
	assert.Equal(t, 1, fetcher.getFetchCount())

	// Full cache hit
	_, err = cache.Get(start, start.Add(2*time.Minute).Add(time.Nanosecond), v2.Ascending, limit)
	require.NoError(t, err)
	assert.Equal(t, 1, fetcher.getFetchCount())

	// [0:00, 3:00) with limit=3 should also have full cache, because there are already 3 items in
	// the cache, so it won't do more fetches for [2:00+ns, 3:00).
	_, err = cache.Get(start, start.Add(3*time.Minute), v2.Ascending, limit)
	require.NoError(t, err)
	assert.Equal(t, 1, fetcher.getFetchCount())
	// However, with descending, it will have to fetch [2:00+ns, 3:00) first (instead of using cache),
	// so this will cause an increase in fetch count.
	_, err = cache.Get(start, start.Add(3*time.Minute), v2.Descending, limit)
	require.NoError(t, err)
	assert.Equal(t, 2, fetcher.getFetchCount())

	// Fetch with descending order and limit
	// Resulting in cache [7:00, 10:00)
	items, err = cache.Get(start, end, v2.Descending, limit)
	require.NoError(t, err)
	require.Len(t, items, limit)
	for i, item := range items {
		expectedTime := end.Add(-time.Minute - time.Duration(i)*time.Minute)
		assert.Equal(t, expectedTime, item.ts)
	}
	assert.Equal(t, 3, fetcher.getFetchCount())

	// Old cache dropped
	// And this result won't be cached (remain as [7:00, 10:00)) as it's strictly older than
	// what's in cache
	_, err = cache.Get(start, start.Add(3*time.Minute), v2.Ascending, limit)
	require.NoError(t, err)
	assert.Equal(t, 4, fetcher.getFetchCount())

	// Full hit new cache
	_, err = cache.Get(start.Add(7*time.Minute), end, v2.Ascending, limit)
	require.NoError(t, err)
	assert.Equal(t, 4, fetcher.getFetchCount())
}

// Test partial overlap with newer range
func TestPartialOverlap_NewerRange(t *testing.T) {
	testCases := []struct {
		name                string
		initialDirection    v2.FetchOrder
		overlapDirection    v2.FetchOrder
		subRangeDirection   v2.FetchOrder
		expectedFetchCounts []int
	}{
		{
			name:                "Ascending-Ascending-Ascending",
			initialDirection:    v2.Ascending,
			overlapDirection:    v2.Ascending,
			subRangeDirection:   v2.Ascending,
			expectedFetchCounts: []int{1, 2, 2},
		},
		{
			name:                "Ascending-Descending-Ascending",
			initialDirection:    v2.Ascending,
			overlapDirection:    v2.Descending,
			subRangeDirection:   v2.Ascending,
			expectedFetchCounts: []int{1, 2, 2},
		},
		{
			name:                "Descending-Ascending-Descending",
			initialDirection:    v2.Descending,
			overlapDirection:    v2.Ascending,
			subRangeDirection:   v2.Descending,
			expectedFetchCounts: []int{1, 2, 2},
		},
		{
			name:                "Descending-Descending-Descending",
			initialDirection:    v2.Descending,
			overlapDirection:    v2.Descending,
			subRangeDirection:   v2.Descending,
			expectedFetchCounts: []int{1, 2, 2},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cache, fetcher, baseTime := setupTestCache(100)

			// Initial fetch [0:00, 5:00)
			start := baseTime
			end := baseTime.Add(5 * time.Minute)
			_, err := cache.Get(start, end, tc.initialDirection, 0)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedFetchCounts[0], fetcher.getFetchCount())

			// Query overlapping range [3:00, 8:00)
			newStart := baseTime.Add(3 * time.Minute)
			newEnd := baseTime.Add(8 * time.Minute)

			items, err := cache.Get(newStart, newEnd, tc.overlapDirection, 0)
			require.NoError(t, err)
			require.Len(t, items, 5)
			assert.Equal(t, tc.expectedFetchCounts[1], fetcher.getFetchCount())

			// Verify items are correct and in order
			for i, item := range items {
				expectedTime := newStart.Add(time.Duration(i) * time.Minute)
				if tc.overlapDirection == v2.Descending {
					expectedTime = newEnd.Add(time.Duration(-1*i-1) * time.Minute)
				}
				assert.Equal(t, expectedTime, item.ts)
			}

			// Query within the extended range - should be a cache hit
			subStart := baseTime
			subEnd := baseTime.Add(8 * time.Minute)
			_, err = cache.Get(subStart, subEnd, tc.subRangeDirection, 0)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedFetchCounts[2], fetcher.getFetchCount())
		})
	}
}

// Test partial overlap with newer range with limit query param
func TestPartialOverlap_NewerRange_WithQueryLimit(t *testing.T) {
	t.Run("newer-range ascending query extends cache", func(t *testing.T) {
		cache, fetcher, baseTime := setupTestCache(100)

		// Initial fetch [0:00, 5:00)
		start := baseTime
		end := baseTime.Add(5 * time.Minute)
		_, err := cache.Get(start, end, v2.Ascending, 0)
		require.NoError(t, err)
		assert.Equal(t, 1, fetcher.getFetchCount())

		// Query overlapping range [3:00, 8:00)
		// With limit=4, it'll cut off at 6:00 (the cache end set to +1ns)
		newStart := baseTime.Add(3 * time.Minute)
		newEnd := baseTime.Add(8 * time.Minute)
		items, err := cache.Get(newStart, newEnd, v2.Ascending, 4)
		require.NoError(t, err)
		require.Len(t, items, 4)
		for i, item := range items {
			assert.Equal(t, newStart.Add(time.Duration(i)*time.Minute), item.ts)
		}
		assert.Equal(t, 2, fetcher.getFetchCount())

		// Querying [0:00, 6:00) will have full cache hit
		_, err = cache.Get(baseTime, baseTime.Add(6*time.Minute), v2.Ascending, 0)
		require.NoError(t, err)
		assert.Equal(t, 2, fetcher.getFetchCount())

		// Querying [0:00, 8:00) will have to fetch DB
		_, err = cache.Get(start, newEnd, v2.Ascending, 0)
		require.NoError(t, err)
		assert.Equal(t, 3, fetcher.getFetchCount())
	})

	t.Run("newer-range ascending query has full cache hit", func(t *testing.T) {
		cache, fetcher, baseTime := setupTestCache(100)

		// Initial fetch [0:00, 5:00)
		start := baseTime
		end := baseTime.Add(5 * time.Minute)
		_, err := cache.Get(start, end, v2.Ascending, 0)
		require.NoError(t, err)
		assert.Equal(t, 1, fetcher.getFetchCount())

		// Query overlapping range [3:00, 8:00)
		// With limit=2, it'll cut off at 4:00, the query can be served out of the cache
		// and there is no DB fetch needed
		newStart := baseTime.Add(3 * time.Minute)
		newEnd := baseTime.Add(8 * time.Minute)
		items, err := cache.Get(newStart, newEnd, v2.Ascending, 2)
		require.NoError(t, err)
		require.Len(t, items, 2)
		for i, item := range items {
			assert.Equal(t, newStart.Add(time.Duration(i)*time.Minute), item.ts)
		}
		assert.Equal(t, 1, fetcher.getFetchCount())

		// Querying [0:00, 5:00) will have full cache hit
		_, err = cache.Get(start, end, v2.Ascending, 0)
		require.NoError(t, err)
		assert.Equal(t, 1, fetcher.getFetchCount())

		// Querying [0:00, 6:00) will have to fetch DB
		_, err = cache.Get(start, end.Add(time.Minute), v2.Ascending, 0)
		require.NoError(t, err)
		assert.Equal(t, 2, fetcher.getFetchCount())
	})

	t.Run("newer-range descending query replaces cache", func(t *testing.T) {
		cache, fetcher, baseTime := setupTestCache(100)

		// Initial fetch [0:00, 5:00)
		start := baseTime
		end := baseTime.Add(5 * time.Minute)
		_, err := cache.Get(start, end, v2.Ascending, 0)
		require.NoError(t, err)
		assert.Equal(t, 1, fetcher.getFetchCount())

		// Query overlapping range [3:00, 8:00), but with descending so it'll fetch the
		// high-end of items [6:00, 8:00) in the range
		newStart := baseTime.Add(3 * time.Minute)
		newEnd := baseTime.Add(8 * time.Minute)
		items, err := cache.Get(newStart, newEnd, v2.Descending, 2)
		require.NoError(t, err)
		require.Len(t, items, 2)
		for i, item := range items {
			assert.Equal(t, newEnd.Add(time.Duration(-1*i-1)*time.Minute), item.ts)
		}
		assert.Equal(t, 2, fetcher.getFetchCount())

		// Querying [6:00, 8:00) will have full cache hit
		_, err = cache.Get(baseTime.Add(6*time.Minute), baseTime.Add(8*time.Minute), v2.Ascending, 0)
		require.NoError(t, err)
		assert.Equal(t, 2, fetcher.getFetchCount())

		// Querying [0:00, 5:00) will have to fetch DB
		_, err = cache.Get(start, end, v2.Ascending, 0)
		require.NoError(t, err)
		assert.Equal(t, 3, fetcher.getFetchCount())
	})

	t.Run("newer-range query causes cache eviction", func(t *testing.T) {
		cache, fetcher, baseTime := setupTestCache(6)

		// Initial fetch [0:00, 5:00)
		start := baseTime
		end := baseTime.Add(5 * time.Minute)
		_, err := cache.Get(start, end, v2.Ascending, 0)
		require.NoError(t, err)
		assert.Equal(t, 1, fetcher.getFetchCount())

		// Query overlapping range [3:00, 8:00)
		// This will find 4 items [4:00, 8:00), which is connected to cache and can extend it
		newStart := baseTime.Add(3 * time.Minute)
		newEnd := baseTime.Add(8 * time.Minute)
		items, err := cache.Get(newStart, newEnd, v2.Descending, 4)
		require.NoError(t, err)
		require.Len(t, items, 4)
		for i, item := range items {
			assert.Equal(t, newEnd.Add(time.Duration(-1*i-1)*time.Minute), item.ts)
		}
		assert.Equal(t, 2, fetcher.getFetchCount())

		// Querying [2:00, 8:00) will have full cache hit
		_, err = cache.Get(baseTime.Add(2*time.Minute), baseTime.Add(8*time.Minute), v2.Ascending, 0)
		require.NoError(t, err)
		assert.Equal(t, 2, fetcher.getFetchCount())

		// Querying [1:00, 8:00) will have to fetch DB (the cache range is [2:00, 8:00))
		_, err = cache.Get(baseTime.Add(1*time.Minute), baseTime.Add(8*time.Minute), v2.Ascending, 0)
		require.NoError(t, err)
		assert.Equal(t, 3, fetcher.getFetchCount())
	})
}

// Test partial overlap with older range
func TestPartialOverlap_OlderRange(t *testing.T) {
	t.Run("older-range descending query extends cache", func(t *testing.T) {
		cache, fetcher, baseTime := setupTestCache(100)

		// Initial fetch [5:00, 10:00)
		start := baseTime.Add(5 * time.Minute)
		end := baseTime.Add(10 * time.Minute)
		_, err := cache.Get(start, end, v2.Descending, 0)
		require.NoError(t, err)
		assert.Equal(t, 1, fetcher.getFetchCount())

		// Query older overlapping range [3:00, 8:00) in descending order
		// With limit=4, it'l l cut off at 4:00 (the cache end set to +1ns)
		// This results in cache [4:00, 10:00)
		newStart := baseTime.Add(3 * time.Minute)
		newEnd := baseTime.Add(8 * time.Minute)
		items, err := cache.Get(newStart, newEnd, v2.Descending, 4)
		require.NoError(t, err)
		require.Len(t, items, 4)
		for i, item := range items {
			assert.Equal(t, newEnd.Add(time.Duration(-i-1)*time.Minute), item.ts)
		}
		assert.Equal(t, 2, fetcher.getFetchCount())

		// Querying [4:00, 10:00) will have full cache hit
		_, err = cache.Get(baseTime.Add(4*time.Minute), end, v2.Ascending, 0)
		require.NoError(t, err)
		assert.Equal(t, 2, fetcher.getFetchCount())

		// Querying [0:00, 8:00) will have to fetch DB
		_, err = cache.Get(baseTime, newEnd, v2.Ascending, 0)
		require.NoError(t, err)
		assert.Equal(t, 3, fetcher.getFetchCount())
	})

	t.Run("older-range descending query has full cache hit", func(t *testing.T) {
		cache, fetcher, baseTime := setupTestCache(100)

		// Initial fetch [5:00, 10:00)
		start := baseTime.Add(5 * time.Minute)
		end := baseTime.Add(10 * time.Minute)
		_, err := cache.Get(start, end, v2.Descending, 0)
		require.NoError(t, err)
		assert.Equal(t, 1, fetcher.getFetchCount())

		// Query overlapping range [3:00, 8:00)
		// With limit=2, it'll just fetch 7:00 and 6:00, which are cached
		// So the cache remains as [5:00, 10:00)
		newStart := baseTime.Add(3 * time.Minute)
		newEnd := baseTime.Add(8 * time.Minute)
		items, err := cache.Get(newStart, newEnd, v2.Descending, 2)
		require.NoError(t, err)
		require.Len(t, items, 2)
		for i, item := range items {
			assert.Equal(t, newEnd.Add(time.Duration(-i-1)*time.Minute), item.ts)
		}
		assert.Equal(t, 1, fetcher.getFetchCount())

		// Querying [5:00, 10:00) will have full cache hit
		_, err = cache.Get(start, end, v2.Ascending, 0)
		require.NoError(t, err)
		assert.Equal(t, 1, fetcher.getFetchCount())

		// Querying [3:00, 8:00) will have to fetch DB
		_, err = cache.Get(newStart, newEnd, v2.Ascending, 0)
		require.NoError(t, err)
		assert.Equal(t, 2, fetcher.getFetchCount())
	})

	t.Run("older-range ascending query has no effect on cache", func(t *testing.T) {
		cache, fetcher, baseTime := setupTestCache(100)

		// Initial fetch [5:00, 10:00)
		start := baseTime.Add(5 * time.Minute)
		end := baseTime.Add(10 * time.Minute)
		_, err := cache.Get(start, end, v2.Descending, 0)
		require.NoError(t, err)
		assert.Equal(t, 1, fetcher.getFetchCount())

		// Query overlapping range [3:00, 8:00)
		// With limit=2, it'll just fetch 3:00 and 4:00, which are disjoint with cache
		// so has no effect
		newStart := baseTime.Add(3 * time.Minute)
		newEnd := baseTime.Add(8 * time.Minute)
		items, err := cache.Get(newStart, newEnd, v2.Ascending, 2)
		require.NoError(t, err)
		require.Len(t, items, 2)
		for i, item := range items {
			assert.Equal(t, newStart.Add(time.Duration(i)*time.Minute), item.ts)
		}
		assert.Equal(t, 2, fetcher.getFetchCount())

		// Querying [5:00, 10:00) will have full cache hit
		_, err = cache.Get(start, end, v2.Ascending, 0)
		require.NoError(t, err)
		assert.Equal(t, 2, fetcher.getFetchCount())

		// Querying [3:00, 8:00) will have to fetch DB
		_, err = cache.Get(newStart, newEnd, v2.Ascending, 0)
		require.NoError(t, err)
		assert.Equal(t, 3, fetcher.getFetchCount())
	})

	t.Run("older-range query causes cache eviction", func(t *testing.T) {
		cache, fetcher, baseTime := setupTestCache(6)

		// Initial fetch [5:00, 10:00)
		start := baseTime.Add(5 * time.Minute)
		end := baseTime.Add(10 * time.Minute)
		_, err := cache.Get(start, end, v2.Descending, 0)
		require.NoError(t, err)
		assert.Equal(t, 1, fetcher.getFetchCount())

		// Query overlapping range [3:00, 8:00)
		// This could have created cache [3:00, 10:00), but with eviction it'll be [4:00, 10:00)
		newStart := baseTime.Add(3 * time.Minute)
		newEnd := baseTime.Add(8 * time.Minute)
		items, err := cache.Get(newStart, newEnd, v2.Ascending, 3)
		require.NoError(t, err)
		require.Len(t, items, 3)
		for i, item := range items {
			assert.Equal(t, newStart.Add(time.Duration(i)*time.Minute), item.ts)
		}
		assert.Equal(t, 2, fetcher.getFetchCount())

		// Querying [4:00, 10:00) will have full cache hit
		_, err = cache.Get(baseTime.Add(4*time.Minute), end, v2.Ascending, 0)
		require.NoError(t, err)
		assert.Equal(t, 2, fetcher.getFetchCount())

		// Querying [3:00, 8:00) will have to fetch DB
		_, err = cache.Get(newStart, newEnd, v2.Ascending, 0)
		require.NoError(t, err)
		assert.Equal(t, 3, fetcher.getFetchCount())
	})
}

// Test partial overlap with both newer and older range with limit query param
func TestPartialOverlap_NewerAndOlderRange_WithQueryLimit(t *testing.T) {
	t.Run("ascending query has no effect on cache", func(t *testing.T) {
		cache, fetcher, baseTime := setupTestCache(100)

		// Initial fetch [5:00, 10:00)
		start := baseTime.Add(5 * time.Minute)
		end := baseTime.Add(10 * time.Minute)
		items, err := cache.Get(start, end, v2.Ascending, 0)
		require.NoError(t, err)
		require.Len(t, items, 5)
		assert.Equal(t, 1, fetcher.getFetchCount())

		// Query a larger range [3:00, 12:00)
		// With limit=2, it will not hit any data in cache [5:00, 10:00)
		newStart := baseTime.Add(3 * time.Minute)
		newEnd := baseTime.Add(12 * time.Minute)
		items, err = cache.Get(newStart, newEnd, v2.Ascending, 2)
		require.NoError(t, err)
		require.Len(t, items, 2)
		for i, item := range items {
			assert.Equal(t, newStart.Add(time.Duration(i)*time.Minute), item.ts)
		}
		assert.Equal(t, 2, fetcher.getFetchCount())

		// Querying [5:00, 10:00) will hit full cache
		items, err = cache.Get(start, end, v2.Descending, 0)
		require.NoError(t, err)
		require.Len(t, items, 5)
		assert.Equal(t, 2, fetcher.getFetchCount())

		items, err = cache.Get(newStart, newEnd, v2.Ascending, 2)
		require.NoError(t, err)
		require.Len(t, items, 2)
		assert.Equal(t, 3, fetcher.getFetchCount())
	})

	t.Run("ascending query extends cache", func(t *testing.T) {
		cache, fetcher, baseTime := setupTestCache(100)

		// Initial fetch [5:00, 10:00)
		start := baseTime.Add(5 * time.Minute)
		end := baseTime.Add(10 * time.Minute)
		items, err := cache.Get(start, end, v2.Ascending, 0)
		require.NoError(t, err)
		require.Len(t, items, 5)
		assert.Equal(t, 1, fetcher.getFetchCount())

		// Query a larger range [3:00, 12:00)
		// With limit=3, it will exhaust [3:00, 5:00) so the results are connected to cache [5:00, 10:00)
		// The resulting cache is [3:00, 10:00)
		newStart := baseTime.Add(3 * time.Minute)
		newEnd := baseTime.Add(12 * time.Minute)
		items, err = cache.Get(newStart, newEnd, v2.Ascending, 3)
		require.NoError(t, err)
		require.Len(t, items, 3)
		for i, item := range items {
			assert.Equal(t, newStart.Add(time.Duration(i)*time.Minute), item.ts)
		}
		assert.Equal(t, 2, fetcher.getFetchCount())

		// Querying [3:00, 10:00) will hit full cache
		items, err = cache.Get(newStart, end, v2.Descending, 0)
		require.NoError(t, err)
		require.Len(t, items, 7)
		assert.Equal(t, 2, fetcher.getFetchCount())

		// Query a larger range [3:00, 12:00)
		// With limit=8, this will cover from 3:00 to  10:00
		items, err = cache.Get(newStart, newEnd, v2.Ascending, 8)
		require.NoError(t, err)
		require.Len(t, items, 8)
		assert.Equal(t, 3, fetcher.getFetchCount())

		// Querying [3:00, 10:00+1ns) will have full cache
		items, err = cache.Get(newStart, baseTime.Add(10*time.Minute).Add(time.Nanosecond), v2.Descending, 0)
		require.NoError(t, err)
		require.Len(t, items, 8)
		assert.Equal(t, 3, fetcher.getFetchCount())
	})

	t.Run("descending query replaces cache", func(t *testing.T) {
		cache, fetcher, baseTime := setupTestCache(100)

		// Initial fetch [5:00, 10:00)
		start := baseTime.Add(5 * time.Minute)
		end := baseTime.Add(10 * time.Minute)
		items, err := cache.Get(start, end, v2.Ascending, 0)
		require.NoError(t, err)
		require.Len(t, items, 5)
		assert.Equal(t, 1, fetcher.getFetchCount())

		// Query a larger range [3:00, 12:00)
		// With limit=2, it'll return 11:00 and 10:00
		// Mathematically this is connected to [5:00, 10:00), but the FeedCache cannot decide
		// as it will get 2 items from [10:00, 12:00) as it asks for 2 -- it may assume there
		// are actually more than 2
		// The resulting cache is [10:00, 12:00)
		newStart := baseTime.Add(3 * time.Minute)
		newEnd := baseTime.Add(12 * time.Minute)
		items, err = cache.Get(newStart, newEnd, v2.Descending, 2)
		require.NoError(t, err)
		require.Len(t, items, 2)
		for i, item := range items {
			assert.Equal(t, newEnd.Add(time.Duration(-i-1)*time.Minute), item.ts)
		}
		assert.Equal(t, 2, fetcher.getFetchCount())

		// Querying [10:00, 12:00) will hit full cache
		items, err = cache.Get(end, newEnd, v2.Descending, 0)
		require.NoError(t, err)
		require.Len(t, items, 2)
		assert.Equal(t, 2, fetcher.getFetchCount())

		// Querying [3:00, 12:00) again with limit=2, should have full cache
		items, err = cache.Get(newStart, newEnd, v2.Descending, 2)
		require.NoError(t, err)
		require.Len(t, items, 2)
		assert.Equal(t, 2, fetcher.getFetchCount())

		// Querying [3:00, 12:00) again without limit will have to fetch DB
		items, err = cache.Get(newStart, newEnd, v2.Ascending, 0)
		require.NoError(t, err)
		require.Len(t, items, 9)
		assert.Equal(t, 3, fetcher.getFetchCount())
	})

	t.Run("descending query extends cache", func(t *testing.T) {
		cache, fetcher, baseTime := setupTestCache(100)

		// Initial fetch [5:00, 10:00)
		start := baseTime.Add(5 * time.Minute)
		end := baseTime.Add(10 * time.Minute)
		items, err := cache.Get(start, end, v2.Ascending, 0)
		require.NoError(t, err)
		require.Len(t, items, 5)
		assert.Equal(t, 1, fetcher.getFetchCount())

		// Query a larger range [3:00, 12:00)
		// With limit=3, it'll return 11:00, 10:00, and 9:00, which are connnected to existing
		// cache [5:00, 10:00)
		// The resulting cache is [5:00, 12:00)
		newStart := baseTime.Add(3 * time.Minute)
		newEnd := baseTime.Add(12 * time.Minute)
		items, err = cache.Get(newStart, newEnd, v2.Descending, 3)
		require.NoError(t, err)
		require.Len(t, items, 3)
		for i, item := range items {
			assert.Equal(t, newEnd.Add(time.Duration(-i-1)*time.Minute), item.ts)
		}
		assert.Equal(t, 2, fetcher.getFetchCount())

		// Querying [5:00, 12:00) will hit full cache
		items, err = cache.Get(start, newEnd, v2.Descending, 0)
		require.NoError(t, err)
		require.Len(t, items, 7)
		assert.Equal(t, 2, fetcher.getFetchCount())

		// With limit=8, it will retrieve backward up to 4:00
		// Resulting cache [4:00, 12:00)
		items, err = cache.Get(newStart, newEnd, v2.Descending, 8)
		require.NoError(t, err)
		require.Len(t, items, 8)
		for i, item := range items {
			assert.Equal(t, newEnd.Add(time.Duration(-i-1)*time.Minute), item.ts)
		}
		assert.Equal(t, 3, fetcher.getFetchCount())

		// Querying [4:00, 12:00) will hit full cache
		items, err = cache.Get(baseTime.Add(4*time.Minute), newEnd, v2.Descending, 0)
		require.NoError(t, err)
		require.Len(t, items, 8)
		assert.Equal(t, 3, fetcher.getFetchCount())
	})

	t.Run("cache eviction", func(t *testing.T) {
		cache, fetcher, baseTime := setupTestCache(6)

		// Initial fetch [5:00, 10:00)
		start := baseTime.Add(5 * time.Minute)
		end := baseTime.Add(10 * time.Minute)
		items, err := cache.Get(start, end, v2.Ascending, 0)
		require.NoError(t, err)
		require.Len(t, items, 5)
		assert.Equal(t, 1, fetcher.getFetchCount())

		// Query a larger range [3:00, 12:00)
		// This could have created cache [5:00, 12:00), but with eviction it'll be [6:00, 12:00)
		newStart := baseTime.Add(3 * time.Minute)
		newEnd := baseTime.Add(12 * time.Minute)
		items, err = cache.Get(newStart, newEnd, v2.Descending, 3)
		require.NoError(t, err)
		require.Len(t, items, 3)
		for i, item := range items {
			assert.Equal(t, newEnd.Add(time.Duration(-i-1)*time.Minute), item.ts)
		}
		assert.Equal(t, 2, fetcher.getFetchCount())

		// Querying [6:00, 12:00) will have full cache hit
		items, err = cache.Get(start.Add(time.Minute), newEnd, v2.Descending, 0)
		require.NoError(t, err)
		require.Len(t, items, 6)
		assert.Equal(t, 2, fetcher.getFetchCount())

		// Querying [5:00, 12:00) will fetch DB to cover 5:00
		items, err = cache.Get(start, newEnd, v2.Descending, 0)
		require.NoError(t, err)
		require.Len(t, items, 7)
		assert.Equal(t, 3, fetcher.getFetchCount())
	})

}

// Test partial overlap with both newer and older range
func TestPartialOverlap_NewerAndOlderRange(t *testing.T) {
	cache, fetcher, baseTime := setupTestCache(100)

	// Initial fetch [5:00, 10:00)
	start := baseTime.Add(5 * time.Minute)
	end := baseTime.Add(10 * time.Minute)
	items, err := cache.Get(start, end, v2.Ascending, 0)
	require.NoError(t, err)
	require.Len(t, items, 5)
	assert.Equal(t, 1, fetcher.getFetchCount())

	// Query a larger range [3:00, 12:00)
	extendedStart := baseTime.Add(3 * time.Minute)
	extendedEnd := baseTime.Add(12 * time.Minute)
	items, err = cache.Get(extendedStart, extendedEnd, v2.Ascending, 0)
	require.NoError(t, err)
	require.Len(t, items, 9)

	// Should have two more fetches (two gaps)
	assert.Equal(t, 3, fetcher.getFetchCount())

	// Verify items are correct and in order
	for i, item := range items {
		expectedTime := extendedStart.Add(time.Duration(i) * time.Minute)
		assert.Equal(t, expectedTime, item.ts)
	}

	// Query within the extended range - should be a cache hit
	_, err = cache.Get(extendedStart, extendedEnd, v2.Ascending, 0)
	require.NoError(t, err)
	assert.Equal(t, 3, fetcher.getFetchCount())
}

// Test cache eviction due to maxItems limit
func TestEviction(t *testing.T) {
	cache, fetcher, baseTime := setupTestCache(3)

	// Fetch 5 minutes worth of data
	start := baseTime
	end := baseTime.Add(5 * time.Minute)
	items, err := cache.Get(start, end, v2.Ascending, 0)
	require.NoError(t, err)
	require.Len(t, items, 5)
	assert.Equal(t, 1, fetcher.getFetchCount())

	// Query the full range again - should be a partial cache hit
	// Only the most recent 3 items should be in cache due to maxItems
	start2 := baseTime
	end2 := baseTime.Add(5 * time.Minute)
	items2, err := cache.Get(start2, end2, v2.Ascending, 0)
	require.NoError(t, err)
	require.Len(t, items2, 5)
	assert.Equal(t, 2, fetcher.getFetchCount()) // Need to fetch older items not in cache

	// Query just the most recent 3 items - should be a cache hit
	recentStart := baseTime.Add(2 * time.Minute)
	recentEnd := baseTime.Add(5 * time.Minute)
	items3, err := cache.Get(recentStart, recentEnd, v2.Ascending, 0)
	require.NoError(t, err)
	require.Len(t, items3, 3)
	assert.Equal(t, 2, fetcher.getFetchCount()) // No new fetch needed
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
			direction := v2.Ascending
			if offset%2 == 0 {
				direction = v2.Descending
			}
			items, err := cache.Get(start, end, direction, 0)
			assert.NoError(t, err)
			assert.NotNil(t, items)
		}(i)
	}

	wg.Wait()
}
