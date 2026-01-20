package relay

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/config"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/relay/limiter"
	"github.com/docker/go-units"
)

var _ config.DocumentedConfig = (*RelayConfig)(nil)

type RelayConfig struct {
	// Configuration for the AWS client
	AWS aws.ClientConfig

	// Configuration for the eth client
	EthClient geth.EthClientConfig

	// The graph indexer configuration
	Graph thegraph.Config

	// Name of the bucket to store blobs
	BucketName string `docs:"required"`

	// Name of the dynamodb table to store blob metadata
	MetadataTableName string `docs:"required"`

	// Object storage backend to use (s3 or oci)
	ObjectStorageBackend string

	// OCI region (only used when object-storage-backend is oci)
	OCIRegion string

	// OCI compartment ID (only used when object-storage-backend is oci)
	OCICompartmentID string

	// OCI namespace (only used when object-storage-backend is oci). If not provided, will be retrieved dynamically
	OCINamespace string

	// Address of the EigenDA directory contract, which points to all other EigenDA contract addresses. This is the
	// only contract entrypoint needed offchain
	EigenDADirectory string

	// Relay keys to use
	RelayKeys []v2.RelayKey `docs:"required"`

	// Port to listen on for gRPC
	GRPCPort int `docs:"required"`

	// Enable prometheus metrics collection
	EnableMetrics bool `docs:"required"`

	// Max size of a gRPC message in bytes
	MaxGRPCMessageSize int

	// Max number of items in the metadata cache
	MetadataCacheSize int

	// Max number of concurrent metadata fetches
	MetadataMaxConcurrency int

	// The size of the blob cache, in bytes
	BlobCacheBytes uint64

	// Max number of concurrent blob fetches
	BlobMaxConcurrency int

	// Size of the chunk cache, in bytes
	ChunkCacheBytes uint64

	// Max number of concurrent chunk fetches
	ChunkMaxConcurrency int

	// Max number of keys to fetch in a single GetChunks request
	MaxKeysPerGetChunksRequest int

	// RateLimits contains configuration for rate limiting.
	RateLimits limiter.Config

	// Max number of items in the authentication key cache
	AuthenticationKeyCacheSize int

	// Disable GetChunks() authentication
	AuthenticationDisabled bool

	// Max age of a GetChunks request
	GetChunksRequestMaxPastAge time.Duration

	// Max future age of a GetChunks request
	GetChunksRequestMaxFutureAge time.Duration

	// Timeouts contains configuration for relay timeouts.
	Timeouts TimeoutConfig

	// The interval at which to refresh the onchain state
	OnchainStateRefreshInterval time.Duration

	// Port to listen on for metrics
	MetricsPort int

	// Enable pprof profiling
	EnablePprof bool

	// Port to listen on for pprof
	PprofHttpPort int

	// Maximum age of a gRPC connection before it is closed. If zero, then the server will not close connections based
	// on age
	MaxConnectionAge time.Duration

	// Grace period after MaxConnectionAge before the connection is forcibly closed
	MaxConnectionAgeGrace time.Duration

	// Maximum time a connection can be idle before it is closed
	MaxIdleConnectionAge time.Duration

	// The output type for logs, must be "json" or "text".
	LogOutputType string

	// Whether to enable color in log output (only applies to text output).
	LogColor bool

	// Address of the OperatorStateRetriever contract.
	OperatorStateRetrieverAddr string

	// Address of the Eigen DA service manager contract.
	EigenDAServiceManagerAddr string
}

