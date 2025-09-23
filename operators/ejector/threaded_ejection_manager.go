package ejector

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigensdk-go/logging"
	geth "github.com/ethereum/go-ethereum/common"
)

// A wrapper around an EjectionManager that handles threading and synchronization.
type ThreadedEjectionManager struct {
	ctx    context.Context
	logger logging.Logger

	// The underlying ejection manager that does the actual work.
	ejectionManager EjectionManager

	// Channel for receiving ejection requests.
	ejectionRequestChan chan geth.Address

	// The period between the background checks for ejection progress.
	period time.Duration
}

// Creates a new ThreadedEjectionManager that wraps the given EjectionManager. Runs a background goroutine
// until the context is cancelled.
func NewThreadedEjectionManager(
	ctx context.Context,
	logger logging.Logger,
	ejectionManager EjectionManager,
	period time.Duration,
) *ThreadedEjectionManager {
	tem := &ThreadedEjectionManager{
		ctx:                 ctx,
		logger:              logger,
		ejectionManager:     ejectionManager,
		ejectionRequestChan: make(chan geth.Address),
		period:              period,
	}
	go tem.mainLoop()
	return tem
}

// EjectValidator begins ejection proceedings for a validator if it is appropriate to do so.
//
// There are several conditions where calling this method will not result in a new ejection being attempted:
//   - There is already an ongoing ejection for the validator.
//   - The validator is in the ejection blacklist (i.e. validators we will never attempt to eject).
//   - A previous attempt at ejecting the validator was made too recently.
func (tem *ThreadedEjectionManager) EjectValidator(validatorAddress geth.Address) error {
	select {
	case <-tem.ctx.Done():
		return fmt.Errorf("context closed: %w", tem.ctx.Err())
	case tem.ejectionRequestChan <- validatorAddress:
		return nil
	}
}

// All modifications to struct state are done in this main loop goroutine.
func (tem *ThreadedEjectionManager) mainLoop() {
	ticker := time.NewTicker(tem.period)
	defer ticker.Stop()

	for {
		select {
		case <-tem.ctx.Done():
			tem.logger.Info("Ejection manager shutting down")
			return
		case request := <-tem.ejectionRequestChan:
			tem.ejectionManager.BeginEjection(request, nil) // TODO
		case <-ticker.C:
			tem.ejectionManager.FinalizeEjections()
		}
	}
}
