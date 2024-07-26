package inmem

import (
	"context"
	"errors"
	"sort"
	"sync"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser/batcher"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/google/uuid"
)

var BatchNotFound = errors.New("batch not found")

type minibatchStore struct {
	// BatchRecords maps batch IDs to batch records
	BatchRecords map[uuid.UUID]*batcher.BatchRecord
	// MinibatchRecords maps batch IDs to a map from minibatch indices to minibatch records
	MinibatchRecords map[uuid.UUID]map[uint]*batcher.MinibatchRecord
	// DispersalRequests maps batch IDs to a map from minibatch indices to dispersal requests
	DispersalRequests map[uuid.UUID]map[uint][]*batcher.DispersalRequest
	// DispersalResponses maps batch IDs to a map from minibatch indices to dispersal responses
	DispersalResponses map[uuid.UUID]map[uint][]*batcher.DispersalResponse

	mu     sync.RWMutex
	logger logging.Logger
}

var _ batcher.MinibatchStore = (*minibatchStore)(nil)

func NewMinibatchStore(logger logging.Logger) batcher.MinibatchStore {
	return &minibatchStore{
		BatchRecords:       make(map[uuid.UUID]*batcher.BatchRecord),
		MinibatchRecords:   make(map[uuid.UUID]map[uint]*batcher.MinibatchRecord),
		DispersalRequests:  make(map[uuid.UUID]map[uint][]*batcher.DispersalRequest),
		DispersalResponses: make(map[uuid.UUID]map[uint][]*batcher.DispersalResponse),

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

func (m *minibatchStore) PutMinibatch(ctx context.Context, minibatch *batcher.MinibatchRecord) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.MinibatchRecords[minibatch.BatchID]; !ok {
		m.MinibatchRecords[minibatch.BatchID] = make(map[uint]*batcher.MinibatchRecord)
	}
	m.MinibatchRecords[minibatch.BatchID][minibatch.MinibatchIndex] = minibatch

	return nil
}

func (m *minibatchStore) GetMinibatch(ctx context.Context, batchID uuid.UUID, minibatchIndex uint) (*batcher.MinibatchRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, ok := m.MinibatchRecords[batchID]; !ok {
		return nil, BatchNotFound
	}
	return m.MinibatchRecords[batchID][minibatchIndex], nil
}

func (m *minibatchStore) GetMinibatches(ctx context.Context, batchID uuid.UUID) ([]*batcher.MinibatchRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, ok := m.MinibatchRecords[batchID]; !ok {
		return nil, nil
	}

	res := make([]*batcher.MinibatchRecord, 0, len(m.MinibatchRecords[batchID]))
	for _, minibatch := range m.MinibatchRecords[batchID] {
		res = append(res, minibatch)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].MinibatchIndex < res[j].MinibatchIndex
	})

	return res, nil
}

func (m *minibatchStore) PutDispersalRequest(ctx context.Context, request *batcher.DispersalRequest) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.DispersalRequests[request.BatchID]; !ok {
		m.DispersalRequests[request.BatchID] = make(map[uint][]*batcher.DispersalRequest)
	}

	if _, ok := m.DispersalRequests[request.BatchID][request.MinibatchIndex]; !ok {
		m.DispersalRequests[request.BatchID][request.MinibatchIndex] = make([]*batcher.DispersalRequest, 0)
	}

	for _, r := range m.DispersalRequests[request.BatchID][request.MinibatchIndex] {
		if r.OperatorID == request.OperatorID {
			// replace existing record
			*r = *request
			return nil
		}
	}

	m.DispersalRequests[request.BatchID][request.MinibatchIndex] = append(m.DispersalRequests[request.BatchID][request.MinibatchIndex], request)

	return nil
}

