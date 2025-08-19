package integration_test

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients"
	clientsv2 "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/metrics"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloaddispersal"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloadretrieval"
	"github.com/Layr-Labs/eigenda/api/clients/v2/relay"
	validatorclientsv2 "github.com/Layr-Labs/eigenda/api/clients/v2/validator"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	routerbindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierRouter"
	verifierv1bindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierV1"

	"github.com/Layr-Labs/eigenda/core"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/Layr-Labs/eigenda/common/testinfra"
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

	testConfig         *deploy.Config
	infraManager       *testinfra.InfraManager
	infraResult        *testinfra.InfraResult

	metadataTableName               = "test-BlobMetadata"
	bucketTableName                 = "test-BucketStore"
	metadataTableNameV2             = "test-BlobMetadata-v2"
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
		// Start testinfra containers
		ctx, infraCancel := context.WithTimeout(context.Background(), 5*time.Minute)

		config := testinfra.DefaultConfig()
		// Enable graph node if subgraphs should be deployed
		deployer, ok := testConfig.GetDeployer(testConfig.EigenDA.Deployer)
		config.GraphNode.Enabled = ok && deployer.DeploySubgraphs

		if inMemoryBlobStore {
			fmt.Println("Using in-memory Blob Store - disabling LocalStack")
			config.LocalStack.Enabled = false
		} else {
			fmt.Println("Using shared Blob Store")
		}

		fmt.Println("Starting testinfra containers")
		manager, result, err := testinfra.StartCustom(ctx, config)
		Expect(err).To(BeNil())
		infraManager = manager
		infraResult = result
		cancel = infraCancel

		// Deploy AWS resources if using LocalStack
		if config.LocalStack.Enabled {
			localstack := infraManager.GetLocalStack()
			Expect(localstack).ToNot(BeNil())

			// Set environment variables for LocalStack compatibility
			localStackURL := localstack.Endpoint()
			_ = os.Setenv("AWS_ENDPOINT_URL", localStackURL)
			_ = os.Setenv("AWS_ACCESS_KEY_ID", "test")
			_ = os.Setenv("AWS_SECRET_ACCESS_KEY", "test")
			_ = os.Setenv("AWS_DEFAULT_REGION", "us-east-1")

			// Extract port from LocalStack URL for compatibility with existing deploy function
			parts := strings.Split(localStackURL, ":")
			port := "4566" // default fallback
			if len(parts) >= 3 {
				port = parts[2]
			}

			// Create a temporary dockertest pool for compatibility
			// The actual pool parameter is not used when we pass nil
			err = deploy.DeployResources(nil, port, metadataTableName, bucketTableName, metadataTableNameV2)
			Expect(err).To(BeNil())
		}

		// Wait for infrastructure to be fully ready
		fmt.Println("Waiting for testinfra containers to be ready")
		err = infraManager.WaitForReady(ctx)
		Expect(err).To(BeNil())

		// Test Anvil connectivity before proceeding
		anvil := infraManager.GetAnvil()
		Expect(anvil).ToNot(BeNil())
		fmt.Printf("Testing Anvil connectivity at %s\n", infraResult.AnvilRPC)
		
		// Debug: Check container logs
		container := anvil.GetContainer()
		logs, err := container.Logs(ctx)
		if err == nil {
			logData := make([]byte, 2048)
			n, _ := logs.Read(logData)
			fmt.Printf("Anvil container logs: %s\n", string(logData[:n]))
			logs.Close()
		}

		// Test JSON-RPC health check
		ctx2, cancel2 := context.WithTimeout(ctx, 30*time.Second)
		defer cancel2()
		for {
				// Use proper JSON-RPC call instead of plain HTTP GET
			jsonRPCPayload := `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`
			resp, err := http.Post(infraResult.AnvilRPC, "application/json", strings.NewReader(jsonRPCPayload))
			if err == nil && resp.StatusCode == 200 {
				resp.Body.Close()
				fmt.Println("Anvil is responding to JSON-RPC requests")
				break
			}
			if resp != nil {
				resp.Body.Close()
			}
			select {
			case <-ctx2.Done():
				fmt.Printf("Timeout waiting for Anvil to respond: %v\n", err)
				Expect(err).To(BeNil())
				return
			case <-time.After(time.Second):
				fmt.Printf("Anvil not ready yet, retrying... (%v)\n", err)
				continue
			}
		}
		
		// Update test config to use testinfra endpoints
		fmt.Printf("Updating RPC URL from %s to %s\n", testConfig.Deployers[0].RPC, infraResult.AnvilRPC)
		testConfig.Deployers[0].RPC = infraResult.AnvilRPC

		// Update Graph URLs if Graph Node is enabled
		var graphURL string
		if config.GraphNode.Enabled && infraResult.GraphNodeURL != "" {
			graphURL = infraResult.GraphNodeURL + "/subgraphs/name/Layr-Labs/eigenda-operator-state"
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

			// Test Graph Node connectivity before proceeding
			fmt.Println("Testing Graph Node connectivity...")
			err := testGraphNodeConnectivity(graphURL, 10, 2*time.Second)
			if err != nil {
				fmt.Printf("⚠️  Graph Node connectivity test failed: %v\n", err)
				fmt.Printf("📋 Debug info:\n")
				fmt.Printf("   - GraphURL: %s\n", graphURL)
				fmt.Printf("   - AdminURL: %s\n", infraResult.GraphNodeAdminURL)
				fmt.Printf("   - IPFS URL: %s\n", infraResult.IPFSURL)
				// Don't fail the test, but provide clear diagnostic info
			}
		}

		loggerConfig := common.DefaultLoggerConfig()
		logger, err = common.NewLogger(loggerConfig)
		Expect(err).To(BeNil())

		fmt.Println("Deploying experiment")
		testConfig.DeployExperiment()
		
		// After contract and subgraph deployment, test connectivity again
		if graphURL != "" {
			fmt.Println("Testing Graph Node connectivity after subgraph deployment...")
			err := testGraphNodeConnectivity(graphURL, 15, 3*time.Second)
			if err != nil {
				fmt.Printf("⚠️  Graph Node still not accessible after deployment: %v\n", err)
			} else {
				fmt.Println("✅ Graph Node confirmed working after deployment")
			}
		}
		pk := testConfig.Pks.EcdsaMap[deployer.Name].PrivateKey
		pk = strings.TrimPrefix(pk, "0x")
		pk = strings.TrimPrefix(pk, "0X")
		ethClient, err = geth.NewMultiHomingClient(geth.EthClientConfig{
			RPCURLs:          []string{testConfig.Deployers[0].RPC},
			PrivateKeyString: pk,
			NumConfirmations: numConfirmations,
			NumRetries:       numRetries,
		}, gethcommon.Address{}, logger)
		Expect(err).To(BeNil())

		rpcClient, err = ethrpc.Dial(testConfig.Deployers[0].RPC)
		Expect(err).To(BeNil())

		fmt.Println("Registering blob versions and relays")
		testConfig.RegisterBlobVersionAndRelays(ethClient)

		fmt.Println("Registering disperser keypair")
		err = testConfig.RegisterDisperserKeypair(ethClient)
		if err != nil {
			panic(err)
		}

		fmt.Println("Starting binaries")
		testConfig.StartBinaries()

		eigenDACertVerifierV1, err = verifierv1bindings.NewContractEigenDACertVerifierV1(gethcommon.HexToAddress(testConfig.EigenDAV1CertVerifier), ethClient)
		Expect(err).To(BeNil())
		err = setupRetrievalClients(testConfig)
		Expect(err).To(BeNil())

		fmt.Println("Building client verification and interaction components")

		certBuilder, err = clientsv2.NewCertBuilder(
			logger,
			gethcommon.HexToAddress(testConfig.EigenDA.OperatorStateRetriever),
			gethcommon.HexToAddress(testConfig.EigenDA.RegistryCoordinator),
			ethClient,
		)

		Expect(err).To(BeNil())

		routerAddressProvider, err := verification.BuildRouterAddressProvider(
			gethcommon.HexToAddress(testConfig.EigenDA.CertVerifierRouter),
			ethClient,
			logger)

		Expect(err).To(BeNil())

		staticAddressProvider := verification.NewStaticCertVerifierAddressProvider(
			gethcommon.HexToAddress(testConfig.EigenDA.CertVerifier))

		Expect(err).To(BeNil())

		staticCertVerifier, err = verification.NewCertVerifier(
			logger,
			ethClient,
			staticAddressProvider)

		Expect(err).To(BeNil())

		routerCertVerifier, err = verification.NewCertVerifier(
			logger,
			ethClient,
			routerAddressProvider)

		Expect(err).To(BeNil())

		eigenDACertVerifierRouter, err = routerbindings.NewContractEigenDACertVerifierRouterTransactor(gethcommon.HexToAddress(testConfig.EigenDA.CertVerifierRouter), ethClient)
		Expect(err).To(BeNil())

		eigenDACertVerifierRouterCaller, err = routerbindings.NewContractEigenDACertVerifierRouterCaller(gethcommon.HexToAddress(testConfig.EigenDA.CertVerifierRouter), ethClient)
		Expect(err).To(BeNil())

		chainID, err := ethClient.ChainID(context.Background())
		Expect(err).To(BeNil())

		deployerTransactorOpts = newTransactOptsFromPrivateKey(pk, chainID)

		err = setupPayloadDisperserWithRouter()
		Expect(err).To(BeNil())

	}
})

