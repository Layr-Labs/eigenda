package controller

import (
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/stretchr/testify/mock"
)

type MockClientManager struct {
	mock.Mock
}

var _ NodeClientManager = (*MockClientManager)(nil)

func (m *MockClientManager) GetClient(host, port string) (clients.NodeClientV2, error) {
	args := m.Called(host, port)
	client, _ := args.Get(0).(clients.NodeClientV2)
	return client, args.Error(1)
}
