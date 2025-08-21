package testinfra

import (
	"time"

	"github.com/Layr-Labs/eigenda/common/testinfra/containers"
	"github.com/Layr-Labs/eigenda/common/testinfra/deployment"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// InfraConfig defines the configuration for test infrastructure components
type InfraConfig struct {
	// Anvil blockchain configuration
	Anvil containers.AnvilConfig `json:"anvil"`

	// LocalStack AWS simulation configuration
	LocalStack containers.LocalStackConfig `json:"localstack"`

	// The Graph node configuration
	GraphNode containers.GraphNodeConfig `json:"graphnode"`

	// EigenDA contract deployment configuration
	EigenDA EigenDAConfig `json:"eigenda"`

	// Global test configuration
	Timeout time.Duration `json:"timeout"`
}

// EigenDAConfig defines configuration for EigenDA contract deployment
type EigenDAConfig struct {
	// Enable EigenDA contract deployment
	Enabled bool `json:"enabled"`

	// Deploy EigenDA contracts
	DeployContracts bool `json:"deploy_contracts"`

	// Register disperser keypair after deployment
	RegisterDisperser bool `json:"register_disperser"`

	// Register blob versions and relays after deployment
	RegisterBlobVersionAndRelays bool `json:"register_blob_version_and_relays"`

	// Contract deployment configuration
	Deployer deployment.ContractDeployer `json:"deployer"`

	// Blob version parameters for registration
	BlobVersionParams []deployment.BlobVersionParam `json:"blob_version_params"`

	// Root path for contract deployment scripts
	RootPath string `json:"root_path"`

	// Private keys for operators and stakers (key: name like "opr0", value: hex private key)
	PrivateKeys map[string]string `json:"private_keys"`

	// Generate disperser keypair using KMS and register it
	GenerateDisperserKeypair bool `json:"generate_disperser_keypair"`

	// Disperser KMS key ID (output from keypair generation)
	DisperserKMSKeyID string `json:"disperser_kms_key_id"`

	// Disperser address (output from keypair generation)
	DisperserAddress string `json:"disperser_address"`
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
		EigenDA: EigenDAConfig{
			Enabled:                      false, // disabled by default
			DeployContracts:              false,
			RegisterDisperser:            false,
			RegisterBlobVersionAndRelays: false,
			RootPath:                     "", // will be set by caller
			BlobVersionParams: []deployment.BlobVersionParam{
				{
					CodingRate:      8,
					NumChunks:       8192,
					MaxNumOperators: 3537,
				},
			},
		},
		Timeout: 5 * time.Minute,
	}
}

// DefaultEigenDAConfig returns a configuration with EigenDA contract deployment enabled
func DefaultEigenDAConfig(rootPath string) InfraConfig {
	config := DefaultConfig()
	config.EigenDA.Enabled = true
	config.EigenDA.DeployContracts = true
	config.EigenDA.RegisterDisperser = true
	config.EigenDA.RegisterBlobVersionAndRelays = true
	config.EigenDA.RootPath = rootPath
	config.EigenDA.Deployer = deployment.ContractDeployer{
		Name:            "default",
		RPC:             "", // Will be set to anvil RPC automatically
		PrivateKey:      "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", // anvil account 0
		DeploySubgraphs: true,
		VerifyContracts: false,
		VerifierURL:     "http://localhost:4000/api",
		Slow:            false,
	}
	return config
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

	// EigenDA contract addresses (populated if EigenDA.Enabled=true)
	EigenDAContracts *deployment.EigenDAContracts `json:"eigenda_contracts,omitempty"`

	// Disperser KMS key ID (populated if disperser keypair is generated)
	DisperserKMSKeyID string `json:"disperser_kms_key_id,omitempty"`

	// Disperser address (populated if disperser keypair is generated)
	DisperserAddress gethcommon.Address `json:"disperser_address,omitempty"`
}
