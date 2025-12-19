package controller

import (
	"fmt"
	"time"

	clients "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/config"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigenda/indexer"
)

var _ config.DocumentedConfig = &ControllerConfig{}

// ControllerConfig contains configuration parameters for the controller.
type ControllerConfig struct {
	// Configuration for logging.
	Log config.SimpleLoggerConfig // TODO(cody.littley): not yet wired into flags but will be soon

	// PullInterval is how frequently the Dispatcher polls for new encoded blobs to batch and dispatch.
	// Must be positive.
	PullInterval time.Duration

	// DisperserID is the unique identifier for this disperser instance.
	DisperserID uint32

	// FinalizationBlockDelay is the number of blocks to wait before using operator state.
	// This provides a hedge against chain reorganizations.
	FinalizationBlockDelay uint64

	// BatchMetadataUpdatePeriod is the interval between attempts to refresh batch metadata
	// (reference block number and operator state).
	// Since this changes at most once per eth block, values shorter than 10 seconds are not useful.
	// In practice, checking every several minutes is sufficient.
	// Must be positive.
	BatchMetadataUpdatePeriod time.Duration

	// AttestationTimeout is the maximum time to wait for a single node to provide a signature.
	// Must be positive.
	AttestationTimeout time.Duration

	// BatchAttestationTimeout is the maximum time to wait for all nodes to provide signatures for a batch.
	// Must be positive and must be longer or equal to the AttestationTimeout.
	BatchAttestationTimeout time.Duration

	// SignatureTickInterval is how frequently attestations are updated in the blob metadata store
	// as signature gathering progresses.
	// Must be positive.
	SignatureTickInterval time.Duration

	// MaxBatchSize is the maximum number of blobs to include in a single batch for dispersal.
	// Must be at least 1.
	MaxBatchSize int32

	// SignificantSigningThresholdFraction is a configurable "important" signing threshold fraction.
	// Used to track signing metrics and understand system performance.
	// If the value is 0, special handling for this threshold is disabled.
	// Must be between 0.0 and 1.0.
	SignificantSigningThresholdFraction float64

	// If true, validators that DON'T have a human-friendly name remapping will be reported as their full validator ID
	// in metrics.
	//
	// If false, validators that DON'T have a human-friendly name remapping will be reported as "0x0" in metrics.
	//
	// NOTE: No matter the value of this field, validators that DO have a human-friendly name remapping will be reported
	// as their remapped name in metrics. If you must reduce metric cardinality by reporting ALL validators as "0x0",
	// you shouldn't define any human-friendly name remappings.
	CollectDetailedValidatorSigningMetrics bool

	// If true, accounts that DON'T have a human-friendly name remapping will be reported as their full account ID
	// in metrics.
	//
	// If false, accounts that DON'T have a human-friendly name remapping will be reported as "0x0" in metrics.
	//
	// NOTE: No matter the value of this field, accounts that DO have a human-friendly name remapping will be reported
	// as their remapped name in metrics. If you must reduce metric cardinality by reporting ALL accounts as "0x0",
	// you shouldn't define any human-friendly name remappings.
	EnablePerAccountBlobStatusMetrics bool

	// NumConcurrentRequests is the size of the worker pool for processing dispersal requests concurrently.
	// Must be at least 1.
	NumConcurrentRequests int

	// NodeClientCacheSize is the maximum number of node clients to cache for reuse.
	// Must be at least 1.
	NodeClientCacheSize int

	// MaxDispersalAge is the maximum age a dispersal request can be before it is discarded.
	// Dispersals older than this duration are marked as Failed and not processed.
	//
	// Age is determined by the BlobHeader.PaymentMetadata.Timestamp field, which is set by the
	// client at dispersal request creation time (in nanoseconds since Unix epoch).
	MaxDispersalAge time.Duration

	// The maximum a blob dispersal's self-reported timestamp can be ahead of the local wall clock time.
	// This is a preventative measure needed to prevent an attacker from sending far future timestamps
	// that result in data being tracked for a long time.
	MaxDispersalFutureAge time.Duration

	// The amount of time to retain signing rate data.
	SigningRateRetentionPeriod time.Duration

	// The duration of each signing rate bucket. Smaller buckets yield more granular data, at the cost of memory
	// and storage overhead.
	SigningRateBucketSpan time.Duration

	// BlobDispersalQueueSize is the maximum number of blobs that can be queued for dispersal.
	BlobDispersalQueueSize uint32

	// BlobDispersalRequestBatchSize is the number of blob metadata items to fetch from the store in a single request.
	// Must be at least 1.
	BlobDispersalRequestBatchSize uint32

	// BlobDispersalRequestBackoffPeriod is the delay between fetch attempts when there are no blobs ready
	// for dispersal.
	BlobDispersalRequestBackoffPeriod time.Duration

	// The period at which signing rate data is flushed to persistent storage.
	SigningRateFlushPeriod time.Duration

	// The name of the DynamoDB table used to store signing rate data.
	SigningRateDynamoDbTableName string `docs:"required"`

	// The name of the DynamoDB table used to store "core" metadata (i.e. blob statuses, signatures, etc.).
	DynamoDBTableName string

	// Whether or not to use subgraph.
	UseGraph bool

	// The contract directory contract address, which is used to derive other EigenDA contract addresses.
	EigenDAContractDirectoryAddress string `docs:"required"`

	// The port on which to expose prometheus metrics.
	MetricsPort int

	// The HTTP path to use for the controller readiness probe.
	ControllerReadinessProbePath string

	// The file path to a yaml file that maps user accounts (i.e. the parties submitting blobs) to human-friendly
	// names, which are used for metrics.
	UserAccountRemappingFilePath string

	// The file path to a yaml file that maps validator IDs to human-friendly names, which are used for metrics.
	ValidatorIdRemappingFilePath string

	// Configures the gRPC server for the controller.
	Server common.GRPCServerConfig

	// Configures the encoding manager (i.e. the interface used to send work to encoders).
	EncodingManager EncodingManagerConfig

	// Configures the indexer.
	Indexer indexer.Config

	// Configures the subgraph client.
	ChainState thegraph.Config

	// Configures the Ethereum client, which is used for talking to the EigenDA contracts.
	EthClientConfig geth.EthClientConfig

	// Configures AWS clients used by the controller.
	AwsClient aws.ClientConfig

	// If true, the disperser will not sign StoreChunks requests before sending them to validators.
	DisperserStoreChunksSigningDisabled bool

	// Configures the dispersal request signer used to sign requests to validators.
	DispersalRequestSigner clients.DispersalRequestSignerConfig

	// Configures healthchecks and heartbeat monitoring for the controller.
	HeartbeatMonitor healthcheck.HeartbeatMonitorConfig

	// Configures the payment authorization system.
	PaymentAuthorization PaymentAuthorizationConfig
}

