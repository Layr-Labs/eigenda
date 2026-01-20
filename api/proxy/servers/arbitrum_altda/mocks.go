package arbitrum_altda

import (
	"context"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// mockEthClient is a simple mock implementation of the IEthClient interface
// for use in tests. It returns deterministic mock blocks for any block hash.
type mockEthClient struct{}

// NewMockEthClient creates a new mock ETH client for testing.
// This client returns deterministic mock blocks with reasonable defaults.
func NewMockEthClient() IEthClient {
	return &mockEthClient{}
}

// BlockByHash returns a mock block with a deterministic block number.
// This implementation always succeeds and returns 0 which is mapped to the
// L1 Inbox Submission block number which forces the verifyCertRBNRecencyCheck call to
// fail
func (m *mockEthClient) BlockByHash(ctx context.Context, hash gethcommon.Hash) (*types.Block, error) {
	header := &types.Header{
		Number: big.NewInt(0),
	}
	return types.NewBlockWithHeader(header), nil
}
