package v2

import (
	"errors"
	"math"
	"sync"
	"time"
)

// SortOrder defines the ordering of data returned by fetchFromDB
type SortOrder int

const (
	Ascending SortOrder = iota
	Descending
)

// TimeRange represents a time interval [start, end) where
// start is inclusive and end is exclusive
type TimeRange struct {
	Start time.Time
	End   time.Time
}

type CacheEntry[T any] struct {
	TimeRange
	Data []*T
}

type FeedCache[T any] struct {
	mu           sync.RWMutex
	segment      *CacheEntry[T]
	maxItems     int
	fetchFromDB  func(start, end time.Time, order SortOrder, limit int) ([]*T, error)
	getTimestamp func(*T) time.Time
}

func NewFeedCache[T any](
	maxItems int,
	fetchFn func(start, end time.Time, order SortOrder, limit int) ([]*T, error),
	timestampFn func(*T) time.Time,
) *FeedCache[T] {
	return &FeedCache[T]{
		maxItems:     maxItems,
		fetchFromDB:  fetchFn,
		getTimestamp: timestampFn,
	}
}

func (tr TimeRange) Overlaps(other TimeRange) bool {
	return tr.Start.Before(other.End) && other.Start.Before(tr.End)
}

// reverseOrder reverses the order of elements in a slice
func reverseOrder[T any](data []*T) []*T {
	result := make([]*T, len(data))
	for i, item := range data {
		result[len(data)-1-i] = item
	}
	return result
}

func (c *FeedCache[T]) Get(start, end time.Time, queryOrder SortOrder, limit int) ([]*T, error) {
	if !start.Before(end) {
		return nil, errors.New("the start must be before end")
	}

	queryRange := TimeRange{Start: start, End: end}

	c.mu.RLock()
	segment := c.segment
	c.mu.RUnlock()

	// Handle no cache or no overlap cases together
	if segment == nil || !queryRange.Overlaps(segment.TimeRange) {
		shouldReplaceCache := segment == nil || !segment.TimeRange.End.After(start)
		return c.handleCacheMiss(start, end, queryOrder, limit, shouldReplaceCache)
	}

	// Handle overlapping case
	if queryOrder == Ascending {
		return c.handleAscendingQuery(start, end, segment, limit)
	} else {
		return c.handleDescendingQuery(start, end, segment, limit)
	}
}

func (c *FeedCache[T]) handleCacheMiss(start, end time.Time, queryOrder SortOrder, limit int, shouldReplaceCache bool) ([]*T, error) {
	// Fetch directly with the requested order and limit
	data, err := c.fetchFromDB(start, end, queryOrder, limit)
	if err != nil {
		return nil, err
	}
	hasMore := len(data) == limit

	// Cache data if it's not empty
	if len(data) > 0 && shouldReplaceCache {
		// Normalize data to ascending order before caching
		dataToCache := data
		if queryOrder == Descending {
			dataToCache = reverseOrder(data)
		}

		var newStart, newEnd time.Time
		if queryOrder == Ascending {
			newStart = start
			newEnd = end
			if hasMore {
				newEnd = c.getTimestamp(dataToCache[len(dataToCache)-1]).Add(time.Nanosecond)
			}
		} else {
			newEnd = end
			newStart = start
			if hasMore {
				newStart = c.getTimestamp(dataToCache[0])
			}
		}

		// If data exceeds maxItems, keep only the most recent ones
		if len(dataToCache) > c.maxItems {
			dataToCache = dataToCache[len(dataToCache)-c.maxItems:] // Keep newest
			newStart = c.getTimestamp(dataToCache[0])
		}

		c.mu.Lock()
		c.segment = &CacheEntry[T]{
			TimeRange: TimeRange{
				Start: newStart,
				End:   newEnd,
			},
			Data: dataToCache,
		}
		c.mu.Unlock()
	}

	return data, nil
}

func (c *FeedCache[T]) handleAscendingQuery(start, end time.Time, segment *CacheEntry[T], limit int) ([]*T, error) {
	var beforeData, afterData []*T
	var beforeHasMore, afterHasMore bool
	var err error

	// For ascending query:
	// 1. beforeData can only be merged if hasMore is false
	// 2. afterData can always be merged (it's certainly connected)

	// Get cached data that's overlapping the query range
	cachedOverlap := c.filterDataInRange(segment.Data, start, end)

	// If no gaps, it's full cache hit, just return filtered cached data
	if !start.Before(segment.Start) && !end.After(segment.End) {
		// Apply limit to cached results
		if limit > 0 && len(cachedOverlap) > limit {
			cachedOverlap = cachedOverlap[:limit] // Take first 'limit' items for ascending
		}

		return cachedOverlap, nil
	}

	// Check if we need data before cached segment
	if start.Before(segment.Start) {
		beforeData, err = c.fetchFromDB(start, segment.Start, Ascending, limit)
		if err != nil {
			return nil, err
		}
		if limit > 0 {
			beforeHasMore = len(beforeData) == limit
		}
	}

	// Check if we need data after cached segment
	remaining := math.MaxInt
	if limit > 0 {
		remaining = limit - len(beforeData) - len(cachedOverlap)
	}
	if remaining > 0 && end.After(segment.End) {
		afterData, err = c.fetchFromDB(segment.End, end, Ascending, remaining)
		if err != nil {
			return nil, err
		}
		afterHasMore = len(afterData) == remaining
	}

	// Combine results: beforeData -> cachedOverlap -> afterData
	numToReturn := len(beforeData) + len(cachedOverlap) + len(afterData)
	if limit > 0 {
		numToReturn = min(numToReturn, limit)
	}
	result := make([]*T, 0, numToReturn)

	beforeItems := min(numToReturn, len(beforeData))
	result = append(result, beforeData[:beforeItems]...)

	if len(result) < numToReturn {
		overlapItems := min(numToReturn-len(result), len(cachedOverlap))
		result = append(result, cachedOverlap[:overlapItems]...)
	}

	if len(result) < numToReturn {
		afterItems := min(numToReturn-len(result), len(afterData))
		result = append(result, afterData[:afterItems]...)
	}

	var newCacheData []*T
	var newStart, newEnd time.Time

	if start.Before(segment.Start) && !beforeHasMore {
		newCacheData = append(beforeData, segment.Data...)
		newCacheData = append(newCacheData, afterData...)
		newStart = start
	} else {
		newCacheData = append(segment.Data, afterData...)
		newStart = segment.Start
	}

	// Ensure we don't exceed maxItems, prioritizing recent data
	if len(newCacheData) > c.maxItems {
		newCacheData = newCacheData[len(newCacheData)-c.maxItems:] // Keep newest
		newStart = c.getTimestamp(newCacheData[0])
	}
	newEnd = end
	if len(afterData) == 0 {
		newEnd = segment.End
	} else if afterHasMore {
		// If the query didn't exhaust the range, we can only be sure a more conservative endpoint
		newEnd = c.getTimestamp(newCacheData[len(newCacheData)-1]).Add(time.Nanosecond)
	}

	// Update cache
	c.mu.Lock()

	c.segment = &CacheEntry[T]{
		TimeRange: TimeRange{
			Start: newStart,
			End:   newEnd,
		},
		Data: newCacheData,
	}

	c.mu.Unlock()

	return result, nil
}

