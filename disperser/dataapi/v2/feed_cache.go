package v2

import (
	"context"
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
	mu      sync.RWMutex
	segment *CircularQueue[T]
	// Async updates to the cache segment
	updateWg *sync.WaitGroup

	fetchFromDB  func(ctx context.Context, start, end time.Time, order FetchOrder, limit int) ([]*T, error)
	getTimestamp func(*T) time.Time
}

func NewFeedCache[T any](
	maxItems int,
	fetchFn func(ctx context.Context, start, end time.Time, order FetchOrder, limit int) ([]*T, error),
	timestampFn func(*T) time.Time,
) *FeedCache[T] {
	return &FeedCache[T]{
		segment:      NewCircularQueue[T](maxItems, timestampFn),
		fetchFromDB:  fetchFn,
		getTimestamp: timestampFn,
		updateWg:     &sync.WaitGroup{},
	}
}

// timeRange represents a time interval [start, end) where start is inclusive and end
// is exclusive.
type timeRange struct {
	start time.Time
	end   time.Time
}

// executionPlan describes the breakdown of a data fetch query [start, end) into sub ranges
// that hits cache and that need DB fetches.
type executionPlan[T any] struct {
	// cacheHit is the data items from the cache segment that overlap the query time range.
	cacheHit []*T
	// before is the sub time range that's prior to the cache segment.
	before *timeRange
	// after is the sub time range that's after the cache segment.
	after *timeRange
}

// executionResult describes execution result of a plan.
type executionResult[T any] struct {
	order  FetchOrder
	before *timeRange
	after  *timeRange

	// The DB fetch results corresponding to `before` and `after` ranges.
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

func (c *FeedCache[T]) Get(ctx context.Context, start, end time.Time, queryOrder FetchOrder, limit int) ([]*T, error) {
	if !start.Before(end) {
		return nil, errors.New("the start must be before end")
	}

	plan := c.makePlan(start, end, queryOrder, limit)

	var result *executionResult[T]
	var err error
	if queryOrder == Ascending {
		result, err = c.executePlanAscending(ctx, plan, limit)
	} else {
		result, err = c.executePlanDescending(ctx, plan, limit)
	}
	if err != nil {
		return nil, err
	}

	// Update the cache segment async
	c.updateWg.Add(1)
	go func() {
		defer c.updateWg.Done()
		c.updateCache(result)
	}()

	return result.result, nil
}

func (c *FeedCache[T]) WaitForCacheUpdates() {
	c.updateWg.Wait()
}

func (c *FeedCache[T]) makePlan(start, end time.Time, queryOrder FetchOrder, limit int) executionPlan[T] {
	c.mu.RLock()
	defer c.mu.RUnlock()

	segment := c.segment
	queryRange := timeRange{start: start, end: end}

	// Handle no cache or no overlap cases together
	if segment.size == 0 || !queryRange.overlaps(segment.timeRange) {
		return executionPlan[T]{
			// The data=nil, so it doesn't matter we fill the `before` or `after`
			after: &queryRange,
		}
	}

	// Get cached data that's overlapping the query range
	cachedOverlap := segment.QueryTimeRange(start, end, queryOrder, limit)

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

func (c *FeedCache[T]) executePlanAscending(ctx context.Context, plan executionPlan[T], limit int) (*executionResult[T], error) {
	var beforeData, afterData []*T
	var beforeHasMore, afterHasMore bool
	var err error

	// Fetch data before cache segment if needed
	if plan.before != nil {
		beforeData, err = c.fetchFromDB(ctx, plan.before.start, plan.before.end, Ascending, limit)
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
			afterData, err = c.fetchFromDB(ctx, plan.after.start, plan.after.end, Ascending, remaining)
			if err != nil {
				return nil, err
			}
			afterHasMore = len(afterData) == remaining
		}
	}

	// Combine results: beforeData -> cacheHit -> afterData
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

func (c *FeedCache[T]) executePlanDescending(ctx context.Context, plan executionPlan[T], limit int) (*executionResult[T], error) {
	var beforeData, afterData []*T
	var beforeHasMore, afterHasMore bool
	var err error

	// Fetch data after cache segment if needed
	if plan.after != nil {
		afterData, err = c.fetchFromDB(ctx, plan.after.start, plan.after.end, Descending, limit)
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
			beforeData, err = c.fetchFromDB(ctx, plan.before.start, plan.before.end, Descending, remaining)
			if err != nil {
				return nil, err
			}
			beforeHasMore = len(beforeData) == remaining
		}
	}

	// Combine results: afterData -> cacheHit -> beforeData
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

	c.mu.Lock()
	defer c.mu.Unlock()

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
		c.segment.MergeTimeRange(beforeData, start, end)
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
		c.segment.MergeTimeRange(afterData, start, end)
	}
}

func reverseOrder[T any](data []*T) []*T {
	result := make([]*T, len(data))
	for i, item := range data {
		result[len(data)-1-i] = item
	}
	return result
}
