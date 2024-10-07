package flags

import (
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigenda/indexer"
	"github.com/urfave/cli"
)

const (
	FlagPrefix   = "batcher"
	envVarPrefix = "BATCHER"
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
	PullIntervalFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "pull-interval"),
		Usage:    "Interval at which to pull from the queue",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "PULL_INTERVAL"),
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
	EncoderSocket = cli.StringFlag{
		Name:     "encoder-socket",
		Usage:    "the http ip:port which the distributed encoder server is listening",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "ENCODER_ADDRESS"),
	}
	EnableMetrics = cli.BoolFlag{
		Name:     common.PrefixFlag(FlagPrefix, "enable-metrics"),
		Usage:    "start metrics server",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "ENABLE_METRICS"),
	}
	UseGraphFlag = cli.BoolFlag{
		Name:     common.PrefixFlag(FlagPrefix, "use-graph"),
		Usage:    "Whether to use the graph node",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "USE_GRAPH"),
	}
	BatchSizeLimitFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "batch-size-limit"),
		Usage:    "the maximum batch size in MiB",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "BATCH_SIZE_LIMIT"),
	}
	/* Optional Flags*/
	MetricsHTTPPort = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "metrics-http-port"),
		Usage:    "the http port which the metrics prometheus server is listening",
		Required: false,
		Value:    "9100",
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "METRICS_HTTP_PORT"),
	}
	IndexerDataDirFlag = cli.StringFlag{
		Name:   common.PrefixFlag(FlagPrefix, "indexer-data-dir"),
		Usage:  "the data directory for the indexer",
		EnvVar: common.PrefixEnvVar(envVarPrefix, "INDEXER_DATA_DIR"),
		Value:  "./data/",
	}
	EncodingTimeoutFlag = cli.DurationFlag{
		Name:     "encoding-timeout",
		Usage:    "connection timeout from grpc call to encoder",
		Required: false,
		Value:    10 * time.Second,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "ENCODING_TIMEOUT"),
	}
	AttestationTimeoutFlag = cli.DurationFlag{
		Name:     "attestation-timeout",
		Usage:    "connection timeout from grpc call to DA nodes for attestation",
		Required: false,
		Value:    20 * time.Second,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "ATTESTATION_TIMEOUT"),
	}
	ChainReadTimeoutFlag = cli.DurationFlag{
		Name:     "chain-read-timeout",
		Usage:    "connection timeout to read from chain",
		Required: false,
		Value:    5 * time.Second,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "CHAIN_READ_TIMEOUT"),
	}
	ChainWriteTimeoutFlag = cli.DurationFlag{
		Name:     "chain-write-timeout",
		Usage:    "connection timeout to write to chain",
		Required: false,
		Value:    90 * time.Second,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "CHAIN_WRITE_TIMEOUT"),
	}
	ChainStateTimeoutFlag = cli.DurationFlag{
		Name:     "chain-state-timeout",
		Usage:    "connection timeout to read state from chain",
		Required: false,
		Value:    15 * time.Second,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "CHAIN_STATE_TIMEOUT"),
	}
	TransactionBroadcastTimeoutFlag = cli.DurationFlag{
		Name:     "transaction-broadcast-timeout",
		Usage:    "timeout to broadcast transaction",
		Required: false,
		Value:    10 * time.Minute,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "TRANSACTION_BROADCAST_TIMEOUT"),
	}
	NumConnectionsFlag = cli.IntFlag{
		Name:     "num-connections",
		Usage:    "maximum number of connections to encoders (defaults to 256)",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "NUM_CONNECTIONS"),
		Value:    256,
	}
	FinalizerIntervalFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "finalizer-interval"),
		Usage:    "Interval at which to check for finalized blobs",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "FINALIZER_INTERVAL"),
		Value:    6 * time.Minute,
	}
	FinalizerPoolSizeFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "finalizer-pool-size"),
		Usage:    "Size of the finalizer workerpool",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "FINALIZER_POOL_SIZE"),
		Value:    4,
	}
	EncodingRequestQueueSizeFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "encoding-request-queue-size"),
		Usage:    "Size of the encoding request queue",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "ENCODING_REQUEST_QUEUE_SIZE"),
		Value:    500,
	}
	SRSOrderFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "srs-order"),
		Usage:    "Size of the encoding request queue",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "SRS_ORDER"),
	}
	MaxNumRetriesPerBlobFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "max-num-retries-per-blob"),
		Usage:    "Maximum number of retries to process a blob before marking the blob as FAILED",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "MAX_NUM_RETRIES_PER_BLOB"),
		Value:    2,
	}
	// This flag is available so that we can manually adjust the number of chunks if desired for testing purposes or for other reasons.
	// For instance, we may want to increase the number of chunks / reduce the chunk size to reduce the amount of data that needs to be
	// downloaded by light clients for DAS.
	TargetNumChunksFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "target-num-chunks"),
		Usage:    "Target number of chunks per blob. If set to zero, the number of chunks will be calculated based on the ratio of the total stake to the minimum stake",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "TARGET_NUM_CHUNKS"),
		Value:    0,
	}
	MaxBlobsToFetchFromStoreFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "max-blobs-to-fetch-from-store"),
		Usage:    "Limit used to specify how many blobs to fetch from store at time when used with dynamodb pagination",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "MAX_BLOBS_TO_FETCH_FROM_STORE"),
		Value:    100,
	}
	FinalizationBlockDelayFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "finalization-block-delay"),
		Usage:    "The block delay to use for pulling operator state in order to ensure the state is finalized",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "FINALIZATION_BLOCK_DELAY"),
		Value:    75,
	}
	EnableGnarkBundleEncodingFlag = cli.BoolFlag{
		Name:     common.PrefixFlag(FlagPrefix, "enable-gnark-bundle-encoding"),
		Usage:    "Enable Gnark bundle encoding for chunks",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "ENABLE_GNARK_BUNDLE_ENCODING"),
	}
	MaxNodeConnectionsFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "max-node-connections"),
		Usage:    "Maximum number of connections to the node. Only used when minibatching is enabled. Defaults to 1024.",
		Required: false,
		Value:    1024,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "MAX_NODE_CONNECTIONS"),
	}
	MaxNumRetriesPerDispersalFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "max-num-retries-per-dispersal"),
		Usage:    "Maximum number of retries to disperse a minibatch. Only used when minibatching is enabled. Defaults to 3.",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "MAX_NUM_RETRIES_PER_DISPERSAL"),
		Value:    3,
	}
)

