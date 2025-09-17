package examples

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/metrics"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloaddispersal"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloadretrieval"
	"github.com/Layr-Labs/eigenda/api/clients/v2/relay"
	"github.com/Layr-Labs/eigenda/api/clients/v2/validator"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/eth/directory"
	"github.com/Layr-Labs/eigenda/core/payments"
	"github.com/Layr-Labs/eigenda/core/payments/clientledger"
	"github.com/Layr-Labs/eigenda/core/payments/ondemand"
	"github.com/Layr-Labs/eigenda/core/payments/reservation"
	"github.com/Layr-Labs/eigenda/core/payments/vault"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// These constants are specific to the EigenDA Sepolia testnet. To execute the provided examples on a different
// network, you will need to set these constants to the correct values, based on the chosen network.
const (
	ethRPCURL         = "https://ethereum-sepolia-rpc.publicnode.com"
	disperserHostname = "disperser-testnet-sepolia.eigenda.xyz"
	// Contract registry address for Sepolia - this allows fetching all other contract addresses
	contractRegistryAddress = "0x9620dC4B3564198554e4D2b06dEFB7A369D90257"
	// CertVerifierRouter is not available from the contract registry and must be specified directly
	certVerifierRouterAddress = "0x58D2B844a894f00b7E6F9F492b9F43aD54Cd4429"
)

func createPayloadDisperser(privateKeyHex string) (*payloaddispersal.PayloadDisperser, error) {
	logger, err := createLogger()
	if err != nil {
		panic(fmt.Sprintf("create logger: %v", err))
	}

	kzgProver, err := createKzgProver()
	if err != nil {
		return nil, fmt.Errorf("create kzg prover: %v", err)
	}

	disperserClient, err := createDisperserClient(logger, privateKeyHex, kzgProver)
	if err != nil {
		return nil, fmt.Errorf("create disperser client: %w", err)
	}

	certVerifier, err := createCertVerifier()
	if err != nil {
		return nil, fmt.Errorf("create cert verifier: %w", err)
	}

	certBuilder, err := createCertBuilder()
	if err != nil {
		return nil, fmt.Errorf("create cert builder: %w", err)
	}

	blockNumMonitor, err := createBlockNumberMonitor()
	if err != nil {
		return nil, fmt.Errorf("create block number monitor: %w", err)
	}

	privateKeyBytes := gethcommon.FromHex(privateKeyHex)
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("to ecdsa: %w", err)
	}
	accountID := crypto.PubkeyToAddress(privateKey.PublicKey)

	ethClient, err := createEthClient(logger)
	if err != nil {
		return nil, fmt.Errorf("create eth client: %w", err)
	}

	contractDirectory, err := createContractDirectory(context.Background(), logger, ethClient)
	if err != nil {
		return nil, fmt.Errorf("create contract directory: %w", err)
	}

	clientLedger, err := createClientLedger(
		context.Background(),
		logger,
		clientledger.ClientLedgerModeReservationAndOnDemand,
		ethClient,
		accountID,
		contractDirectory,
		disperserClient,
	)
	if err != nil {
		return nil, fmt.Errorf("create client ledger: %w", err)
	}

	payloadDisperserConfig := payloaddispersal.PayloadDisperserConfig{
		PayloadClientConfig:    *clients.GetDefaultPayloadClientConfig(),
		DisperseBlobTimeout:    30 * time.Second,
		BlobCompleteTimeout:    30 * time.Second,
		BlobStatusPollInterval: 1 * time.Second,
		ContractCallTimeout:    5 * time.Second,
	}

	return payloaddispersal.NewPayloadDisperser(
		logger,
		payloadDisperserConfig,
		disperserClient,
		blockNumMonitor,
		certBuilder,
		certVerifier,
		clientLedger,
		nil,
	)
}

