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
	clientLock       sync.Mutex
	client           *TestClient
)

// GetTestClient returns a TestClient instance, creating one if it does not exist.
// This uses a global static client... this is icky, but it takes ~1 minute
// to read the SRS points, so it's the lesser of two evils to keep it around.
func GetTestClient(t *testing.T) *TestClient {
	clientLock.Lock()
	defer clientLock.Unlock()

	skipInCI(t)
	if client != nil {
		return client
	}

	testConfig, err := getClientConfig()
	require.NoError(t, err)
	client, err = NewTestClient(testConfig)
	require.NoError(t, err)

	err = setupFilesystem(testConfig)
	require.NoError(t, err)

	return client
}

func skipInCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping test in CI environment")
	}
}

// GetClient is similar to GetTestClient, but is intended to be used outside of a unit test.
// Do not call this method in a unit test, as GetTestClient does things like skip the test
// if it is running in a CI environment.
func GetClient() (*TestClient, error) {
	clientLock.Lock()
	defer clientLock.Unlock()

	if client != nil {
		return client, nil
	}

	testConfig, err := getClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get client config: %v", err)
	}
	client, err = NewTestClient(testConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create test client: %v", err)
	}

	err = setupFilesystem(testConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to setup filesystem: %v", err)
	}

	return client, nil
}

// getClientConfig parses and returns the test client configuration.
func getClientConfig() (*TestClientConfig, error) {

	configFile, err := ResolveTildeInPath(targetConfigFile)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve config file: %v", err)
	}
	configFileBytes, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	config := &TestClientConfig{}
	err = json.Unmarshal(configFileBytes, config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %v", err)
	}

	return config, nil
}

func setupFilesystem(config *TestClientConfig) error {
	// Create the test data directory if it does not exist
	err := os.MkdirAll(config.TestDataPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create test data directory: %v", err)
	}

	// Create the SRS directories if they do not exist
	srsPath, err := config.path(SRSPath)
	if err != nil {
		return fmt.Errorf("failed to resolve SRS path: %v", err)
	}
	err = os.MkdirAll(srsPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create SRS directory: %v", err)
	}
	srsTablesPath, err := config.path(SRSPathSRSTables)
	if err != nil {
		return fmt.Errorf("failed to resolve SRS tables path: %v", err)
	}
	err = os.MkdirAll(srsTablesPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create SRS tables directory: %v", err)
	}

	// If any of the srs files do not exist, download them.
	filePath, err := config.path(SRSPathG1)
	if err != nil {
		return fmt.Errorf("failed to resolve SRS G1 path: %v", err)
	}
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
		if err != nil {
			return fmt.Errorf("failed to download G1 point: %v", err)
		}
	} else {
		if err != nil {
			return fmt.Errorf("failed to stat G1 point: %v", err)
		}
	}

	filePath, err = config.path(SRSPathG2)
	if err != nil {
		return fmt.Errorf("failed to resolve SRS G2 path: %v", err)
	}
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
		if err != nil {
			return fmt.Errorf("failed to download G2 point: %v", err)
		}
	} else {
		if err != nil {
			return fmt.Errorf("failed to stat G2 point: %v", err)
		}
	}

	filePath, err = config.path(SRSPathG2PowerOf2)
	if err != nil {
		return fmt.Errorf("failed to resolve SRS G2 power of 2 path: %v", err)
	}
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
		if err != nil {
			return fmt.Errorf("failed to download G2 power of 2 point: %v", err)
		}
	} else {
		if err != nil {
			return fmt.Errorf("failed to stat G2 power of 2 point: %v", err)
		}
	}

	// Check to see if the private key file exists. If not, stop the test.
	filePath, err = ResolveTildeInPath(config.KeyPath)
	if err != nil {
		return fmt.Errorf("failed to resolve key path: %v", err)
	}
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("private key file %s does not exist. This file should "+
			"contain the private key for the account used in the test, in hex.", filePath)
	}
	return nil
}
