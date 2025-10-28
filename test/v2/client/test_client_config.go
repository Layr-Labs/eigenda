package client

import (
	"fmt"
	"path"

	"github.com/Layr-Labs/eigenda/common/config"
	"github.com/Layr-Labs/eigenda/core/payments/clientledger"
	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/docker/go-units"
)

var _ config.VerifiableConfig = (*TestClientConfig)(nil)

// TestClientConfig is the configuration for the test client.
type TestClientConfig struct {
	// The location where the SRS files can be found.
	SrsPath string `docs:"required"`
	// The private key for the account that is paying for dispersals, in hex format (0x...)
	PrivateKey string `docs:"required"`
	// The disperser's hostname (url or IP address)
	DisperserHostname string `docs:"required"`
	// The disperser's port
	DisperserPort int `docs:"required"`
	// The URL(s) to point the eth client to
	//
	// Either this or EthRpcUrlsVar must be set. If both are set, EthRpcUrls is used.
	EthRpcUrls []string `docs:"required"`
	// The contract address for the EigenDA address directory, where all contract addresses are stored
	ContractDirectoryAddress string `docs:"required"`
	// The URL/IP of a subgraph to use for the chain state
	SubgraphUrl string `docs:"required"`
	// The SRS order to use for the test
	SrsOrder uint64
	// The SRS number to load, increasing this beyond necessary can cause the client to take a long time to start
	SRSNumberToLoad uint64
	// The maximum blob size supported by the EigenDA network
	MaxBlobSize uint64
	// The port to use for metrics (if metrics are being collected)
	MetricsPort int
	// If true, do not start the metrics server.
	DisableMetrics bool
	// The size of the thread pool for read operations.
	ValidatorReadConnectionPoolSize int
	// The size of the thread pool for CPU heavy operations.
	ValidatorReadComputePoolSize int
	// The number of connections to open for each relay.
	RelayConnectionCount uint
	// The number of connections to open for each disperser.
	DisperserConnectionCount uint
	// The port to use for the proxy.
	ProxyPort int
	// Client ledger mode used for payments.
	ClientLedgerPaymentMode string
}

// DefaultTestClientConfig returns a default configuration for the test client. Sets default values for fields
// where default values make sense.
func DefaultTestClientConfig() *TestClientConfig {
	return &TestClientConfig{
		DisperserPort:                   443,
		MaxBlobSize:                     16 * units.MiB,
		SrsOrder:                        268435456,
		MetricsPort:                     9101,
		ValidatorReadConnectionPoolSize: 100,
		ValidatorReadComputePoolSize:    20,
		ProxyPort:                       1234,
		RelayConnectionCount:            8,
		DisperserConnectionCount:        8,
		ClientLedgerPaymentMode:         string(clientledger.ClientLedgerModeLegacy),
	}
}

// ResolveSRSPath returns a path relative to the SRSPath root directory.
func (c *TestClientConfig) ResolveSRSPath(srsFile string) (string, error) {
	root, err := util.SanitizePath(c.SrsPath)
	if err != nil {
		return "", fmt.Errorf("failed to sanitize path: %w", err)
	}
	return path.Join(root, srsFile), nil
}

// Verify implements config.VerifiableConfig.
func (c *TestClientConfig) Verify() error {
	if c.SrsPath == "" {
		return fmt.Errorf("SrsPath must be set")
	}
	if c.PrivateKey == "" {
		return fmt.Errorf("PrivateKey must be set")
	}
	if c.DisperserHostname == "" {
		return fmt.Errorf("DisperserHostname must be set")
	}
	if c.DisperserPort <= 0 || c.DisperserPort > 65535 {
		return fmt.Errorf("DisperserPort must be a valid port number")
	}
	if c.EthRpcUrls == nil || len(c.EthRpcUrls) == 0 {
		return fmt.Errorf("EthRpcUrls must be set and contain at least one URL")
	}
	if c.ContractDirectoryAddress == "" {
		return fmt.Errorf("ContractDirectoryAddress must be set")
	}
	if c.SubgraphUrl == "" {
		return fmt.Errorf("SubgraphUrl must be set")
	}
	if c.SrsOrder == 0 {
		return fmt.Errorf("SrsOrder must be set and greater than 0")
	}
	if c.MaxBlobSize == 0 {
		return fmt.Errorf("MaxBlobSize must be set and greater than 0")
	}
	if c.ValidatorReadConnectionPoolSize <= 0 {
		return fmt.Errorf("ValidatorReadConnectionPoolSize must be set and greater than 0")
	}
	if c.ValidatorReadComputePoolSize <= 0 {
		return fmt.Errorf("ValidatorReadComputePoolSize must be set and greater than 0")
	}
	if c.RelayConnectionCount == 0 {
		return fmt.Errorf("RelayConnectionCount must be set and greater than 0")
	}
	if c.DisperserConnectionCount == 0 {
		return fmt.Errorf("DisperserConnectionCount must be set and greater than 0")
	}
	return nil
}
