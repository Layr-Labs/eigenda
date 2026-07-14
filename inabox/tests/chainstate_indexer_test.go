package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/chainstate"
	"github.com/Layr-Labs/eigenda/chainstate/service"
	"github.com/Layr-Labs/eigenda/chainstate/types"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/stretchr/testify/require"
)

// listOperatorsResponse mirrors the JSON shape returned by GET /api/v1/operators.
type listOperatorsResponse struct {
	Operators []types.Operator `json:"operators"`
	Count     int              `json:"count"`
}

// listQuorumAPKsResponse mirrors the JSON shape returned by GET /api/v1/quorum-apk/history.
type listQuorumAPKsResponse struct {
	QuorumAPKs []types.QuorumAPK `json:"quorum_apks"`
	Count      int               `json:"count"`
}

// listSocketUpdatesResponse mirrors the JSON shape returned by
// GET /api/v1/socket-updates/:operator_id.
type listSocketUpdatesResponse struct {
	SocketUpdates []types.OperatorSocketUpdate `json:"socket_updates"`
	Count         int                          `json:"count"`
}

// statusResponse mirrors the JSON shape returned by GET /api/v1/status.
type statusResponse struct {
	LastIndexedBlock uint64 `json:"last_indexed_block"`
}

// TestChainStateIndexerE2E runs the chainstate indexer against the live inabox
// devnet and verifies, through its REST API, that it indexes the operators that
// inabox registers on-chain at startup. It also exercises the graceful-shutdown
// path by asserting that a final state snapshot is persisted on exit.
//
// NOTE (per CLAUDE.md §3): this test was AI-drafted and must be reviewed for
// whether its assertions encode intended behavior, not merely current behavior.
func TestChainStateIndexerE2E(t *testing.T) {
	// Reuse the suite-wide inabox infrastructure (chain + deployed contracts +
	// registered operators). No per-test harness setup is required.
	require.NotNil(t, globalInfra, "global infrastructure must be initialized by TestMain")
	require.NotEmpty(t, globalInfra.TestConfig.Deployers, "expected at least one deployer")

	rpcURL := globalInfra.TestConfig.Deployers[0].RPC
	require.NotEmpty(t, rpcURL, "deployer RPC URL must be set")
	eigenDADirectory := globalInfra.TestConfig.EigenDA.EigenDADirectory
	require.NotEmpty(t, eigenDADirectory, "EigenDADirectory address must be set")

	// The number of operators inabox registers on-chain at startup.
	expectedOperators := globalInfra.TestConfig.Services.Counts.NumOpr
	require.Greater(t, expectedOperators, 0, "expected inabox to register at least one operator")

	logger := test.GetLogger()

	cfg := &chainstate.IndexerConfig{
		EigenDADirectory: eigenDADirectory,
		// Index from genesis so we capture the operator registrations that
		// happened during inabox startup. StartBlockNumber 0 would start from
		// the *current* block and miss them entirely.
		StartBlockNumber: 1,
		BlockBatchSize:   1000,
		// Poll aggressively so the test converges quickly.
		PollInterval:    500 * time.Millisecond,
		PersistInterval: 1 * time.Second,
		PersistencePath: fmt.Sprintf("%s/chainstate.json", t.TempDir()),
		HTTPPort:        freePort(t),
		LoggerConfig:    *common.DefaultLoggerConfig(),
	}
	secret := &chainstate.IndexerSecretConfig{EthRpcUrls: []string{rpcURL}}

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	svc, err := service.New(ctx, cfg, secret, logger)
	require.NoError(t, err, "failed to create chainstate service")

	require.NoError(t, svc.Start(ctx), "failed to start chainstate service")

	baseURL := fmt.Sprintf("http://localhost:%s/api/v1", cfg.HTTPPort)

	// Surface a fatal API server error (e.g. port already in use) instead of
	// timing out opaquely in the polling loop below.
	select {
	case err := <-svc.Errors():
		require.NoError(t, err, "API server failed to start")
	case <-time.After(500 * time.Millisecond):
	}

	// Poll until the indexer has fully processed the operator registrations.
	//
	// We can't just wait on operator count: within a single indexing batch the
	// indexer first saves operators (from RegistryCoordinator events, with empty
	// quorum membership) and only afterwards populates QuorumIDs (from the
	// BLSApkRegistry OperatorAddedToQuorums events). Polling on count alone races
	// against that second step. Instead we wait until all registered operators
	// carry quorum membership, which only holds once the batch is fully applied.
	var operators, registered []types.Operator
	require.Eventually(t, func() bool {
		var all listOperatorsResponse
		if err := httpGetJSON(ctx, baseURL+"/operators?limit=1000", &all); err != nil {
			return false
		}
		var reg listOperatorsResponse
		if err := httpGetJSON(ctx, baseURL+"/operators?registered=true&limit=1000", &reg); err != nil {
			return false
		}
		if len(all.Operators) < expectedOperators || len(reg.Operators) == 0 {
			return false
		}
		for _, op := range reg.Operators {
			if len(op.QuorumIDs) == 0 {
				return false
			}
		}
		operators = all.Operators
		registered = reg.Operators
		return true
	}, 60*time.Second, 500*time.Millisecond,
		"indexer did not fully index %d operators in time", expectedOperators)

	require.Len(t, operators, expectedOperators, "unexpected operator count")

	// inabox configures more operators than maxOperatorCount, so churn can
	// deregister some at startup. Every indexed operator should at least have a
	// socket (set at registration) and a valid ID, regardless of churn.
	for _, op := range operators {
		require.NotEqual(t, types.Operator{}.ID, op.ID, "operator ID must be set")
		require.NotEmpty(t, op.Socket, "operator socket must be set")
	}

	// The registered operators are the ones that survived churn; each must
	// belong to at least one quorum (guaranteed by the poll condition above).
	for _, op := range registered {
		require.True(t, op.IsRegistered(), "registered filter returned a deregistered operator")
		require.NotEmpty(t, op.QuorumIDs, "registered operator must belong to at least one quorum")
	}

	// GET /operators/:id should return a single operator matching the list entry.
	first := registered[0]
	var fetched types.Operator
	require.NoError(t, httpGetJSON(ctx, baseURL+"/operators/"+first.ID.Hex(), &fetched))
	require.Equal(t, first.ID, fetched.ID)
	require.Equal(t, first.Address, fetched.Address)
	require.Equal(t, first.Socket, fetched.Socket)

	// Each registered operator emits a socket-update event, so the socket-update
	// history for our sampled operator should be non-empty.
	var socketUpdates listSocketUpdatesResponse
	require.NoError(t, httpGetJSON(ctx, baseURL+"/socket-updates/"+first.ID.Hex(), &socketUpdates))
	require.NotEmpty(t, socketUpdates.SocketUpdates, "expected at least one socket update for operator")
	require.Equal(t, first.ID, socketUpdates.SocketUpdates[0].OperatorID)

	// Registering/deregistering operators changes each affected quorum's
	// aggregate public key, so every quorum our registered operators belong to
	// should have at least one APK snapshot carrying an aggregate key.
	quorums := map[uint8]struct{}{}
	for _, op := range registered {
		for _, q := range op.QuorumIDs {
			quorums[q] = struct{}{}
		}
	}
	require.NotEmpty(t, quorums, "registered operators should belong to at least one quorum")
	for q := range quorums {
		var apks listQuorumAPKsResponse
		require.NoError(t, httpGetJSON(ctx,
			fmt.Sprintf("%s/quorum-apk/history?quorum_id=%d", baseURL, q), &apks))
		require.NotEmpty(t, apks.QuorumAPKs, "expected APK history for quorum %d", q)
		latest := apks.QuorumAPKs[len(apks.QuorumAPKs)-1]
		require.Equal(t, q, latest.QuorumID, "APK snapshot quorum mismatch")
		require.NotNil(t, latest.APK, "APK snapshot must include an aggregate key")
	}

	// The status endpoint should report a non-zero last-indexed block.
	var status statusResponse
	require.NoError(t, httpGetJSON(ctx, baseURL+"/status", &status))
	require.Greater(t, status.LastIndexedBlock, uint64(0), "last indexed block should advance")

	// Exercise graceful shutdown: cancelling the context must stop the indexer
	// and flush a final state snapshot to disk before Wait returns.
	cancel()
	svc.Wait()

	persisted := readPersistedSnapshot(t, cfg.PersistencePath)
	require.Len(t, persisted.Operators, expectedOperators,
		"persisted snapshot should contain all indexed operators")
}

// persistedSnapshot is a minimal view of the JSON snapshot written by the
// persister, sufficient to assert that operators were saved.
type persistedSnapshot struct {
	Operators map[string]types.Operator `json:"operators"`
}

func readPersistedSnapshot(t *testing.T, path string) persistedSnapshot {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err, "reading persisted snapshot")
	var snap persistedSnapshot
	require.NoError(t, json.Unmarshal(data, &snap), "unmarshalling persisted snapshot")
	return snap
}

// httpGetJSON performs a GET request against url and decodes the JSON body into
// out. It returns an error on any non-200 response or transport failure.
func httpGetJSON(ctx context.Context, url string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("performing request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d for %s", resp.StatusCode, url)
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}
	return nil
}

// freePort asks the kernel for an available TCP port and returns it as a string.
func freePort(t *testing.T) string {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err, "failed to find a free port")
	defer func() { _ = lis.Close() }()
	_, port, err := net.SplitHostPort(lis.Addr().String())
	require.NoError(t, err, "failed to parse free port")
	return port
}
