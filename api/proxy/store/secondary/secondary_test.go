package secondary_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/metrics"
	"github.com/Layr-Labs/eigenda/api/proxy/store/secondary"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/stretchr/testify/require"
)

func testLogger() logging.Logger {
	return logging.NewTextSLogger(os.Stdout, &logging.SLoggerOptions{})
}

func testMetrics() metrics.Metricer {
	return metrics.NewMetrics(nil)
}

// mockStore implements common.SecondaryStore for testing
type mockStore struct {
	shouldFail bool
	backend    common.BackendType
}

func (m *mockStore) Get(ctx context.Context, key []byte) ([]byte, error) {
	if m.shouldFail {
		return nil, errors.New("mock error")
	}
	return []byte("data"), nil
}

func (m *mockStore) Put(ctx context.Context, key []byte, value []byte) error {
	if m.shouldFail {
		return errors.New("mock error")
	}
	return nil
}

func (m *mockStore) Verify(ctx context.Context, key []byte, value []byte) error {
	return nil
}

func (m *mockStore) BackendType() common.BackendType {
	return m.backend
}

func TestErrorOnInsertFailure(t *testing.T) {
	log := testLogger()
	m := testMetrics()

	t.Run("returns true when enabled", func(t *testing.T) {
		sm := secondary.NewSecondaryManager(log, m, nil, nil, false, true)
		require.True(t, sm.ErrorOnInsertFailure())
	})

	t.Run("returns false when disabled", func(t *testing.T) {
		sm := secondary.NewSecondaryManager(log, m, nil, nil, false, false)
		require.False(t, sm.ErrorOnInsertFailure())
	})
}

func TestHandleRedundantWrites_ErrorOnInsertFailureOFF(t *testing.T) {
	log := testLogger()
	m := testMetrics()
	ctx := context.Background()

	t.Run("succeeds when all writes succeed", func(t *testing.T) {
		store1 := &mockStore{backend: common.S3BackendType}
		sm := secondary.NewSecondaryManager(log, m, []common.SecondaryStore{store1}, nil, false, false)

		err := sm.HandleRedundantWrites(ctx, []byte("commit"), []byte("value"))
		require.NoError(t, err)
	})

	t.Run("succeeds when partial failure and flag OFF", func(t *testing.T) {
		store1 := &mockStore{shouldFail: true, backend: common.S3BackendType}
		store2 := &mockStore{shouldFail: false, backend: common.S3BackendType}
		sm := secondary.NewSecondaryManager(log, m,
			[]common.SecondaryStore{store1},
			[]common.SecondaryStore{store2},
			false, false)

		err := sm.HandleRedundantWrites(ctx, []byte("commit"), []byte("value"))
		require.NoError(t, err, "should not error when flag OFF and at least one write succeeds")
	})

	t.Run("errors when all writes fail", func(t *testing.T) {
		store1 := &mockStore{shouldFail: true, backend: common.S3BackendType}
		sm := secondary.NewSecondaryManager(log, m, []common.SecondaryStore{store1}, nil, false, false)

		err := sm.HandleRedundantWrites(ctx, []byte("commit"), []byte("value"))
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to write blob to any redundant targets")
	})
}

func TestHandleRedundantWrites_ErrorOnInsertFailureON(t *testing.T) {
	log := testLogger()
	m := testMetrics()
	ctx := context.Background()

	t.Run("succeeds when all writes succeed", func(t *testing.T) {
		store1 := &mockStore{backend: common.S3BackendType}
		sm := secondary.NewSecondaryManager(log, m, []common.SecondaryStore{store1}, nil, false, true)

		err := sm.HandleRedundantWrites(ctx, []byte("commit"), []byte("value"))
		require.NoError(t, err)
	})

	t.Run("errors on partial failure when flag ON", func(t *testing.T) {
		store1 := &mockStore{shouldFail: true, backend: common.S3BackendType}
		store2 := &mockStore{shouldFail: false, backend: common.S3BackendType}
		sm := secondary.NewSecondaryManager(log, m,
			[]common.SecondaryStore{store1},
			[]common.SecondaryStore{store2},
			false, true)

		err := sm.HandleRedundantWrites(ctx, []byte("commit"), []byte("value"))
		require.Error(t, err, "should error when flag ON and any write fails")
		require.Contains(t, err.Error(), "failed to write to 1 of 2 secondary targets")
		require.Contains(t, err.Error(), "error-on-secondary-insert-failure=true")
	})

	t.Run("errors when all writes fail", func(t *testing.T) {
		store1 := &mockStore{shouldFail: true, backend: common.S3BackendType}
		sm := secondary.NewSecondaryManager(log, m, []common.SecondaryStore{store1}, nil, false, true)

		err := sm.HandleRedundantWrites(ctx, []byte("commit"), []byte("value"))
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to write blob to any redundant targets")
	})
}
