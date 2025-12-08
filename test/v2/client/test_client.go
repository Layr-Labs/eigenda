package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"strings"
	"time"

	clientsv2 "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/api/clients/v2/dispersal"
	metricsv2 "github.com/Layr-Labs/eigenda/api/clients/v2/metrics"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloadretrieval"
	"github.com/Layr-Labs/eigenda/api/clients/v2/relay"
	"github.com/Layr-Labs/eigenda/api/clients/v2/validator"
	"github.com/Layr-Labs/eigenda/api/clients/v2/validator/mock"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	"github.com/Layr-Labs/eigenda/api/hashing"
	proxycommon "github.com/Layr-Labs/eigenda/api/proxy/common"
	proxyconfig "github.com/Layr-Labs/eigenda/api/proxy/config"
	"github.com/Layr-Labs/eigenda/api/proxy/config/enablement"
	proxyserver "github.com/Layr-Labs/eigenda/api/proxy/servers/rest"
	"github.com/Layr-Labs/eigenda/api/proxy/store"
	"github.com/Layr-Labs/eigenda/api/proxy/store/builder"
	common_eigenda "github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/disperser"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/ratelimit"
	"github.com/Layr-Labs/eigenda/core"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/eth/directory"
	"github.com/Layr-Labs/eigenda/core/payments"
	"github.com/Layr-Labs/eigenda/core/payments/clientledger"
	"github.com/Layr-Labs/eigenda/core/payments/ondemand"
	"github.com/Layr-Labs/eigenda/core/payments/reservation"
	"github.com/Layr-Labs/eigenda/core/payments/vault"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	kzgv1 "github.com/Layr-Labs/eigenda/encoding/v1/kzg"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg/committer"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg/verifier"
	"github.com/Layr-Labs/eigenda/encoding/v2/rs"
	"github.com/Layr-Labs/eigenda/test/random"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/docker/go-units"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	SRSPathG1         = "g1.point"
	SRSPathG2         = "g2.point"
	SRSPathG2Trailing = "g2.trailing.point"
	SRSPathSRSTables  = "SRSTables"
)

// TestClient encapsulates the various clients necessary for interacting with EigenDA.
type TestClient struct {
	config                      *TestClientConfig
	payloadClientConfig         *clientsv2.PayloadClientConfig
	logger                      logging.Logger
	certVerifierAddressProvider clientsv2.CertVerifierAddressProvider
	disperserClientMultiplexer  *dispersal.DisperserClientMultiplexer
	payloadDisperser            *dispersal.PayloadDisperser
	relayClient                 relay.RelayClient
	relayPayloadRetriever       *payloadretrieval.RelayPayloadRetriever
	indexedChainState           core.IndexedChainState
	validatorClient             validator.ValidatorClient
	validatorPayloadRetriever   *payloadretrieval.ValidatorPayloadRetriever
	proxyWrapper                *ProxyWrapper
	// For fetching blobs from the validators without verifying or decoding them. Useful for load testing
	// validator downloads with limited CPU resources.
	onlyDownloadValidatorClient validator.ValidatorClient
	certBuilder                 *clientsv2.CertBuilder
	certVerifier                *verification.CertVerifier
	privateKey                  string
	metricsRegistry             *prometheus.Registry
	metrics                     *TestClientMetrics
}

