package inmem

import (
	"context"
	"fmt"
	"sync"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser/batcher"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/google/uuid"
)

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
		return nil, fmt.Errorf("batch not found")
	}
	return b, nil
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
		return nil, nil
	}
	return m.MinibatchRecords[batchID][minibatchIndex], nil
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

	requests, err := m.GetDispersalRequests(ctx, batchID, minibatchIndex)
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

func (m *minibatchStore) GetDispersalRequests(ctx context.Context, batchID uuid.UUID, minibatchIndex uint) ([]*batcher.DispersalRequest, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, ok := m.DispersalRequests[batchID]; !ok {
		return nil, nil
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

	responses, err := m.GetDispersalResponses(ctx, batchID, minibatchIndex)
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

func (m *minibatchStore) GetDispersalResponses(ctx context.Context, batchID uuid.UUID, minibatchIndex uint) ([]*batcher.DispersalResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, ok := m.DispersalResponses[batchID]; !ok {
		return nil, nil
	}

	return m.DispersalResponses[batchID][minibatchIndex], nil
}

func (m *minibatchStore) GetPendingBatch(ctx context.Context) (*batcher.BatchRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return nil, nil
}
