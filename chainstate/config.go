package chainstate

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/config"
	"github.com/Layr-Labs/eigenda/common/geth"
)

var _ config.DocumentedConfig = (*RootIndexerConfig)(nil)

// RootIndexerConfig is the root configuration for the chainstate indexer.
// It separates public and secret configuration for safety.
type RootIndexerConfig struct {
	Config *IndexerConfig
	Secret *IndexerSecretConfig
}

var _ config.VerifiableConfig = (*IndexerConfig)(nil)

// IndexerConfig contains all public configuration for the chainstate indexer.
type IndexerConfig struct {
	// EigenDADirectory contract address
	EigenDADirectory string `docs:"required"`

	// Starting block number for indexing. If 0, starts from contract deployment block.
	StartBlockNumber uint64

	// Number of blocks to process in each batch during indexing.
	BlockBatchSize uint64

	// Interval between polling for new blocks on the chain.
	PollInterval time.Duration

	// Path to JSON file for persisting indexed state to disk.
	PersistencePath string `docs:"required"`

	// Interval for persisting state snapshots to disk.
	PersistInterval time.Duration

	// Port for the HTTP API server that serves indexed data queries.
	HTTPPort string `docs:"required"`

	// Logging configuration.
	LoggerConfig common.LoggerConfig

	// Ethereum client configuration for connecting to RPC endpoints.
	EthClientConfig geth.EthClientConfig
}

var _ config.VerifiableConfig = (*IndexerSecretConfig)(nil)

// IndexerSecretConfig contains sensitive configuration values.
type IndexerSecretConfig struct {
	// Ethereum RPC endpoint URLs for connecting to the blockchain.
	EthRpcUrls []string `docs:"required"`
}

// DefaultIndexerConfig returns a default configuration with sensible values.
func DefaultIndexerConfig() *IndexerConfig {
	return &IndexerConfig{
		StartBlockNumber: 0,
		BlockBatchSize:   1000,
		PollInterval:     12 * time.Second,
		PersistInterval:  30 * time.Second,
		HTTPPort:         "8080",
		LoggerConfig:     *common.DefaultLoggerConfig(),
		EthClientConfig:  geth.DefaultEthClientConfig(),
	}
}

// DefaultRootIndexerConfig returns a default root configuration.
func DefaultRootIndexerConfig() *RootIndexerConfig {
	return &RootIndexerConfig{
		Config: DefaultIndexerConfig(),
		Secret: &IndexerSecretConfig{},
	}
}

// GetName returns the name of this service for documentation.
func (c *IndexerConfig) GetName() string {
	return "ChainStateIndexer"
}

// GetEnvVarPrefix returns the environment variable prefix for this service.
func (c *IndexerConfig) GetEnvVarPrefix() string {
	return "CHAINSTATE_INDEXER"
}

// GetPackagePaths returns the package paths to scan for documentation.
func (c *IndexerConfig) GetPackagePaths() []string {
	return []string{
		"github.com/Layr-Labs/eigenda/chainstate",
	}
}

// Verify validates the configuration and returns an error if invalid.
func (c *IndexerConfig) Verify() error {
	if c.EigenDADirectory == "" {
		return fmt.Errorf("EigenDA directory address is required")
	}
	if c.PersistencePath == "" {
		return fmt.Errorf("persistence path is required")
	}
	if c.HTTPPort == "" {
		return fmt.Errorf("HTTP port is required")
	}
	if c.BlockBatchSize == 0 {
		return fmt.Errorf("block batch size must be greater than 0")
	}
	if c.PollInterval <= 0 {
		return fmt.Errorf("poll interval must be greater than 0")
	}
	if c.PersistInterval <= 0 {
		return fmt.Errorf("persist interval must be greater than 0")
	}
	return nil
}

// Verify validates the secret configuration.
func (c *IndexerSecretConfig) Verify() error {
	if len(c.EthRpcUrls) == 0 {
		return fmt.Errorf("at least one Ethereum RPC URL is required")
	}
	return nil
}

// GetName returns the name of this service for documentation.
func (c *RootIndexerConfig) GetName() string {
	return c.Config.GetName()
}

// GetEnvVarPrefix returns the environment variable prefix for this service.
func (c *RootIndexerConfig) GetEnvVarPrefix() string {
	return c.Config.GetEnvVarPrefix()
}

// GetPackagePaths returns the package paths to scan for documentation.
func (c *RootIndexerConfig) GetPackagePaths() []string {
	return c.Config.GetPackagePaths()
}

// Verify validates the root configuration.
func (c *RootIndexerConfig) Verify() error {
	if c.Config == nil {
		return fmt.Errorf("config is required")
	}
	if c.Secret == nil {
		return fmt.Errorf("secret config is required")
	}
	if err := c.Config.Verify(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}
	if err := c.Secret.Verify(); err != nil {
		return fmt.Errorf("secret config validation failed: %w", err)
	}
	return nil
}
