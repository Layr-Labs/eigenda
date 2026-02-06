package limiter

// Config is the configuration for the relay rate limiting.
type Config struct {

	// Blob rate limiting

	// Max number of GetBlob operations per second
	MaxGetBlobOpsPerSecond float64

	// Burstiness of the GetBlob rate limiter
	GetBlobOpsBurstiness int

	// Max bandwidth for GetBlob operations in bytes per second
	MaxGetBlobBytesPerSecond float64

	// Burstiness of the GetBlob bandwidth rate limiter
	GetBlobBytesBurstiness int

	// Max number of concurrent GetBlob operations
	MaxConcurrentGetBlobOps int

	// Chunk rate limiting

	// Max number of GetChunk operations per second
	MaxGetChunkOpsPerSecond float64

	// Burstiness of the GetChunk rate limiter
	GetChunkOpsBurstiness int

	// Max bandwidth for GetChunk operations in bytes per second
	MaxGetChunkBytesPerSecond float64

	// Burstiness of the GetChunk bandwidth rate limiter
	GetChunkBytesBurstiness int

	// Max number of concurrent GetChunk operations
	MaxConcurrentGetChunkOps int

	// Client rate limiting for GetChunk operations

	// Max number of GetChunk operations per second per client
	MaxGetChunkOpsPerSecondClient float64

	// Burstiness of the GetChunk rate limiter per client
	GetChunkOpsBurstinessClient int

	// Max bandwidth for GetChunk operations in bytes per second per client
	MaxGetChunkBytesPerSecondClient float64

	// Burstiness of the GetChunk bandwidth rate limiter per client
	GetChunkBytesBurstinessClient int

	// Max number of concurrent GetChunk operations per client
	MaxConcurrentGetChunkOpsClient int
}
