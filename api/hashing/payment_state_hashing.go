package hashing

import (
	"fmt"
	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"golang.org/x/crypto/sha3"
)

// HashGetPaymentStateRequest hashes the given GetPaymentStateRequest.
func HashGetPaymentStateRequest(request *pb.GetPaymentStateRequest) ([]byte, error) {
	hasher := sha3.NewLegacyKeccak256()

	// Hash the account ID
	err := hashByteArray(hasher, []byte(request.GetAccountId()))
	if err != nil {
		return nil, fmt.Errorf("failed to hash account id: %w", err)
	}

	// Hash the timestamp
	hashUint64(hasher, request.GetTimestamp())

	return hasher.Sum(nil), nil
}
