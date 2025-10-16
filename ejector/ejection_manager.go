package ejector

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common/enforce"
	"github.com/Layr-Labs/eigenda/common/ratelimit"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
	geth "github.com/ethereum/go-ethereum/common"
)

// TODO(cody.littley) add metrics

// EjectionManager manages and executes validator ejections.
type EjectionManager interface {

	// Begin ejection proceedings against a validator. May not take action if it is not appropriate to do so.
	BeginEjection(
		validatorAddress geth.Address,
		// For each quorum the validator is a member of, the validator's stake in that quorum as a fraction of 1.0.
		stakeFractions map[core.QuorumID]float64,
	)

	// For all eligible ejections that have been started, check their status and finalize if appropriate.
	FinalizeEjections()
}

var _ EjectionManager = (*ejectionManager)(nil)

// Information tracked for each in-progress ejection.
type inProgressEjection struct {
	// The time when the ejection can be finalized.
	ejectionFinalizationTime time.Time
	// For each quorum the validator is a member of, the validator's stakeFraction in that quorum as a fraction of 1.0.
	stakeFraction map[core.QuorumID]float64
}

// A utility that manages ejections and the ejection lifecycle. An ejection manager is responsible for executing
// ejections, not deciding when it is appropriate to eject. That is to say, this utility does not monitor validator
// signing rates.
type ejectionManager struct {
	ctx    context.Context
	logger logging.Logger

	// The configuration for the ejector.
	config *EjectorConfig

	// Provides the wall clock time.
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

	// Ejections that have been started but not completed, keyed by validator address.
	ejectionsInProgress map[geth.Address]*inProgressEjection

	// The number of consecutive failed ejection attempts, keyed by validator address. If this exceeds a
	// threshold, the validator is added to the ejection blacklist. For the purposes of this counter,
	// we only count failed attempts where we started an ejection, but the validator cancelled it on-chain.
	// Golang errors are not counted towards this total.
	failedEjectionAttempts map[geth.Address]uint32

	// Submits ejection transactions.
	transactor EjectionTransactor

	// The rate limiter for ejection transactions, keyed by quorum ID. Limits the fraction of the stake (out of 1.0)
	// that can be ejected per time period. Since a quorum ID is an 8-bit integer (in smart contracts, no less!),
	// it's safe to assume that the map will not grow too large.
	quorumRateLimits map[core.QuorumID]*ratelimit.LeakyBucket
}

// Create a new ejectionManager.
func NewEjectionManager(
	ctx context.Context,
	logger logging.Logger,
	config *EjectorConfig,
	// A source of time.
	timeSource func() time.Time,
	// Submits ejection transactions.
	transactor EjectionTransactor,
) (EjectionManager, error) {

	em := &ejectionManager{
		ctx:                    ctx,
		config:                 config,
		logger:                 logger,
		timeSource:             timeSource,
		ejectionBlacklist:      make(map[geth.Address]struct{}),
		recentEjectionTimes:    make(map[geth.Address]time.Time),
		ejectionsInProgress:    make(map[geth.Address]*inProgressEjection),
		failedEjectionAttempts: make(map[geth.Address]uint32),
		quorumRateLimits:       make(map[core.QuorumID]*ratelimit.LeakyBucket),
		transactor:             transactor,
	}

	for _, addr := range config.DoNotEjectTheseValidators {
		em.ejectionBlacklist[geth.HexToAddress(addr)] = struct{}{}
	}

	// Set up a throttle for quorum 0. We will always have a quorum 0, and this allows us to check to see
	// if the throttle config is valid. Checking here lets us assume it is valid later on.
	var err error
	em.quorumRateLimits[0], err = ratelimit.NewLeakyBucket(
		config.EjectionThrottle,
		config.EjectionThrottleTimePeriod,
		config.StartEjectionThrottleFull,
		ratelimit.OverfillOncePermitted,
		timeSource())
	if err != nil {
		return nil, fmt.Errorf("failed to create leaky bucket: %w", err)
	}

	return em, nil
}

