package hashing

import (
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenda/common"
	"golang.org/x/crypto/sha3"
)

const DisperseBlobRequestDomain = "disperser.DisperseBlobRequest"

// Creates a hash to anchor a dispersal to the given disperser ID and chain ID
// Returns Keccak256(domain || chainId || disperserId || blobKey).
func ComputeDispersalAnchorHash(
	chainId *big.Int,
	disperserId uint32,
	blobKey [32]byte,
) ([]byte, error) {
	if chainId == nil {
		return nil, fmt.Errorf("chainId is nil")
	}

	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(DisperseBlobRequestDomain))
	hasher.Write(common.ChainIdToBytes(chainId))
	hashUint32(hasher, disperserId)
	hasher.Write(blobKey[:])

	return hasher.Sum(nil), nil
}
