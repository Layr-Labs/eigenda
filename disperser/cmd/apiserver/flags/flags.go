package flags

import (
	"runtime"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/ratelimit"
	"github.com/Layr-Labs/eigenda/disperser/apiserver"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
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
	EnablePaymentMeterer = cli.BoolFlag{
		Name:   common.PrefixFlag(FlagPrefix, "enable-payment-meterer"),
		Usage:  "enable payment meterer",
		EnvVar: common.PrefixEnvVar(envVarPrefix, "ENABLE_PAYMENT_METERER"),
	}
	EnableRatelimiter = cli.BoolFlag{
		Name:   common.PrefixFlag(FlagPrefix, "enable-ratelimiter"),
		Usage:  "enable rate limiter",
		EnvVar: common.PrefixEnvVar(envVarPrefix, "ENABLE_RATELIMITER"),
	}
	ReservationsTableName = cli.StringFlag{
		Name:   common.PrefixFlag(FlagPrefix, "reservations-table-name"),
		Usage:  "name of the dynamodb table to store reservation usages",
		Value:  "reservations",
		EnvVar: common.PrefixEnvVar(envVarPrefix, "RESERVATIONS_TABLE_NAME"),
	}
	OnDemandTableName = cli.StringFlag{
		Name:   common.PrefixFlag(FlagPrefix, "on-demand-table-name"),
		Usage:  "name of the dynamodb table to store on-demand payments",
		Value:  "on_demand",
		EnvVar: common.PrefixEnvVar(envVarPrefix, "ON_DEMAND_TABLE_NAME"),
	}
	GlobalRateTableName = cli.StringFlag{
		Name:   common.PrefixFlag(FlagPrefix, "global-rate-table-name"),
		Usage:  "name of the dynamodb table to store global rate usage. If not provided, a local store will be used",
		Value:  "global_rate",
		EnvVar: common.PrefixEnvVar(envVarPrefix, "GLOBAL_RATE_TABLE_NAME"),
	}
	UpdateInterval = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "update-interval"),
		Usage:    "update interval for refreshing the on-chain state",
		Value:    1 * time.Second,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "UPDATE_INTERVAL"),
		Required: false,
	}
	ChainReadTimeout = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "chain-read-timeout"),
		Usage:    "timeout for reading from the chain",
		Value:    10,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "CHAIN_READ_TIMEOUT"),
		Required: false,
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
	OnchainStateRefreshInterval = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "onchain-state-refresh-interval"),
		Usage:    "The interval at which to refresh the onchain state. This flag is only relevant in v2",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "ONCHAIN_STATE_REFRESH_INTERVAL"),
		Value:    1 * time.Hour,
	}
	MaxNumSymbolsPerBlob = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "max-num-symbols-per-blob"),
		Usage:    "max number of symbols per blob. This flag is only relevant in v2",
		Value:    65_536,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "MAX_NUM_SYMBOLS_PER_BLOB"),
		Required: false,
	}
)

var kzgFlags = []cli.Flag{
	// KZG flags for encoding
	// These are copied from encoding/kzg/cli.go as optional flags for compatibility between v1 and v2 dispersers
	// These flags are only used in v2 disperser
	cli.StringFlag{
		Name:     kzg.G1PathFlagName,
		Usage:    "Path to G1 SRS",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "G1_PATH"),
	},
	cli.StringFlag{
		Name:     kzg.G2PathFlagName,
		Usage:    "Path to G2 SRS. Either this flag or G2_POWER_OF_2_PATH needs to be specified. For operator node, if both are specified, the node uses G2_POWER_OF_2_PATH first, if failed then tries to G2_PATH",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "G2_PATH"),
	},
	cli.StringFlag{
		Name:     kzg.CachePathFlagName,
		Usage:    "Path to SRS Table directory",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "CACHE_PATH"),
	},
	cli.Uint64Flag{
		Name:     kzg.SRSOrderFlagName,
		Usage:    "Order of the SRS",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "SRS_ORDER"),
	},
	cli.Uint64Flag{
		Name:     kzg.SRSLoadingNumberFlagName,
		Usage:    "Number of SRS points to load into memory",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "SRS_LOAD"),
	},
	cli.Uint64Flag{
		Name:     kzg.NumWorkerFlagName,
		Usage:    "Number of workers for multithreading",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "NUM_WORKERS"),
		Value:    uint64(runtime.GOMAXPROCS(0)),
	},
	cli.BoolFlag{
		Name:     kzg.VerboseFlagName,
		Usage:    "Enable to see verbose output for encoding/decoding",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "VERBOSE"),
	},
	cli.BoolFlag{
		Name:     kzg.CacheEncodedBlobsFlagName,
		Usage:    "Enable to cache encoded results",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "CACHE_ENCODED_BLOBS"),
	},
	cli.BoolFlag{
		Name:     kzg.PreloadEncoderFlagName,
		Usage:    "Set to enable Encoder PreLoading",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "PRELOAD_ENCODER"),
	},
	cli.StringFlag{
		Name:     kzg.G2PowerOf2PathFlagName,
		Usage:    "Path to G2 SRS points that are on power of 2. Either this flag or G2_PATH needs to be specified. For operator node, if both are specified, the node uses G2_POWER_OF_2_PATH first, if failed then tries to G2_PATH",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "G2_POWER_OF_2_PATH"),
	},
}

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
	EnablePaymentMeterer,
	BucketStoreSize,
	GrpcTimeoutFlag,
	MaxBlobSize,
	ReservationsTableName,
	OnDemandTableName,
	GlobalRateTableName,
	OnchainStateRefreshInterval,
	MaxNumSymbolsPerBlob,
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
	Flags = append(Flags, kzgFlags...)
}
