package inmem

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/batcher"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/google/uuid"
)

var BatchNotFound = errors.New("batch not found")

type minibatchStore struct {
	// BatchRecords maps batch IDs to batch records
	BatchRecords map[uuid.UUID]*batcher.BatchRecord
	Dispersals   map[uuid.UUID]map[uint][]*batcher.MinibatchDispersal
	// BlobMinibatchMapping maps blob key to a map from batch ID to minibatch records
	BlobMinibatchMapping map[string]map[uuid.UUID]*batcher.BlobMinibatchMapping

	mu     sync.RWMutex
	logger logging.Logger
}

var _ batcher.MinibatchStore = (*minibatchStore)(nil)

func NewMinibatchStore(logger logging.Logger) batcher.MinibatchStore {
	return &minibatchStore{
		BatchRecords:         make(map[uuid.UUID]*batcher.BatchRecord),
		Dispersals:           make(map[uuid.UUID]map[uint][]*batcher.MinibatchDispersal),
		BlobMinibatchMapping: make(map[string]map[uuid.UUID]*batcher.BlobMinibatchMapping),

		logger: logger,
	}
}

func (m *minibatchStore) PutBatch(ctx context.Context, batch *batcher.BatchRecord) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.BatchRecords[batch.ID] = batch

	return nil
}

func (m *minibatchStore) GetBatch(ctx context.Context, batchID uuid.UUID) (*batcher.BatchRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	b, ok := m.BatchRecords[batchID]
	if !ok {
		return nil, BatchNotFound
	}
	return b, nil
}

func (m *minibatchStore) GetBatchesByStatus(ctx context.Context, status batcher.BatchStatus) ([]*batcher.BatchRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var batches []*batcher.BatchRecord
	for _, b := range m.BatchRecords {
		if b.Status == status {
			batches = append(batches, b)
		}
	}
	sort.Slice(batches, func(i, j int) bool {
		return batches[i].CreatedAt.Before(batches[j].CreatedAt)
	})
	return batches, nil
}

func (m *minibatchStore) UpdateBatchStatus(ctx context.Context, batchID uuid.UUID, status batcher.BatchStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	b, ok := m.BatchRecords[batchID]
	if !ok {
		return BatchNotFound
	}
	b.Status = status
	return nil
}

func (m *minibatchStore) MarkBatchFormed(ctx context.Context, batchID uuid.UUID, numMinibatches uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	b, ok := m.BatchRecords[batchID]
	if !ok {
		return BatchNotFound
	}
	b.NumMinibatches = numMinibatches
	b.Status = batcher.BatchStatusFormed
	return nil
}

func (m *minibatchStore) PutDispersal(ctx context.Context, dispersal *batcher.MinibatchDispersal) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.Dispersals[dispersal.BatchID]; !ok {
		m.Dispersals[dispersal.BatchID] = make(map[uint][]*batcher.MinibatchDispersal)
	}

	if _, ok := m.Dispersals[dispersal.BatchID][dispersal.MinibatchIndex]; !ok {
		m.Dispersals[dispersal.BatchID][dispersal.MinibatchIndex] = make([]*batcher.MinibatchDispersal, 0)
	}

	for _, r := range m.Dispersals[dispersal.BatchID][dispersal.MinibatchIndex] {
		if r.OperatorID == dispersal.OperatorID {
			// replace existing record
			*r = *dispersal
			return nil
		}
	}

	m.Dispersals[dispersal.BatchID][dispersal.MinibatchIndex] = append(m.Dispersals[dispersal.BatchID][dispersal.MinibatchIndex], dispersal)

	return nil
}

func (m *minibatchStore) UpdateDispersalResponse(ctx context.Context, dispersal *batcher.MinibatchDispersal, response *batcher.DispersalResponse) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.Dispersals[dispersal.BatchID][dispersal.MinibatchIndex]; !ok {
		return fmt.Errorf("dispersal not found")
	}

	for _, r := range m.Dispersals[dispersal.BatchID][dispersal.MinibatchIndex] {
		if r.OperatorID == dispersal.OperatorID {
			r.Signatures = response.Signatures
			r.RespondedAt = response.RespondedAt
			r.Error = response.Error
			return nil
		}
	}

	return nil
}

func (m *minibatchStore) GetDispersal(ctx context.Context, batchID uuid.UUID, minibatchIndex uint, opID core.OperatorID) (*batcher.MinibatchDispersal, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	requests, err := m.GetDispersalsByMinibatch(ctx, batchID, minibatchIndex)
	if err != nil {
		return nil, err
	}
	for _, r := range requests {
		if r.OperatorID == opID {
			return r, nil
		}
	}
	return nil, nil
}

