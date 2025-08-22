package deployment

import (
	"fmt"
	"math/big"
	"time"

	clientsv2 "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/metrics"
	"github.com/Layr-Labs/eigenda/api/clients/v2/payloaddispersal"
	"github.com/Layr-Labs/eigenda/api/clients/v2/verification"
	"github.com/Layr-Labs/eigenda/common"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
)

// CreateTransactorOpts creates transaction options from a private key
func CreateTransactorOpts(privateKeyHex string, chainID *big.Int) (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	opts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	return opts, nil
}

// PayloadDisperserParams contains parameters for setting up a payload disperser
type PayloadDisperserParams struct {
	Logger               logging.Logger
	EthClient            common.EthClient
	DisperserPrivateKey  string
	DisperserHostname    string
	DisperserPort        string
	DisperseBlobTimeout  time.Duration
	BlobCompleteTimeout  time.Duration
	BlobStatusPollInterval time.Duration
	ContractCallTimeout  time.Duration
	CertBuilder          *clientsv2.CertBuilder
	RouterCertVerifier   *verification.CertVerifier
}

// SetupPayloadDisperser sets up a payload disperser with the given configuration
func SetupPayloadDisperser(params PayloadDisperserParams) (*payloaddispersal.PayloadDisperser, error) {
	// Set up the block monitor
	blockMonitor, err := verification.NewBlockNumberMonitor(params.Logger, params.EthClient, time.Second*1)
	if err != nil {
		return nil, fmt.Errorf("failed to create block monitor: %w", err)
	}

	// Set up the signer
	signer, err := auth.NewLocalBlobRequestSigner(params.DisperserPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create signer: %w", err)
	}

	// Create disperser client config
	disperserClientConfig := &clientsv2.DisperserClientConfig{
		Hostname: params.DisperserHostname,
		Port:     params.DisperserPort,
	}

	// Get account ID
	accountId, err := signer.GetAccountID()
	if err != nil {
		return nil, fmt.Errorf("failed to get account ID: %w", err)
	}

	// Create accountant
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

	// Create disperser client
	disperserClient, err := clientsv2.NewDisperserClient(disperserClientConfig, signer, nil, accountant)
	if err != nil {
		return nil, fmt.Errorf("failed to create disperser client: %w", err)
	}

	// Create payload disperser config
	payloadDisperserConfig := payloaddispersal.PayloadDisperserConfig{
		PayloadClientConfig:    *clientsv2.GetDefaultPayloadClientConfig(),
		DisperseBlobTimeout:    params.DisperseBlobTimeout,
		BlobCompleteTimeout:    params.BlobCompleteTimeout,
		BlobStatusPollInterval: params.BlobStatusPollInterval,
		ContractCallTimeout:    params.ContractCallTimeout,
	}

	// Create payload disperser
	payloadDisperser, err := payloaddispersal.NewPayloadDisperser(
		params.Logger,
		payloadDisperserConfig,
		disperserClient,
		blockMonitor,
		params.CertBuilder,
		params.RouterCertVerifier,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create payload disperser: %w", err)
	}

	return payloadDisperser, nil
}