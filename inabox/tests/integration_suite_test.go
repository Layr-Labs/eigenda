package integration_test

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/big"
	"math/rand"
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
	"github.com/Layr-Labs/eigenda/testbed"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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

	testConfig          *deploy.Config
	localstackContainer *testbed.LocalStackContainer
	localStackPort      string

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

	var loggerConfig *common.LoggerConfig
	if os.Getenv("CI") != "" {
		loggerConfig = common.DefaultLoggerConfig()
	} else {
		loggerConfig = common.DefaultConsoleLoggerConfig()
	}
	logger, err = common.NewLogger(loggerConfig)
	Expect(err).To(BeNil())

	if testConfig.Environment.IsLocal() {
		if !inMemoryBlobStore {
			logger.Info("Using shared Blob Store")
			localStackPort = "4570"

			cfg := testbed.DefaultLocalStackConfig()
			cfg.Services = []string{"s3", "dynamodb", "kms"}
			cfg.Port = localStackPort
			cfg.Host = "0.0.0.0"

			localstackContainer, err = testbed.NewLocalStackContainer(context.Background(), cfg)
			Expect(err).To(BeNil())

			deployConfig := testbed.DeployResourcesConfig{
				LocalStackEndpoint:  fmt.Sprintf("http://%s:%s", cfg.Host, cfg.Port),
				MetadataTableName:   metadataTableName,
				BucketTableName:     bucketTableName,
				V2MetadataTableName: metadataTableNameV2,
				Logger:              logger,
			}
			err = testbed.DeployResources(context.Background(), deployConfig)
			Expect(err).To(BeNil())
		} else {
			logger.Info("Using in-memory Blob Store")
		}

		logger.Info("Starting anvil")
		testConfig.StartAnvil()

		deployer, ok := testConfig.GetDeployer(testConfig.EigenDA.Deployer)
		if ok && deployer.DeploySubgraphs {
			logger.Info("Starting graph node")
			testConfig.StartGraphNode()
		}

		logger.Info("Deploying experiment")
		testConfig.DeployExperiment()
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

		logger.Info("Registering blob versions and relays")
		testConfig.RegisterBlobVersionAndRelays(ethClient)

		logger.Info("Registering disperser keypair")
		err = testConfig.RegisterDisperserKeypair(ethClient, logger)
		if err != nil {
			panic(err)
		}

		logger.Info("Starting binaries")
		testConfig.StartBinaries()

		eigenDACertVerifierV1, err = verifierv1bindings.NewContractEigenDACertVerifierV1(gethcommon.HexToAddress(testConfig.EigenDAV1CertVerifier), ethClient)
		Expect(err).To(BeNil())
		err = setupRetrievalClients(testConfig)
		Expect(err).To(BeNil())

		logger.Info("Building client verification and interaction components")

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
	disperserClient, err := clientsv2.NewDisperserClient(
		logger,
		disperserClientConfig,
		signer,
		nil,
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
		kzgVerifier.Srs.G1,
		metrics.NoopRetrievalMetrics)

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
		kzgVerifier.Srs.G1,
		metrics.NoopRetrievalMetrics)

	return err
}

var _ = AfterSuite(func() {
	if testConfig.Environment.IsLocal() {
		if cancel != nil {
			cancel()
		}

		logger.Info("Stopping binaries")
		testConfig.StopBinaries()

		logger.Info("Stopping anvil")
		testConfig.StopAnvil()

		logger.Info("Stopping graph node")
		testConfig.StopGraphNode()

		if localstackContainer != nil {
			logger.Info("Stopping localstack container")
			_ = localstackContainer.Terminate(context.Background())
		}
	}
})
