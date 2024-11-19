package controller

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

var errNoBlobsToDispatch = errors.New("no blobs to dispatch")

type DispatcherConfig struct {
	PullInterval time.Duration

	FinalizationBlockDelay uint64
	NodeRequestTimeout     time.Duration
	NumRequestRetries      int
}

type Dispatcher struct {
	DispatcherConfig

	blobMetadataStore *blobstore.BlobMetadataStore
	pool              common.WorkerPool
	chainState        core.IndexedChainState
	aggregator        core.SignatureAggregator
	nodeClientManager NodeClientManager
	logger            logging.Logger

	lastUpdatedAt uint64
}

type batchData struct {
	Batch           *corev2.Batch
	BatchHeaderHash [32]byte
	BlobKeys        []corev2.BlobKey
	OperatorState   *core.IndexedOperatorState
}

func NewDispatcher(
	config DispatcherConfig,
	blobMetadataStore *blobstore.BlobMetadataStore,
	pool common.WorkerPool,
	chainState core.IndexedChainState,
	aggregator core.SignatureAggregator,
	nodeClientManager NodeClientManager,
	logger logging.Logger,
) (*Dispatcher, error) {
	return &Dispatcher{
		DispatcherConfig: config,

		blobMetadataStore: blobMetadataStore,
		pool:              pool,
		chainState:        chainState,
		aggregator:        aggregator,
		nodeClientManager: nodeClientManager,
		logger:            logger.With("component", "Dispatcher"),

		lastUpdatedAt: 0,
	}, nil
}

func (d *Dispatcher) Start(ctx context.Context) error {
	err := d.chainState.Start(ctx)
	if err != nil {
		return fmt.Errorf("failed to start chain state: %w", err)
	}

	go func() {
		ticker := time.NewTicker(d.PullInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				sigChan, batchData, err := d.HandleBatch(ctx)
				if err != nil {
					if errors.Is(err, errNoBlobsToDispatch) {
						d.logger.Warn("no blobs to dispatch")
					} else {
						d.logger.Error("failed to process a batch", "err", err)
					}
					continue
				}

				go func() {
					err := d.HandleSignatures(ctx, batchData, sigChan)
					if err != nil {
						d.logger.Error("failed to handle signatures", "err", err)
					}
				}()
			}
		}
	}()

	return nil

}

func (d *Dispatcher) HandleBatch(ctx context.Context) (chan core.SigningMessage, *batchData, error) {
	currentBlockNumber, err := d.chainState.GetCurrentBlockNumber()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get current block number: %w", err)
	}
	referenceBlockNumber := uint64(currentBlockNumber) - d.FinalizationBlockDelay

	// Get a batch of blobs to dispatch
	// This also writes a batch header and blob verification info for each blob in metadata store
	batchData, err := d.NewBatch(ctx, referenceBlockNumber)
	if err != nil {
		return nil, nil, err
	}

	batch := batchData.Batch
	state := batchData.OperatorState
	sigChan := make(chan core.SigningMessage, len(state.IndexedOperators))
	for opID, op := range state.IndexedOperators {
		opID := opID
		op := op
		host, dispersalPort, _, err := core.ParseOperatorSocket(op.Socket)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse operator socket: %w", err)
		}

		client, err := d.nodeClientManager.GetClient(host, dispersalPort)
		if err != nil {
			d.logger.Error("failed to get node client", "operator", opID, "err", err)
			continue
		}

		d.pool.Submit(func() {

			req := &corev2.DispersalRequest{
				OperatorID: opID,
				// TODO: get OperatorAddress
				OperatorAddress: gethcommon.Address{},
				Socket:          op.Socket,
				DispersedAt:     uint64(time.Now().UnixNano()),
				BatchHeader:     *batch.BatchHeader,
			}
			err := d.blobMetadataStore.PutDispersalRequest(ctx, req)
			if err != nil {
				d.logger.Error("failed to put dispersal request", "err", err)
				return
			}

			for i := 0; i < d.NumRequestRetries+1; i++ {
				sig, err := d.sendChunks(ctx, client, batch)
				if err == nil {
					storeErr := d.blobMetadataStore.PutDispersalResponse(ctx, &corev2.DispersalResponse{
						DispersalRequest: req,
						RespondedAt:      uint64(time.Now().UnixNano()),
						Signature:        sig.Bytes(),
						Error:            "",
					})
					if storeErr != nil {
						d.logger.Error("failed to put dispersal response", "err", storeErr)
					}

					sigChan <- core.SigningMessage{
						Signature:            sig,
						Operator:             opID,
						BatchHeaderHash:      batchData.BatchHeaderHash,
						AttestationLatencyMs: 0, // TODO: calculate latency
						Err:                  nil,
					}

					break
				}

				d.logger.Warn("failed to send chunks", "operator", opID, "NumAttempts", i, "err", err)
				time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second) // Wait before retrying
			}
		})
	}

	return sigChan, batchData, nil
}

