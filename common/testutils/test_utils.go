package testutils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
	"testing"
	"time"
)

// InitializeRandom initializes the random number generator. Prints the seed so that the test can be rerun
// deterministically. Replace a call to this method with a call to initializeRandomWithSeed to rerun a test
// with a specific seed.
func InitializeRandom() {
	rand.Seed(uint64(time.Now().UnixNano()))
	seed := rand.Uint64()
	fmt.Printf("Random seed: %d\n", seed)
	rand.Seed(seed)
}

// InitializeRandomWithSeed initializes the random number generator with a specific seed.
func InitializeRandomWithSeed(seed uint64) {
	fmt.Printf("Random seed: %d\n", seed)
	rand.Seed(seed)
}

// AssertEventuallyTrue asserts that a condition is true within a given duration. Repeatably checks the condition.
func AssertEventuallyTrue(t *testing.T, condition func() bool, duration time.Duration) {
	start := time.Now()
	for time.Since(start) < duration {
		if condition() {
			return
		}
		time.Sleep(1 * time.Millisecond)
	}
	assert.True(t, condition(), "Condition did not become true within the given duration")
}

// ExecuteWithTimeout executes a function with a timeout.
// Panics if the function does not complete within the given duration.
func ExecuteWithTimeout(f func(), duration time.Duration) {
	done := make(chan struct{})
	go func() {
		f()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(duration):
		panic("function did not complete within the given duration")
	}
}
