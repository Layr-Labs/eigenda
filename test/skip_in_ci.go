package test

import (
	"os"
	"testing"
)

// Causes the test to be skipped if running in a CI environment. Specifically, skips the test if the "CI" environment
// variable is set. (This variable is set by the GitHub action.)
func SkipInCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping test in CI environment")
	}
}