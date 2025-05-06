package mock

import (
	"github.com/Layr-Labs/eigenda/api/clients/v2/validator/internal"
	grpcnode "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
)

var _ internal.ChunkDeserializer = (*MockChunkDeserializer)(nil)

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
func NewMockChunkDeserializerFactory(deserializer internal.ChunkDeserializer) internal.ChunkDeserializerFactory {
	return func(
		assignments map[core.OperatorID]v2.Assignment,
		verifier encoding.Verifier,
	) internal.ChunkDeserializer {
		return deserializer
	}
}
