package mock

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/stretchr/testify/mock"
)

type SocketStateCacheMock struct {
	mock.Mock

	OperatorSockets core.OperatorSockets
}

var _ core.SocketStateCache = (*SocketStateCacheMock)(nil)

func NewSocketStateCacheMock(operatorSockets core.OperatorSockets) (*SocketStateCacheMock, error) {
	if operatorSockets == nil {
		operatorSockets = make(map[core.OperatorID]*core.OperatorSocket)
	}
	return &SocketStateCacheMock{
		OperatorSockets: operatorSockets,
	}, nil
}

func (s *SocketStateCacheMock) GetOperatorSocket(ctx context.Context, operator core.OperatorID) (string, error) {
	args := s.Called(ctx, operator)
	return args.Get(0).(string), args.Error(1)
}

func (s *SocketStateCacheMock) GetOperatorSockets(ctx context.Context, operators []core.OperatorID) (core.OperatorSockets, error) {
	args := s.Called(ctx, operators)
	if args.Get(0) != nil {
		return args.Get(0).(core.OperatorSockets), args.Error(1)
	}

	// If no mock expectation is set, generate deterministic sockets for each operator
	sockets := make(map[core.OperatorID]*core.OperatorSocket)
	for i, operator := range operators {
		socket := generateSocketFromOperatorID(i, operator)
		sockets[operator] = &socket
	}
	return sockets, nil
}

// generateSocketFromOperatorID creates a deterministic socket based on the operator ID
func generateSocketFromOperatorID(operatorIndex int, id core.OperatorID) core.OperatorSocket {
	host := "0.0.0.0"
	dispersalPort := fmt.Sprintf("3%03v", 2*operatorIndex)
	retrievalPort := fmt.Sprintf("3%03v", 2*operatorIndex+1)
	v2DispersalPort := fmt.Sprintf("3%03v", 2*operatorIndex+2)
	v2RetrievalPort := fmt.Sprintf("3%03v", 2*operatorIndex+3)

	return core.MakeOperatorSocket(host, dispersalPort, retrievalPort, v2DispersalPort, v2RetrievalPort)
}