var requiredFlags = []cli.Flag{
	S3BucketNameFlag,
	DynamoDBTableNameFlag,
	PullIntervalFlag,
	BlsOperatorStateRetrieverFlag,
	EigenDAServiceManagerFlag,
	EncoderSocket,
	EnableMetrics,
	BatchSizeLimitFlag,
	UseGraphFlag,
	SRSOrderFlag,
}

var optionalFlags = []cli.Flag{
	MetricsHTTPPort,
	IndexerDataDirFlag,
	EncodingTimeoutFlag,
	AttestationTimeoutFlag,
	ChainReadTimeoutFlag,
	ChainWriteTimeoutFlag,
	ChainStateTimeoutFlag,
	TransactionBroadcastTimeoutFlag,
	NumConnectionsFlag,
	FinalizerIntervalFlag,
	FinalizerPoolSizeFlag,
	EncodingRequestQueueSizeFlag,
	MaxNumRetriesPerBlobFlag,
	TargetNumChunksFlag,
	MaxBlobsToFetchFromStoreFlag,
	FinalizationBlockDelayFlag,
	MaxNodeConnectionsFlag,
	MaxNumRetriesPerDispersalFlag,
	EnableGnarkBundleEncodingFlag,
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func init() {
	Flags = append(requiredFlags, optionalFlags...)
	Flags = append(Flags, geth.EthClientFlags(envVarPrefix)...)
	Flags = append(Flags, common.LoggerCLIFlags(envVarPrefix, FlagPrefix)...)
	Flags = append(Flags, indexer.CLIFlags(envVarPrefix)...)
	Flags = append(Flags, aws.ClientFlags(envVarPrefix, FlagPrefix)...)
	Flags = append(Flags, thegraph.CLIFlags(envVarPrefix)...)
	Flags = append(Flags, common.KMSWalletCLIFlags(envVarPrefix, FlagPrefix)...)
}
