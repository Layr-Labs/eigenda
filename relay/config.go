package relay

import core "github.com/Layr-Labs/eigenda/core/v2"

// Config is the configuration for the relay Server.
type Config struct {

	// Shards contains the IDs of the relays that this server is willing to serve data for. If empty, the server will
	// serve data for any shard it can.
	Shards []core.RelayKey

	// MetadataCacheSize is the size of the metadata cache. Default is 1024 * 1024.
	MetadataCacheSize int

	// MetadataWorkPoolSize is the size of the metadata work pool. Default is 32.
	MetadataWorkPoolSize int

	// BlobCacheSize is the size of the blob cache. Default is 32.
	BlobCacheSize int

	// BlobWorkPoolSize is the size of the blob work pool. Default is 32.
	BlobWorkPoolSize int

	// ChunkCacheSize is the size of the chunk cache. Default is 32.
	ChunkCacheSize int

	// ChunkWorkPoolSize is the size of the chunk work pool. Default is 32.
	ChunkWorkPoolSize int
}

// DefaultConfig returns the default configuration for the relay Server.
func DefaultConfig() *Config {
	return &Config{
		MetadataCacheSize:    1024 * 1024,
		MetadataWorkPoolSize: 32,
		BlobCacheSize:        32,
		BlobWorkPoolSize:     32,
		ChunkCacheSize:       32,
		ChunkWorkPoolSize:    32,
	}
}
