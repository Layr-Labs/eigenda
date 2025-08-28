package e2e_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients"
	clientsv2 "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloaddispersal"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloadretrieval"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/testinfra"
	"github.com/Layr-Labs/eigenda/common/testinfra/containers"
	routerbindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierRouter"
	verifierv1bindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierV1"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// Global test suite state - initialized once in TestMain
type testSuite struct {
	inMemoryBlobStore bool

	infraManager *testinfra.InfraManager
	infraResult  *testinfra.InfraResult

	logger                          logging.Logger
	ethClient                       common.EthClient
	rpcClient                       common.RPCEthClient
	certBuilder                     *clientsv2.CertBuilder
	routerCertVerifier              *verification.CertVerifier
	staticCertVerifier              *verification.CertVerifier
	eigenDACertVerifierRouter       *routerbindings.ContractEigenDACertVerifierRouterTransactor
	eigenDACertVerifierRouterCaller *routerbindings.ContractEigenDACertVerifierRouterCaller
	eigenDACertVerifierV1           *verifierv1bindings.ContractEigenDACertVerifierV1
	eigenDAV2CertVerifierAddress    string // Store the V2 cert verifier address
	deployerTransactorOpts          *bind.TransactOpts

	retrievalClient clients.RetrievalClient

	relayRetrievalClientV2     *payloadretrieval.RelayPayloadRetriever
	validatorRetrievalClientV2 *payloadretrieval.ValidatorPayloadRetriever
	payloadDisperser           *payloaddispersal.PayloadDisperser
	chainReader                core.Reader

	cancel context.CancelFunc
}

var suite *testSuite

func init() {
	suite = &testSuite{}
}

func TestMain(m *testing.M) {
	flag.Parse()

	// Setup
	if err := setupTestEnvironment(); err != nil {
		fmt.Printf("Failed to setup test environment: %v\n", err)
		os.Exit(1)
	}

	// Run tests
	exitCode := m.Run()

	// Teardown
	teardownTestEnvironment()

	os.Exit(exitCode)
}

