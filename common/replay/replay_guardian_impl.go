package replay

import (
	"fmt"
	"sync"
	"time"

	"github.com/emirpasic/gods/queues/priorityqueue"
)

var _ ReplayGuardian = &replayGuardian{}

// replayGuardian is an implementation of ReplayGuardian.
type replayGuardian struct {
	// The time source. In production use cases, this is likely to just be time.Now.
	timeSource func() time.Time

	// The maximum amount of time that a request's timestamp can be ahead of the local wall clock time.
	maxTimeInFuture time.Duration

	// The maximum amount of time that a request's timestamp can be behind the local wall clock time.
	maxTimeInPast time.Duration

	// A set of hashes that have been observed within the time window.
	observedHashes map[string]struct{}

	// A queue of observed hashes, ordered by request timestamp. Used to prune old hashes.
	expirationQueue *priorityqueue.Queue

	// A mutex to protect the observedHashes and expirationQueue.
	lock sync.Mutex
}

// hashWithTimestamp is a request hash with self-reported timestamp associated with that request.
type hashWithTimestamp struct {
	hash      string
	timestamp time.Time
}

// NewReplayGuardian creates a new ReplayGuardian. This implementation is thread safe.
func NewReplayGuardian(
	timeSource func() time.Time,
	// The maximum amount of time that a request's timestamp can be behind the local wall clock time.
	// Increasing this value permits more leniency in the timestamp of incoming requests, at the potential cost
	// of a higher memory overhead.
	maxTimeInPast time.Duration,
	// The maximum amount of time that a request's timestamp can be ahead of the local wall clock time.
	// Increasing this value permits more leniency in the timestamp of incoming requests, at the potential cost of a
	// higher memory overhead. In theory, if requests are sent with a timestamp exactly at the maximum time in the
	// future, this utility will remember them for a total of (maxTimeInFuture + maxTimeInPast), since that is the
	// amount of time that will need to elapse locally before the request exceeds the maximum age. If maxTimeInFuture
	// is extremely large, then an attacker may be able to cause this utility to be forced to remember a very large
	// amount of data.
	maxTimeInFuture time.Duration,
) (ReplayGuardian, error) {

	if timeSource == nil {
		return nil, fmt.Errorf("timeSource cannot be nil")
	}
	if maxTimeInPast < 0 {
		return nil, fmt.Errorf("maxTimeInPast must not be negative, got %v", maxTimeInPast)
	}
	if maxTimeInFuture < 0 {
		return nil, fmt.Errorf("maxTimeInFuture must not be negative, got %v", maxTimeInFuture)
	}

	return &replayGuardian{
		timeSource:      timeSource,
		maxTimeInFuture: maxTimeInFuture,
		maxTimeInPast:   maxTimeInPast,
		observedHashes:  make(map[string]struct{}),
		expirationQueue: priorityqueue.NewWith(compareHashWithTimestamp),
	}, nil
}

// compareKeyWithExpiration compares two hashWithTimestamp objects by their expiration time. Used to create
// a priority queue that orders the requests in chronological order (i.e. the order in which they will expire).
func compareHashWithTimestamp(a interface{}, b interface{}) int {

	keyA := a.(*hashWithTimestamp)
	keyB := b.(*hashWithTimestamp)

	if keyA.timestamp.Before(keyB.timestamp) {
		return -1
	} else if keyA.timestamp.After(keyB.timestamp) {
		return 1
	}
	return 0
}

func (r *replayGuardian) DetailedVerifyRequest(
	requestHash []byte,
	requestTimestamp time.Time,
) ReplayGuardianStatus {
	r.lock.Lock()
	defer r.lock.Unlock()

	now := r.timeSource()

	// Do maintenance on the observedHashes set and expirationQueue.
	r.pruneObservedHashes(now)

	// Reject requests that fall outside the time window we are tracking.
	status := r.verifyTimestamp(now, requestTimestamp)
	if status != StatusValid {
		return status
	}

	// If we've reached this point, then the request will still be in the observedHashes set if it is a replay.
	if _, ok := r.observedHashes[string(requestHash)]; ok {
		return StatusDuplicate
	}

	// The request is not a replay. Add the hash to the observedHashes set and the expirationQueue.
	r.observedHashes[string(requestHash)] = struct{}{}
	r.expirationQueue.Enqueue(&hashWithTimestamp{
		hash:      string(requestHash),
		timestamp: requestTimestamp,
	})

	return StatusValid
}

func (r *replayGuardian) VerifyRequest(requestHash []byte, requestTimestamp time.Time) error {
	status := r.DetailedVerifyRequest(requestHash, requestTimestamp)

	if status != StatusValid {
		return fmt.Errorf("replay guardian request rejected: %s", status)
	}

	return nil
}

// verifyTimestamp verifies that a request's timestamp is within the acceptable range.
func (r *replayGuardian) verifyTimestamp(now time.Time, requestTimestamp time.Time) ReplayGuardianStatus {
	if requestTimestamp.After(now) {
		// The request has a timestamp that is ahead of the local wall clock time.
		timeInFuture := requestTimestamp.Sub(now)
		if timeInFuture > r.maxTimeInFuture {
			return StatusTooFarInFuture
		}
	} else {
		// The request has a timestamp that is behind the local wall clock time.
		timeInPast := now.Sub(requestTimestamp)
		if timeInPast > r.maxTimeInPast {
			return StatusTooOld
		}
	}
	return StatusValid
}

// pruneObservedHashes removes any hashes from the observedHashes set that have expired. A hash is considered to have
// expired if its expiration time is before the current wall clock time.
func (r *replayGuardian) pruneObservedHashes(now time.Time) {

	// Any timestamp older than this is considered to be expired.
	oldestNonExpiredTimestamp := now.Add(-r.maxTimeInPast)

	for {
		next, ok := r.expirationQueue.Peek()
		if !ok {
			// There are no more things we are tracking.
			return
		}

		timestamp := next.(*hashWithTimestamp).timestamp
		if !timestamp.Before(oldestNonExpiredTimestamp) {
			// It's not yet time to remove this hash.
			return
		}

		// Forget about expired hash.
		r.expirationQueue.Dequeue()
		delete(r.observedHashes, next.(*hashWithTimestamp).hash)
	}
}
