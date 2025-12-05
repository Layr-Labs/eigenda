package disperser

import (
	"context"
	"fmt"
	"slices"
	"sync"
)

// A simple thread-safe mock implementation of DisperserRegistry.
type MockDisperserRegistry struct {
	lock sync.Mutex

	defaultDispersers  []uint32
	onDemandDispersers []uint32
	disperserGrpcUris  map[uint32]string
}

// Creates a new mock with empty state.
func NewMockDisperserRegistry() *MockDisperserRegistry {
	return &MockDisperserRegistry{
		disperserGrpcUris: make(map[uint32]string),
	}
}

// Configures what GetDefaultDispersers will return.
func (r *MockDisperserRegistry) SetDefaultDispersers(dispersers []uint32) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.defaultDispersers = dispersers
}

// Configures what IsOnDemandDisperser will return.
func (r *MockDisperserRegistry) SetOnDemandDispersers(dispersers []uint32) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.onDemandDispersers = dispersers
}

// Configures what GetDisperserGrpcUri will return for a specific disperser.
func (r *MockDisperserRegistry) SetDisperserGrpcUri(disperserID uint32, uri string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.disperserGrpcUris[disperserID] = uri
}

// Returns the list configured via SetDefaultDispersers.
func (r *MockDisperserRegistry) GetDefaultDispersers(ctx context.Context) ([]uint32, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	result := make([]uint32, len(r.defaultDispersers))
	copy(result, r.defaultDispersers)
	return result, nil
}

// Returns whether the specified disperser is configured as an on-demand disperser via SetOnDemandDispersers.
func (r *MockDisperserRegistry) IsOnDemandDisperser(ctx context.Context, disperserID uint32) (bool, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	return slices.Contains(r.onDemandDispersers, disperserID), nil
}

// Returns the URI configured via SetDisperserGrpcUri for the specified disperser.
func (r *MockDisperserRegistry) GetDisperserGrpcUri(ctx context.Context, disperserID uint32) (string, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	uri, exists := r.disperserGrpcUris[disperserID]
	if !exists {
		return "", fmt.Errorf("no gRPC URI configured for disperser ID %d", disperserID)
	}

	return uri, nil
}