func DefaultRelayConfig() *RelayConfig {
	return &RelayConfig{
		AWS:       aws.DefaultClientConfig(),
		EthClient: geth.DefaultEthClientConfig(),
		Graph: thegraph.Config{
			// Endpoint is required and has no default value
			PullInterval: 100 * time.Millisecond,
			MaxRetries:   5,
		},
		ObjectStorageBackend:       "s3",
		MaxGRPCMessageSize:         4 * units.MiB,
		MetadataCacheSize:          units.MiB,
		MetadataMaxConcurrency:     32,
		BlobCacheBytes:             units.GiB,
		BlobMaxConcurrency:         32,
		ChunkCacheBytes:            units.GiB,
		ChunkMaxConcurrency:        32,
		MaxKeysPerGetChunksRequest: 1024,
		RateLimits: limiter.Config{
			MaxGetBlobOpsPerSecond:          1024,
			GetBlobOpsBurstiness:            1024,
			MaxGetBlobBytesPerSecond:        20 * units.MiB,
			GetBlobBytesBurstiness:          20 * units.MiB,
			MaxConcurrentGetBlobOps:         1024,
			MaxGetChunkOpsPerSecond:         1024,
			GetChunkOpsBurstiness:           1024,
			MaxGetChunkBytesPerSecond:       80 * units.MiB,
			GetChunkBytesBurstiness:         800 * units.MiB,
			MaxConcurrentGetChunkOps:        1024,
			MaxGetChunkOpsPerSecondClient:   8,
			GetChunkOpsBurstinessClient:     8,
			MaxGetChunkBytesPerSecondClient: 40 * units.MiB,
			GetChunkBytesBurstinessClient:   400 * units.MiB,
			MaxConcurrentGetChunkOpsClient:  1,
		},
		AuthenticationKeyCacheSize:   1024 * 1024,
		AuthenticationDisabled:       false,
		GetChunksRequestMaxPastAge:   5 * time.Minute,
		GetChunksRequestMaxFutureAge: 5 * time.Minute,
		Timeouts: TimeoutConfig{
			GetChunksTimeout:               20 * time.Second,
			GetBlobTimeout:                 20 * time.Second,
			InternalGetMetadataTimeout:     5 * time.Second,
			InternalGetBlobTimeout:         20 * time.Second,
			InternalGetProofsTimeout:       5 * time.Second,
			InternalGetCoefficientsTimeout: 20 * time.Second,
		},
		OnchainStateRefreshInterval: time.Hour,
		MetricsPort:                 9191,
		EnablePprof:                 false,
		PprofHttpPort:               6060,
		MaxConnectionAge:            5 * time.Minute,
		MaxConnectionAgeGrace:       30 * time.Second,
		MaxIdleConnectionAge:        time.Minute,
		LogOutputType:               string(common.JSONLogFormat),
		LogColor:                    false,
	}
}

func (c *RelayConfig) GetEnvVarPrefix() string {
	return "RELAY"
}

func (c *RelayConfig) GetName() string {
	return "Relay"
}

func (c *RelayConfig) GetPackagePaths() []string {
	return []string{
		"github.com/Layr-Labs/eigenda/common/aws",
		"github.com/Layr-Labs/eigenda/common/geth",
		"github.com/Layr-Labs/eigenda/core/thegraph",
		"github.com/Layr-Labs/eigenda/relay",
		"github.com/Layr-Labs/eigenda/relay/limiter",
	}
}

