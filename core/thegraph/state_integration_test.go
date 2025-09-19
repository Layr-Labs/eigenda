package thegraph_test

import (
	"context"
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/Layr-Labs/eigenda/test/testbed"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/shurcooL/graphql"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/network"
)

var (
	templateName string
	testName     string
	graphUrl     string

	localstackPort      = "4570"
	metadataTableName   = "test-BlobMetadata"
	bucketTableName     = "test-BucketStore"
	metadataTableNameV2 = "test-BlobMetadata-v2"
	testQuorums         = []uint8{0, 1}

	logger = test.GetLogger()
)

func init() {
	flag.StringVar(&templateName, "config", "testconfig-anvil-nochurner.yaml", "Name of the config file (in `inabox/templates`)")
	flag.StringVar(&testName, "testname", "", "Name of the test (in `inabox/testdata`)")
	flag.StringVar(&graphUrl, "graphurl", "http://localhost:8000/subgraphs/name/Layr-Labs/eigenda-operator-state", "")
}

func setupTest(t *testing.T) (
	*testbed.AnvilContainer, *testbed.LocalStackContainer, *testbed.GraphNodeContainer, *deploy.Config,
) {
	t.Helper()

	if testing.Short() {
		t.Skip("Skipping graph indexer integration test in short mode")
	}

	flag.Parse()
	ctx := t.Context()
	rootPath := "../../"

	if testName == "" {
		var err error
		testName, err = deploy.CreateNewTestDirectory(templateName, rootPath)
		require.NoError(t, err, "failed to create test directory")
	}

	testConfig := deploy.NewTestConfig(testName, rootPath)
	testConfig.Deployers[0].DeploySubgraphs = true

	// Create a shared Docker network for all containers
	nw, err := network.New(ctx,
		network.WithDriver("bridge"),
		network.WithAttachable())
	require.NoError(t, err, "failed to create Docker network")
	logger.Info("Created Docker network", "name", nw.Name)

	localstackContainer, err := testbed.NewLocalStackContainerWithOptions(ctx, testbed.LocalStackOptions{
		ExposeHostPort: true,
		HostPort:       localstackPort,
		Services:       []string{"s3", "dynamodb", "kms"},
		Logger:         logger,
		Network:        nw,
	})
	require.NoError(t, err, "failed to start localstack container")

	deployConfig := testbed.DeployResourcesConfig{
		LocalStackEndpoint:  fmt.Sprintf("http://0.0.0.0:%s", localstackPort),
		MetadataTableName:   metadataTableName,
		BucketTableName:     bucketTableName,
		V2MetadataTableName: metadataTableNameV2,
		Logger:              logger,
	}
	err = testbed.DeployResources(ctx, deployConfig)
	require.NoError(t, err, "failed to deploy resources")

	anvilContainer, err := testbed.NewAnvilContainerWithOptions(ctx, testbed.AnvilOptions{
		ExposeHostPort: true,
		HostPort:       "8545",
		Logger:         logger,
		Network:        nw,
	})
	require.NoError(t, err, "failed to start anvil container")
	anvilContainerPort := anvilContainer.RpcURL()
	anvilInternalEndpoint := anvilContainer.InternalEndpoint()
	logger.Info("Anvil RPC URL", "url", anvilContainerPort, "internal", anvilInternalEndpoint)

	logger.Info("Starting graph node")
	graphNodeContainer, err := testbed.NewGraphNodeContainerWithOptions(ctx, testbed.GraphNodeOptions{
		PostgresDB:     "graph-node",
		PostgresUser:   "graph-node",
		PostgresPass:   "let-me-in",
		EthereumRPC:    anvilInternalEndpoint,
		ExposeHostPort: true,
		HostHTTPPort:   "8000",
		HostWSPort:     "8001",
		HostAdminPort:  "8020",
		HostIPFSPort:   "5001",
		Logger:         logger,
		Network:        nw,
	})
	require.NoError(t, err, "failed to start graph node")

	// Update the graph URL to use the new container
	graphUrl = graphNodeContainer.HTTPURL() + "/subgraphs/name/Layr-Labs/eigenda-operator-state"

	logger.Info("Deploying experiment")
	err = testConfig.DeployExperiment()
	require.NoError(t, err, "failed to deploy experiment")

	logger.Info("Starting binaries")
	testConfig.StartBinaries()

	t.Cleanup(func() {
		logger.Info("Stopping containers and services")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		logger.Info("Stop binaries")
		testConfig.StopBinaries()

		logger.Info("Stop graph node")
		_ = graphNodeContainer.Terminate(ctx)

		_ = anvilContainer.Terminate(ctx)
		_ = localstackContainer.Terminate(ctx)

		logger.Info("Removing Docker network")
		_ = nw.Remove(ctx)
	})

	return anvilContainer, localstackContainer, graphNodeContainer, testConfig
}

func TestIndexerIntegration(t *testing.T) {
	ctx := t.Context()
	_, _, _, testConfig := setupTest(t)

	client := mustMakeTestClient(t, testConfig, testConfig.Batcher[0].BATCHER_PRIVATE_KEY, logger)
	tx, err := eth.NewWriter(
		logger, client, testConfig.EigenDA.OperatorStateRetriever, testConfig.EigenDA.ServiceManager)
	require.NoError(t, err, "failed to create eth writer")

	cs := thegraph.NewIndexedChainState(eth.NewChainState(tx, client), graphql.NewClient(graphUrl, nil), logger)
	time.Sleep(5 * time.Second)

	err = cs.Start(ctx)
	require.NoError(t, err, "failed to start indexed chain state")

	headerNum, err := cs.GetCurrentBlockNumber(ctx)
	require.NoError(t, err, "failed to get current block number")

	state, err := cs.GetIndexedOperatorState(ctx, headerNum, testQuorums)
	require.NoError(t, err, "failed to get indexed operator state")
	require.Equal(t, len(testConfig.Operators), len(state.IndexedOperators), "operator count mismatch")
}

func mustMakeTestClient(t *testing.T, env *deploy.Config, privateKey string, logger logging.Logger) common.EthClient {
	t.Helper()

	deployer, ok := env.GetDeployer(env.EigenDA.Deployer)
	require.True(t, ok, "failed to get deployer")

	config := geth.EthClientConfig{
		RPCURLs:          []string{deployer.RPC},
		PrivateKeyString: privateKey,
		NumConfirmations: 0,
		NumRetries:       0,
	}

	client, err := geth.NewClient(config, gethcommon.Address{}, 0, logger)
	require.NoError(t, err, "failed to create eth client")
	return client
}
