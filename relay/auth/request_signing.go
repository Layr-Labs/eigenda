package auth

import (
	pb "github.com/Layr-Labs/eigenda/api/grpc/relay"
	"github.com/Layr-Labs/eigenda/api/hashing"
	"github.com/Layr-Labs/eigenda/core"
)

// SignGetChunksRequest signs the given GetChunksRequest with the given private key. Does not
// write the signature into the request.
func SignGetChunksRequest(keys *core.KeyPair, request *pb.GetChunksRequest) ([]byte, error) {
	hash, err := hashing.HashGetChunksRequest(request)
	if err != nil {
		return nil, err
	}
	signature := keys.SignMessage(([32]byte)(hash))
	return signature.G1Point.Serialize(), nil
}
