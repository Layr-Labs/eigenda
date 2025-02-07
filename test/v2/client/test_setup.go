package client

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/stretchr/testify/require"
)

var (
	targetConfigFile = "../config/environment/preprod.json"
	configLock       sync.Mutex
	config           *TestClientConfig
	clientLock       sync.Mutex
	clientMap        = make(map[string]*TestClient)
	logger           logging.Logger
	metrics          *testClientMetrics
)

// GetConfig returns a TestClientConfig instance, creating one if it does not exist.
func GetConfig(t *testing.T) *TestClientConfig {
	configLock.Lock()
	defer configLock.Unlock()

	skipInCI(t)
	if config != nil {
		return config
	}

	configFile := resolveTildeInPath(t, targetConfigFile)
	configFileBytes, err := os.ReadFile(configFile)
	require.NoError(t, err)

	config = &TestClientConfig{}
	err = json.Unmarshal(configFileBytes, config)
	require.NoError(t, err)

	return config
}

// GetClient returns a TestClient instance, creating one if it does not exist.
// This uses a global static client... this is icky, but it takes ~1 minute
// to read the SRS points, so it's the lesser of two evils to keep it around.
func GetClient(t *testing.T, quorums []core.QuorumID) *TestClient {
	clientLock.Lock()
	defer clientLock.Unlock()

	skipInCI(t)

	testConfig := GetConfig(t)

	quorumsString := ""
	for _, quorum := range quorums {
		quorumsString += string(quorum) + ","
	}
	if clientMap[quorumsString] != nil {
		return clientMap[quorumsString]
	}

	if len(clientMap) == 0 {
		// only do this stuff once
		setupFilesystem(t, testConfig)

		var loggerConfig common.LoggerConfig
		if os.Getenv("CI") != "" {
			loggerConfig = common.DefaultLoggerConfig()
		} else {
			loggerConfig = common.DefaultConsoleLoggerConfig()
		}

		testLogger, err := common.NewLogger(loggerConfig)
		require.NoError(t, err)

		logger = testLogger

		testMetrics := newTestClientMetrics(logger, config.MetricsPort)
		metrics = testMetrics
		testMetrics.start()
	}

	client := NewTestClient(t, logger, metrics, testConfig, quorums)

	clientMap[quorumsString] = client

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
	filePath = resolveTildeInPath(t, config.KeyPath)
	_, err = os.Stat(filePath)
	require.NoError(t, err,
		"private key file %s does not exist. This file should "+
			"contain the private key for the account used in the test, in hex.",
		filePath)
}
