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
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	routerbindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierRouter"
	verifierv1bindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierV1"
	paymentvaultbindings "github.com/Layr-Labs/eigenda/contracts/bindings/PaymentVault"
	"github.com/Layr-Labs/eigenda/core"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	coreeth "github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/eth/directory"
	"github.com/Layr-Labs/eigenda/core/payments"
	"github.com/Layr-Labs/eigenda/core/payments/clientledger"
	"github.com/Layr-Labs/eigenda/core/payments/ondemand"
	"github.com/Layr-Labs/eigenda/core/payments/reservation"
	"github.com/Layr-Labs/eigenda/core/payments/vault"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	verifierv2 "github.com/Layr-Labs/eigenda/encoding/kzg/verifier/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
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
	if err := setupPayloadDisperserForContext(ctx, testCtx, infra); err != nil {
		return nil, fmt.Errorf("failed to setup payload disperser: %w", err)
	}

	if err := setupPaymentVaultTransactor(ctx, testCtx, infra); err != nil {
		return nil, fmt.Errorf("setup payment vault transactor: %w", err)
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

func setupPayloadDisperserForContext(
	ctx context.Context,
	testHarness *TestHarness,
	infra *InfrastructureHarness,
) error {
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
	testHarness.TestAccountID = accountId

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

	// Create ClientLedger based on configured mode
	var clientLedger *clientledger.ClientLedger
	mode := infra.TestConfig.ClientLedgerMode
	if mode != clientledger.ClientLedgerModeLegacy {
		eigenDADirectoryAddr := gethcommon.HexToAddress(infra.TestConfig.EigenDA.EigenDADirectory)
		clientLedger, err = buildClientLedger(
			ctx,
			infra.Logger,
			testHarness.EthClient,
			eigenDADirectoryAddr,
			accountId,
			mode,
			disperserClient,
		)
		if err != nil {
			return fmt.Errorf("build client ledger: %w", err)
		}
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
		clientLedger,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create payload disperser: %w", err)
	}

	return nil
}

func buildClientLedger(
	ctx context.Context,
	logger logging.Logger,
	ethClient common.EthClient,
	eigenDADirectoryAddr gethcommon.Address,
	accountID gethcommon.Address,
	mode clientledger.ClientLedgerMode,
	disperserClient *clientsv2.DisperserClient,
) (*clientledger.ClientLedger, error) {
	contractDirectory, err := directory.NewContractDirectory(ctx, logger, ethClient, eigenDADirectoryAddr)
	if err != nil {
		return nil, fmt.Errorf("new contract directory: %w", err)
	}

	paymentVaultAddr, err := contractDirectory.GetContractAddress(ctx, directory.PaymentVault)
	if err != nil {
		return nil, fmt.Errorf("get PaymentVault address: %w", err)
	}

	paymentVault, err := vault.NewPaymentVault(logger, ethClient, paymentVaultAddr)
	if err != nil {
		return nil, fmt.Errorf("new payment vault: %w", err)
	}

	minNumSymbols, err := paymentVault.GetMinNumSymbols(ctx)
	if err != nil {
		return nil, fmt.Errorf("get min num symbols: %w", err)
	}

	var reservationLedger *reservation.ReservationLedger
	var onDemandLedger *ondemand.OnDemandLedger

	// Build reservation ledger if needed
	needsReservation := mode == clientledger.ClientLedgerModeReservationOnly ||
		mode == clientledger.ClientLedgerModeReservationAndOnDemand
	if needsReservation {
		reservationLedger, err = buildReservationLedger(ctx, paymentVault, accountID, minNumSymbols)
		if err != nil {
			return nil, fmt.Errorf("build reservation ledger: %w", err)
		}
	}

	// Build on-demand ledger if needed
	needsOnDemand := mode == clientledger.ClientLedgerModeOnDemandOnly ||
		mode == clientledger.ClientLedgerModeReservationAndOnDemand
	if needsOnDemand {
		onDemandLedger, err = buildOnDemandLedger(ctx, paymentVault, accountID, minNumSymbols, disperserClient)
		if err != nil {
			return nil, fmt.Errorf("build on-demand ledger: %w", err)
		}
	}

	ledger := clientledger.NewClientLedger(
		ctx,
		logger,
		metrics.NoopAccountantMetrics,
		accountID,
		mode,
		reservationLedger,
		onDemandLedger,
		time.Now,
		paymentVault,
		1*time.Second, // update interval for vault monitoring
	)

	return ledger, nil
}

func buildReservationLedger(
	ctx context.Context,
	paymentVault payments.PaymentVault,
	accountID gethcommon.Address,
	minNumSymbols uint32,
) (*reservation.ReservationLedger, error) {
	reservationData, err := paymentVault.GetReservation(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("get reservation: %w", err)
	}
	if reservationData == nil {
		return nil, fmt.Errorf("no reservation found for account %s", accountID.Hex())
	}

	clientReservation, err := reservation.NewReservation(
		reservationData.SymbolsPerSecond,
		time.Unix(int64(reservationData.StartTimestamp), 0),
		time.Unix(int64(reservationData.EndTimestamp), 0),
		reservationData.QuorumNumbers,
	)
	if err != nil {
		return nil, fmt.Errorf("new reservation: %w", err)
	}

	reservationConfig, err := reservation.NewReservationLedgerConfig(
		*clientReservation,
		minNumSymbols,
		true,
		reservation.OverfillOncePermitted,
		10*time.Second,
	)
	if err != nil {
		return nil, fmt.Errorf("new reservation ledger config: %w", err)
	}

	reservationLedger, err := reservation.NewReservationLedger(*reservationConfig, time.Now())
	if err != nil {
		return nil, fmt.Errorf("new reservation ledger: %w", err)
	}

	return reservationLedger, nil
}

func buildOnDemandLedger(
	ctx context.Context,
	paymentVault payments.PaymentVault,
	accountID gethcommon.Address,
	minNumSymbols uint32,
	disperserClient *clientsv2.DisperserClient,
) (*ondemand.OnDemandLedger, error) {
	pricePerSymbol, err := paymentVault.GetPricePerSymbol(ctx)
	if err != nil {
		return nil, fmt.Errorf("get price per symbol: %w", err)
	}

	totalDeposits, err := paymentVault.GetTotalDeposit(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("get total deposit from vault: %w", err)
	}

	paymentState, err := disperserClient.GetPaymentState(ctx)
	if err != nil {
		return nil, fmt.Errorf("get payment state from disperser: %w", err)
	}

	var cumulativePayment *big.Int
	if paymentState.GetCumulativePayment() == nil {
		cumulativePayment = big.NewInt(0)
	} else {
		cumulativePayment = new(big.Int).SetBytes(paymentState.GetCumulativePayment())
	}

	onDemandLedger, err := ondemand.OnDemandLedgerFromValue(
		totalDeposits,
		new(big.Int).SetUint64(pricePerSymbol),
		minNumSymbols,
		cumulativePayment,
	)
	if err != nil {
		return nil, fmt.Errorf("on-demand ledger from value: %w", err)
	}

	return onDemandLedger, nil
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

func setupPaymentVaultTransactor(
	ctx context.Context,
	testHarness *TestHarness,
	infra *InfrastructureHarness,
) error {
	eigenDADirectoryAddr := gethcommon.HexToAddress(infra.TestConfig.EigenDA.EigenDADirectory)
	contractDirectory, err := directory.NewContractDirectory(
		ctx, infra.Logger, testHarness.EthClient, eigenDADirectoryAddr)
	if err != nil {
		return fmt.Errorf("new contract directory: %w", err)
	}

	paymentVaultAddr, err := contractDirectory.GetContractAddress(ctx, directory.PaymentVault)
	if err != nil {
		return fmt.Errorf("get PaymentVault address: %w", err)
	}

	transactor, err := paymentvaultbindings.NewContractPaymentVaultTransactor(paymentVaultAddr, testHarness.EthClient)
	if err != nil {
		return fmt.Errorf("new PaymentVault transactor: %w", err)
	}

	testHarness.PaymentVaultTransactor = transactor

	return nil
}
