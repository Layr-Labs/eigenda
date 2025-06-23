package dispatcher

import (
	"context"
	"errors"
	"fmt"
	"math/big"
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

	// Retry configuration
	EnableRetryAggregator       bool    // Enable failure aggregation and retry decision logic
	MaxAccountFailurePercentage float64 // Maximum percentage of stake failed for any single account before triggering retry
}

// DefaultRetryConfig returns default retry configuration values
func DefaultRetryConfig() *Config {
	return &Config{
		EnableRetryAggregator:       false, // Disabled by default until retry logic implemented
		MaxAccountFailurePercentage: 50.0,  // 50% of account-related failures
	}
}

// getOperatorStake calculates total stake for an operator across all quorums
func (c *dispatcher) getOperatorStake(state *core.IndexedOperatorState, operatorID core.OperatorID) *big.Int {
	totalStake := big.NewInt(0)
	for _, quorumOperators := range state.Operators {
		if opInfo, exists := quorumOperators[operatorID]; exists {
			totalStake.Add(totalStake, opInfo.Stake)
		}
	}
	return totalStake
}

// filterBlobsByAccountIDs filters out blobs from the specified account IDs
func (c *dispatcher) filterBlobsByAccountIDs(blobs []core.EncodedBlob, excludeAccountIDs []string) []core.EncodedBlob {
	if len(excludeAccountIDs) == 0 {
		return blobs
	}

	// Create a map for faster lookup
	excludeMap := make(map[string]bool)
	for _, accountID := range excludeAccountIDs {
		excludeMap[accountID] = true
	}

	var filteredBlobs []core.EncodedBlob
	excludedCount := 0

	for _, blob := range blobs {
		if blob.BlobHeader != nil && excludeMap[blob.BlobHeader.AccountID] {
			excludedCount++
			if c.logger != nil {
				c.logger.Info("Excluding blob from retry batch",
					"account_id", blob.BlobHeader.AccountID,
				)
			}
			continue
		}
		filteredBlobs = append(filteredBlobs, blob)
	}

	if excludedCount > 0 && c.logger != nil {
		c.logger.Info("Filtered blobs for retry batch",
			"original_count", len(blobs),
			"filtered_count", len(filteredBlobs),
			"excluded_count", excludedCount,
			"excluded_accounts", excludeAccountIDs,
		)
	}

	return filteredBlobs
}

type dispatcher struct {
	*Config

	logger  logging.Logger
	metrics *batcher.DispatcherMetrics
}

func NewDispatcher(
	cfg *Config,
	logger logging.Logger,
	metrics *batcher.DispatcherMetrics,
) *dispatcher {
	return &dispatcher{
		Config:  cfg,
		logger:  logger.With("component", "Dispatcher"),
		metrics: metrics,
	}
}

var _ disperser.Dispatcher = (*dispatcher)(nil)

// DisperseBatch distributes encoded blobs to all indexed operators and tracks failures
func (c *dispatcher) DisperseBatch(
	ctx context.Context,
	state *core.IndexedOperatorState,
	blobs []core.EncodedBlob,
	batchHeader *core.BatchHeader,
) chan core.SigningMessage {
	update := make(chan core.SigningMessage, len(state.IndexedOperators))

	// Send chunks to all operators with integrated failure analysis and logging
	retryDecision := c.sendAllChunks(ctx, state, blobs, batchHeader, update)

	// If retry is needed, create a filtered batch and dispatch again
	if c.EnableRetryAggregator && retryDecision != nil && retryDecision.ShouldRetry && len(retryDecision.TriggeringAccounts) > 0 {
		c.logger.Warn("Initiating batch retry with filtered blobs",
			"triggering_accounts", retryDecision.TriggeringAccounts,
			"original_blob_count", len(blobs),
		)

		// Filter out blobs from triggering accounts
		filteredBlobs := c.filterBlobsByAccountIDs(blobs, retryDecision.TriggeringAccounts)

		// Only retry if we have remaining blobs after filtering
		if len(filteredBlobs) > 0 {
			c.logger.Info("Dispatching retry batch",
				"filtered_blob_count", len(filteredBlobs),
				"excluded_accounts", retryDecision.TriggeringAccounts,
			)

			// Create new update channel for retry batch
			retryUpdate := make(chan core.SigningMessage, len(state.IndexedOperators))

			// Send filtered blobs (without retry aggregator to avoid infinite recursion)
			go func() {
				defer close(retryUpdate)

				// Temporarily disable retry aggregator for the retry attempt
				originalRetryEnabled := c.EnableRetryAggregator
				c.EnableRetryAggregator = false

				c.sendAllChunks(ctx, state, filteredBlobs, batchHeader, retryUpdate)

				// Restore original retry setting
				c.EnableRetryAggregator = originalRetryEnabled
			}()

			// Forward retry results to original update channel
			go func() {
				for retryMsg := range retryUpdate {
					update <- retryMsg
				}
			}()
		} else {
			c.logger.Warn("No blobs remaining after filtering triggering accounts",
				"excluded_accounts", retryDecision.TriggeringAccounts,
			)
		}
	}

	return update
}

