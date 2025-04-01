package v2_test

import (
	"testing"
	"time"

	v2 "github.com/Layr-Labs/eigenda/disperser/dataapi/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// dataPoint represents a simple time-series data point for testing
type dataPoint struct {
	timestamp time.Time
}

// createDataPoints returns a slice of data points with sequential timestamps
func createDataPoints(startTime time.Time, count int) []*dataPoint {
	points := make([]*dataPoint, count)
	current := startTime
	for i := 0; i < count; i++ {
		points[i] = &dataPoint{
			timestamp: current,
		}
		current = current.Add(time.Minute)
	}
	return points
}

// Test MergeTimeRange method
func TestMergeTimeRange(t *testing.T) {
	getTimestamp := func(dp *dataPoint) time.Time { return dp.timestamp }
	now := time.Now()

	tests := []struct {
		name               string
		initialPoints      []*dataPoint
		initialStart       time.Time
		initialEnd         time.Time
		mergePoints        []*dataPoint
		mergeStart         time.Time
		mergeEnd           time.Time
		expectedSize       int
		expectedTimestamps []time.Time
		capacity           int
	}{
		{
			name:               "empty queue initialization",
			initialPoints:      []*dataPoint{},
			mergePoints:        createDataPoints(now, 3),
			mergeStart:         now,
			mergeEnd:           now.Add(3 * time.Minute),
			expectedSize:       3,
			expectedTimestamps: []time.Time{now, now.Add(time.Minute), now.Add(2 * time.Minute)},
			capacity:           5,
		},
		{
			name:          "append connected range",
			initialPoints: createDataPoints(now, 3),
			initialStart:  now,
			initialEnd:    now.Add(3 * time.Minute),
			mergePoints:   createDataPoints(now.Add(3*time.Minute), 2),
			mergeStart:    now.Add(3 * time.Minute),
			mergeEnd:      now.Add(5 * time.Minute),
			expectedSize:  5,
			expectedTimestamps: []time.Time{
				now,
				now.Add(time.Minute),
				now.Add(2 * time.Minute),
				now.Add(3 * time.Minute),
				now.Add(4 * time.Minute),
			},
			capacity: 5,
		},
		{
			name:          "prepend connected range",
			initialPoints: createDataPoints(now.Add(3*time.Minute), 3),
			initialStart:  now.Add(3 * time.Minute),
			initialEnd:    now.Add(6 * time.Minute),
			mergePoints:   createDataPoints(now, 3),
			mergeStart:    now,
			mergeEnd:      now.Add(3 * time.Minute),
			expectedSize:  5, // Limited by capacity
			expectedTimestamps: []time.Time{
				now.Add(time.Minute),
				now.Add(2 * time.Minute),
				now.Add(3 * time.Minute),
				now.Add(4 * time.Minute),
				now.Add(5 * time.Minute),
			}, // Only newest items from prepend
			capacity: 5,
		},
		{
			name:          "replace with newer disconnected range",
			initialPoints: createDataPoints(now, 3),
			initialStart:  now,
			initialEnd:    now.Add(3 * time.Minute),
			mergePoints:   createDataPoints(now.Add(5*time.Minute), 2),
			mergeStart:    now.Add(5 * time.Minute),
			mergeEnd:      now.Add(7 * time.Minute),
			expectedSize:  2,
			expectedTimestamps: []time.Time{
				now.Add(5 * time.Minute),
				now.Add(6 * time.Minute),
			}, // New timestamps from newer range
			capacity: 5,
		},
		{
			name:          "ignore contained range",
			initialPoints: createDataPoints(now, 5),
			initialStart:  now,
			initialEnd:    now.Add(5 * time.Minute),
			mergePoints:   createDataPoints(now.Add(2*time.Minute), 2),
			mergeStart:    now.Add(2 * time.Minute),
			mergeEnd:      now.Add(4 * time.Minute),
			expectedSize:  5, // Unchanged
			expectedTimestamps: []time.Time{
				now,
				now.Add(time.Minute),
				now.Add(2 * time.Minute),
				now.Add(3 * time.Minute),
				now.Add(4 * time.Minute),
			}, // Unchanged
			capacity: 5,
		},
		{
			name:          "extend end range",
			initialPoints: createDataPoints(now, 3),
			initialStart:  now,
			initialEnd:    now.Add(3 * time.Minute),
			mergePoints:   createDataPoints(now.Add(2*time.Minute), 3),
			mergeStart:    now.Add(2 * time.Minute),
			mergeEnd:      now.Add(5 * time.Minute),
			expectedSize:  5,
			expectedTimestamps: []time.Time{
				now,
				now.Add(time.Minute),
				now.Add(2 * time.Minute),
				now.Add(3 * time.Minute),
				now.Add(4 * time.Minute),
			}, // Original plus new items past the end
			capacity: 5,
		},
		{
			name:          "extend start range",
			initialPoints: createDataPoints(now.Add(2*time.Minute), 3),
			initialStart:  now.Add(2 * time.Minute),
			initialEnd:    now.Add(5 * time.Minute),
			mergePoints:   createDataPoints(now, 3),
			mergeStart:    now,
			mergeEnd:      now.Add(3 * time.Minute),
			expectedSize:  5,
			expectedTimestamps: []time.Time{
				now,
				now.Add(time.Minute),
				now.Add(2 * time.Minute),
				now.Add(3 * time.Minute),
				now.Add(4 * time.Minute),
			}, // New items that extend start plus original
			capacity: 5,
		},
		{
			name:          "capacity constraint drops oldest",
			initialPoints: createDataPoints(now, 3),
			initialStart:  now,
			initialEnd:    now.Add(3 * time.Minute),
			mergePoints:   createDataPoints(now.Add(3*time.Minute), 5),
			mergeStart:    now.Add(3 * time.Minute),
			mergeEnd:      now.Add(8 * time.Minute),
			expectedSize:  5,
			expectedTimestamps: []time.Time{
				now.Add(3 * time.Minute),
				now.Add(4 * time.Minute),
				now.Add(5 * time.Minute),
				now.Add(6 * time.Minute),
				now.Add(7 * time.Minute),
			}, // Only newest 5 items fit
			capacity: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue := v2.NewCircularQueue[dataPoint](tt.capacity, getTimestamp)

			// Setup initial state if needed
			if len(tt.initialPoints) > 0 {
				queue.MergeTimeRange(tt.initialPoints, tt.initialStart, tt.initialEnd)
			}

			// Execute the target operation being tested
			queue.MergeTimeRange(tt.mergePoints, tt.mergeStart, tt.mergeEnd)

			// Verify results
			// Note we fetch all items from the queue as the purpose here is to verify the desired
			// cache state
			results := queue.QueryTimeRange(time.Time{}, now.Add(24*time.Hour), v2.Ascending, 0)
			require.Equal(t, len(tt.expectedTimestamps), len(results))
			for i, expectedTime := range tt.expectedTimestamps {
				assert.True(t, expectedTime.Equal(results[i].timestamp),
					"Expected timestamp %v at index %d, got %v", expectedTime, i, results[i].timestamp)
			}
		})
	}
}

