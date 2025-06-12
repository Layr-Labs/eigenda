package load

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
	// A number between 0 and 1.0 that specifies the fraction of blobs that are verified by the validator.
	// If 1.0, all blobs are verified. If 0.0, no blobs are verified. If 0.5, half of the blobs are verified.
	ValidatorVerificationFraction float64
	// The maximum number of parallel blobs submissions in flight.
	SubmissionParallelism uint64
	// The maximum number of parallel blob relay read operations in flight.
	RelayReadParallelism uint64
	// The maximum number of parallel blob validator read operations in flight.
	ValidatorReadParallelism uint64
	// The timeout for each blob dispersal, in seconds.
	DispersalTimeout uint32
	// The timeout for reading a blob from a relay, in seconds. This is the timeout per individual read.
	RelayReadTimeout uint32
	// The timeout for reading a blob from the validators, in seconds. This is the timeout per individual read.
	ValidatorReadTimeout uint32
	// EnablePprof enables the pprof HTTP server for profiling
	EnablePprof bool
	// PprofHttpPort is the port that the pprof HTTP server listens on
	PprofHttpPort int
	// FrequencyAcceleration determines the speed at which the frequency of blob submissions accelerates at startup
	// time, in HZ/s. Frequency will start at 0 and accelerate to the target frequency at this rate. If 0, then
	// the frequency will immediately be set to the target frequency.
	FrequencyAcceleration float64
}

// DefaultLoadGeneratorConfig returns a default configuration for the load generator.
func DefaultLoadGeneratorConfig() *LoadGeneratorConfig {
	return &LoadGeneratorConfig{
		MBPerSecond:                   0.5,
		AverageBlobSizeMB:             1.0,
		BlobSizeStdDev:                0.0,
		RelayReadAmplification:        1.0,
		ValidatorReadAmplification:    1.0,
		ValidatorVerificationFraction: 0.01,
		SubmissionParallelism:         300,
		RelayReadParallelism:          300,
		ValidatorReadParallelism:      300,
		DispersalTimeout:              600,
		RelayReadTimeout:              600,
		ValidatorReadTimeout:          600,
		EnablePprof:                   false,
		PprofHttpPort:                 6060,
		FrequencyAcceleration:         0.0025,
	}
}
