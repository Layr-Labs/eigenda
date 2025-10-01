package integration

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients"
	clientsv2 "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/metrics"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloaddispersal"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloadretrieval"
	"github.com/Layr-Labs/eigenda/api/clients/v2/relay"
	validatorclientsv2 "github.com/Layr-Labs/eigenda/api/clients/v2/validator"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	"github.com/Layr-Labs/eigenda/common/geth"
	routerbindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierRouter"
	verifierv1bindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierV1"
	"github.com/Layr-Labs/eigenda/core"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	coreeth "github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	verifierv2 "github.com/Layr-Labs/eigenda/encoding/kzg/verifier/v2"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

// NewTestHarnessWithSetup creates a fully initialized TestHarness with all components set up.
// This provides a fresh set of clients and verifiers for each test.
func NewTestHarnessWithSetup(infra *InfrastructureHarness) (*TestHarness, error) {
	ctx := context.Background()
	testCtx := &TestHarness{
		NumConfirmations: 1,
		NumRetries:       5,
	}

	// Get deployer's private key
	deployer, ok := infra.TestConfig.GetDeployer(infra.TestConfig.EigenDA.Deployer)
	if !ok {
		return nil, fmt.Errorf("failed to get deployer")
	}

	pk := infra.TestConfig.Pks.EcdsaMap[deployer.Name].PrivateKey
	pk = strings.TrimPrefix(pk, "0x")
	pk = strings.TrimPrefix(pk, "0X")

	// Create Ethereum clients
	var err error
	testCtx.EthClient, err = geth.NewMultiHomingClient(geth.EthClientConfig{
		RPCURLs:          []string{infra.TestConfig.Deployers[0].RPC},
		PrivateKeyString: pk,
		NumConfirmations: testCtx.NumConfirmations,
		NumRetries:       testCtx.NumRetries,
	}, gethcommon.Address{}, infra.Logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create eth client: %w", err)
	}

	testCtx.RPCClient, err = ethrpc.Dial(infra.TestConfig.Deployers[0].RPC)
	if err != nil {
		return nil, fmt.Errorf("failed to create rpc client: %w", err)
	}

	// Force foundry to mine a block since it isn't auto-mining
	err = testCtx.RPCClient.CallContext(ctx, nil, "evm_mine")
	if err != nil {
		return nil, fmt.Errorf("failed to mine block: %w", err)
	}

	// Get chain ID
	testCtx.ChainID, err = testCtx.EthClient.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Create transactor options
	testCtx.DeployerTransactorOpts = newTransactOptsFromPrivateKey(pk, testCtx.ChainID)

	// Create contract bindings
	testCtx.EigenDACertVerifierV1, err = verifierv1bindings.NewContractEigenDACertVerifierV1(
		gethcommon.HexToAddress(infra.TestConfig.EigenDAV1CertVerifier),
		testCtx.EthClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create EigenDA cert verifier V1: %w", err)
	}

	testCtx.EigenDACertVerifierRouter, err = routerbindings.NewContractEigenDACertVerifierRouterTransactor(
		gethcommon.HexToAddress(infra.TestConfig.EigenDA.CertVerifierRouter),
		testCtx.EthClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create router transactor: %w", err)
	}

	testCtx.EigenDACertVerifierRouterCaller, err = routerbindings.NewContractEigenDACertVerifierRouterCaller(
		gethcommon.HexToAddress(infra.TestConfig.EigenDA.CertVerifierRouter),
		testCtx.EthClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create router caller: %w", err)
	}

	// Setup verifiers and cert builder
	if err := setupVerifiersForContext(testCtx, infra); err != nil {
		return nil, fmt.Errorf("failed to setup verifiers: %w", err)
	}

	// Setup retrieval clients
	if err := setupRetrievalClientsForContext(testCtx, infra); err != nil {
		return nil, fmt.Errorf("failed to setup retrieval clients: %w", err)
	}

	// Setup payload disperser
	if err := setupPayloadDisperserForContext(testCtx, infra); err != nil {
		return nil, fmt.Errorf("failed to setup payload disperser: %w", err)
	}

	return testCtx, nil
}

