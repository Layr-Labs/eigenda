package clients

import (
	"context"
	"errors"
	"fmt"

	"github.com/Layr-Labs/eigenda/common"
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

	// GetBlobWithProbe is the same as GetBlob, but it also takes a probe object to capture metrics. No metrics are
	// captured if the probe is nil.
	GetBlobWithProbe(
		ctx context.Context,
		blobKey corev2.BlobKey,
		blobVersion corev2.BlobVersion,
		blobCommitments encoding.BlobCommitments,
		referenceBlockNumber uint64,
		quorumID core.QuorumID,
		probe *common.SequenceProbe,
	) ([]byte, error)
}

type retrievalClient struct {
	logger         logging.Logger
	ethClient      core.Reader
	chainState     core.ChainState
	verifier       encoding.Verifier
	connectionPool *workerpool.WorkerPool
	computePool    *workerpool.WorkerPool
}

// encapsulates a GetChunksReply
type getChunksReply struct {
	OperatorID core.OperatorID
	Err        error
	Reply      *grpcnode.GetChunksReply
}

// work done by decode is passed back to the caller via this struct
type decodeChunksReply struct {
	Blob []byte
	Err  error
}

var _ RetrievalClient = &retrievalClient{}

// NewRetrievalClient creates a new retrieval client.
func NewRetrievalClient(
	logger logging.Logger,
	ethClient core.Reader,
	chainState core.ChainState,
	verifier encoding.Verifier,
// connectionPoolSize limits the maximum number of concurrent network connections
	connectionPoolSize int,
// computePoolSize limits the maximum number of concurrent compute intensive tasks
	computePoolSize int,
) RetrievalClient {
	return &retrievalClient{
		logger:         logger.With("component", "RetrievalClient"),
		ethClient:      ethClient,
		chainState:     chainState,
		verifier:       verifier,
		connectionPool: workerpool.New(connectionPoolSize),
		computePool:    workerpool.New(computePoolSize),
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
	return r.GetBlobWithProbe(ctx, blobKey, blobVersion, blobCommitments, referenceBlockNumber, quorumID, nil)
}

func (r *retrievalClient) GetBlobWithProbe(
	ctx context.Context,
	blobKey corev2.BlobKey,
	blobVersion corev2.BlobVersion,
	blobCommitments encoding.BlobCommitments,
	referenceBlockNumber uint64,
	quorumID core.QuorumID,
	probe *common.SequenceProbe,
) ([]byte, error) {

	// TODO: currently, we download, verify, and decode all chunks.
	//  Instead, we could get away with only downloading 1/(encoding ratio) chunks.

	probe.SetStage("verify_commitment")
	commitmentBatch := []encoding.BlobCommitments{blobCommitments}
	err := r.verifier.VerifyCommitEquivalenceBatch(commitmentBatch)
	if err != nil {
		return nil, err
	}

	probe.SetStage("get_operator_state")
	operatorState, err := r.chainState.GetOperatorStateWithSocket(ctx, uint(referenceBlockNumber), []core.QuorumID{quorumID})
	if err != nil {
		return nil, err
	}
	operators, ok := operatorState.Operators[quorumID]
	if !ok {
		return nil, fmt.Errorf("no quorum with ID: %d", quorumID)
	}

	probe.SetStage("get_blob_versions")
	blobVersions, err := r.ethClient.GetAllVersionedBlobParams(ctx)
	if err != nil {
		return nil, err
	}

	blobParam, ok := blobVersions[blobVersion]
	if !ok {
		return nil, fmt.Errorf("invalid blob version %d", blobVersion)
	}

	probe.SetStage("get_encoding_params")
	encodingParams, err := corev2.GetEncodingParams(blobCommitments.Length, blobParam)
	if err != nil {
		return nil, err
	}

	probe.SetStage("get_assignments")
	assignments, err := corev2.GetAssignments(operatorState, blobParam, quorumID)
	if err != nil {
		return nil, errors.New("failed to get assignments")
	}

	// Submit download requests to the connection pool
	probe.SetStage("connection_pool")
	replyChan := make(chan *getChunksReply, len(operators))
	for opID := range operators {
		// make sure the value doesn't change before being submitted to the pool
		boundOperatorId := opID
		opInfo := operatorState.Operators[quorumID][boundOperatorId]
		r.connectionPool.Submit(func() {
			r.downloadChunks(ctx, boundOperatorId, opInfo, blobKey, quorumID, replyChan)
		})
	}

	// Wait for all download requests to finish
	probe.SetStage("download")
	replies := make([]*getChunksReply, 0, len(operators))
	for i := 0; i < len(operators); i++ {
		select {
		case reply := <-replyChan:
			if reply.Err != nil {
				r.logger.Warn("failed to get chunks from operator",
					"operator", reply.OperatorID.Hex(),
					"err", reply.Err)
				continue
			}
			replies = append(replies, reply)
		case <-ctx.Done():
			return nil, errors.New("context cancelled while waiting for chunks from operators")
		}
	}

	// Submit deserialization and verification requests to the compute pool
	probe.SetStage("compute_pool")
	deserializeChan := make(chan clients.RetrievedChunks, len(operators))
	for _, reply := range replies {
		r.computePool.Submit(func() {
			r.deserializeAndVerifyChunks(
				reply.OperatorID,
				assignments,
				reply.Reply,
				blobCommitments,
				encodingParams,
				deserializeChan)
		})
	}

	probe.SetStage("deserialize_and_verify")
	var chunks []*encoding.Frame
	var indices []encoding.ChunkNumber
	for i := 0; i < len(replies); i++ {
		select {
		case reply := <-deserializeChan:
			if reply.Err != nil {
				deserializeChan <- clients.RetrievedChunks{
					OperatorID: reply.OperatorID,
					Err:        reply.Err,
				}
			} else {
				assignment, ok := assignments[reply.OperatorID]
				if !ok {
					return nil, fmt.Errorf("no assignment to operator %s", reply.OperatorID.Hex())
				}

				assignmentIndices := make([]uint, len(assignment.GetIndices()))
				for i, index := range assignment.GetIndices() {
					assignmentIndices[i] = uint(index)
				}

				chunks = append(chunks, reply.Chunks...)
				indices = append(indices, assignmentIndices...)
			}
		case <-ctx.Done():
			return nil, errors.New("context cancelled while waiting for chunks from operators")
		}
	}

	probe.SetStage("compute_pool")
	decodeResponseChan := make(chan *decodeChunksReply, 1)
	r.computePool.Submit(func() {
		blob, err := r.verifier.Decode(
			chunks,
			indices,
			encodingParams,
			uint64(blobCommitments.Length)*encoding.BYTES_PER_SYMBOL,
		)
		decodeResponseChan <- &decodeChunksReply{
			Blob: blob,
			Err:  err,
		}
	})

	probe.SetStage("decode")
	select {
	case decodeResponse := <-decodeResponseChan:
		return decodeResponse.Blob, decodeResponse.Err
	case <-ctx.Done():
		return nil, errors.New("context cancelled while waiting for decode response")
	}
}

// downloadChunks downloads chunks from the operator using the GetChunks() gRPC.
func (r *retrievalClient) downloadChunks(
	ctx context.Context,
	opID core.OperatorID,
	opInfo *core.OperatorInfo,
	blobKey corev2.BlobKey,
	quorumID core.QuorumID,
	replyChan chan *getChunksReply,
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
		replyChan <- &getChunksReply{
			OperatorID: opID,
			Err:        err,
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
		replyChan <- &getChunksReply{
			OperatorID: opID,
			Err:        err,
		}
		return
	}

	replyChan <- &getChunksReply{
		OperatorID: opID,
		Reply:      reply,
	}
}

// deserializeAndVerifyChunks deserializes the chunks from the GetChunksReply and sends them to the chunksChan.
func (r *retrievalClient) deserializeAndVerifyChunks(
	operatorID core.OperatorID,
	assignments map[core.OperatorID]corev2.Assignment,
	getChunksReply *grpcnode.GetChunksReply,
	blobCommitments encoding.BlobCommitments,
	encodingParams encoding.EncodingParams,
	replyChan chan clients.RetrievedChunks,
) {

	chunks := make([]*encoding.Frame, len(getChunksReply.GetChunks()))
	for i, data := range getChunksReply.GetChunks() {
		chunk, err := new(encoding.Frame).DeserializeGnark(data)
		if err != nil {
			replyChan <- clients.RetrievedChunks{
				OperatorID: operatorID,
				Err:        err,
				Chunks:     nil,
			}
			return
		}

		chunks[i] = chunk
	}

	assignment, ok := assignments[operatorID]
	if !ok {
		replyChan <- clients.RetrievedChunks{
			OperatorID: operatorID,
			Err:        fmt.Errorf("no assignment to operator %s", operatorID.Hex()),
		}
	}

	assignmentIndices := make([]uint, len(assignment.GetIndices()))
	for i, index := range assignment.GetIndices() {
		assignmentIndices[i] = uint(index)
	}

	err := r.verifier.VerifyFrames(chunks, assignmentIndices, blobCommitments, encodingParams)
	if err != nil {
		r.logger.Warn("failed to verify chunks from operator",
			"operator", operatorID.Hex(),
			"err", err)
		replyChan <- clients.RetrievedChunks{
			OperatorID: operatorID,
			Err:        err,
		}
		return
	} else {
		r.logger.Info("verified chunks from operator", "operator", operatorID.Hex())
	}

	replyChan <- clients.RetrievedChunks{
		OperatorID: operatorID,
		Err:        nil,
		Chunks:     chunks,
	}
}
