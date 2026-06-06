package controller

import (
	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/core"
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

func (m *MockClientManager) GetClientForValidator(
	host, port string,
	validatorID core.OperatorID,
) (clients.NodeClient, error) {
	args := m.Called(host, port, validatorID)
	client, _ := args.Get(0).(clients.NodeClient)
	//nolint:wrapcheck // Mock method returns error as-is
	return client, args.Error(1)
}