func (c *FeedCache[T]) handleDescendingQuery(start, end time.Time, segment *CacheEntry[T], limit int) ([]*T, error) {
	var beforeData, afterData []*T
	var beforeHasMore, afterHasMore bool
	var err error
	// beforeHasMore, afterHasMore := false, false

	// For descending query:
	// 1. beforeData can always be merged (it's certainly connected)
	// 2. afterData can only be merged if hasMore is false

	// Get cached data that's overlapping the query range
	cachedOverlap := c.filterDataInRange(segment.Data, start, end)

	// If no gaps, just return filtered cached data in descending order
	if !start.Before(segment.Start) && !end.After(segment.End) {
		// Apply limit to cached results
		if limit > 0 && len(cachedOverlap) > limit {
			cachedOverlap = cachedOverlap[len(cachedOverlap)-limit:] // Take last 'limit' items for descending
		}

		// Return in descending order
		return reverseOrder(cachedOverlap), nil
	}

	// Check if we need data after cached segment
	if end.After(segment.End) {
		afterData, err = c.fetchFromDB(segment.End, end, Descending, limit)
		if err != nil {
			return nil, err
		}
		if limit > 0 {
			afterHasMore = len(afterData) == limit
		}
	}

	// Check if we need data before cached segment
	remaining := math.MaxInt
	if limit > 0 {
		remaining = limit - len(beforeData) - len(cachedOverlap)
	}
	if remaining > 0 && start.Before(segment.Start) {
		beforeData, err = c.fetchFromDB(start, segment.Start, Descending, remaining)
		if err != nil {
			return nil, err
		}
		beforeHasMore = len(beforeData) == remaining
	}

	// Combine results: afterData -> cachedOverlap -> beforeData
	numToReturn := len(beforeData) + len(cachedOverlap) + len(afterData)
	if limit > 0 {
		numToReturn = min(numToReturn, limit)
	}
	result := make([]*T, 0, numToReturn)

	afterItems := min(numToReturn, len(afterData))
	result = append(result, afterData[:afterItems]...)

	if len(result) < numToReturn {
		overlapItems := min(numToReturn-len(result), len(cachedOverlap))
		result = append(result, reverseOrder(cachedOverlap)[:overlapItems]...)
	}

	if len(result) < numToReturn {
		beforeItems := min(numToReturn-len(result), len(beforeData))
		result = append(result, beforeData[:beforeItems]...)
	}

	// Normalize to ascending order for caching
	afterData = reverseOrder(afterData)
	// Normalize to ascending order for caching
	beforeData = reverseOrder(beforeData)

	var newCacheData []*T
	var newStart time.Time

	if afterHasMore {
		// There is a afterData segment which is not connected to the cache, we'll replace the cache
		// with this segment
		newCacheData = afterData
		newStart = c.getTimestamp(newCacheData[0])
	} else {
		newCacheData = append(beforeData, segment.Data...)
		newCacheData = append(newCacheData, afterData...)
		newStart = minTime(start, segment.Start)
		if len(beforeData) == 0 {
			newStart = segment.Start
		} else if beforeHasMore {
			newStart = c.getTimestamp(newCacheData[0])
		}
	}

	// Ensure we don't exceed maxItems, prioritizing recent data
	if len(newCacheData) > c.maxItems {
		newCacheData = newCacheData[len(newCacheData)-c.maxItems:] // Keep newest
		newStart = c.getTimestamp(newCacheData[0])
	}

	// Update cache
	c.mu.Lock()

	c.segment = &CacheEntry[T]{
		TimeRange: TimeRange{
			Start: newStart,
			End:   maxTime(end, segment.End),
		},
		Data: newCacheData,
	}

	c.mu.Unlock()

	return result, nil
}

func (c *FeedCache[T]) filterDataInRange(data []*T, start, end time.Time) []*T {
	result := make([]*T, 0)

	for _, item := range data {
		timestamp := c.getTimestamp(item)

		// Since end is exclusive, we break on >=
		if !timestamp.Before(end) {
			break // since data is in ascending order
		}

		// Since start is inclusive, we include >=
		if !timestamp.Before(start) {
			result = append(result, item)
		}
	}

	return result
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}
