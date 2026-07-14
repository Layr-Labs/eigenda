package chainstate

import (
	"context"
	"testing"

	"github.com/Layr-Labs/eigenda/chainstate/store"
	"github.com/Layr-Labs/eigenda/chainstate/types"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/stretchr/testify/require"
)

func newTestICS(t *testing.T) (*IndexedChainState, store.Store) {
	t.Helper()
	memStore := store.NewMemoryStore()
	ics := &IndexedChainState{
		store:  memStore,
		logger: test.GetLogger(),
	}
	return ics, memStore
}

func completeOperator(t *testing.T, id byte) *types.Operator {
	t.Helper()
	keyPair, err := core.GenRandomBlsKeys()
	require.NoError(t, err)
	return &types.Operator{
		ID:                      core.OperatorID{id},
		BLSPubKeyG1:             keyPair.GetPubKeyG1(),
		BLSPubKeyG2:             keyPair.GetPubKeyG2(),
		Socket:                  "host:1",
		RegisteredAtBlockNumber: 10,
	}
}

func TestGetIndexedOperatorsSkipsIncompleteRecords(t *testing.T) {
	ctx := context.Background()
	ics, memStore := newTestICS(t)

	complete := completeOperator(t, 1)
	require.NoError(t, memStore.SaveOperator(ctx, complete))

	// A record missing BLS keys (e.g. pubkey registration predating the
	// configured start block) must not poison the whole query.
	require.NoError(t, memStore.SaveOperator(ctx, &types.Operator{
		ID:                      core.OperatorID{2},
		Socket:                  "host:2",
		RegisteredAtBlockNumber: 20,
	}))

	// Same for a record missing its socket.
	incomplete := completeOperator(t, 3)
	incomplete.Socket = ""
	require.NoError(t, memStore.SaveOperator(ctx, incomplete))

	operators, err := ics.GetIndexedOperators(ctx, 100)
	require.NoError(t, err, "incomplete records must be skipped, not fail the call")
	require.Len(t, operators, 1)
	require.Contains(t, operators, complete.ID)
}

func TestGetIndexedOperatorsDeregisteredBoundary(t *testing.T) {
	ctx := context.Background()
	ics, memStore := newTestICS(t)

	op := completeOperator(t, 1)
	deregBlock := uint64(100)
	op.DeregisteredAtBlockNumber = &deregBlock
	require.NoError(t, memStore.SaveOperator(ctx, op))

	// Deregistered at block 100: excluded at 100 (dereg block inclusive,
	// mirroring the subgraph's deregistrationBlockNumber_gt filter), included
	// at 99.
	operators, err := ics.GetIndexedOperators(ctx, 100)
	require.NoError(t, err)
	require.Empty(t, operators)

	operators, err = ics.GetIndexedOperators(ctx, 99)
	require.NoError(t, err)
	require.Len(t, operators, 1)
}

func TestGetIndexedOperatorInfoByOperatorId(t *testing.T) {
	ctx := context.Background()
	ics, memStore := newTestICS(t)

	op := completeOperator(t, 1)
	require.NoError(t, memStore.SaveOperator(ctx, op))

	info, err := ics.GetIndexedOperatorInfoByOperatorId(ctx, op.ID, 50)
	require.NoError(t, err)
	require.Equal(t, "host:1", info.Socket)

	_, err = ics.GetIndexedOperatorInfoByOperatorId(ctx, core.OperatorID{9}, 50)
	require.Error(t, err)
}

func TestGetQuorumAPKLatestSnapshotSemantics(t *testing.T) {
	ctx := context.Background()
	ics, memStore := newTestICS(t)

	keyPair, err := core.GenRandomBlsKeys()
	require.NoError(t, err)
	require.NoError(t, memStore.SaveQuorumAPK(ctx, &types.QuorumAPK{
		QuorumID:    0,
		BlockNumber: 50,
		APK:         keyPair.GetPubKeyG1(),
	}))

	// Query above the snapshot block returns the latest snapshot <= block.
	apk, err := ics.getQuorumAPK(ctx, 0, 75)
	require.NoError(t, err)
	require.NotNil(t, apk)

	// Query below the first snapshot yields nil (surfaced as an error by
	// GetIndexedOperatorState, not a silent partial key set).
	apk, err = ics.getQuorumAPK(ctx, 0, 25)
	require.NoError(t, err)
	require.Nil(t, apk)
}