func setupVerifiersForContext(testCtx *TestHarness, infra *InfrastructureHarness) error {
	var err error
	testCtx.CertBuilder, err = clientsv2.NewCertBuilder(
		infra.Logger,
		gethcommon.HexToAddress(infra.TestConfig.EigenDA.OperatorStateRetriever),
		gethcommon.HexToAddress(infra.TestConfig.EigenDA.RegistryCoordinator),
		testCtx.EthClient,
	)
	if err != nil {
		return fmt.Errorf("failed to create cert builder: %w", err)
	}

	routerAddressProvider, err := verification.BuildRouterAddressProvider(
		gethcommon.HexToAddress(infra.TestConfig.EigenDA.CertVerifierRouter),
		testCtx.EthClient,
		infra.Logger)
	if err != nil {
		return fmt.Errorf("failed to build router address provider: %w", err)
	}

	staticAddressProvider := verification.NewStaticCertVerifierAddressProvider(
		gethcommon.HexToAddress(infra.TestConfig.EigenDA.CertVerifier))

	testCtx.StaticCertVerifier, err = verification.NewCertVerifier(
		infra.Logger,
		testCtx.EthClient,
		staticAddressProvider)
	if err != nil {
		return fmt.Errorf("failed to create static cert verifier: %w", err)
	}

	testCtx.RouterCertVerifier, err = verification.NewCertVerifier(
		infra.Logger,
		testCtx.EthClient,
		routerAddressProvider)
	if err != nil {
		return fmt.Errorf("failed to create router cert verifier: %w", err)
	}

	return nil
}

