package controller

import (
	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/stretchr/testify/mock"
)

type MockClientManager struct {
	mock.Mock
}

var _ NodeClientManager = (*MockClientManager)(nil)

func (m *MockClientManager) GetClient(host, port string) (clients.NodeClient, error) {
	args := m.Called(host, port)
	client, _ := args.Get(0).(clients.NodeClient)
	return client, args.Error(1)
}