func (m *minibatchStore) GetDispersalRequest(ctx context.Context, batchID uuid.UUID, minibatchIndex uint, opID core.OperatorID) (*batcher.DispersalRequest, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	requests, err := m.GetMinibatchDispersalRequests(ctx, batchID, minibatchIndex)
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

func (m *minibatchStore) GetMinibatchDispersalRequests(ctx context.Context, batchID uuid.UUID, minibatchIndex uint) ([]*batcher.DispersalRequest, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, ok := m.DispersalRequests[batchID]; !ok {
		return nil, BatchNotFound
	}

	return m.DispersalRequests[batchID][minibatchIndex], nil
}

func (m *minibatchStore) PutDispersalResponse(ctx context.Context, response *batcher.DispersalResponse) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.DispersalResponses[response.BatchID]; !ok {
		m.DispersalResponses[response.BatchID] = make(map[uint][]*batcher.DispersalResponse)
	}

	if _, ok := m.DispersalResponses[response.BatchID][response.MinibatchIndex]; !ok {
		m.DispersalResponses[response.BatchID][response.MinibatchIndex] = make([]*batcher.DispersalResponse, 0)
	}

	for _, r := range m.DispersalResponses[response.BatchID][response.MinibatchIndex] {
		if r.OperatorID == response.OperatorID {
			// replace existing record
			*r = *response
			return nil
		}
	}

	m.DispersalResponses[response.BatchID][response.MinibatchIndex] = append(m.DispersalResponses[response.BatchID][response.MinibatchIndex], response)

	return nil
}

func (m *minibatchStore) GetDispersalResponse(ctx context.Context, batchID uuid.UUID, minibatchIndex uint, opID core.OperatorID) (*batcher.DispersalResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	responses, err := m.GetMinibatchDispersalResponses(ctx, batchID, minibatchIndex)
	if err != nil {
		return nil, err
	}
	for _, r := range responses {
		if r.OperatorID == opID {
			return r, nil
		}
	}
	return nil, nil
}

func (m *minibatchStore) GetMinibatchDispersalResponses(ctx context.Context, batchID uuid.UUID, minibatchIndex uint) ([]*batcher.DispersalResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, ok := m.DispersalResponses[batchID]; !ok {
		return nil, BatchNotFound
	}

	return m.DispersalResponses[batchID][minibatchIndex], nil
}

func (m *minibatchStore) GetLatestFormedBatch(ctx context.Context) (batch *batcher.BatchRecord, minibatches []*batcher.MinibatchRecord, err error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	batches, err := m.GetBatchesByStatus(ctx, batcher.BatchStatusFormed)
	if err != nil {
		return nil, nil, err
	}
	if len(batches) == 0 {
		return nil, nil, nil
	}

	batch = batches[0]
	minibatches, err = m.GetMinibatches(ctx, batches[0].ID)
	if err != nil {
		return nil, nil, err
	}

	return batch, minibatches, nil
}

func (m *minibatchStore) getDispersals(ctx context.Context, batchID uuid.UUID) ([]*batcher.DispersalRequest, []*batcher.DispersalResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, ok := m.DispersalRequests[batchID]; !ok {
		return nil, nil, BatchNotFound
	}

	if _, ok := m.DispersalResponses[batchID]; !ok {
		return nil, nil, BatchNotFound
	}

	requests := make([]*batcher.DispersalRequest, 0)
	for _, reqs := range m.DispersalRequests[batchID] {
		requests = append(requests, reqs...)
	}

	responses := make([]*batcher.DispersalResponse, 0)
	for _, resp := range m.DispersalResponses[batchID] {
		responses = append(responses, resp...)
	}

	return requests, responses, nil
}

func (m *minibatchStore) BatchDispersed(ctx context.Context, batchID uuid.UUID) (bool, error) {
	dispersed := true
	requests, responses, err := m.getDispersals(ctx, batchID)
	if err != nil {
		return false, err
	}

	if len(requests) == 0 || len(responses) == 0 {
		return false, nil
	}

	if len(requests) != len(responses) {
		m.logger.Info("number of minibatch dispersal requests does not match the number of responses", "batchID", batchID, "numRequests", len(requests), "numResponses", len(responses))
		return false, nil
	}

	for _, resp := range responses {
		if resp.RespondedAt.IsZero() {
			dispersed = false
			m.logger.Info("response pending", "batchID", batchID, "minibatchIndex", resp.MinibatchIndex, "operatorID", resp.OperatorID.Hex())
		}
	}

	return dispersed, nil
}
