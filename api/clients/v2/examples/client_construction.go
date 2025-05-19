package examples

import (
	"fmt"
	"io"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloaddispersal"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloadretrieval"
	"github.com/Layr-Labs/eigenda/api/clients/v2/relay"
	"github.com/Layr-Labs/eigenda/api/clients/v2/validator"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// These constants are specific to the EigenDA holesky testnet. To execute the provided examples on a different
// network, you will need to set these constants to the correct values, based on the chosen network.
const (
	ethRPCURL                        = "https://ethereum-holesky-rpc.publicnode.com"
	disperserHostname                = "disperser-testnet-holesky.eigenda.xyz"
	certVerifierRouterAddress        = "0x7F40A8e1B62aa1c8Afed23f6E8bAe0D340A4BC4e"
	registryCoordinatorAddress       = "0x53012C69A189cfA2D9d29eb6F19B32e0A2EA3490"
	blsOperatorStateRetrieverAddress = "0x003497Dd77E5B73C40e8aCbB562C8bb0410320E7"
	eigenDAServiceManagerAddress     = "0xD4A7E1Bd8015057293f0D0A557088c286942e84b"
)

func createPayloadDisperser(privateKey string) (*payloaddispersal.PayloadDisperser, error) {
	logger, err := createLogger()
	if err != nil {
		panic(fmt.Sprintf("create logger: %v", err))
	}

	kzgProver, err := createKzgProver()
	if err != nil {
		return nil, fmt.Errorf("create kzg prover: %v", err)
	}

	disperserClient, err := createDisperserClient(privateKey, kzgProver)
	if err != nil {
		return nil, fmt.Errorf("create disperser client: %w", err)
	}

	certVerifier, err := createGenericCertVerifier()
	if err != nil {
		return nil, fmt.Errorf("create cert verifier: %w", err)
	}

	ethClient, err := createEthClient(logger)
	if err != nil {
		return nil, fmt.Errorf("create eth client: %w", err)
	}

	certBuilder, err := createCertBuilder()
	if err != nil {
		return nil, fmt.Errorf("create cert builder: %w", err)
	}

	blockNumMonitor, err := createBlockNumberMonitor()
	if err != nil {
		return nil, fmt.Errorf("create block number monitor: %w", err)
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
		ethClient,
		disperserClient,
		blockNumMonitor,
		certBuilder,
		certVerifier,
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

	reader, err := createEthReader(logger, ethClient)
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
		kzgVerifier.Srs.G1)
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

	// Create the eth reader
	ethReader, err := createEthReader(logger, ethClient)
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
		kzgVerifier.Srs.G1)
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

func createDisperserClient(privateKey string, kzgProver *prover.Prover) (clients.DisperserClient, error) {
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
		disperserClientConfig,
		signer,
		kzgProver,
		nil)
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

func createGenericCertVerifier() (*verification.CertVerifier, error) {
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

	return verification.NewGenericCertVerifier(
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

	return clients.NewCertBuilder(
		logger,
		gethcommon.HexToAddress(blsOperatorStateRetrieverAddress),
		gethcommon.HexToAddress(registryCoordinatorAddress),
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
		1 * time.Second,
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

func createEthReader(logger logging.Logger, ethClient common.EthClient) (*eth.Reader, error) {
	ethReader, err := eth.NewReader(
		logger,
		ethClient,
		blsOperatorStateRetrieverAddress,
		eigenDAServiceManagerAddress,
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
