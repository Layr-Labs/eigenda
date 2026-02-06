package store

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	"github.com/Layr-Labs/eigenda/chainstate/types"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/ethereum/go-ethereum/common"
)

// MemoryStore is an in-memory implementation of the Store interface.
// It stores all data in memory and can be persisted to disk via snapshots.
type MemoryStore struct {
	mu sync.RWMutex

	// Map of operator ID to operator
	operators map[core.OperatorID]*types.Operator

	// Map of "quorumID:blockNumber" to quorum APK
	quorumAPKs map[string]*types.QuorumAPK

	// List of all ejections
	ejections []*types.OperatorEjection

	// List of all socket updates
	socketUpdates []*types.OperatorSocketUpdate

	// Last block number that was indexed
	lastIndexedBlock uint64
}

// memoryStoreSnapshot is the serializable representation of the memory store.
// Uses value types instead of pointers for JSON serialization.
// Uses string keys for operators map since OperatorID ([32]byte) can't be a JSON key.
type memoryStoreSnapshot struct {
	Operators        map[string]types.Operator    `json:"operators"`
	QuorumAPKs       map[string]types.QuorumAPK   `json:"quorum_apks"`
	Ejections        []types.OperatorEjection     `json:"ejections"`
	SocketUpdates    []types.OperatorSocketUpdate `json:"socket_updates"`
	LastIndexedBlock uint64                       `json:"last_indexed_block"`
}

// NewMemoryStore creates a new in-memory store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		operators:     make(map[core.OperatorID]*types.Operator),
		quorumAPKs:    make(map[string]*types.QuorumAPK),
		ejections:     make([]*types.OperatorEjection, 0),
		socketUpdates: make([]*types.OperatorSocketUpdate, 0),
	}
}

// SaveOperator implements Store.SaveOperator.
func (s *MemoryStore) SaveOperator(ctx context.Context, op *types.Operator) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Make a copy to avoid external mutations
	opCopy := *op
	s.operators[op.ID] = &opCopy
	return nil
}

// GetOperator implements Store.GetOperator.
func (s *MemoryStore) GetOperator(ctx context.Context, id core.OperatorID) (*types.Operator, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	op, exists := s.operators[id]
	if !exists {
		return nil, fmt.Errorf("operator not found: %x", id)
	}

	// Return a copy to prevent external mutations
	opCopy := *op
	return &opCopy, nil
}

