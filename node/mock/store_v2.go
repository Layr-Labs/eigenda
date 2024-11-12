package mock

import (
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/node"
	"github.com/stretchr/testify/mock"
)

// MockStoreV2 is a mock implementation ofStoreV2
type MockStoreV2 struct {
	mock.Mock
}

var _ node.StoreV2 = (*MockStoreV2)(nil)

func NewMockStoreV2() *MockStoreV2 {
	return &MockStoreV2{}
}

func (m *MockStoreV2) StoreBatch(batch *corev2.Batch, rawBundles []*node.RawBundles) ([]kvstore.Key, error) {
	args := m.Called(batch, rawBundles)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]kvstore.Key), args.Error(1)
}

func (m *MockStoreV2) DeleteKeys(keys []kvstore.Key) error {
	args := m.Called(keys)
	return args.Error(0)
}

func (m *MockStoreV2) GetChunks(blobKey corev2.BlobKey, quorum core.QuorumID) ([][]byte, error) {
	args := m.Called(blobKey, quorum)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([][]byte), args.Error(1)
}
