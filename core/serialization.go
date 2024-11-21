package core

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"slices"

	"github.com/Layr-Labs/eigenda/api"
	binding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"

	pb "github.com/Layr-Labs/eigenda/api/grpc/node"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/wealdtech/go-merkletree/v2"
	"github.com/wealdtech/go-merkletree/v2/keccak256"
	"golang.org/x/crypto/sha3"
)

var ErrInvalidCommitment = errors.New("invalid commitment")

func ComputeSignatoryRecordHash(referenceBlockNumber uint32, nonSignerKeys []*G1Point) [32]byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, referenceBlockNumber)
	for _, nonSignerKey := range nonSignerKeys {
		hash := nonSignerKey.GetOperatorID()
		buf = append(buf, hash[:]...)
	}

	var res [32]byte
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(buf)
	copy(res[:], hasher.Sum(nil)[:32])

	return res
}

// SetBatchRoot sets the BatchRoot field of the BatchHeader to the Merkle root of the blob headers in the batch (i.e. the root of the Merkle tree whose leaves are the blob headers)
func (h *BatchHeader) SetBatchRoot(blobHeaders []*BlobHeader) (*merkletree.MerkleTree, error) {
	leafs := make([][]byte, len(blobHeaders))
	for i, header := range blobHeaders {
		leaf, err := header.GetBlobHeaderHash()
		if err != nil {
			return nil, fmt.Errorf("failed to compute blob header hash: %w", err)
		}
		leafs[i] = leaf[:]
	}

	tree, err := merkletree.NewTree(merkletree.WithData(leafs), merkletree.WithHashType(keccak256.New()))
	if err != nil {
		return nil, err
	}

	copy(h.BatchRoot[:], tree.Root())
	return tree, nil
}

func (h *BatchHeader) SetBatchRootFromBlobHeaderHashes(blobHeaderHashes [][32]byte) (*merkletree.MerkleTree, error) {
	leafs := make([][]byte, len(blobHeaderHashes))
	for i, hash := range blobHeaderHashes {
		leafs[i] = hash[:]
	}
	tree, err := merkletree.NewTree(merkletree.WithData(leafs), merkletree.WithHashType(keccak256.New()))
	if err != nil {
		return nil, err
	}

	copy(h.BatchRoot[:], tree.Root())
	return tree, nil
}

func (h *BatchHeader) Encode() ([]byte, error) {
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
		return nil, err
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
		return nil, err
	}

	return bytes, nil
}

// GetBatchHeaderHash returns the hash of the reduced BatchHeader that is used to sign the Batch
// ref: https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/libraries/EigenDAHasher.sol#L65
func (h BatchHeader) GetBatchHeaderHash() ([32]byte, error) {
	headerByte, err := h.Encode()
	if err != nil {
		return [32]byte{}, err
	}

	var headerHash [32]byte
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(headerByte)
	copy(headerHash[:], hasher.Sum(nil)[:32])

	return headerHash, nil
}

