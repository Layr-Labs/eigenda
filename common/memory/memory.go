package memory

import (
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/docker/go-units"
	"github.com/shirou/gopsutil/mem"
)

// Variable to allow mocking in tests
var readFile = os.ReadFile

// potential cgroup paths to check for memory limits
var cgroupPaths = []string{
	"/sys/fs/cgroup/memory.max",
	"/sys/fs/cgroup/memory/memory.limit_in_bytes",
	"/sys/fs/cgroup/memory/docker/memory.limit_in_bytes",
}

// unitSuffixes maps common memory unit suffixes to their byte multipliers
var unitSuffixes = map[string]uint64{
	"kb": units.KiB,
	"mb": units.MiB,
	"gb": units.GiB,
	"tb": units.TiB,
}

// GetMaximumAvailableMemory returns the maximum available memory in bytes, i.e. the maximum quantity of memory that
// this process can allocate before experiencing an out of memory error. Handles artificial limits set by the OS and/or
// docker container.
func GetMaximumAvailableMemory() (uint64, error) {
	// Get the system's total memory first
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	systemTotal := vmStat.Total

	// Check if there's a cgroup limit (for Docker/container environments)
	cgroupLimit, err := getCgroupMemoryLimit()
	if err == nil && cgroupLimit > 0 && cgroupLimit < systemTotal {
		// If there's a valid cgroup limit, use it
		return cgroupLimit, nil
	}

	// If no cgroup limit is found, cgroup returns 0 (indicating no limit),
	// or if the cgroup limit exceeds physical memory,
	// or there was an error reading it, return the total system memory
	return systemTotal, nil
}

// SetGCMemorySafetyBuffer tells the garbage collector to aggressively garbage collect when there is only safetyBuffer
// bytes of memory available. Useful for preventing kubernetes from OOM-killing the process.
func SetGCMemorySafetyBuffer(
	logger logging.Logger,
	safetyBuffer uint64,
) error {

	maxMemory, err := GetMaximumAvailableMemory()
	if err != nil {
		return fmt.Errorf("failed to get maximum available memory: %w", err)
	}

	if safetyBuffer > maxMemory {
		return fmt.Errorf("buffer space %d exceeds maximum available memory %d", safetyBuffer, maxMemory)
	}

	limit := maxMemory - safetyBuffer

	debug.SetMemoryLimit(int64(limit))

	logger.Infof("Detected %.2fGB available memory. "+
		"Setting GC target memory limit to %.2f in order to maintain a safety buffer of %.2fGB",
		float64(maxMemory)/float64(units.GiB),
		float64(limit)/float64(units.GiB),
		float64(safetyBuffer)/float64(units.GiB))

	return nil
}

// getCgroupMemoryLimit attempts to read the memory limit from cgroups
// This is relevant when running in a Docker container or other containerized environment
func getCgroupMemoryLimit() (uint64, error) {
	for _, path := range cgroupPaths {
		if _, err := os.Stat(path); err == nil {
			// File exists, read it
			return readCgroupFile(path)
		}
	}

	// Try to read from the proc status, which can sometimes have container limits
	return readProcStatusMemoryLimit()
}

// readCgroupFile reads and parses a cgroup memory limit file
func readCgroupFile(path string) (uint64, error) {
	data, err := readFile(path)
	if err != nil {
		return 0, err
	}

	// Clean the string and handle "max" value which means no limit
	valueStr := strings.TrimSpace(string(data))
	if valueStr == "max" || valueStr == "-1" {
		return 0, nil // No limit
	}

	// Parse the value
	value, err := strconv.ParseUint(valueStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}

// readProcStatusMemoryLimit attempts to get the memory limit from /proc/self/status
// which can reflect container limits
func readProcStatusMemoryLimit() (uint64, error) {
	data, err := readFile("/proc/self/status")
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Limit:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				valueStr := fields[1]
				valueLower := strings.ToLower(valueStr)

				for unitSuffix, multiplier := range unitSuffixes {
					if strings.HasSuffix(valueLower, unitSuffix) {
						// Remove the unit suffix and parse the numeric value
						numStr := valueStr[:len(valueStr)-len(unitSuffix)]
						value, err := strconv.ParseUint(numStr, 10, 64)
						if err != nil {
							continue // Try next suffix if parsing fails
						}
						return value * multiplier, nil
					}
				}

				// Fallback to the general parser if no explicit unit match was found
				value, err := units.RAMInBytes(valueStr)
				if err != nil {
					return 0, err
				}
				return uint64(value), nil
			}
		}
	}

	return 0, nil
}
