package v2

import (
	"errors"
	"math"
	"sync"
	"time"
)

// SortOrder defines the ordering of data returned by fetchFromDB.
type SortOrder int

const (
	Ascending SortOrder = iota
	Descending
)

// TimeRange represents a time interval [start, end) where start is inclusive and end
// is exclusive.
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// CacheEntry describes a segment of results fetched via fetchFromDB for a time range.
type CacheEntry[T any] struct {
	TimeRange
	Data []*T
}

// FeedCache tracks the most recent segment of results fetched via fetchFromDB.
// If new results (as a segment for the time range of query) are connected to the existing
// cached segment, it'll extend the cache segment.
// If there are more than maxItems in cache, it'll evict the oldest items.
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

func reverseOrder[T any](data []*T) []*T {
	result := make([]*T, len(data))
	for i, item := range data {
		result[len(data)-1-i] = item
	}
	return result
}

type ExecutionPlan[T any] struct {
	// data is the cache hit: matching results that are in time range [before.End, after.Start)
	data   []*T
	before *TimeRange
	after  *TimeRange
}

type ExecutionResult[T any] struct {
	order         SortOrder
	before        *TimeRange
	after         *TimeRange
	beforeData    []*T
	afterData     []*T
	beforeHasMore bool
	afterHasMore  bool
	result        []*T
}

func (c *FeedCache[T]) updateCache(result *ExecutionResult[T]) {
	if result.before == nil && result.after == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if result.order == Ascending {
		c.updateCacheAscending(result)
	} else {
		c.updateCacheDescending(result)
	}
}

func (c *FeedCache[T]) updateCacheAscending(result *ExecutionResult[T]) {
	before, after := result.before, result.after
	beforeData, afterData := result.beforeData, result.afterData

	if c.segment == nil {
		if len(afterData) > 0 {
			c.segment = &CacheEntry[T]{
				TimeRange: TimeRange{
					Start: after.Start,
					End:   after.End,
				},
				Data: afterData,
			}
		}
		return
	}

	var newCacheData []*T
	var newStart, newEnd time.Time

	newCacheData = c.segment.Data
	newStart = c.segment.Start
	if result.before != nil {
		beforeEnd := before.End
		if result.beforeHasMore {
			beforeEnd = c.getTimestamp(beforeData[len(beforeData)-1]).Add(time.Nanosecond)
		}
		if !beforeEnd.Before(c.segment.Start) {
			split := len(beforeData) - 1
			for ; split >= 0; split-- {
				if c.getTimestamp(beforeData[split]).Before(c.segment.Start) {
					break
				}
			}
			if split >= 0 {
				newCacheData = append(beforeData[:split+1], c.segment.Data...)
				newStart = before.Start
			}
		}
	}

	if after != nil {
		if after.End.After(c.segment.End) {
			newEnd = after.End
			split := 0
			for ; split < len(afterData); split++ {
				if c.getTimestamp(afterData[split]).After(c.segment.End) {
					break
				}
			}
			if split < len(afterData) {
				newCacheData = append(newCacheData, afterData[split:]...)
			}
		}
	}

	// Ensure we don't exceed maxItems, prioritizing recent data
	if len(newCacheData) > c.maxItems {
		newCacheData = newCacheData[len(newCacheData)-c.maxItems:] // Keep newest
		newStart = c.getTimestamp(newCacheData[0])
	}

	c.segment = &CacheEntry[T]{
		TimeRange: TimeRange{
			Start: newStart,
			End:   newEnd,
		},
		Data: newCacheData,
	}
}

