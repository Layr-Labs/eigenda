package flags

import (
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/ratelimit"
	"github.com/Layr-Labs/eigenda/disperser/apiserver"
	"github.com/urfave/cli"
)

const (
	FlagPrefix   = "disperser-server"
	envVarPrefix = "DISPERSER_SERVER"
)

var (
	/* Required Flags */
	S3BucketNameFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "s3-bucket-name"),
		Usage:    "Name of the bucket to store blobs",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "S3_BUCKET_NAME"),
	}
	DynamoDBTableNameFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "dynamodb-table-name"),
		Usage:    "Name of the dynamodb table to store blob metadata",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "DYNAMODB_TABLE_NAME"),
	}
	GrpcPortFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "grpc-port"),
		Usage:    "Port at which disperser listens for grpc calls",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "GRPC_PORT"),
	}
	GrpcTimeoutFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "grpc-stream-timeout"),
		Usage:    "Timeout for grpc streams",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "GRPC_STREAM_TIMEOUT"),
		Value:    time.Second * 10,
	}
	BlsOperatorStateRetrieverFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "bls-operator-state-retriever"),
		Usage:    "Address of the BLS Operator State Retriever",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "BLS_OPERATOR_STATE_RETRIVER"),
	}
	EigenDAServiceManagerFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "eigenda-service-manager"),
		Usage:    "Address of the EigenDA Service Manager",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "EIGENDA_SERVICE_MANAGER"),
	}
	/* Optional Flags*/
	DisperserVersionFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "disperser-version"),
		Usage:    "Disperser version. Options are 1 and 2.",
		Required: false,
		Value:    1,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "DISPERSER_VERSION"),
	}
	MetricsHTTPPort = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "metrics-http-port"),
		Usage:    "the http port which the metrics prometheus server is listening",
		Required: false,
		Value:    "9100",
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "METRICS_HTTP_PORT"),
	}
	EnableMetrics = cli.BoolFlag{
		Name:     common.PrefixFlag(FlagPrefix, "enable-metrics"),
		Usage:    "start metrics server",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "ENABLE_METRICS"),
	}
	EnableRatelimiter = cli.BoolFlag{
		Name:   common.PrefixFlag(FlagPrefix, "enable-ratelimiter"),
		Usage:  "enable rate limiter",
		EnvVar: common.PrefixEnvVar(envVarPrefix, "ENABLE_RATELIMITER"),
	}
	BucketTableName = cli.StringFlag{
		Name:   common.PrefixFlag(FlagPrefix, "rate-bucket-table-name"),
		Usage:  "name of the dynamodb table to store rate limiter buckets. If not provided, a local store will be used",
		Value:  "",
		EnvVar: common.PrefixEnvVar(envVarPrefix, "RATE_BUCKET_TABLE_NAME"),
	}
	BucketStoreSize = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "rate-bucket-store-size"),
		Usage:    "size (max number of entries) of the local store to use for rate limiting buckets",
		Value:    100_000,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "RATE_BUCKET_STORE_SIZE"),
		Required: false,
	}
	MaxBlobSize = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "max-blob-size"),
		Usage:    "max blob size disperser is accepting",
		Value:    2_097_152,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "MAX_BLOB_SIZE"),
		Required: false,
	}
)

var requiredFlags = []cli.Flag{
	S3BucketNameFlag,
	DynamoDBTableNameFlag,
	GrpcPortFlag,
	BucketTableName,
	BlsOperatorStateRetrieverFlag,
	EigenDAServiceManagerFlag,
}

var optionalFlags = []cli.Flag{
	DisperserVersionFlag,
	MetricsHTTPPort,
	EnableMetrics,
	EnableRatelimiter,
	BucketStoreSize,
	GrpcTimeoutFlag,
	MaxBlobSize,
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func init() {
	Flags = append(requiredFlags, optionalFlags...)
	Flags = append(Flags, geth.EthClientFlags(envVarPrefix)...)
	Flags = append(Flags, common.LoggerCLIFlags(envVarPrefix, FlagPrefix)...)
	Flags = append(Flags, ratelimit.RatelimiterCLIFlags(envVarPrefix, FlagPrefix)...)
	Flags = append(Flags, aws.ClientFlags(envVarPrefix, FlagPrefix)...)
	Flags = append(Flags, apiserver.CLIFlags(envVarPrefix)...)
}
