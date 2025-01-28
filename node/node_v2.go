// These v2 methods are implemented in this separate file to keep the code organized.
// Note that there is no NodeV2 type and these methods are implemented in the existing Node type.

package node

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/common"

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

		// err := n.validateDispersalRequest(ctx, cert, n.onchainState.Load())
		// if err != nil {
		// 	return nil, nil, fmt.Errorf("failed to validate dispersal request: %v", err)
		// }

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

	bundleChan := make(chan response, len(requests))
	for relayKey := range requests {
		relayKey := relayKey
		req := requests[relayKey]
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

// func (n *Node) validateDispersalRequest(ctx context.Context, blobCert *corev2.BlobCertificate, onchainState *eth.OnchainState) error {
// 	if len(blobCert.Signature) != 65 {
// 		return api.NewErrorInvalidArg(fmt.Sprintf("signature is expected to be 65 bytes, but got %d bytes", len(blobCert.Signature)))
// 	}

// 	blobLength := blobCert.BlobHeader.BlobCommitments.Length
// 	if blobLength == 0 {
// 		return api.NewErrorInvalidArg("blob length must be greater than 0")
// 	}
// 	if blobLength > uint(onchainState.MaxNumSymbolsPerBlob) {
// 		return api.NewErrorInvalidArg("blob length too big")
// 	}
// 	if blobLength != encoding.NextPowerOf2(blobLength) {
// 		return api.NewErrorInvalidArg("invalid blob length, must be a power of 2")
// 	}

// 	blobHeader := blobCert.BlobHeader
// 	if blobHeader.PaymentMetadata == (core.PaymentMetadata{}) {
// 		return api.NewErrorInvalidArg("payment metadata is required")
// 	}

// 	if len(blobHeader.PaymentMetadata.AccountID) == 0 || (blobHeader.PaymentMetadata.ReservationPeriod == 0 && blobHeader.PaymentMetadata.CumulativePayment.Cmp(big.NewInt(0)) == 0) {
// 		return api.NewErrorInvalidArg("invalid payment metadata")
// 	}

// 	quorumNumbers := blobHeader.QuorumNumbers
// 	if len(quorumNumbers) == 0 {
// 		return api.NewErrorInvalidArg("blob header must contain at least one quorum number")
// 	}

// 	if len(quorumNumbers) > int(onchainState.QuorumCount) {
// 		return api.NewErrorInvalidArg(fmt.Sprintf("too many quorum numbers specified: maximum is %d", onchainState.QuorumCount))
// 	}

// 	for _, quorum := range quorumNumbers {
// 		if quorum > corev2.MaxQuorumID || uint8(quorum) >= onchainState.QuorumCount {
// 			return api.NewErrorInvalidArg(fmt.Sprintf("invalid quorum number %d; maximum is %d", quorum, onchainState.QuorumCount))
// 		}
// 	}

// 	// // validate every 32 bytes is a valid field element
// 	// _, err = rs.ToFrArray(blob)
// 	// if err != nil {
// 	// 	s.logger.Error("failed to convert a 32bytes as a field element", "err", err)
// 	// 	return api.NewErrorInvalidArg("encountered an error to convert a 32-bytes into a valid field element, please use the correct format where every 32bytes(big-endian) is less than 21888242871839275222246405745257275088548364400416034343698204186575808495617")
// 	// }

// 	if _, ok := onchainState.BlobVersionParameters.Get(corev2.BlobVersion(blobHeader.BlobVersion)); !ok {
// 		return api.NewErrorInvalidArg(fmt.Sprintf("invalid blob version %d; valid blob versions are: %v", blobHeader.BlobVersion, onchainState.BlobVersionParameters.Keys()))
// 	}

// 	if err := n.authenticator.AuthenticateBlobRequest(blobHeader, blobCert.Signature); err != nil {
// 		return api.NewErrorInvalidArg(fmt.Sprintf("authentication failed: %s", err.Error()))
// 	}

// 	// handle payments and check rate limits
// 	timestamp := blobHeader.PaymentMetadata.Timestamp
// 	cumulativePayment := blobHeader.PaymentMetadata.CumulativePayment
// 	accountID := blobHeader.PaymentMetadata.AccountID

// 	paymentHeader := core.PaymentMetadata{
// 		AccountID:         accountID,
// 		Timestamp:         timestamp,
// 		CumulativePayment: cumulativePayment,
// 	}

// 	if err := n.meterer.MeterRequest(ctx, paymentHeader, blobLength, blobHeader.QuorumNumbers); err != nil {
// 		return api.NewErrorResourceExhausted(err.Error())
// 	}

// 	// should node run prover here? doesn't have data yet though
// 	return nil
// }
