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
	UserAccountRemappingFileFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "user-account-remapping-file"),
		Usage:    "Path to YAML file for mapping account IDs to user-friendly names",
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "USER_ACCOUNT_REMAPPING_FILE"),
		Required: false,
	}
	ValidatorIdRemappingFileFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "validator-id-remapping-file"),
		Usage:    "Path to YAML file for mapping validator IDs to user-friendly names",
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "VALIDATOR_ID_REMAPPING_FILE"),
		Required: false,
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
	MaxDispersalAgeFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "max-dispersal-age"),
		Usage:    "Maximum age a dispersal request can be before it is discarded",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "MAX_DISPERSAL_AGE"),
		Value:    45 * time.Second,
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
	DetailedValidatorMetricsFlag = cli.BoolTFlag{
		Name:     common.PrefixFlag(FlagPrefix, "detailed-validator-metrics"),
		Usage:    "Whether to collect detailed validator metrics",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "DETAILED_VALIDATOR_METRICS"),
	}
	EnablePerAccountBlobStatusMetricsFlag = cli.BoolTFlag{
		Name:     common.PrefixFlag(FlagPrefix, "enable-per-account-blob-status-metrics"),
		Usage:    "Whether to report per-account blob status metrics for unmapped accounts. Accounts with valid name remappings will always use their remapped labels. If false, unmapped accounts will be aggregated under account 0x0. (default: true)",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "ENABLE_PER_ACCOUNT_BLOB_STATUS_METRICS"),
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
	DisperserPrivateKeyFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "disperser-private-key"),
		Usage:    "Private key for signing disperser requests (hex format without 0x prefix, alternative to KMS)",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "DISPERSER_PRIVATE_KEY"),
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
	ControllerHeartbeatMaxStallDurationFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "heartbeat-max-stall-duration"),
		Usage:    "Maximum time allowed between heartbeats before a component is considered stalled",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "HEARTBEAT_MAX_STALL_DURATION"),
		Value:    4 * time.Minute,
	}
	SignificantSigningThresholdFractionFlag = cli.Float64Flag{
		Name: common.PrefixFlag(FlagPrefix, "significant-signing-threshold-fraction"),
		Usage: "Fraction of stake that represents a 'significant' signing threshold. Currently used to track" +
			" metrics to better understand signing behavior.",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "SIGNIFICANT_SIGNING_THRESHOLD_FRACTION"),
		Value:    0.55,
	}
	GrpcPortFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "grpc-port"),
		Usage:    "the port for the controller gRPC server",
		Required: false,
		Value:    "32010",
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "GRPC_PORT"),
	}
	GrpcMaxMessageSizeFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "grpc-max-message-size"),
		Usage:    "maximum size of a gRPC message (in bytes). default: 1MB",
		Required: false,
		Value:    1024 * 1024,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "GRPC_MAX_MESSAGE_SIZE"),
	}
	GrpcMaxIdleConnectionAgeFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "grpc-max-idle-connection-age"),
		Usage:    "maximum time a connection can be idle before it is closed",
		Required: false,
		Value:    5 * time.Minute,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "GRPC_MAX_IDLE_CONNECTION_AGE"),
	}
	GrpcAuthorizationRequestMaxPastAgeFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "grpc-authorization-request-max-past-age"),
		Usage:    "the maximum age of an authorization request in the past that the server will accept",
		Required: false,
		Value:    5 * time.Minute,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "GRPC_AUTHORIZATION_REQUEST_MAX_PAST_AGE"),
	}
	GrpcAuthorizationRequestMaxFutureAgeFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "grpc-authorization-request-max-future-age"),
		Usage:    "the maximum age of an authorization request in the future that the server will accept",
		Required: false,
		Value:    3 * time.Minute,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "GRPC_AUTHORIZATION_REQUEST_MAX_FUTURE_AGE"),
	}
	OnDemandPaymentsTableNameFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "on-demand-payments-table-name"),
		Usage:    "Name of the DynamoDB table for storing on-demand payment state",
		Required: false,
		Value:    "on_demand",
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "ON_DEMAND_PAYMENTS_TABLE_NAME"),
	}
	OnDemandPaymentsLedgerCacheSizeFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "ondemand-payments-ledger-cache-size"),
		Usage:    "Maximum number of on-demand ledgers to keep in the LRU cache",
		Required: false,
		Value:    1024,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "ONDEMAND_PAYMENTS_LEDGER_CACHE_SIZE"),
	}
	ReservationPaymentsLedgerCacheSizeFlag = cli.IntFlag{
		Name: common.PrefixFlag(FlagPrefix, "reservation-payments-ledger-cache-size"),
		Usage: "Initial number of reservation ledgers to keep in the LRU cache. May increase " +
			"dynamically if premature evictions are detected, up to 65,536.",
		Required: false,
		Value:    1024,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "RESERVATION_PAYMENTS_LEDGER_CACHE_SIZE"),
	}
	PaymentVaultUpdateIntervalFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "payment-vault-update-interval"),
		Usage:    "Interval for checking payment vault updates",
		Required: false,
		Value:    30 * time.Second,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "PAYMENT_VAULT_UPDATE_INTERVAL"),
	}
	EnablePerAccountPaymentMetricsFlag = cli.BoolTFlag{
		Name:     common.PrefixFlag(FlagPrefix, "enable-per-account-payment-metrics"),
		Usage:    "Whether to report per-account payment metrics. If false, all metrics will be aggregated under account 0x0.",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "ENABLE_PER_ACCOUNT_PAYMENT_METRICS"),
	}
	DisperserIDFlag = cli.Uint64Flag{
		Name:     common.PrefixFlag(FlagPrefix, "disperser-id"),
		Usage:    "Unique identifier for this disperser instance. The value specified must match the index of the associated pubkey in the disperser registry",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "DISPERSER_ID"),
	}
	SigningRateRetentionPeriodFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "signing-rate-retention-period"),
		Usage:    "The amount of time to retain signing rate data",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "SIGNING_RATE_RETENTION_PERIOD"),
		Value:    14 * 24 * time.Hour,
	}
	SigningRateBucketSpanFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "signing-rate-bucket-span"),
		Usage:    "The duration of each signing rate bucket. Smaller buckets yield more granular data, at the cost of memory and storage overhead",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "SIGNING_RATE_BUCKET_SPAN"),
		Value:    10 * time.Minute,
	}
	BlobDispersalQueueSizeFlag = cli.Uint64Flag{
		Name:     common.PrefixFlag(FlagPrefix, "blob-dispersal-queue-size"),
		Usage:    "Maximum number of blobs that can be queued for dispersal",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "BLOB_DISPERSAL_QUEUE_SIZE"),
		Value:    1024,
	}
	BlobDispersalRequestBatchSizeFlag = cli.Uint64Flag{
		Name:     common.PrefixFlag(FlagPrefix, "blob-dispersal-request-batch-size"),
		Usage:    "Number of blob metadata items to fetch from the store in a single request",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "BLOB_DISPERSAL_REQUEST_BATCH_SIZE"),
		Value:    32,
	}
	BlobDispersalRequestBackoffPeriodFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "blob-dispersal-request-backoff-period"),
		Usage:    "Delay between fetch attempts when the dispersal queue is empty",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "BLOB_DISPERSAL_REQUEST_BACKOFF_PERIOD"),
		Value:    50 * time.Millisecond,
	}
	SigningRateFlushPeriodFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "signing-rate-flush-period"),
		Usage:    "The period at which signing rate data is flushed to persistent storage",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "SIGNING_RATE_FLUSH_PERIOD"),
		Value:    1 * time.Minute,
	}
	SigningRateDynamoDbTableNameFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "signing-rate-dynamodb-table-name"),
		Usage:    "The name of the DynamoDB table used to store signing rate data",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "SIGNING_RATE_DYNAMODB_TABLE_NAME"),
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
	DisperserIDFlag,
	SigningRateDynamoDbTableNameFlag,
}

