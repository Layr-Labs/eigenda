package verify

import (
	"encoding/binary"

	binding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifier"
	"github.com/ethereum/go-ethereum/accounts/abi"
	geth_common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// HashBatchMetadata regenerates a batch data hash
// replicates: https://github.com/Layr-Labs/eigenda-utils/blob/c4cbc9ec078aeca3e4a04bd278e2fb136bf3e6de/src/libraries/EigenDAHasher.sol#L46-L54
func HashBatchMetadata(bh *binding.BatchHeader, sigHash [32]byte, blockNum uint32) (geth_common.Hash, error) {
	batchHeaderType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{
			Name: "blobHeadersRoot",
			Type: "bytes32",
		},
		{
			Name: "quorumNumbers",
			Type: "bytes",
		},
		{
			Name: "signedStakeForQuorums",
			Type: "bytes",
		},
		{
			Name: "referenceBlockNumber",
			Type: "uint32",
		},
	})

	if err != nil {
		return [32]byte{}, err
	}

	arguments := abi.Arguments{
		{
			Type: batchHeaderType,
		},
	}

	s := struct {
		BlobHeadersRoot       [32]byte
		QuorumNumbers         []byte
		SignedStakeForQuorums []byte
		ReferenceBlockNumber  uint32
	}{
		BlobHeadersRoot:       bh.BlobHeadersRoot,
		QuorumNumbers:         bh.QuorumNumbers,
		SignedStakeForQuorums: bh.SignedStakeForQuorums,
		ReferenceBlockNumber:  bh.ReferenceBlockNumber,
	}

	bytes, err := arguments.Pack(s)
	if err != nil {
		return [32]byte{}, err
	}

	headerHash := crypto.Keccak256Hash(bytes)
	return HashBatchHashedMetadata(headerHash, sigHash, blockNum)
}

// HashBatchHashedMetadata hashes the given metadata into the commitment that will be stored in the contract
// replicates: https://github.com/Layr-Labs/eigenda-utils/blob/c4cbc9ec078aeca3e4a04bd278e2fb136bf3e6de/src/libraries/EigenDAHasher.sol#L19-L25
func HashBatchHashedMetadata(batchHeaderHash [32]byte, signatoryRecordHash [32]byte, blockNumber uint32) (geth_common.Hash, error) {
	// since the solidity function uses abi.encodePacked, we need to consolidate the byte space that
	// blockNum occupies to only 4 bytes versus 28 or 256 bits when encoded to abi buffer
	a := make([]byte, 4)
	binary.BigEndian.PutUint32(a, blockNumber)

	bytes32Type, err := abi.NewType("bytes32", "bytes32", nil)
	if err != nil {
		return geth_common.BytesToHash([]byte{}), err
	}

	arguments := abi.Arguments{
		{
			Type: bytes32Type,
		},
		{
			Type: bytes32Type,
		},
	}

	bytes, err := arguments.Pack(batchHeaderHash, signatoryRecordHash)
	if err != nil {
		return [32]byte{}, err
	}

	bytes = append(bytes, a...)
	headerHash := crypto.Keccak256Hash(bytes)

	return headerHash, nil
}

// HashBlobHeader function to hash BlobHeader
func HashBlobHeader(blobHeader BlobHeader) (geth_common.Hash, error) {
	blobHeaderType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "commitment", Type: "tuple", Components: []abi.ArgumentMarshaling{
			{Name: "X", Type: "uint256"},
			{Name: "Y", Type: "uint256"},
		}},
		{Name: "dataLength", Type: "uint32"},
		{Name: "quorumBlobParams", Type: "tuple[]", Components: []abi.ArgumentMarshaling{
			{Name: "quorumNumber", Type: "uint8"},
			{Name: "adversaryThresholdPercentage", Type: "uint8"},
			{Name: "confirmationThresholdPercentage", Type: "uint8"},
			{Name: "chunkLength", Type: "uint32"},
		}},
	})
	if err != nil {
		return geth_common.Hash{}, err
	}

	// Create ABI arguments
	arguments := abi.Arguments{
		{Type: blobHeaderType},
	}

	// Pack the BlobHeader
	bytes, err := arguments.Pack(blobHeader)
	if err != nil {
		return geth_common.Hash{}, err
	}
	// Hash the packed bytes using Keccak256
	hash := crypto.Keccak256Hash(bytes)
	return hash, nil
}

// Function to hash and encode header
func HashEncodeBlobHeader(header BlobHeader) (geth_common.Hash, error) {
	// Hash the BlobHeader
	blobHash, err := HashBlobHeader(header)
	if err != nil {
		return geth_common.Hash{}, err
	}

	finalHash := crypto.Keccak256Hash(blobHash.Bytes())
	return finalHash, nil
}
