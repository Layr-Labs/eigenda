package v2

import (
	"errors"
	"math"
	"sync"
	"time"
)

// FetchOrder defines the ordering of data returned by fetchFromDB.
type FetchOrder int

const (
	Ascending FetchOrder = iota
	Descending
)

// FeedCache tracks the most recent segment of results fetched via fetchFromDB.
// If new results (as a segment for the time range of query) are connected to the existing
// cached segment, it'll extend the cache segment.
// If there are more than maxItems in cache, it'll evict the oldest items.
type FeedCache[T any] struct {
	mu           sync.RWMutex
	segment      *cacheEntry[T]
	maxItems     int
	fetchFromDB  func(start, end time.Time, order FetchOrder, limit int) ([]*T, error)
	getTimestamp func(*T) time.Time
}

func NewFeedCache[T any](
	maxItems int,
	fetchFn func(start, end time.Time, order FetchOrder, limit int) ([]*T, error),
	timestampFn func(*T) time.Time,
) *FeedCache[T] {
	return &FeedCache[T]{
		maxItems:     maxItems,
		fetchFromDB:  fetchFn,
		getTimestamp: timestampFn,
	}
}

// timeRange represents a time interval [start, end) where start is inclusive and end
// is exclusive.
type timeRange struct {
	start time.Time
	end   time.Time
}

// cacheEntry describes a segment of results fetched via fetchFromDB for a time range.
//
// It has the following properties:
// - all data items that are in range [start, end) are in the segment
// - no data items that are outside the range [start, end) are included in the segment
//
// The data items are in ascending order by timestamp.
type cacheEntry[T any] struct {
	timeRange
	data []*T
}

// executionPlan describes the breakdown of a data fetch query [start, end) into sub ranges
// that hits cache and that need DB fetches.
type executionPlan[T any] struct {
	// data is the cache hit, ie. matching results that are in time range [before.end, after.start)
	cacheHit []*T
	before   *timeRange
	after    *timeRange
}

// executionResult describes execution result of a plan.
type executionResult[T any] struct {
	order  FetchOrder
	before *timeRange
	after  *timeRange

	// The DB fetch results corresponding to `before` and `after` range.
	beforeData []*T
	afterData  []*T

	// Whether there are more data items in `before` range or `after` range.
	// This may have false positive but will never have false negative (e.g. if it says
	// beforeHasMore=false, then it's guaranteed that there are no more data items)
	beforeHasMore bool
	afterHasMore  bool

	// The result for the data fetch query.
	result []*T
}

func (tr timeRange) overlaps(other timeRange) bool {
	return tr.start.Before(other.end) && other.start.Before(tr.end)
}

