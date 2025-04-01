package clients

import (
	"context"
	"errors"
	"fmt"

	"github.com/docker/go-units"

	"github.com/Layr-Labs/eigenda/api/clients"
	grpcnode "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gammazero/workerpool"
)

// RetrievalClient is an object that can retrieve blobs from the DA nodes.
// To retrieve a blob from the relay, use RelayClient instead.
type RetrievalClient interface {
	// GetBlob downloads chunks of a blob from operator network and reconstructs the blob.
	GetBlob(
		ctx context.Context,
		blobKey corev2.BlobKey,
		blobVersion corev2.BlobVersion,
		blobCommitments encoding.BlobCommitments,
		referenceBlockNumber uint64,
		quorumID core.QuorumID,
	) ([]byte, error)
}

type retrievalClient struct {
	logger         logging.Logger
	ethClient      core.Reader
	chainState     core.ChainState
	verifier       encoding.Verifier
	numConnections int
}

var _ RetrievalClient = &retrievalClient{}

// NewRetrievalClient creates a new retrieval client.
func NewRetrievalClient(
	logger logging.Logger,
	ethClient core.Reader,
	chainState core.ChainState,
	verifier encoding.Verifier,
	numConnections int,
) RetrievalClient {
	return &retrievalClient{
		logger:         logger.With("component", "RetrievalClient"),
		ethClient:      ethClient,
		chainState:     chainState,
		verifier:       verifier,
		numConnections: numConnections,
	}
}

func (r *retrievalClient) GetBlob(
	ctx context.Context,
	blobKey corev2.BlobKey,
	blobVersion corev2.BlobVersion,
	blobCommitments encoding.BlobCommitments,
	referenceBlockNumber uint64,
	quorumID core.QuorumID,
) ([]byte, error) {

	commitmentBatch := []encoding.BlobCommitments{blobCommitments}
	err := r.verifier.VerifyCommitEquivalenceBatch(commitmentBatch)
	if err != nil {
		return nil, err
	}

	operatorState, err := r.chainState.GetOperatorStateWithSocket(ctx, uint(referenceBlockNumber), []core.QuorumID{quorumID})
	if err != nil {
		return nil, err
	}
	operators, ok := operatorState.Operators[quorumID]
	if !ok {
		return nil, fmt.Errorf("no quorum with ID: %d", quorumID)
	}

	blobVersions, err := r.ethClient.GetAllVersionedBlobParams(ctx)
	if err != nil {
		return nil, err
	}

	blobParam, ok := blobVersions[blobVersion]
	if !ok {
		return nil, fmt.Errorf("invalid blob version %d", blobVersion)
	}

	encodingParams, err := corev2.GetEncodingParams(blobCommitments.Length, blobParam)
	if err != nil {
		return nil, err
	}

	assignments, err := corev2.GetAssignments(operatorState, blobParam, quorumID)
	if err != nil {
		return nil, errors.New("failed to get assignments")
	}

	// Fetch chunks from all operators
	chunksChan := make(chan clients.RetrievedChunks, len(operators))
	pool := workerpool.New(r.numConnections)
	for opID := range operators {
		opInfo := operatorState.Operators[quorumID][opID]
		pool.Submit(func() {
			r.getChunksFromOperator(ctx, opID, opInfo, blobKey, quorumID, chunksChan)
		})
	}

	var chunks []*encoding.Frame
	var indices []encoding.ChunkNumber
	// TODO(ian-shim): if we gathered enough chunks, cancel remaining RPC calls
	for i := 0; i < len(operators); i++ {
		reply := <-chunksChan
		if reply.Err != nil {
			r.logger.Warn("failed to get chunks from operator", "operator", reply.OperatorID.Hex(), "err", reply.Err)
			continue
		}
		assignment, ok := assignments[reply.OperatorID]
		if !ok {
			return nil, fmt.Errorf("no assignment to operator %s", reply.OperatorID.Hex())
		}

		assignmentIndices := make([]uint, len(assignment.GetIndices()))
		for i, index := range assignment.GetIndices() {
			assignmentIndices[i] = uint(index)
		}

		err = r.verifier.VerifyFrames(reply.Chunks, assignmentIndices, blobCommitments, encodingParams)
		if err != nil {
			r.logger.Warn("failed to verify chunks from operator", "operator", reply.OperatorID.Hex(), "err", err)
			continue
		} else {
			r.logger.Info("verified chunks from operator", "operator", reply.OperatorID.Hex())
		}

		chunks = append(chunks, reply.Chunks...)
		indices = append(indices, assignmentIndices...)
	}

	if len(chunks) == 0 {
		return nil, errors.New("failed to retrieve any chunks")
	}

	return r.verifier.Decode(
		chunks,
		indices,
		encodingParams,
		uint64(blobCommitments.Length)*encoding.BYTES_PER_SYMBOL,
	)
}

func (r *retrievalClient) getChunksFromOperator(
	ctx context.Context,
	opID core.OperatorID,
	opInfo *core.OperatorInfo,
	blobKey corev2.BlobKey,
	quorumID core.QuorumID,
	chunksChan chan clients.RetrievedChunks,
) {

	// TODO (cody-littley): this client should be refactored to make requests for smaller quantities of data
	//  in order to avoid hitting the max message size limit. This will allow us to have much smaller
	//  message size limits.

	maxBlobSize := 16 * units.MiB // maximum size of the original blob
	encodingRate := 8             // worst case scenario if one validator has 100% stake
	fudgeFactor := units.MiB      // to allow for some overhead from things like protobuf encoding
	maxMessageSize := maxBlobSize*encodingRate + fudgeFactor

	conn, err := grpc.NewClient(
		core.OperatorSocket(opInfo.Socket).GetV2RetrievalSocket(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMessageSize)),
	)
	defer func() {
		err := conn.Close()
		if err != nil {
			r.logger.Error("failed to close connection", "err", err)
		}
	}()
	if err != nil {
		chunksChan <- clients.RetrievedChunks{
			OperatorID: opID,
			Err:        err,
			Chunks:     nil,
		}
		return
	}

	n := grpcnode.NewRetrievalClient(conn)
	request := &grpcnode.GetChunksRequest{
		BlobKey:  blobKey[:],
		QuorumId: uint32(quorumID),
	}

	reply, err := n.GetChunks(ctx, request)
	if err != nil {
		chunksChan <- clients.RetrievedChunks{
			OperatorID: opID,
			Err:        err,
			Chunks:     nil,
		}
		return
	}

	chunks := make([]*encoding.Frame, len(reply.GetChunks()))
	for i, data := range reply.GetChunks() {
		var chunk *encoding.Frame
		chunk, err = new(encoding.Frame).DeserializeGnark(data)
		if err != nil {
			chunksChan <- clients.RetrievedChunks{
				OperatorID: opID,
				Err:        err,
				Chunks:     nil,
			}
			return
		}

		chunks[i] = chunk
	}
	chunksChan <- clients.RetrievedChunks{
		OperatorID: opID,
		Err:        nil,
		Chunks:     chunks,
	}
}
