package v2

import (
	"fmt"
	"github.com/docker/go-units"
	"github.com/stretchr/testify/require"
	"os"
	"os/exec"
	"sync"
	"testing"
)

var (
	preprodConfig = &TestClientConfig{
		TestDataPath:                  "~/.test-v2",
		DisperserHostname:             "disperser-preprod-holesky.eigenda.xyz",
		DisperserPort:                 443,
		EthRPCURLs:                    []string{"https://ethereum-holesky-rpc.publicnode.com"},
		BLSOperatorStateRetrieverAddr: "0x93545e3b9013CcaBc31E80898fef7569a4024C0C",
		EigenDAServiceManagerAddr:     "0x54A03db2784E3D0aCC08344D05385d0b62d4F432",
		EigenDACertVerifierAddress:    "0xe2C7AfB3c47B800b439b0a3d8EA40ca79759B245",
		SubgraphURL:                   "https://subgraph.satsuma-prod.com/51caed8fa9cb/eigenlabs/eigenda-operator-state-preprod-holesky/version/v0.7.0/api",
		SRSOrder:                      268435456,
		MaxBlobSize:                   16 * units.MiB,
		MinimumSigningPercent:         55,
	}

	lock   sync.Mutex
	client *TestClient

	targetConfig = preprodConfig
)

// getClient returns a TestClient instance, creating one if it does not exist.
// This uses a global static client... this is icky, but it takes ~1 minute
// to read the SRS points, so it's the lesser of two evils to keep it around.
func getClient(t *testing.T) *TestClient {
	lock.Lock()
	defer lock.Unlock()

	skipInCI(t)
	setupFilesystem(t, targetConfig)

	if client == nil {
		client = NewTestClient(t, targetConfig)
	}

	return client
}

func skipInCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping test in CI environment")
	}
}

func setupFilesystem(t *testing.T, config *TestClientConfig) {
	// Create the test data directory if it does not exist
	err := os.MkdirAll(config.TestDataPath, 0755)
	require.NoError(t, err)

	// Create the SRS directories if they do not exist
	err = os.MkdirAll(config.path(t, SRSPath), 0755)
	require.NoError(t, err)
	err = os.MkdirAll(config.path(t, SRSPathSRSTables), 0755)
	require.NoError(t, err)

	// If any of the srs files do not exist, download them.
	filePath := config.path(t, SRSPathG1)
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		command := make([]string, 3)
		command[0] = "wget"
		command[1] = "https://srs-mainnet.s3.amazonaws.com/kzg/g1.point"
		command[2] = "--output-document=" + filePath
		fmt.Printf("executing %s\n", command)

		cmd := exec.Command(command[0], command[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		require.NoError(t, err)
	} else {
		require.NoError(t, err)
	}

	filePath = config.path(t, SRSPathG2)
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		command := make([]string, 3)
		command[0] = "wget"
		command[1] = "https://srs-mainnet.s3.amazonaws.com/kzg/g2.point"
		command[2] = "--output-document=" + filePath
		fmt.Printf("executing %s\n", command)

		cmd := exec.Command(command[0], command[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		require.NoError(t, err)
	} else {
		require.NoError(t, err)
	}

	filePath = config.path(t, SRSPathG2PowerOf2)
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		command := make([]string, 3)
		command[0] = "wget"
		command[1] = "https://srs-mainnet.s3.amazonaws.com/kzg/g2.point.powerOf2"
		command[2] = "--output-document=" + filePath
		fmt.Printf("executing %s\n", command)

		cmd := exec.Command(command[0], command[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		require.NoError(t, err)
	} else {
		require.NoError(t, err)
	}

	// Check to see if the private key file exists. If not, stop the test.
	filePath = config.path(t, KeyPath)
	_, err = os.Stat(filePath)
	require.NoError(t, err,
		"private key file %s does not exist. This file should "+
			"contain the private key for the account used in the test, in hex.",
		filePath)
}
