// These v2 methods are implemented in this separate file to keep the code organized.
// Note that there is no NodeV2 type and these methods are implemented in the existing Node type.

package node

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/Layr-Labs/eigenda/api/clients/v2/relay"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
)

type requestMetadata struct {
	blobShardIndex int
	assignment     corev2.Assignment
}
type relayRequest struct {
	chunkRequests []*relay.ChunkRequestByRange
	metadata      []*requestMetadata
}
type response struct {
	metadata []*requestMetadata
	bundles  [][]byte
	err      error
}

type RawBundle struct {
	BlobCertificate *corev2.BlobCertificate
	Bundle          []byte
}

// Determines where to find the chunks we need to download for a given batch. For each chunk in a batch, there will
// be one or more relays that are responsible for serving that chunk. This function determines which relays to contact
// for each chunk, and sorts the requests by relayID to support batching. Additionally, this method also calculates
// the size of the chunk data that will be downloaded, in bytes.
func (n *Node) DetermineChunkLocations(
	batch *corev2.Batch,
	operatorState *core.OperatorState,
	probe *common.SequenceProbe,
) (downloadSizeInBytes uint64, relayRequests map[corev2.RelayKey]*relayRequest, err error) {

	probe.SetStage("determine_chunk_locations")

	blobVersionParams := n.BlobVersionParams.Load()
	if blobVersionParams == nil {
		return 0, nil, fmt.Errorf("blob version params is nil")
	}

	relayRequests = make(map[corev2.RelayKey]*relayRequest)

	for i, cert := range batch.BlobCertificates {
		blobKey, err := cert.BlobHeader.BlobKey()
		if err != nil {
			return 0, nil, fmt.Errorf("failed to get blob key: %w", err)
		}

		if len(cert.RelayKeys) == 0 {
			return 0, nil, fmt.Errorf("no relay keys in the certificate")
		}
		relayIndex := rand.Intn(len(cert.RelayKeys))
		relayKey := cert.RelayKeys[relayIndex]

		blobParams, ok := blobVersionParams.Get(cert.BlobHeader.BlobVersion)
		if !ok {
			return 0, nil, fmt.Errorf("blob version %d not found", cert.BlobHeader.BlobVersion)
		}

		assgn, err := corev2.GetAssignmentForBlob(operatorState, blobParams, cert.BlobHeader.QuorumNumbers, n.Config.ID)
		if err != nil {
			n.Logger.Errorf("failed to get assignment: %v", err)
			continue
		}

		chunkLength, err := blobParams.GetChunkLength(uint32(cert.BlobHeader.BlobCommitments.Length))
		if err != nil {
			return 0, nil, fmt.Errorf("failed to get chunk length: %w", err)
		}
		downloadSizeInBytes += uint64(assgn.NumChunks() * chunkLength)

		req, ok := relayRequests[relayKey]
		if !ok {
			req = &relayRequest{
				chunkRequests: make([]*relay.ChunkRequestByRange, 0),
				metadata:      make([]*requestMetadata, 0),
			}
			relayRequests[relayKey] = req
		}
		// Chunks from one blob are requested to the same relay
		rangeRequests := convertIndicesToRangeRequests(blobKey, assgn.Indices)
		req.chunkRequests = append(req.chunkRequests, rangeRequests...)

		// Code expects one metadata entry per range request
		for range rangeRequests {
			req.metadata = append(req.metadata, &requestMetadata{
				blobShardIndex: i,
				assignment:     assgn,
			})
		}
	}

	return downloadSizeInBytes, relayRequests, nil
}

