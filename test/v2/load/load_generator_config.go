package load

import (
	"time"
)

// LoadGeneratorConfig is the configuration for the load generator.
type LoadGeneratorConfig struct {
	// The desired number of megabytes bytes per second to write.
	MBPerSecond float64
	// The average size of the blobs to write, in megabytes.
	AverageBlobSizeMB float64
	// The standard deviation of the blob size, in megabytes.
	BlobSizeStdDev float64
	// By default, this utility reads each blob back from each relay once. The number of
	// reads per relay is multiplied by this factor. For example, If this is set to 3,
	// then each blob is read back from each relay 3 times. If less than 1, then this value
	// is treated as a probability. For example, if this is set to 0.5, then each blob is read back
	// from each relay with a 50% chance.
	RelayReadAmplification float64
	// By default, this utility reads chunks once. The number of chunk reads is multiplied
	// by this factor. If this is set to 3, then chunks are read back 3 times. If less than 1,
	// then this value is treated as a probability. For example, if this is set to 0.5, then
	// each chunk is read back from validators with a 50% chance.
	ValidatorReadAmplification float64
	// The maximum number of parallel blobs submissions in flight.
	SubmissionParallelism uint64
	// The maximum number of parallel blob read operations in flight.
	ReadParallelism uint64
	// The timeout for each blob dispersal.
	DispersalTimeout time.Duration
	// The timeout for reading each blob.
	ReadTimeout time.Duration
	// EnablePprof enables the pprof HTTP server for profiling
	EnablePprof bool
	// PprofHttpPort is the port that the pprof HTTP server listens on
	PprofHttpPort int
}
