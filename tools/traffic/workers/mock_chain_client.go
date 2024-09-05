package workers

import (
	"context"
	binding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"math/big"
)

type mockChainClient struct {
	mock mock.Mock
}

func (m *mockChainClient) FetchBatchHeader(
	ctx context.Context,
	serviceManagerAddress common.Address,
	batchHeaderHash []byte,
	fromBlock *big.Int,
	toBlock *big.Int) (*binding.IEigenDAServiceManagerBatchHeader, error) {

	m.mock.Called(serviceManagerAddress, batchHeaderHash, fromBlock, toBlock)
	return &binding.IEigenDAServiceManagerBatchHeader{}, nil
}
