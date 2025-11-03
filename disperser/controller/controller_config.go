package controller

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	clients "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/config"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/core/thegraph"
)

var _ config.DocumentedConfig = &ControllerConfig{}

// ControllerConfig contains configuration parameters for the controller.
// The controller is responsible for batching encoded blobs, dispersing them to DA nodes,
// collecting signatures, and creating attestations.
type ControllerConfig struct {
	// PullInterval is how frequently the controller polls for new encoded blobs to batch and dispatch.
	// Must be positive.
	PullInterval time.Duration

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

	// NumRequestRetries is the number of times to retry dispersing to a node after the initial attempt fails.
	// The current implementation has performance issues, so this should typically be 0 (no retries).
	// Must be non-negative.
	NumRequestRetries int

	// MaxBatchSize is the maximum number of blobs to include in a single batch for dispersal.
	// Must be at least 1.
	MaxBatchSize int32

	// SignificantSigningThresholdPercentage is a configurable "important" signing threshold percentage.
	// Used to track signing metrics and understand system performance.
	// If the value is 0, special handling for this threshold is disabled.
	// Must be between 0 and 100.
	SignificantSigningThresholdPercentage uint8

	// SignificantSigningMetricsThresholds are signing thresholds for metrics reporting.
	// Values should be decimal strings between "0.0" (0% signed) and "1.0" (100% signed).
	// Example: []string{"0.55", "0.67"}
	SignificantSigningMetricsThresholds []string

	// NumConcurrentRequests is the size of the worker pool for processing dispersal requests concurrently.
	// Must be at least 1.
	NumConcurrentRequests int

	// NodeClientCacheSize is the maximum number of node clients to cache for reuse.
	// Must be at least 1.
	NodeClientCacheSize int

	// If true, use the new payment authentication system running on the controller.
	// If false, payment authentication is disabled and request validation will always fail
	EnablePaymentAuthentication bool

	// Enables/disables the gRPC server
	EnableGrpcServer bool

	// Port that the gRPC server listens on
	GrpcPort uint16

	// Maximum size of a gRPC message that the server will accept (in bytes)
	MaxGrpcMessageSize uint32

	// Maximum time a connection can be idle before it is closed.
	MaxIdleConnectionAge time.Duration

	// Maximum age of a request in the past that the server will accept.
	// Requests older than this will be rejected to prevent replay attacks.
	RequestMaxPastAge time.Duration

	// Maximum age of a request in the future that the server will accept.
	// Requests with timestamps too far in the future will be rejected.
	RequestMaxFutureAge time.Duration

	// Configurations for the encoding manager (i.e. for talking to encoders).
	EncodingManagerConfig EncodingManagerConfig

	// The name of the DynamoDB table to use for storing blob metadata.
	DynamoDBTableName string `docs:"required"`

	// Configurations for the Ethereum client.
	EthClientConfig geth.EthClientConfig

	// Configuration for the AWS client. TODO
	AwsClientConfig aws.ClientConfig

	// If true, disables signing of disperser store chunks.
	DisperserStoreChunksSigningDisabled bool

	// Configuration for the dispersal request signer.
	DispersalRequestSignerConfig *clients.DispersalRequestSignerConfig

	// The format for logging output ('text' or 'json').
	LogFormat string

	// The logging level.
	// LevelDebug Level = -4
	// LevelInfo  Level = 0
	// LevelWarn  Level = 4
	// LevelError Level = 8
	LogLevel int

	// The duration between indexer pulls
	IndexerPullInterval time.Duration

	// Configuration for The Graph client.
	ChainStateConfig thegraph.Config

	// Whether to use The Graph for chain state queries.
	UseGraph bool

	// The EigenDA contract directory address.
	EigenDAContractDirectoryAddress string `docs:"required"`

	// The port to expose metrics on.
	MetricsPort int

	// The path for the controller readiness probe.
	ControllerReadinessProbePath string

	// Configuration for the heartbeat monitor.
	HeartbeatMonitorConfig *healthcheck.HeartbeatMonitorConfig

	// Configuration for payment authorization.
	PaymentAuthorizationConfig *PaymentAuthorizationConfig
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
	if c.NumRequestRetries < 0 {
		return fmt.Errorf("NumRequestRetries must be non-negative, got %d", c.NumRequestRetries)
	}
	if c.MaxBatchSize < 1 {
		return fmt.Errorf("MaxBatchSize must be at least 1, got %d", c.MaxBatchSize)
	}
	if c.SignificantSigningThresholdPercentage > 100 {
		return fmt.Errorf(
			"SignificantSigningThresholdPercentage must be between 0 and 100, got %d",
			c.SignificantSigningThresholdPercentage)
	}
	for _, threshold := range c.SignificantSigningMetricsThresholds {
		val, err := strconv.ParseFloat(threshold, 64)
		if err != nil {
			return fmt.Errorf("SignificantSigningMetricsThresholds contains invalid float: %s", threshold)
		}
		if val < 0.0 || val > 1.0 {
			return fmt.Errorf(
				"SignificantSigningMetricsThresholds must be between 0.0 and 1.0, got %s",
				threshold)
		}
	}
	if c.NumConcurrentRequests < 1 {
		return fmt.Errorf("NumConcurrentRequests must be at least 1, got %d", c.NumConcurrentRequests)
	}
	if c.NodeClientCacheSize < 1 {
		return fmt.Errorf("NodeClientCacheSize must be at least 1, got %d", c.NodeClientCacheSize)
	}

	if c.EnableGrpcServer {
		if c.MaxGrpcMessageSize < 0 {
			return fmt.Errorf("max grpc message size must be >= 0, got %d", c.MaxGrpcMessageSize)
		}
		if c.MaxIdleConnectionAge < 0 {
			return fmt.Errorf("max idle connection age must be >= 0, got %v", c.MaxIdleConnectionAge)
		}
		if c.RequestMaxPastAge < 0 {
			return fmt.Errorf("request max past age must be >= 0, got %v", c.RequestMaxPastAge)
		}
		if c.RequestMaxFutureAge < 0 {
			return fmt.Errorf("request max future age must be >= 0, got %v", c.RequestMaxFutureAge)
		}
	}
	return nil
}

