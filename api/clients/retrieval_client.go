package clients

import (
	"context"
	"errors"
	"fmt"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/wealdtech/go-merkletree/v2"

	"github.com/gammazero/workerpool"
	"github.com/wealdtech/go-merkletree/v2/keccak256"
)

// RetrievalClient is an object that can retrieve blobs from the network.
type RetrievalClient interface {

	// RetrieveBlob fetches a blob from the network. This method is equivalent to calling
	// RetrieveBlobChunks to get the chunks and then CombineChunks to recombine those chunks into the original blob.
	RetrieveBlob(
		ctx context.Context,
		batchHeaderHash [32]byte,
		blobIndex uint32,
		referenceBlockNumber uint,
		batchRoot [32]byte,
		quorumID core.QuorumID) ([]byte, error)

	// RetrieveBlobChunks downloads the chunks of a blob from the network but do not recombine them. Use this method
	// if detailed information about which node returned which chunk is needed. Otherwise, use RetrieveBlob.
	RetrieveBlobChunks(
		ctx context.Context,
		batchHeaderHash [32]byte,
		blobIndex uint32,
		referenceBlockNumber uint,
		batchRoot [32]byte,
		quorumID core.QuorumID) (*BlobChunks, error)

	// CombineChunks recombines the chunks into the original blob.
	CombineChunks(chunks *BlobChunks) ([]byte, error)
}

// BlobChunks is a collection of chunks retrieved from the network which can be recombined into a blob.
type BlobChunks struct {
	Chunks           []*encoding.Frame
	Indices          []encoding.ChunkNumber
	EncodingParams   encoding.EncodingParams
	BlobHeaderLength uint
	Assignments      map[core.OperatorID]core.Assignment
	AssignmentInfo   core.AssignmentInfo
}

type retrievalClient struct {
	logger                logging.Logger
	indexedChainState     core.IndexedChainState
	assignmentCoordinator core.AssignmentCoordinator
	nodeClient            NodeClient
	verifier              encoding.Verifier
	numConnections        int
}

// NewRetrievalClient creates a new retrieval client.
func NewRetrievalClient(
	logger logging.Logger,
	chainState core.IndexedChainState,
	assignmentCoordinator core.AssignmentCoordinator,
	nodeClient NodeClient,
	verifier encoding.Verifier,
	numConnections int) (RetrievalClient, error) {

	return &retrievalClient{
		logger:                logger.With("component", "RetrievalClient"),
		indexedChainState:     chainState,
		assignmentCoordinator: assignmentCoordinator,
		nodeClient:            nodeClient,
		verifier:              verifier,
		numConnections:        numConnections,
	}, nil
}

// RetrieveBlob retrieves a blob from the network.
func (r *retrievalClient) RetrieveBlob(
	ctx context.Context,
	batchHeaderHash [32]byte,
	blobIndex uint32,
	referenceBlockNumber uint,
	batchRoot [32]byte,
	quorumID core.QuorumID) ([]byte, error) {

	chunks, err := r.RetrieveBlobChunks(ctx, batchHeaderHash, blobIndex, referenceBlockNumber, batchRoot, quorumID)
	if err != nil {
		return nil, err
	}

	return r.CombineChunks(chunks)
}

// RetrieveBlobChunks retrieves the chunks of a blob from the network but does not recombine them.
func (r *retrievalClient) RetrieveBlobChunks(ctx context.Context,
	batchHeaderHash [32]byte,
	blobIndex uint32,
	referenceBlockNumber uint,
	batchRoot [32]byte,
	quorumID core.QuorumID) (*BlobChunks, error) {

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

	// Validate the blob length
	err = r.verifier.VerifyBlobLength(blobHeader.BlobCommitments)
	if err != nil {
		return nil, err
	}

	// Validate the commitments are equivalent
	commitmentBatch := []encoding.BlobCommitments{blobHeader.BlobCommitments}
	err = r.verifier.VerifyCommitEquivalenceBatch(commitmentBatch)
	if err != nil {
		return nil, err
	}

	assignments, info, err := r.assignmentCoordinator.GetAssignments(indexedOperatorState.OperatorState, blobHeader.Length, quorumHeader)
	if err != nil {
		return nil, errors.New("failed to get assignments")
	}

	// Fetch chunks from all operators
	chunksChan := make(chan RetrievedChunks, len(operators))
	pool := workerpool.New(r.numConnections)
	for opID := range operators {
		opID := opID
		opInfo := indexedOperatorState.IndexedOperators[opID]
		pool.Submit(func() {
			r.nodeClient.GetChunks(ctx, opID, opInfo, batchHeaderHash, blobIndex, quorumID, chunksChan)
		})
	}

	encodingParams := encoding.ParamsFromMins(quorumHeader.ChunkLength, info.TotalChunks)

	var chunks []*encoding.Frame
	var indices []encoding.ChunkNumber
	// TODO(ian-shim): if we gathered enough chunks, cancel remaining RPC calls
	for i := 0; i < len(operators); i++ {
		reply := <-chunksChan
		if reply.Err != nil {
			r.logger.Error("failed to get chunks from operator", "operator", reply.OperatorID.Hex(), "err", reply.Err)
			continue
		}
		assignment, ok := assignments[reply.OperatorID]
		if !ok {
			return nil, fmt.Errorf("no assignment to operator %s", reply.OperatorID.Hex())
		}

		err = r.verifier.VerifyFrames(reply.Chunks, assignment.GetIndices(), blobHeader.BlobCommitments, encodingParams)
		if err != nil {
			r.logger.Error("failed to verify chunks from operator", "operator", reply.OperatorID.Hex(), "err", err)
			continue
		} else {
			r.logger.Info("verified chunks from operator", "operator", reply.OperatorID.Hex())
		}

		chunks = append(chunks, reply.Chunks...)
		indices = append(indices, assignment.GetIndices()...)
	}

	return &BlobChunks{
		Chunks:           chunks,
		Indices:          indices,
		EncodingParams:   encodingParams,
		BlobHeaderLength: blobHeader.Length,
		Assignments:      assignments,
		AssignmentInfo:   info,
	}, nil
}

// CombineChunks recombines the chunks into the original blob.
func (r *retrievalClient) CombineChunks(chunks *BlobChunks) ([]byte, error) {
	return r.verifier.Decode(
		chunks.Chunks,
		chunks.Indices,
		chunks.EncodingParams,
		uint64(chunks.BlobHeaderLength)*encoding.BYTES_PER_SYMBOL)
}
