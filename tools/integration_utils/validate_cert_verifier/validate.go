package validate_cert_verifier

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/api/clients/v2/dispersal"
	"github.com/Layr-Labs/eigenda/api/clients/v2/metrics"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	proxycommon "github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	"github.com/Layr-Labs/eigenda/core/eth/directory"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg/committer"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"
)

func RunCreateAndValidateCertValidation(c *cli.Context) error {
	ctx := context.Background()

	networkStr := c.String("eigenda-network")
	jsonRPCURL := c.String("json-rpc-url")
	signerAuthKey := c.String("signer-auth-key")
	srsPath := c.String("srs-path")
	certVerifierAddrStr := c.String("cert-verifier-address")

	logger, err := createLogger()
	if err != nil {
		return fmt.Errorf("create logger: %w", err)
	}

	// Parse and validate the network
	network, err := proxycommon.EigenDANetworkFromString(networkStr)
	if err != nil {
		return fmt.Errorf("parse network: %w", err)
	}

	// Get network configuration
	disperserHostname := network.GetDisperserAddress()
	eigenDADirectoryAddr := gethcommon.HexToAddress(network.GetEigenDADirectory())

	// Parse cert verifier address override if provided
	var certVerifierAddrOverride *gethcommon.Address
	if certVerifierAddrStr != "" {
		addr := gethcommon.HexToAddress(certVerifierAddrStr)
		certVerifierAddrOverride = &addr
		logger.Info("Using cert verifier address override", "address", addr.Hex())
	}

	logger.Info("Starting validate-cert-verifier tool",
		"network", network,
		"disperserHostname", disperserHostname,
		"eigenDADirectoryAddr", eigenDADirectoryAddr.Hex(),
		"jsonRPCURL", jsonRPCURL)

	// Initialize the payload disperser
	payloadDisperser, ethClient, certVerifierAddr, err := initializePayloadDisperser(
		ctx,
		logger,
		disperserHostname,
		eigenDADirectoryAddr,
		jsonRPCURL,
		signerAuthKey,
		srsPath,
		certVerifierAddrOverride,
	)
	if err != nil {
		return fmt.Errorf("initialize payload disperser: %w", err)
	}
	defer func() {
		if closeErr := payloadDisperser.Close(); closeErr != nil {
			logger.Error("Failed to close payload disperser", "error", closeErr)
		}
	}()

	// Create an arbitrary payload to disperse
	arbitraryData := []byte("This is a test payload for EigenDA cert verification")
	payload := coretypes.Payload(arbitraryData)

	logger.Info("Dispersing payload", "size", len(arbitraryData))

	// Disperse the payload and get back the cert
	cert, err := payloadDisperser.SendPayload(ctx, payload)
	if err != nil {
		return fmt.Errorf("disperse payload: %w", err)
	}

	logger.Info("Payload dispersed successfully")

	// The cert has already been verified via checkDACert inside SendPayload,
	// but let's verify it again explicitly to demonstrate the verification
	certVerifier, err := createCertVerifier(certVerifierAddr, ethClient, logger)
	if err != nil {
		return fmt.Errorf("create cert verifier: %w", err)
	}

	fmt.Println("CertVerifier tests:")

	verifyCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = certVerifier.CheckDACert(verifyCtx, cert)
	if err != nil {
		fmt.Println(fmt.Errorf("checkDACert call failed with an error: %w", err))
	}

	fmt.Println("checkDACert call passed with a valid DA Cert! ✓")

	v3Cert, ok := cert.(*coretypes.EigenDACertV3)
	if !ok {
		return fmt.Errorf("could not cast to V3 cert")
	}

	// modify the merkle root of the batch header and ensure verification fails
	v3Cert.BatchHeader.BatchRoot = gethcommon.Hash{0x1, 0x2, 0x3, 0x4}

	err = certVerifier.CheckDACert(verifyCtx, v3Cert)
	var errInvalidCert *verification.CertVerifierInvalidCertError
	if err == nil {
		fmt.Println(fmt.Errorf("checkDACert call passed but should have failed when given invalid DA Cert"))
	} else if !errors.As(err, &errInvalidCert) {
		fmt.Println(fmt.Errorf("checkDACert call failed with unknown error: %w", err))
	} else {
		fmt.Println("checkDACert call failed with a non-revertable error as expected when given invalid DA Cert! ✓")
	}

	// Print certificate details
	blobKey, err := cert.ComputeBlobKey()
	if err != nil {
		return fmt.Errorf("compute blob key: %w", err)
	}

	// rbn=0 is fine since this uses static provider
	version, err := certVerifier.GetCertVersion(ctx, 0)
	if err != nil {
		return fmt.Errorf("get cert version: %w", err)
	}

	fmt.Println("========================================================")
	fmt.Printf("Cert version: %d\n", version)
	fmt.Printf("Blob key: %s\n", blobKey.Hex())
	fmt.Printf("Reference Block Number: %d\n", cert.ReferenceBlockNumber())
	fmt.Printf("Quorum Numbers: %v\n", cert.QuorumNumbers())

	return nil
}

