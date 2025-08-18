package internal

import (
	"fmt"

	grpcnode "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
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
	verifier encoding.Verifier,
) ChunkDeserializer

var _ ChunkDeserializer = &chunkDeserializer{}

// chunkDeserializer is a standard implementation of the ChunkDeserializer interface.
type chunkDeserializer struct {
	assignments map[core.OperatorID]v2.Assignment
	verifier    encoding.Verifier
}

var _ ChunkDeserializerFactory = NewChunkDeserializer

// NewChunkDeserializer creates a new ChunkDeserializer instance.
func NewChunkDeserializer(
	assignments map[core.OperatorID]v2.Assignment,
	verifier encoding.Verifier,
) ChunkDeserializer {
	return &chunkDeserializer{
		assignments: assignments,
		verifier:    verifier,
	}
}

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

	assignmentIndices := make([]uint, len(assignment.GetIndices()))
	for i, index := range assignment.GetIndices() {
		assignmentIndices[i] = uint(index)
	}

	err := d.verifier.VerifyFrames(chunks, assignmentIndices, *blobCommitments, *encodingParams)
	if err != nil {
		return nil, fmt.Errorf("failed to verify chunks from operator %s: %w", operatorID.Hex(), err)
	}

	return chunks, nil
}
