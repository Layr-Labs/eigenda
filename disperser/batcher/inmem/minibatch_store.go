package inmem

import (
	"fmt"

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
	DispersalRequests map[uuid.UUID]map[uint]*batcher.DispersalRequest
	// DispersalResponses maps batch IDs to a map from minibatch indices to dispersal responses
	DispersalResponses map[uuid.UUID]map[uint]*batcher.DispersalResponse

	logger logging.Logger
}

var _ batcher.MinibatchStore = (*minibatchStore)(nil)

func NewMinibatchStore(logger logging.Logger) batcher.MinibatchStore {
	return &minibatchStore{
		BatchRecords:       make(map[uuid.UUID]*batcher.BatchRecord),
		MinibatchRecords:   make(map[uuid.UUID]map[uint]*batcher.MinibatchRecord),
		DispersalRequests:  make(map[uuid.UUID]map[uint]*batcher.DispersalRequest),
		DispersalResponses: make(map[uuid.UUID]map[uint]*batcher.DispersalResponse),

		logger: logger,
	}
}

func (m *minibatchStore) PutBatch(batch *batcher.BatchRecord) error {
	m.BatchRecords[batch.ID] = batch

	return nil
}

func (m *minibatchStore) GetBatch(batchID uuid.UUID) (*batcher.BatchRecord, error) {
	b, ok := m.BatchRecords[batchID]
	if !ok {
		return nil, fmt.Errorf("batch not found")
	}
	return b, nil
}

func (m *minibatchStore) PutMiniBatch(minibatch *batcher.MinibatchRecord) error {
	if _, ok := m.MinibatchRecords[minibatch.BatchID]; !ok {
		m.MinibatchRecords[minibatch.BatchID] = make(map[uint]*batcher.MinibatchRecord)
	}
	m.MinibatchRecords[minibatch.BatchID][minibatch.MinibatchIndex] = minibatch

	return nil
}

func (m *minibatchStore) GetMiniBatch(batchID uuid.UUID, minibatchIndex uint) (*batcher.MinibatchRecord, error) {
	if _, ok := m.MinibatchRecords[batchID]; !ok {
		return nil, nil
	}
	return m.MinibatchRecords[batchID][minibatchIndex], nil
}

func (m *minibatchStore) PutDispersalRequest(request *batcher.DispersalRequest) error {
	if _, ok := m.DispersalRequests[request.BatchID]; !ok {
		m.DispersalRequests[request.BatchID] = make(map[uint]*batcher.DispersalRequest)
	}
	m.DispersalRequests[request.BatchID][request.MinibatchIndex] = request

	return nil
}

func (m *minibatchStore) GetDispersalRequest(batchID uuid.UUID, minibatchIndex uint) (*batcher.DispersalRequest, error) {
	if _, ok := m.DispersalRequests[batchID]; !ok {
		return nil, nil
	}

	return m.DispersalRequests[batchID][minibatchIndex], nil
}

func (m *minibatchStore) PutDispersalResponse(response *batcher.DispersalResponse) error {
	if _, ok := m.DispersalResponses[response.BatchID]; !ok {
		m.DispersalResponses[response.BatchID] = make(map[uint]*batcher.DispersalResponse)
	}
	m.DispersalResponses[response.BatchID][response.MinibatchIndex] = response

	return nil
}

func (m *minibatchStore) GetDispersalResponse(batchID uuid.UUID, minibatchIndex uint) (*batcher.DispersalResponse, error) {
	if _, ok := m.DispersalResponses[batchID]; !ok {
		return nil, nil
	}

	return m.DispersalResponses[batchID][minibatchIndex], nil
}

func (m *minibatchStore) GetPendingBatch() (*batcher.BatchRecord, error) {
	return nil, nil
}
