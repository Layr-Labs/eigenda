package geth

import (
	"sync"

	"github.com/Layr-Labs/eigensdk-go/logging"
)

type RPCStatistics struct {
	numberRpcFault uint64
	Logger         logging.Logger
	mu             *sync.RWMutex
}

func NewRPCStatistics(logger logging.Logger) *RPCStatistics {
	return &RPCStatistics{
		Logger: logger,
		mu:     &sync.RWMutex{},
	}
}

// ProcessError attributes the error and updates total number of fault for RPC
// It returns if RPC should immediately give up
func (f *RPCStatistics) ProcessError(err error) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	if err == nil {
		return false
	}

	nextEndpoint, action := f.handleError(err)

	if nextEndpoint == NewRPC {
		f.numberRpcFault += 1
	}

	return action == Return
}

func (f *RPCStatistics) GetTotalNumberRpcFault() uint64 {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.numberRpcFault
}
