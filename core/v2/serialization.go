package v2

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/wealdtech/go-merkletree/v2"
	"github.com/wealdtech/go-merkletree/v2/keccak256"
	"golang.org/x/crypto/sha3"
)

type abiG1Commit struct {
	X *big.Int
	Y *big.Int
}
type abiG2Commit struct {
	X [2]*big.Int
	Y [2]*big.Int
}
type abiBlobCommitments struct {
	Commitment       abiG1Commit
	LengthCommitment abiG2Commit
	LengthProof      abiG2Commit
	Length           uint32
}
type abiBlobHeader struct {
	BlobVersion         uint8
	BlobCommitments     abiBlobCommitments
	QuorumNumbers       []byte
	PaymentMetadataHash [32]byte
}

func blobHeaderArgMarshaling() []abi.ArgumentMarshaling {
	return []abi.ArgumentMarshaling{
		{
			Name: "blobVersion",
			Type: "uint8",
		},
		{
			Name: "blobCommitments",
			Type: "tuple",
			Components: []abi.ArgumentMarshaling{
				{
					Name: "commitment",
					Type: "tuple",
					Components: []abi.ArgumentMarshaling{
						{
							Name: "X",
							Type: "uint256",
						},
						{
							Name: "Y",
							Type: "uint256",
						},
					},
				},
				{
					Name: "lengthCommitment",
					Type: "tuple",
					Components: []abi.ArgumentMarshaling{
						{
							Name: "X",
							Type: "uint256[2]",
						},
						{
							Name: "Y",
							Type: "uint256[2]",
						},
					},
				},
				{
					Name: "lengthProof",
					Type: "tuple",
					Components: []abi.ArgumentMarshaling{
						{
							Name: "X",
							Type: "uint256[2]",
						},
						{
							Name: "Y",
							Type: "uint256[2]",
						},
					},
				},
				{
					Name: "length",
					Type: "uint32",
				},
			},
		},
		{
			Name: "quorumNumbers",
			Type: "bytes",
		},
		{
			Name: "paymentMetadataHash",
			Type: "bytes32",
		},
	}
}

func (b *BlobHeader) toABIStruct() (abiBlobHeader, error) {
	paymentHash, err := b.PaymentMetadata.Hash()
	if err != nil {
		return abiBlobHeader{}, err
	}
	return abiBlobHeader{
		BlobVersion: uint8(b.BlobVersion),
		BlobCommitments: abiBlobCommitments{
			Commitment: abiG1Commit{
				X: b.BlobCommitments.Commitment.X.BigInt(new(big.Int)),
				Y: b.BlobCommitments.Commitment.Y.BigInt(new(big.Int)),
			},
			LengthCommitment: abiG2Commit{
				X: [2]*big.Int{
					b.BlobCommitments.LengthCommitment.X.A0.BigInt(new(big.Int)),
					b.BlobCommitments.LengthCommitment.X.A1.BigInt(new(big.Int)),
				},
				Y: [2]*big.Int{
					b.BlobCommitments.LengthCommitment.Y.A0.BigInt(new(big.Int)),
					b.BlobCommitments.LengthCommitment.Y.A1.BigInt(new(big.Int)),
				},
			},
			LengthProof: abiG2Commit{
				X: [2]*big.Int{
					b.BlobCommitments.LengthProof.X.A0.BigInt(new(big.Int)),
					b.BlobCommitments.LengthProof.X.A1.BigInt(new(big.Int)),
				},
				Y: [2]*big.Int{
					b.BlobCommitments.LengthProof.Y.A0.BigInt(new(big.Int)),
					b.BlobCommitments.LengthProof.Y.A1.BigInt(new(big.Int)),
				},
			},
			Length: uint32(b.BlobCommitments.Length),
		},
		QuorumNumbers:       b.QuorumNumbers,
		PaymentMetadataHash: paymentHash,
	}, nil
}

