package testutils

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
	"testing"
	"time"
)

// InitializeRandom initializes the random number generator. If no arguments are provided, then the seed is randomly
// generated. If a single argument is provided, then the seed is fixed to that value.
func InitializeRandom(fixedSeed ...uint64) {

	var seed uint64
	if len(fixedSeed) == 0 {
		rand.Seed(uint64(time.Now().UnixNano()))
		seed = rand.Uint64()
	} else if len(fixedSeed) == 1 {
		seed = fixedSeed[0]
	} else {
		panic("too many arguments, expected exactly one seed")
	}

	fmt.Printf("Random seed: %d\n", seed)
	rand.Seed(seed)
}

// AssertEventuallyTrue asserts that a condition is true within a given duration. Repeatably checks the condition.
func AssertEventuallyTrue(t *testing.T, condition func() bool, duration time.Duration, debugInfo ...any) {
	if len(debugInfo) == 0 {
		debugInfo = []any{"Condition did not become true within the given duration"}
	}

	ticker := time.NewTicker(1 * time.Millisecond)
	select {
	case <-ticker.C:
		if condition() {
			return
		}
	case <-time.After(duration):
		assert.True(t, condition(), debugInfo...)
	}
}

// AssertEventuallyEquals asserts that a function returns a specific value within a given duration.
func AssertEventuallyEquals(t *testing.T, expected any, actual func() any, duration time.Duration, debugInfo ...any) {
	if len(debugInfo) == 0 {
		debugInfo = []any{
			"Expected value did not match actual value within the given duration. Expected: %v, Actual: %v",
			expected,
			actual(),
		}
	}

	condition := func() bool {
		return expected == actual()
	}

	AssertEventuallyTrue(t, condition, duration, debugInfo...)
}

// ExecuteWithTimeout executes a function with a timeout.
// Panics if the function does not complete within the given duration.
func ExecuteWithTimeout(f func(), duration time.Duration, debugInfo ...any) {
	if len(debugInfo) == 0 {
		debugInfo = []any{"Function did not complete within the given duration"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)

	finished := false
	go func() {
		f()
		finished = true
		cancel()
	}()

	<-ctx.Done()

	if !finished {
		panic(fmt.Sprintf(debugInfo[0].(string), debugInfo[1:]...))
	}
}

// RandomBytes generates a random byte slice of a given length.
func RandomBytes(length int) []byte {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return bytes
}

// RandomTime generates a random time.
func RandomTime() time.Time {
	return time.Unix(int64(rand.Int31()), 0)
}

// RandomString generates a random string out of printable ASCII characters.
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
