package hashing

import (
	"fmt"
	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"
)

// HashGetPaymentStateRequestFromFields hashes the given GetPaymentStateRequest from accountId and timestamp
func HashGetPaymentStateRequestFromFields(accountId common.Address, timestamp uint64) ([]byte, error) {
	hasher := sha3.NewLegacyKeccak256()

	// Hash the accountId
	err := hashByteArray(hasher, accountId.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to hash account id: %w", err)
	}

	// Hash the timestamp
	hashUint64(hasher, timestamp)

	return hasher.Sum(nil), nil
}

// HashGetPaymentStateRequestFromRequest hashes the given GetPaymentStateRequest from request
func HashGetPaymentStateRequestFromRequest(request *pb.GetPaymentStateRequest) ([]byte, error) {
	accountId := common.HexToAddress(request.GetAccountId())
	timestamp := request.GetTimestamp()

	return HashGetPaymentStateRequestFromFields(accountId, timestamp)
}
