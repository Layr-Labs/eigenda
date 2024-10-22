package v2

import (
	"context"
	"errors"
	"fmt"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding"
)

var (
	ErrChunkLengthMismatch = errors.New("chunk length mismatch")
	ErrBlobQuorumSkip      = errors.New("blob skipped for a quorum before verification")
)

type BlobShard struct {
	BlobCertificate
	Chunks map[core.QuorumID][]*encoding.Frame
}

// shardValidator implements the validation logic that a DA node should apply to its received data
type ShardValidator struct {
	verifier   encoding.Verifier
	chainState core.ChainState
	operatorID core.OperatorID
}

func NewShardValidator(v encoding.Verifier, cst core.ChainState, operatorID core.OperatorID) *ShardValidator {
	return &ShardValidator{
		verifier:   v,
		chainState: cst,
		operatorID: operatorID,
	}
}

func (v *ShardValidator) validateBlobQuorum(quorum core.QuorumID, blob *BlobShard, operatorState *core.OperatorState) ([]*encoding.Frame, *Assignment, error) {

	// Check if the operator is a member of the quorum
	if _, ok := operatorState.Operators[quorum]; !ok {
		return nil, nil, fmt.Errorf("%w: operator %s is not a member of quorum %d", ErrBlobQuorumSkip, v.operatorID.Hex(), quorum)
	}

	// Get the assignments for the quorum
	assignment, err := GetAssignment(operatorState, blob.Version, quorum, v.operatorID)
	if err != nil {
		return nil, nil, err
	}

	// Validate the number of chunks
	if assignment.NumChunks == 0 {
		return nil, nil, fmt.Errorf("%w: operator %s has no chunks in quorum %d", ErrBlobQuorumSkip, v.operatorID.Hex(), quorum)
	}
	if assignment.NumChunks != uint32(len(blob.Chunks[quorum])) {
		return nil, nil, fmt.Errorf("number of chunks (%d) does not match assignment (%d) for quorum %d", len(blob.Chunks[quorum]), assignment.NumChunks, quorum)
	}

	// Validate the chunkLength against the confirmation and adversary threshold parameters
	chunkLength, err := GetChunkLength(blob.Version, uint32(blob.BlobHeader.Length))
	if err != nil {
		return nil, nil, fmt.Errorf("invalid chunk length: %w", err)
	}

	// Get the chunk length
	chunks := blob.Chunks[quorum]
	for _, chunk := range chunks {
		if uint32(chunk.Length()) != chunkLength {
			return nil, nil, fmt.Errorf("%w: chunk length (%d) does not match quorum header (%d) for quorum %d", ErrChunkLengthMismatch, chunk.Length(), chunkLength, quorum)
		}
	}

	return chunks, &assignment, nil
}

func (v *ShardValidator) ValidateBlobs(ctx context.Context, blobs []*BlobShard, pool common.WorkerPool) error {
	var err error
	subBatchMap := make(map[encoding.EncodingParams]*encoding.SubBatch)
	blobCommitmentList := make([]encoding.BlobCommitments, len(blobs))

	for k, blob := range blobs {
		if len(blob.Chunks) != len(blob.BlobHeader.QuorumNumbers) {
			return fmt.Errorf("number of bundles (%d) does not match number of quorums (%d)", len(blob.Chunks), len(blob.BlobHeader.QuorumNumbers))
		}

		state, err := v.chainState.GetOperatorState(ctx, uint(blob.ReferenceBlockNumber), blob.BlobHeader.QuorumNumbers)
		if err != nil {
			return err
		}

		// Saved for the blob length validation
		blobCommitmentList[k] = blob.BlobHeader.BlobCommitments

		// for each quorum
		for _, quorum := range blob.BlobHeader.QuorumNumbers {
			chunks, assignment, err := v.validateBlobQuorum(quorum, blob, state)
			if err != nil {
				return err
			}
			// TODO: Define params for the blob
			params, err := blob.GetEncodingParams()
			if err != nil {
				return err
			}

			if errors.Is(err, ErrBlobQuorumSkip) {
				continue
			} else if err != nil {
				return err
			} else {
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

func (v *ShardValidator) universalVerifyWorker(params encoding.EncodingParams, subBatch *encoding.SubBatch, out chan error) {

	err := v.verifier.UniversalVerifySubBatch(params, subBatch.Samples, subBatch.NumBlobs)
	if err != nil {
		out <- err
		return
	}

	out <- nil
}

func (v *ShardValidator) VerifyBlobLengthWorker(blobCommitments encoding.BlobCommitments, out chan error) {
	err := v.verifier.VerifyBlobLength(blobCommitments)
	if err != nil {
		out <- err
		return
	}

	out <- nil
}
