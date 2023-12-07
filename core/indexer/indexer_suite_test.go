package indexer_test

import (
	"flag"
	"fmt"
	"testing"

	"github.com/Layr-Labs/eigenda/inabox/config"
	genenv "github.com/Layr-Labs/eigenda/inabox/gen-env"
	"github.com/Layr-Labs/eigenda/inabox/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	templateName string
	testName     string

	lock  *config.ConfigLock
	anvil *testutils.AnvilContainer
)

func init() {
	flag.StringVar(&templateName, "config", "testconfig-anvil-local.yaml", "Name of the config file (in `inabox/templates`)")
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
			testName, err = config.CreateNewTestDirectory(templateName, rootPath)
			if err != nil {
				Expect(err).To(BeNil())
			}
		}

		genenv.GenerateConfigLock(rootPath, testName)
		lock = config.OpenConfigLock(rootPath, testName)
		lock.Config.Deployers[0].DeploySubgraphs = false

		if lock.Config.Environment.IsLocal() {
			fmt.Println("Starting anvil")

			anvil = testutils.NewAnvilContainer(lock)
			anvil.MustStart()
		}
	}

})

var _ = AfterSuite(func() {

	if !testing.Short() && lock.Config.Environment.IsLocal() {
		fmt.Println("Stopping anvil")
		anvil.MustStop()
	}

})
