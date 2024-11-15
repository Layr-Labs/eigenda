package limiter

// Config is the configuration for the relay rate limiting.
type Config struct {

	// Blob rate limiting

	// MaxGetBlobOpsPerSecond is the maximum permitted number of GetBlob operations per second. Default is
	// 1024.
	MaxGetBlobOpsPerSecond float64
	// MaxGetBlobBytesPerSecond is the maximum bandwidth, in bytes, that GetBlob operations are permitted
	// to consume per second. Default is 20MiB/s.
	MaxGetBlobBytesPerSecond float64
	// MaxConcurrentGetBlobOps is the maximum number of concurrent GetBlob operations that are permitted.
	// This is in addition to the rate limits. Default is 1024.
	MaxConcurrentGetBlobOps int

	// Chunk rate limiting

	// MaxGetChunkOpsPerSecond is the maximum permitted number of GetChunk operations per second. Default is
	// 1024.
	MaxGetChunkOpsPerSecond float64
	// MaxGetChunkBytesPerSecond is the maximum bandwidth, in bytes, that GetChunk operations are permitted
	// to consume per second. Default is 20MiB/s.
	MaxGetChunkBytesPerSecond float64
	// MaxConcurrentGetChunkOps is the maximum number of concurrent GetChunk operations that are permitted.
	// Default is 1024.
	MaxConcurrentGetChunkOps int

	// Client rate limiting for GetChunk operations

	// MaxGetChunkOpsPerSecondClient is the maximum permitted number of GetChunk operations per second for a single
	// client. Default is 8.
	MaxGetChunkOpsPerSecondClient float64
	// MaxGetChunkBytesPerSecondClient is the maximum bandwidth, in bytes, that GetChunk operations are permitted
	// to consume per second. Default is 2MiB/s.
	MaxGetChunkBytesPerSecondClient float64
	// MaxConcurrentGetChunkOpsClient is the maximum number of concurrent GetChunk operations that are permitted.
	// Default is 1.
	MaxConcurrentGetChunkOpsClient int
}

// DefaultConfig returns a default rate limit configuration.
func DefaultConfig() *Config {
	return &Config{
		MaxGetBlobOpsPerSecond:   1024,
		MaxGetBlobBytesPerSecond: 20 * 1024 * 1024,
		MaxConcurrentGetBlobOps:  1024,

		MaxGetChunkOpsPerSecond:   1024,
		MaxGetChunkBytesPerSecond: 20 * 1024 * 1024,
		MaxConcurrentGetChunkOps:  1024,

		MaxGetChunkOpsPerSecondClient:   8,
		MaxGetChunkBytesPerSecondClient: 2 * 1024 * 1024,
		MaxConcurrentGetChunkOpsClient:  1,
	}
}
