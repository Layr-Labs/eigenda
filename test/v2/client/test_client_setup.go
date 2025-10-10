package client

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/config"
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
	metrics    *TestClientMetrics
)

const LiveTestPrefix = "LIVE_TEST"

// GetEnvironmentConfigPaths returns a list of paths to the environment config files.
func GetEnvironmentConfigPaths() ([]string, error) {
	// Golang tests are always run with CWD set to the dir in which the test file is located.
	// These relative paths should thus only be used for tests in direct subdirs of `test/v2`,
	// such as `test/v2/live` where they are currently used from.
	// TODO: GetEnvironmentConfigPaths should take a base path as an argument
	// to allow for more flexibility in where the config files are located.
	configDir, err := util.SanitizePath("../config")
	if err != nil {
		return nil, fmt.Errorf("failed to sanitize path: %w", err)
	}

	files, err := os.ReadDir(configDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read environment config directory: %w", err)
	}
	var configPaths []string
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".toml") {
			continue
		}
		configPath := fmt.Sprintf("../config/%s", file.Name())
		configPaths = append(configPaths, configPath)
	}
	if len(configPaths) == 0 {
		return nil, fmt.Errorf("no environment config files found in ../config")
	}
	return configPaths, nil
}

// GetConfig returns a TestClientConfig instance parsed from the config file.
func GetConfig(logger logging.Logger, prefix string, configPath string) (*TestClientConfig, error) {
	configLock.Lock()
	defer configLock.Unlock()

	if cfg, ok := configMap[configPath]; ok {
		return cfg, nil
	}

	cfg, err := config.ParseConfig(logger, DefaultTestClientConfig(), prefix, configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Resolve relative SRS path based on config file location
	if cfg.SrsPath != "" && !filepath.IsAbs(cfg.SrsPath) {
		configDir := filepath.Dir(configPath)
		absPath := filepath.Join(configDir, cfg.SrsPath)
		cfg.SrsPath = filepath.Clean(absPath)
		// to debug this, you can print filepath.Abs(cfg.SrsPath)
	}

	configMap[configPath] = cfg

	return cfg, nil
}

// GetTestClient is the same as GetClient, but also performs a check to ensure that the test is not
// running in a CI environment. If using a TestClient in a unit test, it is critical to use this method
// to ensure that the test is not running in a CI environment.
func GetTestClient(t *testing.T, logger logging.Logger, configPath string) *TestClient {
	skipInCI(t)
	c, err := GetClient(logger, configPath)
	require.NoError(t, err)
	return c
}

// GetClient returns a TestClient instance, creating one if it does not exist.
// This uses a global static client... this is icky, but it takes a long time
// to read the SRS points, so it's the lesser of two evils to keep it around.
func GetClient(logger logging.Logger, configPath string) (*TestClient, error) {
	clientLock.Lock()
	defer clientLock.Unlock()

	if client, ok := clientMap[configPath]; ok {
		return client, nil
	}

	testConfig, err := GetConfig(logger, LiveTestPrefix, configPath)
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
			testMetrics := NewTestClientMetrics(logger, testConfig.MetricsPort)
			metrics = testMetrics
			testMetrics.Start()
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
