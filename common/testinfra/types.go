package testinfra

import (
	"fmt"
	"os"
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
	LocalStack LocalStackInfraConfig `json:"localstack"`

	// The Graph node configuration
	GraphNode containers.GraphNodeConfig `json:"graphnode"`

	// EigenDA contract deployment configuration
	EigenDA EigenDAConfig `json:"eigenda"`

	// Global test configuration
	Timeout time.Duration `json:"timeout"`
}

// LocalStackInfraConfig combines container config with resource deployment settings
type LocalStackInfraConfig struct {
	// Container configuration
	containers.LocalStackConfig

	// Whether to automatically deploy AWS resources (S3 buckets, DynamoDB tables)
	DeployResources bool `json:"deploy_resources"`

	// Resource deployment configuration
	Resources deployment.LocalStackDeploymentConfig `json:"resources"`
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
		LocalStack: LocalStackInfraConfig{
			LocalStackConfig: containers.LocalStackConfig{
				Enabled:  true,
				Services: []string{"s3", "dynamodb", "kms", "secretsmanager"},
				Region:   "us-east-1",
				Debug:    false,
			},
			DeployResources: true, // automatically deploy AWS resources by default
			Resources:       deployment.DefaultLocalStackDeploymentConfig(),
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
	config.EigenDA.GenerateDisperserKeypair = true
	config.EigenDA.Deployer = deployment.ContractDeployer{
		Name:            "default",
		RPC:             "",                                                                 // Will be set to anvil RPC automatically
		PrivateKey:      "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", // anvil account 0
		DeploySubgraphs: true,
		VerifyContracts: false,
		VerifierURL:     "http://localhost:4000/api",
		Slow:            false,
	}

	// Configure LocalStack to deploy resources automatically
	config.LocalStack.DeployResources = true
	config.LocalStack.Resources = deployment.LocalStackDeploymentConfig{
		BucketName:          "test-eigenda-blobstore",
		MetadataTableName:   "test-BlobMetadata",
		BucketTableName:     "test-BucketStore",
		V2MetadataTableName: "test-BlobMetadata-v2",
		V2PaymentPrefix:     "e2e_v2_",
		CreateV2Resources:   true,
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

	// AWS configuration for LocalStack (populated if LocalStack.Enabled=true)
	AWSConfig *AWSTestConfig `json:"aws_config,omitempty"`

	// EigenDA contract addresses (populated if EigenDA.Enabled=true)
	EigenDAContracts *deployment.EigenDAContracts `json:"eigenda_contracts,omitempty"`

	// Disperser KMS key ID (populated if disperser keypair is generated)
	DisperserKMSKeyID string `json:"disperser_kms_key_id,omitempty"`

	// Disperser address (populated if disperser keypair is generated)
	DisperserAddress gethcommon.Address `json:"disperser_address,omitempty"`
}

// AWSTestConfig contains AWS configuration for tests
type AWSTestConfig struct {
	EndpointURL     string `json:"endpoint_url"`
	Region          string `json:"region"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`

	// Deployed resource names
	BucketName          string `json:"bucket_name,omitempty"`
	MetadataTableName   string `json:"metadata_table_name,omitempty"`
	BucketTableName     string `json:"bucket_table_name,omitempty"`
	V2MetadataTableName string `json:"v2_metadata_table_name,omitempty"`
	V2PaymentPrefix     string `json:"v2_payment_prefix,omitempty"`
}

// SetEnvironmentVariables sets AWS environment variables for LocalStack testing
func (c *AWSTestConfig) SetEnvironmentVariables() error {
	if c == nil {
		return fmt.Errorf("AWSTestConfig is nil")
	}

	fmt.Println("Setting AWS environment variables from testinfra")
	_ = os.Setenv("AWS_ENDPOINT_URL", c.EndpointURL)
	_ = os.Setenv("AWS_ACCESS_KEY_ID", c.AccessKeyID)
	_ = os.Setenv("AWS_SECRET_ACCESS_KEY", c.SecretAccessKey)
	_ = os.Setenv("AWS_DEFAULT_REGION", c.Region)

	fmt.Printf("LocalStack resources deployed:\n")
	fmt.Printf("  S3 Bucket: %s\n", c.BucketName)
	fmt.Printf("  Metadata Table: %s\n", c.MetadataTableName)
	fmt.Printf("  Bucket Table: %s\n", c.BucketTableName)
	if c.V2MetadataTableName != "" {
		fmt.Printf("  V2 Metadata Table: %s\n", c.V2MetadataTableName)
		fmt.Printf("  V2 Payment Prefix: %s\n", c.V2PaymentPrefix)
	}

	return nil
}
