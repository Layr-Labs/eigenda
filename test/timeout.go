package test

import (
	"context"
	"fmt"
	"time"
)

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
