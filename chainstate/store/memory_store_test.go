package store

import (
	"context"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/chainstate/types"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func opID(b byte) core.OperatorID {
	var id core.OperatorID
	id[0] = b
	return id
}

func txHash(b byte) common.Hash {
	var h common.Hash
	h[0] = b
	return h
}

func TestSaveEjectionDedupIncludesLogIndex(t *testing.T) {
	ctx := context.Background()
	s := NewMemoryStore()

	// EjectionManager.ejectOperators emits one OperatorEjected per quorum in a
	// single transaction: same operator, block, and tx hash, different log
	// index. Both must be stored.
	base := types.OperatorEjection{
		OperatorID:  opID(1),
		BlockNumber: 100,
		TxHash:      txHash(1),
		EjectedAt:   time.Now(),
	}

	quorum0 := base
	quorum0.QuorumIDs = []core.QuorumID{0}
	quorum0.LogIndex = 3
	require.NoError(t, s.SaveEjection(ctx, &quorum0))

	quorum1 := base
	quorum1.QuorumIDs = []core.QuorumID{1}
	quorum1.LogIndex = 4
	require.NoError(t, s.SaveEjection(ctx, &quorum1))

	ejections, err := s.ListEjections(ctx, nil, 0, 0)
	require.NoError(t, err)
	require.Len(t, ejections, 2, "distinct log indexes must not collapse")

	// Re-indexing the same event (identical identity including log index) must
	// still dedup.
	require.NoError(t, s.SaveEjection(ctx, &quorum0))
	ejections, err = s.ListEjections(ctx, nil, 0, 0)
	require.NoError(t, err)
	require.Len(t, ejections, 2)
}

func TestSaveSocketUpdateDedupIncludesLogIndex(t *testing.T) {
	ctx := context.Background()
	s := NewMemoryStore()

	base := types.OperatorSocketUpdate{
		OperatorID:  opID(1),
		BlockNumber: 100,
		TxHash:      txHash(1),
		UpdatedAt:   time.Now(),
	}

	first := base
	first.Socket = "host:1"
	first.LogIndex = 7
	require.NoError(t, s.SaveSocketUpdate(ctx, &first))

	second := base
	second.Socket = "host:2"
	second.LogIndex = 9
	require.NoError(t, s.SaveSocketUpdate(ctx, &second))

	updates, err := s.ListSocketUpdates(ctx, opID(1), 0, 0)
	require.NoError(t, err)
	require.Len(t, updates, 2, "distinct log indexes must not collapse")

	// Identical identity dedups.
	require.NoError(t, s.SaveSocketUpdate(ctx, &first))
	updates, err = s.ListSocketUpdates(ctx, opID(1), 0, 0)
	require.NoError(t, err)
	require.Len(t, updates, 2)
}

func TestOperatorCloningPreventsAliasing(t *testing.T) {
	ctx := context.Background()
	s := NewMemoryStore()

	op := &types.Operator{
		ID:        opID(1),
		QuorumIDs: []core.QuorumID{0},
	}
	require.NoError(t, s.SaveOperator(ctx, op))

	// Mutating the caller's slice after saving must not reach the store.
	op.QuorumIDs[0] = 99
	stored, err := s.GetOperator(ctx, opID(1))
	require.NoError(t, err)
	require.Equal(t, []core.QuorumID{0}, stored.QuorumIDs)

	// Mutating a returned operator's slice must not reach the store.
	stored.QuorumIDs[0] = 42
	again, err := s.GetOperator(ctx, opID(1))
	require.NoError(t, err)
	require.Equal(t, []core.QuorumID{0}, again.QuorumIDs)

	// Same for ListOperators results.
	listed, err := s.ListOperators(ctx, types.OperatorFilter{}, 0, 0)
	require.NoError(t, err)
	require.Len(t, listed, 1)
	listed[0].QuorumIDs[0] = 42
	again, err = s.GetOperator(ctx, opID(1))
	require.NoError(t, err)
	require.Equal(t, []core.QuorumID{0}, again.QuorumIDs)
}

func TestGetOperatorNotFound(t *testing.T) {
	s := NewMemoryStore()
	_, err := s.GetOperator(context.Background(), opID(1))
	require.ErrorIs(t, err, ErrOperatorNotFound)

	err = s.UpdateOperatorSocket(context.Background(), opID(1), "host:1", 5)
	require.ErrorIs(t, err, ErrOperatorNotFound)

	err = s.DeregisterOperator(context.Background(), opID(1), 5, txHash(1))
	require.ErrorIs(t, err, ErrOperatorNotFound)
}

func TestGetLatestQuorumAPKBoundaries(t *testing.T) {
	ctx := context.Background()
	s := NewMemoryStore()

	for _, block := range []uint64{10, 20} {
		require.NoError(t, s.SaveQuorumAPK(ctx, &types.QuorumAPK{
			QuorumID:    1,
			BlockNumber: block,
		}))
	}

	cases := []struct {
		query uint64
		want  *uint64 // nil means no snapshot expected
	}{
		{query: 5, want: nil},
		{query: 10, want: ptr(uint64(10))},
		{query: 15, want: ptr(uint64(10))},
		{query: 20, want: ptr(uint64(20))},
		{query: 25, want: ptr(uint64(20))},
	}
	for _, tc := range cases {
		apk, err := s.GetLatestQuorumAPK(ctx, 1, tc.query)
		require.NoError(t, err)
		if tc.want == nil {
			require.Nil(t, apk, "query block %d", tc.query)
		} else {
			require.NotNil(t, apk, "query block %d", tc.query)
			require.Equal(t, *tc.want, apk.BlockNumber, "query block %d", tc.query)
		}
	}

	// A quorum with no snapshots returns (nil, nil).
	apk, err := s.GetLatestQuorumAPK(ctx, 2, 100)
	require.NoError(t, err)
	require.Nil(t, apk)
}

