package version

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// Semver represents a semantic version.
type Semver struct {
	major  int
	minor  int
	patch  int
	errata string
}

// NewSemver creates a new Semver instance.
func NewSemver(major int, minor int, patch int, errata string) (*Semver, error) {
	if major < 0 {
		return nil, fmt.Errorf("major version must be non-negative")
	}
	if minor < 0 {
		return nil, fmt.Errorf("minor version must be non-negative")
	}
	if patch < 0 {
		return nil, fmt.Errorf("patch version must be non-negative")
	}

	return &Semver{
		major:  major,
		minor:  minor,
		patch:  patch,
		errata: errata,
	}, nil
}

// Parses a semantic version string and returns a Semver instance.
func SemverFromString(versionStr string) (*Semver, error) {
	var major int
	var minor int
	var patch int
	var errata string

	if strings.Contains(versionStr, "-") {
		// Try with errata
		n, err := fmt.Sscanf(versionStr, "%d.%d.%d-%s", &major, &minor, &patch, &errata)
		if err != nil {
			return nil, fmt.Errorf("invalid version format: %w", err)
		}
		if n != 4 {
			return nil, fmt.Errorf("invalid version format")
		}
	} else {
		var extra string
		n, err := fmt.Sscanf(versionStr, "%d.%d.%d%s", &major, &minor, &patch, &extra)
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("invalid version format: %w", err)
		}
		if n != 3 {
			return nil, fmt.Errorf("invalid version format")
		}
	}

	return NewSemver(major, minor, patch, errata)
}

func (s *Semver) String() string {
	errataStr := ""
	if s.errata != "" {
		errataStr = "-" + s.errata
	}

	return fmt.Sprintf("%d.%d.%d%s", s.major, s.minor, s.patch, errataStr)
}

// Get the major version number.
func (s *Semver) Major() int {
	return s.major
}

// Get the minor version number.
func (s *Semver) Minor() int {
	return s.minor
}

// Get the patch version number.
func (s *Semver) Patch() int {
	return s.patch
}

// Get the errata string.
func (s *Semver) Errata() string {
	return s.errata
}

// Compares two Semver instances for equality.
func (s *Semver) Equals(other *Semver) bool {
	if s == nil || other == nil {
		return false
	}
	return s.major == other.major && s.minor == other.minor && s.patch == other.patch
}

// Compares two Semver instances to see if this one is less than the other. Ignores errata.
func (s *Semver) LessThan(other *Semver) bool {
	if s == nil || other == nil {
		return false
	}
	if s.major != other.major {
		return s.major < other.major
	}
	if s.minor != other.minor {
		return s.minor < other.minor
	}
	if s.patch != other.patch {
		return s.patch < other.patch
	}
	return false
}

// Compares two Semver instances to see if this one is greater than the other. Ignores errata.
func (s *Semver) GreaterThan(other *Semver) bool {
	if s == nil || other == nil {
		return false
	}
	if s.major != other.major {
		return s.major > other.major
	}
	if s.minor != other.minor {
		return s.minor > other.minor
	}
	if s.patch != other.patch {
		return s.patch > other.patch
	}
	return false
}

// Compares two Semver instances to see if this one is greater than or equal to the other. Ignores errata.
func (s *Semver) GreaterThanOrEqual(other *Semver) bool {
	return s.GreaterThan(other) || s.Equals(other)
}

// Compares two Semver instances to see if this one is less than or equal to the other. Ignores errata.
func (s *Semver) LessThanOrEqual(other *Semver) bool {
	return s.LessThan(other) || s.Equals(other)
}