// HandleSignatures receives signatures from operators, validates, and aggregates them
func (d *Dispatcher) HandleSignatures(ctx context.Context, batchData *batchData, sigChan chan core.SigningMessage) error {
	quorumAttestation, err := d.aggregator.ReceiveSignatures(ctx, batchData.OperatorState, batchData.BatchHeaderHash, sigChan)
	if err != nil {
		return fmt.Errorf("failed to receive and validate signatures: %w", err)
	}
	quorums := make([]core.QuorumID, len(quorumAttestation.QuorumResults))
	i := 0
	for quorumID := range quorumAttestation.QuorumResults {
		quorums[i] = quorumID
		i++
	}
	aggSig, err := d.aggregator.AggregateSignatures(ctx, d.chainState, uint(batchData.Batch.BatchHeader.ReferenceBlockNumber), quorumAttestation, quorums)
	if err != nil {
		return fmt.Errorf("failed to aggregate signatures: %w", err)
	}
	err = d.blobMetadataStore.PutAttestation(ctx, &corev2.Attestation{
		BatchHeader:      batchData.Batch.BatchHeader,
		AttestedAt:       uint64(time.Now().UnixNano()),
		NonSignerPubKeys: aggSig.NonSigners,
		APKG2:            aggSig.AggPubKey,
		QuorumAPKs:       aggSig.QuorumAggPubKeys,
		Sigma:            aggSig.AggSignature,
		QuorumNumbers:    quorums,
	})
	if err != nil {
		return fmt.Errorf("failed to put attestation: %w", err)
	}

	err = d.updateBatchStatus(ctx, batchData.BlobKeys, v2.Certified)
	if err != nil {
		return fmt.Errorf("failed to mark blobs as certified: %w", err)
	}

	return nil
}

