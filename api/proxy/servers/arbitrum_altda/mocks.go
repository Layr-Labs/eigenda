package arbitrum_altda

import (
	"context"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// mockEthClient is a simple stub implementation of the IEthClient interface
// used when memstore is enabled to avoid actual Ethereum RPC calls.
// It returns an empty block header where block_number=0 ensuring that the
// recency check will be bypassed.
type mockEthClient struct{}

// NewMockEthClient creates a new stub ETH client for memstore mode.
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
