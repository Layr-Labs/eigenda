package indexer_test

import (
	"flag"
	"fmt"
	"testing"

	"github.com/Layr-Labs/eigenda/inabox/strategies/processes/deploy"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	templateName string
	testName     string

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

		if testConfig.Environment.IsLocal() {
			fmt.Println("Starting anvil")
			testConfig.StartAnvil()

			fmt.Println("Deploying experiment")
			testConfig.DeployExperiment()
		}
	}

})

var _ = AfterSuite(func() {

	if !testing.Short() && testConfig.Environment.IsLocal() {
		fmt.Println("Stopping anvil")
		testConfig.StopAnvil()
	}

})
