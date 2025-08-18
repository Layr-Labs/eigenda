package replay

import "time"

// ReplayGuardian ensures that the same request is not processed more than once. It can be used to do things such
// as protecting against replay attacks or accidental duplicate requests.
type ReplayGuardian interface {

	// VerifyRequest verifies that a request with the given hash and timestamp is not a replay
	// of a previous request. If it cannot be determined if a request is a replay or not,
	// then the request is rejected. Only if it can be guaranteed that the request is not a replay
	// will this method return nil.
	VerifyRequest(
		requestHash []byte,
		requestTimestamp time.Time) error
}
