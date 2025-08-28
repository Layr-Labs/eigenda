package integration_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients"
	clientsv2 "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloaddispersal"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloadretrieval"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/testinfra"
	"github.com/Layr-Labs/eigenda/common/testinfra/containers"
	routerbindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierRouter"
	verifierv1bindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierV1"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*
These global vars are shared across tests in the integration suite to provide
communication entrypoints into the local inabox test environment
TODO: Put these into a testSuite object which is initialized per inabox E2E test. Currently this would only enable

	a client suite per test given the inabox eigenda devnet is only spun-up as a singleton and would be shared across test executions (for now).
*/
var (
	templateName      string
	testName          string
	inMemoryBlobStore bool

	testConfig   *deploy.Config
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
	deployerTransactorOpts          *bind.TransactOpts

	retrievalClient clients.RetrievalClient

	relayRetrievalClientV2     *payloadretrieval.RelayPayloadRetriever
	validatorRetrievalClientV2 *payloadretrieval.ValidatorPayloadRetriever
	payloadDisperser           *payloaddispersal.PayloadDisperser
	numConfirmations           int = 3
	numRetries                     = 0
	chainReader                core.Reader

	cancel context.CancelFunc
)

func init() {
	flag.StringVar(&templateName, "config", "testconfig-anvil-nograph.yaml", "Name of the config file (in `inabox/templates`)")
	flag.StringVar(&testName, "testname", "", "Name of the test (in `inabox/testdata`)")
	flag.BoolVar(&inMemoryBlobStore, "inMemoryBlobStore", false, "whether to use in-memory blob store")
}

func TestInaboxIntegration(t *testing.T) {
	RegisterFailHandler(Fail)

	if testing.Short() {
		t.Skip()
	}

	RunSpecs(t, "Integration Suite")
}

