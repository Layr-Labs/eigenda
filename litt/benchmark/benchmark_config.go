package benchmark

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Layr-Labs/eigenda/litt"
	"github.com/docker/go-units"
)

// BenchmarkConfig is a struct that holds the configuration for the benchmark.
type BenchmarkConfig struct {

	// Configuration for the LittDB instance.
	LittConfig *litt.Config

	// The location where the benchmark stores test metadata.
	MetadataDirectory string

	// The maximum target throughput in MB/s.
	MaximumThroughputMB float64

	// The size of the values in MB.
	ValueSizeMB float64

	// Data is written to the DB in batches and then flushed. This determines the size of those batches, in MB.
	BatchSizeMB float64

	// The maximum number of batches permitted to be in-flight at any given time. If batches can't be written/flushed
	// fast enough, then the benchmark engine will not push data as fast as requested.
	BatchParallelism int

	// The frequency at which the benchmark does cohort garbage collection, in seconds
	CohortGCPeriodSeconds float64

	// The size of the write info channel. Controls the max number of keys to prepare for writing ahead of time.
	WriteInfoChanelSize uint64

	// The size of the read info channel. Controls the max number of keys to prepare for reading ahead of time.
	ReadInfoChanelSize uint64

	// The number of keys in a new cohort.
	CohortSize uint64

	// The time-to-live (TTL) for keys in the database, in hours.
	TTLHours float64

	// If data is within this many minutes of its expiration time, it will not be read.
	ReadSafetyMarginMinutes float64

	// A seed for the random number generator used to generate keys and values. When restarting the benchmark,
	// it's important to always use the same seed.
	Seed int64

	// The size of the pool of random data. Instead of generating random data for each key/value pair
	// (which is expensive), data from this pool is reused. When restarting the benchmark,
	// it's important to always use the same pool size.
	RandomPoolSize uint64
}

// DefaultBenchmarkConfig returns a default BenchmarkConfig with the given data paths.
func DefaultBenchmarkConfig() *BenchmarkConfig {
	return &BenchmarkConfig{
		LittConfig:              litt.DefaultConfigNoPaths(),
		MetadataDirectory:       "~/benchmark",
		MaximumThroughputMB:     10,
		ValueSizeMB:             2.0,
		BatchSizeMB:             32,
		BatchParallelism:        8,
		CohortGCPeriodSeconds:   10.0,
		WriteInfoChanelSize:     1024,
		ReadInfoChanelSize:      1024,
		CohortSize:              1024,
		TTLHours:                1.0,
		ReadSafetyMarginMinutes: 5.0,
		Seed:                    1337,
		RandomPoolSize:          units.GiB,
	}
}

// LoadConfig loads the benchmark configuration from the json file at the given path.
func LoadConfig(path string) (*BenchmarkConfig, error) {
	config := DefaultBenchmarkConfig()

	// Expand home directory if path contains ~
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}
		path = strings.Replace(path, "~", homeDir, 1)
	}

	// Resolve relative paths to absolute paths
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	// Read the file
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Create a decoder that will return an error if there are unmatched fields
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.DisallowUnknownFields()

	// Unmarshal JSON into config struct
	err = decoder.Decode(config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return config, nil
}
