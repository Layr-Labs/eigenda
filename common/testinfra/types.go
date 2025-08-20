package testinfra

import (
	"time"

	"github.com/Layr-Labs/eigenda/common/testinfra/containers"
)

// InfraConfig defines the configuration for test infrastructure components
type InfraConfig struct {
	// Anvil blockchain configuration
	Anvil containers.AnvilConfig `json:"anvil"`

	// LocalStack AWS simulation configuration
	LocalStack containers.LocalStackConfig `json:"localstack"`

	// The Graph node configuration
	GraphNode containers.GraphNodeConfig `json:"graphnode"`

	// Global test configuration
	Timeout time.Duration `json:"timeout"`
}

// DefaultConfig returns a sensible default configuration
func DefaultConfig() InfraConfig {
	return InfraConfig{
		Anvil: containers.AnvilConfig{
			Enabled:   true,
			ChainID:   31337,
			BlockTime: 0, // instant mining
			GasLimit:  30000000,
			GasPrice:  0,
			Accounts:  10,
			Mnemonic:  "test test test test test test test test test test test junk",
		},
		LocalStack: containers.LocalStackConfig{
			Enabled:  true,
			Services: []string{"s3", "dynamodb", "kms", "secretsmanager"},
			Region:   "us-east-1",
			Debug:    false,
		},
		GraphNode: containers.GraphNodeConfig{
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
	AnvilRPC          string `json:"anvil_rpc"`
	AnvilChainID      int    `json:"anvil_chain_id"`
	LocalStackURL     string `json:"localstack_url"`
	GraphNodeURL      string `json:"graphnode_url"`       // GraphQL endpoint (port 8000)
	GraphNodeAdminURL string `json:"graphnode_admin_url"` // Admin endpoint (port 8020)
	IPFSURL           string `json:"ipfs_url"`            // IPFS API endpoint (port 5001)
	PostgresURL       string `json:"postgres_url"`
}
