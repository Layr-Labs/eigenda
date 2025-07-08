package client

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/litt/util"
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

// GetEnvironmentConfigPaths returns a list of paths to the environment config files.
func GetEnvironmentConfigPaths() ([]string, error) {
	configDir, err := util.SanitizePath("../config/environment")
	if err != nil {
		return nil, fmt.Errorf("failed to sanitize path: %w", err)
	}

	files, err := os.ReadDir(configDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read environment config directory: %w", err)
	}
	var configPaths []string
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		configPath := fmt.Sprintf("../config/environment/%s", file.Name())
		configPaths = append(configPaths, configPath)
	}
	if len(configPaths) == 0 {
		return nil, fmt.Errorf("no environment config files found in ../config/environment")
	}
	return configPaths, nil
}

// GetConfig returns a TestClientConfig instance parsed from the config file.
func GetConfig(configPath string) (*TestClientConfig, error) {
	configLock.Lock()
	defer configLock.Unlock()

	if config, ok := configMap[configPath]; ok {
		return config, nil
	}

	configFile, err := util.SanitizePath(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed sanitize path: %w", err)
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

	// The environment variable "CI" will be set when running inside a github action.
	// The environment variable "LIVE_TESTS" will be set when running live tests, which is a specific github action.
	//
	// There are three situations we want to consider:
	//
	// 1. When running a tests locally, we want to run live tests if requested. "CI" will not be set, and so
	//    we will not skip the test.
	// 2. When we are running general unit tests as a github action, we specifically don't want to run live tests.
	//    "CI" will be set, and "LIVE_TESTS" will not be set, so we skip the test.
	// 3. When we are running live tests as a github action, we want to run the test. Both "CI" and "LIVE_TESTS" will
	//    be set, so we do not skip the test.
	if os.Getenv("CI") != "" && os.Getenv("LIVE_TESTS") == "" {
		t.Skip("Skipping test in CI environment")
	}
}
