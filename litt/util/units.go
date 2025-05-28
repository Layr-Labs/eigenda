package util

import "fmt"

var byteUnits = []string{"bytes", "KiB", "MiB", "GiB", "TiB", "PiB"}

// PrettyPrintBytes formats a byte count into a human-readable string with appropriate units.
func PrettyPrintBytes(bytes uint64) string {
	floatBytes := float64(bytes)
	unitIndex := 0

	for floatBytes >= 1024 && unitIndex < len(byteUnits)-1 {
		floatBytes /= 1024
		unitIndex++
	}

	if unitIndex == 0 {
		return fmt.Sprintf("%d %s", bytes, byteUnits[unitIndex])
	} else {
		return fmt.Sprintf("%.2f %s", floatBytes, byteUnits[unitIndex])
	}
}
