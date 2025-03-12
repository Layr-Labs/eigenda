package controller_test

import (
	"testing"

	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/controller"
	"github.com/stretchr/testify/require"
)

func TestBlobSetOperations(t *testing.T) {
	q := controller.NewBlobSet()
	require.Equal(t, 0, q.Size())

	key0 := v2.BlobKey([32]byte{0})
	q.AddBlob(key0)
	require.Equal(t, 1, q.Size())
	require.True(t, q.Contains(key0))

	key1 := v2.BlobKey([32]byte{1})
	q.AddBlob(key1)
	require.Equal(t, 2, q.Size())
	require.True(t, q.Contains(key1))

	q.RemoveBlob(key0)
	require.Equal(t, 1, q.Size())
	require.False(t, q.Contains(key0))
	require.True(t, q.Contains(key1))

	q.RemoveBlob(key1)
	require.Equal(t, 0, q.Size())
	require.False(t, q.Contains(key0))
	require.False(t, q.Contains(key1))

	q.RemoveBlob(key1)
	require.Equal(t, 0, q.Size())
	require.False(t, q.Contains(key0))

	q.AddBlob(key0)
	require.Equal(t, 1, q.Size())
	require.True(t, q.Contains(key0))
}
