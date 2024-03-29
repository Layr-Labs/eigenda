package dispatcher

import (
	"context"
	"errors"
	"time"

	commonpb "github.com/Layr-Labs/eigenda/api/grpc/common"
	"github.com/Layr-Labs/eigenda/api/grpc/node"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigensdk-go/logging"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Timeout time.Duration
}

type dispatcher struct {
	*Config

	logger logging.Logger
}

func NewDispatcher(cfg *Config, logger logging.Logger) *dispatcher {
	return &dispatcher{
		Config: cfg,
		logger: logger,
	}
}

var _ disperser.Dispatcher = (*dispatcher)(nil)

func (c *dispatcher) DisperseBatch(ctx context.Context, state *core.IndexedOperatorState, blobs []core.EncodedBlob, batchHeader *core.BatchHeader) chan core.SignerMessage {
	update := make(chan core.SignerMessage, len(state.IndexedOperators))

	// Disperse
	c.sendAllChunks(ctx, state, blobs, batchHeader, update)

	return update
}

func (c *dispatcher) sendAllChunks(ctx context.Context, state *core.IndexedOperatorState, blobs []core.EncodedBlob, batchHeader *core.BatchHeader, update chan core.SignerMessage) {
	for id, op := range state.IndexedOperators {
		go func(op core.IndexedOperatorInfo, id core.OperatorID) {
			blobMessages := make([]*core.BlobMessage, 0)
			hasAnyBundles := false
			for _, blob := range blobs {
				if _, ok := blob.BundlesByOperator[id]; ok {
					hasAnyBundles = true
				}
				blobMessages = append(blobMessages, &core.BlobMessage{
					BlobHeader: blob.BlobHeader,
					// Bundles will be empty if the operator is not in the quorums blob is dispersed on
					Bundles: blob.BundlesByOperator[id],
				})
			}
			if !hasAnyBundles {
				// Operator is not part of any quorum, no need to send chunks
				update <- core.SignerMessage{
					Err:       errors.New("operator is not part of any quorum"),
					Signature: nil,
					Operator:  id,
				}
				return
			}

			sig, err := c.sendChunks(ctx, blobMessages, batchHeader, &op)
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

func (c *dispatcher) sendChunks(ctx context.Context, blobs []*core.BlobMessage, batchHeader *core.BatchHeader, op *core.IndexedOperatorInfo) (*core.Signature, error) {
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
	request, totalSize, err := GetStoreChunksRequest(blobs, batchHeader)
	if err != nil {
		return nil, err
	}

	opt := grpc.MaxCallSendMsgSize(60 * 1024 * 1024 * 1024)
	c.logger.Debug("sending chunks to operator", "operator", op.Socket, "size", totalSize)
	reply, err := gc.StoreChunks(ctx, request, opt)

	if err != nil {
		return nil, err
	}

	sigBytes := reply.GetSignature()
	sig := &core.Signature{G1Point: new(core.Signature).Deserialize(sigBytes)}
	return sig, nil
}

func GetStoreChunksRequest(blobMessages []*core.BlobMessage, batchHeader *core.BatchHeader) (*node.StoreChunksRequest, int64, error) {
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
		BatchHeader: getBatchHeaderMessage(batchHeader),
		Blobs:       blobs,
	}

	return request, totalSize, nil
}

func getBlobMessage(blob *core.BlobMessage) (*node.Blob, error) {
	if blob.BlobHeader == nil {
		return nil, errors.New("blob header is nil")
	}
	if blob.BlobHeader.Commitment == nil {
		return nil, errors.New("blob header commitment is nil")
	}
	commitData := &commonpb.G1Commitment{
		X: blob.BlobHeader.Commitment.X.Marshal(),
		Y: blob.BlobHeader.Commitment.Y.Marshal(),
	}
	var lengthCommitData, lengthProofData node.G2Commitment
	if blob.BlobHeader.LengthCommitment != nil {
		lengthCommitData.XA0 = blob.BlobHeader.LengthCommitment.X.A0.Marshal()
		lengthCommitData.XA1 = blob.BlobHeader.LengthCommitment.X.A1.Marshal()
		lengthCommitData.YA0 = blob.BlobHeader.LengthCommitment.Y.A0.Marshal()
		lengthCommitData.YA1 = blob.BlobHeader.LengthCommitment.Y.A1.Marshal()
	}
	if blob.BlobHeader.LengthProof != nil {
		lengthProofData.XA0 = blob.BlobHeader.LengthProof.X.A0.Marshal()
		lengthProofData.XA1 = blob.BlobHeader.LengthProof.X.A1.Marshal()
		lengthProofData.YA0 = blob.BlobHeader.LengthProof.Y.A0.Marshal()
		lengthProofData.YA1 = blob.BlobHeader.LengthProof.Y.A1.Marshal()
	}

	quorumHeaders := make([]*node.BlobQuorumInfo, len(blob.BlobHeader.QuorumInfos))

	for i, header := range blob.BlobHeader.QuorumInfos {
		quorumHeaders[i] = &node.BlobQuorumInfo{
			QuorumId:              uint32(header.QuorumID),
			AdversaryThreshold:    uint32(header.AdversaryThreshold),
			ChunkLength:           uint32(header.ChunkLength),
			ConfirmationThreshold: uint32(header.ConfirmationThreshold),
			Ratelimit:             header.QuorumRate,
		}
	}

	data, err := blob.Bundles.Serialize()
	if err != nil {
		return nil, err
	}
	bundles := make([]*node.Bundle, len(quorumHeaders))
	// the ordering of quorums in bundles must be same as in quorumHeaders
	for i, quorumHeader := range quorumHeaders {
		quorum := quorumHeader.QuorumId
		if _, ok := blob.Bundles[uint8(quorum)]; ok {
			bundles[i] = &node.Bundle{
				Chunks: data[quorum],
			}
		} else {
			bundles[i] = &node.Bundle{
				// empty bundle for quorums operators are not part of
				Chunks: make([][]byte, 0),
			}
		}
	}

	return &node.Blob{
		Header: &node.BlobHeader{
			Commitment:       commitData,
			LengthCommitment: &lengthCommitData,
			LengthProof:      &lengthProofData,
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
