package grpc

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Layr-Labs/eigenda/api"
	pb "github.com/Layr-Labs/eigenda/api/grpc/node"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/node"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/wealdtech/go-merkletree"
	"github.com/wealdtech/go-merkletree/keccak256"
	"google.golang.org/protobuf/proto"
)

// GetBatchHeader constructs a core.BatchHeader from a proto of pb.StoreChunksRequest.
// Note the StoreChunksRequest is validated as soon as it enters the node gRPC
// interface, see grpc.Server.validateStoreChunkRequest.
func GetBatchHeader(in *pb.StoreChunksRequest) (*core.BatchHeader, error) {
	var batchRoot [32]byte
	copy(batchRoot[:], in.GetBatchHeader().GetBatchRoot())
	batchHeader := core.BatchHeader{
		ReferenceBlockNumber: uint(in.GetBatchHeader().GetReferenceBlockNumber()),
		BatchRoot:            batchRoot,
	}
	return &batchHeader, nil
}

// GetBlobMessages constructs a core.BlobMessage array from a proto of pb.StoreChunksRequest.
// Note the StoreChunksRequest is validated as soon as it enters the node gRPC
// interface, see grpc.Server.validateStoreChunkRequest.
func GetBlobMessages(in *pb.StoreChunksRequest) ([]*core.BlobMessage, error) {
	blobs := make([]*core.BlobMessage, len(in.GetBlobs()))
	for i, blob := range in.GetBlobs() {
		blobHeader, err := GetBlobHeaderFromProto(blob.GetHeader())

		if err != nil {
			return nil, err
		}
		if len(blob.GetBundles()) != len(blob.GetHeader().GetQuorumHeaders()) {
			return nil, fmt.Errorf("number of quorum headers (%d) does not match number of bundles in blob message (%d)", len(blob.GetHeader().GetQuorumHeaders()), len(blob.GetBundles()))
		}

		bundles := make(map[core.QuorumID]core.Bundle, len(blob.GetBundles()))
		for j, chunks := range blob.GetBundles() {
			quorumID := blob.GetHeader().GetQuorumHeaders()[j].QuorumId
			bundles[uint8(quorumID)] = make([]*encoding.Frame, len(chunks.GetChunks()))
			for k, data := range chunks.GetChunks() {
				chunk, err := new(encoding.Frame).Deserialize(data)
				if err != nil {
					return nil, err
				}
				bundles[uint8(quorumID)][k] = chunk
			}
		}

		blobs[i] = &core.BlobMessage{
			BlobHeader: blobHeader,
			Bundles:    bundles,
		}
	}
	return blobs, nil
}

// GetBlobHeaderFromProto constructs a core.BlobHeader from a proto of pb.BlobHeader.
func GetBlobHeaderFromProto(h *pb.BlobHeader) (*core.BlobHeader, error) {
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

	quorumHeaders := make([]*core.BlobQuorumInfo, len(h.GetQuorumHeaders()))
	for i, header := range h.GetQuorumHeaders() {
		if header.GetQuorumId() > core.MaxQuorumID {
			return nil, api.NewInvalidArgError(fmt.Sprintf("quorum ID must be in range [0, %d], but found %d", core.MaxQuorumID, header.GetQuorumId()))
		}
		if err := core.ValidateSecurityParam(header.GetConfirmationThreshold(), header.GetAdversaryThreshold()); err != nil {
			return nil, err
		}

		quorumHeaders[i] = &core.BlobQuorumInfo{
			SecurityParam: core.SecurityParam{
				QuorumID:              core.QuorumID(header.GetQuorumId()),
				AdversaryThreshold:    uint8(header.GetAdversaryThreshold()),
				ConfirmationThreshold: uint8(header.GetConfirmationThreshold()),
				QuorumRate:            header.GetRatelimit(),
			},
			ChunkLength: uint(header.GetChunkLength()),
		}
	}

	return &core.BlobHeader{
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

// rebuildMerkleTree rebuilds the merkle tree from the blob headers and batch header.
func (s *Server) rebuildMerkleTree(batchHeaderHash [32]byte) (*merkletree.MerkleTree, error) {
	batchHeaderBytes, err := s.node.Store.GetBatchHeader(context.Background(), batchHeaderHash)
	if err != nil {
		return nil, errors.New("failed to get the batch header from Store")
	}

	batchHeader, err := new(core.BatchHeader).Deserialize(batchHeaderBytes)
	if err != nil {
		return nil, err
	}

	blobIndex := 0
	leafs := make([][]byte, 0)
	for {
		blobHeaderBytes, err := s.node.Store.GetBlobHeader(context.Background(), batchHeaderHash, blobIndex)
		if err != nil {
			if errors.Is(err, node.ErrKeyNotFound) {
				break
			}
			return nil, err
		}

		var protoBlobHeader pb.BlobHeader
		err = proto.Unmarshal(blobHeaderBytes, &protoBlobHeader)
		if err != nil {
			return nil, err
		}

		blobHeader, err := GetBlobHeaderFromProto(&protoBlobHeader)
		if err != nil {
			return nil, err
		}

		blobHeaderHash, err := blobHeader.GetBlobHeaderHash()
		if err != nil {
			return nil, err
		}
		leafs = append(leafs, blobHeaderHash[:])
		blobIndex++
	}

	if len(leafs) == 0 {
		return nil, errors.New("no blob header found")
	}

	tree, err := merkletree.NewTree(merkletree.WithData(leafs), merkletree.WithHashType(keccak256.New()))
	if err != nil {
		return nil, err
	}

	if !reflect.DeepEqual(tree.Root(), batchHeader.BatchRoot[:]) {
		return nil, errors.New("invalid batch header")
	}

	return tree, nil
}

// // Constructs a core.SecurityParam from a proto of pb.SecurityParams.
// func GetSecurityParam(p []*pb.SecurityParam) []*core.SecurityParam {
// 	res := make([]*core.SecurityParam, len(p))
// 	for i := range p {
// 		res[i] = &core.SecurityParam{
// 			QuorumID:           core.QuorumID(p[i].GetQuorumId()),
// 			AdversaryThreshold: uint8(p[i].GetAdversaryThreshold()),
// 		}
// 	}
// 	return res
// }

// // Constructs a core.QuorumParam array from a proto of pb.BatchHeader.
// func GetQuorumParams(p *pb.BatchHeader) []core.QuorumParam {
// 	quorum := make([]core.QuorumParam, 0)
// 	for _, param := range p.GetQuorumParams() {
// 		qp := core.QuorumParam{
// 			QuorumID:        core.QuorumID(param.GetQuorumId()),
// 			ConfirmationThreshold: uint8(param.GetQuorumThreshold()),
// 		}
// 		quorum = append(quorum, qp)
// 	}
// 	return quorum
// }