// ListOperators implements Store.ListOperators.
func (s *MemoryStore) ListOperators(ctx context.Context, filter types.OperatorFilter, limit, offset int) ([]*types.Operator, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*types.Operator

	for _, op := range s.operators {
		// Apply filters
		if filter.RegisteredOnly && !op.IsRegistered() {
			continue
		}
		if filter.DeregisteredOnly && op.IsRegistered() {
			continue
		}
		if filter.QuorumID != nil {
			found := false
			for _, qid := range op.QuorumIDs {
				if qid == *filter.QuorumID {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		if filter.MinBlock > 0 && op.RegisteredAtBlockNumber < filter.MinBlock {
			continue
		}
		if filter.MaxBlock > 0 && op.RegisteredAtBlockNumber > filter.MaxBlock {
			continue
		}

		opCopy := *op
		result = append(result, &opCopy)
	}

	// Sort by registration block number for consistent ordering
	sort.Slice(result, func(i, j int) bool {
		return result[i].RegisteredAtBlockNumber < result[j].RegisteredAtBlockNumber
	})

	// Apply pagination
	if offset >= len(result) {
		return []*types.Operator{}, nil
	}
	result = result[offset:]
	if limit > 0 && limit < len(result) {
		result = result[:limit]
	}

	return result, nil
}

// UpdateOperatorSocket implements Store.UpdateOperatorSocket.
func (s *MemoryStore) UpdateOperatorSocket(ctx context.Context, id core.OperatorID, socket string, blockNum uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	op, exists := s.operators[id]
	if !exists {
		return fmt.Errorf("operator not found: %x", id)
	}

	// Since op is a pointer, we can modify it directly
	op.Socket = socket
	return nil
}

// DeregisterOperator implements Store.DeregisterOperator.
func (s *MemoryStore) DeregisterOperator(ctx context.Context, id core.OperatorID, blockNum uint64, txHash common.Hash) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	op, exists := s.operators[id]
	if !exists {
		return fmt.Errorf("operator not found: %x", id)
	}

	// Since op is a pointer, we can modify it directly
	op.DeregisteredAtBlockNumber = &blockNum
	op.DeregisteredTxHash = &txHash
	return nil
}

// SaveQuorumAPK implements Store.SaveQuorumAPK.
func (s *MemoryStore) SaveQuorumAPK(ctx context.Context, apk *types.QuorumAPK) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := fmt.Sprintf("%d:%d", apk.QuorumID, apk.BlockNumber)
	apkCopy := *apk
	s.quorumAPKs[key] = &apkCopy
	return nil
}

// GetQuorumAPK implements Store.GetQuorumAPK.
func (s *MemoryStore) GetQuorumAPK(ctx context.Context, quorumID uint8, blockNum uint64) (*types.QuorumAPK, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := fmt.Sprintf("%d:%d", quorumID, blockNum)
	apk, exists := s.quorumAPKs[key]
	if !exists {
		return nil, fmt.Errorf("quorum APK not found for quorum %d at block %d", quorumID, blockNum)
	}

	apkCopy := *apk
	return &apkCopy, nil
}

// ListQuorumAPKs implements Store.ListQuorumAPKs.
func (s *MemoryStore) ListQuorumAPKs(ctx context.Context, filter types.QuorumAPKFilter) ([]*types.QuorumAPK, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*types.QuorumAPK

	for _, apk := range s.quorumAPKs {
		// Apply filters
		if apk.QuorumID != core.QuorumID(filter.QuorumID) {
			continue
		}
		if filter.BlockNumber > 0 && apk.BlockNumber != filter.BlockNumber {
			continue
		}
		if filter.MinBlock > 0 && apk.BlockNumber < filter.MinBlock {
			continue
		}
		if filter.MaxBlock > 0 && apk.BlockNumber > filter.MaxBlock {
			continue
		}

		apkCopy := *apk
		result = append(result, &apkCopy)
	}

	// Sort by block number
	sort.Slice(result, func(i, j int) bool {
		return result[i].BlockNumber < result[j].BlockNumber
	})

	return result, nil
}

// SaveEjection implements Store.SaveEjection.
func (s *MemoryStore) SaveEjection(ctx context.Context, ejection *types.OperatorEjection) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ejectionCopy := *ejection
	s.ejections = append(s.ejections, &ejectionCopy)
	return nil
}

// ListEjections implements Store.ListEjections.
func (s *MemoryStore) ListEjections(ctx context.Context, operatorID *core.OperatorID, limit, offset int) ([]*types.OperatorEjection, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*types.OperatorEjection

	for _, ej := range s.ejections {
		// Filter by operator ID if specified
		if operatorID != nil && ej.OperatorID != *operatorID {
			continue
		}

		ejCopy := *ej
		result = append(result, &ejCopy)
	}

	// Sort by block number (descending - most recent first)
	sort.Slice(result, func(i, j int) bool {
		return result[i].BlockNumber > result[j].BlockNumber
	})

	// Apply pagination
	if offset >= len(result) {
		return []*types.OperatorEjection{}, nil
	}
	result = result[offset:]
	if limit > 0 && limit < len(result) {
		result = result[:limit]
	}

	return result, nil
}

// SaveSocketUpdate implements Store.SaveSocketUpdate.
func (s *MemoryStore) SaveSocketUpdate(ctx context.Context, update *types.OperatorSocketUpdate) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	updateCopy := *update
	s.socketUpdates = append(s.socketUpdates, &updateCopy)
	return nil
}

