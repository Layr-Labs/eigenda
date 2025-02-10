package client

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/stretchr/testify/require"
)

var (
	targetConfigFile = "../config/environment/preprod.json"
	configLock       sync.Mutex
	config           *TestClientConfig
	clientLock       sync.Mutex
	client           *TestClient
	//clientMap        = make(map[string]*TestClient)
	logger  logging.Logger
	metrics *testClientMetrics
)

func SetTargetConfigFile(file string) {
	clientLock.Lock()
	defer clientLock.Unlock()

	targetConfigFile = file
	client.Stop()
	client = nil // TODO
	//clientMap = make(map[string]*TestClient)
}

// GetConfig returns a TestClientConfig instance, creating one if it does not exist.
func GetConfig() (*TestClientConfig, error) {
	configLock.Lock()
	defer configLock.Unlock()

	if config != nil {
		return config, nil
	}

	configFile, err := ResolveTildeInPath(targetConfigFile)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve tilde in path: %w", err)
	}
	configFileBytes, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config = &TestClientConfig{}
	err = json.Unmarshal(configFileBytes, config)
	if err != nil {
		config = nil
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return config, nil
}

// GetTestClient is the same as GetClient, but also performs a check to ensure that the test is not
// running in a CI environment. If using a TestClient in a unit test, it is critical to use this method
// to ensure that the test is not running in a CI environment.
func GetTestClient(t *testing.T) *TestClient {
	skipInCI(t)
	c, err := GetClient()
	require.NoError(t, err)
	return c
}

// GetClient returns a TestClient instance, creating one if it does not exist.
// This uses a global static client... this is icky, but it takes ~1 minute
// to read the SRS points, so it's the lesser of two evils to keep it around.
func GetClient() (*TestClient, error) {
	clientLock.Lock()
	defer clientLock.Unlock()

	if client != nil {
		return client, nil
	}

	testConfig, err := GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	//if len(clientMap) == 0 { // TODO
	var loggerConfig common.LoggerConfig
	if os.Getenv("CI") != "" {
		loggerConfig = common.DefaultLoggerConfig()
	} else {
		loggerConfig = common.DefaultConsoleLoggerConfig()
	}

	testLogger, err := common.NewLogger(loggerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}
	logger = testLogger

	// only do this stuff once
	err = setupFilesystem(logger, testConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to setup filesystem: %w", err)
	}

	testMetrics := newTestClientMetrics(logger, config.MetricsPort)
	metrics = testMetrics
	testMetrics.start()

	client, err = NewTestClient(logger, metrics, testConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create test client: %w", err)
	}

	return client, nil
}

func skipInCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping test in CI environment")
	}
}

func setupFilesystem(logger logging.Logger, config *TestClientConfig) error {
	// Create the test data directory if it does not exist
	err := os.MkdirAll(config.TestDataPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create test data directory: %w", err)
	}

	// Create the SRS directories if they do not exist
	srsPath, err := config.path(SRSPath)
	if err != nil {
		return fmt.Errorf("failed to get SRS path: %w", err)
	}
	err = os.MkdirAll(srsPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create SRS directory: %w", err)
	}
	srsTablesPath, err := config.path(SRSPathSRSTables)
	if err != nil {
		return fmt.Errorf("failed to get SRS tables path: %w", err)
	}
	err = os.MkdirAll(srsTablesPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create SRS tables directory: %w", err)
	}

	// If any of the srs files do not exist, download them.
	filePath, err := config.path(SRSPathG1)
	if err != nil {
		return fmt.Errorf("failed to get SRS G1 path: %w", err)
	}
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		command := make([]string, 3)
		command[0] = "wget"
		command[1] = "https://srs-mainnet.s3.amazonaws.com/kzg/g1.point"
		command[2] = "--output-document=" + filePath
		logger.Info("executing %s", command)

		cmd := exec.Command(command[0], command[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to download G1 point: %w", err)
		}
	} else {
		if err != nil {
			return fmt.Errorf("failed to check if G1 point exists: %w", err)
		}
	}

	filePath, err = config.path(SRSPathG2)
	if err != nil {
		return fmt.Errorf("failed to get SRS G2 path: %w", err)
	}
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		command := make([]string, 3)
		command[0] = "wget"
		command[1] = "https://srs-mainnet.s3.amazonaws.com/kzg/g2.point"
		command[2] = "--output-document=" + filePath
		logger.Info("executing %s", command)

		cmd := exec.Command(command[0], command[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to download G2 point: %w", err)
		}
	} else {
		if err != nil {
			return fmt.Errorf("failed to check if G2 point exists: %w", err)
		}
	}

	filePath, err = config.path(SRSPathG2PowerOf2)
	if err != nil {
		return fmt.Errorf("failed to get SRS G2 power of 2 path: %w", err)
	}
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		command := make([]string, 3)
		command[0] = "wget"
		command[1] = "https://srs-mainnet.s3.amazonaws.com/kzg/g2.point.powerOf2"
		command[2] = "--output-document=" + filePath
		logger.Info("executing %s", command)

		cmd := exec.Command(command[0], command[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to download G2 power of 2 point: %w", err)
		}
	} else {
		if err != nil {
			return fmt.Errorf("failed to check if G2 power of 2 point exists: %w", err)
		}
	}

	// Check to see if the private key file exists. If not, stop the test.
	filePath, err = ResolveTildeInPath(config.KeyPath)
	if err != nil {
		return fmt.Errorf("failed to resolve tilde in path: %w", err)
	}
	_, err = os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("private key file %s does not exist. This file should "+
			"contain the private key for the account used in the test, in hex: %w", filePath, err)
	}

	return nil
}
