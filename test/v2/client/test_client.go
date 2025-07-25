package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	clientsv2 "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloaddispersal"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloadretrieval"
	"github.com/Layr-Labs/eigenda/api/clients/v2/relay"
	"github.com/Layr-Labs/eigenda/api/clients/v2/validator"
	"github.com/Layr-Labs/eigenda/api/clients/v2/validator/mock"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification/test"
	proxycommon "github.com/Layr-Labs/eigenda/api/proxy/common"
	proxymetrics "github.com/Layr-Labs/eigenda/api/proxy/metrics"
	proxyserver "github.com/Layr-Labs/eigenda/api/proxy/server"
	"github.com/Layr-Labs/eigenda/api/proxy/store"
	"github.com/Layr-Labs/eigenda/api/proxy/store/builder"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/litt/util"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	proxyconfig "github.com/Layr-Labs/eigenda/api/proxy/config"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/core"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/docker/go-units"
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
	certVerifierAddressProvider *test.TestCertVerifierAddressProvider
	disperserClient             clientsv2.DisperserClient
	payloadDisperser            *payloaddispersal.PayloadDisperser
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
	metrics                     *testClientMetrics
}

// NewTestClient creates a new TestClient instance.
func NewTestClient(
	logger logging.Logger,
	metrics *testClientMetrics,
	config *TestClientConfig) (*TestClient, error) {

	if config.SRSNumberToLoad == 0 {
		config.SRSNumberToLoad = config.MaxBlobSize / 32
	}

	// Construct the disperser client

	privateKey, err := loadPrivateKey(config.KeyPath, config.KeyVar)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}

	signer, err := auth.NewLocalBlobRequestSigner(privateKey)
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

	kzgConfig := &kzg.KzgConfig{
		LoadG2Points:    true,
		G1Path:          g1Path,
		G2Path:          g2Path,
		G2TrailingPath:  g2TrailingPath,
		CacheDir:        srsTablesPath,
		SRSOrder:        config.SRSOrder,
		SRSNumberToLoad: config.SRSNumberToLoad,
		NumWorker:       32,
	}

	kzgProver, err := prover.NewProver(kzgConfig, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create KZG prover: %w", err)
	}

	disperserConfig := &clientsv2.DisperserClientConfig{
		Hostname:          config.DisperserHostname,
		Port:              fmt.Sprintf("%d", config.DisperserPort),
		UseSecureGrpcFlag: true,
	}

	disperserClient, err := clientsv2.NewDisperserClient(disperserConfig, signer, kzgProver, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create disperser client: %w", err)
	}
	err = disperserClient.PopulateAccountant(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to populate accountant: %w", err)
	}

	ethRPCUrls, err := loadEthRPCURLs(config.EthRPCURLs, config.EthRPCUrlsVar)
	if err != nil {
		return nil, fmt.Errorf("failed to load Ethereum RPC URLs: %w", err)
	}

	ethClientConfig := geth.EthClientConfig{
		RPCURLs:          ethRPCUrls,
		PrivateKeyString: privateKey,
		NumConfirmations: 0,
		NumRetries:       3,
	}
	ethClient, err := geth.NewMultiHomingClient(ethClientConfig, accountId, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create Ethereum client: %w", err)
	}

	ethReader, err := eth.NewReader(
		logger,
		ethClient,
		config.EigenDADirectory,
		config.BLSOperatorStateRetrieverAddr,
		config.EigenDAServiceManagerAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create Ethereum reader: %w", err)
	}

	certVerifierAddressProvider := &test.TestCertVerifierAddressProvider{}

	certVerifier, err := verification.NewCertVerifier(logger, ethClient, certVerifierAddressProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to create cert verifier: %w", err)
	}

	// TODO (litt3): the PayloadPolynomialForm field included inside this config should be tested with different
	//  values, rather than just using the default. Consider a testing strategy that would exercise both encoding
	//  options.
	payloadClientConfig := clientsv2.GetDefaultPayloadClientConfig()

	payloadDisperserConfig := payloaddispersal.PayloadDisperserConfig{
		PayloadClientConfig: *payloadClientConfig,
		DisperseBlobTimeout: 1337 * time.Hour, // this suite enforces its own timeouts
		BlobCompleteTimeout: 1337 * time.Hour, // this suite enforces its own timeouts
		ContractCallTimeout: 1337 * time.Hour, // this suite enforces its own timeouts
	}

	var registry *prometheus.Registry
	if metrics != nil {
		registry = metrics.registry
	}

	certBuilder, err := clientsv2.NewCertBuilder(logger, gethcommon.HexToAddress(config.BLSOperatorStateRetrieverAddr), ethReader.GetRegistryCoordinatorAddress(), ethClient)
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

	payloadDisperser, err := payloaddispersal.NewPayloadDisperser(
		logger,
		payloadDisperserConfig,
		disperserClient,
		blockMon,
		certBuilder,
		certVerifier,
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
	}

	relayUrlProvider, err := relay.NewRelayUrlProvider(ethClient, ethReader.GetRelayRegistryAddress())
	if err != nil {
		return nil, fmt.Errorf("create relay url provider: %w", err)
	}

	relayClient, err := relay.NewRelayClient(relayConfig, logger, relayUrlProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to create relay client: %w", err)
	}

	verifierKzgConfig := kzgConfig
	verifierKzgConfig.LoadG2Points = false
	blobVerifier, err := verifier.NewVerifier(verifierKzgConfig, nil)
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
		blobVerifier.Srs.G1)
	if err != nil {
		return nil, fmt.Errorf("failed to create relay payload retriever: %w", err)
	}

	// Construct the retrieval client

	chainState := eth.NewChainState(ethReader, ethClient)
	icsConfig := thegraph.Config{
		Endpoint:     config.SubgraphURL,
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
		blobVerifier,
		clientConfig,
		validatorClientMetrics)

	validatorPayloadRetriever, err := payloadretrieval.NewValidatorPayloadRetriever(
		logger,
		*validatorPayloadRetrieverConfig,
		retrievalClient,
		blobVerifier.Srs.G1)
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
		blobVerifier,
		onlyDownloadClientConfig,
		validatorClientMetrics)

	proxyWrapper, err := NewProxyWrapper(context.Background(), logger,
		&proxyconfig.AppConfig{
			SecretConfig: proxycommon.SecretConfigV2{
				SignerPaymentKey: privateKey,
				EthRPCURL:        config.EthRPCURLs[0],
			},
			ServerConfig: proxyserver.Config{
				Host:        "localhost",
				Port:        config.ProxyPort,
				EnabledAPIs: []string{"admin"},
			},
			MetricsServerConfig: proxymetrics.Config{
				Enabled: false, // TODO enable this
			},
			StoreBuilderConfig: builder.Config{
				StoreConfig: store.Config{
					BackendsToEnable: []proxycommon.EigenDABackend{proxycommon.V2EigenDABackend},
					DispersalBackend: proxycommon.V2EigenDABackend,
					AsyncPutWorkers:  32,
				},
				ClientConfigV2: proxycommon.ClientConfigV2{
					DisperserClientCfg: clientsv2.DisperserClientConfig{
						Hostname:          config.DisperserHostname,
						Port:              fmt.Sprintf("%d", config.DisperserPort),
						UseSecureGrpcFlag: true,
					},
					PayloadDisperserCfg: payloaddispersal.PayloadDisperserConfig{
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
					PutTries:                           3,
					MaxBlobSizeBytes:                   16 * units.MiB,
					EigenDACertVerifierOrRouterAddress: config.EigenDACertVerifierAddressQuorums0_1,
					BLSOperatorStateRetrieverAddr:      config.BLSOperatorStateRetrieverAddr,
					EigenDAServiceManagerAddr:          config.EigenDAServiceManagerAddr,
					EigenDADirectory:                   config.EigenDADirectory,
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
		disperserClient:             disperserClient,
		payloadDisperser:            payloadDisperser,
		relayClient:                 relayClient,
		relayPayloadRetriever:       relayPayloadRetriever,
		indexedChainState:           indexedChainState,
		validatorClient:             retrievalClient,
		validatorPayloadRetriever:   validatorPayloadRetriever,
		certBuilder:                 certBuilder,
		onlyDownloadValidatorClient: onlyDownloadValidatorClient,
		certVerifier:                certVerifier,
		privateKey:                  privateKey,
		metricsRegistry:             registry,
		metrics:                     metrics,
		proxyWrapper:                proxyWrapper,
	}, nil
}

// loadPrivateKey loads the private key from the file/env var specified in the config.
func loadPrivateKey(keyPath string, keyVar string) (string, error) {
	var privateKey string
	if keyPath != "" {
		privateKeyFile, err := util.SanitizePath(keyPath)
		if err != nil {
			return "", fmt.Errorf("failed to sanitize path: %w", err)
		}

		exists, err := util.Exists(privateKeyFile)
		if err != nil {
			return "", fmt.Errorf("failed to check if private key file exists: %w", err)
		}
		if exists {
			privateKeyBytes, err := os.ReadFile(privateKeyFile)
			if err != nil {
				return "", fmt.Errorf("failed to read private key file: %w", err)
			}
			privateKey = string(privateKeyBytes)
		}
	}

	if privateKey == "" {
		if keyVar == "" {
			return "", fmt.Errorf("either KeyPath must reference a valid key file or KeyVar must be set")
		}
		privateKey = os.Getenv(keyVar)
		if privateKey == "" {
			return "", fmt.Errorf("key not found in environment variable %s", keyVar)
		}
	}

	return formatPrivateKey(privateKey), nil
}

// loadEthRPCURLs loads the Ethereum RPC URLs from the file/env var specified in the config.
func loadEthRPCURLs(urls []string, urlsVar string) ([]string, error) {
	if len(urls) > 0 {
		return urls, nil
	}

	if urlsVar == "" {
		return nil, fmt.Errorf("either EthRPCURLs or EthRPCUrlsVar must be set")
	}

	ethRPCURLs := os.Getenv(urlsVar)
	if ethRPCURLs == "" {
		return nil, fmt.Errorf("URLs not found in environment variable %s", urlsVar)
	}

	return strings.Split(ethRPCURLs, ","), nil
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

// SetCertVerifierAddress sets the address string which will be returned by the cert verifier address to all users of
// the provider
func (c *TestClient) SetCertVerifierAddress(certVerifierAddress string) {
	c.certVerifierAddressProvider.SetCertVerifierAddress(gethcommon.HexToAddress(certVerifierAddress))
}

// GetDisperserClient returns the test client's disperser client.
func (c *TestClient) GetDisperserClient() clientsv2.DisperserClient {
	return c.disperserClient
}

// GetPayloadDisperser returns the test client's payload disperser.
func (c *TestClient) GetPayloadDisperser() *payloaddispersal.PayloadDisperser {
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
	payloadBytesFromRelayRetriever := payloadFromRelayRetriever.Serialize()
	if !bytes.Equal(payload, payloadBytesFromRelayRetriever) {
		return fmt.Errorf("payloads do not match")
	}

	// read blob from a single quorum (assuming success, otherwise will retry)
	payloadFromValidatorRetriever, err := c.validatorPayloadRetriever.GetPayload(ctx, eigenDAV3Cert)
	if err != nil {
		return fmt.Errorf("failed to get payload from validators: %w", err)
	}
	payloadBytesFromValidatorRetriever := payloadFromValidatorRetriever.Serialize()
	if !bytes.Equal(payload, payloadBytesFromValidatorRetriever) {
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
		*blobKey,
		eigenDAV3Cert.RelayKeys(),
		payload,
		uint32(blobLengthSymbols),
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

	_, err = c.ReadBlobWithProxy(ctx, cert, payload, 0)
	if err != nil {
		return fmt.Errorf("failed to read blob with proxy: %w", err)
	}

	return nil
}

// DispersePayload sends a payload to the disperser. Returns the blob key.
func (c *TestClient) DispersePayload(ctx context.Context, payloadBytes []byte) (coretypes.EigenDACert, error) {
	c.logger.Debugf("Dispersing payload of length %d", len(payloadBytes))
	start := time.Now()
	c.metrics.startOperation("dispersal")

	// Important: don't redefine err. It's used by the deferred function to report success or failure.
	var err error
	defer func() {
		c.metrics.endOperation("dispersal")
		if err == nil {
			c.metrics.reportDispersalSuccess()
			c.metrics.reportDispersalTime(time.Since(start))
		} else {
			c.metrics.reportDispersalFailure()
		}
	}()

	payload := coretypes.NewPayload(payloadBytes)

	var cert coretypes.EigenDACert
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
			c.metrics.reportDispersalSuccess()
		} else {
			c.metrics.reportDispersalFailure()
		}
		c.metrics.reportDispersalTime(time.Since(start))
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

	payload, err := blob.ToPayload(c.payloadClientConfig.PayloadPolynomialForm)
	if err != nil {
		return fmt.Errorf("failed to decode blob: %w", err)
	}

	payloadBytesFromRelay := payload.Serialize()

	if !bytes.Equal(payloadBytesFromRelay, expectedPayload) {
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

		blobLengthSymbols := uint32(header.BlobCommitments.Length)
		var blob *coretypes.Blob
		blob, err = coretypes.DeserializeBlob(retrievedBlobBytes, blobLengthSymbols)
		if err != nil {
			return fmt.Errorf("failed to deserialize blob: %w", err)
		}

		var retrievedPayload *coretypes.Payload
		retrievedPayload, err = blob.ToPayload(c.payloadClientConfig.PayloadPolynomialForm)
		if err != nil {
			return fmt.Errorf("failed to convert blob to payload: %w", err)
		}

		payloadBytes := retrievedPayload.Serialize()
		if !bytes.Equal(payloadBytes, expectedPayloadBytes) {
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

// ReadBlobWithProxy reads a blob from the proxy wrapper and compares it to the expected payload bytes.
// The timeout is ignored if zero. If the proxy wrapper is not enabled, this method returns an error.
func (c *TestClient) ReadBlobWithProxy(ctx context.Context,
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
		return nil, fmt.Errorf("failed to read blob from proxy: %w", err)
	}

	if !bytes.Equal(data, expectedPayloadBytes) {
		return nil, fmt.Errorf("read payload does not match expected payload")
	}

	c.metrics.reportProxyReadTime(time.Since(start))

	return data, nil
}

// GetProxyWrapper returns the proxy wrapper. If the proxy wrapper is not enabled, this method returns an error.
func (c *TestClient) GetProxyWrapper() (*ProxyWrapper, error) {
	if c.proxyWrapper == nil {
		return nil, fmt.Errorf("proxy wrapper is not enabled in the test client configuration")
	}
	return c.proxyWrapper, nil
}
