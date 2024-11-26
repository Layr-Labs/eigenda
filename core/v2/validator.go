package v2

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
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
		if len(blob.Bundles) != len(blob.BlobHeader.QuorumNumbers) {
			return fmt.Errorf("number of bundles (%d) does not match number of quorums (%d)", len(blob.Bundles), len(blob.BlobHeader.QuorumNumbers))
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
			params, err := blob.BlobHeader.GetEncodingParams(blobParams)
			if err != nil {
				return err
			}

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