func createRelayPayloadRetriever() (*payloadretrieval.RelayPayloadRetriever, error) {
	logger, err := createLogger()
	if err != nil {
		return nil, fmt.Errorf("create logger: %w", err)
	}

	ethClient, err := createEthClient(logger)
	if err != nil {
		return nil, fmt.Errorf("create eth client: %w", err)
	}

	contractDirectory, err := createContractDirectory(context.Background(), logger, ethClient)
	if err != nil {
		return nil, fmt.Errorf("create contract directory: %w", err)
	}

	reader, err := createEthReader(logger, ethClient, contractDirectory)
	if err != nil {
		return nil, fmt.Errorf("create eth reader: %w", err)
	}

	relayClient, err := createRelayClient(logger, ethClient, reader.GetRelayRegistryAddress())
	if err != nil {
		return nil, fmt.Errorf("create relay client: %w", err)
	}

	kzgVerifier, err := createKzgVerifier()
	if err != nil {
		return nil, fmt.Errorf("create kzg verifier: %w", err)
	}

	relayPayloadRetrieverConfig := payloadretrieval.RelayPayloadRetrieverConfig{
		PayloadClientConfig: *clients.GetDefaultPayloadClientConfig(),
		RelayTimeout:        5 * time.Second,
	}

	return payloadretrieval.NewRelayPayloadRetriever(
		logger,
		rand.New(rand.NewSource(time.Now().UnixNano())),
		relayPayloadRetrieverConfig,
		relayClient,
		kzgVerifier.G1SRS,
		metrics.NoopRetrievalMetrics)
}

func createValidatorPayloadRetriever() (*payloadretrieval.ValidatorPayloadRetriever, error) {
	logger, err := createLogger()
	if err != nil {
		return nil, fmt.Errorf("create logger: %w", err)
	}

	// Create an EthClient for blockchain interaction
	ethClient, err := createEthClient(logger)
	if err != nil {
		return nil, fmt.Errorf("create eth client: %w", err)
	}

	contractDirectory, err := createContractDirectory(context.Background(), logger, ethClient)
	if err != nil {
		return nil, fmt.Errorf("create contract directory: %w", err)
	}

	// Create the eth reader
	ethReader, err := createEthReader(logger, ethClient, contractDirectory)
	if err != nil {
		return nil, fmt.Errorf("create eth reader: %w", err)
	}

	chainState := eth.NewChainState(ethReader, ethClient)
	kzgVerifier, err := createKzgVerifier()
	if err != nil {
		return nil, fmt.Errorf("create kzg verifier: %w", err)
	}

	clientConfig := validator.DefaultClientConfig()

	// Create the retrieval client for fetching blobs from DA nodes
	retrievalClient := validator.NewValidatorClient(
		logger,
		ethReader,
		chainState,
		kzgVerifier,
		clientConfig,
		nil,
	)

	// Create the ValidatorPayloadRetriever config
	validatorPayloadRetrieverConfig := payloadretrieval.ValidatorPayloadRetrieverConfig{
		PayloadClientConfig: *clients.GetDefaultPayloadClientConfig(),
		RetrievalTimeout:    1 * time.Minute,
	}

	return payloadretrieval.NewValidatorPayloadRetriever(
		logger,
		validatorPayloadRetrieverConfig,
		retrievalClient,
		kzgVerifier.G1SRS,
		metrics.NoopRetrievalMetrics)
}

func createRelayClient(
	logger logging.Logger,
	ethClient common.EthClient,
	relayRegistryAddress gethcommon.Address,
) (relay.RelayClient, error) {
	config := &relay.RelayClientConfig{
		UseSecureGrpcFlag:  true,
		MaxGRPCMessageSize: 100 * 1024 * 1024, // 100 MB message size limit
	}

	relayUrlProvider, err := relay.NewRelayUrlProvider(ethClient, relayRegistryAddress)
	if err != nil {
		return nil, fmt.Errorf("create relay url provider: %w", err)
	}

	return relay.NewRelayClient(
		config,
		logger,
		relayUrlProvider)
}

