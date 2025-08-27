package metadata

import "sync/atomic"

var _ BatchMetadataManager = (*MockBatchMetadataManager)(nil)

// mockBatchMetadataManager is a mock implementation of the BatchMetadataManager interface.
type MockBatchMetadataManager struct {
	// The metadata to return when GetMetadata is called.
	metadata atomic.Pointer[BatchMetadata]
}

// Create a mock BatchMetadataManager that returns canned data. The metadata provided to the constructor will
// be returned by GetBlobMetadata, unless SetMetadata is called to change it.
func NewMockBatchMetadataManager(metadata *BatchMetadata) *MockBatchMetadataManager {
	m := &MockBatchMetadataManager{}
	m.metadata.Store(metadata)
	return m
}

func (m *MockBatchMetadataManager) GetMetadata() *BatchMetadata {
	return m.metadata.Load()
}

// SetMetadata sets the metadata to be returned by GetMetadata.
func (m *MockBatchMetadataManager) SetMetadata(metadata *BatchMetadata) {
	m.metadata.Store(metadata)
}

func (m *MockBatchMetadataManager) Close() {
	// intentional no-op
}