// NewBatch creates a batch of blobs to dispatch
// Warning: This function is not thread-safe
func (d *Dispatcher) NewBatch(ctx context.Context, referenceBlockNumber uint64) (*batchData, error) {
	blobMetadatas, err := d.blobMetadataStore.GetBlobMetadataByStatus(ctx, v2.Encoded, d.lastUpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get blob metadata by status: %w", err)
	}

	if len(blobMetadatas) == 0 {
		return nil, errNoBlobsToDispatch
	}

	state, err := d.GetOperatorState(ctx, blobMetadatas, referenceBlockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get operator state: %w", err)
	}

	keys := make([]corev2.BlobKey, len(blobMetadatas))
	newLastUpdatedAt := d.lastUpdatedAt
	for i, metadata := range blobMetadatas {
		if metadata == nil || metadata.BlobHeader == nil {
			return nil, fmt.Errorf("invalid blob metadata")
		}
		blobKey, err := metadata.BlobHeader.BlobKey()
		if err != nil {
			return nil, fmt.Errorf("failed to get blob key: %w", err)
		}
		keys[i] = blobKey
		if metadata.UpdatedAt > newLastUpdatedAt {
			newLastUpdatedAt = metadata.UpdatedAt
		}
	}

	certs, _, err := d.blobMetadataStore.GetBlobCertificates(ctx, keys)
	if err != nil {
		return nil, fmt.Errorf("failed to get blob certificates: %w", err)
	}

	if len(certs) != len(keys) {
		return nil, fmt.Errorf("blob certificates not found for all blob keys")
	}

	certsMap := make(map[corev2.BlobKey]*corev2.BlobCertificate, len(certs))
	for _, cert := range certs {
		blobKey, err := cert.BlobHeader.BlobKey()
		if err != nil {
			return nil, fmt.Errorf("failed to get blob key: %w", err)
		}

		certsMap[blobKey] = cert
	}

	// Keep the order of certs the same as the order of keys
	for i, key := range keys {
		c, ok := certsMap[key]
		if !ok {
			return nil, fmt.Errorf("blob certificate not found for blob key %s", key.Hex())
		}
		certs[i] = c
	}

	batchHeader := &corev2.BatchHeader{
		BatchRoot:            [32]byte{},
		ReferenceBlockNumber: referenceBlockNumber,
	}

	tree, err := corev2.BuildMerkleTree(certs)
	if err != nil {
		return nil, fmt.Errorf("failed to build merkle tree: %w", err)
	}

	copy(batchHeader.BatchRoot[:], tree.Root())

	batchHeaderHash, err := batchHeader.Hash()
	if err != nil {
		return nil, fmt.Errorf("failed to hash batch header: %w", err)
	}

	err = d.blobMetadataStore.PutBatchHeader(ctx, batchHeader)
	if err != nil {
		return nil, fmt.Errorf("failed to put batch header: %w", err)
	}

	// accumulate verification infos in a map to avoid duplicate entries
	// batch write operation fails if there are duplicate entries
	verificationInfoMap := make(map[corev2.BlobKey]*corev2.BlobVerificationInfo)
	for i, cert := range certs {
		if cert == nil || cert.BlobHeader == nil {
			return nil, fmt.Errorf("invalid blob certificate")
		}
		blobKey, err := cert.BlobHeader.BlobKey()
		if err != nil {
			return nil, fmt.Errorf("failed to get blob key: %w", err)
		}

		merkleProof, err := tree.GenerateProofWithIndex(uint64(i), 0)
		if err != nil {
			return nil, fmt.Errorf("failed to generate merkle proof: %w", err)
		}

		verificationInfoMap[blobKey] = &corev2.BlobVerificationInfo{
			BatchHeader:    batchHeader,
			BlobKey:        blobKey,
			BlobIndex:      uint32(i),
			InclusionProof: core.SerializeMerkleProof(merkleProof),
		}
	}

	verificationInfos := make([]*corev2.BlobVerificationInfo, len(verificationInfoMap))
	i := 0
	for _, v := range verificationInfoMap {
		verificationInfos[i] = v
		i++
	}
	err = d.blobMetadataStore.PutBlobVerificationInfos(ctx, verificationInfos)
	if err != nil {
		return nil, fmt.Errorf("failed to put blob verification infos: %w", err)
	}

	d.lastUpdatedAt = newLastUpdatedAt
	return &batchData{
		Batch: &corev2.Batch{
			BatchHeader:      batchHeader,
			BlobCertificates: certs,
		},
		BatchHeaderHash: batchHeaderHash,
		BlobKeys:        keys,
		OperatorState:   state,
	}, nil
}

// GetOperatorState returns the operator state for the given quorums at the given block number
func (d *Dispatcher) GetOperatorState(ctx context.Context, metadatas []*v2.BlobMetadata, blockNumber uint64) (*core.IndexedOperatorState, error) {
	quorums := make(map[core.QuorumID]struct{}, 0)
	for _, m := range metadatas {
		for _, quorum := range m.BlobHeader.QuorumNumbers {
			quorums[quorum] = struct{}{}
		}
	}

	quorumIds := make([]core.QuorumID, len(quorums))
	i := 0
	for id := range quorums {
		quorumIds[i] = id
		i++
	}

	// GetIndexedOperatorState should return state for valid quorums only
	return d.chainState.GetIndexedOperatorState(ctx, uint(blockNumber), quorumIds)
}

func (d *Dispatcher) sendChunks(ctx context.Context, client clients.NodeClientV2, batch *corev2.Batch) (*core.Signature, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, d.NodeRequestTimeout)
	defer cancel()

	sig, err := client.StoreChunks(ctxWithTimeout, batch)
	if err != nil {
		return nil, fmt.Errorf("failed to store chunks: %w", err)
	}

	return sig, nil
}

func (d *Dispatcher) updateBatchStatus(ctx context.Context, keys []corev2.BlobKey, status v2.BlobStatus) error {
	for _, key := range keys {
		err := d.blobMetadataStore.UpdateBlobStatus(ctx, key, status)
		if err != nil {
			d.logger.Error("failed to update blob status", "blobKey", key.Hex(), "status", status.String(), "err", err)
		}
	}
	return nil
}
