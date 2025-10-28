package relay

import (
	"time"

	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/relay/limiter"
)

// NewTestConfig creates a relay configuration suitable for testing.
// The relayIndex determines the relay key and metrics port.
// The grpcPort is set to 0 by default to let the OS assign a port (can be overridden).
func NewTestConfig(relayIndex int) *Config {
	return &Config{
		RelayKeys:                  []v2.RelayKey{v2.RelayKey(relayIndex)},
		GRPCPort:                   0, // OS assigns port
		MaxGRPCMessageSize:         1024 * 1024 * 300,
		MetadataCacheSize:          1024 * 1024,
		MetadataMaxConcurrency:     32,
		BlobCacheBytes:             32 * 1024 * 1024,
		BlobMaxConcurrency:         32,
		ChunkCacheBytes:            32 * 1024 * 1024,
		ChunkMaxConcurrency:        32,
		MaxKeysPerGetChunksRequest: 1024,
		RateLimits: limiter.Config{
			MaxGetBlobOpsPerSecond:          1024,
			GetBlobOpsBurstiness:            1024,
			MaxGetBlobBytesPerSecond:        20 * 1024 * 1024,
			GetBlobBytesBurstiness:          20 * 1024 * 1024,
			MaxConcurrentGetBlobOps:         1024,
			MaxGetChunkOpsPerSecond:         1024,
			GetChunkOpsBurstiness:           1024,
			MaxGetChunkBytesPerSecond:       20 * 1024 * 1024,
			GetChunkBytesBurstiness:         20 * 1024 * 1024,
			MaxConcurrentGetChunkOps:        1024,
			MaxGetChunkOpsPerSecondClient:   8,
			GetChunkOpsBurstinessClient:     8,
			MaxGetChunkBytesPerSecondClient: 2 * 1024 * 1024,
			GetChunkBytesBurstinessClient:   2 * 1024 * 1024,
			MaxConcurrentGetChunkOpsClient:  1,
		},
		AuthenticationKeyCacheSize:   1024,
		AuthenticationDisabled:       true, // Disabled for testing
		GetChunksRequestMaxPastAge:   5 * time.Minute,
		GetChunksRequestMaxFutureAge: 1 * time.Minute,
		Timeouts: TimeoutConfig{
			GetChunksTimeout:               20 * time.Second,
			GetBlobTimeout:                 20 * time.Second,
			InternalGetMetadataTimeout:     5 * time.Second,
			InternalGetBlobTimeout:         20 * time.Second,
			InternalGetProofsTimeout:       5 * time.Second,
			InternalGetCoefficientsTimeout: 20 * time.Second,
		},
		OnchainStateRefreshInterval: 10 * time.Second,
		MetricsPort:                 9100 + relayIndex,
		EnableMetrics:               true,
		EnablePprof:                 false,
		PprofHttpPort:               0,
		MaxConnectionAge:            0,
		MaxConnectionAgeGrace:       5 * time.Second,
		MaxIdleConnectionAge:        30 * time.Second,
	}
}
