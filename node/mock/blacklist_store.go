package mock

import (
	"context"

	pb "github.com/Layr-Labs/eigenda/api/grpc/validator"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/node"
	"github.com/stretchr/testify/mock"
)

// MockBlacklistStore is a mock implementation of BlacklistStore
type MockBlacklistStore struct {
	mock.Mock
	time node.Time
}

var _ node.BlacklistStore = (*MockBlacklistStore)(nil)

func NewMockBlacklistStore(time node.Time) *MockBlacklistStore {
	if time == nil {
		time = node.DefaultTime
	}
	return &MockBlacklistStore{
		time: time,
	}
}

// SetTime allows updating the time implementation for testing
func (m *MockBlacklistStore) SetTime(time node.Time) {
	m.time = time
}

// GetTime returns the current time implementation
func (m *MockBlacklistStore) GetTime() node.Time {
	return m.time
}

func (m *MockBlacklistStore) HasDisperserID(ctx context.Context, disperserId uint32) bool {
	args := m.Called(ctx, disperserId)
	return args.Bool(0)
}

func (m *MockBlacklistStore) HasKey(ctx context.Context, key []byte) bool {
	args := m.Called(ctx, key)
	return args.Bool(0)
}

func (m *MockBlacklistStore) GetByDisperserID(ctx context.Context, disperserId uint32) (*node.Blacklist, error) {
	args := m.Called(ctx, disperserId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*node.Blacklist), args.Error(1)
}

func (m *MockBlacklistStore) Get(ctx context.Context, key []byte) (*node.Blacklist, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*node.Blacklist), args.Error(1)
}

func (m *MockBlacklistStore) Put(ctx context.Context, key []byte, value []byte) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (m *MockBlacklistStore) AddEntry(ctx context.Context, disperserId uint32, contextId, reason string) error {
	args := m.Called(ctx, disperserId, contextId, reason)
	return args.Error(0)
}

func (m *MockBlacklistStore) IsBlacklisted(ctx context.Context, disperserId uint32) bool {
	args := m.Called(ctx, disperserId)
	return args.Bool(0)
}

func (m *MockBlacklistStore) BlacklistDisperserFromBlobCert(request *pb.StoreChunksRequest, blobCert *corev2.BlobCertificate) error {
	args := m.Called(request, blobCert)
	return args.Error(0)
}
