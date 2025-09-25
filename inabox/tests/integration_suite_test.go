package integration_test

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"math/big"
	"math/rand"
	"net"
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
	pb "github.com/Layr-Labs/eigenda/api/grpc/churner"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	routerbindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierRouter"
	verifierv1bindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierV1"
	"github.com/Layr-Labs/eigenda/core"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	coreeth "github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	verifierv2 "github.com/Layr-Labs/eigenda/encoding/kzg/verifier/v2"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigenda/operators/churner"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/Layr-Labs/eigenda/test/testbed"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"google.golang.org/grpc"
)

// TestHarness contains all test infrastructure needed for integration tests
type TestHarness struct {
	// Infrastructure containers
	AnvilContainer      *testbed.AnvilContainer
	GraphNodeContainer  *testbed.GraphNodeContainer
	LocalstackContainer *testbed.LocalStackContainer
	ChainDockerNetwork  *testcontainers.DockerNetwork

	// Churner components
	ChurnerServer   *grpc.Server
	ChurnerListener net.Listener

	// Test configuration
	TestConfig          *deploy.Config
	TestName            string
	TemplateName        string
	InMemoryBlobStore   bool
	LocalStackPort      string
	MetadataTableName   string
	BucketTableName     string
	MetadataTableNameV2 string

	// Ethereum clients and contracts
	EthClient                       common.EthClient
	RPCClient                       common.RPCEthClient
	CertBuilder                     *clientsv2.CertBuilder
	RouterCertVerifier              *verification.CertVerifier
	StaticCertVerifier              *verification.CertVerifier
	EigenDACertVerifierRouter       *routerbindings.ContractEigenDACertVerifierRouterTransactor
	EigenDACertVerifierRouterCaller *routerbindings.ContractEigenDACertVerifierRouterCaller
	EigenDACertVerifierV1           *verifierv1bindings.ContractEigenDACertVerifierV1
	DeployerTransactorOpts          *bind.TransactOpts

	// Retrieval clients
	RetrievalClient            clients.RetrievalClient
	RelayRetrievalClientV2     *payloadretrieval.RelayPayloadRetriever
	ValidatorRetrievalClientV2 *payloadretrieval.ValidatorPayloadRetriever
	PayloadDisperser           *payloaddispersal.PayloadDisperser

	// Core components
	ChainReader core.Reader

	// Context management
	Cancel context.CancelFunc

	// Configuration constants
	NumConfirmations int
	NumRetries       int

	// Logger
	Logger logging.Logger
}

// Global test harness instance - will be initialized once per test suite
var testHarness *TestHarness

// Configuration constants from command line flags
var (
	templateName      string
	testName          string
	inMemoryBlobStore bool
	logger            = test.GetLogger()
)

func init() {
	flag.StringVar(&templateName, "config", "testconfig-anvil.yaml", "Name of the config file (in `inabox/templates`)")
	flag.StringVar(&testName, "testname", "", "Name of the test (in `inabox/testdata`)")
	flag.BoolVar(&inMemoryBlobStore, "inMemoryBlobStore", false, "whether to use in-memory blob store")
}

func TestMain(m *testing.M) {
	flag.Parse()

	if testing.Short() {
		logger.Info("Skipping inabox integration tests in short mode")
		os.Exit(0)
	}

	// Run suite setup
	if err := setupSuite(); err != nil {
		logger.Error("Setup failed:", err)
		teardownSuite()
		os.Exit(1)
	}

	// Run all tests
	code := m.Run()

	// Run suite teardown
	teardownSuite()

	// Exit with test result code
	os.Exit(code)
}

