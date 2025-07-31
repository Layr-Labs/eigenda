package flags

import (
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/docker/go-units"
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
	RelayKeysFlag = cli.IntSliceFlag{
		Name:     common.PrefixFlag(FlagPrefix, "relay-keys"),
		Usage:    "Relay keys to use",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "RELAY_KEYS"),
	}
	MaxGRPCMessageSizeFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "max-grpc-message-size"),
		Usage:    "Max size of a gRPC message in bytes",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "MAX_GRPC_MESSAGE_SIZE"),
		Value:    4 * units.MiB,
	}
	MetadataCacheSizeFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "metadata-cache-size"),
		Usage:    "Max number of items in the metadata cache",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "METADATA_CACHE_SIZE"),
		Value:    units.MiB,
	}
	MetadataMaxConcurrencyFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "metadata-max-concurrency"),
		Usage:    "Max number of concurrent metadata fetches",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "METADATA_MAX_CONCURRENCY"),
		Value:    32,
	}
	BlobCacheBytes = cli.Uint64Flag{
		Name:     common.PrefixFlag(FlagPrefix, "blob-cache-bytes"),
		Usage:    "The size of the blob cache, in bytes.",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "BLOB_CACHE_SIZE"),
		Value:    units.GiB,
	}
	BlobMaxConcurrencyFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "blob-max-concurrency"),
		Usage:    "Max number of concurrent blob fetches",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "BLOB_MAX_CONCURRENCY"),
		Value:    32,
	}
	ChunkCacheBytesFlag = cli.Int64Flag{
		Name:     common.PrefixFlag(FlagPrefix, "chunk-cache-bytes"),
		Usage:    "Size of the chunk cache, in bytes.",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "CHUNK_CACHE_BYTES"),
		Value:    units.GiB,
	}
	ChunkMaxConcurrencyFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "chunk-max-concurrency"),
		Usage:    "Max number of concurrent chunk fetches",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "CHUNK_MAX_CONCURRENCY"),
		Value:    32,
	}
	MaxKeysPerGetChunksRequestFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "max-keys-per-get-chunks-request"),
		Usage:    "Max number of keys to fetch in a single GetChunks request",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "MAX_KEYS_PER_GET_CHUNKS_REQUEST"),
		Value:    1024,
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
		Value:    20 * units.MiB,
	}
	GetBlobBytesBurstinessFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "get-blob-bytes-burstiness"),
		Usage:    "Burstiness of the GetBlob bandwidth rate limiter",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "GET_BLOB_BYTES_BURSTINESS"),
		Value:    20 * units.MiB,
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
		Value:    80 * units.MiB,
	}
	GetChunkBytesBurstinessFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "get-chunk-bytes-burstiness"),
		Usage:    "Burstiness of the GetChunk bandwidth rate limiter",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "GET_CHUNK_BYTES_BURSTINESS"),
		Value:    800 * units.MiB,
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
		Value:    40 * units.MiB,
	}
	GetChunkBytesBurstinessClientFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "get-chunk-bytes-burstiness-client"),
		Usage:    "Burstiness of the GetChunk bandwidth rate limiter per client",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "GET_CHUNK_BYTES_BURSTINESS_CLIENT"),
		Value:    400 * units.MiB,
	}
	MaxConcurrentGetChunkOpsClientFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "max-concurrent-get-chunk-ops-client"),
		Usage:    "Max number of concurrent GetChunk operations per client",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "MAX_CONCURRENT_GET_CHUNK_OPS_CLIENT"),
		Value:    1,
	}
	BlsOperatorStateRetrieverAddrFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "bls-operator-state-retriever-addr"),
		Usage:    "[Deprecated: use EigenDADirectory instead] Address of the BLS operator state retriever",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "BLS_OPERATOR_STATE_RETRIEVER_ADDR"),
	}
	EigenDAServiceManagerAddrFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "eigen-da-service-manager-addr"),
		Usage:    "Address of the Eigen DA service manager",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "EIGEN_DA_SERVICE_MANAGER_ADDR"),
	}
	EigenDADirectoryFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "eigenda-directory"),
		Usage:    "Address of the EigenDA directory contract, which points to all other EigenDA contract addresses. This is the only contract entrypoint needed offchain.",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "EIGENDA_DIRECTORY"),
	}
	AuthenticationKeyCacheSizeFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "authentication-key-cache-size"),
		Usage:    "Max number of items in the authentication key cache",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "AUTHENTICATION_KEY_CACHE_SIZE"),
		Value:    1024 * 1024,
	}
	AuthenticationTimeoutFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "authentication-timeout"),
		Usage:    "Duration to keep authentication results",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "AUTHENTICATION_TIMEOUT"),
		Value:    0, // TODO(cody-littley) remove this feature
	}
	AuthenticationDisabledFlag = cli.BoolFlag{
		Name:     common.PrefixFlag(FlagPrefix, "authentication-disabled"),
		Usage:    "Disable GetChunks() authentication",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "AUTHENTICATION_DISABLED"),
	}
	GetChunksTimeoutFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "get-chunks-timeout"),
		Usage:    "Timeout for GetChunks()",
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "GET_CHUNKS_TIMEOUT"),
		Required: false,
		Value:    20 * time.Second,
	}
	GetBlobTimeoutFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "get-blob-timeout"),
		Usage:    "Timeout for GetBlob()",
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "GET_BLOB_TIMEOUT"),
		Required: false,
		Value:    20 * time.Second,
	}
	InternalGetMetadataTimeoutFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "internal-get-metadata-timeout"),
		Usage:    "Timeout for internal metadata fetch",
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "INTERNAL_GET_METADATA_TIMEOUT"),
		Required: false,
		Value:    5 * time.Second,
	}
	InternalGetBlobTimeoutFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "internal-get-blob-timeout"),
		Usage:    "Timeout for internal blob fetch",
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "INTERNAL_GET_BLOB_TIMEOUT"),
		Required: false,
		Value:    20 * time.Second,
	}
	InternalGetProofsTimeoutFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "internal-get-proofs-timeout"),
		Usage:    "Timeout for internal proofs fetch",
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "INTERNAL_GET_PROOFS_TIMEOUT"),
		Required: false,
		Value:    5 * time.Second,
	}
	InternalGetCoefficientsTimeoutFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "internal-get-coefficients-timeout"),
		Usage:    "Timeout for internal coefficients fetch",
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "INTERNAL_GET_COEFFICIENTS_TIMEOUT"),
		Required: false,
		Value:    20 * time.Second,
	}
	OnchainStateRefreshIntervalFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "onchain-state-refresh-interval"),
		Usage:    "The interval at which to refresh the onchain state",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "ONCHAIN_STATE_REFRESH_INTERVAL"),
		Value:    1 * time.Hour,
	}
	MetricsPortFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "metrics-port"),
		Usage:    "Port to listen on for metrics",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "METRICS_PORT"),
		Value:    9101,
	}
	EnableMetricsFlag = cli.BoolFlag{
		Name:     common.PrefixFlag(FlagPrefix, "enable-metrics"),
		Usage:    "Enable prometheus metrics collection",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "ENABLE_METRICS"),
	}
	EnablePprofFlag = cli.BoolFlag{
		Name:     common.PrefixFlag(FlagPrefix, "enable-pprof"),
		Usage:    "Enable pprof profiling",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "ENABLE_PPROF"),
	}
	PprofHttpPortFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "pprof-port"),
		Usage:    "Port to listen on for pprof",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "PPROF_PORT"),
		Value:    6060,
	}
	GetChunksRequestMaxPastAgeFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "get-chunks-request-max-past-age"),
		Usage:    "Max age of a GetChunks request",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "GET_CHUNKS_REQUEST_MAX_PAST_AGE"),
		Value:    5 * time.Minute,
	}
	GetChunksRequestMaxFutureAgeFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "get-chunks-request-max-future-age"),
		Usage:    "Max future age of a GetChunks request",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "GET_CHUNKS_REQUEST_MAX_FUTURE_AGE"),
		Value:    5 * time.Minute,
	}
	MaxConnectionAgeFlag = cli.DurationFlag{
		Name: common.PrefixFlag(FlagPrefix, "max-connection-age"),
		Usage: "Maximum age of a gRPC connection before it is closed. " +
			"If zero, then the server will not close connections based on age.",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "MAX_CONNECTION_AGE_SECONDS"),
		Value:    5 * time.Minute,
	}
	MaxConnectionAgeGraceFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "max-connection-age-grace"),
		Usage:    "Grace period after MaxConnectionAge before the connection is forcibly closed.",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "MAX_CONNECTION_AGE_GRACE_SECONDS"),
		Value:    30 * time.Second,
	}
	MaxIdleConnectionAgeFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "max-idle-connection-age"),
		Usage:    "Maximum time a connection can be idle before it is closed.",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "MAX_IDLE_CONNECTION_AGE_SECONDS"),
		Value:    time.Minute,
	}
)

