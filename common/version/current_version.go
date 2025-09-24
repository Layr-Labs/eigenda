package version

import "fmt"

// The semantic version string of the code in this branch.
var version = "2.4.0"

// Call this to override the version string (for example, with a more specific build version).
func SetVersion(versionString string) error {
	oldVersion := version
	version = versionString

	_, err := CurrentVersion()
	if err != nil {
		version = oldVersion
		return fmt.Errorf("invalid version string: %w", err)
	}

	return nil
}

// Get the current version of the code in this branch.
func CurrentVersion() (*Semver, error) {
	semver, err := SemverFromString(version)
	if err != nil {
		return nil, fmt.Errorf("failed to parse current version: %w", err)
	}
	return semver, nil
}
