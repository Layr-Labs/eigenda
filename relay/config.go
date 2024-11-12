package relay

import (
	"fmt"
	"github.com/Layr-Labs/eigenda/common"
	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const (
	// relayEnvPrefix is the prefix for all environment variables used by the relay.
	relayEnvPrefix = "RELAY"
)

// Config is the configuration for the relay Server.
//
// Environment variables are mapped into this struct by taking ake the name of the field in this struct,
// converting to upper snake case, and prepending "RELAY_". For example, "BlobCacheSize" can be set using the
// environment variable "RELAY_BLOB_CACHE_SIZE".
//
// For nested structs, add the name of the struct variable before the field name, separated by a period. For example,
// "Log.Format" can be set using the environment variable "RELAY_LOG_FORMAT".
//
// Slice values can be set using a comma-separated list. For example, "Shards" can be set using the environment
// variable "RELAY_SHARDS='1,2,3,4'".
//
// It is also possible to set the configuration using a configuration file. The path to the configuration file should
// be passed as the first argument to the relay binary, e.g. "bin/relay config.yaml". The structure of the config
// file should mirror the structure of this struct, with keys in the config file matching the field names
// of this struct.
type Config struct {

	// Log is the configuration for the logger. Default is common.DefaultLoggerConfig().
	Log common.LoggerConfig

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
		Log:                  common.DefaultLoggerConfig(),
		MetadataCacheSize:    1024 * 1024,
		MetadataWorkPoolSize: 32,
		BlobCacheSize:        32,
		BlobWorkPoolSize:     32,
		ChunkCacheSize:       32,
		ChunkWorkPoolSize:    32,
	}
}

// LoadConfigWithViper loads the configuration for the relay Server using viper.
func LoadConfigWithViper() (*Config, error) {
	config := DefaultConfig()

	if len(os.Args) > 1 {
		viper.SetConfigFile(os.Args[1])
		err := viper.ReadInConfig()
		if err != nil {
			return nil, fmt.Errorf("fatal error reading config file: %s \n", err)
		}
	}

	viper.SetEnvPrefix(relayEnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))
	viper.AutomaticEnv()

	err := viper.Unmarshal(config)
	if err != nil {
		return nil, fmt.Errorf("fatal error unmarshaling configuration: %s \n", err)
	}

	return config, nil
}
