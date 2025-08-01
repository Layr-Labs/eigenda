package relay

import (
	"time"

	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/relay/limiter"
)

// Config is the configuration for the relay Server.
type Config struct {

	// RelayKeys contains the keys of the relays that this server is willing to serve data for. If empty, the server will
	// serve data for any shard it can.
	RelayKeys []v2.RelayKey

	// GRPCPort is the port that the relay server listens on.
	GRPCPort int

	// MaxGRPCMessageSize is the maximum size of a gRPC message that the server will accept.
	MaxGRPCMessageSize int

	// MetadataCacheSize is the maximum number of items in the metadata cache.
	MetadataCacheSize int

	// MetadataMaxConcurrency puts a limit on the maximum number of concurrent metadata fetches actively running on
	// goroutines.
	MetadataMaxConcurrency int

	// BlobCacheBytes is the maximum size of the blob cache, in bytes.
	BlobCacheBytes uint64

	// BlobMaxConcurrency puts a limit on the maximum number of concurrent blob fetches actively running on goroutines.
	BlobMaxConcurrency int

	// ChunkCacheBytes is the maximum size of the chunk cache, in bytes.
	ChunkCacheBytes uint64

	// ChunkMaxConcurrency is the size of the work pool for fetching chunks. Note that this does not
	// impact concurrency utilized by the s3 client to upload/download fragmented files.
	ChunkMaxConcurrency int

	// MaxKeysPerGetChunksRequest is the maximum number of keys that can be requested in a single GetChunks request.
	MaxKeysPerGetChunksRequest int

	// RateLimits contains configuration for rate limiting.
	RateLimits limiter.Config

	// AuthenticationKeyCacheSize is the maximum number of operator public keys that can be cached.
	AuthenticationKeyCacheSize int

	// AuthenticationDisabled will disable authentication if set to true.
	AuthenticationDisabled bool

	// GetChunksRequestMaxPastAge is the maximum age of a GetChunks request that the server will accept.
	GetChunksRequestMaxPastAge time.Duration

	// GetChunksRequestMaxFutureAge is the maximum future age of a GetChunks request that the server will accept.
	GetChunksRequestMaxFutureAge time.Duration

	// Timeouts contains configuration for relay timeouts.
	Timeouts TimeoutConfig

	// OnchainStateRefreshInterval is the interval at which the onchain state is refreshed.
	OnchainStateRefreshInterval time.Duration

	// MetricsPort is the port that the relay metrics server listens on.
	MetricsPort int

	// EnableMetrics enables the metrics HTTP server for prometheus metrics collection
	EnableMetrics bool

	// EnablePprof enables the pprof HTTP server for profiling
	EnablePprof bool

	// PprofHttpPort is the port that the pprof HTTP server listens on
	PprofHttpPort int

	// The maximum permissible age of a GRPC connection before it is closed. If zero, then the server will not close
	// connections based on age.
	MaxConnectionAge time.Duration

	// When the server closes a connection due to MaxConnectionAgeSeconds, it will wait for this grace period before
	// forcibly closing the connection. This allows in-flight requests to complete.
	MaxConnectionAgeGrace time.Duration

	// MaxIdleConnectionAge is the maximum time a connection can be idle before it is closed.
	MaxIdleConnectionAge time.Duration
}
