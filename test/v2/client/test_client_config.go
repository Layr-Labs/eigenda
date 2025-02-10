package client

import (
	"fmt"
	"path"
)

// TestClientConfig is the configuration for the test client.
type TestClientConfig struct {
	// The location where persistent test data is stored (e.g. SRS files). Often private keys are stored here too.
	TestDataPath string
	// The location where the test client's private key is stored.
	// This is the key for the account that is paying for dispersals.
	KeyPath string
	// The disperser's hostname (url or IP address)
	DisperserHostname string
	// The disperser's port
	DisperserPort int
	// The URL(s) to point the eth client to
	EthRPCURLs []string
	// The contract address for the EigenDA BLS operator state retriever
	BLSOperatorStateRetrieverAddr string
	// The contract address for the EigenDA service manager
	EigenDAServiceManagerAddr string
	// The contract address for the EigenDA cert verifier
	EigenDACertVerifierAddress string
	// The URL/IP of a subgraph to use for the chain state
	SubgraphURL string
	// The SRS order to use for the test
	SRSOrder uint64
	// The SRS number to load, increasing this beyond necessary can cause the client to take a long time to start
	SRSNumberToLoad uint64
	// The maximum blob size supported by the EigenDA network
	MaxBlobSize uint64
	// Required signing percentage for a quorum to be considered valid, out of 100
	MinimumSigningPercent int
	// The port to use for metrics (if metrics are being collected)
	MetricsPort int
}

// Path returns the full path to a file in the test data directory.
func (c *TestClientConfig) Path(elements ...string) (string, error) {
	root, err := ResolveTildeInPath(c.TestDataPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve tilde in path: %w", err)
	}

	combinedElements := make([]string, 0, len(elements)+1)
	combinedElements = append(combinedElements, root)
	combinedElements = append(combinedElements, elements...)

	return path.Join(combinedElements...), nil
}