// Given a list of chunk indices we want to download, create a list of relay requests by range.
// Although indices may not be contiguous, it is safe to assume that they will be "mostly contiguous".
// In practice, we should expect to see at most one continuous range of indices per quorum.
//
// Important: the provided indices MUST be in (mostly) sorted order in order to collapse into ranges correctly.
// Unsorted indices may lead to a very large number of range requests being generated. The current chunk assignment
// logic produces mostly sorted indices, so this is not an issue at present.
//
// Eventually, the assignment logic ought to be refactored to return ranges of chunks instead of individual
// indices, but the required changes are non-trivial.
func convertIndicesToRangeRequests(blobKey corev2.BlobKey, indices []uint32) []*relay.ChunkRequestByRange {
	requests := make([]*relay.ChunkRequestByRange, 0)
	if len(indices) == 0 {
		return requests
	}

	startIndex := indices[0]
	for i := 1; i < len(indices); i++ {
		if indices[i] != indices[i-1]+1 {
			// break in continuity, create a request for the previous range
			request := &relay.ChunkRequestByRange{
				BlobKey: blobKey,
				Start:   startIndex,       // inclusive
				End:     indices[i-1] + 1, // exclusive
			}
			requests = append(requests, request)
			startIndex = indices[i]
		}
	}

	// add the last range
	request := &relay.ChunkRequestByRange{
		BlobKey: blobKey,
		Start:   startIndex,                  // inclusive
		End:     indices[len(indices)-1] + 1, // exclusive
	}
	requests = append(requests, request)

	return requests
}

// This method takes a "download plan" from DetermineChunkLocations() and downloads the chunks from the relays.
// It also deserializes the responses from the relays into BlobShards and RawBundles.
func (n *Node) DownloadChunksFromRelays(
	ctx context.Context,
	batch *corev2.Batch,
	relayRequests map[corev2.RelayKey]*relayRequest,
	probe *common.SequenceProbe,
) (blobShards []*corev2.BlobShard, rawBundles []*RawBundle, err error) {

	blobShards = make([]*corev2.BlobShard, len(batch.BlobCertificates))
	rawBundles = make([]*RawBundle, len(batch.BlobCertificates))
	for i, cert := range batch.BlobCertificates {
		blobShards[i] = &corev2.BlobShard{
			BlobCertificate: cert,
		}
		rawBundles[i] = &RawBundle{
			BlobCertificate: cert,
		}
	}

	relayClient, ok := n.RelayClient.Load().(relay.RelayClient)
	if !ok || relayClient == nil {
		return nil, nil, fmt.Errorf("relay client is not set")
	}

	probe.SetStage("download")

	bundleChan := make(chan response, len(relayRequests))
	for relayKey := range relayRequests {
		relayKey := relayKey
		req := relayRequests[relayKey]
		n.DownloadPool.Submit(func() {
			ctxTimeout, cancel := context.WithTimeout(ctx, n.Config.ChunkDownloadTimeout)
			defer cancel()
			bundles, err := relayClient.GetChunksByRange(ctxTimeout, relayKey, req.chunkRequests)
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

	responses := make([]response, len(relayRequests))
	for i := 0; i < len(relayRequests); i++ {
		responses[i] = <-bundleChan
	}

	probe.SetStage("deserialize")

	for i := 0; i < len(relayRequests); i++ {
		resp := responses[i]
		if resp.err != nil {
			// TODO (cody-littley) this is flaky, and will fail if any relay fails. We should retry failures
			return nil, nil, fmt.Errorf("failed to get chunks from relays: %v", resp.err)
		}

		if len(resp.bundles) != len(resp.metadata) {
			return nil, nil,
				fmt.Errorf("number of bundles and metadata do not match (%d != %d)",
					len(resp.bundles), len(resp.metadata))
		}

		for j, bundle := range resp.bundles {
			metadata := resp.metadata[j]
			var err error
			blobShards[metadata.blobShardIndex].Bundle, err = new(core.Bundle).Deserialize(bundle)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to deserialize bundle: %v", err)
			}
			rawBundles[metadata.blobShardIndex].Bundle = bundle

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
	blobVersionParams := n.BlobVersionParams.Load()
	return n.ValidatorV2.ValidateBlobs(ctx, blobShards, blobVersionParams, n.ValidationPool, operatorState)
}