func setupPayloadDisperserWithRouter() error {
	// Set up the block monitor
	blockMonitor, err := verification.NewBlockNumberMonitor(logger, ethClient, time.Second*1)
	if err != nil {
		return err
	}

	// Set up the PayloadDisperser
	privateKeyHex := "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcded"
	signer, err := auth.NewLocalBlobRequestSigner(privateKeyHex)
	if err != nil {
		return err
	}

	disperserClientConfig := &clientsv2.DisperserClientConfig{
		Hostname: "localhost",
		Port:     "32005",
	}

	accountId, err := signer.GetAccountID()
	if err != nil {
		return fmt.Errorf("error getting account ID: %w", err)
	}

	accountant := clientsv2.NewAccountant(
		accountId,
		nil,
		nil,
		0,
		0,
		0,
		0,
		metrics.NoopAccountantMetrics,
	)
	disperserClient, err := clientsv2.NewDisperserClient(disperserClientConfig, signer, nil, accountant)
	if err != nil {
		return err
	}

	payloadDisperserConfig := payloaddispersal.PayloadDisperserConfig{
		PayloadClientConfig:    *clientsv2.GetDefaultPayloadClientConfig(),
		DisperseBlobTimeout:    2 * time.Minute,
		BlobCompleteTimeout:    2 * time.Minute,
		BlobStatusPollInterval: 1 * time.Second,
		ContractCallTimeout:    5 * time.Second,
	}

	payloadDisperser, err = payloaddispersal.NewPayloadDisperser(
		logger,
		payloadDisperserConfig,
		disperserClient,
		blockMonitor,
		certBuilder,
		routerCertVerifier,
		nil,
	)

	return err
}