func createDisperserClient(
	logger logging.Logger,
	privateKey string,
	kzgProver *prover.Prover,
) (*clients.DisperserClient, error) {
	signer, err := auth.NewLocalBlobRequestSigner(privateKey)
	if err != nil {
		return nil, fmt.Errorf("create blob request signer: %w", err)
	}

	disperserClientConfig := &clients.DisperserClientConfig{
		Hostname:          disperserHostname,
		Port:              "443",
		UseSecureGrpcFlag: true,
	}

	return clients.NewDisperserClient(
		logger,
		disperserClientConfig,
		signer,
		kzgProver,
		nil,
		metrics.NoopDispersalMetrics)
}

func createKzgVerifier() (*verifier.Verifier, error) {
	kzgConfig := createKzgConfig()
	kzgConfig.LoadG2Points = false
	blobVerifier, err := verifier.NewVerifier(&kzgConfig, nil)
	if err != nil {
		return nil, fmt.Errorf("create blob verifier: %w", err)
	}

	return blobVerifier, nil
}

func createKzgProver() (*prover.Prover, error) {
	kzgConfig := createKzgConfig()
	kzgProver, err := prover.NewProver(&kzgConfig, nil)
	if err != nil {
		return nil, err
	}

	return kzgProver, nil
}

func createCertVerifier() (*verification.CertVerifier, error) {
	logger, err := createLogger()
	if err != nil {
		return nil, fmt.Errorf("create logger: %v", err)
	}

	ethClient, err := createEthClient(logger)
	if err != nil {
		return nil, fmt.Errorf("create eth client: %w", err)
	}

	routerAddressProvider, err := verification.BuildRouterAddressProvider(
		gethcommon.HexToAddress(certVerifierRouterAddress),
		ethClient,
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("create router address provider: %w", err)
	}

	return verification.NewCertVerifier(
		logger,
		ethClient,
		routerAddressProvider,
	)
}

func createCertBuilder() (*clients.CertBuilder, error) {
	logger, err := createLogger()
	if err != nil {
		return nil, fmt.Errorf("create logger: %v", err)
	}

	ethClient, err := createEthClient(logger)
	if err != nil {
		return nil, fmt.Errorf("create eth client: %w", err)
	}

	contractDirectory, err := createContractDirectory(context.Background(), logger, ethClient)
	if err != nil {
		return nil, fmt.Errorf("create contract directory: %w", err)
	}

	operatorStateRetrieverAddr, err := contractDirectory.GetContractAddress(
		context.Background(), directory.OperatorStateRetriever)
	if err != nil {
		return nil, fmt.Errorf("get OperatorStateRetriever address: %w", err)
	}

	registryCoordinatorAddr, err := contractDirectory.GetContractAddress(
		context.Background(), directory.RegistryCoordinator)
	if err != nil {
		return nil, fmt.Errorf("get RegistryCoordinator address: %w", err)
	}

	return clients.NewCertBuilder(
		logger,
		operatorStateRetrieverAddr,
		registryCoordinatorAddr,
		ethClient,
	)
}

func createBlockNumberMonitor() (*verification.BlockNumberMonitor, error) {
	logger, err := createLogger()
	if err != nil {
		return nil, fmt.Errorf("create logger: %v", err)
	}

	ethClient, err := createEthClient(logger)
	if err != nil {
		return nil, fmt.Errorf("create eth client: %w", err)
	}

	return verification.NewBlockNumberMonitor(
		logger,
		ethClient,
		1*time.Second,
	)
}

func createEthClient(logger logging.Logger) (*geth.EthClient, error) {
	ethClientConfig := geth.EthClientConfig{
		RPCURLs:          []string{ethRPCURL},
		NumConfirmations: 0,
		NumRetries:       3,
	}

	return geth.NewClient(
		ethClientConfig,
		gethcommon.Address{},
		0,
		logger)
}

func createKzgConfig() kzg.KzgConfig {
	srsPath := "../../../../resources/srs"
	return kzg.KzgConfig{
		LoadG2Points:    true,
		G1Path:          filepath.Join(srsPath, "g1.point"),
		G2Path:          filepath.Join(srsPath, "g2.point"),
		G2TrailingPath:  filepath.Join(srsPath, "g2.trailing.point"),
		CacheDir:        filepath.Join(srsPath, "SRSTables"),
		SRSOrder:        268435456, // must always be this constant, which was used during eigenDA SRS generation
		SRSNumberToLoad: uint64(1<<13) / encoding.BYTES_PER_SYMBOL,
		NumWorker:       4,
	}
}