// NewTestClient creates a new TestClient instance.
func NewTestClient(
	ctx context.Context,
	logger logging.Logger,
	metrics *TestClientMetrics,
	config *TestClientConfig) (*TestClient, error) {

	if config.SRSNumberToLoad == 0 {
		config.SRSNumberToLoad = config.MaxBlobSize / 32
	}

	// Construct the disperser client

	signer, err := auth.NewLocalBlobRequestSigner(config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create signer: %w", err)
	}
	accountId, err := signer.GetAccountID()
	if err != nil {
		return nil, fmt.Errorf("failed to get account ID: %w", err)
	}
	logger.Infof("Account ID: %s", accountId.String())

	g1Path, err := config.ResolveSRSPath(SRSPathG1)
	if err != nil {
		return nil, fmt.Errorf("resolve G1 SRS path: %w", err)
	}
	g2Path, err := config.ResolveSRSPath(SRSPathG2)
	if err != nil {
		return nil, fmt.Errorf("resolve G2 SRS path: %w", err)
	}
	g2TrailingPath, err := config.ResolveSRSPath(SRSPathG2Trailing)
	if err != nil {
		return nil, fmt.Errorf("resolve trailing G2 SRS path: %w", err)
	}
	srsTablesPath, err := config.ResolveSRSPath(SRSPathSRSTables)
	if err != nil {
		return nil, fmt.Errorf("resolve SRS tables path: %w", err)
	}

	// There is special logic for the trailing G2 point file. Some environments won't have a dedicated file for
	// trailing G2 points, and instead will simply have the unabridged G2 points (which definitionally contain the
	// trailing G2 points at the end of the file). If there isn't a trailing G2 point file in the expected location,
	// assume that the environment has access to the entire G2 point file, and pass in "" for the trailing path.
	// If this assumption turns out to be wrong, an error will be thrown when SRS parsing is attempted.
	if _, err := os.Stat(g2TrailingPath); errors.Is(err, os.ErrNotExist) {
		g2TrailingPath = ""
	}

	kzgCommitter, err := committer.NewFromConfig(committer.Config{
		G1SRSPath:         g1Path,
		G2SRSPath:         g2Path,
		G2TrailingSRSPath: g2TrailingPath,
		SRSNumberToLoad:   config.SRSNumberToLoad,
	})
	if err != nil {
		return nil, fmt.Errorf("new committer: %w", err)
	}

	var registry *prometheus.Registry
	if metrics != nil {
		registry = metrics.registry
	}

	accountantMetrics := metricsv2.NewAccountantMetrics(registry)
	dispersalMetrics := metricsv2.NewDispersalMetrics(registry)

	multiplexerConfig := dispersal.DefaultDisperserClientMultiplexerConfig()
	multiplexerConfig.DisperserConnectionCount = config.DisperserConnectionCount
	disperserRegistry := disperser.NewLegacyDisperserRegistry(
		fmt.Sprintf("%s:%d", config.DisperserHostname, config.DisperserPort))

	disperserClientMultiplexer, err := dispersal.NewDisperserClientMultiplexer(
		logger,
		multiplexerConfig,
		disperserRegistry,
		signer,
		kzgCommitter,
		dispersalMetrics,
		rand.New(rand.NewSource(time.Now().UnixNano())),
	)
	if err != nil {
		return nil, fmt.Errorf("create disperser client multiplexer: %w", err)
	}

	ethClientConfig := geth.EthClientConfig{
		RPCURLs:          config.EthRpcUrls,
		PrivateKeyString: config.PrivateKey,
		NumConfirmations: 0,
		NumRetries:       3,
	}
	ethClient, err := geth.NewMultiHomingClient(ethClientConfig, accountId, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create Ethereum client: %w", err)
	}

	chainId, err := ethClient.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get chain ID: %w", err)
	}

	contractDirectoryAddress := gethcommon.HexToAddress(config.ContractDirectoryAddress)
	contractDirectory, err := directory.NewContractDirectory(ctx, logger, ethClient, contractDirectoryAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create contract directory: %w", err)
	}

	operatorStateRetrieverAddress, err := contractDirectory.GetContractAddress(ctx, directory.OperatorStateRetriever)
	if err != nil {
		return nil, fmt.Errorf("failed to get OperatorStateRetriever address from contract directory: %w", err)
	}

	serviceManagerAddress, err := contractDirectory.GetContractAddress(ctx, directory.ServiceManager)
	if err != nil {
		return nil, fmt.Errorf("failed to get ServiceManager address from contract directory: %w", err)
	}

	ethReader, err := eth.NewReader(
		logger,
		ethClient,
		operatorStateRetrieverAddress.Hex(),
		serviceManagerAddress.Hex())
	if err != nil {
		return nil, fmt.Errorf("failed to create Ethereum reader: %w", err)
	}

	routerAddress, err := contractDirectory.GetContractAddress(ctx, directory.CertVerifierRouter)
	if err != nil {
		return nil, fmt.Errorf("failed to get CertVerifierRouter address from contract directory: %w", err)
	}

	certVerifierAddressProvider, err := verification.BuildRouterAddressProvider(routerAddress, ethClient, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create cert verifier address provider: %w", err)
	}

	certVerifier, err := verification.NewCertVerifier(logger, ethClient, certVerifierAddressProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to create cert verifier: %w", err)
	}

	// TODO (litt3): the PayloadPolynomialForm field included inside this config should be tested with different
	//  values, rather than just using the default. Consider a testing strategy that would exercise both encoding
	//  options.
	payloadClientConfig := clientsv2.GetDefaultPayloadClientConfig()

	payloadDisperserConfig := dispersal.PayloadDisperserConfig{
		PayloadClientConfig: *payloadClientConfig,
		DisperseBlobTimeout: 1337 * time.Hour, // this suite enforces its own timeouts
		BlobCompleteTimeout: 1337 * time.Hour, // this suite enforces its own timeouts
		ContractCallTimeout: 1337 * time.Hour, // this suite enforces its own timeouts
	}

	certBuilder, err := clientsv2.NewCertBuilder(logger,
		operatorStateRetrieverAddress,
		ethReader.GetRegistryCoordinatorAddress(),
		ethClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create cert builder: %w", err)
	}

	blockMon, err := verification.NewBlockNumberMonitor(
		logger,
		ethClient,
		time.Second*1,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create block number monitor: %w", err)
	}

	paymentVaultAddr, err := contractDirectory.GetContractAddress(ctx, directory.PaymentVault)
	if err != nil {
		return nil, fmt.Errorf("get PaymentVault address: %w", err)
	}

	clientLedger, err := buildClientLedger(
		ctx,
		logger,
		ethClient,
		paymentVaultAddr,
		accountId,
		clientledger.ClientLedgerMode(config.ClientLedgerPaymentMode),
		disperserClientMultiplexer,
		accountantMetrics,
	)
	if err != nil {
		return nil, fmt.Errorf("build client ledger: %w", err)
	}

	payloadDisperser, err := dispersal.NewPayloadDisperser(
		logger,
		payloadDisperserConfig,
		disperserClientMultiplexer,
		blockMon,
		certBuilder,
		certVerifier,
		clientLedger,
		registry)
	if err != nil {
		return nil, fmt.Errorf("failed to create payload disperser: %w", err)
	}

	// Construct the relay client

	// If the relay client attempts to call GetChunks(), it will use this bogus signer.
	// This is expected to be rejected by the relays, since this client is not authorized to call GetChunks().
	rand := random.NewTestRandom()
	keypair, err := rand.BLS()
	if err != nil {
		return nil, fmt.Errorf("failed to generate BLS keypair: %w", err)
	}

	var fakeSigner relay.MessageSigner = func(ctx context.Context, data [32]byte) (*core.Signature, error) {
		return keypair.SignMessage(data), nil
	}

	relayConfig := &relay.RelayClientConfig{
		UseSecureGrpcFlag:  true,
		MaxGRPCMessageSize: units.GiB,
		OperatorID:         &core.OperatorID{0},
		MessageSigner:      fakeSigner,
		ConnectionPoolSize: config.RelayConnectionCount,
	}

	relayUrlProvider, err := relay.NewRelayUrlProvider(ethClient, ethReader.GetRelayRegistryAddress())
	if err != nil {
		return nil, fmt.Errorf("create relay url provider: %w", err)
	}

	relayClient, err := relay.NewRelayClient(relayConfig, logger, relayUrlProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to create relay client: %w", err)
	}

	kzgConfig := &kzgv1.KzgConfig{
		LoadG2Points:    true,
		G1Path:          g1Path,
		G2Path:          g2Path,
		G2TrailingPath:  g2TrailingPath,
		CacheDir:        srsTablesPath,
		SRSOrder:        config.SrsOrder,
		SRSNumberToLoad: config.SRSNumberToLoad,
		NumWorker:       32,
	}
	verifierKzgConfig := verifier.ConfigFromV1KzgConfig(kzgConfig)
	encoder, err := rs.NewEncoder(logger, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create encoder: %w", err)
	}
	blobVerifier, err := verifier.NewVerifier(verifierKzgConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create blob verifier: %w", err)
	}

	relayPayloadRetrieverConfig := &payloadretrieval.RelayPayloadRetrieverConfig{
		PayloadClientConfig: *payloadClientConfig,
		RelayTimeout:        1337 * time.Hour, // this suite enforces its own timeouts
	}

	relayPayloadRetriever, err := payloadretrieval.NewRelayPayloadRetriever(
		logger,
		rand.Rand,
		*relayPayloadRetrieverConfig,
		relayClient,
		blobVerifier.G1SRS,
		metricsv2.NoopRetrievalMetrics)
	if err != nil {
		return nil, fmt.Errorf("failed to create relay payload retriever: %w", err)
	}

	// Construct the retrieval client

	chainState := eth.NewChainState(ethReader, ethClient)
	icsConfig := thegraph.Config{
		Endpoint:     config.SubgraphUrl,
		PullInterval: 100 * time.Millisecond,
		MaxRetries:   5,
	}
	indexedChainState := thegraph.MakeIndexedChainState(icsConfig, chainState, logger)

	validatorPayloadRetrieverConfig := &payloadretrieval.ValidatorPayloadRetrieverConfig{
		PayloadClientConfig: *payloadClientConfig,
		RetrievalTimeout:    1337 * time.Hour, // this suite enforces its own timeouts
	}

	validatorClientMetrics := validator.NewValidatorClientMetrics(registry)

	clientConfig := validator.DefaultClientConfig()
	clientConfig.ConnectionPoolSize = config.ValidatorReadConnectionPoolSize
	clientConfig.ComputePoolSize = config.ValidatorReadComputePoolSize
	retrievalClient := validator.NewValidatorClient(
		logger,
		ethReader,
		indexedChainState,
		encoder,
		blobVerifier,
		clientConfig,
		validatorClientMetrics)

	validatorPayloadRetriever, err := payloadretrieval.NewValidatorPayloadRetriever(
		logger,
		*validatorPayloadRetrieverConfig,
		retrievalClient,
		blobVerifier.G1SRS,
		metricsv2.NoopRetrievalMetrics)
	if err != nil {
		return nil, fmt.Errorf("failed to create validator payload retriever: %w", err)
	}

	// Create a client that only downloads the blob and does not verify it. Useful for load testing validator downloads
	// with limited CPU resources.
	onlyDownloadClientConfig := validator.DefaultClientConfig()
	onlyDownloadClientConfig.ConnectionPoolSize = config.ValidatorReadConnectionPoolSize
	onlyDownloadClientConfig.ComputePoolSize = config.ValidatorReadComputePoolSize
	onlyDownloadClientConfig.UnsafeChunkDeserializerFactory =
		mock.NewMockChunkDeserializerFactory(&mock.MockChunkDeserializer{})
	onlyDownloadClientConfig.UnsafeBlobDecoderFactory =
		mock.NewMockBlobDecoderFactory(&mock.MockBlobDecoder{})

	onlyDownloadValidatorClient := validator.NewValidatorClient(
		logger,
		ethReader,
		indexedChainState,
		encoder,
		blobVerifier,
		onlyDownloadClientConfig,
		validatorClientMetrics)

	proxyWrapper, err := NewProxyWrapper(ctx, logger,
		&proxyconfig.AppConfig{
			SecretConfig: proxycommon.SecretConfigV2{
				SignerPaymentKey: config.PrivateKey,
				EthRPCURL:        config.EthRpcUrls[0],
			},
			EnabledServersConfig: &enablement.EnabledServersConfig{
				Metric:      false,
				ArbCustomDA: false,
				RestAPIConfig: enablement.RestApisEnabled{
					Admin:               false,
					OpGenericCommitment: true,
					OpKeccakCommitment:  true,
					StandardCommitment:  true,
				},
			},
			RestSvrCfg: proxyserver.Config{
				Host: "localhost",
				Port: config.ProxyPort,
				// TODO (cody.littley) enable proxy metrics
				APIsEnabled: &enablement.RestApisEnabled{
					Admin:               false,
					OpGenericCommitment: true,
					OpKeccakCommitment:  true,
					StandardCommitment:  true,
				},
			},
			StoreBuilderConfig: builder.Config{
				StoreConfig: store.Config{
					BackendsToEnable: []proxycommon.EigenDABackend{proxycommon.V2EigenDABackend},
					DispersalBackend: proxycommon.V2EigenDABackend,
					AsyncPutWorkers:  32,
				},
				ClientConfigV2: proxycommon.ClientConfigV2{
					DisperserClientCfg: dispersal.DisperserClientConfig{
						GrpcUri:           fmt.Sprintf("%s:%d", config.DisperserHostname, config.DisperserPort),
						UseSecureGrpcFlag: true,
						DisperserID:       0,
						// use v0 for now, until all dispersers support v1
						RequestVersion: hashing.DisperseBlobRequestVersion0,
						ChainID:        chainId,
					},
					PayloadDisperserCfg: dispersal.PayloadDisperserConfig{
						PayloadClientConfig:    *payloadClientConfig,
						DisperseBlobTimeout:    5 * time.Minute,
						BlobCompleteTimeout:    5 * time.Minute,
						BlobStatusPollInterval: 1 * time.Second,
						ContractCallTimeout:    5 * time.Second,
					},
					RelayPayloadRetrieverCfg: payloadretrieval.RelayPayloadRetrieverConfig{
						PayloadClientConfig: *payloadClientConfig,
						RelayTimeout:        5 * time.Second,
					},
					ClientLedgerMode:                   clientledger.ParseClientLedgerMode(config.ClientLedgerPaymentMode),
					VaultMonitorInterval:               time.Second * 30,
					PutTries:                           3,
					MaxBlobSizeBytes:                   16 * units.MiB,
					EigenDACertVerifierOrRouterAddress: routerAddress.Hex(),
					EigenDADirectory:                   contractDirectoryAddress.Hex(),
					RetrieversToEnable: []proxycommon.RetrieverType{
						proxycommon.RelayRetrieverType,
						proxycommon.ValidatorRetrieverType,
					},
				},
				KzgConfig: *kzgConfig,
			},
		})
	if err != nil {
		return nil, fmt.Errorf("failed to create proxy wrapper: %w", err)
	}

	return &TestClient{
		config:                      config,
		payloadClientConfig:         payloadClientConfig,
		logger:                      logger,
		certVerifierAddressProvider: certVerifierAddressProvider,
		disperserClientMultiplexer:  disperserClientMultiplexer,
		payloadDisperser:            payloadDisperser,
		relayClient:                 relayClient,
		relayPayloadRetriever:       relayPayloadRetriever,
		indexedChainState:           indexedChainState,
		validatorClient:             retrievalClient,
		validatorPayloadRetriever:   validatorPayloadRetriever,
		certBuilder:                 certBuilder,
		onlyDownloadValidatorClient: onlyDownloadValidatorClient,
		certVerifier:                certVerifier,
		privateKey:                  config.PrivateKey,
		metricsRegistry:             registry,
		metrics:                     metrics,
		proxyWrapper:                proxyWrapper,
	}, nil
}

