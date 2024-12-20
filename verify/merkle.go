package verify

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// ProcessInclusionProof computes the merkle root hash based on the provided leaf and proof, returning the result.
// An error is returned if the proof param is malformed.
//
// index is the index of the leaf in the tree, starting from the bottom left of the tree at 0.
//
// If the proof length is 0, then the leaf hash is returned.
//
// NOTE: this method returning a nil error does NOT indicate that the proof is valid. Rather, it merely indicates that
// the proof was well-formed. The hash returned by this method must be compared to the claimed root hash, to
// determine if the proof is valid.
//
// This method is a reimplementation of the on-chain verification method [processInclusionProofKeccak]
// (https://github.com/Layr-Labs/eigenlayer-contracts/blob/dev/src/contracts/libraries/Merkle.sol#L49-L76)
func ProcessInclusionProof(proof []byte, leaf common.Hash, index uint64) (common.Hash, error) {
	if len(proof)%32 != 0 {
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
		index /= 2
	}

	return computedHash, nil
}
