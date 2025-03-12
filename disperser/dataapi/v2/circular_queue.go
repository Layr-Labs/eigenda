package v2

import (
	"time"
)

// CircularQueue describes a segment of results fetched for a time range.
//
// It has the following properties:
// - all data items that are in range [start, end) are in the segment
// - no data items that are outside the range [start, end) are included in the segment
//
// The new segment of results can be appended or prepended to the queue if the time range
// that they represent is connected to the cached segment. If over capacity, it'll always
// evict the oldest items.
//
// The data items are in ascending order by timestamp.
//
// This implementation is NOT thread-safe. The caller must ensure proper synchronization
// when used across multiple threads.
type CircularQueue[T any] struct {
	timeRange

	items    []*T // circular queue
	head     int  // Index of the oldest element
	size     int  // Current number of elements
	capacity int  // Maximum capacity (queue length)

	getTimestamp func(*T) time.Time // Function to extract timestamp from items
}

// NewCircularQueue creates a new CircularQueue with the specified capacity.
func NewCircularQueue[T any](capacity int, getTimestampFn func(*T) time.Time) *CircularQueue[T] {
	return &CircularQueue[T]{
		items:        make([]*T, capacity),
		head:         0,
		size:         0,
		capacity:     capacity,
		getTimestamp: getTimestampFn,
	}
}

// QueryTimeRange returns cached data that's in time range [start, end).
// If there are more than `limit` elements, it will cut off the results up to `limit`.
//
// The parameters:
//   - [start, end): The inclusive start time and exclusive end time of the query range.
//   - order: The order in which to fetch results (Ascending or Descending).
//     For ascending order, it'll get the oldest `limit` elements in range;
//     for descending order, it'll get the newest `limit` elements in range.
//   - limit: The desired number of elements to return. If limit <= 0, all matching
//     elements are returned.
func (c *CircularQueue[T]) QueryTimeRange(start, end time.Time, order FetchOrder, limit int) []*T {
	if c.size == 0 {
		return []*T{}
	}

	// Find start and end indices of the overlap
	startIdx := -1
	endIdx := -1

	for i := 0; i < c.size; i++ {
		idx := (c.head + i) % c.capacity

		timestamp := c.getTimestamp(c.items[idx])

		// Found the first item at or after the start time
		if startIdx == -1 && !timestamp.Before(start) {
			startIdx = i
		}

		// Found the first item at or past the end time (exclusive)
		if !timestamp.Before(end) {
			endIdx = i
			break
		}
	}

	// If we never found the end, set it to the end of data
	if endIdx == -1 {
		endIdx = c.size
	}
	// No overlap found
	if startIdx == -1 || startIdx >= endIdx {
		return []*T{}
	}

	// Calculate how many items in the overlap
	overlapCount := endIdx - startIdx
	// Apply limit if needed
	if limit > 0 && limit < overlapCount {
		if order == Ascending {
			// For ascending, take first 'limit' items
			endIdx = startIdx + limit
		} else {
			// For descending, take last 'limit' items
			startIdx = endIdx - limit
		}
	}

	// Note: we need to make a copy of the overlap because the cache data can be mutated
	// by other threads after this function returns (within this function, the caller
	// makes sure a reader lock is held). The data is of type *T, so it won't deep copy
	// the data, just the pointers.
	result := make([]*T, endIdx-startIdx)
	for i := 0; i < endIdx-startIdx; i++ {
		idx := (c.head + startIdx + i) % c.capacity
		result[i] = c.items[idx]
	}

	return result
}

