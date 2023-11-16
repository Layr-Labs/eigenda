package integration_test

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/clients"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/logging"
	rollupbindings "github.com/Layr-Labs/eigenda/contracts/bindings/MockRollup"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/encoding"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigenda/inabox/strategies/processes/deploy"
	"github.com/Layr-Labs/eigenda/pkg/encoding/kzgEncoder"
	gcommon "github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest/v3"
	"github.com/shurcooL/graphql"
)

var (
	templateName      string
	testName          string
	inMemoryBlobStore bool

	testConfig         *deploy.Config
	dockertestPool     *dockertest.Pool
	dockertestResource *dockertest.Resource
	localStackPort     string

	metadataTableName = "test-BlobMetadata"
	bucketTableName   = "test-BucketStore"
	logger            common.Logger
	ethClient         common.EthClient
	mockRollup        *rollupbindings.ContractMockRollup
	retrievalClient   clients.RetrievalClient
)

func init() {
	flag.StringVar(&templateName, "config", "testconfig-anvil-nograph.yaml", "Name of the config file (in `inabox/templates`)")
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
		testName, err = deploy.CreateNewTestDirectory(templateName, rootPath)
		if err != nil {
			Expect(err).To(BeNil())
		}
	}

	testConfig = deploy.NewTestConfig(testName, rootPath)
	if testConfig.Environment.IsLocal() {
		if !inMemoryBlobStore {
			fmt.Println("Using shared Blob Store")
			localStackPort = "4570"
			pool, resource, err := deploy.StartDockertestWithLocalstackContainer(localStackPort)
			Expect(err).To(BeNil())
			dockertestPool = pool
			dockertestResource = resource

			err = deploy.DeployResources(pool, localStackPort, metadataTableName, bucketTableName)
			Expect(err).To(BeNil())

		} else {
			fmt.Println("Using in-memory Blob Store")
		}

		fmt.Println("Starting anvil")
		testConfig.StartAnvil()

		if deployer, ok := testConfig.GetDeployer(testConfig.EigenDA.Deployer); ok && deployer.DeploySubgraphs {
			fmt.Println("Starting graph node")
			testConfig.StartGraphNode()
		}

		fmt.Println("Deploying experiment")
		testConfig.DeployExperiment()

		fmt.Println("Starting binaries")
		testConfig.StartBinaries()
	}
	logger, err = logging.GetLogger(logging.DefaultCLIConfig())
	Expect(err).To(BeNil())
	pk := testConfig.Pks.EcdsaMap["default"].PrivateKey
	pk = strings.TrimPrefix(pk, "0x")
	pk = strings.TrimPrefix(pk, "0X")
	ethClient, err = geth.NewClient(geth.EthClientConfig{
		RPCURL:           testConfig.Deployers[0].RPC,
		PrivateKeyString: pk,
	}, logger)
	Expect(err).To(BeNil())
	mockRollup, err = rollupbindings.NewContractMockRollup(gcommon.HexToAddress(testConfig.MockRollup), ethClient)
	Expect(err).To(BeNil())
	err = setupRetrievalClient(testConfig)
	Expect(err).To(BeNil())
})

func setupRetrievalClient(testConfig *deploy.Config) error {
	ethClientConfig := geth.EthClientConfig{
		RPCURL:           testConfig.Deployers[0].RPC,
		PrivateKeyString: "351b8eca372e64f64d514f90f223c5c4f86a04ff3dcead5c27293c547daab4ca", // just random private key
	}
	client, err := geth.NewClient(ethClientConfig, logger)
	if err != nil {
		return err
	}
	tx, err := eth.NewTransactor(logger, client, testConfig.Retriever.RETRIEVER_BLS_OPERATOR_STATE_RETRIVER, testConfig.Retriever.RETRIEVER_EIGENDA_SERVICE_MANAGER)
	if err != nil {
		return err
	}

	cs := eth.NewChainState(tx, client)
	querier := graphql.NewClient(testConfig.Churner.CHURNER_GRAPH_URL, nil)
	ics := thegraph.NewIndexedChainState(cs, querier, logger)
	agn := &core.StdAssignmentCoordinator{}
	nodeClient := clients.NewNodeClient(20 * time.Second)
	srsOrder, err := strconv.Atoi(testConfig.Retriever.RETRIEVER_SRS_ORDER)
	if err != nil {
		return err
	}
	encoder, err := encoding.NewEncoder(encoding.EncoderConfig{
		KzgConfig: kzgEncoder.KzgConfig{
			G1Path:         testConfig.Retriever.RETRIEVER_G1_PATH,
			G2Path:         testConfig.Retriever.RETRIEVER_G2_PATH,
			CacheDir:       testConfig.Retriever.RETRIEVER_CACHE_PATH,
			NumWorker:      1,
			SRSOrder:       uint64(srsOrder),
			Verbose:        true,
			PreloadEncoder: true,
		},
	})
	if err != nil {
		return err
	}

	retrievalClient = clients.NewRetrievalClient(logger, ics, agn, nodeClient, encoder, 10)
	return nil
}

var _ = AfterSuite(func() {
	if testConfig.Environment.IsLocal() {

		fmt.Println("Stopping binaries")
		testConfig.StopBinaries()

		fmt.Println("Stopping anvil")
		testConfig.StopAnvil()

		fmt.Println("Stopping graph node")
		testConfig.StopGraphNode()

		deploy.PurgeDockertestResources(dockertestPool, dockertestResource)
	}
})