// HashBatchHeader returns the hash of the BatchHeader that is used to emit the BatchConfirmed event
// ref: https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/libraries/EigenDAHasher.sol#L57
func HashBatchHeader(batchHeader binding.BatchHeader) ([32]byte, error) {
	// The order here has to match the field ordering of BatchHeader defined in IEigenDAServiceManager.sol
	batchHeaderType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{
			Name: "batchRoot",
			Type: "bytes32",
		},
		{
			Name: "quorumNumbers",
			Type: "bytes",
		},
		{
			Name: "confirmationThresholdPercentages",
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
		BatchRoot                        [32]byte
		QuorumNumbers                    []byte
		ConfirmationThresholdPercentages []byte
		ReferenceBlockNumber             uint32
	}{
		BatchRoot:                        batchHeader.BlobHeadersRoot,
		QuorumNumbers:                    batchHeader.QuorumNumbers,
		ConfirmationThresholdPercentages: batchHeader.SignedStakeForQuorums,
		ReferenceBlockNumber:             uint32(batchHeader.ReferenceBlockNumber),
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

// GetBlobHeaderHash returns the hash of the BlobHeader that is used to sign the Blob
func (h BlobHeader) GetBlobHeaderHash() ([32]byte, error) {
	headerByte, err := h.Encode()
	if err != nil {
		return [32]byte{}, err
	}

	var headerHash [32]byte
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(headerByte)
	copy(headerHash[:], hasher.Sum(nil)[:32])

	return headerHash, nil
}

func (h *BlobHeader) GetQuorumBlobParamsHash() ([32]byte, error) {
	quorumBlobParamsType, err := abi.NewType("tuple[]", "", []abi.ArgumentMarshaling{
		{
			Name: "quorumNumber",
			Type: "uint8",
		},
		{
			Name: "adversaryThresholdPercentage",
			Type: "uint8",
		},
		{
			Name: "quorumThresholdPercentage",
			Type: "uint8",
		},
		{
			Name: "chunkLength",
			Type: "uint32",
		},
	})

	if err != nil {
		return [32]byte{}, err
	}

	arguments := abi.Arguments{
		{
			Type: quorumBlobParamsType,
		},
	}

	type quorumBlobParams struct {
		QuorumNumber                 uint8
		AdversaryThresholdPercentage uint8
		QuorumThresholdPercentage    uint8
		ChunkLength                  uint32
	}

	qbp := make([]quorumBlobParams, len(h.QuorumInfos))
	for i, q := range h.QuorumInfos {
		qbp[i] = quorumBlobParams{
			QuorumNumber:                 q.QuorumID,
			AdversaryThresholdPercentage: q.AdversaryThreshold,
			QuorumThresholdPercentage:    q.ConfirmationThreshold,
			ChunkLength:                  uint32(q.ChunkLength),
		}
	}

	bytes, err := arguments.Pack(qbp)
	if err != nil {
		return [32]byte{}, err
	}

	var res [32]byte
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(bytes)
	copy(res[:], hasher.Sum(nil)[:32])

	return res, nil
}

func (h *BlobHeader) Encode() ([]byte, error) {
	if h.Commitment == nil {
		return nil, ErrInvalidCommitment
	}

	// The order here has to match the field ordering of BlobHeader defined in IEigenDAServiceManager.sol
	blobHeaderType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
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
			Name: "dataLength",
			Type: "uint32",
		},
		{
			Name: "quorumBlobParams",
			Type: "tuple[]",
			Components: []abi.ArgumentMarshaling{
				{
					Name: "quorumNumber",
					Type: "uint8",
				},
				{
					Name: "adversaryThresholdPercentage",
					Type: "uint8",
				},
				{
					Name: "quorumThresholdPercentage",
					Type: "uint8",
				},
				{
					Name: "chunkLength",
					Type: "uint32",
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	arguments := abi.Arguments{
		{
			Type: blobHeaderType,
		},
	}

	type quorumBlobParams struct {
		QuorumNumber                 uint8
		AdversaryThresholdPercentage uint8
		QuorumThresholdPercentage    uint8
		ChunkLength                  uint32
	}

	type commitment struct {
		X *big.Int
		Y *big.Int
	}

	qbp := make([]quorumBlobParams, len(h.QuorumInfos))
	for i, q := range h.QuorumInfos {
		qbp[i] = quorumBlobParams{
			QuorumNumber:                 q.QuorumID,
			AdversaryThresholdPercentage: q.AdversaryThreshold,
			QuorumThresholdPercentage:    q.ConfirmationThreshold,
			ChunkLength:                  uint32(q.ChunkLength),
		}
	}
	slices.SortStableFunc[[]quorumBlobParams](qbp, func(a, b quorumBlobParams) int {
		return int(a.QuorumNumber) - int(b.QuorumNumber)
	})

	s := struct {
		Commitment       commitment
		DataLength       uint32
		QuorumBlobParams []quorumBlobParams
	}{
		Commitment: commitment{
			X: h.Commitment.X.BigInt(new(big.Int)),
			Y: h.Commitment.Y.BigInt(new(big.Int)),
		},
		DataLength:       uint32(h.Length),
		QuorumBlobParams: qbp,
	}

	bytes, err := arguments.Pack(s)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (h *BatchHeader) Serialize() ([]byte, error) {
	return encode(h)
}

func (h *BatchHeader) Deserialize(data []byte) (*BatchHeader, error) {
	err := decode(data, h)
	return h, err
}

func (h *BlobHeader) Serialize() ([]byte, error) {
	return encode(h)
}

func (h *BlobHeader) Deserialize(data []byte) (*BlobHeader, error) {
	err := decode(data, h)

	if !(*bn254.G1Affine)(h.BlobCommitments.Commitment).IsInSubGroup() {
		return nil, fmt.Errorf("in BlobHeader Commitment is not in the subgroup")
	}

	if !(*bn254.G2Affine)(h.BlobCommitments.LengthCommitment).IsInSubGroup() {
		return nil, fmt.Errorf("in BlobHeader LengthCommitment is not in the subgroup")
	}

	if !(*bn254.G2Affine)(h.BlobCommitments.LengthProof).IsInSubGroup() {
		return nil, fmt.Errorf("in BlobHeader LengthProof is not in the subgroup")
	}

	return h, err
}

// GetBatchHeader constructs a core.BatchHeader from a proto of pb.StoreChunksRequest.
// Note the StoreChunksRequest is validated as soon as it enters the node gRPC
// interface, see grpc.Server.validateStoreChunkRequest.
func BatchHeaderFromProtobuf(in *pb.BatchHeader) (*BatchHeader, error) {
	if in == nil || len(in.GetBatchRoot()) == 0 {
		return nil, fmt.Errorf("batch header is nil or empty")
	}
	var batchRoot [32]byte
	copy(batchRoot[:], in.GetBatchRoot())
	return &BatchHeader{
		ReferenceBlockNumber: uint(in.GetReferenceBlockNumber()),
		BatchRoot:            batchRoot,
	}, nil
}

// BlobHeaderFromProtobuf constructs a core.BlobHeader from a proto of pb.BlobHeader.
func BlobHeaderFromProtobuf(h *pb.BlobHeader) (*BlobHeader, error) {
	if h == nil {
		return nil, fmt.Errorf("GetBlobHeaderFromProto: blob header is nil")

	}

	commitX := new(fp.Element).SetBytes(h.GetCommitment().GetX())
	commitY := new(fp.Element).SetBytes(h.GetCommitment().GetY())
	commitment := &encoding.G1Commitment{
		X: *commitX,
		Y: *commitY,
	}

	if !(*bn254.G1Affine)(commitment).IsInSubGroup() {
		return nil, errors.New("commitment is not in the subgroup")
	}

	var lengthCommitment, lengthProof encoding.G2Commitment
	if h.GetLengthCommitment() != nil {
		lengthCommitment.X.A0 = *new(fp.Element).SetBytes(h.GetLengthCommitment().GetXA0())
		lengthCommitment.X.A1 = *new(fp.Element).SetBytes(h.GetLengthCommitment().GetXA1())
		lengthCommitment.Y.A0 = *new(fp.Element).SetBytes(h.GetLengthCommitment().GetYA0())
		lengthCommitment.Y.A1 = *new(fp.Element).SetBytes(h.GetLengthCommitment().GetYA1())
	}

	if !(*bn254.G2Affine)(&lengthCommitment).IsInSubGroup() {
		return nil, errors.New("lengthCommitment is not in the subgroup")
	}

	if h.GetLengthProof() != nil {
		lengthProof.X.A0 = *new(fp.Element).SetBytes(h.GetLengthProof().GetXA0())
		lengthProof.X.A1 = *new(fp.Element).SetBytes(h.GetLengthProof().GetXA1())
		lengthProof.Y.A0 = *new(fp.Element).SetBytes(h.GetLengthProof().GetYA0())
		lengthProof.Y.A1 = *new(fp.Element).SetBytes(h.GetLengthProof().GetYA1())
	}

	if !(*bn254.G2Affine)(&lengthProof).IsInSubGroup() {
		return nil, errors.New("lengthProof is not in the subgroup")
	}

	quorumHeaders := make([]*BlobQuorumInfo, len(h.GetQuorumHeaders()))
	for i, header := range h.GetQuorumHeaders() {
		if header.GetQuorumId() > MaxQuorumID {
			return nil, api.NewErrorInvalidArg(fmt.Sprintf("quorum ID must be in range [0, %d], but found %d", MaxQuorumID, header.GetQuorumId()))
		}
		if err := ValidateSecurityParam(header.GetConfirmationThreshold(), header.GetAdversaryThreshold()); err != nil {
			return nil, err
		}

		quorumHeaders[i] = &BlobQuorumInfo{
			SecurityParam: SecurityParam{
				QuorumID:              QuorumID(header.GetQuorumId()),
				AdversaryThreshold:    uint8(header.GetAdversaryThreshold()),
				ConfirmationThreshold: uint8(header.GetConfirmationThreshold()),
				QuorumRate:            header.GetRatelimit(),
			},
			ChunkLength: uint(header.GetChunkLength()),
		}
	}

	return &BlobHeader{
		BlobCommitments: encoding.BlobCommitments{
			Commitment:       commitment,
			LengthCommitment: &lengthCommitment,
			LengthProof:      &lengthProof,
			Length:           uint(h.GetLength()),
		},
		QuorumInfos: quorumHeaders,
		AccountID:   h.AccountId,
	}, nil
}

func SerializeMerkleProof(proof *merkletree.Proof) []byte {
	proofBytes := make([]byte, 0)
	for _, hash := range proof.Hashes {
		proofBytes = append(proofBytes, hash[:]...)
	}
	return proofBytes
}

func DeserializeMerkleProof(data []byte, index uint64) (*merkletree.Proof, error) {
	proof := &merkletree.Proof{
		Index: index,
	}
	if len(data)%32 != 0 {
		return nil, fmt.Errorf("invalid proof length")
	}
	for i := 0; i < len(data); i += 32 {
		var hash [32]byte
		copy(hash[:], data[i:i+32])
		proof.Hashes = append(proof.Hashes, hash[:])
	}
	return proof, nil
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

func (s OperatorSocket) GetDispersalSocket() string {
	ip, port1, _, err := extractIPAndPorts(string(s))
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s:%s", ip, port1)
}

func (s OperatorSocket) GetRetrievalSocket() string {
	ip, _, port2, err := extractIPAndPorts(string(s))
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s:%s", ip, port2)
}

func extractIPAndPorts(s string) (string, string, string, error) {
	regex := regexp.MustCompile(`^([^:]+):([^;]+);([^;]+)$`)
	matches := regex.FindStringSubmatch(s)

	if len(matches) != 4 {
		return "", "", "", errors.New("input string does not match expected format")
	}

	ip := matches[1]
	port1 := matches[2]
	port2 := matches[3]

	return ip, port1, port2, nil
}
