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
	MaxGetBlobOpsPerSecondFlag = cli.Float64Flag{
		Name:     common.PrefixFlag(FlagPrefix, "max-get-blob-ops-per-second"),
		Usage:    "Max number of GetBlob operations per second",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "MAX_GET_BLOB_OPS_PER_SECOND"),
		Value:    1024,
	}
	GetBlobOpsBurstinessFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "get-blob-ops-burstiness"),
		Usage:    "Burstiness of the GetBlob rate limiter",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "GET_BLOB_OPS_BURSTINESS"),
		Value:    1024,
	}
	MaxGetBlobBytesPerSecondFlag = cli.Float64Flag{
		Name:     common.PrefixFlag(FlagPrefix, "max-get-blob-bytes-per-second"),
		Usage:    "Max bandwidth for GetBlob operations in bytes per second",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "MAX_GET_BLOB_BYTES_PER_SECOND"),
		Value:    20 * 1024 * 1024,
	}
	GetBlobBytesBurstinessFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "get-blob-bytes-burstiness"),
		Usage:    "Burstiness of the GetBlob bandwidth rate limiter",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "GET_BLOB_BYTES_BURSTINESS"),
		Value:    20 * 1024 * 1024,
	}
	MaxConcurrentGetBlobOpsFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "max-concurrent-get-blob-ops"),
		Usage:    "Max number of concurrent GetBlob operations",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "MAX_CONCURRENT_GET_BLOB_OPS"),
		Value:    1024,
	}
	MaxGetChunkOpsPerSecondFlag = cli.Float64Flag{
		Name:     common.PrefixFlag(FlagPrefix, "max-get-chunk-ops-per-second"),
		Usage:    "Max number of GetChunk operations per second",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "MAX_GET_CHUNK_OPS_PER_SECOND"),
		Value:    1024,
	}
	GetChunkOpsBurstinessFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "get-chunk-ops-burstiness"),
		Usage:    "Burstiness of the GetChunk rate limiter",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "GET_CHUNK_OPS_BURSTINESS"),
		Value:    1024,
	}
	MaxGetChunkBytesPerSecondFlag = cli.Float64Flag{
		Name:     common.PrefixFlag(FlagPrefix, "max-get-chunk-bytes-per-second"),
		Usage:    "Max bandwidth for GetChunk operations in bytes per second",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "MAX_GET_CHUNK_BYTES_PER_SECOND"),
		Value:    20 * 1024 * 1024,
	}
	GetChunkBytesBurstinessFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "get-chunk-bytes-burstiness"),
		Usage:    "Burstiness of the GetChunk bandwidth rate limiter",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "GET_CHUNK_BYTES_BURSTINESS"),
		Value:    20 * 1024 * 1024,
	}
	MaxConcurrentGetChunkOpsFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "max-concurrent-get-chunk-ops"),
		Usage:    "Max number of concurrent GetChunk operations",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "MAX_CONCURRENT_GET_CHUNK_OPS"),
		Value:    1024,
	}
	MaxGetChunkOpsPerSecondClientFlag = cli.Float64Flag{
		Name:     common.PrefixFlag(FlagPrefix, "max-get-chunk-ops-per-second-client"),
		Usage:    "Max number of GetChunk operations per second per client",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "MAX_GET_CHUNK_OPS_PER_SECOND_CLIENT"),
		Value:    8,
	}
	GetChunkOpsBurstinessClientFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "get-chunk-ops-burstiness-client"),
		Usage:    "Burstiness of the GetChunk rate limiter per client",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "GET_CHUNK_OPS_BURSTINESS_CLIENT"),
		Value:    8,
	}
	MaxGetChunkBytesPerSecondClientFlag = cli.Float64Flag{
		Name:     common.PrefixFlag(FlagPrefix, "max-get-chunk-bytes-per-second-client"),
		Usage:    "Max bandwidth for GetChunk operations in bytes per second per client",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "MAX_GET_CHUNK_BYTES_PER_SECOND_CLIENT"),
		Value:    2 * 1024 * 1024,
	}
	GetChunkBytesBurstinessClientFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "get-chunk-bytes-burstiness-client"),
		Usage:    "Burstiness of the GetChunk bandwidth rate limiter per client",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "GET_CHUNK_BYTES_BURSTINESS_CLIENT"),
	}
	MaxConcurrentGetChunkOpsClientFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "max-concurrent-get-chunk-ops-client"),
		Usage:    "Max number of concurrent GetChunk operations per client",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "MAX_CONCURRENT_GET_CHUNK_OPS_CLIENT"),
		Value:    1,
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
	MaxGetBlobOpsPerSecondFlag,
	GetBlobOpsBurstinessFlag,
	MaxGetBlobBytesPerSecondFlag,
	GetBlobBytesBurstinessFlag,
	MaxConcurrentGetBlobOpsFlag,
	MaxGetChunkOpsPerSecondFlag,
	GetChunkOpsBurstinessFlag,
	MaxGetChunkBytesPerSecondFlag,
	GetChunkBytesBurstinessFlag,
	MaxConcurrentGetChunkOpsFlag,
	MaxGetChunkOpsPerSecondClientFlag,
	GetChunkOpsBurstinessClientFlag,
	MaxGetChunkBytesPerSecondClientFlag,
	GetChunkBytesBurstinessClientFlag,
	MaxConcurrentGetChunkOpsClientFlag,
}

var Flags []cli.Flag

func init() {
	Flags = append(requiredFlags, optionalFlags...)
	Flags = append(Flags, common.LoggerCLIFlags(envVarPrefix, FlagPrefix)...)
	Flags = append(Flags, aws.ClientFlags(envVarPrefix, FlagPrefix)...)
}
