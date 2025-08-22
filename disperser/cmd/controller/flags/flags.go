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
	FlagPrefix   = "controller"
	envVarPrefix = "CONTROLLER"
)

var (
	DynamoDBTableNameFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "dynamodb-table-name"),
		Usage:    "Name of the dynamodb table to store blob metadata",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "DYNAMODB_TABLE_NAME"),
	}
	EigenDAContractDirectoryAddressFlag = cli.StringFlag{
		Name: common.PrefixFlag(FlagPrefix, "eigenda-contract-directory-address"),
		Usage: "Address of the EigenDA contract directory contract, which points to all other EigenDA " +
			"contract addresses.",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "EIGENDA_CONTRACT_DIRECTORY_ADDRESS"),
	}
	UseGraphFlag = cli.BoolTFlag{
		Name:     common.PrefixFlag(FlagPrefix, "use-graph"),
		Usage:    "Whether to use the graph node",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "USE_GRAPH"),
	}
	IndexerDataDirFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "indexer-data-dir"),
		Usage:    "the data directory for the indexer",
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "INDEXER_DATA_DIR"),
		Required: false,
		Value:    "./data/",
	}
	// EncodingManager Flags
	EncodingPullIntervalFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "encoding-pull-interval"),
		Usage:    "Interval at which to pull from the queue",
		Required: false,
		Value:    2 * time.Second,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "ENCODING_PULL_INTERVAL"),
	}
	AvailableRelaysFlag = cli.IntSliceFlag{
		Name:     common.PrefixFlag(FlagPrefix, "available-relays"),
		Usage:    "List of available relays",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "AVAILABLE_RELAYS"),
	}
	EncoderAddressFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "encoder-address"),
		Usage:    "the http ip:port which the distributed encoder server is listening",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "ENCODER_ADDRESS"),
	}
	EncodingRequestTimeoutFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "encoding-request-timeout"),
		Usage:    "Timeout for encoding requests",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "ENCODING_REQUEST_TIMEOUT"),
		Value:    5 * time.Minute,
	}
	EncodingStoreTimeoutFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "encoding-store-timeout"),
		Usage:    "Timeout for interacting with blob store",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "ENCODING_STORE_TIMEOUT"),
		Value:    15 * time.Second,
	}
	NumEncodingRetriesFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "num-encoding-retries"),
		Usage:    "Number of retries for encoding requests",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "NUM_ENCODING_RETRIES"),
		Value:    3,
	}
	NumRelayAssignmentFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "num-relay-assignment"),
		Usage:    "Number of relays to assign to each encoding request",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "NUM_RELAY_ASSIGNMENT"),
		Value:    2,
	}
	NumConcurrentEncodingRequestsFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "num-concurrent-encoding-requests"),
		Usage:    "Number of concurrent encoding requests",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "NUM_CONCURRENT_ENCODING_REQUESTS"),
		Value:    250,
	}
	MaxNumBlobsPerIterationFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "max-num-blobs-per-iteration"),
		Usage:    "Max number of blobs to encode in a single iteration",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "MAX_NUM_BLOBS_PER_ITERATION"),
		Value:    128,
	}
	OnchainStateRefreshIntervalFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "onchain-state-refresh-interval"),
		Usage:    "Interval at which to refresh the onchain state",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "ONCHAIN_STATE_REFRESH_INTERVAL"),
		Value:    1 * time.Hour,
	}

	// Dispatcher Flags
	DispatcherPullIntervalFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "dispatcher-pull-interval"),
		Usage:    "Interval at which to pull from the queue",
		Required: false,
		Value:    1 * time.Second,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "DISPATCHER_PULL_INTERVAL"),
	}
	AttestationTimeoutFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "attestation-timeout"),
		Usage:    "Timeout for node requests",
		Required: false,
		Value:    45 * time.Second,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "ATTESTATION_TIMEOUT"),
	}
	BatchMetadataUpdatePeriodFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "batch-metadata-update-period"),
		Usage:    "Period at which to update batch metadata",
		Required: false,
		Value:    time.Minute,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "BATCH_METADATA_UPDATE_PERIOD"),
	}
	BatchAttestationTimeoutFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "batch-attestation-timeout"),
		Usage:    "Timeout for batch attestation requests",
		Required: false,
		Value:    55 * time.Second,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "BATCH_ATTESTATION_TIMEOUT"),
	}
	SignatureTickIntervalFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "signature-tick-interval"),
		Usage:    "Interval at which new Attestations will be submitted as signature gathering progresses",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "SIGNATURE_TICK_INTERVAL"),
		Value:    50 * time.Millisecond,
	}
	FinalizationBlockDelayFlag = cli.Uint64Flag{
		Name:     common.PrefixFlag(FlagPrefix, "finalization-block-delay"),
		Usage:    "Number of blocks to wait before finalizing",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "FINALIZATION_BLOCK_DELAY"),
		Value:    75,
	}
	NumRequestRetriesFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "num-request-retries"),
		Usage:    "Number of retries for node requests",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "NUM_REQUEST_RETRIES"),
		Value:    0,
	}
	NumConcurrentDispersalRequestsFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "num-concurrent-dispersal-requests"),
		Usage:    "Number of concurrent dispersal requests",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "NUM_CONCURRENT_DISPERSAL_REQUESTS"),
		Value:    600,
	}
	NodeClientCacheNumEntriesFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "node-client-cache-num-entries"),
		Usage:    "Size (number of entries) of the node client cache",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "NODE_CLIENT_CACHE_NUM_ENTRIES"),
		Value:    400,
	}
	MaxBatchSizeFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "max-batch-size"),
		Usage:    "Max number of blobs to disperse in a batch",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "MAX_BATCH_SIZE"),
		Value:    32,
	}
	MetricsPortFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "metrics-port"),
		Usage:    "Port to expose metrics",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "METRICS_PORT"),
		Value:    9101,
	}
	DisperserStoreChunksSigningDisabledFlag = cli.BoolFlag{
		Name:     common.PrefixFlag(FlagPrefix, "disperser-store-chunks-signing-disabled"),
		Usage:    "Whether to disable signing of store chunks requests",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "DISPERSER_STORE_CHUNKS_SIGNING_DISABLED"),
	}
	DisperserKMSKeyIDFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "disperser-kms-key-id"),
		Usage:    "Name of the key used to sign disperser requests (key must be stored in AWS KMS under this name)",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "DISPERSER_KMS_KEY_ID"),
	}
	ControllerReadinessProbePathFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "controller-readiness-probe-path"),
		Usage:    "File path for the readiness probe; created once the controller is fully started and ready to serve traffic",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "CONTROLLER_READINESS_PROBE_PATH"),
		Value:    "/tmp/controller-ready",
	}
	ControllerHealthProbePathFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "controller-health-probe-path"),
		Usage:    "File path for the liveness (health) probe; updated regularly to indicate the controller is still alive and healthy",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "CONTROLLER_HEALTH_PROBE_PATH"),
		Value:    "/tmp/controller-health",
	}
	SignificantSigningThresholdPercentageFlag = cli.UintFlag{
		Name: common.PrefixFlag(FlagPrefix, "significant-signing-threshold-percentage"),
		Usage: "Percentage of stake that represents a 'significant' signing threshold. Currently used to track" +
			" metrics to better understand signing behavior.",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "SIGNIFICANT_SIGNING_THRESHOLD_PERCENTAGE"),
		Value:    55,
	}
	defaultSigningThresholds                cli.StringSlice = []string{"0.55", "0.67"}
	SignificantSigningMetricsThresholdsFlag                 = cli.StringSliceFlag{
		Name:     common.PrefixFlag(FlagPrefix, "significant-signing-thresholds"),
		Usage:    "Significant signing thresholds for metrics, each must be between 0.0 and 1.0",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "SIGNIFICANT_SIGNING_METRICS_THRESHOLDS"),
		Value:    &defaultSigningThresholds,
	}
)

