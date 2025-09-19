package ejector

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigensdk-go/logging"
	geth "github.com/ethereum/go-ethereum/common"
)

// A utility that manages ejections and the ejection lifecycle. An ejection manager is responsible for executing
// ejections, not deciding when it is appropriate to eject. That is to say, this utility does not monitor validator
// signing rates.
type EjectionManager struct {
	ctx        context.Context
	logger     logging.Logger
	timeSource func() time.Time

	// A set of validators that we will not attempt to eject.
	//
	// There are two ways a validator can end up in this blacklist:
	// 1. specified in configuration
	// 2. we've made many attempts to eject the validator, and each attempt has failed (i.e. the validator is
	//    cancelling the ejection on-chain).
	ejectionBlacklist map[geth.Address]struct{}

	// The timestamps of recent ejection attempts, keyed by validator address.
	recentEjections map[geth.Address]time.Time

	// Ejections that have been started but not completed, keyed by validator address. The value is the
	// time the ejection was started.
	startedEjections map[geth.Address]time.Time

	// The number of consecutive failed ejection attempts, keyed by validator address. If this exceeds a
	// threshold, the validator is added to the ejection blacklist.
	failedEjectionAttempts map[geth.Address]uint32

	// Submits ejection transactions.
	transactor EjectionTransactor

	// The minimum time between starting an ejection and completing it.
	ejectionDelay time.Duration

	// The minimum time between two consecutive ejection attempts for the same validator.
	retryDelay time.Duration

	// Channel for receiving ejection requests.
	ejectionRequestChan chan geth.Address

	// The maximum number of consecutive failed ejection attempts before a validator is blacklisted.
	maxConsecutiveFailedEjectionAttempts uint32

	// The period between the background checks for ejection progress.
	period time.Duration
}

// Create a new EjectionManager.
func NewEjectionManager(
	ctx context.Context,
	logger logging.Logger,
	timeSource func() time.Time,
	transactor EjectionTransactor,
// the minimum time between starting an ejection and completing it
	ejectionDelay time.Duration,
// the minimum time between two consecutive ejection attempts for the same validator
	retryDelay time.Duration,
// the maximum number of consecutive failed ejection attempts before a validator is blacklisted
	maxConsecutiveFailedEjectionAttempts uint32,
// the period between the background checks for ejection progress
	period time.Duration,
) (*EjectionManager, error) {

	em := &EjectionManager{
		ctx:                 ctx,
		logger:              logger,
		timeSource:          timeSource,
		ejectionBlacklist:   make(map[geth.Address]struct{}),
		recentEjections:     make(map[geth.Address]time.Time),
		startedEjections:    make(map[geth.Address]time.Time),
		transactor:          transactor,
		ejectionDelay:       ejectionDelay,
		retryDelay:          retryDelay,
		ejectionRequestChan: make(chan geth.Address, 64),
		period:              period,
	}

	go em.mainLoop()

	return em, nil
}

// EjectValidator begins ejection proceedings for a validator if it is appropriate to do so.
//
// There are several conditions where calling this method will not result in a new ejection being attempted:
//   - There is already an ongoing ejection for the validator.
//   - The validator is in the ejection blacklist (i.e. validators we will never attempt to eject).
//   - A previous attempt at ejecting the validator was made too recently.
func (em *EjectionManager) EjectValidator(validatorAddress geth.Address) error {
	select {
	case <-em.ctx.Done():
		return fmt.Errorf("context closed: %w", em.ctx.Err())
	case em.ejectionRequestChan <- validatorAddress:
		return nil
	}
}

// All modifications to struct state are done in this main loop goroutine.
func (em *EjectionManager) mainLoop() {
	ticker := time.NewTicker(em.period)
	defer ticker.Stop()

	for {
		select {
		case <-em.ctx.Done():
			em.logger.Info("Ejection manager shutting down")
			return
		case request := <-em.ejectionRequestChan:
			em.processEjectionRequest(request)
		case <-ticker.C:
			em.cleanRecentEjections()
			em.finalizeEjections()
		}
	}
}