func setupSuite() error {
	logger.Info("bootstrapping test environment")

	// Initialize the test harness
	testHarness = &TestHarness{
		TestName:            testName,
		TemplateName:        templateName,
		InMemoryBlobStore:   inMemoryBlobStore,
		Logger:              logger,
		NumConfirmations:    1,
		NumRetries:          5,
		MetadataTableName:   "test-BlobMetadata",
		BucketTableName:     "test-BucketStore",
		MetadataTableNameV2: "test-BlobMetadata-v2",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	rootPath := "../../"

	var err error
	if testHarness.TestName == "" {
		testHarness.TestName, err = deploy.CreateNewTestDirectory(templateName, rootPath)
		if err != nil {
			return fmt.Errorf("failed to create test directory: %w", err)
		}
	}

	testHarness.TestConfig = deploy.ReadTestConfig(testHarness.TestName, rootPath)

	if testHarness.TestConfig.Environment.IsLocal() {
		// Create a shared Docker network for all containers
		testHarness.ChainDockerNetwork, err = network.New(context.Background(),
			network.WithDriver("bridge"),
			network.WithAttachable())
		if err != nil {
			return fmt.Errorf("failed to create docker network: %w", err)
		}
		logger.Info("Created Docker network", "name", testHarness.ChainDockerNetwork.Name)

		if !testHarness.InMemoryBlobStore {
			logger.Info("Using shared Blob Store")
			testHarness.LocalStackPort = "4570"
			// Use the timeout context for container creation
			testHarness.LocalstackContainer, err = testbed.NewLocalStackContainerWithOptions(ctx, testbed.LocalStackOptions{
				ExposeHostPort: true,
				HostPort:       testHarness.LocalStackPort,
				Logger:         logger,
				Network:        testHarness.ChainDockerNetwork,
			})
			if err != nil {
				return fmt.Errorf("failed to start localstack: %w", err)
			}

			deployConfig := testbed.DeployResourcesConfig{
				LocalStackEndpoint:  fmt.Sprintf("http://0.0.0.0:%s", testHarness.LocalStackPort),
				MetadataTableName:   testHarness.MetadataTableName,
				BucketTableName:     testHarness.BucketTableName,
				V2MetadataTableName: testHarness.MetadataTableNameV2,
				Logger:              logger,
			}
			err = testbed.DeployResources(ctx, deployConfig)
			if err != nil {
				return fmt.Errorf("failed to deploy resources: %w", err)
			}
		} else {
			logger.Info("Using in-memory Blob Store")
		}

		logger.Info("Starting anvil")
		testHarness.AnvilContainer, err = testbed.NewAnvilContainerWithOptions(ctx, testbed.AnvilOptions{
			ExposeHostPort: true,
			HostPort:       "8545",
			Logger:         logger,
			Network:        testHarness.ChainDockerNetwork,
		})
		if err != nil {
			return fmt.Errorf("failed to start anvil: %w", err)
		}
		anvilInternalEndpoint := testHarness.AnvilContainer.InternalEndpoint()
		logger.Info("Anvil RPC URL", "url", testHarness.AnvilContainer.RpcURL(), "internal", anvilInternalEndpoint)

		deployer, ok := testHarness.TestConfig.GetDeployer(testHarness.TestConfig.EigenDA.Deployer)
		if ok && deployer.DeploySubgraphs {
			logger.Info("Starting graph node")
			testHarness.GraphNodeContainer, err = testbed.NewGraphNodeContainerWithOptions(context.Background(), testbed.GraphNodeOptions{
				PostgresDB:     "graph-node",
				PostgresUser:   "graph-node",
				PostgresPass:   "let-me-in",
				EthereumRPC:    anvilInternalEndpoint,
				ExposeHostPort: true,
				HostHTTPPort:   "8000",
				HostWSPort:     "8001",
				HostAdminPort:  "8020",
				HostIPFSPort:   "5001",
				Logger:         logger,
				Network:        testHarness.ChainDockerNetwork,
			})
			if err != nil {
				return fmt.Errorf("failed to start graph node: %w", err)
			}
		}

		logger.Info("Deploying experiment")
		testHarness.TestConfig.DeployExperiment()
		pk := testHarness.TestConfig.Pks.EcdsaMap[deployer.Name].PrivateKey
		pk = strings.TrimPrefix(pk, "0x")
		pk = strings.TrimPrefix(pk, "0X")
		testHarness.EthClient, err = geth.NewMultiHomingClient(geth.EthClientConfig{
			RPCURLs:          []string{testHarness.TestConfig.Deployers[0].RPC},
			PrivateKeyString: pk,
			NumConfirmations: testHarness.NumConfirmations,
			NumRetries:       testHarness.NumRetries,
		}, gethcommon.Address{}, logger)
		if err != nil {
			return fmt.Errorf("failed to create eth client: %w", err)
		}

		testHarness.RPCClient, err = ethrpc.Dial(testHarness.TestConfig.Deployers[0].RPC)
		if err != nil {
			return fmt.Errorf("failed to create rpc client: %w", err)
		}

		// Force foundry to mine a block since it isn't auto-mining
		err = testHarness.RPCClient.CallContext(ctx, nil, "evm_mine")
		if err != nil {
			return fmt.Errorf("failed to mine block: %w", err)
		}

		logger.Info("Starting churner server")
		err = startChurnerServer(testHarness)
		if err != nil {
			return fmt.Errorf("failed to start churner server: %w", err)
		}
		logger.Info("Churner server started", "port", "32002")

		logger.Info("Starting binaries")
		testHarness.TestConfig.StartBinaries(true) // true = for tests, will skip churner

		testHarness.EigenDACertVerifierV1, err = verifierv1bindings.NewContractEigenDACertVerifierV1(gethcommon.HexToAddress(testHarness.TestConfig.EigenDAV1CertVerifier), testHarness.EthClient)
		if err != nil {
			return fmt.Errorf("failed to create EigenDA cert verifier V1: %w", err)
		}
		err = setupRetrievalClients(testHarness)
		if err != nil {
			return fmt.Errorf("failed to setup retrieval clients: %w", err)
		}

		logger.Info("Building client verification and interaction components")

		testHarness.CertBuilder, err = clientsv2.NewCertBuilder(
			logger,
			gethcommon.HexToAddress(testHarness.TestConfig.EigenDA.OperatorStateRetriever),
			gethcommon.HexToAddress(testHarness.TestConfig.EigenDA.RegistryCoordinator),
			testHarness.EthClient,
		)

		if err != nil {
			return fmt.Errorf("failed to create cert builder: %w", err)
		}

		routerAddressProvider, err := verification.BuildRouterAddressProvider(
			gethcommon.HexToAddress(testHarness.TestConfig.EigenDA.CertVerifierRouter),
			testHarness.EthClient,
			logger)

		if err != nil {
			return fmt.Errorf("failed to build router address provider: %w", err)
		}

		staticAddressProvider := verification.NewStaticCertVerifierAddressProvider(
			gethcommon.HexToAddress(testHarness.TestConfig.EigenDA.CertVerifier))

		// No error to check for NewStaticCertVerifierAddressProvider

		testHarness.StaticCertVerifier, err = verification.NewCertVerifier(
			logger,
			testHarness.EthClient,
			staticAddressProvider)

		if err != nil {
			return fmt.Errorf("failed to create static cert verifier: %w", err)
		}

		testHarness.RouterCertVerifier, err = verification.NewCertVerifier(
			logger,
			testHarness.EthClient,
			routerAddressProvider)

		if err != nil {
			return fmt.Errorf("failed to create router cert verifier: %w", err)
		}

		testHarness.EigenDACertVerifierRouter, err = routerbindings.NewContractEigenDACertVerifierRouterTransactor(gethcommon.HexToAddress(testHarness.TestConfig.EigenDA.CertVerifierRouter), testHarness.EthClient)
		if err != nil {
			return fmt.Errorf("failed to create router transactor: %w", err)
		}

		testHarness.EigenDACertVerifierRouterCaller, err = routerbindings.NewContractEigenDACertVerifierRouterCaller(gethcommon.HexToAddress(testHarness.TestConfig.EigenDA.CertVerifierRouter), testHarness.EthClient)
		if err != nil {
			return fmt.Errorf("failed to create router caller: %w", err)
		}

		chainID, err := testHarness.EthClient.ChainID(ctx)
		if err != nil {
			return fmt.Errorf("failed to get chain ID: %w", err)
		}

		testHarness.DeployerTransactorOpts = newTransactOptsFromPrivateKey(pk, chainID)

		err = setupPayloadDisperserWithRouter(testHarness)
		if err != nil {
			return fmt.Errorf("failed to setup payload disperser: %w", err)
		}

	}
	return nil
}

func setupPayloadDisperserWithRouter(harness *TestHarness) error {
	// Set up the block monitor
	blockMonitor, err := verification.NewBlockNumberMonitor(harness.Logger, harness.EthClient, time.Second*1)
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
	disperserClient, err := clientsv2.NewDisperserClient(
		harness.Logger,
		disperserClientConfig,
		signer,
		nil, // no prover so will query disperser for generating commitments
		accountant,
		metrics.NoopDispersalMetrics,
	)
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

	harness.PayloadDisperser, err = payloaddispersal.NewPayloadDisperser(
		harness.Logger,
		payloadDisperserConfig,
		disperserClient,
		blockMonitor,
		harness.CertBuilder,
		harness.RouterCertVerifier,
		nil,
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

func setupRetrievalClients(harness *TestHarness) error {
	ethClientConfig := geth.EthClientConfig{
		RPCURLs:          []string{harness.TestConfig.Deployers[0].RPC},
		PrivateKeyString: "351b8eca372e64f64d514f90f223c5c4f86a04ff3dcead5c27293c547daab4ca", // just random private key
		NumConfirmations: harness.NumConfirmations,
		NumRetries:       harness.NumRetries,
	}
	var err error
	if harness.EthClient == nil {
		harness.EthClient, err = geth.NewMultiHomingClient(ethClientConfig, gethcommon.Address{}, harness.Logger)
		if err != nil {
			return err
		}
	}
	if harness.RPCClient == nil {
		harness.RPCClient, err = ethrpc.Dial(harness.TestConfig.Deployers[0].RPC)
		if err != nil {
			log.Fatalln("could not start tcp listener", err)
		}
	}
	tx, err := coreeth.NewWriter(
		harness.Logger, harness.EthClient, harness.TestConfig.EigenDA.OperatorStateRetriever, harness.TestConfig.EigenDA.ServiceManager)
	if err != nil {
		return err
	}

	cs := coreeth.NewChainState(tx, harness.EthClient)
	agn := &core.StdAssignmentCoordinator{}
	nodeClient := clients.NewNodeClient(20 * time.Second)
	srsOrder, err := strconv.Atoi(harness.TestConfig.Retriever.RETRIEVER_SRS_ORDER)
	if err != nil {
		return err
	}
	kzgConfig := &kzg.KzgConfig{
		G1Path:          harness.TestConfig.Retriever.RETRIEVER_G1_PATH,
		G2Path:          harness.TestConfig.Retriever.RETRIEVER_G2_PATH,
		CacheDir:        harness.TestConfig.Retriever.RETRIEVER_CACHE_PATH,
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

	harness.RetrievalClient, err = clients.NewRetrievalClient(harness.Logger, cs, agn, nodeClient, kzgVerifier, 10)
	if err != nil {
		return err
	}
	harness.ChainReader, err = coreeth.NewReader(
		harness.Logger,
		harness.EthClient,
		harness.TestConfig.EigenDA.OperatorStateRetriever,
		harness.TestConfig.EigenDA.ServiceManager,
	)
	if err != nil {
		return err
	}

	kzgVerifierV2, err := verifierv2.NewVerifier(verifierv2.KzgConfigFromV1Config(kzgConfig), nil)
	if err != nil {
		return fmt.Errorf("new verifier v2: %w", err)
	}

	clientConfig := validatorclientsv2.DefaultClientConfig()
	retrievalClientV2 := validatorclientsv2.NewValidatorClient(harness.Logger, harness.ChainReader, cs, kzgVerifierV2, clientConfig, nil)

	validatorPayloadRetrieverConfig := payloadretrieval.ValidatorPayloadRetrieverConfig{
		PayloadClientConfig: *clientsv2.GetDefaultPayloadClientConfig(),
		RetrievalTimeout:    1 * time.Minute,
	}

	harness.ValidatorRetrievalClientV2, err = payloadretrieval.NewValidatorPayloadRetriever(
		harness.Logger,
		validatorPayloadRetrieverConfig,
		retrievalClientV2,
		kzgVerifier.G1SRS,
		metrics.NoopRetrievalMetrics)

	if err != nil {
		return err
	}

	relayClientConfig := &relay.RelayClientConfig{
		MaxGRPCMessageSize: 100 * 1024 * 1024, // 100 MB message size limit,
	}

	relayUrlProvider, err := relay.NewRelayUrlProvider(harness.EthClient, harness.ChainReader.GetRelayRegistryAddress())
	if err != nil {
		return err
	}

	relayClient, err := relay.NewRelayClient(relayClientConfig, harness.Logger, relayUrlProvider)
	if err != nil {
		return err
	}

	relayPayloadRetrieverConfig := payloadretrieval.RelayPayloadRetrieverConfig{
		PayloadClientConfig: *clientsv2.GetDefaultPayloadClientConfig(),
		RelayTimeout:        5 * time.Second,
	}

	harness.RelayRetrievalClientV2, err = payloadretrieval.NewRelayPayloadRetriever(
		harness.Logger,
		rand.New(rand.NewSource(time.Now().UnixNano())),
		relayPayloadRetrieverConfig,
		relayClient,
		kzgVerifier.G1SRS,
		metrics.NoopRetrievalMetrics)

	return err
}

func startChurnerServer(harness *TestHarness) error {
	// Get deployer's private key
	var privateKey string
	deployer, ok := harness.TestConfig.GetDeployer(harness.TestConfig.EigenDA.Deployer)
	if ok && deployer.Name != "" {
		privateKey = strings.TrimPrefix(harness.TestConfig.Pks.EcdsaMap[deployer.Name].PrivateKey, "0x")
	}

	// Create churner config directly
	// Create logs directory in the test directory
	logsDir := fmt.Sprintf("testdata/%s/logs", harness.TestName)
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	logFilePath := fmt.Sprintf("%s/churner.log", logsDir)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open churner log file: %w", err)
	}

	churnerConfig := &churner.Config{
		EthClientConfig: geth.EthClientConfig{
			RPCURLs:          []string{"http://localhost:8545"},
			PrivateKeyString: privateKey,
		},
		LoggerConfig: common.LoggerConfig{
			Format:       common.TextLogFormat,
			OutputWriter: io.MultiWriter(os.Stdout, logFile),
			HandlerOpts: logging.SLoggerOptions{
				Level:   slog.LevelDebug,
				NoColor: true,
			},
		},
		MetricsConfig: churner.MetricsConfig{
			HTTPPort:      "9095",
			EnableMetrics: true,
		},
		OperatorStateRetrieverAddr: harness.TestConfig.EigenDA.OperatorStateRetriever,
		EigenDAServiceManagerAddr:  harness.TestConfig.EigenDA.ServiceManager,
		EigenDADirectory:           harness.TestConfig.EigenDA.EigenDADirectory,
		GRPCPort:                   "32002",
		ChurnApprovalInterval:      15 * time.Minute,
		PerPublicKeyRateLimit:      1 * time.Second,
	}

	// Set graph URL if graph node is enabled
	if deployer.DeploySubgraphs && harness.GraphNodeContainer != nil {
		churnerConfig.ChainStateConfig = thegraph.Config{
			Endpoint: "http://localhost:8000/subgraphs/name/Layr-Labs/eigenda-operator-state",
		}
	}

	// Create a logger from the churner config with file output
	churnerLogger, err := common.NewLogger(&churnerConfig.LoggerConfig)
	if err != nil {
		return fmt.Errorf("failed to create churner logger: %w", err)
	}

	// Create geth client
	gethClient, err := geth.NewMultiHomingClient(churnerConfig.EthClientConfig, gethcommon.Address{}, churnerLogger)
	if err != nil {
		return fmt.Errorf("failed to create geth client: %w", err)
	}

	// Create writer
	churnerTx, err := coreeth.NewWriter(
		churnerLogger,
		gethClient,
		churnerConfig.OperatorStateRetrieverAddr,
		churnerConfig.EigenDAServiceManagerAddr)
	if err != nil {
		return fmt.Errorf("failed to create writer: %w", err)
	}

	// Create indexer
	chainState := coreeth.NewChainState(churnerTx, gethClient)
	indexer := thegraph.MakeIndexedChainState(churnerConfig.ChainStateConfig, chainState, churnerLogger)

	// Create churner
	churnerMetrics := churner.NewMetrics(churnerConfig.MetricsConfig.HTTPPort, churnerLogger)
	churnerInstance, err := churner.NewChurner(churnerConfig, indexer, churnerTx, churnerLogger, churnerMetrics)
	if err != nil {
		return fmt.Errorf("failed to create churner: %w", err)
	}

	// Create churner server
	churnerSvr := churner.NewServer(churnerConfig, churnerInstance, churnerLogger, churnerMetrics)
	err = churnerSvr.Start(churnerConfig.MetricsConfig)
	if err != nil {
		return fmt.Errorf("failed to start churner server metrics: %w", err)
	}

	// Create listener only after churner is successfully created
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", churnerConfig.GRPCPort))
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", churnerConfig.GRPCPort, err)
	}
	harness.ChurnerListener = listener

	// Create and start gRPC server
	harness.ChurnerServer = grpc.NewServer(grpc.MaxRecvMsgSize(1024 * 1024 * 300))
	pb.RegisterChurnerServer(harness.ChurnerServer, churnerSvr)
	healthcheck.RegisterHealthServer(pb.Churner_ServiceDesc.ServiceName, harness.ChurnerServer)

	// Start serving in goroutine
	go func() {
		churnerLogger.Info("Starting churner gRPC server", "port", churnerConfig.GRPCPort)
		if err := harness.ChurnerServer.Serve(harness.ChurnerListener); err != nil {
			churnerLogger.Info("Churner gRPC server stopped", "error", err)
		}
	}()

	// TODO: Replace with proper health check endpoint instead of fixed sleep
	time.Sleep(100 * time.Millisecond)
	churnerLogger.Info("Churner server started successfully", "port", churnerConfig.GRPCPort, "logFile", logFilePath)

	return nil
}

func teardownSuite() {
	logger.Info("Tearing down test environment")

	if testHarness == nil || testHarness.TestConfig == nil || !testHarness.TestConfig.Environment.IsLocal() {
		return
	}

	ctx, teardownCancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer teardownCancel()

	if testHarness.Cancel != nil {
		testHarness.Cancel()
	}

	logger.Info("Stopping binaries")
	testHarness.TestConfig.StopBinaries()

	if testHarness.ChurnerServer != nil {
		logger.Info("Stopping churner server")
		testHarness.ChurnerServer.GracefulStop()
		if testHarness.ChurnerListener != nil {
			_ = testHarness.ChurnerListener.Close()
		}
	}

	logger.Info("Stopping anvil")
	if testHarness.AnvilContainer != nil {
		if err := testHarness.AnvilContainer.Terminate(ctx); err != nil {
			logger.Warn("Failed to terminate anvil container", "error", err)
		}
	}

	if testHarness.GraphNodeContainer != nil {
		logger.Info("Stopping graph node")
		_ = testHarness.GraphNodeContainer.Terminate(context.Background())
	}

	if testHarness.ChainDockerNetwork != nil {
		logger.Info("Removing Docker network")
		_ = testHarness.ChainDockerNetwork.Remove(context.Background())
	}

	if testHarness.LocalstackContainer != nil {
		logger.Info("Stopping localstack container")
		if err := testHarness.LocalstackContainer.Terminate(ctx); err != nil {
			logger.Warn("Failed to terminate localstack container", "error", err)
		}
	}

	logger.Info("Teardown completed")
}
