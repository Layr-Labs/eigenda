package hashing

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"
)

// HashGetPaymentStateRequest hashes the given GetPaymentStateRequest from accountId and timestamp
func HashGetPaymentStateRequest(accountId common.Address, timestamp uint64) ([]byte, error) {
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
