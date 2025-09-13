package indexer_test

import (
	"context"
	"flag"
	"testing"

	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/Layr-Labs/eigenda/test/testbed"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	anvilContainer *testbed.AnvilContainer
	templateName   string
	testName       string

	testConfig *deploy.Config
)

func init() {
	flag.StringVar(&templateName, "config", "testconfig-anvil-nochurner.yaml", "Name of the config file (in `inabox/templates`)")
	flag.StringVar(&testName, "testname", "", "Name of the test (in `inabox/testdata`)")
}

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = BeforeSuite(func() {
	By("bootstrapping test environment")

	if !testing.Short() {
		rootPath := "../../"

		if testName == "" {
			var err error
			testName, err = deploy.CreateNewTestDirectory(templateName, rootPath)
			if err != nil {
				Expect(err).To(BeNil())
			}
		}

		testConfig = deploy.NewTestConfig(testName, rootPath)
		testConfig.Deployers[0].DeploySubgraphs = false
		logger := test.GetLogger()

		if testConfig.Environment.IsLocal() {
			logger.Info("Starting anvil")
			var err error
			anvilContainer, err = testbed.NewAnvilContainerWithOptions(context.Background(), testbed.AnvilOptions{
				ExposeHostPort: true, // This will bind container port 8545 to host port 8545
				Logger:         logger,
			})
			if err != nil {
				panic(err)
			}

			logger.Info("Deploying experiment")
			if err := testConfig.DeployExperiment(); err != nil {
				panic(err)
			}
		}
	}

})

var _ = AfterSuite(func() {
	if !testing.Short() && testConfig.Environment.IsLocal() {
		_ = anvilContainer.Terminate(context.Background())
	}

})