func (m *minibatchStore) GetDispersalsByBatchID(ctx context.Context, batchID uuid.UUID) ([]*batcher.MinibatchDispersal, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, ok := m.Dispersals[batchID]; !ok {
		return nil, BatchNotFound
	}

	res := make([]*batcher.MinibatchDispersal, 0)
	for _, reqs := range m.Dispersals[batchID] {
		res = append(res, reqs...)
	}

	return res, nil
}

func (m *minibatchStore) GetDispersalsByMinibatch(ctx context.Context, batchID uuid.UUID, minibatchIndex uint) ([]*batcher.MinibatchDispersal, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, ok := m.Dispersals[batchID]; !ok {
		return nil, BatchNotFound
	}

	return m.Dispersals[batchID][minibatchIndex], nil
}

func (m *minibatchStore) GetBlobMinibatchMappings(ctx context.Context, blobKey disperser.BlobKey) ([]*batcher.BlobMinibatchMapping, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, ok := m.BlobMinibatchMapping[blobKey.String()]; !ok {
		return nil, nil
	}

	res := make([]*batcher.BlobMinibatchMapping, 0)
	for _, blobMinibatchMapping := range m.BlobMinibatchMapping[blobKey.String()] {
		res = append(res, blobMinibatchMapping)
	}

	return res, nil
}

func (m *minibatchStore) GetBlobMinibatchMappingsByBatchID(ctx context.Context, batchID uuid.UUID) ([]*batcher.BlobMinibatchMapping, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := make([]*batcher.BlobMinibatchMapping, 0)

	for _, batchToBlobMinibatchMapping := range m.BlobMinibatchMapping {
		for bID, blobMinibatchMapping := range batchToBlobMinibatchMapping {
			if bID == batchID {
				res = append(res, blobMinibatchMapping)
			}
		}
	}
	return res, nil
}

func (m *minibatchStore) PutBlobMinibatchMappings(ctx context.Context, blobMinibatchMappings []*batcher.BlobMinibatchMapping) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, blobMinibatchMapping := range blobMinibatchMappings {
		if blobMinibatchMapping.BlobKey == nil {
			return errors.New("blob key is nil")
		}
		blobKey := blobMinibatchMapping.BlobKey.String()

		if _, ok := m.BlobMinibatchMapping[blobKey]; !ok {
			m.BlobMinibatchMapping[blobKey] = make(map[uuid.UUID]*batcher.BlobMinibatchMapping)
		}

		m.BlobMinibatchMapping[blobKey][blobMinibatchMapping.BatchID] = blobMinibatchMapping
	}
	return nil
}

func (m *minibatchStore) GetLatestFormedBatch(ctx context.Context) (batch *batcher.BatchRecord, err error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	batches, err := m.GetBatchesByStatus(ctx, batcher.BatchStatusFormed)
	if err != nil {
		return nil, err
	}
	if len(batches) == 0 {
		return nil, nil
	}

	batch = batches[0]

	return batch, nil
}

func (m *minibatchStore) BatchDispersed(ctx context.Context, batchID uuid.UUID, numMinibatches uint) (bool, error) {
	dispersed := true
	dispersals, err := m.GetDispersalsByBatchID(ctx, batchID)
	if err != nil {
		return false, err
	}

	if len(dispersals) == 0 {
		return false, nil
	}

	minibatchIndices := make(map[uint]struct{})
	for _, resp := range dispersals {
		minibatchIndices[resp.MinibatchIndex] = struct{}{}
		if resp.RespondedAt.IsZero() || resp.Error != nil {
			dispersed = false
			m.logger.Info("response pending", "batchID", batchID, "minibatchIndex", resp.MinibatchIndex, "operatorID", resp.OperatorID.Hex())
		}
	}
	if len(minibatchIndices) != int(numMinibatches) {
		m.logger.Info("number of minibatches does not match", "batchID", batchID, "numMinibatches", numMinibatches, "minibatchIndices", len(minibatchIndices))
		return false, nil
	}
	for i := uint(0); i < numMinibatches; i++ {
		if _, ok := minibatchIndices[i]; !ok {
			m.logger.Info("minibatch missing", "batchID", batchID, "minibatchIndex", i, "numMinibatches", numMinibatches)
			return false, nil
		}
	}

	return dispersed, nil
}
