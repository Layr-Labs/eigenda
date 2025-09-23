package ejector

import (
	"context"
	"time"

	"github.com/Layr-Labs/eigenda/common/enforce"
	"github.com/Layr-Labs/eigenda/common/ratelimit"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
	geth "github.com/ethereum/go-ethereum/common"
)

// TODO add metrics

// EjectionManager manages and executes validator ejections.
type EjectionManager interface {

	// Begin ejection proceedings against a validator. May not take action if it is not appropriate to do so.
	BeginEjection(
		validatorAddress geth.Address,
	// For each quorum the validator is a member of, the validator's stake in that quorum as a fraction of 1.0.
		stakes map[core.QuorumID]float64,
	)

	// For all eligible ejections that have been started, check their status and finalize if appropriate.
	FinalizeEjections()
}

var _ EjectionManager = (*ejectionManager)(nil)

// Information tracked for each in-progress ejection.
type inProgressEjection struct {
	// The time the ejection was started.
	ejectionStartTime time.Time
	// For each quorum the validator is a member of, the validator's stake in that quorum as a fraction of 1.0.
	stake map[core.QuorumID]float64
}

// A utility that manages ejections and the ejection lifecycle. An ejection manager is responsible for executing
// ejections, not deciding when it is appropriate to eject. That is to say, this utility does not monitor validator
// signing rates.
type ejectionManager struct {
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
	recentEjectionTimes map[geth.Address]time.Time

	// Ejections that have been started but not completed, keyed by validator address. The value is the
	// time the ejection was started.
	ejectionsInProgress map[geth.Address]*inProgressEjection

	// The number of consecutive failed ejection attempts, keyed by validator address. If this exceeds a
	// threshold, the validator is added to the ejection blacklist. For the purposes of this counter,
	// we only count failed attempts where we started an ejection, but the validator cancelled it on-chain.
	// Golang errors are not counted towards this total.
	failedEjectionAttempts map[geth.Address]uint32

	// Submits ejection transactions.
	transactor EjectionTransactor

	// The minimum time between starting an ejection and completing it.
	ejectionDelay time.Duration

	// The minimum time between two consecutive ejection attempts for the same validator.
	retryDelay time.Duration

	// The maximum number of consecutive failed ejection attempts before a validator is blacklisted.
	maxConsecutiveFailedEjectionAttempts uint32

	// The rate limiter for ejection transactions, keyed by quorum ID. Limits the fraction of the stake (out of 1.0)
	// that can be ejected per time period.
	quorumRateLimits map[core.QuorumID]*ratelimit.LeakyBucket
}

// Create a new ejectionManager.
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
) EjectionManager {

	// TODO set up rate limits

	return &ejectionManager{
		ctx:                                  ctx,
		logger:                               logger,
		timeSource:                           timeSource,
		ejectionBlacklist:                    make(map[geth.Address]struct{}),
		recentEjectionTimes:                  make(map[geth.Address]time.Time),
		ejectionsInProgress:                  make(map[geth.Address]*inProgressEjection),
		failedEjectionAttempts:               make(map[geth.Address]uint32),
		transactor:                           transactor,
		ejectionDelay:                        ejectionDelay,
		retryDelay:                           retryDelay,
		maxConsecutiveFailedEjectionAttempts: maxConsecutiveFailedEjectionAttempts,
	}
}

func (em *ejectionManager) BeginEjection(
	validatorAddress geth.Address,
	stakes map[core.QuorumID]float64,
) {

	// Check to see if the validator is blacklisted.
	if _, blacklisted := em.ejectionBlacklist[validatorAddress]; blacklisted {
		em.logger.Infof("validator %s is blacklisted from ejection, skipping", validatorAddress.Hex())
		return
	}

	// Check to see if we are already in the process of ejecting this validator.
	if _, inProgress := em.ejectionsInProgress[validatorAddress]; inProgress {
		em.logger.Infof("ejection already in progress for validator %s, skipping", validatorAddress.Hex())
		return
	}

	// Check to see if we have recently attempted to eject this validator.
	if _, recentlyEjected := em.recentEjectionTimes[validatorAddress]; recentlyEjected {
		em.logger.Infof("recent ejection attempt for validator %s, skipping", validatorAddress.Hex())
		return
	}

	// Check to see if there is already an ejection in progress on-chain for this validator.
	inProgress, err := em.transactor.IsEjectionInProgress(em.ctx, validatorAddress)
	if err != nil {
		em.logger.Errorf("failed to check ejection status for validator %s: %v", validatorAddress.Hex(), err)
		return
	}

	// Check if we are prevented from starting an ejection by rate limiting.
	// TODO
	//now := em.timeSource()
	//permittedQuorums := make([]core.QuorumID, 0, len(stakes))
	//for qid, stake := range stakes {
	//	leakyBucket := em.getLeakyBucketForQuorum(qid)
	//
	//	err := leakyBucket.Fill(now, )
	//
	//}

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

	em.recentEjectionTimes[validatorAddress] = em.timeSource()
	em.ejectionsInProgress[validatorAddress] = &inProgressEjection{
		ejectionStartTime: em.timeSource(),
		stake:             stakes,
	}
}

func (em *ejectionManager) FinalizeEjections() {

	em.cleanRecentEjections()

	// Note: similar to cleanRecentEjections(), we are iterating a map here. At a certain scale a
	// priority queue would be more efficient, but that optimization is premature at this time.

	cutoff := em.timeSource().Add(-em.ejectionDelay)

	for address, ejection := range em.ejectionsInProgress {
		if ejection.ejectionStartTime.After(cutoff) {
			// Not ready to finalize yet.
			continue
		}
		em.finalizeEjection(address)
	}
}

// Get the leaky bucket for a specific quorum, creating it if it doesn't already exist.
func (em *ejectionManager) getLeakyBucketForQuorum(qid core.QuorumID) *ratelimit.LeakyBucket {
	leakyBucket, ok := em.quorumRateLimits[qid]

	if !ok {
		var err error
		leakyBucket, err = ratelimit.NewLeakyBucket(
			100,
			time.Minute,
			false,
			ratelimit.OverfillNotPermitted, em.timeSource()) // TODO proper config
		em.quorumRateLimits[qid] = leakyBucket

		// TODO check params in constructor
		enforce.NilError(err, "should be impossible, leaky bucket parameters are pre-validated")
	}

	return leakyBucket
}

// cleanRecentEjections removes entries from recentEjectionTimes that are older than the retry delay. We only need
// to remember prior ejections when those ejections prevent us from attempting a new ejection.
func (em *ejectionManager) cleanRecentEjections() {

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
	for addr, ts := range em.recentEjectionTimes {
		if ts.Before(cutoff) {
			delete(em.recentEjectionTimes, addr)
		}
	}
}

// Finalize the ejection for a specific validator.
func (em *ejectionManager) finalizeEjection(address geth.Address) {
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
	delete(em.ejectionsInProgress, address)
	delete(em.failedEjectionAttempts, address)
}

// Handle the case where a previously started ejection is no longer in progress.
func (em *ejectionManager) handleAbortedEjection(address geth.Address) {
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

	delete(em.ejectionsInProgress, address)
	return
}
