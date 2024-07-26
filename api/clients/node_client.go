package clients

import (
	"context"
	"errors"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/node"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding"
	node_utils "github.com/Layr-Labs/eigenda/node/grpc"
	"github.com/wealdtech/go-merkletree/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RetrievedChunks struct {
	OperatorID core.OperatorID
	Chunks     []*encoding.Frame
	Err        error
}

type NodeClient interface {
	GetBlobHeader(ctx context.Context, socket string, batchHeaderHash [32]byte, blobIndex uint32) (*core.BlobHeader, *merkletree.Proof, error)
	GetChunks(ctx context.Context, opID core.OperatorID, opInfo *core.IndexedOperatorInfo, batchHeaderHash [32]byte, blobIndex uint32, quorumID core.QuorumID, chunksChan chan RetrievedChunks)
}

type client struct {
	timeout time.Duration
}

func NewNodeClient(timeout time.Duration) NodeClient {
	return client{
		timeout: timeout,
	}
}

func (c client) GetBlobHeader(
	ctx context.Context,
	socket string,
	batchHeaderHash [32]byte,
	blobIndex uint32,
) (*core.BlobHeader, *merkletree.Proof, error) {
	conn, err := grpc.Dial(
		core.OperatorSocket(socket).GetRetrievalSocket(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}
	defer conn.Close()

	n := node.NewRetrievalClient(conn)
	nodeCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	request := &node.GetBlobHeaderRequest{
		BatchHeaderHash: batchHeaderHash[:],
		BlobIndex:       blobIndex,
	}

	reply, err := n.GetBlobHeader(nodeCtx, request)
	if err != nil {
		return nil, nil, err
	}

	blobHeader, err := node_utils.GetBlobHeaderFromProto(reply.GetBlobHeader())
	if err != nil {
		return nil, nil, err
	}

	proof := &merkletree.Proof{
		Hashes: reply.GetProof().GetHashes(),
		Index:  uint64(reply.GetProof().GetIndex()),
	}

	return blobHeader, proof, nil
}

func (c client) GetChunks(
	ctx context.Context,
	opID core.OperatorID,
	opInfo *core.IndexedOperatorInfo,
	batchHeaderHash [32]byte,
	blobIndex uint32,
	quorumID core.QuorumID,
	chunksChan chan RetrievedChunks,
) {
	conn, err := grpc.Dial(
		core.OperatorSocket(opInfo.Socket).GetRetrievalSocket(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		chunksChan <- RetrievedChunks{
			OperatorID: opID,
			Err:        err,
			Chunks:     nil,
		}
		return
	}

	n := node.NewRetrievalClient(conn)
	nodeCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	request := &node.RetrieveChunksRequest{
		BatchHeaderHash: batchHeaderHash[:],
		BlobIndex:       blobIndex,
		QuorumId:        uint32(quorumID),
	}

	reply, err := n.RetrieveChunks(nodeCtx, request)
	if err != nil {
		chunksChan <- RetrievedChunks{
			OperatorID: opID,
			Err:        err,
			Chunks:     nil,
		}
		return
	}

	chunks := make([]*encoding.Frame, len(reply.GetChunks()))
	for i, data := range reply.GetChunks() {
		var chunk *encoding.Frame
		switch reply.GetEncoding() {
		case node.ChunkEncoding_GNARK:
			chunk, err = new(encoding.Frame).DeserializeGnark(data)
		case node.ChunkEncoding_GOB:
			chunk, err = new(encoding.Frame).Deserialize(data)
		case node.ChunkEncoding_UNKNOWN:
			// For backward compatibility, we fallback the UNKNOWN to GNARK
			chunk, err = new(encoding.Frame).DeserializeGnark(data)
			if err != nil {
				chunksChan <- RetrievedChunks{
					OperatorID: opID,
					Err:        errors.New("UNKNOWN chunk encoding format"),
					Chunks:     nil,
				}
			}
		}
		if err != nil {
			chunksChan <- RetrievedChunks{
				OperatorID: opID,
				Err:        err,
				Chunks:     nil,
			}
			return
		}

		chunks[i] = chunk
	}
	chunksChan <- RetrievedChunks{
		OperatorID: opID,
		Err:        nil,
		Chunks:     chunks,
	}
}
