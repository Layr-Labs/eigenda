package ratelimit

// OverfillBehavior describes how leaky bucket overfills are handled
type OverfillBehavior string

const (
	// Disallows any overfills.
	//
	// If there isn't enough bucket capacity to cover a request, then the request will not be permitted.
	OverfillNotPermitted OverfillBehavior = "overfillNotPermitted"

	// Allows a single overfill.
	//
	// That means that if there is *any* available bucket capacity at all, then a single request will be permitted,
	// and the bucket will be filled above capacity. The next request will be required to wait for the extra to
	// drain before it is permitted.
	OverfillOncePermitted OverfillBehavior = "overfillOncePermitted"
)