func initializePayloadDisperser(
	ctx context.Context,
	logger logging.Logger,
	disperserHostname string,
	eigenDADirectoryAddr gethcommon.Address,
	jsonRPCURL string,
	signerAuthKey string,
	srsPath string,
	certVerifierAddrOverride *gethcommon.Address,
) (*dispersal.PayloadDisperser, *geth.EthClient, gethcommon.Address, error) {

	// Create KZG committer
	kzgCommitter, err := createKzgCommitter(srsPath)
	if err != nil {
		return nil, nil, gethcommon.Address{}, fmt.Errorf("create kzg committer: %w", err)
	}

	// Create Ethereum client
	ethClient, err := createEthClient(logger, jsonRPCURL)
	if err != nil {
		return nil, nil, gethcommon.Address{}, fmt.Errorf("create eth client: %w", err)
	}

	// Create contract directory to fetch addresses
	contractDirectory, err := directory.NewContractDirectory(ctx, logger, ethClient, eigenDADirectoryAddr)
	if err != nil {
		return nil, nil, gethcommon.Address{}, fmt.Errorf("create contract directory: %w", err)
	}

	// Fetch cert verifier address - use override if provided, otherwise fetch from directory
	var certVerifierAddr gethcommon.Address
	if certVerifierAddrOverride != nil {
		certVerifierAddr = *certVerifierAddrOverride
		logger.Info("Using cert verifier address override", "certVerifier", certVerifierAddr.Hex())
	} else {
		certVerifierAddr, err = contractDirectory.GetContractAddress(ctx, directory.CertVerifierRouter)
		if err != nil {
			return nil, nil, gethcommon.Address{}, fmt.Errorf("get cert verifier address: %w", err)
		}
		logger.Info("Fetched cert verifier address from directory", "certVerifier", certVerifierAddr.Hex())
	}

	// Fetch remaining contract addresses from the directory
	operatorStateRetrieverAddr, err := contractDirectory.GetContractAddress(ctx, directory.OperatorStateRetriever)
	if err != nil {
		return nil, nil, gethcommon.Address{}, fmt.Errorf("get operator state retriever address: %w", err)
	}

	registryCoordinatorAddr, err := contractDirectory.GetContractAddress(ctx, directory.RegistryCoordinator)
	if err != nil {
		return nil, nil, gethcommon.Address{}, fmt.Errorf("get registry coordinator address: %w", err)
	}

	logger.Info("Contract addresses configured",
		"certVerifier", certVerifierAddr.Hex(),
		"operatorStateRetriever", operatorStateRetrieverAddr.Hex(),
		"registryCoordinator", registryCoordinatorAddr.Hex())

	// Create cert verifier using static address provider
	certVerifier, err := createCertVerifier(certVerifierAddr, ethClient, logger)
	if err != nil {
		return nil, nil, gethcommon.Address{}, fmt.Errorf("create cert verifier: %w", err)
	}

	// Create cert builder
	certBuilder, err := clients.NewCertBuilder(
		logger,
		operatorStateRetrieverAddr,
		registryCoordinatorAddr,
		ethClient,
	)
	if err != nil {
		return nil, nil, gethcommon.Address{}, fmt.Errorf("new cert builder: %w", err)
	}

	// Create block number monitor
	blockNumMonitor, err := verification.NewBlockNumberMonitor(
		logger,
		ethClient,
		1*time.Second,
	)
	if err != nil {
		return nil, nil, gethcommon.Address{}, fmt.Errorf("create block number monitor: %w", err)
	}

	// Configure payload disperser
	payloadDisperserConfig := dispersal.PayloadDisperserConfig{
		PayloadClientConfig:    *clients.GetDefaultPayloadClientConfig(),
		DisperseBlobTimeout:    60 * time.Second,
		BlobCompleteTimeout:    120 * time.Second,
		BlobStatusPollInterval: 2 * time.Second,
		ContractCallTimeout:    10 * time.Second,
	}

	disperserClientMultiplexer, err := createDisperserClientMultiplexer(
		logger, disperserHostname, signerAuthKey, kzgCommitter)
	if err != nil {
		return nil, nil, gethcommon.Address{}, fmt.Errorf("create disperser client multiplexer: %w", err)
	}

	// Create payload disperser (without client ledger for simplicity - legacy payment mode)
	payloadDisperser, err := dispersal.NewPayloadDisperser(
		logger,
		payloadDisperserConfig,
		disperserClientMultiplexer,
		blockNumMonitor,
		certBuilder,
		certVerifier,
		nil, // clientLedger - nil for legacy payment mode
		nil, // registry - nil for no metrics
	)
	if err != nil {
		return nil, nil, gethcommon.Address{}, fmt.Errorf("new payload disperser: %w", err)
	}

	return payloadDisperser, ethClient, certVerifierAddr, nil
}

