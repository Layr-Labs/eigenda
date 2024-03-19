package geth

import (
	"sync"

	"github.com/Layr-Labs/eigensdk-go/logging"
)

type RPCStatistics struct {
	numberRpcFault uint64
	Logger         logging.Logger
	mu             *sync.Mutex
}

func NewRPCStatistics(logger logging.Logger) *RPCStatistics {
	return &RPCStatistics{
		Logger: logger,
		mu:     &sync.Mutex{},
	}
}

// ProcessError attributes the error and updates total number of fault for RPC
func (f *RPCStatistics) ProcessError(err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if err == nil {
		return
	}

	serverFault := f.handleError(err)

	if serverFault {
		f.numberRpcFault += 1
	}
}

func (f *RPCStatistics) GetTotalNumberRpcFault() uint64 {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.numberRpcFault
}
