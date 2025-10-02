package client

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/Layr-Labs/eigenda/test"
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
	// Golang tests are always run with CWD set to the dir in which the test file is located.
	// These relative paths should thus only be used for tests in direct subdirs of `test/v2`,
	// such as `test/v2/live` where they are currently used from.
	// TODO: GetEnvironmentConfigPaths should take a base path as an argument
	// to allow for more flexibility in where the config files are located.
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

	config := DefaultTestClientConfig()
	err = json.Unmarshal(configFileBytes, config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	// Resolve relative SRS path based on config file location
	if config.SRSPath != "" && !filepath.IsAbs(config.SRSPath) {
		configDir := filepath.Dir(configFile)
		absPath := filepath.Join(configDir, config.SRSPath)
		config.SRSPath = filepath.Clean(absPath)
		// to debug this, you can print filepath.Abs(config.SRSPath)
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

	client, err := NewTestClient(context.Background(), logger, metrics, testConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create test client: %w", err)
	}

	clientMap[configPath] = client

	return client, nil
}

// skipInCI skips the test if running in a CI environment, unless explicitly running live tests in CI.
func skipInCI(t *testing.T) {
	// Normally we want to skip these tests in CI. But if we are explicitly running live tests in CI,
	// do not skip them, even though we are in a CI environment.
	if os.Getenv("LIVE_TESTS") != "" {
		return
	}

	// We aren't running a live test, so skip if in CI.
	test.SkipInCI(t)
}