func (c *FeedCache[T]) Get(start, end time.Time, queryOrder FetchOrder, limit int) ([]*T, error) {
	if !start.Before(end) {
		return nil, errors.New("the start must be before end")
	}

	plan := c.makePlan(start, end, queryOrder, limit)

	var result *executionResult[T]
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

func (c *FeedCache[T]) makePlan(start, end time.Time, queryOrder FetchOrder, limit int) executionPlan[T] {
	c.mu.RLock()
	defer c.mu.RUnlock()

	segment := c.segment
	queryRange := timeRange{start: start, end: end}

	// Handle no cache or no overlap cases together
	if segment == nil || !queryRange.overlaps(segment.timeRange) {
		return executionPlan[T]{
			// The data=nil, so it doesn't matter we fill the `before` or `after`
			after: &queryRange,
		}
	}

	// Get cached data that's overlapping the query range
	cachedOverlap := c.getCacheOverlap(segment.data, start, end)
	// Apply limit to cached results
	if limit > 0 && len(cachedOverlap) > limit {
		if queryOrder == Ascending {
			cachedOverlap = cachedOverlap[:limit] // Take first 'limit' items for ascending
		} else {
			cachedOverlap = cachedOverlap[len(cachedOverlap)-limit:] // Take last 'limit' items for descending
		}
	}

	// The query range is fully contained in cache, it's a full cache hit
	if !start.Before(segment.start) && !end.After(segment.end) {
		return executionPlan[T]{
			cacheHit: cachedOverlap,
		}
	}

	// The query range overlaps the cache segment
	var before, after *timeRange
	if start.Before(segment.start) {
		before = &timeRange{
			start: start,
			end:   segment.start,
		}
	}
	if end.After(segment.end) {
		after = &timeRange{
			start: segment.end,
			end:   end,
		}
	}
	return executionPlan[T]{
		before:   before,
		cacheHit: cachedOverlap,
		after:    after,
	}
}

func (c *FeedCache[T]) executePlanAscending(plan executionPlan[T], limit int) (*executionResult[T], error) {
	var beforeData, afterData []*T
	var beforeHasMore, afterHasMore bool
	var err error

	// Fetch data before cache segment if needed
	if plan.before != nil {
		beforeData, err = c.fetchFromDB(plan.before.start, plan.before.end, Ascending, limit)
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
			remaining = limit - len(beforeData) - len(plan.cacheHit)
		}
		if remaining > 0 {
			afterData, err = c.fetchFromDB(plan.after.start, plan.after.end, Ascending, remaining)
			if err != nil {
				return nil, err
			}
			afterHasMore = len(afterData) == remaining
		}
	}

	// Combine results: beforeData -> cachedOverlap -> afterData
	numToReturn := len(beforeData) + len(plan.cacheHit) + len(afterData)
	if limit > 0 {
		numToReturn = min(numToReturn, limit)
	}
	result := make([]*T, 0, numToReturn)

	beforeItems := min(numToReturn, len(beforeData))
	result = append(result, beforeData[:beforeItems]...)

	if len(result) < numToReturn {
		overlapItems := min(numToReturn-len(result), len(plan.cacheHit))
		result = append(result, plan.cacheHit[:overlapItems]...)
	}

	if len(result) < numToReturn {
		afterItems := min(numToReturn-len(result), len(afterData))
		result = append(result, afterData[:afterItems]...)
	}

	return &executionResult[T]{
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

func (c *FeedCache[T]) executePlanDescending(plan executionPlan[T], limit int) (*executionResult[T], error) {
	var beforeData, afterData []*T
	var beforeHasMore, afterHasMore bool
	var err error

	// Fetch data after cache segment if needed
	if plan.after != nil {
		afterData, err = c.fetchFromDB(plan.after.start, plan.after.end, Descending, limit)
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
			remaining = limit - len(beforeData) - len(plan.cacheHit)
		}
		if remaining > 0 {
			beforeData, err = c.fetchFromDB(plan.before.start, plan.before.end, Descending, remaining)
			if err != nil {
				return nil, err
			}
			beforeHasMore = len(beforeData) == remaining
		}
	}

	// Combine results: afterData -> cachedOverlap -> beforeData
	numToReturn := len(beforeData) + len(plan.cacheHit) + len(afterData)
	if limit > 0 {
		numToReturn = min(numToReturn, limit)
	}
	result := make([]*T, 0, numToReturn)

	afterItems := min(numToReturn, len(afterData))
	result = append(result, afterData[:afterItems]...)

	if len(result) < numToReturn {
		overlapItems := min(numToReturn-len(result), len(plan.cacheHit))
		result = append(result, reverseOrder(plan.cacheHit)[:overlapItems]...)
	}

	if len(result) < numToReturn {
		beforeItems := min(numToReturn-len(result), len(beforeData))
		result = append(result, beforeData[:beforeItems]...)
	}

	return &executionResult[T]{
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

func (c *FeedCache[T]) updateCache(result *executionResult[T]) {
	if result.before == nil && result.after == nil {
		return
	}

	before, after := result.before, result.after
	beforeData, afterData := result.beforeData, result.afterData

	if len(beforeData) > 0 {
		start, end := before.start, before.end
		if result.order == Ascending {
			if result.beforeHasMore {
				end = c.getTimestamp(beforeData[len(beforeData)-1]).Add(time.Nanosecond)
			}
		} else {
			beforeData = reverseOrder(beforeData)
			if result.beforeHasMore {
				start = c.getTimestamp(beforeData[0])
			}
		}
		c.mergeCache(beforeData, start, end)
	}

	if len(afterData) > 0 {
		start, end := after.start, after.end
		if result.order == Ascending {
			if result.afterHasMore {
				end = c.getTimestamp(afterData[len(afterData)-1]).Add(time.Nanosecond)
			}
		} else {
			afterData = reverseOrder(afterData)
			if result.afterHasMore {
				start = c.getTimestamp(afterData[0])
			}
		}
		c.mergeCache(afterData, start, end)
	}
}

func (c *FeedCache[T]) mergeCache(data []*T, start, end time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	segment := c.segment

	// No cache yet
	if segment == nil {
		c.setCacheSegment(data, start, end)
		return
	}

	// No overlap with the cache
	if !segment.overlaps(timeRange{start: start, end: end}) {
		// Two special cases: non-overlapping but intervals are connected
		if start.Equal(segment.end) {
			c.setCacheSegment(append(segment.data, data...), segment.start, end)
		}
		if end.Equal(segment.start) {
			c.setCacheSegment(append(data, segment.data...), start, segment.end)
		}

		// If it's a disconnected newer segment, it should replace existing cache
		if start.After(segment.end) {
			c.setCacheSegment(data, start, end)
		}

		return
	}

	// It's a sub range contained in existing cache, do nothing
	if !start.Before(segment.start) && !end.After(segment.end) {
		return
	}

	// The data is newer than cache, extend the cache
	if end.After(segment.end) {
		split := 0
		for ; split < len(data); split++ {
			if !c.getTimestamp(data[split]).Before(c.segment.end) {
				break
			}
		}
		if split < len(data) {
			newData := append(segment.data, data[split:]...)
			c.setCacheSegment(newData, segment.start, end)
		}
		return
	}

	// Now we must have start.Before(segment.start) && segment.start.Before(end)
	split := len(data) - 1
	for ; split >= 0; split-- {
		if c.getTimestamp(data[split]).Before(segment.start) {
			break
		}
	}
	if split >= 0 {
		newData := append(data[:split+1], segment.data...)
		c.setCacheSegment(newData, start, segment.end)
	}
}

func (c *FeedCache[T]) setCacheSegment(data []*T, start, end time.Time) {
	// Ensure we don't exceed maxItems, prioritizing recent data
	if len(data) > c.maxItems {
		data = data[len(data)-c.maxItems:] // Keep newest
		start = c.getTimestamp(data[0])
	}
	c.segment = &cacheEntry[T]{
		timeRange: timeRange{
			start: start,
			end:   end,
		},
		data: data,
	}
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

func reverseOrder[T any](data []*T) []*T {
	result := make([]*T, len(data))
	for i, item := range data {
		result[len(data)-1-i] = item
	}
	return result
}
