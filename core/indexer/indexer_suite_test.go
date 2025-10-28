package indexer_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/Layr-Labs/eigenda/test/testbed"
)

var (
	anvilContainer  *testbed.AnvilContainer
	templateName    string
	testName        string
	headerStoreType string

	testConfig *deploy.Config
)

func init() {
	flag.StringVar(&templateName, "config", "testconfig-anvil.yaml", "Name of the config file (in `inabox/templates`)")
	flag.StringVar(&testName, "testname", "", "Name of the test (in `inabox/testdata`)")
	flag.StringVar(&headerStoreType, "headerStore", "leveldb",
		"The header store implementation to be used (inmem, leveldb)")
}

func TestMain(m *testing.M) {
	flag.Parse()

	if testing.Short() {
		fmt.Println("Skipping integration tests in short mode")
		os.Exit(0)
	}

	rootPath := "../../"
	logger := test.GetLogger()

	if testName == "" {
		var err error
		testName, err = deploy.CreateNewTestDirectory(templateName, rootPath)
		if err != nil {
			logger.Fatal("Failed to create test directory:", err)
		}
	}

	testConfig = deploy.ReadTestConfig(testName, rootPath)
	testConfig.Deployers[0].DeploySubgraphs = false

	if testConfig.Environment.IsLocal() {
		logger.Info("Starting anvil")
		var err error
		anvilContainer, err = testbed.NewAnvilContainerWithOptions(context.Background(), testbed.AnvilOptions{
			ExposeHostPort: true, // This will bind container port 8545 to host port 8545
			Logger:         logger,
		})
		if err != nil {
			logger.Fatal("Failed to start anvil container:", err)
		}

		logger.Info("Deploying experiment")
		if err := testConfig.DeployExperiment(); err != nil {
			logger.Fatal("Failed to deploy experiment:", err)
		}
	}

	code := m.Run()

	// Cleanup
	if testConfig != nil && testConfig.Environment.IsLocal() && anvilContainer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_ = anvilContainer.Terminate(ctx)
	}

	os.Exit(code)
}
