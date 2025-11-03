package internal

import (
	"fmt"

	grpcnode "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg/verifier"
)

// A ChunkDeserializer is responsible for deserializing binary chunks. Will only return chunks if they are valid.
type ChunkDeserializer interface {

	// DeserializeAndVerify deserializes the binary chunks as received from a validator and verifies them.
	DeserializeAndVerify(
		blobKey v2.BlobKey,
		operatorID core.OperatorID,
		getChunksReply *grpcnode.GetChunksReply,
		blobCommitments *encoding.BlobCommitments,
		encodingParams *encoding.EncodingParams,
	) ([]*encoding.Frame, error)
}

// ChunkDeserializerFactory is a function that creates a new ChunkDeserializer instance.
type ChunkDeserializerFactory func(
	assignments map[core.OperatorID]v2.Assignment,
	verifier *verifier.Verifier,
) ChunkDeserializer

var _ ChunkDeserializer = &chunkDeserializer{}

// chunkDeserializer is a standard implementation of the ChunkDeserializer interface.
type chunkDeserializer struct {
	assignments map[core.OperatorID]v2.Assignment
	verifier    *verifier.Verifier
}

var _ ChunkDeserializerFactory = NewChunkDeserializer

// NewChunkDeserializer creates a new ChunkDeserializer instance.
func NewChunkDeserializer(
	assignments map[core.OperatorID]v2.Assignment,
	verifier *verifier.Verifier,
) ChunkDeserializer {
	return &chunkDeserializer{
		assignments: assignments,
		verifier:    verifier,
	}
}

// assume all the chunks come from one blob. In theory, universal verification
// works as long as all chunk lengths are equal, and we can find the right root of
// unities.
func (d *chunkDeserializer) DeserializeAndVerify(
	_ v2.BlobKey, // used for unit tests
	operatorID core.OperatorID,
	getChunksReply *grpcnode.GetChunksReply,
	blobCommitments *encoding.BlobCommitments,
	encodingParams *encoding.EncodingParams,
) ([]*encoding.Frame, error) {

	chunks := make([]*encoding.Frame, len(getChunksReply.GetChunks()))
	for i, data := range getChunksReply.GetChunks() {
		chunk, err := new(encoding.Frame).DeserializeGnark(data)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize chunk from operator %s: %w", operatorID.Hex(), err)
		}

		chunks[i] = chunk
	}

	assignment := d.assignments[operatorID]

	assignmentIndices := make([]core.ChunkNumber, len(assignment.GetIndices()))
	for i, index := range assignment.GetIndices() {
		assignmentIndices[i] = core.ChunkNumber(index)
	}

	samples := make([]encoding.Sample, len(chunks))
	for ind := range chunks {
		samples[ind] = encoding.Sample{
			Commitment:      blobCommitments.Commitment,
			Chunk:           chunks[ind],
			AssignmentIndex: assignmentIndices[ind],
			BlobIndex:       0, // there is only 1 blob
		}
	}

	// verify all chunks for operator using universal verification, it reduces the complexity from
	// n*m to n + m, where n is the number of chunks, and m is the length of each chunk in field elements
	// For theory, see https://ethresear.ch/t/a-universal-verification-equation-for-data-availability-sampling/13240
	err := d.verifier.UniversalVerifySubBatch(
		*encodingParams,
		samples,
		1, // only verify one blob
	)
	if err != nil {
		return nil, fmt.Errorf("failed to verify chunks from operator %s: %w", operatorID.Hex(), err)
	}

	return chunks, nil
}
