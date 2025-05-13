package memory

import (
	"os"
	"strconv"
	"strings"

	"github.com/docker/go-units"
	"github.com/shirou/gopsutil/mem"
)

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

// getCgroupMemoryLimit attempts to read the memory limit from cgroups
// This is relevant when running in a Docker container or other containerized environment
func getCgroupMemoryLimit() (uint64, error) {
	// Check cgroup v2 first
	cgroup2Path := "/sys/fs/cgroup/memory.max"
	if _, err := os.Stat(cgroup2Path); err == nil {
		return readCgroupFile(cgroup2Path)
	}

	// Check cgroup v1
	cgroup1Path := "/sys/fs/cgroup/memory/memory.limit_in_bytes"
	if _, err := os.Stat(cgroup1Path); err == nil {
		return readCgroupFile(cgroup1Path)
	}

	// Also check the Docker-specific path in cgroup v1
	dockerCgroupPath := "/sys/fs/cgroup/memory/docker/memory.limit_in_bytes"
	if _, err := os.Stat(dockerCgroupPath); err == nil {
		return readCgroupFile(dockerCgroupPath)
	}

	// Try to read from the proc status, which can sometimes have container limits
	return readProcStatusMemoryLimit()
}

// readCgroupFile reads and parses a cgroup memory limit file
func readCgroupFile(path string) (uint64, error) {
	data, err := os.ReadFile(path)
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
	data, err := os.ReadFile("/proc/self/status")
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Limit:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				// Parse with support for units like kB
				valueStr := fields[1]
				if strings.HasSuffix(valueStr, "kB") {
					valueStr = strings.TrimSuffix(valueStr, "kB")
					value, err := strconv.ParseUint(valueStr, 10, 64)
					if err != nil {
						return 0, err
					}
					return value * 1024, nil
				}

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