func (c *FeedCache[T]) updateCacheDescending(result *ExecutionResult[T]) {
	before, after := result.before, result.after
	beforeData, afterData := result.beforeData, result.afterData

	if c.segment == nil {
		if len(afterData) > 0 {
			c.segment = &CacheEntry[T]{
				TimeRange: TimeRange{
					Start: after.Start,
					End:   after.End,
				},
				Data: reverseOrder(afterData),
			}
		}
		return
	}

	var newCacheData []*T
	var newStart, newEnd time.Time

	// Normalize to ascending order for caching
	afterData = reverseOrder(afterData)
	beforeData = reverseOrder(beforeData)

	newCacheData = c.segment.Data
	newStart = c.segment.Start
	if before != nil {
		beforeStart := before.Start
		if result.beforeHasMore {
			beforeStart = c.getTimestamp(beforeData[0])
		}
		if beforeStart.Before(c.segment.Start) {
			newStart = beforeStart
			split := len(beforeData) - 1
			for ; split >= 0; split-- {
				if c.getTimestamp(beforeData[split]).Before(c.segment.Start) {
					break
				}
			}
			if split >= 0 {
				newCacheData = append(beforeData[:split+1], c.segment.Data...)
			}
		}
	}

	newEnd = c.segment.End
	if after != nil {
		afterStart := after.Start
		if result.afterHasMore {
			afterStart = c.getTimestamp(afterData[0])
		}

		// The afterData segment which is not connected to the cache, we'll replace the cache
		// with this afterData
		if afterStart.After(c.segment.End) {
			c.segment = &CacheEntry[T]{
				TimeRange: TimeRange{
					Start: afterStart,
					End:   after.End,
				},
				Data: afterData,
			}
			return
		}

		if after.End.After(c.segment.End) {
			newEnd = after.End
			split := 0
			for ; split < len(afterData); split++ {
				if !c.getTimestamp(afterData[split]).Before(c.segment.End) {
					break
				}
			}
			if split < len(afterData) {
				newCacheData = append(newCacheData, afterData[split:]...)
			}
		}
	}

	// Ensure we don't exceed maxItems, prioritizing recent data
	if len(newCacheData) > c.maxItems {
		newCacheData = newCacheData[len(newCacheData)-c.maxItems:] // Keep newest
		newStart = c.getTimestamp(newCacheData[0])
	}

	c.segment = &CacheEntry[T]{
		TimeRange: TimeRange{
			Start: newStart,
			End:   newEnd,
		},
		Data: newCacheData,
	}
}

func (c *FeedCache[T]) Get(start, end time.Time, queryOrder SortOrder, limit int) ([]*T, error) {
	if !start.Before(end) {
		return nil, errors.New("the start must be before end")
	}

	plan := c.makePlan(start, end, queryOrder, limit)

	var result *ExecutionResult[T]
	var err error
	if queryOrder == Ascending {
		result, err = c.executePlanAscending(plan, limit)
	} else {
		result, err = c.executePlanDescending(plan, limit)
	}
	if err != nil {
		return nil, err
	}

	// TODO(jianxiao): make this run async so the results can return asap; this needs to get
	// unit tests work.
	c.updateCache(result)

	return result.result, nil
}

