package testinfra

import (
	"fmt"
	"os"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloadretrieval"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/testinfra/containers"
	"github.com/Layr-Labs/eigenda/common/testinfra/deployment"
	"github.com/Layr-Labs/eigenda/core"
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

	// Retrieval clients configuration
	RetrievalClients RetrievalClientsConfig `json:"retrieval_clients"`

	// Payload disperser configuration
	PayloadDisperser PayloadDisperserConfig `json:"payload_disperser"`

	// Churner configuration
	Churner ChurnerConfig `json:"churner"`

	// Batcher configuration
	Batcher BatcherConfig `json:"batcher"`

	// Encoder configuration
	Encoder EncoderConfig `json:"encoder"`
	
	// Operators configuration
	Operators OperatorsConfig `json:"operators"`
}

// ChurnerConfig defines configuration for the churner service
type ChurnerConfig struct {
	// Enable churner service
	Enabled bool `json:"enabled"`

	// Container configuration (if using containerized churner)
	containers.ChurnerConfig
}

// BatcherConfig defines configuration for the batcher service
type BatcherConfig struct {
	// Enable batcher service
	Enabled bool `json:"enabled"`

	// Container configuration (if using containerized batcher)
	containers.BatcherConfig
}

// EncoderConfig defines configuration for the encoder service
type EncoderConfig struct {
	// Enable encoder service
	Enabled bool `json:"enabled"`

	// Container configuration (if using containerized encoder)
	containers.EncoderConfig
}

// OperatorsConfig defines configuration for operator nodes
type OperatorsConfig struct {
	// Enable operators
	Enabled bool `json:"enabled"`
	
	// Number of operators to deploy
	Count int `json:"count"`
	
	// Maximum number of operators to actually run (for testing partial network)
	MaxOperatorCount int `json:"max_operator_count"`
	
	// Base operator configuration (will be customized per operator)
	BaseConfig containers.OperatorConfig `json:"base_config"`
	
	// Stakes for each operator (quorum -> operator index -> stake amount)
	Stakes map[uint32]map[int]int `json:"stakes"`
}

// RetrievalClientsConfig defines configuration for setting up retrieval clients
type RetrievalClientsConfig struct {
	// Enable retrieval clients setup
	Enabled bool `json:"enabled"`

	// RPC URL for the retrieval clients (defaults to Anvil RPC if not specified)
	RPC string `json:"rpc"`

	// SRS parameters for KZG verifier
	SRSOrder         string `json:"srs_order"`
	G1Path           string `json:"g1_path"`
	G2Path           string `json:"g2_path"`
	G2PowerOf2Path   string `json:"g2_power_of_2_path"`
	CachePath        string `json:"cache_path"`

	// Contract addresses (populated from EigenDAContracts if not specified)
	OperatorStateRetriever string `json:"operator_state_retriever"`
	ServiceManager         string `json:"service_manager"`
}

// PayloadDisperserConfig defines configuration for setting up payload disperser
type PayloadDisperserConfig struct {
	// Enable payload disperser setup
	Enabled bool `json:"enabled"`

	// Disperser private key for signing requests
	DisperserPrivateKey string `json:"disperser_private_key"`

	// Disperser client configuration
	DisperserHostname string `json:"disperser_hostname"`
	DisperserPort     string `json:"disperser_port"`

	// Timeout configurations
	DisperseBlobTimeout    time.Duration `json:"disperse_blob_timeout"`
	BlobCompleteTimeout    time.Duration `json:"blob_complete_timeout"`
	BlobStatusPollInterval time.Duration `json:"blob_status_poll_interval"`
	ContractCallTimeout    time.Duration `json:"contract_call_timeout"`
}

// RetrievalClientsComponents contains all the retrieval clients
type RetrievalClientsComponents struct {
	EthClient                  common.EthClient                           `json:"-"` // Don't serialize
	RPCClient                  common.RPCEthClient                        `json:"-"` // Don't serialize
	RetrievalClient            clients.RetrievalClient                    `json:"-"` // Don't serialize
	ChainReader                core.Reader                                `json:"-"` // Don't serialize
	RelayRetrievalClientV2     *payloadretrieval.RelayPayloadRetriever    `json:"-"` // Don't serialize
	ValidatorRetrievalClientV2 *payloadretrieval.ValidatorPayloadRetriever `json:"-"` // Don't serialize
}

// PayloadDisperserComponents contains all the payload disperser components
type PayloadDisperserComponents struct {
	PayloadDisperser       interface{} `json:"-"` // Use interface{} to avoid circular import with payloaddispersal.PayloadDisperser
	DeployerTransactorOpts interface{} `json:"-"` // Use interface{} to avoid circular import with bind.TransactOpts
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

	// Configure PayloadDisperser with default values
	config.EigenDA.PayloadDisperser = PayloadDisperserConfig{
		Enabled:                true,
		DisperserPrivateKey:    "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcded",
		DisperserHostname:      "localhost",
		DisperserPort:          "32005",
		DisperseBlobTimeout:    2 * time.Minute,
		BlobCompleteTimeout:    2 * time.Minute,
		BlobStatusPollInterval: 1 * time.Second,
		ContractCallTimeout:    5 * time.Second,
	}

	// Configure Churner with default values
	config.EigenDA.Churner = ChurnerConfig{
		Enabled: true,
		ChurnerConfig: containers.DefaultChurnerConfig(),
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

	// Cert verification components (populated if EigenDA contracts are deployed)
	CertVerification *deployment.CertVerificationComponents `json:"cert_verification,omitempty"`

	// Retrieval clients (populated if EigenDA contracts are deployed and RetrievalClients.Enabled=true)
	RetrievalClients *RetrievalClientsComponents `json:"retrieval_clients,omitempty"`

	// Payload disperser components (populated if EigenDA contracts are deployed and PayloadDisperser.Enabled=true)
	PayloadDisperser *PayloadDisperserComponents `json:"payload_disperser,omitempty"`

	// Churner service URLs (populated if Churner.Enabled=true)
	ChurnerURL         string `json:"churner_url,omitempty"`
	ChurnerInternalURL string `json:"churner_internal_url,omitempty"`

	// Encoder service URLs (populated if Encoder.Enabled=true)
	EncoderURL         string `json:"encoder_url,omitempty"`
	EncoderInternalURL string `json:"encoder_internal_url,omitempty"`

	// Batcher service URLs (populated if Batcher.Enabled=true)
	BatcherURL         string `json:"batcher_url,omitempty"`
	BatcherInternalURL string `json:"batcher_internal_url,omitempty"`
	
	// Operator addresses (populated if Operators.Enabled=true)
	// Map from operator ID to their addresses
	OperatorAddresses map[int]string `json:"operator_addresses,omitempty"`
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