// formatPrivateKey formats the private key by removing leading/trailing whitespace and "0x" prefix.
func formatPrivateKey(privateKey string) string {
	privateKey = strings.Trim(privateKey, "\n \t")
	privateKey, _ = strings.CutPrefix(privateKey, "0x")
	return privateKey
}

// GetConfig returns the test client's configuration.
func (c *TestClient) GetConfig() *TestClientConfig {
	return c.config
}

// GetLogger returns the test client's logger.
func (c *TestClient) GetLogger() logging.Logger {
	return c.logger
}

// GetDisperserClient returns the test client's disperser client multiplexer.
func (c *TestClient) GetDisperserClientMultiplexer() *dispersal.DisperserClientMultiplexer {
	return c.disperserClientMultiplexer
}

// GetPayloadDisperser returns the test client's payload disperser.
func (c *TestClient) GetPayloadDisperser() *dispersal.PayloadDisperser {
	return c.payloadDisperser
}

// GetRelayClient returns the test client's relay client.
func (c *TestClient) GetRelayClient() relay.RelayClient {
	return c.relayClient
}

// GetRelayPayloadRetriever returns the test client's relay payload retriever.
func (c *TestClient) GetRelayPayloadRetriever() *payloadretrieval.RelayPayloadRetriever {
	return c.relayPayloadRetriever
}

