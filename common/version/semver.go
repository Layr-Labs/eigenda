package version

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

var _ fmt.Stringer = (*Semver)(nil)

// Semver represents a semantic version.
type Semver struct {
	major  uint64
	minor  uint64
	patch  uint64
	errata string
}

// NewSemver creates a new Semver instance.
func NewSemver(major uint64, minor uint64, patch uint64, errata string) *Semver {
	return &Semver{
		major:  major,
		minor:  minor,
		patch:  patch,
		errata: errata,
	}
}

// Parses a semantic version string and returns a Semver instance.
//
// Requires the string to have the following format: X.Y.Z[-errata], where X, Y, and Z are
// non-negative integers, and errata is an optional arbitrary string. Note that if
// errata is present, it must be preceded by a hyphen, e.g. "1.2.3-alpha.1".
func SemverFromString(versionStr string) (*Semver, error) {
	var major uint64
	var minor uint64
	var patch uint64
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

		// "extra" will catch any trailing characters after the last integer. If we have trailing characters, they
		// should always be preceded by a hyphen. Since in this branch we don't have a hyphen, consider any trailing
		// characters to be an error.
		var extra string

		n, err := fmt.Sscanf(versionStr, "%d.%d.%d%s", &major, &minor, &patch, &extra)
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("invalid version format: %w", err)
		}
		if n != 3 || extra != "" {
			return nil, fmt.Errorf("invalid version format")
		}
	}

	return NewSemver(major, minor, patch, errata), nil
}

func (s *Semver) String() string {
	errataStr := ""
	if s.errata != "" {
		errataStr = "-" + s.errata
	}

	return fmt.Sprintf("%d.%d.%d%s", s.major, s.minor, s.patch, errataStr)
}

// Get the major version number.
func (s *Semver) Major() uint64 {
	return s.major
}

// Get the minor version number.
func (s *Semver) Minor() uint64 {
	return s.minor
}

// Get the patch version number.
func (s *Semver) Patch() uint64 {
	return s.patch
}

// Get the errata string.
func (s *Semver) Errata() string {
	return s.errata
}

// Compares two Semver instances. Returns -1 if a < b, 1 if a > b, and 0 if a == b.
// Panics if either a or b is nil. Ignores the errata field.
func SemverComparator(a *Semver, b *Semver) int {
	if a.major > b.major {
		return 1
	}
	if a.major < b.major {
		return -1
	}
	if a.minor > b.minor {
		return 1
	}
	if a.minor < b.minor {
		return -1
	}
	if a.patch > b.patch {
		return 1
	}
	if a.patch < b.patch {
		return -1
	}
	return 0
}

// Compares two Semver instances for equality. Ignores errata.
func (s *Semver) Equals(other *Semver) bool {
	return SemverComparator(s, other) == 0
}

// Compares two Semver instances to see if this one is less than the other. Ignores errata.
func (s *Semver) LessThan(other *Semver) bool {
	return SemverComparator(s, other) == -1
}

// Compares two Semver instances to see if this one is greater than the other. Ignores errata.
func (s *Semver) GreaterThan(other *Semver) bool {
	return SemverComparator(s, other) == 1
}

// Compares two Semver instances to see if this one is greater than or equal to the other. Ignores errata.
func (s *Semver) GreaterThanOrEqual(other *Semver) bool {
	return SemverComparator(s, other) >= 0
}

// Compares two Semver instances to see if this one is less than or equal to the other. Ignores errata.
func (s *Semver) LessThanOrEqual(other *Semver) bool {
	return SemverComparator(s, other) <= 0
}

// Compares two Semver instances for strict equality, including errata.
func (s *Semver) StrictEquals(other *Semver) bool {
	return s.Equals(other) && s.errata == other.errata
}