func (em *ejectionManager) BeginEjection(
	validatorAddress geth.Address,
	stakeFractions map[core.QuorumID]float64,
) {

	// Sanity check stake fractions.
	if !em.areStakeFractionsValid(validatorAddress, stakeFractions) {
		return
	}

	// Check to see if the validator is blacklisted.
	if _, blacklisted := em.ejectionBlacklist[validatorAddress]; blacklisted {
		em.logger.Infof("validator %s is blacklisted from ejection, will not begin ejection",
			validatorAddress.Hex())
		return
	}

	// Check to see if we are already in the process of ejecting this validator.
	if _, ejectionAlreadyBeingTracked := em.ejectionsInProgress[validatorAddress]; ejectionAlreadyBeingTracked {
		em.logger.Infof("ejection already in progress for validator %s, will not begin ejection",
			validatorAddress.Hex())
		return
	}

	// Check to see if we have recently attempted to eject this validator.
	if _, recentlyEjected := em.recentEjectionTimes[validatorAddress]; recentlyEjected {
		em.logger.Infof("recent ejection attempt for validator %s, will not begin ejection",
			validatorAddress.Hex())
		return
	}

	// Check to see if there is already an ejection in progress on-chain for this validator.
	ejectionStartedOnchain, err := em.transactor.IsEjectionInProgress(em.ctx, validatorAddress)
	if err != nil {
		em.logger.Errorf("failed to check ejection status for validator %s, will not begin ejection: %v",
			validatorAddress.Hex(), err)
		return
	}
	if ejectionStartedOnchain {
		// An ejection is already in progress onchain. Record it, and we can try to finalize it later.
		em.logger.Infof("ejection already in progress on-chain for validator %s, "+
			"will not begin ejection but will attempt to finalize",
			validatorAddress.Hex())

		em.scheduleFutureEjectionFinalization(validatorAddress, stakeFractions)
		return
	}

	// Check if we are prevented from starting an ejection by rate limiting.
	allowedByRateLimits := em.checkRateLimits(validatorAddress, stakeFractions)
	if !allowedByRateLimits {
		// Rate limiting prevents us from starting an ejection at this time.
		// checkRateLimits() will have logged the reason, since it has more context.
		return
	}

	// Start a new ejection.
	err = em.transactor.StartEjection(em.ctx, validatorAddress)
	if err != nil {
		em.logger.Errorf("failed to start ejection for validator %s: %v", validatorAddress.Hex(), err)
		em.cleanUpFailedEjection(validatorAddress, stakeFractions)
		return
	}
	em.logger.Infof("started ejection proceedings against %s", validatorAddress.Hex())

	em.scheduleFutureEjectionFinalization(validatorAddress, stakeFractions)
}

// Mark that an ejection has been started and must be finished in the future.
func (em *ejectionManager) scheduleFutureEjectionFinalization(
	validatorAddress geth.Address,
	stakeFractions map[core.QuorumID]float64,
) {
	em.recentEjectionTimes[validatorAddress] = em.timeSource()
	em.ejectionsInProgress[validatorAddress] = &inProgressEjection{
		ejectionFinalizationTime: em.timeSource().Add(em.config.EjectionFinalizationDelay),
		stakeFraction:            stakeFractions,
	}
}

// Check that the stake fractions are all valid (i.e. in the range (0.0, 1.0]), returning true if they are valid,
// and false otherwise.
func (em *ejectionManager) areStakeFractionsValid(
	validatorAddress geth.Address,
	stakeFractions map[core.QuorumID]float64,
) bool {
	for qid, stake := range stakeFractions {
		if stake <= 0.0 {
			em.logger.Errorf(
				"validator %s has non-positive stake %.4f in quorum %d, will not begin ejection",
				validatorAddress.Hex(), stake, qid)
			return false
		}
		if stake > 1.0 {
			em.logger.Errorf(
				"validator %s has stake %.4f > 1.0 in quorum %d, will not begin ejection",
				validatorAddress.Hex(), stake, qid)
			return false
		}
	}
	return true
}

func (em *ejectionManager) FinalizeEjections() {
	em.cleanRecentEjections()

	// Note: similar to cleanRecentEjections(), we are iterating a map here. At a certain scale a
	// priority queue would be more efficient, but that optimization is premature at this time.

	now := em.timeSource()

	for address, ejection := range em.ejectionsInProgress {
		if now.After(ejection.ejectionFinalizationTime) {
			ejected := em.finalizeEjection(address)

			if !ejected {
				em.cleanUpFailedEjection(address, ejection.stakeFraction)
			}
		}
	}
}