func (c *FeedCache[T]) makePlan(start, end time.Time, queryOrder SortOrder, limit int) ExecutionPlan[T] {
	c.mu.RLock()
	defer c.mu.RUnlock()

	segment := c.segment
	queryRange := TimeRange{Start: start, End: end}

	// Handle no cache or no overlap cases together
	if segment == nil || !queryRange.Overlaps(segment.TimeRange) {
		return ExecutionPlan[T]{
			// The data=nil, so it doesn't matter we fill the `before` or `after`
			after: &queryRange,
		}
	}

	// Get cached data that's overlapping the query range
	cachedOverlap := c.getCacheOverlap(segment.Data, start, end)
	// Apply limit to cached results
	if limit > 0 && len(cachedOverlap) > limit {
		if queryOrder == Ascending {
			cachedOverlap = cachedOverlap[:limit] // Take first 'limit' items for ascending
		} else {
			cachedOverlap = cachedOverlap[len(cachedOverlap)-limit:] // Take last 'limit' items for descending
		}
	}

	// The query range is fully contained in cache, it's a full cache hit
	if !start.Before(segment.Start) && !end.After(segment.End) {
		return ExecutionPlan[T]{
			data: cachedOverlap,
		}
	}

	// The query range overlaps the cache segment
	var before, after *TimeRange
	if start.Before(segment.Start) {
		before = &TimeRange{
			Start: start,
			End:   segment.Start,
		}
	}
	if end.After(segment.End) {
		after = &TimeRange{
			Start: segment.End,
			End:   end,
		}
	}
	return ExecutionPlan[T]{
		before: before,
		data:   cachedOverlap,
		after:  after,
	}
}
func (c *FeedCache[T]) executePlanAscending(plan ExecutionPlan[T], limit int) (*ExecutionResult[T], error) {
	var beforeData, afterData []*T
	var beforeHasMore, afterHasMore bool
	var err error

	// Fetch data before cache segment if needed
	if plan.before != nil {
		beforeData, err = c.fetchFromDB(plan.before.Start, plan.before.End, Ascending, limit)
		if err != nil {
			return nil, err
		}
		if limit > 0 {
			beforeHasMore = len(beforeData) == limit
		}
	}

	// Fetch data after cache segment if needed
	if plan.after != nil {
		remaining := math.MaxInt
		if limit > 0 {
			remaining = limit - len(beforeData) - len(plan.data)
		}
		if remaining > 0 {
			afterData, err = c.fetchFromDB(plan.after.Start, plan.after.End, Ascending, remaining)
			if err != nil {
				return nil, err
			}
			afterHasMore = len(afterData) == remaining
		}
	}

	// Combine results: beforeData -> cachedOverlap -> afterData
	numToReturn := len(beforeData) + len(plan.data) + len(afterData)
	if limit > 0 {
		numToReturn = min(numToReturn, limit)
	}
	result := make([]*T, 0, numToReturn)

	beforeItems := min(numToReturn, len(beforeData))
	result = append(result, beforeData[:beforeItems]...)

	if len(result) < numToReturn {
		overlapItems := min(numToReturn-len(result), len(plan.data))
		result = append(result, plan.data[:overlapItems]...)
	}

	if len(result) < numToReturn {
		afterItems := min(numToReturn-len(result), len(afterData))
		result = append(result, afterData[:afterItems]...)
	}

	return &ExecutionResult[T]{
		order:         Ascending,
		before:        plan.before,
		after:         plan.after,
		beforeData:    beforeData,
		afterData:     afterData,
		beforeHasMore: beforeHasMore,
		afterHasMore:  afterHasMore,
		result:        result,
	}, nil
}

func (c *FeedCache[T]) executePlanDescending(plan ExecutionPlan[T], limit int) (*ExecutionResult[T], error) {
	var beforeData, afterData []*T
	var beforeHasMore, afterHasMore bool
	var err error

	// Fetch data after cache segment if needed
	if plan.after != nil {
		afterData, err = c.fetchFromDB(plan.after.Start, plan.after.End, Descending, limit)
		if err != nil {
			return nil, err
		}
		if limit > 0 {
			afterHasMore = len(afterData) == limit
		}
	}

	// Fetch data before cache segment if needed
	if plan.before != nil {
		remaining := math.MaxInt
		if limit > 0 {
			remaining = limit - len(beforeData) - len(plan.data)
		}
		if remaining > 0 {
			beforeData, err = c.fetchFromDB(plan.before.Start, plan.before.End, Descending, remaining)
			if err != nil {
				return nil, err
			}
			beforeHasMore = len(beforeData) == remaining
		}
	}

	// Combine results: afterData -> cachedOverlap -> beforeData
	numToReturn := len(beforeData) + len(plan.data) + len(afterData)
	if limit > 0 {
		numToReturn = min(numToReturn, limit)
	}
	result := make([]*T, 0, numToReturn)

	afterItems := min(numToReturn, len(afterData))
	result = append(result, afterData[:afterItems]...)

	if len(result) < numToReturn {
		overlapItems := min(numToReturn-len(result), len(plan.data))
		result = append(result, reverseOrder(plan.data)[:overlapItems]...)
	}

	if len(result) < numToReturn {
		beforeItems := min(numToReturn-len(result), len(beforeData))
		result = append(result, beforeData[:beforeItems]...)
	}

	return &ExecutionResult[T]{
		order:         Descending,
		before:        plan.before,
		after:         plan.after,
		beforeData:    beforeData,
		afterData:     afterData,
		beforeHasMore: beforeHasMore,
		afterHasMore:  afterHasMore,
		result:        result,
	}, nil
}

func (c *FeedCache[T]) getCacheOverlap(data []*T, start, end time.Time) []*T {
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
