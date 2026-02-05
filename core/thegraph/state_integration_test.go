package thegraph_test

import (
	"flag"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	inaboxtests "github.com/Layr-Labs/eigenda/inabox/tests"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/shurcooL/graphql"
	"github.com/stretchr/testify/require"
)

var (
	templateName string
	testName     string
	graphUrl     string
	testQuorums  = []uint8{0, 1}
	logger       = test.GetLogger()
)

func init() {
	flag.StringVar(&templateName, "config", "testconfig-anvil-nochurner.yaml", "Name of the config file (in `inabox/templates`)")
	flag.StringVar(&testName, "testname", "", "Name of the test (in `inabox/testdata`)")
	flag.StringVar(&graphUrl, "graphurl", "http://localhost:8000/subgraphs/name/Layr-Labs/eigenda-operator-state", "")
}

func setupTest(t *testing.T) *inaboxtests.InfrastructureHarness {
	t.Helper()

	if testing.Short() {
		t.Skip("Skipping graph indexer integration test in short mode")
	}

	flag.Parse()

	// Setup infrastructure using the centralized function
	config := &inaboxtests.InfrastructureConfig{
		TemplateName:     templateName,
		TestName:         testName,
		Logger:           logger,
		RootPath:         "../../",
		DisableDisperser: true,
	}

	// Start all the necessary infrastructure like anvil, graph node, and eigenda components
	// TODO(dmanc): We really only need to register operators on chain, maybe add some sort of
	// configuration to allow that mode.
	infraHarness, err := inaboxtests.SetupInfrastructure(t.Context(), config)
	require.NoError(t, err, "failed to setup global infrastructure")

	// Update the graph URL to use the container from infrastructure
	if infraHarness.ChainHarness.GraphNode != nil {
		graphUrl = infraHarness.ChainHarness.GraphNode.HTTPURL() + "/subgraphs/name/Layr-Labs/eigenda-operator-state"
	}

	t.Cleanup(func() {
		logger.Info("Tearing down test infrastructure")
		inaboxtests.TeardownInfrastructure(infraHarness)
	})

	return infraHarness
}

func TestIndexerIntegration(t *testing.T) {
	ctx := t.Context()
	infraHarness := setupTest(t)

	client := infraHarness.ChainHarness.EthClient
	tx, err := eth.NewWriter(
		// TODO(dmanc): Expose the operator state retriever and service manager addresses in the infrastructure harness
		// or use the contract directory. Then we can remove the dependency on the test config.
		logger,
		client,
		infraHarness.TestConfig.EigenDA.OperatorStateRetriever,
		infraHarness.TestConfig.EigenDA.ServiceManager,
	)
	require.NoError(t, err, "failed to create eth writer")

	cs := thegraph.NewIndexedChainState(eth.NewChainState(tx, client), graphql.NewClient(graphUrl, nil), logger)
	time.Sleep(5 * time.Second)

	err = cs.Start(ctx)
	require.NoError(t, err, "failed to start indexed chain state")

	headerNum, err := cs.GetCurrentBlockNumber(ctx)
	require.NoError(t, err, "failed to get current block number")

	state, err := cs.GetIndexedOperatorState(ctx, headerNum, testQuorums)
	require.NoError(t, err, "failed to get indexed operator state")
	require.Equal(
		t,
		len(infraHarness.OperatorHarness.ServersV2),
		len(state.IndexedOperators),
		"operator count mismatch",
	)
}
