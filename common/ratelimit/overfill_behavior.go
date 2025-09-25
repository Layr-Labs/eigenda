package ratelimit

// OverfillBehavior describes how leaky bucket overfills are handled
type OverfillBehavior string

const (
	// Disallows any overfills.
	//
	// If there isn't enough bucket capacity to cover a dispersal, then the dispersal will not be permitted.
	OverfillNotPermitted OverfillBehavior = "overfillNotPermitted"

	// Allows a single overfill.
	//
	// That means that if there is *any* available bucket capacity at all, then a single dispersal will be permitted,
	// and the bucket will be filled above capacity. Then, the user will have to wait for the extra to drain before
	// making another dispersal.
	OverfillOncePermitted OverfillBehavior = "overfillOncePermitted"
)
