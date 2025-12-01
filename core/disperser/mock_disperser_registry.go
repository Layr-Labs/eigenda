package disperser

import (
	"context"
	"fmt"
	"sync"
)

// A simple thread-safe mock implementation of DisperserRegistry.
type MockDisperserRegistry struct {
	lock sync.Mutex

	defaultDispersers    []uint32
	defaultDispersersErr error

	onDemandDispersers    []uint32
	onDemandDispersersErr error

	// Map from disperser ID to gRPC URI
	disperserGrpcUris map[uint32]string
	// Map from disperser ID to error (if we want GetDisperserGrpcUri to fail for specific IDs)
	disperserGrpcUriErrs map[uint32]error
}

// Creates a new mock with empty state.
func NewMockDisperserRegistry() *MockDisperserRegistry {
	return &MockDisperserRegistry{
		disperserGrpcUris:    make(map[uint32]string),
		disperserGrpcUriErrs: make(map[uint32]error),
	}
}

// Configures what GetDefaultDispersers will return.
func (r *MockDisperserRegistry) SetDefaultDispersers(dispersers []uint32, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.defaultDispersers = dispersers
	r.defaultDispersersErr = err
}

// Configures what GetOnDemandDispersers will return.
func (r *MockDisperserRegistry) SetOnDemandDispersers(dispersers []uint32, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.onDemandDispersers = dispersers
	r.onDemandDispersersErr = err
}

// Configures what GetDisperserGrpcUri will return for a specific disperser.
func (r *MockDisperserRegistry) SetDisperserGrpcUri(disperserID uint32, uri string, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if err != nil {
		r.disperserGrpcUriErrs[disperserID] = err
	} else {
		r.disperserGrpcUris[disperserID] = uri
	}
}

// Returns the list configured via SetDefaultDispersers.
func (r *MockDisperserRegistry) GetDefaultDispersers(ctx context.Context) ([]uint32, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.defaultDispersersErr != nil {
		return nil, r.defaultDispersersErr
	}
	// Return a copy to avoid external modifications
	result := make([]uint32, len(r.defaultDispersers))
	copy(result, r.defaultDispersers)
	return result, nil
}

// Returns the list configured via SetOnDemandDispersers.
func (r *MockDisperserRegistry) GetOnDemandDispersers(ctx context.Context) ([]uint32, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.onDemandDispersersErr != nil {
		return nil, r.onDemandDispersersErr
	}
	// Return a copy to avoid external modifications
	result := make([]uint32, len(r.onDemandDispersers))
	copy(result, r.onDemandDispersers)
	return result, nil
}

// Returns the URI configured via SetDisperserGrpcUri for the specified disperser.
func (r *MockDisperserRegistry) GetDisperserGrpcUri(ctx context.Context, disperserID uint32) (string, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if err, exists := r.disperserGrpcUriErrs[disperserID]; exists {
		return "", err
	}

	uri, exists := r.disperserGrpcUris[disperserID]
	if !exists {
		return "", fmt.Errorf("no gRPC URI configured for disperser ID %d", disperserID)
	}

	return uri, nil
}
