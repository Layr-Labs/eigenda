package v2

import (
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
	fetchFromDB  func(start, end time.Time) ([]*T, error)
	getTimestamp func(*T) time.Time
	order        SortOrder
}

func NewFeedCache[T any](
	maxItems int,
	fetchFn func(start, end time.Time) ([]*T, error),
	timestampFn func(*T) time.Time,
	order SortOrder,
) *FeedCache[T] {
	return &FeedCache[T]{
		maxItems:     maxItems,
		fetchFromDB:  fetchFn,
		getTimestamp: timestampFn,
		order:        order,
	}
}

func (tr TimeRange) Overlaps(other TimeRange) bool {
	return tr.Start.Before(other.End) && other.Start.Before(tr.End)
}

// reverseOrder reverses the order of elements in a slice
func (c *FeedCache[T]) reverseOrder(data []*T) []*T {
	result := make([]*T, len(data))
	for i, item := range data {
		result[len(data)-1-i] = item
	}
	return result
}

func (c *FeedCache[T]) Get(start, end time.Time) ([]*T, error) {
	queryRange := TimeRange{Start: start, End: end}

	c.mu.RLock()
	segment := c.segment
	c.mu.RUnlock()

	// Handle no cache or no overlap cases together
	if segment == nil || !queryRange.Overlaps(segment.TimeRange) {
		data, err := c.fetchFromDB(start, end)
		if err != nil {
			return nil, err
		}

		// Only cache if this is newer than existing cache
		if segment == nil || queryRange.Start.After(segment.TimeRange.Start) {
			// Normalize data to ascending order before caching
			dataToCache := data
			if c.order == Descending {
				dataToCache = c.reverseOrder(data)
			}

			// If data exceeds maxItems, keep only the most recent ones
			if len(dataToCache) > c.maxItems {
				dataToCache = dataToCache[len(dataToCache)-c.maxItems:] // Keep newest
				start = c.getTimestamp(dataToCache[0])
			}

			c.mu.Lock()
			c.segment = &CacheEntry[T]{
				TimeRange: TimeRange{
					Start: start,
					End:   end,
				},
				Data: dataToCache,
			}
			c.mu.Unlock()
		}

		return data, nil
	}

	// Handle overlapping case
	var beforeData, afterData []*T
	var err error

	// Check if we need data before cached segment
	if start.Before(segment.Start) {
		beforeData, err = c.fetchFromDB(start, segment.Start)
		if err != nil {
			return nil, err
		}
		// Normalize to ascending order for caching
		if c.order == Descending {
			beforeData = c.reverseOrder(beforeData)
		}
	}

	// Check if we need data after cached segment
	if end.After(segment.End) {
		afterData, err = c.fetchFromDB(segment.End, end)
		if err != nil {
			return nil, err
		}
		// Normalize to ascending order for caching
		if c.order == Descending {
			afterData = c.reverseOrder(afterData)
		}
	}

	// If no gaps, just return filtered cached data in the requested order
	if len(beforeData) == 0 && len(afterData) == 0 {
		cachedResult := c.filterDataInRange(segment.Data, start, end)
		if c.order == Descending {
			return c.reverseOrder(cachedResult), nil
		}
		return cachedResult, nil
	}

	// Calculate total size for result
	totalSize := len(beforeData) + len(afterData)
	cachedInRange := c.countDataInRange(segment.Data, start, end)

	// Process data for cache (always in ascending order)
	newCache := make([]*T, 0, totalSize+len(c.segment.Data))
	if beforeData != nil {
		newCache = append(newCache, beforeData...)
	}
	newCache = append(newCache, segment.Data...)
	if afterData != nil {
		newCache = append(newCache, afterData...)
	}

	// Result for the user includes only data in the requested range
	// Result = beforeData + filterDataInrange(cached data) + afterData
	result := make([]*T, 0, totalSize+cachedInRange)
	if beforeData != nil {
		result = append(result, beforeData...)
	}
	result = append(result, c.filterDataInRange(segment.Data, start, end)...)
	if afterData != nil {
		result = append(result, afterData...)
	}

	// Update cache with extended range
	c.mu.Lock()
	// If newCache would exceed maxItems, trim oldest data
	if len(newCache) > c.maxItems {
		newCache = newCache[len(newCache)-c.maxItems:] // Keep newest
	}

	c.segment = &CacheEntry[T]{
		TimeRange: TimeRange{
			Start: c.getTimestamp(newCache[0]),
			End:   maxTime(segment.End, end),
		},
		Data: newCache,
	}
	c.mu.Unlock()

	// Return the result in the requested order
	if c.order == Descending {
		return c.reverseOrder(result), nil
	}
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

func (c *FeedCache[T]) countDataInRange(data []*T, start, end time.Time) int {
	count := 0

	for _, item := range data {
		timestamp := c.getTimestamp(item)

		if !timestamp.Before(end) {
			break // since data is in ascending order

		}

		if !timestamp.Before(start) {
			count++
		}
	}

	return count
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}