var requiredFlags = []cli.Flag{
	DynamoDBTableNameFlag,
	UseGraphFlag,
	EncodingPullIntervalFlag,
	AvailableRelaysFlag,
	EncoderAddressFlag,

	DispatcherPullIntervalFlag,
	AttestationTimeoutFlag,
	BatchAttestationTimeoutFlag,
}

var optionalFlags = []cli.Flag{
	IndexerDataDirFlag,
	EncodingRequestTimeoutFlag,
	EncodingStoreTimeoutFlag,
	NumEncodingRetriesFlag,
	NumRelayAssignmentFlag,
	NumConcurrentEncodingRequestsFlag,
	MaxNumBlobsPerIterationFlag,
	OnchainStateRefreshIntervalFlag,
	SignatureTickIntervalFlag,
	FinalizationBlockDelayFlag,
	NumRequestRetriesFlag,
	NumConcurrentDispersalRequestsFlag,
	NodeClientCacheNumEntriesFlag,
	MaxBatchSizeFlag,
	MetricsPortFlag,
	DisperserStoreChunksSigningDisabledFlag,
	DisperserKMSKeyIDFlag,
	ControllerReadinessProbePathFlag,
	ControllerHealthProbePathFlag,
	SignificantSigningThresholdPercentageFlag,
	SignificantSigningMetricsThresholdsFlag,
	EigenDAContractDirectoryAddressFlag,
	BatchMetadataUpdatePeriodFlag,
}

var Flags []cli.Flag

func init() {
	Flags = append(requiredFlags, optionalFlags...)
	Flags = append(Flags, geth.EthClientFlags(envVarPrefix)...)
	Flags = append(Flags, common.LoggerCLIFlags(envVarPrefix, FlagPrefix)...)
	Flags = append(Flags, indexer.CLIFlags(envVarPrefix)...)
	Flags = append(Flags, aws.ClientFlags(envVarPrefix, FlagPrefix)...)
	Flags = append(Flags, thegraph.CLIFlags(envVarPrefix)...)
}