var requiredFlags = []cli.Flag{
	GRPCPortFlag,
	BucketNameFlag,
	MetadataTableNameFlag,
	RelayKeysFlag,
	EnableMetricsFlag,
}

var optionalFlags = []cli.Flag{
	MaxGRPCMessageSizeFlag,
	MetadataCacheSizeFlag,
	MetadataMaxConcurrencyFlag,
	BlobCacheBytes,
	BlobMaxConcurrencyFlag,
	ChunkCacheBytesFlag,
	ChunkMaxConcurrencyFlag,
	MaxKeysPerGetChunksRequestFlag,
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
	AuthenticationKeyCacheSizeFlag,
	AuthenticationTimeoutFlag,
	AuthenticationDisabledFlag,
	GetChunksTimeoutFlag,
	GetBlobTimeoutFlag,
	InternalGetMetadataTimeoutFlag,
	InternalGetBlobTimeoutFlag,
	InternalGetProofsTimeoutFlag,
	InternalGetCoefficientsTimeoutFlag,
	OnchainStateRefreshIntervalFlag,
	MetricsPortFlag,
	EnablePprofFlag,
	PprofHttpPortFlag,
	GetChunksRequestMaxPastAgeFlag,
	GetChunksRequestMaxFutureAgeFlag,
	EigenDADirectoryFlag,
	BlsOperatorStateRetrieverAddrFlag,
	EigenDAServiceManagerAddrFlag,
	MaxConnectionAgeFlag,
	MaxConnectionAgeGraceFlag,
	MaxIdleConnectionAgeFlag,
}

var Flags []cli.Flag

func init() {
	Flags = append(requiredFlags, optionalFlags...)
	Flags = append(Flags, common.LoggerCLIFlags(envVarPrefix, FlagPrefix)...)
	Flags = append(Flags, aws.ClientFlags(envVarPrefix, FlagPrefix)...)
	Flags = append(Flags, geth.EthClientFlags(envVarPrefix)...)
	Flags = append(Flags, thegraph.CLIFlags(envVarPrefix)...)
}
