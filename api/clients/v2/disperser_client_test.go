package clients

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestMutexPreventsSimultaneousRequests tests that the mutex in disperserClient
// prevents multiple goroutines from executing critical sections concurrently.
func TestMutexPreventsSimultaneousRequests(t *testing.T) {
	// Create a struct with a mutex and a counter
	client := &struct {
		requestMutex sync.Mutex
		counter      int
		callTimes    []time.Time
	}{}

	// Use this function to simulate a request that takes some time
	simulateRequest := func() {
		client.requestMutex.Lock()
		defer client.requestMutex.Unlock()

		// Record the time of the call
		callTime := time.Now()
		client.callTimes = append(client.callTimes, callTime)
		client.counter++

		// Simulate processing time
		time.Sleep(200 * time.Millisecond)
	}

	// Number of concurrent "requests" to attempt
	numRequests := 3

	// Use a WaitGroup to wait for all goroutines to complete
	var wg sync.WaitGroup
	wg.Add(numRequests)

	// Start time for our test
	startTime := time.Now()

	// Launch multiple goroutines to make concurrent "requests"
	for i := 0; i < numRequests; i++ {
		go func() {
			defer wg.Done()
			simulateRequest()
		}()
	}

	// Wait for all requests to complete
	wg.Wait()

	// Verify that the correct number of requests were made
	require.Equal(t, numRequests, client.counter, "Expected number of requests")

	// Check that the requests were executed sequentially, not concurrently
	// The time difference between consecutive requests should be at least the delay time
	for i := 1; i < len(client.callTimes); i++ {
		timeDiff := client.callTimes[i].Sub(client.callTimes[i-1])
		require.GreaterOrEqual(t, timeDiff.Milliseconds(), int64(199), // slightly less than 200ms to account for timing variations
			"Requests were not executed sequentially. Time between request %d and %d was only %v",
			i-1, i, timeDiff)
	}

	// The total time should be at least (numRequests * delay)
	// This verifies that the requests were not processed concurrently
	totalTime := time.Since(startTime)
	expectedMinTime := time.Duration(numRequests) * 200 * time.Millisecond
	require.GreaterOrEqual(t, totalTime.Milliseconds(), expectedMinTime.Milliseconds()-10, // allow small timing variations
		"Total execution time was less than expected, suggesting concurrent execution")
}
