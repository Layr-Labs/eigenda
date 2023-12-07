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
	"github.com/Layr-Labs/eigenda/inabox/config"
	genenv "github.com/Layr-Labs/eigenda/inabox/gen-env"
	"github.com/Layr-Labs/eigenda/inabox/testutils"
	"github.com/shurcooL/graphql"
	"github.com/stretchr/testify/assert"
)

var (
	templateName string
	testName     string
	graphUrl     string
	eigenDA      *testutils.EigenDA
	lock         *config.ConfigLock
)

func init() {
	flag.StringVar(&templateName, "config", "testconfig-anvil.yaml", "Name of the config file (in `inabox/templates`)")
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
		testName, err = config.CreateNewTestDirectory(templateName, rootPath)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Starting EigenDA")
	lock = genenv.GenerateConfigLock(rootPath, testName)
	genenv.GenerateDockerCompose(lock)
	genenv.CompileDockerCompose(rootPath, testName)
	eigenDA = testutils.NewEigenDA(rootPath, testName)
	eigenDA.MustStart()
}

func teardown() {
	fmt.Println("Stopping EigenDA")
	// eigenDA.MustStop()
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

	client := mustMakeTestClient(t, lock, lock.Envs.Batcher.BATCHER_PRIVATE_KEY, logger)
	tx, err := eth.NewTransactor(logger, client, lock.Config.EigenDA.OperatorStateRetreiver, lock.Config.EigenDA.ServiceManager)
	assert.NoError(t, err)

	cs := thegraph.NewIndexedChainState(eth.NewChainState(tx, client), graphql.NewClient(graphUrl, nil), logger)
	time.Sleep(5 * time.Second)

	err = cs.Start(context.Background())
	assert.NoError(t, err)

	headerNum, err := cs.GetCurrentBlockNumber()
	assert.NoError(t, err)

	state, err := cs.GetIndexedOperatorState(context.Background(), headerNum, quorums)
	assert.NoError(t, err)
	assert.Equal(t, len(lock.Operators), len(state.IndexedOperators))
}

func mustMakeTestClient(t *testing.T, env *config.ConfigLock, privateKey string, logger common.Logger) common.EthClient {
	deployer, ok := env.Config.GetDeployer(env.Config.EigenDA.Deployer)
	var rpc string
	if deployer.LocalAnvil {
		rpc = "http://localhost:8545"
	} else {
		rpc = deployer.RPC
	}
	assert.True(t, ok)

	config := geth.EthClientConfig{
		RPCURL:           rpc,
		PrivateKeyString: privateKey,
	}

	client, err := geth.NewClient(config, logger)
	assert.NoError(t, err)
	return client
}
