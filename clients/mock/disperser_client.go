package mock

import (
	"context"

	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/clients"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/stretchr/testify/mock"
)

type MockDisperserClient struct {
	mock.Mock
}

var _ clients.DisperserClient = (*MockDisperserClient)(nil)

func NewMockDisperserClient() *MockDisperserClient {
	return &MockDisperserClient{}
}

func (c *MockDisperserClient) DisperseBlobAuthenticated(ctx context.Context, data []byte, securityParams []*core.SecurityParam) (*disperser.BlobStatus, []byte, error) {
	args := c.Called(data, securityParams)
	var status *disperser.BlobStatus
	if args.Get(0) != nil {
		status = (args.Get(0)).(*disperser.BlobStatus)
	}
	var key []byte
	if args.Get(1) != nil {
		key = (args.Get(1)).([]byte)
	}
	var err error
	if args.Get(2) != nil {
		err = (args.Get(2)).(error)
	}
	return status, key, err
}

func (c *MockDisperserClient) DisperseBlob(ctx context.Context, data []byte, securityParams []*core.SecurityParam) (*disperser.BlobStatus, []byte, error) {
	args := c.Called(data, securityParams)
	var status *disperser.BlobStatus
	if args.Get(0) != nil {
		status = (args.Get(0)).(*disperser.BlobStatus)
	}
	var key []byte
	if args.Get(1) != nil {
		key = (args.Get(1)).([]byte)
	}
	var err error
	if args.Get(2) != nil {
		err = (args.Get(2)).(error)
	}
	return status, key, err
}

func (c *MockDisperserClient) GetBlobStatus(ctx context.Context, key []byte) (*disperser_rpc.BlobStatusReply, error) {
	args := c.Called(key)
	var reply *disperser_rpc.BlobStatusReply
	if args.Get(0) != nil {
		reply = (args.Get(0)).(*disperser_rpc.BlobStatusReply)
	}
	var err error
	if args.Get(1) != nil {
		err = (args.Get(1)).(error)
	}
	return reply, err
}