var _ = BeforeSuite(func() {
	By("bootstrapping test environment")

	// When running the inabox tests we typically call go test ./tests -v -config=testconfig-anvil.yaml
	// from the inabox directory so the rootPath is two levels up.
	rootPath := "../../"

	var err error
	if testName == "" {
		testName, err = deploy.CreateNewTestDirectory(templateName, rootPath)
		if err != nil {
			Expect(err).To(BeNil())
		}
	}

	testConfig = deploy.NewTestConfig(testName, rootPath)
	if testConfig.Environment.IsLocal() {
		// Start testinfra containers with EigenDA contract deployment
		ctx, infraCancel := context.WithTimeout(context.Background(), 5*time.Minute)

		// Use the new testinfra EigenDA API for complete orchestration
		config := testinfra.DefaultEigenDAConfig(rootPath)

		// Configure based on test requirements
		deployer, ok := testConfig.GetDeployer(testConfig.EigenDA.Deployer)
		config.GraphNode.Enabled = ok && deployer.DeploySubgraphs

		// Load private keys from secrets directory before starting testinfra
		// This ensures operators get funded during contract deployment
		err = testConfig.LoadPrivateKeys()
		if err != nil {
			fmt.Printf("Failed to load private keys: %v\n", err)
			Expect(err).To(BeNil())
		}

		// Pass the loaded private keys to testinfra for contract deployment
		config.EigenDA.PrivateKeys = make(map[string]string)
		for name, keyInfo := range testConfig.Pks.EcdsaMap {
			config.EigenDA.PrivateKeys[name] = keyInfo.PrivateKey
		}

		// Configure retrieval clients
		config.EigenDA.RetrievalClients.Enabled = true
		config.EigenDA.RetrievalClients.SRSOrder = "10000"

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
		config.EigenDA.Operators.Count = 4            // Start 4 operators as specified in testconfig
		config.EigenDA.Operators.MaxOperatorCount = 3 // But limit to 3 running (for testing churn)

		// Enable relays (requirement: 4 relays as per inabox config)
		config.EigenDA.Relays.Enabled = true
		config.EigenDA.Relays.Count = 4 // Start 4 relays as required by inabox

		if inMemoryBlobStore {
			fmt.Println("Using in-memory Blob Store - disabling LocalStack")
			config.LocalStack.Enabled = false
		} else {
			fmt.Println("Using shared Blob Store with LocalStack")
		}

		fmt.Println("Starting testinfra containers with EigenDA contract deployment")
		manager, result, err := testinfra.StartCustom(ctx, config)
		Expect(err).To(BeNil())
		infraManager = manager
		infraResult = result
		cancel = infraCancel

		// Update test config to use testinfra endpoints
		fmt.Printf("Updating RPC URL from %s to %s\n", testConfig.Deployers[0].RPC, infraResult.AnvilRPC)
		testConfig.Deployers[0].RPC = infraResult.AnvilRPC

		// Also update the globals to ensure the RPC is propagated to all services
		if testConfig.Services.Variables == nil {
			testConfig.Services.Variables = make(map[string]map[string]string)
		}
		if testConfig.Services.Variables["globals"] == nil {
			testConfig.Services.Variables["globals"] = make(map[string]string)
		}
		testConfig.Services.Variables["globals"]["CHAIN_RPC"] = infraResult.AnvilRPC

		// Also propagate the LocalStack endpoint if it's available
		if awsEndpoint := os.Getenv("AWS_ENDPOINT_URL"); awsEndpoint != "" {
			testConfig.Services.Variables["globals"]["AWS_ENDPOINT_URL"] = awsEndpoint
			testConfig.SetLocalstackEndpoint(awsEndpoint)
			testConfig.SetLocalstackRegion("us-east-1")
		}

		// Set environment variable to indicate dispersers are provided by testinfra
		os.Setenv("DISPERSERS_PROVIDED", "true")

		// Also export disperser URLs if needed by other components
		if infraResult.DisperserURLs != nil {
			if url, ok := infraResult.DisperserURLs["1"]; ok {
				os.Setenv("DISPERSER_V1_URL", url)
			}
			if url, ok := infraResult.DisperserURLs["2"]; ok {
				os.Setenv("DISPERSER_V2_URL", url)
			}
		}

		// Update contract addresses from testinfra deployment
		if infraResult.EigenDAContracts != nil {
			fmt.Println("Using EigenDA contracts deployed by testinfra")
			contracts := infraResult.EigenDAContracts
			testConfig.EigenDA.ServiceManager = contracts.ServiceManager
			testConfig.EigenDA.OperatorStateRetriever = contracts.OperatorStateRetriever
			testConfig.EigenDA.RegistryCoordinator = contracts.RegistryCoordinator
			testConfig.EigenDAV1CertVerifier = contracts.EigenDAV1CertVerifier
			testConfig.EigenDAV2CertVerifier = contracts.EigenDAV2CertVerifier
			testConfig.EigenDA.CertVerifierRouter = contracts.EigenDACertVerifierRouter
			// Set CertVerifier to V2 for v2 tests
			testConfig.EigenDA.CertVerifier = contracts.EigenDAV2CertVerifier
		}

		// Update Graph URLs if Graph Node is enabled
		if config.GraphNode.Enabled && infraResult.GraphNodeURL != "" {
			graphURL := infraResult.GraphNodeURL + "/subgraphs/name/Layr-Labs/eigenda-operator-state"
			fmt.Printf("Setting Graph URL to %s\n", graphURL)
			testConfig.GraphURL = graphURL

			// Set admin URL for subgraph deployment
			fmt.Printf("Setting Graph Admin URL to %s\n", infraResult.GraphNodeAdminURL)
			testConfig.GraphAdminURL = infraResult.GraphNodeAdminURL

			// Set IPFS URL for subgraph deployment
			if infraResult.IPFSURL != "" {
				fmt.Printf("Setting IPFS URL to %s\n", infraResult.IPFSURL)
				testConfig.IPFSURL = infraResult.IPFSURL
			}
		}

		loggerConfig := common.DefaultLoggerConfig()
		logger, err = common.NewLogger(loggerConfig)
		Expect(err).To(BeNil())

		// Create eth client for remaining operations
		deployer = testConfig.Deployers[0] // Use first deployer
		pk := testConfig.Pks.EcdsaMap[deployer.Name].PrivateKey
		pk = strings.TrimPrefix(pk, "0x")
		ethClient, err = geth.NewMultiHomingClient(geth.EthClientConfig{
			RPCURLs:          []string{testConfig.Deployers[0].RPC},
			PrivateKeyString: pk,
			NumConfirmations: numConfirmations,
			NumRetries:       numRetries,
		}, gethcommon.Address{}, logger)
		Expect(err).To(BeNil())

		rpcClient, err = ethrpc.Dial(testConfig.Deployers[0].RPC)
		Expect(err).To(BeNil())

		// Use disperser keypair generated by testinfra
		if infraResult.DisperserKMSKeyID != "" {
			testConfig.DisperserKMSKeyID = infraResult.DisperserKMSKeyID
			testConfig.DisperserAddress = infraResult.DisperserAddress
			fmt.Printf("✅ Using disperser keypair from testinfra: KMS Key ID: %s, Address: %s\n",
				testConfig.DisperserKMSKeyID, testConfig.DisperserAddress.Hex())
		}

		// Pass churner URL from testinfra to GenerateAllVariables if available
		if infraResult.ChurnerURL != "" {
			// Use external URL for local binaries
			os.Setenv("CHURNER_URL", infraResult.ChurnerURL)
			fmt.Printf("✅ Using churner from testinfra: %s\n", infraResult.ChurnerURL)
			// Verify it was set
			fmt.Printf("✅ CHURNER_URL env var is now: %s\n", os.Getenv("CHURNER_URL"))
		} else {
			fmt.Printf("⚠️ No ChurnerURL in infraResult\n")
		}

		// Pass encoder URLs from testinfra if available
		if infraResult.EncoderURLs != nil && len(infraResult.EncoderURLs) > 0 {
			// Use external URLs for local binaries
			// Note: When batcher is also containerized, it uses internal URL directly
			if url, ok := infraResult.EncoderURLs["1"]; ok {
				os.Setenv("ENCODER_URL", url)
				fmt.Printf("✅ Using encoder v1 from testinfra: %s\n", url)
			}
			if url, ok := infraResult.EncoderURLs["2"]; ok {
				os.Setenv("ENCODER_V2_URL", url)
				fmt.Printf("✅ Using encoder v2 from testinfra: %s\n", url)
			}
		} else {
			fmt.Printf("⚠️ No EncoderURLs in infraResult\n")
		}

		// Pass batcher info from testinfra if available
		if infraResult.BatcherURL != "" {
			// Batcher doesn't expose a service URL, only metrics, so we flag its existence
			os.Setenv("BATCHER_PROVIDED", "true")
			fmt.Printf("✅ Using batcher from testinfra (metrics: %s)\n", infraResult.BatcherURL)
		} else {
			fmt.Printf("⚠️ No BatcherURL in infraResult\n")
		}

		// Pass controller info from testinfra if available
		if infraResult.ControllerMetricsURL != "" {
			// Controller is containerized, so we flag its existence
			os.Setenv("CONTROLLER_PROVIDED", "true")
			fmt.Printf("✅ Using controller from testinfra (metrics: %s)\n", infraResult.ControllerMetricsURL)
		} else {
			fmt.Printf("⚠️ No ControllerMetricsURL in infraResult\n")
		}

		// Pass operators info from testinfra if available
		if infraResult.OperatorAddresses != nil && len(infraResult.OperatorAddresses) > 0 {
			os.Setenv("OPERATORS_PROVIDED", "true")
			fmt.Printf("✅ Using %d operators from testinfra\n", len(infraResult.OperatorAddresses))
			for id, addr := range infraResult.OperatorAddresses {
				fmt.Printf("  - Operator %d: %s\n", id, addr)
			}
		} else {
			fmt.Printf("⚠️ No operators in infraResult\n")
		}

		// Pass relays info from testinfra if available
		if infraResult.RelayURLs != nil && len(infraResult.RelayURLs) > 0 {
			os.Setenv("RELAYS_PROVIDED", "true")
			fmt.Printf("✅ Using %d relays from testinfra\n", len(infraResult.RelayURLs))
			for id, url := range infraResult.RelayURLs {
				fmt.Printf("  - Relay %d: %s\n", id, url)
				// Set environment variables for each relay
				os.Setenv(fmt.Sprintf("RELAY_%d_URL", id), url)
			}
		} else {
			fmt.Printf("⚠️ No relays in infraResult\n")
		}

		// Generate all config variables for the binaries. Depends on the test config being set
		// with the output of the test infra deployment.
		// TODO: This generate variables method is very complex, we should simplify it.
		testConfig.GenerateAllVariables()

		fmt.Println("Starting EigenDA binaries")
		testConfig.StartBinaries()

		// Use cert verification components from testinfra
		Expect(infraResult.CertVerification).ToNot(BeNil(), "cert verification components must be initialized by testinfra")
		eigenDACertVerifierV1 = infraResult.CertVerification.EigenDACertVerifierV1
		certBuilder = infraResult.CertVerification.CertBuilder
		routerCertVerifier = infraResult.CertVerification.RouterCertVerifier
		staticCertVerifier = infraResult.CertVerification.StaticCertVerifier
		eigenDACertVerifierRouter = infraResult.CertVerification.EigenDACertVerifierRouter
		eigenDACertVerifierRouterCaller = infraResult.CertVerification.EigenDACertVerifierRouterCaller

		// Use retrieval clients from testinfra
		Expect(infraResult.RetrievalClients).ToNot(BeNil(), "retrieval clients must be initialized by testinfra")
		fmt.Println("Using retrieval clients from testinfra")
		ethClient = infraResult.RetrievalClients.EthClient
		rpcClient = infraResult.RetrievalClients.RPCClient
		retrievalClient = infraResult.RetrievalClients.RetrievalClient
		chainReader = infraResult.RetrievalClients.ChainReader
		relayRetrievalClientV2 = infraResult.RetrievalClients.RelayRetrievalClientV2
		validatorRetrievalClientV2 = infraResult.RetrievalClients.ValidatorRetrievalClientV2

		// Use payload disperser from testinfra
		Expect(infraResult.PayloadDisperser).ToNot(BeNil(), "payload disperser must be initialized by testinfra")
		fmt.Println("Using payload disperser from testinfra")
		payloadDisperser = infraResult.PayloadDisperser.PayloadDisperser.(*payloaddispersal.PayloadDisperser)
		deployerTransactorOpts = infraResult.PayloadDisperser.DeployerTransactorOpts.(*bind.TransactOpts)

	}
})

var _ = AfterSuite(func() {
	if testConfig.Environment.IsLocal() {
		if cancel != nil {
			cancel()
		}

		fmt.Println("Stopping binaries")
		testConfig.StopBinaries()

		// Force cleanup as a failsafe in case normal cleanup fails
		fmt.Println("Performing failsafe cleanup of any remaining processes")
		testConfig.ForceStopBinaries()

		if infraManager != nil {
			fmt.Println("Stopping testinfra containers")
			ctx, stopCancel := context.WithTimeout(context.Background(), 1*time.Minute)
			defer stopCancel()
			err := infraManager.Stop(ctx)
			if err != nil {
				fmt.Printf("Error stopping infrastructure: %v\n", err)
			}
		}
	}
})