var _ config.VerifiableConfig = &ControllerConfig{}

func DefaultControllerConfig() *ControllerConfig {
	return &ControllerConfig{
		Log:                                 *config.DefaultSimpleLoggerConfig(),
		PullInterval:                        1 * time.Second,
		FinalizationBlockDelay:              75,
		AttestationTimeout:                  45 * time.Second,
		BatchMetadataUpdatePeriod:           time.Minute,
		BatchAttestationTimeout:             55 * time.Second,
		SignatureTickInterval:               50 * time.Millisecond,
		MaxBatchSize:                        32,
		SignificantSigningThresholdFraction: 0.55,
		NumConcurrentRequests:               600,
		NodeClientCacheSize:                 400,
		MaxDispersalAge:                     45 * time.Second,
		MaxDispersalFutureAge:               45 * time.Second,
		SigningRateRetentionPeriod:          14 * 24 * time.Hour, // 2 weeks
		SigningRateBucketSpan:               10 * time.Minute,
		BlobDispersalQueueSize:              1024,
		BlobDispersalRequestBatchSize:       32,
		BlobDispersalRequestBackoffPeriod:   50 * time.Millisecond,
		SigningRateFlushPeriod:              1 * time.Minute,
	}
}

func (c *ControllerConfig) Verify() error {
	if c.PullInterval <= 0 {
		return fmt.Errorf("PullInterval must be positive, got %v", c.PullInterval)
	}
	if c.BatchMetadataUpdatePeriod <= 0 {
		return fmt.Errorf("BatchMetadataUpdatePeriod must be positive, got %v", c.BatchMetadataUpdatePeriod)
	}
	if c.AttestationTimeout <= 0 {
		return fmt.Errorf("AttestationTimeout must be positive, got %v", c.AttestationTimeout)
	}
	if c.BatchAttestationTimeout <= 0 {
		return fmt.Errorf("BatchAttestationTimeout must be positive, got %v", c.BatchAttestationTimeout)
	}
	if c.BatchAttestationTimeout < c.AttestationTimeout {
		return fmt.Errorf("BatchAttestationTimeout must be longer than AttestationTimeout, got %v < %v",
			c.BatchAttestationTimeout, c.AttestationTimeout)
	}
	if c.SignatureTickInterval <= 0 {
		return fmt.Errorf("SignatureTickInterval must be positive, got %v", c.SignatureTickInterval)
	}
	if c.MaxBatchSize < 1 {
		return fmt.Errorf("MaxBatchSize must be at least 1, got %d", c.MaxBatchSize)
	}
	if c.SignificantSigningThresholdFraction > 1.0 || c.SignificantSigningThresholdFraction < 0.0 {
		return fmt.Errorf(
			"SignificantSigningThresholdFraction must be between 0.0 and 1.0, got %f",
			c.SignificantSigningThresholdFraction)
	}
	if c.NumConcurrentRequests < 1 {
		return fmt.Errorf("NumConcurrentRequests must be at least 1, got %d", c.NumConcurrentRequests)
	}
	if c.NodeClientCacheSize < 1 {
		return fmt.Errorf("NodeClientCacheSize must be at least 1, got %d", c.NodeClientCacheSize)
	}
	if c.MaxDispersalAge <= 0 {
		return fmt.Errorf("MaxDispersalAge must be positive, got %v", c.MaxDispersalAge)
	}
	if c.MaxDispersalFutureAge <= 0 {
		return fmt.Errorf("MaxDispersalFutureAge must be positive, got %v", c.MaxDispersalFutureAge)
	}
	if c.SigningRateRetentionPeriod <= 0 {
		return fmt.Errorf("SigningRateRetentionPeriod must be positive, got %v", c.SigningRateRetentionPeriod)
	}
	if c.SigningRateBucketSpan <= 0 {
		return fmt.Errorf("SigningRateBucketSpan must be positive, got %v", c.SigningRateBucketSpan)
	}
	if c.BlobDispersalQueueSize < 1 {
		return fmt.Errorf("BlobDispersalQueueSize must be at least 1, got %d", c.BlobDispersalQueueSize)
	}
	if c.BlobDispersalRequestBatchSize < 1 {
		return fmt.Errorf("BlobDispersalRequestBatchSize must be at least 1, got %d", c.BlobDispersalRequestBatchSize)
	}
	if c.BlobDispersalRequestBackoffPeriod <= 0 {
		return fmt.Errorf("BlobDispersalRequestBackoffPeriod must be positive, got %v",
			c.BlobDispersalRequestBackoffPeriod)
	}
	if c.SigningRateFlushPeriod <= 0 {
		return fmt.Errorf("SigningRateFlushPeriod must be positive, got %v", c.SigningRateFlushPeriod)
	}
	if c.SigningRateDynamoDbTableName == "" {
		return fmt.Errorf("SigningRateDynamoDbTableName must not be empty")
	}
	if err := c.DispersalRequestSigner.Verify(); err != nil {
		return fmt.Errorf("invalid dispersal request signer config: %w", err)
	}
	if err := c.EncodingManager.Verify(); err != nil {
		return fmt.Errorf("invalid encoding manager config: %w", err)
	}
	if err := c.PaymentAuthorization.Verify(); err != nil {
		return fmt.Errorf("invalid payment authorization config: %w", err)
	}
	if err := c.Log.Verify(); err != nil {
		return fmt.Errorf("invalid logger config: %w", err)
	}
	return nil
}