func setupTestEnvironment() error {
	fmt.Println("Setting up e2e test environment")

	// When running the e2e tests from the e2e directory, the rootPath is two levels up.
	rootPath := "../../"

	var err error

	// Start testinfra containers with EigenDA contract deployment
	ctx, infraCancel := context.WithTimeout(context.Background(), 5*time.Minute)
	suite.cancel = infraCancel

	// Use the new testinfra EigenDA API for complete orchestration
	config := testinfra.DefaultEigenDAConfig(rootPath)

	// Enable GraphNode (required for churner and other components)
	config.GraphNode.Enabled = true

	// Configure retrieval clients
	config.EigenDA.RetrievalClients.Enabled = true
	config.EigenDA.RetrievalClients.SRSOrder = "10000"
	// KZG paths will be automatically set by testinfra using its own resources

	// Enable churner service
	config.EigenDA.Churner.Enabled = true

	// Enable both v1 and v2 encoders
	encoderV1Config := testinfra.EncoderConfig{
		Enabled:       true,
		EncoderConfig: containers.DefaultEncoderV1Config(),
	}
	encoderV1Config.GRPCPort = "34000" // v1 on port 34000

	encoderV2Config := testinfra.EncoderConfig{
		Enabled:       true,
		EncoderConfig: containers.DefaultEncoderV2Config(),
	}
	encoderV2Config.GRPCPort = "34001" // v2 on port 34001

	config.EigenDA.Encoders = []testinfra.EncoderConfig{encoderV1Config, encoderV2Config}

	// Enable both v1 and v2 dispersers
	disperserV1Config := testinfra.DisperserConfig{
		Enabled:         true,
		DisperserConfig: containers.DefaultDisperserConfig(1), // v1 disperser
	}
	disperserV1Config.GRPCPort = "32003"

	disperserV2Config := testinfra.DisperserConfig{
		Enabled:         true,
		DisperserConfig: containers.DefaultDisperserConfig(2), // v2 disperser
	}
	disperserV2Config.GRPCPort = "32005"

	config.EigenDA.Dispersers = []testinfra.DisperserConfig{disperserV1Config, disperserV2Config}

	config.EigenDA.Batcher.Enabled = true
	config.EigenDA.Batcher.BatcherConfig = containers.DefaultBatcherConfig()

	// Enable controller service
	config.EigenDA.Controller.Enabled = true
	config.EigenDA.Controller.ControllerConfig = containers.DefaultControllerConfig()

	// Enable operators
	config.EigenDA.Operators.Enabled = true
	config.EigenDA.Operators.Count = 4            // Start 4 operators
	config.EigenDA.Operators.MaxOperatorCount = 3 // But limit to 3 running (for testing churn)

	// Enable relays
	config.EigenDA.Relays.Enabled = true
	config.EigenDA.Relays.Count = 4 // Start 4 relays

	if suite.inMemoryBlobStore {
		fmt.Println("Using in-memory Blob Store - disabling LocalStack")
		config.LocalStack.Enabled = false
	} else {
		fmt.Println("Using shared Blob Store with LocalStack")
	}

	fmt.Println("Starting testinfra containers with EigenDA contract deployment")
	manager, result, err := testinfra.StartCustom(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to start testinfra: %w", err)
	}
	suite.infraManager = manager
	suite.infraResult = result

	// Log contract addresses from testinfra deployment and store V2 cert verifier address
	if suite.infraResult.EigenDAContracts != nil {
		fmt.Println("✅ Using EigenDA contracts deployed by testinfra")
		suite.eigenDAV2CertVerifierAddress = suite.infraResult.EigenDAContracts.EigenDAV2CertVerifier
	}

	// Log Graph URLs if Graph Node is enabled
	if config.GraphNode.Enabled && suite.infraResult.GraphNodeURL != "" {
		fmt.Printf("✅ Graph Node available at: %s\n", suite.infraResult.GraphNodeURL)
		if suite.infraResult.GraphNodeAdminURL != "" {
			fmt.Printf("✅ Graph Admin available at: %s\n", suite.infraResult.GraphNodeAdminURL)
		}
		if suite.infraResult.IPFSURL != "" {
			fmt.Printf("✅ IPFS available at: %s\n", suite.infraResult.IPFSURL)
		}
	}

	loggerConfig := common.DefaultLoggerConfig()
	suite.logger, err = common.NewLogger(loggerConfig)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	// Log disperser keypair info from testinfra
	if suite.infraResult.DisperserKMSKeyID != "" {
		fmt.Printf("✅ Using disperser keypair from testinfra: KMS Key ID: %s, Address: %s\n",
			suite.infraResult.DisperserKMSKeyID, suite.infraResult.DisperserAddress.Hex())
	}

	// Log which components are provided by testinfra
	if suite.infraResult.ChurnerURL != "" {
		fmt.Printf("✅ Using churner from testinfra: %s\n", suite.infraResult.ChurnerURL)
	}

	if len(suite.infraResult.EncoderURLs) > 0 {
		for version, url := range suite.infraResult.EncoderURLs {
			fmt.Printf("✅ Using encoder v%s from testinfra: %s\n", version, url)
		}
	}

	if suite.infraResult.BatcherURL != "" {
		fmt.Printf("✅ Using batcher from testinfra (metrics: %s)\n", suite.infraResult.BatcherURL)
	}

	if suite.infraResult.ControllerMetricsURL != "" {
		fmt.Printf("✅ Using controller from testinfra (metrics: %s)\n", suite.infraResult.ControllerMetricsURL)
	}

	if len(suite.infraResult.OperatorAddresses) > 0 {
		fmt.Printf("✅ Using %d operators from testinfra\n", len(suite.infraResult.OperatorAddresses))
	}

	if len(suite.infraResult.RelayURLs) > 0 {
		fmt.Printf("✅ Using %d relays from testinfra\n", len(suite.infraResult.RelayURLs))
	}

	// Use cert verification components from testinfra
	if suite.infraResult.CertVerification == nil {
		return fmt.Errorf("cert verification components must be initialized by testinfra")
	}
	suite.eigenDACertVerifierV1 = suite.infraResult.CertVerification.EigenDACertVerifierV1
	suite.certBuilder = suite.infraResult.CertVerification.CertBuilder
	suite.routerCertVerifier = suite.infraResult.CertVerification.RouterCertVerifier
	suite.staticCertVerifier = suite.infraResult.CertVerification.StaticCertVerifier
	suite.eigenDACertVerifierRouter = suite.infraResult.CertVerification.EigenDACertVerifierRouter
	suite.eigenDACertVerifierRouterCaller = suite.infraResult.CertVerification.EigenDACertVerifierRouterCaller

	// Use retrieval clients from testinfra
	if suite.infraResult.RetrievalClients == nil {
		return fmt.Errorf("retrieval clients must be initialized by testinfra")
	}
	fmt.Println("Using retrieval clients from testinfra")
	suite.ethClient = suite.infraResult.RetrievalClients.EthClient
	suite.rpcClient = suite.infraResult.RetrievalClients.RPCClient
	suite.retrievalClient = suite.infraResult.RetrievalClients.RetrievalClient
	suite.chainReader = suite.infraResult.RetrievalClients.ChainReader
	suite.relayRetrievalClientV2 = suite.infraResult.RetrievalClients.RelayRetrievalClientV2
	suite.validatorRetrievalClientV2 = suite.infraResult.RetrievalClients.ValidatorRetrievalClientV2

	// Use payload disperser from testinfra
	if suite.infraResult.PayloadDisperser == nil {
		return fmt.Errorf("payload disperser must be initialized by testinfra")
	}
	fmt.Println("Using payload disperser from testinfra")
	suite.payloadDisperser = suite.infraResult.PayloadDisperser.PayloadDisperser.(*payloaddispersal.PayloadDisperser)
	suite.deployerTransactorOpts = suite.infraResult.PayloadDisperser.DeployerTransactorOpts.(*bind.TransactOpts)

	return nil
}

func teardownTestEnvironment() {
	if suite.cancel != nil {
		suite.cancel()
	}

	if suite.infraManager != nil {
		fmt.Println("Stopping testinfra containers")
		ctx, stopCancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer stopCancel()
		err := suite.infraManager.Stop(ctx)
		if err != nil {
			fmt.Printf("Error stopping infrastructure: %v\n", err)
		}
	}
}