func (b *BlobHeader) BlobKey() (BlobKey, error) {
	blobHeaderType, err := abi.NewType("tuple", "", blobHeaderArgMarshaling())
	if err != nil {
		return [32]byte{}, err
	}

	arguments := abi.Arguments{
		{
			Type: blobHeaderType,
		},
	}

	s, err := b.toABIStruct()
	if err != nil {
		return [32]byte{}, err
	}

	bytes, err := arguments.Pack(s)
	if err != nil {
		return [32]byte{}, err
	}

	var headerHash [32]byte
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(bytes)
	copy(headerHash[:], hasher.Sum(nil)[:32])

	return headerHash, nil
}

func (c *BlobCertificate) Hash() ([32]byte, error) {
	if c.BlobHeader == nil {
		return [32]byte{}, fmt.Errorf("blob header is nil")
	}

	blobCertType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{
			Name:       "blobHeader",
			Type:       "tuple",
			Components: blobHeaderArgMarshaling(),
		},
		{
			Name: "relayKeys",
			Type: "uint16[]",
		},
	})
	if err != nil {
		return [32]byte{}, err
	}

	arguments := abi.Arguments{
		{
			Type: blobCertType,
		},
	}

	bh, err := c.BlobHeader.toABIStruct()
	if err != nil {
		return [32]byte{}, err
	}
	s := struct {
		BlobHeader abiBlobHeader
		RelayKeys  []RelayKey
	}{
		BlobHeader: bh,
		RelayKeys:  c.RelayKeys,
	}

	bytes, err := arguments.Pack(s)
	if err != nil {
		return [32]byte{}, err
	}

	var blobCertHash [32]byte
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(bytes)
	copy(blobCertHash[:], hasher.Sum(nil)[:32])

	return blobCertHash, nil
}

func (c *BlobCertificate) Serialize() ([]byte, error) {
	return encode(c)
}

func DeserializeBlobCertificate(data []byte) (*BlobCertificate, error) {
	var c BlobCertificate
	err := decode(data, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// GetBatchHeaderHash returns the hash of the batch header
func (h BatchHeader) Hash() ([32]byte, error) {
	var headerHash [32]byte

	// The order here has to match the field ordering of ReducedBatchHeader defined in IEigenDAServiceManager.sol
	// ref: https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/interfaces/IEigenDAServiceManager.sol#L43
	batchHeaderType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{
			Name: "blobHeadersRoot",
			Type: "bytes32",
		},
		{
			Name: "referenceBlockNumber",
			Type: "uint32",
		},
	})
	if err != nil {
		return headerHash, err
	}

	arguments := abi.Arguments{
		{
			Type: batchHeaderType,
		},
	}

	s := struct {
		BlobHeadersRoot      [32]byte
		ReferenceBlockNumber uint32
	}{
		BlobHeadersRoot:      h.BatchRoot,
		ReferenceBlockNumber: uint32(h.ReferenceBlockNumber),
	}

	bytes, err := arguments.Pack(s)
	if err != nil {
		return headerHash, err
	}

	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(bytes)
	copy(headerHash[:], hasher.Sum(nil)[:32])

	return headerHash, nil
}

func (h BatchHeader) Serialize() ([]byte, error) {
	return encode(h)
}

func DeserializeBatchHeader(data []byte) (*BatchHeader, error) {
	var h BatchHeader
	err := decode(data, &h)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func BuildMerkleTree(certs []*BlobCertificate) (*merkletree.MerkleTree, error) {
	leafs := make([][]byte, len(certs))
	for i, cert := range certs {
		leaf, err := cert.Hash()
		if err != nil {
			return nil, fmt.Errorf("failed to compute blob header hash: %w", err)
		}
		leafs[i] = leaf[:]
	}

	tree, err := merkletree.NewTree(merkletree.WithData(leafs), merkletree.WithHashType(keccak256.New()))
	if err != nil {
		return nil, err
	}

	return tree, nil
}

func encode(obj any) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(obj)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decode(data []byte, obj any) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(obj)
	if err != nil {
		return err
	}
	return nil
}
