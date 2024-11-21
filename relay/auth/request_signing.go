package auth

import (
	"encoding/binary"
	pb "github.com/Layr-Labs/eigenda/api/grpc/relay"
	"github.com/Layr-Labs/eigenda/core"
	"golang.org/x/crypto/sha3"
)

var (
	iByte = []byte{0x69}
	rByte = []byte{0x72}
)

// HashGetChunksRequest hashes the given GetChunksRequest.
func HashGetChunksRequest(request *pb.GetChunksRequest) []byte {

	// Protobuf serialization is non-deterministic, so we can't just hash the
	// serialized bytes. Instead, we have to define our own hashing function.

	hasher := sha3.NewLegacyKeccak256()

	hasher.Write(request.GetOperatorId())
	for _, chunkRequest := range request.GetChunkRequests() {
		if chunkRequest.GetByIndex() != nil {
			getByIndex := chunkRequest.GetByIndex()
			hasher.Write(iByte)
			hasher.Write(getByIndex.BlobKey)
			for _, index := range getByIndex.ChunkIndices {
				indexBytes := make([]byte, 4)
				binary.BigEndian.PutUint32(indexBytes, index)
				hasher.Write(indexBytes)
			}
		} else {
			getByRange := chunkRequest.GetByRange()
			hasher.Write(rByte)
			hasher.Write(getByRange.BlobKey)

			startBytes := make([]byte, 4)
			binary.BigEndian.PutUint32(startBytes, getByRange.StartIndex)
			hasher.Write(startBytes)

			endBytes := make([]byte, 4)
			binary.BigEndian.PutUint32(endBytes, getByRange.EndIndex)
			hasher.Write(endBytes)
		}
	}

	return hasher.Sum(nil)
}

// SignGetChunksRequest signs the given GetChunksRequest with the given private key. Does not
// write the signature into the request.
func SignGetChunksRequest(keys *core.KeyPair, request *pb.GetChunksRequest) []byte {
	hash := HashGetChunksRequest(request)
	signature := keys.SignMessage(([32]byte)(hash))
	return signature.G1Point.Serialize()
}
