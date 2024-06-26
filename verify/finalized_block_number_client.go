package verify

import (
	"context"
	"encoding/json"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

type FinalizedBlockClient struct {
	c *rpc.Client
}

// Dial connects a client to the given URL.
func Dial(rawurl string) (*FinalizedBlockClient, error) {
	return DialContext(context.Background(), rawurl)
}

// DialContext connects a client to the given URL with context.
func DialContext(ctx context.Context, rawurl string) (*FinalizedBlockClient, error) {
	c, err := rpc.DialContext(ctx, rawurl)
	if err != nil {
		return nil, err
	}
	return NewFinalizedBlockClient(c), nil
}

// NewFinalizedBlockClient creates a client that uses the given RPC client.
func NewFinalizedBlockClient(c *rpc.Client) *FinalizedBlockClient {
	return &FinalizedBlockClient{c}
}

// Close closes the underlying RPC connection.
func (ec *FinalizedBlockClient) Close() {
	ec.c.Close()
}

// Client gets the underlying RPC client.
func (ec *FinalizedBlockClient) Client() *rpc.Client {
	return ec.c
}

func (c *FinalizedBlockClient) GetBlock(ctx context.Context, method string, args ...interface{}) (*types.Block, error) {
	var raw json.RawMessage
	err := c.c.CallContext(ctx, &raw, method, args...)
	if err != nil {
		return nil, err
	}

	// Decode header and transactions.
	var head *types.Header
	if err := json.Unmarshal(raw, &head); err != nil {
		return nil, err
	}
	// When the block is not found, the API returns JSON null.
	if head == nil {
		return nil, ethereum.NotFound
	}

	return types.NewBlockWithHeader(head), nil
}
