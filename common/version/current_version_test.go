package version

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCurrentVersionIsValid(t *testing.T) {
	_, err := CurrentVersion()
	require.NoError(t, err)
}

// If the current branch has the format "release/SEMVER", then verify that the current version matches SEMVER.
func TestCurrentVersion(t *testing.T) {
	// Get the current git branch name
	branch, err := getBranchName()
	if err != nil {
		t.Skipf("Cannot get current branch name: %v", err)
		return
	}

	// Check if branch follows the release/SEMVER pattern
	const releasePrefix = "release/"
	if !strings.HasPrefix(branch, releasePrefix) {
		t.Skipf("Current branch '%s' is not a release branch, skipping version check", branch)
		return
	}

	// Extract the expected version from the branch name
	expectedVersionStr := branch[len(releasePrefix):]
	expectedVersion, err := SemverFromString(expectedVersionStr)
	if err != nil {
		t.Fatalf("Branch name contains invalid semver '%s': %v", expectedVersionStr, err)
	}

	// Get the actual current version
	actualVersion, err := CurrentVersion()
	require.NoError(t, err)

	// Verify they match
	require.True(t, actualVersion.Equals(expectedVersion),
		"Current version %s does not match branch version %s",
		actualVersion.String(), expectedVersion.String())
}

// getBranchName returns the current git branch name
func getBranchName() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("get current branch name: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}
