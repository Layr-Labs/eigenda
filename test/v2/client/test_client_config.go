package client

import (
	"fmt"
	"path"

	"github.com/Layr-Labs/eigenda/litt/util"
)

// TestClientConfig is the configuration for the test client.
type TestClientConfig struct {
	// The location where the SRS files can be found.
	SRSPath string
	// The location where the test client's private key is stored. This is the key for the account that is
	// paying for dispersals.
	//
	// Either this or KeyVar must be set. If both are set, KeyPath is used.
	KeyPath string
	// The environment variable that contains the private key for the account that is paying for dispersals.
	//
	// This is used if KeyPath is not set.
	KeyVar string
	// The disperser's hostname (url or IP address)
	DisperserHostname string
	// The disperser's port
	DisperserPort int
	// The URL(s) to point the eth client to
	//
	// Either this or EthRPCURLsVar must be set. If both are set, EthRPCURLs is used.
	EthRPCURLs []string
	// The environment variable that contains the URL(s) to point the eth client to. Use a comma-separated list.
	//
	// Either this or EthRPCURLs must be set. If both are set, EthRPCURLs is used.
	EthRPCUrlsVar string
	// The contract address for the EigenDA address directory, where all contract addresses are stored
	//
	// Currently the EigenDA address directory is just used to look up BLSOperatorStateRetrieverAddr and EigenDAServiceManagerAddr.
	// In a later PR, ensure all addresses are populated into the directory, and use it for all contract address lookups.
	EigenDADirectory string
	// The contract address for the EigenDA BLS operator state retriever
	BLSOperatorStateRetrieverAddr string
	// The contract address for the EigenDA service manager
	EigenDAServiceManagerAddr string
	// The contract address for the EigenDA cert verifier, which specifies required quorums 0 and 1
	//
	// If this value is not set, that tests utilizing it will be skipped
	EigenDACertVerifierAddressQuorums0_1 string
	// The contract address for the EigenDA cert verifier, which specifies required quorums 0, 1, and 2
	//
	// If this value is not set, that tests utilizing it will be skipped
	EigenDACertVerifierAddressQuorums0_1_2 string
	// The contract address for the EigenDA cert verifier, which specifies required quorum 2
	//
	// If this value is not set, that tests utilizing it will be skipped
	EigenDACertVerifierAddressQuorums2 string
	// The URL/IP of a subgraph to use for the chain state
	SubgraphURL string
	// The SRS order to use for the test
	SRSOrder uint64
	// The SRS number to load, increasing this beyond necessary can cause the client to take a long time to start
	SRSNumberToLoad uint64
	// The maximum blob size supported by the EigenDA network
	MaxBlobSize uint64
	// The port to use for metrics (if metrics are being collected)
	MetricsPort int
	// If true, do not start the metrics server.
	DisableMetrics bool
	// The size of the thread pool for read operations on the relay.
	ValidatorReadConnectionPoolSize int
	// The size of the thread pool for CPU heavy operations.
	ValidatorReadComputePoolSize int
	// The port to use for the proxy.
	ProxyPort int
}

// TODO create default config

// ResolveSRSPath returns a path relative to the SRSPath root directory.
func (c *TestClientConfig) ResolveSRSPath(srsFile string) (string, error) {
	root, err := util.SanitizePath(c.SRSPath)
	if err != nil {
		return "", fmt.Errorf("failed to sanitize path: %w", err)
	}
	return path.Join(root, srsFile), nil
}
