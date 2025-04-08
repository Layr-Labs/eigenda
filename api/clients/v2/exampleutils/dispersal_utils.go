package exampleutils

import (
	"crypto/rand"
	"io"
	"path/filepath"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloaddispersal"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// CreatePayloadDisperser creates a PayloadDisperser with necessary values configured for a basic test
func CreatePayloadDisperser(privateKey string) (*payloaddispersal.PayloadDisperser, error) {
	// Create a logger with a null output for examples to avoid polluting test output
	config := common.DefaultLoggerConfig()
	config.OutputWriter = io.Discard // Send logs to /dev/null
	logger, err := common.NewLogger(config)
	if err != nil {
		return nil, err
	}

	disperserClient, err := createDisperserClient(privateKey)
	if err != nil {
		return nil, err
	}

	certVerifier, err := createCertVerifier(privateKey, logger)
	if err != nil {
		return nil, err
	}

	return payloaddispersal.NewPayloadDisperser(
		logger,
		createPayloadDisperserConfig(),
		disperserClient,
		certVerifier,
		nil,
	)
}

// CreateRandomPayload creates a payload with random data of the specified size
func CreateRandomPayload(size int) (*coretypes.Payload, error) {
	payloadBytes := make([]byte, size)
	_, err := rand.Read(payloadBytes)
	if err != nil {
		return nil, err
	}
	return coretypes.NewPayload(payloadBytes), nil
}

func createKZGConfig() *kzg.KzgConfig {
	srsPath := "../../../../resources/srs"

	return &kzg.KzgConfig{
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

func createDisperserClientConfig() *clients.DisperserClientConfig {
	return &clients.DisperserClientConfig{
		Hostname:          "disperser-testnet-holesky.eigenda.xyz",
		Port:              "443",
		UseSecureGrpcFlag: true,
	}
}

func createEthClientConfig(privateKey string) geth.EthClientConfig {
	return geth.EthClientConfig{
		RPCURLs:          []string{"https://ethereum-holesky-rpc.publicnode.com"},
		PrivateKeyString: privateKey,
		NumConfirmations: 0,
		NumRetries:       3,
	}
}

func createPayloadDisperserConfig() payloaddispersal.PayloadDisperserConfig {
	payloadClientConfig := clients.GetDefaultPayloadClientConfig()
	return payloaddispersal.PayloadDisperserConfig{
		PayloadClientConfig:    *payloadClientConfig,
		DisperseBlobTimeout:    2 * time.Minute,
		BlobCertifiedTimeout:   2 * time.Minute,
		BlobStatusPollInterval: 1 * time.Second,
		ContractCallTimeout:    5 * time.Second,
	}
}

func createDisperserClient(privateKey string) (clients.DisperserClient, error) {
	signer, err := auth.NewLocalBlobRequestSigner(privateKey)
	if err != nil {
		return nil, err
	}

	kzgProver, err := prover.NewProver(createKZGConfig(), nil)
	if err != nil {
		return nil, err
	}

	// Create and return disperser client
	return clients.NewDisperserClient(
		createDisperserClientConfig(),
		signer,
		kzgProver,
		nil)
}

func createEthClient(privateKey string, logger logging.Logger) (*geth.EthClient, error) {
	signer, err := auth.NewLocalBlobRequestSigner(privateKey)
	if err != nil {
		return nil, err
	}

	accountID, err := signer.GetAccountID()
	if err != nil {
		return nil, err
	}

	return geth.NewClient(
		createEthClientConfig(privateKey),
		accountID,
		0,
		logger)
}

// createCertVerifier creates a certificate verifier
func createCertVerifier(privateKey string, logger logging.Logger) (*verification.CertVerifier, error) {
	ethClient, err := createEthClient(privateKey, logger)
	if err != nil {
		return nil, err
	}

	return verification.NewCertVerifier(
		logger,
		ethClient,
		verification.NewStaticCertVerifierAddressProvider(
			gethcommon.HexToAddress("0xFe52fE1940858DCb6e12153E2104aD0fDFbE1162")))
}
