package test

import (
	"context"
	binding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"sync"
)

type mockChainClient struct {
	lock  *sync.Mutex
	Count uint
}

func newMockChainClient(lock *sync.Mutex) *mockChainClient {
	return &mockChainClient{
		lock: lock,
	}

}

func (m *mockChainClient) FetchBatchHeader(
	ctx context.Context,
	serviceManagerAddress common.Address,
	batchHeaderHash []byte,
	fromBlock *big.Int,
	toBlock *big.Int) (*binding.IEigenDAServiceManagerBatchHeader, error) {

	m.lock.Lock()
	defer m.lock.Unlock()

	m.Count++

	return &binding.IEigenDAServiceManagerBatchHeader{}, nil
}
