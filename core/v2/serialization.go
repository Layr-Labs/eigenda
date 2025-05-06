package v2

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/big"
	"slices"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding"
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
	DataLength       uint32
}

// ComputeBlobKey accepts as parameters the elements which contribute to the hash of a BlobHeader. It computes the
// hash and returns the result, which represents a BlobKey.
//
// This function exists so that the BlobKey can be computed without first constructing a BlobHeader object. Since
// the BlobHeader contains the full payment metadata, and payment metadata isn't stored on chain, it isn't always
// possible to reconstruct from the data available.
//
// The hashing structure here must ALWAYS match the hashing structure that we perform onchain:
// https://github.com/Layr-Labs/eigenda/blob/a6dd724acdf732af483fd2d9a86325febe7ebdcd/contracts/src/libraries/EigenDAHasher.sol#L119
func ComputeBlobKey(
	blobVersion BlobVersion,
	blobCommitments encoding.BlobCommitments,
	quorumNumbers []core.QuorumID,
	paymentMetadataHash [32]byte,
) ([32]byte, error) {
	versionType, err := abi.NewType("uint16", "", nil)
	if err != nil {
		return [32]byte{}, err
	}
	quorumNumbersType, err := abi.NewType("bytes", "", nil)
	if err != nil {
		return [32]byte{}, err
	}
	commitmentType, err := abi.NewType(
		"tuple", "", []abi.ArgumentMarshaling{
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
				Name: "dataLength",
				Type: "uint32",
			},
		})
	if err != nil {
		return [32]byte{}, err
	}
	arguments := abi.Arguments{
		{
			Type: versionType,
		},
		{
			Type: quorumNumbersType,
		},
		{
			Type: commitmentType,
		},
	}
	// Sort the quorum numbers to ensure the hash is consistent
	sortedQuorums := make([]core.QuorumID, len(quorumNumbers))
	copy(sortedQuorums, quorumNumbers)
	slices.Sort(sortedQuorums)
	packedBytes, err := arguments.Pack(
		blobVersion,
		sortedQuorums,
		abiBlobCommitments{
			Commitment: abiG1Commit{
				X: blobCommitments.Commitment.X.BigInt(new(big.Int)),
				Y: blobCommitments.Commitment.Y.BigInt(new(big.Int)),
			},
			// Most cryptography library serializes a G2 point by having
			// A0 followed by A1 for both X, Y field of G2. However, ethereum
			// precompile assumes an ordering of A1, A0. We choose
			// to conform with Ethereum order when serializing a blobHeaderV2
			// for instance, gnark, https://github.com/Consensys/gnark-crypto/blob/de0d77f2b4d520350bc54c612828b19ce2146eee/ecc/bn254/marshal.go#L1078
			// Ethereum, https://eips.ethereum.org/EIPS/eip-197#definition-of-the-groups
			LengthCommitment: abiG2Commit{
				X: [2]*big.Int{
					blobCommitments.LengthCommitment.X.A1.BigInt(new(big.Int)),
					blobCommitments.LengthCommitment.X.A0.BigInt(new(big.Int)),
				},
				Y: [2]*big.Int{
					blobCommitments.LengthCommitment.Y.A1.BigInt(new(big.Int)),
					blobCommitments.LengthCommitment.Y.A0.BigInt(new(big.Int)),
				},
			},
			// Same as above
			LengthProof: abiG2Commit{
				X: [2]*big.Int{
					blobCommitments.LengthProof.X.A1.BigInt(new(big.Int)),
					blobCommitments.LengthProof.X.A0.BigInt(new(big.Int)),
				},
				Y: [2]*big.Int{
					blobCommitments.LengthProof.Y.A1.BigInt(new(big.Int)),
					blobCommitments.LengthProof.Y.A0.BigInt(new(big.Int)),
				},
			},
			DataLength: uint32(blobCommitments.Length),
		},
	)
	if err != nil {
		return [32]byte{}, err
	}

	var headerHash [32]byte
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(packedBytes)
	copy(headerHash[:], hasher.Sum(nil)[:32])

	blobKeyType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{
			Name: "blobHeaderHash",
			Type: "bytes32",
		},
		{
			Name: "paymentMetadataHash",
			Type: "bytes32",
		},
	})
	if err != nil {
		return [32]byte{}, err
	}

	arguments = abi.Arguments{
		{
			Type: blobKeyType,
		},
	}

	s2 := struct {
		BlobHeaderHash      [32]byte
		PaymentMetadataHash [32]byte
	}{
		BlobHeaderHash:      headerHash,
		PaymentMetadataHash: paymentMetadataHash,
	}

	packedBytes, err = arguments.Pack(s2)
	if err != nil {
		return [32]byte{}, err
	}

	var blobKey [32]byte
	hasher = sha3.NewLegacyKeccak256()
	hasher.Write(packedBytes)
	copy(blobKey[:], hasher.Sum(nil)[:32])

	return blobKey, nil
}

// BlobKey computes the BlobKey of the BlobHeader.
//
// A BlobKey simply the hash of the BlobHeader
func (b *BlobHeader) BlobKey() (BlobKey, error) {
	blobHeaderWithoutPayment, err := b.GetBlobHeaderWithoutPayment()
	if err != nil {
		return BlobKey{}, fmt.Errorf("get blob header without payment: %w", err)
	}

	return blobHeaderWithoutPayment.BlobKey()
}

func (b *BlobHeader) GetBlobHeaderWithoutPayment() (*BlobHeaderWithoutPayment, error) {
	paymentMetadataHash, err := b.PaymentMetadata.Hash()
	if err != nil {
		return nil, fmt.Errorf("hash payment metadata: %w", err)
	}

	return &BlobHeaderWithoutPayment{
		BlobVersion:         b.BlobVersion,
		BlobCommitments:     b.BlobCommitments,
		QuorumNumbers:       b.QuorumNumbers,
		PaymentMetadataHash: paymentMetadataHash,
	}, nil
}

func (b *BlobHeaderWithoutPayment) BlobKey() (BlobKey, error) {
	return ComputeBlobKey(
		b.BlobVersion,
		b.BlobCommitments,
		b.QuorumNumbers,
		b.PaymentMetadataHash,
	)
}

func (c *BlobCertificate) Hash() ([32]byte, error) {
	if c.BlobHeader == nil {
		return [32]byte{}, fmt.Errorf("blob header is nil")
	}

	blobKeyType, err := abi.NewType("bytes32", "", nil)
	if err != nil {
		return [32]byte{}, err
	}

	signatureType, err := abi.NewType("bytes", "", nil)
	if err != nil {
		return [32]byte{}, err
	}

	relayKeysType, err := abi.NewType("uint32[]", "", nil)
	if err != nil {
		return [32]byte{}, err
	}

	arguments := abi.Arguments{
		{
			Type: blobKeyType,
		},
		{
			Type: signatureType,
		},
		{
			Type: relayKeysType,
		},
	}

	blobKey, err := c.BlobHeader.BlobKey()
	if err != nil {
		return [32]byte{}, err
	}

	bytes, err := arguments.Pack(blobKey, c.Signature, c.RelayKeys)
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