// DefaultControllerConfig returns a ControllerConfig with default values.
func DefaultControllerConfig() *ControllerConfig {
	return &ControllerConfig{
		PullInterval:                          1 * time.Second,
		FinalizationBlockDelay:                75,
		AttestationTimeout:                    45 * time.Second,
		BatchMetadataUpdatePeriod:             time.Minute,
		BatchAttestationTimeout:               55 * time.Second,
		SignatureTickInterval:                 50 * time.Millisecond,
		NumRequestRetries:                     0,
		MaxBatchSize:                          32,
		SignificantSigningThresholdPercentage: 55,
		SignificantSigningMetricsThresholds:   []string{"0.55", "0.67"},
		NumConcurrentRequests:                 600,
		NodeClientCacheSize:                   400,
		EnablePaymentAuthentication:           true,
		EnableGrpcServer:                      true,
		GrpcPort:                              32010,
		MaxGrpcMessageSize:                    1024 * 1024, // 1 MB
		MaxIdleConnectionAge:                  5 * time.Minute,
		RequestMaxPastAge:                     5 * time.Minute,
		RequestMaxFutureAge:                   3 * time.Minute,
		EncodingManagerConfig:                 *DefaultEncodingManagerConfig(),
		EthClientConfig: geth.EthClientConfig{
			NumConfirmations: 10,
			NumRetries:       2,
		},
		AwsClientConfig: aws.ClientConfig{
			FragmentParallelismFactor:   8,
			FragmentParallelismConstant: 0,
		},
		DisperserStoreChunksSigningDisabled: false,
		DispersalRequestSignerConfig:        &clients.DispersalRequestSignerConfig{},
		LogFormat:                           "json",
		LogLevel:                            int(slog.LevelDebug),
		IndexerPullInterval:                 1 * time.Second,
		ChainStateConfig: thegraph.Config{ // TODO what defaults?
			PullInterval: 2 * time.Second,
			MaxRetries:   5,
		},
		UseGraph:                     true,
		MetricsPort:                  9101,
		ControllerReadinessProbePath: "/tmp/controller-ready",
		HeartbeatMonitorConfig: &healthcheck.HeartbeatMonitorConfig{
			FilePath:         "/tmp/controller-health",
			MaxStallDuration: 4 * time.Minute,
		},
		PaymentAuthorizationConfig: DefaultPaymentAuthorizationConfig(),
	}
}

func (c *ControllerConfig) GetEnvVarPrefix() string {
	return "CONTROLLER"
}

func (c *ControllerConfig) GetName() string {
	return "Controller"
}

func (c *ControllerConfig) GetPackagePaths() []string {
	return []string{
		"github.com/Layr-Labs/eigenda/disperser/controller",
		"github.com/Layr-Labs/eigenda/common",
	}
}
