package integration_test

import (
	"context"
	"flag"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/logging"
	rollupbindings "github.com/Layr-Labs/eigenda/contracts/bindings/MockRollup"
	"github.com/Layr-Labs/eigenda/inabox/config"
	genenv "github.com/Layr-Labs/eigenda/inabox/gen-env"
	gcommon "github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
)

var (
	templateName      string
	testName          string
	inMemoryBlobStore bool

	testConfig    *config.Config
	compose       tc.ComposeStack
	composeCancel func()

	logger     common.Logger
	ethClient  common.EthClient
	mockRollup *rollupbindings.ContractMockRollup
	// retrievalClient clients.RetrievalClient
)

func init() {
	flag.StringVar(&templateName, "config", "testconfig-anvil.yaml", "Name of the config file (in `inabox/templates`)")
	flag.StringVar(&testName, "testname", "", "Name of the test (in `inabox/testdata`)")
	flag.BoolVar(&inMemoryBlobStore, "inMemoryBlobStore", false, "whether to use in-memory blob store")
}

func TestInaboxIntegration(t *testing.T) {
	RegisterFailHandler(Fail)

	if testing.Short() {
		t.Skip()
	}

	RunSpecs(t, "Integration Suite")
}

var _ = BeforeSuite(func() {
	By("bootstrapping test environment")

	rootPath := "../../"

	var err error
	if testName == "" {
		testName, err = config.CreateNewTestDirectory(templateName, rootPath)
		Expect(err).To(BeNil())
	}

	testConfig = config.OpenConfig(filepath.Join(rootPath, "inabox/testdata", testName, "config.yaml"))
	lock := genenv.GenerateConfigLock(rootPath, testName)
	genenv.GenerateDockerCompose(lock)
	genenv.CompileDockerCompose(rootPath, testName)

	StartEigenDA(rootPath, testName)

	pk := lock.Config.Pks.EcdsaMap["default"].PrivateKey
	pk = strings.TrimPrefix(pk, "0x")

	logger, err = logging.GetLogger(logging.DefaultCLIConfig())
	Expect(err).To(BeNil())

	ethClient = NewEthClient(pk)

	mockRollup, err = rollupbindings.NewContractMockRollup(gcommon.HexToAddress(testConfig.MockRollup), ethClient)
	Expect(err).To(BeNil())
})

var _ = AfterSuite(func() {
	// if testConfig.Environment.IsLocal() && compose != nil {
	// 	composeCancel()
	// 	err := compose.Down(context.Background(), tc.RemoveOrphans(true), tc.RemoveImagesLocal, tc.RemoveVolumes(true))
	// 	Expect(err).To(BeNil())
	// }
})

func NewEthClient(pk string) *geth.EthClient {
	ethClient, err := geth.NewClient(geth.EthClientConfig{
		RPCURL:           "http://localhost:8545",
		PrivateKeyString: pk,
	}, logger)
	Expect(err).To(BeNil())
	return ethClient
}

func StartEigenDA(rootPath, testName string) {
	composeFilePath := filepath.Join(rootPath, "inabox/testdata", testName, "docker-compose.yml")
	var err error
	compose, err = tc.NewDockerCompose(composeFilePath)
	Expect(err).To(BeNil())

	var ctx context.Context
	ctx, composeCancel = context.WithCancel(context.Background())
	err = compose.Up(ctx, tc.Wait(true), tc.RemoveOrphans(true))
	Expect(err).To(BeNil())
}
