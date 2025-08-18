package common

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Helper utility functions //

func ContainsDuplicates[P comparable](s []P) bool {
	seen := make(map[P]struct{})
	for _, v := range s {
		if _, ok := seen[v]; ok {
			return true
		}
		seen[v] = struct{}{}
	}
	return false
}

func Contains[P comparable](s []P, e P) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func ParseBytesAmount(s string) (uint64, error) {
	s = strings.TrimSpace(s)

	// Extract numeric part and unit
	numStr := s
	unit := ""
	for i, r := range s {
		if !('0' <= r && r <= '9' || r == '.') { //nolint:staticcheck // QF1001 cleaner this way than applying DeMorgan's law
			numStr = s[:i]
			unit = s[i:]
			break
		}
	}

	// Convert numeric part to float64
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid numeric value: %w", err)
	}

	unit = strings.ToLower(strings.TrimSpace(unit))

	// Convert to uint64 based on the unit (case-insensitive)
	switch unit {
	case "b", "":
		return uint64(num), nil
	case "kib":
		return uint64(num * 1024), nil
	case "mib":
		return uint64(num * 1024 * 1024), nil
	default:
		return 0, fmt.Errorf("unsupported unit: %s", unit)
	}
}

// EigenDABackend is an enum representing various eigenDA backends
type EigenDABackend uint8

const (
	V1EigenDABackend EigenDABackend = iota + 1
	V2EigenDABackend
)

// Used when marshalling the proxy config and logging to stdout at proxy startup.
// []uint8 gets marshalled as a base64 string by default, which is unreadable.
// This makes it so that it'll be marshalled as an array of strings instead.
func (e EigenDABackend) MarshalJSON() ([]byte, error) {
	return json.Marshal(EigenDABackendToString(e))
}

type InvalidBackendError struct {
	Backend string
}

func (e InvalidBackendError) Error() string {
	return fmt.Sprintf("invalid backend option: %s", e.Backend)
}

// StringToEigenDABackend converts a string to EigenDABackend enum value.
// It returns an [InvalidBackendError] if the input string does not match any known backend,
// which is automatically converted to a 400 Bad Request error by the error middleware.
func StringToEigenDABackend(inputString string) (EigenDABackend, error) {
	inputString = strings.ToUpper(strings.TrimSpace(inputString))

	switch inputString {
	case "V1":
		return V1EigenDABackend, nil
	case "V2":
		return V2EigenDABackend, nil
	default:
		return 0, InvalidBackendError{Backend: inputString}
	}
}

// EigenDABackendToString converts an EigenDABackend enum to its string representation
func EigenDABackendToString(backend EigenDABackend) string {
	switch backend {
	case V1EigenDABackend:
		return "V1"
	case V2EigenDABackend:
		return "V2"
	default:
		return "unknown"
	}
}
