package thegraph_test

import (
	"context"
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/logging"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/shurcooL/graphql"
	"github.com/stretchr/testify/assert"
)

var (
	templateName string
	testName     string
	graphUrl     string
	testConfig   *deploy.Config
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

	fmt.Println("Starting anvil")
	testConfig.StartAnvil()

	fmt.Println("Starting graph node")
	testConfig.StartGraphNode()

	fmt.Println("Deploying experiment")
	testConfig.DeployExperiment()

	fmt.Println("Starting binaries")
	testConfig.StartBinaries()

}

func teardown() {
	fmt.Println("Stoping anvil")
	testConfig.StopAnvil()

	fmt.Println("Stop graph node")
	testConfig.StopGraphNode()

	fmt.Println("Stop binaries")
	testConfig.StopBinaries()
}

func TestIndexerIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip graph indexer integrations test in short mode")
	}
	setup()
	defer teardown()

	logger, err := logging.GetLogger(logging.Config{
		StdLevel:  "debug",
		FileLevel: "debug",
	})
	assert.NoError(t, err)

	client := mustMakeTestClient(t, testConfig, testConfig.Batcher[0].BATCHER_PRIVATE_KEY, logger)
	tx, err := eth.NewTransactor(logger, client, testConfig.EigenDA.OperatorStateRetreiver, testConfig.EigenDA.ServiceManager)
	assert.NoError(t, err)

	cs := thegraph.NewIndexedChainState(eth.NewChainState(tx, client), graphql.NewClient(graphUrl, nil), logger)
	time.Sleep(5 * time.Second)

	err = cs.Start(context.Background())
	assert.NoError(t, err)

	headerNum, err := cs.GetCurrentBlockNumber()
	assert.NoError(t, err)

	state, err := cs.GetIndexedOperatorState(context.Background(), headerNum, quorums)
	assert.NoError(t, err)
	assert.Equal(t, len(testConfig.Operators), len(state.IndexedOperators))
}

func mustMakeTestClient(t *testing.T, env *deploy.Config, privateKey string, logger common.Logger) common.EthClient {
	deployer, ok := env.GetDeployer(env.EigenDA.Deployer)
	assert.True(t, ok)

	config := geth.EthClientConfig{
		RPCURL:           deployer.RPC,
		PrivateKeyString: privateKey,
		NumConfirmations: 0,
	}

	client, err := geth.NewClient(config, logger)
	assert.NoError(t, err)
	return client
}
