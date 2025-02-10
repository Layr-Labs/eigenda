package client

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	v2 "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
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
	gethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	SRSPath           = "srs"
	SRSPathG1         = SRSPath + "/g1.point"
	SRSPathG2         = SRSPath + "/g2.point"
	SRSPathG2PowerOf2 = SRSPath + "/g2.point.powerOf2"
	SRSPathSRSTables  = SRSPath + "/SRSTables"
)

// TestClient encapsulates the various clients necessary for interacting with EigenDA.
type TestClient struct {
	Config                    *TestClientConfig
	Logger                    logging.Logger
	DisperserClient           clients.DisperserClient
	PayloadDisperser          *clients.PayloadDisperser
	RelayClient               clients.RelayClient
	RelayPayloadRetriever     *clients.RelayPayloadRetriever
	indexedChainState         core.IndexedChainState
	RetrievalClient           clients.RetrievalClient
	ValidatorPayloadRetriever *clients.ValidatorPayloadRetriever
	CertVerifier              *verification.CertVerifier
	PrivateKey                string
	MetricsRegistry           *prometheus.Registry
	metrics                   *testClientMetrics
	quorums                   []core.QuorumID
	blobCodec                 codecs.BlobCodec
}

// ResolveTildeInPath resolves the tilde (~) in the given path to the user's home directory.
func ResolveTildeInPath(path string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	return strings.Replace(path, "~", homeDir, 1), nil
}

// path returns the full path to a file in the test data directory.
func (c *TestClientConfig) path(elements ...string) (string, error) {
	root, err := ResolveTildeInPath(c.TestDataPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve tilde in path: %w", err)
	}

	combinedElements := make([]string, 0, len(elements)+1)
	combinedElements = append(combinedElements, root)
	combinedElements = append(combinedElements, elements...)

	return path.Join(combinedElements...), nil
}

