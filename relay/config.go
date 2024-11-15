package relay

import (
	"fmt"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
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
// Environment variables are mapped into this struct by taking the name of the field in this struct,
// converting to upper case, and prepending "RELAY_". For example, "BlobCacheSize" can be set using the
// environment variable "RELAY_BLOBCACHESIZE".
//
// For nested structs, add the name of the struct variable before the field name, separated by an underscore.
// For example, "Log.Format" can be set using the environment variable "RELAY_LOG_FORMAT".
//
// Slice values can be set using a comma-separated list. For example, "RelayIDs" can be set using the environment
// variable "RELAY_RELAYIDS='1,2,3,4'".
//
// It is also possible to set the configuration using a configuration file. The path to the configuration file should
// be passed as the first argument to the relay binary, e.g. "bin/relay config.yaml". The structure of the config
// file should mirror the structure of this struct, with keys in the config file matching the field names
// of this struct.
type Config struct {

	// Log is the configuration for the logger. Default is common.DefaultLoggerConfig().
	Log common.LoggerConfig

	// Configuration for the AWS client. Default is aws.DefaultClientConfig().
	AWS aws.ClientConfig

	// RelayIDs contains the IDs of the relays that this server is willing to serve data for. If empty, the server will
	// serve data for any shard it can.
	RelayIDs []core.RelayKey

	// GRPCPort is the port that the relay server listens on. Default is 50051. // TODO what is a good port?
	GRPCPort int

	// BucketName is the name of the S3 bucket that stores blobs. Default is "relay".
	BucketName string

	// MetadataTableName is the name of the DynamoDB table that stores metadata. Default is "metadata".
	MetadataTableName string

	// MaxGRPCMessageSize is the maximum size of a gRPC message that the server will accept.
	// Default is 1024 * 1024 * 300 (300 MiB).
	MaxGRPCMessageSize int

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
}

// DefaultConfig returns the default configuration for the relay Server.
func DefaultConfig() *Config {
	return &Config{
		Log:                    common.DefaultLoggerConfig(),
		AWS:                    *aws.DefaultClientConfig(),
		GRPCPort:               50051,
		MaxGRPCMessageSize:     1024 * 1024 * 300,
		BucketName:             "relay",
		MetadataTableName:      "metadata",
		MetadataCacheSize:      1024 * 1024,
		MetadataMaxConcurrency: 32,
		BlobCacheSize:          32,
		BlobMaxConcurrency:     32,
		ChunkCacheSize:         32,
		ChunkMaxConcurrency:    32,
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