// GetIndexedChainState returns the test client's indexed chain state.
func (c *TestClient) GetIndexedChainState() core.IndexedChainState {
	return c.indexedChainState
}

// GetValidatorClient returns the test client's validator client.
func (c *TestClient) GetValidatorClient() validator.ValidatorClient {
	return c.validatorClient
}

// GetValidatorPayloadRetriever returns the test client's validator payload retriever.
func (c *TestClient) GetValidatorPayloadRetriever() *payloadretrieval.ValidatorPayloadRetriever {
	return c.validatorPayloadRetriever
}

// GetCertVerifier returns the test client's cert verifier.
func (c *TestClient) GetCertVerifier() *verification.CertVerifier {
	return c.certVerifier
}

// GetCertBuilder returns the test client's cert builder.
func (c *TestClient) GetCertBuilder() *clientsv2.CertBuilder {
	return c.certBuilder
}

// GetPrivateKey returns the test client's private key.
func (c *TestClient) GetPrivateKey() string {
	return c.privateKey
}

// GetMetricsRegistry returns the test client's metrics registry.
func (c *TestClient) GetMetricsRegistry() *prometheus.Registry {
	return c.metricsRegistry
}

// Stop stops the test client.
func (c *TestClient) Stop() {
	c.metrics.stop()
	if c.proxyWrapper != nil {
		if err := c.proxyWrapper.Stop(); err != nil {
			c.logger.Errorf("failed to stop proxy wrapper: %v", err)
		}
	}
}

