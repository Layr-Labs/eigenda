package common

import (
	"fmt"
	"strings"

	"golang.org/x/exp/constraints"
)

// the name of a unit step
type unitStep struct {
	// name of the unit step
	name string
	// multiply by this number to get the previous unit step. For example, if this unit is "KiB", the step is 1024.
	// Taking a number of kilobytes and multiplying by 1024 gives you the number of bytes.
	multiple uint64
}

// byteUnits is a list of units for bytes, in increasing order of size.
var byteSteps = []unitStep{
	{"bytes", 1},
	{"KiB", 1024},
	{"MiB", 1024},
	{"GiB", 1024},
	{"TiB", 1024},
	{"PiB", 1024},
	{"EiB", 1024},
}

var timeSteps = []unitStep{
	{"ns", 1},
	{"μs", 1000},
	{"ms", 1000},
	{"s", 1000},
	{"minutes", 60},
	{"hours", 60},
	{"days", 24},
	{"years", 365}, // I don't care that this is imprecise, I refuse to mess with leap years.
}

// prettyPrintUnit formats a quantity in a human-readable way using the provided unit steps. The quantity
// is assumed to be in the smallest supported unit (e.g., bytes, nanoseconds, etc.).
func prettyPrintUnit(quantity uint64, steps []unitStep) string {

	if quantity < steps[1].multiple {
		// Edge case, print without a decimal point if we have the smallest unit.
		return fmt.Sprintf("%d %s", quantity, steps[0].name)
	}

	unit := steps[0].name
	floatQuantity := float64(quantity)

	for i := 1; i < len(steps); i++ {
		if floatQuantity >= float64(steps[i].multiple) {
			floatQuantity /= float64(steps[i].multiple)
			unit = steps[i].name
		} else {
			// We've found the appropriate unit.
			break
		}
	}

	return fmt.Sprintf("%.2f %s", floatQuantity, unit)
}

// PrettyPrintBytes formats a byte count into a human-readable string with appropriate units.
func PrettyPrintBytes(bytes uint64) string {
	return prettyPrintUnit(bytes, byteSteps)
}

// PrettyPrintTime formats a time duration in nanoseconds into a human-readable string with appropriate units.
func PrettyPrintTime(nanoseconds uint64) string {
	return prettyPrintUnit(nanoseconds, timeSteps)
}

// TODO write unit tests for this

// CommaOMatic converts any integer type to a string with thousands separators.
func CommaOMatic[T constraints.Integer](value T) string {
	s := fmt.Sprintf("%d", value) // works for both signed and unsigned

	// Handle negatives without disturbing digit grouping.
	start := 0
	if len(s) > 0 && s[0] == '-' {
		start = 1
	}

	n := len(s) - start
	if n <= 3 {
		return s
	}

	var b strings.Builder
	// + n/3 commas is a safe upper bound
	b.Grow(len(s) + n/3)

	// Preserve sign
	if start == 1 {
		b.WriteByte('-')
	}

	// First group may be 1–3 digits depending on length
	first := n % 3
	if first == 0 {
		first = 3
	}
	b.WriteString(s[start : start+first])

	// Write remaining groups as ,XYZ
	for i := start + first; i < len(s); i += 3 {
		b.WriteByte(',')
		b.WriteString(s[i : i+3])
	}

	return b.String()
}
