package geth_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/common/geth"
	damock "github.com/Layr-Labs/eigenda/common/mock"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	privateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	rpcURLs    = []string{"http://127.0.0.1/abcd", "https://www.da:9000/abcd", "https://a-b-c.A.B.C/dddd"}
)

type JsonError struct{}

func (j *JsonError) Error() string  { return "json error" }
func (j *JsonError) ErrorCode() int { return -32000 }

func makeTestMultihomingClient(numRetries int, designatedError error) (*geth.MultiHomingClient, error) {
	logger := logging.NewNoopLogger()

	ethClientCfg := geth.EthClientConfig{
		RPCURLs:          rpcURLs,
		PrivateKeyString: privateKey,
		NumConfirmations: 0,
		NumRetries:       numRetries,
	}

	mockClient := geth.MultiHomingClient{}
	controller, err := geth.NewFailoverController(logger, rpcURLs)
	if err != nil {
		return nil, err
	}

	mockClient.Logger = logger
	mockClient.NumRetries = ethClientCfg.NumRetries
	mockClient.FailoverController = controller

	for i := 0; i < len(rpcURLs); i++ {
		mockEthClient := &damock.MockEthClient{}
		mockEthClient.On("ChainID", mock.Anything).Return(big.NewInt(0), designatedError)
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

func make500Error() error {
	return rpc.HTTPError{
		StatusCode: 500,
		Status:     "INTERNAL_ERROR",
		Body:       []byte{},
	}
}

func TestMultihomingClient_UrlDomain(t *testing.T) {
	client, err := makeTestMultihomingClient(2, nil)
	require.Nil(t, err)
	urlDomains := client.FailoverController.UrlDomains
	fmt.Println("urlDomains", urlDomains)
	require.Equal(t, urlDomains[0], "127.0.0.1")
	require.Equal(t, urlDomains[1], "www.da")
	require.Equal(t, urlDomains[2], "a-b-c.A.B.C")
}

func TestMultihomingClientSenderFaultZeroRetry(t *testing.T) {
	// 4xx attributes to sender's fault, RPC should not rotate
	statusCodes := []int{401, 499}
	for _, sc := range statusCodes {

		httpRespError := rpc.HTTPError{
			StatusCode: sc,
			Status:     "INTERNAL_ERROR",
			Body:       []byte{},
		}

		client, _ := makeTestMultihomingClient(0, httpRespError)

		index, _ := client.GetRPCInstance()
		require.Equal(t, index, 0)

		makeFailureCall(t, client, 10)

		// given error is 401, 409, when failure arises above, current rpc will be reused
		index, _ = client.GetRPCInstance()
		require.Equal(t, index, 0)
	}

	// 4xx attributes to remote server fault, RPC should rotate
	statusCodes = []int{403, 429}
	for _, sc := range statusCodes {

		httpRespError := rpc.HTTPError{
			StatusCode: sc,
			Status:     "INTERNAL_ERROR",
			Body:       []byte{},
		}

		client, _ := makeTestMultihomingClient(1, httpRespError)

		index, _ := client.GetRPCInstance()
		require.Equal(t, index, 0)

		makeFailureCall(t, client, 1)

		// given num retry is 1, when failure arises, current rpc should becomes the next one
		index, _ = client.GetRPCInstance()
		require.Equal(t, index, 2)
	}

	// 2xx attributes to sender's fault with JSON RPC fault, RPC should not rotate
	rpcError := JsonError{}

	client, _ := makeTestMultihomingClient(2, &rpcError)

	index, _ := client.GetRPCInstance()
	require.Equal(t, index, 0)

	makeFailureCall(t, client, 10)

	// given num retry is 0, when failure arises above, current rpc should becomes the next one
	index, _ = client.GetRPCInstance()
	require.Equal(t, index, 1)

}

func TestMultihomingClientRPCFaultZeroRetry(t *testing.T) {
	httpRespError := make500Error()
	client, _ := makeTestMultihomingClient(0, httpRespError)

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

func TestMultihomingClientRPCFaultOneRetry(t *testing.T) {
	httpRespError := make500Error()

	client, _ := makeTestMultihomingClient(1, httpRespError)

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

func TestMultihomingClientRPCFaultTwoRetry(t *testing.T) {
	httpRespError := make500Error()
	client, _ := makeTestMultihomingClient(2, httpRespError)

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
