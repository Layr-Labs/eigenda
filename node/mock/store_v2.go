package mock

import (
	"github.com/Layr-Labs/eigenda/common/kvstore"
	"github.com/Layr-Labs/eigenda/node"
	"github.com/stretchr/testify/mock"
)

// MockStoreV2 is a mock implementation of StoreV2
type MockStoreV2 struct {
	mock.Mock
}

var _ node.ValidatorStore = (*MockStoreV2)(nil)

func NewMockStoreV2() *MockStoreV2 {
	return &MockStoreV2{}
}

func (m *MockStoreV2) StoreBatch(batchData []*node.BundleToStore) (uint64, error) {
	args := m.Called(batchData)
	if args.Get(0) == nil {
		return 0, args.Error(1)
	}
	return 0, args.Error(1)
}

func (m *MockStoreV2) DeleteKeys(keys []kvstore.Key) error {
	args := m.Called(keys)
	return args.Error(0)
}

func (m *MockStoreV2) GetBundleData(bundleKey []byte) ([]byte, error) {
	args := m.Called(bundleKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockStoreV2) Stop() error {
	return nil
}
