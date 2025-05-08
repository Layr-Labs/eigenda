package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloaddispersal"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloadretrieval"
	"github.com/Layr-Labs/eigenda/api/clients/v2/relay"
	"github.com/Layr-Labs/eigenda/api/clients/v2/validator"
	"github.com/Layr-Labs/eigenda/api/clients/v2/validator/mock"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification/test"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
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
	payloadClientConfig         *clients.PayloadClientConfig
	logger                      logging.Logger
	certVerifierAddressProvider *test.TestCertVerifierAddressProvider
	disperserClient             clients.DisperserClient
	payloadDisperser            *payloaddispersal.PayloadDisperser
	relayClient                 relay.RelayClient
	relayPayloadRetriever       *payloadretrieval.RelayPayloadRetriever
	indexedChainState           core.IndexedChainState
	validatorClient             validator.ValidatorClient
	validatorPayloadRetriever   *payloadretrieval.ValidatorPayloadRetriever
	// For fetching blobs from the validators without verifying or decoding them. Useful for load testing
	// validator downloads with limited CPU resources.
	onlyDownloadValidatorClient validator.ValidatorClient
	certVerifier                *verification.GenericCertVerifier
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

	disperserConfig := &clients.DisperserClientConfig{
		Hostname:          config.DisperserHostname,
		Port:              fmt.Sprintf("%d", config.DisperserPort),
		UseSecureGrpcFlag: true,
	}

	disperserClient, err := clients.NewDisperserClient(disperserConfig, signer, kzgProver, nil)
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

	certVerifierAddressProvider := &test.TestCertVerifierAddressProvider{}

	certVerifier, err := verification.NewGenericCertVerifier(logger, ethClient, certVerifierAddressProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to create cert verifier: %w", err)
	}

	// TODO (litt3): the PayloadPolynomialForm field included inside this config should be tested with different
	//  values, rather than just using the default. Consider a testing strategy that would exercise both encoding
	//  options.
	payloadClientConfig := clients.GetDefaultPayloadClientConfig()

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

	certBuilder, err := clients.NewCertBuilder(logger, gethcommon.HexToAddress(config.BLSOperatorStateRetrieverAddr), gethcommon.HexToAddress(config.EigenDARegistryCoordinatorAddress), ethClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create cert builder: %w", err)
	}

	blockMon, err := verification.NewBlockNumberMonitor(
		logger,
		ethClient,
		time.Second * 12,
	)

	payloadDisperser, err := payloaddispersal.NewPayloadDisperser(
		logger,
		payloadDisperserConfig,
		ethClient,
		disperserClient,
		blockMon,
		certBuilder,
		certVerifier,
		registry)
	if err != nil {
		return nil, fmt.Errorf("failed to create payload disperser: %w", err)
	}

	// Construct the relay client

	ethReader, err := eth.NewReader(
		logger,
		ethClient,
		config.BLSOperatorStateRetrieverAddr,
		config.EigenDAServiceManagerAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create Ethereum reader: %w", err)
	}

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
		onlyDownloadValidatorClient: onlyDownloadValidatorClient,
		certVerifier:                certVerifier,
		privateKey:                  privateKey,
		metricsRegistry:             registry,
		metrics:                     metrics,
	}, nil
}

