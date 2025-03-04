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
	configLock sync.Mutex
	clientLock sync.Mutex
	configMap  = make(map[string]*TestClientConfig)
	clientMap  = make(map[string]*TestClient)
	logger     logging.Logger
	metrics    *testClientMetrics
)

const (
	PreprodEnv = "../config/environment/preprod.json"
	TestnetEnv = "../config/environment/testnet.json"

	G1URL         = "https://srs-mainnet.s3.amazonaws.com/kzg/g1.point"
	G2URL         = "https://srs-mainnet.s3.amazonaws.com/kzg/g2.point"
	G2PowerOf2URL = "https://srs-mainnet.s3.amazonaws.com/kzg/g2.point.powerOf2"
)

// GetConfig returns a TestClientConfig instance parsed from the config file.
func GetConfig(configPath string) (*TestClientConfig, error) {
	configLock.Lock()
	defer configLock.Unlock()

	if config, ok := configMap[configPath]; ok {
		return config, nil
	}

	configFile, err := ResolveTildeInPath(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve tilde in path: %w", err)
	}
	configFileBytes, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := &TestClientConfig{}
	err = json.Unmarshal(configFileBytes, config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	configMap[configPath] = config

	return config, nil
}

// GetTestClient is the same as GetClient, but also performs a check to ensure that the test is not
// running in a CI environment. If using a TestClient in a unit test, it is critical to use this method
// to ensure that the test is not running in a CI environment.
func GetTestClient(t *testing.T, configPath string) *TestClient {
	skipInCI(t)
	c, err := GetClient(configPath)
	require.NoError(t, err)
	return c
}

// GetClient returns a TestClient instance, creating one if it does not exist.
// This uses a global static client... this is icky, but it takes a long time
// to read the SRS points, so it's the lesser of two evils to keep it around.
func GetClient(configPath string) (*TestClient, error) {
	clientLock.Lock()
	defer clientLock.Unlock()

	if client, ok := clientMap[configPath]; ok {
		return client, nil
	}

	testConfig, err := GetConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	if len(clientMap) == 0 {
		// do one time setup

		var loggerConfig common.LoggerConfig
		if os.Getenv("CI") != "" {
			loggerConfig = common.DefaultLoggerConfig()
		} else {
			loggerConfig = common.DefaultConsoleLoggerConfig()
		}

		logger, err = common.NewLogger(loggerConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create logger: %w", err)
		}

		// only do this stuff once
		err = setupFilesystem(testConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to setup filesystem: %w", err)
		}

		testMetrics := newTestClientMetrics(logger, testConfig.MetricsPort)
		metrics = testMetrics
		testMetrics.start()
	}

	client, err := NewTestClient(logger, metrics, testConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create test client: %w", err)
	}

	clientMap[configPath] = client

	return client, nil
}

func skipInCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping test in CI environment")
	}
}

// ensureSRSFileIsPresent checks if a file exists at the given path. If it does not, it downloads the file from the
// given URL into the given path.
func ensureSRSFileIsPresent(
	config *TestClientConfig,
	filePath string,
	url string) error {

	var err error
	filePath, err = config.ResolveSRSPath(filePath)
	if err != nil {
		return fmt.Errorf("failed to resolve SRS path: %w", err)
	}

	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		command := make([]string, 3)
		command[0] = "wget"
		command[1] = url
		command[2] = "--output-document=" + filePath
		logger.Info("executing %s", command)

		cmd := exec.Command(command[0], command[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to download %s: %w", url, err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check if file exists: %w", err)
	}

	return nil
}

func setupFilesystem(config *TestClientConfig) error {
	// Create the test data directory if it does not exist
	srsPath, err := ResolveTildeInPath(config.SRSPath)
	if err != nil {
		return fmt.Errorf("failed to resolve tilde in path: %w", err)
	}
	err = os.MkdirAll(srsPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create test data directory: %w", err)
	}

	// Create the SRS directories if they do not exist
	srsTablesPath, err := config.ResolveSRSPath(SRSPathSRSTables)
	if err != nil {
		return fmt.Errorf("failed to resolve SRS tables path: %w", err)
	}
	err = os.MkdirAll(srsTablesPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create SRS tables directory: %w", err)
	}

	// If any of the srs files do not exist, download them.
	err = ensureSRSFileIsPresent(config, SRSPathG1, G1URL)
	if err != nil {
		return fmt.Errorf("failed to locate G1 point: %w", err)
	}
	err = ensureSRSFileIsPresent(config, SRSPathG2, G2URL)
	if err != nil {
		return fmt.Errorf("failed to locate G2 point: %w", err)
	}
	err = ensureSRSFileIsPresent(config, SRSPathG2PowerOf2, G2PowerOf2URL)
	if err != nil {
		return fmt.Errorf("failed to locate G2 power of 2 point: %w", err)
	}

	return nil
}
