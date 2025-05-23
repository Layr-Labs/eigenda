package encodingload

// EncodingLoadGeneratorConfig is the configuration for the encoding load generator.
type EncodingLoadGeneratorConfig struct {
	// The desired number of megabytes bytes per second to encode.
	MBPerSecond float64
	// The average size of the blobs to encode, in megabytes.
	AverageBlobSizeMB float64
	// The standard deviation of the blob size, in megabytes.
	BlobSizeStdDev float64
	// The maximum blob size in bytes.
	MaxBlobSize float64
	// The maximum number of parallel blobs in flight.
	MaxParallelism uint64
	// EnablePprof enables the pprof HTTP server for profiling
	EnablePprof bool
	// PprofHttpPort is the port that the pprof HTTP server listens on
	PprofHttpPort int
}