func newTransactOptsFromPrivateKey(privateKeyHex string, chainID *big.Int) *bind.TransactOpts {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("invalid private key: %v", err)
	}

	opts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("failed to create transactor: %v", err)
	}

	return opts
}

func setupRetrievalClients(testConfig *deploy.Config) error {
	ethClientConfig := geth.EthClientConfig{
		RPCURLs:          []string{testConfig.Deployers[0].RPC},
		PrivateKeyString: "351b8eca372e64f64d514f90f223c5c4f86a04ff3dcead5c27293c547daab4ca", // just random private key
		NumConfirmations: numConfirmations,
		NumRetries:       numRetries,
	}
	var err error
	if ethClient == nil {
		ethClient, err = geth.NewMultiHomingClient(ethClientConfig, gethcommon.Address{}, logger)
		if err != nil {
			return err
		}
	}
	if rpcClient == nil {
		rpcClient, err = ethrpc.Dial(testConfig.Deployers[0].RPC)
		if err != nil {
			log.Fatalln("could not start tcp listener", err)
		}
	}
	tx, err := eth.NewWriter(
		logger, ethClient, testConfig.EigenDA.OperatorStateRetriever, testConfig.EigenDA.ServiceManager)
	if err != nil {
		return err
	}

	cs := eth.NewChainState(tx, ethClient)
	agn := &core.StdAssignmentCoordinator{}
	nodeClient := clients.NewNodeClient(20 * time.Second)
	srsOrder, err := strconv.Atoi(testConfig.Retriever.RETRIEVER_SRS_ORDER)
	if err != nil {
		return err
	}
	kzgConfig := &kzg.KzgConfig{
		G1Path:          testConfig.Retriever.RETRIEVER_G1_PATH,
		G2Path:          testConfig.Retriever.RETRIEVER_G2_PATH,
		G2PowerOf2Path:  testConfig.Retriever.RETRIEVER_G2_POWER_OF_2_PATH,
		CacheDir:        testConfig.Retriever.RETRIEVER_CACHE_PATH,
		SRSOrder:        uint64(srsOrder),
		SRSNumberToLoad: uint64(srsOrder),
		NumWorker:       1,
		PreloadEncoder:  false,
		LoadG2Points:    true,
	}

	kzgVerifier, err := verifier.NewVerifier(kzgConfig, nil)
	if err != nil {
		return err
	}

	retrievalClient, err = clients.NewRetrievalClient(logger, cs, agn, nodeClient, kzgVerifier, 10)
	if err != nil {
		return err
	}
	chainReader, err = eth.NewReader(
		logger,
		ethClient,
		testConfig.EigenDA.OperatorStateRetriever,
		testConfig.EigenDA.ServiceManager,
	)
	if err != nil {
		return err
	}

	clientConfig := validatorclientsv2.DefaultClientConfig()
	retrievalClientV2 := validatorclientsv2.NewValidatorClient(logger, chainReader, cs, kzgVerifier, clientConfig, nil)

	validatorPayloadRetrieverConfig := payloadretrieval.ValidatorPayloadRetrieverConfig{
		PayloadClientConfig: *clientsv2.GetDefaultPayloadClientConfig(),
		RetrievalTimeout:    1 * time.Minute,
	}

	validatorRetrievalClientV2, err = payloadretrieval.NewValidatorPayloadRetriever(
		logger,
		validatorPayloadRetrieverConfig,
		retrievalClientV2,
		kzgVerifier.Srs.G1)

	if err != nil {
		return err
	}

	relayClientConfig := &relay.RelayClientConfig{
		MaxGRPCMessageSize: 100 * 1024 * 1024, // 100 MB message size limit,
	}

	relayUrlProvider, err := relay.NewRelayUrlProvider(ethClient, chainReader.GetRelayRegistryAddress())
	if err != nil {
		return err
	}

	relayClient, err := relay.NewRelayClient(relayClientConfig, logger, relayUrlProvider)
	if err != nil {
		return err
	}

	relayPayloadRetrieverConfig := payloadretrieval.RelayPayloadRetrieverConfig{
		PayloadClientConfig: *clientsv2.GetDefaultPayloadClientConfig(),
		RelayTimeout:        5 * time.Second,
	}

	relayRetrievalClientV2, err = payloadretrieval.NewRelayPayloadRetriever(
		logger,
		rand.New(rand.NewSource(time.Now().UnixNano())),
		relayPayloadRetrieverConfig,
		relayClient,
		kzgVerifier.Srs.G1)

	return err
}