// Test QueryTimeRange method
func TestQueryTimeRange(t *testing.T) {
	getTimestamp := func(dp *dataPoint) time.Time { return dp.timestamp }
	now := time.Now()

	tests := []struct {
		name               string
		points             []*dataPoint
		queryStart         time.Time
		queryEnd           time.Time
		order              v2.FetchOrder
		limit              int
		expectedTimestamps []time.Time
		capacity           int
	}{
		{
			name:               "empty queue",
			points:             []*dataPoint{},
			queryStart:         now,
			queryEnd:           now.Add(5 * time.Minute),
			order:              v2.Ascending,
			limit:              0,
			expectedTimestamps: []time.Time{},
			capacity:           5,
		},
		{
			name:       "exact range match",
			points:     createDataPoints(now, 5),
			queryStart: now,
			queryEnd:   now.Add(5 * time.Minute),
			order:      v2.Ascending,
			limit:      0,
			expectedTimestamps: []time.Time{
				now,
				now.Add(time.Minute),
				now.Add(2 * time.Minute),
				now.Add(3 * time.Minute),
				now.Add(4 * time.Minute),
			},
			capacity: 5,
		},
		{
			name:       "partial range match",
			points:     createDataPoints(now, 5),
			queryStart: now.Add(2 * time.Minute),
			queryEnd:   now.Add(4 * time.Minute),
			order:      v2.Ascending,
			limit:      0,
			expectedTimestamps: []time.Time{
				now.Add(2 * time.Minute),
				now.Add(3 * time.Minute),
			},
			capacity: 5,
		},
		{
			name:               "no range overlap",
			points:             createDataPoints(now, 5),
			queryStart:         now.Add(6 * time.Minute),
			queryEnd:           now.Add(8 * time.Minute),
			order:              v2.Ascending,
			limit:              0,
			expectedTimestamps: []time.Time{},
			capacity:           5,
		},
		{
			name:       "limit ascending",
			points:     createDataPoints(now, 5),
			queryStart: now,
			queryEnd:   now.Add(5 * time.Minute),
			order:      v2.Ascending,
			limit:      3,
			expectedTimestamps: []time.Time{
				now,
				now.Add(time.Minute),
				now.Add(2 * time.Minute),
			},
			capacity: 5,
		},
		{
			name:       "limit descending",
			points:     createDataPoints(now, 5),
			queryStart: now,
			queryEnd:   now.Add(5 * time.Minute),
			order:      v2.Descending,
			limit:      3,
			expectedTimestamps: []time.Time{
				now.Add(2 * time.Minute),
				now.Add(3 * time.Minute),
				now.Add(4 * time.Minute),
			},
			capacity: 5,
		},
		{
			name:       "limit larger than range",
			points:     createDataPoints(now, 3),
			queryStart: now,
			queryEnd:   now.Add(3 * time.Minute),
			order:      v2.Ascending,
			limit:      10,
			expectedTimestamps: []time.Time{
				now,
				now.Add(time.Minute),
				now.Add(2 * time.Minute),
			},
			capacity: 5,
		},
		{
			name:       "zero limit returns all",
			points:     createDataPoints(now, 5),
			queryStart: now,
			queryEnd:   now.Add(5 * time.Minute),
			order:      v2.Ascending,
			limit:      0,
			expectedTimestamps: []time.Time{
				now,
				now.Add(time.Minute),
				now.Add(2 * time.Minute),
				now.Add(3 * time.Minute),
				now.Add(4 * time.Minute),
			},
			capacity: 5,
		},
		{
			name:       "negative limit returns all",
			points:     createDataPoints(now, 5),
			queryStart: now,
			queryEnd:   now.Add(5 * time.Minute),
			order:      v2.Ascending,
			limit:      -1,
			expectedTimestamps: []time.Time{
				now,
				now.Add(time.Minute),
				now.Add(2 * time.Minute),
				now.Add(3 * time.Minute),
				now.Add(4 * time.Minute),
			},
			capacity: 5,
		},
		{
			name:       "start time equals data point time",
			points:     createDataPoints(now, 5),
			queryStart: now.Add(2 * time.Minute),
			queryEnd:   now.Add(5 * time.Minute),
			order:      v2.Ascending,
			limit:      0,
			expectedTimestamps: []time.Time{
				now.Add(2 * time.Minute),
				now.Add(3 * time.Minute),
				now.Add(4 * time.Minute),
			},
			capacity: 5,
		},
		{
			name:       "end time equals data point time (exclusive)",
			points:     createDataPoints(now, 5),
			queryStart: now,
			queryEnd:   now.Add(3 * time.Minute),
			order:      v2.Ascending,
			limit:      0,
			expectedTimestamps: []time.Time{
				now,
				now.Add(time.Minute),
				now.Add(2 * time.Minute),
			},
			capacity: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue := v2.NewCircularQueue[dataPoint](tt.capacity, getTimestamp)

			// Setup initial data
			if len(tt.points) > 0 {
				// Add a ns to make it exclusive
				queue.MergeTimeRange(tt.points, tt.points[0].timestamp,
					tt.points[len(tt.points)-1].timestamp.Add(time.Nanosecond))
			}

			// Execute the query
			results := queue.QueryTimeRange(tt.queryStart, tt.queryEnd, tt.order, tt.limit)

			// Verify results
			require.Equal(t, len(tt.expectedTimestamps), len(results))
			for i, expectedTime := range tt.expectedTimestamps {
				assert.True(t, expectedTime.Equal(results[i].timestamp),
					"Expected timestamp %v at index %d, got %v", expectedTime, i, results[i].timestamp)
			}
		})
	}
}
