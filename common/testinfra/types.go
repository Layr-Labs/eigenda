package testinfra

import (
	"time"
)

// InfraConfig defines the configuration for test infrastructure components
type InfraConfig struct {
	// Anvil blockchain configuration
	Anvil AnvilConfig `json:"anvil"`

	// LocalStack AWS simulation configuration
	LocalStack LocalStackConfig `json:"localstack"`

	// The Graph node configuration
	GraphNode GraphNodeConfig `json:"graphnode"`

	// Global test configuration
	Timeout time.Duration `json:"timeout"`
}

// AnvilConfig configures the Anvil blockchain container
type AnvilConfig struct {
	Enabled   bool   `json:"enabled"`
	ChainID   int    `json:"chain_id"`
	BlockTime int    `json:"block_time"` // seconds between blocks, 0 for instant mining
	GasLimit  uint64 `json:"gas_limit"`
	GasPrice  uint64 `json:"gas_price"`
	Accounts  int    `json:"accounts"`   // number of pre-funded accounts
	Mnemonic  string `json:"mnemonic"`   // custom mnemonic for deterministic accounts
	Fork      string `json:"fork"`       // fork from this RPC URL
	ForkBlock uint64 `json:"fork_block"` // fork from specific block
}

// LocalStackConfig configures the LocalStack AWS simulation container
type LocalStackConfig struct {
	Enabled  bool     `json:"enabled"`
	Services []string `json:"services"` // AWS services to enable: s3, dynamodb, kms, secretsmanager
	Region   string   `json:"region"`
	Debug    bool     `json:"debug"`
}

// GraphNodeConfig configures The Graph node container
type GraphNodeConfig struct {
	Enabled      bool   `json:"enabled"`
	PostgresDB   string `json:"postgres_db"`
	PostgresUser string `json:"postgres_user"`
	PostgresPass string `json:"postgres_pass"`
	EthereumRPC  string `json:"ethereum_rpc"` // will be set to Anvil RPC if Anvil is enabled
	IPFSEndpoint string `json:"ipfs_endpoint"`
}

// DefaultConfig returns a sensible default configuration
func DefaultConfig() InfraConfig {
	return InfraConfig{
		Anvil: AnvilConfig{
			Enabled:   true,
			ChainID:   31337,
			BlockTime: 0, // instant mining
			GasLimit:  30000000,
			GasPrice:  0,
			Accounts:  10,
			Mnemonic:  "test test test test test test test test test test test junk",
		},
		LocalStack: LocalStackConfig{
			Enabled:  true,
			Services: []string{"s3", "dynamodb", "kms", "secretsmanager"},
			Region:   "us-east-1",
			Debug:    false,
		},
		GraphNode: GraphNodeConfig{
			Enabled:      false, // disabled by default due to complexity
			PostgresDB:   "graph",
			PostgresUser: "graph",
			PostgresPass: "graph",
			IPFSEndpoint: "", // will use embedded IPFS if empty
		},
		Timeout: 5 * time.Minute,
	}
}

// InfraResult contains connection information for started infrastructure
type InfraResult struct {
	AnvilRPC      string `json:"anvil_rpc"`
	AnvilChainID  int    `json:"anvil_chain_id"`
	LocalStackURL string `json:"localstack_url"`
	GraphNodeURL  string `json:"graphnode_url"`
	PostgresURL   string `json:"postgres_url"`
}
