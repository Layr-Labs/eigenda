package core

import (
	"errors"
	"fmt"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/encoding"
)

var (
	ErrChunkLengthMismatch = errors.New("chunk length mismatch")
	ErrBlobQuorumSkip      = errors.New("blob skipped for a quorum before verification")
)

type ShardValidator interface {
	ValidateBatch(*BatchHeader, []*BlobMessage, *OperatorState, common.WorkerPool) error
	ValidateBlobs(blobs []*BlobMessage, operatorState *OperatorState, pool common.WorkerPool) error
	UpdateOperatorID(OperatorID)
}

// shardValidator implements the validation logic that a DA node should apply to its received data
type shardValidator struct {
	verifier   encoding.Verifier
	assignment AssignmentCoordinator
	chainState ChainState
	operatorID OperatorID
}

func NewShardValidator(v encoding.Verifier, asgn AssignmentCoordinator, cst ChainState, operatorID OperatorID) ShardValidator {
	return &shardValidator{
		verifier:   v,
		assignment: asgn,
		chainState: cst,
		operatorID: operatorID,
	}
}

func (v *shardValidator) validateBlobQuorum(quorumHeader *BlobQuorumInfo, blob *BlobMessage, operatorState *OperatorState) ([]*encoding.Frame, *Assignment, *encoding.EncodingParams, error) {
	if err := ValidateSecurityParam(uint32(quorumHeader.ConfirmationThreshold), uint32(quorumHeader.AdversaryThreshold)); err != nil {
		return nil, nil, nil, err
	}

	// Check if the operator is a member of the quorum
	if _, ok := operatorState.Operators[quorumHeader.QuorumID]; !ok {
		return nil, nil, nil, fmt.Errorf("%w: operator %s is not a member of quorum %d", ErrBlobQuorumSkip, v.operatorID.Hex(), quorumHeader.QuorumID)
	}

	// Get the assignments for the quorum
	assignment, info, err := v.assignment.GetOperatorAssignment(operatorState, blob.BlobHeader, quorumHeader.QuorumID, v.operatorID)
	if err != nil {
		return nil, nil, nil, err
	}

	// Validate the number of chunks
	if assignment.NumChunks == 0 {
		return nil, nil, nil, fmt.Errorf("%w: operator %s has no chunks in quorum %d", ErrBlobQuorumSkip, v.operatorID.Hex(), quorumHeader.QuorumID)
	}
	if assignment.NumChunks != uint(len(blob.Bundles[quorumHeader.QuorumID])) {
		return nil, nil, nil, fmt.Errorf("number of chunks (%d) does not match assignment (%d) for quorum %d", len(blob.Bundles[quorumHeader.QuorumID]), assignment.NumChunks, quorumHeader.QuorumID)
	}

	// Validate the chunkLength against the confirmation and adversary threshold parameters
	ok, err := v.assignment.ValidateChunkLength(operatorState, blob.BlobHeader.Length, quorumHeader)
	if err != nil || !ok {
		return nil, nil, nil, fmt.Errorf("invalid chunk length: %w", err)
	}

	// Get the chunk length
	chunks := blob.Bundles[quorumHeader.QuorumID]
	for _, chunk := range chunks {
		if uint(chunk.Length()) != quorumHeader.ChunkLength {
			return nil, nil, nil, fmt.Errorf("%w: chunk length (%d) does not match quorum header (%d) for quorum %d", ErrChunkLengthMismatch, chunk.Length(), quorumHeader.ChunkLength, quorumHeader.QuorumID)
		}
	}

	// Check the received chunks against the commitment
	params := encoding.ParamsFromMins(quorumHeader.ChunkLength, info.TotalChunks)
	if params.ChunkLength != uint64(quorumHeader.ChunkLength) {
		return nil, nil, nil, fmt.Errorf("%w: chunk length from encoding parameters (%d) does not match quorum header (%d)", ErrChunkLengthMismatch, params.ChunkLength, quorumHeader.ChunkLength)
	}

	return chunks, &assignment, &params, nil
}

func (v *shardValidator) UpdateOperatorID(operatorID OperatorID) {
	v.operatorID = operatorID
}

func (v *shardValidator) ValidateBatch(batchHeader *BatchHeader, blobs []*BlobMessage, operatorState *OperatorState, pool common.WorkerPool) error {
	headers := make([]*BlobHeader, len(blobs))
	for i, blob := range blobs {
		headers[i] = blob.BlobHeader
	}
	err := ValidateBatchHeaderRoot(batchHeader, headers)
	if err != nil {
		return err
	}

	return v.ValidateBlobs(blobs, operatorState, pool)
}

func (v *shardValidator) ValidateBlobs(blobs []*BlobMessage, operatorState *OperatorState, pool common.WorkerPool) error {
	var err error
	subBatchMap := make(map[encoding.EncodingParams]*encoding.SubBatch)
	blobCommitmentList := make([]encoding.BlobCommitments, len(blobs))

	for k, blob := range blobs {
		if len(blob.Bundles) != len(blob.BlobHeader.QuorumInfos) {
			return fmt.Errorf("number of bundles (%d) does not match number of quorums (%d)", len(blob.Bundles), len(blob.BlobHeader.QuorumInfos))
		}

		// Saved for the blob length validation
		blobCommitmentList[k] = blob.BlobHeader.BlobCommitments

		// for each quorum
		for _, quorumHeader := range blob.BlobHeader.QuorumInfos {
			chunks, assignment, params, err := v.validateBlobQuorum(quorumHeader, blob, operatorState)
			if errors.Is(err, ErrBlobQuorumSkip) {
				continue
			} else if err != nil {
				return err
			} else {
				// Check the received chunks against the commitment
				blobIndex := 0
				subBatch, ok := subBatchMap[*params]
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
					subBatchMap[*params] = &encoding.SubBatch{
						Samples:  samples,
						NumBlobs: 1,
					}
				} else {
					subBatch.Samples = append(subBatch.Samples, samples...)
					subBatch.NumBlobs += 1
				}
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
			v.VerifyBlobLengthWorker(blobCommitments, out)
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

func (v *shardValidator) VerifyBlobLengthWorker(blobCommitments encoding.BlobCommitments, out chan error) {
	err := v.verifier.VerifyBlobLength(blobCommitments)
	if err != nil {
		out <- err
		return
	}

	out <- nil
}

func ValidateBatchHeaderRoot(batchHeader *BatchHeader, blobHeaders []*BlobHeader) error {
	// Check the batch header root
	derivedHeader := &BatchHeader{}

	_, err := derivedHeader.SetBatchRoot(blobHeaders)
	if err != nil {
		return fmt.Errorf("failed to compute batch header root: %w", err)
	}

	if batchHeader.BatchRoot != derivedHeader.BatchRoot {
		return fmt.Errorf("batch header root does not match computed root")
	}

	return nil

}