// ListSocketUpdates implements Store.ListSocketUpdates.
func (s *MemoryStore) ListSocketUpdates(ctx context.Context, operatorID core.OperatorID, limit, offset int) ([]*types.OperatorSocketUpdate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*types.OperatorSocketUpdate

	for _, update := range s.socketUpdates {
		if update.OperatorID == operatorID {
			updateCopy := *update
			result = append(result, &updateCopy)
		}
	}

	// Sort by block number (descending - most recent first)
	sort.Slice(result, func(i, j int) bool {
		return result[i].BlockNumber > result[j].BlockNumber
	})

	// Apply pagination
	if offset >= len(result) {
		return []*types.OperatorSocketUpdate{}, nil
	}
	result = result[offset:]
	if limit > 0 && limit < len(result) {
		result = result[:limit]
	}

	return result, nil
}

// GetLastIndexedBlock implements Store.GetLastIndexedBlock.
func (s *MemoryStore) GetLastIndexedBlock(ctx context.Context) (uint64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.lastIndexedBlock, nil
}

// SetLastIndexedBlock implements Store.SetLastIndexedBlock.
func (s *MemoryStore) SetLastIndexedBlock(ctx context.Context, blockNum uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lastIndexedBlock = blockNum
	return nil
}

// Snapshot implements Store.Snapshot.
func (s *MemoryStore) Snapshot() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Convert pointer maps/slices to value types for JSON serialization
	// Convert OperatorID keys to hex strings since byte arrays can't be JSON keys
	operators := make(map[string]types.Operator, len(s.operators))
	for id, op := range s.operators {
		operators[id.Hex()] = *op
	}

	quorumAPKs := make(map[string]types.QuorumAPK, len(s.quorumAPKs))
	for key, apk := range s.quorumAPKs {
		quorumAPKs[key] = *apk
	}

	ejections := make([]types.OperatorEjection, len(s.ejections))
	for i, ej := range s.ejections {
		ejections[i] = *ej
	}

	socketUpdates := make([]types.OperatorSocketUpdate, len(s.socketUpdates))
	for i, upd := range s.socketUpdates {
		socketUpdates[i] = *upd
	}

	snapshot := memoryStoreSnapshot{
		Operators:        operators,
		QuorumAPKs:       quorumAPKs,
		Ejections:        ejections,
		SocketUpdates:    socketUpdates,
		LastIndexedBlock: s.lastIndexedBlock,
	}

	return json.Marshal(snapshot)
}

// Restore implements Store.Restore.
func (s *MemoryStore) Restore(data []byte) error {
	var snapshot memoryStoreSnapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return fmt.Errorf("failed to unmarshal snapshot: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Convert value types to pointer maps/slices
	// Convert hex string keys back to OperatorID
	s.operators = make(map[core.OperatorID]*types.Operator, len(snapshot.Operators))
	for idStr, op := range snapshot.Operators {
		idBytes, err := hex.DecodeString(idStr)
		if err != nil {
			return fmt.Errorf("failed to decode operator ID %q: %w", idStr, err)
		}
		if len(idBytes) != 32 {
			return fmt.Errorf("invalid operator ID length %q: expected 32 bytes, got %d", idStr, len(idBytes))
		}
		var id core.OperatorID
		copy(id[:], idBytes)
		opCopy := op
		s.operators[id] = &opCopy
	}

	s.quorumAPKs = make(map[string]*types.QuorumAPK, len(snapshot.QuorumAPKs))
	for key, apk := range snapshot.QuorumAPKs {
		apkCopy := apk
		s.quorumAPKs[key] = &apkCopy
	}

	s.ejections = make([]*types.OperatorEjection, len(snapshot.Ejections))
	for i, ej := range snapshot.Ejections {
		ejCopy := ej
		s.ejections[i] = &ejCopy
	}

	s.socketUpdates = make([]*types.OperatorSocketUpdate, len(snapshot.SocketUpdates))
	for i, upd := range snapshot.SocketUpdates {
		updCopy := upd
		s.socketUpdates[i] = &updCopy
	}

	s.lastIndexedBlock = snapshot.LastIndexedBlock

	return nil
}