// sendAllChunks distributes chunks to operators and aggregates failures with stake tracking
func (c *dispatcher) sendAllChunks(
	ctx context.Context,
	state *core.IndexedOperatorState,
	blobs []core.EncodedBlob,
	batchHeader *core.BatchHeader,
	update chan core.SigningMessage,
) *RetryDecision {
	// Only create failure aggregator if retry analysis is enabled
	var failureAggregator *FailureAggregator
	if c.EnableRetryAggregator {
		failureAggregator = NewFailureAggregator(c.logger)
	}

	for id, op := range state.IndexedOperators {
		operatorStake := c.getOperatorStake(state, id)
		if failureAggregator != nil {
			failureAggregator.AddTotalStake(operatorStake)
		}

		go func(op core.IndexedOperatorInfo, id core.OperatorID, stake *big.Int) {
			blobMessages := make([]*core.EncodedBlobMessage, 0)
			hasAnyBundles := false
			batchHeaderHash, err := batchHeader.GetBatchHeaderHash()
			if err != nil {
				update <- core.SigningMessage{
					Err:                  fmt.Errorf("failed to get batch header hash: %w", err),
					Signature:            nil,
					Operator:             id,
					BatchHeaderHash:      [32]byte{},
					AttestationLatencyMs: -1,
				}
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
				// Track operator failure with batch meterer error parsing if aggregator enabled
				if failureAggregator != nil {
					operatorFailure := failureAggregator.createOperatorFailure(id, op.Socket, stake, err)
					failureAggregator.AddFailure(operatorFailure)

					// Log batch meterer errors for account-level monitoring
					failureAggregator.LogBatchMeterError(operatorFailure, id, op.Socket)
				}

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
		}, id, operatorStake)
	}

	// Perform failure analysis and logging if aggregator is enabled
	var retryDecision *RetryDecision
	if c.EnableRetryAggregator && failureAggregator != nil {
		retryDecision = failureAggregator.ShouldRetryBatch(c.MaxAccountFailurePercentage)
		failureAggregator.LogRetryDecision(retryDecision)
		failureAggregator.LogFailureStatistics()
	}

	return retryDecision
}

func (c *dispatcher) sendChunks(
	ctx context.Context,
	blobs []*core.EncodedBlobMessage,
	batchHeader *core.BatchHeader,
	op *core.IndexedOperatorInfo,
) (*core.Signature, error) {
	// TODO Add secure Grpc

	conn, err := grpc.NewClient(
		core.OperatorSocket(op.Socket).GetV1DispersalSocket(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		c.logger.Warn("Disperser cannot connect to operator dispersal socket",
			"dispersal_socket", core.OperatorSocket(op.Socket).GetV1DispersalSocket(),
			"err", err)
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
	c.logger.Debug("sending chunks to operator",
		"operator", op.Socket,
		"num blobs", len(blobs),
		"size", totalSize,
		"request message size", proto.Size(request),
		"request serialization time", time.Since(start),
		"use Gnark chunk encoding", c.EnableGnarkBundleEncoding)
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

func GetStoreChunksRequest(
	blobMessages []*core.EncodedBlobMessage,
	batchHeader *core.BatchHeader,
	useGnarkBundleEncoding bool,
) (*node.StoreChunksRequest, int64, error) {
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

func GetStoreBlobsRequest(
	blobMessages []*core.EncodedBlobMessage,
	batchHeader *core.BatchHeader,
	useGnarkBundleEncoding bool,
) (*node.StoreBlobsRequest, int64, error) {
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