func (c *RelayConfig) Verify() error {
	if c.BucketName == "" {
		return fmt.Errorf("invalid bucket name: %s", c.BucketName)
	}
	if c.MetadataTableName == "" {
		return fmt.Errorf("invalid metadata table name: %s", c.MetadataTableName)
	}
	if c.ObjectStorageBackend != "s3" && c.ObjectStorageBackend != "oci" {
		return fmt.Errorf("invalid object storage backend: %s", c.ObjectStorageBackend)
	}
	if c.ObjectStorageBackend == "oci" {
		if c.OCIRegion == "" {
			return fmt.Errorf("invalid OCI region: %s", c.OCIRegion)
		}
		if c.OCICompartmentID == "" {
			return fmt.Errorf("invalid OCI compartment ID: %s", c.OCICompartmentID)
		}
		if c.OCINamespace == "" {
			return fmt.Errorf("invalid OCI namespace: %s", c.OCINamespace)
		}
	}
	if c.EigenDADirectory == "" {
		return fmt.Errorf("invalid EigenDA directory address: %s", c.EigenDADirectory)
	}
	if len(c.RelayKeys) == 0 {
		return fmt.Errorf("invalid relay keys: %v", c.RelayKeys)
	}

	if c.GRPCPort <= 0 || c.GRPCPort > 65535 {
		return fmt.Errorf("invalid gRPC port: %d", c.GRPCPort)
	}

	if c.MaxGRPCMessageSize <= 0 {
		return fmt.Errorf("invalid max gRPC message size: %d", c.MaxGRPCMessageSize)
	}

	if c.MetadataCacheSize <= 0 {
		return fmt.Errorf("invalid metadata cache size: %d", c.MetadataCacheSize)
	}

	if c.MetadataMaxConcurrency <= 0 {
		return fmt.Errorf("invalid metadata max concurrency: %d", c.MetadataMaxConcurrency)
	}

	if c.BlobCacheBytes <= 0 {
		return fmt.Errorf("invalid blob cache bytes: %d", c.BlobCacheBytes)
	}

	if c.BlobMaxConcurrency <= 0 {
		return fmt.Errorf("invalid blob max concurrency: %d", c.BlobMaxConcurrency)
	}

	if c.ChunkCacheBytes <= 0 {
		return fmt.Errorf("invalid chunk cache bytes: %d", c.ChunkCacheBytes)
	}

	if c.ChunkMaxConcurrency <= 0 {
		return fmt.Errorf("invalid chunk max concurrency: %d", c.ChunkMaxConcurrency)
	}

	if c.MaxKeysPerGetChunksRequest <= 0 {
		return fmt.Errorf("invalid max keys per GetChunks request: %d", c.MaxKeysPerGetChunksRequest)
	}

	if c.Timeouts.GetChunksTimeout <= 0 {
		return fmt.Errorf("invalid GetChunks timeout: %s", c.Timeouts.GetChunksTimeout)
	}

	if c.Timeouts.GetBlobTimeout <= 0 {
		return fmt.Errorf("invalid GetBlob timeout: %s", c.Timeouts.GetBlobTimeout)
	}

	if c.Timeouts.InternalGetMetadataTimeout <= 0 {
		return fmt.Errorf("invalid InternalGetMetadata timeout: %s", c.Timeouts.InternalGetMetadataTimeout)
	}

	if c.Timeouts.InternalGetBlobTimeout <= 0 {
		return fmt.Errorf("invalid InternalGetBlob timeout: %s", c.Timeouts.InternalGetBlobTimeout)
	}

	if c.Timeouts.InternalGetProofsTimeout <= 0 {
		return fmt.Errorf("invalid InternalGetProofs timeout: %s", c.Timeouts.InternalGetProofsTimeout)
	}

	if c.Timeouts.InternalGetCoefficientsTimeout <= 0 {
		return fmt.Errorf("invalid InternalGetCoefficients timeout: %s", c.Timeouts.InternalGetCoefficientsTimeout)
	}

	if c.OnchainStateRefreshInterval <= 0 {
		return fmt.Errorf("invalid onchain state refresh interval: %s", c.OnchainStateRefreshInterval)
	}

	if c.MetricsPort <= 0 || c.MetricsPort > 65535 {
		return fmt.Errorf("invalid metrics port: %d", c.MetricsPort)
	}

	if c.EnablePprof {
		if c.PprofHttpPort <= 0 || c.PprofHttpPort > 65535 {
			return fmt.Errorf("invalid pprof HTTP port: %d", c.PprofHttpPort)
		}
	}

	if c.MaxConnectionAge < 0 {
		return fmt.Errorf("invalid max connection age: %s", c.MaxConnectionAge)
	}

	if c.MaxConnectionAgeGrace < 0 {
		return fmt.Errorf("invalid max connection age grace: %s", c.MaxConnectionAgeGrace)
	}

	if c.MaxIdleConnectionAge < 0 {
		return fmt.Errorf("invalid max idle connection age: %s", c.MaxIdleConnectionAge)
	}

	return nil
}
