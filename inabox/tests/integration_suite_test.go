package integration_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/Layr-Labs/eigenda/test"
)

// Global infrastructure that is shared across all tests
var globalInfra *InfrastructureHarness

// Configuration constants from command line flags
var (
	templateName      string
	testName          string
	inMemoryBlobStore bool
	logger            = test.GetLogger()
)

func init() {
	flag.StringVar(&templateName, "config", "testconfig-anvil.yaml", "Name of the config file (in `inabox/templates`)")
	flag.StringVar(&testName, "testname", "", "Name of the test (in `inabox/testdata`)")
	flag.BoolVar(&inMemoryBlobStore, "inMemoryBlobStore", false, "whether to use in-memory blob store")
}

func TestMain(m *testing.M) {
	flag.Parse()

	if testing.Short() {
		logger.Info("Skipping inabox integration tests in short mode")
		os.Exit(0)
	}

	// Run suite setup
	if err := setupSuite(); err != nil {
		logger.Error("Setup failed:", err)
		teardownSuite()
		os.Exit(1)
	}

	// Run all tests
	code := m.Run()

	// Run suite teardown
	teardownSuite()

	// Exit with test result code
	os.Exit(code)
}

func setupSuite() error {
	logger.Info("bootstrapping test environment")

	// Setup the global infrastructure
	var err error
	globalInfra, err = SetupGlobalInfrastructure(templateName, testName, inMemoryBlobStore, logger)
	if err != nil {
		return fmt.Errorf("failed to setup global infrastructure: %w", err)
	}

	return nil
}

func teardownSuite() {
	logger.Info("Tearing down test environment")

	// Teardown the global infrastructure
	if globalInfra != nil {
		TeardownGlobalInfrastructure(globalInfra)
	}

	logger.Info("Teardown completed")
}
