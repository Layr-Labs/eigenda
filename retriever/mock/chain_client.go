package mock

import (
	"context"
	"math/big"

	binding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"
	"github.com/Layr-Labs/eigenda/retriever/eth"
	gcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
)

type MockChainClient struct {
	mock.Mock
}

var _ eth.ChainClient = (*MockChainClient)(nil)

func NewMockChainClient() *MockChainClient {
	return &MockChainClient{}
}

func (c *MockChainClient) FetchBatchHeader(ctx context.Context, serviceManagerAddress gcommon.Address, batchHeaderHash []byte, fromBlock *big.Int, toBlock *big.Int) (*binding.BatchHeader, error) {
	args := c.Called()
	return args.Get(0).(*binding.BatchHeader), args.Error(1)
}