// loadPrivateKey loads the private key from the file/env var specified in the config.
func loadPrivateKey(keyPath string, keyVar string) (string, error) {
	if keyPath != "" {
		privateKeyFile, err := ResolveTildeInPath(keyPath)
		if err != nil {
			return "", fmt.Errorf("failed to resolve tilde in path: %w", err)
		}
		privateKey, err := os.ReadFile(privateKeyFile)
		if err != nil {
			return "", fmt.Errorf("failed to read private key file: %w", err)
		}

		return formatPrivateKey(string(privateKey)), nil
	}

	if keyVar == "" {
		return "", fmt.Errorf("either KeyPath or KeyVar must be set")
	}
	privateKey := os.Getenv(keyVar)
	if privateKey == "" {
		return "", fmt.Errorf("key not found in environment variable %s", keyVar)
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
func (c *TestClient) GetDisperserClient() clients.DisperserClient {
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
func (c *TestClient) GetCertVerifier() *verification.GenericCertVerifier {
	return c.certVerifier
}

// GetCertBuilder returns the test client's cert builder.
func (c *TestClient) GetCertBuilder() *clients.CertBuilder {
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
}

// DisperseAndVerify sends a payload to the disperser. Waits until the payload is confirmed and then reads
// it back from the relays and the validators.
func (c *TestClient) DisperseAndVerify(ctx context.Context, payload []byte) error {
	start := time.Now()
	eigenDACert, err := c.DispersePayload(ctx, payload)
	if err != nil {
		return fmt.Errorf("failed to disperse payload: %w", err)
	}
	c.metrics.reportCertificationTime(time.Since(start))

	blobKey, err := eigenDACert.ComputeBlobKey()
	if err != nil {
		return fmt.Errorf("failed to compute blob key: %w", err)
	}

	// read blob from a single relay (assuming success, otherwise will retry)
	payloadFromRelayRetriever, err := c.relayPayloadRetriever.GetPayload(ctx, eigenDACert)
	if err != nil {
		return fmt.Errorf("failed to get payload from relay: %w", err)
	}
	payloadBytesFromRelayRetriever := payloadFromRelayRetriever.Serialize()
	if !bytes.Equal(payload, payloadBytesFromRelayRetriever) {
		return fmt.Errorf("payloads do not match")
	}

	// read blob from a single quorum (assuming success, otherwise will retry)
	payloadFromValidatorRetriever, err := c.validatorPayloadRetriever.GetPayload(ctx, eigenDACert)
	if err != nil {
		return fmt.Errorf("failed to get payload from validators: %w", err)
	}
	payloadBytesFromValidatorRetriever := payloadFromValidatorRetriever.Serialize()
	if !bytes.Equal(payload, payloadBytesFromValidatorRetriever) {
		return fmt.Errorf("payloads do not match")
	}

	commitment, err := eigenDACert.Commitments()
	if err != nil {
		return fmt.Errorf("failed to parse blob commitments: %w", err)
	}

	blobLengthSymbols := commitment.Length

	// read blob from ALL relays
	err = c.ReadBlobFromRelays(
		ctx,
		*blobKey,
		eigenDACert.RelayKeys(),
		payload,
		blobLengthSymbols,
		0)
	if err != nil {
		return fmt.Errorf("failed to read blob from relays: %w", err)
	}


	// read blob from ALL quorums
	err = c.ReadBlobFromValidators(
		ctx,
		*blobKey,
		eigenDACert.BlobVersion(),
		*commitment,
		eigenDACert.QuorumNumbers(),
		eigenDACert.ReferenceBlockNumber(),
		payload,
		0,
		true)
	if err != nil {
		return fmt.Errorf("failed to read blob from validators: %w", err)
	}

	return nil
}

// DispersePayload sends a payload to the disperser. Returns the blob key.
func (c *TestClient) DispersePayload(ctx context.Context, payloadBytes []byte) (coretypes.EigenDACert, error) {
	c.logger.Debugf("Dispersing payload of length %d", len(payloadBytes))
	start := time.Now()

	payload := coretypes.NewPayload(payloadBytes)

	cert, err := c.GetPayloadDisperser().SendPayload(ctx, payload)

	if err != nil {
		return nil, fmt.Errorf("failed to disperse payload, %s", err)
	}

	c.metrics.reportDispersalTime(time.Since(start))

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

	start := time.Now()

	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	blobBytesFromRelay, err := c.relayClient.GetBlob(ctx, relayKey, key)
	if err != nil {
		return fmt.Errorf("failed to read blob from relay: %w", err)
	}

	c.metrics.reportRelayReadTime(time.Since(start), relayKey)

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
	blobKey corev2.BlobKey,
	blobVersion corev2.BlobVersion,
	blobCommitments encoding.BlobCommitments,
	quorums []core.QuorumID,
	referenceBlockNumber uint32,
	expectedPayloadBytes []byte,
	timeout time.Duration,
	validateAndDecode bool) error {

	for _, quorumID := range quorums {
		err := c.ReadBlobFromValidatorsInQuorum(
			ctx,
			blobKey,
			blobVersion,
			blobCommitments,
			quorumID,
			referenceBlockNumber,
			expectedPayloadBytes,
			timeout,
			validateAndDecode)
		if err != nil {
			return fmt.Errorf("failed to read blob from validators in quorum %d: %w", quorumID, err)
		}
	}

	return nil
}

// ReadBlobFromValidatorsInQuorum reads a blob from the validators in a specific quorum and compares it to
// the given payload.
func (c *TestClient) ReadBlobFromValidatorsInQuorum(
	ctx context.Context,
	blobKey corev2.BlobKey,
	blobVersion corev2.BlobVersion,
	blobCommitments encoding.BlobCommitments,
	quorumID core.QuorumID,
	referenceBlockNumber uint32,
	expectedPayloadBytes []byte,
	timeout time.Duration,
	validateAndDecode bool) error {

	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	if validateAndDecode {
		start := time.Now()

		retrievedBlobBytes, err := c.validatorClient.GetBlob(
			ctx,
			blobKey,
			blobVersion,
			blobCommitments,
			uint64(referenceBlockNumber),
			quorumID)
		if err != nil {
			return fmt.Errorf("failed to read blob from validators, %s", err)
		}

		c.metrics.reportValidatorReadTime(time.Since(start), quorumID)

		blobLengthSymbols := uint32(blobCommitments.Length)
		blob, err := coretypes.DeserializeBlob(retrievedBlobBytes, blobLengthSymbols)
		if err != nil {
			return fmt.Errorf("failed to deserialize blob: %w", err)
		}

		retrievedPayload, err := blob.ToPayload(c.payloadClientConfig.PayloadPolynomialForm)
		if err != nil {
			return fmt.Errorf("failed to convert blob to payload: %w", err)
		}

		payloadBytes := retrievedPayload.Serialize()
		if !bytes.Equal(payloadBytes, expectedPayloadBytes) {
			return fmt.Errorf("payloads do not match")
		}
	} else {

		// Just download the blob without validating or decoding. Don't report timing metrics for this operation.

		_, err := c.onlyDownloadValidatorClient.GetBlob(
			ctx,
			blobKey,
			blobVersion,
			blobCommitments,
			uint64(referenceBlockNumber),
			quorumID)
		if err != nil {
			return fmt.Errorf("failed to read blob from validators: %w", err)
		}
	}

	return nil
}
