package controller_test

import (
	"github.com/Layr-Labs/eigenda/api/clients/mock"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"testing"

	"github.com/Layr-Labs/eigenda/disperser/controller"
	"github.com/stretchr/testify/require"
)

func TestNodeClientManager(t *testing.T) {
	rand := random.NewTestRandom(t)

	_, private := rand.ECDSA()
	requestSigner := mock.NewStaticRequestSigner(private)

	m, err := controller.NewNodeClientManager(2, requestSigner, nil)
	require.NoError(t, err)

	client0, err := m.GetClient("localhost", "0000")
	require.NoError(t, err)
	require.NotNil(t, client0)

	client1, err := m.GetClient("localhost", "0000")
	require.NoError(t, err)
	require.NotNil(t, client1)

	require.Same(t, client0, client1)

	// fill up the cache
	client2, err := m.GetClient("localhost", "1111")
	require.NoError(t, err)
	require.NotNil(t, client2)

	// evict client0
	client3, err := m.GetClient("localhost", "2222")
	require.NoError(t, err)
	require.NotNil(t, client3)

	// accessing client0 again should create new client
	client4, err := m.GetClient("localhost", "0000")
	require.NoError(t, err)
	require.NotNil(t, client0)

	require.NotSame(t, client0, client4)
}
