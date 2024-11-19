package main

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/relay"
	"github.com/Layr-Labs/eigenda/relay/cmd/flags"
	"github.com/urfave/cli"
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

	// BucketName is the name of the S3 bucket that stores blobs. Default is "relay".
	BucketName string

	// MetadataTableName is the name of the DynamoDB table that stores metadata. Default is "metadata".
	MetadataTableName string

	// RelayConfig is the configuration for the relay.
	RelayConfig relay.Config
}

func NewConfig(ctx *cli.Context) (Config, error) {
	loggerConfig, err := common.ReadLoggerCLIConfig(ctx, flags.FlagPrefix)
	if err != nil {
		return Config{}, err
	}
	awsClientConfig := aws.ReadClientConfig(ctx, flags.FlagPrefix)
	relayIDs := ctx.IntSlice(flags.RelayIDsFlag.Name)
	if len(relayIDs) == 0 {
		return Config{}, fmt.Errorf("no relay IDs specified")
	}
	config := Config{
		Log:               *loggerConfig,
		AWS:               awsClientConfig,
		BucketName:        ctx.String(flags.BucketNameFlag.Name),
		MetadataTableName: ctx.String(flags.MetadataTableNameFlag.Name),
		RelayConfig: relay.Config{
			RelayIDs:               make([]core.RelayKey, len(relayIDs)),
			GRPCPort:               ctx.Int(flags.GRPCPortFlag.Name),
			MaxGRPCMessageSize:     ctx.Int(flags.MaxGRPCMessageSizeFlag.Name),
			MetadataCacheSize:      ctx.Int(flags.MetadataCacheSizeFlag.Name),
			MetadataMaxConcurrency: ctx.Int(flags.MetadataMaxConcurrencyFlag.Name),
			BlobCacheSize:          ctx.Int(flags.BlobCacheSizeFlag.Name),
			BlobMaxConcurrency:     ctx.Int(flags.BlobMaxConcurrencyFlag.Name),
			ChunkCacheSize:         ctx.Int(flags.ChunkCacheSizeFlag.Name),
			ChunkMaxConcurrency:    ctx.Int(flags.ChunkMaxConcurrencyFlag.Name),
		},
	}
	for i, id := range relayIDs {
		config.RelayConfig.RelayIDs[i] = core.RelayKey(id)
	}
	return config, nil
}
