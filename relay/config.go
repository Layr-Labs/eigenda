package relay

import (
	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/relay/limiter"
)

// Config is the configuration for the relay Server.
type Config struct {

	// RelayIDs contains the IDs of the relays that this server is willing to serve data for. If empty, the server will
	// serve data for any shard it can.
	RelayIDs []core.RelayKey

	// MetadataCacheSize is the maximum number of items in the metadata cache. Default is 1024 * 1024.
	MetadataCacheSize int

	// MetadataMaxConcurrency puts a limit on the maximum number of concurrent metadata fetches actively running on
	// goroutines. Default is 32.
	MetadataMaxConcurrency int

	// BlobCacheSize is the maximum number of items in the blob cache. Default is 32.
	BlobCacheSize int

	// BlobMaxConcurrency puts a limit on the maximum number of concurrent blob fetches actively running on goroutines.
	// Default is 32.
	BlobMaxConcurrency int

	// ChunkCacheSize is the maximum number of items in the chunk cache. Default is 32.
	ChunkCacheSize int

	// ChunkMaxConcurrency is the size of the work pool for fetching chunks. Default is 32. Note that this does not
	// impact concurrency utilized by the s3 client to upload/download fragmented files.
	ChunkMaxConcurrency int

	// MaxKeysPerGetChunksRequest is the maximum number of keys that can be requested in a single GetChunks request.
	// Default is 1024. // TODO should this be the max batch size? What is that?
	MaxKeysPerGetChunksRequest int

	// RateLimits contains configuration for rate limiting.
	RateLimits limiter.Config
}

// DefaultConfig returns the default configuration for the relay Server.
func DefaultConfig() *Config {
	return &Config{
		MetadataCacheSize:          1024 * 1024,
		MetadataMaxConcurrency:     32,
		BlobCacheSize:              32,
		BlobMaxConcurrency:         32,
		ChunkCacheSize:             32,
		ChunkMaxConcurrency:        32,
		MaxKeysPerGetChunksRequest: 1024,
		RateLimits:                 *limiter.DefaultConfig(),
	}
}
