package controller

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common/config"
)

// ControllerConfig contains configuration parameters for the controller.
type ControllerConfig struct {
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

	// The amount of time to retain signing rate data.
	SigningRateRetentionPeriod time.Duration

	// The duration of each signing rate bucket. Smaller buckets yield more granular data, at the cost of memory
	// and storage overhead.
	SigningRateBucketSpan time.Duration
}

var _ config.VerifiableConfig = &ControllerConfig{}

func DefaultDispatcherConfig() *ControllerConfig {
	return &ControllerConfig{
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
		SigningRateRetentionPeriod:          14 * 24 * time.Hour, // 2 weeks
		SigningRateBucketSpan:               10 * time.Minute,
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
	if c.SigningRateRetentionPeriod <= 0 {
		return fmt.Errorf("SigningRateRetentionPeriod must be positive, got %v", c.SigningRateRetentionPeriod)
	}
	if c.SigningRateBucketSpan <= 0 {
		return fmt.Errorf("SigningRateBucketSpan must be positive, got %v", c.SigningRateBucketSpan)
	}
	return nil
}