// DisperseAndVerify sends a payload to the disperser. Waits until the payload is confirmed and then reads
// it back from the relays and the validators.
func (c *TestClient) DisperseAndVerify(ctx context.Context, payload []byte) error {
	eigenDACert, err := c.DispersePayload(ctx, payload)
	if err != nil {
		return fmt.Errorf("failed to disperse payload: %w", err)
	}

	eigenDAV3Cert, ok := eigenDACert.(*coretypes.EigenDACertV3)
	if !ok {
		return fmt.Errorf("expected EigenDACertV3, got %T", eigenDACert)
	}

	blobKey, err := eigenDAV3Cert.ComputeBlobKey()
	if err != nil {
		return fmt.Errorf("failed to compute blob key: %w", err)
	}

	// read blob from a single relay (assuming success, otherwise will retry)
	payloadFromRelayRetriever, err := c.relayPayloadRetriever.GetPayload(ctx, eigenDAV3Cert)
	if err != nil {
		return fmt.Errorf("failed to get payload from relay: %w", err)
	}
	if !bytes.Equal(payload, payloadFromRelayRetriever) {
		return fmt.Errorf("payloads do not match")
	}

	// read blob from a single quorum (assuming success, otherwise will retry)
	payloadFromValidatorRetriever, err := c.validatorPayloadRetriever.GetPayload(ctx, eigenDAV3Cert)
	if err != nil {
		return fmt.Errorf("failed to get payload from validators: %w", err)
	}
	if !bytes.Equal(payload, payloadFromValidatorRetriever) {
		return fmt.Errorf("payloads do not match")
	}

	commitment, err := eigenDAV3Cert.Commitments()
	if err != nil {
		return fmt.Errorf("failed to parse blob commitments: %w", err)
	}

	blobLengthSymbols := commitment.Length

	// read blob from ALL relays
	err = c.ReadBlobFromRelays(
		ctx,
		blobKey,
		eigenDAV3Cert.RelayKeys(),
		payload,
		blobLengthSymbols,
		0)
	if err != nil {
		return fmt.Errorf("failed to read blob from relays: %w", err)
	}

	blobHeader, err := eigenDAV3Cert.BlobHeader()
	if err != nil {
		return fmt.Errorf("failed to get blob header from cert: %w", err)
	}

	// read blob from ALL quorums
	err = c.ReadBlobFromValidators(
		ctx,
		blobHeader,
		eigenDAV3Cert.BatchHeader.ReferenceBlockNumber,
		payload,
		0,
		true)
	if err != nil {
		return fmt.Errorf("failed to read blob from validators: %w", err)
	}

	return nil
}