func setupRetrievalClientsForContext(testHarness *TestHarness, infraHarness *InfrastructureHarness) error {
	tx, err := coreeth.NewWriter(
		infraHarness.Logger,
		testHarness.EthClient,
		infraHarness.TestConfig.EigenDA.OperatorStateRetriever,
		infraHarness.TestConfig.EigenDA.ServiceManager)
	if err != nil {
		return fmt.Errorf("failed to create writer: %w", err)
	}

	cs := coreeth.NewChainState(tx, testHarness.EthClient)
	agn := &core.StdAssignmentCoordinator{}
	nodeClient := clients.NewNodeClient(20 * time.Second)

	srsOrder, err := strconv.Atoi(infraHarness.TestConfig.Retriever.RETRIEVER_SRS_ORDER)
	if err != nil {
		return fmt.Errorf("failed to parse SRS order: %w", err)
	}

	kzgConfig := &kzg.KzgConfig{
		G1Path:          infraHarness.TestConfig.Retriever.RETRIEVER_G1_PATH,
		G2Path:          infraHarness.TestConfig.Retriever.RETRIEVER_G2_PATH,
		CacheDir:        infraHarness.TestConfig.Retriever.RETRIEVER_CACHE_PATH,
		SRSOrder:        uint64(srsOrder),
		SRSNumberToLoad: uint64(srsOrder),
		NumWorker:       1,
		PreloadEncoder:  false,
		LoadG2Points:    true,
	}

	kzgVerifier, err := verifier.NewVerifier(kzgConfig, nil)
	if err != nil {
		return fmt.Errorf("failed to create kzg verifier: %w", err)
	}

	testHarness.RetrievalClient, err = clients.NewRetrievalClient(
		infraHarness.Logger, cs, agn, nodeClient, kzgVerifier, 10)
	if err != nil {
		return fmt.Errorf("failed to create retrieval client: %w", err)
	}

	testHarness.ChainReader, err = coreeth.NewReader(
		infraHarness.Logger,
		testHarness.EthClient,
		infraHarness.TestConfig.EigenDA.OperatorStateRetriever,
		infraHarness.TestConfig.EigenDA.ServiceManager,
	)
	if err != nil {
		return fmt.Errorf("failed to create chain reader: %w", err)
	}

	// Setup V2 retrieval clients
	kzgVerifierV2, err := verifierv2.NewVerifier(verifierv2.KzgConfigFromV1Config(kzgConfig), nil)
	if err != nil {
		return fmt.Errorf("new verifier v2: %w", err)
	}

	clientConfig := validatorclientsv2.DefaultClientConfig()
	retrievalClientV2 := validatorclientsv2.NewValidatorClient(
		infraHarness.Logger, testHarness.ChainReader, cs, kzgVerifierV2, clientConfig, nil)

	validatorPayloadRetrieverConfig := payloadretrieval.ValidatorPayloadRetrieverConfig{
		PayloadClientConfig: *clientsv2.GetDefaultPayloadClientConfig(),
		RetrievalTimeout:    1 * time.Minute,
	}

	testHarness.ValidatorRetrievalClientV2, err = payloadretrieval.NewValidatorPayloadRetriever(
		infraHarness.Logger,
		validatorPayloadRetrieverConfig,
		retrievalClientV2,
		kzgVerifier.G1SRS,
		metrics.NoopRetrievalMetrics)
	if err != nil {
		return fmt.Errorf("failed to create validator payload retriever: %w", err)
	}

	// Setup relay client
	relayClientConfig := &relay.RelayClientConfig{
		MaxGRPCMessageSize: 100 * 1024 * 1024, // 100 MB message size limit
	}

	relayUrlProvider, err := relay.NewRelayUrlProvider(
		testHarness.EthClient, testHarness.ChainReader.GetRelayRegistryAddress())
	if err != nil {
		return fmt.Errorf("failed to create relay URL provider: %w", err)
	}

	relayClient, err := relay.NewRelayClient(relayClientConfig, infraHarness.Logger, relayUrlProvider)
	if err != nil {
		return fmt.Errorf("failed to create relay client: %w", err)
	}

	relayPayloadRetrieverConfig := payloadretrieval.RelayPayloadRetrieverConfig{
		PayloadClientConfig: *clientsv2.GetDefaultPayloadClientConfig(),
		RelayTimeout:        5 * time.Second,
	}

	testHarness.RelayRetrievalClientV2, err = payloadretrieval.NewRelayPayloadRetriever(
		infraHarness.Logger,
		rand.New(rand.NewSource(time.Now().UnixNano())),
		relayPayloadRetrieverConfig,
		relayClient,
		kzgVerifier.G1SRS,
		metrics.NoopRetrievalMetrics)
	if err != nil {
		return fmt.Errorf("failed to create relay payload retriever: %w", err)
	}

	return nil
}

func setupPayloadDisperserForContext(testHarness *TestHarness, infra *InfrastructureHarness) error {
	// Set up the block monitor
	blockMonitor, err := verification.NewBlockNumberMonitor(infra.Logger, testHarness.EthClient, time.Second*1)
	if err != nil {
		return fmt.Errorf("failed to create block number monitor: %w", err)
	}

	// Set up the PayloadDisperser
	privateKeyHex := "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcded"
	signer, err := auth.NewLocalBlobRequestSigner(privateKeyHex)
	if err != nil {
		return fmt.Errorf("failed to create blob request signer: %w", err)
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
		infra.Logger,
		disperserClientConfig,
		signer,
		nil, // no prover so will query disperser for generating commitments
		accountant,
		metrics.NoopDispersalMetrics,
	)
	if err != nil {
		return fmt.Errorf("failed to create disperser client: %w", err)
	}

	payloadDisperserConfig := payloaddispersal.PayloadDisperserConfig{
		PayloadClientConfig:    *clientsv2.GetDefaultPayloadClientConfig(),
		DisperseBlobTimeout:    2 * time.Minute,
		BlobCompleteTimeout:    2 * time.Minute,
		BlobStatusPollInterval: 1 * time.Second,
		ContractCallTimeout:    5 * time.Second,
	}

	testHarness.PayloadDisperser, err = payloaddispersal.NewPayloadDisperser(
		infra.Logger,
		payloadDisperserConfig,
		disperserClient,
		blockMonitor,
		testHarness.CertBuilder,
		testHarness.RouterCertVerifier,
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create payload disperser: %w", err)
	}

	return nil
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
