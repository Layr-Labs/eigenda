package core

import (
	"errors"
)

var (
	ErrChunkLengthMismatch = errors.New("chunk length mismatch")
	ErrInvalidHeader       = errors.New("invalid header")
	ErrBlobQuorumSkip      = errors.New("blob skipped for a quorum before verification")
)

type ChunkValidator interface {
	ValidateBatch([]*BlobMessage, *OperatorState) error
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

// preprocessBlob for each Quorum
func (v *chunkValidator) preprocessBlob(quorumHeader *BlobQuorumInfo, blob *BlobMessage, operatorState *OperatorState) ([]*Chunk, *Assignment, *EncodingParams, error) {
	if quorumHeader.AdversaryThreshold >= quorumHeader.QuorumThreshold {
		return nil, nil, nil, errors.New("invalid header: quorum threshold does not exceed adversary threshold")
	}

	// Check if the operator is a member of the quorum
	if _, ok := operatorState.Operators[quorumHeader.QuorumID]; !ok {
		return nil, nil, nil, ErrBlobQuorumSkip
	}

	// Get the assignments for the quorum
	assignment, info, err := v.assignment.GetOperatorAssignment(operatorState, quorumHeader.QuorumID, quorumHeader.QuantizationFactor, v.operatorID)
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

	chunkLength, err := v.assignment.GetChunkLengthFromHeader(operatorState, quorumHeader)
	if err != nil {
		return nil, nil, nil, err
	}

	// Validate the chunkLength against the quorum and adversary threshold parameters
	numOperators := uint(len(operatorState.Operators[quorumHeader.QuorumID]))
	minChunkLength, err := v.assignment.GetMinimumChunkLength(numOperators, blob.BlobHeader.BlobCommitments.Length, quorumHeader.QuantizationFactor, quorumHeader.QuorumThreshold, quorumHeader.AdversaryThreshold)
	if err != nil {
		return nil, nil, nil, err
	}
	params, err := GetEncodingParams(minChunkLength, info.TotalChunks)
	if err != nil {
		return nil, nil, nil, err
	}

	if params.ChunkLength != chunkLength {
		return nil, nil, nil, errors.New("number of chunks does not match assignment")
	}

	// Get the chunk length
	chunks := blob.Bundles[quorumHeader.QuorumID]
	for _, chunk := range chunks {
		if uint(chunk.Length()) != chunkLength {
			return nil, nil, nil, ErrChunkLengthMismatch
		}
	}

	// Validate the chunk length
	if chunkLength*quorumHeader.QuantizationFactor*numOperators != quorumHeader.EncodedBlobLength {
		return nil, nil, nil, ErrInvalidHeader
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
		chunks, assignment, params, err := v.preprocessBlob(quorumHeader, blob, operatorState)
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

func (v *chunkValidator) ValidateBatch(blobs []*BlobMessage, operatorState *OperatorState) error {
	subBatchMap := make(map[EncodingParams]*SubBatch)

	for _, blob := range blobs {
		if len(blob.Bundles) != len(blob.BlobHeader.QuorumInfos) {
			return errors.New("number of bundles does not match number of quorums")
		}

		// Validate the blob length
		err := v.encoder.VerifyBlobLength(blob.BlobHeader.BlobCommitments)
		if err != nil {
			return err
		}
		// for each quorum
		for _, quorumHeader := range blob.BlobHeader.QuorumInfos {
			// Check if the operator is a member of the quorum
			chunks, assignment, params, err := v.preprocessBlob(quorumHeader, blob, operatorState)
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
	numSubBatch := len(subBatchMap)
	out := make(chan error, numSubBatch)
	for params, subBatch := range subBatchMap {
		params := params
		subBatch := subBatch
		go v.universalVerifyWorker(params, subBatch, out)
	}

	for i := 0; i < numSubBatch; i++ {
		err := <-out
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *chunkValidator) universalVerifyWorker(params EncodingParams, subBatch *SubBatch, out chan error) {

	err := v.encoder.UniversalVerifyChunks(params, subBatch.Samples, subBatch.NumBlobs)
	if err != nil {
		out <- err
		return
	}

	out <- nil
}
