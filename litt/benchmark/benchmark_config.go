package benchmark

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Layr-Labs/eigenda/litt"
)

// BenchmarkConfig is a struct that holds the configuration for the benchmark.
type BenchmarkConfig struct {

	// Configuration for the LittDB instance.
	LittConfig *litt.Config

	// The location where the benchmark stores test metadata.
	MetadataDirectory string

	// The maximum target throughput in MB/s.
	MaximumThroughputMB float64

	// The average size of the values in MB.
	AverageValueSizeMB float64

	// The standard deviation of the size of the values in MB.
	ValueSizeStandardDeviationMB float64

	// Data is written to the DB in batches and then flushed. This determines the size of those batches, in MB.
	BatchSizeMB float64

	// The maximum number of batches permitted to be in-flight at any given time. If batches can't be written/flushed
	// fast enough, then the benchmark engine will not push data as fast as requested.
	BatchParallelism int
}

// DefaultBenchmarkConfig returns a default BenchmarkConfig with the given data paths.
func DefaultBenchmarkConfig() *BenchmarkConfig {
	return &BenchmarkConfig{
		LittConfig:                   litt.DefaultConfigNoPaths(),
		MetadataDirectory:            "~/benchmark",
		MaximumThroughputMB:          10.0,
		AverageValueSizeMB:           2.0,
		ValueSizeStandardDeviationMB: 0.25,
		BatchSizeMB:                  10,
		BatchParallelism:             10,
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
