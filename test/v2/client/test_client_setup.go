package client

import (
	"encoding/json"
	"fmt"
	"os"
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

		// TODO (cody.littley): add a setting to enable colored logging
		loggerConfig := common.DefaultTextLoggerConfig()

		logger, err = common.NewLogger(loggerConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create logger: %w", err)
		}

		if !testConfig.DisableMetrics {
			testMetrics := newTestClientMetrics(logger, testConfig.MetricsPort)
			metrics = testMetrics
			testMetrics.start()
		}
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