// NewTestClient creates a new TestClient instance.
func NewTestClient(
	logger logging.Logger,
	metrics *testClientMetrics,
	config *TestClientConfig,
	quorums []core.QuorumID) (*TestClient, error) {

	if config.SRSNumberToLoad == 0 {
		config.SRSNumberToLoad = config.MaxBlobSize / 32
	}

	// Construct the disperser client

	privateKeyFile, err := ResolveTildeInPath(config.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve tilde in path: %w", err)
	}
	privateKey, err := os.ReadFile(privateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	privateKeyString := string(privateKey)
	privateKeyString = strings.Trim(privateKeyString, "\n \t")
	privateKeyString, _ = strings.CutPrefix(privateKeyString, "0x")

	signer, err := auth.NewLocalBlobRequestSigner(privateKeyString)
	if err != nil {
		return nil, fmt.Errorf("failed to create signer: %w", err)
	}
	signerAccountId, err := signer.GetAccountID()
	if err != nil {
		return nil, fmt.Errorf("failed to get account ID: %w", err)
	}
	accountId := gethcommon.HexToAddress(signerAccountId)
	logger.Infof("Account ID: %s", accountId.String())

	g1Path, err := config.path(SRSPathG1)
	if err != nil {
		return nil, fmt.Errorf("failed to get path to G1 file: %w", err)
	}
	g2Path, err := config.path(SRSPathG2)
	if err != nil {
		return nil, fmt.Errorf("failed to get path to G2 file: %w", err)
	}
	g2PowerOf2Path, err := config.path(SRSPathG2PowerOf2)
	if err != nil {
		return nil, fmt.Errorf("failed to get path to G2 power of 2 file: %w", err)
	}
	srsTablesPath, err := config.path(SRSPathSRSTables)
	if err != nil {
		return nil, fmt.Errorf("failed to get path to SRS tables: %w", err)
	}

	kzgConfig := &kzg.KzgConfig{
		LoadG2Points:    true,
		G1Path:          g1Path,
		G2Path:          g2Path,
		G2PowerOf2Path:  g2PowerOf2Path,
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
	} else {
		logger.Warn("Accountant populated") // TODO delete
	}

	ethClientConfig := geth.EthClientConfig{
		RPCURLs:          config.EthRPCURLs,
		PrivateKeyString: privateKeyString,
		NumConfirmations: 0,
		NumRetries:       3,
	}
	ethClient, err := geth.NewMultiHomingClient(ethClientConfig, accountId, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create Ethereum client: %w", err)
	}

	certVerifier, err := verification.NewCertVerifier(
		logger,
		ethClient,
		config.EigenDACertVerifierAddress,
		time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to create cert verifier: %w", err)
	}

	payloadClientConfig := clients.GetDefaultPayloadClientConfig()
	payloadClientConfig.EigenDACertVerifierAddr = config.EigenDACertVerifierAddress

	payloadDisperserConfig := &clients.PayloadDisperserConfig{
		PayloadClientConfig: *payloadClientConfig,
		Quorums:             quorums,
		DisperseBlobTimeout: 1337 * time.Hour, // this suite enforces its own timeouts
	}

	blobCodec, err := codecs.CreateCodec(codecs.PolynomialFormEval, payloadDisperserConfig.BlobEncodingVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to create blob codec: %w", err)
	}

	payloadDisperser, err := clients.NewPayloadDisperser(
		logger,
		*payloadDisperserConfig,
		blobCodec,
		disperserClient,
		certVerifier)
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

	relayURLS, err := ethReader.GetRelayURLs(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get relay URLs: %w", err)
	}

	// If the relay client attempts to call GetChunks(), it will use this bogus signer.
	// This is expected to be rejected by the relays, since this client is not authorized to call GetChunks().
	rand := random.NewTestRandom(nil)
	keypair := rand.BLS()
	var fakeSigner clients.MessageSigner = func(ctx context.Context, data [32]byte) (*core.Signature, error) {
		return keypair.SignMessage(data), nil
	}

	relayConfig := &clients.RelayClientConfig{
		Sockets:            relayURLS,
		UseSecureGrpcFlag:  true,
		MaxGRPCMessageSize: units.GiB,
		OperatorID:         &core.OperatorID{0},
		MessageSigner:      fakeSigner,
	}
	relayClient, err := clients.NewRelayClient(relayConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create relay client: %w", err)
	}

	blobVerifier, err := verifier.NewVerifier(kzgConfig, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create blob verifier: %w", err)
	}

	relayPayloadRetrieverConfig := &clients.RelayPayloadRetrieverConfig{
		PayloadClientConfig: *payloadClientConfig,
		RelayTimeout:        1337 * time.Hour, // this suite enforces its own timeouts
	}

	relayPayloadRetriever, err := clients.NewRelayPayloadRetriever(
		logger,
		rand.Rand,
		*relayPayloadRetrieverConfig,
		relayClient,
		blobCodec,
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

	validatorPayloadRetrieverConfig := &clients.ValidatorPayloadRetrieverConfig{
		PayloadClientConfig:           *payloadClientConfig,
		MaxConnectionCount:            20,
		BlsOperatorStateRetrieverAddr: config.BLSOperatorStateRetrieverAddr,
		EigenDAServiceManagerAddr:     config.EigenDAServiceManagerAddr,
		RetrievalTimeout:              1337 * time.Hour, // this suite enforces its own timeouts
	}

	retrievalClient := clients.NewRetrievalClient(
		logger,
		ethReader,
		indexedChainState,
		blobVerifier,
		int(validatorPayloadRetrieverConfig.MaxConnectionCount))

	validatorPayloadRetriever, err := clients.NewValidatorPayloadRetriever(
		logger,
		*validatorPayloadRetrieverConfig,
		blobCodec,
		retrievalClient,
		blobVerifier.Srs.G1)
	if err != nil {
		return nil, fmt.Errorf("failed to create validator payload retriever: %w", err)
	}

	return &TestClient{
		Config:                    config,
		Logger:                    logger,
		DisperserClient:           disperserClient,
		PayloadDisperser:          payloadDisperser,
		RelayClient:               relayClient,
		RelayPayloadRetriever:     relayPayloadRetriever,
		indexedChainState:         indexedChainState,
		RetrievalClient:           retrievalClient,
		ValidatorPayloadRetriever: validatorPayloadRetriever,
		CertVerifier:              certVerifier,
		PrivateKey:                privateKeyString,
		MetricsRegistry:           metrics.registry,
		metrics:                   metrics,
		quorums:                   quorums,
		blobCodec:                 blobCodec,
	}, nil
}

// Stop stops the test client.
func (c *TestClient) Stop() {
	c.metrics.stop()
}

// DisperseAndVerify sends a payload to the disperser. Waits until the payload is confirmed and then reads
// it back from the relays and the validators.
func (c *TestClient) DisperseAndVerify(
	ctx context.Context,
	payload []byte,
	salt uint32) error {

	start := time.Now()
	eigenDACert, err := c.DispersePayload(ctx, payload, salt)
	if err != nil {
		return fmt.Errorf("failed to disperse payload: %w", err)
	}
	c.metrics.reportCertificationTime(time.Since(start))

	blobKey, err := eigenDACert.ComputeBlobKey()
	if err != nil {
		return fmt.Errorf("failed to compute blob key: %w", err)
	}

	// read blob from a single relay (assuming success, otherwise will retry)
	payloadBytesFromRelayRetriever, err := c.RelayPayloadRetriever.GetPayload(ctx, eigenDACert)
	if err != nil {
		return fmt.Errorf("failed to read blob from relay: %w", err)
	}
	if !bytes.Equal(payload, payloadBytesFromRelayRetriever) {
		return fmt.Errorf("payloads do not match")
	}

	// read blob from a single quorum (assuming success, otherwise will retry)
	payloadBytesFromValidatorRetriever, err := c.ValidatorPayloadRetriever.GetPayload(ctx, eigenDACert)
	if err != nil {
		return fmt.Errorf("failed to read blob from validators: %w", err)
	}
	if !bytes.Equal(payload, payloadBytesFromValidatorRetriever) {
		return fmt.Errorf("payloads do not match")
	}

	// read blob from ALL relays
	err = c.ReadBlobFromRelays(ctx, *blobKey, eigenDACert.BlobInclusionInfo.BlobCertificate.RelayKeys, payload)
	if err != nil {
		return fmt.Errorf("failed to read blob from relays: %w", err)
	}

	blobHeader := eigenDACert.BlobInclusionInfo.BlobCertificate.BlobHeader
	commitment, err := verification.BlobCommitmentsBindingToInternal(&blobHeader.Commitment)
	if err != nil {
		return fmt.Errorf("failed to convert blob commitments: %w", err)
	}

	// read blob from ALL quorums
	err = c.ReadBlobFromValidators(
		ctx,
		*blobKey,
		blobHeader.Version,
		*commitment,
		eigenDACert.BlobInclusionInfo.BlobCertificate.BlobHeader.QuorumNumbers,
		payload)
	if err != nil {
		return fmt.Errorf("failed to read blob from validators: %w", err)
	}

	return nil
}

// DispersePayload sends a payload to the disperser. Returns the blob key.
func (c *TestClient) DispersePayload(
	ctx context.Context,
	payload []byte,
	salt uint32) (*verification.EigenDACert, error) {

	c.Logger.Debugf("Dispersing payload of length %d", len(payload))
	start := time.Now()

	cert, err := c.PayloadDisperser.SendPayload(ctx, payload, salt)

	if err != nil {
		return nil, fmt.Errorf("failed to disperse payload: %w", err)
	}
	c.metrics.reportDispersalTime(time.Since(start))

	return cert, nil
}

// VerifyBlobCertification verifies that the blob has been properly certified by the network.
func (c *TestClient) VerifyBlobCertification(
	key corev2.BlobKey,
	expectedQuorums []core.QuorumID,
	signedBatch *v2.SignedBatch,
	inclusionInfo *v2.BlobInclusionInfo) error {

	blobCert := inclusionInfo.BlobCertificate
	if blobCert == nil {
		return fmt.Errorf("missing blob certificate")
	}
	if len(blobCert.RelayKeys) == 0 {
		return fmt.Errorf("missing relay keys")
	}

	// make sure the returned header hash matches the expected blob key
	bh, err := corev2.BlobHeaderFromProtobuf(blobCert.BlobHeader)
	if err != nil {
		return fmt.Errorf("failed to convert blob header: %w", err)
	}
	computedBlobKey, err := bh.BlobKey()
	if err != nil {
		return fmt.Errorf("failed to compute blob key: %w", err)
	}
	if !bytes.Equal(key[:], computedBlobKey[:]) {
		return fmt.Errorf("blob key mismatch: expected %x, got %x", key, computedBlobKey)
	}

	// verify that expected quorums are present
	quorumSet := make(map[core.QuorumID]struct{}, len(expectedQuorums))
	for _, quorumNumber := range signedBatch.Attestation.QuorumNumbers {
		quorumSet[core.QuorumID(quorumNumber)] = struct{}{}
	}
	// There may be other quorums in the batch. No biggie as long as the expected ones are there.
	if len(expectedQuorums) > len(quorumSet) {
		return fmt.Errorf("missing quorums: expected %v, got %v", expectedQuorums, quorumSet)
	}
	for expectedQuorum := range quorumSet {
		if _, ok := quorumSet[expectedQuorum]; !ok {
			return fmt.Errorf("missing quorum %d", expectedQuorum)
		}
	}

	// Check the signing percentages
	signingPercents := make(map[core.QuorumID]int, len(signedBatch.Attestation.QuorumNumbers))
	for i, quorumNumber := range signedBatch.Attestation.QuorumNumbers {
		percent := int(signedBatch.Attestation.QuorumSignedPercentages[i])
		signingPercents[core.QuorumID(quorumNumber)] = percent
	}
	for _, quorum := range expectedQuorums {
		percent, ok := signingPercents[quorum]
		if !ok {
			return fmt.Errorf("missing quorum %d", quorum)
		}
		if percent < 0 || percent > 100 {
			return fmt.Errorf("invalid signing percentage %d", percent)
		}
		if percent < c.Config.MinimumSigningPercent {
			return fmt.Errorf("quorum %d signed by only %d%%", quorum, percent)
		}
	}

	// On-chain verification
	err = c.CertVerifier.VerifyCertV2FromSignedBatch(context.Background(), signedBatch, inclusionInfo)
	if err != nil {
		return fmt.Errorf("failed to verify cert: %w", err)
	}

	return nil
}

// ReadBlobFromRelays reads a blob from the relays and compares it to the given payload.
func (c *TestClient) ReadBlobFromRelays(
	ctx context.Context,
	key corev2.BlobKey,
	relayKeys []corev2.RelayKey,
	expectedPayload []byte) error {

	for _, relayID := range relayKeys {
		start := time.Now()

		c.Logger.Debugf("Reading blob from relay %d", relayID)
		blobFromRelay, err := c.RelayClient.GetBlob(ctx, relayID, key)
		if err != nil {
			return fmt.Errorf("failed to read blob from relay: %w", err)
		}

		c.metrics.reportRelayReadTime(time.Since(start), relayID)

		payloadBytesFromRelay, err := c.blobCodec.DecodeBlob(blobFromRelay)
		if err != nil {
			return fmt.Errorf("failed to decode blob: %w", err)
		}

		if !bytes.Equal(payloadBytesFromRelay, expectedPayload) {
			return fmt.Errorf("payloads do not match")
		}
	}

	return nil
}

// ReadBlobFromValidators reads a blob from the validators and compares it to the given payload.
func (c *TestClient) ReadBlobFromValidators(
	ctx context.Context,
	blobKey corev2.BlobKey,
	blobVersion corev2.BlobVersion,
	blobCommitments encoding.BlobCommitments,
	quorums []core.QuorumID,
	expectedPayload []byte) error {

	currentBlockNumber, err := c.indexedChainState.GetCurrentBlockNumber()
	if err != nil {
		return fmt.Errorf("failed to get current block number: %w", err)
	}

	for _, quorumID := range quorums {
		c.Logger.Debugf("Reading blob from validators for quorum %d", quorumID)

		start := time.Now()

		retrievedBlob, err := c.RetrievalClient.GetBlob(
			ctx,
			blobKey,
			blobVersion,
			blobCommitments,
			uint64(currentBlockNumber),
			quorumID)
		if err != nil {
			return fmt.Errorf("failed to read blob from validators: %w", err)
		}

		c.metrics.reportValidatorReadTime(time.Since(start), quorumID)

		retrievedPayload, err := c.blobCodec.DecodeBlob(retrievedBlob)
		if err != nil {
			return fmt.Errorf("failed to decode blob: %w", err)
		}
		if !bytes.Equal(retrievedPayload, expectedPayload) {
			return fmt.Errorf("payloads do not match")
		}
	}

	return nil
}