// Similar to DisperseAndVerify, but uses the proxy instead of using the clients directly.
func (c *TestClient) DisperseAndVerifyWithProxy(ctx context.Context, payload []byte) error {
	cert, err := c.DispersePayloadWithProxy(ctx, payload)
	if err != nil {
		return fmt.Errorf("failed to disperse payload with proxy: %w", err)
	}

	_, err = c.ReadPayloadWithProxy(ctx, cert, payload, 0)
	if err != nil {
		return fmt.Errorf("failed to read payload with proxy: %w", err)
	}

	return nil
}

// DispersePayload sends a payload to the disperser. Returns the blob key.
func (c *TestClient) DispersePayload(ctx context.Context, payloadBytes []byte) (cert coretypes.EigenDACert, err error) {
	c.logger.Debugf("Dispersing payload of length %d", len(payloadBytes))
	start := time.Now()
	c.metrics.startOperation("dispersal")

	// Important: don't redefine err. It's used by the deferred function to report success or failure.

	defer func() {
		c.metrics.endOperation("dispersal")
		if err == nil {
			c.metrics.reportDispersalSuccess()
			c.metrics.reportDispersalTime(time.Since(start))
		} else {
			c.metrics.reportDispersalFailure()
		}
	}()

	payload := coretypes.Payload(payloadBytes)
	cert, err = c.GetPayloadDisperser().SendPayload(ctx, payload)

	if err != nil {
		return nil, fmt.Errorf("failed to disperse payload, %s", err)
	}

	return cert, nil
}

// DispersePayloadWithProxy sends a payload to the proxy wrapper, which then disperses it to EigenDA. Returns the cert
// in byte format, since that's what the proxy returns.
func (c *TestClient) DispersePayloadWithProxy(ctx context.Context, payloadBytes []byte) (cert []byte, err error) {
	if c.proxyWrapper == nil {
		return nil, fmt.Errorf("proxy wrapper not initialized")
	}
	c.logger.Debugf("Dispersing payload of length %d with proxy", len(payloadBytes))

	start := time.Now()
	c.metrics.startOperation("dispersal")

	// Important: don't redefine err. It's used by the deferred function to report success or failure.
	defer func() {
		c.metrics.endOperation("dispersal")
		if err == nil {
			c.metrics.reportDispersalTime(time.Since(start))
			c.metrics.reportDispersalSuccess()
		} else {
			c.metrics.reportDispersalFailure()
		}
	}()

	cert, err = c.proxyWrapper.SendPayload(ctx, payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to send payload via proxy: %w", err)
	}

	return cert, nil
}

// ReadBlobFromRelays reads a blob from the relays and compares it to the given payload.
//
// The timeout provided is a timeout for each individual relay read, not all reads as a whole.
func (c *TestClient) ReadBlobFromRelays(
	ctx context.Context,
	key corev2.BlobKey,
	relayKeys []corev2.RelayKey,
	expectedPayload []byte,
	blobLengthSymbols uint32,
	timeout time.Duration) error {

	for _, relayID := range relayKeys {
		err := c.ReadBlobFromRelay(ctx, key, relayID, expectedPayload, blobLengthSymbols, timeout)

		if err != nil {
			return fmt.Errorf("failed to read blob from relay %d: %w", relayID, err)
		}
	}

	return nil
}

