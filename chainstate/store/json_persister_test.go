package store

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Layr-Labs/eigenda/chainstate/types"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/stretchr/testify/require"
)

func TestPersisterSaveLoadRoundTrip(t *testing.T) {
	ctx := context.Background()
	logger := test.GetLogger()
	path := filepath.Join(t.TempDir(), "state.json")

	s := NewMemoryStore()
	require.NoError(t, s.SaveOperator(ctx, &types.Operator{
		ID:        opID(1),
		Socket:    "host:1",
		QuorumIDs: []core.QuorumID{0},
	}))
	require.NoError(t, s.SetLastIndexedBlock(ctx, 42))

	require.NoError(t, NewJSONPersister(s, path, logger).Save(ctx))

	restored := NewMemoryStore()
	require.NoError(t, NewJSONPersister(restored, path, logger).Load(ctx))

	op, err := restored.GetOperator(ctx, opID(1))
	require.NoError(t, err)
	require.Equal(t, "host:1", op.Socket)

	last, err := restored.GetLastIndexedBlock(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(42), last)
}

func TestPersisterLoadMissingFileIsFreshStart(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "does-not-exist.json")

	s := NewMemoryStore()
	require.NoError(t, NewJSONPersister(s, path, test.GetLogger()).Load(ctx))

	last, err := s.GetLastIndexedBlock(ctx)
	require.NoError(t, err)
	require.Zero(t, last)
}

func TestPersisterLoadCorruptFileErrors(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "state.json")
	require.NoError(t, os.WriteFile(path, []byte("{not json"), 0644))

	s := NewMemoryStore()
	err := NewJSONPersister(s, path, test.GetLogger()).Load(ctx)
	require.Error(t, err, "corrupt state must fail loudly, not half-restore")
}