func testGraphNodeConnectivity(graphURL string, maxRetries int, retryInterval time.Duration) error {
	client := &http.Client{Timeout: 10 * time.Second}
	
	// Test queries - try multiple approaches
	testQueries := []struct {
		name  string
		query string
		url   string
	}{
		{"GraphQL root", `{"query": "{__schema{queryType{name}}}"}`, "/graphql"},
		{"Subgraph meta", `{"query": "{_meta{block{number}}}"}`, "/graphql"},
		{"Health check", "", "/"},
	}
	
	for i := 0; i < maxRetries; i++ {
		fmt.Printf("Testing Graph Node connectivity (attempt %d/%d)...\n", i+1, maxRetries)
		
		for _, test := range testQueries {
			testURL := graphURL + test.url
			fmt.Printf("  Testing %s at %s\n", test.name, testURL)
			
			var req *http.Request
			var err error
			
			if test.query != "" {
				req, err = http.NewRequest("POST", testURL, strings.NewReader(test.query))
				if err != nil {
					fmt.Printf("    ❌ Failed to create request: %v\n", err)
					continue
				}
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, err = http.NewRequest("GET", testURL, nil)
				if err != nil {
					fmt.Printf("    ❌ Failed to create request: %v\n", err)
					continue
				}
			}
			
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("    ❌ Connection failed: %v\n", err)
				continue
			}
			defer resp.Body.Close()
			
			body, _ := io.ReadAll(resp.Body)
			bodyPreview := string(body)
			if len(bodyPreview) > 100 {
				bodyPreview = bodyPreview[:100] + "..."
			}
			fmt.Printf("    Status: %d, Body preview: %s\n", resp.StatusCode, bodyPreview)
			
			if resp.StatusCode == 200 {
				fmt.Printf("✅ Graph Node is accessible via %s\n", test.name)
				return nil
			}
		}
		
		if i < maxRetries-1 {
			fmt.Printf("⏱️  Waiting %v before retry...\n", retryInterval)
			time.Sleep(retryInterval)
		}
	}
	
	return fmt.Errorf("graph node at %s is not accessible after %d attempts", graphURL, maxRetries)
}

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
