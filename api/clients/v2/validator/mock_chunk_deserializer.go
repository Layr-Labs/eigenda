package validator

import (
	grpcnode "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
)

var _ ChunkDeserializer = (*MockChunkDeserializer)(nil)

type MockChunkDeserializer struct {
	// A lambda function to be called when DeserializeAndVerify is called.
	DeserializeAndVerifyFunction func(
		blobKey v2.BlobKey,
		operatorID core.OperatorID,
		getChunksReply *grpcnode.GetChunksReply,
		blobCommitments *encoding.BlobCommitments,
		encodingParams *encoding.EncodingParams,
	) ([]*encoding.Frame, error)
}

func (m *MockChunkDeserializer) DeserializeAndVerify(
	blobKey v2.BlobKey,
	operatorID core.OperatorID,
	getChunksReply *grpcnode.GetChunksReply,
	blobCommitments *encoding.BlobCommitments,
	encodingParams *encoding.EncodingParams,
) ([]*encoding.Frame, error) {
	if m.DeserializeAndVerifyFunction == nil {
		return nil, nil
	}
	return m.DeserializeAndVerifyFunction(blobKey, operatorID, getChunksReply, blobCommitments, encodingParams)
}

// NewMockChunkDeserializerFactory creates a new ChunkDeserializerFactory that returns the provided deserializer.
func NewMockChunkDeserializerFactory(deserializer ChunkDeserializer) ChunkDeserializerFactory {
	return func(
		assignments map[core.OperatorID]v2.Assignment,
		verifier encoding.Verifier,
	) ChunkDeserializer {
		return deserializer
	}
}
