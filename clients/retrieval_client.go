package clients

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
	"github.com/gammazero/workerpool"
	"github.com/wealdtech/go-merkletree"
	"github.com/wealdtech/go-merkletree/keccak256"

	coreindexer "github.com/Layr-Labs/eigenda/core/indexer"
	"github.com/Layr-Labs/eigenda/indexer"
)

type RetrievalClient interface {
	RetrieveBlob(
		ctx context.Context,
		batchHeaderHash [32]byte,
		blobIndex uint32,
		referenceBlockNumber uint,
		batchRoot [32]byte,
		quorumID core.QuorumID) ([]byte, error)
}

type retrievalClient struct {
	logger                common.Logger
	indexedChainState     core.IndexedChainState
	assignmentCoordinator core.AssignmentCoordinator
	nodeClient            NodeClient
	encoder               core.Encoder
	numConnections        int
}

var _ RetrievalClient = (*retrievalClient)(nil)

func NewRetrievalClient(
	logger common.Logger,
	chainState core.ChainState,
	indexer indexer.Indexer,
	assignmentCoordinator core.AssignmentCoordinator,
	nodeClient NodeClient,
	encoder core.Encoder,
	numConnections int,
) (*retrievalClient, error) {

	indexedState, err := coreindexer.NewIndexedChainState(
		chainState,
		indexer,
	)
	if err != nil {
		return nil, err
	}
	return &retrievalClient{
		logger:                logger,
		indexedChainState:     indexedState,
		assignmentCoordinator: assignmentCoordinator,
		nodeClient:            nodeClient,
		encoder:               encoder,
		numConnections:        numConnections,
	}, nil
}

func (r *retrievalClient) RetrieveBlob(
	ctx context.Context,
	batchHeaderHash [32]byte,
	blobIndex uint32,
	referenceBlockNumber uint,
	batchRoot [32]byte,
	quorumID core.QuorumID) ([]byte, error) {
	indexedOperatorState, err := r.indexedChainState.GetIndexedOperatorState(ctx, referenceBlockNumber, []core.QuorumID{quorumID})
	if err != nil {
		return nil, err
	}
	operators, ok := indexedOperatorState.Operators[quorumID]
	if !ok {
		return nil, fmt.Errorf("no quorum with ID: %d", quorumID)
	}

	// Get blob header from any operator
	var blobHeader *core.BlobHeader
	var proof *merkletree.Proof
	var proofVerified bool
	for opID := range operators {
		opInfo := indexedOperatorState.IndexedOperators[opID]
		blobHeader, proof, err = r.nodeClient.GetBlobHeader(ctx, opInfo.Socket, batchHeaderHash, blobIndex)
		if err != nil {
			// try another operator
			r.logger.Warn("failed to dial operator while fetching BlobHeader, trying different operator", "operator", opInfo.Socket, "err", err)
			continue
		}

		blobHeaderHash, err := blobHeader.GetBlobHeaderHash()
		if err != nil {
			r.logger.Warn("got invalid blob header, trying different operator", "operator", opInfo.Socket, "err", err)
			continue
		}
		proofVerified, err = merkletree.VerifyProofUsing(blobHeaderHash[:], false, proof, [][]byte{batchRoot[:]}, keccak256.New())
		if err != nil {
			r.logger.Warn("got invalid blob header proof, trying different operator", "operator", opInfo.Socket, "err", err)
			continue
		}
		if !proofVerified {
			r.logger.Warn("failed to verify blob header against given proof, trying different operator", "operator", opInfo.Socket)
			continue
		}

		break
	}
	if blobHeader == nil || proof == nil || !proofVerified {
		return nil, fmt.Errorf("failed to get blob header from all operators (header hash: %s, index: %d)", batchHeaderHash, blobIndex)
	}

	var quorumHeader *core.BlobQuorumInfo
	for _, header := range blobHeader.QuorumInfos {
		if header.QuorumID == quorumID {
			quorumHeader = header
			break
		}
	}

	if quorumHeader == nil {
		return nil, fmt.Errorf("no quorum header for quorum %d", quorumID)
	}

	assignements, info, err := r.assignmentCoordinator.GetAssignments(indexedOperatorState.OperatorState, quorumID, uint(quorumHeader.QuantizationFactor))
	if err != nil {
		return nil, fmt.Errorf("failed to get assignments")
	}

	// Fetch chunks from all operators
	chunksChan := make(chan RetrievedChunks, len(operators))
	pool := workerpool.New(r.numConnections)
	for opID := range operators {
		opID := opID
		opInfo := indexedOperatorState.IndexedOperators[opID]
		pool.Submit(func() {
			r.nodeClient.GetChunks(ctx, opID, opInfo, batchHeaderHash, blobIndex, quorumID, chunksChan)
			// TODO(ian-shim): validate chunks received from nodes
		})
	}

	var chunks []*core.Chunk
	var indices []core.ChunkNumber
	// TODO(ian-shim): if we gathered enough chunks, cancel remaining RPC calls
	for i := 0; i < len(operators); i++ {
		reply := <-chunksChan
		if reply.Err != nil {
			continue
		}
		assignment, ok := assignements[reply.OperatorID]
		if !ok {
			return nil, fmt.Errorf("no assignment to operator %v", reply.OperatorID)
		}

		chunks = append(chunks, reply.Chunks...)
		indices = append(indices, assignment.GetIndices()...)
	}

	chunkLength, err := r.assignmentCoordinator.GetChunkLengthFromHeader(indexedOperatorState.OperatorState, quorumHeader)
	if err != nil {
		return nil, err
	}

	encodingParams, err := core.GetEncodingParams(chunkLength, info.TotalChunks)
	if err != nil {
		return nil, err
	}

	return r.encoder.Decode(chunks, indices, encodingParams, uint64(blobHeader.Length)*bn254.BYTES_PER_COEFFICIENT)
}
