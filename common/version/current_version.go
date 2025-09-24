package version

import "fmt"

// The semantic version string of the code in this branch.
const version = "2.4.0"

// Get the current version of the code in this branch.
func CurrentVersion() (*Semver, error) {
	semver, err := SemverFromString(version)
	if err != nil {
		return nil, fmt.Errorf("failed to parse current version: %w", err)
	}
	return semver, nil
}
