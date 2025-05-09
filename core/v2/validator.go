package v2

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	pbvalidator "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ErrChunkLengthMismatch = errors.New("chunk length mismatch")
	ErrBlobQuorumSkip      = errors.New("blob skipped for a quorum before verification")
)

type ShardValidator interface {
	ValidateBatchHeader(ctx context.Context, header *BatchHeader, blobCerts []*BlobCertificate) error
	ValidateBlobs(ctx context.Context, blobs []*BlobShard, blobVersionParams *BlobVersionParameterMap, pool common.WorkerPool, state *core.OperatorState) error
}

type BlobShard struct {
	*BlobCertificate
	Bundles core.Bundles
}

// shardValidator implements the validation logic that a DA node should apply to its received data
type shardValidator struct {
	verifier   encoding.Verifier
	operatorID core.OperatorID
	logger     logging.Logger
}

var _ ShardValidator = (*shardValidator)(nil)

func NewShardValidator(v encoding.Verifier, operatorID core.OperatorID, logger logging.Logger) *shardValidator {
	return &shardValidator{
		verifier:   v,
		operatorID: operatorID,
		logger:     logger,
	}
}

func (v *shardValidator) validateBlobQuorum(quorum core.QuorumID, blob *BlobShard, blobParams *core.BlobVersionParameters, operatorState *core.OperatorState) ([]*encoding.Frame, *Assignment, error) {

	// Check if the operator is a member of the quorum
	if _, ok := operatorState.Operators[quorum]; !ok {
		return nil, nil, fmt.Errorf("%w: operator %s is not a member of quorum %d", ErrBlobQuorumSkip, v.operatorID.Hex(), quorum)
	}

	// Get the assignments for the quorum
	assignment, err := GetAssignment(operatorState, blobParams, quorum, v.operatorID)
	if err != nil {
		return nil, nil, err
	}

	// Validate the number of chunks
	if assignment.NumChunks == 0 {
		return nil, nil, fmt.Errorf("%w: operator %s has no chunks in quorum %d", ErrBlobQuorumSkip, v.operatorID.Hex(), quorum)
	}
	if assignment.NumChunks != uint32(len(blob.Bundles[quorum])) {
		return nil, nil, fmt.Errorf("number of chunks (%d) does not match assignment (%d) for quorum %d", len(blob.Bundles[quorum]), assignment.NumChunks, quorum)
	}

	// Get the chunk length
	chunkLength, err := GetChunkLength(uint32(blob.BlobHeader.BlobCommitments.Length), blobParams)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid chunk length: %w", err)
	}

	chunks := blob.Bundles[quorum]
	for _, chunk := range chunks {
		if uint32(chunk.Length()) != chunkLength {
			return nil, nil, fmt.Errorf("%w: chunk length (%d) does not match quorum header (%d) for quorum %d", ErrChunkLengthMismatch, chunk.Length(), chunkLength, quorum)
		}
	}

	return chunks, &assignment, nil
}

func (v *shardValidator) ValidateBatchHeader(ctx context.Context, header *BatchHeader, blobCerts []*BlobCertificate) error {
	if header == nil {
		return fmt.Errorf("batch header is nil")
	}

	if len(blobCerts) == 0 {
		return fmt.Errorf("no blob certificates")
	}

	tree, err := BuildMerkleTree(blobCerts)
	if err != nil {
		return fmt.Errorf("failed to build merkle tree: %v", err)
	}

	if !bytes.Equal(tree.Root(), header.BatchRoot[:]) {
		return fmt.Errorf("batch root does not match")
	}

	return nil
}

func filterOperatorStateQuorums(operatorState *core.OperatorState, quorums []core.QuorumID) map[core.QuorumID]struct{} {
	filtered := make(map[core.QuorumID]struct{})
	for _, quorum := range quorums {
		if _, ok := operatorState.Operators[quorum]; ok {
			filtered[quorum] = struct{}{}
		}
	}

	return filtered
}