func (c *ControllerConfig) GetEnvVarPrefix() string {
	return "CONTROLLER"
}

func (c *ControllerConfig) GetName() string {
	return "Controller"
}

// clients "github.com/Layr-Labs/eigenda/api/clients/v2"
// 	"github.com/Layr-Labs/eigenda/common"
// 	"github.com/Layr-Labs/eigenda/common/aws"
// 	"github.com/Layr-Labs/eigenda/common/config"
// 	"github.com/Layr-Labs/eigenda/common/geth"
// 	"github.com/Layr-Labs/eigenda/common/healthcheck"
// 	"github.com/Layr-Labs/eigenda/core/thegraph"
// 	"github.com/Layr-Labs/eigenda/indexer"

func (c *ControllerConfig) GetPackagePaths() []string {
	return []string{
		"github.com/Layr-Labs/eigenda/disperser/controller",
		"github.com/Layr-Labs/eigenda/common/config",
		"github.com/Layr-Labs/eigenda/common",
		"github.com/Layr-Labs/eigenda/indexer",
		"github.com/Layr-Labs/eigenda/core/thegraph",
		"github.com/Layr-Labs/eigenda/common/geth",
		"github.com/Layr-Labs/eigenda/common/aws",
		"github.com/Layr-Labs/eigenda/common/healthcheck",
		"github.com/Layr-Labs/eigenda/api/clients/v2",
		"github.com/Layr-Labs/eigenda/core/payments/ondemand/ondemandvalidation",
		"github.com/Layr-Labs/eigenda/core/payments/reservation/reservationvalidation",
	}
}
