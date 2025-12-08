package hashing

import (
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenda/common"
	"golang.org/x/crypto/sha3"
)

// Version of the protobuf DisperseBlobRequest
type DisperseBlobRequestVersion = uint32

const (
	// The version from before request versioning was introduced.
	// Use of this version indicates that legacy hashing was used: signature is over blobKey.
	DisperseBlobRequestVersion0 DisperseBlobRequestVersion = 0

	// Introduces a new hash that the signature is over.
	// hash(domain || chainID || disperserID || blobKey).
	DisperseBlobRequestVersion1 DisperseBlobRequestVersion = 1
)

const DisperseBlobRequestDomain = "disperser.DisperseBlobRequest"

// Computes the hash for signing/verifying DisperseBlobRequest.
// useNewHashVersion determines the hashing algorithm:
// - false: hash = blobKey (legacy, for backward compatibility)
// - true: hash = Keccak256(domain || chainId || disperserId || blobKey)
func HashDisperseBlobRequest(
	useNewHashVersion bool,
	blobKey [32]byte,
	disperserId uint32,
	chainId *big.Int,
) ([]byte, error) {
	if !useNewHashVersion {
		// Legacy hashing: just the blobKey
		return blobKey[:], nil
	}

	if chainId == nil {
		return nil, fmt.Errorf("chainId is required for new hash version")
	}

	hasher := sha3.NewLegacyKeccak256()

	hasher.Write([]byte(DisperseBlobRequestDomain))
	hasher.Write(common.ChainIdToBytes(chainId))
	hashUint32(hasher, disperserId)
	hasher.Write(blobKey[:])

	return hasher.Sum(nil), nil
}
