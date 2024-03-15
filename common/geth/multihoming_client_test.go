package geth_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/common/geth"
	damock "github.com/Layr-Labs/eigenda/common/mock"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	privateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	rpcURLs    = []string{"127.0.0.1", "127.0.0.2", "127.0.0.3"}
)

func makeTestMultihomingClient(t *testing.T, numRetries int) (*geth.MultiHomingClient, error) {
	logger := logging.NewNoopLogger()

	ethClientCfg := geth.EthClientConfig{
		RPCURLs:          rpcURLs,
		PrivateKeyString: privateKey,
		NumConfirmations: 0,
		NumRetries:       numRetries,
	}

	mockClient := geth.MultiHomingClient{}
	controller := geth.NewFailoverController(len(rpcURLs), logger)

	//mockClient.rpcUrls = rpcURLs
	mockClient.Logger = logger
	mockClient.NumRetries = ethClientCfg.NumRetries
	mockClient.FailoverController = controller

	for i := 0; i < len(rpcURLs); i++ {
		mockEthClient := &damock.MockEthClient{}
		mockEthClient.On("ChainID", mock.Anything).Return(big.NewInt(0), ethereum.NotFound)
		mockClient.RPCs = append(mockClient.RPCs, mockEthClient)
	}

	return &mockClient, nil
}

func makeFailureCall(t *testing.T, client *geth.MultiHomingClient, numCall int) {
	for i := 0; i < numCall; i++ {
		ctx := context.Background()
		_, err := client.ChainID(ctx)
		require.NotNil(t, err)
	}
}

func TestMultihomingClientZeroRetry(t *testing.T) {
	client, _ := makeTestMultihomingClient(t, 0)

	index, _ := client.GetRPCInstance()
	require.Equal(t, index, 0)

	makeFailureCall(t, client, 1)

	// given num retry is 0, when failure arises above, current rpc should becomes the next one
	index, _ = client.GetRPCInstance()
	require.Equal(t, index, 1)

	makeFailureCall(t, client, 1)

	index, _ = client.GetRPCInstance()
	require.Equal(t, index, 2)

	makeFailureCall(t, client, 1)

	index, _ = client.GetRPCInstance()
	require.Equal(t, index, 0)
}

func TestMultihomingClientONERetry(t *testing.T) {
	client, _ := makeTestMultihomingClient(t, 1)

	index, _ := client.GetRPCInstance()
	require.Equal(t, index, 0)

	makeFailureCall(t, client, 1)

	// given num retry is 1, when failure arises above, two rpc are used, current rpc should becomes 2
	index, _ = client.GetRPCInstance()
	require.Equal(t, index, 2)

	makeFailureCall(t, client, 1)

	index, _ = client.GetRPCInstance()
	require.Equal(t, index, 1)

	makeFailureCall(t, client, 1)

	index, _ = client.GetRPCInstance()
	require.Equal(t, index, 0)
}

func TestMultihomingClientTWORetry(t *testing.T) {
	client, _ := makeTestMultihomingClient(t, 2)

	index, _ := client.GetRPCInstance()
	require.Equal(t, index, 0)

	makeFailureCall(t, client, 1)

	// given num retry is 2, when failure arises above, three rpc are used, current rpc should becomes 0
	index, _ = client.GetRPCInstance()
	require.Equal(t, index, 0)

	makeFailureCall(t, client, 1)

	index, _ = client.GetRPCInstance()
	require.Equal(t, index, 0)

	makeFailureCall(t, client, 1)

	index, _ = client.GetRPCInstance()
	require.Equal(t, index, 0)
}
