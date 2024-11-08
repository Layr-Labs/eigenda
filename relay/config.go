package relay

import core "github.com/Layr-Labs/eigenda/core/v2"

// TODO make sure each of these is actually used

// Config is the configuration for the relay Server.
type Config struct {

	// Shards contains the IDs of the relays that this server is willing to serve data for. If empty, the server will
	// serve data for any shard it can.
	Shards []core.RelayKey

	// MaximumBlobKeyLimit defines the maximum number of blobs a request for blobs/chunks may touch. It does
	// not matter how little data is requested from each blob, the total number of blobs accessed by a single
	// request must not exceed this limit. Default is 32.
	MaximumBlobKeyLimit int

	// MaximumGetBlobsByteCount is the maximum number of bytes that can be requested in a single GetBlobs call.
	// Default is 1GB.
	MaximumGetBlobsByteCount int

	// MetadataCacheSize is the size of the metadata cache. Default is 1024 * 1024.
	MetadataCacheSize int

	// MetadataWorkPoolSize is the size of the metadata work pool. Default is 32.
	MetadataWorkPoolSize int

	// BlobCacheSize is the size of the blob cache. Default is 32.
	// TODO what is the largest blob we support? Is 32 too big?
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
		MaximumBlobKeyLimit:      32,
		MaximumGetBlobsByteCount: 1024 * 1024 * 1024,
		MetadataCacheSize:        1024 * 1024,
		MetadataWorkPoolSize:     32,
		BlobCacheSize:            32,
		BlobWorkPoolSize:         32,
		ChunkCacheSize:           32,
		ChunkWorkPoolSize:        32,
	}
}