// ReadBlobFromRelay reads a blob from the relay and compares it to the given payload.
func (c *TestClient) ReadBlobFromRelay(
	ctx context.Context,
	key corev2.BlobKey,
	relayKey corev2.RelayKey,
	expectedPayload []byte,
	blobLengthSymbols uint32,
	timeout time.Duration,
) error {

	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	// Important: don't redefine err. It's used by the deferred function to report success or failure.
	var err error

	c.metrics.startOperation("relay_read")
	start := time.Now()

	defer func() {
		c.metrics.endOperation("relay_read")
		if err == nil {
			c.metrics.reportRelayReadSuccess()
			c.metrics.reportRelayReadTime(time.Since(start), relayKey)
		} else {
			c.metrics.reportRelayReadFailure()
		}
	}()

	blobBytesFromRelay, err := c.relayClient.GetBlob(ctx, relayKey, key)
	if err != nil {
		return fmt.Errorf("failed to read blob from relay: %w", err)
	}

	blob, err := coretypes.DeserializeBlob(blobBytesFromRelay, blobLengthSymbols)
	if err != nil {
		return fmt.Errorf("failed to deserialize blob: %w", err)
	}

	payloadFromRelay, err := blob.ToPayload(c.payloadClientConfig.PayloadPolynomialForm)
	if err != nil {
		return fmt.Errorf("failed to decode blob: %w", err)
	}

	if !bytes.Equal(payloadFromRelay, expectedPayload) {
		return fmt.Errorf("payloads do not match")
	}

	return nil
}

// ReadBlobFromValidators reads a blob from the validators and compares it to the given payload.
//
// The timeout provided is a timeout for each read from a quorum, not all reads as a whole.
func (c *TestClient) ReadBlobFromValidators(
	ctx context.Context,
	header *corev2.BlobHeaderWithHashedPayment,
	referenceBlockNumber uint32,
	expectedPayloadBytes []byte,
	timeout time.Duration,
	validateAndDecode bool,
) error {

	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	// Important: don't redefine err. It's used by the deferred function to report success or failure.
	var err error

	c.metrics.startOperation("validator_read")
	start := time.Now()

	defer func() {
		c.metrics.endOperation("validator_read")
		if err == nil {
			if validateAndDecode {
				// Only report timing if we actually do the full operation. Skip report if we only download the blob.
				c.metrics.reportValidatorReadTime(time.Since(start))
			}
			c.metrics.reportValidatorReadSuccess()
		} else {
			c.metrics.reportValidatorReadFailure()
		}
	}()

	if validateAndDecode {
		var retrievedBlobBytes []byte
		retrievedBlobBytes, err = c.validatorClient.GetBlob(
			ctx,
			header,
			uint64(referenceBlockNumber))
		if err != nil {
			return fmt.Errorf("failed to read blob from validators, %s", err)
		}

		blobLengthSymbols := header.BlobCommitments.Length
		var blob *coretypes.Blob
		blob, err = coretypes.DeserializeBlob(retrievedBlobBytes, blobLengthSymbols)
		if err != nil {
			return fmt.Errorf("failed to deserialize blob: %w", err)
		}

		var retrievedPayload coretypes.Payload
		retrievedPayload, err = blob.ToPayload(c.payloadClientConfig.PayloadPolynomialForm)
		if err != nil {
			return fmt.Errorf("failed to convert blob to payload: %w", err)
		}

		if !bytes.Equal(retrievedPayload, expectedPayloadBytes) {
			return fmt.Errorf("payloads do not match")
		}
	} else {

		// Just download the blob without validating or decoding. Don't report timing metrics for this operation.

		_, err = c.onlyDownloadValidatorClient.GetBlob(
			ctx,
			header,
			uint64(referenceBlockNumber))
		if err != nil {
			return fmt.Errorf("failed to read blob from validators: %w", err)
		}
	}

	return nil
}

// ReadPayloadWithProxy reads a payload from the proxy wrapper and compares it to the expected payload bytes.
// The timeout is ignored if zero. If the proxy wrapper is not enabled, this method returns an error.
func (c *TestClient) ReadPayloadWithProxy(
	ctx context.Context,
	cert []byte,
	expectedPayloadBytes []byte,
	timeout time.Duration,
) ([]byte, error) {

	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	// Important: don't redefine err. It's used by the deferred function to report success or failure.
	var err error

	start := time.Now()
	c.metrics.startOperation("proxy_read")

	defer func() {
		c.metrics.endOperation("proxy_read")
		if err == nil {
			c.metrics.reportProxyReadSuccess()
			c.metrics.reportProxyReadTime(time.Since(start))
		} else {
			c.metrics.reportProxyReadFailure()
		}
	}()

	var data []byte
	data, err = c.proxyWrapper.GetPayload(ctx, cert)
	if err != nil {
		return nil, fmt.Errorf("failed to read payload from proxy: %w", err)
	}

	if !bytes.Equal(data, expectedPayloadBytes) {
		return nil, fmt.Errorf("read payload does not match expected payload")
	}

	return data, nil
}