// MergeTimeRange merges a new segment of results representing time range [start, end) to
// the existing cache.
//
// Behavior:
//   - If the queue is empty, initializes it with the provided data.
//   - If the time ranges don't overlap but are connected, appends or prepends data as appropriate.
//   - If the new time range is disconnected but newer, replaces the queue contents.
//   - If the new time range overlaps, extends the range as needed.
//   - If the new time range is entirely contained within the existing range, does nothing.
//
// This method handles these cases to ensure the time range invariant is maintained,
// while prioritizing newer data when capacity constraints are encountered.
func (c *CircularQueue[T]) MergeTimeRange(items []*T, start, end time.Time) {
	if len(items) == 0 {
		return
	}

	if c.size == 0 {
		c.reset(items)
		c.start, c.end = maxTimestamp(start, c.headTimestamp()), end
		return
	}

	if !c.overlaps(timeRange{start: start, end: end}) {
		// Two special cases: non-overlapping but internewItems are connected
		if start.Equal(c.end) {
			c.appendItems(items)
			c.start, c.end = maxTimestamp(c.start, c.headTimestamp()), end
		}
		if end.Equal(c.start) {
			c.prependItems(items)
			// Note c.end unchanged
			c.start = maxTimestamp(start, c.headTimestamp())
		}

		// If it's a disconnected newer segment, it should replace existing cache
		if start.After(c.end) {
			c.reset(items)
			c.start, c.end = maxTimestamp(start, c.headTimestamp()), end
		}

		return
	}

	// It's a sub range contained in existing cache, do nothing
	if !start.Before(c.start) && !end.After(c.end) {
		return
	}

	// The data is newer than cache, extend the cache
	if end.After(c.end) {
		split := 0
		for ; split < len(items); split++ {
			if !c.getTimestamp(items[split]).Before(c.end) {
				break
			}
		}
		if split < len(items) {
			c.appendItems(items[split:])
			c.start, c.end = maxTimestamp(c.start, c.headTimestamp()), end
		}
		return
	}

	// Now we must have start.Before(segment.start) && segment.start.Before(end)
	split := len(items) - 1
	for ; split >= 0; split-- {
		if c.getTimestamp(items[split]).Before(c.start) {
			break
		}
	}
	if split >= 0 {
		c.prependItems(items[:split+1])
		// Note c.end unchanged
		c.start = maxTimestamp(start, c.headTimestamp())
	}
}

// headTimestamp returns the timestamp of the head element in the queue.
// Assumes the queue is not empty (c.size > 0) and c.items[c.head] is not nil (ensured by
// the caller).
func (c *CircularQueue[T]) headTimestamp() time.Time {
	return c.getTimestamp(c.items[c.head])
}

// reset initializes the cache with the given data, limiting to capacity
// This method resets the circular queue and adds at most capacity elements,
// prioritizing the most recent (latest timestamp) elements if needed.
// If newItems is empty, this method does nothing.
func (c *CircularQueue[T]) reset(newItems []*T) {
	if len(newItems) == 0 {
		return
	}

	// Reset the circular queue
	c.head = 0
	c.size = 0

	// Determine how many data points to use (up to capacity)
	numToAdd := len(newItems)
	startIdx := 0
	if numToAdd > c.capacity {
		// Only add the most recent points that fit in the capacity
		startIdx = len(newItems) - c.capacity
		numToAdd = c.capacity
	}

	// Add data points directly to the queue without function calls
	for i := 0; i < numToAdd; i++ {
		c.items[i] = newItems[startIdx+i]
	}
	// Update size
	c.size = numToAdd
}

// prependItems adds multiple elements to the front of the queue.
// Elements must be in ascending time order (oldest to newest).
// This never drops newer elements to make room for older ones.
func (c *CircularQueue[T]) prependItems(newItems []*T) {
	if len(newItems) == 0 {
		return
	}

	// If queue is empty, just initialize with the data
	if c.size == 0 {
		c.reset(newItems)
		return
	}

	// Calculate how many elements we can actually add
	// We never drop newer elements to make room for older ones
	spaceAvailable := c.capacity - c.size
	numToAdd := len(newItems)
	if numToAdd > spaceAvailable {
		numToAdd = spaceAvailable
	}

	// Queue is full, no room to add older elements
	if numToAdd <= 0 {
		return
	}

	// Only add the newest numToAdd elements from newItems
	// This means we take the last numToAdd elements from the array
	startIdx := len(newItems) - numToAdd

	// Add elements one by one to the front, starting with the newest
	// to preserve ascending time order in the queue
	for i := len(newItems) - 1; i >= startIdx; i-- {
		// Move head back and increase size
		c.head = (c.head - 1 + c.capacity) % c.capacity
		c.items[c.head] = newItems[i]
	}

	c.size += numToAdd
}

// appendItems adds multiple elements to the back of the queue.
// Elements must be in ascending time order (oldest to newest).
// Drops oldest elements if necessary to make room for newer ones.
func (c *CircularQueue[T]) appendItems(newItems []*T) {
	if len(newItems) == 0 {
		return
	}

	// If queue is empty, just initialize with the data
	if c.size == 0 {
		c.reset(newItems)
		return
	}

	// If new data exceeds capacity, use only the newest portion
	if len(newItems) >= c.capacity {
		c.reset(newItems)
		return
	}

	// Calculate if we need to drop oldest elements
	totalSize := c.size + len(newItems)
	overflow := totalSize - c.capacity
	if overflow > 0 {
		// We need to drop some oldest elements
		c.head = (c.head + overflow) % c.capacity
		c.size -= overflow
	}

	// Add new elements to the back
	for _, val := range newItems {
		idx := (c.head + c.size) % c.capacity
		c.items[idx] = val
		c.size++
	}
}

func maxTimestamp(t1, t2 time.Time) time.Time {
	if t1.Before(t2) {
		return t2
	}
	return t1
}
