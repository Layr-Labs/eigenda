package thegraph_test

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigenda/testbed"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/shurcooL/graphql"
	"github.com/stretchr/testify/assert"
)

var (
	anvilContainer      *testbed.AnvilContainer
	localstackContainer *testbed.LocalStackContainer
	templateName        string
	testName            string
	graphUrl            string
	testConfig          *deploy.Config

	localstackPort      = "4570"
	metadataTableName   = "test-BlobMetadata"
	bucketTableName     = "test-BucketStore"
	metadataTableNameV2 = "test-BlobMetadata-v2"

	logger = testutils.GetLogger()
)

func init() {
	flag.StringVar(&templateName, "config", "testconfig-anvil-nochurner.yaml", "Name of the config file (in `inabox/templates`)")
	flag.StringVar(&testName, "testname", "", "Name of the test (in `inabox/testdata`)")
	flag.StringVar(&graphUrl, "graphurl", "http://localhost:8000/subgraphs/name/Layr-Labs/eigenda-operator-state", "")
}

func setup() {
	if testing.Short() {
		return
	}

	rootPath := "../../"

	if testName == "" {
		var err error
		testName, err = deploy.CreateNewTestDirectory(templateName, rootPath)
		if err != nil {
			panic(err)
		}
	}

	testConfig = deploy.NewTestConfig(testName, rootPath)
	testConfig.Deployers[0].DeploySubgraphs = true
	var err error
	localstackContainer, err = testbed.NewLocalStackContainerWithOptions(context.Background(), testbed.LocalStackOptions{
		ExposeHostPort: true,
		HostPort:       localstackPort,
		Services:       []string{"s3", "dynamodb", "kms"},
		Logger:         logger,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Deploying LocalStack resources")
	deployConfig := testbed.DeployResourcesConfig{
		LocalStackEndpoint:  fmt.Sprintf("http://0.0.0.0:%s", localstackPort),
		MetadataTableName:   metadataTableName,
		BucketTableName:     bucketTableName,
		V2MetadataTableName: metadataTableNameV2,
		Logger:              logger,
	}
	err = testbed.DeployResources(context.Background(), deployConfig)
	if err != nil {
		panic(err)
	}

	anvilContainer, err = testbed.NewAnvilContainerWithOptions(context.Background(), testbed.AnvilOptions{
		ExposeHostPort: true,
		HostPort:       "8545",
		Logger:         logger,
	})
	if err != nil {
		panic(err)
	}

	logger.Info("Starting graph node")
	testConfig.StartGraphNode()

	logger.Info("Deploying experiment")
	if err := testConfig.DeployExperiment(); err != nil {
		panic(err)
	}

	pk := testConfig.Pks.EcdsaMap["default"].PrivateKey
	pk = strings.TrimPrefix(pk, "0x")
	pk = strings.TrimPrefix(pk, "0X")
	ethClient, err := geth.NewMultiHomingClient(geth.EthClientConfig{
		RPCURLs:          []string{testConfig.Deployers[0].RPC},
		PrivateKeyString: pk,
		NumConfirmations: 0,
		NumRetries:       1,
	}, gethcommon.Address{}, logger)
	if err != nil {
		panic(err)
	}
	testConfig.RegisterBlobVersionAndRelays(ethClient)

	logger.Info("Registering disperser keypair")
	err = testConfig.RegisterDisperserKeypair(ethClient)
	if err != nil {
		panic(err)
	}

	logger.Info("Starting binaries")
	testConfig.StartBinaries()
}

func teardown() {
	_ = localstackContainer.Terminate(context.Background())
	_ = anvilContainer.Terminate(context.Background())

	logger.Info("Stop graph node")
	testConfig.StopGraphNode()

	logger.Info("Stop binaries")
	testConfig.StopBinaries()
}

func TestIndexerIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip graph indexer integrations test in short mode")
	}
	setup()
	defer teardown()

	client := mustMakeTestClient(t, testConfig, testConfig.Batcher[0].BATCHER_PRIVATE_KEY, logger)
	tx, err := eth.NewWriter(
		logger, client, testConfig.EigenDA.OperatorStateRetriever, testConfig.EigenDA.ServiceManager)
	assert.NoError(t, err)

	cs := thegraph.NewIndexedChainState(eth.NewChainState(tx, client), graphql.NewClient(graphUrl, nil), logger)
	time.Sleep(5 * time.Second)

	err = cs.Start(context.Background())
	assert.NoError(t, err)

	headerNum, err := cs.GetCurrentBlockNumber(context.Background())
	assert.NoError(t, err)

	state, err := cs.GetIndexedOperatorState(context.Background(), headerNum, quorums)
	assert.NoError(t, err)
	assert.Equal(t, len(testConfig.Operators), len(state.IndexedOperators))
}

func mustMakeTestClient(t *testing.T, env *deploy.Config, privateKey string, logger logging.Logger) common.EthClient {
	deployer, ok := env.GetDeployer(env.EigenDA.Deployer)
	assert.True(t, ok)

	config := geth.EthClientConfig{
		RPCURLs:          []string{deployer.RPC},
		PrivateKeyString: privateKey,
		NumConfirmations: 0,
		NumRetries:       0,
	}

	client, err := geth.NewClient(config, gethcommon.Address{}, 0, logger)
	assert.NoError(t, err)
	return client
}