func createDisperserClientMultiplexer(
	logger logging.Logger,
	disperserHostName string,
	privateKey string,
	kzgCommitter *committer.Committer,
) (*dispersal.DisperserClientMultiplexer, error) {
	signer, err := auth.NewLocalBlobRequestSigner(privateKey)
	if err != nil {
		return nil, fmt.Errorf("create blob request signer: %w", err)
	}

	hostname, portStr, err := net.SplitHostPort(disperserHostName)
	if err != nil {
		return nil, fmt.Errorf("parse disperser host: %w", err)
	}

	portUint64, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("parse disperser port: %w", err)
	}

	multiplexerConfig := dispersal.DefaultDisperserClientMultiplexerConfig()
	connectionInfo := &clients.DisperserConnectionInfo{
		Hostname: hostname,
		Port:     uint16(portUint64),
	}
	disperserRegistry := clients.NewLegacyDisperserRegistry(connectionInfo)

	return dispersal.NewDisperserClientMultiplexer(
		logger,
		multiplexerConfig,
		disperserRegistry,
		signer,
		kzgCommitter,
		metrics.NoopDispersalMetrics,
		8,
	), nil
}

func createEthClient(logger logging.Logger, rpcURL string) (*geth.EthClient, error) {
	ethClientConfig := geth.EthClientConfig{
		RPCURLs:          []string{rpcURL},
		NumConfirmations: 0,
		NumRetries:       3,
	}

	client, err := geth.NewClient(
		ethClientConfig,
		gethcommon.Address{},
		0,
		logger)
	if err != nil {
		return nil, fmt.Errorf("new eth client: %w", err)
	}
	return client, nil
}

func createCertVerifier(
	certVerifierAddress gethcommon.Address,
	ethClient common.EthClient,
	logger logging.Logger,
) (*verification.CertVerifier, error) {
	// Use static address provider since we're given a specific cert verifier address
	addressProvider := verification.NewStaticCertVerifierAddressProvider(certVerifierAddress)

	verifier, err := verification.NewCertVerifier(logger, ethClient, addressProvider)
	if err != nil {
		return nil, fmt.Errorf("new cert verifier: %w", err)
	}
	return verifier, nil
}

func createKzgCommitter(srsPath string) (*committer.Committer, error) {
	config := committer.Config{
		G1SRSPath:         srsPath + "/g1.point",
		G2SRSPath:         srsPath + "/g2.point",
		G2TrailingSRSPath: srsPath + "/g2.trailing.point",
		SRSNumberToLoad:   8192 / 32, // 8192 / encoding.BYTES_PER_SYMBOL
	}

	committer, err := committer.NewFromConfig(config)
	if err != nil {
		return nil, fmt.Errorf("new kzg committer from config: %w", err)
	}
	return committer, nil
}

func createLogger() (logging.Logger, error) {
	config := common.DefaultTextLoggerConfig()
	logger, err := common.NewLogger(config)
	if err != nil {
		return nil, fmt.Errorf("create new logger: %w", err)
	}

	return logger, nil
}
