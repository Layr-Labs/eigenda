package core

import (
	"errors"
	"fmt"

	"github.com/Layr-Labs/eigenda/common"
)

var (
	ErrChunkLengthMismatch = errors.New("chunk length mismatch")
	ErrInvalidHeader       = errors.New("invalid header")
	ErrBlobQuorumSkip      = errors.New("blob skipped for a quorum before verification")
)

type ChunkValidator interface {
	ValidateBatch([]*BlobMessage, *OperatorState, common.WorkerPool) error
	ValidateBlob(*BlobMessage, *OperatorState) error
	UpdateOperatorID(OperatorID)
}

// chunkValidator implements the validation logic that a DA node should apply to its recieved chunks
type chunkValidator struct {
	encoder    Encoder
	assignment AssignmentCoordinator
	chainState ChainState
	operatorID OperatorID
}

func NewChunkValidator(enc Encoder, asgn AssignmentCoordinator, cst ChainState, operatorID OperatorID) ChunkValidator {
	return &chunkValidator{
		encoder:    enc,
		assignment: asgn,
		chainState: cst,
		operatorID: operatorID,
	}
}

func (v *chunkValidator) validateBlobQuorum(quorumHeader *BlobQuorumInfo, blob *BlobMessage, operatorState *OperatorState) ([]*Chunk, *Assignment, *EncodingParams, error) {
	if quorumHeader.AdversaryThreshold >= quorumHeader.QuorumThreshold {
		return nil, nil, nil, errors.New("invalid header: quorum threshold does not exceed adversary threshold")
	}

	// Check if the operator is a member of the quorum
	if _, ok := operatorState.Operators[quorumHeader.QuorumID]; !ok {
		return nil, nil, nil, ErrBlobQuorumSkip
	}

	// Get the assignments for the quorum
	assignment, info, err := v.assignment.GetOperatorAssignment(operatorState, blob.BlobHeader, quorumHeader.QuorumID, v.operatorID)
	if err != nil {
		return nil, nil, nil, err
	}

	// Validate the number of chunks
	if assignment.NumChunks == 0 {
		return nil, nil, nil, ErrBlobQuorumSkip
	}
	if assignment.NumChunks != uint(len(blob.Bundles[quorumHeader.QuorumID])) {
		return nil, nil, nil, errors.New("number of chunks does not match assignment")
	}

	// Validate the chunkLength against the quorum and adversary threshold parameters
	ok, err := v.assignment.ValidateChunkLength(operatorState, blob.BlobHeader, quorumHeader.QuorumID)
	if err != nil || !ok {
		return nil, nil, nil, fmt.Errorf("invalid chunk length: %w", err)
	}

	// Get the chunk length
	chunks := blob.Bundles[quorumHeader.QuorumID]
	for _, chunk := range chunks {
		if uint(chunk.Length()) != quorumHeader.ChunkLength {
			return nil, nil, nil, ErrChunkLengthMismatch
		}
	}

	// Check the received chunks against the commitment
	params, err := GetEncodingParams(quorumHeader.ChunkLength, info.TotalChunks)
	if err != nil {
		return nil, nil, nil, err
	}

	return chunks, &assignment, &params, nil
}

func (v *chunkValidator) ValidateBlob(blob *BlobMessage, operatorState *OperatorState) error {
	if len(blob.Bundles) != len(blob.BlobHeader.QuorumInfos) {
		return errors.New("number of bundles does not match number of quorums")
	}

	// Validate the blob length
	err := v.encoder.VerifyBlobLength(blob.BlobHeader.BlobCommitments)
	if err != nil {
		return err
	}

	for _, quorumHeader := range blob.BlobHeader.QuorumInfos {
		// preprocess validation info
		chunks, assignment, params, err := v.validateBlobQuorum(quorumHeader, blob, operatorState)
		if err == ErrBlobQuorumSkip {
			continue
		} else if err != nil {
			return err
		} else {
			// Check the received chunks against the commitment
			err = v.encoder.VerifyChunks(chunks, assignment.GetIndices(), blob.BlobHeader.BlobCommitments, *params)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (v *chunkValidator) UpdateOperatorID(operatorID OperatorID) {
	v.operatorID = operatorID
}

func (v *chunkValidator) ValidateBatch(blobs []*BlobMessage, operatorState *OperatorState, pool common.WorkerPool) error {
	subBatchMap := make(map[EncodingParams]*SubBatch)
	blobCommitmentList := make([]BlobCommitments, len(blobs))

	for k, blob := range blobs {
		if len(blob.Bundles) != len(blob.BlobHeader.QuorumInfos) {
			return errors.New("number of bundles does not match number of quorums")
		}

		// Saved for the blob length validation
		blobCommitmentList[k] = blob.BlobHeader.BlobCommitments

		// for each quorum
		for _, quorumHeader := range blob.BlobHeader.QuorumInfos {
			chunks, assignment, params, err := v.validateBlobQuorum(quorumHeader, blob, operatorState)
			if err == ErrBlobQuorumSkip {
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
				samples := make([]Sample, len(chunks))
				for ind := range chunks {
					samples[ind] = Sample{
						Commitment:      blob.BlobHeader.BlobCommitments.Commitment,
						Chunk:           chunks[ind],
						AssignmentIndex: uint(indices[ind]),
						BlobIndex:       blobIndex,
					}
				}

				// update subBatch
				if !ok {
					subBatchMap[*params] = &SubBatch{
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

	for i := 0; i < numResult; i++ {
		err := <-out
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *chunkValidator) universalVerifyWorker(params EncodingParams, subBatch *SubBatch, out chan error) {

	err := v.encoder.UniversalVerifySubBatch(params, subBatch.Samples, subBatch.NumBlobs)
	if err != nil {
		out <- err
		return
	}

	out <- nil
}

func (v *chunkValidator) VerifyBlobLengthWorker(blobCommitments BlobCommitments, out chan error) {
	err := v.encoder.VerifyBlobLength(blobCommitments)
	if err != nil {
		out <- err
		return
	}

	out <- nil
}
