// These v2 methods are implemented in this separate file to keep the code organized.
// Note that there is no NodeV2 type and these methods are implemented in the existing Node type.

package node

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/tracing"
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

func (n *Node) DownloadBundles(
	ctx context.Context,
	batch *corev2.Batch,
	operatorState *core.OperatorState,
	probe *common.SequenceProbe,
) ([]*corev2.BlobShard, []*RawBundles, error) {

	probe.SetStage("prepare_to_download")

	ctx, span := tracing.TraceOperation(ctx, "DownloadBundles")
	defer span.End()

	relayClient, ok := n.RelayClient.Load().(clients.RelayClient)
	if !ok || relayClient == nil {
		return nil, nil, fmt.Errorf("relay client is not set")
	}

	blobVersionParams := n.BlobVersionParams.Load()
	if blobVersionParams == nil {
		return nil, nil, fmt.Errorf("blob version params is nil")
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
			blobParams, ok := blobVersionParams.Get(cert.BlobHeader.BlobVersion)
			if !ok {
				return nil, nil, fmt.Errorf("blob version %d not found", cert.BlobHeader.BlobVersion)
			}

			if _, ok := operatorState.Operators[quorum]; !ok {
				// operator is not part of the quorum or the quorum is not valid
				n.Logger.Debug("operator is not part of the quorum or the quorum is not valid",
					"quorum", quorum)
				continue
			}

			assgn, err := corev2.GetAssignment(operatorState, blobParams, quorum, n.Config.ID)
			if err != nil {
				n.Logger.Errorf("failed to get assignment: %v", err)
				continue
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

	probe.SetStage("download")
	// TODO (cody-littley) metric for the time until start of download
	// TODO (cody-littley) metric for the time of individual downloads

	bundleChan := make(chan response, len(requests))
	for relayKey := range requests {
		relayKey := relayKey
		req := requests[relayKey]
		n.DownloadPool.Submit(func() {
			ctxTimeout, cancel := context.WithTimeout(ctx, n.Config.ChunkDownloadTimeout)
			defer cancel()

			workerCtx, workerSpan := tracing.TraceOperation(ctxTimeout, "GetChunksByRange")
			bundles, err := relayClient.GetChunksByRange(workerCtx, relayKey, req.chunkRequests)
			workerSpan.End()

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

	responses := make([]response, len(requests))
	for i := 0; i < len(requests); i++ {
		responses[i] = <-bundleChan
	}

	probe.SetStage("deserialize")

	var err error
	for i := 0; i < len(requests); i++ {
		resp := responses[i]
		if resp.err != nil {
			// TODO (cody-littley) this is flaky, and will fail if any relay fails. We should retry failures
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
	ctx, span := tracing.TraceOperation(ctx, "ValidateBatchV2")
	defer span.End()

	if n.ValidatorV2 == nil {
		return fmt.Errorf("store v2 is not set")
	}

	if err := n.ValidatorV2.ValidateBatchHeader(ctx, batch.BatchHeader, batch.BlobCertificates); err != nil {
		return fmt.Errorf("failed to validate batch header: %v", err)
	}
	pool := workerpool.New(n.Config.NumBatchValidators)
	blobVersionParams := n.BlobVersionParams.Load()
	return n.ValidatorV2.ValidateBlobs(ctx, blobShards, blobVersionParams, pool, operatorState)
}

func (n *Node) InitTracingV2(ctx context.Context) error {
	// Initialize tracing if enabled
	if n.Config.Tracing.Enabled {
		tracingCfg := tracing.TracingConfig{
			Enabled:     n.Config.Tracing.Enabled,
			ServiceName: n.Config.Tracing.ServiceName + "-v2",
			Endpoint:    n.Config.Tracing.Endpoint,
			SampleRatio: n.Config.Tracing.SampleRatio,
		}

		telemetryShutdown, err := tracing.InitTelemetry(ctx, tracingCfg)
		if err != nil {
			n.Logger.Error("Failed to initialize V2 tracing", "err", err)
			// Continue with startup even if tracing fails
		} else {
			n.Logger.Info("Enabled tracing for Node V2", "endpoint", n.Config.Tracing.Endpoint)
			// Add cleanup handler for tracing when node shuts down
			go func() {
				<-ctx.Done()
				if err := telemetryShutdown(context.Background()); err != nil {
					n.Logger.Error("Failed to shutdown V2 telemetry", "err", err)
				}
			}()
		}
	}
	return nil
}
