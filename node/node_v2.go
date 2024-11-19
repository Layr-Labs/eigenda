// These v2 methods are implemented in this separate file to keep the code organized.
// Note that there is no NodeV2 type and these methods are implemented in the existing Node type.

package node

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/gammazero/workerpool"
)

type requestMetadata struct {
	blobShardIndex int
	quorum         core.QuorumID
}
type relayRequest struct {
	chunkRequests []*clients.ChunkRequestByRange
	metadata      []*requestMetadata
}
type response struct {
	metadata []*requestMetadata
	bundles  [][]byte
	err      error
}

type RawBundles struct {
	BlobCertificate *corev2.BlobCertificate
	Bundles         map[core.QuorumID][]byte
}

func (n *Node) DownloadBundles(ctx context.Context, batch *corev2.Batch, operatorState *core.OperatorState) ([]*corev2.BlobShard, []*RawBundles, error) {
	if n.RelayClient == nil {
		return nil, nil, fmt.Errorf("relay client is not set")
	}

	blobShards := make([]*corev2.BlobShard, len(batch.BlobCertificates))
	rawBundles := make([]*RawBundles, len(batch.BlobCertificates))
	requests := make(map[corev2.RelayKey]*relayRequest)
	for i, cert := range batch.BlobCertificates {
		blobKey, err := cert.BlobHeader.BlobKey()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get blob key: %v", err)
		}

		if len(cert.RelayKeys) == 0 {
			return nil, nil, fmt.Errorf("no relay keys in the certificate")
		}
		blobShards[i] = &corev2.BlobShard{
			BlobCertificate: cert,
			Bundles:         make(map[core.QuorumID]core.Bundle),
		}
		rawBundles[i] = &RawBundles{
			BlobCertificate: cert,
			Bundles:         make(map[core.QuorumID][]byte),
		}
		relayIndex := rand.Intn(len(cert.RelayKeys))
		relayKey := cert.RelayKeys[relayIndex]
		for _, quorum := range cert.BlobHeader.QuorumNumbers {
			assgn, err := corev2.GetAssignment(operatorState, batch.BlobCertificates[0].BlobHeader.BlobVersion, quorum, n.Config.ID)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to get assignments: %v", err)
			}

			req, ok := requests[relayKey]
			if !ok {
				req = &relayRequest{
					chunkRequests: make([]*clients.ChunkRequestByRange, 0),
					metadata:      make([]*requestMetadata, 0),
				}
				requests[relayKey] = req
			}
			// Chunks from one blob are requested to the same relay
			req.chunkRequests = append(req.chunkRequests, &clients.ChunkRequestByRange{
				BlobKey: blobKey,
				Start:   assgn.StartIndex,
				End:     assgn.StartIndex + assgn.NumChunks,
			})
			req.metadata = append(req.metadata, &requestMetadata{
				blobShardIndex: i,
				quorum:         quorum,
			})
		}
	}

	pool := workerpool.New(len(requests))
	bundleChan := make(chan response, len(requests))
	for relayKey := range requests {
		relayKey := relayKey
		req := requests[relayKey]
		pool.Submit(func() {
			bundles, err := n.RelayClient.GetChunksByRange(ctx, relayKey, req.chunkRequests)
			if err != nil {
				n.Logger.Errorf("failed to get chunks from relays: %v", err)
				bundleChan <- response{
					metadata: nil,
					bundles:  nil,
					err:      err,
				}
				return
			}
			bundleChan <- response{
				metadata: req.metadata,
				bundles:  bundles,
				err:      nil,
			}
		})
	}
	pool.StopWait()

	var err error
	for i := 0; i < len(requests); i++ {
		resp := <-bundleChan
		if resp.err != nil {
			return nil, nil, fmt.Errorf("failed to get chunks from relays: %v", resp.err)
		}
		for i, bundle := range resp.bundles {
			metadata := resp.metadata[i]
			blobShards[metadata.blobShardIndex].Bundles[metadata.quorum], err = new(core.Bundle).Deserialize(bundle)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to deserialize bundle: %v", err)
			}
			rawBundles[metadata.blobShardIndex].Bundles[metadata.quorum] = bundle
		}
	}

	return blobShards, rawBundles, nil
}

func (n *Node) ValidateBatchV2(
	ctx context.Context,
	batch *corev2.Batch,
	blobShards []*corev2.BlobShard,
	operatorState *core.OperatorState,
) error {
	if n.ValidatorV2 == nil {
		return fmt.Errorf("store v2 is not set")
	}

	if err := n.ValidatorV2.ValidateBatchHeader(ctx, batch.BatchHeader, batch.BlobCertificates); err != nil {
		return fmt.Errorf("failed to validate batch header: %v", err)
	}
	pool := workerpool.New(n.Config.NumBatchValidators)
	return n.ValidatorV2.ValidateBlobs(ctx, blobShards, pool, operatorState)
}