// Check if we are prevented from starting an ejection by rate limiting. If we are prevented from starting
// an ejection in any quorum, we revert all fills and return false. If we are permitted to start an ejection
// in all quorums, we return true and debit the leaky buckets for each quorum.
func (em *ejectionManager) checkRateLimits(
	validatorAddress geth.Address,
	stakeFractions map[core.QuorumID]float64,
) bool {

	now := em.timeSource()
	permittedQuorums := make([]core.QuorumID, 0, len(stakeFractions))
	for qid, stake := range stakeFractions {
		leakyBucket := em.getLeakyBucketForQuorum(qid)

		allowed, err := leakyBucket.Fill(now, stake)

		// The only way we can get an error here is if time moves backwards.
		enforce.NilError(err, "should be impossible")

		if !allowed {
			// We are prevented by rate limiting from starting an ejection in this quorum.
			// We will need to undo all previous fills before bailing out.
			for _, quorumID := range permittedQuorums {
				stakeToUndo := stakeFractions[quorumID]
				leakyBucketToUndo := em.getLeakyBucketForQuorum(quorumID)
				err = leakyBucketToUndo.RevertFill(now, stakeToUndo)
				enforce.NilError(err, "should be impossible")
			}

			em.logger.Warnf("rate limit prevents ejection of validator %s in quorum %d, skipping",
				validatorAddress.Hex(), qid)
			return false
		}
		permittedQuorums = append(permittedQuorums, qid)
	}

	return true
}

// Refund the rate limit fills for each quorum. This should be called if we fail to finalize an ejection.
// Also removes the ejection from ejectionsInProgress.
func (em *ejectionManager) cleanUpFailedEjection(
	validatorAddress geth.Address,
	stakeFractions map[core.QuorumID]float64,
) {
	now := em.timeSource()
	for qid, stake := range stakeFractions {
		leakyBucket := em.getLeakyBucketForQuorum(qid)
		err := leakyBucket.RevertFill(now, stake)
		enforce.NilError(err, "should be impossible")
	}

	delete(em.ejectionsInProgress, validatorAddress)
}

// Get the leaky bucket for a specific quorum, creating it if it doesn't already exist.
func (em *ejectionManager) getLeakyBucketForQuorum(qid core.QuorumID) *ratelimit.LeakyBucket {
	leakyBucket, ok := em.quorumRateLimits[qid]

	if !ok {
		var err error
		leakyBucket, err = ratelimit.NewLeakyBucket(
			em.config.EjectionThrottle,
			em.config.EjectionThrottleTimePeriod,
			em.config.StartEjectionThrottleFull,
			ratelimit.OverfillOncePermitted,
			em.timeSource())
		em.quorumRateLimits[qid] = leakyBucket

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

	cutoff := em.timeSource().Add(-em.config.EjectionRetryDelay)
	for addr, ts := range em.recentEjectionTimes {
		if ts.Before(cutoff) {
			delete(em.recentEjectionTimes, addr)
		}
	}
}

// Finalize the ejection for a specific validator. Returns true if the ejection was finalized, false otherwise.
func (em *ejectionManager) finalizeEjection(address geth.Address) bool {
	// Check to see if the ejection is still in progress.
	inProgress, err := em.transactor.IsEjectionInProgress(em.ctx, address)
	if err != nil {
		em.logger.Errorf("failed to check ejection status for validator %s, will not finalize ejection: %v",
			address.Hex(), err)
		return false
	}

	if !inProgress {
		// Either the validator cancelled the ejection or another ejector finalized it for us.
		em.handleAbortedEjection(address)
		return false
	}

	// Complete the ejection.
	err = em.transactor.CompleteEjection(em.ctx, address)
	if err != nil {
		// We failed to eject, give up for now.
		em.logger.Errorf("failed to complete ejection for validator %s: %v", address.Hex(), err)
		return false
	}

	em.logger.Infof("successfully completed ejection for validator %s", address.Hex())
	// If we return before we get here, it's the responsibility of the caller to refund the rate limits
	// and remove the in-progress ejection.
	delete(em.ejectionsInProgress, address)
	delete(em.failedEjectionAttempts, address)

	return true
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
		if em.failedEjectionAttempts[address] >= em.config.MaxConsecutiveFailedEjectionAttempts {
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
}