// processEjectionRequest processes a single ejection request.
func (em *EjectionManager) processEjectionRequest(validatorAddress geth.Address) {
	// Check to see if the validator is blacklisted.
	if _, blacklisted := em.ejectionBlacklist[validatorAddress]; blacklisted {
		em.logger.Infof("validator %s is blacklisted from ejection, skipping", validatorAddress.Hex())
		return
	}

	// Check to see if we are already in the process of ejecting this validator.
	if _, inProgress := em.startedEjections[validatorAddress]; inProgress {
		em.logger.Infof("ejection already in progress for validator %s, skipping", validatorAddress.Hex())
		return
	}

	// Check to see if we have recently attempted to eject this validator.
	if _, recentlyEjected := em.recentEjections[validatorAddress]; recentlyEjected {
		em.logger.Infof("recent ejection attempt for validator %s, skipping", validatorAddress.Hex())
		return
	}

	// Check to see if there is already an ejection in progress on-chain for this validator.
	inProgress, err := em.transactor.IsEjectionInProgress(em.ctx, validatorAddress)
	if err != nil {
		em.logger.Errorf("failed to check ejection status for validator %s: %v", validatorAddress.Hex(), err)
		return
	}

	if inProgress {
		// An ejection is already in progress. Record it, and we can try to finalize it later.
		em.logger.Infof("ejection already in progress on-chain for validator %s", validatorAddress.Hex())
	} else {
		// Start a new ejection.
		err = em.transactor.StartEjection(em.ctx, validatorAddress)
		if err != nil {
			em.logger.Errorf("failed to start ejection for validator %s: %v", validatorAddress.Hex(), err)
			return
		}
	}

	em.recentEjections[validatorAddress] = em.timeSource()
	em.startedEjections[validatorAddress] = em.timeSource()
}

// cleanRecentEjections removes entries from recentEjections that are older than the retry delay. We only need
// to remember prior ejections when those ejections prevent us from attempting a new ejection.
func (em *EjectionManager) cleanRecentEjections() {

	// Note: iterating this entire map is not as efficient as a priority queue. However, there are two mitigating
	// factors that make this less than optimal approach acceptable.
	//
	// 1. The total number of validators has a moderately small upper bound (i.e. 2,000). Cheap for an O(n) operation,
	//    and each step is just a map lookup and a time comparison.
	// 2. This method is called infrequently (e.g. every 5 minutes).
	//
	// With this in mind, I have decided to keep the implementation simple for now.

	// Another possible optimization if this code ever becomes a hotspot is to execute eth transactions on
	// background goroutines, so that this loop is not blocked on network calls. Premature at current scale.

	cutoff := em.timeSource().Add(-em.retryDelay)
	for addr, ts := range em.recentEjections {
		if ts.Before(cutoff) {
			delete(em.recentEjections, addr)
		}
	}
}

// For each ejection that that was started long enough ago to be eligible for finalization, check the status
// and finalize if appropriate.
func (em *EjectionManager) finalizeEjections() {

	// Note: similar to cleanRecentEjections(), we are iterating a map here. At a certain scale a
	// priority queue would be more efficient, but that optimization is premature at this time.

	cutoff := em.timeSource().Add(-em.ejectionDelay)

	for address, ejectionStartedTimestamp := range em.startedEjections {
		if ejectionStartedTimestamp.After(cutoff) {
			// Not ready to finalize yet.
			continue
		}
		em.finalizeEjection(address)
	}
}

// Finalize the ejection for a specific validator.
func (em *EjectionManager) finalizeEjection(address geth.Address) {
	// Check to see if the ejection is still in progress.
	inProgress, err := em.transactor.IsEjectionInProgress(em.ctx, address)
	if err != nil {
		em.logger.Errorf("failed to check ejection status for validator %s: %v", address.Hex(), err)
		return
	}

	if !inProgress {
		// Either the validator cancelled the ejection or another ejector finalized it for us.
		em.handleAbortedEjection(address)
		return
	}

	// Complete the ejection.
	err = em.transactor.CompleteEjection(em.ctx, address)
	if err != nil {
		// We failed to eject. Leave the ejection in progress so we can try again later.
		em.logger.Errorf("failed to complete ejection for validator %s: %v", address.Hex(), err)
		return
	}

	em.logger.Infof("successfully completed ejection for validator %s", address.Hex())
	delete(em.startedEjections, address)
}

// Handle the case where a previously started ejection is no longer in progress.
func (em *EjectionManager) handleAbortedEjection(address geth.Address) {
	isPresent, err := em.transactor.IsValidatorPresentInAnyQuorum(em.ctx, address)
	if err != nil {
		em.logger.Errorf("failed to check quorum presence for validator %s: %v", address.Hex(), err)
		return
	}

	if isPresent {
		// The validator cancelled the ejection. Increment the failed attempt counter.
		em.logger.Warnf("ejection for validator %s was cancelled", address.Hex())
		em.failedEjectionAttempts[address]++
		if em.failedEjectionAttempts[address] >= em.maxConsecutiveFailedEjectionAttempts {
			em.logger.Errorf(
				"Validator %s has exceeded maximum consecutive failed ejection attempts, "+
					"adding to blacklist. No further attempts will be made to eject.", address.Hex())
			em.ejectionBlacklist[address] = struct{}{}
			delete(em.failedEjectionAttempts, address)
		} else {
			em.logger.Infof("validator %s has %d consecutive failed ejection attempts",
				address.Hex(), em.failedEjectionAttempts[address])
		}
	} else {
		// A different ejector finalized the ejection for us, or the validator was removed from all quorums by
		// some other mechanism. Either way, we are done here.
		em.logger.Infof("validator %s no longer present in any quorum, ejection complete", address.Hex())
	}

	delete(em.startedEjections, address)
	return
}