func (v *shardValidator) ValidateBlobs(ctx context.Context, blobs []*BlobShard, blobVersionParams *BlobVersionParameterMap, pool common.WorkerPool, state *core.OperatorState) error {
	if len(blobs) == 0 {
		return fmt.Errorf("no blobs")
	}

	if blobVersionParams == nil {
		return fmt.Errorf("blob version params is nil")
	}

	var err error
	subBatchMap := make(map[encoding.EncodingParams]*encoding.SubBatch)
	blobCommitmentList := make([]encoding.BlobCommitments, len(blobs))

	for k, blob := range blobs {
		relevantQuorums := filterOperatorStateQuorums(state, blob.BlobHeader.QuorumNumbers)
		if len(blob.Bundles) != len(relevantQuorums) {
			return fmt.Errorf("number of bundles (%d) does not match number of relevant quorums (%d)", len(blob.Bundles), len(relevantQuorums))
		}

		// Saved for the blob length validation
		blobCommitmentList[k] = blob.BlobHeader.BlobCommitments

		// for each quorum
		for _, quorum := range blob.BlobHeader.QuorumNumbers {
			blobParams, ok := blobVersionParams.Get(blob.BlobHeader.BlobVersion)
			if !ok {
				return fmt.Errorf("blob version %d not found", blob.BlobHeader.BlobVersion)
			}
			chunks, assignment, err := v.validateBlobQuorum(quorum, blob, blobParams, state)
			if errors.Is(err, ErrBlobQuorumSkip) {
				v.logger.Warn("Skipping blob for quorum", "quorum", quorum, "err", err)
				continue
			} else if err != nil {
				return err
			}

			// TODO: Define params for the blob
			params, err := GetEncodingParams(blob.BlobHeader.BlobCommitments.Length, blobParams)

			if err != nil {
				return err
			}

			// Check the received chunks against the commitment
			blobIndex := 0
			subBatch, ok := subBatchMap[params]
			if ok {
				blobIndex = subBatch.NumBlobs
			}

			indices := assignment.GetIndices()
			samples := make([]encoding.Sample, len(chunks))
			for ind := range chunks {
				samples[ind] = encoding.Sample{
					Commitment:      blob.BlobHeader.BlobCommitments.Commitment,
					Chunk:           chunks[ind],
					AssignmentIndex: uint(indices[ind]),
					BlobIndex:       blobIndex,
				}
			}

			// update subBatch
			if !ok {
				subBatchMap[params] = &encoding.SubBatch{
					Samples:  samples,
					NumBlobs: 1,
				}
			} else {
				subBatch.Samples = append(subBatch.Samples, samples...)
				subBatch.NumBlobs += 1
			}
		}
	}

	// Parallelize the universal verification for each subBatch
	numResult := len(subBatchMap) + len(blobCommitmentList)
	// create a channel to accept results, we don't use stop
	out := make(chan error, numResult)

	// parallelize subBatch verification
	for params, subBatch := range subBatchMap {
		params := params
		subBatch := subBatch
		pool.Submit(func() {
			v.universalVerifyWorker(params, subBatch, out)
		})
	}

	// parallelize length proof verification
	for _, blobCommitments := range blobCommitmentList {
		blobCommitments := blobCommitments
		pool.Submit(func() {
			v.verifyBlobLengthWorker(blobCommitments, out)
		})
	}
	// check if commitments are equivalent
	err = v.verifier.VerifyCommitEquivalenceBatch(blobCommitmentList)
	if err != nil {
		return err
	}

	for i := 0; i < numResult; i++ {
		err := <-out
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *shardValidator) universalVerifyWorker(params encoding.EncodingParams, subBatch *encoding.SubBatch, out chan error) {

	err := v.verifier.UniversalVerifySubBatch(params, subBatch.Samples, subBatch.NumBlobs)
	if err != nil {
		out <- err
		return
	}

	out <- nil
}

func (v *shardValidator) verifyBlobLengthWorker(blobCommitments encoding.BlobCommitments, out chan error) {
	err := v.verifier.VerifyBlobLength(blobCommitments)
	if err != nil {
		out <- err
		return
	}

	out <- nil
}

// GetOperatorVerboseState returns the verbose state of all operators within the supplied quorums including their node info.
// The returned state is for the block number supplied.
func GetOperatorVerboseState(ctx context.Context, stakesWithSocket core.OperatorStakesWithSocket, quorums []core.QuorumID, blockNumber uint32) (core.OperatorStateVerbose, error) {

	loggerConfig := common.DefaultLoggerConfig()
	logger, err := common.NewLogger(loggerConfig)
	if err != nil {
		panic(err)
	}
	logger = logger.With("component", "GetOperatorVerboseState")
	logger.Infof("GetOperatorVerboseState: Starting for quorums %v at block %d", quorums, blockNumber)

	quorumBytes := make([]byte, len(quorums))
	for ind, quorum := range quorums {
		quorumBytes[ind] = byte(uint8(quorum))
	}

	// result is a struct{Operators [][]opstateretriever.OperatorStateRetrieverOperator; Sockets [][]string}
	// Operators is a [][]*opstateretriever.OperatorStake with the same length and order as quorumBytes, and then indexed by operator index
	// Sockets is a [][]string with the same length and order as quorumBytes, and then indexed by operator index
	// By contract definition, Operators and Sockets are parallel arrays
	logger.Infof("GetOperatorVerboseState: Calling GetOperatorStateWithSocket for %d quorums", len(quorums))
	state := make(core.OperatorStateVerbose, len(quorums))
	totalOperators := 0
	successfulNodeInfoFetches := 0
	failedNodeInfoFetches := 0

	for _, quorumID := range quorums {
		state[quorumID] = make(map[core.OperatorIndex]core.OperatorInfoVerbose, len(stakesWithSocket[quorumID]))
		logger.Infof("GetOperatorVerboseState: Processing %d operators for quorum %d", len(stakesWithSocket[quorumID]), quorumID)

		for j, op := range stakesWithSocket[quorumID] {
			totalOperators++
			operatorIndex := core.OperatorIndex(j)

			logger.Infof("GetOperatorVerboseState: Fetching node info for operator %s at %s", op.OperatorID, op.Socket)
			nodeVersion, err := GetNodeInfoFromSocket(ctx, op.Socket)
			if err != nil {
				failedNodeInfoFetches++
				logger.Warnf("Failed to fetch node version from %s for operator %s: %s", op.Socket, op.OperatorID, err)

				// Instead of failing completely, continue with nil NodeInfo for this operator
				state[quorumID][operatorIndex] = core.OperatorInfoVerbose{
					OperatorID: op.OperatorID,
					Socket:     op.Socket,
					Stake:      op.Stake,
					NodeInfo:   nil,
				}
				continue
			}
			successfulNodeInfoFetches++

			logger.Infof("GetOperatorVerboseState: Got node info %+v for operator %s", nodeVersion, op.OperatorID)
			state[quorumID][operatorIndex] = core.OperatorInfoVerbose{
				OperatorID: op.OperatorID,
				Socket:     op.Socket,
				Stake:      op.Stake,
				NodeInfo:   nodeVersion,
			}
		}
	}

	logger.Infof("GetOperatorVerboseState: Completed with %d total operators, %d successful node info fetches, %d failed fetches",
		totalOperators, successfulNodeInfoFetches, failedNodeInfoFetches)
	return state, nil
}

// GetNodeInfoFromSocket pings the operator's endpoint and returns NodeInfoReply
func GetNodeInfoFromSocket(ctx context.Context, socket core.OperatorSocket) (*pbvalidator.GetNodeInfoReply, error) {

	loggerConfig := common.DefaultLoggerConfig()
	logger, err := common.NewLogger(loggerConfig)
	if err != nil {
		panic(err)
	}
	logger = logger.With("component", "GetNodeInfoFromSocket")
	logger.Info("Getting node info from endpoint", "socket", socket)

	// host, _, _, v2DispersalPort, _, err := core.ParseOperatorSocket(socket)
	// if err != nil {
	// 	return nil, err
	// }
	// endpoint := fmt.Sprintf("%s:%s", host, v2DispersalPort)
	endpoint := socket.GetV2DispersalSocket()
	logger.Info("Getting node info from endpoint", "endpoint", endpoint)

	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("Failed to create GRPC client", "endpoint", endpoint, "error", err)
		return nil, err
	}
	defer conn.Close()

	client := pbvalidator.NewDispersalClient(conn)
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	logger.Info("Sending GetNodeInfo request", "endpoint", endpoint)
	resp, err := client.GetNodeInfo(ctxTimeout, &pbvalidator.GetNodeInfoRequest{})
	if err != nil {
		logger.Error("Failed to get node info", "endpoint", endpoint, "error", err)
		return nil, err
	}

	logger.Info("Successfully got node info", "endpoint", endpoint, "version", resp.Semver)
	return resp, nil
}
