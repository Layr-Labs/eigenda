package mock

import (
	"context"

	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/stretchr/testify/mock"
	"github.com/wealdtech/go-merkletree/v2"
)

type MockNodeClient struct {
	mock.Mock
}

var _ clients.NodeClient = (*MockNodeClient)(nil)

func NewNodeClient() *MockNodeClient {
	return &MockNodeClient{}
}

func (c *MockNodeClient) GetBlobHeader(ctx context.Context, socket string, batchHeaderHash [32]byte, blobIndex uint32) (*core.BlobHeader, *merkletree.Proof, error) {
	args := c.Called(socket, batchHeaderHash, blobIndex)
	var hashes [][]byte
	if args.Get(1) != nil {
		hashes = (args.Get(1)).([][]byte)
	}

	var index uint64
	if args.Get(2) != nil {
		index = (args.Get(2)).(uint64)
	}

	var err error = nil
	if args.Get(3) != nil {
		err = args.Get(3).(error)
	}

	proof := &merkletree.Proof{
		Hashes: hashes,
		Index:  index,
	}
	return (args.Get(0)).(*core.BlobHeader), proof, err
}

func (c *MockNodeClient) GetChunks(
	ctx context.Context,
	opID core.OperatorID,
	opInfo *core.IndexedOperatorInfo,
	batchHeaderHash [32]byte,
	blobIndex uint32,
	quorumID core.QuorumID,
	chunksChan chan clients.RetrievedChunks,
) {
	args := c.Called(opID, opInfo, batchHeaderHash, blobIndex)
	encodedBlob := (args.Get(0)).(core.EncodedBlob)
	chunks, err := encodedBlob.EncodedBundlesByOperator[opID][quorumID].ToFrames()
	if err != nil {
		chunksChan <- clients.RetrievedChunks{
			OperatorID: opID,
			Err:        err,
			Chunks:     nil,
		}

	}
	chunksChan <- clients.RetrievedChunks{
		OperatorID: opID,
		Err:        nil,
		Chunks:     chunks,
	}
}
