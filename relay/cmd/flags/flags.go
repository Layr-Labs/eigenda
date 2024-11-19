package flags

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/urfave/cli"
)

const (
	FlagPrefix   = "relay"
	envVarPrefix = "RELAY"
)

var (
	GRPCPortFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "grpc-port"),
		Usage:    "Port to listen on for gRPC",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "GRPC_PORT"),
	}
	BucketNameFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "bucket-name"),
		Usage:    "Name of the s3 bucket to store blobs",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "BUCKET_NAME"),
	}
	MetadataTableNameFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "metadata-table-name"),
		Usage:    "Name of the dynamodb table to store blob metadata",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "METADATA_TABLE_NAME"),
	}
	RelayIDsFlag = cli.IntSliceFlag{
		Name:     common.PrefixFlag(FlagPrefix, "relay-ids"),
		Usage:    "Relay IDs to use",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "RELAY_IDS"),
	}
	MaxGRPCMessageSizeFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "max-grpc-message-size"),
		Usage:    "Max size of a gRPC message in bytes",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "MAX_GRPC_MESSAGE_SIZE"),
		Value:    1024 * 1024 * 300,
	}
	MetadataCacheSizeFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "metadata-cache-size"),
		Usage:    "Max number of items in the metadata cache",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "METADATA_CACHE_SIZE"),
		Value:    1024 * 1024,
	}
	MetadataMaxConcurrencyFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "metadata-max-concurrency"),
		Usage:    "Max number of concurrent metadata fetches",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "METADATA_MAX_CONCURRENCY"),
		Value:    32,
	}
	BlobCacheSizeFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "blob-cache-size"),
		Usage:    "Max number of items in the blob cache",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "BLOB_CACHE_SIZE"),
		Value:    32,
	}
	BlobMaxConcurrencyFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "blob-max-concurrency"),
		Usage:    "Max number of concurrent blob fetches",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "BLOB_MAX_CONCURRENCY"),
		Value:    32,
	}
	ChunkCacheSizeFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "chunk-cache-size"),
		Usage:    "Max number of items in the chunk cache",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "CHUNK_CACHE_SIZE"),
		Value:    32,
	}
	ChunkMaxConcurrencyFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "chunk-max-concurrency"),
		Usage:    "Max number of concurrent chunk fetches",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "CHUNK_MAX_CONCURRENCY"),
		Value:    32,
	}
)

var requiredFlags = []cli.Flag{
	GRPCPortFlag,
	BucketNameFlag,
	MetadataTableNameFlag,
	RelayIDsFlag,
}

var optionalFlags = []cli.Flag{
	MaxGRPCMessageSizeFlag,
	MetadataCacheSizeFlag,
	MetadataMaxConcurrencyFlag,
	BlobCacheSizeFlag,
	BlobMaxConcurrencyFlag,
	ChunkCacheSizeFlag,
	ChunkMaxConcurrencyFlag,
}

var Flags []cli.Flag

func init() {
	Flags = append(requiredFlags, optionalFlags...)
	Flags = append(Flags, common.LoggerCLIFlags(envVarPrefix, FlagPrefix)...)
	Flags = append(Flags, aws.ClientFlags(envVarPrefix, FlagPrefix)...)
}
