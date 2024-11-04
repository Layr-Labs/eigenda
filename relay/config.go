package relay

// Config is the configuration for the relay Server.
type Config struct {

	// MetadataCacheSize is the size of the metadata cache. Default is 1024 * 1024.
	MetadataCacheSize int

	MaxGetBlobsSize int

	MaxGetChunksSize int
}

func DefaultConfig() *Config {
	return &Config{
		MetadataCacheSize: 1024 * 1024,
	}
}
