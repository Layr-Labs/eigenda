package dispatcher

import (
	"context"
	"errors"
	"fmt"
	"time"

	commonpb "github.com/Layr-Labs/eigenda/api/grpc/common"
	"github.com/Layr-Labs/eigenda/api/grpc/node"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/batcher"
	"github.com/Layr-Labs/eigensdk-go/logging"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

type Config struct {
	Timeout                   time.Duration
	EnableGnarkBundleEncoding bool
}

type dispatcher struct {
	*Config

	logger  logging.Logger
	metrics *batcher.DispatcherMetrics
}

func NewDispatcher(cfg *Config, logger logging.Logger, metrics *batcher.DispatcherMetrics) *dispatcher {
	return &dispatcher{
		Config:  cfg,
		logger:  logger.With("component", "Dispatcher"),
		metrics: metrics,
	}
}

var _ disperser.Dispatcher = (*dispatcher)(nil)

func (c *dispatcher) DisperseBatch(ctx context.Context, state *core.IndexedOperatorState, blobs []core.EncodedBlob, batchHeader *core.BatchHeader) chan core.SigningMessage {
	update := make(chan core.SigningMessage, len(state.IndexedOperators))

	// Disperse
	c.sendAllChunks(ctx, state, blobs, batchHeader, update)

	return update
}

func (c *dispatcher) sendAllChunks(ctx context.Context, state *core.IndexedOperatorState, blobs []core.EncodedBlob, batchHeader *core.BatchHeader, update chan core.SigningMessage) {
	for id, op := range state.IndexedOperators {
		go func(op core.IndexedOperatorInfo, id core.OperatorID) {
			blobMessages := make([]*core.EncodedBlobMessage, 0)
			hasAnyBundles := false
			batchHeaderHash, err := batchHeader.GetBatchHeaderHash()
			if err != nil {
				return
			}
			for _, blob := range blobs {
				if _, ok := blob.EncodedBundlesByOperator[id]; ok {
					hasAnyBundles = true
				}
				blobMessages = append(blobMessages, &core.EncodedBlobMessage{
					BlobHeader: blob.BlobHeader,
					// Bundles will be empty if the operator is not in the quorums blob is dispersed on
					EncodedBundles: blob.EncodedBundlesByOperator[id],
				})
			}
			if !hasAnyBundles {
				// Operator is not part of any quorum, no need to send chunks
				update <- core.SigningMessage{
					Err:                  errors.New("operator is not part of any quorum"),
					Signature:            nil,
					Operator:             id,
					BatchHeaderHash:      batchHeaderHash,
					AttestationLatencyMs: -1,
				}
				return
			}

			requestedAt := time.Now()
			sig, err := c.sendChunks(ctx, blobMessages, batchHeader, &op)
			latencyMs := float64(time.Since(requestedAt).Milliseconds())
			if err != nil {
				update <- core.SigningMessage{
					Err:                  err,
					Signature:            nil,
					Operator:             id,
					BatchHeaderHash:      batchHeaderHash,
					AttestationLatencyMs: latencyMs,
				}
				c.metrics.ObserveLatency(id.Hex(), false, latencyMs)
			} else {
				update <- core.SigningMessage{
					Signature:            sig,
					Operator:             id,
					BatchHeaderHash:      batchHeaderHash,
					AttestationLatencyMs: latencyMs,
					Err:                  nil,
				}
				c.metrics.ObserveLatency(id.Hex(), true, latencyMs)
			}

		}(core.IndexedOperatorInfo{
			PubkeyG1: op.PubkeyG1,
			PubkeyG2: op.PubkeyG2,
			Socket:   op.Socket,
		}, id)
	}
}

func (c *dispatcher) sendChunks(ctx context.Context, blobs []*core.EncodedBlobMessage, batchHeader *core.BatchHeader, op *core.IndexedOperatorInfo) (*core.Signature, error) {
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
	start := time.Now()
	request, totalSize, err := GetStoreChunksRequest(blobs, batchHeader, c.EnableGnarkBundleEncoding)
	if err != nil {
		return nil, err
	}
	c.logger.Debug("sending chunks to operator", "operator", op.Socket, "num blobs", len(blobs), "size", totalSize, "request message size", proto.Size(request), "request serialization time", time.Since(start), "use Gnark chunk encoding", c.EnableGnarkBundleEncoding)
	opt := grpc.MaxCallSendMsgSize(60 * 1024 * 1024 * 1024)
	reply, err := gc.StoreChunks(ctx, request, opt)

	if err != nil {
		return nil, err
	}

	sigBytes := reply.GetSignature()
	point, err := new(core.Signature).Deserialize(sigBytes)
	if err != nil {
		return nil, err
	}
	sig := &core.Signature{G1Point: point}
	return sig, nil
}

// SendBlobsToOperator sends blobs to an operator via the node's StoreBlobs endpoint
// It returns the signatures of the blobs sent to the operator in the same order as the blobs
// with nil values for blobs that were not attested by the operator
func (c *dispatcher) SendBlobsToOperator(ctx context.Context, blobs []*core.EncodedBlobMessage, batchHeader *core.BatchHeader, op *core.IndexedOperatorInfo) ([]*core.Signature, error) {
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
	start := time.Now()
	request, totalSize, err := GetStoreBlobsRequest(blobs, batchHeader, c.EnableGnarkBundleEncoding)
	if err != nil {
		return nil, err
	}
	c.logger.Debug("sending chunks to operator", "operator", op.Socket, "num blobs", len(blobs), "size", totalSize, "request message size", proto.Size(request), "request serialization time", time.Since(start), "use Gnark chunk encoding", c.EnableGnarkBundleEncoding)
	opt := grpc.MaxCallSendMsgSize(60 * 1024 * 1024 * 1024)
	reply, err := gc.StoreBlobs(ctx, request, opt)

	if err != nil {
		return nil, err
	}

	signaturesInBytes := reply.GetSignatures()
	signatures := make([]*core.Signature, 0, len(signaturesInBytes))
	for _, sigBytes := range signaturesInBytes {
		sig := sigBytes.GetValue()
		if sig != nil {
			point, err := new(core.Signature).Deserialize(sig)
			if err != nil {
				return nil, err
			}
			signatures = append(signatures, &core.Signature{G1Point: point})
		} else {
			signatures = append(signatures, nil)
		}
	}
	return signatures, nil
}

func (c *dispatcher) AttestBatch(ctx context.Context, state *core.IndexedOperatorState, blobHeaderHashes [][32]byte, batchHeader *core.BatchHeader) (chan core.SigningMessage, error) {
	batchHeaderHash, err := batchHeader.GetBatchHeaderHash()
	if err != nil {
		return nil, err
	}
	responseChan := make(chan core.SigningMessage, len(state.IndexedOperators))

	for id, op := range state.IndexedOperators {
		go func(op core.IndexedOperatorInfo, id core.OperatorID) {
			conn, err := grpc.Dial(
				core.OperatorSocket(op.Socket).GetDispersalSocket(),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
			if err != nil {
				c.logger.Error("disperser cannot connect to operator dispersal socket", "socket", core.OperatorSocket(op.Socket).GetDispersalSocket(), "err", err)
				return
			}
			defer conn.Close()

			nodeClient := node.NewDispersalClient(conn)

			requestedAt := time.Now()
			sig, err := c.SendAttestBatchRequest(ctx, nodeClient, blobHeaderHashes, batchHeader, &op)
			latencyMs := float64(time.Since(requestedAt).Milliseconds())
			if err != nil {
				responseChan <- core.SigningMessage{
					Err:                  err,
					Signature:            nil,
					Operator:             id,
					BatchHeaderHash:      batchHeaderHash,
					AttestationLatencyMs: latencyMs,
				}
				c.metrics.ObserveLatency(id.Hex(), false, latencyMs)
			} else {
				responseChan <- core.SigningMessage{
					Signature:            sig,
					Operator:             id,
					BatchHeaderHash:      batchHeaderHash,
					AttestationLatencyMs: latencyMs,
					Err:                  nil,
				}
				c.metrics.ObserveLatency(id.Hex(), true, latencyMs)
			}
		}(core.IndexedOperatorInfo{
			PubkeyG1: op.PubkeyG1,
			PubkeyG2: op.PubkeyG2,
			Socket:   op.Socket,
		}, id)
	}

	return responseChan, nil
}

func (c *dispatcher) SendAttestBatchRequest(ctx context.Context, nodeDispersalClient node.DispersalClient, blobHeaderHashes [][32]byte, batchHeader *core.BatchHeader, op *core.IndexedOperatorInfo) (*core.Signature, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()
	// start := time.Now()
	hashes := make([][]byte, len(blobHeaderHashes))
	for i, hash := range blobHeaderHashes {
		hashes[i] = hash[:]
	}

	request := &node.AttestBatchRequest{
		BatchHeader:      getBatchHeaderMessage(batchHeader),
		BlobHeaderHashes: hashes,
	}

	c.logger.Debug("sending AttestBatch request to operator", "operator", op.Socket, "numBlobs", len(blobHeaderHashes), "requestMessageSize", proto.Size(request), "referenceBlockNumber", batchHeader.ReferenceBlockNumber)
	opt := grpc.MaxCallSendMsgSize(60 * 1024 * 1024 * 1024)
	reply, err := nodeDispersalClient.AttestBatch(ctx, request, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to send AttestBatch request to operator %s: %w", core.OperatorSocket(op.Socket).GetDispersalSocket(), err)
	}

	sigBytes := reply.GetSignature()
	point, err := new(core.Signature).Deserialize(sigBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize signature: %w", err)
	}
	return &core.Signature{G1Point: point}, nil
}

func GetStoreChunksRequest(blobMessages []*core.EncodedBlobMessage, batchHeader *core.BatchHeader, useGnarkBundleEncoding bool) (*node.StoreChunksRequest, int64, error) {
	blobs := make([]*node.Blob, len(blobMessages))
	totalSize := int64(0)
	for i, blob := range blobMessages {
		var err error
		blobs[i], err = getBlobMessage(blob, useGnarkBundleEncoding)
		if err != nil {
			return nil, 0, err
		}
		totalSize += getBundlesSize(blob)
	}

	request := &node.StoreChunksRequest{
		BatchHeader: getBatchHeaderMessage(batchHeader),
		Blobs:       blobs,
	}

	return request, totalSize, nil
}

func GetStoreBlobsRequest(blobMessages []*core.EncodedBlobMessage, batchHeader *core.BatchHeader, useGnarkBundleEncoding bool) (*node.StoreBlobsRequest, int64, error) {
	blobs := make([]*node.Blob, len(blobMessages))
	totalSize := int64(0)
	for i, blob := range blobMessages {
		var err error
		blobs[i], err = getBlobMessage(blob, useGnarkBundleEncoding)
		if err != nil {
			return nil, 0, err
		}
		totalSize += getBundlesSize(blob)
	}

	request := &node.StoreBlobsRequest{
		Blobs:                blobs,
		ReferenceBlockNumber: uint32(batchHeader.ReferenceBlockNumber),
	}

	return request, totalSize, nil
}

func getBlobMessage(blob *core.EncodedBlobMessage, useGnarkBundleEncoding bool) (*node.Blob, error) {
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

	var err error
	bundles := make([]*node.Bundle, len(quorumHeaders))
	if useGnarkBundleEncoding {
		// the ordering of quorums in bundles must be same as in quorumHeaders
		for i, quorumHeader := range quorumHeaders {
			quorum := quorumHeader.QuorumId
			if chunksData, ok := blob.EncodedBundles[uint8(quorum)]; ok {
				if chunksData.Format != core.GnarkChunkEncodingFormat {
					chunksData, err = chunksData.ToGnarkFormat()
					if err != nil {
						return nil, err
					}
				}
				bundleBytes, err := chunksData.FlattenToBundle()
				if err != nil {
					return nil, err
				}
				bundles[i] = &node.Bundle{
					Bundle: bundleBytes,
				}
			} else {
				bundles[i] = &node.Bundle{
					// empty bundle for quorums operators are not part of
					Bundle: make([]byte, 0),
				}
			}
		}
	} else {
		// the ordering of quorums in bundles must be same as in quorumHeaders
		for i, quorumHeader := range quorumHeaders {
			quorum := quorumHeader.QuorumId
			if chunksData, ok := blob.EncodedBundles[uint8(quorum)]; ok {
				if chunksData.Format != core.GobChunkEncodingFormat {
					chunksData, err = chunksData.ToGobFormat()
					if err != nil {
						return nil, err
					}
				}
				bundles[i] = &node.Bundle{
					Chunks: chunksData.Chunks,
				}
			} else {
				bundles[i] = &node.Bundle{
					// empty bundle for quorums operators are not part of
					Chunks: make([][]byte, 0),
				}
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

func getBundlesSize(blob *core.EncodedBlobMessage) int64 {
	size := int64(0)
	for _, bundle := range blob.EncodedBundles {
		size += int64(bundle.Size())
	}
	return size
}
