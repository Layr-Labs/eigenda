package integration_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/chainstate"
	"github.com/Layr-Labs/eigenda/common"
	coreeth "github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	integration "github.com/Layr-Labs/eigenda/inabox/tests"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

// dialEthClient dials the given RPC URL and registers cleanup to close it.
func dialEthClient(t *testing.T, rpcURL string) *ethclient.Client {
	t.Helper()
	client, err := ethclient.Dial(rpcURL)
	require.NoError(t, err, "failed to dial Ethereum RPC")
	t.Cleanup(client.Close)
	return client
}

// TestChainStateSubgraphParity asserts that the chainstate indexer's
// core.IndexedChainState implementation returns the same indexed operator state
// as the operator-state subgraph it is intended to replace. Both are run
// against the same inabox devnet and share the same on-chain ChainState, so the
// comparison isolates the indexed data (aggregate public keys and per-operator
// BLS keys / sockets).
//
// The subgraph is the reference: this is a conformance test, not a symmetric
// one. Fields the subgraph does not track (e.g. QuorumAPK.TotalStake) are out of
// scope and not compared.
//
// NOTE (per CLAUDE.md §3): AI-drafted. The assertions encode assumptions about
// how the two implementations should agree (notably the "latest snapshot <=
// block" APK semantics); review before relying on them.
func TestChainStateSubgraphParity(t *testing.T) {
	require.NotNil(t, globalInfra, "global infrastructure must be initialized by TestMain")
	require.NotNil(t, globalInfra.ChainHarness.GraphNode, "subgraph must be deployed (deploySubgraphs: true)")
	require.NotEmpty(t, globalInfra.TestConfig.Deployers, "expected at least one deployer")

	logger := test.GetLogger()
	ctx := t.Context()

	rpcURL := globalInfra.TestConfig.Deployers[0].RPC
	require.NotEmpty(t, rpcURL, "deployer RPC URL must be set")
	eigenDADirectory := globalInfra.TestConfig.EigenDA.EigenDADirectory
	require.NotEmpty(t, eigenDADirectory, "EigenDADirectory address must be set")

	// A per-test harness gives us a ready-built on-chain reader and eth client.
	harness, err := integration.NewTestHarnessWithSetup(globalInfra)
	require.NoError(t, err, "failed to create test harness")
	defer harness.Cleanup()

	// The on-chain ChainState is shared by both implementations, so the on-chain
	// half of the interface is identical by construction and only the indexed
	// half is under test.
	onchainChainState := coreeth.NewChainState(harness.ChainReader, harness.EthClient)

	// Build the subgraph-backed reference implementation.
	subgraphURL := globalInfra.ChainHarness.GraphNode.HTTPURL() +
		"/subgraphs/name/Layr-Labs/eigenda-operator-state"
	subgraphICS := thegraph.MakeIndexedChainState(
		thegraph.Config{Endpoint: subgraphURL, PullInterval: 100 * time.Millisecond, MaxRetries: 5},
		onchainChainState,
		logger,
	)
	require.NoError(t, subgraphICS.Start(ctx), "failed to start subgraph indexed chain state")

	// Build the chainstate-backed implementation under test.
	cfg := &chainstate.IndexerConfig{
		EigenDADirectory: eigenDADirectory,
		StartBlockNumber: 1,
		BlockBatchSize:   1000,
		PollInterval:     500 * time.Millisecond,
		PersistInterval:  1 * time.Second,
		PersistencePath:  fmt.Sprintf("%s/chainstate.json", t.TempDir()),
		HTTPPort:         freePort(t),
		LoggerConfig:     *common.DefaultLoggerConfig(),
	}
	indexer, err := chainstate.NewIndexer(ctx, cfg, dialEthClient(t, rpcURL), logger)
	require.NoError(t, err, "failed to create indexer")

	indexerCtx, cancelIndexer := context.WithCancel(ctx)
	chainstateICS := chainstate.NewIndexedChainState(indexer, onchainChainState, logger)
	require.NoError(t, chainstateICS.Start(indexerCtx), "failed to start chainstate indexed chain state")
	defer func() {
		cancelIndexer()
		indexer.Wait()
	}()

	// Choose a reference block that both sides have indexed. Use the current
	// chain head, then wait for the chainstate indexer to catch up to it. The
	// subgraph is already consumed by the running inabox services, so it is
	// caught up; we additionally gate on it having APK data below.
	referenceBlock, err := onchainChainState.GetCurrentBlockNumber(ctx)
	require.NoError(t, err, "failed to get current block number")

	require.Eventually(t, func() bool {
		last, err := indexer.GetStore().GetLastIndexedBlock(ctx)
		return err == nil && last >= uint64(referenceBlock)
	}, 60*time.Second, 500*time.Millisecond, "chainstate indexer did not catch up to block %d", referenceBlock)

	// Determine the quorums to compare from the on-chain operator state at the
	// reference block.
	quorums := activeQuorums(t, ctx, onchainChainState, referenceBlock)
	require.NotEmpty(t, quorums, "expected at least one active quorum at the reference block")

	// The subgraph indexes asynchronously; wait until it has APK data for every
	// quorum before comparing, so we do not race a stale subgraph.
	require.Eventually(t, func() bool {
		state, err := subgraphICS.GetIndexedOperatorState(ctx, referenceBlock, quorums)
		return err == nil && len(state.AggKeys) == len(quorums)
	}, 60*time.Second, time.Second, "subgraph did not index APKs for all quorums at block %d", referenceBlock)

	subgraphState, err := subgraphICS.GetIndexedOperatorState(ctx, referenceBlock, quorums)
	require.NoError(t, err, "subgraph GetIndexedOperatorState failed")

	chainstateState, err := chainstateICS.GetIndexedOperatorState(ctx, referenceBlock, quorums)
	require.NoError(t, err, "chainstate GetIndexedOperatorState failed")

	// Compare aggregate public keys per quorum.
	require.Equal(t, len(subgraphState.AggKeys), len(chainstateState.AggKeys),
		"different number of quorum aggregate public keys")
	for quorum, subgraphAPK := range subgraphState.AggKeys {
		chainstateAPK, ok := chainstateState.AggKeys[quorum]
		require.True(t, ok, "chainstate missing aggregate public key for quorum %d", quorum)
		require.True(t, subgraphAPK.G1Affine.Equal(chainstateAPK.G1Affine),
			"aggregate public key mismatch for quorum %d", quorum)
	}

	// Compare the indexed operator sets (BLS keys and socket) keyed by ID.
	require.Equal(t, len(subgraphState.IndexedOperators), len(chainstateState.IndexedOperators),
		"different number of indexed operators")
	for operatorID, subgraphOp := range subgraphState.IndexedOperators {
		chainstateOp, ok := chainstateState.IndexedOperators[operatorID]
		require.True(t, ok, "chainstate missing indexed operator %s", operatorID.Hex())
		require.Equal(t, subgraphOp.Socket, chainstateOp.Socket,
			"socket mismatch for operator %s", operatorID.Hex())
		require.True(t, subgraphOp.PubkeyG1.G1Affine.Equal(chainstateOp.PubkeyG1.G1Affine),
			"pubkey G1 mismatch for operator %s", operatorID.Hex())
		require.True(t, subgraphOp.PubkeyG2.G2Affine.Equal(chainstateOp.PubkeyG2.G2Affine),
			"pubkey G2 mismatch for operator %s", operatorID.Hex())
	}
}

// activeQuorums returns the quorum IDs that have at least one operator in the
// on-chain operator state at the given block. It probes the standard quorums
// (0 and 1) used by the inabox devnet.
func activeQuorums(
	t *testing.T,
	ctx context.Context,
	chainState *coreeth.ChainState,
	blockNumber uint,
) []uint8 {
	t.Helper()
	candidates := []uint8{0, 1}
	state, err := chainState.GetOperatorState(ctx, blockNumber, candidates)
	require.NoError(t, err, "failed to get on-chain operator state")

	var quorums []uint8
	for _, q := range candidates {
		if len(state.Operators[q]) > 0 {
			quorums = append(quorums, q)
		}
	}
	return quorums
}
