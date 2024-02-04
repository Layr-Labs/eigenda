package dispatcher

import (
	"context"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/node"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Timeout time.Duration
}

type dispatcher struct {
	*Config

	logger common.Logger
}

func NewDispatcher(cfg *Config, logger common.Logger) *dispatcher {
	return &dispatcher{
		Config: cfg,
		logger: logger,
	}
}

var _ disperser.Dispatcher = (*dispatcher)(nil)

func (c *dispatcher) DisperseBatch(ctx context.Context, state *core.IndexedOperatorState, blobs []core.EncodedBlob, header *core.BatchHeader) chan core.SignerMessage {
	update := make(chan core.SignerMessage, len(state.IndexedOperators))

	// Disperse
	c.sendAllChunks(ctx, state, blobs, header, update)

	return update
}

func (c *dispatcher) sendAllChunks(ctx context.Context, state *core.IndexedOperatorState, blobs []core.EncodedBlob, header *core.BatchHeader, update chan core.SignerMessage) {
	for id, op := range state.IndexedOperators {
		go func(op core.IndexedOperatorInfo, id core.OperatorID) {
			blobMessages := make([]*core.BlobMessage, len(blobs))
			for i, blob := range blobs {
				blobMessages[i] = blob[id]
			}
			sig, err := c.sendChunks(ctx, blobMessages, header, &op)
			if err != nil {
				update <- core.SignerMessage{
					Err:       err,
					Signature: nil,
					Operator:  id,
				}
			} else {
				update <- core.SignerMessage{
					Signature: sig,
					Operator:  id,
					Err:       nil,
				}
			}

		}(core.IndexedOperatorInfo{
			PubkeyG1: op.PubkeyG1,
			PubkeyG2: op.PubkeyG2,
			Socket:   op.Socket,
		}, id)
	}
}

func (c *dispatcher) sendChunks(ctx context.Context, blobs []*core.BlobMessage, header *core.BatchHeader, op *core.IndexedOperatorInfo) (*core.Signature, error) {
	// TODO Add secure Grpc

	conn, err := grpc.Dial(
		core.OperatorSocket(op.Socket).GetDispersalSocket(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		c.logger.Warn("Disperser cannot connect to operator dispersal socket", "dispersal_socket", core.OperatorSocket(op.Socket).GetDispersalSocket(), "err", err)
		return nil, err
	}
	defer conn.Close()

	gc := node.NewDispersalClient(conn)
	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()

	request, totalSize, err := GetStoreChunksRequest(blobs, header)
	if err != nil {
		return nil, err
	}

	opt := grpc.MaxCallSendMsgSize(1024 * 1024 * 1024)
	c.logger.Debug("sending chunks to operator", "operator", op.Socket, "size", totalSize)
	reply, err := gc.StoreChunks(ctx, request, opt)

	if err != nil {
		return nil, err
	}

	sigBytes := reply.GetSignature()
	sig := &core.Signature{G1Point: new(core.Signature).Deserialize(sigBytes)}
	return sig, nil
}

func GetStoreChunksRequest(blobMessages []*core.BlobMessage, header *core.BatchHeader) (*node.StoreChunksRequest, int64, error) {
	blobs := make([]*node.Blob, len(blobMessages))
	totalSize := int64(0)
	for i, blob := range blobMessages {
		var err error
		blobs[i], err = getBlobMessage(blob)
		if err != nil {
			return nil, 0, err
		}
		totalSize += blob.BlobHeader.EncodedSizeAllQuorums()
	}

	request := &node.StoreChunksRequest{
		BatchHeader: getBatchHeaderMessage(header),
		Blobs:       blobs,
	}

	return request, totalSize, nil
}

func getBlobMessage(blob *core.BlobMessage) (*node.Blob, error) {
	commitData, err := blob.BlobHeader.Commitment.Serialize()
	if err != nil {
		return nil, err
	}

	lengthCommitData, err := blob.BlobHeader.LengthCommitment.Serialize()
	if err != nil {
		return nil, err
	}

	lengthProofData, err := blob.BlobHeader.LengthProof.Serialize()
	if err != nil {
		return nil, err
	}

	quorumHeaders := make([]*node.BlobQuorumInfo, len(blob.BlobHeader.QuorumInfos))

	for i, header := range blob.BlobHeader.QuorumInfos {
		quorumHeaders[i] = &node.BlobQuorumInfo{
			QuorumId:           uint32(header.QuorumID),
			AdversaryThreshold: uint32(header.AdversaryThreshold),
			ChunkLength:        uint32(header.ChunkLength),
			QuorumThreshold:    uint32(header.QuorumThreshold),
			Ratelimit:          header.QuorumRate,
		}
	}

	data, err := blob.Bundles.Serialize()
	if err != nil {
		return nil, err
	}
	bundles := make([]*node.Bundle, len(blob.Bundles))
	// the ordering of quorums in bundles must be same as in quorumHeaders
	for i, quorumHeader := range quorumHeaders {
		quorum := quorumHeader.QuorumId
		bundles[i] = &node.Bundle{
			Chunks: data[quorum],
		}
	}

	return &node.Blob{
		Header: &node.BlobHeader{
			Commitment:       commitData,
			LengthCommitment: lengthCommitData,
			LengthProof:      lengthProofData,
			Length:           uint32(blob.BlobHeader.Length),
			QuorumHeaders:    quorumHeaders,
		},
		Bundles: bundles,
	}, nil
}

func getBatchHeaderMessage(header *core.BatchHeader) *node.BatchHeader {

	return &node.BatchHeader{
		BatchRoot:            header.BatchRoot[:],
		ReferenceBlockNumber: uint32(header.ReferenceBlockNumber),
	}
}
