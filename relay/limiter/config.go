package limiter

// Config is the configuration for the relay rate limiting.
type Config struct {

	// Blob rate limiting

	// MaxGetBlobOpsPerSecond is the maximum permitted number of GetBlob operations per second. Default is
	// 1024.
	MaxGetBlobOpsPerSecond float64
	// The burstiness of the MaxGetBlobOpsPerSecond rate limiter. This is the maximum burst size that happen within
	// a short time window. Default is 1024.
	GetBlobOpsBurstiness int

	// MaxGetBlobBytesPerSecond is the maximum bandwidth, in bytes, that GetBlob operations are permitted
	// to consume per second. Default is 20MiB/s.
	MaxGetBlobBytesPerSecond float64
	// The burstiness of the MaxGetBlobBytesPerSecond rate limiter. This is the maximum burst size that happen within
	// a short time window. Default is 20MiB.
	GetBlobBytesBurstiness int

	// MaxConcurrentGetBlobOps is the maximum number of concurrent GetBlob operations that are permitted.
	// This is in addition to the rate limits. Default is 1024.
	MaxConcurrentGetBlobOps int

	// Chunk rate limiting

	// MaxGetChunkOpsPerSecond is the maximum permitted number of GetChunk operations per second. Default is
	// 1024.
	MaxGetChunkOpsPerSecond float64
	// The burstiness of the MaxGetChunkOpsPerSecond rate limiter. This is the maximum burst size that happen within
	// a short time window. Default is 1024.
	GetChunkOpsBurstiness int

	// MaxGetChunkBytesPerSecond is the maximum bandwidth, in bytes, that GetChunk operations are permitted
	// to consume per second. Default is 20MiB/s.
	MaxGetChunkBytesPerSecond float64
	// The burstiness of the MaxGetChunkBytesPerSecond rate limiter. This is the maximum burst size that happen within
	// a short time window. Default is 20MiB.
	GetChunkBytesBurstiness int

	// MaxConcurrentGetChunkOps is the maximum number of concurrent GetChunk operations that are permitted.
	// Default is 1024.
	MaxConcurrentGetChunkOps int

	// Client rate limiting for GetChunk operations

	// MaxGetChunkOpsPerSecondClient is the maximum permitted number of GetChunk operations per second for a single
	// client. Default is 8.
	MaxGetChunkOpsPerSecondClient float64
	// The burstiness of the MaxGetChunkOpsPerSecondClient rate limiter. This is the maximum burst size that happen
	// within a short time window. Default is 8.
	GetChunkOpsBurstinessClient int

	// MaxGetChunkBytesPerSecondClient is the maximum bandwidth, in bytes, that GetChunk operations are permitted
	// to consume per second. Default is 2MiB/s.
	MaxGetChunkBytesPerSecondClient float64
	// The burstiness of the MaxGetChunkBytesPerSecondClient rate limiter. This is the maximum burst size that happen
	// within a short time window. Default is 2MiB.
	GetChunkBytesBurstinessClient int

	// MaxConcurrentGetChunkOpsClient is the maximum number of concurrent GetChunk operations that are permitted.
	// Default is 1.
	MaxConcurrentGetChunkOpsClient int
}