// GetProxyWrapper returns the proxy wrapper. If the proxy wrapper is not enabled, this method returns an error.
func (c *TestClient) GetProxyWrapper() (*ProxyWrapper, error) {
	if c.proxyWrapper == nil {
		return nil, fmt.Errorf("proxy wrapper is not enabled in the test client configuration")
	}
	return c.proxyWrapper, nil
}

func (c *TestClient) EstimateGasAndReportCheckDACert(
	ctx context.Context,
	eigenDAV3Cert *coretypes.EigenDACertV3,
) (uint64, error) {
	gas, err := c.certVerifier.EstimateGasCheckDACert(ctx, eigenDAV3Cert)
	if err != nil {
		return 0, fmt.Errorf("failed to estimate gas for CheckDACert call: %w", err)
	}

	c.metrics.reportEstimateGasCheckDACert(gas)
	return gas, nil
}

func buildClientLedger(
	ctx context.Context,
	logger logging.Logger,
	ethClient common_eigenda.EthClient,
	paymentVaultAddr gethcommon.Address,
	accountID gethcommon.Address,
	mode clientledger.ClientLedgerMode,
	disperserClientMultiplexer *dispersal.DisperserClientMultiplexer,
	accountantMetrics metricsv2.AccountantMetricer,
) (*clientledger.ClientLedger, error) {
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
	switch mode {
	case clientledger.ClientLedgerModeReservationOnly:
		reservationLedger, err = buildReservationLedger(ctx, paymentVault, accountID, minNumSymbols)
		if err != nil {
			return nil, fmt.Errorf("build reservation ledger: %w", err)
		}
	case clientledger.ClientLedgerModeOnDemandOnly:
		cumulativePayment, err := getCumulativePayment(ctx, disperserClientMultiplexer)
		if err != nil {
			return nil, fmt.Errorf("get cumulative payment: %w", err)
		}
		onDemandLedger, err = buildOnDemandLedger(ctx, paymentVault, accountID, minNumSymbols, cumulativePayment)
		if err != nil {
			return nil, fmt.Errorf("build on-demand ledger: %w", err)
		}

	case clientledger.ClientLedgerModeReservationAndOnDemand:
		reservationLedger, err = buildReservationLedger(ctx, paymentVault, accountID, minNumSymbols)
		if err != nil {
			return nil, fmt.Errorf("build reservation ledger: %w", err)
		}
		cumulativePayment, err := getCumulativePayment(ctx, disperserClientMultiplexer)
		if err != nil {
			return nil, fmt.Errorf("get cumulative payment: %w", err)
		}
		onDemandLedger, err = buildOnDemandLedger(ctx, paymentVault, accountID, minNumSymbols, cumulativePayment)
		if err != nil {
			return nil, fmt.Errorf("build on-demand ledger: %w", err)
		}

	default:
		return nil, fmt.Errorf("unexpected client ledger mode: %s", mode)
	}

	ledger := clientledger.NewClientLedger(
		ctx,
		logger,
		accountantMetrics,
		accountID,
		mode,
		reservationLedger,
		onDemandLedger,
		time.Now,
		paymentVault,
		30*time.Second,
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
		ratelimit.OverfillOncePermitted,
		// TODO(litt3): once the checkpointed onchain config registry is ready, that should be used
		// instead of hardcoding. At that point, this field will be removed from the config struct
		// entirely, and the value will be fetched dynamically at runtime.
		60*time.Second,
	)
	if err != nil {
		return nil, fmt.Errorf("new reservation ledger config: %w", err)
	}

	reservationLedger, err := reservation.NewReservationLedger(*reservationConfig, time.Now)
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
	cumulativePayment *big.Int,
) (*ondemand.OnDemandLedger, error) {
	pricePerSymbol, err := paymentVault.GetPricePerSymbol(ctx)
	if err != nil {
		return nil, fmt.Errorf("get price per symbol: %w", err)
	}

	totalDeposits, err := paymentVault.GetTotalDeposit(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("get total deposit from vault: %w", err)
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

func getCumulativePayment(
	ctx context.Context,
	disperserClientMultiplexer *dispersal.DisperserClientMultiplexer,
) (*big.Int, error) {
	disperserClient, err := disperserClientMultiplexer.GetDisperserClient(ctx, time.Now(), true)
	if err != nil {
		return nil, fmt.Errorf("get disperser client: %w", err)
	}

	paymentState, err := disperserClient.GetPaymentState(ctx)
	if err != nil {
		return nil, fmt.Errorf("get payment state: %w", err)
	}

	if paymentState.GetCumulativePayment() == nil {
		return big.NewInt(0), nil
	}
	return new(big.Int).SetBytes(paymentState.GetCumulativePayment()), nil
}