var optionalFlags = []cli.Flag{
	IndexerDataDirFlag,
	UserAccountRemappingFileFlag,
	ValidatorIdRemappingFileFlag,
	EncodingRequestTimeoutFlag,
	EncodingStoreTimeoutFlag,
	NumEncodingRetriesFlag,
	NumRelayAssignmentFlag,
	NumConcurrentEncodingRequestsFlag,
	MaxNumBlobsPerIterationFlag,
	OnchainStateRefreshIntervalFlag,
	MaxDispersalAgeFlag,
	SignatureTickIntervalFlag,
	FinalizationBlockDelayFlag,
	NumConcurrentDispersalRequestsFlag,
	NodeClientCacheNumEntriesFlag,
	MaxBatchSizeFlag,
	MetricsPortFlag,
	DisperserStoreChunksSigningDisabledFlag,
	DisperserKMSKeyIDFlag,
	DisperserPrivateKeyFlag,
	ControllerReadinessProbePathFlag,
	ControllerHealthProbePathFlag,
	ControllerHeartbeatMaxStallDurationFlag,
	SignificantSigningThresholdFractionFlag,
	EigenDAContractDirectoryAddressFlag,
	BatchMetadataUpdatePeriodFlag,
	GrpcPortFlag,
	GrpcMaxMessageSizeFlag,
	GrpcMaxIdleConnectionAgeFlag,
	GrpcAuthorizationRequestMaxPastAgeFlag,
	GrpcAuthorizationRequestMaxFutureAgeFlag,
	OnDemandPaymentsTableNameFlag,
	OnDemandPaymentsLedgerCacheSizeFlag,
	ReservationPaymentsLedgerCacheSizeFlag,
	PaymentVaultUpdateIntervalFlag,
	EnablePerAccountPaymentMetricsFlag,
	DetailedValidatorMetricsFlag,
	EnablePerAccountBlobStatusMetricsFlag,
	SigningRateRetentionPeriodFlag,
	SigningRateBucketSpanFlag,
	BlobDispersalQueueSizeFlag,
	BlobDispersalRequestBatchSizeFlag,
	BlobDispersalRequestBackoffPeriodFlag,
	SigningRateFlushPeriodFlag,
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
