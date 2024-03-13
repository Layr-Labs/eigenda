package geth

import (
	"sync"

	"github.com/Layr-Labs/eigensdk-go/logging"
)

type FailoverController struct {
	NumberFault     uint64
	NumberSuccess   uint64
	SwitchTrigger   int
	NumberOfBackups int
	currentRPCIndex int
	Logger          logging.Logger
	mu              *sync.Mutex
}

func NewFailoverController(numBackup int, switchTrigger int, logger logging.Logger) *FailoverController {
	return &FailoverController{
		NumberFault:     0,
		NumberSuccess:   0,
		currentRPCIndex: 0,
		SwitchTrigger:   switchTrigger,
		NumberOfBackups: numBackup,
		Logger:          logger,
		mu:              &sync.Mutex{},
	}
}

// To use the Failover controller, one must insert this function
// after every call that uses RPC.
// This function attribute the error and update statistics
func (f *FailoverController) ProcessError(err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if err == nil {
		f.NumberSuccess += 1
		return
	}

	fault := HandleError(err)
	if fault == SenderFault {
		return
	} else if fault == RPCFault {
		f.updateRPCFault(err)
		return
	} else if fault == Ok {
		return
	} else { // TooManyRequest
		f.updateRPCFault(err)
		return
	}
}

// update rpc fault
func (f *FailoverController) updateRPCFault(err error) {
	f.NumberFault += 1
	f.Logger.Error("RPC fault", "error", err)
}

// return two values
// boolean indicates if it is primary
// integer
func (f *FailoverController) GetClientIndex() (bool, uint64) {
	f.mu.Lock()
	defer f.mu.Unlock()
	index := (f.NumberFault / uint64(f.SwitchTrigger)) % uint64(f.NumberOfBackups+1)
	if uint64(f.currentRPCIndex) != index {
		f.Logger.Info("Change RPC index", "index", index)
		f.currentRPCIndex = int(index)
	}

	if index == 0 {
		return true, 0
	} else {
		return false, index
	}
}
