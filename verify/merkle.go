package verify

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// ProcessInclusionProof processes the Merkle proof
func ProcessInclusionProof(proof []byte, leaf common.Hash, index uint64) (common.Hash, error) {
	if len(proof) == 0 || len(proof)%32 != 0 {
		return common.Hash{}, errors.New("proof length should be a multiple of 32 bytes or 256 bits")
	}

	computedHash := leaf
	for i := 0; i < len(proof); i += 32 {
		var proofElement common.Hash
		copy(proofElement[:], proof[i:i+32])

		var combined []byte
		if index%2 == 0 { // right
			combined = append(computedHash.Bytes(), proofElement.Bytes()...)
		} else { // left
			combined = append(proofElement.Bytes(), computedHash.Bytes()...)
		}

		computedHash = crypto.Keccak256Hash(combined)
		index = index / 2
	}

	return computedHash, nil
}