func createEthReader(
	logger logging.Logger,
	ethClient common.EthClient,
	contractDirectory *directory.ContractDirectory,
) (*eth.Reader, error) {
	operatorStateRetrieverAddr, err := contractDirectory.GetContractAddress(
		context.Background(), directory.OperatorStateRetriever)
	if err != nil {
		return nil, fmt.Errorf("get OperatorStateRetriever address: %w", err)
	}

	serviceManagerAddr, err := contractDirectory.GetContractAddress(context.Background(), directory.ServiceManager)
	if err != nil {
		return nil, fmt.Errorf("get ServiceManager address: %w", err)
	}

	ethReader, err := eth.NewReader(
		logger,
		ethClient,
		operatorStateRetrieverAddr.Hex(),
		serviceManagerAddr.Hex(),
	)
	if err != nil {
		return nil, fmt.Errorf("new reader: %w", err)
	}

	return ethReader, nil
}

func createLogger() (logging.Logger, error) {
	config := common.DefaultLoggerConfig()
	config.OutputWriter = io.Discard // Send logs to /dev/null
	logger, err := common.NewLogger(config)
	if err != nil {
		return nil, fmt.Errorf("create new logger: %w", err)
	}

	return logger, nil
}

func createContractDirectory(
	ctx context.Context,
	logger logging.Logger,
	ethClient common.EthClient,
) (*directory.ContractDirectory, error) {
	directoryAddress := gethcommon.HexToAddress(contractRegistryAddress)
	contractDirectory, err := directory.NewContractDirectory(ctx, logger, ethClient, directoryAddress)
	if err != nil {
		return nil, fmt.Errorf("new contract directory: %w", err)
	}
	return contractDirectory, nil
}

func createClientLedger(
	ctx context.Context,
	logger logging.Logger,
	mode clientledger.ClientLedgerMode,
	ethClient common.EthClient,
	accountID gethcommon.Address,
	contractDirectory *directory.ContractDirectory,
	disperserClient *clients.DisperserClient,
) (*clientledger.ClientLedger, error) {
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

	now := time.Now()

	reservationLedger, err = createReservationLedger(ctx, paymentVault, accountID, now, minNumSymbols)
	if err != nil {
		return nil, fmt.Errorf("create reservation ledger: %w", err)
	}
	onDemandLedger, err = createOnDemandLedger(ctx, paymentVault, accountID, minNumSymbols, disperserClient)
	if err != nil {
		return nil, fmt.Errorf("create on-demand ledger: %w", err)
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
		5*time.Minute,
	)

	return ledger, nil
}

func createReservationLedger(
	ctx context.Context,
	paymentVault payments.PaymentVault,
	accountID gethcommon.Address,
	now time.Time,
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
		time.Minute,
	)
	if err != nil {
		return nil, fmt.Errorf("new reservation ledger config: %w", err)
	}

	reservationLedger, err := reservation.NewReservationLedger(*reservationConfig, now)
	if err != nil {
		return nil, fmt.Errorf("new reservation ledger: %w", err)
	}

	return reservationLedger, nil
}

func createOnDemandLedger(
	ctx context.Context,
	paymentVault payments.PaymentVault,
	accountID gethcommon.Address,
	minNumSymbols uint32,
	disperserClient *clients.DisperserClient,
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

	if paymentState.GetCumulativePayment() == nil {
		return nil, errors.New("received nil cumulative payment from disperser")
	}

	onDemandLedger, err := ondemand.OnDemandLedgerFromValue(
		totalDeposits,
		new(big.Int).SetUint64(pricePerSymbol),
		minNumSymbols,
		new(big.Int).SetBytes(paymentState.GetCumulativePayment()),
	)
	if err != nil {
		return nil, fmt.Errorf("new on-demand ledger: %w", err)
	}

	return onDemandLedger, nil
}
