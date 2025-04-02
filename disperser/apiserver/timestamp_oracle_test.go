package apiserver_test

import (
	"sync"
	"testing"

	"github.com/Layr-Labs/eigenda/disperser/apiserver"
	"github.com/stretchr/testify/assert"
)

func TestGetUniqueTimestamp(t *testing.T) {
	t.Run("Monotonicity", func(t *testing.T) {
		oracle, _ := apiserver.NewTimestampOracle(1, 10)

		var timestamps []uint64
		iterations := 1000

		for i := 0; i < iterations; i++ {
			timestamps = append(timestamps, oracle.GetUniqueTimestamp())
		}

		for i := 1; i < len(timestamps); i++ {
			if timestamps[i] <= timestamps[i-1] {
				assert.True(t, timestamps[i-1] < timestamps[i])
			}
		}
	})

	t.Run("Rapid sequential calls", func(t *testing.T) {
		oracle, _ := apiserver.NewTimestampOracle(1, 10)

		// Make many rapid calls to ensure uniqueness
		timestamps := make(map[uint64]bool)
		iterations := 10000

		for i := 0; i < iterations; i++ {
			ts := oracle.GetUniqueTimestamp()
			assert.False(t, timestamps[ts])
			timestamps[ts] = true
		}
	})

	t.Run("Concurrent calls", func(t *testing.T) {
		oracle, _ := apiserver.NewTimestampOracle(1, 10)
		var wg sync.WaitGroup

		// Collect timestamps from concurrent goroutines
		timestampsChan := make(chan uint64, 1000)
		numGoroutines := 10
		callsPerGoroutine := 100

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < callsPerGoroutine; j++ {
					timestampsChan <- oracle.GetUniqueTimestamp()
				}
			}()
		}

		// Wait for all goroutines to finish
		wg.Wait()
		close(timestampsChan)

		// Check for uniqueness
		timestamps := make(map[uint64]bool)
		for ts := range timestampsChan {
			assert.False(t, timestamps[ts])
			timestamps[ts] = true
		}

		// Verify we got the expected number of timestamps
		expectedCount := numGoroutines * callsPerGoroutine
		assert.Equal(t, expectedCount, len(timestamps))
	})
}