func ptr[T any](v T) *T {
	return &v
}

func TestPaginateClamping(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}

	require.Equal(t, []int{1, 2, 3, 4, 5}, paginate(items, 0, 0), "no limit, no offset")
	require.Equal(t, []int{1, 2}, paginate(items, 2, 0))
	require.Equal(t, []int{3, 4}, paginate(items, 2, 2))
	require.Equal(t, []int{5}, paginate(items, 10, 4), "limit past end clamps")
	require.Empty(t, paginate(items, 2, 5), "offset at end")
	require.Empty(t, paginate(items, 2, 99), "offset past end")
	require.Empty(t, paginate(items, 2, -1), "negative offset")
}

func TestSnapshotRestoreRoundTrip(t *testing.T) {
	ctx := context.Background()
	s := NewMemoryStore()

	deregBlock := uint64(80)
	deregTx := txHash(9)
	require.NoError(t, s.SaveOperator(ctx, &types.Operator{
		ID:                        opID(1),
		Address:                   common.HexToAddress("0x1234"),
		Socket:                    "host:1",
		RegisteredAtBlockNumber:   50,
		RegisteredTxHash:          txHash(2),
		DeregisteredAtBlockNumber: &deregBlock,
		DeregisteredTxHash:        &deregTx,
		QuorumIDs:                 []core.QuorumID{0, 1},
	}))
	require.NoError(t, s.SaveQuorumAPK(ctx, &types.QuorumAPK{
		QuorumID:    1,
		BlockNumber: 60,
		UpdatedAt:   time.Unix(1000, 0).UTC(),
	}))
	ejection := &types.OperatorEjection{
		OperatorID:  opID(1),
		QuorumIDs:   []core.QuorumID{1},
		BlockNumber: 70,
		LogIndex:    4,
		TxHash:      txHash(3),
		EjectedAt:   time.Unix(2000, 0).UTC(),
	}
	require.NoError(t, s.SaveEjection(ctx, ejection))
	update := &types.OperatorSocketUpdate{
		OperatorID:  opID(1),
		Socket:      "host:1",
		BlockNumber: 50,
		LogIndex:    2,
		TxHash:      txHash(2),
		UpdatedAt:   time.Unix(3000, 0).UTC(),
	}
	require.NoError(t, s.SaveSocketUpdate(ctx, update))
	require.NoError(t, s.SetLastIndexedBlock(ctx, 100))

	data, err := s.Snapshot()
	require.NoError(t, err)

	restored := NewMemoryStore()
	require.NoError(t, restored.Restore(data))

	op, err := restored.GetOperator(ctx, opID(1))
	require.NoError(t, err)
	require.Equal(t, "host:1", op.Socket)
	require.Equal(t, []core.QuorumID{0, 1}, op.QuorumIDs)
	require.NotNil(t, op.DeregisteredAtBlockNumber)
	require.Equal(t, uint64(80), *op.DeregisteredAtBlockNumber)

	last, err := restored.GetLastIndexedBlock(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(100), last)

	// The APK block index must be rebuilt so lookups still work.
	apk, err := restored.GetLatestQuorumAPK(ctx, 1, 65)
	require.NoError(t, err)
	require.NotNil(t, apk)
	require.Equal(t, uint64(60), apk.BlockNumber)

	ejections, err := restored.ListEjections(ctx, nil, 0, 0)
	require.NoError(t, err)
	require.Len(t, ejections, 1)
	require.Equal(t, uint(4), ejections[0].LogIndex)

	// The dedup key sets are rebuilt from the records (not persisted), so
	// re-saving a previously seen event must still be a no-op.
	require.NoError(t, restored.SaveEjection(ctx, ejection))
	ejections, err = restored.ListEjections(ctx, nil, 0, 0)
	require.NoError(t, err)
	require.Len(t, ejections, 1)

	require.NoError(t, restored.SaveSocketUpdate(ctx, update))
	updates, err := restored.ListSocketUpdates(ctx, opID(1), 0, 0)
	require.NoError(t, err)
	require.Len(t, updates, 1)
	require.Equal(t, uint(2), updates[0].LogIndex)
}

// TestSnapshotConcurrentWithWrites exercises Snapshot racing writer methods;
// run with -race to catch aliasing between the snapshot copy and live state.
func TestSnapshotConcurrentWithWrites(t *testing.T) {
	ctx := context.Background()
	s := NewMemoryStore()
	require.NoError(t, s.SaveOperator(ctx, &types.Operator{
		ID:        opID(1),
		QuorumIDs: []core.QuorumID{0},
	}))

	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < 100; i++ {
			op, err := s.GetOperator(ctx, opID(1))
			if err != nil {
				return
			}
			op.QuorumIDs = append(op.QuorumIDs, core.QuorumID(i%256))
			if err := s.SaveOperator(ctx, op); err != nil {
				return
			}
		}
	}()

	for i := 0; i < 100; i++ {
		_, err := s.Snapshot()
		require.NoError(t, err)
	}
	<-done
}
