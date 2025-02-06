package client

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"os/exec"
	"sync"
	"testing"
)

var (
	targetConfigFile = "../config/environment/preprod.json"
	configLock       sync.Mutex
	config           *TestClientConfig
	clientLock       sync.Mutex
	client           *TestClient
)

// GetConfig returns a TestClientConfig instance, creating one if it does not exist.
func GetConfig(t *testing.T) *TestClientConfig {
	configLock.Lock()
	defer configLock.Unlock()

	skipInCI(t)
	if config != nil {
		return config
	}

	configFile, err := ResolveTildeInPath(targetConfigFile)
	require.NoError(t, err)
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
func GetClient(t *testing.T) *TestClient {
	clientLock.Lock()
	defer clientLock.Unlock()

	skipInCI(t) // TODO create non-test variant
	if client != nil {
		return client
	}

	testConfig := GetConfig(t)
	var err error
	client, err = NewTestClient(testConfig)
	require.NoError(t, err)

	setupFilesystem(t, testConfig)

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
	srsPath, err := config.path(SRSPath)
	require.NoError(t, err)
	err = os.MkdirAll(srsPath, 0755)
	require.NoError(t, err)
	srsTablesPath, err := config.path(SRSPathSRSTables)
	require.NoError(t, err)
	err = os.MkdirAll(srsTablesPath, 0755)
	require.NoError(t, err)

	// If any of the srs files do not exist, download them.
	filePath, err := config.path(SRSPathG1)
	require.NoError(t, err)
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

	filePath, err = config.path(SRSPathG2)
	require.NoError(t, err)
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

	filePath, err = config.path(SRSPathG2PowerOf2)
	require.NoError(t, err)
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
	filePath, err = ResolveTildeInPath(config.KeyPath)
	require.NoError(t, err)
	_, err = os.Stat(filePath)
	require.NoError(t, err,
		"private key file %s does not exist. This file should "+
			"contain the private key for the account used in the test, in hex.",
		filePath)
}
