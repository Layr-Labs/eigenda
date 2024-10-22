package mock

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/Layr-Labs/eigenda/api/clients"
	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/stretchr/testify/mock"
)

type BlobKey struct {
	BatchHeaderHash []byte
	BlobIndex       uint32
}

type MockDisperserClient struct {
	mock.Mock
	mockRequestIDStore map[string][]byte
	mockRetrievalStore map[string][]byte
}

var _ clients.DisperserClient = (*MockDisperserClient)(nil)

func NewMockDisperserClient() *MockDisperserClient {
	return &MockDisperserClient{
		mockRequestIDStore: make(map[string][]byte),
		mockRetrievalStore: make(map[string][]byte),
	}
}

func (c *MockDisperserClient) DisperseBlobAuthenticated(ctx context.Context, data []byte, quorums []uint8) (*disperser.BlobStatus, []byte, error) {
	// do data validation for the benefit of high-level client tests
	_, err := rs.ToFrArray(data)
	if err != nil {
		return nil, nil, err
	}

	args := c.Called(data, quorums)
	var status *disperser.BlobStatus
	if args.Get(0) != nil {
		status = (args.Get(0)).(*disperser.BlobStatus)
	}

	var key []byte
	if args.Get(1) != nil {
		key = (args.Get(1)).([]byte)
	}

	if args.Get(2) != nil {
		err = (args.Get(2)).(error)
	}

	keyStr := base64.StdEncoding.EncodeToString(key)
	c.mockRequestIDStore[keyStr] = data

	return status, key, err
}

func (c *MockDisperserClient) DisperseBlob(ctx context.Context, data []byte, quorums []uint8) (*disperser.BlobStatus, []byte, error) {
	args := c.Called(data, quorums)
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

	keyStr := base64.StdEncoding.EncodeToString(key)
	c.mockRequestIDStore[keyStr] = data

	return status, key, err
}

// TODO: implement in the subsequent PR
func (c *MockDisperserClient) DispersePaidBlob(ctx context.Context, data []byte, quorums []uint8) (*disperser.BlobStatus, []byte, error) {
	return nil, nil, nil
}

func (c *MockDisperserClient) GetBlobStatus(ctx context.Context, key []byte) (*disperser_rpc.BlobStatusReply, error) {
	args := c.Called(key)
	var reply *disperser_rpc.BlobStatusReply
	if args.Get(0) != nil {
		reply = (args.Get(0)).(*disperser_rpc.BlobStatusReply)
		if reply.Status == disperser_rpc.BlobStatus_FINALIZED {
			retrievalKey := fmt.Sprintf("%s-%d", base64.StdEncoding.EncodeToString(reply.Info.BlobVerificationProof.BatchMetadata.BatchHeaderHash), reply.Info.BlobVerificationProof.BlobIndex)
			requestIDKey := base64.StdEncoding.EncodeToString(key)
			c.mockRetrievalStore[retrievalKey] = c.mockRequestIDStore[requestIDKey]
		}
	}
	var err error
	if args.Get(1) != nil {
		err = (args.Get(1)).(error)
	}

	return reply, err
}

func (c *MockDisperserClient) RetrieveBlob(ctx context.Context, batchHeaderHash []byte, blobIndex uint32) ([]byte, error) {
	args := c.Called(batchHeaderHash, blobIndex)
	var blob []byte
	if args.Get(0) != nil {
		blob = (args.Get(0)).([]byte)
	} else {
		keyStr := fmt.Sprintf("%s-%d", base64.StdEncoding.EncodeToString(batchHeaderHash), blobIndex)
		blob = c.mockRetrievalStore[keyStr]
	}

	var err error
	if args.Get(1) != nil {
		err = (args.Get(1)).(error)
	}
	return blob, err
}

func (c *MockDisperserClient) Close() error {
	args := c.Called()
	return args.Error(0)
}
