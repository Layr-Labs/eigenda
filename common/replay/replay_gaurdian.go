package replay

import "time"

// ReplayGuardian ensures that the same request is not processed more than once. It can be used to do things such
// as protecting against replay attacks or accidental duplicate requests.
type ReplayGuardian interface {

	// VerifyRequest verifies that a request with the given hash and timestamp is not a replay
	// of a previous request. If it cannot be determined if a request is a replay or not,
	// then the request is rejected. Only if it can be guaranteed that the request is not a replay
	// will this method return nil.
	//
	// In order to be a verified unique request, the following conditions must be met:
	// - the request's timestamp must be no more than X minutes ahead of the local wall clock time
	// - the request's timestamp must be no more than Y minutes behind the local wall clock time
	// - the request's hash must not have been previously observed (hashes are remembered until they are Y in the past)
	VerifyRequest(
		requestHash []byte,
		requestTimestamp time.Time) error

	// The same as VerifyRequest, but returns a detailed status code instead of an error.
	DetailedVerifyRequest(
		requestHash []byte,
		requestTimestamp time.Time) ReplayGuardianStatus
}

// ReplayGuardianStatus indicates the result of a replay guardian check.
type ReplayGuardianStatus string

const (
	// The request is not a duplicate and is within the acceptable time range.
	StatusValid ReplayGuardianStatus = "Valid"
	// The request is too old to be accepted.
	StatusTooOld ReplayGuardianStatus = "TooOld"
	// The request is too far in the future to be accepted.
	StatusTooFarInFuture ReplayGuardianStatus = "TooFarInFuture"
	// The request is a duplicate of a previously seen request.
	StatusDuplicate ReplayGuardianStatus = "Duplicate"
)
