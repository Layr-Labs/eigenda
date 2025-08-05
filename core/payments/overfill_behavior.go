package payments

// OverdraftBehavior describes how leaky bucket overdrafts are handled
type OverdraftBehavior string

const (
	// Disallows any overdrafts.
	//
	// If there isn't enough bucket capacity to cover a dispersal, then the dispersal will not be permitted.
	OverdraftNotPermitted OverdraftBehavior = "overdraftNotPermitted"

	// Allows a single overdraft.
	//
	// That means that if there is *any* available bucket capacity at all, then a single dispersal will be permitted,
	// and the bucket will be filled above capacity. Then, the user will have to wait for the extra to be drain before
	// making another dispersal.
	OverdraftOncePermitted OverdraftBehavior = "overdraftOncePermitted"
)
